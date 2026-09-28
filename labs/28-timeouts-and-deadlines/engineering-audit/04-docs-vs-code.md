# Docs vs Code Comparison

Target Lab: `labs/28-timeouts-and-deadlines`

## Documentation Review

### `README.md`
- Claims component packages: `internal/deadline`, `internal/retry`, `internal/circuit`, `internal/idempotency`. Verified: All 4 packages exist and implement claimed features.
- Claims execution commands `go test ./...`, `go test -race ./...`, `go run ./cmd/demo`. Verified: All 3 commands run without errors or warnings.

### `engineering/01-design.md` & `engineering/02-implementation-notes.md`
- Design specifies full jitter exponential backoff formula: `sleep = random_between(0, min(MaxBackoff, BaseBackoff * 2^attempt))`. Code (`internal/retry/retry.go:46-47`) matches formula exactly.
- Design specifies state machine (`CLOSED` -> `OPEN` -> `HALF_OPEN` -> `CLOSED`). Code (`internal/circuit/circuit.go`) implements exact state model and thread-safe transitions.
- Design specifies thread-safe idempotency key storage with TTL. Code (`internal/idempotency/idempotency.go`) matches.

### `cmd/demo/main.go` Output Verification
- Claimed Demo Output in execution report vs actual `go run ./cmd/demo`:
  - Demo 1: Context deadline propagation result: `context deadline exceeded` (MATCH)
  - Demo 2: Exponential backoff with full jitter attempts 1..3 success (MATCH)
  - Demo 3: Circuit Breaker state transitions (`CLOSED` -> `OPEN` -> fast rejection -> `HALF_OPEN` -> `CLOSED`) (MATCH)
  - Demo 4: Idempotence protection first charge vs deduplicated retry (MATCH)

## Assessment
- DOC_CODE_MISMATCH: NONE
- TEST_CLAIM_MISMATCH: NONE
- RESEARCH_IMPLEMENTATION_MISMATCH: NONE
- FAKE_DEMO: NONE
- FAKE_BENCHMARK: NONE

Documentation and execution results match the codebase 100%.
