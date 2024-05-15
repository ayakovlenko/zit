package identity

import (
	"fmt"
	"os"
	"zit/internal/config"
	"zit/internal/git"

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

		if err := git.EnsureGitDir(); err != nil {
			return err
		}

		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		confPath, err := config.LocateConfFile(
			fs,
			userHomeDir,
			os.Getenv(config.EnvVarName),
		)
		if err != nil {
			return err
		}

		conf, err := config.Load(confPath)
		if err != nil {
			return err
		}

		host, err := git.RemoteURL("origin")
		if err != nil {
			if _, ok := err.(*git.ErrNoRemoteURL); ok {
				fmt.Printf(`Error: %s

Add remote URL so that zit could use it for choosing the correct git identity as
defined in the configuration file:

    git remote add origin <url>
`, err)
				os.Exit(1)
			} else {
				return err
			}
		}

		repo, err := git.ExtractRepoInfo(host)
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

		dryRun := cCtx.Bool(dryRunFlag)

		if !dryRun {
			if err := git.SetConfig("--local", "user.name", cred.Name); err != nil {
				return err
			}
			if err := git.SetConfig("--local", "user.email", cred.Email); err != nil {
				return err
			}
		}

		if dryRun {
			fmt.Printf("[dry-run] ")
		}

		fmt.Printf("set user: %s <%s>\n", cred.Name, cred.Email)

		return nil
	},
}
