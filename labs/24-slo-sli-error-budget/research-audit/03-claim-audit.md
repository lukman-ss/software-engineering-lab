# Claim Audit: SLO, SLI & Error Budget

## Claim 1: Distinctions between SLI, SLO, SLA
Claim: SLI is a quantitative measure of service quality; SLO is a target value for an SLI; SLA is a contract with explicit penalties/consequences for missing SLOs.
Location: `research/03-evidence.md` (Evidence 1), `research/05-report.md` (Finding 1)
Evidence Provided: Quotes from Google SRE Book Chapter 4.
Source: https://sre.google/sre-book/service-level-objectives/
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Canonical definitions from the originators of SRE.

---

## Claim 2: Common SLI Types by System Archetype
Claim: User-facing systems focus on availability, latency, throughput; storage systems on latency, availability, durability; big data systems on throughput, end-to-end latency; all systems on correctness.
Location: `research/03-evidence.md` (Evidence 2 & 3), `research/05-report.md` (Finding 2)
Evidence Provided: Direct quotes from Google SRE Book Chapter 4.
Source: https://sre.google/sre-book/service-level-objectives/
Source Actually Supports Claim: YES
Classification: FACT / BEST_PRACTICE
Severity: LOW
Notes: Fully supported by primary text.

---

## Claim 3: Time-based vs Aggregate Availability Calculation
Claim: Availability can be calculated as time-based `uptime / (uptime + downtime)` or aggregate `successful requests / total requests`. Distributed systems generally prefer aggregate (request success rate / yield).
Location: `research/03-evidence.md` (Evidence 4), `research/05-report.md` (Finding 3)
Evidence Provided: Direct quotes from Google SRE Book Chapter 3.
Source: https://sre.google/sre-book/embracing-risk/
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Accurately distinguishes between aggregate yield and uptime time-slices.

---

## Claim 4: Availability Downtime Windows
Claim: 99% availability equals 7.2 hours/month; 99.9% equals 43.2 minutes/month; 99.99% equals 4.32 minutes/month.
Location: `research/03-evidence.md` (Evidence 5 & 21), `research/05-report.md` (Finding 4)
Evidence Provided: Table from Google SRE Book Appendix A.
Source: https://sre.google/sre-book/availability-table/
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: The research explicitly caught and documented the arithmetic discrepancy in the topic specification (which stated 7h 18m and 4m 23s respectively).

---

## Claim 5: Error Budget Definition and Governance Mechanism
Claim: Error budget equals 100% - SLO. It acts as an objective throttle for deployment risk: if budget is positive, release velocity can proceed; if exhausted, releases freeze or slow while stability engineering takes priority.
Location: `research/03-evidence.md` (Evidence 6 & 18), `research/05-report.md` (Finding 5 & 8)
Evidence Provided: Quotes from Google SRE Book Chapter 3.
Source: https://sre.google/sre-book/embracing-risk/
Source Actually Supports Claim: YES
Classification: FACT / IMPLEMENTATION-SPECIFIC
Severity: LOW
Notes: Verified against primary SRE text.

---

## Claim 6: Latency Distributions and Percentiles vs Averages
Claim: Averages obscure tail latencies and changes across workloads; high percentiles (P95, P99, P99.9) must be used.
Location: `research/03-evidence.md` (Evidence 10), `research/05-report.md` (Finding 7)
Evidence Provided: Quotes and Figure 4-1 reference from Google SRE Book Chapter 4.
Source: https://sre.google/sre-book/service-level-objectives/
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Core statistical principle in SRE telemetry.

---

## Claim 7: Non-linear Cost of Increasing Reliability (100x per nine)
Claim: Improving reliability by an additional nine may cost ~100x more due to redundant hardware and developer opportunity costs.
Location: `research/03-evidence.md` (Evidence 17), `research/05-report.md` (Finding 11)
Evidence Provided: Direct quote from Google SRE Book Chapter 3.
Source: https://sre.google/sre-book/embracing-risk/
Source Actually Supports Claim: YES
Classification: INTERPRETATION / HEURISTIC
Severity: MEDIUM
Notes: Quoted accurately from the SRE Book, but rightly flagged in `06-open-questions.md` as an empirical heuristic without published formal econometric models.

---

## Claim 8: Burn Rate Alerting Thresholds
Claim: Burn rate measures error budget consumption rate; Datadog uses 2-hour window where burn rate 1-6 is elevated and >6 is critical.
Location: `research/03-evidence.md` (Evidence 12), `research/05-report.md` (Finding 8)
Evidence Provided: Datadog documentation excerpt.
Source: https://docs.datadoghq.com/service_level_objectives/
Source Actually Supports Claim: YES
Classification: IMPLEMENTATION-SPECIFIC
Severity: MEDIUM
Notes: Supported by Datadog docs, but not an SRE universal constant. The research accurately labels it as implementation-specific in `04-contradictions.md` and `06-open-questions.md`.

---

## Claim 9: Error Budget Remaining Calculation Formula
Claim: Remaining error budget percentage is calculated as `100 * (current status - target) / (100 - target)`.
Location: `research/03-evidence.md` (Evidence 13), `research/05-report.md` (Finding 12)
Evidence Provided: Datadog documentation.
Source: https://docs.datadoghq.com/service_level_objectives/
Source Actually Supports Claim: YES
Classification: IMPLEMENTATION-SPECIFIC / MATHEMATICAL FORMULATION
Severity: LOW
Notes: Derived mathematically from `(allowed_error - actual_error) / allowed_error * 100`.

---

## Claim 10: Symptom-Based Alerting
Claim: Alerts should page on user-visible symptoms rather than infrastructure causes to reduce alert fatigue and accurately reflect SLO impact.
Location: `research/03-evidence.md` (Evidence 14 & 20), `research/05-report.md` (Finding 9)
Evidence Provided: Prometheus Alerting Best Practices, Google SRE Book Chapter 6.
Source: https://prometheus.io/docs/practices/alerting/, https://sre.google/sre-book/monitoring-distributed-systems/
Source Actually Supports Claim: YES
Classification: FACT / BEST_PRACTICE
Severity: LOW
Notes: Supported by both Google SRE and Prometheus official guidelines.
