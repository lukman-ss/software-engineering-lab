# Engineering Code Audit

## Finding 1: Core Mathematical Formulas & Bit Array Sizing

Location: `internal/bloom/bloom.go:85-102`
Claimed Behavior: Sizing adheres to standard Bloom filter formulas:
$m = \lceil -n \ln(\epsilon) / (\ln 2)^2 \rceil$, $k = \max(1, \lfloor (m/n)\ln 2 \rceil)$.
Observed Implementation:
- `optimalM` uses `math.Ceil(-float64(n) * math.Log(epsilon) / (math.Log(2) * math.Log(2)))`.
- `optimalK` computes `math.Round(float64(m) / float64(n) * math.Log(2))` and clamps minimum to 1.
- `words := (m + 63) / 64` allocates appropriate `uint64` slices.
Assessment: PASS
Severity: LOW
Notes: Sizing mathematically matches research specification.

## Finding 2: Double Hashing Derivation

Location: `internal/bloom/bloom.go:61-83`
Claimed Behavior: Kirsch-Mitzenmacher double-hashing technique derives $k$ probe indices from 2 independent hashes with no asymptotic degradation in false-positive rate.
Observed Implementation:
- `baseHashes` uses standard library `hash/fnv` (FNV-1a 64-bit) for $h_1$ and 64-bit rotate-XOR ($h_1 \oplus (h_1 \gg 17 \mid h_1 \ll 47)$) for $h_2$.
- $h_2 == 0$ fallback correctly avoids zero step size.
- `doubleHash(h1, h2, i, m)` evaluates `(h1 + i*h2) % m`.
Assessment: PASS
Severity: LOW
Notes: Empirical FP tests verify that probe dispersion meets theoretical guarantees without external libraries.

## Finding 3: Thread-Safety & Concurrency Management

Location: `internal/bloom/bloom.go:104-134`, `internal/store/lsm.go:34-84`, `internal/store/cache.go:18-89`
Claimed Behavior: Safe concurrent access under readers and writers with race detection compliance.
Observed Implementation:
- `SyncFilter` wraps `Filter` with `sync.RWMutex` where `Add` acquires `Lock` and `Check`/`M`/`K` acquire `RLock`.
- `LSMStore` uses `sync.RWMutex` for segment slice mutations/reads and `sync/atomic` for `diskReads`.
- `Cache` uses `sync.RWMutex` for map access, `sync/atomic` for `backendCalls`, and delegates filter operations to `SyncFilter`.
Assessment: PASS
Severity: LOW
Notes: `go test -race ./...` completes cleanly with zero race warnings.

## Finding 4: Segment Skipping & I/O Reduction

Location: `internal/store/lsm.go:56-75`
Claimed Behavior: If segment Bloom filter returns false, segment data lookup is bypassed and disk read counter is not incremented.
Observed Implementation:
- `if seg.filter != nil && !seg.filter.Check(keyBytes)` triggers `continue`, skipping simulated disk read (`atomic.AddUint64(&s.diskReads, 1)`).
- If filter is nil or returns true, lookup proceeds and disk counter increments.
Assessment: PASS
Severity: LOW
Notes: Correctly models LSM read path and SSTable filter skipping.

## Finding 5: Cache Penetration Admission Gate

Location: `internal/store/cache.go:54-79`
Claimed Behavior: Queries for absent keys blocked by Bloom filter return false immediately without querying the backend or incrementing backend call count.
Observed Implementation:
- In `Get`, `c.filter.Check([]byte(key))` returning `false` returns `"", false` immediately.
- Backend fetch `c.backend.Fetch(key)` and `atomic.AddUint64(&c.backendCalls, 1)` are strictly avoided on negative filter check.
Assessment: PASS
Severity: LOW
Notes: Accurately reflects cache penetration defense patterns.
