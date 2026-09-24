// posturereport combines the open findings from every Hexward tool a company
// runs into one security-posture report — a single score, the priorities to fix
// first, and an executive + technical view that prints to PDF.
//
//	posturereport -dir /var/lib/hexward      # auto-discover the tools' databases
//	posturereport -dbs certlight.db,asm.db    # or list them explicitly
//	posturereport -out posture.html           # write the report once (cron/monthly)
//
// It reads the databases read-only, runs no scans, and changes nothing.
package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/nizartuanku/posture-report/license"
	"github.com/nizartuanku/posture-report/posture"
	"github.com/nizartuanku/posture-report/report"
	"github.com/nizartuanku/posture-report/sources"
	"github.com/nizartuanku/posture-report/web"
)

// issuerPublicKeyB64 is baked in at build time by the release process.
// Empty → free edition.
var issuerPublicKeyB64 = ""

const moduleID = "posturereport"

// postureTierLimits: free reads 3 tools, Pro 20, Team unlimited (MaxTargets = source cap).
var postureTierLimits = map[license.Tier]license.Limits{
	license.TierFree: {MaxTargets: 3},
	license.TierPro:  {MaxTargets: 20},
	license.TierTeam: {MaxTargets: 0},
}

func main() {
	listen := flag.String("listen", "127.0.0.1:8432", "dashboard listen address")
	title := flag.String("title", "Security Posture", "report title (usually the company name)")
	dbsCSV := flag.String("dbs", "", "comma-separated Hexward database paths")
	dir := flag.String("dir", "", "directory to auto-discover Hexward .db files")
	out := flag.String("out", "", "write the report to this HTML file once and exit (cron/monthly mode)")
	licFile := flag.String("license", "posturereport-license.key", "license key file")
	aiURL := flag.String("ai-assist-url", os.Getenv("POSTURE_REPORT_AI_ASSIST_URL"), "optional hexward-ai sidecar URL for AI-narrated explanations, e.g. http://127.0.0.1:8435 (off when empty)")
	aiKeyFile := flag.String("ai-assist-key-file", os.Getenv("POSTURE_REPORT_AI_ASSIST_KEY_FILE"), "API key file for a dedicated AI host or your own OpenAI-compatible endpoint (Pro/Team)")
	aiLang := flag.String("ai-assist-lang", os.Getenv("POSTURE_REPORT_AI_ASSIST_LANG"), "language of AI explanations: en (default) or id")
	aiNoThinking := flag.Bool("ai-assist-no-thinking", os.Getenv("POSTURE_REPORT_AI_ASSIST_NO_THINKING") == "1", "disable reasoning mode (Qwen3 enterprise profiles)")
	flag.Parse()

	var pub ed25519.PublicKey
	if issuerPublicKeyB64 != "" {
		if b, err := base64.StdEncoding.DecodeString(issuerPublicKeyB64); err == nil {
			pub = ed25519.PublicKey(b)
		}
	}
	key := ""
	if b, err := os.ReadFile(*licFile); err == nil {
		key = strings.TrimSpace(string(b))
	}
	act := license.Activate(pub, moduleID, key, time.Now())
	srcCap := postureTierLimits[act.Tier].MaxTargets

	paths := gatherPaths(*dbsCSV, *dir)

	load := func() posture.Report {
		used := paths
		if srcCap > 0 && len(used) > srcCap {
			used = used[:srcCap]
		}
		var items []posture.Item
		for _, p := range used {
			its, err := sources.Read(p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "posturereport: skipping %s: %v\n", p, err)
				continue
			}
			items = append(items, its...)
		}
		return posture.Build(*title, items, time.Now())
	}

	notice := act.Notice
	if srcCap > 0 && len(paths) > srcCap {
		notice = strings.TrimSpace(notice + fmt.Sprintf(" Free edition reads %d tools; %d are configured — upgrade to include them all.", srcCap, len(paths)))
	}

	if *out != "" {
		if err := os.WriteFile(*out, []byte(report.HTML(load())), 0o644); err != nil {
			fmt.Fprintln(os.Stderr, "posturereport: "+err.Error())
			os.Exit(1)
		}
		fmt.Printf("Wrote %s\n", *out)
		return
	}

	server := &web.Server{Load: load, Tier: string(act.Tier), Notice: notice}
	aiAssist, aiErr := web.NewAIAssist(web.AIConfig{URL: *aiURL, KeyFile: *aiKeyFile, Language: *aiLang, NoThinking: *aiNoThinking})
	if aiErr != nil {
		fmt.Fprintln(os.Stderr, "posture-report: "+aiErr.Error())
		os.Exit(2)
	}
	server.AI = aiAssist
	if aiAssist != nil {
		fmt.Fprintf(os.Stderr, "posture-report: AI Assist on — explanations from %s (language %s)\n", aiAssist.Endpoint, aiAssist.Language)
	}
	fmt.Printf("Posture Report — %s edition (%d tool database(s) configured)\n", act.Tier, len(paths))
	fmt.Printf("Dashboard: http://%s\n", *listen)
	if err := http.ListenAndServe(*listen, server.Handler()); err != nil {
		fmt.Fprintln(os.Stderr, "posturereport: "+err.Error())
		os.Exit(1)
	}
}

// gatherPaths merges explicit -dbs entries with -dir discovery, de-duplicated,
// order-stable (explicit first).
func gatherPaths(csv, dir string) []string {
	seen := map[string]bool{}
	var out []string
	add := func(p string) {
		p = strings.TrimSpace(p)
		if p != "" && !seen[p] {
			seen[p] = true
			out = append(out, p)
		}
	}
	for _, p := range strings.Split(csv, ",") {
		add(p)
	}
	if dir != "" {
		if found, err := sources.Discover(dir); err == nil {
			for _, p := range found {
				add(p)
			}
		} else {
			fmt.Fprintf(os.Stderr, "posturereport: -dir %s: %v\n", dir, err)
		}
	}
	return out
}
