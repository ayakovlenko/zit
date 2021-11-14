package config

import (
	"gopkg.in/yaml.v2"
)

func parseYaml(s string) (*ConfigV2, error) {
	var config ConfigV2
	if err := yaml.Unmarshal([]byte(s), &config); err != nil {
		return nil, err
	}
	return &config, nil
}
