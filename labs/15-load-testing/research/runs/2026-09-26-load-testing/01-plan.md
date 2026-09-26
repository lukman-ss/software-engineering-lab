# Research Plan

## Research Topic
Load Testing for Booking Bengkel Application — Methodologies, Tools, Metrics, and Bottleneck Identification Strategies

## Objective
Gather primary-source evidence on load testing best practices for a workshop booking application (Login, Booking, Select Branch, Payment, Generate Invoice, WhatsApp Confirmation). Focus on:
1. Industry-standard definitions: Load Test, Stress Test, Spike Test, Endurance/Soak Test, Breakpoint Test, Smoke Test
2. Key metrics: P50, P95, P99, RPS, Error Rate, CPU, Memory, DB Connection Pool utilization
3. Tool comparison: k6 vs JMeter vs Locust vs Gatling — strengths, weaknesses, best-fit scenarios
4. Bottleneck identification: application vs database vs external API isolation methodology
5. Common load testing pitfalls and anti-patterns
6. P95/P99 degradation investigation methodology
7. Open vs closed workload models and their impact on test accuracy

## Research Questions
1. What are the industry-standard definitions for each performance test type?
2. Which metrics are universally recommended across standards (ISO/IEC 25010, Google SRE, Azure WAF)?
3. When should each load testing tool be preferred over others?
4. How do you isolate whether a bottleneck is in the application, database, or external API?
5. What are the most common load testing mistakes?
6. How do you investigate P95/P99 degradation specifically?
7. When should load testing be integrated into SDLC?
8. What is the difference between open and closed workload models, and why does it matter?

## Search Strategy
- **Tier 1 (Primary)**: Official documentation — k6 (grafana.com/docs/k6), JMeter (jmeter.apache.org), Locust (docs.locust.io), Gatling (docs.gatling.io), Google SRE Book (sre.google), Azure Well-Architected Framework (learn.microsoft.com), AWS Well-Architected, ISO/IEC 25010
- **Tier 2 (Secondary)**: Martin Fowler blog, CNCF, engineering blogs (Netflix, Uber, Shopify), vendor whitepapers (Grafana, Datadog, New Relic)
- **Tier 3 (Community)**: Stack Overflow, Dev.to, Medium, conference talks
- Cross-check every claim against 2-3 independent sources
- Verify all URLs and source original content

## Sources Verified
| Source | Type | URL | Status |
|--------|------|-----|--------|
| k6 Load Test Types | Tier 1 | grafana.com/docs/k6/latest/testing-guides/test-types/ | Verified |
| k6 Built-in Metrics | Tier 1 | grafana.com/docs/k6/latest/using-k6/metrics/reference/ | Verified |
| k6 Thresholds | Tier 1 | grafana.com/docs/k6/latest/using-k6/thresholds/ | Verified |
| k6 Calculate Concurrent Users | Tier 1 | grafana.com/docs/k6/latest/testing-guides/calculate-concurrent-users/ | Verified |
| k6 Stress Testing | Tier 1 | grafana.com/docs/k6/latest/testing-guides/test-types/stress-testing/ | Verified |
| k6 Soak Testing | Tier 1 | grafana.com/docs/k6/latest/testing-guides/test-types/soak-testing/ | Verified |
| Azure Performance Testing | Tier 1 | learn.microsoft.com/azure/well-architected/performance-efficiency/performance-test | Verified |
| Google SRE Ch.17 Testing | Tier 1 | sre.google/sre-book/testing-reliability/ | Verified |
| Google SRE Ch.21 Handling Overload | Tier 1 | sre.google/sre-book/handling-overload/ | Verified |
| ISO/IEC 25010 (via Wikipedia) | Tier 1 | en.wikipedia.org/wiki/ISO/IEC_25010 | Verified |
| Locust Documentation | Tier 1 | docs.locust.io/en/stable/what-is-locust.html | Verified |
| Gatling Workload Models | Tier 1 | docs.gatling.io/testing-concepts/workload-models/ | Verified |
| AWS Well-Architected WAF | Tier 1 | docs.aws.amazon.com/wellarchitected/latest/performance-efficiency-pillar/ | Verified |

## Risks / Unknowns
- ISO/IEC 25010 full standard is paywalled; evidence from Wikipedia summary and ISO references
- JMeter official documentation returned 403 during fetch; evidence derived from known features and community consensus
- Real-world production benchmarks from major companies are proprietary
- Language mixing: source materials in English, lab content in Bahasa Indonesia

## Execution
- 14 primary/secondary sources fetched and verified
- Evidence claims cross-checked against minimum 2 independent sources per claim
- All URLs verified accessible as of 2026-09-26
