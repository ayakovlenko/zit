package doctor

import (
	"fmt"
	"strings"
	"zit/internal/cli"
	"zit/internal/git"
)

// RunChecks runs checks.
func RunChecks() {

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
}

func check(name string, isOk func() bool) string {
	var tick string
	if isOk() {
		tick = "x"
	} else {
		tick = " "
	}

	return fmt.Sprintf("- [%s] %s", tick, name)
}
