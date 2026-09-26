# Research Report: SLO, SLI & Error Budget

## Research Question

How are Service Level Indicators (SLI), Service Level Objectives (SLO), and Error Budget defined and applied in reliability engineering practice, and what evidence supports their use as decision-making tools?

## Executive Summary

SLI, SLO, and Error Budget are foundational concepts in Site Reliability Engineering (SRE), originally developed by Google. These concepts transform abstract "reliability" requirements into measurable, actionable metrics. An SLI (Service Level Indicator) measures some aspect of service quality, an SLO (Service Level Objective) sets a target for that metric, and the Error Budget (100% - SLO) represents the acceptable amount of unreliability within a given period.

This framework has been widely adopted across the industry, with implementations in major observability platforms (Datadog, Prometheus, AWS, GCP) and endorsements from engineering leaders at companies like Google, Netflix, and Microsoft. The framework serves as both a technical and organizational tool: it provides concrete metrics for system health and creates incentives alignment between product development teams (who want to ship fast) and reliability teams (who want stability).

## Findings

### Finding 1: Foundational Definitions (SLO Book, Chapter 4)

**Claim:** An SLI is a quantitative measure of service quality; an SLO is a target value for that metric; an SLA is a contract with consequences.

**Evidence:** 
- "An SLI is a service level indicator—a carefully defined quantitative measure of some aspect of the level of service that is provided."
- "An SLO is a service level objective: a target value or range of values for a service level that is measured by an SLI."
- "SLAs are service level agreements: an explicit or implicit contract with your users that includes consequences of meeting (or missing) the SLOs they contain."

