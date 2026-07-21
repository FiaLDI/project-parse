package domain

import "time"

// ProjectContext is the immutable-enough snapshot passed to plugins.
// Plugins must not mutate shared maps; treat as read-only.
type ProjectContext struct {
	Root      string
	Files     *FileIndex
	StartedAt time.Time
}

// MarkerExists reports whether a scan marker basename was recorded.
func (pc ProjectContext) MarkerExists(name string) bool {
	if pc.Files == nil {
		return false
	}
	for _, m := range pc.Files.Markers() {
		if m == name {
			return true
		}
	}
	return false
}

// FileCount returns the number of indexed filesystem entries.
func (pc ProjectContext) FileCount() int {
	if pc.Files == nil {
		return 0
	}
	return pc.Files.Len()
}
