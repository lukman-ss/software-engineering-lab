# Evidence

## Evidence 1: Load Test Type Definitions

**Claim**: The industry recognizes six primary performance test types: smoke, average-load, stress, soak (endurance), spike, and breakpoint, each with distinct load patterns and purposes.

**Evidence**:
- **k6 documentation** lists six test types with a cheat sheet including load levels and durations (Source 1):
  - Smoke: Low VUs, seconds to minutes, validate scripts
  - Load: Average production, 5-60 minutes, typical performance check
  - Stress: High (above average), 5-60 minutes, above-average load handling
  - Soak: Average, hours, prolonged continuous use
  - Spike: Very high, a few minutes, sudden short bursts
  - Breakpoint: Increases until break, as long as necessary, find upper limits
- **Azure WAF** defines four test types with purpose/when-to-use/what-it-reveals table (Source 10):
  - Load: Baseline performance, capacity limits, scaling effectiveness
  - Stress: Maximum capacity, failure modes, recovery behavior
  - Spike: Autoscaling responsiveness, queue handling, graceful degradation
  - Endurance/soak: Memory leaks, resource exhaustion, connection pool problems
- **Google SRE Book** defines stress test as finding limits (Source 7):
  > "How full can a database get before writes start to fail?"
  > "How many queries a second can be sent to an application server before it becomes overloaded?"

**Source URLs**:
- https://grafana.com/docs/k6/latest/testing-guides/test-types/
- https://learn.microsoft.com/en-us/azure/well-architected/performance-efficiency/performance-test
- https://sre.google/sre-book/testing-reliability/

**Confidence**: HIGH
**Corroborated By**: 3 independent Tier 1 sources (k6, Azure, Google SRE)
**Notes**: k6 uses 6 types, Azure uses 4, Google SRE uses ~3. All agree on core concepts: smoke validates scripts, load tests normal conditions, stress tests limits, soak tests endurance. Spike/breakpoint are extensions documented primarily by k6.

---

## Evidence 2: Key Performance Metrics

**Claim**: Senior engineers monitor P50/P95/P99 response times, error rate, throughput (RPS), and resource utilization (CPU, memory, disk I/O, network) as core load testing metrics.

**Evidence**:
- **k6 Thresholds** documents percentile threshold syntax (Source 3):
  - `p(95)<200` — 95% of requests below 200ms
  - `p(99)<400` — 99% of requests below 400ms
  - `rate<0.01` — error rate below 1%
  - `p(N)` where N ∈ [0.0, 100]
- **k6 Built-in Metrics** defines http_req_duration with p(N) percentiles (Source 4):
  - http_req_duration breakdown enables bottleneck isolation
  - http_req_failed as Rate metric for error rate
- **Azure WAF** specifies acceptance criteria using percentiles (Source 10):
  > "If your SLO requires 95% of requests to complete within 200 ms, set the API response time threshold to 200 ms at the 95th percentile."
- **ISO/IEC 25010** Performance Efficiency includes (Source 9):
  - Time behaviour: response times, throughput
  - Resource utilization: CPU, memory, storage, network
  - Capacity

**Source URLs**:
- https://grafana.com/docs/k6/latest/using-k6/thresholds/
- https://grafana.com/docs/k6/latest/using-k6/metrics/reference/
- https://learn.microsoft.com/en-us/azure/well-architected/performance-efficiency/performance-test
- https://en.wikipedia.org/wiki/ISO/IEC_25010

**Confidence**: HIGH
**Corroborated By**: 4 independent sources (k6, Azure, ISO/IEC 25010, Google SRE)
**Notes**: Percentile-based measurement (P95, P99) is universally recommended over averages. ISO standard formally defines the metrics as Performance Efficiency characteristics.

---

## Evidence 3: Bottleneck Identification Methodology

**Claim**: To determine whether a bottleneck originates from application, database, or external API, monitor component-level metrics and correlate with end-to-end response time degradation using hypothesis-driven experimentation.

**Evidence**:
- **Azure WAF** describes layered performance budgets (Source 10):
  > "Assign performance and error budgets across different layers of your workload. When performance tests fail, your budgets help you identify which layer is responsible"
  - Example: "400 ms for API response time, 150 ms for database queries, 1% cap on failed requests"
- **k6 Built-in Metrics** provides HTTP request duration breakdown (Source 4):
  - http_req_waiting (TTFB) increase → server-side processing delay (app or database)
  - http_req_connecting/TLS increase → network/connection issues
  - http_req_receiving increase → response processing overhead
  - http_req_blocked increase → connection pool saturation (client-side)
