# Engineering Audit Plan

Target Lab: `labs/39-bloom-filters`
Implementation Files:
- `internal/bloom/bloom.go`
- `internal/store/lsm.go`
- `internal/store/cache.go`
Tests:
- `internal/bloom/bloom_test.go`
- `internal/store/lsm_test.go`
Executable/Demo:
- `cmd/demo/main.go`
Approved Research Inputs:
- `labs/39-bloom-filters/research/05-report.md`
- `labs/39-bloom-filters/research-audit/07-verdict.md`
- `labs/39-bloom-filters/engineering/01-design.md`
- `labs/39-bloom-filters/engineering/02-implementation-notes.md`

Main Claims To Verify:
1. Zero false negatives guarantee across inserted elements.
2. Empirical false positive rate matches theoretical $\epsilon$ ($\approx 1\%$) for configured $n=10,000$ and $m/n \approx 9.59$ bits/element.
3. Optimal sizing formulas $m = \lceil -n \ln(\epsilon) / (\ln 2)^2 \rceil$ and $k = \text{round}((m/n) \ln 2)$ implemented accurately.
4. LSM store reduces simulated disk reads by $>99\%$ on absent keys when per-segment Bloom filters are active.
5. Cache admission gate prevents cache penetration attacks on absent keys without starving or corrupting valid keys.
6. Thread-safe operations via `SyncFilter` are data-race-free under concurrent read/write workloads.
7. No third-party dependencies used (pure standard library).

Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Poor hash independence in double-hashing causing elevated false positive rates.
- Concurrency data race on `Filter` or `SyncFilter` bit slices during parallel writes.
- Off-by-one or bit-shifting overflow in 64-bit word indexing.
- Test assertions too loose or falsified (e.g. mock results hardcoded).
