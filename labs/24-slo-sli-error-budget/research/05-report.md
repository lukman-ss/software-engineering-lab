# Research Report

## Research Question

How are SLO (Service Level Objective), SLI (Service Level Indicator), and Error Budget defined and applied in practice by authoritative Site Reliability Engineering sources, and what evidence supports the formulas and guidance used in labs/24-slo-sli-error-budget?

## Executive Summary

The concepts of SLI, SLO, and Error Budget are well-established in the SRE discipline. Google's own SRE Book and SRE Workbook are the most authoritative sources, corroborated by vendor documentation (Datadog) and standards (OpenSLO). All key claims in the lab topic material are supported by primary sources. The main nuances are: (1) time-window choice (rolling vs. calendar) is a documented open decision; (2) burn-rate alert parameters are starting points requiring tuning; (3) downtime-per-month figures depend on whether a 30-day or 30.44-day month is assumed (minor). No material factual errors were found in the lab material, and all core formulas are verified against at least two independent authoritative sources.

## Findings

### Finding 1: Canonical Definitions (SLI / SLO / SLA / Error Budget)

Claim: An **SLI** is a quantitative measure of service level provided; an **SLO** is a target for that measure; an **SLI** is a target value or range; an **SLA** is an agreement with consequences; an **error budget** is the allowable failure rate derived from an SLO.

Evidence:
- From Google SRE Book Ch.4: "An SLI is a service level *indicator*—a carefully defined quantitative measure of some aspect of the level of service that is provided." and "An SLO is a *service level objective*: a target value or range of values for a service level that is measured by an SLI."
- "SLAs are service level *agreements*: an explicit or implicit contract with your users that includes consequences of meeting (or missing) the SLOs they contain."
- Google SRE Workbook Ch.2: "an SLI ... should be treated as the ratio of two numbers: the number of good events divided by the total number of events." and "The error budget gives the number of allowed bad events."

Sources:
- https://sre.google/sre-book/service-level-objectives/
- https://sre.google/workbook/implementing-slos/

Confidence: HIGH (multiple authoritative sources agree)

### Finding 2: SLI Types and Selection Criteria

Claim: Common SLI types include availability, latency, quality, freshness, correctness, coverage, and durability; SLIs should be user-centric, not infrastructure metrics; percentiles (P95/P99) are preferred over averages.

Evidence:
- Google SRE Book Ch.4 lists request-driven (availability/latency/throughput), storage (latency/availability/durability), big data (throughput/end-to-end latency), and universal correctness SLIs.
- "Most metrics are better thought of as *distributions* rather than averages... A simple average can obscure these tail latencies."
- SRE Workbook Ch.2 Table 2-1: availability, latency, quality for request-driven; freshness, correctness, coverage for pipeline; durability for storage.
- "You should not use CPU < 80% as an SLO user metric"—this guidance is present in SRE Book (CPU/RAM as diagnostic signals), corroborated by SRE Workbook.

Sources:
- https://sre.google/sre-book/service-level-objectives/
- https://sre.google/workbook/implementing-slos/
- https://prometheus.io/docs/practices/alerting/

Confidence: HIGH

### Finding 3: Error Budget Calculation and the 100% - SLO Formula

Claim: Error budget = 100% - SLO target. The allowed failure count = total requests × (error budget %). Datadog adds an explicit formula: error budget remaining = 100 × (current - target)/(100 - target).

Evidence:
- SRE Workbook: "if you have a 99.9% success ratio SLO, then a service that receives 3 million requests over a four-week period had a budget of 3,000 (0.1%) errors."
- Embracing Risk: "a quarterly error budget based on the service's SLO" — 99.999% ⇒ 0.001% budget.
- Datadog: error budget remaining formula = 100 × (current status - target)/(100 - target).

Sources:
- https://sre.google/workbook/implementing-slos/
- https://sre.google/sre-book/embracing-risk/
- https://docs.datadoghq.com/service_level_objectives/

Confidence: HIGH

Cross-check: Lab exercise example — Payment Webhook 200,000 requests, SLO 99.99%. Calculation by hand: 200,000 × (100 - 99.99)% = 200,000 × 0.0001 = 20 allowed failures. 15 failures consumed = 15/20 = 75% of budget; 25% remaining. Matches lab expectation.

### Finding 4: Why Not 100%? Cost/Benefit of Reliability "Nines"

Claim: Each additional nine of availability costs disproportionately more while marginal utility to users approaches zero.

Evidence:
- SRE Book Ch.3: "increasing reliability is worse for a service... cost does not increase linearly as reliability increments—an incremental improvement in reliability may cost 100x more."
- "as you go from 99% to 99.9% to 99.99% reliability, each extra nine comes at an increased cost, but the marginal utility to your customers steadily approaches zero."
- "100% reliability means you can never update or improve your service."
- Availability Table (Appendix A): per-year downtime for each nines level.

Sources:
- https://sre.google/sre-book/embracing-risk/
- https://sre.google/sre-book/availability-table/

Confidence: HIGH

### Finding 5: Availability "Nines" Downtime Table

Claim: 99% ≈ 7.2h/month, 99.9% ≈ 43.2m/month, 99.99% ≈ 4.32m/month; 99.999% ≈ 25.9s/month.

