package git

import "fmt"

// exitCoder matches exec.ExitError and our mock
type exitCoder interface {
	ExitCode() int
}

func IsGitDir(client GitClient) (bool, error) {
	out, err := client.Exec("rev-parse", "--is-inside-work-tree")
	if err != nil {
		// Exit code 128 means not a git repository: treat as "false"
		if exitErr, ok := err.(exitCoder); ok && exitErr.ExitCode() == 128 {
			return false, nil
		}
		return false, err
	}
	return out == "true", nil
}

func SetLocalConfig(gitClient GitClient, key, value string) error {
	if _, err := gitClient.Exec("config", "--local", key, value); err != nil {
		return fmt.Errorf("set local config %s=%s: %v", key, value, err)
	}
	return nil
}
