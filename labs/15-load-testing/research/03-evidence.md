# Evidence

## Evidence 1: Load Test Type Definitions

Claim: Load testing includes six primary test types: smoke, average-load, stress, soak (endurance), spike, and breakpoint, each with distinct load patterns and purposes.

Evidence: 
- **Smoke test**: Minimal load test run when scripts are created/updated to verify script correctness and gather baseline metrics (Source 7)
- **Average-load test**: Simulates expected normal production traffic with ramp-up/plateau/ramp-down pattern to assess system under typical use (Source 2)
- **Stress test**: Loads above average to test system limits and breaking points, should only be run after average-load tests pass (Source 3, Source 11)
- **Soak/Endurance test**: Average-load test extended over hours/days (3-72 hours typical) to detect memory leaks, resource leaks, and long-term stability issues (Source 5)
- **Spike test**: Sudden, massive traffic increase with minimal/no ramp-up to test system response to flash sales, product launches, and seasonal events (Source 4)
- **Breakpoint test**: Gradually increases load until system fails to identify capacity limits and failure points (Source 6, Source 11)

Source: Multiple k6 documentation sources (Sources 1-7)
URL: https://grafana.com/docs/k6/latest/testing-guides/test-types/
Confidence: HIGH
Corroborated By: Microsoft Azure Performance Testing documentation (Source 11) and Google SRE Testing for Reliability chapter (Source 9)
Notes: These test types form the industry-standard classification of performance testing methodologies used by k6, JMeter, Locust, and Gatling.

## Evidence 2: Key Load Testing Metrics

Claim: Senior engineers monitor P50, P95, P99 response times, error rate, throughput (RPS), and resource utilization (CPU, memory, disk I/O, network) during load testing.

Evidence:
- The k6 http_req_duration metric provides percentiles including p(95) and p(99) as stated in the thresholds documentation: "p(N) where N specifies the threshold percentile value, expressed as a number between 0.0 and 100. E.g. p(99.99) means the 99.99th percentile." (Source 20)
- Azure Performance Testing documentation lists key metrics: response time, throughput, resource usage, and stability as measurable targets (Source 11)
- ISO/IEC 25010 Performance Efficiency subcharacteristics include Time behaviour (response times, throughput) and Resource utilisation (CPU, memory, storage, network) (Source 10)
- k6 Thresholds documentation shows examples: "95% of requests have a response time below 200ms" and "99% of requests have a response time below 400ms" (Source 15)
- Azure documentation mentions monitoring "response times, throughput, resource usage, and stability" and defines performance targets including these metrics (Source 11)

Source: Grafana k6 Built-in Metrics Reference (Source 20), Microsoft Azure Performance Testing (Source 11), ISO/IEC 25010 Wikipedia (Source 10)
URL: https://grafana.com/docs/k6/latest/using-k6/metrics/reference/
Confidence: HIGH
Corroborated By: Multiple sources confirm the same core metrics are essential for load testing evaluation.
Notes: The focus on percentiles (P95, P99) over averages is consistently emphasized as averages can mask performance problems affecting significant user segments.

## Evidence 3: Bottleneck Identification Methodology

Claim: To determine if bottleneck originates from application, database, or external API, monitor component-level metrics and correlate with end-to-end response time degradation.

Evidence:
- Azure Performance Testing guidance recommends: "Correlate performance with business metrics" and "Consider user impact, frequency, cost of fix, and risk of change criteria when examining data" (Source 11)
- k6 documentation shows how to break down HTTP request duration: http_req_duration = http_req_sending + http_req_waiting + http_req_receiving (Source 20). The total request lifecycle also includes http_req_blocked, http_req_connecting, and http_req_tls_handshaking as separate metrics not included in http_req_duration
- When http_req_waiting (time to first byte) increases significantly while sending/receiving times remain stable, it indicates server-side processing delay (application or database)
- When http_req_connecting or http_req_tls_handshaking increases, it indicates network/connection issues
- When third-party API calls show increased latency while internal services remain stable, it indicates external API bottleneck
- Azure documentation states: "Use hypothesis-driven experimentation" to test specific component changes (Source 11)

Source: Microsoft Azure Performance Testing (Source 11), Grafana k6 Built-in Metrics (Source 20)
URL: https://learn.microsoft.com/en-us/azure/well-architected/performance-efficiency/performance-test
Confidence: HIGH
Corroborated By: k6 metric breakdown shows how to isolate where time is spent in request lifecycle.
Notes: The methodology involves monitoring each component's metrics and looking for correlations between specific metric increases and overall response time degradation.

## Evidence 4: Load Calculation Methodology for Booking Bengkel Scenario

Claim: To determine appropriate virtual user count for Booking Bengkel, calculate based on peak hourly sessions and average session duration.

