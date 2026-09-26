# 03 Claim Audit

## Claim 1
Claim: Load testing consists of four standardized primary test types: Average-load, Stress, Spike, and Soak/Endurance testing.
Location: `05-report.md` (Finding 1), `03-evidence.md` (Evidence 1)
Evidence Provided: Detailed definitions, duration profiles, ramp patterns from k6 docs and Azure Well-Architected Framework.
Source: Source 1, 2, 3, 4, 10, 15
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Fully supported across multiple Tier 1 authoritative sources (k6 docs, Azure Well-Architected, Google SRE).

## Claim 2
Claim: Key metrics to monitor during load tests are RED method (Rate, Errors, Duration) and 4 Golden Signals (Latency, Traffic, Errors, Saturation), coupled with backend resource metrics (CPU, Memory, DB connections, Queue length).
Location: `05-report.md` (Finding 2), `03-evidence.md` (Evidence 2, 6, 16, 22)
Evidence Provided: k6 metric definitions (`http_reqs`, `http_req_duration`, `http_req_failed`), Google SRE Chapter 6 (Four Golden Signals), Grafana RED blog, Azure performance testing guidance.
Source: Source 1, 3, 5, 11, 12, 15, 21
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Supported by vendor and standard industry literature.

## Claim 3
Claim: Percentiles (P95, P99) are essential; average response time masks tail latency and severe outliers.
Location: `05-report.md` (Finding 3), `03-evidence.md` (Evidence 7)
Evidence Provided: Google SRE tail latency analysis ("1% of requests might easily take 5 seconds"), k6 percentile aggregation documentation.
Source: Source 5, 6, 7, 12
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Strongly supported by Google SRE and k6 engineering documentation.

## Claim 4
Claim: Thresholds codify SLOs as automated pass/fail criteria (e.g., `p(95)<200`, `rate<0.01`).
Location: `05-report.md` (Finding 4), `03-evidence.md` (Evidence 8, 25)
Evidence Provided: k6 thresholds specification with exact aggregation syntax, Azure Well-Architected threshold recommendations.
Source: Source 6, 15, 16
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Directly supported by k6 and Azure docs.

## Claim 5
Claim: Tool comparison: k6 is JS-based and API-focused; Locust is Python-based using greenlets; JMeter is multi-protocol (HTTP, JDBC, JMS); Gatling offers async non-blocking execution.
Location: `05-report.md` (Finding 5), `03-evidence.md` (Evidence 12, 13, 14, 15)
Evidence Provided: k6 docs, Locust docs, Azure Load Testing documentation for JMeter, Gatling vendor page.
Source: Source 1, 7, 8, 9, 15, 16, 25
Source Actually Supports Claim: YES (k6, Locust, JMeter); PARTIAL (Gatling)
Classification: FACT / IMPLEMENTATION-SPECIFIC
Severity: MEDIUM
Notes: k6 and Locust are verified from primary docs. JMeter is verified via Azure docs. Gatling DSL and low-level async performance claims are only partially verified due to primary doc 403. Report properly notes this limitation.

## Claim 6
Claim: Bottleneck identification requires correlating API metrics with backend resource metrics (CPU, memory, DB connections, queues).
Location: `05-report.md` (Finding 6), `03-evidence.md` (Evidence 16)
Evidence Provided: Azure Well-Architected layer budgeting, Google SRE overload analysis, k6 soak testing guidance.
Source: Source 1, 11, 15
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Universal consensus across SRE and performance engineering documentation.

## Claim 7
Claim: External dependencies should be tested with real calls in sandbox/staging to measure actual latency, while high-volume stress tests should use mocks/stubs to avoid rate limits and cost.
Location: `05-report.md` (Finding 7), `03-evidence.md` (Evidence 17)
Evidence Provided: Azure Well-Architected recommendation ("include actual third-party API calls... mocking hides real-world performance problems") qualified by rate limiting and operational cost considerations.
Source: Source 15
Source Actually Supports Claim: YES
Classification: INTERPRETATION
Severity: LOW
Notes: Balanced interpretation preventing overgeneralization.

## Claim 8
Claim: Test data must be production-realistic in volume and diversity; small dummy datasets fail to trigger database indexing and lock contention bottlenecks.
Location: `05-report.md` (Finding 8), `03-evidence.md` (Evidence 18)
Evidence Provided: Azure Well-Architected guidelines on realistic data sets and edge cases.
Source: Source 15
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Directly supported by Azure architecture guidance.

## Claim 9
Claim: Load testing should start early and integrate into CI/CD pipelines to prevent performance regression.
Location: `05-report.md` (Finding 9), `03-evidence.md` (Evidence 10)
Evidence Provided: Azure deployment pipeline recommendations, Google SRE release engineering test gates.
Source: Source 7, 15, 17
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Strongly supported by cloud and SRE literature.

## Claim 10
Claim: Specific triage workflow for P95 spike at 800 VUs: correlate CPU/Memory/DB connection pool saturation to isolate the bottleneck layer.
Location: `03-evidence.md` (Evidence 29), `05-report.md` (Limitations item 6)
Evidence Provided: Synthesized reasoning from Evidence 7, 16, and 17.
Source: None single source; synthesized.
Source Actually Supports Claim: PARTIAL
Classification: HYPOTHESIS / INTERPRETATION
Severity: MEDIUM
Notes: The research report explicitly identifies this as synthesized best-practice guidance rather than an established single-source fact. Fully transparent.
