# Docs vs Code Audit

Lab: labs/15-load-testing

## Sources Compared
- README.md (project overview and instructions)
- engineering/01-design.md (design doc)
- engineering/02-implementation-notes.md (implementation notes)
- Code: internal/server/server.go, internal/loadtest/*.go, cmd/demo/main.go
- Tests: internal/loadtest/metrics_test.go, tests/loadtest_test.go
- Demo output: actual `go run ./cmd/demo` runs (2 samples)

## Methodology
Each finding records the document location, the claimed behavior, the observed implementation/demo/test, and an assessment. Severity levels per GAP types: LOW, MEDIUM, HIGH, CRITICAL. Only discrepancies between documentation and implementation/test/demo are recorded. Per pipeline override, research/ and content/ are not audited.

---

## Finding 1
Location: engineering/01-design.md:32 (Component #2 LoadTester description)
Claim: "LoadTester: Concurrency orchestrator generating HTTP traffic with specified virtual users (VUs) and iterations."
Observed: The `loadtest.Config` struct has `VUs` and `Duration` fields; no "iterations" field exists. The runner executes requests for a duration, not a fixed iteration count. The demo uses `Duration: testDuration`.
Assessment: DOC_CODE_MISMATCH (terminology only)
Severity: LOW
Notes: The term "iterations" appears to be a documentation artifact; the implementation and demo are time-based. This does not affect correctness or the demonstration of load vs stress.

## Finding 2
Location: engineering/01-design.md:33 (Component #3 MetricsAggregator description)
Claim: "MetricsAggregator: Thread-safe latency collector sorting durations to derive accurate percentiles."
Observed: There is no `MetricsAggregator` struct. Latency collection is performed per-VU in `runner.go` (each goroutine appends to a private slice). After all VUs finish, slices are concatenated and passed to `CalculateMetrics` (a pure function) which sorts and computes percentiles. The system is thread-safe due to lack of shared mutable state during the run, but there is no active "collector" with internal locking.
Assessment: DOC_CODE_MISMATCH (terminology/abstraction level)
Severity: LOW
Notes: The description slightly overstates the structure; the function `CalculateMetrics` correctly sorts and computes percentiles as claimed. The ponytail in metrics.go confirms the exact-sort approach.

## Finding 3
Location: engineering/02-implementation-notes.md:16 (Implementation-Specific Choices)
Claim: "Wait durations in server mock are fixed (20ms), making the tail latency strictly a function of queuing time when VUs exceed max DB connections."
Observed: `server.New(cfg)` defaults `DBQueryDuration` to 10 * time.Millisecond if `cfg.DBQueryDuration <= 0`. The demo in `cmd/demo/main.go` sets `DBQueryDuration: 20 * time.Millisecond`. Thus the wait duration is configurable, not fixed at 20ms in the server implementation; it is fixed only for the demo run.
Assessment: DOC_CODE_MISMATCH (contextual accuracy)
Severity: LOW
Notes: The claim holds true for the demo as executed, but the server implementation itself allows configuration. The wording in the notes is accurate when read in the context of the demo, but could be misinterpreted as a server-wide constant.

## Finding 4
Location: engineering/01-design.md:22 (Expected Behavior)
Claim: "**Smoke Load (low VUs)**: All requests process within normal latency limits. Average and P95 are close. Error rate is 0%."
Observed: Demo output (2 sample runs):
- Run 1: Smoke Avg 21.25ms, P50 21.21ms, P95 21.37ms, P99 22.28ms, Errors 0.
- Run 2: Smoke Avg 21.14ms, P50 21.21ms, P95 21.55ms, P99 21.77ms, Errors 0.
Avg and P50/P95/P99 are within ~1-2ms of each other (jitter from scheduling/GC). Error rate 0%. Matches claim.
Assessment: PASS
Severity: N/A (match)

## Finding 5
Location: engineering/01-design.md:22 (Expected Behavior)
Claim: "**Stress Load (high VUs)**: Concurrent requests exceed server resource capacity (connection pool limit). Tail requests queue up. P95 and P99 latency spikes significantly, while average latency degrades less severely, proving the masking effect of averages."
Observed: Demo output:
- Run 1: Stress Avg 739ms, P50 669ms, P95 1.35s, P99 1.58s (P95/P99 >> Avg/P50).
- Run 2: Stress Avg 470ms, P50 586ms, P95 774ms, P99 1.16s (P95/P99 > Avg/P50).
In both runs, P95 and P99 significantly exceed the average (and smoke latencies), demonstrating tail spikes and the averaging masking effect. Errors 0. Matches claim.
Assessment: PASS
Severity: N/A (match)

## Finding 6
Location: engineering/01-design.md:23 (Failure Scenario)
Claim: "Under excessive concurrent load, queuing behind a constrained resource (connection pool) causes high latency and severe tail degradation for the 95th and 99th percentiles."
Observed: Same as Finding 5; the latency spikes are directly attributable to queuing on the semaphore (max 5 connections). The design's mechanism (semaphore) is verified in code.
Assessment: PASS
Severity: N/A (match)

## Finding 7
Location: engineering/01-design.md:24-25 (Success Criteria)
Claim: "Automated benchmarks and tests execute without external dependencies."
Observed: All tests use `httptest.NewServer`; no external services required. `go test -v ./...` succeeds offline.
Assessment: PASS
Severity: N/A (match)

## Finding 8
Location: engineering/01-design.md:24-25 (Success Criteria)
Claim: "Load generator computes Min, Max, Average, P50, P90, P95, and P99 latencies accurately."
Observed: `TestCalculateMetrics` validates exact values for a known sequence. Demo output shows plausible Min/Max/Avg/Pxx relationships (Min <= P50 <= P95 <= P99 <= Max) with rough alignment to the 20ms base (e.g., smoke Min~21ms, stress Min~470ms due to queuing).
Assessment: PASS
Severity: N/A (match)

## Finding 9
Location: engineering/01-design.md:24-25 (Success Criteria)
Claim: "Demonstration clearly contrasts Smoke test metrics against Stress test metrics."
Observed: Demo prints two tables with clear labels; stress latencies are 20x-50x higher than smoke for tail percentiles.
Assessment: PASS
Severity: N/A (match)

## Finding 10
Location: engineering/01-design.md:24-25 (Success Criteria)
Claim: "All tests pass with zero race conditions (`go test -race ./...`)."
Observed: Verified by actual execution: `go test -race ./...` reports zero races for all packages.
Assessment: PASS
Severity: N/A (match)

## Summary

Three LOW-severity DOC_CODE_MISMATCH findings were identified, all relating to minor terminology or abstraction level discrepancies in the design and implementation notes. No TEST_CLAIM_MISMATCH or RESEARCH_IMPLEMENTATION_MISMATCH (per override) was found. The documentation accurately captures the core behavior, structure, and success criteria; the code and demo fulfill all claims.