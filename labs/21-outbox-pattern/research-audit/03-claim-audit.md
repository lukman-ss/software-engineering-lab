# Claim Audit: Transactional Outbox Pattern

Target Lab: `labs/21-outbox-pattern`
Research Set Under Audit: `research/2026-09-26-outbox-pattern/`
Date: 2026-09-26

---

## Claim 1: Dual-Write Problem & 2PC Infeasibility

Claim:
Service updating a database and publishing to a message broker cannot do both atomically without risk of half-success. Distributed 2PC transactions spanning DB + Kafka/RabbitMQ are unsupported or highly impractical.

Location:
`05-report.md` (Finding 1), `03-evidence.md` (Evidence 1 & 2)

Evidence Provided:
Microservices.io transactional-outbox pattern statement; Debezium 2019 blog explaining why Kafka cannot enlist in XA transactions.

Source:
- https://microservices.io/patterns/data/transactional-outbox.html
- https://debezium.io/blog/2019/02/19/reliable-microservices-data-exchange-with-the-outbox-pattern/

Source Actually Supports Claim: YES

Classification: FACT

Severity: HIGH (Core architectural invariant)

Notes: Fully supported.

---

## Claim 2: Database Transaction Boundaries Cannot Encompass External Side Effects

Claim:
Wrapping broker pushes inside a database transaction block (`DB::transaction(...)`) still causes inconsistencies if the DB rolls back or connection drops after broker publish.

Location:
`05-report.md` (Finding 2), `03-evidence.md` (Evidence 2)

Evidence Provided:
Microservices.io and standard ACID transactional scope definitions.

Source:
- https://microservices.io/patterns/data/transactional-outbox.html

Source Actually Supports Claim: YES

Classification: FACT

Severity: HIGH

Notes: Fully supported.

---

## Claim 3: Outbox Core Atomic Write & Async Relay

Claim:
Writing business aggregate change and an outbox event record within a single database transaction guarantees all-or-nothing persistence. An asynchronous relay process subsequently reads the outbox and publishes to the broker.

Location:
`05-report.md` (Finding 3), `03-evidence.md` (Evidence 3, 4, 5)

Evidence Provided:
Microservices.io pattern definition, component taxonomy (Sender, Database, Message Outbox, Message Relay).

Source:
- https://microservices.io/patterns/data/transactional-outbox.html

Source Actually Supports Claim: YES

Classification: FACT

Severity: HIGH

Notes: Fully supported.

---

## Claim 4: Canonical Outbox Schema vs Universal Generalization

Claim:
Debezium's canonical outbox schema includes `id`, `aggregatetype`, `aggregateid`, `type`, and `payload`. This represents Debezium's implementation conventions (and default SMT mappings), not a mandatory universal standard across all possible outbox implementations.

Location:
`05-report.md` (Finding 4), `03-evidence.md` (Evidence 6)

Evidence Provided:
Debezium Outbox Event Router documentation and Debezium 2019 blog.

Source:
- https://debezium.io/documentation/reference/stable/transformations/outbox-event-router.html
- https://debezium.io/blog/2019/02/19/reliable-microservices-data-exchange-with-the-outbox-pattern/

Source Actually Supports Claim: YES

Classification: IMPLEMENTATION-SPECIFIC (Properly qualified in report)

Severity: MEDIUM

Notes: The report explicitly disclaims that this is a universal standard, correctly identifying it as Debezium's default convention.

---

## Claim 5: Relay Alternatives (Polling Publisher vs Log Tailing CDC)

Claim:
Message relay can be implemented via Polling Publisher (portable across SQL DBs, challenges in polling frequency and ordering) or Transaction Log Tailing / CDC (high performance, low latency, database-engine specific, requires deduplication handling).

Location:
`05-report.md` (Finding 5), `03-evidence.md` (Evidence 8, 9)

Evidence Provided:
Microservices.io polling-publisher and transaction-log-tailing definitions, Debezium Postgres WAL implementation.

Source:
- https://microservices.io/patterns/data/polling-publisher.html
- https://microservices.io/patterns/data/transaction-log-tailing.html
- https://debezium.io/blog/2019/02/19/reliable-microservices-data-exchange-with-the-outbox-pattern/

Source Actually Supports Claim: YES

Classification: FACT / ARCHITECTURAL TRADEOFF

Severity: HIGH

Notes: Accurately reflects both patterns and trade-offs.

---

## Claim 6: End-to-End Delivery Semantics (At-Least-Once & Idempotent Consumer)

Claim:
Transactional Outbox provides at-least-once delivery semantics due to possible relay crashes between publish and status update/deletion. Therefore, consumer idempotency (e.g. via `consumed_messages` / message log tracking) is mandatory.

Location:
`05-report.md` (Finding 6), `03-evidence.md` (Evidence 7, 10), `02-sources.md` (Source 9)

Evidence Provided:
Microservices.io results section, Debezium `ConsumedMessage` implementation, EIP Idempotent Receiver pattern.

Source:
- https://microservices.io/patterns/data/transactional-outbox.html
- https://debezium.io/blog/2019/02/19/reliable-microservices-data-exchange-with-the-outbox-pattern/
- https://www.enterpriseintegrationpatterns.com/patterns/messaging/IdempotentReceiver.html

Source Actually Supports Claim: YES

Classification: FACT

Severity: CRITICAL

Notes: Accurately rejects magical "end-to-end exactly-once" without consumer deduplication.

---

## Claim 7: Kafka Transactions Scope Boundary

Claim:
Kafka Exactly-Once Semantics (EOS) coordinates multi-partition / multi-topic stream processing (Kafka-to-Kafka), but does not solve the upstream DB-to-Kafka dual-write problem.

Location:
`05-report.md` (Finding 7), `04-contradictions.md` (Nuance 1)

Evidence Provided:
Confluent blog on Kafka Transactions.

Source:
- https://www.confluent.io/blog/transactions-apache-kafka/

Source Actually Supports Claim: YES

Classification: FACT

Severity: HIGH

Notes: Essential boundary distinction properly articulated.

---

## Claim 8: Operational Guidance (Payload Sizing, Table Retention, DLQ, Monitoring)

Claim:
- Keep payload minimal (reference keys + essential state).
- Clean up processed records via retention purge or CDC ephemeral row technique.
- Employ DLQ for unprocessable messages.
- Monitor oldest unprocessed event age and queue lag. Note: Numeric thresholds (e.g., 2s vs 47m) are educational examples from lab specs, NOT industry benchmarks.

Location:
`05-report.md` (Findings 8, 9), `03-evidence.md` (Evidence 11, 12, 13, 16)

Evidence Provided:
Debezium blog (DLQ, ephemeral row, event evolution), Lab spec (practice guidance).

Source:
- https://debezium.io/blog/2019/02/19/reliable-microservices-data-exchange-with-the-outbox-pattern/
- Lab specification

Source Actually Supports Claim: YES (with explicit caveats on numeric examples)

Classification: INTERPRETATION / PRACTICE GUIDANCE

Severity: MEDIUM

Notes: Report explicitly identifies arbitrary numeric threshold examples as lab-spec derived, adhering strictly to anti-hallucination rules.
