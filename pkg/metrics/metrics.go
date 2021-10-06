package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	EventLogsSyncCounter prometheus.Counter
}

func SetupMetrics() *Metrics {
	m := &Metrics{
		EventLogsSyncCounter: promauto.NewCounter(prometheus.CounterOpts{
			Name: "aivenaudit_event_logs_synced_total",
			Help: "The total number of synchronized event logs",
		}),
	}

	prometheus.MustRegister(m.EventLogsSyncCounter)

	return m
}

func Handlers(mux *http.ServeMux) {
	mux.Handle("/metrics", promhttp.Handler())
}
