package app

import (
	"context"
	"errors"
	"testing"

	"github.com/FiaLDI/project-parse/internal/config"
)

func TestScanRequiresScanner(t *testing.T) {
	a := New(config.Default(), nil, Deps{})
	_, err := a.Scan(context.Background(), ".")
	if !errors.Is(err, ErrDependencyMissing) {
		t.Fatalf("expected ErrDependencyMissing, got %v", err)
	}
}

func TestListPluginsRequiresRegistry(t *testing.T) {
	a := New(config.Default(), nil, Deps{})
	_, err := a.ListPlugins()
	if !errors.Is(err, ErrDependencyMissing) {
		t.Fatalf("expected ErrDependencyMissing, got %v", err)
	}
}

func TestDoctorCriticalPath(t *testing.T) {
	a := New(config.Default(), nil, Deps{})
	res, err := a.Doctor(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if res.HasCriticalFailures() {
		t.Fatalf("unexpected critical failures: %+v", res.Checks)
	}
}
