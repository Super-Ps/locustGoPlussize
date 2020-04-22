# coding: utf8

import six
from itertools import chain

from flask import request, Response
from locust import Locust, TaskSet, task, web, runners, stats as locust_stats
from prometheus_client import Metric, REGISTRY, exposition

# This locustfile adds an external web endpoint to the locust master, and makes it serve as a prometheus exporter.
# Runs it as a normal locustfile, then points prometheus to it.
# locust -f prometheus_exporter.py --master

# Lots of code taken from [mbolek's locust_exporter](https://github.com/mbolek/locust_exporter), thx mbolek!


class LocustCollector(object):
    registry = REGISTRY

    def collect(self):
        # locust_runner is not None, it indicates that test started.
        if runners.locust_runner:

            stats = []

            for s in chain(locust_stats.sort_stats(runners.locust_runner.request_stats),
                           [runners.locust_runner.stats.total]):
                stats.append({
                    "method": s.method,
                    "name": s.name,
                    "num_requests": s.num_requests,
                    "num_failures": s.num_failures,
                    "avg_response_time": s.avg_response_time,
                    "min_response_time": s.min_response_time or 0,
                    "max_response_time": s.max_response_time,
                    "current_rps": s.current_rps,
                    "median_response_time": s.median_response_time,
                    "avg_content_length": s.avg_content_length,
                })

            errors = [e.to_dict() for e in six.itervalues(runners.locust_runner.errors)]

            metric = Metric('locust_user_count', 'Swarmed users', 'gauge')
            metric.add_sample('locust_user_count', value=runners.locust_runner.user_count, labels={})
            yield metric

            metric = Metric('locust_errors', 'Locust requests errors', 'gauge')
            for err in errors:
                metric.add_sample('locust_errors', value=err['occurrences'],
                                  labels={'path': err['name'], 'method': err['method']})
            yield metric

            is_distributed = isinstance(runners.locust_runner, runners.MasterLocustRunner)
            if is_distributed:
                metric = Metric('locust_slave_count', 'Locust number of slaves', 'gauge')
                metric.add_sample('locust_slave_count', value=len(runners.locust_runner.clients.values()), labels={})
                yield metric

            metric = Metric('locust_fail_ratio', 'Locust failure ratio', 'gauge')
            metric.add_sample('locust_fail_ratio', value=runners.locust_runner.stats.total.fail_ratio, labels={})
            yield metric

            metric = Metric('locust_state', 'State of the locust swarm', 'gauge')
            metric.add_sample('locust_state', value=1, labels={'state': runners.locust_runner.state})
            yield metric

            stats_metrics = ['avg_content_length', 'avg_response_time', 'current_rps', 'max_response_time',
                             'median_response_time', 'min_response_time', 'num_failures', 'num_requests']

            for mtr in stats_metrics:
                mtype = 'gauge'
                if mtr in ['num_requests', 'num_failures']:
                    mtype = 'counter'
                metric = Metric('locust_requests_' + mtr, 'Locust requests ' + mtr, mtype)
                for stat in stats:
                    if 'Total' not in stat['name']:
                        metric.add_sample('locust_requests_' + mtr, value=stat[mtr],
                                          labels={'path': stat['name'], 'method': stat['method']})
                yield metric


@web.app.route("/export/prometheus")
def prometheus_exporter():
    registry = REGISTRY
    encoder, content_type = exposition.choose_encoder(request.headers.get('Accept'))
    if 'name[]' in request.args:
        registry = REGISTRY.restricted_registry(request.args.get('name[]'))
    body = encoder(registry)
    return Response(body, content_type=content_type)


REGISTRY.register(LocustCollector())


class DummyTaskSet(TaskSet):
    @task(20)
    def hello(self):
        pass


class Dummy(Locust):
    task_set = DummyTaskSet
