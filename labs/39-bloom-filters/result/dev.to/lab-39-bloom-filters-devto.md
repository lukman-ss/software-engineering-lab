---
title: "Bloom Filters: Probabilistic Data Structures for I/O Reduction and Cache Penetration Defense"
published: false
description: "How Bloom filters prune 99% of disk reads in LSM-trees and block cache penetration attacks — with zero false negatives, empirical FP verification, and a pure Go standard library implementation."
tags: go, datastructures, performance, architecture
---

# Bloom Filters: Probabilistic Data Structures for I/O Reduction and Cache Penetration Defense

An absent-key query in an LSM-tree forces redundant disk lookups across every SSTable segment. Unchecked non-existent queries also penetrate caching layers, exhausting database connections and inflating tail latency.

> Probabilistic filters prevent unnecessary backend dives. One false positive is better than a thousand disk reads.

Lab 39 implements an in-memory Bloom filter in Go standard library to demonstrate deterministic I/O pruning and cache defense.

## Disk I/O Mahal, Cache Penetration Mengancam

Disk I/O is the most expensive operation in database engines; every segment lookup consumes disk heads and system bandwidth. When non-existent keys bypass cache, they penetrate downstream backends linearly, triggering connection pool exhaustion and tail latency spikes.

In this lab, 10,000 absent queries without filtering trigger 100,000 segment reads across 10 SSTables. With a Bloom filter configured at 1% error rate, reads plummet to 988—a 99.01% reduction.

## Bloom Filter: Definisi dan Jaminan Asimetris

A Bloom filter is a space-efficient probabilistic data structure that provides asymmetric membership guarantees. It guarantees zero false negatives: an element present in the set will never be reported absent.

Conversely, false positives are possible with rate ε:
- **Definitely Not in Set**: A `false` result confirms key absence with 100% certainty.
- **Possibly in Set**: A `true` result indicates the key exists or hash bits collided. Standard filters do not support deletion.

## Fondasi Matematika Ukuran Bit Array dan Hash Optimal

The memory footprint *m* and probe count *k* are determined analytically from expected element capacity *n* and error rate target ε using closed formulas:

```
m = ceil(-n × ln(ε) / (ln 2)²)
k = round((m/n) × ln 2)
```

For *n* = 10,000 and ε = 0.01:
- `m` = 95,851 bits (~11.7 KB)
- `k` = 7 hash probes
- `m/n` = 9.59 bits per element

| Target FP Rate ε | Bits per Element | Probes k |
| --- | --- | --- |
| 0.10 (10%) | 4.79 bits | 3 |
| 0.01 (1%) | 9.59 bits | 7 |
| 0.001 (0.1%) | 14.38 bits | 10 |

## Optimasi Double Hashing Kirsch-Mitzenmacher

Executing *k* distinct cryptographic hash functions per key degrades CPU throughput. Kirsch and Mitzenmacher (2006) proved that two 64-bit base hashes linearly combined as `(h1 + i × h2) mod m` simulate *k* independent hashes without loss of accuracy.

```go
func baseHashes(key []byte) (uint64, uint64) {
    h := fnv.New64a()
    _, _ = h.Write(key)
    h1 := h.Sum64()
    h2 := h1 ^ (h1>>17 | h1<<47)
    if h2 == 0 { h2 = 0xdeadbeefdeadbeef }
    return h1, h2
}
```

This computation runs entirely in CPU registers with zero memory allocations.

## Arsitektur Bit Array Word-Aligned dan API Inti

To maximize memory locality and CPU throughput, the bit array is stored as a 64-bit word-aligned slice `[]uint64`.

```go
type Filter struct {
    bits []uint64
    m    uint
    k    uint
}
```

Bit index manipulation operates via division and modulo masks:
- Set: `f.bits[pos/64] |= 1 << (pos % 64)`
- Check: `if f.bits[pos/64] & (1 << (pos % 64)) == 0 { return false }`

Check performs short-circuit evaluation, returning immediately upon finding the first zero bit.

## Arsitektur Penyimpanan LSM-Tree Berbasis Segment Filter

LSM-Tree databases store data in multiple immutable SSTable segments. Looking up absent keys forces linear disk scanning across older segments.

```go
func (s *LSMStore) Get(key string) (string, bool) {
    for i := len(s.segments) - 1; i >= 0; i-- {
        seg := s.segments[i]
        if seg.filter != nil && !seg.filter.Check(keyBytes) {
            continue // Disk read completely skipped
        }
        atomic.AddUint64(&s.diskReads, 1)
        // load data block...
    }
}
```

Attaching a small Bloom filter to each segment header prunes absent queries prior to physical disk access.

## Pengurangan 99% Disk I/O pada Pencarian Kunci Absen

