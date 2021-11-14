package git

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	giturls "github.com/whilp/git-urls"
)

// ErrNoRemoteURL defines an error returned when the remote URL is not set.
type ErrNoRemoteURL struct {
	name string
}

func (e *ErrNoRemoteURL) Error() string {
	return fmt.Sprintf("remote %q is not set", e.name)
}

func git(args ...string) (string, error) {
	theCmd := exec.Command("git", args...)

	bout, err := theCmd.CombinedOutput()
	sout := strings.TrimSpace(string(bout))

	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return sout, err
		}

		return sout, fmt.Errorf(
			"failed to execute %+v:\n%s",
			theCmd,
			sout,
		)
	}

	return sout, nil
}

// RemoteURL gets git remote URL by remote name.
func RemoteURL(name string) (string, error) {
	out, err := git("remote", "get-url", name)
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

// SetConfig TODO
func SetConfig(scope, key, value string) error {
	_, err := git("config", scope, key, value)
	if err != nil {
		return err
	}
	return nil
}

// GetConfig TODO
func GetConfig(scope, key string) (string, error) {
	out, err := git("config", scope, key)

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

var ownerRepoPattern = regexp.MustCompile(`\/?(.*?)\/(.*?)(\.git)?$`)
var repoOnlyPattern = regexp.MustCompile(`\/?(.*?)(\.git)?$`)

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

// IsGitDir checks if dir is a git directory
func IsGitDir(dir string) (bool, error) {
	if _, err := git("status"); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() == 128 {
				return false, nil
			}
		}

		return false, err
	}

	return true, nil
}
