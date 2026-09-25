# Implementation Notes

## Files Added
- `internal/compat/model.go`: Domain models, legacy and modern DTOs, enriched responses.
- `internal/compat/store.go`: Thread-safe in-memory store simulating relational tables `users` and `user_phones`, dual-write mutations, and column dropping.
- `internal/compat/flags.go`: Feature flags managing write modes (`WriteLegacyOnly`, `WriteDual`, `WriteNewOnly`) and read modes (`ReadLegacyOnly`, `ReadFallback`, `ReadNewOnly`).
- `internal/compat/metrics.go`: Observability counters tracking legacy/new read traffic, dual writes, errors, backfill progress, and data drift.
- `internal/compat/backfill.go`: Resumable, idempotent batch backfill worker with checkpoint state.
- `internal/compat/service.go`: Core domain service coordinating database mutations, reads, fallback reads, and data reconciliation.
- `internal/compat/handler.go`: HTTP handler providing `/api/v1/users` and `/api/v2/users` with RFC 8594 `Deprecation` and `Sunset` headers.
- `internal/compat/service_test.go`: Unit tests for serialization compatibility, idempotent backfills, fallback read lazy migration, and contract enforcement.
- `tests/migration_test.go`: Lifecycle integration tests for the full 5-stage deployment sequence and rollback scenarios.
- `tests/concurrency_test.go`: Parallel tests validating thread safety of reads, writes, and backfills under `-race`.
- `cmd/demo/main.go`: Interactive runnable CLI demonstration of the complete Expand -> Migrate -> Contract lifecycle.
- `engineering/01-design.md`: Engineering design document.
- `engineering/02-implementation-notes.md`: Implementation choices, trade-offs, and boundaries.
- `engineering/03-execution-result.md`: Recorded build, test, race, and demo execution logs.

## Core Design Decisions
1. **Parallel Change (Expand-Migrate-Contract)**:
   - Split structural changes into 3 distinct operational phases.
   - Preserved legacy `phone` string field while adding new `phones` array field in both API payloads and storage schemas.
2. **Resumable and Idempotent Backfill**:
   - Backfill processes records using a chunked ID cursor (`last_processed_id`).
   - Idempotency guaranteed by checking existing phone entries before insertion, ensuring rerun safety without duplication.
3. **Fallback Reading (Dual Read)**:
   - Modern read requests for un-migrated records fall back to legacy `users.phone` and lazily backfill into `user_phones`.
4. **Data Reconciliation**:
   - Built-in drift detection audit comparing legacy columns against child table primary records before advancing deployment stages.
5. **Contract Guard**:
   - Application of Contract phase checks metric counters to ensure legacy traffic has dropped to zero before dropping columns.

## Implementation-Specific Choices
- **In-Memory Store**: Used `sync.RWMutex` with Go maps to model relational database behavior (tables, foreign keys, nullable columns, dropping columns) without external DBMS dependencies (`ponytail: in-memory mock storage; replace with database/sql for persistent store.`).
- **HTTP Deprecation Pipeline**: Implemented standard `Deprecation: true` and `Sunset: <date>` HTTP response headers in `internal/compat/handler.go` to signal deprecation windows to legacy callers.

## Known Limitations
- In-memory storage does not persist data across process restarts; replacing with `database/sql` is required for disk persistence.
- Cross-table transactions in production PostgreSQL require explicit ACID transactions (`BEGIN ... COMMIT`) or transactional outbox pattern when scaling beyond single-instance databases.

## Trade-offs
- **Dual-Write Latency vs Consistency**: Synchronous dual-writes increase write latency during the migration phase, but provide instantaneous zero-data-loss rollback capability.
- **In-Memory Concurrency vs Relational Locks**: Go mutexes simulate row/table locking without database engine overhead.

## What Is Demonstrated
- Zero-downtime structural schema migration (1:1 to 1:N).
- Backward-compatible JSON serialization consuming additive fields across legacy (V1) and modern (V2) clients.
- Safe rollback without data loss during dual-write phase.
- Data drift detection through reconciliation.
- Contract phase enforcement triggered by zero-traffic observability verification.

## What Is Not Demonstrated
- Distributed transaction coordination across separate microservices (e.g. Outbox pattern / Kafka-based event streaming).
- PostgreSQL DDL lock timeouts (`SET lock_timeout = '2s'`) and `CREATE INDEX CONCURRENTLY` runtime execution.
