package main

import (
	"zit/cli"
	"zit/clone"
	"zit/doctor"
	"zit/identity"
	"zit/version"

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
		doctor.DoctorCmd,
		clone.CloneCmd,
	)
}
