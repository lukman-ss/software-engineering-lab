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
- `research/05-report.md`
- `research-audit/07-verdict.md` (APPROVED)

Main Claims To Verify:
1. Zero False Negatives: Querying inserted elements returns true 100% of the time.
2. Controlled False Positive Rate: Empirical false positive rate matches theoretical target $\epsilon \approx (1 - e^{-kn/m})^k$ within acceptable statistical bounds.
3. Optimal Sizing: Calculated parameters $m$ (bit array size) and $k$ (hash count) match Bloom filter mathematical formulas.
4. LSM-Tree I/O Reduction: Segment skipping via Bloom filters reduces simulated disk lookups for absent keys by >95% (target ~99% at $\epsilon=0.01$).
5. Cache Penetration Defense: Absent key queries are filtered before querying the expensive backend, reducing cache misses/penetration backend hits.
6. Concurrency Safety: `SyncFilter` safely supports concurrent reads and writes without data races under `go test -race`.
7. Standard Library Purity: No unapproved third-party dependencies.

Commands To Run:
- `rtk go test ./...`
- `rtk go test -v ./...`
- `rtk go test -race ./...`
- `rtk go run ./cmd/demo`

Primary Risks:
- Mathematical rounding or overflow errors in $m$ or $k$ calculation.
- Poor dispersion/correlation in double hashing ($h_1 + i \cdot h_2 \pmod m$) using FNV-1a + rotate-XOR causing inflated false positive rates.
- Concurrency bugs or race conditions in `SyncFilter` or `store.Cache` / `store.LSMStore`.
- Divergence between documented design in `engineering/01-design.md` and actual code implementation.
