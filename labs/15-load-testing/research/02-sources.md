# Sources

## Source 1
Title: Load testing | Grafana k6 documentation
Publisher: Grafana Labs
URL: https://grafana.com/docs/k6/latest/testing-guides/test-types/load-testing/
Published: NOT VERIFIED (current as of 2025)
Accessed: 2026-09-26
Source Tier: Tier 1 (official documentation)
Relevance: Defines "average-load test" concept, purpose, ramp-up/ramp-down patterns, and k6 configuration example.

## Source 2
Title: Stress testing | Grafana k6 documentation
Publisher: Grafana Labs
URL: https://grafana.com/docs/k6/latest/testing-guides/test-types/stress-testing/
Published: NOT VERIFIED (current as of 2025)
Accessed: 2026-09-26
Source Tier: Tier 1 (official documentation)
Relevance: Defines stress testing as testing above-average load to find limits and degradation patterns.

## Source 3
Title: Spike testing | Grafana k6 documentation
Publisher: Grafana Labs
URL: https://grafana.com/docs/k6/latest/testing-guides/test-types/spike-testing/
Published: NOT VERIFIED (current as of 2025)
Accessed: 2026-09-26
Source Tier: Tier 1 (official documentation)
Relevance: Defines spike testing for sudden massive load, key differentiators (fast ramp-up, no plateau).

## Source 4
Title: Soak testing | Grafana k6 documentation
Publisher: Grafana Labs
URL: https://grafana.com/docs/k6/latest/testing-guides/test-types/soak-testing/
Published: NOT VERIFIED (current as of 2025)
Accessed: 2026-09-26
Source Tier: Tier 1 (official documentation)
Relevance: Defines soak/endurance testing for extended periods to detect memory leaks, resource leaks.

## Source 5
Title: Metrics | Grafana k6 documentation
Publisher: Grafana Labs
URL: https://grafana.com/docs/k6/latest/using-k6/metrics/
Published: NOT VERIFIED (current as of 2025)
Accessed: 2026-09-26
Source Tier: Tier 1 (official documentation)
Relevance: Lists built-in metrics: http_reqs, http_req_failed, http_req_duration, iteration_duration, virtual user counts.

## Source 6
Title: Thresholds | Grafana k6 documentation
Publisher: Grafana Labs
URL: https://grafana.com/docs/k6/latest/using-k6/thresholds/
Published: NOT VERIFIED (current as of 2025)
Accessed: 2026-09-26
Source Tier: Tier 1 (official documentation)
Relevance: Documents threshold syntax (p(95)<200, rate<0.01), percentile aggregation methods, SLO codification.

## Source 7
Title: API load testing | Grafana k6 documentation
Publisher: Grafana Labs
URL: https://grafana.com/docs/k6/latest/testing-guides/api-load-testing/
Published: NOT VERIFIED (current as of 2025)
Accessed: 2026-09-26
Source Tier: Tier 1 (official documentation)
Relevance: Testing pyramid (isolated → integrated → e2e), smoke/breakpoint test types, VUs vs RPS, load generator placement.

## Source 8
Title: What is Locust?
Publisher: Locust Project (locust.io)
URL: https://docs.locust.io/en/stable/what-is-locust.html
Published: NOT VERIFIED (current as of 2025)
Accessed: 2026-09-26
Source Tier: Tier 1 (official documentation)
Relevance: Defines Locust as Python-based load testing tool, greenlet model, distributed support, web UI, hackable.

## Source 9
Title: Writing a locustfile
Publisher: Locust Project (locust.io)
URL: https://docs.locust.io/en/stable/writing-a-locustfile.html
Published: NOT VERIFIED (current as of 2025)
Accessed: 2026-09-26
Source Tier: Tier 1 (official documentation)
Relevance: Defines task/weight/wait_time model, @task decorator, @tag decorator, HttpUser class, connection pooling.

