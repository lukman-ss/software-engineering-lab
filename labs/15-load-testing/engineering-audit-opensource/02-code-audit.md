# Code Audit

Scope: `internal/loadtest/runner.go`, `internal/loadtest/metrics.go`, `internal/server/server.go`, `cmd/demo/main.go`. Verified against live execution (`go build`, `go test -race -count=1 -v`, `go run ./cmd/demo`).

## Finding 1 — Build & Module

Location: `go.mod:3`
Claimed Behavior: Module compiles with Go toolchain, standard library only.
Observed Implementation: `go build ./...` exits 0 with no output. No third-party dependencies. `go.mod` requires `go 1.26.7`.
Assessment: PASS
Severity: LOW
Notes: README states "Go 1.22+" as minimum; `go.mod` pins 1.26.7. Higher minor version satisfies the stated minimum. No contradiction in behavior.

## Finding 2 — Concurrency model of `Runner.Run`

Location: `internal/loadtest/runner.go:44-102`
Claimed Behavior: VUs run concurrently; results aggregated safely after completion (race-free).
Observed Implementation: Spawns `cfg.VUs` goroutines. Each goroutine writes **only** to its own `results[vuID]` slot on a pre-allocated slice (index disjoint per goroutine). Main goroutine reads slices only after `wg.Wait()`. `sync.WaitGroup.Add`/`Done`/`Wait` establish happens-before edges per the Go memory model; the race detector confirms no races.
Assessment: PASS
Severity: MEDIUM
Notes: `bytes.NewReader(r.cfg.Body)` is called per-request from a shared immutable byte slice — concurrent reads are safe.

## Finding 3 — Context cancellation / timeout handling in runner

Location: `internal/loadtest/runner.go:55-56, 65-99`
Claimed Behavior: Runner respects both a deadline (request creation with `NewRequestWithContext`) and the load-test duration (`ctx` timeout).
Observed Implementation: `context.WithTimeout(ctx, cfg.Duration)` bounds the whole run. Each request is built with the same `ctx`; on `client.Do` error, the code only increments `errs` when `ctx.Err() == nil` — i.e. deadline/cancel errors are not double-counted as request errors. On `ctx.Done()`, the goroutine stores its partial result and returns.
Assessment: PASS
Severity: MEDIUM
Notes: Correct separation of "test duration expired" from "real request failures."

## Finding 4 — Error propagation & counting

Location: `internal/loadtest/runner.go:71-74, 82-88, 92-96`
Claimed Behavior: Bad request creation, transport errors, and HTTP 5xx are all counted as errors; only successful (<400) responses contribute latencies.
Observed Implementation: `http.NewRequestWithContext` errors → `errs++` and `continue`. `client.Do` errors (excluding context-cancellation) → `errs++`. Status `>= 400` → `errs++` (no latency recorded). Otherwise latency appended. `Result` = `CalculateMetrics(allLatencies, totalErrs, totalDuration)`.
Assessment: PASS
Severity: LOW
Notes: `resp.Body` drained and closed (`io.Copy(io.Discard, …)` then `Close)` — no connection leak.

## Finding 5 — Metrics math (percentiles & aggregates)

Location: `internal/loadtest/metrics.go:23-73`
Claimed Behavior: Min/Max/Avg/P50/P90/P95/P99 computed exactly from a sorted latency slice via nearest-rank.
Observed Implementation: `total = len(latencies) + errors`; `RPS = total / duration.Seconds()`. Sorts a copy, sums, takes `sorted[0]`/`sorted[len-1]`, `Avg = sum/len`. `percentile` uses `idx = int(float64(len-1) * pct/100)` and returns `sorted[idx]`.
Assessment: PASS
Severity: MEDIUM
Notes: Verified against `TestCalculateMetrics`: with i=1..100ms, P50→idx49→50ms, P95→idx94→95ms, P99→idx98→99ms, Avg→50.5ms. Boundary-correct for the test's expectations. No division-by-zero (guards `total==0` and `len(latencies)==0`).

## Finding 6 — Server resource exhaustion simulation

Location: `internal/server/server.go:31-82`
Claimed Behavior: `POST /booking` simulates a bounded DB-connection pool; requests beyond capacity queue; context cancellation aborts cleanly.
Observed Implementation: `semaphore = make(chan struct{}, MaxDBConnections)`. Handler acquires via `select { case semaphore<-struct{}: case <-ctx.Done(): return }`. On context cancellation before acquire, returns without consuming/producing a slot (the `defer` draining is registered only after successful acquire). After acquire, waits `DBQueryDuration` (timer or context-done). `activeReq` mutated via `atomic`.
Assessment: PASS
Severity: MEDIUM
Notes: `MaxDBConnections <= 0` defaults to 5; `DBQueryDuration <= 0` defaults to 10ms. Method guard returns 405 for non-POST.

## Finding 7 — Demo reproducibility & physical consistency

Location: `cmd/demo/main.go:16-71`
Claimed Behavior: Smoke (2 VUs ≤ 5 conn) finishes fast; stress (50 VUs ≫ 5 conn) shows queuing tail latency.
Observed Implementation: `httptest` server with pool=5, 20ms query. Smoke prints ~188 req / ~21ms P95 / ~94 RPS. Stress prints ~473 req / ~210ms P95 / ~236 RPS. Independent re-run matches recorded `engineering/03-execution-result.md` within timing tolerance.
Assessment: PASS
Severity: LOW
Notes: Back-of-envelope consistency: 5 slots × (1000ms/20ms)=250 RPS cap; stress observed ~236 RPS, and 50 VUs competing for 5 slots ⇒ ~9 queued × 20ms ≈ 180ms wait + 20ms service ≈ 200ms avg. Output physically coherent.

## Finding 8 — Demo table formatting

Location: `cmd/demo/main.go:61-72`
Claimed Behavior: Nicely aligned metric table.
Observed Implementation: `tabwriter` with padding=2, minwidth=0. Works.
Assessment: PASS
Severity: LOW
Notes: `w.Flush()` called. No error returned/handled but stdout writes are best-effort here — acceptable for a demo.

## Finding 9 — `P90Latency` field never computed

Location: `internal/loadtest/metrics.go:18`
Claimed Behavior: Result struct exposes P90.
Observed Implementation: `P90Latency` field set but **never assigned** in `CalculateMetrics` (only P50/P95/P99). Always zero.
Assessment: WARNING
Severity: MEDIUM
Notes: Design doc `engineering/01-design.md:21` lists P90 among required latencies. Field exists but stays 0 — silent inconsistency. Does not affect claimed demo/test assertions (neither asserts P90). Fix candidate: add `res.P90Latency = percentile(sorted, 90)`.
