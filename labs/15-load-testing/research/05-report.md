# Research Report

## Research Question
How should an engineering team approach load testing, measure system limits, analyze percentiles, and plan tests for high-criticality systems like "Booking Bengkel"?

## Executive Summary
Load testing determines the boundaries and degradation patterns of a system under load rather than merely proving it is fast. Effective load testing incorporates various patterns (smoke, load, stress, spike, soak), monitors full-stack metrics (client-side latency and server-side resource saturation), and evaluates percentiles (P95, P99) to detect latency outliers that simple averages hide. For transactional architectures like "Booking Bengkel", focus should target state-mutating endpoints, starting with low-concurrency smoke tests, ramping up incrementally, and isolating components using observability.

## Findings

### Finding 1
Claim: Average response time is misleading; percentiles (P95, P99) are essential for discovering latency issues in real workloads.
Evidence: Outliers are diluted by averages. In a dataset where 95% of users receive fast responses (e.g., 100ms) but 5% wait 3000ms, the average remains acceptable, while P95 immediately flags the experience of the slower subset.
Sources: Grafana k6 Documentation, Industry Best Practices
Confidence: HIGH

### Finding 2
Claim: Incremental testing starting with smoke tests prevents invalidating test results.
Evidence: Grafana k6 docs note that smoke tests with low VUs (e.g., 2–5 VUs) for seconds or minutes verify script integrity and baseline health before stress or peak tests are run.
Sources: Types of load testing (Grafana Labs)
URL: https://k6.io/docs/test-types/
Confidence: HIGH

### Finding 3
Claim: Comprehensive bottleneck detection requires correlating client-side latency with server-side resource metrics.
Evidence: Microsoft Azure Load Testing highlights that test results must pair engine-side response times and RPS with server-side metrics (CPU, memory, database reads/writes, connection pools) to distinguish network, compute, database, and downstream API bottlenecks.
Sources: What is Azure Load Testing? (Microsoft Learn)
URL: https://learn.microsoft.com/en-us/azure/load-testing/overview-what-is-azure-load-testing
Confidence: HIGH

### Finding 4 (Booking Bengkel Plan Analysis)
Claim: Critical endpoints and failure modes must be isolated systematically.
Evidence:
1. **Critical Endpoints:** State-changing/transactional endpoints (`POST /booking`, `POST /payment`, `POST /invoice`) and authentication (`POST /login`) take priority over read-only lookups (`GET /branches`).
2. **Initial Virtual Users:** 5 to 10 VUs for smoke testing, scaling to normal expected traffic (e.g., 50–100 VUs) before stress testing.
3. **Metrics Monitored:** RPS, P95/P99 Latency, Error Rate (4xx/5xx), CPU, Memory, DB Connection Pool saturation, and external API latency.
4. **Isolating Bottlenecks:**
   - *App level:* High CPU/memory or thread pool exhaustion while DB/external latency is low.
   - *Database level:* DB CPU spikes, slow query logs, connection pool reaching 100% capacity, or lock wait timeouts.
   - *Third-party API level:* App threads blocked in external HTTP calls (e.g., WhatsApp webhook / payment gateway), high external response time while local CPU/DB are idle.
5. **Investigating P95 jump (300ms to 2.5s at 800 VUs):**
   - Check APM/traces for where time is spent.
   - Inspect DB connection pool metrics and active locks.
   - Check if downstream APIs (Payment/WhatsApp) are throttling or timing out.
Sources: Grafana Labs & Microsoft Learn Best Practices
Confidence: HIGH

## Areas of Agreement
- Average metrics conceal performance degradations.
- Test environments should closely mimic production resources and configurations.
- Tests should progress sequentially from low to high loads (Smoke -> Average Load -> Stress/Spike -> Soak).

## Areas of Disagreement
- None identified in the reviewed primary documentation.

## Limitations
- Specific external API quotas and production hardware configurations for the "Booking Bengkel" system were unspecified, requiring standard architectural sizing assumptions.

## Conclusion
Load testing provides actionable insights on when, where, and why a system degrades. Testing must focus on realistic flows, rely on P95/P99 latency, and correlate client metrics with server-side resource constraints to prevent production outages.
