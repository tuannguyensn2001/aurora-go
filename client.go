package core

import (
	"context"
	"log/slog"
)

type client struct {
	storage *storage
	engine  *engine
	logger  *slog.Logger
}

type ClientOptions struct {
	Logger *slog.Logger
}

func NewClient(storage *storage, opts ClientOptions) *client {
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
	err := c.storage.Start(ctx)
	if err != nil {
		return err
	}
	c.engine.bootstrap()
	return nil
}

func (c *client) GetParameter(ctx context.Context, parameterName string, attribute *attribute, opts ...parameterOption) *resolvedValue {
	parameter, err := c.storage.GetParameterConfig(ctx, parameterName)
	if err != nil {
		c.logger.Error("failed to get parameter from storage", "error", err)
		return NewResolvedValue(nil, false)
	}

	return c.engine.evaluateParameter(ctx, parameterName, parameter, attribute)
}
