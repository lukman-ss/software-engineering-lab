# Docs vs Code

## Source Documents Reviewed
- README.md (lab overview)
- engineering/01-design.md (Concept, Expected Behavior, Success Criteria, Architecture)
- engineering/02-implementation-notes.md (Files added, core decisions, limitations, trade-offs, what is demonstrated/not demonstrated)
- engineering/03-execution-result.md (Build, Tests, Race Detector, Demo, Final Status)

## Claims From README
1. **Structure**: `cmd/demo` entry point; `internal/server` mock service with constrained connection pool; `internal/loadtest` concurrency runner and percentile calculator; `tests` automated integration tests; `engineering/` design/notes/results.
   - Verified: tree matches exactly. PASS.
2. **Demo Command**: `go run ./cmd/demo` runs a comparative Smoke vs Stress test.
   - Verified: runs successfully, prints two metric tables. PASS.
3. **Test Commands**: `go test -v ./...` and `go test -race ./...`.
   - Verified: both succeed with all tests passing and no races. PASS.
4. **Requirements**: Go 1.22+.
   - Verified: go.mod specifies `go 1.22`. PASS.

## Claims From engineering/01-design.md
### Concept To Prove (lines 7-11)
- Load testing reveals boundaries, saturation, degradation patterns. → Demo shows Stress P95/P99 >> Smoke; tests assert stress P95 > smoke P95. PASS.
- Average response time conceals tail latency; percentiles (P95, P99) necessary. → Implementation computes P50/P90/P95/P99; demo prints them. Stress Avg 477ms vs P95 710ms shows tail > average. PASS.
- Incremental test stages (smoke vs stress) differentiate baseline vs resource exhaustion. → Demo runs smoke then stress; test contrasts 1 vs 10 VUs. PASS.
- Downstream resource saturation (e.g., DB pool) causes non-linear latency degradation for tail requests. → Server semaphore limits DB slots; when VUs > max, requests queue. Stress P95 blowup vs smoke confirms. PASS.

### Expected Behavior (lines 12-15)
- Smoke Load (low VUs): all requests within normal latency; Avg and P95 close; error rate 0%. → Demo Smoke: Avg 21ms, P95 22ms, Errors 0. Test `TestLoadTest_SmokeVsStress` smoke half implicitly has low errors (not asserted but observed). PASS.
- Stress Load (high VUs): concurrent requests exceed server capacity; tail requests queue; P95/P99 spike significantly; Avg degrades less severely (masking). → Demo Stress: Avg 477ms, P95 710ms, P99 1148ms (P95 > Avg, P99 > P95). PASS.
- Failure Scenario: excessive load → queuing → high latency + severe tail degradation for P95/P99. → Same evidence; PASS.
- Success Criteria: (1) automated benchmarks/exec without external deps; (2) load generator computes Min/Max/Avg/P50/P90/P95/P99; (3) demo contrasts Smoke vs Stress metrics; (4) all tests pass with `go test -race ./...` zero races.
   - (1) stdlib only, no external deps in go.mod. PASS.
   - (2) load generator computes all seven statistics (internal/loadtest/Result includes them; demo prints a subset). PASS.
   - (3) demo prints comparative tables. PASS.
   - (4) tests pass with -race. PASS.

### Architecture (lines 25-28) and Components (lines 30-34)
Matches file layout: server (BookingServer semaphore), loadtest (LoadTester VUs + MetricsAggregator), demo (executable). PASS.

### Test Strategy (lines 35-39)
- Unit tests: statistical calculations → metrics_test.go (PASS)
- Integration tests: actual HTTP load against mock server → loadtest_test.go (PASS)
- Concurrency test: race detector → `go test -race ./...` clean (PASS)

### Execution Plan (lines 40-48)
Steps 1-8 all occurred; engineering notes and code confirm. PASS.

### Implementation Decisions (lines 49-51)
- stdlib net/http/sync only, no third-party deps → PASS.
- Ponytail: exact sorting for percentiles (instead of HdrHistogram) → exact sorting used (metrics.go sort.Slice); tests pass; documented as ponytail. PASS.

