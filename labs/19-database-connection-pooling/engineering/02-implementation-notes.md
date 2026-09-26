# Implementation Notes

## Files Added
- `internal/pool/mockdb.go`: A mock `database/sql/driver` implementation simulating max connections and handshake latency.
- `internal/pool/service.go`: Business logic illustrating safe (external IO separate from DB lock) and unsafe (external IO during DB lock) practices.
- `tests/pool_test.go`: Suite proving overhead, exhaustion, starvation, and concurrent success.
- `cmd/demo/main.go`: Executable printing demonstration.

## Core Design Decisions
- Utilized a custom `database/sql/driver` rather than a live PostgreSQL instance. This ensures deterministic tests for network overhead (`connectDelay`) and hard backend connection bounds (`maxConnections`) without managing external infrastructure.
- `sql.DB` inherently manages a connection pool. Modifying `SetMaxOpenConns` and `SetMaxIdleConns` natively replicates client-side connection pooling semantics.

## Implementation-Specific Choices
- Wait states and connection limits are simulated with memory structures (`atomic.Int32`) rather than actual Postgres processes. 
- A 10ms handshake penalty is enforced in the demo's mock driver to visually amplify the real-world TCP/TLS database connection penalty.

## Known Limitations
- Real-world `pg_stat_activity` waiting states and CPU context switching overhead cannot be truly simulated without a live heavy PostgreSQL installation under load.
- Memory allocation per backend process (which usually leads to out-of-memory under high connection count) is not demonstrated.

## Trade-offs
- A mock database limits fidelity to true PostgreSQL behavior but guarantees repeatable, fast execution across CI environments without flakiness.

## What Is Demonstrated
- The latency penalty of opening direct unpooled connections versus reusing them.
- Connection refusal when client pools collectively exceed the backend database's safe limits.
- Connection pool starvation and application timeout caused by leaked connections held during long network calls.

## What Is Not Demonstrated
- Complex proxy logic (like PgBouncer transaction vs session mode).
- Query performance degradation due to internal DB lock contention.
