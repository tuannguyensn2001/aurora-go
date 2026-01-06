package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/tuannguyensn2001/aurora-go/auroratype"
)

// Storage stores configuration in memory.
type strategy struct {
	config map[string]auroratype.Parameter
	mu     sync.RWMutex
}

// NewStorage creates a new in-memory storage.
func NewStrategy() *strategy {
	return &strategy{
		config: make(map[string]auroratype.Parameter),
	}
}

// Save stores the configuration in memory.
func (m *strategy) Save(ctx context.Context, config map[string]auroratype.Parameter) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.config = config
	return nil
}

// Get retrieves a parameter from memory.
func (m *strategy) Get(ctx context.Context, parameterName string) (auroratype.Parameter, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	val, ok := m.config[parameterName]
	if !ok {
		return auroratype.Parameter{}, errors.New("parameter not found")
	}

	return val, nil
}
