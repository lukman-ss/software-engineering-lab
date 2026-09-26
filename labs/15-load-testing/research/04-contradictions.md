# Contradictions

## Contradiction 1
Claim: "k6 is a load testing tool optimized for API testing" vs "Azure Load Testing supports JMeter and Locust but NOT k6."

Source A (k6 Docs): "Grafana k6 is an open-source, developer-friendly, and extensible performance testing tool... k6 is optimized for minimal resource consumption and designed for running high-load performance tests such as spike, stress, or soak tests." (Source 1, 7)

Source B (Azure Load Testing): "Azure Load Testing currently does not support other testing frameworks than Apache JMeter and Locust." (Source 16)

ASSESSMENT: Not contradictory — they describe different contexts. k6 is a standalone open-source tool that anyone can use. Azure Load Testing is a specific Azure cloud service that only supports JMeter and Locust (no k6). A team can choose to run k6 self-hosted OR use Azure's JMeter/Locust-managed service.

## Contradiction 2
Claim: "Test with mocked external dependencies for faster, predictable tests" vs "Test with REAL external API calls to understand end-to-end latency."

Source A (general best practice implied in many guides): Many load testing guides suggest using mocks/stubs to isolate the system under test and reduce noise.

Source B (Azure Well-Architected): "When testing under load, include actual third-party API calls. Mocking external dependencies makes tests run faster and more predictably, but it hides real-world performance problems. If your app depends on a payment processor API, test with real calls to understand end-to-end latency." (Source 15)

ASSESSMENT: Resolved by contextual qualification. Real calls are essential for production-like validation in controlled sandbox environments; mocks/stubs with artificial latency are appropriate for high-volume stress tests to avoid rate limits, financial costs, and ToS violations. Research report now clarifies this distinction (Finding 7 revised 2026-09-26).

## Contradiction 3
Claim: "Use laptop/test environment that mirrors production" vs "Only way to get accurate performance data is production testing."

Source A (Azure Well-Architected PE:03): "Your test environment should mirror production as closely as practical... For mission-critical workloads, match production exactly across: Compute SKUs and configurations, Autoscaling settings, Caching configurations, Network conditions, External dependencies." (Source 15)

Source B (Azure Well-Architected PE:06): "Test environments can't fully replicate real-world conditions that affect performance... Production tests expose problems that only surface under actual usage and provide accurate baselines for future optimization." (Source 15)

ASSESSMENT: Not contradictory — layered approach. Source A advocates for staging/test environment that mirrors production as closely as possible. Source B acknowledges staging can't replicate everything (real user patterns, true network latency, real traffic mix). Both recommend production testing as a final validation step. The guidance is: start with staging, then run controlled tests in production (off-peak) to validate.

## Contradiction 4
Claim: "Stress test should be much higher than average (50-100% or more)" vs "Load should be higher than average but no fixed percentage."

Source A (k6 Stress Testing): "Load should be higher than what the system experiences on average. Some testers might have default targets for stress tests—say an increase upon average load by 50 or 100 percent—there's no fixed percentage... The load simulated in a Stress test depends on the stressful situations that the system may be subject to. Sometimes this may be just a few percentage points above that average. Other times, it can be 50 to 100% higher, as mentioned. Some stressful situations can be twice, triple, or even orders of magnitude higher." (Source 2)

ASSESSMENT: Not contradictory — same source clarifies both points. The guide acknowledges that while some teams use 50-100% as rule-of-thumb, the actual percentage depends on the specific use case and risk scenarios (a few % to orders of magnitude).

## Contradiction 5
Claim: "k6 uses JavaScript" vs "Locust uses Python" vs "JMeter uses Java/XML" — are these really different or same?

Source A (k6): Uses JavaScript (Source 1, 5, 6, 7)
Source B (Locust): Uses Python (Source 8, 9)
Source C (JMeter): Uses Java/XML (Source 16, 22)

ASSESSMENT: Not contradictory — they are different tools with different language ecosystems. The topic specification lists all four (k6, JMeter, Locust, Gatling) as popular tools, implying they serve similar purposes with different tradeoffs. Each tool's language choice is intentional and affects usability:
- k6: JavaScript — web developers familiar with JS can write tests
- Locust: Python — popular language for developers, easy scripting
- JMeter: Java/XML — traditional enterprise testing tool, GUI-based

## Contradiction 6
Claim: "Average-load test plateau should be 5x ramp-up duration" vs "Soak test duration should be hours/days."

Source A (k6 Load Testing): "Aim for an average duration at least five times longer than the ramp-up to assess the performance trend over a significant period of time." (Source 1)

Source B (k6 Soak Testing): "Some typical values are 3, 4, 8, 12, 24, and 48 to 72 hours." (Source 4)

ASSESSMENT: Not contradictory — they describe different test types with different goals. Average-load test plateau is minutes-to-hours; soak test is hours-to-days. Source B explicitly notes the key difference: "The soak test differs from an average-load test in test duration." A soak test uses average load but extends the plateau dramatically (hours/days vs minutes).

## No Material Contradictions Discovered

After thorough review, all apparent contradictions are resolved:
- Context differences (standalone tool vs managed service)
- Complementary approaches (mock vs real dependencies for different phases)
- Layered strategy (staging + controlled production)
- Clarified guidance (no fixed percentage but range guidance)
- Different tools with different language choices
- Different test types with different duration goals

All sources converge on the same core principles:
1. Multiple test types serve different purposes (load, stress, spike, soak)
2. Backend monitoring is critical for bottleneck identification
3. Realistic data and realistic environments (including real dependencies) are essential
4. Percentiles (P95/P99) matter more than averages
5. Thresholds should be based on SLOs and codified pass/fail criteria