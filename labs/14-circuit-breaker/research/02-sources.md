# Sources

## Source 1
Title: Circuit Breaker
Publisher: Martin Fowler
URL: https://martinfowler.com/bliki/CircuitBreaker.html
Published: 6 March 2014
Accessed: 2026-09-26
Source Tier: 1 (Primary / Original Pattern Author)
Relevance: Original definition of Circuit Breaker pattern, CLOSED/OPEN/HALF-OPEN states, failure threshold, reset timeout, state transitions, and monitoring recommendations.

## Source 2
Title: Circuit Breaker Pattern
Publisher: Microsoft Azure Architecture Center
URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker
Published: 2025-02-05
Accessed: 2026-09-26
Source Tier: 1 (Official Cloud Documentation)
Relevance: Detailed specification of Closed/Open/Half-Open states, failure counter reset behavior, time-based vs count-based windows, half-open success counter, when-to-use/when-not-to-use guidance, concurrency considerations.

## Source 3
Title: Retry Pattern
Publisher: Microsoft Azure Architecture Center
URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/retry
Published: 2024-07-18
Accessed: 2026-09-26
Source Tier: 1 (Official Cloud Documentation)
Relevance: Distinguishes Retry from Circuit Breaker, describes retry strategies (cancel/immediate/delayed), idempotency concerns, exponential backoff guidance, interaction patterns with Circuit Breaker.

## Source 4
Title: Bulkhead Pattern
Publisher: Microsoft Azure Architecture Center
URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/bulkhead
Published: 2026-03-19
Accessed: 2026-09-26
Source Tier: 1 (Official Cloud Documentation)
Relevance: Explains isolation via pools (thread pools, connection limits, CPU slices), distinguishes Bulkhead from Circuit Breaker, describes how combining patterns provides comprehensive resilience.

## Source 5
Title: CircuitBreaker
Publisher: Resilience4j
URL: https://resilience4j.readme.io/docs/circuitbreaker
Accessed: 2026-09-26
Source Tier: 1 (Reference Implementation Documentation)
Relevance: Implementation details: sliding windows (count-based/time-based), failureRateThreshold, permittedNumberOfCallsInHalfOpenState, waitDurationInOpenState, thread safety via AtomicReference, distinction between circuit breaker and bulkhead.

## Source 6
Title: Handling Overload
Publisher: Google (SRE Book)
URL: https://sre.google/sre-book/handling-overload/
Published: 2017
Accessed: 2026-09-26
Source Tier: 1 (Authoritative Industry Reference)
Relevance: Cascading failure mechanics, client-side throttling (adaptive throttling), criticality levels (CRITICAL_PLUS, CRITICAL, SHEDDABLE_PLUS, SHEDDABLE), retry budgets, per-customer quotas, load shedding strategies.

## Source 7
Title: Addressing Cascading Failures
Publisher: Google (SRE Book)
URL: https://sre.google/sre-book/addressing-cascading-failures/
Published: 2017
Accessed: 2026-09-26
Source Tier: 1 (Authoritative Industry Reference)
Relevance: Detailed treatment of cascading failure propagation, blast radius containment, bulkhead patterns in distributed systems, load shedding and circuit breaking as complementary mechanisms.

## Source 8
Title: Timeouts, Retries, and Backoff with Jitter
Publisher: AWS Builders Library
URL: https://aws.amazon.com/builders-library/timeouts-retries-and-backoff-with-jitter/
Accessed: 2026-09-26
Source Tier: 1 (Authoritative Industry Reference)
Relevance: Full jitter algorithm, exponential backoff with decorrelation, timeout sizing guidance, retry budget concept, distinction between transient and non-transient failures.

## Source 9
Title: Circuit Breaker Pattern: When to Use and How It Works
Publisher: Microsoft (Archived Content)
URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker
Published: 2025-02-05
Accessed: 2026-09-26
Source Tier: 1 (Official Cloud Documentation)
Relevance: Cross-referenced with Source 2; confirmed consistency across Microsoft documentation versions.

## Source 10
Title: Go HTTP Client Timeout Behavior
Publisher: Go Standard Library
URL: https://pkg.go.dev/net/http#Client
Accessed: 2026-09-26
Source Tier: 1 (Official Documentation)
Relevance: Go http.Client.Timeout documentation; explicit timeout field vs idle conn timeout; no default timeout means requests can hang indefinitely without explicit configuration.