package middleware

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/unistack-org/micro/v3/logger"
	metrics "github.com/vielendanke/file-service/internal/app/fileservice/commons/metrics"
)

var (
	mw *incomingInstrumentation
)

func init() {
	mw = &incomingInstrumentation{
		duration: metrics.GetOrMakeHistogramVec(
			prometheus.HistogramOpts{
				Namespace: metrics.NS,
				Name:      "external_service_response_time_seconds",
				Help:      "Request time duration.",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"method", "protocol", "handler"},
		),
		requests: metrics.GetOrMakeCounterVec(
			prometheus.CounterOpts{
				Namespace: metrics.NS,
				Name:      "external_service_requests_total",
				Help:      "Total number of requests received.",
			},
			[]string{"code", "method", "protocol", "handler"},
		),
		requestSize: metrics.GetOrMakeHistogramVec(
			prometheus.HistogramOpts{
				Namespace: metrics.NS,
				Name:      "external_service_requests_size_histogram_bytes",
				Help:      "Request size in bytes.",
				Buckets:   []float64{100, 1000, 2000, 5000, 10000},
			},
			[]string{"method", "protocol", "handler"},
		),
		responseSize: metrics.GetOrMakeHistogramVec(
			prometheus.HistogramOpts{
				Namespace: metrics.NS,
				Name:      "external_service_response_size_histogram_bytes",
				Help:      "Response size in bytes.",
				Buckets:   []float64{100, 1000, 2000, 5000, 10000},
			},
			[]string{"code", "method", "protocol", "handler"},
		),
		inflight: metrics.GetOrMakeGauge(
			prometheus.GaugeOpts{
				Namespace: metrics.NS,
				Name:      "external_service_requests_in_flight",
				Help:      "Number of http requests which are currently running.",
			},
		),
	}
}

func HttpMetricsWrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler := "unknown"
		route := mux.CurrentRoute(r)
		handler, err := route.GetPathTemplate()
		if err != nil {
			logger.Info(r.Context(), "Cannot get PathTemplate: %v", err)
		}

		protocol := r.Proto
		method := r.Method
		mw.requestSize.WithLabelValues(method, protocol, handler).Observe(computeApproximateRequestSize(r))
		timer := prometheus.NewTimer(mw.duration.WithLabelValues(method, protocol, handler))
		sw := &statusWriter{http.StatusOK, w, w.(http.Hijacker)}
		mw.inflight.Inc()
		defer mw.inflight.Dec()
		next.ServeHTTP(sw, r)
		// TODO: not implemented	:	ins.responseSize
		timer.ObserveDuration()
		mw.requests.WithLabelValues(fmt.Sprintf("%d", sw.status), method, protocol, handler).Inc()
	})
}

type statusWriter struct {
	status int
	http.ResponseWriter
	http.Hijacker
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func computeApproximateRequestSize(r *http.Request) float64 {
	s := 0
	if r.URL != nil {
		s += len(r.URL.String())
	}

	s += len(r.Method)
	s += len(r.Proto)
	for name, values := range r.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}
	s += len(r.Host)

	if r.ContentLength != -1 {
		s += int(r.ContentLength)
	}

	return float64(s)
}
