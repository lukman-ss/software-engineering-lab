# Audit Plan: Research Audit for Saga Pattern (Lab 29)

## Target Lab
`labs/29-saga-pattern`

## Scope of Audit
PIPELINE OVERRIDE: Research audit only.
Files under audit:
- `research/01-plan.md`
- `research/02-sources.md`
- `research/03-evidence.md`
- `research/04-contradictions.md`
- `research/05-report.md`
- `research/06-open-questions.md`

## Claims To Verify
1. 2PC is infeasible/suboptimal for microservices with database-per-service pattern due to availability and latency tradeoffs.
2. Saga decomposes a distributed transaction into a sequence of local transactions.
3. Two primary coordination approaches exist: Choreography (event-driven) and Orchestration (central coordinator), each with distinct tradeoffs.
4. Compensating transactions are business-level semantic corrective actions rather than automatic database rollbacks.
5. Sagas categorize transactions into Compensable, Pivot (point of no return), and Retryable transactions.
6. Sagas lack native ACID Isolation ('I'), exposing systems to anomalies (dirty reads, lost updates, fuzzy reads).
7. Countermeasures exist to mitigate isolation anomalies: semantic lock, commutative updates, pessimistic view, reread values, version files, risk-based concurrency.
8. Local transactions and message publication must be atomically coordinated (e.g., Transactional Outbox / Event Sourcing), and steps/compensations require idempotency.
9. Failure of compensating transactions presents a hard boundary requiring operational reconciliation (DLQ, monitoring, manual intervention).

## Sources To Audit
1. Microsoft Azure Architecture Center — Saga Design Pattern (`https://learn.microsoft.com/en-us/azure/architecture/patterns/saga`)
2. Chris Richardson / Microservices.io — Pattern: Saga (`https://microservices.io/patterns/data/saga.html`)
3. Chris Richardson / Manning Publications — Microservices Patterns Chapter 4 (`https://livebook.manning.com/book/microservices-patterns/chapter-4/143`)
4. Hector Garcia-Molina & Kenneth Salem — Sagas Paper (ACM SIGMOD 1987) (`https://dl.acm.org/doi/10.1145/62224.62226`)

## Primary Risks
- Overreliance on high-level documentation without formal mathematical/database proofs.
- Lack of empirical performance benchmark figures (p95 latency, throughput comparisons vs 2PC) in free authoritative documentation.
- Nuance around failure handling when compensating transactions themselves fail.

## Audit Strategy
1. Independently fetch and verify all URLs cited in `02-sources.md` and `03-evidence.md`.
2. Map each technical claim in `05-report.md` and `03-evidence.md` against verified primary sources.
3. Check for internal contradictions between research documents.
4. Verify open questions and research gaps are honestly reported.
5. Generate verdict according to auditor severity model.
