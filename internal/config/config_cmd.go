package config

import (
	"fmt"
	"os"
	"path/filepath"
	"zit/internal/app"
	"zit/pkg/xdg"

	"github.com/urfave/cli/v2"
)

const sampleConfig = `---
users:
  work: &work_user
    name: "John Doe"
    email: "john.doe@corp.com"

  personal:
    github_user: &personal_github_user
      name: "JD42"
      email: "JD42@users.noreply.github.com"
      sign:
        key: "~/.ssh/id_ed25519_github.pub"
        format: "ssh"

    gitlab_user: &personal_gitlab_user
      name: "JD42"
      email: "786972-JD42@users.noreply.gitlab.com"

hosts:

  github.com:
    default: *personal_github_user
    overrides:
      - owner: "corp"
        user: *work_user

  gitlab.com:
    default: *personal_gitlab_user
`

func ConfigCmd(appConfig app.Config) *cli.Command {
	return &cli.Command{
		Name:  "config",
		Usage: "Configuration management",
		Subcommands: []*cli.Command{
			configInitCmd(appConfig),
			configPathCmd(appConfig),
			configShowCmd(appConfig),
		},
	}
}

func configInitCmd(appConfig app.Config) *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "Initialize configuration file",
		Action: func(_ *cli.Context) error {
			confPath := xdg.LocateConfig(
				appConfig.AppName(),
				appConfig.UserHomeDir(),
				appConfig.XDGHomePathFromEnv(),
				appConfig.ConfigFilename(),
			)

			err := os.MkdirAll(filepath.Dir(confPath), 0755)
			if err != nil {
				return err
			}

			f, err := appConfig.FS().Create(confPath)
			if err != nil {
				return err
			}
			defer f.Close()

			if _, err := fmt.Fprint(f, sampleConfig); err != nil {
				return err
			}

			return nil
		},
	}
}

func configPathCmd(appConfig app.Config) *cli.Command {
	return &cli.Command{
		Name:  "path",
		Usage: "Show path to configuration file",
		Action: func(_ *cli.Context) error {
			confPath, err := LocateConfFile(appConfig)

			if err != nil {
				if _, ok := err.(*ConfigNotFoundError); ok {
					_, _ = fmt.Fprintf(os.Stderr, "error: config not found; run:\n\n    zit config init\n")

					os.Exit(1)
				}

				return err
			}

			_, _ = fmt.Fprintln(os.Stdout, confPath)

			return nil
		},
	}
}

func configShowCmd(appConfig app.Config) *cli.Command {
	return &cli.Command{
		Name:  "show",
		Usage: "Show configuration file contents",
		Action: func(_ *cli.Context) error {
			confPath, err := LocateConfFile(appConfig)

			if err != nil {
				if _, ok := err.(*ConfigNotFoundError); ok {
					_, _ = fmt.Fprintf(os.Stderr, "error: config not found; run:\n\n    zit config init\n")

					os.Exit(1)
				}

				return err
			}

			bs, err := os.ReadFile(confPath)
			if err != nil {
				return err
			}

			_, _ = os.Stdout.Write(bs)

			return nil
		},
	}
}
