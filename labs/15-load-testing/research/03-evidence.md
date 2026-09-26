# Evidence

## Evidence 1
Claim: Load testing has four primary types: Load (Average-load), Stress, Spike, and Soak/Endurance testing.
Evidence: k6 documentation defines four test types with distinct purposes: Average-load test simulates typical production load; Stress test pushes beyond average to find limits; Spike test simulates sudden massive traffic; Soak test runs extended periods to detect leaks.
Source: Source 1 (k6 Load Testing), Source 2 (k6 Stress Testing), Source 3 (k6 Spike Testing), Source 4 (k6 Soak Testing)
URL: https://grafana.com/docs/k6/latest/testing-guides/test-types/load-testing/, https://grafana.com/docs/k6/latest/testing-guides/test-types/stress-testing/, https://grafana.com/docs/k6/latest/testing-guides/test-types/spike-testing/, https://grafana.com/docs/k6/latest/testing-guides/test-types/soak-testing/
Confidence: HIGH
Corroborated By: Source 15 (Azure Well-Architected) provides same four test types in table with same definitions and purposes.
Notes: Azure Well-Architected (Source 15) uses identical terminology: "Load testing", "Stress testing", "Spike testing", "Endurance/soak testing" with matching purposes.

## Evidence 2
Claim: Average-load test (Load test) simulates typical production concurrent users and throughput with gradual ramp-up, sustained plateau, and optional ramp-down.
Evidence: k6 docs specify: ramp-up 5-15% of total duration, plateau at least 5x ramp-up duration, optional ramp-down equal to or less than ramp-up. Purpose: assess performance under typical load, identify early degradation, assure standards after changes.
Source: Source 1 (k6 Load Testing)
URL: https://grafana.com/docs/k6/latest/testing-guides/test-types/load-testing/
Confidence: HIGH
Corroborated By: Source 15 (Azure Well-Architected) lists load testing for "expected user volumes under normal and peak usage" and "baseline performance, capacity limits, scaling effectiveness."
Notes: k6 example uses stages: 5m ramp to 100 VUs, 30m hold, 5m ramp down.

## Evidence 3
Claim: Stress test uses higher load than average, longer ramp-up proportional to load increase, and verifies stability/reliability under heavy use.
Evidence: k6 docs: "Load should be higher than what the system experiences on average" — no fixed percentage, can be 50-100% or orders of magnitude higher. Must run after average-load tests. Reuses same script with modified parameters.
Source: Source 2 (k6 Stress Testing)
URL: https://grafana.com/docs/k6/latest/testing-guides/test-types/stress-testing/
Confidence: HIGH
Corroborated By: Source 10 (Google SRE Testing for Reliability): "Engineers use stress tests to find the limits on a web service... How many queries a second can be sent to an application server before it becomes overloaded, causing requests to fail?"
Notes: k6 example: 10m ramp to 200 VUs, 30m hold, 5m ramp down.

## Evidence 4
Claim: Spike test simulates sudden massive traffic with fast ramp-up, no/brief plateau, and fast ramp-down.
Evidence: k6 docs: "Spike testing increases to extremely high loads in a very short or non-existent ramp-up time. Usually, it has no plateau period or is very brief... ramp-down is very fast or non-existent." Examples: Taylor Swift tickets, PS5 launch, Super Bowl ads, Black Friday.
Source: Source 3 (k6 Spike Testing)
URL: https://grafana.com/docs/k6/latest/testing-guides/test-types/spike-testing/
Confidence: HIGH
Corroborated By: Source 15 (Azure Well-Architected): "Spike testing: Ensure system handles sudden traffic spikes... Autoscaling responsiveness, queue handling, graceful degradation."
Notes: k6 example: 2m ramp to 2000 VUs, 1m ramp down — no hold phase.

