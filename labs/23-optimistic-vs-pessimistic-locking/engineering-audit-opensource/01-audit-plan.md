# Engineering Audit Plan

Target Lab: labs/23-optimistic-vs-pessimistic-locking
Implementation Files: 
  - internal/inventory/model.go
  - internal/inventory/store.go
  - internal/inventory/service.go
  - cmd/demo/main.go
Tests: tests/locking_test.go
Executable/Demo: cmd/demo/main.go
Approved Research Inputs: N/A (per pipeline override - implementation and tests audit only)
Main Claims To Verify:
  1. Naive read-modify-write produces lost update anomaly under concurrency.
  2. Pessimistic locking (SELECT ... FOR UPDATE) prevents lost update and ensures correct final stock.
  3. Optimistic locking (version check) detects conflicts and prevents lost update, with retry logic ensuring eventual consistency.
  4. Atomic single-statement update (UPDATE ... WHERE stock >= qty) prevents lost update without locks.
  3. All strategies handle error conditions (insufficient stock, invalid quantity) correctly.
Commands To Run:
  - go vet ./...
  - go test ./...
  - go test -race ./...
  - go run ./cmd/demo
Primary Risks:
  - Data races on atomic counters if reads overlap with writes (mitigated by synchronization in current code).
  - Missing test coverage for error cases (negative/zero quantity, insufficient stock for non-pessimistic strategies).
  - Potential flakiness in TestOptimisticLockingConflict due to dependence on timing for conflict generation.
  - Version guard redundancy due to global store mutex (s.mu) serializing writes, though correctness unaffected.