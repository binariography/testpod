package http

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusMiddleware struct {
	counter  *prometheus.CounterVec
	duration *prometheus.HistogramVec
}

func NewMetricMiddleware() *PrometheusMiddleware {

	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: "http",
			Name:      "requests_total",
			Help:      "The total number of HTTP requests.",
		},
		[]string{"status"},
	)

	// duration is partitioned by the HTTP method and handler. It uses custom
	// buckets based on the expected request duration.
	duration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_duration_seconds",
			Help:    "A histogram of latencies for requests.",
			Buckets: []float64{.25, .5, 1, 2.5, 5, 10},
		},
		[]string{"handler", "method"},
	)

	prometheus.MustRegister(counter, duration)

	return &PrometheusMiddleware{
		counter:  counter,
		duration: duration,
	}

}

func (p *PrometheusMiddleware) Handler(next http.Handler) http.Handler {
	middleware := promhttp.InstrumentHandlerDuration(p.duration.MustCurryWith(prometheus.Labels{"handler": "pull"}), promhttp.InstrumentHandlerCounter(p.counter,
		next))

	return middleware
}
