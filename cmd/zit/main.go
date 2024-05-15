package main

import (
	"fmt"
	"os"
	"zit/internal/doctor"
	"zit/internal/identity"
	"zit/internal/version"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "zit",
		Usage: "git identity manager",
		Commands: []*cli.Command{
			version.VersionCmd,
			doctor.DoctorCmd,
			identity.SetCmd,
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
