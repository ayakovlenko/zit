package config

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

const (
	yamlFormat    = "yaml"
	jsonnetFormat = "jsonnet"
	otherFormat   = "other"
)

func Load(filename string) (*ConfigV2, error) {
	var format string
	if strings.HasSuffix(filename, ".yaml") {
		format = yamlFormat
	} else if strings.HasSuffix(filename, ".jsonnet") {
		format = jsonnetFormat
	} else {
		format = otherFormat
	}

	if format != yamlFormat && format != jsonnetFormat {
		return nil, fmt.Errorf("unsupported config format")
	}

	r, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	switch format {
	case yamlFormat:
		buf := new(bytes.Buffer)
		buf.ReadFrom(r)
		confStr := buf.String()
		return parseYaml(confStr)
	case jsonnetFormat:
		fmt.Println("WARN: Jsonnet configs are deprecated and going to be unsupported in future versions. Migrate to YAML format.")
		hostMap, err := readHostMap(filename, r)
		if err != nil {
			return nil, err
		}
		configV2 := toV2(hostMap)
		return configV2, nil
	default:
		return nil, fmt.Errorf("something went horribly wrong")
	}
}

func toV2(hostMap *HostMap) *ConfigV2 {
	configV2 := ConfigV2{
		Hosts: map[string]HostV2{},
	}
	for host, hostConfig := range *hostMap {
		configV2.Hosts[host] = HostV2(hostConfig)
	}
	return &configV2
}
