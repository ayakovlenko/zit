package config

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var ConfigShowJsonCmd = &cobra.Command{
	Use:   "show-json",
	Short: "???", // TODO
	Run: func(cmd *cobra.Command, args []string) {
		path, err := LocateConfFile()
		if err != nil {
			fmt.Printf("Error locating conf file: %s\n", err)
			os.Exit(1)
		}

		file, err := os.Open(path)
		if err != nil {
			fmt.Printf("Error opening conf file: %s\n", err)
			os.Exit(1)
		}

		confJSON, err := readConfJSON(path, file)
		if err != nil {
			fmt.Printf("Error reading conf file: %s\n", err)
			os.Exit(1)
		}

		fmt.Println(confJSON)
	},
}
