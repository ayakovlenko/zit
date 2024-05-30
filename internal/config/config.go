package config

import (
	"fmt"
)

// EnvVarName TODO
const EnvVarName = "ZIT_CONFIG"

// ErrConfigNotFound TODO
type ErrConfigNotFound struct {
	EnvVar bool
	Path   string
}

func (err *ErrConfigNotFound) Error() string {
	envVar := ""
	if err.EnvVar {
		envVar = " defined in " + EnvVarName + " variable"
	}
	return fmt.Sprintf("config file%s is not found at %q", envVar, err.Path)
}

// HostMap TODO
type HostMap map[string]HostConfig

// HostConfig TODO
type HostConfig struct {
	Default   *User      `json:"default" yaml:"default"`
	Overrides []Override `json:"overrides" yaml:"overrides"`
}

// User TODO
type User struct {
	Name    string   `json:"name" yaml:"name"`
	Email   string   `json:"email" yaml:"email"`
	Signing *Signing `json:"sign" yaml:"sign"`
}

type Signing struct {
	Key    string `json:"key" yaml:"key"`
	Format string `json:"format" yaml:"format"`
}

// Override TODO
type Override struct {
	Owner string `json:"owner" yaml:"owner"`
	Repo  string `json:"repo,omitempty" yaml:"repo"`
	User  User   `json:"user" yaml:"user"`
}

type ConfigRoot struct {
	Hosts map[string]HostConfig `yaml:"hosts"`
}

// Get TODO
func (c *ConfigRoot) Get(host string) (*HostConfig, error) {
	hostConf, ok := (c.Hosts)[host]
	if !ok {
		return nil, fmt.Errorf("cannot find config for host %q", host)
	}

	return &hostConf, nil
}
