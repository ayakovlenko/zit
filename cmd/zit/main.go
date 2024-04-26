package main

import (
	"fmt"
	"os"
	"zit/internal/cli"
	"zit/internal/identity"
	"zit/internal/version"

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

func init() {
	rootCmd.AddCommand(
		identity.SetCmd,
		version.VersionCmd,
		cli.DoctorCmd,
	)
}
