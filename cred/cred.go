package cred

import (
	"fmt"
	"os"
	"zit/cli"
	"zit/config"
	"zit/git"

	"github.com/spf13/cobra"
)

// SetCredCmd TODO
var SetCredCmd = &cobra.Command{
	Use:   "zit",
	Short: "git identity manager",
	Run: func(cmd *cobra.Command, args []string) {

		confPath, err := config.LocateConfFile()
		cli.PrintlnExit(err)

		confFile, err := os.Open(confPath)
		cli.PrintlnExit(err)

		hostMap, err := config.ReadHostMap(confPath, confFile)
		cli.PrintlnExit(err)

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
				cli.PrintlnExit(err)
			}
		}

		repo, err := git.ExtractRepoInfo(host)
		cli.PrintlnExit(err)

		conf, err := hostMap.Get((*repo).Host)
		cli.PrintlnExit(err)

		cred := findBestMatch(*conf, *repo)
		if cred == nil {
			cli.PrintlnExit(fmt.Errorf("cannot find a match for host %q", (*repo).Host))
		}

		cli.PrintlnExit(
			git.SetConfig("--local", "user.name", cred.name),
		)
		cli.PrintlnExit(
			git.SetConfig("--local", "user.email", cred.email),
		)

		fmt.Printf("set user: %s <%s>\n", cred.name, cred.email)
	},
}

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
