package git

func IsGitDir(client GitClient) (bool, error) {
	out, err := client.Exec("rev-parse", "--is-inside-work-tree")
	if err != nil {
		return false, err
	}
	return out == "true", nil
}
