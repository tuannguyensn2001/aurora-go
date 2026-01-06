package static

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/tuannguyensn2001/aurora-go/auroratype"
	"gopkg.in/yaml.v2"
)

type staticStorage struct {
	filePath string
	config   map[string]auroratype.Parameter
}

func NewStorage(filePath string) *staticStorage {

	return &staticStorage{
		filePath: filePath,
	}
}

func (s *staticStorage) Start(ctx context.Context) error {

	f, err := os.Open(s.filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	var config map[string]auroratype.Parameter
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
func (s *staticStorage) GetParameterConfig(ctx context.Context, parameterName string) (auroratype.Parameter, error) {
	val, ok := s.config[parameterName]
	if !ok {
		return auroratype.Parameter{}, errors.New("parameter not found")
	}

	return val, nil

}