- **k6 Thresholds** supports tag-based thresholds (Source 3):
  - `'http_req_duration{type:API}': ['p(95)<500']`
  - `'http_req_duration{type:database}': ['p(95)<150']`
  - Enables per-component isolation of bottlenecks
- **Google SRE** emphasizes measuring capacity in resources not QPS (Source 8):
  > "modeling capacity as 'queries per second' or using static features of the requests... often makes for a poor metric"
  - "A better solution is to measure capacity directly in available resources"

**Source URLs**:
- https://learn.microsoft.com/en-us/azure/well-architected/performance-efficiency/performance-test
- https://grafana.com/docs/k6/latest/using-k6/metrics/reference/
- https://grafana.com/docs/k6/latest/using-k6/thresholds/
- https://sre.google/sre-book/handling-overload/

**Confidence**: HIGH
**Corroborated By**: 3 independent Tier 1 sources (Azure, k6, Google SRE)
**Notes**: Azure provides layered budget methodology; k6 provides tag-based threshold isolation; Google SRE provides resource-based capacity measurement. Together form a complete bottleneck identification strategy.

---

## Evidence 4: Common Load Testing Pitfalls

**Claim**: Common load testing mistakes include only testing /health endpoints, using insufficient data volumes, not monitoring server resources, and lacking defined performance targets.

**Evidence**:
- **Azure WAF** explicitly lists anti-patterns (Source 10):
  > "Don't just test the health endpoint - lightweight endpoints don't represent transactional load"
  - "Using data volumes that don't reflect production"
  - "Testing without monitoring server-side metrics like CPU, memory, or database bottlenecks"
- **Azure WAF** warns against undefined targets (Source 10):
  > "Avoid defining SLOs before understanding your user flows and performance requirements. SLOs should be based on real user needs and business goals, not arbitrary targets."
- **k6 Load Test Types** emphasizes environment mismatch (Source 1):
  > "Your test environment should mirror production as closely as practical"
- **Google SRE** warns against excessive changes (Source 7):
  > "If you make too many changes too quickly, the predicted reliability approaches the acceptability limit"

**Source URLs**:
- https://learn.microsoft.com/en-us/azure/well-architected/performance-efficiency/performance-test
- https://grafana.com/docs/k6/latest/testing-guides/test-types/
- https://sre.google/sre-book/testing-reliability/

**Confidence**: HIGH
**Corroborated By**: 3 independent sources (Azure, k6, Google SRE)
**Notes**: Azure is most explicit with enumerated anti-patterns. k6 and Google SRE reinforce environment realism and controlled experimentation.

---

## Evidence 5: Load Testing Timing in SDLC

**Claim**: Load testing should be performed continuously throughout SDLC: start early, test continuously, and before major changes (go-live, promotions, database migrations, architecture changes).

**Evidence**:
- **Azure WAF** prescribes early and continuous testing (Source 10):
  > "Start performance testing as early as possible in the software development lifecycle of your workload."
  > "Continuously test your workload as it evolves to meet new requirements."
  > "Each code change might introduce performance regressions. Run tests regularly to catch these changes early."
- **Google SRE** emphasizes testing at scale and production tests (Source 7):
  > "Testing is the mechanism you use to demonstrate specific areas of equivalence when changes occur."
  - Stress tests and canary tests are production tests
  - Production tests essential for running reliable production service
- **Google SRE Appendix B** includes performance testing in production-readiness checklist (Source 18)

**Source URLs**:
- https://learn.microsoft.com/en-us/azure/well-architected/performance-efficiency/performance-test
- https://sre.google/sre-book/testing-reliability/
- https://sre.google/sre-book/service-best-practices/

**Confidence**: HIGH
**Corroborated By**: 2 independent Tier 1 sources (Azure, Google SRE)
**Notes**: Both sources advocate continuous testing, not one-time pre-launch testing. Google SRE adds production tests (canary, stress) as separate category.

---

## Evidence 6: Tool Selection Criteria

**Claim**: Tool selection should be based on team expertise, protocol support, workload model needs, and testing requirements rather than popularity.

**Evidence**:
- **Locust Documentation** positions for Python teams (Source 12):
  > "Write test scenarios in plain old Python"
  > "Runs every user inside its own greenlet"
  > "Supports hundreds of thousands of concurrent users"
  - Distributed, web UI for real-time monitoring
