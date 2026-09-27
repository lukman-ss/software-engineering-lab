# 06 — Research Gap Analysis

Comprehensive evaluation of research gaps, limitations, and unverified areas identified in `labs/24-slo-sli-error-budget/research/`.

---

## Gap 1: Single-Vendor / Mono-Vendor Source Dominance (Google SRE Ecosystem)

Type:
SCOPE_ERROR / VENDOR_BIAS

Severity:
MEDIUM

Location:
`research/02-sources.md`, `research/05-report.md` (Limitations)

Problem:
7 of the 9 cited sources (77.7%) are published by Google (Google SRE Book, SRE Workbook, Google Cloud Docs). While Google is the pioneer and primary authority on SRE, over-reliance on a single vendor risks treating Google-internal operational practices (e.g. 28-day window, 14.4x burn rate) as universal software engineering standards.

Required Revision:
None for initial approval, provided the research report explicitly highlights vendor limitations (which it does in `05-report.md:200-205`).

Can Be Approved Without Fix:
YES

---

## Gap 2: Empirical Basis of "70% Outages from Change"

Type:
UNVERIFIED_CLAIM

Severity:
MEDIUM

Location:
`research/03-evidence.md` (Evidence 16), `research/04-contradictions.md` (Contradiction 5)

Problem:
The claim that "70% of outages are caused by changes" comes from Google SRE Workbook Appendix B without an accompanying empirical study or external industry benchmark.

Required Revision:
Ensure the claim is framed strictly as an internal observation by Google rather than an industry-wide fact. (Research report already handled this in Findings & Limitations).

Can Be Approved Without Fix:
YES

---

## Gap 3: Non-Linear Cost Increase Formula (~100x per Nine)

Type:
OVERGENERALIZATION

Severity:
LOW

Location:
`research/03-evidence.md` (Evidence 6), `research/06-open-questions.md`

Problem:
The assertion that each additional "nine" of availability increases cost by 100x is a qualitative heuristic rather than a closed-form mathematical equation.

Required Revision:
Acknowledge as a heuristic rule of thumb. (Already properly contextualized in `06-open-questions.md`).

Can Be Approved Without Fix:
YES

---

## Gap 4: Datadog-Specific Error Budget Remaining Formula

Type:
IMPLEMENTATION_GAP

Severity:
LOW

Location:
`research/03-evidence.md` (Evidence 14), `research/05-report.md` (Finding 13)

Problem:
`error_budget_remaining = 100 * (current_status - target) / (100 - target)` is vendor-specific to Datadog's SLO engine. Google SRE docs define error budget directly as ratio of bad requests to total budget.

Required Revision:
Annotate as a vendor-specific UI/calculation formula. (Already annotated as MEDIUM confidence / vendor-specific).

Can Be Approved Without Fix:
YES
