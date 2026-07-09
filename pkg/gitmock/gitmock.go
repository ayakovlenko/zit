package gitmock

import (
	"fmt"
	"strings"
	"zit/pkg/git"
)

type MockClient struct {
	commands map[string]mockResult
}

type mockResult struct {
	output string
	err    error
}

func NewMockGitClient() *MockClient {
	return &MockClient{
		commands: make(map[string]mockResult),
	}
}

func (m *MockClient) Exec(args ...string) (string, error) {
	cmd := strings.Join(args, " ")

	if res, ok := m.commands[cmd]; ok {
		return res.output, res.err
	}

	return "", fmt.Errorf("command %q not found in mock", cmd)
}

func (m *MockClient) AddCommand(args []string, output string, err error) {
	argsKey := strings.Join(args, " ")
	m.commands[argsKey] = mockResult{
		output: output,
		err:    err,
	}
}

func (m *MockClient) AddExitError(args []string, output string, exitCode int) {
	argsKey := strings.Join(args, " ")
	m.commands[argsKey] = mockResult{
		output: output,
		err:    &mockExitError{exitCode: exitCode, stderr: output},
	}
}

// Verify MockClient implements git.GitClient at compile time.
var _ git.GitClient = (*MockClient)(nil)

type mockExitError struct {
	exitCode int
	stderr   string
}

func (e *mockExitError) Error() string {
	return fmt.Sprintf("exit status %d", e.exitCode)
}

func (e *mockExitError) ExitCode() int {
	return e.exitCode
}
