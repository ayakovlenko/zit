package config

import (
	"testing"
	"zit/internal/app"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocateConfig(t *testing.T) {
	t.Parallel()

	t.Run("get XDG-style YAML config from XDG_CONFIG_HOME if it exists", func(t *testing.T) {
		t.Parallel()

		var err error

		appConfig := app.NewConfig(
			newMapFS(),
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
			newMapFS(),
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
			newMapFS(),
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

	t.Run("return error when env var path does not exist", func(t *testing.T) {
		t.Parallel()

		appConfig := app.NewConfig(newMapFS(), "/home", "/missing/config.yaml", "")

		_, err := LocateConfFile(appConfig)
		require.Error(t, err)

		var notFound *ConfigNotFoundError
		require.ErrorAs(t, err, &notFound)
		assert.True(t, notFound.EnvVar)
	})

	t.Run("return error when no config file exists anywhere", func(t *testing.T) {
		t.Parallel()

		appConfig := app.NewConfig(newMapFS(), "/home", "", "")

		_, err := LocateConfFile(appConfig)
		require.Error(t, err)

		var notFound *ConfigNotFoundError
		require.ErrorAs(t, err, &notFound)
		assert.False(t, notFound.EnvVar)
	})

	t.Run("get YAML config if env var is defined", func(t *testing.T) {
		t.Parallel()

		var err error

		appConfig := app.NewConfig(
			newMapFS(),
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

		_, err := Load(app.OsFS{}, "test_data/config_00.txt")

		if err != ErrUnsupportedConfigFormat {
			t.Errorf("want: ErrUnsupportedConfigFormat; have: %+v", err)
		}
	})

	t.Run("return error when file does not exist", func(t *testing.T) {
		t.Parallel()

		_, err := Load(newMapFS(), "/missing/config.yaml")

		require.Error(t, err)
	})

	t.Run("parse YAML from injected filesystem", func(t *testing.T) {
		t.Parallel()

		mfs := newMapFS()
		mfs.writeFile("/home/.zit/config.yaml", []byte(`
hosts:
  github.example.com:
    default:
      name: "Jane Doe"
      email: "jane@example.com"
`))

		conf, err := Load(mfs, "/home/.zit/config.yaml")
		require.NoError(t, err)

		host, err := conf.Get("github.example.com")
		require.NoError(t, err)

		assert.Equal(t, "Jane Doe", host.Default.Name)
		assert.Equal(t, "jane@example.com", host.Default.Email)
	})

	t.Run("simple YAML config", func(t *testing.T) {
		t.Parallel()

		config, err := Load(app.OsFS{}, "test_data/config_01.yaml")
		require.NoError(t, err)

		host, err := config.Get("github.corp.com")
		require.NoError(t, err)

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
