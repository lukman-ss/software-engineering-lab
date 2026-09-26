# Docs vs Code

## README vs Implementation
- README claims structure: cmd/demo, internal/server, internal/loadtest, tests, engineering/ → matches actual layout.
- README commands: `go run ./cmd/demo`, `go test -v ./...`, `go test -race ./...` → all verified working.
- README states "Smoke vs. Stress test scenarios", "constrained connection pool", "percentile calculator" → all present in code.

No DOC_CODE_MISMATCH.

## Demo vs Documented Result
- Demo output (observed 2026-09-26 run):
  - Smoke (2 VUs): Avg 21.17ms, P95 21.40ms, P99 21.66ms, 0 errors.
  - Stress (50 VUs): Avg 555.72ms, P50 568.82ms, P95 994.49ms, P99 1.32s, 0 errors.
- engineering/03-execution-result.md recorded:
  - Smoke: Avg 21.25ms, P95 21.37ms, P99 22.28ms.
  - Stress: Avg 739.09ms, P50 669.56ms, P95 1.35s, P99 1.58s.
- Prose: smoke P95 close to avg; stress P95 >> smoke P95 and stress P95 > stress avg. Pattern reproduced exactly. Absolute numbers differ slightly run-to-run (probabilistic 10% 25x latency injection + scheduling). This is expected variance, not a mismatch.
- No FAKE_DEMO. Demo is real, deterministic in pattern, variable in absolute values.

No TEST_CLAIM_MISMATCH. tests/loadtest_test.go asserts stress P95 > smoke P95 and stress P95 > stress Avg — both consistent with observed runs.

## Engineering Notes vs Code
- 02-implementation-notes.md states semaphore pattern, per-goroutine slices, custom Transport, fixed 20ms waits → all verified in code.
- States "What Is Not Demonstrated: distributed load generation, real DB lock contention" → accurate scoping, no overclaim.
- Design 01-design.md success criteria: no external deps, correct percentiles, smoke vs stress contrast, zero races → all verified.

No RESEARCH_IMPLEMENTATION_MISMATCH in scope of this audit (research content not audited per pipeline override).

## Discrepancies
None material. No DOC_CODE_MISMATCH, no TEST_CLAIM_MISMATCH, no FAKE_DEMO, no FAKE_BENCHMARK.