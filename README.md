# Posture Report

**One security-posture score across every Hexward tool you run — self-hosted, read-only, prints to PDF.**

You run several Hexward tools (TLS, attack surface, canaries, CVEs, firewall audit, logs, DMARC, M365/Workspace posture). Each is great on its own — but a manager wants one answer: *how are we doing, and what do we fix first?* Posture Report reads the open findings from every tool's database (read-only), folds them into a single posture score, and produces one report with two views: an **Executive** page (score, plain-language summary, the handful of things to fix first) and a **Technical** page (every open finding with its remediation). It runs no scans and changes nothing.

```
posturereport -dir /var/lib/hexward       # auto-discover the tools' databases
posturereport -dbs certlight.db,asm.db    # or list them
posturereport -out posture.html           # write the report once (cron/monthly)
```

Dashboard on 127.0.0.1:8432; the full report is at /report (Print → PDF).

## Editions
Free reads 3 tools · Pro 20 · Team unlimited. Pro/Team: **whop.com/nizar-tuanku/posture-report?utm_source=github** — part of the **Hexward Essentials (SMB)** bundle.

**Whop sells paid licences only.** Free: github.com/nizartuanku/posture-report — this repository is the free edition, Apache-2.0, no time limit; nothing on Whop is free, so try it here first.

Free edition is Apache-2.0. Part of the Hexward line: **whop.com/nizar-tuanku**

## AI Assist (optional)

Posture Report can explain a finding in plain language with a small language model that runs on
your own hardware. It is off by default. Turn it on by starting a
[hexward-ai](https://github.com/nizartuanku/hexward-ai) sidecar and pointing Posture Report at it:

```sh
posturereport -ai-assist-url http://127.0.0.1:8435
```

The dashboard gains a **Fix first** list (the same top priorities as the report), and each item gets an **✨ Explain** button. The model writes what the finding means and
what to verify before you act. It also gets a fixed disclaimer.

- **The engine still decides.** The model receives one finding after the Hexward tool that owns it has produced it; Posture Report only reads those databases and never writes to them.
  It cannot add, remove, re-score or close a finding. If the sidecar is off, slow or broken,
  the button shows a short note and nothing else changes.
- **What leaves the process.** One finding: its check, title, target, severity, status,
  remediation and a sanitised copy of its evidence. Keys that look like secrets (password,
  token, secret, private, credential, cookie, session, signature and similar) are dropped
  first. Nothing goes to the internet. The sidecar runs where you run it.
- **Editions.** The free edition works with a sidecar on the same host. That is the `lab`
  profile, SmolLM3-3B. Pro and Team can also use one dedicated AI host for several products,
  or your own OpenAI-compatible endpoint, through `-ai-assist-key-file`. The recommended
  profile there is `smb` (Phi-4-mini-instruct). Enterprise uses Qwen3 or your own endpoint.
- **Language.** English is the supported language in this release. `-ai-assist-lang id`
  (Bahasa Indonesia) remains as an unsupported preview. More languages will be added based on
  demand.
- **Honest limit.** Small local models sometimes add general background that is not in the
  evidence. For example, they may name a well-known attack, and that background can be wrong.
  Treat the explanation as a starting point. The finding, its evidence and its fix text remain
  the record, which is why every explanation carries the "verify against raw findings" line.
- **Speed.** On a CPU-only machine an explanation takes about 15–50 seconds, depending on the
  model. Measurements are in hexward-ai's `docs/TIERS.md`.

Environment equivalents: `POSTURE_REPORT_AI_ASSIST_URL`, `POSTURE_REPORT_AI_ASSIST_KEY_FILE`,
`POSTURE_REPORT_AI_ASSIST_LANG`, `POSTURE_REPORT_AI_ASSIST_NO_THINKING=1`.