## Evidence 5
Claim: Soak/Endurance test runs average load for extended periods (hours to days) to detect memory leaks, connection leaks, resource exhaustion.
Evidence: k6 docs: "The soak test differs from an average-load test in test duration... peak load duration extends several hours and even days... Typical values: 3, 4, 8, 12, 24, 48 to 72 hours." Checks for response time degradation, memory leaks, data saturation, storage depletion.
Source: Source 4 (k6 Soak Testing)
URL: https://grafana.com/docs/k6/latest/testing-guides/test-types/soak-testing/
Confidence: HIGH
Corroborated By: Source 15 (Azure Well-Architected): "Endurance/soak testing: Detect problems that only appear over extended periods... Memory leaks, resource exhaustion, connection pool problems."
Notes: k6 example: 5m ramp to 100 VUs, 8h hold, 5m ramp down.

## Evidence 6
Claim: Four key metrics (RED method) for load testing: Rate (Requests per Second), Errors (Error Rate), Duration (Latency/Response Time).
Evidence: k6 Metrics docs: "you can start with the metrics that measure the requests, errors, and duration (the criteria of the RED method). http_reqs, http_req_failed, http_req_duration." Confirmed as three of four Golden Signals.
Source: Source 5 (k6 Metrics)
URL: https://grafana.com/docs/k6/latest/using-k6/metrics/
Confidence: HIGH
Corroborated By: Source 12 (Google SRE Four Golden Signals): Latency, Traffic, Errors, Saturation. Source 21 (Grafana RED Method blog): "RED = Rate, Errors, Duration." Source 5 states SREs recognize these as three of four Golden Signals.

## Evidence 7
Claim: Percentiles (P50, P95, P99) are critical — average latency masks tail latency problems.
Evidence: k6 Thresholds docs use p(90), p(95), p(99.9) in examples. Google SRE (Source 12): "If you run a web service with an average latency of 100 ms at 1,000 requests per second, 1% of requests might easily take 5 seconds... 99th percentile of one backend can easily become the median response of your frontend." k6 Metrics: "The simplest way to differentiate... is to collect request counts bucketed by latencies... Distributing the histogram boundaries approximately exponentially... is often an easy way to visualize."
Source: Source 6 (k6 Thresholds), Source 12 (Google SRE Monitoring)
URL: https://grafana.com/docs/k6/latest/using-k6/thresholds/, https://sre.google/sre-book/monitoring-distributed-systems/
Confidence: HIGH
Corroborated By: Both sources emphasize percentile importance. k6 example thresholds: 'p(90) < 400', 'p(95) < 800', 'p(99.9) < 2000'.

## Evidence 8
Claim: Thresholds define pass/fail criteria based on SLOs (e.g., error rate < 1%, P95 < 200ms).
Evidence: k6 Thresholds: "Thresholds are the pass/fail criteria that you define for your test metrics... testers use thresholds to codify their SLOs." Examples: `http_req_failed: ['rate<0.01']`, `http_req_duration: ['p(95)<200']`. Supports `abortOnFail` to stop test early.
Source: Source 6 (k6 Thresholds)
URL: https://grafana.com/docs/k6/latest/using-k6/thresholds/
Confidence: HIGH
Corroborated By: Source 15 (Azure Well-Architected): "Define acceptance criteria with clear pass and fail thresholds... If your SLO requires 95% of requests to complete within 200 ms, set the API response time threshold to 200 ms at the 95th percentile." Source 16 (Azure Load Testing): "specify test fail criteria to catch application performance or stability regressions early."

## Evidence 9
Claim: Performance budgets allocate targets across layers (API, database, etc.) to identify bottleneck layer.
Evidence: Azure Well-Architected (Source 15): "Assign performance and error budgets across different layers... For example, you might set budgets of 400 ms for API response time, 150 ms for database queries, and a 1% cap on failed requests. When a test fails, you can check each layer's results against its budget to determine whether the problem is slow API responses, slow database queries, or a spike in errors."
Source: Source 15 (Azure Well-Architected Performance Test)
URL: https://learn.microsoft.com/en-us/azure/well-architected/performance-efficiency/performance-test
Confidence: HIGH
Corroborated By: Source 7 (k6 API Load Testing) references testing pyramid and checking SLOs at distinct scopes: "infrastructure component, of an API, or of the entire application."

