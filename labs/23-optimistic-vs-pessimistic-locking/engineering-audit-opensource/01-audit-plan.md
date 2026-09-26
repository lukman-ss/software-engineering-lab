# Engineering Audit Plan

Target Lab: labs/23-optimistic-vs-pessimistic-locking
Implementation Files:
- internal/inventory/model.go
- internal/inventory/service.go
- internal/inventory/store.go
- cmd/demo/main.go
Tests:
- tests/locking_test.go
Executable/Demo: go run ./cmd/demo
Approved Research Inputs: research/05-report.md (verified claims 1-9 in research-audit/03-claim-audit.md)
Main Claims To Verify:
1. Naive read-modify-write demonstrates lost update anomaly under concurrent goroutines
2. Pessimistic locking prevents concurrency conflicts and maintains exact inventory
3. Optimistic locking detects version conflicts and safely rejects stale writes
4. Optimistic locking with retry converges successfully under contention
5. Atomic single-statement operations guarantee consistency under concurrent load
6. No data corruption or race conditions exist
Commands To Run:
- go test ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Probabilistic nature of concurrent tests (lost update and conflict detection rely on timing)
- Missing edge case tests (error conditions, retry exhaustion)
- TestOptimisticLockingWithRetry does not assert all goroutines succeeded
- No verification of research claims vs implementation alignment (per pipeline override)