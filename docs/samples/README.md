# Sample run

This is a real Posture Report output, not a mockup. It reads a genuine
CertLight SQLite database — produced by scanning five public TLS test hosts
for real — and computes the score the same way it would for any Hexward
database handed to it. Nothing in [`sample-report.html`](sample-report.html)
or [`sample-summary.json`](sample-summary.json) was hand-edited.

**Source:** a lab instance of Posture Report (free edition) pointed at a
`certlight.db` from a lab run of CertLight 0.1.3 against
`expired.badssl.com`, `self-signed.badssl.com`, `incomplete-chain.badssl.com`,
`wrong.host.badssl.com` and `tls-v1-1.badssl.com` — the same real scan
documented in [CertLight's own `docs/samples/`](https://github.com/nizartuanku/certlight/tree/main/docs/samples).
This sample shows exactly how Posture Report aggregates and scores findings
it did not itself produce — CertLight's engine is the only source of the
findings; Posture Report only reads, tallies and reports.

## What Posture Report computed

```
GET /api/summary
{
    "title": "Hexward Sample Posture",
    "tier": "free",
    "score": 15,
    "rating": "At Risk",
    "open_total": 10,
    "tools": 1,
    "sources": [{"product": "CertLight", "open": 10}],
    "counts": {"critical": 1, "high": 3, "medium": 6}
}
```

The score is not a guess — it is 100 minus a weighted penalty, capped at
100, read straight from `posture/posture.go`: critical = 25, high = 12,
medium = 4, low = 1. This run's 10 open findings (1 critical, 3 high, 6
medium, 0 low) add up to exactly `1×25 + 3×12 + 6×4 = 85`, so
`100 − 85 = 15` — **"At Risk"** (the "Needs Attention" / "At Risk" boundary
is 50; this run is well below it).

`sample-report.html` is the full report Posture Report serves at `/report`
— the same priority ordering a real client-facing report would use, with
the most severe findings first (`TopPriorities`, capped at 6) and the
one-paragraph summary the report opens with: *"10 open finding(s) across 1
Hexward tool(s); 1 critical and 3 high-severity issue(s) need attention
now."*

## What this sample does not show

Only one source database (`certlight.db`) was supplied, so `tools: 1` and
there is no cross-product correlation to show — a real deployment would
typically point `-dbs` at several Hexward products' databases at once
(`certlight.db,asm.db,auditlight-jobs/`) and this report would combine all
of them into one score, exactly as the main README's quick-start example
shows. AI Assist was not exercised in this sample (no `-ai-assist-url`);
the main README's "AI Assist" section covers what that adds.
