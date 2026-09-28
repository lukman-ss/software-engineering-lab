# Evidence

## Evidence 1: SLI Definition

Claim: An SLI is a service level indicator—a carefully defined quantitative measure of some aspect of the level of service that is provided.
Evidence: "An SLI is a service level *indicator*—a carefully defined quantitative measure of some aspect of the level of service that is provided."
Source: Service Level Objectives (Google SRE Book)
URL: https://sre.google/sre-book/service-level-objectives/
Confidence: HIGH
Corroborated By: Source 2 (SRE Workbook), Source 10 (Datadog)
Notes: Most common SLIs include request latency, error rate, system throughput.

## Evidence 2: SLO Definition

Claim: An SLO is a service level objective: a target value or range of values for a service level that is measured by an SLI.
Evidence: "An SLO is a *service level objective*: a target value or range of values for a service level that is measured by an SLI."
Source: Service Level Objectives (Google SRE Book)
URL: https://sre.google/sre-book/service-level-objectives/
Confidence: HIGH
Corroborated By: Source 2, Source 10
Notes: Natural structure is SLI ≤ target or lower bound ≤ SLI ≤ upper bound.

## Evidence 3: Error Budget Concept

Claim: An error budget is the allowed amount of unreliability (100% - SLO) that can be spent on innovation while maintaining reliability targets.
Evidence: "If you have a 99.9% success ratio SLO, then a service that receives 3 million requests over a four-week period had a budget of 3,000 (0.1%) errors over that period."
Source: Implementing SLOs (Google SRE Workbook)
URL: https://sre.google/workbook/implementing-slos/
Confidence: HIGH
Corroborated By: Source 3 (Embracing Risk), Source 10 (Datadog)
Notes: Error budget = 100% - SLO; allows product velocity vs reliability trade-off decisions.

## Evidence 4: SLI Must Be User-Centric

Claim: SLIs should measure user experience, not internal infrastructure metrics like CPU utilization.
Evidence: "SLOs that are meaningful, understood, and represented in metrics, you can configure alerting to notify an on-caller only when there are actionable, specific threats to the error budget."
Source: Alerting on SLOs (Google SRE Workbook)
URL: https://sre.google/workbook/alerting-on-slos/
Confidence: HIGH
Corroborated By: Source 3, Source 9 (Prometheus)
Notes: CPU is infrastructure metric; user cares about success, latency, correctness—not internal utilization.

## Evidence 5: Percentiles Over Averages for Latency SLIs

Claim: Use percentile latency (P95/P99) not average latency for user-facing SLIs because averages hide tail latency.
Evidence: "Most metrics are better thought of as *distributions* rather than averages. For example, for a latency SLI, some requests will be serviced quickly, while others will invariably take longer—sometimes much longer. A simple average can obscure these tail latencies."
Source: Service Level Objectives (Google SRE Book)
URL: https://sre.google/sre-book/service-level-objectives/
Confidence: HIGH
Corroborated By: Source 6 (Home Depot case study), Source 10 (Datadog)
Notes: Figure 4-1 shows 50th vs 99th percentile latency; P99 can be 20x slower than median.

## Evidence 6: Burn Rate Alerting Parameters

Claim: Recommended SLO-based alerting uses multi-window, multi-burn-rate: page on 2% budget in 1h (burn rate 14.4) or 5% in 6h (burn rate 6); ticket on 10% in 3d (burn rate 1).
Evidence: "We recommend the parameters listed in Table 5-8 as the starting point for your SLO-based alerting configuration."
Source: Alerting on SLOs (Google SRE Workbook)
URL: https://sre.google/workbook/alerting-on-slos/
Confidence: HIGH
Corroborated By: Source 4 (same), implicit in Source 10 (Datadog burn rate indicator)
Notes: Table 5-8 specifies: Page: 1h/5m @ 14.4x burn rate (2% budget); Page: 6h/30m @ 6x burn rate (5% budget); Ticket: 3d/6h @ 1x burn rate (10% budget).

## Evidence 7: Error Budget Calculation Formula

Claim: Error budget remaining % = 100 × (current status - target) / (100 - target)
Evidence: "The remaining error budget is displayed as a percentage and is calculated using the following formula: $$\text"error budget remaining" = 100 * {\text"current status" - \text" target"} / { 100 - \text"target"}$$"
Source: Service Level Objectives (Datadog)
URL: https://docs.datadoghq.com/service_level_objectives/
Confidence: HIGH
Corroborated By: Source 2 (Implementing SLOs), Source 3 (Embracing Risk)
Notes: Example: 95% current status on 90% SLO → (95-90)/(100-90) = 5/10 = 50% error budget remaining.

## Evidence 8: Availability "Nines" Downtime

