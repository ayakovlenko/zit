package main

import (
	"zit/cli"
	"zit/cred"
	"zit/doctor"
	"zit/version"
)

func main() {
	cli.PrintlnExit(cred.SetCredCmd.Execute())
}

func init() {
	cred.SetCredCmd.AddCommand(
		version.VersionCmd,
		doctor.DoctorCmd,
	)
}
