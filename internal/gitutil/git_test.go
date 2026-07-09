package gitutil

import (
	"errors"
	"testing"

	"zit/pkg/git"
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

		assertEquals(t, want, have)
	})

	t.Run("ssh string without .git suffix", func(t *testing.T) {
		url := "git@github.com:golang/go"

		want := &RepoInfo{
			"github.com",
			"golang",
			"go",
		}
		have, _ := ExtractRepoInfo(url)

		assertEquals(t, want, have)
	})

	t.Run("https string", func(t *testing.T) {
		url := "https://github.com/golang/go.git"

		want := &RepoInfo{
			"github.com",
			"golang",
			"go",
		}
		have, _ := ExtractRepoInfo(url)

		assertEquals(t, want, have)
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

			assertEquals(t, want, have)
		})

		t.Run("ssh string", func(t *testing.T) {
			url := "git@github.com:hackerloft/hackerloft.github.io.git"

			want := &RepoInfo{
				"github.com",
				"hackerloft",
				"hackerloft.github.io",
			}
			have, _ := ExtractRepoInfo(url)

			assertEquals(t, want, have)
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

			assertEquals(t, want, have)
		})

		t.Run("ssh string", func(t *testing.T) {
			url := "git@gist.github.com:2b50de6071556f306e1952b227a47292.git"

			want := &RepoInfo{
				"gist.github.com",
				"",
				"2b50de6071556f306e1952b227a47292",
			}
			have, _ := ExtractRepoInfo(url)

			assertEquals(t, want, have)
		})
	})
}

func TestRemoteURL(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		gitClient := git.NewMockGitClient()
		gitClient.AddCommand(
			[]string{"remote", "get-url", "origin"},
			"git@github.com:user/repo.git",
			nil,
		)

		url, err := RemoteURL(gitClient, "origin")

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if url != "git@github.com:user/repo.git" {
			t.Errorf("want: git@github.com:user/repo.git; have: %s", url)
		}
	})

	t.Run("exit code 2: remote not set", func(t *testing.T) {
		gitClient := git.NewMockGitClient()
		gitClient.AddExitError(
			[]string{"remote", "get-url", "origin"},
			"error: No such remote 'origin'",
			2,
		)

		_, err := RemoteURL(gitClient, "origin")

		if err == nil {
			t.Error("expected error, got nil")
		}
		if !errors.Is(err, &ErrNoRemoteURL{"origin"}) && err.Error() != `remote "origin" is not set` {
			t.Errorf("expected ErrNoRemoteURL, got: %v", err)
		}
	})

	t.Run("exit code 128: return underlying error", func(t *testing.T) {
		gitClient := git.NewMockGitClient()
		gitClient.AddExitError(
			[]string{"remote", "get-url", "origin"},
			"fatal: not a git repository",
			128,
		)

		_, err := RemoteURL(gitClient, "origin")

		if err == nil {
			t.Error("expected error, got nil")
		}
		if err.Error() != "fatal: not a git repository" {
			t.Errorf("expected git error message, got: %v", err)
		}
	})
}
