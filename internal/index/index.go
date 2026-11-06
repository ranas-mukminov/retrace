// Package index is the local search index: a SQLite FTS5 database that
// lives in the user's data directory and never leaves the machine.
package index

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

// Result is one search hit.
type Result struct {
	Path    string
	Snippet string
}

func dbPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".local", "share", "retrace")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	return filepath.Join(dir, "retrace.db"), nil
}

func openDB() (*sql.DB, error) {
	p, err := dbPath()
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", p)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS docs(path TEXT PRIMARY KEY, mtime INTEGER);
		CREATE VIRTUAL TABLE IF NOT EXISTS fts USING fts5(path UNINDEXED, content, tokenize='porter unicode61');`)
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

var textExts = map[string]bool{
	".txt": true, ".md": true, ".markdown": true, ".org": true,
	".rst": true, ".adoc": true, ".log": true,
}

// IndexDir indexes text files under root, skipping unchanged files.
func IndexDir(root string) error {
	db, err := openDB()
	if err != nil {
		return err
	}
	defer db.Close()

	indexed, skipped := 0, 0
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // skip unreadable entries
		}
		if d.IsDir() {
			if strings.HasPrefix(d.Name(), ".") && path != root {
				return filepath.SkipDir
			}
			return nil
		}
		if !textExts[strings.ToLower(filepath.Ext(path))] {
			skipped++
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		mtime := info.ModTime().Unix()
		var old int64
		if err := db.QueryRow(`SELECT mtime FROM docs WHERE path = ?`, path).Scan(&old); err == nil && old == mtime {
			return nil // unchanged
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		if !isText(data) {
			skipped++
			return nil
		}
		if _, err := db.Exec(`INSERT OR REPLACE INTO docs(path, mtime) VALUES(?, ?)`, path, mtime); err != nil {
			return nil
		}
		if _, err := db.Exec(`DELETE FROM fts WHERE path = ?`, path); err != nil {
			return nil
		}
		if _, err := db.Exec(`INSERT INTO fts(path, content) VALUES(?, ?)`, path, string(data)); err != nil {
			return nil
		}
		indexed++
		return nil
	})
	if err != nil {
		return err
	}
	fmt.Printf("indexed %d files (%d skipped)\n", indexed, skipped)
	return nil
}

func isText(data []byte) bool {
	// Reject content with NUL bytes: binary files mislabeled as text.
	for _, b := range data {
		if b == 0 {
			return false
		}
	}
	return true
}

// ftsQuery turns a user's words into a safe FTS5 query: every term must
// appear, order does not matter, no FTS5 syntax leaks through.
func ftsQuery(query string) string {
	var terms []string
	for _, w := range strings.Fields(query) {
		terms = append(terms, `"`+strings.ReplaceAll(w, `"`, `""`)+`"`)
	}
	return strings.Join(terms, " AND ")
}

// Search returns up to 20 hits ranked by bm25, with a text snippet each.
func Search(query string) ([]Result, error) {
	db, err := openDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(`SELECT path, snippet(fts, 1, '', '', ' … ', 32) FROM fts WHERE fts MATCH ? ORDER BY bm25(fts) LIMIT 20`, ftsQuery(query))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Result
	for rows.Next() {
		var r Result
		if err := rows.Scan(&r.Path, &r.Snippet); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// Count returns the number of indexed documents.
func Count() (int, error) {
	db, err := openDB()
	if err != nil {
		return 0, err
	}
	defer db.Close()
	var n int
	err = db.QueryRow(`SELECT COUNT(*) FROM docs`).Scan(&n)
	return n, err
}
