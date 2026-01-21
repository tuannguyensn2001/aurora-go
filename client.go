package core

import (
	"context"
	"log/slog"
	"time"

	"github.com/tuannguyensn2001/aurora-go/auroratype"
	"github.com/tuannguyensn2001/aurora-go/core/evaluator"
	"github.com/tuannguyensn2001/aurora-go/experiment"
)

type ClientOptions struct {
	Logger          *slog.Logger
	MetricsRecorder MetricsRecorder
}

type ParameterOption func(*parameterOptions)

type parameterOptions struct {
	strategy Storage
}

func WithStrategy(s Storage) ParameterOption {
	return func(o *parameterOptions) {
		o.strategy = s
	}
}

type Client struct {
	storage          *fetcherStorage
	engine           *engine
	experimentEngine *experiment.Engine
	logger           *slog.Logger
	recorder         MetricsRecorder
}

func NewClient(storage *fetcherStorage, opts ClientOptions) *Client {
	eng := newEngine()
	eng.bootstrap()

	expEngine := experiment.NewEngine()
	expEngine.Bootstrap()

	logger := opts.Logger
	if logger == nil {
		logger = slog.Default()
	}

	recorder := opts.MetricsRecorder
	if recorder == nil {
		recorder = NewNoopRecorder()
	}

	if storage.logger == nil {
		storage.logger = logger
	}

	return &Client{
		storage:          storage,
		engine:           eng,
		experimentEngine: expEngine,
		logger:           logger,
		recorder:         recorder,
	}
}

func (c *Client) Start(ctx context.Context) error {
	c.logger.Info("Starting Aurora client")
	return c.storage.Start(ctx)
}

func (c *Client) GetParameter(ctx context.Context, parameterName string, attribute *attribute, opts ...ParameterOption) *resolvedValue {
	c.logger.Debug("Getting parameter", "parameter", parameterName)

	start := time.Now()
	defer func() {
		duration := time.Since(start).Nanoseconds()
		c.recorder.Histogram("get_parameter_latency", float64(duration), []string{})
	}()

	storageTag := "default"
	var paramOpts parameterOptions
	if len(opts) > 0 {
		for _, opt := range opts {
			opt(&paramOpts)
		}
	}
	if paramOpts.strategy != nil {
		storageTag = "custom"
	}

	var strg Storage
	if paramOpts.strategy != nil {
		strg = paramOpts.strategy
	} else {
		strg = c.storage
	}

	if c.experimentEngine != nil {
		experiments, err := strg.GetExperiments(ctx)
		if err == nil && len(experiments) > 0 {
			attrMap := make(map[string]any)
			if attribute != nil {
				for k, v := range attribute.vals {
					attrMap[k] = v
				}
			}

			result := c.experimentEngine.Evaluate(ctx, experiments, parameterName, attrMap)
			if result.Matched {
				value := result.Values[parameterName]
				c.recorder.Count("experiment_matched", 1, []string{"experiment:" + result.ExperimentID, "variant:" + result.VariantKey})
				return NewResolvedValue(value, true)
			}
		}
	}

	var config auroratype.Parameter
	var err error

	config, err = strg.Get(ctx, parameterName)

	if err != nil {
		c.logger.Error("Failed to get parameter config", "parameter", parameterName, "error", err)
		c.recorder.Count("get_parameter", 1, []string{"status:not_found", "storage:" + storageTag})
		return NewResolvedValue(nil, false)
	}

	result := c.engine.evaluateParameter(ctx, parameterName, config, attribute)

	if result.matched {
		c.recorder.Count("get_parameter", 1, []string{"status:resolved", "storage:" + storageTag})
	} else {
		c.recorder.Count("get_parameter", 1, []string{"status:fallback", "storage:" + storageTag})
	}

	return result
}

func (c *Client) RegisterOperator(name string, fn func(a, b any) bool) {
	c.logger.Info("Registering custom operator", "operator", name)
	c.engine.registerOperator(evaluator.Operator(name), fn)
}
