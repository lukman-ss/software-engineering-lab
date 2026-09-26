# Claim Audit

## Claim 1: The Dual-Write Problem and Non-Atomicity Across Heterogeneous Systems

Claim:
Updating a database and publishing a message cannot be made atomic using traditional distributed transactions (2PC) or database transactions alone.

Location:
`research/05-report.md`: Section "Finding 1: The Dual-Write Problem and Failed Transactions"
`research/03-evidence.md`: Evidence 1 & Evidence 2

Evidence Provided:
Traditional database transactions only control database operations, not external message brokers or HTTP endpoints. 2PC is often unsupported or impractical across brokers and datastores.

Source:
Microservices.io (`https://microservices.io/patterns/data/transactional-outbox.html`) & Debezium Blog (`https://debezium.io/blog/2019/02/19/reliable-microservices-data-exchange-with-the-outbox-pattern/`)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Both Richardson and Morling explicitly note that 2PC cannot span most brokers (like Kafka/RabbitMQ) and DB transactions do not roll back broker sends.

---

## Claim 2: Outbox Pattern Atomicity Guarantee

Claim:
The Outbox pattern guarantees that messages are sent if and only if the database transaction commits.

Location:
`research/05-report.md`: Section "Finding 2: Outbox Pattern Guarantees Atomicity"
`research/03-evidence.md`: Evidence 3

Evidence Provided:
Both the business entity update and message record are inserted into the same database transaction. A separate relay process reads and delivers the message asynchronously.

Source:
Microservices.io (`https://microservices.io/patterns/data/transactional-outbox.html`)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Accurately reflects the core mechanism and theoretical guarantees of the Transactional Outbox pattern.

---

## Claim 3: Message Relay Implementation Patterns

Claim:
Two patterns exist for implementing the message relay: Transaction Log Tailing and Polling Publisher.

Location:
`research/05-report.md`: Section "Finding 3: Message Relay Implementation Alternatives"
`research/03-evidence.md`: Evidence 5

Evidence Provided:
Log tailing captures outbox table writes via database transaction log (binlog/WAL), while Polling queries the outbox table periodically.

Source:
Microservices.io (`https://microservices.io/patterns/data/transaction-log-tailing.html`, `https://microservices.io/patterns/data/polling-publisher.html`) & Debezium Blog

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Canonical classification used across distributed systems literature and pattern catalogs.

---

## Claim 4: At-Least-Once Delivery and Consumer Idempotency

Claim:
Outbox pattern provides at-least-once delivery, requiring consumers to be idempotent to handle duplicate messages.

Location:
`research/05-report.md`: Section "Finding 4: Idempotent Consumer Requirement"
`research/03-evidence.md`: Evidence 6

Evidence Provided:
Relay crashes or network retries after publishing but before marking processed result in re-delivery. Consumers track message UUIDs/eventIds to ignore duplicates.

Source:
Microservices.io & Debezium Blog

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Both sources explicitly stress that at-least-once is the real-world delivery semantic and consumer deduplication is mandatory.

---

## Claim 5: Outbox Table Structure and Payload Design (Thin vs Fat Events)

Claim:
Canonical outbox table design uses `id`, `aggregatetype`, `aggregateid`, `type`, `payload`. Thin events vs fat events represent trade-offs between RPC fetch overhead and storage/network bloat.

Location:
`research/05-report.md`: Section "Finding 5: Outbox Table Structure and Payload Design"
`research/03-evidence.md`: Evidence 4 & Evidence 7 & Evidence 8

Evidence Provided:
Debezium Outbox SMT specification expects `id`, `aggregatetype`, `aggregateid`, `type`, `payload`. The Debezium blog demonstrates full order payload (fat events), while the lab spec highlights hazards of storing 500-field objects (thin events alternative).

Source:
Debezium Documentation (Outbox Event Router) & Debezium Blog

Source Actually Supports Claim:
YES

Classification:
INTERPRETATION

Severity:
LOW

Notes:
The report accurately frames thin vs fat events as design choices rather than universal rules.

---

## Claim 6: Cleanup Strategy and Monitoring Metrics

Claim:
Outbox systems require operational cleanup (archival/deletion) to prevent unbounded table growth. Metrics such as unprocessed count and oldest event age are critical; numeric thresholds (e.g. 2s normal vs 47m abnormal) are illustrative SLO examples, not universal operational constants.

Location:
`research/05-report.md`: Section "Finding 6: Operational Requirements: Cleanup and Monitoring"
`research/03-evidence.md`: Evidence 9 & Evidence 10

Evidence Provided:
Debezium CDC pattern deletes row immediately or requires compaction/deletion; polling publishers require periodic cleanup. Monitoring thresholds depend on service SLAs.

Source:
Topic Specification & Debezium Blog

Source Actually Supports Claim:
YES

Classification:
EXAMPLE

Severity:
LOW

Notes:
The report explicitly qualifies the numbers as illustrative SLA examples and avoids presenting them as universal constants.
