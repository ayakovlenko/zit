package version

import (
	"context"
	"fmt"
	"strings"

	"github.com/urfave/cli/v3"
)

func VersionCmd(version string) *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: "Print version",
		Action: func(_ context.Context, _ *cli.Command) error {
			fmt.Println(strings.TrimSpace(version))

			return nil
		},
	}
}
