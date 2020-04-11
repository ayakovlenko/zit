package git

import "testing"

func TestExtractHostNameFromRemote(t *testing.T) {

	t.Run("ssh string", func(t *testing.T) {
		url := "git@github.com:golang/go.git"

		want := RepoInfo{
			"github.com",
			"golang",
			"go",
		}
		have, _ := ExtractRepoInfo(url)

		if *have != want {
			t.Errorf("want: %s; have: %s", want, *have)
		}
	})

	t.Run("ssh string without .git suffix", func(t *testing.T) {
		url := "git@github.com:golang/go"

		want := RepoInfo{
			"github.com",
			"golang",
			"go",
		}
		have, _ := ExtractRepoInfo(url)

		if *have != want {
			t.Errorf("want: %s; have: %s", want, *have)
		}
	})

	t.Run("https string", func(t *testing.T) {
		url := "https://github.com/golang/go.git"

		want := RepoInfo{
			"github.com",
			"golang",
			"go",
		}
		have, _ := ExtractRepoInfo(url)

		if *have != want {
			t.Errorf("want: %s; have: %s", want, *have)
		}
	})
}
