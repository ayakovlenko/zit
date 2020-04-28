package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "???", // TODO
	Run: func(cmd *cobra.Command, args []string) {
		fs := afero.NewOsFs()

		filename, err := defaultConfFileLocation()
		if err != nil {
			fmt.Printf("error while getting default conf file location: %s\n", err)
			os.Exit(1)
		}

		if err := createConfig(fs, filename, exampleConfig); err != nil {
			fmt.Printf("error while creating a conf file: %s\n", err)
			os.Exit(1)
		}
	},
}

const exampleConfig = `local User(name, email) = { name: name, email: email };

local user = {
  personal: User('jdoe', 'jdoe@users.noreply.github.com'), // Example user
  work: User('John Doe', 'john.doe@corp.com'), // Example user
};

// This is just an example.
// Feel free to delete it.
local example = {
  'github.com': {  // Public GitHub
    default: user.personal,
    overrides: [
      {  // Your employer's organization at github.com
        owner: 'corp',
        user: user.work,
      },
    ],
  },
  'github.corp.com': {  // Your employer's GitHub Enterprise
    default: user.work,
  },
}

// This is your real config.
local config = {}

config
`

func createConfig(fs afero.Fs, filename, contents string) error {
	dir := filepath.Dir(filename)

	exists, err := afero.DirExists(fs, dir)
	if err != nil {
		return err
	}

	if !exists {
		if err := fs.MkdirAll(dir, 0700); err != nil {
			return err
		}
	}

	exists, err = afero.Exists(fs, filename)
	if err != nil {
		return err
	}

	if exists {
		return ErrConfFileExists
	}

	if err := afero.WriteFile(
		fs, filename, []byte(contents), 0644,
	); err != nil {
		return err
	}

	return nil
}
