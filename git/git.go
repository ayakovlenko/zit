package git

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	giturls "github.com/whilp/git-urls"
)

func gitCommand(args []string) (string, error) {
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

// RemoteHost TODO
func RemoteHost(name string) (string, error) {
	out, err := gitCommand([]string{"remote", "get-url", name})
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() == 128 {
				return "", fmt.Errorf("remote %q not set", name)
			}
		}
		return out, err
	}

	return out, nil
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

var ownerRepoPattern = regexp.MustCompile(`\/?(.*)\/([^.]*)(\.git)?$`)

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
