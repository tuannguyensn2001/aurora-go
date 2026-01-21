package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/tuannguyensn2001/aurora-go/auroratype"
)

type MetricsRecorder interface {
	Count(metricName string, count int, tags []string)
	Histogram(metricName string, value float64, tags []string)
}

type Storage struct {
	config      map[string]auroratype.Parameter
	experiments []auroratype.Experiment
	mu          sync.RWMutex
	recorder    MetricsRecorder
}

func NewStorage() *Storage {
	return &Storage{
		config:      make(map[string]auroratype.Parameter),
		experiments: make([]auroratype.Experiment, 0),
		recorder:    &noopRecorder{},
	}
}

type noopRecorder struct{}

func (n *noopRecorder) Count(metricName string, count int, tags []string)         {}
func (n *noopRecorder) Histogram(metricName string, value float64, tags []string) {}

func (m *Storage) Save(ctx context.Context, config map[string]auroratype.Parameter) error {
	start := time.Now()
	defer func() {
		duration := float64(time.Since(start).Nanoseconds())
		m.recorder.Histogram("storage_save_latency", duration, nil)
	}()

	m.mu.Lock()
	defer m.mu.Unlock()

	m.config = config
	return nil
}

func (m *Storage) Get(ctx context.Context, parameterName string) (auroratype.Parameter, error) {
	start := time.Now()
	defer func() {
		duration := float64(time.Since(start).Nanoseconds())
		m.recorder.Histogram("storage_get_latency", duration, nil)
	}()

	m.mu.RLock()
	defer m.mu.RUnlock()

	val, ok := m.config[parameterName]
	if !ok {
		m.recorder.Count("storage_get_total", 1, []string{"status:miss"})
		return auroratype.Parameter{}, errors.New("parameter not found")
	}

	m.recorder.Count("storage_get_total", 1, []string{"status:hit"})
	return val, nil
}

func (m *Storage) SaveExperiments(ctx context.Context, experiments []auroratype.Experiment) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.experiments = experiments
	return nil
}

func (m *Storage) GetExperiments(ctx context.Context) ([]auroratype.Experiment, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.experiments, nil
}

func NewStrategy() *Storage {
	return NewStorage()
}
