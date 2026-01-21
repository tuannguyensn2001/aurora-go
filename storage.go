package core

import (
	"context"
	"time"

	"github.com/tuannguyensn2001/aurora-go/auroratype"
	"github.com/tuannguyensn2001/aurora-go/storage/memory"
)

type fetcherStorage struct {
	fetcher  Fetcher
	interval time.Duration
	strategy Storage
	recorder MetricsRecorder
}

func WithStorage(strategy Storage) func(s *fetcherStorage) {
	return func(s *fetcherStorage) {
		s.strategy = strategy
	}
}

func WithInterval(interval time.Duration) func(opts *fetcherStorage) {
	return func(opts *fetcherStorage) {
		opts.interval = interval
	}
}

func WithMetricsRecorder(recorder MetricsRecorder) func(opts *fetcherStorage) {
	return func(opts *fetcherStorage) {
		opts.recorder = recorder
	}
}

func NewFetcherStorage(fetcher Fetcher, opts ...func(opts *fetcherStorage)) *fetcherStorage {
	storage := &fetcherStorage{
		fetcher:  fetcher,
		interval: 1 * time.Minute,
		strategy: memory.NewStorage(),
		recorder: NewNoopRecorder(),
	}

	for _, opt := range opts {
		opt(storage)
	}

	return storage
}

func (w *fetcherStorage) Start(ctx context.Context) error {
	if err := w.sync(ctx); err != nil {
		return err
	}

	if w.fetcher.IsStatic() {
		return nil
	}

	go w.poll(ctx)

	return nil
}

func (w *fetcherStorage) sync(ctx context.Context) error {
	config, err := w.fetcher.Fetch(ctx)
	if err != nil {
		w.recorder.Count(MetricStorageSyncTotal, 1, []string{"status:error"})
		return err
	}

	err = w.strategy.Save(ctx, config)
	if err != nil {
		w.recorder.Count(MetricStorageSyncTotal, 1, []string{"status:error"})
		return err
	}

	experiments, err := w.fetcher.FetchExperiments(ctx)
	if err != nil {
		w.recorder.Count(MetricStorageSyncTotal, 1, []string{"status:error"})
		return err
	}

	if experiments != nil {
		err = w.strategy.SaveExperiments(ctx, experiments)
		if err != nil {
			w.recorder.Count(MetricStorageSyncTotal, 1, []string{"status:error"})
			return err
		}
	}

	w.recorder.Count(MetricStorageSyncTotal, 1, []string{"status:success"})
	return nil
}

func (w *fetcherStorage) poll(ctx context.Context) {
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

func (w *fetcherStorage) Get(ctx context.Context, parameterName string) (auroratype.Parameter, error) {
	return w.strategy.Get(ctx, parameterName)
}

func (w *fetcherStorage) Save(ctx context.Context, config map[string]auroratype.Parameter) error {
	return w.strategy.Save(ctx, config)
}

func (w *fetcherStorage) GetExperiments(ctx context.Context) ([]auroratype.Experiment, error) {
	return w.strategy.GetExperiments(ctx)
}

func (w *fetcherStorage) SaveExperiments(ctx context.Context, experiments []auroratype.Experiment) error {
	return w.strategy.SaveExperiments(ctx, experiments)
}
