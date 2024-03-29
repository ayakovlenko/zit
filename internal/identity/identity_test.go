package identity

import (
	"testing"
	"zit/internal/config"
	"zit/internal/git"

	"github.com/stretchr/testify/assert"
)

func TestFindBestMatch(t *testing.T) {

	defaultUser := config.User{
		Name:  "john doe",
		Email: "john.doe@gmail.com",
	}

	otherUser := config.User{
		Name:  "ion popescu",
		Email: "ion.popescu@gmail.com",
	}

	corpUser1 := config.User{
		Name:  "john doe",
		Email: "john.doe@corporate.com",
	}

	corpUser2 := config.User{
		Name:  "john doe",
		Email: "john.doe@corporate2.com",
	}

	corpUser3 := config.User{
		Name:  "john doe",
		Email: "john.doe@corporate3.com",
	}

	t.Run("match default user", func(t *testing.T) {
		conf := config.HostV2{
			Default: &defaultUser,
			Overrides: []config.Override{
				{
					Owner: "corporation",
					Repo:  "",
					User:  corpUser1,
				},
			},
		}

		repoInfo := git.RepoInfo{
			Host:  "github.com",
			Owner: "johndoe",
			Name:  "repo",
		}

		want := &defaultUser
		have := findBestMatch(conf, repoInfo)

		assert.Equal(t, want, have)
	})

	t.Run("match owner override", func(t *testing.T) {
		conf := config.HostV2{
			Default: &defaultUser,
			Overrides: []config.Override{
				{
					Owner: "corporation",
					Repo:  "",
					User:  corpUser1,
				},
			},
		}

		repoInfo := git.RepoInfo{
			Host:  "github.com",
			Owner: "corporation",
			Name:  "repo",
		}

		want := &corpUser1
		have := findBestMatch(conf, repoInfo)

		assert.Equal(t, want, have)
	})

	t.Run("match name override", func(t *testing.T) {
		repoName := "override"
		conf := config.HostV2{
			Default: &defaultUser,
			Overrides: []config.Override{
				{
					Owner: defaultUser.Name,
					Repo:  repoName,
					User:  otherUser,
				},
			},
		}

		repoInfo := git.RepoInfo{
			Host:  "github.com",
			Owner: defaultUser.Name,
			Name:  repoName,
		}

		want := &otherUser
		have := findBestMatch(conf, repoInfo)

		assert.Equal(t, want, have)
	})

	t.Run("match repo override", func(t *testing.T) {
		conf := config.HostV2{
			Default: &defaultUser,
			Overrides: []config.Override{
				{
					Owner: "",
					Repo:  "gist",
					User:  corpUser1,
				},
			},
		}

		repoInfo := git.RepoInfo{
			Host:  "gist.github.com",
			Owner: "",
			Name:  "gist",
		}

		want := &corpUser1
		have := findBestMatch(conf, repoInfo)

		assert.Equal(t, want, have)
	})

	t.Run("match repo and owner override", func(t *testing.T) {
		conf := config.HostV2{
			Default: &defaultUser,
			Overrides: []config.Override{
				{
					Owner: "1",
					Repo:  "1",
					User:  corpUser1,
				},
				{
					Owner: "2",
					Repo:  "2",
					User:  corpUser2,
				},
				{
					Owner: "3",
					Repo:  "3",
					User:  corpUser3,
				},
			},
		}

		repoInfo := git.RepoInfo{
			Host:  "github.com",
			Owner: "2",
			Name:  "2",
		}

		want := &corpUser2
		have := findBestMatch(conf, repoInfo)

		assert.Equal(t, want, have)
	})
}
