# Engineering Audit Plan

Target Lab: labs/23-optimistic-vs-pessimistic-locking
Implementation Files:
- internal/inventory/model.go
- internal/inventory/store.go
- internal/inventory/service.go
- cmd/demo/main.go
Tests: tests/locking_test.go
Executable/Demo: cmd/demo/main.go
Approved Research Inputs: README.md, engineering/01-design.md, engineering/02-implementation-notes.md (for context, but audit focuses on code and tests)
Main Claims To Verify:
1. Naive read-modify-write (NaiveDeduct) produces lost update anomaly under concurrency.
2. Pessimistic locking (PessimisticDeduct) prevents lost update by holding row lock during read-modify-write.
3. Optimistic locking (OptimisticDeduct) detects conflicts via version check and returns error; with retry it eventually succeeds.
4. Atomic single-statement update (AtomicDeduct) prevents lost update without locks by checking stock >= qty in same statement.
5. All implementations handle error cases (insufficient stock, invalid quantity, and not found correctly.
Commands To Run:
- go test ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- The store uses in-memory mutexes to simulate database row locks; must ensure locking granularity is correct.
- Optimistic locking retry logic may have livelock or starvation under high contention (but bounded by maxRetries).
- Atomic operation in simulation is just a mutex-protected check-and-update; must verify it matches the claimed SQL behavior.
- Test for naive lost update may be flaky if timing doesn't cause overlap; but test uses sleep to increase window.
- Demo output must match actual execution (not hardcoded).