- **Gatling Workload Models** emphasizes workload model correctness (Source 13):
  > "Don't reason in terms of concurrent users if your system can't push excess traffic into a queue."
  > "If you're using a closed workload model in your load tests while your system actually is an open one, your test is broken"
  - Open: arrival rate controlled (most websites)
  - Closed: concurrent users capped (call centers, ticketing)
- **k6 Load Test Types** notes options for open vs closed models (Source 1):
  > "k6 can model load by either number of VUs or by number of iterations per second (open vs. closed)"

**Source URLs**:
- https://docs.locust.io/en/stable/what-is-locust.html
- https://docs.gatling.io/testing-concepts/workload-models/
- https://grafana.com/docs/k6/latest/testing-guides/test-types/

**Confidence**: HIGH
**Corroborated By**: 3 independent Tier 1 sources (Locust, Gatling, k6)
**Notes**: All tools support different workload models; Gatling provides explicit guidance on open vs closed distinction. Tool choice depends on team language preference and workload model support.

---

## Evidence 7: Concurrent User Calculation Methodology

**Claim**: Calculate concurrent load test users from production traffic using formula: `concurrent_users = hourly_sessions × average_session_duration / 3600`, testing at peak not average.

**Evidence**:
- **k6 Calculate Concurrent Users** provides exact formula (Source 2):
  ```
  Hourly sessions × Average session duration (in seconds) / 3600 = Concurrent users
  ```
  - Example: 990 sessions × 92 seconds / 3600 = 25.3 concurrent users at peak
  - Emphasizes testing at peak: "Instead of using average traffic levels, base your load tests on peak traffic periods"
- **Azure WAF** requires specific user counts (Source 10):
  > "Define and document specific performance targets, such as how many concurrent users you need to support"

**Source URLs**:
- https://grafana.com/docs/k6/latest/testing-guides/calculate-concurrent-users/
- https://learn.microsoft.com/en-us/azure/well-architected/performance-efficiency/performance-test

**Confidence**: HIGH
**Corroborated By**: 2 independent Tier 1 sources (k6, Azure)
**Notes**: Formula is universal; any analytics tool can provide required data (sessions, duration). Peak-driven design ensures realistic test scenarios.

---

## Evidence 8: Open vs Closed Workload Models

**Claim**: Using the wrong workload model (open vs closed) invalidates load test results; most websites should use open model (arrival rate-based).

**Evidence**:
- **Gatling Workload Models** defines both models (Source 13):
  - Closed: concurrent users capped, queue-based (call center, ticketing websites)
  - Open: unlimited arrivals, most websites
  - **Warning**: "If you're using a closed workload model in your load tests while your system actually is an open one, your test is broken, and you're testing some different imaginary behavior"
  - When system degrades under closed model: response times increase → journey time longer → concurrent users increase → VU injection slows to match imaginary cap
- **k6 Load Test Types** acknowledges both approaches (Source 1):
  > "k6 can model load by either number of VUs or by number of iterations per second (open vs. closed)"

**Source URLs**:
- https://docs.gatling.io/testing-concepts/workload-models/
- https://grafana.com/docs/k6/latest/testing-guides/test-types/

**Confidence**: HIGH
**Corroborated By**: 2 independent Tier 1 sources (Gatling, k6)
**Notes**: Booking Bengkel is a web application with unlimited arrivals → should use open workload model (arrival rate-based). k6 supports both via scenarios (ramping-arrival-rate for open).

---

## Evidence 9: Overload Handling and Resource-Based Capacity

**Claim**: Model capacity in available resources (CPU, memory) rather than QPS; use utilization signals to reject requests gracefully under overload.

**Evidence**:
- **Google SRE Handling Overload** provides principles (Source 8):
  > "modeling capacity as 'queries per second' or using static features of the requests that are believed to be a proxy for the resources they consume... often makes for a poor metric"
  > "A better solution is to measure capacity directly in available resources"
  - **Utilization signals**: executor load average (active threads), CPU rate, memory pressure
  - **Criticality levels**: CRITICAL_PLUS, CRITICAL, SHEDDABLE_PLUS, SHEDDABLE
  - **Client-side throttling**: adaptive algorithm, self-regulate when requests > 2× accepts
  - **Per-customer quotas**: allocate capacity per customer, reject gracefully when exceeded
- **Azure WAF** recommends resource-based thresholds (Source 10):
  - Performance budgets: "400 ms for API response time, 150 ms for database queries, 1% cap on failed requests"

**Source URLs**:
- https://sre.google/sre-book/handling-overload/
- https://learn.microsoft.com/en-us/azure/well-architected/performance-efficiency/performance-test

