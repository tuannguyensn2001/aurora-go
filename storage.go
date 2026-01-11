package core

import (
	"context"
	"time"

	"github.com/tuannguyensn2001/aurora-go/auroratype"
	"github.com/tuannguyensn2001/aurora-go/storage/memory"
)

// WrapperStorage combines a Fetcher and Storage to create a complete storage solution.
// It fetches data from the Fetcher and stores it in the Storage.
type storage struct {
	fetcher  Fetcher
	interval time.Duration
	strategy Storage
}

func WithStorage(strategy Storage) func(storage *storage) {
	return func(storage *storage) {
		storage.strategy = strategy
	}
}

func WithInterval(interval time.Duration) func(opts *storage) {
	return func(opts *storage) {
		opts.interval = interval
	}
}

// NewWrapperStorage creates a new WrapperStorage that combines fetching and storage strategies.
func NewStorage(fetcher Fetcher, opts ...func(opts *storage)) *storage {
	storage := &storage{
		fetcher:  fetcher,
		interval: 1 * time.Minute,
		strategy: memory.NewStrategy(),
	}

	for _, opt := range opts {
		opt(storage)
	}

	return storage
}

// Start initializes the storage by fetching data and starting background polling.
func (w *storage) Start(ctx context.Context) error {
	// Initial fetch and save
	if err := w.sync(ctx); err != nil {
		return err
	}

	// Start background polling
	if w.fetcher.IsStatic() {
		return nil
	}

	go w.poll(ctx)

	return nil
}

func (w *storage) sync(ctx context.Context) error {
	config, err := w.fetcher.Fetch(ctx)
	if err != nil {
		return err
	}

	return w.strategy.Save(ctx, config)
}

func (w *storage) poll(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_ = w.sync(ctx)
		}
	}
}

// GetParameterConfig retrieves a parameter from the storage.
func (w *storage) GetParameterConfig(ctx context.Context, parameterName string) (auroratype.Parameter, error) {
	return w.strategy.Get(ctx, parameterName)
}