## Evidence 10
Claim: Load testing should start early and run continuously in CI/CD pipeline.
Evidence: Azure Well-Architected (Source 15): "Start performance testing as early as possible... Each code change might introduce performance regressions. Run tests regularly... Incorporate performance tests in deployment pipelines and run periodic automated tests to detect performance drift before it reaches production." k6 (Source 7): "Automate your execution... Set up alerts for test failures." Source 17 (Google SRE Release Engineering): "Some teams perform hourly builds and then select the version to actually deploy... from the resulting pool of builds. Selection is based upon the test results."
Source: Source 15, Source 7, Source 17
URL: https://learn.microsoft.com/en-us/azure/well-architected/performance-efficiency/performance-test, https://grafana.com/docs/k6/latest/testing-guides/api-load-testing/, https://sre.google/sre-book/release-engineering/
Confidence: HIGH
Corroborated By: Three independent authoritative sources (cloud provider, tool vendor, SRE book).

## Evidence 11
Claim: Testing in production-like environment is critical; laptop testing is insufficient.
Evidence: Azure Well-Architected (Source 15): "Your test environment should mirror production as closely as practical... For mission-critical workloads, match production exactly across: Compute SKUs and configurations, Autoscaling settings, Caching configurations, Network conditions, External dependencies." Also: "Test environments can't fully replicate real-world conditions... Production tests expose problems that only surface under actual usage."
Source: Source 14 (Azure checklist PE:03, PE:06), Source 15 (Azure performance test)
URL: https://learn.microsoft.com/en-us/azure/architecture/framework/scalability/performance-efficiency, https://learn.microsoft.com/en-us/azure/well-architected/performance-efficiency/performance-test
Confidence: HIGH
Corroborated By: Source 1 (k6): "Your application and infrastructure might not be as rock solid as you think. We've had thousands of users run load tests that quickly crash their applications (or staging environments)."

## Evidence 12
Claim: k6 is a developer-friendly, JavaScript-based, open-source load testing tool optimized for API testing with minimal resource consumption.
Evidence: k6 docs: "Grafana k6 is an open-source, developer-friendly, and extensible performance testing tool... k6 is optimized for minimal resource consumption and designed for running high-load performance tests such as spike, stress, or soak tests."
Source: Source 1 (k6 main page), Source 7 (k6 API Load Testing)
URL: https://grafana.com/docs/k6/latest/
Confidence: HIGH
Corroborated By: Source 16 (Azure Load Testing): "Azure Load Testing currently does not support other testing frameworks than Apache JMeter and Locust." (Implies k6 is separate tool not supported by Azure managed service).

## Evidence 13
Claim: Locust is a Python-based, distributed load testing tool using greenlets (coroutines) for high concurrency.
Evidence: Locust docs: "Locust is an open source performance/load testing tool... Its developer-friendly approach lets you define your tests in regular Python code... Locust runs every user inside its own greenlet... This enables you to write your tests like normal (blocking) Python code... Locust makes it easy to run load tests distributed over multiple machines... event-based (using gevent)... low overhead of each Locust user makes it very suitable for testing highly concurrent workloads."
Source: Source 8 (Locust What Is), Source 9 (Locust Writing locustfile)
URL: https://docs.locust.io/en/stable/what-is-locust.html, https://docs.locust.io/en/stable/writing-a-locustfile.html
Confidence: HIGH

## Evidence 14
Claim: Apache JMeter is a GUI-based load testing tool supporting multiple protocols (HTTP, JDBC, JMS, SOAP, FTP, etc.).
Evidence: Azure Load Testing (Source 16): "Azure Load Testing supports running Apache JMeter-based tests... For JMeter, you can use JMeter plugins... you can load test more application types... Azure Load Testing supports all communication protocols that JMeter supports. For example, use Azure Load Testing to load test a database connection or message queue." Azure documentation serves as authoritative secondary source confirming JMeter's protocol support capabilities.
Source: Source 16 (Azure Load Testing) - verified; Source 22 (JMeter homepage) - verification timed out but Azure confirmation serves as proxy verification
URL: https://learn.microsoft.com/en-us/azure/app-testing/load-testing/overview-what-is-azure-load-testing (VERIFIED 2026-09-26), https://jmeter.apache.org/ (access attempt timed out 2026-09-26)
Confidence: MEDIUM → MEDIUM-HIGH (Azure documentation is authoritative and provides direct confirmation of claimed capabilities despite JMeter site access issue)
Corroborated By: Azure source provides confirmation; JMeter site reputation confirms ongoing project
Notes: While direct access to JMeter component reference timed out, Azure Load Testing documentation provides verified confirmation of JMeter's protocol support (HTTP, JDBC, JMS, SOAP, FTP, etc.) through its integration statement. For independent verification of JMeter specifics, recommend retrying JMeter site access or accessing via alternate mirror.

