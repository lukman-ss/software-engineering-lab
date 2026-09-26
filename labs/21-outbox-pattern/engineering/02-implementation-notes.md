# Implementation Notes

## Files Added
- `cmd/demo/main.go`: Interactive executable demonstrating dual-write failure, transactional outbox atomicity, and idempotent message consumption.
- `internal/outbox/model.go`: Domain structures (`Order`, `OutboxMessage`) and statuses.
- `internal/outbox/db.go`: Thread-safe transactional memory database supporting transactional staging, atomic commit, rollback, and status updates.
- `internal/outbox/broker.go`: Thread-safe mock message broker supporting fault-injection (`SetFailNext`) to test resilience.
- `internal/outbox/service.go`: Business logic illustrating naive dual-write vs atomic outbox writes.
- `internal/outbox/relay.go`: Asynchronous polling relay worker querying `PENDING` outbox records and updating status to `PROCESSED`.
- `internal/outbox/consumer.go`: Downstream consumer tracking processed message IDs for idempotency.
- `tests/outbox_test.go`: Suite covering happy path, tx rollback, duplicate handling, dual-write flaw, and concurrent writes.

## Core Design Decisions
- **In-Memory Transactional DB**: Standardized transactional boundary (`Tx` struct with staged maps) to demonstrate exact database transactional guarantees without external DB server setup.
- **Polling Publisher Relay**: Asynchronous background polling goroutine fetching pending records, publishing to broker, and marking as processed.
- **Consumer ID Tracking**: Deduplication map keyed by event UUID ensuring idempotency on duplicate deliveries.

## Implementation-Specific Choices
- Custom `MockBroker` with configurable network failures to test dual-write inconsistency and relay recovery.
- Event payload marshaled into standard JSON format (`OrderCreated`).

## Known Limitations
- High-throughput DB table partitioning is omitted for clarity.
- Log-tailing (CDC/Debezium) is not implemented; polling publisher pattern is used instead.

## Trade-offs
- **Polling Latency vs Operational Complexity**: Polling relies on configurable ticker interval, introducing minor delivery latency compared to transaction log tailing, but requires zero database plugin configuration.

## What Is Demonstrated
- Dual-write failure where database commits but broker fails.
- Atomic persistence of domain entity (`Order`) and `OutboxMessage`.
- At-least-once message dispatch via outbox polling relay.
- Idempotent deduplication on consumer side.

## What Is Not Demonstrated
- Distributed transaction log tailing (Debezium / Postgres WAL tailing).
- Dead-letter queue (DLQ) routing for non-retryable poison payloads.
