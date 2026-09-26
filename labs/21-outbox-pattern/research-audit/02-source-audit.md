# Source Audit

## Source 1

Claimed Title: Pattern: Transactional outbox
Claimed Publisher: microservices.io
URL: https://microservices.io/patterns/data/transactional-outbox.html

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. URL is active and contains the canonical pattern definition by Chris Richardson.

Assessment:
PASS

---

## Source 2

Claimed Title: Pattern: Transaction log tailing
Claimed Publisher: microservices.io
URL: https://microservices.io/patterns/data/transaction-log-tailing.html

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Details transaction log mining/tailing (MySQL binlog, Postgres WAL, DynamoDB streams).

Assessment:
PASS

---

## Source 3

Claimed Title: Pattern: Polling publisher
Claimed Publisher: microservices.io
URL: https://microservices.io/patterns/data/polling-publisher.html

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Details polling-based outbox message relay alternative.

Assessment:
PASS

---

## Source 4

Claimed Title: Pattern: Saga
Claimed Publisher: microservices.io
URL: https://microservices.io/patterns/data/saga.html

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Contextual source explaining multi-service transactions and dual-write problem necessity.

Assessment:
PASS

---

## Source 5

Claimed Title: Reliable Microservices Data Exchange With The Outbox Pattern
Claimed Publisher: Debezium (Gunnar Morling)
URL: https://debezium.io/blog/2019/02/19/reliable-microservices-data-exchange-with-the-outbox-pattern/

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Provides concrete implementation details, table schema (`id`, `aggregatetype`, `aggregateid`, `type`, `payload`), and consumer idempotency logic via `MessageLog`.

Assessment:
PASS

---

## Source 6

Claimed Title: Outbox Event Router
Claimed Publisher: Debezium Documentation
URL: https://debezium.io/documentation/reference/stable/transformations/outbox-event-router.html

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Official reference documentation for Debezium Outbox Event Router SMT.

Assessment:
PASS
