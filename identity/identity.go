package identity

import (
	"zit/config"
	"zit/git"
)

func findBestMatch(conf config.Config, repo git.RepoInfo) (cred *credentials) {
	if conf.Default != nil {
		cred = &credentials{
			conf.Default.Name,
			conf.Default.Email,
		}
	}

	if conf.Overrides != nil {
		for _, override := range conf.Overrides {
			if override.Repo != nil {
				if override.Owner == repo.Owner && *override.Repo == repo.Name {
					cred = &credentials{
						override.User.Name,
						override.User.Email,
					}
					break
				} else {
					continue
				}
			}

			if override.Owner == repo.Owner {
				cred = &credentials{
					override.User.Name,
					override.User.Email,
				}
				break
			}
		}
	}

	return
}

type credentials struct {
	name  string
	email string
}
