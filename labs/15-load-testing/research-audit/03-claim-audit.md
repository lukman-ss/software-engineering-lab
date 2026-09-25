# Claim Audit

## Claim 1

Claim: Load testing validates system behavior under expected and peak workloads.
Location: research/03-evidence.md (Evidence 1)
Evidence Provided: Direct quote from Grafana k6 docs.
Source: https://k6.io/docs/test-types/
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Supported by primary source.

## Claim 2

Claim: Starting with a smoke test before large-scale load testing is a best practice.
Location: research/03-evidence.md (Evidence 2)
Evidence Provided: Direct quote from Grafana k6 docs.
Source: https://k6.io/docs/test-types/
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Clear methodology standard.

## Claim 3

Claim: Performance testing requires monitoring both client-side and server-side metrics to find bottlenecks.
Location: research/03-evidence.md (Evidence 3)
Evidence Provided: Direct quote from Microsoft Learn Azure Load Testing.
Source: https://learn.microsoft.com/en-us/azure/load-testing/overview-what-is-azure-load-testing
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Accurately reflected in documentation.

## Claim 4

Claim: Stress testing evaluates the system by increasing load above normal levels to find failure points.
Location: research/03-evidence.md (Evidence 4)
Evidence Provided: Direct quote from Grafana k6 docs.
Source: https://k6.io/docs/test-types/
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Matches test-type definitions.

## Claim 5

Claim: Average response time is misleading; percentiles (P95, P99) are essential for discovering latency issues in real workloads.
Location: research/05-report.md (Finding 1)
Evidence Provided: Explanation of outlier dilution by mathematical averages.
Source: Grafana k6 Documentation / Industry Best Practices
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Core performance engineering principle.

## Claim 6

Claim: Incremental testing starting with smoke tests prevents invalidating test results.
Location: research/05-report.md (Finding 2)
Evidence Provided: k6 smoke test guidelines (2-5 VUs).
Source: https://k6.io/docs/test-types/
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Standard practice.

## Claim 7

Claim: Comprehensive bottleneck detection requires correlating client-side latency with server-side resource metrics.
Location: research/05-report.md (Finding 3)
Evidence Provided: Microsoft Azure Load Testing engine-side vs server-side metrics overview.
Source: https://learn.microsoft.com/en-us/azure/load-testing/overview-what-is-azure-load-testing
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Standard observability correlation.

## Claim 8

Claim: Specific plan and diagnosis patterns for the Booking Bengkel scenario (critical endpoints, VU ranges, bottleneck isolation, and P95 spike investigation).
Location: research/05-report.md (Finding 4)
Evidence Provided: Logical derivation based on performance testing standards.
Source: Grafana Labs & Microsoft Learn Best Practices
Source Actually Supports Claim: PARTIAL
Classification: INTERPRETATION
Severity: LOW
Notes: Booking Bengkel is a specific case study; recommendations derive logically from best practices rather than being explicitly cited inside Grafana/Azure docs. This is properly contextualized as an analytical plan.
