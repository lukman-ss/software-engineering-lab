# Engineering Audit Plan

Target Lab: labs/18-deadlock
Implementation Files: internal/bank/account.go, internal/transfer/transfer.go
Tests: tests/transfer_test.go
Executable/Demo: cmd/demo/main.go
Approved Research Inputs: research/runs/2026-09-25-deadlock/05-report.md
Main Claims To Verify:
1. Deadlock causes circular wait (Finding 1)
2. Deadlock victim aborts (Finding 2)
3. Lock ordering prevents deadlock (Finding 3)
4. Long transactions increase deadlock probability (Finding 4)
5. Retry mechanism recovers deadlocks (Finding 5)
Commands To Run:
1. go test ./...
2. go test -race ./...
3. go run ./cmd/demo
Primary Risks:
1. Inaccurate simulation of database locking
2. Deadlock not deterministically triggered
3. Race conditions in test/demo orchestration
