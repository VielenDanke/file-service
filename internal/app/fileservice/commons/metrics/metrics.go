package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	DefaultMetricsPath = "/metrics"
)

var (
	registerer = prometheus.DefaultRegisterer
	gatherer   = prometheus.DefaultGatherer

	DefaultSummaryObjectives = map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}
	vc                       = newCache()
)

// Handler returns HTTP handler for serving metrics
func Handler() http.Handler {
	return promhttp.Handler()
}

// Gatherer returns default gatherer interface from the Prometheus.
func Gatherer() prometheus.Gatherer {
	return gatherer
}

// Registerer returns default registry from the Prometheus.
// It is the most common way to use the client.
func Registerer() prometheus.Registerer {
	return registerer
}

// MustRegister registers a new metric in the registry. If the metric fail to register it will panic.
func MustRegister(collectors ...prometheus.Collector) {
	registerer.MustRegister(collectors...)
}

// Unregister metric from the registry.
func Unregister(collector prometheus.Collector) bool {
	return registerer.Unregister(collector)
}

// NewCounter creates a new Counter with predefined namespace
func NewCounter(name, help string) prometheus.Counter {
	return prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: NS,
			Name:      name,
			Help:      help,
		},
	)
}

// NewGauge creates a new Gauge with predefined namespace
func NewGauge(name, help string) prometheus.Gauge {
	return prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: NS,
			Name:      name,
			Help:      help,
		},
	)
}

// NewHistogram creates a new Histogram with predefined namespace
func NewHistogram(name, help string, buckets []float64) prometheus.Histogram {
	return prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: NS,
			Name:      name,
			Help:      help,
			Buckets:   buckets,
		},
	)
}

// NewSummary creates a new Summary with predefined namespace
func NewSummary(name, help string) prometheus.Summary {
	return prometheus.NewSummary(
		prometheus.SummaryOpts{
			Namespace:  NS,
			Name:       name,
			Help:       help,
			MaxAge:     time.Minute,
			Objectives: DefaultSummaryObjectives,
		},
	)
}

// NewCounterVec creates a new CounterVec with predefined namespace
func NewCounterVec(name, help string, labelValues []string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: NS,
			Name:      name,
			Help:      help,
		},
		labelValues,
	)
}

// NewGaugeVec creates a new GaugeVec with predefined namespace
func NewGaugeVec(name, help string, labelValues []string) *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: NS,
			Name:      name,
			Help:      help,
		},
		labelValues,
	)
}

// NewHistogramVec creates a new HistogramVec with predefined namespace
func NewHistogramVec(name, help string, buckets []float64, labelValues []string) *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: NS,
			Name:      name,
			Help:      help,
			Buckets:   buckets,
		},
		labelValues,
	)
}

// NewSummaryVec creates a new SummaryVec with predefined namespace
func NewSummaryVec(name, help string, labelValues []string) *prometheus.SummaryVec {
	return prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  NS,
			Name:       name,
			Help:       help,
			MaxAge:     time.Minute,
			Objectives: DefaultSummaryObjectives,
		},
		labelValues,
	)
}

// MustRegisterCounter creates and registers new Counter with predefined namespace
// Panics if metrics with same name already registered
func MustRegisterCounter(name, help string) prometheus.Counter {
	collector := NewCounter(name, help)
	MustRegister(collector)

	return collector
}

// MustRegisterGauge creates and registers new Gauge with predefined namespace
// Panics if metrics with same name already registered
func MustRegisterGauge(name, help string) prometheus.Gauge {
	collector := NewGauge(name, help)
	MustRegister(collector)

	return collector
}

// MustRegisterHistogram creates and registers new Histogram with predefined namespace
// Panics if metrics with same name already registered
func MustRegisterHistogram(name, help string, buckets []float64) prometheus.Histogram {
	collector := NewHistogram(name, help, buckets)
	MustRegister(collector)

	return collector
}

// MustRegisterSummary creates and registers new Summary with predefined namespace
// Panics if metrics with same name already registered
func MustRegisterSummary(name, help string) prometheus.Summary {
	collector := NewSummary(name, help)
	MustRegister(collector)

	return collector
}

// MustRegisterCounterVec creates and registers new CounterVec with predefined namespace
// Panics if metrics with same name already registered
func MustRegisterCounterVec(name, help string, labelValues []string) *prometheus.CounterVec {
	collector := NewCounterVec(name, help, labelValues)
	MustRegister(collector)

	return collector
}

// MustRegisterGaugeVec creates and registers new GaugeVec with predefined namespace
// Panics if metrics with same name already registered
func MustRegisterGaugeVec(name, help string, labelValues []string) *prometheus.GaugeVec {
	collector := NewGaugeVec(name, help, labelValues)
	MustRegister(collector)

	return collector
}

// MustRegisterHistogramVec creates and registers new HistogramVec with predefined namespace
// Panics if metrics with same name already registered
func MustRegisterHistogramVec(name, help string, buckets []float64, labelValues []string) *prometheus.HistogramVec {
	collector := NewHistogramVec(name, help, buckets, labelValues)
	MustRegister(collector)

	return collector
}

// MustRegisterSummaryVec creates and registers new SummaryVec with predefined namespace
// Panics if metrics with same name already registered
func MustRegisterSummaryVec(name, help string, labelValues []string) *prometheus.SummaryVec {
	collector := NewSummaryVec(name, help, labelValues)
	MustRegister(collector)

	return collector
}
