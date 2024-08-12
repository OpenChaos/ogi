## Ogi

[![Go Report Card](https://goreportcard.com/badge/OpenChaos/ogi)](https://goreportcard.com/report/OpenChaos/ogi)

### What It Is?

After a redesign.. **Ogi** is a thin core utility that's a "context free", "utility customizable" workflow runner.

In its basic design, it would allow one to use available (or separately/privately written) flows at any/all stages of `Consumer`, `Transformer` & `Producer`.

* `Consumer` is the init point, that spins up any kind of data retrieval or listens to published requests. Then hands them over to configured `Transformer`. Like reading a file line-by-line, reading records from a database, subscribing to a kafka topic, etc.

* `Transformer` has the flow for any checks/edits required on received data; any pre-final task spins. Then the desired form of received data is handed over to configured `Producer`. Like checking if keys are valid before Producer persists it, create a JSON based on it that can be sent as a Request by Producer, etc.

* `Producer` contains the final action desired on the desired state of data/request received via Consumer-to-Transformer. It could be parsist to a file, send to a Kafka topic, add a database entry, make a custom HTTP Request, etc.

> * Allows custom components for Consumer/Transformer/Producer to be dynamically loaded via [Golang Plug-ins](https://pkg.go.dev/plugin).
>
> * Each Ogi run can pick a suitable type of each based on configuration; allowing all sorts of mix-match.
>
> * `Producer` can return an info response to `Transformer`; and it to `Consumer`. So even bidirectional flows are possible.
>
> * Each flow gets logged with a `Message ID`; so one can debug/err-check/observe events.

*Repo for stable plug-ins is [OpenChaos/ogi-plugin-umbrella](https://github.com/OpenChaos/ogi-plugin-umbrella).*

![ogi means a japanese folding hand-fan](docs/ogi.png "ogi means a japanese folding hand-fan")

---

### BELOW CONTENT IS WIP

#### How it works?

It contains 3 primary components. A consumer, transformer and producer.

Ogi initializes Consumer and let it handle the flow to transformer, or if required directly producer. That flow is internal to consumer used and not mandated. Similar internal flow freedom is granted to transformer and producer. Like, if required your producer can have multiple outputs anywhere from Kafka, S3 to something like e-mail.

All 3, `consumer`, `transformer` and `producer` are instantiated as per config and thus any combination of available types could be brought into play.

Consumer, Transformer and Producer support usage of `golang plugin`, so separately managed and developed constructs could be used in combination as well.
Since they are loaded as per configuration provided identifier, one can have combination of say built-in Kafka consumer with built-in kubernetes-log transformer but custom external plug-in of Google Cloud Datastore for cold storage of logs.
This also gives capability to write a producer sending output to more than one output sinks in same flow to achieve replication.

---

* [available constructs](./docs/available-constructs.md) lists available consumer, transformer and producer

* [design](./docs/design.md) details the internal structure

* [example workflows available](./docs/example-workflows.md)

---

#### Quickstart

```
## create ogi-conf.env with required configurations in $PWD
## help could be taken from lnik above
#
## if any plugin used, that '*.so' should be present in $PWD

wget -c https://github.com/OpenChaos/ogi/releases/download/v1.0/ogi-linux-amd64

docker run -it --env-file ogi-conf.env $PWD:/opt/ogi ubuntu:16.04 /opt/ogi/ogi-linux-amd64
```

[set of configurations to make ogi work to specific behavior](./docs/config-set.md)

> _this uses [golang plugins](https://golang.org/pkg/plugin/) for extensibility, currently supported on linux, utilize docker to run if using something else_
> [ogi-v1.0-linux-amd64](https://github.com/OpenChaos/ogi/releases/download/v1.0/ogi-linux-amd64) could run from a Docker on non-linux platform

* latest release: [v1.0](https://github.com/OpenChaos/ogi/releases/tag/v1.0)

---

### What it was?

> A utility to enable flexible [ETL](https://en.wikipedia.org/wiki/Extract,_transform,_load) scenarios.
>
> Initially written to fan-out bulk topic `labels[app:appname]` tagged logs pushed from Kubernetes to Kafka, into `app` specific topics.
>
> Able to be scaled up using Kubernetes/Nomad/Mesos elastic deployments, inherently without context.

---