In benchmark scenario two, 10,000 absent queries were executed against an LSM store consisting of 10 segments holding 20,000 keys.

| Experiment Metric | Without Filter | With Bloom Filter |
| --- | --- | --- |
| Data Segments | 10 segments | 10 segments |
| Stored Keys | 20,000 records | 20,000 records |
| Absent Queries | 10,000 lookups | 10,000 lookups |
| Total Segment Reads | 100,000 reads | 988 reads |
| I/O Reduction | Baseline (0%) | 99.01% reduction |

The 0.988% leak rate aligns with the theoretical 1.00% false positive parameter.

## Anatomi Serangan Cache Penetration pada Backend

Cache penetration happens when malicious or random queries for non-existent identifiers bypass in-memory caches.

Normal caching mechanics:
1. Client requests non-existent key.
2. Cache reports miss, forwards request downstream.
3. Database executes disk queries, returns empty result.
4. Because null is often un-cached, subsequent requests hit database repeatedly.

In our lab attack simulation, 5,000 random absent queries produced 5,000 direct backend calls under unprotected cache configurations.

## Bloom Filter sebagai Admission Gate pada Arsitektur Cache

Placing a thread-safe Bloom filter (`SyncFilter`) in front of the cache creates an admission barrier.

```go
func (c *Cache) Get(key string) (string, bool) {
    if c.filter != nil && !c.filter.Check([]byte(key)) {
        return "", false // Rejected before contacting backend
    }
    // Check local memory cache, then fetch backend on miss
}
```

Absent keys evaluate to false within ~150 nanoseconds in memory, preventing backend connection spikes.

## Pembuktian Invarian Nol False Negative pada 10.000 Kunci

Zero false negatives is an absolute invariant. If a filter reports an existing item absent, the system suffers silent data loss.

```go
func TestNoFalseNegatives(t *testing.T) {
    bf := New(10000, 0.01)
    for i := 0; i < 10000; i++ {
        bf.Add([]byte(keys[i]))
    }
    for _, k := range keys {
        if !bf.Check([]byte(k)) {
            t.Fatalf("False negative detected for key: %s", k)
        }
    }
}
```

The test asserted 10,000 inserted keys; zero were missed (0.00% false negative rate).

## Validasi Tingkat False Positive Empiris pada 50.000 Kueri

`TestFalsePositiveRate` populated 10,000 keys and probed 50,000 absent keys against theoretical tolerance:

| Metric | Theoretical Target | Measured Output |
| --- | --- | --- |
| Filter Capacity *n* | 10,000 elements | 10,000 elements |
| Target FP Rate ε | 0.0100 (1.00%) | 0.0098 - 0.00996 (~0.99%) |
| Absent Key Queries | 50,000 queries | 50,000 unique queries |
| Detected False Positives | ~500 hits | 489 - 498 hits |
| Tolerance Limit | 2.5 × target (2.5%) | PASS |

Empirical measurements closely match the theoretical formula ε ≈ (1 - e^(-kn/m))^k.

## Verifikasi Presisi Dimensi Memori dan Jumlah Hash

`TestOptimalSizing` locks implementation correctness against floating-point rounding errors:

```go
func TestOptimalSizing(t *testing.T) {
    const n = 1000
    const eps = 0.01
    bf := New(n, eps)
    expectedM := uint(math.Ceil(-float64(n)*math.Log(eps)/(math.Log(2)*math.Log(2))))
    expectedK := uint(math.Round(float64(expectedM)/float64(n)*math.Log(2)))
    // Assert bf.M() == expectedM && bf.K() == expectedK
}
```

For *n*=1,000 and ε=0.01: `expectedM` = 9,586 bits, `expectedK` = 7 probes. Result matches theoretical calculations precisely.

## Keamanan Akses Konkuren dengan Race Detector Bersih

High-throughput admission gates face concurrent worker access. While `Filter` remains lock-free for immutable segments, `SyncFilter` wraps it with `sync.RWMutex`:

```go
type SyncFilter struct {
    mu sync.RWMutex
    f  *Filter
}
```

A benchmark test running 8 concurrent writers (1,000 keys each) and 16 concurrent readers (8,000 checks) passed clean under `go test -race`.

## Pengujian Mitigasi Cache Penetration pada Beban Serangan

`TestCachePenetrationMitigation` compares unprotected versus protected cache behavior under 1,000 attack queries:
- Unprotected Cache: 1,000 backend calls (100% penetration rate).
- Protected Cache: 0 backend calls (100% filtration).

Known valid keys (`user:101`) continued serving directly from memory with zero backend calls, confirming admission gate defense without blocking valid traffic.

## Tabel Komparasi Performa Komprehensif Arsitektur Lab

