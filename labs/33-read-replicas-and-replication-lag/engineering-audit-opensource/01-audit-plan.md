# Engineering Audit Plan

Target Lab:
`labs/33-read-replicas-and-replication-lag`

Implementation Files:
- `cmd/demo/main.go`
- `internal/cluster/cluster.go`
- `internal/router/router.go`
- `go.mod`

Tests:
- `tests/replication_test.go`

Executable/Demo:
- `go run ./cmd/demo` (verified during audit, real output recorded)

Approved Research Inputs:
- `research/01-plan.md`, `research/02-sources.md`, `research/03-evidence.md`, `research/04-contradictions.md`, `research/05-report.md`, `research/06-open-questions.md`
- `research-audit/07-verdict.md` (research audit verdict; treated as approved research baseline)

Main Claims To Verify:
1. Async replication produces observable stale reads via naive round-robin replica routing.
2. Sticky session routing sends reads to primary within `StickyDuration` and to replica pool after TTL expiry.
3. Causal token / LSN routing waits for replica catch-up and falls back to primary on timeout.
4. Lag-aware routing excludes replicas exceeding `MaxLSNDiff` and falls back to primary.
5. Sync replication guarantees immediate replica freshness at the cost of write latency.
6. Concurrency safety under concurrent read/write load with race detector.
7. Demo output is real and matches claimed behavior.
8. README matches implementation and test commands.

Commands To Run:
- `go build ./...`
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Race conditions in `Node.WaitForLSN` goroutine lifecycle and `sync.Cond` usage.
- `ReadWithToken` fallback path silently swallows wait errors and may return stale primary data without surfacing the timeout.
- `ReadLagAware` LSN-diff comparison uses `primaryLSN >= appliedLSN` which is always true for valid LSNs; the `else` branch (replica ahead of primary) is unreachable — dead code, not a bug but indicates the threshold logic only handles lag, never negative lag.
- `Cluster.Write` holds `c.mu` RLock while sending on `replica.walChannel` with `default` drop — under burst, WAL entries can be silently dropped, causing permanent replica lag divergence with no error surfaced.
- Sync replication `Write` sleeps inside the write path holding `c.mu` RLock, blocking all other writes during the sleep (write serialization via RWMutex contention).
- `ReadWithStickySession` falls back to `ReadLagAware` (replica) after TTL even if the replica has not caught up — session freshness guarantee ends at TTL regardless of actual lag.
- Test `TestReadWithToken_LSN` asserts `elapsed < 150ms` wait but the replica lag is 200ms; the assertion `elapsed < 150ms` would FAIL if the token path actually waited. In practice the test passes because the token path finds a replica already caught up (race-dependent). This is a fragile/flaky test assertion.