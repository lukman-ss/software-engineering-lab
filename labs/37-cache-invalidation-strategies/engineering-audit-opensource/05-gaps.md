# Gaps

## Gap 1 — MISSING_TEST (MEDIUM)
Service Get/Update error & fallback paths covered in code but not tested:
XFetch stale-fallback-on-recompute-error (stampede.go:159-164), SWR sync-fetch on
beyond-stale-window error, SingleFlight flight-func error propagation, CacheAside/Update
DB-write failure, WriteBehind lost async Write error.
Why it matters: failure handling is otherwise unproven; a future refactor could silently
invert the `if ok` fallback. Add tests with a failing/mock-err DB injected.

## Gap 2 — MISSING_TEST (MEDIUM)
XFetch `Get` end-to-end not unit-tested. Only `ShouldRecompute` pure fn tested. Early-refresh
behavior is demo-only.
Why it matters: the formula is proven but the Get() branch wiring (GetRaw → now.After →
remaining → recompute-or-serve, plus stale fallback) is unverified by the test suite.

## Gap 3 — MISSING_EDGE_CASE / weak assertion (MEDIUM)
SWR test value-only, timing-coupled (30ms + 50ms Sleeps); does not assert RevalidateCount or
the single-in-flight-per-key guard (stampede.go:228-231); does not assert stale value held
during async refresh.
Why it matters: test can pass while revalidation logic regresses; also flake-prone under load.

## Gap 4 — MISSING_TEST (MEDIUM)
Write-Behind overflow drop and Close() drain not asserted; only post-sleep write count checked.
Why it matters: the silent-drop path (patterns.go:156-160 default branch) and graceful
shutdown drain (128-131) are documented behavior — should have an overflow test + drain test.

## Gap 5 — RACE_CONDITION (latent) (MEDIUM)
`XFetchService.SetRandFunc` (stampede.go:109-111) writes randFunc with no synchronization
while `getRand` (113-120) reads it with no lock on the randFunc path. Safe only because tests
set it before concurrent Gets; concurrent SetRandFunc would race.
Why it matters: not exercised by `-race`, but real bug if ever used concurrently. Fix: protect
randFunc read/write with s.mu, or document single-writer invariant.

## Gap 6 — MISSING_TEST (MEDIUM)
No context-cancellation or timeout test. MockDB honours ctx.Done (repo.go:37-42, 56-61) but
services only seen with context.Background(); cancellation never challenged.
Why it matters: ctx-propagation path is unproven.

## Gap 7 — DOC_CODE_MISMATCH (LOW)
engineering/01-design.md §Architecture lists `jitter.go: TTL jitter calculation`; no such file —
jitter implemented in store.go. Stale doc reference only.

## Gap 8 — UNKNOWN_RESULT / timing risk (LOW)
SWR and WriteBehind rely on real wall-clock Sleeps; documented as known limitation. No
determinism risk found in observed runs, but flakiness acknowledged. Acceptable for lab;
monitor under load.

## Summary of severity
MEDIUM: 1,2,3,4,5,6
LOW: 7,8
CRITICAL/HIGH: (none) — no fabricated result, no broken core behavior, no race in exercised
concurrency.