| Evaluation Scenario | Without Bloom Filter | With Bloom Filter | Efficiency Gain |
| --- | --- | --- | --- |
| LSM Disk Reads (10k absent) | 100,000 reads | 988 reads | 99.01% I/O reduction |
| Cache Penetration (5k attacks) | 5,000 hits | 0 hits | 100% backend protection |
| Lookup Probe Latency | 0 ns | ~150 ns | Sub-microsecond RAM access |
| Memory Overhead (10k items) | 0 bytes | 11.7 KB | Minimal (9.59 bits/element) |

Trading 11.7 KB of RAM and 150 ns CPU overhead eliminates milliseconds of disk seek time.

## Verifikasi Perilaku Filter Kosong: Penolakan Aman Awal

Boundary testing confirms that a newly allocated filter cleanly rejects all queries prior to key insertion:

```go
func TestEmptyFilter(t *testing.T) {
    bf := New(100, 0.01)
    if bf.Check([]byte("any-key")) {
        t.Errorf("Empty filter should not return true for arbitrary key")
    }
}
```

Zero-initialized bit arrays ensure the initial state behaves as a closed gate, preventing false hits on empty segments.

## Analisis Trade-off Memori versus False Positive Rate

Configuring epsilon requires balancing RAM allocation against tolerable disk and backend leakage:

| Target ε | Bits/Item | Probes k | RAM / 1M Keys | Workload Target |
| --- | --- | --- | --- | --- |
| 0.10 (10%) | 4.79 | 3 | 599 KB | Coarse disk pruning |
| 0.01 (1%) | 9.59 | 7 | 1.17 MB | Standard (RocksDB) |
| 0.001 (0.1%) | 14.38 | 10 | 1.75 MB | Strict cache gate |
| 0.0001 (0.01%) | 19.17 | 13 | 2.34 MB | High hash overhead |

The 1% rate represents the industrial sweet spot between memory consumption and filtering efficacy.

## Karakteristik Fungsi Hash: Efisiensi FNV-1a versus Keamanan

FNV-1a was selected because it is fast, allocation-free, and natively available in Go standard library (`hash/fnv`). It provides uniform dispersion for typical natural identifiers (UUIDs, auto-increment IDs).

However, FNV-1a is non-cryptographic. Adversarial clients can craft hash collisions to artificially saturate bit positions. Publicly exposed admission gates should adopt keyed hashes like SipHash.

## Batasan Fundamental: Ketiadaan Dukungan Operasi Delete

Standard Bloom filters cannot delete elements. Clearing a bit to remove key *A* risks clearing shared bits set by keys *B* and *C*, causing fatal false negatives.

In Lab 39, this limitation matches LSM-Tree segment immutability: filters are created once per segment and discarded during SSTable compaction. When dynamic deletions are mandatory, Counting Bloom Filters or Cuckoo Filters are required.

## Pertimbangan Operasional untuk Deployment di Lingkungan Produksi

Deploying probabilistic filters into production requires explicit operational safeguards:
1. **Capacity Saturation**: If element count exceeds design capacity *n*, bit saturation approaches 1.0, collapsing filtering power. Provision 20-30% headroom.
2. **Deterministic Serialization**: Persisting filters to SSTable headers requires recording *m*, *k*, and raw words with consistent hash seeding.
3. **Telemetry**: Monitor rejected vs passed query ratios to detect saturation drift early.

## Rekapitulasi Verifikasi Suite Pengujian dan Tolok Ukur

All assertions and race conditions verified clean across both packages:
- `TestNoFalseNegatives`: PASS (0 false negatives / 10k keys)
- `TestFalsePositiveRate`: PASS (~0.99% empirical vs 1.00% target)
- `TestOptimalSizing`: PASS (m=9,586, k=7 exact match)
- `TestSyncFilterConcurrency`: PASS (race-free with 24 goroutines)
- `TestLSMStoreWithAndWithoutFilter`: PASS (99.01% disk I/O reduction)
- `TestCachePenetrationMitigation`: PASS (100% penetration blocked)
- Benchmarks: ~500 ns Add, ~150 ns Check

## Panduan Reproduksi Mandiri dan Verifikasi Eksekusi Lab

To reproduce all lab benchmarks and test suites locally:

```sh
cd labs/39-bloom-filters
go test -v ./...                 # Run full unit & integration suite
go test -race ./...              # Verify race detector cleanliness
go run ./cmd/demo                # Execute 3-scenario CLI simulation
```

Requires Go 1.22+ standard library with zero external dependencies.

## Source

[Software Engineering Lab 39 — Bloom Filters](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/39-bloom-filters), commit `c977525a4a31ba811006df2bbc04fd031ee11765`. Code excerpts belong in their original package context.
