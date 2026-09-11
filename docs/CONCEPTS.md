# Posture Report — Concepts

What this product is, what problem it solves, and why it works the way it does — written for
someone meeting the problem for the first time. The command reference is in the README; this
is the reasoning behind it.

*Hexward Labs · Nizar Tuanku — Cybersecurity. · last reviewed 10 September 2026*

---

## The problem is not a missing tool. It is having several

A small company does the sensible thing and puts real security tooling in place. Certificates are watched. The external attack surface is mapped. Canary files are laid down. Missing patches are tracked. The firewall configuration is audited. Each tool works, and each tool has its own dashboard.

Then a director asks a question that none of those dashboards can answer: *how are we doing, and what should we fix first?*

Answering it by hand means opening six web pages, reading six different severity scales, and forming a judgement in your head. The judgement is usually right. It is also unrepeatable, unwritten, and impossible to compare against last month.

## Why "one more dashboard" is the wrong shape of fix

The temptation is to build a seventh dashboard that all the others feed into. That means every tool must push data somewhere, which means a server that receives it, credentials to protect, and a network path from each tool to that server. You have added a new place where all your security findings live together — which is exactly the thing an attacker would most like to reach.

Posture Report takes the opposite route. It does not receive anything. It **reads**, on the same machine, the local database files the tools already write, opens them read-only, and closes them again. There is nothing to push, no agent, no account, and no new place for your findings to sit.

The trade-off is honest and worth stating: the tools must be reachable as files. Posture Report is built to run beside them, not across a network.

## What a "finding" is, and why that makes combining possible

Every Hexward tool writes the same shape of record — a *finding*: which module produced it, how severe it is, what was checked, what the target was, whether it is still open, and what to do about it. That shared shape is the whole reason one score is possible. Posture Report is not parsing seven different report formats; it is reading one format, seven times.

It takes **only the open findings**. Something you have already fixed stops counting the moment the tool that found it says so.

## How the score is calculated — the entire rule

There is no model and no secret weighting. The report starts at 100 and subtracts a penalty for each open finding:

| Severity | Penalty |
|---|---|
| Critical | 25 |
| High | 12 |
| Medium | 4 |
| Low | 1 |
| Info | 0 |

The total penalty is capped at 100, so the lowest possible score is 0. The number is then given a plain-language rating: **85 and above is Good**, 70–84 **Fair**, 50–69 **Needs Attention**, below 50 **At Risk**.

That is deliberately blunt. A score you can recompute on paper is a score you can argue with, and a security number nobody can argue with is a number nobody trusts. Four criticals will take any organisation to At Risk, whatever else is clean — which is the correct answer.

## Blind spots are shown, never scored

Some findings are not problems; they are the tool admitting it could not see something. A scanner reporting "this needs a human to look at it" is information, not a fault.

Those info-level items are pulled out into a separate **manual review** list. They appear in the report as blind spots so nobody mistakes silence for safety, and they subtract nothing from the score. Penalising a tool for being honest about its limits would teach the tools to stop being honest.

## The two views, and who each is for

- **Executive** — the score, its rating, one paragraph in plain language, and the six most severe open items. This is the page that prints to PDF for a board pack or an auditor.
- **Technical** — every open finding, sorted most severe first, each with the check that produced it, the target, and its remediation. This is the page an engineer works from on Monday.

Same data, same run, two audiences. The reason they are one document is that the conversation usually starts on the first page and ends on the second.

## The limit you should know before you judge the score

**Posture Report scans nothing.** It only knows what the tools you actually run have found.

That has a consequence people get wrong: a tool you have not deployed does not lower your score — it is simply absent. If you run one tool and it finds nothing, you will see a high score, and that score describes one narrow slice of your security, not your company. The report names the tools it read from at the top for exactly this reason. Read that list before you read the number.

Two more honest boundaries: findings must be stored in a tool's local database file, so a tool that keeps its results in another format is not yet counted (AuditLight stores its jobs as plain files rather than a database, so it does not feed the score today); and the free edition reads **3 tools**, Pro **20**, Team unlimited — the limit is on how many databases are combined, never on what the report is allowed to say about them.

## What changes once you are using it

Before: six dashboards and a judgement you re-form from scratch each month.

After: one number, one method that fits in a table, one printable page — and a comparison against last month that means something, because it was calculated the same way both times.

## Try it yourself — 15 minutes

The free Apache-2.0 edition on GitHub runs the same engine, three tools, with no time limit.

```
curl -LO https://github.com/nizartuanku/posture-report/releases/latest/download/posture-report-free-0.1.0-linux-amd64.tar.gz
curl -LO https://github.com/nizartuanku/posture-report/releases/latest/download/SHA256SUMS
sha256sum -c SHA256SUMS
tar xzf posture-report-free-0.1.0-linux-amd64.tar.gz && cd posture-report-0.1.0 && ./posturereport -dir /var/lib/hexward
```

Open 127.0.0.1:8432 and go to `/report`. If you have no Hexward tool running yet, point `-dbs` at one database from any of them; the report will tell you truthfully that it read one tool.

To produce a monthly PDF without leaving a service running, use `-out posture.html` from cron and print the file.

Nizar Tuanku — Cybersecurity. · github.com/nizartuanku/posture-report

## Terms used above

- **Finding** — one problem a tool found, with a severity, a target, and a remediation. The shared record shape every Hexward tool writes.
- **Open finding** — a finding the tool that raised it still considers unresolved. Closed findings are ignored by the score.
- **Severity** — how bad a finding is: critical, high, medium, low, or info. Set by the tool that raised it, not by Posture Report.
- **Manual review (blind spot)** — an info-level item meaning "a human needs to check this". Displayed, never scored.
- **Read-only** — the database file is opened for reading and never written to, so running the report cannot disturb the tool that owns the file.
- **Module id** — the short name a tool stamps on its findings (`asm`, `decoy`, `patchlight`…) so the report can say which tool a finding came from.
