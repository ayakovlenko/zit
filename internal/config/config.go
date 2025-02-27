package config

import (
	"fmt"
)

// EnvVarName TODO
const EnvVarName = "ZIT_CONFIG"

// ConfigNotFoundError TODO
type ConfigNotFoundError struct {
	EnvVar bool
	Path   string
}

func (err *ConfigNotFoundError) Error() string {
	envVar := ""
	if err.EnvVar {
		envVar = " defined in " + EnvVarName + " variable"
	}

	return fmt.Sprintf("config file%s is not found at %s", envVar, err.Path)
}

// HostMap TODO
type HostMap map[string]HostConfig

// HostConfig TODO
type HostConfig struct {
	Default   *User      `yaml:"default"`
	Overrides []Override `yaml:"overrides"`
}

// User TODO
type User struct {
	Name    string   `yaml:"name"`
	Email   string   `yaml:"email"`
	Signing *Signing `yaml:"sign"`
}

type Signing struct {
	Key    string `yaml:"key"`
	Format string `yaml:"format"`
}

// Override TODO
type Override struct {
	Owner string `yaml:"owner"`
	Repo  string `yaml:"repo"`
	User  User   `yaml:"user"`
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
