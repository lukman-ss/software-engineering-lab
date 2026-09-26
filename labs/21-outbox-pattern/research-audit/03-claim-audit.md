# Claim Audit

## Claim 1

Claim:
Updating a database and publishing a message cannot be made atomic using traditional distributed transactions (2PC) easily or database transactions alone; DB transactions cannot rollback external brokers.

Location:
`research/03-evidence.md` (Evidence 1 & 2), `research/05-report.md` (Finding 1)

Evidence Provided:
Cites microservices.io Transactional Outbox stating 2PC is not viable due to broker/db lack of support or unwanted coupling. Cites architectural reality that database rollback does not undo external broker writes.

Source:
https://microservices.io/patterns/data/transactional-outbox.html

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Accurately reflects fundamental distributed systems transaction boundary limitation.

---

## Claim 2

Claim:
The Outbox pattern guarantees that messages are published if and only if the database transaction commits.

Location:
`research/03-evidence.md` (Evidence 3), `research/05-report.md` (Finding 2)

Evidence Provided:
Cites microservices.io: service stores message in database as part of transaction that updates business entities; separate process relays messages.

Source:
https://microservices.io/patterns/data/transactional-outbox.html

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Atomicity is guaranteed at persistence level; actual message delivery to consumers is asynchronous/eventual.

---

## Claim 3

Claim:
There are two canonical patterns for implementing the message relay: Transaction Log Tailing and Polling Publisher.

Location:
`research/03-evidence.md` (Evidence 5), `research/05-report.md` (Finding 3)

Evidence Provided:
Cites microservices.io Polling Publisher and Transaction Log Tailing pattern catalog entries.

Source:
https://microservices.io/patterns/data/transaction-log-tailing.html, https://microservices.io/patterns/data/polling-publisher.html

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Consensus architectural taxonomy.

---

## Claim 4

Claim:
Outbox pattern provides at-least-once delivery; consumers must be idempotent to handle duplicates caused by relay retries or crashes.

Location:
`research/03-evidence.md` (Evidence 6), `research/05-report.md` (Finding 4)

Evidence Provided:
Cites microservices.io: "The Message relay might publish a message more than once... consumer must be idempotent". Cites Debezium blog duplicate detection mechanism via UUID tracking.

Source:
https://microservices.io/patterns/data/transactional-outbox.html, https://debezium.io/blog/2019/02/19/reliable-microservices-data-exchange-with-the-outbox-pattern/

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Properly captures at-least-once delivery reality and explicit rejection of magical exactly-once assumptions.

---

## Claim 5

Claim:
Canonical outbox table design uses `id` (UUID), `aggregatetype` (string), `aggregateid` (string), `type` (string), and `payload` (json/jsonb).

Location:
`research/03-evidence.md` (Evidence 4 & 7), `research/05-report.md` (Finding 5)

Evidence Provided:
Cites Debezium reference documentation and Debezium Outbox pattern blog post.

Source:
https://debezium.io/documentation/reference/stable/transformations/outbox-event-router.html

Source Actually Supports Claim:
YES

Classification:
IMPLEMENTATION-SPECIFIC

Severity:
LOW

Notes:
While originating from Debezium's SMT convention, this has become the de-facto standard schema in relational outbox architectures. Correctly categorized in research as implementation standard rather than universal SQL requirement.

---

## Claim 6

Claim:
Numeric thresholds such as oldest unprocessed event age around 2 seconds (normal) vs 47 minutes (abnormal).

Location:
`research/03-evidence.md` (Evidence 10), `research/05-report.md` (Finding 6)

Evidence Provided:
Labs specification operational guidelines.

Source:
Topic specification (`labs/21-outbox-pattern`)

Source Actually Supports Claim:
PARTIAL

Classification:
EXAMPLE

Severity:
LOW

Notes:
The research author explicitly qualified this in `03-evidence.md` and `05-report.md` with: "These values are illustrative SLO examples, not universal operational constants" and noted it as MEDIUM confidence. Handled appropriately without making false universal claims.
