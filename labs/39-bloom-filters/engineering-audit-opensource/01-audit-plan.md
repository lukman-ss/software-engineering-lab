# Engineering Audit Plan

Target Lab: labs/39-bloom-filters
Implementation Files:
- `internal/bloom/bloom.go`
- `internal/store/lsm.go`
- `internal/store/cache.go`
- `cmd/demo/main.go`

Tests:
- `internal/bloom/bloom_test.go`
- `internal/store/lsm_test.go`

Executable/Demo:
- `cmd/demo/main.go`

Approved Research Inputs:
- `research/05-report.md`
- `research/runs/2026-09-28-bloom-filters/05-report.md`
- `research-revision/02-changes-made.md`
- `research-audit/07-verdict.md` (Status: APPROVED)

Main Claims To Verify:
1. Zero false negatives guarantee across tested keys ($P(\text{false negative}) = 0$).
2. Theoretical sizing calculations:
   - $m = \lceil -n \cdot \ln(\epsilon) / (\ln 2)^2 \rceil$
   - $k = \text{round}((m/n) \cdot \ln 2)$
   - bits per element at $\epsilon = 0.01 \approx 9.59$.
3. Empirical false-positive rate aligns with theoretical target $\epsilon = 0.01$ within statistical bounds.
4. LSM-tree disk I/O reduction $> 99\%$ for absent-key lookups.
5. Cache penetration mitigation drops backend calls on absent keys to near zero.
6. Thread safety of `SyncFilter` under concurrent reads and writes (`-race` clean).
7. Demo runs without mock shortcuts or hardcoded outputs.

Commands To Run:
```bash
go build ./...
go test -v -count=1 ./...
go test -race -count=1 ./...
go run ./cmd/demo
```

Primary Risks:
- Hash quality risk: Stdlib FNV-1a with rotate-XOR double hashing could degrade false positive rate or cause clustering under adversarial patterns.
- Concurrency hazard: `SyncFilter` synchronization leaks or race conditions during concurrent Add/Check.
- Boundary conditions: $n=0$, $\epsilon \le 0$, $\epsilon \ge 1$, empty filter behavior, large $n$ integer overflows.
- Doc/code discrepancy: Mismatch between claims in README/design notes and actual implementation constants or behavior.
