package app

import (
	"io/fs"
	"os"
)

// FS is the filesystem interface used by the app.
type FS interface {
	Stat(name string) (fs.FileInfo, error)
	Create(name string) (*os.File, error)
	MkdirAll(path string, perm os.FileMode) error
	ReadFile(name string) ([]byte, error)
}

// OsFS is an FS backed by the real operating system.
type OsFS struct{}

func (OsFS) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(name)
}

func (OsFS) Create(name string) (*os.File, error) {
	return os.Create(name)
}

func (OsFS) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (OsFS) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

type Config interface {
	AppName() string
	AppVersion() string
	ConfigFilename() string
	FS() FS
	UserHomeDir() string
	ConfigPathFromEnv() string
	XDGHomePathFromEnv() string
}

const (
	appName        = "zit"
	appVersion     = "v3.1.2"
	configFilename = "config.yaml"
)

func NewConfig( //nolint: ireturn
	fs FS,
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
	fs                 FS
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

func (c *config) FS() FS { //nolint: ireturn
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
