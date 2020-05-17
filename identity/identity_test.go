package identity

import (
	"reflect"
	"testing"
	"zit/config"
	"zit/git"
)

func TestFindBestMatch(t *testing.T) {

	conf := config.Config{
		Default: nil,
		Overrides: []config.Override{
			{
				Owner: "",
				Repo:  nil,
				User: config.User{
					Name:  "john doe",
					Email: "john.doe@gmail.com",
				},
			},
		},
	}

	repoInfo := git.RepoInfo{
		Host:  "gist.github.com",
		Owner: "",
		Name:  "29274a722b2591e603f6551706ed05b2",
	}

	want := &credentials{
		name:  "john doe",
		email: "john.doe@gmail.com",
	}

	have := findBestMatch(conf, repoInfo)

	if !reflect.DeepEqual(want, have) {
		t.Errorf("want: %s; have: %s", want, have)
	}
}
