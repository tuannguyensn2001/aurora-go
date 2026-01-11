package core

import (
	"context"
	"log/slog"
	"time"

	"github.com/tuannguyensn2001/aurora-go/auroratype"
)

// ClientOptions contains configuration options for the Client.
type ClientOptions struct {
	Logger          *slog.Logger
	MetricsRecorder MetricsRecorder
}

// ParameterOption is a functional option for GetParameter.
type ParameterOption func(*parameterOptions)

type parameterOptions struct {
	strategy Storage
}

// WithStrategy sets a custom storage strategy for retrieving the parameter.
// This is useful for use cases requiring strong consistency
// instead of the default eventually consistent storage.
//
// Example:
//
//	customStorage := myCustomStorage{}
//	client.GetParameter(ctx, "key", nil, WithStrategy(customStorage))
func WithStrategy(s Storage) ParameterOption {
	return func(o *parameterOptions) {
		o.strategy = s
	}
}

// Client is the main entry point for Aurora configuration management.
// It provides methods to retrieve parameters and register custom operators.
type Client struct {
	storage  *storage
	engine   *engine
	logger   *slog.Logger
	recorder MetricsRecorder
}

// NewClient creates a new Aurora client with the given storage and options.
func NewClient(storage *storage, opts ClientOptions) *Client {
	engine := newEngine()
	engine.bootstrap()

	logger := opts.Logger
	if logger == nil {
		logger = slog.Default()
	}

	recorder := opts.MetricsRecorder
	if recorder == nil {
		recorder = NewNoopRecorder()
	}

	return &Client{
		storage:  storage,
		engine:   engine,
		logger:   logger,
		recorder: recorder,
	}
}

// Start initializes the client by starting the storage layer.
func (c *Client) Start(ctx context.Context) error {
	c.logger.Info("Starting Aurora client")
	return c.storage.Start(ctx)
}

// GetParameter retrieves a parameter value based on the given attributes.
// It evaluates rules and constraints to determine the appropriate value.
//
// By default, parameters are retrieved from the default storage (eventually consistent).
// For strong consistency requirements, use WithStrategy to provide a custom storage:
//
//	customStorage := myCustomStorage{}
//	client.GetParameter(ctx, "key", nil, WithStrategy(customStorage))
func (c *Client) GetParameter(ctx context.Context, parameterName string, attribute *attribute, opts ...ParameterOption) *resolvedValue {
	c.logger.Debug("Getting parameter", "parameter", parameterName)

	start := time.Now()
	defer func() {
		duration := time.Since(start).Nanoseconds()
		c.recorder.Histogram("get_parameter_latency", float64(duration), []string{})
	}()

	var paramOpts parameterOptions
	if len(opts) > 0 {
		for _, opt := range opts {
			opt(&paramOpts)
		}
	}

	storageTag := "default"
	if paramOpts.strategy != nil {
		storageTag = "custom"
	}

	var config auroratype.Parameter
	var err error

	if paramOpts.strategy != nil {
		config, err = paramOpts.strategy.Get(ctx, parameterName)
	} else {
		config, err = c.storage.GetParameterConfig(ctx, parameterName)
	}

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

// RegisterOperator allows users to register a custom operator.
// The operator function takes two values (a, b) and returns a boolean result.
//
// Example:
//
//	client.RegisterOperator("startsWith", func(a, b any) bool {
//	    strA, okA := a.(string)
//	    strB, okB := b.(string)
//	    if !okA || !okB {
//	        return false
//	    }
//	    return strings.HasPrefix(strA, strB)
//	})
func (c *Client) RegisterOperator(name string, fn func(a, b any) bool) {
	c.logger.Info("Registering custom operator", "operator", name)
	c.engine.registerOperator(Operator(name), fn)
}
