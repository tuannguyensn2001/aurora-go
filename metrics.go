package core

// MetricsRecorder is an interface for recording metrics during GetParameter calls.
// Implement this interface to integrate with your observability backend.
//
// Example Prometheus implementation:
//
//	type prometheusRecorder struct {
//	    counter   *prometheus.CounterVec
//	    histogram prometheus.Histogram
//	}
//
//	func (p *prometheusRecorder) Count(metricName string, count int, tags []string) {
//	    if metricName != "get_parameter" {
//	        return
//	    }
//	    labels := make(prometheus.Labels)
//	    for _, tag := range tags {
//	        parts := strings.SplitN(tag, ":", 2)
//	        if len(parts) == 2 {
//	            labels[parts[0]] = parts[1]
//	        }
//	    }
//	    p.counter.With(labels).Add(float64(count))
//	}
//
//	func (p *prometheusRecorder) Histogram(metricName string, value float64, tags []string) {
//	    if metricName != "get_parameter_latency" {
//	        return
//	    }
//	    p.histogram.Observe(value)
//	}
type MetricsRecorder interface {
	Count(metricName string, count int, tags []string)
	Histogram(metricName string, value float64, tags []string)
}

type noopRecorder struct{}

func (n *noopRecorder) Count(metricName string, count int, tags []string)         {}
func (n *noopRecorder) Histogram(metricName string, value float64, tags []string) {}

func NewNoopRecorder() MetricsRecorder {
	return &noopRecorder{}
}

const (
	MetricStatusNotFound = "not_found"
	MetricStatusResolved = "resolved"
	MetricStatusFallback = "fallback"
)

const (
	MetricS3FetchLatency = "s3_fetch_latency"
	MetricS3FetchTotal   = "s3_fetch_total"

	MetricStorageSyncTotal   = "storage_sync_total"
	MetricStorageSaveLatency = "storage_save_latency"
	MetricStorageGetLatency  = "storage_get_latency"
	MetricStorageGetTotal    = "storage_get_total"
)
