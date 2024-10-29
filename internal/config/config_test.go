package config

import (
	"testing"
	"zit/internal/app"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocateConfig(t *testing.T) {
	t.Parallel()

	t.Run("get XDG-style YAML config from XDG_CONFIG_HOME if it exists", func(t *testing.T) {
		t.Parallel()

		var err error

		appConfig := app.NewConfig(
			afero.NewMemMapFs(),
			"/home",
			"",
			"/some/custom/xdg/config",
		)

		fs := appConfig.FS()

		_, err = fs.Create("/some/custom/xdg/config/zit/config.yaml")
		require.NoError(t, err)

		// should be ignored
		_, err = fs.Create("/home/.config/zit/config.yaml")
		require.NoError(t, err)

		// should be ignored
		_, err = fs.Create("/home/.zit/config.yaml")
		require.NoError(t, err)

		have, err := LocateConfFile(appConfig)
		require.NoError(t, err)

		want := "/some/custom/xdg/config/zit/config.yaml"

		assert.Equal(t, want, have)
	})

	t.Run("get default XDG-style YAML config if it exists", func(t *testing.T) {
		t.Parallel()

		var err error

		appConfig := app.NewConfig(
			afero.NewMemMapFs(),
			"/home",
			"",
			"",
		)

		fs := appConfig.FS()

		_, err = fs.Create("/home/.config/zit/config.yaml")
		require.NoError(t, err)

		// should be ignored
		_, err = fs.Create("/home/.zit/config.yaml")
		require.NoError(t, err)

		have, err := LocateConfFile(appConfig)
		require.NoError(t, err)

		want := "/home/.config/zit/config.yaml"

		assert.Equal(t, want, have)
	})

	t.Run("get default YAML config if it exists", func(t *testing.T) {
		t.Parallel()

		var err error

		appConfig := app.NewConfig(
			afero.NewMemMapFs(),
			"/home",
			"",
			"",
		)

		fs := appConfig.FS()

		_, err = fs.Create("/home/.zit/config.yaml")
		require.NoError(t, err)

		have, err := LocateConfFile(appConfig)
		require.NoError(t, err)

		want := "/home/.zit/config.yaml"

		assert.Equal(t, want, have)
	})

	t.Run("get YAML config if env var is defined", func(t *testing.T) {
		t.Parallel()

		var err error

		appConfig := app.NewConfig(
			afero.NewMemMapFs(),
			"/home",
			"/some/custom/config/path/config.yaml",
			"",
		)

		fs := appConfig.FS()

		_, err = fs.Create("/some/custom/config/path/config.yaml")
		require.NoError(t, err)

		// should be ignored
		_, err = fs.Create("/home/.config/zit/config.yaml")
		require.NoError(t, err)

		// should be ignored
		_, err = fs.Create("/home/.zit/config.yaml")
		require.NoError(t, err)

		have, err := LocateConfFile(appConfig)
		require.NoError(t, err)

		want := "/some/custom/config/path/config.yaml"

		assert.Equal(t, want, have)
	})
}

func TestLoad(t *testing.T) {
	t.Parallel()

	t.Run("unsupported config", func(t *testing.T) {
		t.Parallel()

		_, err := Load("test_data/config_00.txt")

		if err != ErrUnsupportedConfigFormat {
			t.Errorf("want: ErrUnsupportedConfigFormat; have: %+v", err)
		}
	})

	t.Run("simple YAML config", func(t *testing.T) {
		t.Parallel()

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
