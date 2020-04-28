package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/google/go-jsonnet"
	"github.com/spf13/cobra"
)

var ErrConfFileExists = errors.New("conf file already exists")

// HostMap TODO
type HostMap map[string]Config

// Get TODO
func (hm *HostMap) Get(host string) (*Config, error) {
	conf, ok := (*hm)[host]
	if !ok {
		return nil, fmt.Errorf("cannot find config for host %q", host)
	}

	return &conf, nil
}

// Config TODO
type Config struct {
	Default   *User      `json:"default"`
	Overrides []Override `json:"overrides"`
}

// User TODO
type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Override TODO
type Override struct {
	Owner string  `json:"owner"`
	Repo  *string `json:"repo,omitempty"`
	User  User    `json:"user"`
}

func readConfJSON(filename string, r io.Reader) (string, error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	confStr := buf.String()

	vm := jsonnet.MakeVM()
	confJSON, err := vm.EvaluateSnippet(filename, confStr)
	if err != nil {
		return "", err
	}

	return confJSON, nil
}

// ReadHostMap TODO
func ReadHostMap(filename string, r io.Reader) (*HostMap, error) {
	confJSON, err := readConfJSON(filename, r)

	if err != nil {
		return nil, err
	}

	var hostMap HostMap
	if err := json.Unmarshal([]byte(confJSON), &hostMap); err != nil {
		return nil, err
	}

	return &hostMap, nil
}

// LocateConfFile locates the path of the configuration file.
func LocateConfFile() (string, error) {
	fileExists := func(filename string) bool {
		info, err := os.Stat(filename)
		if os.IsNotExist(err) {
			return false
		}
		return !info.IsDir()
	}

	confPath, err := defaultConfFileLocation()
	if err != nil {
		return "", err
	}

	if !fileExists(confPath) {
		return "", fmt.Errorf("config file not found at %s", confPath)
	}

	return confPath, nil
}

func defaultConfFileLocation() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return path.Join(home, ".zit", "config.jsonnet"), nil
}

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "???", // TODO
}

func init() {
	ConfigCmd.AddCommand(
		configPathCmd,
		configShowJsonCmd,
		configInitCmd,
	)
}
