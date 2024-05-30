package doctor

import (
	"fmt"
	"strings"
	"zit/internal/gitutil"
	"zit/pkg/git"
)

type check struct {
	Name      string
	CheckFunc func() (bool, error)
	// FixFunc   func() (bool, error)
}

func runChecks() error {
	gitClient := git.NewGitClient()

	useConfigOnly := check{
		Name: "git config --global user.useConfigOnly true",
		CheckFunc: func() (bool, error) {
			out, err := gitutil.GetConfig(gitClient, "--global", "user.useConfigOnly")
			if err != nil {
				return false, err
			}
			return out == "true", nil
		},
	}

	globalUserName := check{
		Name: "git config --unset-all --global user.name",
		CheckFunc: func() (bool, error) {
			out, err := gitutil.GetConfig(gitClient, "--global", "user.name")
			if err != nil {
				return false, err
			}
			return out == "", nil
		},
	}

	globalEmail := check{
		Name: "git config --unset-all --global user.email",
		CheckFunc: func() (bool, error) {
			out, err := gitutil.GetConfig(gitClient, "--global", "user.email")
			if err != nil {
				return false, err
			}
			return out == "", nil
		},
	}

	systemUserName := check{
		Name: "git config --unset-all --system user.name",
		CheckFunc: func() (bool, error) {
			out, err := gitutil.GetConfig(gitClient, "--system", "user.name")
			if err != nil {
				return false, err
			}
			return out == "", nil
		},
	}

	systemEmail := check{
		Name: "git config --unset-all --system user.email",
		CheckFunc: func() (bool, error) {
			out, err := gitutil.GetConfig(gitClient, "--system", "user.email")
			if err != nil {
				return false, err
			}
			return out == "", nil
		},
	}

	var checks = []check{
		useConfigOnly,
		globalUserName,
		globalEmail,
		systemUserName,
		systemEmail,
	}

	outs := []string{}
	for _, check := range checks {
		ok, err := check.CheckFunc()
		if err != nil {
			return err
		}
		outs = append(outs, fmtResult(check.Name, ok))
	}
	fmt.Println(strings.Join(outs, "\n"))
	return nil
}

// format check run result as:
//
// "- [x] check name" for a successful check
//
// "- [ ] check name" for a failed check
func fmtResult(name string, ok bool) string {
	tick := " "
	if ok {
		tick = "x"
	}
	return fmt.Sprintf("- [%s] %s", tick, name)
}
