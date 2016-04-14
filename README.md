# Prometheus Marathon task count exporter

This exporter for Prometheus only exposes one metric (fow now), the number of tasks for each app in a Marathon cluster. It has one dimension, the name of the app.

This exporter is not meant to report low-level metrics on Marathon. If you are looking for this, you should look at the [Mesos expoer](https://github.com/prometheus-junkyard/mesos_exporter) or the [Marathon exporter](https://github.com/prometheus-junkyard/mesos_exporter).

## Example output

```
marathon_task_count{task="/production/app-01"} 10
marathon_task_count{task="/production/app-02"} 3
marathon_task_count{task="/staging/app-03"} 0
```

## Build

```
$ go get github.com/apognu/marathon_task_exporter
$ go build github.com/apognu/marathon_task_exporter
```

## Usage

```
$ ./marathon_task_exporter -help
Usage of ./marathon_task_exporter:
  -marathon.url string
        Marathon instance URL
  -web.listen-address string
        Address to listen on for HTTP interface (default ":9091")
  -web.telemetry-path string
        Path under which to expose metrics (default "/metrics")

$ ./marathon_task_exporter -marathon.url=http://my.marathon.tld
```

### Alerting

You can use Prometheus's Alert Manager with this exporter, by first grouping by the _task_ label in your route configuration:

```
route:
  group_by: [ 'task' ]
```

An example alert rule could look like this:

```
ALERT ProductionAPIInstanceCount
  IF marathon_task_count{task = "/production/.+"} < 3
  FOR 30s
  ANNOTATIONS {
    summary = "Support Marathon Task count",
    description = "Number of tasks for *{{$labels.task}}* has changed -> *{{$value}}*"
  }
```

## What's next?

 * Try to include Mesos labels as metric dimensions
 * Take app health checks into account to count only live tasks
 * Handle HTTP authentication on Marathon