## Claims From engineering/02-implementation-notes.md
### Core Design Decisions (lines 10-14)
- Hand-rolled percentile vs HdrHistogram: confirmed (metrics.go). PASS.
- Semaphore pattern in server to simulate DB pool bottlenecks → server.go uses buffered channel semaphore; TestServer_MaxDBConnectionsBound validates bound. PASS.
- Per-goroutine slices in load generator to avoid mutex contention → runner.go per-VU results[vuID]; -race clean. PASS.

### Implementation-Specific Choices (lines 15-18)
- Demo fixed query duration 20ms; server configurable DBQueryDuration; tail latency strictly queuing time when VUs exceed max DB connections. → PARTIAL: tail latency is queuing OR 10% random slowdown (activeReq > maxDBConnections). The note overclaims; it is not *strictly* queuing. LOW DOC-CODE mismatch (see Code Audit Finding 6).
- Load generator loops with `time.Since` rather than pre-generating requests → runner.go line 65-100 loop uses `time.Since` for duration; no pregen. PASS.

### Known Limitations (lines 19-22)
- Percentile exact-sort memory scales linearly with request count; unsuitable for long/million-RPS but fine for 2-second lab → metrics.go sorts all latencies; comment notes slice sorting fine for <1M; acceptable. PASS.
- Simulated external network latency omitted to isolate connection pool bottleneck → server has no network; client uses custom Transport with huge idle conns to avoid being bottleneck. PASS.

### Trade-offs (lines 23-25)
- stdlib net/http client carries own connection-pool limits; overridden via custom Transport to ensure load generator not bottleneck → runner.go lines 31-40 sets MaxIdleConns/PerHost to 1000. PASS.

### What Is Demonstrated (lines 26-30)
- Measuring P50, P95, P99 percentiles → metrics.go + demo prints them. PASS.
- Smoke load linear (fast Avg, fast P95) → demo Smoke Avg/P95 ~21ms. PASS.
- Stress load queues behind saturation point (connection pool limit), forcing P95 severe degradation → demo Stress P95 710ms vs Smoke P95 22ms. PASS.

### What Is Not Demonstrated (lines 31-33)
- Distributed load generation (multi-node) → engineering notes state this is not demonstrated; code has no such capability. PASS (honest).
- Real database lock contention or CPU saturation (mocked via semaphore and time.Sleep) → server uses semaphore + time.Sleep (timer) for DB delay; no real DB or CPU load. PASS.

## Claims From engineering/03-execution-result.md
### Build
- `go build ./...` success → verified. PASS.
### Tests
- `go test -v ./...` output listed 6 test functions but actual run shows 9 functions (omitted: TestCalculateMetrics_SingleSample, TestCalculateMetrics_RPS, TestServer_MaxDBConnectionsBound, TestLoadTest_SuccessAndErrorInvariant). The document says:
```
?   	github.com/lukman/software-engineering-lab/labs/15-load-testing/cmd/demo	[no test files]
=== RUN   TestCalculateMetrics
...
PASS
ok  	github.com/lukman/software-engineering-lab/labs/15-load-testing/internal/loadtest	0.347s
?   	github.com/lukman/software-engineering-lab/labs/15-load-testing/internal/server	[no test files]
=== RUN   TestLoadTest_SmokeVsStress
...
PASS
ok  	github.com/lukman/software-engineering-lab/labs/15-load-testing/tests	1.657s
```
It lists only the two package summaries and omits four tests that were added after the doc was written. This is a DOC vs CODE mismatch: the execution result under-reports the test suite. LOW severity.
### Race Detector
- `go test -race ./...` success → verified. PASS.
### Demo
- Sample run showing Smoke (2 VUs) vs Stress (50 VUs) with Stress P95/P99 >> Smoke. My actual run reproduced the invariant: Stress P95 710ms vs Smoke P95 22ms; Stress P99 1148ms vs Smoke P99 25ms. The doc notes "Output metrics vary per run..." and the key invariant holds. PASS.
### Final Engineering Status
- READY_FOR_ENGINEERING_AUDIT → audit in progress.

## Summary of DOC_CODE_MISMATCH findings
1. engineering/03-execution-result.md test list is stale (omits 4 tests added later). LOW.
2. engineering/02-implementation-notes.md claim "Tail latency is strictly a function of queuing time when VUs exceed max DB connections" is inaccurate due to the 10% random slowdown active when `activeReq > MaxDBConnections`. LOW.
3. README and engineering notes correctly describe the structure, commands, and success criteria; no other mismatches found.