package git

import (
	"fmt"
	"os/exec"
	"strings"
)

type GitClient interface {
	Exec(args ...string) (string, error)
}

type realGitClient struct{}

func (r *realGitClient) Exec(args ...string) (string, error) {
	theCmd := exec.Command("git", args...)

	bout, err := theCmd.CombinedOutput()
	sout := strings.TrimSpace(string(bout))

	if err != nil {
		// exit code is not 0 but we might still care about the output, e.g.
		// `git rev-parse --is-inside-work-tree` returns 128 but we only need
		// to check if the output is string "true"
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

func NewGitClient() GitClient {
	return &realGitClient{}
}
