# Engineering Audit Plan

Target Lab: labs/27-database-constraints

## Implementation Files
- `internal/dberr/errors.go` — SQLSTATE error type taxonomy (`23502`, `23503`, `23505`, `23514`, `23001`, `23P01`, `40001`), `ConstraintError`, `IsConstraintViolation`, `MapToDomainError`.
- `internal/model/model.go` — `User`, `Order` domain structs.
- `internal/engine/engine.go` — In-memory relational engine: `InsertUserUnsafe`, `InsertUser` (NOT NULL / CHECK / UNIQUE / PARTIAL UNIQUE), `InsertOrder` (NOT NULL / CHECK / FOREIGN KEY), `SoftDeleteUser`, `GetUser`, `GetUsersCount`.
- `internal/store/store.go` — `UnsafeStore` (app-level read-then-write + `InsertUserUnsafe`), `SafeStore` (`RegisterUser`, `RegisterUserPartial`, `CreateOrder`).
- `cmd/demo/main.go` — Standalone demo exercising all constraint paths plus a 50-goroutine concurrency stress test.

## Tests
- `internal/store/store_test.go` — 8 tests: `TestNotNullConstraints`, `TestCheckConstraints`, `TestUniqueConstraint`, `TestForeignKeyConstraint`, `TestPartialUniqueIndex`, `TestConcurrentRegistration_Safe_EnforcesUniqueness`, `TestConcurrentRegistration_Unsafe_SuffersRaceCondition`, `TestErrorClassification`.

## Executable/Demo
- `go run ./cmd/demo` — interactive constraint-violation + concurrency demonstration.

## Approved Research Inputs (implementation alignment targets only)
- Research claims: NOT NULL / CHECK / UNIQUE / FOREIGN KEY / PARTIAL UNIQUE constraints enforced at storage layer; concurrency race eliminated by DB UNIQUE constraint (SQLSTATE 23505); error mapping to domain errors.

## Main Claims To Verify
1. Each constraint type (NOT NULL, CHECK, UNIQUE, FK, partial unique) is enforced and produces the documented SQLSTATE code.
2. `SafeStore` concurrency guarantees exactly 1 winner / N-1 unique violations under 50 simultaneous goroutines (SQLSTATE 23505).
3. `UnsafeStore` loses the race (produces >1 duplicate row) — i.e. the demo actually exhibits the bug it claims to contrast with.
4. `SoftDeleteUser` + partial index allows one active row, rejects duplicate actives, permits soft-deleted duplicates.
5. Error mapping (`MapToDomainError`) maps SQLSTATE codes to domain-friendly messages.
6. `go test -race` is clean (no data races on the in-memory maps/atomics).
7. Demo output matches claimed results (1 success, 49 rejections, integrity intact).

## Commands To Run
- `go build ./...`
- `go vet ./...`
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`
- `go test -race -count=5 -run "Concurrent" ./internal/store/` (repeatability of concurrency claims)

## Primary Risks
- Concurrency tests are timing-dependent (artificial `time.Sleep` in `UnsafeStore`). Risk of flakiness / non-determinism.
- `MapToDomainError` is never exercised by an automated test (coverage gap).
- Design doc references `internal/domain`, `internal/service`, `internal/db`, `internal/errors` packages that do not exist — documentation drift.
- `SoftDeleteUser` cross-index interaction between `emailIndex` (full UNIQUE path) and `activeEmails` (partial path) not covered; mixed registration could violate uniqueness.
- Engineering `03-execution-result.md` lists only 7 tests; omits `TestConcurrentRegistration_Unsafe_SuffersRaceCondition` — mismatch with actual test binary.
