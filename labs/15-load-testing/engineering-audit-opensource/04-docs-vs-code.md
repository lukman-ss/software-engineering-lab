# Docs vs Code

Compare: README, engineering notes, claimed results vs implementation and observed behavior. Implementation-only audit; not a judgment on research prose.

## D1 — README demo command

Claimed (README.md):
- `go run ./cmd/demo` runs the comparative smoke vs stress demo.
Observed: Command succeeds; produces smoke (~21ms) and stress (~210ms) tables on first and second runs (recorded execution result also shows similar values).
Status: MATCH

## D2 — README test commands

Claimed (README.md):
- `go test -v ./...`
- `go test -race ./...`
Observed: Both pass. `go test -race -count=1 -v ./...` run during this audit passed all 2 unit tests and 5 integration tests with no data races.
Status: MATCH

## D3 — Engineering design: P90 in success criteria

Claimed (`engineering/01-design.md:21`):
- "Load generator computes Min, Max, Average, P50, P90, P95, and P99 latencies accurately."
Observed: `metrics.go:18` declares `P90Latency`; `CalculateMetrics` never assigns it (always 0). `metrics_test.go` does not assert P90. Demo does not print P90.
Status: DOC_CODE_MISMATCH (P90 promised, never computed)
Severity: MEDIUM — does not invalidate core claims (smoke-vs-stress divergence, percentile math for P50/P95/P99).

## D4 — Engineering design: P99 "accuracy" claim

Claimed (`engineering/01-design.md:21`, `02-implementation-notes.md:27`):
- Percentiles computed "accurately"; P95/P99 degrade under saturation.
Observed: Nearest-rank `idx=int((len-1)*pct/100)` is deterministic and tested for P50/P95/P99 (exact for 100 samples). P99 for 100 samples → idx98 → 99ms. Correct under nearest-rank convention.
Status: MATCH

## D5 — Implementation notes: queuing claim

Claimed (`engineering/02-implementation-notes.md:27-29`):
- "Stress load queuing behind a saturation point (connection pool limit), forcing P95 to severely degrade."
Observed: Demo reproduces: smoke P95≈21ms, stress P95≈212ms (10x). Integration test asserts same ordering and passed.
Status: MATCH — real, not fabricated.

## D6 — Execution result vs re-run

Claimed (`engineering/03-execution-result.md:51-81`):
- Records `go test -v`, `go test -race`, `go run ./cmd/demo` outputs with smoke ~21ms / stress ~200ms+ latencies.
Observed (this audit): `go test -count=1` PASS; `go test -race -count=1 -v` PASS (same 7 tests); `go run ./cmd/demo` smoke ~21ms / stress ~201ms — same shape as recorded.
Status: MATCH — recorded demo output is real and reproducible, not fabricated.

## D7 — "No external dependencies" claim

Claimed (`engineering/01-design.md:50`):
- "Standard library `net/http` and `sync` used exclusively; no third-party dependencies."
Observed: Confirmed — `go.mod` has zero `require` directives; imports are stdlib only.
Status: MATCH

## D8 — Limitation disclosure

Claimed (`engineering/02-implementation-notes.md:19-21,30-33`):
- Admits: exact-sort percentiles unsuitable for million-RPS; network latency omitted; distributed generation and real DB contention not demonstrated.
Observed: Code matches — fixed 20ms timer, semaphore mock, in-memory sort.
Status: MATCH — honest scoping, no overclaim on these points.

## D9 — Test claims vs test code

Claimed (README:9): "tests: Automated integration tests validating metric accuracy and stress-induced latency growth."
Observed: `tests/loadtest_test.go` validates exactly that (smoke-vs-stress P95 ordering + error counting). `internal/loadtest/metrics_test.go` validates metric accuracy.
Status: MATCH

## Summary of mismatches

- DOC_CODE_MISMATCH: P90 promised in design, never computed or tested (D3). Only mismatch found.
- No TEST_CLAIM_MISMATCH: tests do what README claims.
- No FAKE_DEMO / FAKE_BENCHMARK: demo output verified live and physically consistent (D7 in code audit).
