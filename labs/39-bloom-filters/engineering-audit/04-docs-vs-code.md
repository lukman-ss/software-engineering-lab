# Docs vs Code Audit

## Documentation Claims vs Implementation Reality

| Claim in README / Design / Research | Code Location | Observed Reality | Audit Status |
|---|---|---|---|
| Zero false negatives guarantee | `internal/bloom/bloom.go:42-53` & `bloom_test.go:12-28` | Tested across 10,000 keys; 0 false negatives observed | MATCH |
| Optimal parameter calculation $m = \lceil -n \ln(\epsilon) / (\ln 2)^2 \rceil$ | `internal/bloom/bloom.go:86-91` | Exact math formula used via `math.Log` and `math.Ceil` | MATCH |
| Optimal hash count $k = \text{round}((m/n)\ln 2)$ | `internal/bloom/bloom.go:94-102` | Exact math formula used via `math.Round` with $\min(k) = 1$ | MATCH |
| Double-hashing scheme for probe generation | `internal/bloom/bloom.go:65-84` | Double-hashing $h_i = (h_1 + i \cdot h_2) \pmod m$ implemented using stdlib `hash/fnv` | MATCH |
| LSM disk I/O reduction > 99% | `internal/store/lsm.go:56-74` & `cmd/demo/main.go:85` | Demo outputs 99.01% reduction (100,000 reads down to 988 reads) | MATCH |
| Cache penetration defense stops absent queries | `internal/store/cache.go:55-58` & `cmd/demo/main.go:122` | Unprotected cache had 5000/5000 hits; protected cache had 0 hits | MATCH |
| Concurrency safety via `SyncFilter` | `internal/bloom/bloom.go:104-134` | `sync.RWMutex` correctly guards operations and passes `-race` | MATCH |
| Zero external dependencies | `go.mod` | Pure stdlib (`hash/fnv`, `math`, `sync`, `sync/atomic`) | MATCH |

## Discrepancy Checks

- `DOC_CODE_MISMATCH`: None found.
- `TEST_CLAIM_MISMATCH`: None found.
- `RESEARCH_IMPLEMENTATION_MISMATCH`: None found.
- `FAKE_DEMO`: None found. Demo runs real live computations without hardcoded values.
- `FAKE_BENCHMARK`: None found. Benchmarks run standard Go testing harness.
