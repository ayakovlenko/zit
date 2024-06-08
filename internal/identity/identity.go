package identity

import (
	"fmt"
	"zit/internal/config"
	"zit/internal/gitutil"
	"zit/pkg/git"
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

// setIdentity sets identity in a given repository based on a chosen identity.
func setIdentity(
	cred config.User,
	gitClient git.GitClient,
	dryRun bool,
) error {
	if !dryRun {
		if err := gitutil.SetConfig(gitClient, "--local", "user.name", cred.Name); err != nil {
			return err
		}

		if err := gitutil.SetConfig(gitClient, "--local", "user.email", cred.Email); err != nil {
			return err
		}
	}

	fmt.Printf("set user: %s <%s>\n", cred.Name, cred.Email)

	sign := cred.Signing
	if sign == nil {
		return nil
	}

	if !dryRun {
		if err := gitutil.SetConfig(gitClient, "--local", "commit.gpgsign", "true"); err != nil {
			return err
		}

		if err := gitutil.SetConfig(gitClient, "--local", "user.signingKey", sign.Key); err != nil {
			return err
		}

		if err := gitutil.SetConfig(gitClient, "--local", "gpg.format", sign.Format); err != nil {
			return err
		}
	}

	fmt.Printf("set signing key: %s key at %s\n", sign.Format, sign.Key)

	return nil
}
