# Docs vs Code Analysis

Target Lab: labs/23-optimistic-vs-pessimistic-locking

## Comparison Matrix

| Item / Claim | Location in Docs / README | Location in Code / Execution | Alignment Status | Notes |
|---|---|---|---|---|
| Strategy 1: Pessimistic Locking (`SELECT ... FOR UPDATE`) | README.md:4, 01-design.md:8 | `store.go:93`, `service.go:20`, `locking_test.go:40` | MATCH | Exact row-locking behavior verified |
| Strategy 2: Optimistic Locking (Version check + retry) | README.md:5, 01-design.md:9 | `store.go:120`, `service.go:28`, `locking_test.go:85` | MATCH | Version comparison + jittered backoff verified |
| Strategy 3: Atomic Single-Statement Operations | README.md:6, 01-design.md:10 | `store.go:155`, `service.go:47`, `locking_test.go:147` | MATCH | Single-statement conditional decrement verified |
| Concurrency lost update demonstration | README.md:3, 01-design.md:13 | `store.go:65`, `service.go:16`, `locking_test.go:10` | MATCH | Unsynchronized read-modify-write lost update verified |
| Project commands: `go test -v ./...`, `go test -race ./...`, `go run ./cmd/demo` | README.md:37-49 | `locking_test.go`, `cmd/demo/main.go` | MATCH | All commands execute clean without errors |

## Mismatch Audits

- **DOC_CODE_MISMATCH**: None found.
- **TEST_CLAIM_MISMATCH**: None found.
- **RESEARCH_IMPLEMENTATION_MISMATCH**: None found. Research report (`research/05-report.md`) defines the 3 strategies and lost update anomaly; code implements them directly.
