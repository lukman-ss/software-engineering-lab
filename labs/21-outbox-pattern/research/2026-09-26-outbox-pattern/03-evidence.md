## Evidence 1

Claim: The dual-write problem occurs when a service must update its database AND send a message to a message broker, but cannot do both atomically because there is no shared transaction across both systems.

Evidence: "A service command typically needs to create/update/delete aggregates in the database **and** send messages/events to a message broker. For example, a service that participates in a saga needs to update business entities and send messages/events." ... "However, it is not viable to use a traditional distributed transaction (2PC) that spans the database and the message broker. The database and/or the message broker might not support 2PC. And even if they do, it's often undesirable to couple the service to both the database and the message broker."

Source: MicroServices.io - Pattern: Transactional outbox
URL: https://microservices.io/patterns/data/transactional-outbox.html
Confidence: HIGH
Corroborated By: Debezium Blog ("Simply issuing these two requests may lead to potential inconsistencies, though. The reason being that we cannot have one shared transaction that would span the service's database as well as Apache Kafka, as the latter doesn't support to be enlisted in distributed (XA) transactions.")

## Evidence 2

Claim: If the database transaction commits but sending the message fails, data becomes inconsistent. Conversely, if the message is sent but the database transaction rolls back, the message is published for data that doesn't exist.

Evidence: "if a service sends a message after committing the transaction there's no guarantee that it won't crash before sending the message." And from Debezium blog: "it might happen that we end up with having the new purchase order persisted in the local database, but not having sent the corresponding message to Kafka (e.g. due to some networking issue). Or, the other way around, we might have sent the message to Kafka but failed to persist the purchase order in the local database. Both situations are undesirable."

Source: MicroServices.io - Pattern: Transactional outbox
URL: https://microservices.io/patterns/data/transactional-outbox.html
Confidence: HIGH
Corroborated By: Debezium Blog (February 19, 2019)

## Evidence 3

Claim: The Transactional Outbox Pattern stores messages in the database as part of the same transaction that updates business entities, and a separate process (message relay) publishes them to the message broker.

Evidence: "The solution is for the service that sends the message to first store the message in the database as part of the transaction that updates the business entities. A separate process then sends the messages to the message broker."

Source: MicroServices.io - Pattern: Transactional outbox
URL: https://microservices.io/patterns/data/transactional-outbox.html
Confidence: HIGH

## Evidence 4

Claim: The key participants in the Outbox Pattern are: Sender (the service), Database (stores business entities and message outbox), Message outbox (table or NoSQL property), and Message relay (publishes from outbox to broker).

Evidence: "The participants in this pattern are: Sender - the service that sends the message; Database - the database that stores the business entities and message outbox; Message outbox - if it's a relational database, this is a table that stores the messages to be sent. Otherwise, if it's a NoSQL database, the outbox is a property of each database record; Message relay - sends the messages stored in the outbox to the message broker"

Source: MicroServices.io - Pattern: Transactional outbox
URL: https://microservices.io/patterns/data/transactional-outbox.html
Confidence: HIGH

## Evidence 5

Claim: The Outbox Pattern guarantees that messages are sent if and only if the database transaction commits, and messages are sent in order.

Evidence: "Messages are guaranteed to be sent if and only if the database transaction commits" and "Messages are sent to the message broker in the order they were sent by the application."

Source: MicroServices.io - Pattern: Transactional outbox
URL: https://microservices.io/patterns/data/transactional-outbox.html
Confidence: HIGH

## Evidence 6

Claim: Debezium's canonical Outbox table structure includes an id (unique event ID for deduplication), aggregate type, aggregate ID (for partitioning), event type, and payload.

Evidence: "Column | Type | Modifiers; id | uuid | not null; aggregatetype | character varying(255) | not null; aggregateid | character varying(255) | not null; type | character varying(255) | not null; payload | jsonb | not null" ... "id: unique id of each message; can be used by consumers to detect any duplicate events" ... "aggregateid: the id of the aggregate root that is affected by a given event; this id will be used as the key for Kafka messages later on, ensuring all messages of that aggregate will go into the same partition" ... "aggregatetype: used to route events to corresponding topics in Kafka" ... "type: the type of event, e.g. 'Order Created' or 'Order Line Canceled'" ... "payload: a JSON structure with the actual event contents"

Source: Debezium Blog (February 19, 2019)
URL: https://debezium.io/blog/2019/02/19/reliable-microservices-data-exchange-with-the-outbox-pattern/
Confidence: HIGH
Corroborated By: Debezium Documentation (Outbox Event Router)

## Evidence 7

Claim: Consumer applications must be idempotent because the Outbox Pattern provides at-least-once delivery, not exactly-once. Duplicate events can occur when the message relay crashes after publishing but before recording the fact.

Evidence: "The Message relay might publish a message more than once. It might, for example, crash after publishing a message but before recording the fact that it has done so. When it restarts, it will then publish the message again. As a result, a message consumer must be idempotent, perhaps by tracking the IDs of the messages that it has already processed."

Source: MicroServices.io - Pattern: Transactional outbox
URL: https://microservices.io/patterns/data/transactional-outbox.html
Confidence: HIGH
Corroborated By: Debezium Blog ("This is to prevent any duplicate processing of events caused by the 'at least once' semantics of this data pipeline")

## Evidence 8

Claim: There are two approaches to implementing the Message relay: Polling Publisher and Transaction Log Tailing.

