from flask import Flask, jsonify, request, abort
from elasticsearch import Elasticsearch
import json
from datetime import datetime
from urllib.parse import urlparse, parse_qs
import time
from dlivecat import logfunc, es_search, es_update, es_count, es
from elasticsearch import Elasticsearch
from pyfasttext import FastText
import re
import emoji
import langid
import requests
import random

KEYWORD_LIST_URL = "https://www.ilivenet.com/manual.txt"

# wget https://dl.fbaipublicfiles.com/fasttext/supervised-models/lid.176.bin
model = FastText('lid.176.bin')

app = Flask(__name__)


def is_ascii(s):
    return all(ord(c) < 128 for c in s)

def remove_meaningless_string(text):
    # url是沒意義的所以予以刪除
    RE_URL = re.compile("(https?://[^\s]+)", flags=re.UNICODE)
    # 移除emoji
    text = emoji.get_emoji_regexp().sub('', text)
    text = RE_URL.sub(r'', text)
    text = text.replace('\n', ' ')
    text = text + '\n'
    return text


def detect_language(string):
    #先移除掉沒意義的字串後再開始偵測語言，提高準確性
    string = remove_meaningless_string(string)
    langid_res = langid.classify(string)[0]
    # langid在判別中日文上比較優秀
    if langid_res == "zh":
        return langid_res
    fasttext_res = model.predict_single(string, k=1)
    if not fasttext_res:
        return langid_res
    fasttext_res = fasttext_res[0]
    if fasttext_res == 'ja':
        if langid_res == 'ja':
            return 'ja'
        else:
            return langid_res

    return fasttext_res


def query_elastic(q, fr=0, sz=50):
    body = {
        "size": sz,
        "from": fr,
        "query": {
            "bool": {
                "must_not": [
                        {"match_phrase": {"status": "invalid"}},
                 ],
                "should": [
                    {
                        "match": {
                            "title": {
                                "query": q,
                                "boost": 2,
                                "minimum_should_match": "90%",
                            }
                        }
                    }, {
                        "match": {
                            "description": {
                                "query": q,
                                "minimum_should_match": "70%",
                            }
                        }
                    }, {
                        "match": {
                            "tags": {
                                "query": q,
                                "boost": 4,
                                "minimum_should_match": "80%",
                            }
                        }
                    }, {
                        "match": {
                            "host": {
                                "query": q,
                                "boost": 4,
                                "minimum_should_match": "80%",
                            }
                        }
                    }, {
                        "match": {
                            "platform": {
                                "query": q,
                                "boost": 3,
                                "minimum_should_match": "80%",
                            }
                        }
                    }
                ],
                "minimum_should_match": 1,
            }
        }
    }
    # 針對非英文再加上match_phrase
    if not is_ascii(q):
        body["query"]["bool"]["should"].extend([
            {
                "match_phrase": {
                    "title": {
                        "query": q,
                        "boost": 2,
                        "slop": int(len(q) * 0.4) + 1
                    }
                }
            }, {
                "match_phrase": {
                    "description": {
                        "query": q,
                        "slop": int(len(q) * 0.6) + 1
                    }
                }
            }, {
                "match_phrase": {
                    "tags": {
                        "query": q,
                        "boost": 4,
                        "slop": int(len(q) * 0.2) + 1
                    }
                }
            }, {
                "match_phrase": {
                    "host": {
                        "query": q,
                        "boost": 4,
                        "slop": int(len(q) * 0.2) + 1
                    }
                }
            }, {
                "match_phrase": {
                    "platform": {
                        "query": q,
                        "boost": 3,
                        "slop": int(len(q) * 0.4) + 1
                    }
                }
            },
        ])

    body["sort"] = [
        "_score",
        {"published": {"order": "desc"}},
        {"timestamp": {"order": "desc"}},
    ]

    res = es_search(body)
    if not res:
        return False

    # if no results
    if len(res['hits']['hits']) == 0:
        body = {
            "size": sz,
            "from": fr,
            "query": {
                "bool": {
                    "must_not": [
                        {"match_phrase": {"status": "invalid"}},
                    ],
                    "should": [
                        {"match_phrase": {"platform": "twitch"}},
                        {"match_phrase": {"platform": "youtube"}}
                    ]
                }
            },
            "sort": [
                {"viewers": {"order": "desc"}},
                {"timestamp": {"order": "desc"}},
            ]
        }
        res = es_search(body)
        if not res:
            return False
        res["found"] = False
    else:
        res["found"] = True
    return res


