# Source Audit

Target Lab: `labs/21-outbox-pattern`

---

## Source 1

Claimed Title: Outbox Event Router (Debezium Documentation)
Claimed Publisher: Debezium Community
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
- None. Documentation confirms expected table schema (`id`, `aggregatetype`, `aggregateid`, `type`, `payload`), SMT configuration options (`route.by.field`, `route.topic.replacement`), and behavior on UPDATE/DELETE.

Assessment:
PASS

---

## Source 2

Claimed Title: Reliable Microservices Data Exchange With the Outbox Pattern
Claimed Publisher: Debezium Blog (Gunnar Morling - Debezium project lead)
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
- None. Article explicitly explains the dual-write problem, why distributed transactions (XA/2PC) are not feasible with systems like Kafka, outbox table schema, transaction log tailing with CDC, and consumer-side duplicate exclusion via `MessageLog` table.

Assessment:
PASS

---

## Source 3

Claimed Title: Pattern: Transactional Outbox
Claimed Publisher: Microservices.io (Chris Richardson)
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
- None. Canonical architectural pattern definition. Verifies forces (no 2PC, commit atomicity, order preservation), solution structure (Sender, DB, Message Outbox, Message Relay), at-least-once delivery issues, and requirement for idempotent consumers.

Assessment:
PASS

---

## Source 4

Claimed Title: Pattern: Polling Publisher
Claimed Publisher: Microservices.io (Chris Richardson)
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
- None. Confirms polling as alternative relay mechanism for SQL databases, along with trade-offs (polling overhead and ordering complexities).

Assessment:
PASS

---

## Source 5

Claimed Title: Pattern: Saga
Claimed Publisher: Microservices.io (Chris Richardson)
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
- None. Provides context on why transactional messaging/outbox is essential when coordinating multi-service transactions without 2PC.

Assessment:
PASS

---

## Source 6

Claimed Title: Debezium Architecture
Claimed Publisher: Debezium Community
URL: https://debezium.io/documentation/reference/stable/architecture.html

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Authoritative overview of Debezium connectors reading DBMS transaction logs (PostgreSQL WAL, MySQL binlog) via Kafka Connect / Debezium Server / Embedded Engine.

Assessment:
PASS