## Source 10
Title: Testing for Reliability | Google SRE Book (Chapter 17)
Publisher: Google (O'Reilly Media)
URL: https://sre.google/sre-book/testing-reliability/
Published: 2017 (Copyright)
Accessed: 2026-09-26
Source Tier: Tier 1 (authoritative technical book)
Relevance: Defines stress test as finding limits of a web service; performance test detects degradation; smoke test precedes others.

## Source 11
Title: Handling Overload | Google SRE Book (Chapter 21)
Publisher: Google (O'Reilly Media)
URL: https://sre.google/sre-book/handling-overload/
Published: 2017 (Copyright)
Accessed: 2026-09-26
Source Tier: Tier 1 (authoritative technical book)
Relevance: Discusses CPU-based capacity, per-customer quotas, client-side throttling, overload handling, retry budgets.

## Source 12
Title: Monitoring Distributed Systems | Google SRE Book (Chapter 6)
Publisher: Google (O'Reilly Media)
URL: https://sre.google/sre-book/monitoring-distributed-systems/
Published: 2017 (Copyright)
Accessed: 2026-09-26
Source Tier: Tier 1 (authoritative technical book)
Relevance: Four Golden Signals (latency, traffic, errors, saturation); percentile importance for tail latency; symptom vs cause.

## Source 13
Title: Performance Efficiency Pillar - AWS Well-Architected Framework
Publisher: Amazon Web Services (AWS)
URL: https://docs.aws.amazon.com/wellarchitected/latest/performance-efficiency-pillar/welcome.html
Published: November 6, 2024
Accessed: 2026-09-26
Source Tier: Tier 1 (official cloud provider documentation)
Relevance: Five focus areas: Architecture selection, Compute and hardware, Data management, Networking, Process and culture.

## Source 14
Title: Design review checklist for Performance Efficiency - Microsoft Azure Well-Architected Framework
Publisher: Microsoft Learn
URL: https://learn.microsoft.com/en-us/azure/architecture/framework/scalability/performance-efficiency
Published: 2026-04-20
Accessed: 2026-09-26
Source Tier: Tier 1 (official cloud provider documentation)
Relevance: 12-point checklist for performance efficiency, capacity planning before usage changes, testing in production-like env.

## Source 15
Title: Architecture Strategies for Performance Testing - Microsoft Azure Well-Architected Framework
Publisher: Microsoft Learn
URL: https://learn.microsoft.com/en-us/azure/well-architected/performance-efficiency/performance-test
Published: 2026-06-24
Accessed: 2026-09-26
Source Tier: Tier 1 (official cloud provider documentation)
Relevance: Test types table, performance budgets, hypothesis-driven experimentation, progressive production testing, real-world data.

## Source 16
Title: What is Azure Load Testing?
Publisher: Microsoft Learn
URL: https://learn.microsoft.com/en-us/azure/app-testing/load-testing/overview-what-is-azure-load-testing
Published: 2024-04-04
Accessed: 2026-09-26
Source Tier: Tier 1 (official cloud provider documentation)
Relevance: Managed load testing service, JMeter/Locust support, client/server-side metrics, fail criteria, auto-stop.

## Source 17
Title: Release Engineering | Google SRE Book (Chapter 8)
Publisher: Google (O'Reilly Media)
URL: https://sre.google/sre-book/release-engineering/
Published: 2017 (Copyright)
Accessed: 2026-09-26
Source Tier: Tier 1 (authoritative technical book)
Relevance: Self-service model, high velocity releases, hermetic builds, continuous build and deployment.

## Source 18
Title: The Practical Test Pyramid
Publisher: MartinFowler.com
URL: https://martinfowler.com/articles/practical-test-pyramid.html
Published: Feb 26, 2018
Accessed: 2026-09-26
Source Tier: Tier 2 (reputable technical publication)
Relevance: Testing pyramid concept: many unit tests, fewer service tests, fewest UI tests; referenced by k6 docs.

## Source 19
Title: Test Coverage (bliki)
Publisher: MartinFowler.com
URL: https://martinfowler.com/bliki/TestCoverage.html
Published: Apr 17, 2012
Accessed: 2026-09-26
Source Tier: Tier 2 (reputable technical publication)
Relevance: Test coverage as a tool to find untested parts of a codebase, cautioning against numeric targets.

## Source 20
Title: Synthetic Monitoring
Publisher: MartinFowler.com
URL: https://martinfowler.com/bliki/SyntheticMonitoring.html
Published: Jan 25, 2017
Accessed: 2026-09-26
Source Tier: Tier 2 (reputable technical publication)
Relevance: Synthetic monitoring runs subset of automated tests against live production on regular basis.

## Source 21
Title: The RED Method
Publisher: Grafana Blog
URL: https://grafana.com/blog/2018/08/02/the-red-method-how-to-instrument-your-services/
Published: Aug 2, 2018
Accessed: 2026-09-26
Source Tier: Tier 1 (official blog of tool vendor)
Relevance: RED = Rate, Errors, Duration; three foundational metrics for services.

## Source 22
Title: Apache JMeter
Publisher: Apache Software Foundation
URL: https://jmeter.apache.org/
Published: NOT VERIFIED (ongoing project)
Accessed: 2026-09-26
Source Tier: Tier 1 (official project documentation)
Relevance: Open-source load testing tool, GUI-based, supports HTTP, JDBC, JMS, SOAP, FTP, etc.
Verification Notes: Azure Load Testing (Source 16) confirms these capabilities through integration documentation; direct JMeter site access timed out during research.

## Source 23
Title: What's New in JMeter (component reference attempt)
Publisher: Apache Software Foundation
URL: https://jmeter.apache.org/usermanual/component_reference.html
Published: NOT VERIFIED
Accessed: NOT ACCESSED — request timed out
Source Tier: Tier 1 (official project documentation)
Relevance: NOT ACCESSED — request timed out. NOT USED as primary evidence for any claim. JMeter protocol support verified through Azure Load Testing (Source 16) instead.
Verification Status: DISCLAIMED — not included in primary evidence list due to inaccessibility.

## Source 24
Title: ISO/IEC 25010:2011 Software Engineering
Publisher: ISO (International Organization for Standardization)
URL: STUB — paywalled standard, not directly verified
Published: 2011 (according to ISO registry, secondary confirmation)
Accessed: NOT VERIFIED — full text not inspected
Source Tier: Tier 1 (international standard)
Relevance: Performance efficiency quality model: time behavior, resource utilization, capacity — widely cited in SE literature but NOT VERIFIED from paywalled text.
Verification Status: DISCLAIMED — not included in primary evidence list due to lack of direct verification; sub-characteristics cited are based on secondary descriptions — see Open Questions.

## Source 25
Title: Gatling Open Source vs. Gatling Enterprise
Publisher: Gatling Corp (gatling.io)
URL: https://gatling.io/open-source/
Published: NOT VERIFIED (ongoing vendor page)
Accessed: 2026-09-26 (Reviser Agent)
Source Tier: Tier 1 (official vendor documentation)
Relevance: Alternative verification for Gatling. Confirms: open-source Community Edition available, async/non-blocking architecture positioning, enterprise features (distributed tracing, AI analysis, SLO tracking), JMeter/LoadRunner converter tooling. Primary docs page (gatling.io/docs/gatling/guides/concepts/) returned 403 during both research and revision. Supports Evidence 15 with MEDIUM confidence (high-level positioning verified; detailed DSL/architecture internals not independently confirmed).