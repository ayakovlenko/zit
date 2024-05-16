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

// ---

type gitTuple struct {
	s string
	e error
}

type mockGitClient struct {
	commands map[string]gitTuple
}

func NewMockGitClient() *mockGitClient {
	return &mockGitClient{
		commands: make(map[string]gitTuple),
	}
}

func (m *mockGitClient) Exec(args ...string) (string, error) {
	cmd := strings.Join(args, " ")

	if t, ok := m.commands[cmd]; ok {
		return t.s, t.e
	}

	return "", fmt.Errorf("command %q not found in mock", cmd)
}

func (m *mockGitClient) AddCommand(args []string, s string, e error) {
	m.commands[strings.Join(args, " ")] = gitTuple{s, e}
}