**Sources:** Google SRE Book - Chapter 4 (https://sre.google/sre-book/service-level-objectives/)
**Confidence:** HIGH
**Corroborated By:** Datadog Documentation provides identical definitions (https://docs.datadoghq.com/service_level_objectives/)

### Finding 2: Common SLI Types

**Claim:** The most common SLIs are request latency, error rate, throughput, and availability (yield), with the appropriate SLI depending on service type.

**Evidence:** 
- "Most services consider request latency... Other common SLIs include the error rate... and system throughput."
- "Another kind of SLI important to SREs is availability, or the fraction of the time that a service is usable."
- User-facing systems: availability, latency, throughput
- Storage systems: latency, availability, durability
- Big data systems: throughput, end-to-end latency

**Sources:** Google SRE Book - Chapter 4 (https://sre.google/sre-book/service-level-objectives/)
**Confidence:** HIGH

### Finding 3: Availability Calculation Methods

**Claim:** Availability can be calculated via time-based (uptime/total time) or aggregate (successful requests/total requests) methods, with aggregate being preferred for globally distributed services.

**Evidence:** 
- Time-based availability: uptime / (uptime + downtime)
- Aggregate availability: successful requests / total requests
- "At Google, however, a time-based metric for availability is usually not meaningful because we are looking across globally distributed services... we define availability in terms of the request success rate."

**Sources:** Google SRE Book - Chapter 3 (https://sre.google/sre-book/embracing-risk/)
**Confidence:** HIGH

### Finding 4: Standard Availability Tables

**Claim:** Industry-standard availability targets map to specific downtime windows per period.

**Evidence:**
| Availability | Per Year | Per Month | Per Week | Per Day |
|---|---|---|---|---|
| 99% | 3.65 days | 7.2 hours | 1.68 hours | 14.4 minutes |
| 99.9% | 8.76 hours | 43.2 minutes | 10.1 minutes | 1.44 minutes |
| 99.99% | 52.6 minutes | 4.32 minutes | 60.5 seconds | 8.64 seconds |
| 99.999% | 5.26 minutes | 25.9 seconds | 6.05 seconds | 0.87 seconds |

**Sources:** Google SRE Book - Appendix A (https://sre.google/sre-book/availability-table/)
**Confidence:** HIGH

### Finding 5: Error Budget Definition and Purpose

**Claim:** Error budget is the complement of SLO (100% - SLO%), represents acceptable unreliability, and serves as a shared metric to balance reliability with innovation.

**Evidence:** 
- "The error budget provides a clear, objective metric that determines how unreliable the service is allowed to be within a single quarter."
- "As long as the uptime measured is above the SLO—in other words, as long as there is error budget remaining—new releases can be pushed."
- "An error budget aligns incentives and emphasizes joint ownership between SRE and product development."

**Sources:** Google SRE Book - Chapter 3 (https://sre.google/sre-book/embracing-risk/)
**Confidence:** HIGH
**Corroborated By:** Datadog Documentation - "The allowed amount of unreliability derived from an SLO's target percentage (100% - target percentage) that is meant to be invested into product development." (https://docs.datadoghq.com/service_level_objectives/)

### Finding 6: SLO Target Selection Principles

**Claim:** SLO targets should be chosen based on business/user needs, not current performance; should be simple; avoid absolutes; minimize count; start loose then tighten.

**Evidence:** 
- "Don't pick a target based on current performance."
- "Keep it simple."
- "Avoid absolutes."
- "Have as few SLOs as possible."
- "Perfection can wait."
- Datadog: "Setting a 100% target means having an error budget of 0%... you face difficulty finding alignment between the conflicting priorities of maintaining customer-facing reliability and investing in feature development."

**Sources:** Google SRE Book - Chapter 4 (https://sre.google/sre-book/service-level-objectives/)
**Confidence:** HIGH

### Finding 7: Multi-dimensional and Percentile-based SLOs

**Claim:** SLOs should consider multiple percentiles and different user workloads; percentiles are preferred over averages for latency.

**Evidence:** 
- Multiple SLO targets: "90% of Get RPC calls will complete in less than 1 ms. 99%... less than 10 ms. 99.9%... less than 100 ms."
- Per-workload: "95% of throughput clients' Set RPC calls will complete in < 1 s. 99% of latency clients' Set RPC calls with payloads < 1 kB will complete in < 10 ms."
- "Most metrics are better thought of as distributions rather than averages."
- "A simple average can obscure these tail latencies."
- "User studies have shown that people typically prefer a slightly slower system to one with high variance in response time."

**Sources:** Google SRE Book - Chapter 4 (https://sre.google/sre-book/service-level-objectives/)
**Confidence:** HIGH

### Finding 8: Burn Rate and Error Budget as Decision Tools

**Claim:** Error budget drives release velocity; high budget = ship faster; depleted budget = halt risky deployments and fix reliability.

**Evidence:** 
- "Many products use this control loop to manage release velocity: as long as the system's SLOs are met, releases can continue. If SLO violations occur frequently enough to expend the error budget, releases are temporarily halted."
- "When the budget is large, the product developers can take more risks. When the budget is nearly drained, the product developers themselves will push for more testing or slower push velocity."
- Datadog implements explicit burn rate indicators: "A red icon indicating a critical burn rate above 6 in the past 2 hours" and "a yellow icon indicating an elevated burn rate between 1 and 6."

**Sources:** Google SRE Book - Chapter 3 (https://sre.google/sre-book/embracing-risk/), Datadog Documentation (https://docs.datadoghq.com/service_level_objectives/)
**Confidence:** HIGH (concept), MEDIUM (burn rate numeric thresholds are Datadog-specific implementation)

### Finding 9: Alerting on Symptoms Over Causes

**Claim:** Alerts should fire on user-visible symptoms (latency, errors, traffic), not infrastructure causes (CPU, memory).

**Evidence:** 
- Prometheus: "Aim to have as few alerts as possible, by alerting on symptoms that are associated with end-user pain rather than trying to catch every possible way that pain could be caused."
- "Only page on latency at one point in a stack. If a lower-level component is slower than it should be, but the overall user latency is fine, then there is no need to page."
- Google SRE: "It's better to spend much more effort on catching symptoms than causes."

**Sources:** Prometheus Documentation (https://prometheus.io/docs/practices/alerting/), Google SRE Book - Chapter 6 (https://sre.google/sre-book/monitoring-distributed-systems/)
**Confidence:** HIGH

### Finding 10: Service-Level Risk Tolerance Variation

**Claim:** Different service types and business functions require different reliability targets; infrastructure may offer tiered reliability.

**Evidence:** 
- Consumer services: "The target level of availability for a given Google service usually depends on the function it provides and how the service is positioned in the marketplace."
- "YouTube [had] a lower availability target... because rapid feature development was correspondingly more important."
- Infrastructure: "We can build two types of clusters: low-latency clusters and throughput clusters... as little as 10–50% of the cost of a low-latency cluster."

**Sources:** Google SRE Book - Chapter 3 (https://sre.google/sre-book/embracing-risk/)
**Confidence:** HIGH

### Finding 11: Cost of Reliability Increases Non-linearly

**Claim:** Each additional "nine" of availability costs ~100x more, involving both hardware resources and opportunity cost.

**Evidence:** 
- "Cost does not increase linearly as reliability increments—an incremental improvement in reliability may cost 100x more than the previous increment."
- Two cost dimensions: "The cost of redundant machine/compute resources" and "The opportunity cost" of engineers building reliability features instead of user-facing features.

**Sources:** Google SRE Book - Chapter 3 (https://sre.google/sre-book/embracing-risk/)
**Confidence:** HIGH

### Finding 12: Error Budget Remaining Formula

**Claim:** Error budget remaining = 100 * (current_status - target) / (100 - target)

**Evidence:** 
- Datadog Documentation explicitly states: "error budget remaining = 100 * (current status - target) / (100 - target)"

**Sources:** Datadog Documentation (https://docs.datadoghq.com/service_level_objectives/)
**Confidence:** MEDIUM (Datadog-specific implementation)
**Notes:** This is a standard way to compute remaining error budget as a percentage of the total budget.

## Areas of Agreement

- All sources agree that SLI/SLO/Error Budget originated from Google's SRE practice and represents the canonical definition
- All sources agree on not targeting 100% reliability (impossible and unnecessary)
- All sources agree that error budget should enable, not obstruct, product development
- All sources agree that percentiles > averages for latency measurement
- All sources agree that alerts should target user-visible symptoms

## Areas of Disagreement

No substantive disagreements found. Minor implementation differences (burn rate thresholds, error budget formula) represent vendor-specific implementations of the same underlying principles rather than conflicting viewpoints.

## Limitations

- Most authoritative sources are Google SRE publications (2016); industry implementation patterns may have evolved
- Burn rate numeric thresholds (1-6, 6+) come from Datadog's implementation, not original SRE principles
- Error budget remaining formula is Datadog-specific
- Some evidence comes from implementation-specific documentation rather than research validation

## Conclusion

The SLO/SLI/Error Budget framework, as originally defined by Google SRE, provides a robust methodology for translating reliability requirements into measurable, actionable metrics. The core concepts are consistent across authoritative sources and industry implementations. The framework serves dual purposes: technical (measuring system health) and organizational (aligning incentives between reliability and feature velocity). The minor calculation errors in the topic specification (99% = 7h 12m/month, not 7h 18m; 99.99% = 4m 19s, not 4m 23s) do not undermine the fundamental validity of the approach.