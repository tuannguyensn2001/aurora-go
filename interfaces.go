package core

//go:generate $GOPATH/bin/mockery --name=Fetcher
//go:generate $GOPATH/bin/mockery --name=Storage

import (
	"context"

	"github.com/tuannguyensn2001/aurora-go/auroratype"
)

// Fetcher is responsible for fetching configuration data from a source.
// Implementations: aurorafetcher.S3Fetcher, aurorafetcher.FileFetcher
type Fetcher interface {
	Fetch(ctx context.Context) (map[string]auroratype.Parameter, error)
	IsStatic() bool
}

// Storage is responsible for storing and retrieving configuration data.
// Implementations: aurorastorage.MemoryStorage
type Storage interface {
	Save(ctx context.Context, config map[string]auroratype.Parameter) error
	Get(ctx context.Context, parameterName string) (auroratype.Parameter, error)
}
