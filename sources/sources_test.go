package sources

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// makeDB writes a minimal Sentinel-schema database with the given rows.
func makeDB(t *testing.T, name string, rows [][2]string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := db.Exec(`CREATE TABLE findings(
		id TEXT PRIMARY KEY, fingerprint TEXT, module TEXT, target TEXT, check_id TEXT,
		title TEXT, severity TEXT, status TEXT, remediation TEXT, evidence TEXT,
		group_id TEXT, first_seen TIMESTAMP, last_seen TIMESTAMP, resolved_at TIMESTAMP, absent_streak INT)`); err != nil {
		t.Fatal(err)
	}
	for i, r := range rows {
		module, status := r[0], r[1]
		_, err := db.Exec(`INSERT INTO findings(id,fingerprint,module,target,check_id,title,severity,status,remediation,first_seen,last_seen)
			VALUES(?,?,?,?,?,?,?,?,?,datetime('now'),datetime('now'))`,
			name+string(rune('a'+i)), "fp", module, "t", "chk", "a problem", "high", status, "fix it")
		if err != nil {
			t.Fatal(err)
		}
	}
	return path
}

func TestReadReturnsOnlyOpenGroupedByModule(t *testing.T) {
	path := makeDB(t, "certwatch.db", [][2]string{
		{"certwatch", "open"},
		{"certwatch", "open"},
		{"certwatch", "resolved"}, // must be excluded
	})
	items, err := Read(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("want 1 module group, got %d", len(items))
	}
	if items[0].Module != "certwatch" || items[0].Product != "CertLight" {
		t.Errorf("module/product wrong: %+v", items[0])
	}
	if len(items[0].Findings) != 2 {
		t.Errorf("want 2 open findings (resolved excluded), got %d", len(items[0].Findings))
	}
}

func TestReadNonSentinelDBErrors(t *testing.T) {
	path := filepath.Join(t.TempDir(), "empty.db")
	db, _ := sql.Open("sqlite3", path)
	db.Exec(`CREATE TABLE other(x INT)`)
	db.Close()
	if _, err := Read(path); err == nil {
		t.Fatal("a database without a findings table should error, not silently return nothing")
	}
}

func TestDiscover(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "a.db"), []byte("x"), 0o644)
	os.WriteFile(filepath.Join(dir, "b.sqlite"), []byte("x"), 0o644)
	os.WriteFile(filepath.Join(dir, "notes.txt"), []byte("x"), 0o644)
	found, err := Discover(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(found) != 2 {
		t.Errorf("want 2 db files discovered, got %d: %v", len(found), found)
	}
}
