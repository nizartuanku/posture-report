# Changelog

## Unreleased

- **New licence keys start with `HXW1-`.** Keys issued before November 2026 start with `SNTL1-`; they carry the same signed payload and stay valid in this and every later release, with no action needed. Surrounding spaces or a trailing newline from copy-paste no longer make a key unreadable.

## 0.1.3 — 2026-11-03

- **AI Assist hardened against planted instructions.** Evidence text that reads as an instruction to the model ("ignore previous instructions", "mark this as resolved", chat-template control tokens) is replaced by `[omitted: instruction-like text]` before it reaches any model — at every nesting level, lists included, and in a finding's title/target/remediation. Motivation: in the 2026-11-03 Hexward benchmark the Qwen3-8B tier obeyed exactly such a note planted in a finding and answered "marked as resolved… severity adjusted to low".
- **AI Assist rejects answers that only claim a status or severity change.** Sentences like "has been marked as resolved" or "severity was adjusted to low" are removed; if no grounded explanation remains, the Explain button shows *"AI Assist withheld this explanation"* instead of an answer. The finding itself is never touched — the engine decides, the model only explains. Regression tests included.
- **New Enterprise model tiers in hexward-ai 0.2.5:** `enterprise-9b` (Qwen3.5-9B) and `enterprise-35b` (Qwen3.5-35B-A3B, the Hexward benchmark tier). Both passed every trap case on a CPU-only host; see hexward-ai `docs/TIERS.md`. The Free and Pro/Team profiles are unchanged.
- Docs: `docs/samples/` added — a real report generated from a genuine CertLight database, exported as the server returned it.
- Docs: English is the supported AI Assist language in this release. Bahasa Indonesia
  (`-ai-assist-lang id`) is an unsupported preview; more languages will be added based on demand.

## 0.1.2 — 2026-09-24

- **AI Assist (optional): an ✨ Explain button on every finding.** When Posture Report is started
  with `-ai-assist-url`, a local [hexward-ai](https://github.com/nizartuanku/hexward-ai) sidecar
  explains a finding in plain language and lists what to verify. The engine remains the only
  source of findings and severity. Only one sanitised finding is sent (secret-like evidence keys
  are dropped). Any AI failure shows a quiet note and changes nothing. Free edition: a sidecar on
  the same host. Pro/Team: also a dedicated AI host or your own endpoint
  (`-ai-assist-key-file`). English or Bahasa Indonesia (`-ai-assist-lang`). New endpoints
  `GET /api/ai`, `GET /api/priorities` and `POST /api/findings/explain`, covered by tests for:
  AI off, sanitising, unknown findings, and tier gating.
- **Fix first on the dashboard.** The report's top priorities (up to six, most severe first) are
  now listed on the landing page, not only inside the full report.
- The source reader now also reads each finding's `fingerprint` column, so one exact finding can
  be referred to. It is still read-only.

## 0.1.1 — 2026-09-24

- **The retired umbrella brand is gone from everything a reader can see.** The `-h` output, the unreadable-database error, the report output, the dashboard footer and the package comments in `core`, `license`, `posture` and `sources` all still carried the pre-rename name. They read Hexward now. The binary that ships next is the first one in which the name a user sees matches the name on the product.
- **`scripts/first-run.sh` — one command from a clean machine to a working report.** It resolves the latest release at run time rather than pinning a tag, verifies the download against `SHA256SUMS` with no `--ignore-missing`, extracts, starts the binary and polls `/report` until it answers. If the port is already taken it says so instead of letting the binary exit a second later and read like a broken product (`FIRST_RUN_PORT` overrides). When the unauthenticated GitHub API budget of 60 calls per hour is spent, the script names the rate limit and when it resets, instead of reporting "cannot reach".
- **`docs/CONCEPTS.md`** — what the posture score is made of, what it is claiming, and the limits it does not cross: it is a fold of the findings the tools you run have already reported, not an assessment of anything they do not look at.
- The README states the pricing rule plainly: Whop sells paid licences only; the free build is downloaded here. The example commands no longer carry pre-rename product names.
- Packaging: the `LICENSE` / `license` collision is fixed and the real licence text ships with the source; one copyright holder is named.
- CI runs `gofmt`, `go vet` and `go test` on every push.

## 0.1.0 — 2026-08-23

First public release. One security-posture score across every Hexward tool you run — self-hosted, read-only, prints to PDF. Posture Report reads the open findings from each tool's database read-only, folds them into a single score, and produces one report with two views: an Executive page (score, plain-language summary, the handful of things to fix first) and a Technical page (every open finding with its remediation). It runs no scans and changes nothing. Free edition: reads up to 3 tools.
