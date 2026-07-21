package domain

import "testing"

func TestFileIndexBasics(t *testing.T) {
	idx := NewFileIndex()
	idx.Add(FileMeta{RelPath: "go.mod", Name: "go.mod", Ext: ".mod"})
	idx.Add(FileMeta{RelPath: "cmd/main.go", Name: "main.go", Ext: ".go"})
	idx.SetMarkers([]string{"go.mod"})

	tests := []struct {
		name string
		got  any
		want any
	}{
		{name: "len", got: idx.Len(), want: 2},
		{name: "has_go_mod", got: idx.Has("go.mod"), want: true},
		{name: "missing", got: idx.Has("package.json"), want: false},
		{name: "by_name", got: len(idx.ByName("main.go")), want: 1},
		{name: "by_ext", got: len(idx.ByExt(".go")), want: 1},
		{name: "markers", got: len(idx.Markers()), want: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("got %#v want %#v", tt.got, tt.want)
			}
		})
	}
}

func TestProjectContextHelpers(t *testing.T) {
	idx := NewFileIndex()
	idx.SetMarkers([]string{"package.json"})
	pc := ProjectContext{Root: ".", Files: idx}

	if !pc.MarkerExists("package.json") {
		t.Fatal("expected package.json marker")
	}
	if pc.MarkerExists("go.mod") {
		t.Fatal("unexpected go.mod marker")
	}
	if pc.FileCount() != 0 {
		t.Fatalf("expected 0 files, got %d", pc.FileCount())
	}
}
