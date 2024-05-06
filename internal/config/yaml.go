package config

import (
	"gopkg.in/yaml.v2"
)

func parseYaml(s string) (*ConfigRoot, error) {
	var config ConfigRoot
	if err := yaml.Unmarshal([]byte(s), &config); err != nil {
		return nil, err
	}
	return &config, nil
}
