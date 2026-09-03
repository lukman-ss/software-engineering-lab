package observability

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Standard latency buckets suitable for sub-millisecond up to several seconds
	httpDurationBuckets = []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10}
	depDurationBuckets  = []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10}
)

type PrometheusCollector struct {
	HTTPRequestsTotal         *prometheus.CounterVec
	HTTPRequestDuration       *prometheus.HistogramVec
	HTTPRequestErrorsTotal    *prometheus.CounterVec
	HTTPInFlightRequests      *prometheus.GaugeVec
	DependencyDurationSeconds *prometheus.HistogramVec
	DependencyErrorsTotal     *prometheus.CounterVec
}

func NewPrometheusCollector(reg prometheus.Registerer) *PrometheusCollector {
	factory := promauto.With(reg)

	return &PrometheusCollector{
		HTTPRequestsTotal: factory.NewCounterVec(
			prometheus.CounterOpts{
				Name: "lab07_http_requests_total",
				Help: "Total number of HTTP requests processed by method, route, and status class.",
			},
			[]string{"method", "route", "status_class"},
		),
		HTTPRequestDuration: factory.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "lab07_http_request_duration_seconds",
				Help:    "HTTP request duration in seconds.",
				Buckets: httpDurationBuckets,
			},
			[]string{"method", "route", "status_class"},
		),
		HTTPRequestErrorsTotal: factory.NewCounterVec(
			prometheus.CounterOpts{
				Name: "lab07_http_request_errors_total",
				Help: "Total number of HTTP request errors by method, route, and status class.",
			},
			[]string{"method", "route", "status_class"},
		),
		HTTPInFlightRequests: factory.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "lab07_http_in_flight_requests",
				Help: "Current number of in-flight HTTP requests.",
			},
			[]string{"route"},
		),
		DependencyDurationSeconds: factory.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "lab07_dependency_duration_seconds",
				Help:    "Dependency operation latency in seconds.",
				Buckets: depDurationBuckets,
			},
			[]string{"component", "operation", "outcome"},
		),
		DependencyErrorsTotal: factory.NewCounterVec(
			prometheus.CounterOpts{
				Name: "lab07_dependency_errors_total",
				Help: "Total number of dependency errors by component, operation, and outcome.",
			},
			[]string{"component", "operation", "outcome"},
		),
	}
}
