package identity

import (
	"zit/internal/config"
	"zit/internal/gitutil"
)

func findBestMatch(conf config.HostConfig, repo gitutil.RepoInfo) *config.User {
	var user *config.User

	if conf.Default != nil {
		user = conf.Default
	}

	if conf.Overrides == nil {
		return user
	}

	for _, override := range conf.Overrides {
		if override.Repo != "" {
			if override.Owner == repo.Owner && override.Repo == repo.Name {
				return &override.User
			}
		}

		if override.Owner == repo.Owner {
			return &override.User
		}
	}

	return user
}
