package doctor

import (
	"fmt"
	"strings"
	"zit/cli"
	"zit/git"

	"github.com/spf13/cobra"
)

// DoctorCmd TODO
var DoctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check git setup for potential problems",
	Run: func(cmd *cobra.Command, args []string) {

		out := []string{
			check(
				"git config --global user.useConfigOnly true",
				func() bool {
					out, err := git.GetConfig("--global", "user.useConfigOnly")
					cli.PrintlnExit(err)
					return out == "true"
				},
			),
			check(
				"git config --unset-all --global user.name",
				func() bool {
					out, err := git.GetConfig("--global", "user.name")
					cli.PrintlnExit(err)
					return out == ""
				},
			),
			check(
				"git config --unset-all --global user.email",
				func() bool {
					out, err := git.GetConfig("--global", "user.email")
					cli.PrintlnExit(err)
					return out == ""
				},
			),
			check(
				"git config --unset-all --system user.name",
				func() bool {
					out, err := git.GetConfig("--system", "user.name")
					cli.PrintlnExit(err)
					return out == ""
				},
			),
			check(
				"git config --unset-all --system user.email",
				func() bool {
					out, err := git.GetConfig("--system", "user.email")
					cli.PrintlnExit(err)
					return out == ""
				},
			),
		}

		fmt.Println(strings.Join(out, "\n"))
	},
}

func check(name string, check func() bool) string {
	var tick string
	if ok := check(); ok {
		tick = "x"
	} else {
		tick = " "
	}

	return fmt.Sprintf("- [%s] %s", tick, name)
}
