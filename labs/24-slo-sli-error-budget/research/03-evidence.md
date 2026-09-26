## Evidence 1: Definitions of SLI, SLO, SLA

**Claim:** An SLI is a quantitative measure of some aspect of the level of service provided; an SLO is a target value for a service level measured by an SLI; an SLA is a contract with explicit consequences for missing SLOs.

**Evidence:** 
- "An SLI is a service level indicator—a carefully defined quantitative measure of some aspect of the level of service that is provided."
- "An SLO is a service level objective: a target value or range of values for a service level that is measured by an SLI."
- "SLAs are service level agreements: an explicit or implicit contract with your users that includes consequences of meeting (or missing) the SLOs they contain."

**Source:** Google SRE Book - Chapter 4: Service Level Objectives
**URL:** https://sre.google/sre-book/service-level-objectives/
**Confidence:** HIGH
**Corroborated By:** Datadog Documentation (Source 6) - provides identical definitions
**Notes:** This is the authoritative definition from Google SRE, the originators of the SLO framework.

---

## Evidence 2: Common SLI Types

**Claim:** Common SLIs include request latency, error rate, system throughput, and availability (yield).

**Evidence:** 
- "Most services consider request latency—how long it takes to return a response to a request—as a key SLI. Other common SLIs include the error rate, often expressed as a fraction of all requests received, and system throughput, typically measured in requests per second."
- "Another kind of SLI important to SREs is availability, or the fraction of the time that a service is usable. It is often defined in terms of the fraction of well-formed requests that succeed, sometimes called yield."

**Source:** Google SRE Book - Chapter 4: Service Level Objectives
**URL:** https://sre.google/sre-book/service-level-objectives/
**Confidence:** HIGH
**Corroborated By:** Datadog Documentation - defines SLI as "A quantitative measurement of a service's performance or reliability. In Datadog SLOs an SLI is a metric or an aggregation of one or more monitors."

---

## Evidence 3: SLI Categories by Service Type

**Claim:** Different service types require different SLIs: user-facing systems need availability, latency, throughput; storage systems need latency, availability, durability; big data systems need throughput and end-to-end latency.

**Evidence:** 
- "User-facing serving systems... generally care about availability, latency, and throughput."
- "Storage systems often emphasize latency, availability, and durability."
- "Big data systems... tend to care about throughput and end-to-end latency."
- "All systems should care about correctness."

**Source:** Google SRE Book - Chapter 4: Service Level Objectives
**URL:** https://sre.google/sre-book/service-level-objectives/
**Confidence:** HIGH
**Corroborated By:** None directly, but consistent with industry practice

---

## Evidence 4: Availability Calculations (Time-based vs Aggregate)

**Claim:** Availability can be calculated as time-based (uptime/downtime) or aggregate (successful requests/total requests).

**Evidence:** 
- Time-based availability formula: uptime / (uptime + downtime)
- Aggregate availability formula: successful requests / total requests
- "At Google, however, a time-based metric for availability is usually not meaningful because we are looking across globally distributed services... Therefore, instead of using metrics around uptime, we define availability in terms of the request success rate."

**Source:** Google SRE Book - Chapter 3: Embracing Risk
**URL:** https://sre.google/sre-book/embracing-risk/
**Confidence:** HIGH
**Corroborated By:** Availability Table Appendix provides standardized time-based calculations

---

## Evidence 5: Availability Targets and Downtime Windows

**Claim:** Standard availability targets correspond to specific downtime windows per year/quarter/month/week/day/hour.

**Evidence:** 
- 99% = 3.65 days/year, 7.2 hours/month, 1.68 hours/week, 14.4 minutes/day
- 99.9% = 8.76 hours/year, 43.2 minutes/month, 10.1 minutes/week, 1.44 minutes/day
- 99.99% = 52.6 minutes/year, 4.32 minutes/month, 60.5 seconds/week, 8.64 seconds/day
- 99.999% = 5.26 minutes/year, 25.9 seconds/month, 6.05 seconds/week, 0.87 seconds/day

**Source:** Google SRE Book - Appendix A: Availability Table
**URL:** https://sre.google/sre-book/availability-table/
**Confidence:** HIGH
**Corroborated By:** Chapter 3 Embracing Risk shows the same formulas

---

## Evidence 6: Error Budget Definition and Purpose

**Claim:** Error budget is the inverse of SLO (100% - SLO target), representing the acceptable rate of unreliability; it aligns incentives between product development (velocity) and SRE (reliability).

