package achim

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func ParseYAMLFile[T any](path string) (*T, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf(`read YAML from file "%s": %w`, path, err)
	}
	var object T
	if err = yaml.Unmarshal(raw, &object); err != nil {
		return nil, fmt.Errorf(`parse YAML file "%s": %w`, path, err)
	}
	return &object, nil
}
