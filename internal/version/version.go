package version

import (
	_ "embed" // embed
)

//go:embed version.txt
var version string
