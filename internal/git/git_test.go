package git

import "testing"

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
