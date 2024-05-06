package main

import (
	"fmt"
	"os"
	"zit/internal/cli"
	"zit/internal/config"
	"zit/internal/identity"
	"zit/internal/version"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "zit",
	Short: "git identity manager",
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Config commands",
}

var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Config path",
	RunE: func(cmd *cobra.Command, args []string) error {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		envVar := os.Getenv(config.EnvVarName)

		configFile, err := config.LocateConfFile(afero.OsFs{}, homeDir, envVar)
		if err != nil {
			return err
		}

		fmt.Println(configFile)

		return nil
	},
}

func init() {
	configCmd.AddCommand(
		configPathCmd,
	)

	rootCmd.AddCommand(
		identity.SetCmd,
		version.VersionCmd,
		cli.DoctorCmd,
		configCmd,
	)
}
