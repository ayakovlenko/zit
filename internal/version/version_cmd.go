package version

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

var VersionCmd = &cli.Command{
	Name:  "version",
	Usage: "Print version",
	Action: func(_ *cli.Context) error {
		fmt.Println(strings.TrimSpace(version))
		return nil
	},
}