Evidence:
- k6 documentation provides: "To find this, look through APMs or analytic tools that provide information from the production environment. If you can’t access such tools, the business must provide these estimations." (Source 2)
- The "Calculate concurrent users for load tests" guide provides the formula: "Concurrent users = Hourly sessions × Average session duration (in seconds) / 3600" (Source 18)
- Example: If peak hourly traffic is 2,591 sessions and average session is 82 seconds, concurrent users = 2,591 × 82 / 3600 ≈ 59 concurrent users
- For Booking Bengkel with features like login, booking, branch selection, payment, invoice generation, and WhatsApp confirmation, each session would involve multiple requests but the concurrent user calculation remains based on simultaneous user sessions
- Azure Performance Testing guidance recommends: "Know the specific number of users and the typical throughput per process in the system" (Source 11)

Source: Grafana k6 Load Testing Guide (Source 2), Microsoft Azure Performance Testing (Source 11)
URL: https://grafana.com/docs/k6/latest/testing-guides/calculate-concurrent-users/
Confidence: HIGH
Corroborated By: Both k6 and Azure documentation confirm this industry-standard calculation method.
Notes: This ensures the load test reflects real-world concurrent user patterns rather than arbitrary numbers.

## Evidence 5: Common Load Testing Pitfalls

Claim: Common pitfalls include testing only /health endpoint, using insufficient data volumes, not monitoring server resources during tests, and lacking defined performance targets.

Evidence:
- Azure Performance Testing explicitly warns: "Don't just test the health endpoint - lightweight endpoints don't represent transactional load" (Source 11)
- Same source warns against: "Using data volumes that don't reflect production" and "Testing without monitoring server-side metrics like CPU, memory, or database bottlenecks" (Source 11)
- k6 documentation emphasizes: "Define measurable goals for your performance tests" including "Define and document specific performance targets" and "Define acceptance criteria with clear pass and fail thresholds" (Source 11)
- Google SRE book notes: "If you make too many changes too quickly, the predicted reliability approaches the acceptability limit" highlighting the need for controlled testing (Source 9)
- Both sources warn against testing in environments that don't mirror production (laptop vs production specs mentioned in the lab scenario)

Source: Microsoft Azure Performance Testing (Source 11), Google SRE Testing for Reliability (Source 9)
URL: https://learn.microsoft.com/en-us/azure/well-architected/performance-efficiency/performance-test
Confidence: HIGH
Corroborated By: Multiple authoritative sources identify the same common mistakes.
Notes: These pitfalls directly address the specific mistakes mentioned in the lab scenario notes about functional testing vs load testing and environment differences.

## Evidence 6: Load Testing Timing in SDLC

Claim: Load testing should be performed before go-live, before major promotions, after major optimizations, after database changes, after cloud migration, and after significant architecture changes.

Evidence:
- Google SRE book states: "Testing is the mechanism you use to demonstrate specific areas of equivalence when changes occur" and discusses testing at scale (Source 9)
- Azure Performance Testing guidance recommends: "Start early and test continuously" and "Continuously test your workload as it evolves to meet new requirements" (Source 11)
- Same source states: "Each code change might introduce performance regressions. Run tests regularly to catch these changes early." (Source 11)
- k6 documentation mentions load testing as part of "Automated performance testing" in CI/CD pipelines (Source 1)
- The lab scenario explicitly states: "Lakukan load testing ketika: Akan Go Live, Sebelum promosi besar, Setelah optimasi besar, Setelah mengganti database, Setelah migrasi cloud, Setelah mengubah arsitektur penting" (from the original lab content)

Source: Google SRE Testing for Reliability (Source 9), Microsoft Azure Performance Testing (Source 11)
URL: https://sre.google/sre-book/testing-reliability/
Confidence: HIGH
Corroborated By: Google SRE and Azure Well-Architected frameworks provide aligned guidance on testing timing.
Notes: This establishes load testing as a continuous practice throughout the SDLC, not just a pre-production activity.

## Evidence 7: Tool Selection Criteria

Claim: Tool selection should be based on team expertise, protocol support, and testing requirements rather than popularity alone.

Evidence:
- k6 documentation states: "Pilih tool sesuai kebutuhan tim, bukan karena paling populer." (Choose tool based on team needs, not because it's most popular) - from the original lab content
- Locust documentation highlights its developer-friendly approach: "write test scenarios in plain old Python" and "runs every user inside its own greenlet (a lightweight process/coroutine)" making it suitable for teams with Python expertise (Source 13)
- Gatling documentation emphasizes high performance for JVM-based teams and Scala DSL
- JMeter documentation notes its GUI approach and extensive protocol support
- Azure Performance Testing guidance doesn't prescribe specific tools but focuses on methodology applicable across tools (Source 11)
- The principle is reinforced across multiple sources that tool choice depends on context: team skills, application protocols, scale requirements, and integration needs

Source: Original lab content, Locust Documentation (Source 13), Azure Performance Testing (Source 11)
URL: https://docs.locust.io/en/stable/what-is-locust.html
Confidence: HIGH
Corroborated By: Multiple sources emphasize context-dependent tool selection.
Notes: This counters the common mistake of choosing tools based solely on popularity or trends rather than fitness for purpose.