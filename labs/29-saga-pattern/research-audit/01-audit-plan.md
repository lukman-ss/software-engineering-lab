# 01 - Audit Plan: Saga Pattern Research

## Target Lab
`labs/29-saga-pattern`

## Scope
Research files audit only (Pipeline Override: research files under `labs/29-saga-pattern/research/`).

## Files Reviewed
- `labs/29-saga-pattern/research/01-plan.md`
- `labs/29-saga-pattern/research/02-sources.md`
- `labs/29-saga-pattern/research/03-evidence.md`
- `labs/29-saga-pattern/research/04-contradictions.md`
- `labs/29-saga-pattern/research/05-report.md`
- `labs/29-saga-pattern/research/06-open-questions.md`

## Claims To Verify
1. Saga pattern origin (Garcia-Molina & Salem 1987 ACM SIGMOD).
2. Fundamental definition: sequence of local transactions + compensating transactions.
3. Coordination strategies: Choreography (decentralized/event-driven) vs Orchestration (centralized coordinator).
4. Compensating transactions as semantic undo operations (not automatic ACID rollback).
5. Taxonomy of steps: compensable, pivot (point of no return), retryable transactions.
6. 2PC vs Saga tradeoffs: blocking/isolation vs availability/resilience/eventual consistency.
7. Dual-write problem and Transactional Outbox pattern mitigation (CDC/Debezium).
8. Idempotency requirement for message consumers & participants (e.g. `PROCESSED_MESSAGES` table, idempotency keys).
9. Lack of ACID Isolation: data anomalies (dirty reads, lost updates, fuzzy reads) and 6 countermeasures.
10. Failure modes: compensation itself failing, lack of standard recovery protocol without manual/DLQ intervention.
11. Real-world orchestration frameworks: AWS Step Functions, Temporal, Eventuate Tram.

## Primary Risks
- Relying on single-source formalizations (e.g. Microsoft Azure Architecture Center's 6 specific countermeasures or pivot taxonomy) as universal standards.
- Incomplete verbatim validation of the 1987 paper due to PDF compression/formatting.
- Over-generalizing workflow-level "exactly-once" guarantees when participants operate under at-least-once delivery semantics.
- Conflating theoretical recovery protocols with production operational realities (human intervention / dead-letter queues).

## Audit Strategy
1. Perform comprehensive source inspection across Tier 1 and Tier 2 citations in `02-sources.md`.
2. Map every claim in `03-evidence.md` and `05-report.md` against cited sources, checking for unsupported generalizations or scope errors.
3. Verify handling of contradictions in `04-contradictions.md` and gap recording in `06-open-questions.md`.
4. Render an evidence-based verdict in accordance with the Auditor Agent Quality Gates.
