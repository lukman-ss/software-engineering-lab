# Claim Audit

Target Lab: `labs/21-outbox-pattern`

---

## Claim 1

Claim: Updating a database and publishing a message cannot be made atomic using traditional distributed transactions (2PC) or database transactions alone.

Location: `research/05-report.md` (Finding 1), `research/03-evidence.md` (Evidence 1, 2)

Evidence Provided: Database transactions only rollback relational database mutations. External systems (message queues, Redis, HTTP endpoints) do not participate in local DB rollback. 2PC is often unsupported or causes tight operational coupling.

Source: `https://microservices.io/patterns/data/transactional-outbox.html`, `https://debezium.io/blog/2019/02/19/reliable-microservices-data-exchange-with-the-outbox-pattern/`

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes: Accurately reflects distributed systems constraints and dual-write mechanics.

---

## Claim 2

Claim: The Outbox pattern guarantees that messages are sent if and only if the database transaction commits.

Location: `research/05-report.md` (Finding 2), `research/03-evidence.md` (Evidence 3)

Evidence Provided: Storing the message in an outbox table within the same DB transaction boundary ensures both succeed or both abort. The separate message relay asynchronously guarantees eventual publication.

Source: `https://microservices.io/patterns/data/transactional-outbox.html`, `https://debezium.io/blog/2019/02/19/reliable-microservices-data-exchange-with-the-outbox-pattern/`

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes: Well-documented canonical guarantee of the Transactional Outbox pattern.

---

## Claim 3

Claim: Two patterns exist for implementing the message relay: Transaction Log Tailing (CDC) and Polling Publisher.

Location: `research/05-report.md` (Finding 3), `research/03-evidence.md` (Evidence 5)

Evidence Provided: Relay can either poll the outbox table or tail database transaction logs (WAL/binlog). Trade-offs regarding latency, database load, and portability are cited.

Source: `https://microservices.io/patterns/data/transaction-log-tailing.html`, `https://microservices.io/patterns/data/polling-publisher.html`, `https://debezium.io/documentation/reference/stable/architecture.html`

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes: Accurately categorizes the relay mechanisms.

---

## Claim 4

Claim: Outbox pattern provides at-least-once delivery, requiring consumers to be idempotent to handle duplicate messages.

Location: `research/05-report.md` (Finding 4), `research/03-evidence.md` (Evidence 6)

Evidence Provided: Relay crashes after publishing but before recording completion or offset commit cause redelivery. Consumers must track message/event UUIDs.

Source: `https://microservices.io/patterns/data/transactional-outbox.html`, `https://debezium.io/blog/2019/02/19/reliable-microservices-data-exchange-with-the-outbox-pattern/`

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes: Critical nuance properly emphasized: outbox guarantees at-least-once, not magic exactly-once.

---

## Claim 5

Claim: Outbox table design typically uses columns `id`, `aggregateid`, `aggregatetype`, `type`, `payload`, supporting both thin and fat event payloads.

Location: `research/05-report.md` (Finding 5), `research/03-evidence.md` (Evidence 4, 7, 8)

Evidence Provided: Debezium outbox event router specification and real-world implementations standardise on these columns. Thin vs fat tradeoffs (DB size vs consumer RPC call) documented.

Source: `https://debezium.io/documentation/reference/stable/transformations/outbox-event-router.html`, `https://debezium.io/blog/2019/02/19/reliable-microservices-data-exchange-with-the-outbox-pattern/`

Source Actually Supports Claim:
YES

Classification:
IMPLEMENTATION-SPECIFIC

Severity:
LOW

Notes: Appropriately identified as common/reference design rather than universal SQL standard.

---

## Claim 6

Claim: Outbox systems require operational cleanup and monitoring of key metrics (unprocessed event count, oldest event age, publish failure rate).

Location: `research/05-report.md` (Finding 6), `research/03-evidence.md` (Evidence 9, 10)

Evidence Provided: Uncleaned tables grow indefinitely with high event volume. Oldest unprocessed event age is key latency/lag health indicator. Threshold numbers (2s normal vs 47m critical) are presented as illustrative SLA examples.

Source: Topic specification, operational best practices

Source Actually Supports Claim:
YES

Classification:
EXAMPLE

Severity:
LOW

Notes: The report explicitly qualifies numeric thresholds as illustrative SLO examples rather than universal operational constants.
