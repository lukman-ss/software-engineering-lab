# Engineering Design

Target Lab: labs/39-bloom-filters
Research Status: APPROVED (research-audit/07-verdict.md)

## Concept To Prove

A Bloom filter is a probabilistic set-membership structure that:
1. Never produces false negatives ("definitely not in set" is always correct).
2. May produce false positives ("possibly in set" may be wrong).
3. Achieves ~10 bits per element for a 1 % false-positive rate.
4. Reduces unnecessary I/O in an LSM-tree-like store and prevents cache penetration.

Mathematical basis (from approved research):
- ε ≈ (1 − e^(−kn/m))^k
- Optimal k = (m/n) · ln 2
- m/n ≈ −1.44 · log₂ ε  (bits per element for target ε)

At ε = 0.01: k ≈ 7, m/n ≈ 9.585 bits.

## Expected Behavior

| Operation | Result |
|-----------|--------|
| Add(x), then Check(x) | always true (no false negatives) |
| Check(y) for y never added | true with probability ε, false otherwise |
| Measured FP rate at ε=0.01 | converges to ~1 % over large sample |

## Failure Scenario

Without a Bloom filter:
- Every lookup in the LSM store hits disk even for absent keys.
- Cache penetration: repeated queries for absent keys bypass the cache and hammer the backend.

With a Bloom filter:
- Absent keys are rejected before disk access (with ≥ 99 % accuracy at ε = 0.01).
- Cache layer uses the filter as an admission gate; malicious absent-key queries are stopped.

## Success Criteria

- [ ] `Filter.Add` / `Filter.Check` obey the no-false-negative guarantee across all tests.
- [ ] Measured FP rate ≤ 2 × ε (within statistical noise on 10 000 samples).
- [ ] LSM store demo shows significant reduction in "disk access" count with filter enabled.
- [ ] Cache penetration demo shows backend calls drop to ≈ ε fraction when filter is active.
- [ ] All tests pass with `go test ./...`.
- [ ] Race detector passes with `go test -race ./...`.

## Architecture

```
labs/39-bloom-filters/
├── internal/
│   ├── bloom/
│   │   ├── bloom.go          # Filter: bit array + k hash functions
│   │   └── bloom_test.go     # unit + property tests + benchmarks
│   └── store/
│       ├── lsm.go            # Toy LSM store with optional Bloom filter per segment
│       ├── lsm_test.go
│       ├── cache.go          # Cache with Bloom filter anti-penetration gate
│       └── cache_test.go
├── cmd/demo/
│   └── main.go               # Prints observable numbers from all three scenarios
├── go.mod
└── README.md
```

## Components

### bloom.Filter
- Bit array backed by `[]uint64` (word-aligned, 64 bits per word).
- Double-hashing (Kirsch-Mitzenmacher): derive k positions from two base hashes.
  - h_i(x) = (h1(x) + i·h2(x)) mod m
  - Base hashes: FNV-1a (h1) and a shifted XOR variant (h2) — stdlib only, no imports.
- Constructor `New(n uint, epsilon float64)` computes optimal m and k.
- Thread-safe variant `SyncFilter` wraps with `sync.RWMutex`.

### store.LSMStore
- Toy representation: one or more immutable "segments" (slices of key strings).
- Each segment may carry a `*bloom.Filter`.
- `Get(key)` without filter: linear scan all segments (disk hit counted).
- `Get(key)` with filter: skip segment if filter says "definitely not".
- Exposes `DiskHits()` counter for demo comparison.

### store.Cache
- LRU-like in-memory map (implementation decision: simple `map` with capacity, no eviction complexity needed).
- Optional Bloom filter as admission gate:
  - On Get(key): if filter says "definitely not in any seen key set", skip backend entirely.
  - Tracks `BackendCalls` counter.

## Test Strategy

1. **No false negatives** – add 10 000 random keys, verify all return true.
2. **False positive rate** – sample 100 000 unseen keys; measure empirical FP ≤ 2 × ε.
3. **Optimal parameters** – verify computed m and k match formulas within 1 ULP.
4. **LSM disk hit reduction** – assert DiskHits(with filter) < DiskHits(without filter).
5. **Cache penetration** – assert BackendCalls(with filter) < BackendCalls(without filter).
6. **Concurrency** – 100 goroutines concurrent Add/Check on SyncFilter; -race passes.

## Execution Plan

```bash
cd labs/39-bloom-filters
go mod tidy
go test ./...
go test -race ./...
go run ./cmd/demo
```

## Implementation Decisions

1. **FNV-1a for hashing**: Standard library has no MurmurHash3; FNV-1a is in stdlib (hash/fnv) and approved by research as suitable.
2. **Double hashing for k probes**: Kirsch-Mitzenmacher technique derives k positions from two independent hashes. Noted as open question in research; chosen as the canonical efficient approach. No extra dependency.
3. **[]uint64 bit array**: Word-granularity, cache-friendly. Bit access via index/mask arithmetic.
4. **Toy LSM/Cache**: Demonstrates the I/O pattern, not production-grade persistence.
5. **SyncFilter wraps Filter**: Adds `sync.RWMutex`; Read acquires RLock, Write acquires Lock.
