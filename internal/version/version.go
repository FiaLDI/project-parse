package version

// Values are overwritten via -ldflags at build time.
var (
	Version   = "0.1.0-dev"
	Commit    = "none"
	BuildDate = "unknown"
)

// String returns a human-readable version line.
func String() string {
	return Version + " (commit=" + Commit + ", built=" + BuildDate + ")"
}
