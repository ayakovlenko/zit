package config

import (
	"testing"

	"github.com/spf13/afero"
)

func TestLocateConfig(t *testing.T) {
	t.Run("get XDG-style YAML config from XDG_CONFIG_HOME if it exists", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		_, _ = fs.Create("/some/custom/xdg/config/zit/config.yaml")
		// should be ignored
		_, _ = fs.Create("/home/.config/zit/config.yaml")
		// should be ignored
		_, _ = fs.Create("/home/.zit/config.yaml")

		have, _ := LocateConfFile(fs, "/home", "", "/some/custom/xdg/config")
		want := "/some/custom/xdg/config/zit/config.yaml"

		if have != want {
				t.Errorf("want: %s, have: %s", want, have)
		}
	})
		
	t.Run("get default XDG-style YAML config if it exists", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		_, _ = fs.Create("/home/.config/zit/config.yaml")
		// should be ignored
		_, _ = fs.Create("/home/.zit/config.yaml")

		have, _ := LocateConfFile(fs, "/home", "", "")
		want := "/home/.config/zit/config.yaml"

		if have != want {
				t.Errorf("want: %s, have: %s", want, have)
		}
	})

	t.Run("get default YAML config if it exists", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		_, _ = fs.Create("/home/.zit/config.yaml")

		have, _ := LocateConfFile(fs, "/home", "", "")
		want := "/home/.zit/config.yaml"

		if have != want {
			t.Errorf("want: %s, have: %s", want, have)
		}
	})

	t.Run("get YAML config if env var is defined", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		_, _ = fs.Create("/some/custom/config/path/config.yaml")
		// should be ignored
		_, _ = fs.Create("/home/.config/zit/config.yaml")
		// should be ignored
		_, _ = fs.Create("/home/.zit/config.yaml")

		have, _ := LocateConfFile(fs, "/home", "/some/custom/config/path/config.yaml", "")
        want := "/some/custom/config/path/config.yaml"

		if have != want {
			t.Errorf("want: %s, have: %s", want, have)
		}
	})
}

func TestLoad(t *testing.T) {
	t.Run("unsupported config", func(t *testing.T) {
		_, err := Load("test_data/config_00.txt")

		if err != ErrUnsupportedConfigFormat {
			t.Errorf("want: ErrUnsupportedConfigFormat; have: %+v", err)
		}
	})

	t.Run("simple YAML config", func(t *testing.T) {
		config, _ := Load("test_data/config_01.yaml")

		host, _ := config.Get("github.corp.com")

		name := host.Default.Name
		email := host.Default.Email

		if name != "John Doe" {
			t.Errorf("want: John Doe; have: %s", name)
		}

		if email != "john.doe@corp.com" {
			t.Errorf("want: john.doe@corp.com; have: %s", email)
		}
	})
}
