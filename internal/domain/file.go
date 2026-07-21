package domain

import "time"

// FileMeta describes a single filesystem entry discovered during scan.
type FileMeta struct {
	Path    string    `json:"path"`
	RelPath string    `json:"rel_path"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"mod_time"`
	Ext     string    `json:"ext"`
	Name    string    `json:"name"`
	IsDir   bool      `json:"is_dir"`
}

// FileIndex is a path-keyed catalogue of discovered files.
type FileIndex struct {
	byRelPath map[string]FileMeta
	markers   []string
}

// NewFileIndex creates an empty file index.
func NewFileIndex() *FileIndex {
	return &FileIndex{
		byRelPath: make(map[string]FileMeta),
		markers:   nil,
	}
}

// Add inserts or replaces a file entry.
func (idx *FileIndex) Add(meta FileMeta) {
	if idx.byRelPath == nil {
		idx.byRelPath = make(map[string]FileMeta)
	}
	idx.byRelPath[meta.RelPath] = meta
}

// Get returns a file by relative path.
func (idx *FileIndex) Get(relPath string) (FileMeta, bool) {
	if idx == nil {
		return FileMeta{}, false
	}
	meta, ok := idx.byRelPath[relPath]
	return meta, ok
}

// Has reports whether a relative path exists in the index.
func (idx *FileIndex) Has(relPath string) bool {
	_, ok := idx.Get(relPath)
	return ok
}

// Len returns the number of indexed entries.
func (idx *FileIndex) Len() int {
	if idx == nil {
		return 0
	}
	return len(idx.byRelPath)
}

// All returns a snapshot of all indexed files.
func (idx *FileIndex) All() []FileMeta {
	if idx == nil {
		return nil
	}
	out := make([]FileMeta, 0, len(idx.byRelPath))
	for _, meta := range idx.byRelPath {
		out = append(out, meta)
	}
	return out
}

// ByName returns files whose basename equals name.
func (idx *FileIndex) ByName(name string) []FileMeta {
	if idx == nil {
		return nil
	}
	var out []FileMeta
	for _, meta := range idx.byRelPath {
		if meta.Name == name {
			out = append(out, meta)
		}
	}
	return out
}

// ByExt returns files with the given extension (including leading dot, e.g. ".go").
func (idx *FileIndex) ByExt(ext string) []FileMeta {
	if idx == nil {
		return nil
	}
	var out []FileMeta
	for _, meta := range idx.byRelPath {
		if meta.Ext == ext {
			out = append(out, meta)
		}
	}
	return out
}

// SetMarkers records well-known marker file basenames discovered during scan.
func (idx *FileIndex) SetMarkers(markers []string) {
	idx.markers = append([]string(nil), markers...)
}

// Markers returns discovered marker basenames.
func (idx *FileIndex) Markers() []string {
	if idx == nil {
		return nil
	}
	return append([]string(nil), idx.markers...)
}
