package identity

import (
	"zit/internal/config"
	"zit/internal/git"
)

func findBestMatch(conf config.HostV2, repo git.RepoInfo) (user *config.User) {
	if conf.Default != nil {
		user = &config.User{
			Name:  conf.Default.Name,
			Email: conf.Default.Email,
		}
	}

	if conf.Overrides != nil {
		for _, override := range conf.Overrides {
			if override.Repo != "" {
				if override.Owner == repo.Owner && override.Repo == repo.Name {
					user = &config.User{
						Name:  override.User.Name,
						Email: override.User.Email,
					}
					break
				} else {
					continue
				}
			}

			if override.Owner == repo.Owner {
				user = &config.User{
					Name:  override.User.Name,
					Email: override.User.Email,
				}
				break
			}
		}
	}

	return
}