def get_platform_data(platform, fr=0, sz=30, language="", exclude_language=[]):
    query = {
        "bool": {
            "must": [
                {"match_phrase": {"status": "live"}},
            ],
        }
    }
    if not platform == "all":
        query["bool"]["filter"] = [{"match_phrase": {"platform": platform}}]
    if language:
        query["bool"]["filter"].append({"match_phrase": {"language": language}})
    if exclude_language:
        if not "must_not" in query["bool"] or not query["bool"]["must_not"]:
            query["bool"]["must_not"] = []
        elif not type(query["bool"]["must_not"]) == list:
            query["bool"]["must_not"] = [query["bool"]["must_not"]]
        for l in exclude_language:
            query["bool"]["must_not"].append({"match_phrase": {"language": l}})

    # body = {
    #     "from": fr,
    #     "size": sz,
    #     "query": query,
    #     "sort": [
    #         {"viewers": {"order": "desc"}},
    #         {"timestamp": {"order": "desc"}},
    #         {"published": {"order": "desc"}},
    #     ]
    # }
    body = {
        "from": fr,
        "size": sz,
        "query": {
            "function_score": {
                "query": query,
                "random_score": {
                    "seed": str(int(time.mktime(datetime.now().timetuple()))),
                    "field": "_seq_no"
                },
                "boost": "5",
                "boost_mode": "replace"
            }
        },
    }
    res = es_search(body)
    return res


def get_channel_data(channel):
    body = {
        "size": 1,
        "query": {
            "bool": {
                "must": [
                    {"match_phrase": {"channel": channel}},
                    {"match_phrase": {"status": "live"}},
                ],
            }
        },
        "sort": [
            {"timestamp": {"order": "desc"}},
            {"published": {"order": "desc"}},
        ]
    }
    res = es_search(body)
    return res


def trans_to_smallcase_key(form):
    new_dic = {}
    for k in form.keys():
        s_k = k.lower()
        new_dic[s_k] = form[k]
    return new_dic


def get_parameters_from_url(request):
    parsed = urlparse(request.url)
    return parse_qs(parsed.query)


def get_random_streams(fr=0, sz=30):
    body = {
        "from": fr,
        "size": sz,
        "query": {
            "function_score": {
                "query": {
                    "bool": {
                        "filter": [
                            {"match_phrase": {"status": "live"}},
                        ],
                    }
                },
                "random_score": {
                    "seed": str(int(time.mktime(datetime.now().timetuple()))),
                    "field": "_seq_no"
                },
                "boost": "5",
                "boost_mode": "replace"
            }
        }
    }

    return es_search(body=body)



@app.route("/", methods=['GET'])
def search():
    qs = get_parameters_from_url(request)
    res = {}
    fr = 0 if not qs.get("from", []) else qs["from"][0]
    sz = 30 if not qs.get("size", []) else qs["size"][0]
    lang = "" if not qs.get("language", []) else qs["language"][0]
    if 'q' in qs:
        res = query_elastic(qs['q'][0], fr, sz)
    elif 'platform' in qs:
        res = get_platform_data(qs["platform"][0], fr, sz, lang)
    elif 'channel' in qs:
        res = get_channel_data(qs["channel"][0])
    logfunc("'GET'", request.url)
    return jsonify(res)


