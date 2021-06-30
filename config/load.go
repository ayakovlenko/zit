package config

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

func Load(filename string) (*ConfigV2, error) {
	isYaml := strings.HasSuffix(filename, ".yaml")
	isJsonnet := strings.HasSuffix(filename, ".jsonnet")

	if !isYaml && !isJsonnet {
		return nil, fmt.Errorf("unsupported config format")
	}

	r, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	if isYaml {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r)
		confStr := buf.String()
		return parseYaml(confStr)
	} else if isJsonnet {
		hostMap, err := readHostMap(filename, r)
		if err != nil {
			return nil, err
		}
		configV2 := toV2(hostMap)
		return configV2, nil
	} else {
		return nil, fmt.Errorf("something went horribly wrong")
	}
}
