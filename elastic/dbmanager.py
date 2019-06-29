import time
import threading
import datetime
import requests
from requests.exceptions import HTTPError as RequestsHTTPError
from diffimg import diff
import os
from dlivecat import logfunc, es_search, es_update
import random
import urllib.request
from urllib.error import HTTPError
from eserver import query_elastic, get_random_streams
import json

from elasticsearch import Elasticsearch
es = Elasticsearch(timeout=60, max_retries=5, retry_on_timeout=True)

TOLERANT__PUBLISHED_TIMEDELTA = datetime.timedelta(days=7)
TOLERANT__TIMESTAMP_TIMEDELTA = datetime.timedelta(minutes=30)
KEYWORD_LIST_URL = "https://www.ilivenet.com/manual.txt"
KEYWORD_RESULTS_TEXT_FILE = "keywords_in_front_page.txt"

class ExpireDataThread(threading.Thread):
    def run(self):
        print(self.name, " starts!!")
        try:
            while True:
                body = {
                    "size": 1000,
                    "query": {
                        "bool":{
                            "must":[
                                {"range": {"timestamp": {"lt": datetime.datetime.now() - TOLERANT__TIMESTAMP_TIMEDELTA}}}
                            ],
                            "must_not": [
                                {"match_phrase": {"status": "invalid"}},
                            ]
                        }
                    },
                }
                results = es_search(body=body)
                if not results:
                    continue
                for hit in results['hits']['hits']:
                    es_update(hit['_id'], {"script": {"source": "ctx._source.status='invalid'"}})
                logfunc(self.name, "Mark {} data as invalid".format(len(results['hits']['hits'])))
                time.sleep(5)

        except KeyboardInterrupt:
            print("Forced Stop.")

        except Exception as e:
            logfunc(e)



class BaseCheckThumbnailsThread(threading.Thread):
    compared_img_name = "compared_{}_{}.jpg"
    sample_img = "sample_{}.jpg"
    platform = ""

    def __init__(self, name_no=0):
        super().__init__()
        self.daemon = True
        self.name_no = name_no
        self.name = "Thread-{}-{}".format(self.platform, name_no)
        self.compared_img_name = self.compared_img_name.format(self.platform, name_no)
        self.sample_img = self.sample_img.format(self.platform)

    def compare_img(self):
        res = diff(self.sample_img, self.compared_img_name)
        try:
            os.remove(self.compared_img_name)
        except Exception as e:
            logfunc(self.name, e)
        # If two image are same, then return True
        if res > 0.0:
            return False
        return True

    def download(self, url):
        urllib.request.urlretrieve(url, self.compared_img_name)

    def delete_data(self, hit):
        if not es_update(hit['_id'], {"script": {"source": "ctx._source.status='invalid'"}}):
           logfunc("Can't delete", hit['_id'])
        else:
            pass
            #logfunc(self.name, "Delete", hit['_source']["host"])

    def process(self, hit):
        self.download(hit['_source']["thumbnails"])
        if self.compare_img():
            self.delete_data(hit)

    def run(self):
        print(self.name, " starts!!")
        try:
            while True:
                # body = {
                #     "size": 3000,
                #     "query": {
                #         "bool": {
                #             "filter": [
                #                 {"match_phrase": {"status": "live"}},
                #                 {"match_phrase": {"platform": self.platform}},
                #                # {"range": {"timestamp": {"gt": datetime.datetime.now() - datetime.timedelta(minutes=10)}}}
                #             ]
                #         }
                #     },
                #     "sort": [
                #         {"viewers": {"order": "desc"}},
                #         {"timestamp": {"order": "desc"}},
                #         {"published": {"order": "desc"}},
                #     ],
                #     "_source": ["_id", "thumbnails", "host"],
                # }
                body = {
                    "size": 100,
                    "query": {
                        "function_score": {
                            "query": {
                                "bool": {
                                    "filter": [
                                        {"match_phrase": {"status": "live"}},
                                        {"match_phrase": {"platform": self.platform}},
                                    ]
                                }
                            },
                            "random_score": {
                                "seed": str(int(time.mktime(datetime.datetime.now().timetuple()))+int(random.random()*100000000)*self.name_no),
                                "field": "_seq_no"
                            },
                            "boost": "5",
                            "boost_mode": "replace"
                        }
                    },

                    "_source": ["_id", "thumbnails", "host"],
                }

                results = es_search(body)
                if not results:
                    time.sleep(10)
                    continue

                for hit in results["hits"]["hits"]:
                    try:
                        if not hit['_source']["thumbnails"]:
                            continue
                        self.process(hit)
                    except KeyError:
                        logfunc(self.name, "No Key 'thumbnails'")
                    except RequestsHTTPError as e:
                        logfunc(self.name, e)
                    except Exception as e:
                        logfunc(self.name, e)
                time.sleep(1)

        except KeyBoardInterrupt:
            if os.path.isfile(self.compared_img_name):
                os.remove(self.compared_img_name)


class CheckTwitchThumbnailsThread(BaseCheckThumbnailsThread):
    platform = "twitch"


class CheckYoutubeThumbnailsThread(BaseCheckThumbnailsThread):
    platform = "youtube"
    def process(self, hit):
        try:
            self.download(hit['_source']["thumbnails"])
        except HTTPError:
            self.delete_data(hit)


class ManageHomePageStreamsDisplayThread(threading.Thread):
    def __init__(self):
        super().__init__()
        self.daemon = True

    def run(self):
        try:
            while True:
                try:
                    r = requests.get(KEYWORD_LIST_URL)
                    keywords = json.loads(r.content)['keyword']
                except requests.exceptions.ConnectionError:
                    logfunc("Can't connect to ", KEYWORD_LIST_URL)
                    keywords = []
                if len(keywords):
                    keyword_dic = {}
                    for k in keywords:
                        keyword_dic[k] = query_elastic(k, sz=5)
                    keyword_results = []
                    while len(keyword_results) < 12:
                        error_counter = 0
                        for k in keywords:
                            if not len(keyword_dic[k]["hits"]["hits"]):
                                error_counter += 1
                                continue
                            ci = random.choice(range(len(keyword_dic[k]["hits"]["hits"])))
                            keyword_results.append(keyword_dic[k]["hits"]["hits"][ci])
                            del(keyword_dic[k]["hits"]["hits"][ci])
                            if len(keyword_results) == 12:
                                break
                        if error_counter == len(keywords):
                            logfunc("No enough results to display in cover page!")
                            break
                else:
                    keyword_results = get_random_streams(sz=12)["hits"]["hits"]

                with open(KEYWORD_RESULTS_TEXT_FILE, "w") as f:
                    json.dump(keyword_results, f)
                time.sleep(0)
        except KeyboardInterrupt:
            pass


expirethread = ExpireDataThread()
expirethread.start()

for i in range(10):
    check_twitch_thread = CheckTwitchThumbnailsThread(i)
    check_twitch_thread.start()

check_youtube_thread = CheckYoutubeThumbnailsThread()
check_youtube_thread.start()

managefrontpagedisplay_thread = ManageHomePageStreamsDisplayThread()
managefrontpagedisplay_thread.start()

