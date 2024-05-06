package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

func Load(fs afero.Fs, filename string) (*ConfigRoot, error) {
	format := formatFromFilename(filename)

	if format != yamlFormat && format != jsonnetFormat {
		return nil, ErrUnsupportedConfigFormat
	}

	switch format {
	case yamlFormat:
		r, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer r.Close()
		buf := new(bytes.Buffer)
		buf.ReadFrom(r)
		confStr := buf.String()
		return parseYaml(confStr)
	case jsonnetFormat:
		fmt.Println("WARN: Jsonnet configs are deprecated and going to be unsupported in future versions. Migrate to YAML format.")
		r, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer r.Close()
		return parseJsonnet(filename, r)
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

	// if ZIT_CONFIG is not set, try default location
	envVarDefined := confPathFromEnv != ""
	if envVarDefined {
		confPath = confPathFromEnv
	} else {
		confPath = path.Join(userHomeDir, ".zit", "config.jsonnet")
	}

	if !fileExists(fs, confPath) {
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

func parseJsonnet(filename string, r io.Reader) (*ConfigRoot, error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	confStr := buf.String()

	vm := jsonnet.MakeVM()
	confJSON, err := vm.EvaluateAnonymousSnippet(filename, confStr)
	if err != nil {
		return nil, err
	}

	var Hosts HostMap
	if err := json.Unmarshal([]byte(confJSON), &Hosts); err != nil {
		return nil, err
	}

	return &ConfigRoot{
		Hosts,
	}, nil
}

func parseYaml(s string) (*ConfigRoot, error) {
	var config ConfigRoot
	if err := yaml.Unmarshal([]byte(s), &config); err != nil {
		return nil, err
	}
	return &config, nil
}
