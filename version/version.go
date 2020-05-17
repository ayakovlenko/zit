package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "v2.0.1"

// VersionCmd is a command that prints the app version.
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}
