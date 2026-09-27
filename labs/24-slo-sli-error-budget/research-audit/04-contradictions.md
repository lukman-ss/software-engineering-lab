# 04 — Contradictions

Audit of contradictions and discrepancies identified in `labs/24-slo-sli-error-budget/research/04-contradictions.md` and `labs/24-slo-sli-error-budget/research/05-report.md`.

---

## Contradiction 1: Monthly Downtime Basis (30 Days vs Gregorian Average Month)

Statement A:
Topic specification / general rule of thumb states 99% availability equals ~7 hours 18 minutes/month (438 minutes).
Location: Topic Spec / `research/04-contradictions.md`

Statement B:
Google SRE Book Appendix A Table 1-1 states 99% availability equals 7.2 hours/month = 7 hours 12 minutes/month (432 minutes).
Location: Google SRE Book Appendix A (`https://sre.google/sre-book/availability-table/`)

Type:
SOURCE_CONFLICT / METHODOLOGICAL_DIFFERENCE

Impact:
LOW. Variance of 6 minutes (1.4%) stems purely from assuming a strict 30-day month (30 * 24 = 720 hrs = 43,200 min -> 1% = 432 min) versus an average Gregorian month of 30.44 days (30.4375 * 24 = 730.5 hrs = 43,830 min -> 1% = 438.3 min).

Assessment:
PASS with clarification. Research document correctly identifies the root cause and notes that both calculations are mathematically valid under their stated assumptions.

---

## Contradiction 2: SLO Measurement Window (30 Days vs 28 Days / 4 Weeks)

Statement A:
Topic lab specification uses a 30-day rolling measurement window.
Location: Topic Spec / `research/04-contradictions.md`

Statement B:
Google SRE Workbook Ch. 2 recommends a 4-week (28-day) rolling window as a general-purpose interval.
Location: Google SRE Workbook Ch. 2 (`https://sre.google/workbook/implementing-slos/`)

Type:
PRACTICE_VARIATION

Impact:
LOW. 28 days guarantees an equal number of weekends/weekdays per window, while 30 days aligns better with monthly business/billing cycles.

Assessment:
PASS with clarification. Both are valid rolling window strategies. Research document transparently explains the operational trade-offs.

---

## Contradiction 3: Burn Rate Alerting Thresholds (Google SRE vs Datadog)

Statement A:
Google SRE Workbook Ch. 5 specifies 14.4x burn rate (1h/5m window) and 6x burn rate (6h/30m window) for pages, and 1x burn rate (3d/6h window) for tickets.
Location: Google SRE Workbook Ch. 5 Table 5-8

Statement B:
Datadog SLO documentation uses a 2-hour window with elevated burn rate (1-6) and critical burn rate (>6).
Location: Datadog SLO Documentation (`https://docs.datadoghq.com/service_level_objectives/`)

Type:
VENDOR_IMPLEMENTATION_DIFFERENCE

Impact:
LOW. Underneath vendor-specific UI/alert settings, both models adhere strictly to multi-window budget consumption rate principles.

Assessment:
PASS. Research report accurately labels these as vendor-specific variations of the same underlying mathematical concept.

---

## Contradiction 4: Empirical Validity of "70% Outage from Change"

Statement A:
Google SRE Workbook Appendix B states "Changes are a major source of instability, representing roughly 70% of our outages".
Location: Google SRE Workbook Appendix B (`https://sre.google/workbook/error-budget-policy/`)

Statement B:
No external academic, cross-industry, or vendor empirical studies corroborate the exact figure of 70%.
Location: `research/04-contradictions.md` (Contradiction 5) & `research/06-open-questions.md`

Type:
UNVERIFIED_GENERALIZATION

Impact:
MEDIUM if stated as an absolute law of software engineering; LOW if contextualized as a Google internal observation.

Assessment:
PASS with restriction. Research report explicitly limits the claim to Google's internal experience and cautions against treating it as a universal industry constant.

---

## Overall Consistency Summary

No material contradictions that invalidate the research core were found. Core concepts (SLI, SLO, SLA, Error Budget, 100% wrong target, ratio SLIs, percentile latency, symptom-based alerting) are 100% consistent across all Tier 1 and Tier 2 sources.