Evidence:
- Google SRE Book Appendix A Table 1-1:
  - 99% → 7.2 hours/month, 14.4 min/day
  - 99.9% → 43.2 minutes/month, 1.44 min/day
  - 99.99% → 4.32 minutes/month, 8.64 sec/day
  - 99.999% → 25.9 seconds/month, 0.87 sec/day

Source: https://sre.google/sre-book/availability-table/

Confidence: HIGH

### Finding 6: Burn Rate Alerting and Multiwindow Configuration

Claim: Burn rate = speed of error budget consumption; recommended alerting uses multiwindow, multi-burn-rate with specific parameters.

Evidence:
- SRE Workbook Ch.5: "Burn rate is how fast, relative to the SLO, the service consumes the error budget. With an SLO of 99.9% over 30 days, a constant 0.1% error rate uses exactly all of the error budget: a burn rate of 1."
- Table 5-8: Page on 2% budget in 1h (burn rate 14.4); Page on 5% in 6h (burn rate 6); Ticket on 10% in 3d (burn rate 1).
- Multiwindow technique: alert fires when both long window (e.g., 1h) and short window (5m) exceed threshold.

Sources:
- https://sre.google/workbook/alerting-on-slos/
- https://docs.datadoghq.com/service_level_objectives/ (burn rate indicator: red if >6 in 2h; yellow if 1–6 in 2h)

Confidence: HIGH (core concept); MEDIUM (vendor-specific thresholds vary)

### Finding 7: Real-World SLO Adoption (Evernote, Home Depot)

Claim: Organizations adopt SLOs to align product and operations teams; The Home Depot scaled SLOs to 800 services in < 1 year using a VALET framework (Volume, Availability, Latency, Errors, Tickets).

Evidence:
- Evernote: "We introduced SLOs... 99.95% uptime measured over a monthly window, set for certain services and methods."
- "After introducing SLOs, the relationship between our operations and development teams has subtly but markedly improved."
- Home Depot: "we were tracking SLOs for 800 services, with about 50 new services per month being registered with VALET."

Sources:
- https://sre.google/workbook/slo-engineering-case-studies/

Confidence: HIGH

### Finding 8: Error Budget Policy as Release Decision Tool

Claim: When error budget is exhausted, the policy is to pause releases/ prioritize reliability work; the error budget makes reliability a decision metric, not just a dashboard number.

Evidence:
- SRE Book Ch.4: "The SLO violation rate can be compared against the error budget... with the gap used as an input to the process that decides when to roll out new releases."
- SRE Workbook Ch.2: Error budget policy actions: dev team prioritizes reliability bugs; production freeze; or full focus on reliability until re-budgeted.
- "As long as the uptime measured is above the SLO—in other words, as long as there is error budget remaining—new releases can be pushed."

Sources:
- https://sre.google/sre-book/service-level-objectives/
- https://sre.google/workbook/implementing-slos/

Confidence: HIGH

## Areas of Agreement

1. SLI = ratio-based measure (good/total events) preferred for tooling consistency and intuitive 0-100% scale.
2. SLO must be set below 100% to create a meaningful error budget.
3. Percentiles (P95/P99) preferred over averages for latency measurement.
4. User-centric SLIs preferred over infrastructure metrics (CPU/RAM) for SLOs.
5. Error budget serves as the alignment mechanism between product (velocity) and SRE (reliability).
6. Not all endpoints should have the same SLO—criticality/business impact should drive differentiation.
7. SLOs require multiple dimensions (availability + latency, not just availability).
8. OpenSLO standard exists for vendor-neutral, declarative SLO-as-code.

## Areas of Disagreement

1. Time window: rolling (Google's recommendation, 4 weeks) vs. calendar month (Evernote's choice) vs. undecided (Home Depot). No authoritative resolution.
2. Burn rate threshold parameters: differ slightly between Google SRE Workbook (multi-window) and Datadog's single-2-hour-window indicator. Both agree burn rate is the right concept.
3. Month-length assumption in downtime calculations: 30-day (Google table) vs. 30.44-day (lab material) — affects 99% and 99.99% figures by ~1-2%.

## Limitations

- Google's internal practices (Borgmon) may not map 1:1 to open-source tooling (Prometheus, Grafana).
- OpenSLO specification adoption is documented primarily on its own site; independent industry adoption numbers are not verified.
- Vendor documentation (Datadog, Prometheus) reflects implementation choices, not universal standards.
- Lab material uses 30.44-day months implicitly; this assumption is not stated.
- No Indonesian-language SRE sources were sought; definitions rely on English authoritative sources.

## Conclusion

The lab topic material's definitions and formulas are accurate and well-supported. The error budget formula (100% - SLO), downtime-per-nines table, burn rate concept, and multi-dimensional SLO approach all have HIGH confidence from primary sources. The lab's downtime figures (7h18m, 4m23s) are arithmetically correct under an average-month assumption, though slightly differ from Google's 30-day convention. Implementation recommendations (rolling 4-week window, specific burn rate thresholds) are starting points requiring tuning to organizational context. The lab's central thesis—SLOs make reliability a number that can drive decisions—is fully validated by primary sources.