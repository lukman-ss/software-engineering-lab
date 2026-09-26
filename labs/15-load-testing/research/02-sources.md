# Sources

## Source 1

Title: Load test types
Publisher: Grafana k6 Documentation
URL: https://grafana.com/docs/k6/latest/testing-guides/test-types/
Published: Continuous (latest version)
Accessed: 2026-09-26
Source Tier: Tier 1
Relevance: Authoritative documentation defining all six load test types (smoke, load, stress, soak, spike, breakpoint) with load patterns, duration recommendations, k6 configuration examples, and key considerations for each type.

## Source 2

Title: Average-load testing
Publisher: Grafana k6 Documentation
URL: https://grafana.com/docs/k6/latest/testing-guides/test-types/load-testing/
Published: Continuous (latest version)
Accessed: 2026-09-26
Source Tier: Tier 1
Relevance: Defines "average-load test" as simulation of typical production traffic, including ramp-up/plateau/ramp-down pattern, recommended stage durations (ramp-up 5-15% of total), and the importance of production-derived user counts.

## Source 3

Title: Stress testing
Publisher: Grafana k6 Documentation
URL: https://grafana.com/docs/k6/latest/testing-guides/test-types/stress-testing/
Published: Continuous (latest version)
Accessed: 2026-09-26
Source Tier: Tier 1
Relevance: Defines stress testing as load above average, emphasizing it should only be run after average-load tests pass, and that the load level depends on the system's risk profile (no fixed percentage like 50 or 100%).

## Source 4

Title: Spike testing
Publisher: Grafana k6 Documentation
URL: https://grafana.com/docs/k6/latest/testing-guides/test-types/spike-testing/
Published: Continuous (latest version)
Accessed: 2026-09-26
Source Tier: Tier 1
Relevance: Defines spike testing as sudden, massive traffic increase with minimal/no ramp-up, used for flash sales, product launches, and seasonal events. Emphasizes backend monitoring, key-process focus, and "run, tune, repeat" methodology.

## Source 5

Title: Soak testing
Publisher: Grafana k6 Documentation
URL: https://grafana.com/docs/k6/latest/testing-guides/test-types/soak-testing/
Published: Continuous (latest version)
Accessed: 2026-09-26
Source Tier: Tier 1
Relevance: Defines soak/endurance testing as average-load test extended over hours/days (typical values: 3, 4, 8, 12, 24, 48-72 hours), used to detect memory leaks, resource leaks, data saturation, and storage depletion.

## Source 6

Title: Breakpoint testing
Publisher: Grafana k6 Documentation
URL: https://grafana.com/docs/k6/latest/testing-guides/test-types/breakpoint-testing/
Published: Continuous (latest version)
Accessed: 2026-09-26
Source Tier: Tier 1
Relevance: Defines breakpoint testing as gradual load increase to find system limits/capacity, also known as capacity, point load, or limit testing. Notes it can find the cloud account bill limit if not turned off in elastic environments.

## Source 7

Title: Smoke testing
Publisher: Grafana k6 Documentation
URL: https://grafana.com/docs/k6/latest/testing-guides/test-types/smoke-testing/
Published: Continuous (latest version)
Accessed: 2026-09-26
Source Tier: Tier 1
Relevance: Defines smoke testing as minimal load test (2-20 VUs, seconds to minutes) run whenever a script is created/updated, before more extensive tests. Serves to validate script correctness and gather baseline metrics.

## Source 8

Title: Ramping arrival rate
Publisher: Grafana k6 Documentation
URL: https://grafana.com/docs/k6/latest/using-k6/scenarios/executors/ramping-arrival-rate/
Published: Continuous (latest version)
Accessed: 2026-09-26
Source Tier: Tier 1
Relevance: Documents the ramping-arrival-rate executor for open-model load, with stages config, startRate, timeUnit, preAllocatedVUs, maxVUs options. Critical for breakpoint testing where load must increase regardless of system degradation.

## Source 9

