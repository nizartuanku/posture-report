// Package sources reads the open findings out of another Sentinel tool's SQLite
// database, read-only. Every tool in the line writes the same findings schema
// (that's the whole point of Sentinel Core), so Posture Report can open each
// tool's db beside it, read what it found, and never write a byte. One db file
// may hold more than one module; Read groups by module into posture.Items.
package sources

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/nizartuanku/posture-report/core"
	"github.com/nizartuanku/posture-report/posture"
)

// openDB is swappable in tests; production opens the file read-only.
var openDB = func(path string) (*sql.DB, error) {
	return sql.Open("sqlite3", "file:"+path+"?mode=ro&_busy_timeout=2000")
}

// Read returns the OPEN findings in one Sentinel database, grouped by module.
// A file that has no findings table (not a Sentinel db) is a user-facing error,
// so the caller can warn and skip it rather than failing the whole report.
func Read(path string) ([]posture.Item, error) {
	db, err := openDB(path)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	return readDB(db, path)
}

func readDB(db *sql.DB, path string) ([]posture.Item, error) {
	rows, err := db.Query(`
SELECT module, target, check_id, title, severity, status, remediation, evidence,
       first_seen, last_seen
FROM findings WHERE status = 'open'`)
	if err != nil {
		return nil, fmt.Errorf("%s: not a readable Sentinel database (%w)", filepath.Base(path), err)
	}
	defer rows.Close()

	byModule := map[string][]core.Finding{}
	for rows.Next() {
		var (
			f        core.Finding
			sev, st  string
			evidence sql.NullString
			first    time.Time
			last     time.Time
		)
		if err := rows.Scan(&f.Module, &f.Target, &f.Check, &f.Title, &sev, &st,
			&f.Remediation, &evidence, &first, &last); err != nil {
			return nil, err
		}
		f.Severity = core.Severity(sev)
		f.Status = core.FindingStatus(st)
		f.FirstSeen, f.LastSeen = first, last
		if evidence.Valid && evidence.String != "" {
			_ = json.Unmarshal([]byte(evidence.String), &f.Evidence)
		}
		byModule[f.Module] = append(byModule[f.Module], f)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	modules := make([]string, 0, len(byModule))
	for m := range byModule {
		modules = append(modules, m)
	}
	sort.Strings(modules)

	items := make([]posture.Item, 0, len(modules))
	for _, m := range modules {
		items = append(items, posture.Item{
			Product:  posture.ProductName(m),
			Module:   m,
			Findings: byModule[m],
		})
	}
	return items, nil
}

// Discover finds Sentinel database files in a directory (non-recursive). Used by
// the "-dir" convenience flag so a single-box deployment needs no per-file list.
func Discover(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if ext := filepath.Ext(name); ext == ".db" || ext == ".sqlite" {
			out = append(out, filepath.Join(dir, name))
		}
	}
	sort.Strings(out)
	return out, nil
}
