package config

import (
	"testing"

	"github.com/spf13/afero"
)

func TestLocateConfig(t *testing.T) {

	fs := afero.NewMemMapFs()
	_, _ = fs.Create("/home/.zit/config.jsonnet")

	have, _ := LocateConfFile(fs, "/home", "")
	want := "/home/.zit/config.jsonnet"

	if have != want {
		t.Errorf("want: %s, have: %s", want, have)
	}
}
