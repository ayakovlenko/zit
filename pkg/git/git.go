package git

import "fmt"

func IsGitDir(client GitClient) (bool, error) {
	out, err := client.Exec("rev-parse", "--is-inside-work-tree")
	if err != nil {
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
