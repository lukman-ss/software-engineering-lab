# Docs vs Code

## README.md
- Claims 4 components: `internal/deadline`, `internal/retry`, `internal/circuit`, `internal/idempotency`
- ✅ All 4 exist with these paths
- Claims tests: `go test ./...`, `go test -race ./...`
- ✅ Both commands work and pass
- Claims demo: `go run ./cmd/demo`
- ✅ Demo runs successfully, output matches 03-execution-result.md

## engineering/01-design.md
- Concept: Context deadline propagation, backoff+jitter, circuit breaking, idempotency ✓ All implemented
- Arch: 5 components match code paths
- Components section: `Deadline`, `Retrier`, `CircuitBreaker`, `IdempotencyStore` — match exported APIs (exact names: `ExecuteWithBudget`, `Retrier`, `Breaker`, `Store`; minor doc-code naming variance, acceptable)
- Success criteria: ctx abort ✓, race pass ✓, tests ✓ — all verified
- Impl decisions: stdlib-only ✓ (imports: context/sync/time/math/rand/v2), no external deps ✓
- Decision: `pkg`-level tradeoffs (no Redis/DB) disclosed ✓

## engineering/02-implementation-notes.md
- Files added: All listed files exist
- Jitter algorithm: `sleep=rand.Float64()*min(maxBackoff, base*2^(attempt-1))` — matches code ✓
- Context propagation: "Sub-calls inherit parent context, ctx timeout overrides local budget" — matches `ExecuteWithBudget` (childCtx inherit, timeout = min(parent, budget)) ✓
- Circuit: "Transitions into Open when reaching consecutive failures" ✓
- Idempotency: "Memory store with TTL mapping request tokens to responses" ✓
- Known limitations: "In-memory idempotency does not persist across crashes" ✓ (disclosed)
- Tradeoffs: stdlib-only ✓
- Not demonstrated: "Distributed coordination across processes/network (e.g. gRPC headers)" — disclosed ✓

## engineering/03-execution-result.md
- Claims:
  - Tests: all packages `ok` (timings 0.3-0.5s)
  - Race: all packages `ok` (timings 1.1-1.5s)
  - Demo: 4-section output with specific values
- ✅ Audit re-ran all 3 commands. Results match exactly:
  - `go test ./...` → all `ok` (same package list, no test files for cmd/demo)
  - `go test -race ./...` → all `ok`
  - `go run ./cmd/demo` → identical output (DeadlineExceeded, attempts 1-3, CLOSED/OPEN/HALF_OPEN, DEDUPLICATED)
- No FAKE_BENCHMARK, no UNVERIFIED_RESULT

## Issues Found
- None (all docs match code/tests/demo)
- Minor: No mismatch. Commit to record.
