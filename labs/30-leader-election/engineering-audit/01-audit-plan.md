# Engineering Audit Plan

Target Lab: `labs/30-leader-election`
Implementation Files:
- `internal/coordinator/coordinator.go`
- `internal/candidate/candidate.go`
- `internal/storage/storage.go`
- `cmd/demo/main.go`
Tests:
- `tests/election_test.go`
Executable/Demo:
- `cmd/demo/main.go`
Approved Research Inputs:
- `research/01-plan.md`
- `research/02-sources.md`
- `research/03-mechanisms-and-leases.md`
- `research/04-fencing-and-redlock.md`
- `research/05-system-comparison-and-best-practices.md`
- `research-revision/03-revision-result.md`
Main Claims To Verify:
1. Lease-based leader election acquires exclusive lock with TTL and issues strictly monotonically increasing fencing tokens.
2. Heartbeat keeps lease alive; failure/pause longer than TTL allows standby candidate promotion with higher fencing token.
3. Fenced storage rejects stale tokens (`token <= lastSeenToken`) with explicit error, preventing split-brain writes.
4. Concurrency safety across concurrent candidate campaigns under race detector.
5. README and engineering notes match actual code and execution behaviors.
Commands To Run:
- `go test -v -count=1 ./...`
- `go test -race -count=1 ./...`
- `go run ./cmd/demo`
Primary Risks:
- Clock/timer race conditions in lease expiration logic.
- False positive test assertions that do not strictly check stale rejection.
- Discrepancies between logged demo outputs and actual execution.
