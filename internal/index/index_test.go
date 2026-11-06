package index_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ranas-mukminov/retrace/internal/index"
)

func TestIndexAndSearch(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "notes.md"), []byte("# Trip\nThe lighthouse at Kvirike is worth the climb.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "todo.txt"), []byte("buy train tickets\ncall the plumber\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := index.IndexDir(dir); err != nil {
		t.Fatal(err)
	}
	if n, err := index.Count(); err != nil || n != 2 {
		t.Fatalf("count = %d, err = %v, want 2", n, err)
	}

	results, err := index.Search("lighthouse")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("got %d results, want 1", len(results))
	}
	if got := results[0].Path; got != filepath.Join(dir, "notes.md") {
		t.Fatalf("path = %q, want notes.md", got)
	}

	// Re-indexing unchanged files adds nothing.
	if err := index.IndexDir(dir); err != nil {
		t.Fatal(err)
	}
	if n, err := index.Count(); err != nil || n != 2 {
		t.Fatalf("count after re-index = %d, err = %v, want 2", n, err)
	}

	// A query with no matches returns an empty slice, not an error.
	results, err = index.Search("nonexistentword")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		t.Fatalf("got %d results, want 0", len(results))
	}
}
