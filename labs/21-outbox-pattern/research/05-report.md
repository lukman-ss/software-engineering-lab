# Research Report

## Research Question
How does the Transactional Outbox pattern solve the dual-write problem where database commit succeeds but event/queue publishing fails?

## Executive Summary
The Transactional Outbox pattern solves the dual-write problem by moving the second write (message publishing) into the same transactional boundary as the first write (database update). Instead of updating the database and publishing to a message broker in two separate operations, the pattern stores the message in an "outbox" table within the same database transaction as the business data update. This ensures atomicity: either both the business data and the message are persisted, or neither is. A separate, asynchronous process (the message relay) then reads from the outbox table and publishes messages to the message broker. If the relay fails, messages remain safely in the outbox table for later retry. This pattern eliminates the possibility of having committed database changes without corresponding messages.

## Findings

### Finding 1: The Dual-Write Problem and Failed Transactions

Claim: Updating a database and publishing a message cannot be made atomic using traditional distributed transactions (2PC) or database transactions alone.

Evidence: Traditional database transactions only control database operations, not external systems like Redis, RabbitMQ, Kafka, or HTTP endpoints. Without 2PC, there is no way to atomically commit changes across the database and message broker. If using 2PC, the database and/or message broker might not support it, and coupling the service to both resources is often undesirable.

Sources:
- microservices.io Transactional Outbox: "However, it is not viable to use a traditional distributed transaction (2PC) that spans the database and the message broker. The database and/or the message broker might not support 2PC."
- Labs specification: "Transaction database hanya mengontrol database. Ia tidak bisa melakukan: ROLLBACK Redis, ROLLBACK RabbitMQ, ROLLBACK Kafka, ROLLBACK HTTP Request"

Confidence: HIGH — multiple authoritative sources agree.

### Finding 2: Outbox Pattern Guarantees Atomicity

Claim: The Outbox pattern guarantees that messages are sent if and only if the database transaction commits.

Evidence: By storing the message in an outbox table as part of the same database transaction that updates business entities, both operations succeed or fail together. The message relay then asynchronously publishes messages from the outbox table. This ensures there is never a state where database changes are committed but messages are lost.

Sources:
- microservices.io Transactional Outbox: "The solution is for the service that sends the message to first store the message in the database as part of the transaction that updates the business entities. A separate process then sends the messages to the message broker."
- Debezium Blog: "By synchronously writing to the PurchaseOrder table, the source service benefits from 'read your own writes' semantics. At the same time, we get reliable, asynchronous, eventually consistent data propagation to other services via Apache Kafka."

Confidence: HIGH — canonical pattern definition and implementation validation.

### Finding 3: Message Relay Implementation Alternatives

Claim: Two patterns exist for implementing the message relay: Transaction Log Tailing and Polling Publisher.

Evidence: The message relay can either tail the database transaction log (capturing INSERT operations on the outbox table in real-time via binlog/WAL/streams) or periodically poll the outbox table for unprocessed records. Each has tradeoffs: log-tailing is lower latency but database-specific; polling works with any SQL database but has higher latency and ordering challenges.

Sources:
- microservices.io Transaction Log Tailing: "Tail the database transaction log and publish each message/event inserted into the outbox to the message broker."
- microservices.io Polling Publisher: "Publish messages by polling the database’s outbox table."
- Debezium Blog: "Debezium tails the transaction log ("write-ahead log", WAL) of the order service’s Postgres database in order to capture any new events in the outbox table and propagates them to Apache Kafka."

Confidence: HIGH — pattern catalog documentation and real-world implementation.

### Finding 4: Idempotent Consumer Requirement

Claim: Outbox pattern provides at-least-once delivery, requiring consumers to be idempotent to handle duplicate messages.

Evidence: The message relay might publish a message more than once (e.g., if it crashes after publishing but before marking the message as processed). Therefore, consumers must track processed message IDs (e.g., via UUID or eventId) and ignore duplicates. This is necessary because message brokers can deliver messages more than once.

