# 03 — Claim Audit

Audit of claims extracted from `labs/24-slo-sli-error-budget/research/03-evidence.md` and `labs/24-slo-sli-error-budget/research/05-report.md`.

---

## Claim 1: SLI Definition
Claim: "An SLI is a service level indicator—a carefully defined quantitative measure of some aspect of the level of service that is provided."
Location: `03-evidence.md` (Evidence 1), `05-report.md` (Finding 1)
Evidence Provided: Direct quote from Google SRE Book Ch. 4.
Source: Google SRE Book Ch. 4 (`https://sre.google/sre-book/service-level-objectives/`)
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Industry-standard canonical definition.

---

## Claim 2: SLI as Ratio of Good Events / Total Events
Claim: "we generally recommend treating the SLI as the ratio of two numbers: the number of good events divided by the total number of events."
Location: `03-evidence.md` (Evidence 3), `05-report.md` (Finding 2)
Evidence Provided: Direct quote from Google SRE Workbook Ch. 2.
Source: Google SRE Workbook Ch. 2 (`https://sre.google/workbook/implementing-slos/`)
Source Actually Supports Claim: YES
Classification: FACT / BEST_PRACTICE
Severity: LOW
Notes: Supported by both Google SRE and Datadog metric-based SLOs.

---

## Claim 3: SLO Definition and Multi-Target Thresholds
Claim: "An SLO is a service level objective: a target value or range of values for a service level that is measured by an SLI." Can include multiple percentiles/thresholds.
Location: `03-evidence.md` (Evidence 3, 22), `05-report.md` (Finding 3)
Evidence Provided: Direct quote from Google SRE Book Ch. 4.
Source: Google SRE Book Ch. 4
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Standard definition across SRE literature.

---

## Claim 4: SLA vs SLO Distinction
Claim: SLAs are contracts with explicit business/financial consequences; SLOs do not have explicit external consequences. "If there is no explicit consequence, then you are almost certainly looking at an SLO."
Location: `03-evidence.md` (Evidence 4), `05-report.md` (Finding 4)
Evidence Provided: Direct quote and explanation from Google SRE Book Ch. 4.
Source: Google SRE Book Ch. 4
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Accurately reflects industry distinction.

---

## Claim 5: Downtime Table Calculations
Claim: 99% = 7.2 hrs/mo (432 min), 99.9% = 43.2 min/mo, 99.99% = 4.32 min/mo (259 sec), 99.999% = 25.9 sec/mo.
Location: `03-evidence.md` (Evidence 5), `05-report.md` (Finding 5)
Evidence Provided: Google SRE Book Appendix A Table 1-1.
Source: Google SRE Book Appendix A (`https://sre.google/sre-book/availability-table/`)
Source Actually Supports Claim: YES
Classification: FACT (under 30-day month assumption)
Severity: LOW
Notes: Minor variance when using 30.44-day average month (e.g. 7h18m for 99%) is transparently addressed in research.

---

## Claim 6: 100% Reliability Is the Wrong Target
Claim: 100% is the wrong target due to non-zero component failure probability, user device/network unreliability, non-linear cost increase (~100x per nine), and changes being the primary cause of outages.
Location: `03-evidence.md` (Evidence 6), `05-report.md` (Finding 6)
Evidence Provided: Google SRE Workbook Ch. 2 & Book Ch. 3.
Source: Google SRE Workbook Ch. 2, Google SRE Book Ch. 3
Source Actually Supports Claim: YES
Classification: INTERPRETATION / HEURISTIC
Severity: MEDIUM
Notes: Conceptual arguments are fully supported. Numeric heuristic (~100x cost) is a rule of thumb, not an empirical law.

---

## Claim 7: Error Budget Definition
Claim: "An error budget is 1 minus the SLO of the service. A 99.9% SLO service has a 0.1% error budget."
Location: `03-evidence.md` (Evidence 7), `05-report.md` (Finding 7)
Evidence Provided: Direct quote from Google SRE Book Ch. 3.
Source: Google SRE Book Ch. 3 (`https://sre.google/sre-book/embracing-risk/`)
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Mathematical definition widely accepted across industry.

---

