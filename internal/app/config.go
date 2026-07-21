package app

import (
	"errors"

	"github.com/FiaLDI/project-parse/internal/config"
)

// ConfigInitOptions controls default config file creation.
type ConfigInitOptions struct {
	Path  string
	Force bool
}

// InitConfig writes a default parser.yaml configuration file.
func InitConfig(opts ConfigInitOptions) (string, error) {
	return config.WriteDefault(opts.Path, opts.Force)
}

// IsConfigExistsErr reports whether err is a config-already-exists error.
func IsConfigExistsErr(err error) bool {
	return errors.Is(err, config.ErrConfigExists)
}
