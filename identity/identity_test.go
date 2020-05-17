package identity

import (
	"reflect"
	"testing"
	"zit/config"
	"zit/git"
)

func TestFindBestMatch(t *testing.T) {

	t.Run("match default user", func(t *testing.T) {
		conf := config.Config{
			Default: &config.User{
				Name:  "john doe",
				Email: "john.doe@gmail.com",
			},
			Overrides: []config.Override{
				{
					Owner: "corporation",
					Repo:  "",
					User: config.User{
						Name:  "john doe",
						Email: "john.doe@corporate.com",
					},
				},
			},
		}

		repoInfo := git.RepoInfo{
			Host:  "github.com",
			Owner: "johndoe",
			Name:  "repo",
		}

		want := &credentials{
			name:  "john doe",
			email: "john.doe@gmail.com",
		}

		have := findBestMatch(conf, repoInfo)

		if !reflect.DeepEqual(want, have) {
			t.Errorf("want: %s; have: %s", want, have)
		}
	})

	t.Run("match owner override", func(t *testing.T) {
		conf := config.Config{
			Default: &config.User{
				Name:  "john doe",
				Email: "john.doe@gmail.com",
			},
			Overrides: []config.Override{
				{
					Owner: "corporation",
					Repo:  "",
					User: config.User{
						Name:  "john doe",
						Email: "john.doe@corporation.com",
					},
				},
			},
		}

		repoInfo := git.RepoInfo{
			Host:  "github.com",
			Owner: "corporation",
			Name:  "repo",
		}

		want := &credentials{
			name:  "john doe",
			email: "john.doe@corporation.com",
		}

		have := findBestMatch(conf, repoInfo)

		if !reflect.DeepEqual(want, have) {
			t.Errorf("want: %s; have: %s", want, have)
		}
	})

	t.Run("match repo override", func(t *testing.T) {
		conf := config.Config{
			Default: &config.User{
				Name:  "john doe",
				Email: "john.doe@gmail.com",
			},
			Overrides: []config.Override{
				{
					Owner: "",
					Repo:  "gist",
					User: config.User{
						Name:  "john doe",
						Email: "john.doe@corporation.com",
					},
				},
			},
		}

		repoInfo := git.RepoInfo{
			Host:  "gist.github.com",
			Owner: "",
			Name:  "gist",
		}

		want := &credentials{
			name:  "john doe",
			email: "john.doe@corporation.com",
		}

		have := findBestMatch(conf, repoInfo)

		if !reflect.DeepEqual(want, have) {
			t.Errorf("want: %s; have: %s", want, have)
		}
	})

	t.Run("match repo and owner override", func(t *testing.T) {
		conf := config.Config{
			Default: &config.User{
				Name:  "john doe",
				Email: "john.doe@gmail.com",
			},
			Overrides: []config.Override{
				{
					Owner: "1",
					Repo:  "1",
					User: config.User{
						Name:  "john doe 1",
						Email: "john.doe@corporation.com",
					},
				},
				{
					Owner: "2",
					Repo:  "2",
					User: config.User{
						Name:  "john doe 2",
						Email: "john.doe@corporation.com",
					},
				},
				{
					Owner: "3",
					Repo:  "3",
					User: config.User{
						Name:  "john doe 3",
						Email: "john.doe@corporation.com",
					},
				},
			},
		}

		repoInfo := git.RepoInfo{
			Host:  "github.com",
			Owner: "2",
			Name:  "2",
		}

		want := &credentials{
			name:  "john doe 2",
			email: "john.doe@corporation.com",
		}

		have := findBestMatch(conf, repoInfo)

		if !reflect.DeepEqual(want, have) {
			t.Errorf("want: %s; have: %s", want, have)
		}
	})
}
