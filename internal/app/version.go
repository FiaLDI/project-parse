package app

import "github.com/FiaLDI/project-parse/internal/version"

// Version returns the build version string.
func (a *App) Version() string {
	return version.String()
}
