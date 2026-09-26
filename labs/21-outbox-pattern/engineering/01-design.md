# Engineering Design

Target Lab: labs/21-outbox-pattern
Research Status: APPROVED

## Concept To Prove
The Transactional Outbox pattern guarantees atomic persistence of business operations and outgoing event records within a single database transaction, resolving the dual-write problem. An asynchronous polling relay worker then reads pending events, dispatches them to a broker/consumer, and updates status, providing at-least-once delivery semantics coupled with idempotent consumer handling.

## Expected Behavior
1. Creating an order atomically inserts the order record into the `orders` table and an outbox event into the `outbox_events` table within a single SQL transaction.
2. If database transaction rolls back, neither the order nor outbox event is persisted.
3. Polling relay fetches pending outbox events, dispatches them to the message broker, and marks outbox events as `PROCESSED`.
4. Consumer tracks processed event IDs (`message_log` table / set) to enforce idempotency when processing duplicate events.
5. Outbox cleanup job removes or archives processed outbox records after retention interval.

## Failure Scenario
1. **Direct broker write failure (dual-write without outbox)**: DB commit succeeds but broker write fails, leading to state inconsistency.
2. **Relay crash / Network error during publish**: Relay dispatches message but fails before marking outbox record as `PROCESSED`. On retry, duplicate message is sent. Consumer detects duplicate via event ID and ignores payload execution.
3. **Transaction rollback**: Invalid order parameters cause DB transaction rollback. No event is added to outbox, zero messages dispatched.

## Success Criteria
- 100% atomicity between business state and outbox state.
- Zero lost events under broker network disconnect / relay retries.
- Zero duplicate processing by idempotent consumers despite at-least-once relay delivery.
- Cleanup worker successfully purges processed events.
- All unit and concurrency tests pass with zero data races (`go test -race ./...`).

## Architecture
```text
[ Client ] -> [ Order Service ]
                    |
                    v (Single DB Tx)
      +-----------------------------+
      |  SQLite Database            |
      |  - orders table             |
      |  - outbox_events table      |
      +-----------------------------+
                    ^
                    | Poll (FOR UPDATE / Lock)
            [ Outbox Relay ]
                    |
                    v Publish
            [ Message Broker ]
                    |
                    v Deliver
           [ Idempotent Consumer ] -> (Consumer DB message_log)
```

## Components
- `SQLite DB`: Embedded SQL engine storing business data (`orders`) and outbox queue (`outbox_events`).
- `OrderService`: Business service performing atomic order creation + outbox insertion.
- `OutboxRelay`: Background worker polling pending events, delivering to broker, marking processed, with retry / cleanup support.
- `MessageBroker`: In-memory broker interface simulating async event transport with configurable fault injection (failures, latency, retries).
- `IdempotentConsumer`: Event subscriber keeping an processed event registry to prevent duplicate processing.

## Test Strategy
1. **Unit Tests**:
   - Atomic Tx: Order creation + Outbox event insertion.
   - Rollback handling on invalid order.
2. **Relay & Idempotency Tests**:
   - Outbox polling dispatch and marking processed.
   - Duplicate message delivery handling by idempotent consumer.
   - Broker failure retry mechanism.
3. **Concurrency & Race Detector**:
   - Multiple concurrent order creation requests and relay polling cycles.
   - Go race detector verification (`go test -race`).

## Execution Plan
1. Create SQLite DB schema (`orders`, `outbox_events`, `consumer_log`).
2. Implement outbox models, DB repository, OrderService, OutboxRelay, and IdempotentConsumer.
3. Write test suite in `tests/outbox_test.go`.
4. Implement `cmd/demo/main.go` demonstrating dual-write comparison vs transactional outbox.
5. Run tests & demo; output execution results to `engineering/03-execution-result.md`.

## Implementation Decisions
- **Database Engine**: `modernc.org/sqlite` or stdlib `database/sql` with in-memory SQLite / mock DB driver to ensure zero CGO dependencies and portable test execution.
- **Relay Mechanism**: Polling Publisher model with configurable interval and batch size.
- **Event Structure**: JSON-serialized payloads with unique UUID event identifiers.