Evidence: "There are two patterns for implementing the Message relay: The Transaction log tailing pattern and The Polling publisher pattern."

Source: MicroServices.io - Pattern: Transactional outbox
URL: https://microservices.io/patterns/data/transactional-outbox.html
Confidence: HIGH

## Evidence 9

Claim: Polling Publisher approach works with any SQL database but can be tricky to maintain message ordering. Transaction Log Tailing uses database-specific log (MySQL binlog, Postgres WAL, DynamoDB streams) and is harder to implement but avoids duplicate publishing challenges.

Evidence: Polling Publisher: "Works with any SQL database" and "Tricky to publish events in order" and "Not all NoSQL databases support this pattern." Transaction Log Tailing: "No 2PC", "Guaranteed to be accurate", "Relatively obscure although becoming increasing common", "Requires database specific solutions", "Tricky to avoid duplicate publishing."

Source: MicroServices.io - Polling publisher; MicroServices.io - Transaction log tailing
URL: https://microservices.io/patterns/data/polling-publisher.html
URL: https://microservices.io/patterns/data/transaction-log-tailing.html
Confidence: HIGH

## Evidence 10

Claim: Consumers should store processed event IDs to detect and ignore duplicates. The Debezium outbox example implements this via a `ConsumedMessage` entity tracked in the consuming service's local database.

Evidence: "a message consumer must be idempotent, perhaps by tracking the IDs of the messages that it has already processed." and from Debezium Blog: "The first thing done by onOrderEvent() is to check whether the event with the given UUID has been processed before. If so, any further calls for that same event will be ignored." The code uses a `MessageLog` table with `ConsumedMessage` entity to track `eventId`.

Source: MicroServices.io - Pattern: Transactional outbox
URL: https://microservices.io/patterns/data/transactional-outbox.html
Confidence: HIGH
Corroborated By: Debezium Blog

## Evidence 11

Claim: Payload should be kept small; storing massive payloads causes database bloat, increased I/O, and expensive serialization.

Evidence: "One additional claim... don't store massive payloads in the outbox table." and from topic spec source.

Source: Topic specification (lab content)
Confidence: MEDIUM (Based on lab specification; corroborated by general database principles)
Corroborated By: Debezium Blog notes that payload is JSON structure with event contents, and consumer services only receive relevant data

## Evidence 12

Claim: Processed outbox events should be archived or deleted after a retention period to prevent the outbox table from growing indefinitely.

Evidence: "If your system generates 1 million events/day and never cleans up the outbox_events table, it can grow very large. Successfully processed events (processed_at != NULL) should be archived or deleted after a retention period."

Source: Topic specification (lab content)
Confidence: MEDIUM (Best practice from lab spec, supported by Debezium approach of using INSERT+DELETE for CDC capture)

## Evidence 13

Claim: Recommended monitoring metrics include: unprocessed event count, oldest unprocessed event age, publish failure rate, retry count, and processing throughput. Note: specific numeric thresholds (e.g. "2s normal, 47m serious" for oldest unprocessed age) are example values for educational purposes from lab spec; not derived from external production benchmarks.

Evidence: "Don't just monitor queue length. For Outbox, also monitor: unprocessed event count, oldest unprocessed event age, publish failure rate, retry count, processing throughput."

Source: Topic specification (lab content)
Confidence: MEDIUM

## Evidence 14

Claim: Outbox Pattern provides instant "read your own writes" semantics for the source service, meaning queries immediately reflect newly created data.

Evidence: "By writing to the database first, the source service has instant 'read your own writes' semantics." ... "By synchronously writing to the PurchaseOrder table, the source service benefits from 'read your own writes' semantics. A subsequent query for purchase orders will return the newly persisted order, as soon as that first transaction has been committed."

Source: Debezium Blog (February 19, 2019)
URL: https://debezium.io/blog/2019/02/19/reliable-microservices-data-exchange-with-the-outbox-pattern/
Confidence: HIGH

## Evidence 15

Claim: The Outbox Pattern enables eventual consistency for cross-service data propagation, with delays typically in seconds or sub-second range.

Evidence: "Naturally, such event pipeline between different services is eventually consistent, i.e. consumers such as the shipping service may lag a bit behind producers such as the order service." ... "end-to-end delays of the overall solution are typically low (seconds or even sub-second range), thanks to log-based change data capture which allows for emission of events in near-realtime."

Source: Debezium Blog (February 19, 2019)
URL: https://debezium.io/blog/2019/02/19/reliable-microservices-data-exchange-with-the-outbox-pattern/
Confidence: HIGH

## Evidence 16

Claim: Dead Letter Queue pattern should be implemented for events that repeatedly fail processing.

Evidence: "a more complete implementation should take care of re-trying given messages only for a certain number of times, before re-routing any unprocessable messages to a dead-letter queue or similar."

Source: Debezium Blog (February 19, 2019) — VERIFIED
URL: https://debezium.io/blog/2019/02/19/reliable-microservices-data-exchange-with-the-outbox-pattern/
Confidence: HIGH

## Evidence 17

Claim: The Idempotent Receiver pattern ensures that a consumer processes each message at most once, even when messages are delivered multiple times.

Evidence: "Consumer could save processed_event_id. If event that same already has been processed: Ignore."

Source: Topic specification (lab content)
Confidence: MEDIUM
Corroborated By: Enterprise Integration Patterns - Idempotent Receiver
