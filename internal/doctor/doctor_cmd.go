package doctor

import (
	"context"

	"github.com/urfave/cli/v3"
)

var DoctorCmd = &cli.Command{
	Name:  "doctor",
	Usage: "Check git setup for potential problems",
	Action: func(_ context.Context, _ *cli.Command) error {
		return runChecks()
	},
}
