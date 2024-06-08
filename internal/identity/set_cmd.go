package identity

import (
	"fmt"
	"os"
	"zit/internal/config"
	"zit/internal/gitutil"
	"zit/pkg/git"

	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

const dryRunFlag = "dry-run"

// SetCmd is a command that sets git identity based on the configuration file.
var SetCmd = &cli.Command{
	Name:  "set",
	Usage: "Set git identity",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  dryRunFlag,
			Value: false,
			Usage: "run without applying configurations",
		},
	},
	Action: func(cCtx *cli.Context) error {
		fs := afero.NewOsFs()

		dryRun := cCtx.Bool(dryRunFlag)

		gitClient := git.NewGitClient()

		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("cannot user home dir: %v", err)
		}

		configPathFromEnv := os.Getenv(config.EnvVarName)

		return setIdentityAction(fs, gitClient, userHomeDir, dryRun, configPathFromEnv)
	},
}

func setIdentityAction(
	fs afero.Fs,
	gitClient git.GitClient,
	userHomeDir string,
	dryRun bool,
	configPathFromEnv string,
) error {
	if err := gitutil.EnsureGitDir(gitClient); err != nil {
		return err
	}

	confPath, err := config.LocateConfFile(
		fs,
		userHomeDir,
		configPathFromEnv,
	)
	if err != nil {
		return err
	}

	conf, err := config.Load(confPath)
	if err != nil {
		return err
	}

	host, err := gitutil.RemoteURL(gitClient, "origin")
	if err != nil {
		if _, ok := err.(*gitutil.ErrNoRemoteURL); ok {
			fmt.Printf(`Error: %s

Add remote URL so that zit could use it for choosing the correct git identity as
defined in the configuration file:

git remote add origin <url>
`, err)
			os.Exit(1) // TODO: return "FriendlyError" instead of os.Exit
		} else {
			return err
		}
	}

	repo, err := gitutil.ExtractRepoInfo(host)
	if err != nil {
		return err
	}

	hostConf, err := conf.Get(repo.Host)
	if err != nil {
		return err
	}

	cred := findBestMatch(*hostConf, *repo)
	if cred == nil {
		return fmt.Errorf("cannot find a match for host %q", repo.Host)
	}

	if !dryRun {
		if err := gitutil.SetConfig(gitClient, "--local", "user.name", cred.Name); err != nil {
			return err
		}

		if err := gitutil.SetConfig(gitClient, "--local", "user.email", cred.Email); err != nil {
			return err
		}

		if sign := cred.Signing; sign != nil {
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
	}

	if dryRun {
		fmt.Printf("[dry-run]\n")
	}

	fmt.Printf("set user: %s <%s>\n", cred.Name, cred.Email)
	if sign := cred.Signing; sign != nil {
		fmt.Printf("set signing key: %s key at %s\n", sign.Format, sign.Key)
	}

	return nil
}
