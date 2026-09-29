# Code Audit

Target Lab: labs/39-bloom-filters

## Finding 1

Location: `internal/bloom/bloom.go:22-31, 86-102`
Claimed Behavior: Sizing calculations match classical Bloom filter formulas ($m = \lceil -n \ln(\epsilon) / (\ln 2)^2 \rceil$, $k = \text{round}((m/n) \ln 2)$) and allocate sufficient 64-bit words.
Observed Implementation: `optimalM` uses `math.Ceil(-float64(n) * math.Log(epsilon) / (math.Log(2)*math.Log(2)))`. `optimalK` computes `(m/n)*ln(2)` rounded, clamped to $\ge 1$. Word count is `(m + 63) / 64`.
Assessment: PASS
Severity: LOW
Notes: Formula matches the literature accurately.

## Finding 2

Location: `internal/bloom/bloom.go:65-77, 80-83`
Claimed Behavior: Kirsch-Mitzenmacher double hashing scheme generates $k$ hash probes using two base hashes without requiring $k$ distinct hash functions.
Observed Implementation: `baseHashes` computes 64-bit FNV-1a for $h_1$ and derives $h_2$ using 64-bit rotate-XOR permute ($h_1 \oplus (h_1 \gg 17 \mid h_1 \ll 47)$). `doubleHash` computes $(h_1 + i \cdot h_2) \pmod m$.
Assessment: PASS
Severity: LOW
Notes: Valid implementation of Kirsch-Mitzenmacher optimization using standard library `hash/fnv`.

## Finding 3

Location: `internal/bloom/bloom.go:34-53`
Claimed Behavior: Add sets bits at $k$ probe locations; Check verifies all $k$ probe locations are set to 1.
Observed Implementation: Add sets `f.bits[pos/64] |= 1 << (pos % 64)`. Check returns `false` early if `f.bits[pos/64] & (1 << (pos % 64)) == 0`. If all $k$ bits are set, returns `true`.
Assessment: PASS
Severity: LOW
Notes: Bitwise indexing and early-exit membership testing are correct.

## Finding 4

Location: `internal/bloom/bloom.go:104-134`
Claimed Behavior: `SyncFilter` provides thread-safe access to underlying `Filter`.
Observed Implementation: `SyncFilter` embeds `sync.RWMutex`. `Add` acquires `sf.mu.Lock()`, `Check` acquires `sf.mu.RLock()`. `M()` and `K()` also acquire read locks.
Assessment: PASS
Severity: LOW
Notes: Synchronization is complete and free of data races.

## Finding 5

Location: `internal/store/lsm.go:56-74`
Claimed Behavior: LSM-Tree store skips disk read (simulated atomic counter) if segment Bloom filter returns false on lookup.
Observed Implementation: `Get(key)` iterates segments in reverse chronological order. If `seg.filter != nil && !seg.filter.Check(keyBytes)`, the segment is bypassed. Otherwise `atomic.AddUint64(&s.diskReads, 1)` is incremented.
Assessment: PASS
Severity: LOW
Notes: Accurate simulation of LSM SSTable bloom filter skip optimization.

## Finding 6

Location: `internal/store/cache.go:54-79`
Claimed Behavior: Cache checks Bloom admission filter before delegating cache misses to backend database.
Observed Implementation: If `c.filter != nil && !c.filter.Check([]byte(key))`, `Get` returns `"", false` immediately without calling `c.backend.Fetch` or incrementing `backendCalls`.
Assessment: PASS
Severity: LOW
Notes: Correct modeling of anti-penetration admission gate.
