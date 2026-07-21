package scanner

import (
	"path"
	"strings"
)

// globMatcher matches slash-separated relative paths against include/exclude globs.
// Supports *, ?, and ** (including patterns like **/node_modules/**).
type globMatcher struct {
	patterns []string
}

func newGlobMatcher(patterns []string) *globMatcher {
	norm := make([]string, 0, len(patterns))
	for _, p := range patterns {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		p = strings.ReplaceAll(p, "\\", "/")
		norm = append(norm, p)
	}
	return &globMatcher{patterns: norm}
}

// Match reports whether relPath matches any configured pattern.
func (m *globMatcher) Match(relPath string) bool {
	if m == nil || len(m.patterns) == 0 {
		return false
	}
	rel := strings.ReplaceAll(relPath, "\\", "/")
	rel = strings.TrimPrefix(rel, "./")
	for _, pattern := range m.patterns {
		if matchGlob(pattern, rel) {
			return true
		}
	}
	return false
}

// matchGlob implements a small subset of doublestar semantics for ** / * / ?.
func matchGlob(pattern, name string) bool {
	pattern = path.Clean("/" + pattern)
	name = path.Clean("/" + name)
	pattern = strings.TrimPrefix(pattern, "/")
	name = strings.TrimPrefix(name, "/")
	return matchParts(splitPath(pattern), splitPath(name))
}

func splitPath(p string) []string {
	if p == "" || p == "." {
		return nil
	}
	return strings.Split(p, "/")
}

func matchParts(pat, name []string) bool {
	for len(pat) > 0 {
		if pat[0] == "**" {
			pat = pat[1:]
			if len(pat) == 0 {
				return true
			}
			for i := 0; i <= len(name); i++ {
				if matchParts(pat, name[i:]) {
					return true
				}
			}
			return false
		}
		if len(name) == 0 {
			return false
		}
		if !matchSegment(pat[0], name[0]) {
			return false
		}
		pat = pat[1:]
		name = name[1:]
	}
	return len(name) == 0
}

func matchSegment(pat, seg string) bool {
	if pat == "*" {
		return true
	}
	// Simple ? and * within a single path segment.
	i, j := 0, 0
	star := -1
	match := 0
	for j < len(seg) {
		if i < len(pat) && (pat[i] == '?' || pat[i] == seg[j]) {
			i++
			j++
			continue
		}
		if i < len(pat) && pat[i] == '*' {
			star = i
			match = j
			i++
			continue
		}
		if star != -1 {
			i = star + 1
			match++
			j = match
			continue
		}
		return false
	}
	for i < len(pat) && pat[i] == '*' {
		i++
	}
	return i == len(pat)
}