Claim: 99.9% availability allows 43.2 minutes downtime per month; 99.99% allows 4.32 minutes per month.
Evidence: From Availability Table: 99.9% → 43.2 minutes/month; 99.99% → 4.32 minutes/month
Source: Availability Table (Google SRE Book)
URL: https://sre.google/sre-book/availability-table/
Confidence: HIGH
Corroborated By: Source 3 (Embracing Risk discusses "nines"), basic calculation verification
Notes: 100% - 99.9% = 0.1% downtime; 0.1% × 30 days × 24h × 60m = 43.2 minutes.

## Evidence 9: SLO Must Be Below 100%

Claim: Setting SLO at 100% eliminates error budget and prevents innovation/reliability balance.
Evidence: "Setting a 100% target means having an error budget of 0% since error budget is equal to 100%—SLO target. Without error budget representing acceptable risk, you face difficulty finding alignment between the conflicting priorities of maintaining customer-facing reliability and investing in feature development."
Source: Service Level Objectives (Datadog)
URL: https://docs.datadoghq.com/service_level_objectives/
Confidence: HIGH
Corroborated By: Source 2, Source 3
Notes: Also causes division-by-zero in alert evaluations.

## Evidence 10: OpenSLO Declarative Format

Claim: OpenSLO provides vendor-neutral YAML specification for SLOs as code.
Evidence: "OpenSLO is a service level objective (SLO) language that declaratively defines reliability and performance targets using a simple YAML specification."
Source: OpenSLO Specification
URL: https://openslo.github.io/OpenSLO/
Confidence: HIGH
Corroborated By: Source 8 (self), implicit in Source 6 (case studies show need for standardization)
Notes: Apache 2.0 license; enables Git workflow integration; schema versions v1, v2alpha.

## Evidence 11: Evernote SLO Adoption Journey

Claim: Evernote adopted 99.95% monthly uptime SLO after moving to GCP, improved dev/ops alignment.
Evidence: "Our first SLOs document contained the following: A definition of the SLOs: This was an uptime measure: 99.95% uptime measured over a monthly window..."
Source: SLO Engineering Case Studies (Google SRE Workbook - Evernote)
URL: https://sre.google/workbook/slo-engineering-case-studies/
Confidence: HIGH
Corroborated By: Source 6 (Home Depot case study), Source 2 (implementing SLOs guidance)
Notes: After 9 months, on SLO v3; shared SLO dashboards with Google CRE team; aligned priorities.

## Evidence 12: Home Depot VALET Framework

Claim: The Home Depot created VALET (Volume, Availability, Latency, Errors, Tickets) SLO framework, scaled to 800 services in <1 year.
Evidence: "After tracking SLOs for about 50 services at the beginning of the year, by the end of the year we were tracking SLOs for 800 services, with about 50 new services per month being registered with VALET."
Source: SLO Engineering Case Studies (Google SRE Workbook - Home Depot)
URL: https://sre.google/workbook/slo-engineering-case-studies/
Confidence: HIGH
Corroborated By: Source 6 (self), Source 2 (implementing SLOs)
Notes: VALET became part of annual performance reviews; simplified metrics for business owners (99.5% = MVP, 99.9% = adequate, 99.95% = selling systems, 99.99% = shared infra).

## Evidence 13: Alert on Symptoms Not Causes

Claim: Alert on user-visible symptoms (latency, error rate) not root causes (CPU, disk) to avoid alert fatigue.
Evidence: "Aim to have as few alerts as possible, by alerting on symptoms that are associated with end-user pain rather than trying to catch every possible way that pain could be caused."
Source: Alerting Best Practices (Prometheus)
URL: https://prometheus.io/docs/practices/alerting/
Confidence: HIGH
Corroborated By: Source 4, Source 10
Notes: Only page on latency at one point in stack; for error rates, page on user-visible errors.

## Evidence 14: SLO Decision Matrix

Claim: SLO compliance, toil, and customer satisfaction determine action: tighten/loosen SLO, adjust toil, change release velocity.
Evidence: Table 2-5 shows decision matrix: e.g., SLO Met + Low Toil + High Sat → (a) increase velocity or (b) step back; SLO Missed + High Toil + Low Sat → Offload toil and fix product.
Source: Implementing SLOs (Google SRE Workbook)
URL: https://sre.google/workbook/implementing-slos/
Confidence: HIGH
Corroborated By: Source 2 (self), logical extension of error budget policy
Notes: Provides actionable framework for SLO-based prioritization beyond just "stop releases".

## Evidence 15: Low-Traffic Service Alerting Challenges

Claim: Low-traffic services can cause alerting noise; single failed request may represent high burn rate.
Evidence: "If a system receives 10 requests per hour, then a single failed request results in an hourly error rate of 10%. For a 99.9% SLO, this request constitutes a 1,000x burn rate and would page immediately..."
Source: Alerting on SLOs (Google SRE Workbook)
URL: https://sre.google/workbook/alerting-on-slos/
Confidence: HIGH
Corroborated By: Source 4 (self), Source 9 (Prometheus mentions low-traffic considerations)
Notes: Solutions: generate artificial traffic, combine services, modify client retry logic, lower SLO or increase window.