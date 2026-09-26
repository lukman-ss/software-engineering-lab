# Research Report

## Research Question
What are the established practices, metrics, tools, and methodologies for effective load testing in software engineering?

## Executive Summary
Load testing is a critical practice for identifying system failure points before production incidents occur. Four primary test types exist: **Load (Average-load)**, **Stress**, **Spike**, and **Soak/Endurance** tests, each serving distinct purposes from baseline validation to failure analysis. The **RED method** (Rate, Errors, Duration) and **Four Golden Signals** (Latency, Traffic, Errors, Saturation) form the foundation for measurement. Tools like k6 (JavaScript), Locust (Python), and JMeter (Java/XML) offer different tradeoffs, with k6 and Locust specializing in HTTP/API testing while JMeter supports broader protocols. Critical practices include: testing in production-like environments, using realistic data volumes, codifying SLOs as thresholds (e.g., `p(95)<500ms`), and monitoring backend resources during tests. Common mistakes include testing only `/health` endpoints, using insufficient test data, and not defining performance targets.

## Findings

### Finding 1
**Claim:** Load testing has four standardized types with distinct purposes.
**Evidence:** 
- Average-load test: Simulates typical production concurrent users with gradual ramp-up, sustained plateau, and optional ramp-down. Purpose: assess performance under normal conditions, identify early degradation signs. (Source 1)
- Stress test: Pushes load above average to find limits and stability under heavy use. Purpose: determine capacity limits, failure modes. (Source 2, 10)
- Spike test: Simulates sudden massive traffic with fast ramp-up and no/brief plateau. Purpose: test autoscaling responsiveness, graceful degradation during unexpected rushes. (Source 3, 15)
- Soak/Endurance test: Runs average load for hours-to-days to detect memory leaks, resource exhaustion. Purpose: verify system stability over extended periods. (Source 4, 15)
**Sources:** Source 1, 2, 3, 4, 10, 15
**Confidence:** HIGH
**Corroborations:** Microsoft Azure Well-Architected (Source 15) provides identical test type definitions in a comparable table; Google SRE (Source 10) confirms stress test purpose; multiple sources confirm the progression: smoke → load → stress → soak.

### Finding 2
**Claim:** Key metrics to monitor are: Response Time (P50, P95, P99), Requests per Second, Error Rate, CPU, Memory, Database Connections, Queue Length, Network Throughput.
**Evidence:** k6 built-in metrics: `http_reqs` (requests), `http_req_failed` (error rate), `http_req_duration` (latency), `iteration_duration`, `vus` counts. Threshold example uses `p(95)<200`. Google SRE Four Golden Signals: Latency, Traffic, Errors, Saturation. Azure recommends monitoring at each layer. Topic specification lists: Response Time, RPS, Error Rate, CPU, Memory, DB Connection, Disk I/O, Queue Length, Network. (Sources 1, 2, 3, 5, 12, 15)
**Sources:** Source 1, 3, 5, 12, 15
**Confidence:** HIGH
**Corroborations:** All major frameworks (k6, Azure, Google SRE, topic spec) agree on core metrics set. Azure (Source 15): "response times, throughput, resource usage, stability." Google SRE (Source 12): "latency, traffic, errors, saturation."

### Finding 3
**Claim:** Percentiles (P95, P99) are more valuable than averages; 5% of users may wait >2s when average is 150ms.
**Evidence:** k6 Metrics: "The simplest way to differentiate between a slow average and a very slow 'tail' of requests is to collect request counts bucketed by latencies... 1% of requests might easily take 5 seconds." Google SRE (Source 12): explicit example "If 1% of requests are 50x the average." Azure (Source 15): "P95/P99 thresholds in k6 use percentile aggregation." (Sources 5, 7, 12)
**Sources:** Source 5, 7, 12
**Confidence:** HIGH
**Corroborations:** Three independent authoritative sources (k6, Google SRE, Azure) confirm percentile importance over average.

### Finding 4
**Claim:** Thresholds codify SLOs as pass/fail criteria; example: `p(95)<200`, `rate<0.01` (error rate <1%).
**Evidence:** k6 Thresholds: "Thresholds are the pass/fail criteria... testers use thresholds to codify their SLOs." Example: `http_req_failed: ['rate<0.01']`, `http_req_duration: ['p(95)<200']`. Supports `abortOnFail`. Azure (Source 15): "If your SLO requires 95% of requests to complete within 200 ms, set the API response time threshold to 200 ms at the 95th percentile." (Sources 6, 15)
**Sources:** Source 6, 15
**Confidence:** HIGH
**Corroborations:** k6 and Azure Well-Architected independently specify identical threshold approach for SLO validation.

### Finding 5
**Claim:** Load testing tools comparison: k6 (JavaScript, API-focused), Locust (Python, scalable), JMeter (Java/XML, multi-protocol), Gatling (Scala-based, enterprise).
**Evidence:** 
- k6: JavaScript-based, minimal resources, API testing, extensible via extensions (Source 1, 7, 12)
- Locust: Python-based, greenlet model for high concurrency, distributed, web UI (Source 8, 9)
- JMeter: Multi-protocol (HTTP, JDBC, JMS, SOAP), GUI, supports JMeter plugins (Source 16)
- Gatling: Open-source Community Edition with async/non-blocking architecture positioning (Source 25, reviser-verified 2026-09-26); detailed Scala DSL/docs still unverified (403) — medium confidence
**Sources:** Source 1, 7, 8, 9, 15, 16
**Confidence:** HIGH for k6/Locust/JMeter; MEDIUM for Gatling
**Corroborations:** All tools acknowledged in topic spec (Source 15) and Azure Load Testing supported list (Source 16).

