package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var EventLogsSyncCounter = promauto.NewCounter(prometheus.CounterOpts{
	Name: "aivenaudit_event_logs_synced_total",
	Help: "The total number of synchronized event logs",
})

func SetupPrometheus() {
	http.Handle("/metrics", promhttp.Handler())
}