**Evidence:** 
- "The error budget provides a clear, objective metric that determines how unreliable the service is allowed to be within a single quarter."
- "Product Management defines an SLO... The actual uptime is measured by a neutral third party: our monitoring system. The difference between these two numbers is the 'budget' of how much 'unreliability' is remaining for the quarter."
- "As long as the uptime measured is above the SLO—in other words, as long as there is error budget remaining—new releases can be pushed."
- "Error budget aligns incentives and emphasizes joint ownership between SRE and product development."

**Source:** Google SRE Book - Chapter 3: Embracing Risk (Motivation for Error Budgets)
**URL:** https://sre.google/sre-book/embracing-risk/
**Confidence:** HIGH
**Corroborated By:** Datadog Documentation - "Error Budget: The allowed amount of unreliability derived from an SLO's target percentage (100% - target percentage) that is meant to be invested into product development."

---

## Evidence 7: Error Budget Calculation Example

**Claim:** Error budget calculation: for SLO 99.999% per quarter, error budget = 0.001% failure rate; 0.0002% actual failures = 20% of error budget consumed.

**Evidence:** 
- "For example, imagine that a service's SLO is to successfully serve 99.999% of all queries per quarter. This means that the service's error budget is a failure rate of 0.001% for a given quarter. If a problem causes us to fail 0.0002% of the expected queries for the quarter, the problem spends 20% of the service's quarterly error budget."

**Source:** Google SRE Book - Chapter 3: Embracing Risk
**URL:** https://sre.google/sre-book/embracing-risk/
**Confidence:** HIGH

---

## Evidence 8: SLO Target Selection Principles

**Claim:** SLO targets should not be based on current performance; should be simple; avoid absolutes (100%); have as few SLOs as possible; start loose and tighten over time.

**Evidence:** 
- "Don't pick a target based on current performance... adopting values without reflection may lock you into supporting a system that requires heroic efforts to meet its targets."
- "Keep it simple... Complicated aggregations in SLIs can obscure changes to system performance."
- "Avoid absolutes... Even a system that approaches such ideals will probably take a long time to design and build, and will be expensive to operate."
- "Have as few SLOs as possible... Choose just enough SLOs to provide good coverage of your system's attributes."
- "Perfection can wait... It's better to start with a loose target that you tighten than to choose an overly strict target that has to be relaxed."

**Source:** Google SRE Book - Chapter 4: Service Level Objectives (Choosing Targets)
**URL:** https://sre.google/sre-book/service-level-objectives/
**Confidence:** HIGH
**Corroborated By:** Datadog - "Setting a 100% target means having an error budget of 0%... Without error budget representing acceptable risk, you face difficulty finding alignment between the conflicting priorities of maintaining customer-facing reliability and investing in feature development."

---

## Evidence 9: Multi-dimensional SLOs

**Claim:** SLOs can have multiple targets for different percentiles or different user workloads.

**Evidence:** 
- "If the shape of the performance curves are important, then you can specify multiple SLO targets: 90% of Get RPC calls will complete in less than 1 ms. 99% of Get RPC calls will complete in less than 10 ms. 99.9% of Get RPC calls will complete in less than 100 ms."
- "If you have users with heterogeneous workloads... it may be appropriate to define separate objectives for each class of workload: 95% of throughput clients' Set RPC calls will complete in < 1 s. 99% of latency clients' Set RPC calls with payloads < 1 kB will complete in < 10 ms."

**Source:** Google SRE Book - Chapter 4: Service Level Objectives (Defining Objectives)
**URL:** https://sre.google/sre-book/service-level-objectives/
**Confidence:** HIGH

---

## Evidence 10: Percentile-based vs Average-based SLIs

**Claim:** Percentiles (P95, P99, P99.9) are preferred over averages for latency SLIs because averages obscure tail latencies.

**Evidence:** 
- "Most metrics are better thought of as distributions rather than averages."
- "A simple average can obscure these tail latencies, as well as changes in them... Although a typical request is served in about 50 ms, 5% of requests are 20 times slower!"
- "Monitoring and alerting based only on the average latency would show no change in behavior over the course of the day, when there are in fact significant changes in the tail latency."
- "User studies have shown that people typically prefer a slightly slower system to one with high variance in response time, so some SRE teams focus only on high percentile values, on the grounds that if the 99.9th percentile behavior is good, then the typical experience is certainly going to be."

