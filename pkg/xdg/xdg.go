package xdg

import "path/filepath"

// LocateConfig locates the configuration file for the given application name
// according to XDG Base Directory Specification v0.8:
// https://specifications.freedesktop.org/basedir-spec/latest/.
func LocateConfig(appName string, homeDir string, xdgConfigHomeVal string, configFilename string) string {
	if xdgConfigHomeVal != "" {
		return filepath.Join(xdgConfigHomeVal, appName, configFilename)
	}

	return filepath.Join(homeDir, ".config", appName, configFilename)
}
