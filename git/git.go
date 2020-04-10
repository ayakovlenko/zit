package git

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	giturls "github.com/whilp/git-urls"
)

func gitCommand(args []string) (*string, error) {
	theCmd := exec.Command("git", args...)

	out, err := theCmd.CombinedOutput()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return nil, err
		}

		return nil, fmt.Errorf(
			"failed to execute %+v:\n%s",
			theCmd,
			string(out),
		)
	}

	res := strings.TrimSpace(string(out))
	return &res, nil
}

// RemoteHost TODO
func RemoteHost(name string) (*string, error) {
	out, err := gitCommand([]string{"remote", "get-url", name})
	if err != nil {
		return nil, err
	}

	remoteURL := strings.TrimSpace(string(*out))

	return &remoteURL, nil
}

// SetConfig TODO
func SetConfig(scope, key, value string) error {
	_, err := gitCommand([]string{"config", scope, key, value})
	if err != nil {
		return err
	}
	return nil
}

// GetConfig TODO
func GetConfig(scope, key string) (string, error) {
	out, err := gitCommand([]string{"config", scope, key})

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() == 1 {
				return "", nil
			}
		}

		return "", err
	}

	return *out, nil
}

// RepoInfo TODO
type RepoInfo struct {
	Host  string
	Owner string
	Name  string
}

var ownerRepoPattern = regexp.MustCompile(`\/?(.*)\/(.*)\.git$`)

// ExtractRepoInfo TODO
func ExtractRepoInfo(remote string) (*RepoInfo, error) {
	u, err := giturls.Parse(remote)
	if err != nil {
		return nil, err
	}

	match := ownerRepoPattern.FindStringSubmatch(u.Path)

	res := RepoInfo{
		u.Hostname(),
		match[1],
		match[2],
	}

	return &res, nil
}
