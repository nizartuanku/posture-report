package posture

import (
	"testing"
	"time"

	"github.com/nizartuanku/posture-report/core"
)

func f(module, check, title string, sev core.Severity) core.Finding {
	return core.Finding{Module: module, Check: check, Title: title, Severity: sev,
		Remediation: "do the thing", Target: module + ":x", Status: core.StatusOpen}
}

func sample() []Item {
	return []Item{
		{Product: "TenantWatch", Module: "tenantwatch", Findings: []core.Finding{
			f("tenantwatch", "tenant.admin-mfa", "Admin without MFA", core.SeverityHigh),
			f("tenantwatch", "tenant.external-sharing", "Anyone-link sharing", core.SeverityMedium),
			f("tenantwatch", "tenant.manual-review", "Manual review: sign-in activity", core.SeverityInfo),
		}},
		{Product: "CertWatch", Module: "certwatch", Findings: []core.Finding{
			f("certwatch", "tls.expiry", "Cert expiring", core.SeverityLow),
		}},
		{Product: "RuleHawk", Module: "rulehawk", Findings: []core.Finding{
			f("rulehawk", "fw.permissive", "any->any allow", core.SeverityCritical),
		}},
	}
}

func TestBuildScoresAndSorts(t *testing.T) {
	rep := Build("Contoso", sample(), time.Unix(0, 0))
	// penalty = crit25 + high12 + med4 + low1 = 42 → score 58.
	if rep.Score != 58 {
		t.Errorf("score = %d, want 58", rep.Score)
	}
	if rep.Rating != "Needs Attention" {
		t.Errorf("rating = %q", rep.Rating)
	}
	if rep.OpenTotal != 4 {
		t.Errorf("open total = %d, want 4 (info excluded)", rep.OpenTotal)
	}
	if len(rep.ManualReviews) != 1 {
		t.Errorf("manual reviews = %d, want 1", len(rep.ManualReviews))
	}
	// Top priority must be the critical, most-severe first.
	if len(rep.TopPriorities) == 0 || rep.TopPriorities[0].Severity != core.SeverityCritical {
		t.Errorf("top priority should be critical, got %+v", rep.TopPriorities)
	}
	if rep.Counts["critical"] != 1 || rep.Counts["high"] != 1 {
		t.Errorf("counts wrong: %+v", rep.Counts)
	}
}

func TestCleanPosture(t *testing.T) {
	rep := Build("Clean", []Item{{Product: "CertWatch", Module: "certwatch"}}, time.Unix(0, 0))
	if rep.Score != 100 || rep.Rating != "Good" {
		t.Errorf("clean posture should score 100/Good, got %d/%s", rep.Score, rep.Rating)
	}
	if rep.OpenTotal != 0 {
		t.Errorf("open total = %d", rep.OpenTotal)
	}
}

func TestScoreFloorsAtZero(t *testing.T) {
	var many []core.Finding
	for i := 0; i < 10; i++ {
		many = append(many, f("x", "c", "crit", core.SeverityCritical))
	}
	rep := Build("Bad", []Item{{Product: "X", Module: "x", Findings: many}}, time.Unix(0, 0))
	if rep.Score != 0 || rep.Rating != "At Risk" {
		t.Errorf("10 crits should floor score to 0/At Risk, got %d/%s", rep.Score, rep.Rating)
	}
}

func TestProductName(t *testing.T) {
	if ProductName("tenantwatch") != "TenantWatch" || ProductName("asm") != "Attack Surface Monitor" {
		t.Error("known module names wrong")
	}
	if ProductName("newthing") != "Newthing" {
		t.Errorf("unknown module fallback wrong: %q", ProductName("newthing"))
	}
}
