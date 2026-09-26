1. **Atomic persistence is guaranteed within a single transaction** — Order and outbox message are written to the same database transaction. Either both persist or neither does.

2. **Rollback discards both records** — If any error occurs before commit, `Rollback()` ensures no order or outbox record is persisted.

3. **Relay provides asynchronous, decoupled dispatch** — A background worker periodically polls pending outbox messages and publishes them to the broker.

4. **At-least-once delivery requires idempotent consumers** — Relay may publish the same message multiple times if it crashes before marking it processed. Consumer deduplication by event ID prevents duplicate processing.

5. **Dual-write problem is real and demonstrable** — Writing to database then publishing outside the transaction leaves the system inconsistent when broker fails.

6. **In-memory simulation vs production realities** — The lab uses in-memory maps for clarity. Production requires persistent storage, cleanup, monitoring, and retry policies.

7. **Polling publisher is simpler than CDC** — The implementation uses polling instead of transaction log tailing (Debezium-style). Polling has latency but works with any SQL database.

8. **Thread-safe implementation verified** — Concurrent writes under `go test -race ./...` pass with zero race conditions detected.

9. **Status transitions are explicit** — Outbox messages start as `PENDING`, transition to `PROCESSED` only after successful broker publish.

10. **No exactly-once semantics** — Outbox guarantees atomic persistence but does not eliminate at-least-once delivery. Idempotent consumers are mandatory.
