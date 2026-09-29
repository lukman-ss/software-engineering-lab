# Code Audit

## Finding 1

Location: `internal/bloom/bloom.go:22-31, 86-102`
Claimed Behavior: Optimal $m$ and $k$ computation adhering to Bloom (1970) formulas.
Observed Implementation: `optimalM` computes `uint(math.Ceil(-float64(n) * math.Log(epsilon) / (math.Log(2) * math.Log(2))))`. `optimalK` computes `uint(math.Round(float64(m) / float64(n) * math.Log(2)))` with a lower bound guard `k >= 1`. Word allocation allocates `(m + 63) / 64` `uint64` words.
Assessment: PASS
Severity: LOW
Notes: Mathematical implementation matches research report equation (1) and (2) exactly. Word sizing cleanly handles ceiling division for 64-bit alignment.

## Finding 2

Location: `internal/bloom/bloom.go:34-53, 65-84`
Claimed Behavior: Kirsch-Mitzenmacher double hashing $h_i(x) = (h_1(x) + i \cdot h_2(x)) \pmod m$ using standard library hashes, guaranteeing zero false negatives and bit indexing without bounds overflow.
Observed Implementation: Uses `hash/fnv` (FNV-1a 64-bit) for $h_1$ and a 64-bit circular rotate-XOR for $h_2$. Computes `pos = (h1 + i*h2) % m`. Bit manipulation correctly uses `pos / 64` for word indexing and `1 << (pos % 64)` for mask.
Assessment: PASS
Severity: LOW
Notes: No bit indexing out-of-bounds error possible because `pos < m` and bit array length is `(m + 63) / 64`.

## Finding 3

Location: `internal/bloom/bloom.go:104-134`
Claimed Behavior: Concurrent thread safety for Bloom filter reads and writes.
Observed Implementation: `SyncFilter` wraps `*Filter` with `sync.RWMutex`. `Add` acquires full `Lock()`, while `Check`, `M()`, and `K()` acquire `RLock()`.
Assessment: PASS
Severity: LOW
Notes: Concurrency safety verified under `-race` with 8 parallel writer goroutines and 16 reader goroutines.

## Finding 4

Location: `internal/store/lsm.go:56-74`
Claimed Behavior: LSM segments skip disk lookups when segment Bloom filter returns false, and increment atomic `diskReads` when filter returns true or is absent.
Observed Implementation: Iterates segments in reverse (newest to oldest). If `seg.filter != nil && !seg.filter.Check(keyBytes)` is met, the loop continues without incrementing `diskReads`. When candidate found or filter returns true, `atomic.AddUint64(&s.diskReads, 1)` is called and segment map is checked.
Assessment: PASS
Severity: LOW
Notes: Clean simulation model reflecting production LSM-tree segment probe patterns.

## Finding 5

Location: `internal/store/cache.go:49-79`
Claimed Behavior: Cache gate rejects absent keys before checking map or calling backend if Bloom filter indicates key is definitely absent.
Observed Implementation: `c.filter.Check([]byte(key))` returns false triggers early exit without backend increment or backend call. Cache miss after positive filter result calls backend and populates both map and filter atomically/safely.
Assessment: PASS
Severity: LOW
Notes: Synchronization protects `data` map under `sync.RWMutex` while filter calls `SyncFilter.Add` / `SyncFilter.Check`.

## Finding 6

Location: `go.mod`
Claimed Behavior: Zero external dependencies (standard library only).
Observed Implementation: `go.mod` contains only `module bloomfilters` and `go 1.22`. No `require` directives.
Assessment: PASS
Severity: LOW
Notes: Verified compliant with stdlib-only policy.
