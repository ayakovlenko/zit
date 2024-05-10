package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/google/go-jsonnet"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

const (
	yamlFormat    = "yaml"
	jsonnetFormat = "jsonnet"
	otherFormat   = "other"
)

var ErrUnsupportedConfigFormat = fmt.Errorf("unsupported config format")

func Load(filename string) (*ConfigRoot, error) {
	format := formatFromFilename(filename)

	if format != yamlFormat && format != jsonnetFormat {
		return nil, ErrUnsupportedConfigFormat
	}

	switch format {
	case yamlFormat:
		contents, err := os.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		return parseYaml(contents)
	case jsonnetFormat:
		fmt.Println("WARN: Jsonnet configs are deprecated and going to be unsupported in future versions. Migrate to YAML format.")
		return parseJsonnet(filename)
	default:
		return nil, fmt.Errorf("something went horribly wrong")
	}
}

func formatFromFilename(filename string) string {
	if strings.HasSuffix(filename, ".yaml") {
		return yamlFormat
	} else if strings.HasSuffix(filename, ".jsonnet") {
		return jsonnetFormat
	}
	return otherFormat
}

// LocateConfFile locates the path of the configuration file.
func LocateConfFile(fs afero.Fs, userHomeDir, confPathFromEnv string) (string, error) {
	var confPath string

	jsonnetDefault := path.Join(userHomeDir, ".zit", "config.jsonnet")
	yamlDefault := path.Join(userHomeDir, ".zit", "config.yaml")

	// if ZIT_CONFIG is not set, try default location
	envVarDefined := confPathFromEnv != ""
	if envVarDefined {
		confPath = confPathFromEnv
	} else if fileExists(fs, jsonnetDefault) {
		confPath = path.Join(userHomeDir, ".zit", "config.jsonnet")
	} else if fileExists(fs, yamlDefault) {
		confPath = path.Join(userHomeDir, ".zit", "config.yaml")
	} else {
		return "", &ErrConfigNotFound{
			EnvVar: envVarDefined,
			Path:   confPath,
		}
	}

	return confPath, nil
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

func parseJsonnet(filename string) (*ConfigRoot, error) {
	vm := jsonnet.MakeVM()
	confJSON, err := vm.EvaluateFile(filename)
	if err != nil {
		return nil, err
	}

	var Hosts HostMap
	if err := json.Unmarshal([]byte(confJSON), &Hosts); err != nil {
		return nil, err
	}

	return &ConfigRoot{Hosts}, nil
}

func parseYaml(contents []byte) (*ConfigRoot, error) {
	var config ConfigRoot
	if err := yaml.Unmarshal(contents, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
