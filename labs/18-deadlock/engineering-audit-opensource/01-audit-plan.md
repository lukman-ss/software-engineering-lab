# Engineering Audit Plan

Target Lab: labs/18-deadlock
Implementation Files:
- internal/bank/account.go
- internal/transfer/transfer.go
- cmd/demo/main.go
Tests: tests/transfer_test.go
Executable/Demo: cmd/demo/main.go
Approved Research Inputs: Research status is APPROVED (per engineering/01-design.md)
Main Claims To Verify:
1. Deadlock as circular wait condition (naive transfers deadlock)
2. Deadlock handling via abort (deadlock victim)
3. Deadlock prevention using Lock Ordering (ordered transfers never deadlock)
4. Transaction duration impacting deadlock probability (longer transactions increase deadlock frequency)
5. Application-level retry mechanism to recover from deadlock aborts
Commands To Run:
- go test ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Timing sensitivity in tests/demo (though values are tuned)
- Potential data races in lock/unlock operations
- Retry mechanism may livelock or starve under contention