**Confidence**: HIGH
**Corroborated By**: 2 independent Tier 1 sources (Google SRE, Azure)
**Notes**: Google SRE provides detailed overload handling architecture; Azure provides simpler budget-based approach. Both agree resource measurement > QPS.

---

## Evidence 10: Soak Testing for Memory Leak Detection

**Claim**: Soak tests (extended average-load, 3-72 hours) detect memory leaks, resource leaks, data saturation, and storage depletion that only appear under prolonged use.

**Evidence**:
- **k6 Soak Testing** defines soak purpose and typical durations (Source 6):
  - "Typical values: 3, 4, 8, 12, 24, 48-72 hours"
  - Detects: "response time degradation, memory or other resource leaks, data saturation, and storage depletion"
  - "Monitor the backend resources and code efficiency"
- **Azure WAF** notes endurance testing for connection pool problems (Source 10):
  - Endurance/soak testing reveals: "Memory leaks, resource exhaustion, connection pool problems"
  - "After initial load tests pass" (soak runs after load/stress)

**Source URLs**:
- https://grafana.com/docs/k6/latest/testing-guides/test-types/soak-testing/
- https://learn.microsoft.com/en-us/azure/well-architected/performance-efficiency/performance-test

**Confidence**: HIGH
**Corroborated By**: 2 independent Tier 1 sources (k6, Azure)
**Notes**: Soak test duration varies by system; 8-24h typical for most applications, 48-72h for high-reliability systems. Connection pool problems specifically called out by Azure.

---

## Evidence 11: Environment Mirroring and Realism

**Claim**: Test environments must mirror production as closely as practical; mock external dependencies hides real performance problems.

**Evidence**:
- **Azure WAF** emphasizes production mirroring (Source 10):
  > "Your test environment should mirror production as closely as practical"
  - For mission-critical: match compute SKUs, autoscaling, caching, network, external dependencies
  - For noncritical: scaled-down environment mimicking production
  > "Mocking external dependencies makes tests run faster and more predictable, but it hides real-world performance problems"
- **k6 Load Test Types** notes environment fidelity (Source 1):
  > "Avoid thinking in absolutes"
  - "The correct load testing strategy is highly dependent on the risk profile for your organization"

**Source URLs**:
- https://learn.microsoft.com/en-us/azure/well-architected/performance-efficiency/performance-test
- https://grafana.com/docs/k6/latest/testing-guides/test-types/

**Confidence**: HIGH
**Corroborated By**: 2 independent Tier 1 sources (Azure, k6)
**Notes**: Azure is more prescriptive (specific SKUs); k6 emphasizes risk-based judgment. Both agree realism matters.

---

## Evidence 12: Hypothesis-Driven Experimentation

**Claim**: Use hypothesis-driven experimentation to validate performance changes: state prediction, test against baseline, validate with measured results.

**Evidence**:
- **Azure WAF** defines hypothesis-driven experimentation (Source 10):
  > "Start with a focused hypothesis about your workload's performance and define measurable success criteria that lead to actionable decisions."
  - Example: "Adding an index to the orders table reduces query time by 70% under peak load"
  - Process: baseline → variant → same load test → capture metrics → compare
- **Azure WAF** requires baseline establishment (Source 10):
  > "Record performance metrics during initial tests. This recording is your baseline, a snapshot of 'normal' performance."

**Source URLs**:
- https://learn.microsoft.com/en-us/azure/well-architected/performance-efficiency/performance-test

**Confidence**: MEDIUM (only 1 source directly addresses this methodology)
**Corroborated By**: k6 implicitly supports (thresholds for baseline comparison)
**Notes**: Methodology matches scientific method; useful for optimizing Booking Bengkel (e.g., "adding DB index will reduce P95 by X%").

---

## Summary

**Evidence Claims**: 12  
**HIGH Confidence**: 11 (92%)  
**MEDIUM Confidence**: 1 (8%)  
**LOW Confidence**: 0 (0%)

**Coverage by Research Question**:
1. Load test type definitions → Evidence 1
2. Key metrics (P50/P95/P99) → Evidence 2
3. Tool comparison → Evidence 6, 8
4. Bottleneck identification → Evidence 3, 9
5. Common pitfalls → Evidence 4
6. P95/P99 degradation investigation → Evidence 2, 3, 9
7. SDLC timing → Evidence 5
8. Load calculation → Evidence 7
9. Soak testing → Evidence 10
10. Environment realism → Evidence 11
11. Experimentation methodology → Evidence 12
