# Evidence

## Evidence 1
Claim: Load testing validates system behavior under expected and peak workloads.
Evidence: "Load testing is a subset of performance testing that generally looks for how a system responds to normal and peak usage. You're looking for slow response times, errors, crashes, and other issues to determine how many users and transactions the system can accommodate..."
Source: Types of load testing (Grafana Labs)
URL: https://k6.io/docs/test-types/
Confidence: HIGH
Corroborated By: Microsoft Learn Azure Load Testing Documentation
Notes: Confirms fundamental definition and goal.

## Evidence 2
Claim: Starting with a smoke test before large-scale load testing is a best practice.
Evidence: "Start with a smoke test. Before beginning larger tests, validate that your load testing scripts work as expected and that your system performs well with a few users."
Source: Types of load testing (Grafana Labs)
URL: https://k6.io/docs/test-types/
Confidence: HIGH
Corroborated By: Standard testing methodologies.
Notes: Answers the question on how many virtual users to use initially (a minimal amount to verify logic).

## Evidence 3
Claim: Performance testing requires monitoring both client-side and server-side metrics to find bottlenecks.
Evidence: "Client-side metrics give you details reported by the test engine... request response time, or the number of requests per second. Server-side metrics provide information about your Azure application components... type of HTTP responses, or container resource consumption."
Source: What is Azure Load Testing?
URL: https://learn.microsoft.com/en-us/azure/load-testing/overview-what-is-azure-load-testing
Confidence: HIGH
Corroborated By: Grafana Labs
Notes: Helps isolate application vs DB vs third-party API issues.

## Evidence 4
Claim: Stress testing evaluates the system by increasing load above normal levels to find failure points.
Evidence: "Stress tests verify the stability and reliability of the system under heavier than normal usage... This test determines how much the performance degrades with the extra load and whether the system survives it."
Source: Types of load testing (Grafana Labs)
URL: https://k6.io/docs/test-types/
Confidence: HIGH
Corroborated By: General performance engineering standards.
Notes: Differentiates stress testing from standard load testing.
