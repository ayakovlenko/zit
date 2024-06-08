package gitutil

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"zit/pkg/git"

	giturls "github.com/mojotx/git-urls"
)

// ErrNoRemoteURL defines an error returned when the remote URL is not set.
type ErrNoRemoteURL struct {
	name string
}

func (e *ErrNoRemoteURL) Error() string {
	return fmt.Sprintf("remote %q is not set", e.name)
}

// RemoteURL gets git remote URL by remote name.
func RemoteURL(gitClient git.GitClient, name string) (string, error) {
	out, err := gitClient.Exec("remote", "get-url", name)
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() == 128 {
				return "", &ErrNoRemoteURL{name}
			}
		}
		return out, err
	}

	return out, nil
}

// GetConfig TODO
func GetConfig(gitClient git.GitClient, scope, key string) (string, error) {
	out, err := gitClient.Exec("config", scope, key)

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() == 1 {
				return "", nil
			}
		}

		return out, err
	}

	return out, nil
}

// RepoInfo TODO
type RepoInfo struct {
	Host  string
	Owner string
	Name  string
}

var (
	ownerRepoPattern = regexp.MustCompile(`\/?(.*?)\/(.*?)(\.git)?$`)
	repoOnlyPattern  = regexp.MustCompile(`\/?(.*?)(\.git)?$`)
)

// ExtractRepoInfo extracts repository information, such as the repository owner
// (username or organization name), the repository name, and the git host of the
// repository) from remote URL.
func ExtractRepoInfo(remoteURL string) (*RepoInfo, error) {
	u, err := giturls.Parse(remoteURL)
	if err != nil {
		return nil, err
	}

	var owner string
	var repo string

	match := ownerRepoPattern.FindStringSubmatch(u.Path)
	if match == nil {
		match = repoOnlyPattern.FindStringSubmatch(u.Path)

		if match != nil {
			owner = ""
			repo = match[1]
		} else {
			return nil, fmt.Errorf("remote url doesn't match any pattern: %s", remoteURL)
		}
	} else {
		owner = match[1]
		repo = match[2]
	}

	res := RepoInfo{
		u.Hostname(),
		owner,
		repo,
	}

	return &res, nil
}

func EnsureGitDir(gitClient git.GitClient) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	ok, err := git.IsGitDir(gitClient)
	if err != nil {
		return err
	}

	if !ok {
		fmt.Fprintf(os.Stderr, `Error: %q is not a git directory

Make sure you are executing zit inside a git directory.

If you are, perhaps you have forgotten to initialize a new repository? In this
case, run:

    git init
`, dir)
		os.Exit(1)
	}

	return nil
}
