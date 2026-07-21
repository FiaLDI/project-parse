package app

import "errors"

// ErrNotImplemented is returned by use-cases not yet wired in early stages.
var ErrNotImplemented = errors.New("not implemented")

// ErrDependencyMissing is returned when a required collaborator was not injected.
var ErrDependencyMissing = errors.New("required dependency is not configured")
