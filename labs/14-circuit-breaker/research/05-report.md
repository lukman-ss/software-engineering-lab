# Research Report

## Research Question
How does the Circuit Breaker pattern halt cascading failures, preserve caller resources, and coordinate with Timeouts, Retries, Bulkheads, and Fallbacks?

## Executive Summary
Network and downstream dependency failures in microservices are inevitable. When a service experiences latency or total failure, callers holding connections until timeouts expire exhaust OS threads, socket file descriptors, and memory pools. Circuit breakers intercept calls, tracking error rates. Once a threshold is breached, the circuit trips to OPEN, fast-failing future calls in nanoseconds without reaching the network. After a cooldown period, the breaker enters HALF_OPEN, permitting canary probe traffic to safely determine health before restoring normal CLOSED routing.

## Findings

### Finding 1: Cascading Failure Dynamics
Claim: Unmitigated latency causes caller thread and connection exhaustion, dragging down upstream systems.
Evidence: Microsoft Azure Architecture Center notes blocked requests hold memory, threads, and DB connections, leading to domino-style failure across services.
Sources: Microsoft Azure Architecture Center, Martin Fowler.
Confidence: HIGH

### Finding 2: The Three-State Machine
Claim: Circuit Breakers transition deterministically between CLOSED, OPEN, and HALF_OPEN.
Evidence: When CLOSED, failures increment a counter; threshold breach transitions to OPEN. OPEN fast-fails. Cooldown timer expiration transitions to HALF_OPEN. Successful probe transitions to CLOSED; failed probe resets cooldown in OPEN.
Sources: Martin Fowler, Microsoft Azure Architecture Center.
Confidence: HIGH

### Finding 3: Interaction with Retries and Bulkheads
Claim: Retries without circuit breaking amplify downstream degradation ("retry storm"); Bulkheads isolate resource allocations so failures don't consume unrelated worker pools.
Evidence: AWS Builder's Library highlights retry amplification and the necessity of backoff and jitter. Circuit breakers act as a stopgap to silence retries when a system is in distress.
Sources: AWS Builder's Library (Source 12), Microsoft Azure Architecture Center (Source 3); see also AWS Architecture Blog on backoff and jitter (Source 6).
Confidence: HIGH

## Areas of Agreement
All primary sources agree that Circuit Breakers must fail fast when OPEN, isolate blast radius, protect downstream recovery, and provide canary probes during HALF_OPEN.

## Areas of Disagreement
Measurement strategy for failure thresholds: raw consecutive failures (simpler, deterministic) vs rolling-window error rate percentages (more robust in production). Production systems (Netflix Hystrix - Source 2, Sony gobreaker - Source 7) use rolling windows; this lab uses consecutive failures for determinism (see research/04-contradictions.md for framing).

## Limitations
In-memory circuit breakers track state on a per-instance basis without shared coordination across a scaled-out fleet. In a deployment of multiple container instances (e.g., 100 replicas), each instance independently transitions through states, meaning 100 simultaneous probes could strike a recovering downstream dependency at the end of the cooldown period unless synchronized via a centralized store or coordinated rate limits.

**Source**: cep21/circuit (Source 8) — implementation context; Microsoft Azure Architecture Center (Source 3) — notes concurrency considerations for multi-instance deployments

## Conclusion
Circuit breakers do not fix the dependency; they prevent caller self-destruction and grant the dependency space to recover.
