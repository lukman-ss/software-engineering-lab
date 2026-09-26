# Research Plan

## Research Topic
Transactional Outbox Pattern — Database Sudah Commit, Tapi Event/Queue Gagal Dikirim

## Objective
Produce structured research covering:
1. Dual-write problem and why distributed transactions fail
2. Transactional Outbox pattern core concepts
3. Message relay implementations (polling, transaction log tailing)
4. Idempotent consumer requirements
5. Event payload design and outbox table structure
6. Monitoring and operational concerns
7. Anti-patterns and common mistakes

## Research Questions
1. What is the dual-write problem and why does it occur?
2. Why can't database transactions span database and message broker?
3. How does the Outbox pattern guarantee atomic updates?
4. What are the two message relay implementations and their tradeoffs?
5. Why must consumers be idempotent? What causes duplicate events?
6. What outbox table structure is recommended? What fields are essential?
7. What monitoring metrics indicate Outbox system health?
8. What are common anti-patterns and how to avoid them?

## Search Strategy
- Primary: microservices.io Transactional Outbox (original pattern definition)
- Primary: Debezium Outbox Event Router documentation (reference implementation)
- Primary: Debezium blog on Outbox pattern (implementation details)
- Primary: microservices.io Saga and Transactional messaging patterns
- Secondary: Microsoft Azure Architecture patterns

## Expected Primary Sources
| Source | URL | Tier |
|--------|-----|------|
| microservices.io Transactional Outbox | https://microservices.io/patterns/data/transactional-outbox.html | 1 |
| Debezium Outbox Event Router | https://debezium.io/documentation/reference/stable/transformations/outbox-event-router.html | 1 |
| Debezium Blog: Reliable Microservices Data Exchange | https://debezium.io/blog/2019/02/19/reliable-microservices-data-exchange-with-the-outbox-pattern/ | 1 |
| microservices.io Polling Publisher | https://microservices.io/patterns/data/polling-publisher.html | 1 |
| microservices.io Transaction Log Tailing | https://microservices.io/patterns/data/transaction-log-tailing.html | 1 |

## Risks / Unknowns
- No Microsoft Azure Architecture Center page found for Transactional Outbox
- Some external blog URLs may be inaccessible or redirected
- Must distinguish between generic pattern and Debezium-specific implementation
- No universal numeric recommendations — all config values are environment-specific
