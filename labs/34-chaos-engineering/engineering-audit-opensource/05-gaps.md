# Gaps

## GAP-1
Type: DOC_CODE_MISMATCH
Location: engineering/01-design.md vs code layout
Description: Design references `pkg/fault`, `pkg/monitor`, `pkg/circuitbreaker`, `pkg/experiment`. Code uses `internal/...`.
Severity: LOW
Status: Non-blocking. README structure section correctly lists `internal/...`.

## GAP-2
Type: IMPLEMENTATION_OVERCLAIM
Location: engineering/01-design.md Test Strategy vs internal/fault/injector.go
Description: Design claims probabilistic triggering. Injector has unused `errorRate` field, deterministic `forceError` only.
Severity: LOW
Status: Non-blocking. README does not claim probabilistic injection. Implementation notes accurately describe deterministic behavior.

## GAP-3
Type: MISSING_TEST
Location: tests/chaos_test.go vs internal/experiment/runner.go
Description: No explicit test for natural COMPLETED path (timeout without breach) or context-cancel abort. Auto-abort path tested, completion shown only in demo.
Severity: LOW
Status: Non-blocking. Core abort behavior proven.

## GAP-4
Type: MISSING_EDGE_CASE
Location: internal/experiment/runner.go Run()
Description: Concurrent Run() calls on same instance unguarded. Single-call usage in tests/demo safe.
Severity: LOW
Status: Non-blocking. Scoped to documented usage.

No HIGH or CRITICAL gaps. No FAKE_DEMO, FAKE_BENCHMARK, BROKEN_IMPLEMENTATION, RACE_CONDITION, UNHANDLED_ERROR.
