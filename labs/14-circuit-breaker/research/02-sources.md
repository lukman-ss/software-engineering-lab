# Sources

## Source 1
Title: Circuit Breaker
Publisher: Martin Fowler (martinfowler.com)
URL: https://martinfowler.com/bliki/CircuitBreaker.html
Published: 2014-03-06
Accessed: 2026-09-25
Source Tier: Tier 2 (expert technical article, cites Nygard Release It!)
Relevance: Canonical pattern definition; three states (closed/open/half-open), failure counting, self-resetting probe

## Source 2
Title: How it Works — Hystrix Wiki
Publisher: Netflix / Hystrix
URL: https://github.com/Netflix/Hystrix/wiki/How-it-Works
Published: 2017-07-03 (last wiki edit by Matt Jacobs, 32 revisions)
Accessed: 2026-09-25
Source Tier: Tier 1 (original production implementation)
Relevance: Circuit-breaker tripping logic (request volume threshold + error percentage), sleep window, half-open single-probe, bulkhead via thread pools/semaphores, fallback

## Source 3
Title: Circuit Breaker pattern — Azure Architecture Center
Publisher: Microsoft Learn
URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker
Published: 2025-02-05 (ms.date; updated 2026-09-26)
Accessed: 2026-09-25
Source Tier: Tier 1 (official cloud provider documentation)
Relevance: Three-state machine, fail-fast in Open, Half-Open limited probes, failure counter time-based reset, monitoring, manual override, accelerated breaking on 429/503

## Source 4
Title: CircuitBreaker — gobreaker
Publisher: Sony (github.com/sony/gobreaker)
URL: https://github.com/sony/gobreaker
Published: 2016–present
Accessed: 2026-09-25
Source Tier: Tier 1 (Go library source)
Relevance: Go reference impl; Settings {Name, MaxRequests, Interval, Timeout, ReadyToTrip, OnStateChange}; mutex-protected state ; Half-Open MaxRequests probes; Interval rolling window

## Source 5
Title: Handling Overload (SRE Book Ch. 21)
Publisher: Google SRE
URL: https://sre.google/sre-book/handling-overload/
Published: 2016 (SRE Book first edition)
Accessed: 2026-09-25
Source Tier: Tier 1 (SRE practices)
Relevance: Client-side throttling, per-customer limits, criticality, cascading failure from retry amplification, shedding vs caller-side protection

## Source 6
Title: Exponential Backoff And Jitter — AWS Architecture Blog
Publisher: Amazon Web Services (Marc Brooker)
URL: https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/
Published: 2015-03-04 (update note 2023-05)
Accessed: 2026-09-25
Source Tier: Tier 1 (cloud provider engineering blog)
Relevance: Full/equal/decorrelated jitter, retry storm mitigation, proof jitter reduces contention vs plain exponential backoff

## Source 7
Title: Circuit — cep21/circuit
Publisher: cep21
URL: https://github.com/cep21/circuit
Published: 2017–present
Accessed: 2026-09-25
Source Tier: Tier 1 (Go library source)
Relevance: Hystrix-like Go impl; zero-goroutine, context-aware, configurable open/close logic, metrics/SLO tracking, expvar

## Source 8
Title: Bulkhead pattern — Azure Architecture Center
Publisher: Microsoft Learn
URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/bulkhead
Published: 2026-03-19 (ms.date)
Accessed: 2026-09-25
Source Tier: Tier 1
Relevance: Bulkhead as isolation (connection pools, partitions, cells) distinct from circuit breaker; combine with retry/CB/throttling

## Source 9
Title: Managing Load (SRE Workbook Ch. 11)
Publisher: Google SRE
URL: https://sre.google/workbook/managing-load/
Published: 2018 (SRE Workbook)
Accessed: 2026-09-25
Source Tier: Tier 1
Relevance: Load shedding vs breaker vs balancing; Dressy case study (shedding + balancing feedback loop), Pokémon GO thundering-herd retry sync, kill switches

## Source 10
Title: Retry pattern — Azure Architecture Center
Publisher: Microsoft Learn
URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/retry
Published: 2024-07-18 (ms.date)
Accessed: 2026-09-25
Source Tier: Tier 1
Relevance: Transient fault retry; cancel / immediate retry / retry-after-delay; exponential increment; idempotency; combine Retry + Circuit Breaker with sensitivity to CB exceptions

## Source 11
Title: Queue-Based Load Leveling pattern — Azure Architecture Center
Publisher: Microsoft Learn
URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/queue-based-load-leveling
Published: 2026-06-09 (ms.date)
Accessed: 2026-09-25
Source Tier: Tier 1
Relevance: Async decoupling via queue for non-critical flows; idempotency, dead-letter, ordering limits; enables degraded async processing when sync dependency down

## Source 12
Title: Timeouts, retries, and backoff with jitter — AWS Builder's Library
Publisher: Amazon Web Services
URL: https://builder.aws.com/content/3EumjoZascWd1oZiEgL8ORlv3qE/timeouts-retries-and-backoff-with-jitter/
Note: Original URL (aws.amazon.com/builders-library/...) redirects 301 to this AWS Builder Center URL.
Published: 2020 (Builder's Library)
Accessed: 2026-09-25
Source Tier: Tier 1
Relevance: Retry amplification under degraded dependency; jitter to prevent synchronized spikes; referenced by Source 6 — NOT VERIFIED via direct fetch (AWS Builder Center returned placeholder); claim supported by Source 6 + Source 10 instead

## Source 13
Title: Addressing Cascading Failures (SRE Book Ch. 22)
Publisher: Google SRE
URL: https://sre.google/sre-book/addressing-cascading-failures/
Published: 2016
Accessed: 2026-09-25
Source Tier: Tier 1
Relevance: Cascading failure mechanics; overload propagation; referenced via Source 5 chapter navigation — direct fetch not performed; listed as related primary
