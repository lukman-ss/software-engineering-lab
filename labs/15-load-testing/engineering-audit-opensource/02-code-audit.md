# Code Audit

Scope: implementation + tests only (research/content out of scope per pipeline override).

## Finding 1

Location: `internal/loadtest/runner.go:44-115`
Claimed Behavior: N VUs run concurrently; results aggregated without lock contention; run bounded by Duration.
Observed Implementation: Per-VU `vuResult` slots indexed by goroutine ID, merged after `wg.Wait` (happens-before safe). `context.WithTimeout(ctx, cfg.Duration)` bounds execution. Body drained via `io.Copy(io.Discard)` + `Close` for conn reuse. Custom transport `MaxIdleConnsPerHost: 1000` avoids client-side pooling bottleneck.
Assessment: PASS
Severity: LOW
Notes: `go test -race` clean across two full runs. No shared mutable state between VUs.

## Finding 2

Location: `internal/loadtest/runner.go:84-98`
Claimed Behavior: Transport failures and HTTP >= 400 counted as errors; in-flight cancellations at deadline excluded.
Observed Implementation: `client.Do` error counted only if `ctx.Err() == nil`. Status >= 400 increments `errs`, success latencies only recorded. `NewRequest` error path same guard.
Assessment: PASS
Severity: LOW
Notes: Proven by `TestLoadTest_ErrorCount` (500s all errors) and `TestLoadTest_DialError` (refused conn all errors). Excluding deadline-canceled in-flight requests from totals is reasonable; means Total < attempts at boundary.

## Finding 3

Location: `internal/loadtest/metrics.go:23-73`
Claimed Behavior: Accurate Min/Max/Avg/P50/P90/P95/P99.
Observed Implementation: Defensive empty (`Result{}`) and all-errors (counts+RPS, zero latencies) paths. Copy-sort, then min/max/avg via integer division. `percentile` = floor nearest-rank `sorted[int((n-1)*pct/100)]`.
Assessment: PASS
Severity: LOW
Notes: Hand-verified: 1..100ms gives P50=50ms (idx 49), P95=95ms (idx 94), P99=99ms (idx 98); unit test asserts exactly these. Floor-index slightly underestimates vs linear interpolation — disclosed approximation class, self-consistent, fine at lab scale. `ponytail:` ceiling comment present.

## Finding 4

Location: `internal/server/server.go:51-89`
Claimed Behavior: Bounded connection pool causes queueing + tail degradation under stress.
Observed Implementation: Buffered-channel semaphore cap `MaxDBConnections`, ctx-aware acquire/release, `time.Timer` (not bare Sleep) with ctx abort, atomic `activeReq` tracking. Overload branch: `activeReq > MaxDB` triggers 25x duration at 10% probability.
Assessment: PASS
Severity: LOW
Notes: Live demo reproduced claim: stress P95 772ms vs smoke P95 21.5ms (~36x). Timer+`Stop` correct, no leak. Global `math/rand` safe for concurrent use (locked source), auto-seeded on Go 1.22 — nondeterministic values, deterministic mechanism.

## Finding 5

Location: `internal/loadtest/runner.go:26-29`
Claimed Behavior: (implicit) Config always usable.
Observed Implementation: Only `VUs <= 0` defaulted to 1. Empty URL → `NewRequest` errors every spin until deadline (hot loop counting errors). Zero/negative Duration → instant expiry, empty `Result{}`. Empty Method → GET silently.
Assessment: WARNING
Severity: LOW
Notes: No validation tests exist for these inputs. Lab-scope demo never hits them. Harden with URL/method check returning zero-Result or error — optional, non-blocking.

## Finding 6

Location: `internal/server/server.go:85`
Claimed Behavior: 201 JSON booking response.
Observed Implementation: `_ = json.NewEncoder(w).Encode(...)` — encode error discarded.
Assessment: WARNING
Severity: LOW
Notes: Encoding a static struct to `ResponseWriter` fails only if client went away, in which case nothing actionable. Benign; noted for completeness.

## Finding 7

Location: `cmd/demo/main.go:1-72`
Claimed Behavior: Smoke (2 VUs) vs stress (50 VUs) contrast against 5-conn pool, 20ms query, 2s each.
Observed Implementation: Matches claim exactly. `httptest.NewServer`, tabwriter output of Total/Success/Errors/RPS/Avg/P50/P95/P99.
Assessment: PASS
Severity: LOW
Notes: Real execution 2026-09-26: smoke 186 req / 0 err / P95 21.5ms; stress 178 req / 0 err / P95 772ms / P99 1.13s. Shape matches recorded `engineering/03-execution-result.md` (values vary, invariant holds). No fabrication.
