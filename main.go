package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type Exporter struct {
	MarathonURL string
	mutex       sync.RWMutex
	client      *http.Client

	taskCount *prometheus.GaugeVec
}

func NewExporter(marathonURL string) *Exporter {
	return &Exporter{
		MarathonURL: marathonURL,
		client:      http.DefaultClient,
		taskCount: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "marathon_task",
			Name:      "count",
			Help:      "Number of task running on Marathon",
		}, []string{"task"}),
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.taskCount.Describe(ch)
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	resp, err := http.Get(fmt.Sprintf("%s/v2/apps", e.MarathonURL))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	apps := &MarathonResponse{}
	json.NewDecoder(resp.Body).Decode(apps)

	metrics := make(map[string]*TaskMetric)
	for _, a := range apps.Apps {
		metrics[a.ID] = &TaskMetric{}
	}

	resp, err = http.Get(fmt.Sprintf("%s/v2/tasks", e.MarathonURL))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	tasks := &MarathonResponse{}
	json.NewDecoder(resp.Body).Decode(tasks)
	for _, t := range tasks.Tasks {
		metrics[t.ID].Count++
	}

	e.taskCount.Reset()
	for t, m := range metrics {
		e.taskCount.
			With(prometheus.Labels{"task": t}).
			Set(float64(m.Count))

		e.taskCount.Collect(ch)
	}
}

func main() {
	var (
		bindAddress = flag.String("web.listen-address", ":9091", "Address to listen on for HTTP interface")
		metricsPath = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics")
		marathonURL = flag.String("marathon.url", "", "Marathon instance URL")
	)
	flag.Parse()

	if *marathonURL == "" {
		log.Fatal("you should provide Marathon's URL")
	}

	exporter := NewExporter(*marathonURL)
	prometheus.MustRegister(exporter)

	http.Handle(*metricsPath, prometheus.Handler())
	log.Fatal(http.ListenAndServe(*bindAddress, nil))
}
