# Gaps

| # | Type | Location | Severity | Description |
|---|------|----------|----------|-------------|
| 1 | MISSING_TEST | circuit_breaker.go | LOW | No explicit unit test asserting that two distinct generations cannot overlap in HalfOpen to exceed HalfOpenMaxCalls under concurrency. TestHalfOpenThrottlesExcessCalls covers single concurrent probe; true flood (N simultaneous) not asserted with atomic count. Stress smoke in TestInterleavedConcurrentTransitions is weak. |
| 2 | MISSING_TEST | circuit_breaker.go | LOW | HalfOpenMaxCalls >1 not explicitly tested with concurrent successful probes; only ==1 path asserted. Could hide counter overflow. |
| 3 | DOC_CODE_MISMATCH | engineering/02-implementation-notes.md | LOW | References `checkStateTransitionLocked`; actual function is `advanceLocked`. Name stale, description accurate. |
| 4 | DOC_CODE_MISMATCH | engineering/03-execution-result.md | LOW | Scenario 1 without-breaker rows show `state=CLOSED`; code CheckoutWithoutBreaker omits state (HasBreaker=false). README and live demo correct; execution-result stale. |
| 5 | DOC_CODE_MISMATCH | engineering/03-execution-result.md | LOW | Test listing stale: documents 11 unit + 1 integration tests; repo now has 16 + 2. Execution counts themselves PASS; doc list incomplete. |
| 6 | MISSING_TEST | internal/payment | LOW | payment package has no unit tests of its own (client error wrapping, status handling). Correctness relied on breaker+checkout tests. Scope-limited lab, but client.go untested directly. |
| 7 | MISSING_TEST | internal/checkout | LOW | checkout Service has no dedicated unit tests; exercised only via integration. Minor for a facade. |
| 8 | MISSING_EDGE_CASE | circuit_breaker.go | LOW | Zero-value Config with all fields unset still works (New fallbacks) but not asserted directly; DefaultConfig path covered, explicit zero-Config path not. |

No HIGH/CRITICAL gaps. No RACE_CONDITION observed (race detector clean). No UNHANDLED_ERROR, BROKEN_IMPLEMENTATION, MISSING_EDGE_CASE beyond normal lab scope, FAKE_DEMO, or FAKE_BENCHMARK. No RESEARCH_MISMATCH (research not in scope per override). No IMPLEMENTATION_OVERCLAIM (README honestly declares omitted metrics, consecutive failure counting limitation).
