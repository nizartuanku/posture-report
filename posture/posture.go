// Package posture is the intelligence of Posture Report: it takes the open
// findings from every Sentinel tool a company runs (TLS, attack surface,
// canaries, CVEs, firewall, logs, DMARC, tenant posture, …) and folds them into
// one executive picture — a single security posture score, the handful of
// things to fix first, and the blind spots worth reviewing by hand. It is a
// pure aggregator: feed it findings, get a Report. No I/O, fully testable.
package posture

import (
	"sort"
	"strings"
	"time"

	"github.com/nizartuanku/posture-report/core"
)

// Item is one source tool's contribution: its open findings.
type Item struct {
	Product  string // display name, e.g. "CertWatch"
	Module   string // module id, e.g. "certwatch"
	Findings []core.Finding
}

// ProductSummary is one tool's line in the report's source table.
type ProductSummary struct {
	Product string
	Module  string
	Counts  map[string]int
	Open    int
}

// Report is the whole combined picture.
type Report struct {
	Title         string
	GeneratedAt   time.Time
	Sources       []ProductSummary
	Counts        map[string]int
	OpenTotal     int
	Score         int    // 0..100
	Rating        string // Good | Fair | Needs Attention | At Risk
	TopPriorities []core.Finding
	ManualReviews []core.Finding
	AllFindings   []core.Finding // severity-sorted, for the technical model
	Empty         bool
}

// severity weights for the posture score. Critical dominates; info never
// penalises (info is context / manual-review, not a live problem).
var weight = map[core.Severity]int{
	core.SeverityCritical: 25,
	core.SeverityHigh:     12,
	core.SeverityMedium:   4,
	core.SeverityLow:      1,
	core.SeverityInfo:     0,
}

func rank(s core.Severity) int {
	switch s {
	case core.SeverityCritical:
		return 4
	case core.SeverityHigh:
		return 3
	case core.SeverityMedium:
		return 2
	case core.SeverityLow:
		return 1
	default:
		return 0
	}
}

// isManualReview marks the honest "couldn't assess this — review by hand" notes
// the tools emit (info severity, check ending in manual-review). They are shown
// as blind spots, not counted as live findings or scored against.
func isManualReview(f core.Finding) bool {
	return f.Severity == core.SeverityInfo ||
		strings.Contains(f.Check, "manual-review") ||
		strings.Contains(strings.ToLower(f.Title), "manual review")
}

// Build folds the per-tool items into one report as of now.
func Build(title string, items []Item, now time.Time) Report {
	if strings.TrimSpace(title) == "" {
		title = "Security Posture"
	}
	rep := Report{
		Title:       title,
		GeneratedAt: now,
		Counts:      map[string]int{},
	}

	penalty := 0
	for _, it := range items {
		sum := ProductSummary{
			Product: it.Product, Module: it.Module,
			Counts: map[string]int{},
		}
		for _, f := range it.Findings {
			if isManualReview(f) {
				rep.ManualReviews = append(rep.ManualReviews, f)
				continue
			}
			sum.Counts[string(f.Severity)]++
			sum.Open++
			rep.Counts[string(f.Severity)]++
			rep.OpenTotal++
			rep.AllFindings = append(rep.AllFindings, f)
			penalty += weight[f.Severity]
		}
		rep.Sources = append(rep.Sources, sum)
	}

	// Score: start at 100, subtract the (capped) weighted penalty.
	if penalty > 100 {
		penalty = 100
	}
	rep.Score = 100 - penalty
	rep.Rating = ratingFor(rep.Score)

	// Severity-sorted views.
	sort.SliceStable(rep.AllFindings, func(i, j int) bool {
		return rank(rep.AllFindings[i].Severity) > rank(rep.AllFindings[j].Severity)
	})
	// Top priorities: the most severe handful.
	for _, f := range rep.AllFindings {
		if len(rep.TopPriorities) >= 6 {
			break
		}
		rep.TopPriorities = append(rep.TopPriorities, f)
	}

	rep.Empty = rep.OpenTotal == 0 && len(rep.ManualReviews) == 0 && len(items) == 0
	return rep
}

func ratingFor(score int) string {
	switch {
	case score >= 85:
		return "Good"
	case score >= 70:
		return "Fair"
	case score >= 50:
		return "Needs Attention"
	default:
		return "At Risk"
	}
}

// ProductName maps a Sentinel module id to its display name. Unknown modules
// fall back to a title-cased id so a new tool still renders sensibly.
func ProductName(module string) string {
	switch module {
	case "certwatch":
		return "CertLight"
	case "asm":
		return "Attack Surface Monitor"
	case "decoy":
		return "Decoy"
	case "patchlight":
		return "Patchlight"
	case "rulehawk":
		return "RuleHawk"
	case "loglight":
		return "Loglight"
	case "dmarcwatch":
		return "DmarcWatch"
	case "ruleforge":
		return "RuleForge"
	case "tenantwatch":
		return "TenantWatch"
	case "":
		return "Unknown"
	default:
		return strings.ToUpper(module[:1]) + module[1:]
	}
}
