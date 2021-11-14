package identity

import (
	"fmt"
	"os"
	"zit/internal/cli"
	"zit/internal/config"
	"zit/internal/git"

	"github.com/spf13/cobra"
)

const dryRunFlag = "dry-run"

// SetCmd is a command that sets git identity based on the configuration file.
var SetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set git identity",
	Run: func(cmd *cobra.Command, args []string) {
		ensureGitDir()

		confPath, err := config.LocateConfFile()
		cli.PrintlnExit(err)

		conf, err := config.Load(confPath)
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

		hostConf, err := conf.Get((*repo).Host)
		cli.PrintlnExit(err)

		cred := findBestMatch(*hostConf, *repo)
		if cred == nil {
			cli.PrintlnExit(fmt.Errorf("cannot find a match for host %q", (*repo).Host))
		}

		dryRun, err := cmd.Flags().GetBool(dryRunFlag)
		cli.PrintlnExit(err)

		if !dryRun {
			cli.PrintlnExit(
				git.SetConfig("--local", "user.name", cred.Name),
			)
			cli.PrintlnExit(
				git.SetConfig("--local", "user.email", cred.Email),
			)
		}

		fmt.Printf("set user: %s <%s>\n", cred.Name, cred.Email)
	},
}

func init() {
	SetCmd.Flags().Bool(dryRunFlag, false, "dry run")
}

func ensureGitDir() {
	dir, err := os.Getwd()
	cli.PrintlnExit(err)

	ok, err := git.IsGitDir(dir)
	cli.PrintlnExit(err)

	if !ok {
		fmt.Printf(`Error: %q is not a git directory

Make sure you are executing zit inside a git directory.

If you are, perhaps you have forgotten to initialize a new repository? In this
case, run:

    git init
`, dir)
		os.Exit(1)
	}
}
