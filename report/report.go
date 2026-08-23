// Package report renders a posture.Report as one self-contained HTML page with
// two models in it: an Executive view (a single score, the business-language
// summary, and the handful of things to fix first) and a Technical view (every
// open finding with its remediation, grouped by severity). The page is fully
// standalone — no external assets — so it prints straight to PDF and can be
// emailed or filed for an audit.
package report

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/nizartuanku/posture-report/core"
	"github.com/nizartuanku/posture-report/posture"
)

type bar struct {
	Label string
	Class string
	Count int
	Pct   int
}

type view struct {
	Rep        posture.Report
	Generated  string
	Summary    string
	Bars       []bar
	ToolCount  int
	ScoreClass string
}

// HTML renders the full two-model report page.
func HTML(rep posture.Report) string {
	v := view{
		Rep:       rep,
		Generated: rep.GeneratedAt.Format("2 January 2006, 15:04 MST"),
		Summary:   summarise(rep),
		ToolCount: len(rep.Sources),
	}
	sevs := []struct {
		s     core.Severity
		label string
		class string
	}{
		{core.SeverityCritical, "Critical", "crit"},
		{core.SeverityHigh, "High", "high"},
		{core.SeverityMedium, "Medium", "med"},
		{core.SeverityLow, "Low", "low"},
	}
	max := 1
	for _, x := range sevs {
		if c := rep.Counts[string(x.s)]; c > max {
			max = c
		}
	}
	for _, x := range sevs {
		c := rep.Counts[string(x.s)]
		v.Bars = append(v.Bars, bar{Label: x.label, Class: x.class, Count: c, Pct: c * 100 / max})
	}
	switch rep.Rating {
	case "Good":
		v.ScoreClass = "low"
	case "Fair":
		v.ScoreClass = "med"
	case "Needs Attention":
		v.ScoreClass = "high"
	default:
		v.ScoreClass = "crit"
	}

	var b bytes.Buffer
	if err := tmpl.Execute(&b, v); err != nil {
		return "<pre>report render error: " + template.HTMLEscapeString(err.Error()) + "</pre>"
	}
	return b.String()
}

// summarise writes the one-paragraph executive sentence from the numbers.
func summarise(rep posture.Report) string {
	if rep.OpenTotal == 0 {
		if len(rep.ManualReviews) > 0 {
			return fmt.Sprintf("No open findings across %d connected Sentinel tool(s). %d area(s) are flagged for manual review — worth a look, but nothing is actively failing.", len(rep.Sources), len(rep.ManualReviews))
		}
		return fmt.Sprintf("No open findings across %d connected Sentinel tool(s). Posture is clean this period — keep the tools running so it stays that way.", len(rep.Sources))
	}
	c := rep.Counts["critical"]
	h := rep.Counts["high"]
	var lead string
	switch {
	case c > 0:
		lead = fmt.Sprintf("%d critical and %d high-severity issue(s) need attention now", c, h)
	case h > 0:
		lead = fmt.Sprintf("%d high-severity issue(s) need attention this month", h)
	default:
		lead = "no critical or high-severity issues — the open items are lower-priority hygiene"
	}
	return fmt.Sprintf("%d open finding(s) across %d Sentinel tool(s); %s. The priorities below are ordered so the highest-risk items are fixed first.", rep.OpenTotal, len(rep.Sources), lead)
}

var funcs = template.FuncMap{
	"sevClass": func(s core.Severity) string {
		switch s {
		case core.SeverityCritical:
			return "crit"
		case core.SeverityHigh:
			return "high"
		case core.SeverityMedium:
			return "med"
		case core.SeverityLow:
			return "low"
		default:
			return "info"
		}
	},
	"product": func(module string) string { return posture.ProductName(module) },
}

var tmpl = template.Must(template.New("report").Funcs(funcs).Parse(pageHTML))
