package metrics

import (
	"strings"
	"time"
)

// NS is a common namespace for the all metrics specified in nnm services
const NS = "sberbank"

// Labels for using them explicitly in metric calls for avoiding typos in the label values
const (
	LabelCode      = "code"
	LabelDB        = "db"
	LabelEndpoint  = "endpoint"
	LabelHandler   = "handler"
	LabelHost      = "host"
	LabelIsError   = "is_error"
	LabelMethod    = "method"
	LabelNamespace = "namespace"
	LabelOperation = "operation"
	LabelQuery     = "query"
	LabelSet       = "set"
	LabelStatus    = "status"
	LabelProtocol  = "protocol"
)

// Response status possible label values.
const (
	StatusOk              = "ok"
	StatusDegraded        = "degraded"
	StatusClientError     = "client_error"
	StatusError           = "error"
	StatusTimeout         = "timeout"
	StatusNoResponse      = "no_response"
	StatusCanceled        = "canceled"
	StatusTooManyRequests = "too_many_requests"
	StatusNotFound        = "not_found"
)

// HTTPCodeToStatus convert HTTP response code to the status string accordingly with nnm metrics standard.
func HTTPCodeToStatus(httpCode int) string {
	switch {
	case httpCode >= 500:
		return StatusError
	case httpCode >= 400 && httpCode != 404:
		return StatusClientError
	case httpCode == 404:
		return StatusNotFound
	case httpCode >= 200:
		return StatusOk
	}
	return StatusOk
}

// Seconds calculates the duration seconds instead on nanoseconds that default for time.Duration.
func Seconds(duration time.Duration) float64 {
	return float64(duration) / float64(time.Second)
}

// SinceSeconds just wraps time.Since() with converting result to seconds
func SinceSeconds(started time.Time) float64 {
	return float64(time.Since(started)) / float64(time.Second)
}

// IsError is a trivial helper for minimize repetitive checks for error
// values. It passing appropriate numbers to metrics.
func IsError(err error) string {
	if err != nil {
		return "1"
	}
	return "0"
}

// BoolToStr is a trivial helper for for minimize repetitive checks
// for boolean values. It passing appropriate numbers to metrics.
func BoolToStr(val bool) string {
	if val {
		return "1"
	}
	return "0"
}

// MakePathTag makes the tag for a metric accordingly with recommendations of Prometheus and for compatibility with Grafana.
// It replaces slashes with underscores and replaces empty path (usual for GET /) with "root" word.
func MakePathTag(path string) string {
	if path == "/" {
		return "root"
	}
	return strings.ToLower(strings.Replace(strings.Trim(path, "/"), "/", "_", -1))
}