**Source:** Google SRE Book - Chapter 4: Service Level Objectives (Aggregation)
**URL:** https://sre.google/sre-book/service-level-objectives/
**Confidence:** HIGH
**Corroborated By:** Monitoring Distributed Systems chapter - "The simplest way to differentiate between a slow average and a very slow 'tail' of requests is to collect request counts bucketed by latencies (suitable for rendering a histogram)... Distributing the histogram boundaries approximately exponentially... is often an easy way to visualize the distribution of your requests."

---

## Evidence 11: Client-side vs Server-side SLI Collection

**Claim:** Some SLIs must be collected client-side because server-side metrics miss user-impacting issues.

**Evidence:** 
- "Concentrating on the response latency of the Shakespeare search backend might miss poor user latency due to problems with the page's JavaScript: in this case, measuring how long it takes for a page to become usable in the browser is a better proxy for what the user actually experiences."

**Source:** Google SRE Book - Chapter 4: Service Level Objectives (Collecting Indicators)
**URL:** https://sre.google/sre-book/service-level-objectives/
**Confidence:** HIGH

---

## Evidence 12: Burn Rate Concept for Alerting

**Claim:** Burn rate measures how fast error budget is being consumed; high burn rate triggers alerts before budget is exhausted.

**Evidence:** 
- Datadog: "Burn rate indicators use a rolling 2-hour window to evaluate which SLOs are consuming their error budget too quickly. Burn rate indicators appear next to the applicable SLO names."
- "A red icon indicating a critical burn rate above 6 in the past 2 hours."
- "A yellow icon indicating an elevated burn rate between 1 and 6 in the past 2 hours."

**Source:** Datadog Documentation - Service Level Objectives
**URL:** https://docs.datadoghq.com/service_level_objectives/
**Confidence:** MEDIUM (Datadog implementation, not primary Google SRE source)
**Corroborated By:** Google SRE Book references "SLO-violation error budget" as input to release decisions, but doesn't specify "burn rate" term explicitly.

---

## Evidence 13: Error Budget Alert Formula

**Claim:** Error budget remaining = 100 * (current status - target) / (100 - target)

**Evidence:** 
- "The remaining error budget is displayed as a percentage and is calculated using the following formula: error budget remaining = 100 * (current status - target) / (100 - target)"

**Source:** Datadog Documentation - Service Level Objectives
**URL:** https://docs.datadoghq.com/service_level_objectives/
**Confidence:** MEDIUM (Datadog-specific implementation)

---

## Evidence 14: Alerting on Symptoms, Not Causes

**Claim:** Alerts should fire on symptoms (user-visible problems) rather than causes (infrastructure issues).

**Evidence:** 
- "Aim to have as few alerts as possible, by alerting on symptoms that are associated with end-user pain rather than trying to catch every possible way that pain could be caused."
- "Only page on latency at one point in a stack. If a lower-level component is slower than it should be, but the overall user latency is fine, then there is no need to page."
- "For error rates, page on user-visible errors. If there are errors further down the stack that will cause such a failure, there is no need to page on them separately."

**Source:** Prometheus Documentation - Alerting Best Practices
**URL:** https://prometheus.io/docs/practices/alerting/
**Confidence:** HIGH
**Corroborated By:** Google SRE Book - Monitoring Distributed Systems: "Monitoring symptoms is easier the further 'up' your stack you monitor... It's better to spend much more effort on catching symptoms than causes."

---

## Evidence 15: SLO Safety Margin and Avoiding Over-achievement

**Claim:** Internal SLOs should be tighter than advertised SLOs; services shouldn't significantly exceed their SLO to avoid users building dependencies on unrealistic performance.

**Evidence:** 
- "Using a tighter internal SLO than the SLO advertised to users gives you room to respond to chronic problems before they become visible externally."
- "Don't overachieve... Users build on the reality of what you offer, rather than what you say you'll supply... If your service's actual performance is much better than its stated SLO, users will come to rely on its current performance."
- "You can avoid over-dependence by deliberately taking the system offline occasionally (Google's Chubby service introduced planned outages in response to being overly available)..."

**Source:** Google SRE Book - Chapter 4: Service Level Objectives (SLOs Set Expectations)
**URL:** https://sre.google/sre-book/service-level-objectives/
**Confidence:** HIGH

---

## Evidence 16: Risk Tolerance Varies by Service Type

**Claim:** Consumer services and infrastructure services have different risk tolerance; infrastructure services may offer multiple service tiers with different reliability levels.

