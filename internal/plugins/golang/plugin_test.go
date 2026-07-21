package golang_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/FiaLDI/project-parse/internal/cache"
	"github.com/FiaLDI/project-parse/internal/domain"
	"github.com/FiaLDI/project-parse/internal/plugins/golang"
)

func TestAnalyzeGoMod(t *testing.T) {
	dir := t.TempDir()
	mod := filepath.Join(dir, "go.mod")
	if err := os.WriteFile(mod, []byte("module example.com/demo\ngo 1.24\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	idx := domain.NewFileIndex()
	idx.Add(domain.FileMeta{
		Path: mod, RelPath: "go.mod", Name: "go.mod", Ext: ".mod",
	})
	idx.SetMarkers([]string{"go.mod"})
	pctx := domain.ProjectContext{Root: dir, Files: idx}

	p := golang.New(cache.New(cache.Options{}))
	if !p.Supports(pctx) {
		t.Fatal("expected support")
	}
	res, err := p.Analyze(context.Background(), pctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 1 {
		t.Fatalf("findings=%+v", res.Findings)
	}
	if res.Findings[0].Attributes["module"] != "example.com/demo" {
		t.Fatalf("attrs=%v", res.Findings[0].Attributes)
	}
}
