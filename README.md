# Posture Report

**One security-posture score across every Hexward tool you run — self-hosted, read-only, prints to PDF.**

You run several Hexward tools (TLS, attack surface, canaries, CVEs, firewall audit, logs, DMARC, M365/Workspace posture). Each is great on its own — but a manager wants one answer: *how are we doing, and what do we fix first?* Posture Report reads the open findings from every tool's database (read-only), folds them into a single posture score, and produces one report with two views: an **Executive** page (score, plain-language summary, the handful of things to fix first) and a **Technical** page (every open finding with its remediation). It runs no scans and changes nothing.

```
posturereport -dir /var/lib/sentinel      # auto-discover the tools' databases
posturereport -dbs certwatch.db,asm.db    # or list them
posturereport -out posture.html           # write the report once (cron/monthly)
```

Dashboard on 127.0.0.1:8432; the full report is at /report (Print → PDF).

## Editions
Free reads 3 tools · Pro 20 · Team unlimited. Pro/Team: **whop.com/nizar-tuanku/posture-report?utm_source=github** — part of the **Hexward Essentials (SMB)** bundle.

Free edition is Apache-2.0. Part of the Hexward line: **whop.com/nizar-tuanku**
