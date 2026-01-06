package file

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/tuannguyensn2001/aurora-go/auroratype"
	"gopkg.in/yaml.v2"
)

// Fetcher fetches configuration from a local file.
type fetcher struct {
	filePath string
}

// Options configures the Fetcher.
type Options struct {
	FilePath string
}

// NewFetcher creates a new file-based Fetcher.
func New(opts Options) *fetcher {
	return &fetcher{
		filePath: opts.FilePath,
	}
}

// Fetch retrieves configuration data from a local file.
func (f *fetcher) Fetch(ctx context.Context) (map[string]auroratype.Parameter, error) {
	file, err := os.Open(f.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config map[string]auroratype.Parameter

	if strings.HasSuffix(f.filePath, ".yaml") || strings.HasSuffix(f.filePath, ".yml") {
		err = yaml.NewDecoder(file).Decode(&config)
	} else if strings.HasSuffix(f.filePath, ".json") {
		err = json.NewDecoder(file).Decode(&config)
	} else {
		return nil, errors.New("unsupported file format: must be .yaml, .yml, or .json")
	}

	if err != nil {
		return nil, err
	}

	return config, nil
}

func (f *fetcher) IsStatic() bool {
	return true
}
