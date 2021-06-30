package clone

import (
	"fmt"
	"os"
	"zit/cli"
	"zit/config"
	"zit/git"
	"zit/identity"

	"github.com/spf13/cobra"
)

// CloneCmd TODO
var CloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone git repo and do all the setup",
	RunE: func(cmd *cobra.Command, args []string) error {
		url := args[0]

		repo, err := git.ExtractRepoInfo(url)
		if err != nil {
			return err
		}

		confPath, err := config.LocateConfFile()
		cli.PrintlnExit(err)

		confFile, err := os.Open(confPath)
		cli.PrintlnExit(err)

		hostMap, err := config.ReadHostMap(confPath, confFile)
		cli.PrintlnExit(err)

		conf, err := hostMap.Get((*repo).Host)
		cli.PrintlnExit(err)

		cred := identity.FindBestMatch(*conf, *repo)
		if cred == nil {
			cli.PrintlnExit(fmt.Errorf("cannot find a match for host %q", (*repo).Host))
		}

		// clone repo
		if err := git.Clone(url); err != nil {
			return err
		}

		// cd repo-name
		if err := os.Chdir(repo.Name); err != nil {
			return err
		}

		// zit set
		cli.PrintlnExit(
			git.SetConfig("--local", "user.name", cred.Name),
		)
		cli.PrintlnExit(
			git.SetConfig("--local", "user.email", cred.Email),
		)

		fmt.Printf("cloned %s\n", url)
		fmt.Printf("set user: %s <%s>\n", cred.Name, cred.Email)

		return nil
	},
}
