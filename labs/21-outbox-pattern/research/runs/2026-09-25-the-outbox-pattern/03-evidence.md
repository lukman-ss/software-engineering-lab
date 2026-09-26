# 03-evidence.md

## Evidence 1

Claim:
Database transactions alone cannot guarantee atomicity across both a database update and message broker publish without distributed transactions (2PC), leading to the dual-write problem.

Evidence:
"The command must atomically update the database and send messages in order to avoid data inconsistencies and bugs. However, it is not viable to use a traditional distributed transaction (2PC) that spans the database and the message broker... But without using 2PC, sending a message in the middle of a transaction is not reliable."

Source:
Microservices.io

URL:
https://microservices.io/patterns/data/transactional-outbox.html

Confidence:
HIGH

Corroborated By:
AWS Prescriptive Guidance ("If the flight database update fails but the notification is sent out, the payment service will process the payment based on the event notification... If the flight service fails after committing the transaction, this might result in the event notification not being sent.") and Debezium Blog ("we cannot have one shared transaction that would span the service’s database as well as Apache Kafka, as the latter doesn’t support to be enlisted in distributed (XA) transactions").

Notes:
Two-phase commit (2PC) is widely rejected in modern distributed systems due to performance overhead and tight coupling.

## Evidence 2

Claim:
The Outbox pattern resolves the dual-write problem by saving the domain state and the outgoing event in the same local database transaction.

Evidence:
"The solution is for the service that sends the message to first store the message in the database as part of the transaction that updates the business entities. A separate process then sends the messages to the message broker."

Source:
Microservices.io

URL:
https://microservices.io/patterns/data/transactional-outbox.html

Confidence:
HIGH

Corroborated By:
AWS Prescriptive Guidance ("When the flight table is updated, the outbox table is also updated in the same transaction.") and Debezium Blog ("When receiving a request for placing a purchase order, not only an INSERT into the PurchaseOrder table is done, but, as part of the same transaction, also a record representing the event to be sent is inserted into that outbox table.").

Notes:
This guarantees atomicity: either both the domain state change and the event are persisted, or neither is.

## Evidence 3

Claim:
A separate message relay reads events from the outbox and publishes them to the message broker, commonly implemented via Polling Publisher or Transaction Log Tailing / Change Data Capture (CDC).

Evidence:
"There are two patterns for implementing the Message relay: The Transaction log tailing pattern, The Polling publisher pattern."

Source:
Microservices.io

URL:
https://microservices.io/patterns/data/transactional-outbox.html

Confidence:
HIGH

Corroborated By:
AWS Prescriptive Guidance (describes both scheduled polling and CDC with DynamoDB Streams/Kinesis) and Debezium Blog (describes log-based Change Data Capture with Debezium and Postgres WAL).

Notes:
Polling is simpler to set up but can suffer from latency and database load; log-based CDC provides near-real-time streaming with less query overhead on the database.

## Evidence 4

Claim:
The Outbox pattern provides at-least-once message delivery, making consumer idempotency mandatory.

Evidence:
"The Message relay might publish a message more than once. It might, for example, crash after publishing a message but before recording the fact that it has done so. When it restarts, it will then publish the message again. As a result, a message consumer must be idempotent, perhaps by tracking the IDs of the messages that it has already processed."

Source:
Microservices.io

URL:
https://microservices.io/patterns/data/transactional-outbox.html

Confidence:
HIGH

Corroborated By:
AWS Prescriptive Guidance ("The events processing service might send out duplicate messages or events, so we recommend that you make the consuming service idempotent by tracking the processed messages.") and Debezium Blog ("This is to prevent any duplicate processing of events caused by the 'at least once' semantics of this data pipeline.").

Notes:
Consuming services often use an idempotency key (such as the event UUID) stored in an inbox table or deduplication log to ignore duplicate events.

## Evidence 5

Claim:
Event payloads should avoid unnecessarily large representations to prevent database bloat, serialization costs, and I/O degradation.

Evidence:
NOT VERIFIED

Source:
N/A

URL:
N/A

Confidence:
LOW

Corroborated By:
None

Notes:
While widely acknowledged as an engineering best practice (referencing an ID vs. sending massive bloated aggregates), explicit quantitative data or strong design restrictions on payload size were not detailed in the retrieved primary references beyond general warnings and aggregate designs.
