package identity

import (
	"fmt"
	"os"
	"zit/internal/app"
	"zit/internal/config"
	"zit/internal/gitutil"
	"zit/pkg/git"

	"github.com/urfave/cli/v2"
)

const dryRunFlag = "dry-run"

// SetCmd is a command that sets git identity based on the configuration file.
func SetCmd(appConfig app.Config) *cli.Command {
	return &cli.Command{
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
			gitClient := git.NewGitClient()

			if err := gitutil.EnsureGitDir(gitClient); err != nil {
				return err
			}

			confPath, err := config.LocateConfFile(appConfig)
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

			return setIdentity(
				gitClient,
				*cred,
				cCtx.Bool(dryRunFlag),
			)
		},
	}
}
