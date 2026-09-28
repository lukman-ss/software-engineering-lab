# Documentation vs Code Audit

## DOC_CODE_MISMATCH
Location: `engineering/01-design.md:66` (Test Strategy list) vs `tests/replication_test.go`
Claim: Design doc lists a test named `TestReadYourOwnWrites_LSNToken`.
Code: Actual test is named `TestReadWithToken_LSN`.
Assessment: DOC_CODE_MISMATCH
Severity: LOW
Notes: Minor naming mismatch; test covers same functionality.

## DOC_CODE_MISMATCH
Location: `README.md:30-33` (Run Demo command) vs `cmd/demo/main.go`
Claim: `go run ./cmd/demo` is the demo command.
Code: `cmd/demo/main.go` exists and runs successfully.
Assessment: PASS (no mismatch)
Severity: NONE
Notes: Matches exactly.

## DOC_CODE_MISMATCH
Location: `engineering/03-execution-result.md:52-73` (Demo output) vs actual `go run ./cmd/demo` run
Claim: Demo output lines: `[Write Primary] Wrote key 'user:profile:123', Primary LSN: 1`, `[Naive Read] Node: replica-2, Applied LSN: 0, Val: '', Err: key not found`, etc.
Observed: Actual run produced:
```
[Write Primary] Wrote key 'user:profile:123', Primary LSN: 1
[Naive Read] Node: replica-2, Applied LSN: 0, Val: '', Err: key not found
[Sticky Read (within 500ms)] Node: primary, LSN: 1, Val: '{"name":"Alice","tier":"premium"}', Err: <nil>
...
[Sync Write] LSN: 1, Write Duration: 101.222958ms
[Sync Naive Read] Node: replica-2, LSN: 1, Val: '{"name":"Bob"}', Err: <nil>
```
Assessment: PASS (output matches verbatim except for duration which varies realistically)
Severity: NONE
Notes: Demo is real; recorded output corresponds to an actual run.

## RESEARCH_IMPLEMENTATION_MISMATCH
(Per pipeline override: Do not audit research/content in this stage. Skip.)
Claim: Not evaluated.
Assessment: NOT_APPLICABLE
Severity: NONE
Notes: Pipeline specifies audit implementation and tests only; research audit is separate.

## DOC_CODE_MISMATCH
Location: `README.md:12-16` (Concepts Implemented) vs internal implementation
Claim: Implements:
- Asynchronous vs Synchronous Replication (`remote_apply`)
- Stale Read Detection & Mitigation
- Time-Based Sticky Routing
- Causal Token / Minimum LSN Routing
- Lag-Aware Dynamic Routing & Primary Fallback
Code: Each concept is implemented:
- Async/Sync: `cluster.ReplicationMode`, `Write` branching.
- Stale Read: demonstrated by `TestNaiveReplicationLag_StaleRead` and demo step 1.
- Sticky Routing: `ReadWithStickySession`.
- Causal Token: `ReadWithToken`.
- Lag-Aware: `ReadLagAware`, `TestReplicaLagThreshold_Fallback`.
Assessment: PASS
Severity: NONE
Notes: README claims are all substantiated by code and tests.

## DOC_CODE_MISMATCH
Location: `README.md:21-29` (Running the Lab) vs actual commands
Claim: `go test -v ./...`, `go test -race ./...`, `go run ./cmd/demo`.
Code: All three commands succeed (see audit execution).
Assessment: PASS
Severity: NONE
Notes: Commands documented and verified.

## Summary
- DOC_CODE_MISMATCH: one minor naming inconsistency in design doc vs test name (LOW).
- All other documentation claims match implementation and execution.
- No evidence of fake benchmark, fake demo, or fabricated results.
- Research audit artifacts exist but are out of scope per pipeline override.