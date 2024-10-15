package config

import (
	"fmt"
	"os"
	"path"
	"strings"

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
func LocateConfFile(fs afero.Fs, userHomeDir, confPathFromEnv string, xdgHomeFromEnv string) (string, error) {
	// try ZIT_CONFIG location
	if confPathFromEnv != "" {
		if fileExists(fs, confPathFromEnv) {
			return confPathFromEnv, nil
		}
		return "", &ErrConfigNotFound{
			EnvVar: true,
			Path:   "'" + confPathFromEnv + "'",
		}
	}

	// try default XDG-style locations:
	// - $XDG_CONFIG_HOME/zit/config.yaml (if XDG_CONFIG_HOME is set)
	// - $HOME/.config/zit/config.yaml
	if xdgHomeFromEnv == "" {
		xdgHomeFromEnv = path.Join(userHomeDir, ".config")
	}
	xdgYamlDefault := path.Join(xdgHomeFromEnv, "zit", "config.yaml")
	if fileExists(fs, xdgYamlDefault) {
		return xdgYamlDefault, nil
	}

	// try default dotfile location
	// $HOME/.zit/config.yaml
	yamlDefault := path.Join(userHomeDir, ".zit", "config.yaml")
	if fileExists(fs, yamlDefault) {
		return yamlDefault, nil
	}

	// we ran out of default options
	return "", &ErrConfigNotFound{
		EnvVar: false,
		Path:   "neither " + xdgYamlDefault + " nor " + yamlDefault,
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
