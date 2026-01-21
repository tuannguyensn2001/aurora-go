package core

import (
	"context"

	"github.com/tuannguyensn2001/aurora-go/auroratype"
)

type Fetcher interface {
	Fetch(ctx context.Context) (map[string]auroratype.Parameter, error)
	FetchExperiments(ctx context.Context) ([]auroratype.Experiment, error)
	IsStatic() bool
}

type Storage interface {
	Save(ctx context.Context, config map[string]auroratype.Parameter) error
	Get(ctx context.Context, parameterName string) (auroratype.Parameter, error)
	SaveExperiments(ctx context.Context, experiments []auroratype.Experiment) error
	GetExperiments(ctx context.Context) ([]auroratype.Experiment, error)
}
