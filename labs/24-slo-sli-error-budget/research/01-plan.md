# Research Plan

## Research Topic
SLO, SLI & Error Budget — Service Level Objectives, Service Level Indicators, and Error Budget concepts in Site Reliability Engineering

## Objective
Investigate and collect authoritative evidence on SLO/SLI/Error Budget concepts, their practical application, industry standards, calculation methods, and best practices for reliability engineering decision-making.

## Research Questions

1. **Definitions & Core Concepts**
   - What are the canonical definitions of SLI, SLO, and Error Budget?
   - How do they relate to SLA (Service Level Agreement)?
   - What are the mathematical formulations?

2. **Industry Standards & Frameworks**
   - What does Google's SRE book / workbook say?
   - What are the OpenSLO / OpenTelemetry standards?
   - What do CNCF / TAG Observability recommend?

3. **SLI Selection & Design**
   - What are the recommended SLI types (availability, latency, quality, freshness, correctness, durability)?
   - How to choose SLIs based on user experience vs infrastructure metrics?
   - What are anti-patterns in SLI selection?

4. **SLO Target Setting**
   - How to determine appropriate SLO targets based on business impact?
   - What are common SLO tiers (99%, 99.9%, 99.99%, 99.999%) and their implications?
   - How to handle different SLOs for different endpoints/criticality levels?

5. **Error Budget Calculation & Burn Rate**
   - What are the standard formulas for error budget calculation?
   - How does burn rate alerting work (multi-window, multi-burn-rate)?
   - What are the recommended alerting thresholds?

6. **Practical Implementation**
   - How to implement SLOs in practice (Prometheus, Grafana, Datadog, etc.)?
   - What are the common pitfalls and how to avoid them?
   - How to use error budgets for release decisions and prioritization?

7. **Case Studies & Real-World Applications**
   - How do companies (Google, Netflix, Slack, etc.) implement these concepts?
   - What are the measurable outcomes?

## Search Strategy

1. Primary sources: Google SRE Book, SRE Workbook, official documentation
2. Standards: OpenSLO specification, OpenTelemetry semantic conventions
3. Industry publications: ACM Queue, USENIX ;login:, conferences (SREcon, Velocity)
4. Vendor documentation: Prometheus, Grafana, Datadog, Honeycomb, Nobl9
5. Technical blogs from recognized experts

## Expected Primary Sources

- **Tier 1**: Google SRE Book (2016), SRE Workbook (2018), OpenSLO spec, OpenTelemetry docs, CNCF TAG Observability whitepapers
- **Tier 2**: SREcon presentations, USENIX papers, major tech company engineering blogs
- **Tier 3**: Community discussions, personal blogs (for discovery only)

## Risks / Unknowns

- Some Google SRE content may be dated (2016-2018); need to verify current practices
- OpenSLO is relatively new (2021+); adoption may be limited
- Vendor-specific implementations may differ from standards
- Indonesian language context for the lab - need to verify terminology alignment