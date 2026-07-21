package ports

// Registry holds plugins and applies enable/disable policy.
type Registry interface {
	Register(p Plugin)
	// Enabled returns plugins allowed by the given policy lists.
	// If enabled is empty, all registered plugins except those in disabled are returned.
	Enabled(enabled, disabled []string) []Plugin
	List(enabled, disabled []string) []PluginMeta
}
