package main

import (
	"fmt"
	"os"
	"zit/internal/app"
	"zit/internal/config"
	"zit/internal/doctor"
	"zit/internal/identity"
	"zit/internal/version"
	"zit/pkg/xdg"

	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

func main() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	appConfig := app.NewConfig(
		afero.NewOsFs(),
		userHomeDir,
		os.Getenv(config.EnvVarName),
		os.Getenv(xdg.ConfigHome),
	)

	app := &cli.App{ //nolint: exhaustruct
		Name:  appConfig.AppName(),
		Usage: "git identity manager",
		Commands: []*cli.Command{
			version.VersionCmd(appConfig.AppVersion()),
			doctor.DoctorCmd,
			identity.SetCmd(appConfig),
			config.ConfigCmd(appConfig),
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
