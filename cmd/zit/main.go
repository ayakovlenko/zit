package main

import (
	"zit/internal/cli"
	"zit/internal/identity"
	"zit/internal/version"

	"github.com/spf13/cobra"
)

func main() {
	cli.PrintlnExit(rootCmd.Execute())
}

var rootCmd = &cobra.Command{
	Use:   "zit",
	Short: "git identity manager",
}

func init() {
	rootCmd.AddCommand(
		identity.SetCmd,
		version.VersionCmd,
		cli.DoctorCmd,
	)
}
