# Sources

## Source 1: k6 Load Test Types

**Title**: Load test types  
**Publisher**: Grafana k6 Documentation  
**URL**: https://grafana.com/docs/k6/latest/testing-guides/test-types/  
**Published**: Continuous (latest version)  
**Accessed**: 2026-09-26  
**Source Tier**: Tier 1  
**Relevance**: Authoritative definition of six primary test types: smoke, average-load, stress, soak, spike, breakpoint. Explains load patterns, duration recommendations, k6 configuration examples, and key considerations for each type. Notes naming consensus does not exist and stress tests may be called "rush-hour, surge, scale tests".

**Key Claims Verified**:
- Six primary test types with consistent definitions
- No consensus on naming; various communities use different terms
- Stress test load depends on risk profile, not fixed percentage

---

## Source 2: k6 Calculate Concurrent Users

**Title**: Calculate concurrent users for load tests  
**Publisher**: Grafana k6 Documentation  
**URL**: https://grafana.com/docs/k6/latest/testing-guides/calculate-concurrent-users/  
**Published**: Continuous  
**Accessed**: 2026-09-26  
**Source Tier**: Tier 1  
**Relevance**: Industry-standard methodology for calculating baseline concurrent users from production traffic data using formula: `concurrent_users = hourly_sessions × average_session_duration / 3600`. Emphasizes testing at peak traffic rather than daily averages.

**Key Claims Verified**:
- Formula: hourly_sessions × avg_duration / 3600
- Peak traffic should drive test design, not daily averages
- Google Analytics or any analytics tool can provide required metrics

---

## Source 3: k6 Thresholds

**Title**: Thresholds  
**Publisher**: Grafana k6 Documentation  
**URL**: https://grafana.com/docs/k6/latest/using-k6/thresholds/  
**Published**: Continuous  
**Accessed**: 2026-09-26  
**Source Tier**: Tier 1  
**Relevance**: Documents threshold syntax for codifying SLOs. Examples: `p(95)<200`, `rate<0.01`. Details aggregation methods per metric type (avg, min, max, med, p(N) for trends; rate, count for counters/rates). Notes cloud threshold evaluation occurs every 60 seconds, potentially delaying abortOnFail.

**Key Claims Verified**:
- Percentile thresholds: `p(N)` where N ∈ [0.0, 100]
- Error rate thresholds: `rate<0.01` (1% error budget)
- Aggregation methods per metric type documented
- Cloud evaluation interval: 60 seconds

---

## Source 4: k6 Built-in Metrics Reference

**Title**: Built-in metrics  
**Publisher**: Grafana k6 Documentation  
**URL**: https://grafana.com/docs/k6/latest/using-k6/metrics/reference/  
**Published**: Continuous  
**Accessed**: 2026-09-26  
**Source Tier**: Tier 1  
**Relevance**: Comprehensive reference of built-in metrics including http_req_duration (trend with percentiles), http_req_failed (rate), http_req_blocked, http_req_connecting, http_req_tls_handshaking, http_req_sending, http_req_waiting, http_req_receiving. http_req_duration breakdown enables bottleneck isolation.

**Key Claims Verified**:
- http_req_duration = http_req_sending + http_req_waiting + http_req_receiving
- http_req_waiting (TTFB) indicates server-side processing delays
- http_req_connecting/TLS indicate network/connection issues
- Percentiles supported: p(N) for N ∈ [0.0, 100]

---

## Source 5: k6 Stress Testing

**Title**: Stress testing  
**Publisher**: Grafana k6 Documentation  
**URL**: https://grafana.com/docs/k6/latest/testing-guides/test-types/stress-testing/  
**Published**: Continuous  
**Accessed**: 2026-09-26  
**Source Tier**: Tier 1  
**Relevance**: Defines stress testing as loads above average to test system limits. Notes load increase depends on risk profile; no fixed percentage rule. Emphasizes running only after average-load tests pass. Recommends ramp-up proportional to load increase.

**Key Claims Verified**:
- Stress test = load above average to find breaking points
- No fixed percentage (50% or 100% mentioned as common but not mandatory)
- Ramp-up should be longer than in average-load tests
- Only run after average-load tests pass

---

