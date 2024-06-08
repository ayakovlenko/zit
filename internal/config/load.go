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
func LocateConfFile(fs afero.Fs, userHomeDir, confPathFromEnv string) (string, error) {
	var confPath string

	yamlDefault := path.Join(userHomeDir, ".zit", "config.yaml")

	// if ZIT_CONFIG is not set, try default location
	envVarDefined := confPathFromEnv != ""

	if !envVarDefined {
		return yamlDefault, nil
	}

	if !fileExists(fs, confPathFromEnv) {
		return "", &ErrConfigNotFound{
			EnvVar: envVarDefined,
			Path:   confPath,
		}
	}

	return confPathFromEnv, nil
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