@app.route("/platform_page", methods=['GET'])
def get_platform_page():
    qs = get_parameters_from_url(request)
    if "size" in qs:
        size = qs["size"][0]
    else:
        size = 6

    def request_es(platform, lang_list):
        resp = []
        for lang in lang_list:
            res = get_platform_data(platform, sz=size, language=lang)
            if not res:
                res = []
            resp.append({"platform": platform,
                         "language": lang,
                         "data": res["hits"]["hits"]})

        res = get_platform_data(platform, sz=size, exclude_language=lang_list)
        if not res:
            res = []
        resp.append({"platform": platform,
                     "language": "others",
                     "data": res["hits"]["hits"]})
        return resp

    response = []
    # YouTube
    language_list = ["en", 'zh', 'ja']
    response.extend(request_es("YouTube", language_list))
    # Twitch
    language_list = ["en", 'zh', 'ja']
    response.extend(request_es("Twitch", language_list))
    # Facebook
    language_list = ["en", 'zh', 'ja']
    response.extend(request_es("Facebook", language_list))
    # 17直播
    language_list = ["en", 'zh', 'ja']
    response.extend(request_es("17直播", language_list))
    # 西瓜直播
    language_list = ["en", 'zh', 'ja']
    response.extend(request_es("西瓜直播", language_list))

    logfunc("'GET'", request.url)
    return jsonify(response)


@app.route("/update_click_through", methods=['GET'])
def update_videos_click_through():
    qs = get_parameters_from_url(request)
    try:
        video_url = qs["videourl"][0]
    except KeyError:
        abort(400)
    res = es_search(body={
        "query": {
            "bool": {
                "must": [
                    {"match_phrase": {"videourl": video_url}},
                ],
                "filter": [
                    {"match_phrase": {"status": "live"}},
                ]
            }
        },
        "_source": "_id",
    })
    if not res:
        return False

    if len(res['hits']['hits']) == 0:
        return "Can't find corresponding data"
    else:
        res = es_update(_id=res["hits"]["hits"][0]["_id"],
                        body={"script": {
                            "source": "if(ctx._source.containsKey(\"click_through\")){ctx._source.s+=params.count} else{ctx._source.click_through=1}",
                            "lang": "painless",
                            "params": {
                                "count": 1
                            }}})
        if not res:
            return False
    return 'ok'


@app.route("/add", methods=['POST'])
def create_or_update_doc():
    if request.is_json:
        form = request.get_json()
    else:
        form = request.form.copy()
    form = trans_to_smallcase_key(form)

    try:
        if not form["host"] or not form["platform"] or not form["title"] or not form["published"]:
            return abort(400)
    except KeyError:
        return abort(400)

    # 以videourl當作unique ID
    res = es_search(body={
        "query": {
            "bool": {
                "must": {"match_phrase": {"videourl": form["videourl"]}},
                "filter": [
                    {"match_phrase": {"host": form["host"]}},
                    {"match_phrase": {"platform": form["platform"]}}
                ],
                "must_not": [
                    {"match_phrase": {"status": "invalid"}},
                ]
            },
        },
        "_source": "_id",
        "sort": {"timestamp": {"order": "desc"}},
    })
    if not res:
        return False

    form["timestamp"] = datetime.now().strftime("%Y-%m-%dT%H:%M:%SZ")
    form["click_through"] = 0
    if not "status" in form or not form["status"]:
        form["status"] = "live"
    try:
        published_time = time.strptime(form["published"], "%Y-%m-%dT%H:%M:%SZ")
    except ValueError:
        published_time = time.strptime(form["published"], "%Y-%m-%dT%H:%M:%S+0000")
    viewers = 0 if not form.get("viewers", None) else int(form["viewers"])
    form["popular_rate"] = int(viewers * 10000000000 / (time.mktime(datetime.now().timetuple())-time.mktime(published_time)))

    try:
        if len(res['hits']['hits']) == 0:
            # Classify video's language
            if not "language" in form or not form["language"]:
                test_string = form["title"] + " " + form["description"] + " " + form["host"]
                form["language"] = detect_language(test_string)
            # Create data
            es.index(index="livestreams", body=form)

        else:
            res = es_update(_id=res["hits"]["hits"][0]["_id"],
                            body={"doc": {"timestamp": form["timestamp"],
                                          "description": form["description"],
                                          "status": form["status"],
                                          "popular_rate": form["popular_rate"]}
                                  }
                            )
            if not res:
                return abort(500)

    except Exception as e:
        logfunc("'Wrong Request'")
        print(form)
        print(e)
        abort(400)
    return 'ok'


