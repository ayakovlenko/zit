package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractHostNameFromRemote(t *testing.T) {
	assertEquals := func(t *testing.T, want, have *RepoInfo) {
		t.Helper()

		if *have != *want {
			t.Errorf("want: %s; have: %s", want, have)
		}
	}

	t.Run("ssh string", func(t *testing.T) {
		url := "git@github.com:golang/go.git"

		want := &RepoInfo{
			"github.com",
			"golang",
			"go",
		}
		have, _ := ExtractRepoInfo(url)

		assertEquals(t, have, want)
	})

	t.Run("ssh string without .git suffix", func(t *testing.T) {
		url := "git@github.com:golang/go"

		want := &RepoInfo{
			"github.com",
			"golang",
			"go",
		}
		have, _ := ExtractRepoInfo(url)

		assertEquals(t, have, want)
	})

	t.Run("https string", func(t *testing.T) {
		url := "https://github.com/golang/go.git"

		want := &RepoInfo{
			"github.com",
			"golang",
			"go",
		}
		have, _ := ExtractRepoInfo(url)

		assertEquals(t, have, want)
	})

	t.Run("dots in the name", func(t *testing.T) {
		t.Run("https string", func(t *testing.T) {
			url := "https://github.com/hackerloft/hackerloft.github.io.git"

			want := &RepoInfo{
				"github.com",
				"hackerloft",
				"hackerloft.github.io",
			}
			have, _ := ExtractRepoInfo(url)

			assertEquals(t, have, want)
		})

		t.Run("ssh string", func(t *testing.T) {
			url := "git@github.com:hackerloft/hackerloft.github.io.git"

			want := &RepoInfo{
				"github.com",
				"hackerloft",
				"hackerloft.github.io",
			}
			have, _ := ExtractRepoInfo(url)

			assertEquals(t, have, want)
		})
	})

	t.Run("gist.github.com", func(t *testing.T) {
		t.Run("https string", func(t *testing.T) {
			url := "https://gist.github.com/2b50de6071556f306e1952b227a47292.git"

			want := &RepoInfo{
				"gist.github.com",
				"",
				"2b50de6071556f306e1952b227a47292",
			}
			have, _ := ExtractRepoInfo(url)

			assertEquals(t, have, want)
		})

		t.Run("ssh string", func(t *testing.T) {
			url := "git@gist.github.com:2b50de6071556f306e1952b227a47292.git"

			want := &RepoInfo{
				"gist.github.com",
				"",
				"2b50de6071556f306e1952b227a47292",
			}
			have, _ := ExtractRepoInfo(url)

			assertEquals(t, have, want)
		})
	})
}

func TestIsGitDir(t *testing.T) {
	t.Run("is git directory", func(t *testing.T) {
		gitClient := NewMockGitClient()

		gitClient.AddCommand(
			[]string{"rev-parse", "--is-inside-work-tree"},
			"true",
			nil,
		)

		ok, _ := IsGitDir(gitClient)

		assert.True(t, ok)
	})

	t.Run("is not git directory", func(t *testing.T) {
		gitClient := NewMockGitClient()

		gitClient.AddCommand(
			[]string{"rev-parse", "--is-inside-work-tree"},
			"fatal: not a git repository (or any of the parent directories): .git",
			nil,
		)

		ok, _ := IsGitDir(gitClient)

		assert.False(t, ok)
	})
}
