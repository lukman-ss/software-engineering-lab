# Sources

## Source 1
Title: Circuit Breaker Pattern — Martin Fowler
Publisher: Martin Fowler
URL: https://martinfowler.com/bliki/CircuitBreaker.html
Published: 2014-12-23
Accessed: 2026-09-25
Source Tier: Tier 2 (expert technical article)
Relevance: Primary definition of the pattern, states, and transitions

## Source 2
Title: Hystrix Circuit Breaker — Netflix
Publisher: Netflix
URL: https://github.com/Netflix/Hystrix/wiki/How-it-Works
Published: 2012-2018 (archived)
Accessed: 2026-09-25
Source Tier: Tier 1 (original implementation documentation)
Relevance: Original production implementation details, configuration parameters

## Source 3
Title: Circuit Breaker Pattern — Microsoft Azure Architecture Center
Publisher: Microsoft
URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker
Published: 2023-08-29
Accessed: 2026-09-25
Source Tier: Tier 1 (official cloud provider documentation)
Relevance: Authoritative pattern description, states, configuration, considerations

## Source 4
Title: AWS Well-Architected Framework — Reliability Pillar
Publisher: Amazon Web Services
URL: https://docs.aws.amazon.com/wellarchitected/latest/reliability-pillar/
Published: 2024
Accessed: 2026-09-25
Source Tier: Tier 1 (official cloud provider framework)
Relevance: Production best practices for resilience, circuit breaker usage

## Source 5
Title: Google SRE Book — Handling Overload
Publisher: Google
URL: https://sre.google/sre-book/handling-overload/
Published: 2016
Accessed: 2026-09-25
Source Tier: Tier 1 (original SRE practices)
Relevance: Circuit breaker in context of overload protection, cascade failure

## Source 6
Title: Retry with Exponential Backoff — AWS Architecture Blog
Publisher: Amazon Web Services
URL: https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/
Published: 2015-03-10
Accessed: 2026-09-25
Source Tier: Tier 1 (official cloud provider engineering blog)
Relevance: Retry strategies, jitter, interaction with circuit breaker

## Source 7
Title: gobreaker — Sony
Publisher: Sony
URL: https://github.com/sony/gobreaker
Published: 2016-present
Accessed: 2026-09-25
Source Tier: Tier 1 (Go library source code)
Relevance: Go implementation reference, state machine, thread safety

## Source 8
Title: Circuit — Go Resilience (cep21)
Publisher: cep21
URL: https://github.com/cep21/circuit
Published: 2019-present
Accessed: 2026-09-25
Source Tier: Tier 1 (Go library source code)
Relevance: Modern Go implementation, configuration options, metrics

## Source 9
Title: Bulkhead Pattern — Microsoft Azure Architecture Center
Publisher: Microsoft
URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/bulkhead
Published: 2023-08-29
Accessed: 2026-09-25
Source Tier: Tier 1 (official cloud provider documentation)
Relevance: Distinction from circuit breaker, isolation patterns

## Source 10
Title: Managing Load — Google SRE Workbook
Publisher: Google
URL: https://sre.google/workbook/managing-load/
Published: 2018
Accessed: 2026-09-25
Source Tier: Tier 1 (original SRE practices)
Relevance: Load shedding vs circuit breaker, priority-based shedding, load balancing interaction, Dressy case study

## Source 11
Title: Timeout and Retry Patterns — Microsoft Azure
Publisher: Microsoft
URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/retry
Published: 2023-08-29
Accessed: 2026-09-25
Source Tier: Tier 1 (official cloud provider documentation)
Relevance: Timeout vs retry vs circuit breaker interaction

## Source 12
Title: Timeouts, retries, and backoff with jitter — AWS Builder's Library
Publisher: Amazon Web Services
URL: https://aws.amazon.com/builders-library/timeouts-retries-and-backoff-with-jitter/
Published: 2020
Accessed: 2026-09-25
Source Tier: Tier 1 (official cloud provider engineering guidance)
Relevance: Retry amplification and retry storms under dependency degradation, backoff and jitter to prevent synchronized retry spikes

## Source 13
Title: Queue-Based Load Leveling — Microsoft Azure Architecture Center
Publisher: Microsoft
URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/queue-based-load-leveling
Published: 2023-08-29
Accessed: 2026-09-25
Source Tier: Tier 1 (official cloud provider documentation)
Relevance: Decoupling non-critical flows via queues, buffering to enable fallback/degraded behavior and support idempotent processing