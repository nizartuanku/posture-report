package report

// pageHTML is the self-contained two-model report. No external assets, prints
// cleanly to PDF, theme-aware. Driven by the `view` struct in report.go.
const pageHTML = `<!DOCTYPE html>
<html lang="en"><head><meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{{.Rep.Title}} — Posture Report</title>
<style>
:root{--blue:#4A86D6;--teal:#2FE0C0;--ink:#0f1720;--ink2:#3a4653;--muted:#6b7885;--line:#e3e8ee;--bg:#f5f7fa;--card:#fff;--shadow:0 1px 3px rgba(16,24,40,.06);
--crit:#c1121f;--high:#e8590c;--med:#d9a406;--low:#2f8f5b;--info:#4A86D6;
--crit-bg:#fdecee;--high-bg:#fdefe6;--med-bg:#fbf6e3;--low-bg:#e9f6ee;--info-bg:#eaf1fb;
--font:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,Helvetica,Arial,sans-serif;--mono:ui-monospace,"SF Mono",Menlo,Consolas,monospace;}
*{box-sizing:border-box}body{margin:0;background:var(--bg);color:var(--ink);font-family:var(--font);font-size:15px;line-height:1.55}
.wrap{max-width:900px;margin:0 auto;padding:0 20px 80px}
.topbar{position:sticky;top:0;z-index:10;background:rgba(255,255,255,.9);backdrop-filter:blur(8px);border-bottom:1px solid var(--line)}
.tin{max-width:900px;margin:0 auto;padding:10px 20px;display:flex;align-items:center;gap:12px}
.brand{display:flex;align-items:center;gap:9px;font-weight:700;letter-spacing:-.2px}.brand small{font-weight:500;color:var(--muted)}
.toggle{margin-left:auto;display:inline-flex;background:#eef2f7;border-radius:10px;padding:3px}
.toggle button{border:0;background:transparent;padding:7px 15px;border-radius:8px;font:inherit;font-weight:600;font-size:13.5px;color:var(--ink2);cursor:pointer}
.toggle button.active{background:#fff;color:var(--blue);box-shadow:var(--shadow)}
.btn{border:1px solid var(--line);background:#fff;border-radius:8px;padding:7px 12px;font:inherit;font-size:13px;font-weight:600;color:var(--ink2);cursor:pointer}
.rhead{background:var(--card);border:1px solid var(--line);border-radius:16px;box-shadow:var(--shadow);padding:24px 26px;margin-top:16px}
.rhead h1{margin:0 0 3px;font-size:23px;letter-spacing:-.4px}.rhead .sub{color:var(--muted);font-size:14px}
.meta{display:grid;grid-template-columns:repeat(auto-fit,minmax(140px,1fr));gap:12px 20px;border-top:1px solid var(--line);margin-top:16px;padding-top:14px}
.meta .k{font-size:11px;text-transform:uppercase;letter-spacing:.6px;color:var(--muted);font-weight:700}.meta .v{font-size:14px;font-weight:600}
section{margin-top:26px}.s-title{display:flex;align-items:center;gap:10px;font-size:12px;font-weight:800;letter-spacing:1.1px;text-transform:uppercase;color:var(--muted);margin:0 0 12px}.s-title::after{content:"";flex:1;height:1px;background:var(--line)}
.card{background:var(--card);border:1px solid var(--line);border-radius:14px;box-shadow:var(--shadow);padding:22px 24px}
.exec-top{display:grid;grid-template-columns:190px 1fr;gap:24px;align-items:center}
.gauge{text-align:center}.gauge .num{font-size:54px;font-weight:800;line-height:1;letter-spacing:-1.5px}.gauge .den{font-size:16px;color:var(--muted);font-weight:600}
.gauge .rating{margin-top:8px;display:inline-block;font-weight:700;font-size:13px;padding:4px 12px;border-radius:20px}
.num.crit,.rating.crit{color:var(--crit)}.rating.crit{background:var(--crit-bg)}.num.high,.rating.high{color:var(--high)}.rating.high{background:var(--high-bg)}
.num.med,.rating.med{color:var(--med)}.rating.med{background:var(--med-bg)}.num.low,.rating.low{color:var(--low)}.rating.low{background:var(--low-bg)}
.riskrow{display:grid;grid-template-columns:78px 1fr 34px;align-items:center;gap:12px;font-size:13px;margin:9px 0}.riskrow .lab{font-weight:700}
.track{height:10px;background:#eef2f7;border-radius:6px;overflow:hidden}.fill{height:100%;border-radius:6px}.cnt{text-align:right;font-weight:700}
.crit .fill,.fill.crit{background:var(--crit)}.high .fill,.fill.high{background:var(--high)}.med .fill,.fill.med{background:var(--med)}.low .fill,.fill.low{background:var(--low)}
.lab.crit{color:var(--crit)}.lab.high{color:var(--high)}.lab.med{color:var(--med)}.lab.low{color:var(--low)}
.prio{display:grid;gap:12px}.pitem{border:1px solid var(--line);border-radius:12px;padding:14px 16px}
.pitem .top{display:flex;gap:10px;align-items:flex-start}.sev{font-size:10.5px;font-weight:800;letter-spacing:.5px;padding:3px 8px;border-radius:6px;text-transform:uppercase;white-space:nowrap}
.sev.crit{background:var(--crit-bg);color:var(--crit)}.sev.high{background:var(--high-bg);color:var(--high)}.sev.med{background:var(--med-bg);color:var(--med)}.sev.low{background:var(--low-bg);color:var(--low)}.sev.info{background:var(--info-bg);color:var(--info)}
.pitem .ttl{flex:1;font-weight:600}.pitem .prod{font-size:11px;font-weight:700;color:var(--blue)}
.pitem .rem{font-size:13.5px;color:var(--ink2);margin-top:7px}.pitem .tgt{font-family:var(--mono);font-size:11.5px;color:var(--muted);margin-top:5px}
table.tt{width:100%;border-collapse:collapse;font-size:13.5px}table.tt th,table.tt td{text-align:left;padding:9px 10px;border-bottom:1px solid var(--line)}
table.tt th{font-size:10.5px;text-transform:uppercase;letter-spacing:.6px;color:var(--muted)}table.tt td.n{text-align:right;font-variant-numeric:tabular-nums;font-weight:700}
.pill{display:inline-block;min-width:22px;text-align:center;font-size:12px;font-weight:700;padding:1px 7px;border-radius:20px;margin-left:4px}
.pill.crit{background:var(--crit-bg);color:var(--crit)}.pill.high{background:var(--high-bg);color:var(--high)}.pill.med{background:var(--med-bg);color:var(--med)}.pill.low{background:var(--low-bg);color:var(--low)}
.note{background:var(--info-bg);border:1px solid #cfe0f7;border-radius:10px;padding:11px 14px;font-size:13px;color:#28477a;margin-top:12px}
.finding{background:var(--card);border:1px solid var(--line);border-left:4px solid var(--line);border-radius:12px;box-shadow:var(--shadow);margin-bottom:12px;padding:14px 16px}
.finding.crit{border-left-color:var(--crit)}.finding.high{border-left-color:var(--high)}.finding.med{border-left-color:var(--med)}.finding.low{border-left-color:var(--low)}.finding.info{border-left-color:var(--info)}
.finding .top{display:flex;gap:10px;align-items:flex-start}.finding .rem{font-size:13.5px;color:var(--ink2);margin-top:7px}.finding .tgt{font-family:var(--mono);font-size:11.5px;color:var(--muted);margin-top:5px}
.empty{text-align:center;color:var(--muted);padding:34px;border:1px dashed var(--line);border-radius:12px}
.foot{margin-top:32px;text-align:center;color:var(--muted);font-size:12px;line-height:1.7}
.hide{display:none!important}
@media print{.topbar{display:none!important}.view{display:block!important}.card,.rhead,.finding{box-shadow:none;break-inside:avoid}section{break-inside:avoid}.wrap{padding:0}body{background:#fff}}
@media(prefers-color-scheme:dark){:root:not([data-theme=light]){--bg:#0d1117;--card:#161b22;--ink:#e6edf3;--ink2:#b7c2cd;--muted:#8b98a5;--line:#26303b}
:root:not([data-theme=light]) .topbar{background:rgba(13,17,23,.9)}:root:not([data-theme=light]) .toggle{background:#21262d}:root:not([data-theme=light]) .toggle button.active{background:#0d1117}
:root:not([data-theme=light]) .track,:root:not([data-theme=light]) .btn{background:#21262d}}
</style></head><body>
<div class="topbar"><div class="tin">
  <span class="brand"><svg width="22" height="24" viewBox="0 0 22 24"><path d="M11 1l9 5.2v10.6L11 23l-9-5.2V6.2z" fill="none" stroke="#4A86D6" stroke-width="1.6"/><path d="M11 6.5l4.3 2.5v5L11 16.5 6.7 14V9z" fill="#2FE0C0"/></svg>Posture Report <small>· Sentinel</small></span>
  <div class="toggle"><button id="t-exec" class="active" onclick="show('exec')">Executive</button><button id="t-tech" onclick="show('tech')">Technical</button></div>
  <button class="btn" onclick="window.print()">Print / PDF</button>
</div></div>
<div class="wrap">
  <div class="rhead">
    <h1>{{.Rep.Title}}</h1>
    <div class="sub">Combined security posture across {{.ToolCount}} Sentinel tool(s)</div>
    <div class="meta">
      <div><div class="k">Generated</div><div class="v">{{.Generated}}</div></div>
      <div><div class="k">Open findings</div><div class="v">{{.Rep.OpenTotal}}</div></div>
      <div><div class="k">Tools</div><div class="v">{{.ToolCount}}</div></div>
      <div><div class="k">Manual review</div><div class="v">{{len .Rep.ManualReviews}}</div></div>
    </div>
  </div>

  <div id="view-exec" class="view">
    <section><div class="s-title">Posture</div><div class="card">
      <div class="exec-top">
        <div class="gauge"><div><span class="num {{.ScoreClass}}">{{.Rep.Score}}</span><span class="den">/100</span></div><div class="rating {{.ScoreClass}}">{{.Rep.Rating}}</div></div>
        <div><p style="margin:0">{{.Summary}}</p></div>
      </div>
    </div></section>

    <section><div class="s-title">Risk distribution</div><div class="card">
      {{range .Bars}}<div class="riskrow {{.Class}}"><span class="lab {{.Class}}">{{.Label}}</span><div class="track"><div class="fill {{.Class}}" style="width:{{.Pct}}%"></div></div><span class="cnt">{{.Count}}</span></div>{{end}}
    </div></section>

    {{if .Rep.TopPriorities}}
    <section><div class="s-title">Fix these first</div><div class="prio">
      {{range .Rep.TopPriorities}}<div class="pitem"><div class="top"><span class="sev {{sevClass .Severity}}">{{.Severity}}</span><span class="ttl">{{.Title}}</span></div>
        <div class="prod">{{product .Module}}</div><div class="rem">→ {{.Remediation}}</div><div class="tgt">{{.Target}}</div></div>{{end}}
    </div></section>
    {{end}}

    {{if .Rep.ManualReviews}}
    <section><div class="s-title">Blind spots — review by hand</div><div class="card">
      <div style="font-size:13.5px;color:var(--ink2)">{{len .Rep.ManualReviews}} area(s) could not be assessed automatically with the access granted. They are listed so the gap is visible — silence is never treated as "all clear".</div>
      <div class="note">{{range .Rep.ManualReviews}}• {{product .Module}}: {{.Title}}<br>{{end}}</div>
    </div></section>
    {{end}}

    <section><div class="s-title">By tool</div><div class="card">
      <table class="tt"><thead><tr><th>Tool</th><th>Open</th><th style="text-align:right">Breakdown</th></tr></thead><tbody>
      {{range .Rep.Sources}}<tr><td>{{.Product}}</td><td class="n">{{.Open}}</td><td style="text-align:right">
        {{if index .Counts "critical"}}<span class="pill crit">{{index .Counts "critical"}}</span>{{end}}
        {{if index .Counts "high"}}<span class="pill high">{{index .Counts "high"}}</span>{{end}}
        {{if index .Counts "medium"}}<span class="pill med">{{index .Counts "medium"}}</span>{{end}}
        {{if index .Counts "low"}}<span class="pill low">{{index .Counts "low"}}</span>{{end}}
        {{if not .Open}}<span style="color:var(--muted)">clean</span>{{end}}</td></tr>{{end}}
      </tbody></table>
    </div></section>
  </div>

  <div id="view-tech" class="view hide">
    <section><div class="s-title">All open findings ({{.Rep.OpenTotal}})</div>
      {{if .Rep.AllFindings}}{{range .Rep.AllFindings}}
      <div class="finding {{sevClass .Severity}}"><div class="top"><span class="sev {{sevClass .Severity}}">{{.Severity}}</span><span class="ttl" style="flex:1;font-weight:600">{{.Title}}</span><span class="prod" style="font-size:11px;font-weight:700;color:var(--blue)">{{product .Module}}</span></div>
        <div class="rem">→ {{.Remediation}}</div><div class="tgt">{{.Check}} · {{.Target}}</div></div>
      {{end}}{{else}}<div class="empty">No open findings across the connected tools.</div>{{end}}
    </section>
    <section><div class="s-title">Methodology</div><div class="card">
      <div style="font-size:13.5px;color:var(--ink2)">Posture Report reads the open findings from each connected Sentinel tool's local database, read-only, and combines them. It runs no scans of its own and changes nothing. The score weights findings by severity (critical &gt; high &gt; medium &gt; low); info-level "manual review" items are shown as blind spots but never scored against.</div>
    </div></section>
  </div>

  <div class="foot">Posture Report · Sentinel — self-hosted security tooling · generated locally, no data leaves your network</div>
</div>
<script>function show(w){var e=w=='exec';document.getElementById('view-exec').classList.toggle('hide',!e);document.getElementById('view-tech').classList.toggle('hide',e);document.getElementById('t-exec').classList.toggle('active',e);document.getElementById('t-tech').classList.toggle('active',!e);window.scrollTo({top:0,behavior:'smooth'})}</script>
</body></html>`
