package config

import (
	"testing"

	"github.com/spf13/afero"
)

func TestCreateConfig(t *testing.T) {
	filename := "/Users/name/.zit/config.jsonnet"

	t.Run("create conf file", func(t *testing.T) {
		fs := afero.NewMemMapFs()

		err := createConfig(fs, filename, "test")

		if err != nil {
			t.Errorf("got error: %s", err)
		}

		bs, err := afero.ReadFile(fs, filename)
		if err != nil {
			t.Errorf("got error: %s", err)
		}

		want := "test"
		have := string(bs)
		if have != want {
			t.Errorf("want: %s; have: %s", want, have)
		}
	})

	t.Run("fail if conf file exists", func(t *testing.T) {
		fs := afero.NewMemMapFs()

		if _, err := fs.Create(filename); err != nil {
			t.Errorf("got error: %s", err)
		}

		err := createConfig(fs, filename, "test")

		if err != ErrConfFileExists {
			t.Errorf("want: %s; have: %s", ErrConfFileExists, err)
		}
	})
}
