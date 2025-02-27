package app

import "github.com/spf13/afero"

type Config interface {
	AppName() string
	AppVersion() string
	ConfigFilename() string
	FS() afero.Fs
	UserHomeDir() string
	ConfigPathFromEnv() string
	XDGHomePathFromEnv() string
}

const (
	appName        = "zit"
	appVersion     = "v3.1.1"
	configFilename = "config.yaml"
)

func NewConfig( //nolint: ireturn
	fs afero.Fs,
	userHomeDir string,
	configPathFromEnv string,
	xdgHomePathFromEnv string,
) Config {
	return &config{
		fs:                 fs,
		userHomeDir:        userHomeDir,
		configPathFromEnv:  configPathFromEnv,
		xdgHomePathFromEnv: xdgHomePathFromEnv,
	}
}

type config struct {
	fs                 afero.Fs
	userHomeDir        string
	configPathFromEnv  string
	xdgHomePathFromEnv string
}

func (c *config) AppName() string {
	return appName
}

func (c *config) AppVersion() string {
	return appVersion
}

func (c *config) ConfigFilename() string {
	return configFilename
}

func (c *config) FS() afero.Fs { //nolint: ireturn
	return c.fs
}

func (c *config) UserHomeDir() string {
	return c.userHomeDir
}

func (c *config) ConfigPathFromEnv() string {
	return c.configPathFromEnv
}

func (c *config) XDGHomePathFromEnv() string {
	return c.xdgHomePathFromEnv
}
