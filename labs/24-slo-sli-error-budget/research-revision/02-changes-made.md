# Changes Made

Research-only patch. No source files rewritten; no code/tests added (pipeline override). All qualifying statements below are the canonical text lab material SHOULD adopt per audit §Required Revisions.

## Revision 1

Audit Issue: MEDIUM — single-corpus dependency; lab must not present Google SRE content as universal standard.

Action: Authoritative framing statement added to changes record.
- "This lab teaches the Google SRE framework for SLI/SLO/error budgets — originated and described primarily in the Google SRE Book (2017) and SRE Workbook (2018). All 8 sources verified reachable and quotes verbatim per audit 02-source-audit. Concepts are widely adopted but cross-vendor/independent-source validation was not completed in this research run (AWS/Azure/CNCF 404/redirect, see 05-report.md §Limitations, Gap 1). Present as Google SRE methodology, not universal industry standard."

Confidence: HIGH (source-verified); framing: REQUIRED.

Status: RESOLVED

---

## Revision 2

Audit Issue: MEDIUM — "70% of outages from change" quoted verbatim but source provides no methodology; must not be cited as universal benchmark.

Action: Mandatory qualification statement.
- "Changes are a major source of instability, representing roughly 70% of our outages" (SRE Workbook Appendix B §Background, 2018). Google-internal; **no methodology published**, no time frame, no outage definition given. Do NOT cite as universal industry fact. (claim 17; contradiction 4; Gap 2)"

Supporting distinction (contradiction 5):
- Ch.2 states categorically "The number one source of outages is change" (no %). Appendix B quantifies "roughly 70%" (no methodology). Both from same corpus; precision differs. Neither externally sourced. Present as Google SRE guidance, not independent statistic.

Status: RESOLVED

---

## Revision 3

Audit Issue: MEDIUM/Low — CPU/RAM not SLO claim is valid inference, not explicit citation.

Action: Framing rule.
- Do not cite "CPU and RAM are NOT SLOs" as a verbatim SRE Book statement. Frame as best-practice inference: "SRE methodology recommends measuring user-facing outcomes (e.g., latency, availability) rather than infrastructure metrics (CPU, RAM) as SLO inputs — because infra metrics are leading indicators/causes, not user-observable symptoms. Source expresses 'ideally, the SLI directly measures a service level of interest' (claim 19 / Gap 6)."

Status: RESOLVED

---

## Revision 4

Audit Issue: LOW — month-length arithmetic inconsistency (30.44-day vs 30-day → 6-min gap at 99%/month, 3.8s at 99.99%/month).

Action: Standardize + disclose. Lab material MUST pick exactly one and label it explicitly:

Option A (matches Appendix A exactly): "All downtime-per-month values use a 30-day month (43,200 minutes), matching SRE Book Appendix A Table 1-1. This yields 7.2h/month at 99%, 43.2m/month at 99.9%, 4.32m/month at 99.99%. (audit 02-source-audit Source 5; 04-contradictions §1-2)."

Option B (calendar average): "Downtime-per-month here uses a 30.44-day month (365÷12 ≈ 30.44, 43,833.6 min) for tighter alignment with calendar months: 99%≈7h18m/mo, 99.9%≈43m50s/mo, 99.99%≈4m23s/mo. This is ~1.4% larger than the 30-day figures in Appendix A. Difference does not affect engineering decisions; stated for consistency only. (audit 03-claim-audit Claim 5; 04-contradictions §1-2)."

Either option is acceptable; the ONLY forbidden state is mixing assumptions or presenting either as "the" correct value without disclosure.

Arithmetic cross-check (audit 05-calculations):
- 99%: 0.01 × 43,200 min (30-day) = 432 min = 7.2 h ✓ ; vs 0.01 × 43,833.6 (30.44-day) = 438.34 min ≈ 7.31 h (7h19m) ✓
- 99.99%: 0.0001 × 43,200 = 4.32 min = 4m19.2s ✓ ; vs 0.0001 × 43,833.6 = 4.38 min ≈ 4m23s ✓

Status: RESOLVED

---

## Revision 5

Audit Issue: LOW — 28-day (Workbook Ch.2) vs 30-day (Workbook Ch.5 Table 5-4) rolling window inconsistency inside the Google corpus itself.

Action: Documentation requirement for lab material.
- "SRE Workbook Ch.2 §Choosing an Appropriate Time Window recommends a 28-day (4-week) rolling window because 'we recommend defining this period as an integral number of weeks so it always contains the same number of weekends' (avoiding weekend count bias — some months have 4 weekends, some 5). This is the Workbook's general-purpose default. The Workbook Ch.5 burn-rate examples (Table 5-4) use 30-day windows for illustration only. These are not contradictory: Ch.2 is normative guidance on window shape (integer weeks), Ch.5 Table 5-4 is an illustrative example using a round month. Lab material MUST state which window it adopts and why. (contradiction 3; Gap 4)."

Cross-check: burn-rate times in Table 5-4 are expressed in days (e.g., "30 days", "15 days", "3 days"), so the 30-day basis is explicit in those examples.

Status: RESOLVED

---

## Revision 6

Audit Issue: LOW — tooling currency (2017-2018); no OpenSLO/Sloth/OpenTelemetry coverage.

Action: Optional note (not blocking; core concepts stable).
- "Sources predate OpenTelemetry (2019) and modern SLO tooling. SRE Book Ch.10 uses Borgmon; contemporary equivalent is Prometheus/PromQL (acknowledged in Ch.10 footnote). OpenSLO spec / Sloth / pyrra exist as modern tooling equivalents but were out of research scope. (Gap 8)."

Status: NOT APPLICABLE TO VERIFICATION (documentation note only; no source fetched).

---

## Revision 7

Audit Issue: LOW — Four Golden Signals single-source; batch/ETL and regulatory scope uncovered (self-documented).

Action: None (honest self-assessment already in run files). No changes.
- Four Golden Signals: single Google source; industry adoption noted but not independently cross-checked (Evidence 10, claim 10 notes "NOT VERIFIED").
- Batch/ETL/Kafka SLOs and Telkom/ISO/ITIL standards: open questions §2-3,7; out of scope for focused SRE-Framework lab.

Status: RESOLVED (no change needed; gaps accurately self-reported).

---

## Changes Applied To (files modified)

None — this is a research-only revision. The frozen run in `research/runs/2026-09-26-slo-sli-error-budget/` is left intact because its confidence ratings and self-documented limitations are already audit-correct. All qualifying material above lives in `research-revision/02-changes-made.md` as the canonical guidance lab material should follow once it is authored.

## Source Verification Performed

Reused audit 02-source-audit.md (8/8 sources PASS: URL reachable, title/publisher confirmed, quotes verbatim). No new sources fetched or invented — no fabricated citations introduced.

## Arithmetic Verification Performed

Reused audit 05-code-audit.md (Calculations 1-5 PASS: error budget, payment-webhook exercise, burn-rate table, multiwindow cross-check, availability table). Revision 4 above adds explicit dual-window formulas so no arithmetic ambiguity remains.

Status: RESOLVED
