package version

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	_ "embed" // embed
)

//go:embed version.txt
var version string

// Tag returns application version tag.
func Tag() string {
	return strings.TrimSpace(version)
}

// VersionCmd is a command that prints the app version.
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Tag())
	},
}
