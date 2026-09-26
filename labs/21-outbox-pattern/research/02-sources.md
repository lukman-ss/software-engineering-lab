# Sources

## Source 1
Title: Pattern: Transactional outbox
Publisher: microservices.io
URL: https://microservices.io/patterns/data/transactional-outbox.html
Source Tier: 1 (Primary / Original Pattern Author)
Relevance: Original definition of the Transactional Outbox pattern, Problem with dual-write, Forces, Solution, Benefits, Drawbacks, and Issues including duplicate delivery and idempotency requirements.

## Source 2
Title: Pattern: Transaction log tailing
Publisher: microservices.io
URL: https://microservices.io/patterns/data/transaction-log-tailing.html
Source Tier: 1 (Primary / Pattern Implementation)
Relevance: Details Transaction Log Tailing pattern as one of two implementations for the Message relay in the Outbox pattern, with MySQL binlog, Postgres WAL, DynamoDB streams as mechanisms.

## Source 3
Title: Pattern: Polling publisher
Publisher: microservices.io
URL: https://microservices.io/patterns/data/polling-publisher.html
Source Tier: 1 (Primary / Pattern Implementation)
Relevance: Details Polling publisher pattern as alternative to Transaction Log Tailing for the Message relay in the Outbox pattern.

## Source 4
Title: Pattern: Saga
Publisher: microservices.io
URL: https://microservices.io/patterns/data/saga.html
Source Tier: 1 (Primary / Related Pattern)
Relevance: Explains why Transactional Outbox is needed as alternative to 2PC, describes the context where dual-write problems occur in distributed transactions.

## Source 5
Title: Reliable Microservices Data Exchange With The Outbox Pattern
Publisher: Debezium (Gunnar Morling)
URL: https://debezium.io/blog/2019/02/19/reliable-microservices-data-exchange-with-the-outbox-pattern/
Source Tier: 1 (Authoritative Implementation)
Relevance: Complete implementation details including outbox table structure with id, aggregateid, aggregatetype, type, payload columns, CDI event integration, CDC-based event capture, duplicate detection in consumers, and message logging.

## Source 6
Title: Outbox Event Router
Publisher: Debezium Documentation
URL: https://debezium.io/documentation/reference/stable/transformations/outbox-event-router.html
Source Tier: 1 (Reference Implementation Documentation)
Relevance: Detailed configuration of the Outbox Event Router SMT including expected outbox table columns, payload serialization (JSON/Avro), topic routing configuration, field placement options, and distributed tracing support.