# Engineering Design

Target Lab: labs/21-outbox-pattern
Research Status: APPROVED

## Concept To Prove
The Transactional Outbox pattern guarantees that state changes to an entity and corresponding domain events are written atomically within a single local transaction, eliminating the dual-write problem. An asynchronous Message Relay then publishes these events to a message broker with at-least-once delivery, requiring consumers to be idempotent.

## Expected Behavior
1. In a single local transaction, an entity update and an outbox event are saved.
2. If the transaction commits, both records persist. If it rolls back, neither persists.
3. A decoupled Message Relay periodically polls or extracts unprocessed events from the outbox table.
4. The Message Relay dispatches events to an external message broker.
5. Upon successful delivery to the broker, the relay marks the event as processed (or deletes it).
6. Downstream consumers track event IDs to handle duplicate deliveries gracefully (idempotency).

## Failure Scenario
1. **Dual-write failure**: Direct write to DB succeeds, but broker is unavailable or write fails, leaving system inconsistent. (Outbox prevents this by avoiding direct broker writes in the application transaction).
2. **Transaction Rollback**: An error occurs before commit; both entity change and outbox message are discarded. No spurious event is sent.
3. **Relay Crash / Duplicate Delivery**: Relay publishes to broker, but crashes before updating outbox status. On restart, it republishes the message. The idempotent consumer detects the duplicate ID and ignores it.

## Success Criteria
- Demonstration of atomic write (entity + outbox message).
- Demonstration of rollback safety (no outbox record on abort).
- Working decoupled relay picking up outbox messages and sending to mock broker.
- Demonstration of at-least-once delivery and consumer-side idempotent deduplication.
- Thread-safe implementations with zero race conditions.

## Architecture
- `Database`: Thread-safe mock relational DB supporting transactions (`Tx`) with tables:
  - `Entities` (e.g., Orders or Users)
  - `Outbox` (Message ID, Payload, Status, CreatedAt)
- `Broker`: Thread-safe mock message broker/queue receiving published messages.
- `OrderService`: Writes order state and outbox record inside the same transaction.
- `Relay`: Asynchronous polling worker that queries pending outbox records, publishes to broker, and updates status.
- `Consumer`: Reads from broker and uses an idempotency key set to avoid double-processing.

## Components
- `internal/outbox/db.go`: In-memory transactional DB simulation (`DB`, `Tx`).
- `internal/outbox/broker.go`: Thread-safe broker interface & implementation.
- `internal/outbox/service.go`: Business service implementing order creation with transactional outbox.
- `internal/outbox/relay.go`: Polling message relay.
- `internal/outbox/consumer.go`: Idempotent consumer.

## Test Strategy
- Unit & integration tests in `tests/outbox_test.go`:
  - Happy path: Save entity + outbox, relay dispatches to broker, consumer processes it.
  - Rollback path: Transaction rollback discards outbox entry; relay sends nothing.
  - Duplicate delivery / Idempotency: Relay publishes duplicate message, consumer safely deduplicates.
  - Concurrency: Multiple concurrent orders and relays without data race (`-race`).

## Execution Plan
1. Set up Go module.
2. Implement components in `internal/outbox`.
3. Create automated tests in `tests/outbox_test.go`.
4. Create executable demo in `cmd/demo/main.go`.
5. Run tests, race detector, demo.
6. Record execution results in `03-execution-result.md` and implementation notes in `02-implementation-notes.md`.

## Implementation Decisions
- Implementation Decision: In-memory transaction model used to simulate relational DB transaction semantics (BEGIN, COMMIT, ROLLBACK) without external CGO/SQLite dependencies.
- Implementation Decision: Polling publisher pattern chosen for the relay mechanism as outlined in the research report.
- Implementation Decision: Consumer idempotency implemented via message ID deduplication cache.
