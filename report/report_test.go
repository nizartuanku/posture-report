package report

import (
	"strings"
	"testing"
	"time"

	"github.com/nizartuanku/posture-report/core"
	"github.com/nizartuanku/posture-report/posture"
)

func TestHTMLContainsKeyContent(t *testing.T) {
	items := []posture.Item{
		{Product: "RuleHawk", Module: "rulehawk", Findings: []core.Finding{
			{Module: "rulehawk", Check: "fw.permissive", Title: "any->any allow rule",
				Severity: core.SeverityCritical, Remediation: "tighten the rule", Target: "fw1", Status: core.StatusOpen},
		}},
		{Product: "TenantWatch", Module: "tenantwatch", Findings: []core.Finding{
			{Module: "tenantwatch", Check: "tenant.manual-review", Title: "Manual review: sign-in activity",
				Severity: core.SeverityInfo, Remediation: "review", Target: "m365:x", Status: core.StatusOpen},
		}},
	}
	rep := posture.Build("Contoso Nusantara", items, time.Unix(1_700_000_000, 0))
	html := HTML(rep)

	for _, want := range []string{
		"Contoso Nusantara", "Executive", "Technical",
		"any-&gt;any allow rule", // html-escaped
		"tighten the rule", "RuleHawk",
		"Blind spots", "Manual review: sign-in activity",
		"Fix these first", "class=\"sev crit\"",
	} {
		if !strings.Contains(html, want) {
			t.Errorf("report HTML missing %q", want)
		}
	}
	// The report must be a full standalone document.
	if !strings.HasPrefix(html, "<!DOCTYPE html>") || !strings.Contains(html, "</html>") {
		t.Error("report is not a standalone HTML document")
	}
}

func TestHTMLEmptyState(t *testing.T) {
	rep := posture.Build("Clean Co", []posture.Item{{Product: "CertWatch", Module: "certwatch"}}, time.Unix(0, 0))
	html := HTML(rep)
	if !strings.Contains(html, "No open findings") {
		t.Error("clean report should say there are no open findings")
	}
	if !strings.Contains(html, ">100<") {
		t.Error("clean report should show score 100")
	}
}
