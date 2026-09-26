# Evidence

## Evidence 1: The Dual-Write Problem

Claim: When a service updates a database and then publishes a message, or vice versa, there is no guarantee both operations succeed together.

Evidence: Without distributed transactions (2PC), there is no atomic guarantee across database and message broker. Either the database update succeeds and message fails, or message succeeds and database fails.

Source: microservices.io Transactional Outbox
URL: https://microservices.io/patterns/data/transactional-outbox.html
Confidence: HIGH — canonical pattern description.

## Evidence 2: Why Traditional DB Transaction Does Not Work

Claim: Database transaction cannot control external systems like Redis, RabbitMQ, Kafka, or HTTP endpoints.

Evidence: `DB::transaction()` only controls database operations. It cannot perform `ROLLBACK Redis`, `ROLLBACK RabbitMQ`, or `ROLLBACK HTTP Request`. If queue publish fails after DB commit, the event is lost.

Source: Topic specification (lab description)
URL: labs/21-outbox-pattern
Confidence: HIGH — operational reality of distributed systems.

## Evidence 3: Outbox Pattern Solution

Claim: Store the message in the database as part of the same transaction that updates business entities. A separate process then sends the messages to the message broker.

Evidence: The participants are:
- Sender: service that sends the message
- Database: stores business entities and message outbox
- Message outbox: table storing messages to be sent
- Message relay: sends messages from outbox to message broker

Because both operations are in the same database transaction, they are atomic: both succeed or both fail.

Source: microservices.io Transactional Outbox
URL: https://microservices.io/patterns/data/transactional-outbox.html
Confidence: HIGH — pattern definition.

## Evidence 4: Outbox Table Structure

Claim: Outbox table typically has columns: id (uuid), aggregateid, aggregatetype, type, payload (json/jsonb).

Evidence: Debezium example outbox table structure:
- id: unique id of each message, can be used by consumers to detect duplicate events
- aggregatetype: the type of the aggregate root (e.g., "Order", "Customer")
- aggregateid: id of the aggregate root, used as Kafka message key for partition ordering
- type: the type of event (e.g., "OrderCreated")
- payload: JSON structure with actual event contents

Source: Debezium Blog — Reliable Microservices Data Exchange
URL: https://debezium.io/blog/2019/02/19/reliable-microservices-data-exchange-with-the-outbox-pattern/
Confidence: HIGH — real-world implementation.

## Evidence 5: Message Relay Implementations

Claim: There are two patterns for implementing the Message relay: Transaction Log Tailing and Polling Publisher.

Evidence:
- Transaction Log Tailing: tails the database transaction log (MySQL binlog, Postgres WAL, DynamoDB streams), captures INSERTs to outbox table, and publishes to message broker
- Polling Publisher: periodically queries outbox table for unprocessed records and publishes them to message broker

Source: microservices.io
URLs: https://microservices.io/patterns/data/transaction-log-tailing.html, https://microservices.io/patterns/data/polling-publisher.html
Confidence: HIGH — pattern catalog documentation.

## Evidence 6: Idempotent Consumer Requirement

Claim: The message relay might publish a message more than once (e.g., crash after publish but before marking as processed). Therefore, consumers must be idempotent.

Evidence: Consumer tracks `processed_event_id` or message UUID. If event already processed, ignore. This is necessary because message brokers can deliver messages more than once (at-least-once delivery guarantee).

Source: microservices.io Transactional Outbox
URL: https://microservices.io/patterns/data/transactional-outbox.html
Confidence: HIGH — acknowledged drawback of the pattern.

## Evidence 7: Debezium Outbox Event Router

Claim: Debezium Outbox Event Router SMT transforms outbox table changes into Kafka messages with routing based on aggregate type.

Evidence: The SMT expects outbox table columns: id, aggregateid, aggregatetype, type, payload. Configuration:
- `route.by.field=aggregatetype`
- `route.topic.replacement=outbox.event.${routedByValue}`
- Message key = aggregateid for partition ordering
- Message header = id for duplicate detection

Source: Debezium Documentation — Outbox Event Router
URL: https://debezium.io/documentation/reference/stable/transformations/outbox-event-router.html
Confidence: HIGH — reference implementation.

## Evidence 8: Payload Size Recommendation

Claim: Avoid storing large objects in outbox; store minimal identifiers instead.

Evidence: Full object payloads (500 fields) cause database bloat, high I/O, expensive serialization. Better to store event_id and aggregate_id only; consumer can fetch full object if needed.

Source: Topic specification (lab description)
URL: labs/21-outbox-pattern
Confidence: HIGH — best practice guidance.

## Evidence 9: Cleanup Strategy

Claim: Outbox table must be archived or deleted after events are successfully processed.

Evidence: With 1 million events/day, outbox_events table grows indefinitely. Events with `processed_at IS NOT NULL` are candidates for archival or deletion after retention period.

Source: Topic specification (lab description)
URL: labs/21-outbox-pattern
Confidence: HIGH — operational necessity.

Evidence 10: Monitoring Metrics

Claim: Monitor unprocessed event count, oldest unprocessed event age, publish failure rate, retry count, processing throughput.

Evidence: "Oldest Outbox Event" metric is critical. Example SLA thresholds: normal oldest-event age around 2 seconds. Abnormal (47 minutes) indicates serious problem with publisher or message broker. These values are illustrative SLO examples, not universal operational constants.

Source: Topic specification (lab description)
URL: labs/21-outbox-pattern
Confidence: MEDIUM — operational guidance from lab author, presented as illustrative example.