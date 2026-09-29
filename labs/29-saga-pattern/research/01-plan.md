# Research Plan: Saga Pattern — Mengelola Transaksi Terdistribusi Tanpa 2PC

## Research Topic

Saga Pattern — distributed transaction management without Two-Phase Commit (2PC), including Choreography Saga, Orchestration Saga, Compensating Transactions, Dual-Write Problem, Event-Driven Architecture, Microservices consistency and resilience.

## Objective

Investigate and verify the technical claims made in the lab topic specification:

1. Saga Pattern originated from Garcia-Molina & Salem (1987) as a database recovery mechanism, later adapted for distributed microservices.
2. Sagas decompose distributed transactions into a sequence of local transactions with compensating transactions for rollback.
3. Two coordination approaches exist: Choreography (event-driven, decentralized) and Orchestration (centralized coordinator).
4. Compensating transactions must be semantically equivalent "undo" operations, not automatic database rollback.
5. The Dual-Write Problem (atomic update of database + message broker) is a prerequisite challenge that must be solved for reliable sagas.
6. Sagas sacrifice isolation (ACID "I") in favor of availability; countermeasures exist for data anomalies.
7. The Transactional Outbox Pattern resolves the dual-write problem by writing events to an outbox table in the same DB transaction and using CDC for propagation.
8. Idempotency is mandatory for saga participants; duplicate detection via processed-message IDs is the standard mechanism.
9. 2PC suffers from blocking, coordinator SPOF, and poor fit for microservices; sagas trade atomicity for resilience.
10. Real-world implementations exist in AWS Step Functions, Temporal, Camunda, Eventuate Tram, and Spring framework ecosystem.

## Research Questions

### RQ1: Origin and Formal Definition of Sagas
- Who introduced the saga concept and in what context (1987 ACM SIGMOD paper)?
- What was the original problem domain (long-lived database transactions)?
- How has the definition evolved for modern distributed systems?

### RQ2: Choreography vs Orchestration Tradeoffs
- What are the structural differences between choreography and orchestration sagas?
- Under what conditions is each approach preferable?
- What are the specific failure modes unique to each approach?

### RQ3: Compensating Transactions — Properties and Design Constraints
- What properties must a compensating transaction satisfy (semantic undo, idempotency, etc.)?
- What happens when a compensating transaction itself fails?
- How do real frameworks handle compensation failure scenarios?

### RQ4: Sagas vs 2PC — Technical Comparison
- What are the availability, latency, and coupling tradeoffs between 2PC and sagas?
- What consistency guarantees does each provide (strong vs eventual)?
- When is 2PC still the correct choice despite its drawbacks?

### RQ5: Dual-Write Problem and Reliable Event Publication
- What is the dual-write problem and why does it break consistency?
- How does the Transactional Outbox Pattern (CDC-based) solve it?
- What alternatives exist (event sourcing, polling publisher, etc.)?

### RQ6: Data Anomalies and Isolation Countermeasures
- What specific data anomalies can occur in sagas (dirty reads, lost updates, fuzzy reads)?
- What countermeasures exist (semantic lock, commutative updates, pessimistic view, reread values, version files)?
- How do these countermeasures trade off against system complexity?

### RQ7: Idempotency and Duplicate Detection
- Why is idempotency mandatory for saga participants?
- What mechanisms achieve duplicate detection (PROCESSED_MESSAGES table, event UUID, idempotency keys)?
- How do at-least-once delivery semantics interact with saga compensation logic?

### RQ8: Real-World Saga Frameworks
- How do AWS Step Functions implement orchestration sagas with Catch/Retry?
- How does Temporal's deterministic replay eliminate the orchestrator SPOF problem?
- What patterns do Eventuate Tram Sagas provide for local message publishing?

## Search Strategy

- Primary: Garcia-Molina & Salem 1987 paper (Cornell/ACM), Microsoft Azure Architecture Center, AWS Step Functions documentation, Temporal documentation.
- Secondary: Chris Richardson's microservices.io, Debezium/Outbox pattern blog, industry engineering publications.
- Keywords: `saga pattern`, `compensating transaction`, `choreography vs orchestration`, `dual-write problem`, `transactional outbox`, `2PC vs saga`, `idempotent consumer`, `eventual consistency microservices`.

## Expected Primary Sources

1. Garcia-Molina, H. & Salem, K. (1987). "Sagas." ACM SIGMOD Record, 16(3). DOI: 10.1145/62224.62226
2. Microsoft Azure Architecture Center — Saga Design Pattern (learn.microsoft.com/en-us/azure/architecture/patterns/saga)
3. AWS Step Functions — Saga Pattern (docs.aws.amazon.com/prescriptive-guidance/latest/modernization-data-persistence/saga-pattern.html)
4. Chris Richardson — Microservices Patterns: Saga (microservices.io/patterns/data/saga.html)
5. Temporal Blog — Saga Design Pattern Explained (temporal.io/blog/saga-pattern-made-easy)
6. Temporal Blog — Choreography vs Orchestration (temporal.io/blog/to-choreograph-or-orchestrate-your-saga-that-is-the-question)
7. Debezium Blog — Transactional Outbox Pattern (debezium.io/blog/2019/02/19/reliable-microservices-data-exchange-with-the-outbox-pattern/)
8. Chris Richardson — Idempotent Consumer Pattern (microservices.io/patterns/data/idempotent-consumer.html)

## Risks / Unknowns

- The original 1987 Garcia-Molina paper is a scanned PDF; exact page-level quotes may require careful extraction.
- Framework-specific implementation details (Camunda, Axon, Seata) not deeply explored; focus is on patterns not tooling.
- Edge cases around compensation failure (what happens if the compensation also fails) are theoretically described but lack comprehensive empirical evidence across all frameworks.
- The topic of exactly-once semantics vs at-least-once + idempotency is nuanced and may require further research for production-grade systems.
