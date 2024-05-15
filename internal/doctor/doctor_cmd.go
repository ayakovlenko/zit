package doctor

import "github.com/urfave/cli/v2"

var DoctorCmd = &cli.Command{
	Name:  "doctor",
	Usage: "Check git setup for potential problems",
	Action: func(_ *cli.Context) error {
		return runChecks()
	},
}
