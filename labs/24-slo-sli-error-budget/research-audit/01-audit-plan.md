# 01 — Audit Plan

## Target Lab
`labs/24-slo-sli-error-budget`

## Objective & Scope
PIPELINE OVERRIDE:
- Audit research files (`labs/24-slo-sli-error-budget/research/`) ONLY.
- Do NOT audit implementation, tests, or source code in this stage.
- Do NOT modify research files.
- Deliver audit output strictly to `labs/24-slo-sli-error-budget/research-audit/`.

## Files Reviewed
- `labs/24-slo-sli-error-budget/research/01-plan.md`
- `labs/24-slo-sli-error-budget/research/02-sources.md`
- `labs/24-slo-sli-error-budget/research/03-evidence.md`
- `labs/24-slo-sli-error-budget/research/04-contradictions.md`
- `labs/24-slo-sli-error-budget/research/05-report.md`
- `labs/24-slo-sli-error-budget/research/06-open-questions.md`

## Major Claims To Verify
1. Canonical definitions of SLI, SLO, SLA, Error Budget (Google SRE Book Ch. 4, Ch. 3).
2. Availability downtime calculations for 99%, 99.9%, 99.99%, 99.999% (Google SRE Book Appendix A Table 1-1).
3. 100% reliability is wrong target & cost non-linear (~100x per nine) (Google SRE Workbook Ch. 2, Book Ch. 3).
4. Error Budget = 1 - SLO formula and 10M request example (Google SRE Book Ch. 3).
5. Burn rate alerting recommendations (14.4x/1h, 6x/6h, 1x/3d) and short window = 1/12 long window (Google SRE Workbook Ch. 5).
6. SLI best practices (ratio good/total, percentile over average, user-facing vs CPU/RAM, endpoint-specific SLOs) (Google SRE Book Ch. 4, SRE Workbook Ch. 5).
7. Claimed statistic: "70% of outages caused by changes" (Google SRE Workbook Appendix B).
8. Datadog Error Budget Remaining formula: `100 * (current_status - target) / (100 - target)` (Datadog SLO Docs).

## Audit Strategy
1. **Source Audit (`02-source-audit.md`)**: Inspect all 9 cited sources in `02-sources.md` using live WebFetch / knowledge verification to confirm accessibility, publisher accuracy, relevance, and tiering.
2. **Claim Audit (`03-claim-audit.md`)**: Extract major claims from `03-evidence.md` and `05-report.md`, mapping evidence to source text and classifying severity/support level.
3. **Contradictions Audit (`04-contradictions.md`)**: Audit internal & external contradictions identified in `04-contradictions.md` (e.g., 30-day vs 28-day window, monthly downtime basis, burn rate thresholds, 70% change outage statistic).
4. **Code Audit (`05-code-audit.md`)**: Formally state N/A under Pipeline Override (Research Audit only).
5. **Gap Analysis (`06-gaps.md`)**: Identify missing/weak sources, overgeneralizations, vendors-specific limitations, and unverified stats.
6. **Final Verdict (`07-verdict.md`)**: Synthesize findings into quality gates and determine verdict (`APPROVED`, `APPROVED_WITH_WARNINGS`, `NEEDS_REVISION`, `REJECTED`).

## Primary Risks
- Mono-vendor bias (over-reliance on Google SRE documentation).
- Unverified empirical claims presented as universal constants (e.g., "70% outages caused by changes", "100x cost per nine").
- Mathematical ambiguity regarding monthly downtime calculations (30 days vs 30.44 days).
