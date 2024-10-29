package config

import (
	"fmt"
	"os"
	"path"
	"strings"
	"zit/internal/app"
	"zit/pkg/xdg"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

const (
	yamlFormat  = "yaml"
	otherFormat = "other"
)

var ErrUnsupportedConfigFormat = fmt.Errorf("unsupported config format")

func Load(filename string) (*ConfigRoot, error) {
	format := formatFromFilename(filename)

	if format != yamlFormat {
		return nil, ErrUnsupportedConfigFormat
	}

	switch format {
	case yamlFormat:
		contents, err := os.ReadFile(filename)
		if err != nil {
			return nil, err
		}

		return parseYaml(contents)
	default:
		return nil, fmt.Errorf("something went horribly wrong")
	}
}

func formatFromFilename(filename string) string {
	if strings.HasSuffix(filename, ".yaml") {
		return yamlFormat
	}

	return otherFormat
}

// LocateConfFile locates the path of the configuration file.
func LocateConfFile(appConfig app.Config) (string, error) {
	fs := appConfig.FS()

	userHomeDir := appConfig.UserHomeDir()

	confPathFromEnv := appConfig.ConfigPathFromEnv()

	xdgHomeFromEnv := appConfig.XDGHomePathFromEnv()

	// try ZIT_CONFIG location
	if confPathFromEnv != "" {
		if fileExists(fs, confPathFromEnv) {
			return confPathFromEnv, nil
		}

		return "", &ConfigNotFoundError{
			EnvVar: true,
			Path:   "'" + confPathFromEnv + "'",
		}
	}

	xdgConfigFile := xdg.LocateConfig(
		appConfig.AppName(),
		userHomeDir,
		xdgHomeFromEnv,
		appConfig.ConfigFilename(),
	)

	if fileExists(fs, xdgConfigFile) {
		return xdgConfigFile, nil
	}

	// try default dotfile location
	// $HOME/.zit/config.yaml
	yamlDefault := path.Join(userHomeDir, "."+appConfig.AppName(), "config.yaml")
	if fileExists(fs, yamlDefault) {
		return yamlDefault, nil
	}

	// we ran out of default options
	return "", &ConfigNotFoundError{
		EnvVar: false,
		Path:   "neither " + xdgConfigFile + " nor " + yamlDefault,
	}
}

func fileExists(fs afero.Fs, filename string) bool {
	fileInfo, err := fs.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	if fileInfo.IsDir() {
		return false
	}

	return true
}

func parseYaml(contents []byte) (*ConfigRoot, error) {
	var config ConfigRoot
	if err := yaml.Unmarshal(contents, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
