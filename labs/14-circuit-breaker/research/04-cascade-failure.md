# 04 Cascade Failure

## Definition
A fault in one component triggers a chain reaction that progressively disables additional components until a large portion of the system is non-functional.

**Source**: Google SRE Book, "Addressing Cascading Failures" (Ch. 22, 2016) — https://sre.google/sre-book/addressing-cascading-failures/ — Tier 1 — "When everything is healthy the request flow can look like this... When one of many backend systems becomes latent it can block the entire user request"; "If the host application is not isolated from these external failures, it risks being taken down with them." (Martin Fowler cross-ref).

## Mechanism Without Breaker
```
Downstream slow/down
↓ caller requests block until timeout
↓ caller threads / pool slots / memory held
↓ caller rejects new requests / queues
↓ upchain callers also block → domino
```

### Resource Exhaustion Vector (Verified)
Azure Circuit Breaker pattern: "If a service is busy... blocked requests might hold critical system resources, such as memory, threads, and database connections. This problem can exhaust resources... causing even more cascading failures across the system."
URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker — 2025-02-05 — Confidence HIGH.

### SRE Quantitative Model (Verified)
Hystrix Wiki: "for an application that depends on 30 services where each service has 99.99% uptime, 99.99^30 = 99.7% uptime. 0.3% of 1 billion requests = 3,000,000 failures. 2+ hours downtime/month even if all dependencies have excellent uptime."
URL: https://github.com/Netflix/Hystrix/wiki — 2017-07-03 — Confidence HIGH — Corroborated by Azure Cascade Failure description.

## How Breaker Interrupts Cascade (Verified)
1. Threshold reached → OPEN
2. Fail fast → no threads/connections consumed by failing call
3. Downstream gets breathing room; caller preserves capacity
4. Cooldown → HALF_OPEN probing — if healthy, CLOSED, normal flow resumes

**Source**: Azure Circuit Breaker — "fail fast instead of queueing"; Hystrix Wiki — "stop cascading failures in a complex distributed system".

## Quantitative (Illustrative, NOT production data)
- Without CB: each request holds worker for full timeout; 100 concurrent blocked = 100 workers lost
- With CB (after N failures): subsequent calls return in microseconds; downstream call count plateaus

Per lab note: durations/counts from actual execution, not pre-baked figures.

## Pattern Relationship Table

| Pattern | Prevents | Mechanism | Source |
|---|---|---|---|
| Circuit Breaker | Downstream failure propagation | fail fast after threshold | Azure Circuit Breaker (2025-02-05) |
| Timeout | individual request hanging | bound max wait | Azure Circuit Breaker; Google SRE Workbook Ch.11 (2018) |
| Bulkhead | cross-dependency contamination | isolate resource pools | Azure Bulkhead (2026-03-19) |
| Load Shedding | overload from excess traffic | drop low-criticality requests | Google SRE Managing Load (2018); Dressy case study |
| Retry | transient failures | re-attempt, bounded + jittered | Azure Retry pattern (2024-07-18); AWS Backoff+jitter (2015) |
