package aurora

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type staticStorage struct {
	filePath string
	config   map[string]parameter
}

func NewStaticStorage(filePath string) *staticStorage {

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
	var config map[string]parameter
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
func (s *staticStorage) getParameterConfig(ctx context.Context, parameterName string) (parameter, error) {
	val, ok := s.config[parameterName]
	if !ok {
		return parameter{}, errors.New("parameter not found")
	}

	return val, nil

}