## Claim 8: Multi-Window Multi-Burn-Rate Alerting
Claim: Multi-window multi-burn-rate alerting combines long and short windows (short = 1/12 long) to balance precision, recall, detection time, and reset time (e.g., 14.4x/1h+5m for 2% budget page, 6x/6h+30m for 5% budget page, 1x/3d+6h for 10% budget ticket).
Location: `03-evidence.md` (Evidence 9), `05-report.md` (Finding 9)
Evidence Provided: Google SRE Workbook Ch. 5 Table 5-8 and Figure 5-6.
Source: Google SRE Workbook Ch. 5 (`https://sre.google/workbook/alerting-on-slos/`)
Source Actually Supports Claim: YES
Classification: IMPLEMENTATION-SPECIFIC / BEST_PRACTICE
Severity: LOW
Notes: Accurately cites Google SRE Workbook Chapter 5.

---

## Claim 9: Percentiles Over Averages for Latency
Claim: Averages obscure tail latencies; distributions should be measured using high percentiles (P50, P90, P95, P99).
Location: `03-evidence.md` (Evidence 10), `05-report.md` (Finding 10)
Evidence Provided: Google SRE Book Ch. 4 Figure 4-1 & text.
Source: Google SRE Book Ch. 4
Source Actually Supports Claim: YES
Classification: FACT / BEST_PRACTICE
Severity: LOW
Notes: Universally recognized in distributed systems engineering.

---

## Claim 10: Infrastructure Metrics (CPU/RAM) Are Not User SLOs
Claim: SLOs should measure user experience (latency, error rate, throughput), while CPU and memory are diagnostic symptoms / causes.
Location: `03-evidence.md` (Evidence 11, 21), `05-report.md` (Finding 11, 15)
Evidence Provided: Google SRE Book Ch. 4, Ch. 6, Prometheus Alerting docs.
Source: Google SRE Book Ch. 4, Prometheus Alerting Docs
Source Actually Supports Claim: YES
Classification: BEST_PRACTICE
Severity: LOW
Notes: Fully aligned with Google and CNCF Prometheus guidelines.

---

## Claim 11: Endpoint-Specific SLO Bucketing
Claim: Different endpoints require different SLOs based on business criticality (e.g. CRITICAL 99.99%, HIGH_FAST 99.9%, LOW 99%).
Location: `03-evidence.md` (Evidence 12), `05-report.md` (Finding 12)
Evidence Provided: Google SRE Workbook Ch. 5 Table 5-10.
Source: Google SRE Workbook Ch. 5
Source Actually Supports Claim: YES
Classification: FACT / PATTERN
Severity: LOW
Notes: Directly documented in Google SRE Workbook Ch. 5.

---

## Claim 12: Datadog Error Budget Remaining Formula
Claim: `error_budget_remaining = 100 * (current_status - target) / (100 - target)`
Location: `03-evidence.md` (Evidence 14), `05-report.md` (Finding 13)
Evidence Provided: Datadog Documentation.
Source: Datadog SLO docs (`https://docs.datadoghq.com/service_level_objectives/`)
Source Actually Supports Claim: YES
Classification: IMPLEMENTATION-SPECIFIC
Severity: LOW
Notes: Properly classified in research as Datadog-specific formula.

---

## Claim 13: "Changes represent roughly 70% of our outages"
Claim: Changes are responsible for ~70% of production outages.
Location: `03-evidence.md` (Evidence 16), `04-contradictions.md` (Contradiction 5), `05-report.md` (Limitations)
Evidence Provided: Google SRE Workbook Appendix B text.
Source: Google SRE Workbook Appendix B (`https://sre.google/workbook/error-budget-policy/`)
Source Actually Supports Claim: PARTIAL
Classification: HYPOTHESIS / INTERNAL_STATISTIC
Severity: MEDIUM
Notes: Correctly flagged by Research Agent as a Google internal observation without published empirical dataset or external validation. Properly constrained in report limitations.

---

## Claim 14: Sample Calculation Verification
Claim: Webhook with 200,000 requests and 99.99% SLO yields 20 allowed errors; 15 errors consume 75% of budget.
Location: `03-evidence.md` (Evidence 20), `04-contradictions.md` (Verifikasi Hitungan)
Evidence Provided: Manual arithmetic and SRE Book aggregate availability formula.
Source: Internal calculation based on SRE definitions.
Source Actually Supports Claim: YES
Classification: EXAMPLE / MATHEMATICAL_PROOF
Severity: LOW
Notes: Exact arithmetic: 200,000 * (1 - 0.9999) = 20; 15 / 20 = 0.75 (75%).
