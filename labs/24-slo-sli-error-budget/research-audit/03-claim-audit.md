# Claim Audit: Research for 24-slo-sli-error-budget

## Claim 1: SLI, SLO, SLA Definitions
Claim: SLI is a quantitative measure of service level; SLO is a target value for an SLI; SLA is a contract with consequences.
Location: `research/05-report.md: Finding 1`, `research/03-evidence.md: Evidence 1`
Evidence Provided: Direct quotes from Google SRE Book Chapter 4.
Source: https://sre.google/sre-book/service-level-objectives/
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Authoritative and accurate.

---

## Claim 2: Common SLI Types by Service Category
Claim: Common SLIs are latency, error rate, throughput, and availability; user-facing, storage, and big data systems prioritize distinct SLIs.
Location: `research/05-report.md: Finding 2`, `research/03-evidence.md: Evidence 2 & 3`
Evidence Provided: Quotes from Google SRE Book Chapter 4.
Source: https://sre.google/sre-book/service-level-objectives/
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Supported by original text.

---

## Claim 3: Availability Calculation Methods (Time-based vs Aggregate)
Claim: Availability can be computed as time-based (uptime / (uptime + downtime)) or aggregate request success rate (successful requests / total requests).
Location: `research/05-report.md: Finding 3`, `research/03-evidence.md: Evidence 4`
Evidence Provided: Direct references from Google SRE Book Chapter 3.
Source: https://sre.google/sre-book/embracing-risk/
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Correctly highlights why aggregate is preferred in distributed services.

---

## Claim 4: Availability Target to Downtime Mapping
Claim: Standard availability percentages map to precise downtime budgets (e.g., 99% = 7.2 hours/month; 99.9% = 43.2 minutes/month; 99.99% = 4.32 minutes/month).
Location: `research/05-report.md: Finding 4`, `research/03-evidence.md: Evidence 5 & 21`
Evidence Provided: Appendix A Availability Table numbers.
Source: https://sre.google/sre-book/availability-table/
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: The research report correctly flags calculation inaccuracies in typical informal approximations (e.g. 7.2 hours = 7h 12m, not 7h 18m).

---

## Claim 5: Error Budget Definition and Governance
Claim: Error budget is (100% - SLO), represents acceptable unreliability, and acts as a control loop governing release velocity and risk management.
Location: `research/05-report.md: Finding 5 & 8`, `research/03-evidence.md: Evidence 6 & 18`
Evidence Provided: Google SRE Book Chapter 3 quotes.
Source: https://sre.google/sre-book/embracing-risk/
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Core SRE principle verified.

---

## Claim 6: Percentile-based Latency SLIs over Averages
Claim: Percentiles (P95, P99, P99.9) are required over averages for latency SLIs because arithmetic averages hide tail latency and distribution skews.
Location: `research/05-report.md: Finding 7`, `research/03-evidence.md: Evidence 10`
Evidence Provided: Google SRE Book Chapter 4 aggregation section.
Source: https://sre.google/sre-book/service-level-objectives/
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Solid foundation verified.

---

## Claim 7: Symptom-Based Alerting
Claim: Alerting should focus on user-visible symptoms (latency, errors) rather than low-level infrastructure causes (CPU/memory).
Location: `research/05-report.md: Finding 9`, `research/03-evidence.md: Evidence 14`
Evidence Provided: Prometheus Alerting documentation and Google SRE Book Chapter 6.
Source: https://prometheus.io/docs/practices/alerting/, https://sre.google/sre-book/monitoring-distributed-systems/
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Verified across both sources.

---

## Claim 8: Non-Linear Cost of Reliability (100x per Nine)
Claim: Each incremental improvement in reliability may cost ~100x more due to redundancy and engineering opportunity costs.
Location: `research/05-report.md: Finding 11`, `research/03-evidence.md: Evidence 17`
Evidence Provided: Quote from Google SRE Book Chapter 3.
Source: https://sre.google/sre-book/embracing-risk/
Source Actually Supports Claim: YES
Classification: INTERPRETATION
Severity: MEDIUM
Notes: The claim is in the SRE book as an illustrative heuristic/rule of thumb, but lacks empirical dataset citations. Research properly categorized this limitation in `06-open-questions.md`.

---

## Claim 9: Vendor-Specific Error Budget Remaining Formula and Burn Rate
Claim: Error budget remaining is computed as `100 * (current status - target) / (100 - target)` and burn rate levels (1-6 elevated, >6 critical).
Location: `research/05-report.md: Finding 8 & 12`, `research/03-evidence.md: Evidence 12 & 13`
Evidence Provided: Datadog documentation quotes.
Source: https://docs.datadoghq.com/service_level_objectives/
Source Actually Supports Claim: YES
Classification: IMPLEMENTATION-SPECIFIC
Severity: MEDIUM
Notes: The research report correctly notes that these formulas and thresholds are Datadog implementation details rather than universal SRE standards.
