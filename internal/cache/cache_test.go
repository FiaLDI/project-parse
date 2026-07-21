package cache

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadOnce(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a.txt")
	if err := os.WriteFile(path, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}

	c := New(Options{})
	for i := 0; i < 3; i++ {
		data, err := c.Read(path)
		if err != nil {
			t.Fatalf("read %d: %v", i, err)
		}
		if string(data) != "hello" {
			t.Fatalf("got %q", data)
		}
	}
	if got := c.DiskReads(path); got != 1 {
		t.Fatalf("disk reads = %d, want 1", got)
	}
}

func TestReadInvalidatesOnChange(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a.txt")
	if err := os.WriteFile(path, []byte("v1"), 0o644); err != nil {
		t.Fatal(err)
	}

	c := New(Options{})
	if _, err := c.Read(path); err != nil {
		t.Fatal(err)
	}

	// Ensure mtime/size change is observable.
	if err := os.WriteFile(path, []byte("v2-longer"), 0o644); err != nil {
		t.Fatal(err)
	}

	data, err := c.Read(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "v2-longer" {
		t.Fatalf("got %q", data)
	}
	if got := c.DiskReads(path); got != 2 {
		t.Fatalf("disk reads = %d, want 2", got)
	}
}

func TestMaxFileBytes(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "big.bin")
	if err := os.WriteFile(path, []byte("0123456789"), 0o644); err != nil {
		t.Fatal(err)
	}

	c := New(Options{MaxFileBytes: 4})
	if _, err := c.Read(path); err == nil {
		t.Fatal("expected size error")
	}
}

func TestExistsAndStat(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "x")
	if err := os.WriteFile(path, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	c := New(Options{})
	if !c.Exists(path) {
		t.Fatal("expected exists")
	}
	if c.Exists(filepath.Join(dir, "missing")) {
		t.Fatal("expected missing")
	}
	info, err := c.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Size() != 1 {
		t.Fatalf("size=%d", info.Size())
	}
}
