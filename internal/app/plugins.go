package app

import (
	"fmt"

	"github.com/FiaLDI/project-parse/internal/ports"
)

// ListPlugins returns plugin metadata according to current config policy.
func (a *App) ListPlugins() ([]ports.PluginMeta, error) {
	if a.registry == nil {
		return nil, fmt.Errorf("%w: registry", ErrDependencyMissing)
	}
	return a.registry.List(a.cfg.Plugins.Enabled, a.cfg.Plugins.Disabled), nil
}
