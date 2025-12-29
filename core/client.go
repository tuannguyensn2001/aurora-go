package core

import (
	"context"
	"log/slog"
)

type iStorage interface {
	getParameterConfig(ctx context.Context, parameterName string) (Parameter, error)
	start(ctx context.Context) error
}

type client struct {
	storage iStorage
	engine  *engine
	logger  *slog.Logger
}

type ClientOptions struct {
	Logger *slog.Logger
}

func NewClient(storage iStorage, opts ClientOptions) *client {
	logger := slog.Default()
	if opts.Logger != nil {
		logger = opts.Logger
	}
	engine := newEngine()
	return &client{
		storage: storage,
		engine:  engine,
		logger:  logger,
	}
}

type parameterOpts struct {
}

type parameterOption func(*parameterOpts)

func (c *client) Start(ctx context.Context) error {
	err := c.storage.start(ctx)
	if err != nil {
		return err
	}
	c.engine.bootstrap()
	return nil
}

func (c *client) GetParameter(ctx context.Context, parameterName string, attribute *attribute, opts ...parameterOption) *resolvedValue {
	parameter, err := c.storage.getParameterConfig(ctx, parameterName)
	if err != nil {
		c.logger.Error("failed to get parameter from storage", "error", err)
		return NewResolvedValue(nil, false)
	}

	return c.engine.evaluateParameter(ctx, parameterName, parameter, attribute)
}