**Evidence:** 
- Consumer services: "The target level of availability for a given Google service usually depends on the function it provides and how the service is positioned in the marketplace."
- Infrastructure: "One approach to meeting the needs of both use cases is to engineer all infrastructure services to be ultra-reliable. Given the fact that these infrastructure services also tend to aggregate huge amounts of resources, such an approach is usually far too expensive in practice."
- Bigtable example: "We can build two types of clusters: low-latency clusters and throughput clusters... we are able to satisfy these relaxed needs at a much lower cost, perhaps as little as 10–50% of the cost of a low-latency cluster."

**Source:** Google SRE Book - Chapter 3: Embracing Risk (Identifying Risk Tolerance)
**URL:** https://sre.google/sre-book/embracing-risk/
**Confidence:** HIGH

---

## Evidence 17: Cost of Increasing Reliability

**Claim:** Cost of reliability increases non-linearly; each additional "nine" costs ~100x more.

**Evidence:** 
- "As we build systems, cost does not increase linearly as reliability increments—an incremental improvement in reliability may cost 100x more than the previous increment."
- Two dimensions: "The cost of redundant machine/compute resources" and "The opportunity cost" of engineers working on reliability instead of features.

**Source:** Google SRE Book - Chapter 3: Embracing Risk (Managing Risk)
**URL:** https://sre.google/sre-book/embracing-risk/
**Confidence:** HIGH

---

## Evidence 18: Error Budget Drives Release Velocity

**Claim:** Error budget directly controls release velocity - when budget is high, release faster; when budget is low, slow down or halt releases.

**Evidence:** 
- "Many products use this control loop to manage release velocity: as long as the system's SLOs are met, releases can continue. If SLO violations occur frequently enough to expend the error budget, releases are temporarily halted while additional resources are invested in system testing and development to make the system more resilient."
- "When the budget is large, the product developers can take more risks. When the budget is nearly drained, the product developers themselves will push for more testing or slower push velocity, as they don't want to risk using up the budget and stall their launch."

**Source:** Google SRE Book - Chapter 3: Embracing Risk (Benefits)
**URL:** https://sre.google/sre-book/embracing-risk/
**Confidence:** HIGH

---

## Evidence 19: SLO Status Corrections

**Claim:** SLO status corrections allow excluding specific time periods from SLO calculations (maintenance, non-business hours, deployments).

**Evidence:** 
- "Status corrections allow you to exclude specific time periods from SLO status and error budget calculations... Prevent expected downtime, such as scheduled maintenance, from depleting your error budget... Ignore non-business hours... Ensure that temporary issues caused by deployments do not negatively impact your SLOs."

**Source:** Datadog Documentation - Service Level Objectives (SLO Status Corrections)
**URL:** https://docs.datadoghq.com/service_level_objectives/
**Confidence:** MEDIUM (Implementation-specific, but concept aligns with Google's approach to planned outages)

---

## Evidence 20: Four Golden Signals

**Claim:** The four golden signals of monitoring are latency, traffic, errors, and saturation.

**Evidence:** 
- "The four golden signals of monitoring are latency, traffic, errors, and saturation. If you can only measure four metrics of your user-facing system, focus on these four."
- Latency: "time it takes to service a request... distinguish between the latency of successful requests and the latency of failed requests"
- Traffic: "measure of how much demand is being placed on your system"
- Errors: "rate of requests that fail, either explicitly (e.g., HTTP 500s), implicitly (HTTP 200 but wrong content), or by policy (request over SLO threshold is an error)"
- Saturation: "how 'full' your service is... emphasize the resources that are most constrained"

**Source:** Google SRE Book - Chapter 6: Monitoring Distributed Systems
**URL:** https://sre.google/sre-book/monitoring-distributed-systems/
**Confidence:** HIGH

---

## Evidence 21: Availability Calculation Error in Topic Specification

**Claim:** The topic specification claims "99% → ~7 jam 18 menit/bulan" and "99.9% → ~43 menit/bulan" - need to verify against authoritative source.

**Evidence:** 
- Google Availability Table: 99% = 7.2 hours/month = 7 hours 12 minutes (not 7 hours 18 minutes)
- Google Availability Table: 99.9% = 43.2 minutes/month (matches ~43 minutes)
- 99.99% = 4.32 minutes/month = 4 minutes 19 seconds (topic says 4 menit 23 detik)

**Source:** Google SRE Book - Appendix A: Availability Table
**URL:** https://sre.google/sre-book/availability-table/
**Confidence:** HIGH
**Notes:** Topic specification has minor calculation errors for 99% and 99.99% targets.