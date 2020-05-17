package identity

import (
	"zit/config"
	"zit/git"
)

func findBestMatch(conf config.Config, repo git.RepoInfo) (user *config.User) {
	if conf.Default != nil {
		user = &config.User{
			conf.Default.Name,
			conf.Default.Email,
		}
	}

	if conf.Overrides != nil {
		for _, override := range conf.Overrides {
			if override.Repo != "" {
				if override.Owner == repo.Owner && override.Repo == repo.Name {
					user = &config.User{
						override.User.Name,
						override.User.Email,
					}
					break
				} else {
					continue
				}
			}

			if override.Owner == repo.Owner {
				user = &config.User{
					override.User.Name,
					override.User.Email,
				}
				break
			}
		}
	}

	return
}