Sources:
- microservices.io Transactional Outbox: "The Message relay might publish a message more than once... As a result, a message consumer must be idempotent, perhaps by tracking the IDs of the messages that it has already processed."
- Debezium Blog: "Propagating the event UUID as a Kafka message header allows for an efficient detection and exclusion of duplicates in the consumer."

Confidence: HIGH — acknowledged limitation across all sources.

### Finding 5: Outbox Table Structure and Payload Design

Claim: Outbox table design has two common approaches: thin events (minimal identifiers) and fat events (full state payloads), each with different trade-offs.

Evidence: Thin events store only essential fields (event ID, aggregate ID, type) and require consumers to query the source service for full state when needed. Fat events (Event-Carried State Transfer) embed the complete domain state in the payload, avoiding downstream RPC overhead but increasing storage and network bandwidth. The Debezium outbox implementation demonstrates the fat event pattern with full JSON objects containing line items, prices, dates, and customer details. Table structure typically includes: id (uuid), aggregateid, aggregatetype, type, payload (json/jsonb).

Sources:
- Debezium Blog outbox table structure showing: id, aggregatetype, aggregateid, type, payload
- Debezium Blog canonical example showing fat events with full order state (line items, prices, dates, customer details)
- Labs specification: "Jangan Simpan Payload Raksasa" warning against storing entire objects with 500 fields
- Debezium Outbox Event Router documentation detailing expected table columns

Confidence: HIGH — sources confirm both design approaches exist with different trade-offs.

### Finding 6: Operational Requirements: Cleanup and Monitoring

Claim: Outbox systems require cleanup of processed events and monitoring of key metrics to detect problems.

Evidence: Without cleanup, the outbox_events table grows indefinitely (1M events/day → very large table after one year). Processed events (processed_at != NULL) should be archived or deleted after a retention period. Critical metrics include: unprocessed event count, oldest unprocessed event age, publish failure rate, retry count, and processing throughput. Example SLA thresholds: normal oldest-event age around 2 seconds, while 47+ minutes indicates serious problems. These are illustrative examples, not universal constants, and must be tuned to specific service-level objectives.

Sources:
- Labs specification: "Outbox Juga Harus Dibersihkan" and monitoring section detailing key metrics
- Debezium Blog: MessageLog table tracking processed events with UUID and timestamp
- microservices.io: Implicit in pattern requiring message relay to mark events as processed

Confidence: MEDIUM — operational guidance from lab author and implementation examples.

## Areas of Agreement
All sources agree on:
1. Dual-write problem exists when updating database and publishing message separately
2. Traditional database transactions cannot span database + message broker
3. Outbox pattern stores message in database transaction for atomicity
4. Two message relay implementations: polling and transaction log tailing
5. Pattern provides at-least-once delivery (not exactly-once)
6. Consumers must be idempotent to handle duplicate messages
7. Outbox table requires id, aggregateid, aggregatetype, type, payload columns

## Limitations
- No universal numeric recommendations for timeouts, retry counts, or table sizes — all must be tuned to specific SLAs and traffic patterns
- Cleanup strategies vary (archive vs delete) and require operational procedures
- Debezium-focused sources emphasize CDC-based implementation, while pattern catalog treats relay patterns more generally

## Conclusion
The Transactional Outbox pattern solves the dual-write problem by ensuring database updates and message publishing happen within the same atomic transaction boundary. By storing messages in an outbox table as part of the business transaction, the pattern guarantees that messages are sent if and only if database changes commit. While this introduces eventual consistency for message delivery (at-least-once with possible duplicates), it eliminates the risk of inconsistent state where database changes exist without corresponding messages. Successful implementation requires idempotent consumers, appropriate outbox table design, reliable message relay implementation (polling or log-tailing), operational cleanup procedures, and monitoring of key health metrics.