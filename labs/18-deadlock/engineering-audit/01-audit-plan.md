# Engineering Audit Plan

Target Lab: labs/18-deadlock
Implementation Files: `internal/bank/account.go`, `internal/transfer/transfer.go`
Tests: `tests/transfer_test.go`
Executable/Demo: `cmd/demo/main.go`
Approved Research Inputs:
- Deadlock occurs via circular wait condition (Coffman conditions).
- Systems resolve deadlocks by aborting one transaction as victim.
- Lock ordering prevents circular wait and avoids deadlocks entirely.
- Transaction duration correlates directly with deadlock probability.
- Application-level retries recover aborted transactions.
Main Claims To Verify:
- Concurrent bidirectional naive transfers trigger deadlocks and return `ErrDeadlock`.
- Lock-ordered concurrent bidirectional transfers complete successfully without deadlocks.
- Application retry recovers from aborted deadlock attempts.
- Longer transaction duration increases deadlock occurrence compared to zero delay.
- Concurrency is safe and race detector passes without data races.
Commands To Run:
- `go test ./...`
- `go test -race ./...`
- `go run ./cmd/demo`
Primary Risks:
- Simulation may fail to produce deadlocks reliably under fast schedulers if timeouts/delays are misconfigured.
- Flaky test execution if concurrency timings are too tight.
- Race conditions during shared balance updates if locks are improperly managed.
