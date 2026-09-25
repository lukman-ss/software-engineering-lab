# Sources

## Source 1
Title: Circuit Breaker
Publisher: Martin Fowler (martinfowler.com)
URL: https://martinfowler.com/bliki/CircuitBreaker.html
Published: 2014-03-06
Accessed: 2026-09-25
Source Tier: Tier 1 (Foundational Software Architecture Reference)
Relevance: Formalizes the pattern, introduces three states (CLOSED, OPEN, HALF_OPEN), and defines fail-fast and probe semantics.

## Source 2
Title: Circuit Breaker Pattern - Azure Architecture Center
Publisher: Microsoft Learn
URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker
Published: 2025-02-05 (Updated 2026-09-19)
Accessed: 2026-09-25
Source Tier: Tier 1 (Cloud Vendor Reference Architecture)
Relevance: Detailed state machine diagrams, interaction with Retry pattern, half-open traffic rate limiting, observability, and transient fault separation.

## Source 3
Title: Timeouts, retries, and backoff with jitter
Publisher: AWS Architecture Center / AWS Builder's Library
URL: https://aws.amazon.com/builders-library/timeouts-retries-and-backoff-with-jitter/
Published: 2019-11-01
Accessed: 2026-09-25
Source Tier: Tier 1 (Cloud Vendor Reliability Engineering Guide)
Relevance: Analyzes retry storms, exponential backoff, jitter, client timeout bounding, and how retries without circuit protection amplify downstream outages.
