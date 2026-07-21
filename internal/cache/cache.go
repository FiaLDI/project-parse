package cache

import (
	"fmt"
	"io/fs"
	"os"
	"sync"
	"time"
)

// Options configures cache behaviour.
type Options struct {
	// MaxFileBytes caps how many bytes Read will load. 0 means unlimited.
	MaxFileBytes int64
}

type entry struct {
	size    int64
	modTime time.Time
	data    []byte
	err     error
}

// Cache is a thread-safe read-once file content cache.
type Cache struct {
	opts Options

	mu      sync.RWMutex
	entries map[string]*entry
	reads   map[string]int // filesystem read attempts per path
}

// New creates an empty file cache.
func New(opts Options) *Cache {
	return &Cache{
		opts:    opts,
		entries: make(map[string]*entry),
		reads:   make(map[string]int),
	}
}

// Exists reports whether path exists on disk (does not cache content).
func (c *Cache) Exists(path string) bool {
	_, err := os.Lstat(path)
	return err == nil
}

// Stat returns file info from the OS.
func (c *Cache) Stat(path string) (fs.FileInfo, error) {
	return os.Stat(path)
}

// Read returns file contents. Identical path+size+mtime hits reuse the cached bytes
// without a second filesystem read of the body.
func (c *Cache) Read(path string) ([]byte, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		return nil, fmt.Errorf("cache: %s is a directory", path)
	}

	size := info.Size()
	mod := info.ModTime()

	c.mu.RLock()
	if e, ok := c.entries[path]; ok && e.size == size && e.modTime.Equal(mod) {
		data, err := e.data, e.err
		c.mu.RUnlock()
		if err != nil {
			return nil, err
		}
		return clone(data), nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	if e, ok := c.entries[path]; ok && e.size == size && e.modTime.Equal(mod) {
		if e.err != nil {
			return nil, e.err
		}
		return clone(e.data), nil
	}

	c.reads[path]++
	var data []byte
	if c.opts.MaxFileBytes > 0 && size > c.opts.MaxFileBytes {
		err = fmt.Errorf("cache: file %s exceeds max size (%d > %d)", path, size, c.opts.MaxFileBytes)
	} else {
		data, err = os.ReadFile(path)
	}
	c.entries[path] = &entry{
		size:    size,
		modTime: mod,
		data:    data,
		err:     err,
	}
	if err != nil {
		return nil, err
	}
	return clone(data), nil
}

// DiskReads returns how many times the filesystem body was read for path.
// Intended for tests verifying the single-read guarantee.
func (c *Cache) DiskReads(path string) int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.reads[path]
}

// Invalidate drops a cached entry for path.
func (c *Cache) Invalidate(path string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, path)
}

func clone(b []byte) []byte {
	if b == nil {
		return nil
	}
	out := make([]byte, len(b))
	copy(out, b)
	return out
}
