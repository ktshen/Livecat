from flask import Flask, jsonify, request, abort
from elasticsearch import Elasticsearch
import json
from datetime import datetime
from urllib.parse import urlparse, parse_qs



app = Flask(__name__)
es = Elasticsearch()

def is_ascii(s):
    return all(ord(c) < 128 for c in s)


def es_search(q):
    if is_ascii(q):
        body = {
            "size": 50,
            "query":{
                "bool": {
                    "should": [
                    {
                        "match":{
                            "title": {
                                "query": q,
                                "boost": 2,
                                "minimum_should_match": "90%",
                            }
                        }
                    },{
                        "match":{
                            "description": {
                                "query": q,
                                "minimum_should_match": "70%",
                            }
                        }
                    },{
                        "match":{
                            "tags": {
                                "query": q,
                                "boost": 4,
                                "minimum_should_match": "80%",
                            }
                        }
                    },{
                        "match":{
                            "host": {
                                "query": q,
                                "boost": 4,
                                "minimum_should_match": "80%",
                            }
                        }
                    },{
                        "match":{
                            "platform": {
                                "query": q,
                                "boost": 3,
                                "minimum_should_match": "80%",
                            }
                        }
                    }
                   ],
                }
            }
        }
    else:
        body = {
            "size": 50,
            "query":{
                "bool": {
                    "should": [
                    {
                        "match_phrase":{
                            "title": {
                                "query": q,
                                "boost": 2,
                                "slop": int(len(q) * 0.4)
                            }
                        }
                    },{
                        "match_phrase":{
                            "description": {
                                "query": q,
                                "slop": int(len(q) * 0.6)
                            }
                        }
                    },{
                        "match_phrase":{
                            "tags": {
                                "query": q,
                                "boost": 4,
                                "slop": int(len(q) * 0.2)
                            }
                        }
                    },{
                        "match_phrase":{
                            "host": {
                                "query": q,
                                "boost": 4,
                                "slop": int(len(q) * 0.2)
                            }
                        }
                    },{
                        "match_phrase":{
                            "platform": {
                                "query": q,
                                "boost": 3,
                                "slop": int(len(q) * 0.4)
                            }
                        }
                    }
                   ],
                }
            }
        }
    body["sort"] = [
        "_score",
        {"timestamp": {"order": "desc"}},
        {"published": {"order": "desc"}},
    ]
    res = es.search(index="livestreams", body=body)

    # if no results
    if res['hits']['total'] == 0:
        body = {
            "size": 20,
            "query": {
                "bool": {
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
        res = es.search(index="livestreams", body=body)
    return res


def get_platform_data(platform, fr, sz):
    if platform == "all":
        query = {"match_all": {}}
    else:
        query = {"match_phrase": {"platform": platform}}
    body = {
        "from": fr,
        "size": sz,
        "query": query,
        "sort": [
            {"timestamp": {"order": "desc"}},
            {"published": {"order": "desc"}},
        ]
    }
    res = es.search(index="livestreams", body=body)
    return res


def get_channel_data(channel):
    body = {
        "size": 1,
        "query": {
            "match_phrase": {"channel": channel}
        },
        "sort": [
            {"timestamp": {"order": "desc"}},
            {"published": {"order": "desc"}},
        ]
    }
    res = es.search(index="livestreams", body=body)
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


@app.route("/", methods=['GET'])
def search():
    qs = get_parameters_from_url(request)
    res = {}
    if 'q' in qs:
        res = es_search(qs['q'][0])
    elif 'platform' in qs:
        fr = 0 if not qs.get("from", []) else qs["from"][0]
        sz = 20 if not qs.get("size", []) else qs["size"][0]
        res = get_platform_data(qs["platform"][0], fr, sz)
    elif 'channel' in qs:
        res = get_channel_data(qs["channel"][0])
    print(datetime.now(), " ", "'GET' ", request.url)
    return jsonify(res)

@app.route("/home_page", methods=['GET'])
def get_home_page():
    qs = get_parameters_from_url(request)
    if "size" in qs:
        size = qs["size"][0]
    else:
        size = 6

    def request_es(platform, sort_list):
        res = es.search(index="livestreams", body={
            "size": size,
            "query": {
                "match_phrase": {"platform": platform}
            },
            "sort": sort_list
        })
        return res["hits"]["hits"]

    response = []
    # Youtube
    response.extend(request_es("youtube", [{"viewers": {"order": "desc"}},{"timestamp": {"order": "desc"}}]))
    # Twitch
    response.extend(request_es("twitch", [{"viewers": {"order": "desc"}}, {"timestamp": {"order": "desc"}}]))
    # Facebook
    response.extend(request_es("facebook", [{"published": {"order": "desc"}}, {"timestamp": {"order": "desc"}}]))
    #17直播
    response.extend(request_es("17直播", [{"viewers": {"order": "desc"}}, {"timestamp": {"order": "desc"}}]))
    #西瓜直播
    response.extend(request_es("西瓜直播", [{"viewers": {"order": "desc"}}, {"timestamp": {"order": "desc"}}]))
    return jsonify(response)


@app.route("/update_click_through", methods=['GET'])
def update_videos_click_through():
    qs = get_parameters_from_url(request)
    try:
        video_url = qs["videourl"][0]
    except KeyError:
        abort(400)
    res = es.search(index="livestreams", body={
        "query":{
            "match_phrase": {"videourl": video_url}
        },
        "_source": "_id",
    })
    if res['hits']['total'] == 0:
        return "Can't find corresponding data"
    else:
        es.update(index="livestreams", id=res["hits"]["hits"][0]["_id"], doc_type='_doc',
                  body={"script" : {
                            "source": "if(ctx._source.containsKey(\"click_through\")){ctx._source.click_through+=params.count} else{ctx._source.click_through=1}",
                            "lang": "painless",
                            "params" : {
                                "count" : 1
                            }
                        }
                  }
        )
    return 'ok'



@app.route("/add", methods=['POST'])
def create_or_update_doc():
    if request.is_json:
        form = request.get_json()
    else:
        form = request.form.copy()
    form = trans_to_smallcase_key(form)

    try:
        if not form["host"] or not form["platform"]:
            return abort(400)
    except KeyError:
        return abort(400)
    # If host and platform both match a document, then I assume that this document should be updated
    res = es.search(index="livestreams", body={
        "query":{
            "bool" :{
                "must": [
                    {"match_phrase": {"host": form["host"]}},
                    {"match_phrase": {"platform": form["platform"]}}
                ],
                "should": [
                    {"match_phrase": {"title": form["title"]}},
                    {"match_phrase": {"description": form["description"]}},
                ]
            },
        },
        "_source": "_id",
        "sort": {"timestamp": {"order": "desc"}},
    })

    form["timestamp"] = datetime.now().strftime("%Y-%m-%dT%H:%M:%SZ")
    form["click_through"] = 0

    try:
        if res['hits']['total'] == 0:

            es.index(index="livestreams", body=form, doc_type='_doc')
        else:
            es.update(index="livestreams", id=res["hits"]["hits"][0]["_id"], doc_type='_doc',
                      body={"doc": {"timestamp": form["timestamp"],
                                    "title": form["title"],
                                    "description": form["description"]}})
    except Exception as e:
        print(datetime.now(), " 'Wrong Request' ")
        print(form)
        print(e)
        abort(400)
    return 'ok'


