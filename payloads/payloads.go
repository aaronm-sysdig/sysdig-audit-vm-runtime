package payloads

var K8sLiveJson string = `
{
    "requests": [
        {
            "format": {
                "type": "data"
            },
            "time": {
                "from": 1697169000000000,
                "to": 1697255400000000,
                "sampling": 600000000
            },
            "metrics": {
                "k0": "kubernetes.workload.name",
                "k1": "kubernetes.namespace.name",
                "k2": "kubernetes.cluster.name"
            },
            "group": {
                "aggregations": {},
                "groupAggregations": {},
                "by": [
                    {
                        "metric": "k0"
                    },
                    {
                        "metric": "k1"
                    },
                    {
                        "metric": "k2"
                    }
                ],
                "configuration": {
                    "groups": [
                        {
                            "groupBy": [
                                {
                                    "metric": "kubernetes.cluster.name"
                                },
                                {
                                    "metric": "kubernetes.namespace.name"
                                },
                                {
                                    "metric": "kubernetes.workload.name"
                                }
                            ]
                        }
                    ]
                }
            },
            "paging": {
                "from": 0,
                "to": 999
            },
            "sort": [
                {
                    "k2": "asc",
                    "k1": "asc",
                    "k0": "asc"
                }
            ],
            "scope": ""
        }
    ]
}`
