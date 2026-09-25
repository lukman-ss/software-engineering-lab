# Engineering Design

Target Lab: labs/13-backward-compatibility
Research Status: APPROVED_WITH_WARNINGS

## Concept To Prove
Prove the **Expand -> Migrate -> Contract (Parallel Change)** pattern for database schema and API evolution:
1. Systems can transition from a 1:1 legacy schema (`users.phone`) to a 1:N schema (`user_phones` table and array payload) without service downtime or breaking legacy consumers.
2. Dual-write ensures data consistency across both representations during transition.
3. Batched, idempotent, and resumable backfill safely migrates historical data without table lockup or data corruption.
4. Fallback reading (dual-read) enables reads before backfill completes.
5. Observability and feature flags enable safe canary switching and instant rollback without data loss.
6. Contract phase safely drops legacy schema and endpoints only after legacy access reaches zero.

## Expected Behavior
- **Legacy Consumer (V1)**: Receives and parses payload containing string `phone`. Continues working unchanged across Expand and Migrate phases. Receives deprecation headers (`Deprecation: true`, `Sunset: <date>`).
- **New Consumer (V2)**: Receives and parses payload containing `phones` array with primary flags.
- **Dual-Write**: New mutations write atomically to both legacy `users` table and new `user_phones` table under feature flag control.
- **Backfill Worker**: Iterates legacy records in configurable batches, checkpointing the last processed ID. Repeated runs are idempotent.
- **Fallback Read**: When reading an un-backfilled user in V2 mode, system falls back to legacy field and lazily hydrates or serves data without errors.
- **Rollback Safety**: If V2 code is rolled back to V1 while dual-write is active, V1 continues serving and writing to legacy fields without data loss.
- **Contract Enforcement**: Contract phase shuts down legacy write and removes legacy field after zero legacy traffic is observed.

## Failure Scenario
1. **Destructive Alteration**: Direct drop/rename of legacy field breaks legacy V1 consumers immediately (deserialization crash / HTTP 500).
2. **Missing Backfill (Premature Read Switch)**: Switching reads to new schema before backfill finishes causes missing phone records (data starvation) unless fallback reading is active.
3. **Dual-Write Drift / Desynchronization**: Unsynchronized writes to legacy and new stores cause diverging state. Reconciler detects discrepancies.
4. **Premature Contract / Rollback after Dual-Write stopped**: Stopping dual-write while legacy consumers or rollbacks still occur leads to silent data loss on the legacy side.

## Success Criteria
1. Automated tests pass covering:
   - Happy path: V1 and V2 consumers interacting across Expand, Migrate, and Contract phases.
   - Failure path: Data drift detection, non-idempotent duplicate prevention, premature contract detection.
   - Edge cases: Empty records, missing phones, batch boundaries in backfill.
   - Concurrency: Parallel reads, writes, and backfills pass `-race` without data races or deadlocks.
2. Demo binary executes cleanly, demonstrating the entire lifecycle step-by-step with real data and metric verification.

## Architecture
```text
[ Client V1 (Legacy) ]     [ Client V2 (Modern) ]
          │                         │
          ▼                         ▼
┌──────────────────────────────────────────────────┐
│              HTTP / Service Layer                │
│  - Transformation / Deprecation Pipeline         │
│  - Feature Flags (WriteMode, ReadMode)           │
│  - Observability Metrics (Legacy/New Traffic)    │
└─────────┬────────────────────────────────┬───────┘
          │                                │
          ▼                                ▼
┌──────────────────┐             ┌─────────────────┐
│ Legacy Storage   │◄──Backfill──┤ Modern Storage  │
│ (users.phone)    │   Worker    │ (user_phones)   │
└──────────────────┘             └─────────────────┘
```

## Components
1. `Storage`: Thread-safe in-memory database simulating relational tables `users` (legacy column `phone`) and `user_phones` (child table with `user_id`, `number`, `is_primary`). Supports atomic dual-writes and table contract (drop column).
2. `FeatureFlagManager`: Controls migration stages:
   - `WriteMode`: `LegacyOnly`, `DualWrite`, `NewOnly`.
   - `ReadMode`: `LegacyOnly`, `FallbackRead`, `NewOnly`.
   - `ContractApplied`: Boolean flag indicating legacy schema dropped.
3. `MetricsCollector`: Thread-safe counters for `LegacyReadHits`, `NewReadHits`, `DualWriteCount`, `DualWriteErrors`, `DriftDetectedCount`.
4. `BackfillWorker`: Resumable, batch-oriented data migration worker using ID checkpoints and idempotent upsert logic.
5. `CompatService`: Domain facade orchestrating reads, writes, reconciliation audits, and deprecation lifecycle.
6. `HTTPHandler`: Standard HTTP handlers returning JSON contracts, deprecation headers (`Deprecation: @epoch`, `Sunset: @epoch`), and status codes.

## Test Strategy
- Unit tests (`internal/compat/service_test.go`): Validate serialization backward compatibility, dual-write logic, backfill batching, checkpoint resuming, and data reconciliation.
- Integration & Lifecycle tests (`tests/migration_test.go`): Execute the full 5-stage deployment sequence (Expand -> DualWrite -> Backfill -> ReadSwitch -> Contract) and verify rollback safety.
- Concurrency tests (`tests/concurrency_test.go`): Run concurrent writes, reads, and backfill workers under `go test -race` to prove thread safety.

## Execution Plan
1. Implement domain model, storage, feature flags, metrics, service, and backfill in `internal/compat/`.
2. Implement HTTP handlers with deprecation headers in `internal/compat/handler.go`.
3. Implement end-to-end demo binary in `cmd/demo/main.go`.
4. Write comprehensive tests in `internal/compat/service_test.go`, `tests/migration_test.go`, and `tests/concurrency_test.go`.
5. Execute `go test ./...` and `go test -race ./...`.
6. Execute `go run ./cmd/demo`.
7. Document implementation notes in `engineering/02-implementation-notes.md`.
8. Document command outputs in `engineering/03-execution-result.md`.
9. Update README.md to match code.

## Implementation Decisions
- **In-Memory Thread-Safe Storage**: Standard library `sync.RWMutex` with map-based tables to eliminate external DB dependencies (Postgres/Docker) while preserving exact relational mechanics (foreign keys, atomic dual-writes, nullable fields, column drop simulation). Marked with `ponytail: in-memory mock storage; replace with database/sql for persistent store.`
- **Batch Backfill Mechanism**: Implemented with explicit batch size and stateful checkpoint tracker to simulate chunked DB cursor updates without locking.
- **Standard HTTP Headers**: RFC 8594-aligned `Deprecation` and `Sunset` headers implemented for HTTP API compatibility.
