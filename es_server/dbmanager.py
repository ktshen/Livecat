from elasticsearch import Elasticsearch, ConflictError
import time
import threading
import datetime
import requests
from requests.exceptions import HTTPError
from diffimg import diff
import os

TOLERANT__PUBLISHED_TIMEDELTA = datetime.timedelta(days=7)
TOLERANT__TIMESTAMP_TIMEDELTA = datetime.timedelta(hours=1)

es = Elasticsearch()

def logfunc(*string):
    print(datetime.datetime.now(), " :", *string)


class ExpireDataThread(threading.Thread):
    def run(self):
        try:
            while True:
                try:
                    body = {
                        "query": {
                            "bool":{
                                "should" :[
                                    {"range": {"timestamp": {"lt": datetime.datetime.now() - TOLERANT__TIMESTAMP_TIMEDELTA}}},
                                ]

                            }
                        }}
                    results = es.delete_by_query(index="livestreams", doc_type="_doc", body=body)
                    logfunc("Clear ", results['deleted'], " of data in es")
                    time.sleep(300)
                except ConflictError as e:
                    logfunc(datetime.datetime.now(), " ", str(e))
                    time.sleep(10)


        except KeyboardInterrupt:
            print("Forced Stop.")

        except Exception as e:
            print(e)


class CheckTwitchValidityThread(threading.Thread):
    def __init__(self):
        super().__init__()
        self.daemon = True
        self.compared_img_name = "compared.jpg"
        self.sample_img = "sample.jpg"

    def compare_img(self):
        res = diff(self.sample_img, self.compared_img_name)
        try:
            os.remove(self.compared_img_name)
        except Exception as e:
            print(e)
        # If two image are same, then return True
        if res > 0.0:
            return False
        return True

    def download(self, url):
        with open(self.compared_img_name, 'wb') as handle:
            response = requests.get(url, stream=True)
            if not response.ok:
                raise HTTPError("Someting goes wrong, can't get the image file.")

            for block in response.iter_content(1024):
                if not block:
                    break
                handle.write(block)

    def run(self):
        print("CheckTwitchValidityThread starts!!!")
        while True:
            body = {
                "size": 3000,
                "query": {
                    "bool": {
                        "must": [
                            {"match_phrase": {"platform": "twitch"}},
                           # {"range": {"timestamp": {"gt": datetime.datetime.now() - datetime.timedelta(minutes=10)}}}
                        ]
                    }
                },
                "sort": [
                    {"viewers": {"order": "desc"}},
                    {"timestamp": {"order": "desc"}},
                    {"published": {"order": "desc"}},
                ],
                "_source": ["_id", "thumbnails", "host"],
            }
            try:
                results = es.search(index="livestreams", body=body)
            except ConflictError as e:
                logfunc(datetime.datetime.now(), " TwitchThread ", str(e))
                time.sleep(10)
                continue

            for hit in results["hits"]["hits"]:
                try:
                    if not hit['_source']["thumbnails"]:
                        continue
                    self.download(hit['_source']["thumbnails"])
                    if self.compare_img():
                        es.delete(index="livestreams", doc_type="_doc", id=hit["_id"])
                        print(datetime.datetime.now(), " Delete ", hit['_source']["host"])
                    time.sleep(0.0005)
                except KeyError:
                    print("No Key 'thumbnails'")
                except HTTPError as e:
                    print(e)
                except Exception as e:
                    print(e)
            time.sleep(60)


expirethread = ExpireDataThread()
check_twitch_thread = CheckTwitchValidityThread()
check_twitch_thread.start()
expirethread.start()

