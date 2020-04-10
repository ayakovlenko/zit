package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

// VersionCmd TODO
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v1.0.0+1")
	},
}