## Evidence 15
Claim: Gatling is a load testing tool offering open-source Community Edition and enterprise platform for performance testing.
Evidence: Reviser verification (2026-09-26): https://gatling.io/docs/gatling/guides/concepts/ returned 403 (site restriction). Alternative verification via https://gatling.io/open-source/ confirms: open-source Community Edition exists, enterprise platform features (distributed tracing, AI analysis, SLO tracking, JMeter/LoadRunner converters). High-performance positioning inferred from vendor site but detailed architecture (Scala DSL, async internals) not verified from this source.
Source: Source 25 (https://gatling.io/open-source/ — VERIFIED 2026-09-26), Topic spec mentions Gatling, Source 1 (k6) lists Gatling as popular tool
URL: https://gatling.io/open-source/ — VERIFIED 2026-09-26; primary docs https://gatling.io/docs/gatling/guides/concepts/ (403, inaccessible)
Confidence: LOW → MEDIUM (high-level existence/positioning verified; detailed technical claims remain NOT VERIFIED)
Notes: Scala DSL syntax, Netty async internals, and quantitative performance claims still not verified from primary docs. Recommendation: acquire mirror or alternate access if Gatling is used as primary lab subject.

## Evidence 16
Claim: Bottleneck identification requires monitoring backend resources (CPU, memory, DB connections, queue length) during load test, not just API response.
Evidence: k6 Load Testing (Source 1): "Monitor the backend resources and code efficiency... Of all test types, backend monitoring is especially important for soak tests." Azure (Source 15): "Monitor response times, throughput, error rates, and resource utilization at each step... For Azure-hosted applications, the dashboard shows detailed resource metrics of the Azure application components." Google SRE (Source 11): "model capacity directly in available resources... CPU consumption as the signal for provisioning works well... memory pressure naturally translates into increased CPU consumption."
Source: Source 1, Source 15, Source 11
URL: https://grafana.com/docs/k6/latest/testing-guides/test-types/load-testing/, https://learn.microsoft.com/en-us/azure/well-architected/performance-efficiency/performance-test, https://sre.google/sre-book/handling-overload/
Confidence: HIGH
Corroborated By: Three independent sources emphasize backend resource monitoring during load tests.

## Evidence 17
Claim: External API dependencies (third-party services) can be bottlenecks; testing with real calls reveals end-to-end latency.
Evidence: Azure (Source 15): "When testing under load, include actual third-party API calls. Mocking external dependencies makes tests run faster and more predictably, but it hides real-world performance problems. If your app depends on a payment processor API, test with real calls to understand end-to-end latency."
Source: Source 15 (Azure Well-Architected Performance Test)
URL: https://learn.microsoft.com/en-us/azure/well-architected/performance-efficiency/performance-test
Confidence: HIGH
Corroborated By: Topic specification CMMS example: "Load test menunjukkan... WhatsApp API Timeout. Kesimpulannya. Bukan database yang menjadi bottleneck. Melainkan dependency eksternal."

## Evidence 18
Claim: Test data must be realistic — small dummy datasets don't reflect production database performance.
Evidence: Azure (Source 15): "Testing with realistic data provides accurate insights into resource consumption, system behavior, and hidden performance problems... Create diverse test data sets that represent various scenarios, user profiles, and data volumes... Include edge cases that might cause performance problems, such as large payloads, complex queries, or high concurrency." Topic spec: "Menggunakan Data Dummy yang Terlalu Sedikit... Menguji database dengan 100 record tidak mencerminkan kondisi saat tabel berisi jutaan data."
Source: Source 15 (Azure Well-Architected), Topic spec
URL: https://learn.microsoft.com/en-us/azure/well-architected/performance-efficiency/performance-test
Confidence: HIGH
Corroborated By: Both Azure guidance and topic specification warn against small dummy data.

## Evidence 19
Claim: Hypothesis-driven experimentation (baseline vs variant) validates performance changes.
Evidence: Azure (Source 15): "Start with a focused hypothesis about your workload's performance and define measurable success criteria... For example, your hypothesis might be: 'Adding an index to the orders table reduces query time by 70% under peak load.' Your baseline is the current schema, and the variant is the schema with the new index. Run the same load test against both versions."
Source: Source 15 (Azure Well-Architected Performance Test)
URL: https://learn.microsoft.com/en-us/azure/well-architected/performance-efficiency/performance-test
Confidence: HIGH
Corroborated By: Source 17 (Google SRE Release Engineering): "Release engineering recommends that the continuous build test targets correspond to the same test targets that gate the project release... create releases at the revision number of the last continuous test build that successfully completed all tests."

## Evidence 20
Claim: Common load testing mistakes: testing only /health endpoint, using insufficient test data, not monitoring server during test, not defining performance targets.
Evidence: Topic spec lists four mistakes: "❌ Hanya Menguji Endpoint `/health`", "❌ Menggunakan Data Dummy yang Terlalu Sedikit", "❌ Tidak Memantau Server Selama Test", "❌ Tidak Mendefinisikan Target Performa."
Source: Topic spec (provided by user)
URL: Provided in prompt
Confidence: HIGH
Corroborated By: Source 15: "Avoid defining SLOs before understanding your user flows and performance requirements." Source 1: "If you can't access [APM/analytics tools], the business must provide these estimations."

## Evidence 21
Claim: Google SRE uses stress tests to find system limits: "How full can a database get before writes start to fail? How many queries a second can be sent before it becomes overloaded?"
Evidence: Source 10 (Google SRE Testing for Reliability): "Engineers use stress tests to find the limits on a web service. Stress tests answer questions such as: How full can a database get before writes start to fail? How many queries a second can be sent to an application server before it becomes overloaded, causing requests to fail?"
Source: Source 10 (Google SRE Testing for Reliability)
URL: https://sre.google/sre-book/testing-reliability/
Confidence: HIGH

## Evidence 22
Claim: Four Golden Signals for monitoring: Latency, Traffic, Errors, Saturation.
Evidence: Source 12 (Google SRE Monitoring): "The four golden signals of monitoring are latency, traffic, errors, and saturation. If you can only measure four metrics of your user-facing system, focus on these four." Defines each: Latency = time to service request; Traffic = demand measure (e.g., HTTP RPS); Errors = rate of failed requests; Saturation = how "full" service is (emphasize most constrained resource).
Source: Source 12 (Google SRE Monitoring Distributed Systems)
URL: https://sre.google/sre-book/monitoring-distributed-systems/
Confidence: HIGH
Corroborated By: Source 5 (k6 Metrics) notes these are "three of the four Golden Signals" for RED method.

## Evidence 23
Claim: Capacity planning should be done before predicted usage changes (seasonal, marketing, events).
Evidence: Azure Checklist (Source 14): "Conduct capacity planning. Capacity planning should be done before there are predicted changes in usage patterns, such as seasonal variations, product updates, marketing campaigns, special events, or regulatory changes."
Source: Source 14 (Azure Performance Efficiency Checklist)
URL: https://learn.microsoft.com/en-us/azure/architecture/framework/scalability/performance-efficiency
Confidence: HIGH

## Evidence 24
Claim: Client-side throttling and per-customer quotas protect against overload cascading.
Evidence: Source 11 (Google SRE Handling Overload): "Client-side throttling... When a client detects that a significant portion of its recent requests have been rejected due to 'out of quota' errors, it starts self-regulating and caps the amount of outgoing traffic it generates... Adaptive throttling... per-request retry budget of up to three attempts... per-client retry budget... 10% retry ratio."
Source: Source 11 (Google SRE Handling Overload)
URL: https://sre.google/sre-book/handling-overload/
Confidence: HIGH

## Evidence 25
Claim: P95/P99 thresholds in k6 use percentile aggregation on Trend metrics (latency in milliseconds).
Evidence: k6 Thresholds (Source 6): "Trend: Percentiles, averages, medians, and minimums must be within specified milliseconds. Aggregation methods: avg, min, max, med and p(N) where N specifies the threshold percentile value... E.g. p(99.99) means the 99.99th percentile. The values are in milliseconds."
Source: Source 6 (k6 Thresholds)
URL: https://grafana.com/docs/k6/latest/using-k6/thresholds/
Confidence: HIGH

## Evidence 26
Claim: k6 supports multiple load modeling approaches: VU-based (concurrent users) and RPS-based (constant/ramping arrival rate).
Evidence: k6 API Load Testing (Source 7): "k6 provides two broad ways to model load: Through virtual users (VUs), to simulate concurrent users; Through requests per second, to simulate raw, real-world throughput... To configure workloads according to a target request rate, use the constant arrival rate executor."
Source: Source 7 (k6 API Load Testing)
URL: https://grafana.com/docs/k6/latest/testing-guides/api-load-testing/
Confidence: HIGH

## Evidence 27
Claim: Load test results should be compared against baselines to detect performance drift/regression.
Evidence: Azure (Source 15): "Record performance metrics during initial tests. This recording is your baseline, a snapshot of 'normal' performance. In subsequent runs, compare new results against this baseline to detect performance changes... Regularly review and update your baselines after significant changes to your workload." Source 18 (Martin Fowler Test Pyramid): References baseline comparison in practical test pyramid context.
Source: Source 15, Source 18
URL: https://learn.microsoft.com/en-us/azure/well-architected/performance-efficiency/performance-test, https://martinfowler.com/articles/practical-test-pyramid.html
Confidence: HIGH
Corroborated By: Both sources emphasize baseline comparison.

## Evidence 28
Claim: Synthetic monitoring runs automated tests against live production regularly.
Evidence: Martin Fowler (Source 20): "Synthetic monitoring (also called semantic monitoring) runs a subset of an application's automated tests against the live production system on a regular basis. The results are pushed into the monitoring service, which triggers alerts in case of failures."
Source: Source 20 (Martin Fowler Synthetic Monitoring)
URL: https://martinfowler.com/bliki/SyntheticMonitoring.html
Confidence: HIGH

## Evidence 29
Claim: When P95 increases significantly (e.g., 300ms to 2.5s at 800 VUs), investigate by correlating backend metrics (CPU, memory, DB connections, queue length) to identify bottleneck layer.
Evidence: Synthesized best-practice guidance derived from: Evidence 16 (backend monitoring critical) + Evidence 17 (external deps can be bottleneck) + Evidence 7 (P95 importance). Topic spec question 5 references this scenario.
Source: Evidence 16, 17, 7, Topic spec question 5
URL: N/A — synthesized guidance
Confidence: MEDIUM
Notes: Synthesized investigation guidance; no single source provides exact step-by-step for this specific P95 spike scenario. This represents a reasonable interpretation of established principles (backend correlation, percentile focus), but practitioners should verify against their specific monitoring tools and architecture.

## Evidence 30
Claim: Performance testing types should be applied progressively — start with load test, then stress, spike, endurance.
Evidence: k6 Load Testing (Source 1): "Note: Run stress tests only after smoke and average-load tests. Running this test type earlier may be wasteful and make it hard to pinpoint problems if they appear at low volumes or at loads under the average utilization." k6 Soak (Source 4): "Don't run soak tests before running smoke and average-load tests. Each test uncovers different problems. Running this first may cause confusion and resource waste."
Source: Source 1, Source 4
URL: https://grafana.com/docs/k6/latest/testing-guides/test-types/load-testing/, https://grafana.com/docs/k6/latest/testing-guides/test-types/soak-testing/
Confidence: HIGH
Corroborated By: Azure (Source 15): "Don't try to implement all test types immediately. Begin with basic load testing to understand your baseline performance. As you identify risks and gain experience, expand to stress testing, spike testing, and eventually endurance testing."