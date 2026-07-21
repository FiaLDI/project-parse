package scanner_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/FiaLDI/project-parse/internal/ports"
	"github.com/FiaLDI/project-parse/internal/scanner"
)

func TestScanIndexesAndExcludes(t *testing.T) {
	root := t.TempDir()
	write := func(rel, body string) {
		t.Helper()
		path := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	write("go.mod", "module example")
	write("cmd/main.go", "package main")
	write("node_modules/pkg/index.js", "ignored")
	write("web/node_modules/x.js", "ignored")
	write(".git/config", "ignored")
	write("package.json", `{"name":"demo"}`)
	write(".github/workflows/ci.yml", "name: ci")

	s := scanner.New()
	pctx, err := s.Scan(context.Background(), root, ports.ScanOptions{
		Exclude: []string{
			"**/node_modules/**",
			"**/.git/**",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		path string
		want bool
	}{
		{name: "go_mod", path: "go.mod", want: true},
		{name: "main_go", path: "cmd/main.go", want: true},
		{name: "package_json", path: "package.json", want: true},
		{name: "workflow", path: ".github/workflows/ci.yml", want: true},
		{name: "root_node_modules", path: "node_modules/pkg/index.js", want: false},
		{name: "nested_node_modules", path: "web/node_modules/x.js", want: false},
		{name: "git", path: ".git/config", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pctx.Files.Has(tt.path); got != tt.want {
				t.Fatalf("Has(%s)=%v want %v (count=%d markers=%v)", tt.path, got, tt.want, pctx.FileCount(), pctx.Files.Markers())
			}
		})
	}

	for _, marker := range []string{"go.mod", "package.json", ".github/workflows"} {
		if !pctx.MarkerExists(marker) {
			t.Fatalf("missing marker %q in %v", marker, pctx.Files.Markers())
		}
	}
}

func TestScanIncludeFilter(t *testing.T) {
	root := t.TempDir()
	files := []string{"a.go", "b.txt", "dir/c.go"}
	for _, rel := range files {
		path := filepath.Join(root, rel)
		_ = os.MkdirAll(filepath.Dir(path), 0o755)
		if err := os.WriteFile(path, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	s := scanner.New()
	pctx, err := s.Scan(context.Background(), root, ports.ScanOptions{
		Include: []string{"**/*.go"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if pctx.FileCount() != 2 {
		t.Fatalf("file count=%d want 2; markers=%v", pctx.FileCount(), pctx.Files.Markers())
	}
	if pctx.Files.Has("b.txt") {
		t.Fatal("b.txt should be excluded by include filter")
	}
}

func TestScanRootMustBeDirectory(t *testing.T) {
	root := t.TempDir()
	file := filepath.Join(root, "not-a-dir")
	if err := os.WriteFile(file, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := scanner.New().Scan(context.Background(), file, ports.ScanOptions{})
	if err == nil {
		t.Fatal("expected error for non-directory root")
	}
}
