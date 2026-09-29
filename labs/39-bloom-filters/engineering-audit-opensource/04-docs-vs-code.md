# Docs vs Code Audit

Target Lab: labs/39-bloom-filters

## Comparisons

### 1. Mathematical Formulas & Sizing
- **Claim in README & Design**:
  - $m = \lceil -n \cdot \ln(\epsilon) / (\ln 2)^2 \rceil$
  - $k = \text{round}((m/n) \cdot \ln 2)$
  - Bits/elem at $\epsilon = 0.01 \approx 9.59$ bits ($m = 95,851$ for $n = 10,000$).
- **Implementation in `internal/bloom/bloom.go`**:
  - Exactly implements these formulas in `optimalM` and `optimalK`.
  - For $n = 10,000, \epsilon = 0.01$, `bf.M()` returns $95,851$, `bf.K()` returns $7$.
- **Assessment**: PASS.

### 2. False Negative & False Positive Behavior
- **Claim**: 0% false negatives; false positive rate approximately $\epsilon$ (1.00%).
- **Observed**:
  - `TestNoFalseNegatives`: 0 false negatives out of 10,000 keys.
  - `TestFalsePositiveRate`: 489 false positives out of 50,000 queries (0.978% - 0.996%).
  - Demo: 0 false negatives across 10,000 keys; 498 false positives across 50,000 keys (0.996%).
- **Assessment**: PASS.

### 3. LSM-Tree Disk I/O Reduction
- **Claim**: Absent-key lookups across segments reduce disk reads by $> 99\%$.
- **Observed in Demo**:
  - Reads without filter: 100,000
  - Reads with filter: 988
  - Reduction: 99.01%
- **Assessment**: PASS.

### 4. Cache Penetration Mitigation
- **Claim**: Unprotected cache allows 100% backend hits on absent keys; Bloom filter protected cache eliminates or drops penetration to $\approx 0\%$.
- **Observed in Demo & Tests**:
  - Unprotected: 5,000 hits out of 5,000 malicious queries (100%).
  - Protected: 0 hits (0.00%).
- **Assessment**: PASS.

### 5. Dependency Claims
- **Claim**: Standard library only (`hash/fnv`, `math`, `sync`). No external dependencies.
- **Observed `go.mod`**:
  - Module `bloomfilters`, `go 1.22`, zero `require` directives.
- **Assessment**: PASS.

## Discrepancies Found
None. Documentation accurately describes implementation, tests, and demo output.
