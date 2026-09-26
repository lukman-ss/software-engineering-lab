# Docs vs Code Comparison

## Finding 1
Document: README.md, line 6
Claim: "Entry point running comparative Smoke vs. Stress test scenarios."
Code: cmd/demo/main.go runs smoke test (2 VUs) then stress test (50 VUs) and prints results
Assessment: PASS
Severity: LOW
Notes: The demo exactly matches the documented behavior.

## Finding 2
Document: README.md, line 8
Claim: "Mock service with constrained connection pool capacity to demonstrate saturation."
Code: internal/server/server.go uses a buffered channel as semaphore with capacity = MaxDBConnections
Assessment: PASS
Severity: LOW
Notes: Correct implementation of connection pool bottleneck.

## Finding 3
Document: README.md, line 9
Claim: "Built-in concurrency runner and percentile calculator."
Code: internal/loadtest/runner.go (concurrency orchestrator) and internal/loadtest/metrics.go (percentile calculator)
Assessment: PASS
Severity: LOW
Notes: Matches documentation exactly.

## Finding 4
Document: engineering/01-design.md, line 26-28
Claim: Architecture: internal/server (HTTP server with constrained connection pool), internal/loadtest (load test harness), cmd/demo (executable running smoke then stress test)
Code: Exactly matches the claimed architecture
Assessment: PASS
Severity: LOW
Notes: No deviation from documented architecture.

## Finding 5
Document: engineering/01-design.md, line 31-33
Claim: Components: BookingServer (HTTP handler with configurable concurrency limit), LoadTester (generates HTTP traffic), MetricsAggregator (thread-safe latency collector)
Code: Server = BookingServer, Runner = LoadTester, CalculateMetrics = MetricsAggregator
Assessment: PASS
Severity: LOW
Notes: Component mapping is accurate.

## Finding 6
Document: engineering/01-design.md, line 50
Claim: "Hand-rolled percentile calculation over standard arrays rather than adding HdrHistogram dependency."
Code: internal/loadtest/metrics.go uses sort.Slice on a copied array
Assessment: PASS
Severity: LOW
Notes: Matches documented decision.

## Finding 7
Document: engineering/01-design.md, line 51
Claim: "Semaphore pattern (buffered channel) used in the server handler to simulate database connection pool bottlenecks"
Code: internal/server/server.go line 41: semaphore := make(chan struct{}, cfg.MaxDBConnections)
Assessment: PASS
Severity: LOW
Notes: Exactly matches documentation.

## Finding 8
Document: engineering/01-design.md, line 52
Claim: "Per-goroutine slices in load generator to avoid mutex contention, aggregating once on completion"
Code: internal/loadtest/runner.go line 48-53: type vuResult { latencies []time.Duration; errors int }, line 107-112: aggregation after wg.Wait()
Assessment: PASS
Severity: LOW
Notes: Exactly matches documented implementation.

## Finding 9
Document: engineering/01-design.md, line 54
Claim: "Simulated external network latency is omitted to isolate the connection pool bottleneck"
Code: No artificial network delay is added in either server or client
Assessment: PASS
Severity: LOW
Notes: Confirmed by inspection.

## Finding 10
Document: engineering/01-design.md, line 55
Claim: "Used standard net/http client which carries its own connection pooling limits; overridden using a custom http.Transport"
Code: internal/loadtest/runner.go line 31-34: Custom http.Transport with MaxIdleConns=1000, MaxIdleConnsPerHost=1000
Assessment: PASS
Severity: LOW
Notes: Exactly matches documented decision and implementation match.

## Finding 11
Document: engineering/01-design.md, line 23-24 (Expected Behavior)
Claim: "**Smoke Load (low VUs)**: All requests process within normal latency limits. Average and P95 are close. Error rate is 0%." "**Stress Load (high VUs)**: Concurrent requests exceed server resource capacity (connection pool limit). Tail requests queue up. P95 and P99 latency spikes significantly, while average latency degrades less severely, proving the masking effect of averages."
Code: Demo output shows:
  Smoke: Avg=21.2ms, P95=21.6ms (close), Errors=0
  Stress: Avg=612ms, P95=1105ms, P99=1486ms (significant tail spike)
Assessment: PASS
Severity: LOW
Notes: The demo output precisely matches the expected behavior described.

## Finding 12
Document: engineering/02-implementation-notes.md, line 16
Claim: "Wait durations in server mock are fixed (20ms), making the tail latency strictly a function of queuing time when VUs exceed max DB connections."
Code: internal/server/server.go lines 68-73: Base duration = DBQueryDuration, but with 10% chance of 25x duration when activeReq > MaxDBConnections
Assessment: DOC_CODE_MISMATCH
Severity: MEDIUM
Notes: The documentation claims tail latency is "strictly a function of queuing time" but the implementation adds random slow queries (25x duration 10% of the time when over capacity). This is a mismatch between documented claims and actual code behavior. While the random behavior is realistic, it violates the explicit claim made in the notes.

## Finding 13
Document: engineering/02-implementation-notes.md, line 20-21
Claim: "Percentile algorithm is exact-sort, memory footprint scales linearly with request count; unsuitable for hour-long multi-million RPS benchmarks, but perfect for a 2-second lab test."
Code: internal/loadtest/metrics.go lines 45-49: Copy slice, sort.Slice
Assessment: PASS
Severity: LOW
Notes: Matches documentation exactly.

## Finding 14
Document: engineering/02-implementation-notes.md, line 31
Claim: "Measuring P50, P95, and P99 percentiles manually."
Code: internal/loadtest/metrics.go calculates P50, P90, P95, P99
Assessment: PASS
Severity: LOW
Notes: Correct.

## Finding 15
Document: engineering/02-implementation-notes.md, line 32
Claim: "Smoke load executing linearly (fast average, fast P95)."
Code: Demo output shows smoke test has low latency and close P95/Avg
Assessment: PASS
Severity: LOW
Notes: Validated by demo.

## Finding 16
Document: engineering/02-implementation-notes.md, line 33
Claim: "Stress load queuing behind a saturation point (connection pool limit), forcing P95 to severely degrade."
Code: Demo output shows stress P95 (1105ms) >> smoke P95 (21.6ms)
Assessment: PASS
Severity: LOW
Notes: Validated by demo.

## Finding 17
Document: README.md, lines 20-23
Claim: Instructions for running tests and race detector: go test -v ./..., go test -race ./...
Code: These commands work and produce successful output
Assessment: PASS
Severity: LOW
Notes: Documentation matches executable commands.

## Finding 18
Document: README.md, line 12
Claim: "Go 1.22+"
Code: go.mod specifies go 1.22
Assessment: PASS
Severity: LOW
Notes: Version requirement matches.