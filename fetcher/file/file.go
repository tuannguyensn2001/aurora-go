package file

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/tuannguyensn2001/aurora-go/auroratype"
	"gopkg.in/yaml.v2"
)

type Fetcher struct {
	filePath            string
	experimentsFilePath string
	static              bool
}

type Options struct {
	FilePath            string
	ExperimentsFilePath string
	Static              bool
}

func New(opts Options) *Fetcher {
	return &Fetcher{
		filePath:            opts.FilePath,
		experimentsFilePath: opts.ExperimentsFilePath,
		static:              opts.Static,
	}
}

func (f *Fetcher) Fetch(ctx context.Context) (map[string]auroratype.Parameter, error) {
	if f.filePath == "" {
		return make(map[string]auroratype.Parameter), nil
	}

	file, err := os.Open(f.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config map[string]auroratype.Parameter

	ext := strings.ToLower(filepath.Ext(f.filePath))
	if ext == ".yaml" || ext == ".yml" {
		err = yaml.NewDecoder(file).Decode(&config)
	} else if ext == ".json" {
		err = json.NewDecoder(file).Decode(&config)
	} else {
		return nil, errors.New("unsupported file format: must be .yaml, .yml, or .json")
	}

	if err != nil {
		return nil, err
	}

	return config, nil
}

func (f *Fetcher) FetchExperiments(ctx context.Context) ([]auroratype.Experiment, error) {
	expFilePath := f.experimentsFilePath

	if expFilePath == "" && f.filePath != "" {
		dir := filepath.Dir(f.filePath)
		expFilePath = filepath.Join(dir, "experiments.yaml")
	}

	if expFilePath == "" {
		return nil, nil
	}

	data, err := os.ReadFile(expFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var config struct {
		Experiments []auroratype.Experiment `yaml:"experiments"`
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return config.Experiments, nil
}

func (f *Fetcher) IsStatic() bool {
	return f.static
}
