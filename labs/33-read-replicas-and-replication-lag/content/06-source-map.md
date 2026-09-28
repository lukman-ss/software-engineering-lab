# Source Map

## Anomali Replication Lag (Stale Read)

Research:
- `research/05-report.md` — Finding 1: "Async Replication Causes Variable Lag"
- `research/06-open-questions.md` — Question 1 (Optimal Sticky Routing Duration)

Implementation:
- `internal/cluster/cluster.go` — `replicaWorker()` (goroutine per replica dengan configurable lag delay)
- `internal/router/router.go` — `ReadNaive()` (round-robin ke replica tanpa pengecekan lag)

Tests:
- `tests/replication_test.go` — `TestNaiveReplicationLag_StaleRead` (membuktikan stale read terjadi, lalu berhasil setelah lag)

Demo:
- `cmd/demo/main.go` — bagian 1 (async replication lag anomaly)

---

## Sticky Session Routing (Time-Based)

Research:
- `research/05-report.md` — Finding 3: "Session Guarantees via LSN/ClusterTime", Conclusion point 3
- `research/06-open-questions.md` — Question 1 (heuristic-based durations)
- `research-audit/06-gaps.md` — Gap 1: sticky window 5s bersifat heuristik

Implementation:
- `internal/router/router.go` — `ReadWithStickySession()` (cek `time.Since(lastWriteTime) < StickyDuration`)
- `internal/router/router.go` — `Write()` (simpan `SessionState` ke `sync.Map`)

Tests:
- `tests/replication_test.go` — `TestStickySessionRouting` (primary dalam TTL, replica setelah TTL, stale untuk session lain)

Demo:
- `cmd/demo/main.go` — bagian 2 (sticky session read-your-own-writes)

---

## Causal Token / Minimum LSN Routing

Research:
- `research/05-report.md` — Finding 3: "Session Guarantees via LSN/ClusterTime"
- `research/02-sources.md` — Source 4 (MongoDB causal sessions), Source 10 (Terry et al. session guarantees)

Implementation:
- `internal/router/router.go` — `ReadWithToken()` (pilih replica ≥ minLSN atau blocking wait)
- `internal/cluster/cluster.go` — `WaitForLSN()` (sync.Cond + context timeout)

Tests:
- `tests/replication_test.go` — `TestReadWithToken_LSN` (menunggu ~200ms, readLSN ≥ lsn)
- `tests/replication_test.go` — `TestWaitForLSN_ContextTimeout` (DeadlineExceeded)

Demo:
- `cmd/demo/main.go` — bagian 3 (causal token / LSN wait)

---

## Lag-Aware Routing & Primary Fallback

Research:
- `research/05-report.md` — Finding 4: "Middleware Supports Automatic Splitting", Finding 5: monitoring metrics

Implementation:
- `internal/router/router.go` — `ReadLagAware()` (filter ΔLSN ≤ MaxLSNDiff, fallback `primary (fallback-lag)`)

Tests:
- `tests/replication_test.go` — `TestReplicaLagThreshold_Fallback` (replica lag 10s, MaxLSNDiff=2 → fallback primary)

---

## Synchronous Replication (remote_apply)

Research:
- `research/05-report.md` — Finding 5: "Synchronous Replication Eliminates Lag but Increases Latency"
- `research/02-sources.md` — Source 1 (PostgreSQL `synchronous_commit`), Source 3 (Azure sync replication), Source 7 (Vitess semi-sync)

Implementation:
- `internal/cluster/cluster.go` — `Write()` sync path: `time.Sleep(lag)` per replica → apply → broadcast

Tests:
- `tests/replication_test.go` — `TestSynchronousReplication_Freshness` (ReadNaive langsung fresh di replica)

Demo:
- `cmd/demo/main.go` — bagian 4 (sync write + naive read)

---

## Concurrency & Thread Safety

Research:
- `research/06-open-questions.md` — Question 3 (LSN polling overhead)
- `research-audit/06-gaps.md` — Gap 2 (operational overhead of LSN check)

Implementation:
- `internal/cluster/cluster.go` — `sync.RWMutex`, `sync.Cond`, `atomic.Uint64`, buffered `walChannel`
- `internal/router/router.go` — `sync.Map` (session map), `atomic.Uint64` (round-robin index)

Tests:
- `tests/replication_test.go` — `TestConcurrentAccess_RaceFree` (15 goroutines × 20 ops, `go test -race` passes)

---

## Implementation Notes (diluar source code)

Engineering docs:
- `engineering/01-design.md` — architecture diagram, test strategy, implementation decisions
- `engineering/02-implementation-notes.md` — file list, design decisions, known limitations, trade-offs
- `engineering/03-execution-result.md` — build, test, race detector, demo output

Engineering audit:
- `engineering-audit/06-verdict.md` — `APPROVED`, zero warnings, zero blocking issues
