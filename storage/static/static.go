package static

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/tuannguyensn2001/aurora-go/core"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type staticStorage struct {
	filePath string
	config   map[string]core.Parameter
}

func NewStorage(filePath string) *staticStorage {

	return &staticStorage{
		filePath: filePath,
	}
}

func (s *staticStorage) start(ctx context.Context) error {

	f, err := os.Open(s.filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	var config map[string]core.Parameter
	if strings.HasSuffix(s.filePath, ".yaml") {
		err = yaml.NewDecoder(f).Decode(&config)
		if err != nil {
			return err
		}
	} else if strings.HasSuffix(s.filePath, ".json") {
		err = json.NewDecoder(f).Decode(&config)
		if err != nil {
			return err
		}
	} else {
		return errors.New("invalid file path")
	}

	s.config = config

	return nil
}
func (s *staticStorage) getParameterConfig(ctx context.Context, parameterName string) (core.Parameter, error) {
	val, ok := s.config[parameterName]
	if !ok {
		return core.Parameter{}, errors.New("parameter not found")
	}

	return val, nil

}
