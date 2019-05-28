from elasticsearch import Elasticsearch, ConflictError
from urllib3.exceptions import TimeoutError, NewConnectionError
import datetime

es = Elasticsearch(timeout=60, max_retries=5, retry_on_timeout=True)

def logfunc(*string):
    print(datetime.datetime.now(), end='  ')
    for i in string:
        print(i, end=' ')
    print()


def connect_es_decorator(func):
    def with_attempt(*args, **kwargs):
        success = False
        retries_counter = 0
        results = None
        while not success and retries_counter < 3:
            try:
                results = func(*args, **kwargs)
                success = True
            except (TimeoutError, NewConnectionError) as e:
                logfunc(e)
                retries_counter += 1
        if not success:
            return False
        if not results:
            return True
        return results
    return with_attempt


@connect_es_decorator
def es_search(body):
    return es.search(index='livestreams',  body=body)

@connect_es_decorator
def es_update(_id, body):
    es.update(index="livestreams", id=_id, body=body)

@connect_es_decorator
def es_count(body):
    return es.count(index="livestreams", body=body)


