from random import choice
from locust import HttpLocust, TaskSet


ids = []


def api_create_url(l):
    with l.client.post('/api/v1/url/', '{"url": "http://google.com/"}') as resp:
        ids.append(resp.json()['id'])

def api_retrieve_url(l):
    if len(ids) > 0:
        l.client.get('/api/v1/url/{}'.format(choice(ids)))

def redirect_to_long_url(l):
    if len(ids) > 0:
        with l.client.get('/{}'.format(choice(ids)), allow_redirects=False, catch_response=True) as resp:
            if resp.status_code == 301:
                resp.success()
            else:
                resp.failure('expect redirect: {}'.format(resp.status_code))


class GoshBehaviour(TaskSet):
    tasks = {
        api_create_url: 5,
        api_retrieve_url: 4,
        redirect_to_long_url: 4,
    }


class Gosh(HttpLocust):
    task_set = GoshBehaviour
    min_wait = 10
    max_wait = 9000
