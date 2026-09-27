# Documentation vs Code Verification

Target Lab: labs/23-optimistic-vs-pessimistic-locking

## Comparison Matrix

| Component | README / Docs Claim | Code Implementation | Status |
|---|---|---|---|
| Directory Structure | Lists `cmd/demo/main.go`, `internal/inventory/`, `tests/locking_test.go`, `engineering/`, `go.mod`, `README.md` | Structure matches README exactly | MATCH |
| How to Run Commands | `go test -v ./...`<br>`go test -race ./...`<br>`go run ./cmd/demo` | Executed all commands; all compile and pass as specified | MATCH |
| Concurrency Strategy 1 | Naive Read-Modify-Write lost update | `NaiveDeduct` in `store.go` creates lost update | MATCH |
| Concurrency Strategy 2 | Pessimistic Locking (`SELECT ... FOR UPDATE`) | `PessimisticDeduct` uses per-row `sync.Mutex` | MATCH |
| Concurrency Strategy 3 | Optimistic Locking (Version check in `WHERE`) | `OptimisticDeduct` checks `p.Version` guard | MATCH |
| Concurrency Strategy 4 | Jittered Exponential Backoff Retry | `DeductOptimisticWithRetry` in `service.go` | MATCH |
| Concurrency Strategy 5 | Atomic Single-Statement Operations | `AtomicDeduct` in `store.go` | MATCH |

## Discrepancy Log

No discrepancies found between README claims, design notes, test suite, and actual code implementation.
