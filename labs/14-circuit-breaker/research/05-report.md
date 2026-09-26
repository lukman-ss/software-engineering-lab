# Research Report

## Research Question
How does the Circuit Breaker pattern halt cascading failures, preserve caller resources, and coordinate with Timeouts, Retries, Bulkheads, and Fallbacks?

## Executive Summary
Downstream latency or failure causes callers holding connections until timeouts expire to exhaust threads, sockets, and memory. Circuit breakers track error rates; once a threshold is breached the circuit trips to OPEN and fast-fails future calls in microseconds without network I/O. After a cooldown the breaker enters HALF_OPEN, permitting limited probe traffic to verify recovery before restoring CLOSED routing. Timeout bounds single attempts, jittered retry handles transient blips, bulkhead isolates resource pools, shedding drops inbound overload — breaker stops repeated outbound calls to a known-failing dependency.

## Findings

### Finding 1
Claim: Unmitigated latency causes caller thread/connection exhaustion and domino-style cascading failure.
Evidence: "If a service is busy, failure in one part of the system might lead to cascading failures. [...] this strategy can block concurrent requests to the same operation until the time-out period expires. These blocked requests might hold critical system resources, such as memory, threads, and database connections."
Sources:
- Microsoft Azure Architecture Center, Circuit Breaker Pattern (2025-02-05) — https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker
- Martin Fowler, Circuit Breaker (2014-03-06) — https://martinfowler.com/bliki/CircuitBreaker.html ("you can run out of critical resources leading to cascading failures across multiple systems")
- Netflix Hystrix Wiki (2017-07-03) — https://github.com/Netflix/Hystrix/wiki ("99.99^30 = 99.7% uptime ... 2+ hours downtime/month even if all dependencies have excellent uptime")
Confidence: HIGH

### Finding 2
Claim: Circuit breakers transition deterministically between CLOSED, OPEN, and HALF_OPEN with fail-fast in OPEN and limited probes in HALF_OPEN.
Evidence: "You can implement the proxy as a state machine that includes the following states... Closed... Open... Half-Open" / "Open: The request from the application fails immediately and an exception is returned" / "Half-Open: A limited number of requests from the application are allowed to pass through and invoke the operation. If these requests are successful... switches to the Closed state."
Sources:
- Microsoft Azure Architecture Center, Circuit Breaker Pattern (2025-02-05) — https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker
- Martin Fowler, Circuit Breaker (2014-03-06) — https://martinfowler.com/bliki/CircuitBreaker.html
- Netflix Hystrix Wiki, How it Works (2017-07-03) — https://github.com/Netflix/Hystrix/wiki/How-it-Works (volume threshold + error % → OPEN; sleep window → single HALF_OPEN probe)
Confidence: HIGH

### Finding 3
Claim: Retry differs from breaker (retry repeats expecting success; breaker prevents likely-failing calls); unbounded retry amplifies outages and requires jittered backoff plus retry budgets.
Evidence: "The Retry pattern enables an application to retry an operation with the expectation that it eventually succeeds. The Circuit Breaker pattern prevents an application from performing an operation that's likely to fail." / "An aggressive retry policy with minimal delay between attempts, and a large number of retries, could further degrade a busy service that's running close to or at capacity."
Sources:
- Microsoft Azure, Circuit Breaker pattern (2025-02-05) — https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker
- Microsoft Azure, Retry pattern (2024-07-18) — https://learn.microsoft.com/en-us/azure/architecture/patterns/retry
- AWS Architecture Blog, Exponential Backoff And Jitter (2015-03-04) — https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/ (Full/Equal/Decorrelated jitter; no-jitter backoff "clear loser")
- Google SRE, Handling Overload Ch.21 (2016) — https://sre.google/sre-book/handling-overload/ (per-request budget 3 attempts, per-client 10% retry ratio)
Confidence: HIGH

### Finding 4
Claim: Bulkhead isolates resource pools per dependency and is distinct from circuit breaker; load shedding drops inbound requests server-side while breaker stops outbound calls caller-side.
Evidence: "Isolate the elements of an application into pools so that if one element fails, the others continue to function." / "Hystrix employs the bulkhead pattern to isolate dependencies from each other and to limit concurrent access to any one of them."
Sources:
- Microsoft Azure, Bulkhead pattern (2026-03-19) — https://learn.microsoft.com/en-us/azure/architecture/patterns/bulkhead
- Netflix Hystrix Wiki, How it Works — Isolation (2017-07-03) — https://github.com/Netflix/Hystrix/wiki/How-it-Works
- Google SRE Workbook, Managing Load Ch.11 (2018) — https://sre.google/workbook/managing-load/ (shedding vs balancing; Dressy feedback-loop case; Pokémon GO thundering-herd retry sync)
Confidence: HIGH

### Finding 5
Claim: Production breaker observability centers on state, failure/timeout/rejection counters, error %, latency histograms, and state-change alerts; Go reference impls expose counts + transition callbacks.
Evidence: "Circuit breakers should provide clear observability into both failed and successful requests" / "reports successes, failures, rejections, and timeouts to the circuit breaker, which maintains a rolling set of counters".
Sources:
- Microsoft Azure, Circuit Breaker pattern (2025-02-05) — https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker
- Netflix Hystrix Wiki (2017-07-03) — https://github.com/Netflix/Hystrix/wiki/How-it-Works
- Sony gobreaker — https://github.com/sony/gobreaker (Counts struct + OnStateChange)
- cep21/circuit — https://github.com/cep21/circuit (MetricEventStream, expvar, SLO tracking)
Confidence: HIGH

## Areas of Agreement
All primary sources agree: breaker must fail fast when OPEN, isolate blast radius, protect downstream recovery window, probe via limited HALF_OPEN traffic, surface state changes for monitoring, and combine with timeout + bounded retry + bulkhead rather than replacing them.

## Areas of Disagreement
Threshold measurement: consecutive-failure count (Fowler basic impl; gobreaker default 5) vs rolling-window error % with volume gate (Hystrix; Azure time-bounded count). Production favors rolling window; lab uses consecutive count for determinism. See 04-contradictions.md. No disagreement on state semantics.

## Limitations
- In-memory breaker is per-process; N replicas probe independently after cooldown (thundering-herd risk across fleet) unless coordinated externally.
- AWS Builder's Library page fetch returned placeholder; its claims cross-checked via AWS Architecture Blog + Azure Retry instead (see 02-sources.md Source 12).
- No universal threshold/cooldown values verified; tuning is per-dependency SLO.
- SRE Ch.22 (Addressing Cascading Failures) listed as related primary via chapter navigation, not direct fetch.

## Conclusion
Circuit breakers do not fix the dependency; they prevent caller self-destruction and grant the dependency space to recover. Layered with timeout, jittered bounded retry, bulkhead isolation, and async queue decoupling for non-critical flows, the breaker limits blast radius so callers fail fast, preserve resources, degrade gracefully, and re-admit traffic only after verified probes.
