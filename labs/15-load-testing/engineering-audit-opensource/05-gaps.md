# Gaps

| ID | Gap Type | Location | Severity | Summary |
|----|----------|----------|----------|---------|
| G1 | DOC_CODE_MISMATCH | engineering/03-execution-result.md test section | LOW | Execution-result doc lists only `TestCalculateMetrics/TestCalculateMetrics_Empty/TestCalculateMetrics_Invariants` and `TestLoadTest_SmokeVsStress/TestLoadTest_ErrorCount/TestServer_MethodNotAllowed/TestLoadTest_DialError/TestServer_ContextCanceled`, omitting the four tests added later: `TestCalculateMetrics_SingleSample`, `TestCalculateMetrics_RPS`, `TestServer_MaxDBConnectionsBound`, `TestLoadTest_SuccessAndErrorInvariant`. Actual `go test -v` shows all 9 passing. |
| G2 | DOC_CODE_MISMATCH | engineering/02-implementation-notes.md (Implementation-Specific Choices) | LOW | Note claims "Tail latency is strictly a function of queuing time when VUs exceed max DB connections." Code (server.go:74-78) adds a 10% chance of `DBQueryDuration * 25` slowdown whenever `activeReq > MaxDBConnections`, independent of semaphore queuing. Tail latency is queuing OR random slowdown, not strictly queuing. |
| G3 | DOC_CODE_MISMATCH | cmd/demo/main.go printResults | LOW | `Result` includes `P90Latency` but the demo prints only P50/P95/P99 (README/research Finding 2 lists P50/P95/P99 prominently). Engineering 01-design success criterion says generator computes P50/P90/P95/P99; all are computed in code but P90 is not displayed. |
| G4 | MISSING_EDGE_CASE | internal/loadtest/metrics.go percentile | LOW | `percentile()` uses lower-interpolation `idx = int((len-1)*pct/100)`. It is deterministic and tested for ascending series + invariants, but the definition is not documented as the "lower" variant, nor tested against out-of-order/random input. No functional correctness gap for the lab's scale. |
| G5 | UNHANDLED_ERROR (semantic) | internal/loadtest/runner.go:91-98 | MEDIUM | Latency is recorded only for successful (2xx) responses. Transport errors and HTTP >= 400 responses are not timed, so P50/P95/P99 under represent response time for error-heavy workloads (research Finding 2 expects percentile response times for all requests). In the demo all responses are 201, so Smoke/Stress percentiles are unaffected; gap only matters when errors carry latency. Explicitly acknowledged below. |
| G6 | IMPLEMENTATION_OVERCLAIM | engineering/01-design.md (line 16) | LOW | Design states "Average and P95 are close" for smoke load. Real smoke runs show Avg ~21.3ms and P95 ~22.2ms — close in this case, but the claim has no margin. Cosmetic; not actionable as a defect. |
| G7 | RACE_CONDITION risk (low likelihood) | tests/loadtest_test.go:142-172 polling loop | LOW | `TestServer_MaxDBConnectionsBound` polls `srv.ActiveConnections()` from a goroutine at 1ms ticker while the channel is being mutated. `len(channel)` is not a documented-atomic operation, yet it is widely treated as safe in practice; it relies on Go runtime internals. `-race` did not flag it, but a stricter auditor would treat `len()` on a mutating channel as a latent race. No failure observed. |
| G8 | MISSING_TEST | internal/server/server.go | LOW | Server package has no dedicated unit tests (only tested via integration in tests/). `go test` shows `internal/server [no test files]`. Adequately covered via integration, but server-only logic (tail-slowdown branch) is untested in isolation. |
| G9 | RACE_CONDITION (mitigated) | internal/server/server.go:5 `math/rand` | LOW | Uses global `math/rand` which is concurrency-safe (guarded internally) but the recommended API is `math/rand/v2` or a local source. Not a defect; no race. |
| G10 | MISSING_TEST | runner.go VU cancellation loop | LOW | The load generator's `select { case <-ctx.Done(): ... default: }` cancellation path is not directly unit-tested (only server-side cancellation via `TestServer_ContextCanceled`). Functional under -race. |

## Non-blocking, accepted
G4, G6, G7, G8, G9, G10.

## Requires attention (not failing, but documented)
G1, G2, G3 (doc mismatches - update the stale docs), G5 (latency-excludes-errors — a design scoping decision; acceptable for this lab where all responses succeed).

## No fabricated results
Demo output in engineering/03-execution-result.md is a documented "sample run" and explicitly states "actual values vary per run." My reproduced run produced different numeric values (186 vs 188 requests, etc.) but the SAME structural invariant (Stress P95/P99 >> Smoke). No fake benchmark/result. FAKE_DEMO/FAKE_BENCHMARK = none.

## No research mismatch requiring code change
Research provides general industry context (test types, percentile metrics, bottleneck isolation, pitfalls, SDLC timing). The implementation realizes the applicable subset (smoke vs stress, percentile metrics, DB-pool saturation isolation) using stdlib only, consistent with engineering design decision "Automated benchmarks and tests execute without external dependencies." RESEARCH_IMPLEMENTATION_MISMATCH = none.