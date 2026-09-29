# Source Register: Saga Pattern Research

## Source 1: Garcia-Molina & Salem — "Sagas" (ACM SIGMOD 1987)

- **Title:** Sagas
- **Publisher:** ACM (Association for Computing Machinery)
- **URL:** https://www.cs.cornell.edu/andru/cs711/2002fa/reading/sagas.pdf
- **DOI:** 10.1145/62224.62226
- **Published:** 1987
- **Accessed:** 2026-09-29
- **Source Tier:** Tier 1 (Academic Paper / Original Source)
- **Relevance:** Original definition of the saga concept; introduced compensating transactions for long-lived database transactions.

## Source 2: Microsoft Azure Architecture Center — Saga Design Pattern

- **Title:** Saga Design Pattern
- **Publisher:** Microsoft
- **URL:** https://learn.microsoft.com/en-us/azure/architecture/patterns/saga
- **Published:** 2025-02-25 (last updated)
- **Accessed:** 2026-09-29
- **Source Tier:** Tier 1 (Official Company Documentation)
- **Relevance:** Comprehensive description of saga pattern including choreography vs orchestration, compensable/pivot/retryable transactions, data anomaly countermeasures.

## Source 3: Chris Richardson — Microservices.io: Saga Pattern

- **Title:** Pattern: Saga
- **Publisher:** Chris Richardson / Microservices.io
- **URL:** https://microservices.io/patterns/data/saga.html
- **Published:** 2024-2026 (ongoing updates)
- **Accessed:** 2026-09-29
- **Source Tier:** Tier 2 (Reputable Technical Publication)
- **Relevance:** Authoritative reference on saga pattern in microservices context; defines choreography/orchestration, related patterns (outbox, event sourcing, idempotent consumer).

## Source 4: AWS Step Functions — Saga Pattern

- **Title:** Saga pattern
- **Publisher:** Amazon Web Services
- **URL:** https://docs.aws.amazon.com/prescriptive-guidance/latest/modernization-data-persistence/saga-pattern.html
- **Published:** 2023+ (AWS Prescriptive Guidance)
- **Accessed:** 2026-09-29
- **Source Tier:** Tier 1 (Official Company Documentation)
- **Relevance:** Describes orchestration saga with AWS Step Functions; defines Catch/Retry for compensation; identifies complexity tradeoffs.

## Source 5: Temporal Blog — Saga Design Pattern Explained

- **Title:** Saga design pattern explained: Benefits, use cases, and implementation
- **Publisher:** Temporal Technologies
- **URL:** https://temporal.io/blog/saga-pattern-made-easy
- **Published:** 2023-05-24
- **Accessed:** 2026-09-29
- **Source Tier:** Tier 2 (Reputable Technical Publication)
- **Relevance:** Practical implementation of sagas using Temporal workflows; demonstrates compensation pattern, idempotency keys, code examples in Java and Go.

## Source 6: Temporal Blog — Choreography vs Orchestration

- **Title:** To choreograph or orchestrate your saga, that is the question
- **Publisher:** Temporal Technologies
- **URL:** https://temporal.io/blog/to-choreograph-or-orchestrate-your-saga-that-is-the-question
- **Published:** 2023-07-13
- **Accessed:** 2026-09-29
- **Source Tier:** Tier 2 (Reputable Technical Publication)
- **Relevance:** Direct comparison of choreography and orchestration approaches; discusses tradeoffs including SPOF, debuggability, and when to use each.

## Source 7: Debezium Blog — Transactional Outbox Pattern

- **Title:** Reliable Microservices Data Exchange With the Outbox Pattern
- **Publisher:** Gunnar Morling / Debezium
- **URL:** https://debezium.io/blog/2019/02/19/reliable-microservices-data-exchange-with-the-outbox-pattern/
- **Published:** 2019-02-19
- **Accessed:** 2026-09-29
- **Source Tier:** Tier 1 (Original Authoritative Source — CDC/Outbox domain)
- **Relevance:** Explains the dual-write problem; provides outbox table schema; demonstrates CDC-based event propagation with Debezium; addresses idempotent consumption via event UUIDs.

## Source 8: Chris Richardson — Idempotent Consumer Pattern

- **Title:** Pattern: Idempotent Consumer
- **Publisher:** Chris Richardson / Microservices.io
- **URL:** https://microservices.io/patterns/data/idempotent-consumer.html
- **Published:** 2020-2026 (ongoing)
- **Accessed:** 2026-09-29
- **Source Tier:** Tier 2 (Reputable Technical Publication)
- **Relevance:** Defines idempotency requirement for message consumers; describes PROCESSED_MESSAGES table approach for duplicate detection.

## Source 9: AWS Step Functions Documentation (Welcome / Concepts)

- **Title:** What is Step Functions?
- **Publisher:** Amazon Web Services
- **URL:** https://docs.aws.amazon.com/step-functions/latest/dg/welcome.html
- **Published:** 2026+ (ongoing)
- **Accessed:** 2026-09-29
- **Source Tier:** Tier 1 (Official Company Documentation)
- **Relevance:** Describes Standard vs Express workflows; exactly-once vs at-least-once execution semantics; state machine model underlying orchestration sagas.
