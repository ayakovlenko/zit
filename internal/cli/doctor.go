package cli

import (
	"zit/internal/doctor"

	"github.com/spf13/cobra"
)

var DoctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check git setup for potential problems",
	RunE: func(cmd *cobra.Command, args []string) error {
		return doctor.RunChecks()
	},
}
