package identity

import (
	"zit/internal/config"
	"zit/internal/gitutil"
)

func findBestMatch(conf config.HostConfig, repo gitutil.RepoInfo) *config.User {
	var user *config.User

	if conf.Default != nil {
		user = &config.User{
			Name:  conf.Default.Name,
			Email: conf.Default.Email,
		}
	}

	if conf.Overrides == nil {
		return user
	}

	for _, override := range conf.Overrides {
		if override.Repo != "" {
			if override.Owner == repo.Owner && override.Repo == repo.Name {
				return &config.User{
					Name:  override.User.Name,
					Email: override.User.Email,
				}
			}
		}

		if override.Owner == repo.Owner {
			return &config.User{
				Name:  override.User.Name,
				Email: override.User.Email,
			}
		}
	}

	return user
}
