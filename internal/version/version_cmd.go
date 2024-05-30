package version

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

func VersionCmd(version string) *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: "Print version",
		Action: func(_ *cli.Context) error {
			fmt.Println(strings.TrimSpace(version))
			return nil
		},
	}
}
