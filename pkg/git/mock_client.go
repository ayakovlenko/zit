package git

import (
	"fmt"
	"strings"
)

type mockGitClient struct {
	commands map[string][2]interface{}
}

func NewMockGitClient() *mockGitClient {
	return &mockGitClient{
		commands: make(map[string][2]interface{}),
	}
}

func (m *mockGitClient) Exec(args ...string) (string, error) {
	cmd := strings.Join(args, " ")

	if tup, ok := m.commands[cmd]; ok {
		var ret string = tup[0].(string)
		var err error

		if isNil := tup[1] == nil; !isNil {
			err = tup[1].(error)
		}

		return ret, err
	}

	return "", fmt.Errorf("command %q not found in mock", cmd)
}

func (m *mockGitClient) AddCommand(args []string, ret string, err error) {
	argsKey := strings.Join(args, " ")
	m.commands[argsKey] = [2]interface{}{
		ret,
		err,
	}
}
