from elasticsearch import Elasticsearch, ConflictError
from urllib3.exceptions import TimeoutError, NewConnectionError
import datetime

es = Elasticsearch(timeout=60, max_retries=5, retry_on_timeout=True)

def logfunc(*string):
    print(datetime.datetime.now(), end='  ')
    for i in string:
        print(i, end=' ')
    print()

def es_search(body):
    success = False
    retries_counter = 0
    results = {}
    while not success and retries_counter < 3:
        try:
            results = es.search(index='livestreams',  body=body)
            success = True
        except (TimeoutError, NewConnectionError) as e:
            logfunc(e)
            retries_counter += 1
    if not success:
        return False
    return results


def es_update(_id, body):
    success = False
    retries_counter = 0
    while not success and retries_counter < 3:
        try:
            es.update(index="livestreams", id=_id, body=body)
            success = True
        except (ConflictError, TimeoutError, NewConnectionError) as e:
            logfunc(e)
            retries_counter += 1
    return success