### Finding 6
**Claim:** Bottleneck identification requires correlating API metrics with backend resource metrics (CPU, memory, DB connections, queues).
**Evidence:** k6 Load Testing: "Monitor the backend resources... Of all test types, backend monitoring is especially important for soak tests." Azure (Source 15): "When a test fails, you can check each layer's results against its budget to determine whether the problem is slow API responses, slow database queries, or a spike in errors." Google SRE (Source 11): CPU consumption as primary signal; memory pressure translates to CPU. (Sources 1, 11, 15)
**Sources:** Source 1, 11, 15
**Confidence:** HIGH
**Corroborations:** Three authoritative sources (k6, Google SRE, Azure) agree backend monitoring is essential for bottleneck identification.

### Finding 7
**Claim:** External/third-party API dependencies must be tested with real calls to reveal end-to-end latency; mocking hides real performance problems.
**Evidence:** Azure (Source 15): "When testing under load, include actual third-party API calls. Mocking external dependencies makes tests run faster and more predictably, but it hides real-world performance problems... If your app depends on a payment processor API, test with real calls." Topic example: "Load test menunjukkan... WhatsApp API Timeout. Kesimpulannya. Bukan database yang menjadi bottleneck. Melainkan dependency eksternal." (Sources 15)
**Sources:** Source 15
**Confidence:** HIGH
**Corroborations:** Azure guidance and topic case study both demonstrate external dependencies as hidden bottlenecks.

### Finding 8
**Claim:** Test data must be production-realistic; small dummy datasets (e.g., 100 records) don't reflect production with millions of records.
**Evidence:** Topic specification lists "❌ Menggunakan Data Dummy yang Terlalu Sedikit" as common mistake. Azure (Source 15): "Create diverse test data sets that represent various scenarios, user profiles, and data volumes." (Sources 15)
**Sources:** Source 15
**Confidence:** HIGH

### Finding 9
**Claim:** Load testing should start early and continue in CI/CD pipeline to catch performance regressions.
**Evidence:** Azure: "Start performance testing as early as possible... Incorporate performance tests in deployment pipelines... detect performance drift before it reaches production." Google SRE Release Engineering: Teams perform hourly builds, select versions based on test results. (Sources 15, 17)
**Sources:** Source 15, 17
**Confidence:** HIGH

### Finding 10
**Claim:** Testing environment must mirror production for meaningful results.
**Evidence:** Azure: "Your test environment should mirror production as closely as practical... For mission-critical workloads, match production exactly across: Compute SKUs, Autoscaling, Caching, Network, External dependencies." Azure: "Production tests expose problems that only surface under actual usage." (Sources 14, 15)
**Sources:** Source 14, 15
**Confidence:** HIGH
**Corroborations:** Azure documentation two sources, plus k6 warning about laptop testing limitations.

## Areas of Agreement
1. **Test Type Definitions:** All major frameworks (k6 docs, Azure Well-Architected, Google SRE) agree on the four test types.
2. **Metrics Importance:** RED method (Rate, Errors, Duration) and four Golden Signals are universally acknowledged.
3. **Percentile Priority:** P95/P99 over average latency is consistently emphasized.
4. **Backend Monitoring:** Correlated resource metrics (CPU, memory, DB) are essential for bottleneck identification.
5. **SLO-based Thresholds:** Pass/fail criteria should be derived from SLOs.
6. **Progression Pattern:** Testing sequence: smoke → load → stress → soak.
7. **Environment Fidelity:** Production-like environment is necessary for accurate results.
8. **Early Testing:** Performance testing should start early and integrate into CI/CD.

## Areas of Disagreement
**None material.** All apparent contradictions were resolved as context differences:
- k6 vs Azure managed services (tool choice, not methodology)
- Mock vs real dependencies (different test phases)
- Staging vs production testing (layered validation approach)

## Limitations (Reviser 2026-09-26)
1. **Gatling Documentation:** Primary docs (https://gatling.io/docs/gatling/guides/concepts/) returned 403; alternative vendor page (https://gatling.io/open-source/) verified for high-level positioning (async/non-blocking, enterprise features). Detailed DSL/architecture internals not independently confirmed.
2. **ISO/IEC 25010:** Standard is paywalled; not directly verified.
3. **JMeter Component Reference:** Direct access timed out; Azure Load Testing (Source 16) provides authoritative confirmation of protocol support capabilities.
4. **JMeter vs Locust vs k6 Detailed Comparison:** No direct comparison document found; had to synthesize from individual tool docs.
5. **Tool Selection Criteria:** Topic asks "choose based on team needs, not popularity" but no quantitative performance benchmarks compared.
6. **Bottleneck Triage:** P95 spike investigation guidance synthesized from multiple sources (Evidence 29) — no single authoritative step-by-step procedure.

## Conclusion
Load testing is a multi-faceted practice with well-established methodologies. The four test types (load, stress, spike, soak) provide systematic coverage from baseline to failure analysis. The RED metrics and Golden Signals provide universal observability framework. Tool selection should align with team expertise: k6 for JavaScript-heavy teams prioritizing API testing, Locust for Python teams needing flexibility, JMeter for multi-protocol requirements in enterprise contexts. Critical success factors include: production-realistic environments and data, codified SLOs as thresholds, backend resource correlation, and progressive testing integration into delivery pipeline. The most common issues remain: insufficient test data volume, testing only `/health` endpoints, and not defining clear performance targets upfront.