@app.route("/total_streams", methods=['GET'])
def total_streams():
    qs = get_parameters_from_url(request)
    body = {
        "query": {
            "bool" :{
                "filter": [
                    {"match_phrase": {"status": "live"}},
                ],
            }
        }
    }
    if "platform" in qs:
       body["query"]["bool"]["must"] = [{"match_phrase": {"platform": qs["platform"][0]}}]
    response = es_count(body=body)
    if not response:
        abort(400)
    return jsonify(response)


@app.route("/hot_page", methods=['GET'])
def get_top_popular_rate():
    body = {
        "size": 48,
        "query": {
            "bool": {
                "filter": [
                    {"match_phrase": {"status": "live"}},
                ],
            }
        },
        "sort" :[
            {"popular_rate": {"order": "desc"}},
            {"published": {"order": "desc"}},
        ]
    }
    response = es_search(body=body)
    if not response:
        abort(400)
    return jsonify(response)


@app.route("/top_viewers", methods=['GET'])
def get_top_viewers():
    body = {
        "size": 48,
        "query": {
            "bool": {
                "filter": [
                    {"match_phrase": {"status": "live"}},
                ],
            }
        },
        "sort" :[
            {"popular_rate": {"order": "desc"}},
            {"published": {"order": "desc"}},
        ]
    }
    response = es_search(body=body)
    if not response:
        abort(400)
    return jsonify(response)


@app.route("/home_page", methods=['GET'])
def home_page():
    results = get_random_streams()
    response = {"subscriptions": results["hits"]["hits"][0:4],  # Subscriptions
                "upcoming": results["hits"]["hits"][4:5],       # Upcoming Stream
                "recommended": results["hits"]["hits"][5:9],    # Recommended
                "today": results["hits"]["hits"][9:13],
                "within_72_hours": results["hits"]["hits"][13:17]}
    # Most Viewed
    body = {
        "size": 4,
        "query": {
            "bool": {
                "filter": [
                    {"match_phrase": {"status": "live"}},
                ],
            }
        },
        "sort": [
            {"viewers": {"order": "desc"}},
            {"published": {"order": "desc"}},
        ]
    }
    results = es_search(body=body)
    response["most_viewed"] = results["hits"]["hits"][0:4]
    # Hot
    body = {
        "size": 4,
        "query": {
            "bool": {
                "filter": [
                    {"match_phrase": {"status": "live"}},
                ],
            }
        },
        "sort": [
            {"popular_rate": {"order": "desc"}},
            {"published": {"order": "desc"}},
        ]
    }
    results = es_search(body=body)
    response["hot"] = results["hits"]["hits"][0:4]

    return jsonify(response)


@app.route("/cover_page", methods=['GET'])
def cover_page():
    try:
        r = requests.get(KEYWORD_LIST_URL)
        keywords = json.loads(r.content)['keyword']
    except requests.exceptions.ConnectionError:
        logfunc("Can't connect to ", KEYWORD_LIST_URL)
        keywords = []
    if len(keywords):
        keyword_results = []
        for k in keywords:
            result = query_elastic(k, sz=12)
            keyword_results.extend(result["hits"]["hits"])
        keyword_results = random.shuffle(keyword_results)
    else:
        keyword_results = get_random_streams(sz=12)["hits"]["hits"]
    return jsonify(keyword_results)