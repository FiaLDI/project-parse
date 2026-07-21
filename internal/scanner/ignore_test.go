package scanner

import "testing"

func TestGlobMatcher(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		path    string
		want    bool
	}{
		{name: "exact", pattern: "go.mod", path: "go.mod", want: true},
		{name: "star_ext", pattern: "*.go", path: "main.go", want: true},
		{name: "star_ext_nested_fail", pattern: "*.go", path: "cmd/main.go", want: false},
		{name: "double_star_file", pattern: "**/*.go", path: "cmd/main.go", want: true},
		{name: "node_modules_tree", pattern: "**/node_modules/**", path: "web/node_modules/pkg/index.js", want: true},
		{name: "node_modules_dir", pattern: "**/node_modules/**", path: "web/node_modules", want: true},
		{name: "git_dir_children", pattern: "**/.git/**", path: ".git/config", want: true},
		{name: "question", pattern: "file.?", path: "file.c", want: true},
		{name: "no_match", pattern: "**/vendor/**", path: "src/main.go", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newGlobMatcher([]string{tt.pattern})
			if got := m.Match(tt.path); got != tt.want {
				t.Fatalf("Match(%q,%q)=%v want %v", tt.pattern, tt.path, got, tt.want)
			}
		})
	}
}

func TestDirExcludedByPattern(t *testing.T) {
	tests := []struct {
		pattern string
		dir     string
		want    bool
	}{
		{pattern: "**/node_modules/**", dir: "node_modules", want: true},
		{pattern: "**/node_modules/**", dir: "apps/web/node_modules", want: true},
		{pattern: "**/node_modules/**", dir: "src", want: false},
		{pattern: "**/.git/**", dir: ".git", want: true},
		{pattern: "vendor/**", dir: "vendor", want: true},
		{pattern: "vendor/**", dir: "src", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.pattern+"_"+tt.dir, func(t *testing.T) {
			if got := dirExcludedByPattern(tt.pattern, tt.dir); got != tt.want {
				t.Fatalf("got %v want %v", got, tt.want)
			}
		})
	}
}
