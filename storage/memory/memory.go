package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/tuannguyensn2001/aurora-go/auroratype"
)

const (
	MetricStorageSaveLatency = "storage_save_latency"
	MetricStorageGetLatency  = "storage_get_latency"
	MetricStorageGetTotal    = "storage_get_total"
)

type MetricsRecorder interface {
	Count(metricName string, count int, tags []string)
	Histogram(metricName string, value float64, tags []string)
}

// Storage stores configuration in memory.
type strategy struct {
	config   map[string]auroratype.Parameter
	mu       sync.RWMutex
	recorder MetricsRecorder
}

// NewStorage creates a new in-memory storage.
func NewStrategy() *strategy {
	return &strategy{
		config:   make(map[string]auroratype.Parameter),
		recorder: &noopRecorder{},
	}
}

type noopRecorder struct{}

func (n *noopRecorder) Count(metricName string, count int, tags []string)         {}
func (n *noopRecorder) Histogram(metricName string, value float64, tags []string) {}

// Save stores the configuration in memory.
func (m *strategy) Save(ctx context.Context, config map[string]auroratype.Parameter) error {
	start := time.Now()
	defer func() {
		duration := float64(time.Since(start).Nanoseconds())
		m.recorder.Histogram(MetricStorageSaveLatency, duration, nil)
	}()

	m.mu.Lock()
	defer m.mu.Unlock()

	m.config = config
	return nil
}

// Get retrieves a parameter from memory.
func (m *strategy) Get(ctx context.Context, parameterName string) (auroratype.Parameter, error) {
	start := time.Now()
	defer func() {
		duration := float64(time.Since(start).Nanoseconds())
		m.recorder.Histogram(MetricStorageGetLatency, duration, nil)
	}()

	m.mu.RLock()
	defer m.mu.RUnlock()

	val, ok := m.config[parameterName]
	if !ok {
		m.recorder.Count(MetricStorageGetTotal, 1, []string{"status:miss"})
		return auroratype.Parameter{}, errors.New("parameter not found")
	}

	m.recorder.Count(MetricStorageGetTotal, 1, []string{"status:hit"})
	return val, nil
}
