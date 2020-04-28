package config

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "???", // TODO
	Run: func(cmd *cobra.Command, args []string) {
		path, err := LocateConfFile()
		if err != nil {
			fmt.Printf("Error locating conf file: %s\n", err)
			os.Exit(1)
		}

		fmt.Print(path)
	},
}