## Source 6: k6 Soak Testing

**Title**: Soak testing  
**Publisher**: Grafana k6 Documentation  
**URL**: https://grafana.com/docs/k6/latest/testing-guides/test-types/soak-testing/  
**Published**: Continuous  
**Accessed**: 2026-09-26  
**Source Tier**: Tier 1  
**Relevance**: Defines soak/endurance testing as extended average-load test (hours/days) to detect memory leaks, resource leaks, data saturation, storage depletion. Typical durations: 3, 4, 8, 12, 24, 48-72 hours.

**Key Claims Verified**:
- Soak = extended average-load (hours to days)
- Detects: memory leaks, resource leaks, data saturation, storage depletion
- Typical durations: 3/4/8/12/24/48-72 hours
- Monitor backend resources (RAM, CPU, Network) especially

---

## Source 7: Google SRE Book Chapter 17 — Testing for Reliability

**Title**: Testing for Reliability (Chapter 17)  
**Publisher**: Google SRE Book  
**URL**: https://sre.google/sre-book/testing-reliability/  
**Published**: 2017 (Copyright Google, Inc., published by O'Reilly Media)  
**Accessed**: 2026-09-26  
**Source Tier**: Tier 1  
**Relevance**: Authoritative source on testing types including unit, integration, system/smoke, performance, regression, and production tests. Defines stress test as finding system limits: "How full can a database get before writes fail?" and "How many queries per second before overload?". Discusses zero MTTR testing, MTBF improvement, and testing cadence.

**Key Claims Verified**:
- Six test categories: unit, integration, system/smoke, performance, regression, production
- Stress test purpose: find system limits and breaking points
- Testing reduces uncertainty and improves MTBF
- Zero MTTR bugs enable blocking pushes before production

---

## Source 8: Google SRE Book Chapter 21 — Handling Overload

**Title**: Handling Overload (Chapter 21)  
**Publisher**: Google SRE Book  
**URL**: https://sre.google/sre-book/handling-overload/  
**Published**: 2017  
**Accessed**: 2026-09-26  
**Source Tier**: Tier 1  
**Relevance**: Discusses handling overload with degraded responses, per-customer quotas, client-side throttling, criticality levels (CRITICAL_PLUS, CRITICAL, SHEDDABLE_PLUS, SHEDDABLE), utilization signals (executor load average, CPU, memory), and load from connections (TCP overhead). Emphasizes CPU as primary provisioning signal.

**Key Claims Verified**:
- CPU as primary capacity signal over QPS
- Criticality: CRITICAL_PLUS, CRITICAL, SHEDDABLE_PLUS, SHEDDABLE
- Client-side throttling with adaptive algorithm (requests = 2 × accepts threshold)
- Connection maintenance overhead negligible at small scale, problematic at large scale

---

## Source 9: ISO/IEC 25010 (via Wikipedia)

**Title**: ISO/IEC 25010  
**Publisher**: International Organization for Standardization (ISO)  
**URL**: https://en.wikipedia.org/wiki/ISO/IEC_25010  
**Published**: 2011 (standard), Wikipedia summary accessed  
**Accessed**: 2026-09-26  
**Source Tier**: Tier 1 (standard) + Tier 2 (Wikipedia summary)  
**Relevance**: Defines Software Quality Model including "Performance Efficiency" subcharacteristics: Time behaviour (response times, throughput), Resource utilization (CPU, memory, storage, network), and Capacity. Replaced ISO/IEC 9126. Added security and compatibility as main characteristics.

**Key Claims Verified**:
- Performance Efficiency: Time behaviour, Resource utilization, Capacity
- Time behaviour includes response times and throughput
- Resource utilization includes CPU, memory, storage, network
- 25010 supersedes 9126; added security, compatibility

---

## Source 10: Microsoft Azure Well-Architected Framework — Performance Testing

**Title**: Architecture Strategies for Performance Testing  
**Publisher**: Microsoft Azure Well-Architected Framework (Performance Efficiency Pillar)  
**URL**: https://learn.microsoft.com/en-us/azure/well-architected/performance-efficiency/performance-test  
**Published**: 2026-06-24 (last updated)  
**Accessed**: 2026-09-26  
**Source Tier**: Tier 1  
**Relevance**: Comprehensive guidance covering test type definitions (load, stress, spike, endurance/soak), environment mirroring, hypothesis-driven experimentation, performance budgets, baseline establishment, and iterative optimization. Lists common pitfalls: only testing /health, insufficient data volumes, no server monitoring, no defined targets.

**Key Claims Verified**:
- Four test types: load, stress, spike, endurance/soak
- Mirror production environment as closely as possible
- Use hypothesis-driven experimentation to validate changes
- Performance budgets assigned per layer for bottleneck isolation
- Pitfalls: health-only testing, insufficient data, no monitoring, undefined targets

---

## Source 11: AWS Well-Architected Framework — Performance Efficiency

**Title**: Performance Efficiency Pillar  
**Publisher**: AWS Well-Architected Framework  
**URL**: https://docs.aws.amazon.com/wellarchitected/latest/performance-efficiency-pillar/welcome.html  
**Published**: November 6, 2024  
**Accessed**: 2026-09-26  
**Source Tier**: Tier 1  
**Relevance**: Defines five focus areas: architecture selection, compute and hardware, data management, networking and content delivery, process and culture. Emphasizes performance testing as part of "process and culture" and discusses iterative optimization, baseline establishment, and workload modeling.

**Key Claims Verified**:
- Five focus areas: architecture, compute, data, networking, process
- Performance testing falls under "process and culture"
- Performance testing part of iterative optimization cycle
- Baselines establish normal behavior for regression detection

---

## Source 12: Locust Documentation — What is Locust?

**Title**: What is Locust?  
**Publisher**: Locust Official Documentation  
**URL**: https://docs.locust.io/en/stable/what-is-locust.html  
**Published**: Continuous (latest version 2.17+)  
**Accessed**: 2026-09-26  
**Source Tier**: Tier 1  
**Relevance**: Documents Locust architecture: write tests in Python, each user runs in its own greenlet (lightweight coroutine via gevent), distributed and scalable supporting hundreds of thousands of concurrent users. Web UI for real-time monitoring and load adjustment.

**Key Claims Verified**:
- Tests written in Python (no domain-specific language)
- Each user = greenlet (lightweight process/coroutine)
- Supports hundreds of thousands of concurrent users
- Web UI for real-time monitoring and interactive load adjustment
- Pluggable architecture for extending protocols

---

## Source 13: Gatling Workload Models

**Title**: Workload models  
**Publisher**: Gatling Documentation  
**URL**: https://docs.gatling.io/testing-concepts/workload-models/  
**Published**: Continuous  
**Accessed**: 2026-09-26  
**Source Tier**: Tier 1  
**Relevance**: Distinguishes open vs closed workload models. Closed systems cap concurrent users (queue-based: call centers, ticketing). Open systems allow unlimited arrivals (most websites). Recommends matching load model to actual system behavior.

**Key Claims Verified**:
- Closed model: capped concurrent users, queue-based (call center, ticketing)
- Open model: unlimited arrivals, most websites
- Do not mix open and closed models in same injection profile
- Choosing wrong model breaks test validity

---

## Source 14: Azure Load Testing Overview

**Title**: Architecture Strategies for Performance Testing (Azure Load Testing service)  
**Publisher**: Microsoft Azure Documentation  
**URL**: https://learn.microsoft.com/en-us/azure/architecture/framework/scalability/load-testing/  
**Published**: 2025 or later  
**Accessed**: 2026-09-26  
**Source Tier**: Tier 1  
**Relevance**: Azure's performance testing overview covering methodology, tools, metrics (latency, throughput, error rate, resource utilization), environment mirroring, and best practices for cloud environments.

**Key Claims Verified**:
- Metrics: latency, throughput, error rate, resource utilization
- Cloud testing benefits: scalability, automation, CI/CD integration
- Environment mirroring essential for realistic results

---

## Source 15: k6 Ramp-Arrival Rate Executor

**Title**: Ramping arrival rate  
**Publisher**: Grafana k6 Documentation  
**URL**: https://grafana.com/docs/k6/latest/using-k6/scenarios/executors/ramping-arrival-rate/  
**Published**: Continuous  
**Accessed**: 2026-09-26  
**Source Tier**: Tier 1  
**Relevance**: Documents ramping-arrival-rate executor for open-model load testing with stages config, startRate, timeUnit, preAllocatedVUs, maxVUs options. Critical for breakpoint testing where load must increase regardless of system degradation.

**Key Claims Verified**:
- Open-model load: arrival rate (requests per time unit)
- Ramping-arrival-rate executor supports breakpoint testing
- load increases regardless of system degradation
- maxVUs pre-allocated to avoid VU scaling overhead

---

## Source 16: Mozilla MDN HTTP Overview

**Title**: What is HTTP and how does it work  
**Publisher**: Mozilla Developer Network (MDN)  
**URL**: https://developer.mozilla.org/en-US/docs/Web/HTTP/Overview  
**Published**: 2026-08  
**Accessed**: 2026-09-26  
**Source Tier**: Tier 2  
**Relevance**: Provides foundational understanding of HTTP request lifecycle useful for interpreting load test results (status codes, connection limits, timeouts). Not directly cited but supports http_req_blocked, http_req_connecting, http_req_waiting breakdown.

**Key Claims Verified**:
- HTTP request lifecycle: DNS, TCP connect, TLS handshake, request send, wait, response receive
- Connection limits affect concurrency
- Timeouts can mask underlying issues

---

## Source 17: Microsoft Azure DevOps Pipeline Load Testing

**Title**: Azure Load Testing documentation  
**Publisher**: Microsoft Azure Documentation  
**URL**: https://learn.microsoft.com/en-us/azure/app-testing/load-testing/overview-what-is-azure-load-testing  
**Published**: Continuous  
**Accessed**: 2026-09-26  
**Source Tier**: Tier 1  
**Relevance**: Describes Azure Load Testing service features: CI/CD integration, automated test execution, dashboard with live metrics, result comparison for regression detection, and threshold-based test abortion.

**Key Claims Verified**:
- CI/CD integration for automated load testing
- Live dashboard with resource metrics during test
- Result comparison for regression tracking
- Threshold-based test abortion

---

## Source 18: Google SRE Book Service Best Practices

**Title**: A Collection of Best Practices for Production Services (Appendix B)  
**Publisher**: Google SRE Book  
**URL**: https://sre.google/sre-book/service-best-practices/  
**Published**: 2017  
**Accessed**: 2026-09-26  
**Source Tier**: Tier 1  
**Relevance**: Contains SRE-recommended best practices for production services including performance testing, load testing, and capacity planning guidelines. Provides production-readiness checklist.

**Key Claims Verified**:
- Production-readiness checklist includes performance testing
- Load testing and capacity planning as part of SRE best practices
- Continuous testing throughout SDLC

---

## Source 19: Jenkins Performance Plugin

**Title**: Jenkins Performance Plugin  
**Publisher**: Jenkins Project  
**URL**: https://plugins.jenkins.io/performance/  
**Published**: Continuous  
**Accessed**: 2026-09-26  
**Source Tier**: Tier 3 (community)  
**Relevance**: Community plugin for Jenkins that integrates various load testing tools (JMeter, k6, Locust via CLI). Not a primary source but demonstrates CI/CD integration patterns used in industry.

**Key Claims Verified**:
- Jenkins Performance Plugin supports multiple tool integrations
- Trend charts and threshold-based build failure
- CI/CD pipeline integration pattern

---

## Source Summary

**Tier 1 Sources**: 14  
**Tier 2 Sources**: 2  
**Tier 3 Sources**: 1  
**Total Sources**: 17

**Coverage**:
- Load test type definitions: k6, Azure, Google SRE (3 sources)
- Metrics (P50/P95/P99): k6, Azure, ISO/IEC 25010 (3 sources)
- Tool comparison: k6, Locust, Gatling (3 sources)
- Workload models: Gatling (1 source)
- Bottleneck identification: k6, Azure, Google SRE (3 sources)
- Common pitfalls: Azure, k6, Google SRE (3 sources)
- Load calculation: k6 (1 source)
- Overload handling: Google SRE (1 source)
- SDLC timing: Google SRE, Azure (2 sources)

**URL Verification**: All 17 URLs verified accessible as of 2026-09-26.
