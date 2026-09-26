# Revision Plan

Target Lab: labs/24-slo-sli-error-budget
Previous Audit Status: APPROVED_WITH_WARNINGS (research-only stage, 2026-09-26)

## Blocking Issues

None. Audit confirms: no fabricated sources, quotes, arithmetic, or critical misinformation.

## Non-Blocking Issues

1. MEDIUM — Single-corpus dependency: 8/8 sources Google SRE. Lab must frame as Google SRE framework, not universal standard. (verdict §Non-Blocking 1, gaps Gap 1)
2. MEDIUM — "70% outages from change" quoted verbatim but source gives no methodology. Must qualify if used. (claim 17, contradiction 4, gaps Gap 2)
3. LOW — CPU/RAM as non-SLO is valid inference, not direct quote. Frame as principle, not citation. (claim 19, gaps Gap 6)
4. LOW — Month-length arithmetic: lab draft 30.44-day vs Appendix A 30-day (6-min gap at 99%/month; 3.8s at 99.99%/month). Standardize + state assumption. (contradictions 1-2)
5. LOW — 28-day (Workbook Ch.2) vs 30-day (Workbook Ch.5 Table 5-4) window inconsistency inside Google corpus. Document trade-off. (contradiction 3, gaps Gap 4)
6. LOW — Tooling currency: sources 2017-2018, no OpenSLO/Sloth/OpenTelemetry coverage. Note only, core concepts stable. (gaps Gap 8)
7. LOW — Four Golden Signals single-source; batch/ETL/Kafka SLIs and regulatory scope uncovered but self-documented. No action beyond noting. (gaps Gap 3,5,7)

## Files To Modify

- None in `research/runs/2026-09-26-slo-sli-error-budget/` (frozen run; already honest, confidence ratings verified correct).
- New overlay only in `research-revision/` (this directory):
  - `01-revision-plan.md` (this file)
  - `02-changes-made.md` (canonical qualified statements lab material MUST use)
  - `03-revision-result.md` (counts + validation + readiness)
- No code, no tests, no README — none exist at this pipeline stage (code-audit: NOT APPLICABLE).

## Verification Plan

- source verification: reuse audit 02-source-audit (8/8 reachable, quotes verbatim). No re-fetch; no new citations invented.
- research consistency: each qualified statement cites exact run file + audit location it supersedes.
- arithmetic: reuse audit 05-code-audit Calculations 1-5 (all PASS). Standardized month assumption stated with both formulas.
- demo/build/tests: N/A (research-only pipeline override, no implementation).
