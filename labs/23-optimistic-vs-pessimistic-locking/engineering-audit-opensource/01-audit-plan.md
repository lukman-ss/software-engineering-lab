# Engineering Audit Plan

Target Lab: labs/23-optimistic-vs-pessimistic-locking
Implementation Files:
- internal/inventory/model.go
- internal/inventory/store.go
- internal/inventory/service.go
- cmd/demo/main.go
Tests:
- tests/locking_test.go
Executable/Demo: cmd/demo/main.go
Approved Research Inputs: engineering/01-design.md (not audited in this stage per pipeline override)
Main Claims To Verify:
1. Naive read-modify-write demonstrates lost updates under concurrent goroutines.
2. Pessimistic locking prevents concurrency conflicts entirely by blocking writers on a per-row lock.
3. Optimistic locking detects version mismatch and safely rejects stale writes without data corruption; application retry loop eventually converges successfully.
4. Atomic single-statement updates serialize safely at the database statement level without manual transaction lock blocks.
Commands To Run:
- go test ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Race conditions in Store implementation
- Incorrect locking behavior in pessimistic/optimistic/atomic methods
- Tests not properly validating claimed behavior (false positives/negatives)
- Demo output not matching actual implementation behavior
- Missing edge case handling (e.g., zero/negative quantities, non-existent products)