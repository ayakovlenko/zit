package config

import "fmt"

type ConfigV2 struct {
	Hosts map[string]HostV2 `yaml:"hosts"`
}

type HostV2 struct {
	Default   *User      `yaml:"default"`
	Overrides []Override `yaml:"overrides"`
}

// Get TODO
func (c *ConfigV2) Get(host string) (*HostV2, error) {
	hostConf, ok := (*&c.Hosts)[host]
	if !ok {
		return nil, fmt.Errorf("cannot find config for host %q", host)
	}

	return &hostConf, nil
}
