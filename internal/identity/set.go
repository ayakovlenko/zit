package identity

import (
	"fmt"
	"os"
	"zit/internal/config"
	"zit/internal/git"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

const dryRunFlag = "dry-run"

// SetCmd is a command that sets git identity based on the configuration file.
var SetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set git identity",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ensureGitDir(); err != nil {
			return err
		}

		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		// check ZIT_CONFIG env variable
		envVar, _ := os.LookupEnv(config.EnvVarName)

		confPath, err := config.LocateConfFile(afero.NewOsFs(), home, envVar)
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

		hostConf, err := conf.Get((*repo).Host)
		if err != nil {
			return err
		}

		cred := findBestMatch(*hostConf, *repo)
		if cred == nil {
			return fmt.Errorf("cannot find a match for host %q", (*repo).Host)
		}

		dryRun, err := cmd.Flags().GetBool(dryRunFlag)
		if err != nil {
			return err
		}

		if !dryRun {
			if err := git.SetConfig("--local", "user.name", cred.Name); err != nil {
				return err
			}
			if err := git.SetConfig("--local", "user.email", cred.Email); err != nil {
				return err
			}
		}

		fmt.Printf("set user: %s <%s>\n", cred.Name, cred.Email)

		return nil
	},
}

func init() {
	SetCmd.Flags().Bool(dryRunFlag, false, "dry run")
}

func ensureGitDir() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	ok, err := git.IsGitDir(dir)
	if err != nil {
		return err
	}

	if !ok {
		fmt.Fprintf(os.Stderr, `Error: %q is not a git directory

Make sure you are executing zit inside a git directory.

If you are, perhaps you have forgotten to initialize a new repository? In this
case, run:

    git init
`, dir)
		os.Exit(1)
	}

	return nil
}
