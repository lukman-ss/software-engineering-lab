# 03 - Claim Audit

Target Lab: `labs/24-slo-sli-error-budget`
Audit Date: 2026-09-28

---

## Claim 1

Claim:
An SLI is a carefully defined quantitative measure of some aspect of the level of service provided, structured as a ratio of good events to total events.

Location:
`research/03-evidence.md` (Evidence 1) & `research/05-report.md` (Finding 1)

Evidence Provided:
Google SRE Book Ch.4 & Google SRE Workbook Ch.2 quotes.

Source:
`https://sre.google/sre-book/service-level-objectives/`

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Fully verified by primary sources.

---

## Claim 2

Claim:
SLO target must be set strictly below 100% because 100% eliminates error budget, prevents velocity, and leads to division-by-zero errors in mathematical models.

Location:
`research/03-evidence.md` (Evidence 9) & `research/05-report.md` (Finding 4)

Evidence Provided:
Datadog SLO documentation & Google SRE Book Ch.3 quotes.

Source:
`https://sre.google/sre-book/embracing-risk/` & `https://docs.datadoghq.com/service_level_objectives/`

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Fundamental SRE principle. Fully verified.

---

## Claim 3

Claim:
Percentiles (P95/P99) should be used instead of arithmetic mean for latency SLIs because average latency hides tail latency.

Location:
`research/03-evidence.md` (Evidence 5) & `research/05-report.md` (Finding 2)

Evidence Provided:
Google SRE Book Ch.4 quotes and Home Depot case study.

Source:
`https://sre.google/sre-book/service-level-objectives/`

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Standard statistical recommendation for distributed systems latency.

---

## Claim 4

Claim:
SLIs must be user-centric (measuring request success, latency, correctness) and not based on internal infrastructure utilization metrics like CPU or memory.

Location:
`research/03-evidence.md` (Evidence 4) & `research/05-report.md` (Finding 2)

Evidence Provided:
Google SRE Workbook Ch.5 and Prometheus Alerting Best Practices.

Source:
`https://sre.google/workbook/alerting-on-slos/` & `https://prometheus.io/docs/practices/alerting/`

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Corroborated across Google and CNCF Prometheus guidelines.

---

## Claim 5

Claim:
Recommended burn rate alerting parameters are: Page on 2% budget in 1h (burn rate 14.4) or 5% in 6h (burn rate 6); Ticket on 10% in 3d (burn rate 1).

Location:
`research/03-evidence.md` (Evidence 6) & `research/05-report.md` (Finding 6)

Evidence Provided:
Google SRE Workbook Ch.5 Table 5-8.

Source:
`https://sre.google/workbook/alerting-on-slos/`

Source Actually Supports Claim:
YES

Classification:
IMPLEMENTATION-SPECIFIC

Severity:
MEDIUM

Notes:
Source presents Table 5-8 as "the starting point for your SLO-based alerting configuration", requiring tuning for specific organizational context and traffic volume.

---

## Claim 6

Claim:
Error budget remaining percentage is calculated as `100 * (current status - target) / (100 - target)`.

Location:
`research/03-evidence.md` (Evidence 7) & `research/05-report.md` (Finding 3)

Evidence Provided:
Datadog Service Level Objectives documentation.

Source:
`https://docs.datadoghq.com/service_level_objectives/`

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Valid mathematical formulation for remaining error budget percentage.

---

## Claim 7

Claim:
Availability "Nines" allowable downtime per month is 7.2 hours for 99%, 43.2 minutes for 99.9%, and 4.32 minutes for 99.99%.

Location:
`research/03-evidence.md` (Evidence 8) & `research/05-report.md` (Finding 5)

Evidence Provided:
Google SRE Book Appendix A Availability Table.

Source:
`https://sre.google/sre-book/availability-table/`

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Exact match for Google's 30-day month convention.

---

## Claim 8

Claim:
OpenSLO is an industry standard specification for defining SLOs declaratively in YAML.

Location:
`research/03-evidence.md` (Evidence 10) & `research/05-report.md` (Finding 8)

Evidence Provided:
OpenSLO specification documentation.

Source:
`https://openslo.github.io/OpenSLO/`

Source Actually Supports Claim:
PARTIAL

Classification:
HYPOTHESIS / INTERPRETATION

Severity:
MEDIUM

Notes:
OpenSLO exists and has an open YAML spec, but evidence for wide industry adoption outside its own repository is limited.
