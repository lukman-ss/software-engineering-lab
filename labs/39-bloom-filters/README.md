# Lab 39: Bloom Filters

Bloom filters are probabilistic set-membership data structures that use a compact bit array and multiple hash functions to answer: "Is this element in the set?"

The answer is either:
- **Definitely not in the set** — always correct (zero false negatives).
- **Possibly in the set** — correct with probability `1 - epsilon` (false positives occur at rate `epsilon`).

## Mathematical Foundation

| Parameter | Formula |
|-----------|---------|
| Bit array size (m) | `ceil(-n * ln(epsilon) / (ln 2)^2)` |
| Optimal hash count (k) | `round((m/n) * ln 2)` |
| Bits per element at 1% FP | ~9.59 (`m/n ≈ −1.44 * log₂(epsilon)`) |

At `n=10,000, epsilon=0.01`: `m=95,851 bits (≈11.7 KB)`, `k=7`.

## What This Lab Demonstrates

1. **Core Bloom filter correctness**: zero false negatives across 10,000 inserted keys; empirical FP rate ≈ 0.99% (target 1.00%).
2. **LSM-tree disk I/O reduction**: absent-key lookups across 10 segments drop from 100,000 disk reads (no filter) to ~988 reads (with Bloom filter) — a **>99% reduction**.
3. **Cache penetration mitigation**: 1,000 adversarial absent-key queries reach the backend 1,000 times without protection, versus ~0 times with Bloom filter admission gate.

## Running the Lab

```bash
cd labs/39-bloom-filters

# Run tests
go test ./...

# Run with race detector
go test -race ./...

# Run demo
go run ./cmd/demo
```

## Structure

```
labs/39-bloom-filters/
├── cmd/demo/main.go              # Runnable demonstration
├── internal/
│   ├── bloom/
│   │   ├── bloom.go              # Filter + SyncFilter implementation
│   │   └── bloom_test.go         # Unit tests and benchmarks
│   └── store/
│       ├── lsm.go                # Simulated LSM-tree with per-segment filters
│       ├── lsm_test.go           # LSM disk read + cache penetration tests
│       └── cache.go              # Cache with Bloom filter admission gate
├── go.mod
└── engineering/
    ├── 01-design.md
    ├── 02-implementation-notes.md
    └── 03-execution-result.md
```

## Key Guarantees

| Property | Value |
|----------|-------|
| False negative rate | **0%** (absolute guarantee) |
| False positive rate | ~`epsilon` (configurable; default 1%) |
| Memory at 1% FP | ~9.59 bits/element |
| Deletion support | Not supported (use Counting Bloom or Cuckoo filter) |
| Thread safety | `SyncFilter` variant using `sync.RWMutex` |

## Dependencies

Standard library only (`hash/fnv`, `math`, `sync`). No third-party packages.