Title: Testing for Reliability (Chapter 17)
Publisher: Google SRE (Site Reliability Engineering) Book
URL: https://sre.google/sre-book/testing-reliability/
Published: 2017 (Copyright Google, Inc., published by O'Reilly Media)
Accessed: 2026-09-26
Source Tier: Tier 1
Relevance: Authoritative source on testing types (unit, integration, system/smoke, performance, regression) and production tests (stress test for finding system limits, canary test). Defines stress testing purpose: "how full can a database get before writes start to fail?" and "how many queries a second can be sent to an application server before it becomes overloaded?"

## Source 10

Title: ISO/IEC 25010
Publisher: International Organization for Standardization (ISO)
URL: https://en.wikipedia.org/wiki/ISO/IEC_25010
Published: 2011 (standard), Wikipedia article accessed for definition summary
Accessed: 2026-09-26
Source Tier: Tier 1 (standard) + Tier 2 (Wikipedia summary)
Relevance: Defines the Software Quality Model including "Performance Efficiency" with subcharacteristics: Time behaviour (response times, throughput), Resource utilization (CPU, memory, storage, network), and Capacity. Provides industry-standard definition of performance-related quality attributes.

## Source 11

Title: Architecture Strategies for Performance Testing
Publisher: Microsoft Azure Well-Architected Framework (Performance Efficiency Pillar)
URL: https://learn.microsoft.com/en-us/azure/well-architected/performance-efficiency/performance-test
Published: 2026-06-24 (last updated)
Accessed: 2026-09-26
Source Tier: Tier 1
Relevance: Comprehensive guidance on performance testing including test type definitions (load, stress, spike, endurance/soak), terminology (performance targets, thresholds, budgets, acceptance criteria), environment considerations (mirroring production), production testing strategies, and iterative optimization process.

## Source 12

Title: Performance Efficiency Pillar
Publisher: AWS Well-Architected Framework
URL: https://docs.aws.amazon.com/wellarchitected/latest/performance-efficiency-pillar/welcome.html
Published: November 6, 2024 (whitepaper edition)
Accessed: 2026-09-26
Source Tier: Tier 1
Relevance: Defines the five focus areas of performance efficiency: architecture selection, compute and hardware, data management, networking and content delivery, process and culture. Emphasizes that performance testing is part of the "process and culture" focus area.

## Source 13

Title: What is Locust?
Publisher: Locust Official Documentation
URL: https://docs.locust.io/en/stable/what-is-locust.html
Published: Continuous (latest version 2.17+)
Accessed: 2026-09-26
Source Tier: Tier 1
Relevance: Authoritative source describing Locust's architecture: write tests in Python, run every user in its own greenlet (lightweight process/coroutine via gevent), distributed and scalable supporting hundreds of thousands of concurrent users, web-based UI for real-time monitoring.

## Source 14

Title: Your first test
Publisher: Locust Official Documentation
URL: https://docs.locust.io/en/stable/quickstart.html
Published: Continuous (latest version)
Accessed: 2026-09-26
Source Tier: Tier 1
Relevance: Documents Locust's HttpUser class, task decorator, wait_time functions, user weight system, and how to define user behavior. Shows how to simulate realistic concurrent users with Think time between tasks.

## Source 15

Title: A Collection of Best Practices for Production Services
Publisher: Google SRE Book (Appendix B)
URL: https://sre.google/sre-book/service-best-practices/
Published: 2017 (Copyright Google, Inc., published by O'Reilly, Media)
Accessed: 2026-09-26
Source Tier: Tier 1
Relevance: Contains SRE-recommended best practices for production services including performance testing, load testing, and capacity planning guidelines. Provides production-readiness checklist.

## Source 16

Title: Introduction to the Spring IoC Container and Beans
Publisher: VMware Spring Framework Documentation
URL: https://docs.spring.io/spring-framework/reference/core/beans/introduction.html
Published: Continuous (v7.0.x)
Accessed: 2026-09-26
Source Tier: Tier 1
Relevance: While focused on DI, this was referenced in prior lab research; not directly relevant to load testing topic. Excluded from active evidence.

## Source 17

Title: Smoke testing (k6)
Publisher: Grafana k6 Documentation
URL: https://grafana.com/docs/k6/latest/testing-guides/test-types/smoke-testing/
Published: Continuous
Accessed: 2026-09-26
Source Tier: Tier 1
Relevance: Defines smoke test as running with 2-20 VUs for seconds to minutes, validating scripts before running larger tests, gathering baseline metrics.

## Source 18

Title: Calculate concurrent users for load tests
Publisher: Grafana k6 Documentation
URL: https://grafana.com/docs/k6/latest/testing-guides/calculate-concurrent-users/
Published: Continuous
Accessed: 2026-09-26
Source Tier: Tier 1
Relevance: Provides methodology for calculating concurrent VUs from real production traffic data (peak sessions per second × average session duration). Critical for answering the "how many VUs" question for the Booking Bengkel scenario.

## Source 19

Title: Performance testing guidance
Publisher: Microsoft Azure Documentation (Performance Efficiency Pillar)
URL: https://learn.microsoft.com/en-us/azure/architecture/framework/scalability/load-testing/
Published: 2025 or later
Accessed: 2026-09-26
Source Tier: Tier 1
Relevance: Azure's performance testing and load testing overview covering methodology, tools, metrics (latency, throughput, error rate, resource utilization), and best practices for performance testing in cloud environments.

## Source 20

Title: Application and infrastructure monitoring with k6
Publisher: Grafana k6 Documentation
URL: https://grafana.com/docs/k6/latest/testing-guides/test-types/load-testing/#results-analysis
Published: Continuous
Accessed: 2026-09-26
Source Tier: Tier 1
Relevance: Discusses results analysis methodology including interpreting trends across load stages.

## Source 21

Title: What is HTTP and how does it work
Publisher: Mozilla Developer Network (MDN) — Web docs
URL: https://developer.mozilla.org/en-US/docs/Web/HTTP/Overview
Published: 2026-08
Accessed: 2026-09-26
Source Tier: Tier 2
Relevance: Provides foundational understanding of HTTP request lifecycle for interpreting HTTP-level load test results (status codes, connection limits, timeouts).
