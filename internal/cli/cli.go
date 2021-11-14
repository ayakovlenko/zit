package cli

import (
	"fmt"
	"os"
)

// PrintlnExit TODO
func PrintlnExit(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
