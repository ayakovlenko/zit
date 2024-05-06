package config

import (
	"testing"

	"github.com/spf13/afero"
)

func TestLocateConfig(t *testing.T) {

	t.Run("get Jsonnet config if it exists", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		_, _ = fs.Create("/home/.zit/config.jsonnet")

		have, _ := LocateConfFile(fs, "/home", "")
		want := "/home/.zit/config.jsonnet"

		if have != want {
			t.Errorf("want: %s, have: %s", want, have)
		}
	})

	t.Run("get YAML config if env var is defined", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		_, _ = fs.Create("/home/.zit/config.yaml")

		have, _ := LocateConfFile(fs, "/home", "/home/.zit/config.yaml")
		want := "/home/.zit/config.yaml"

		if have != want {
			t.Errorf("want: %s, have: %s", want, have)
		}
	})
}
