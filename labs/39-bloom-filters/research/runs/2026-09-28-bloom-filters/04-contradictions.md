# 04-contradictions.md

## Disagreements / Variations Discovered

### Disagreement 1: Standard FP Rate Formula Approximation vs. Strict Independence Assumption

**SOURCE A:** Wikipedia (Standard Analysis):  
Formula: $\varepsilon = (1 - e^{-kn/m})^k$  
Assumes independent probabilities for each bit tested.

**SOURCE B:** Goel and Gupta (As cited in Wikipedia):  
Formula: $\varepsilon \leq \left(1 - e^{-\frac{k(n + 0.5)}{m - 1}}\right)^k$  
Proves a rigorous upper bound without the independence assumption.

**ASSESSMENT:** The difference between the standard approximation and Goel-Gupta upper bound is less than half an extra element and at most 1 less bit. For engineering purposes ($n \ge 10^5$), the standard approximation $\varepsilon = (1 - e^{-kn/m})^k$ is universally used in industry libraries (Google Guava, RedisBloom, RocksDB).

---

### Disagreement 2: Standard Bloom Filter vs. Cache-Efficient / Blocked Bloom Filter Performance

**SOURCE A:** Standard Bloom Filter (Bloom 1970, RedisBloom):  
Probe bits are uniformly distributed across the entire $m$-bit array. Requires $k$ cache line accesses per query when $m$ exceeds CPU cache size.

**SOURCE B:** RocksDB Full Filter (Cache-Local Bloom Filter, Putze et al. 2007):  
Probes are restricted to a single CPU cache line (64 bytes / 512 bits) to avoid CPU cache misses.

**ASSESSMENT:** RocksDB's block-local probe increases false positive rate slightly (factor of 1.13x higher due to key clustering variance per shard) but speeds up query latency significantly by guaranteeing at most 1 CPU cache miss per check. This is an explicit trade-off between strict FP optimality and hardware CPU cache locality.

---

### Disagreement 3: Bits Per Key for 1% False Positive Rate

**SOURCE A:** Theoretical Optimal Formula:  
$m/n = -\ln(0.01) / (\ln 2)^2 \approx 9.585$ bits per key.

**SOURCE B:** RocksDB Documentation (Practical Implementation):  
Recommends ~9.9 to 10 bits per key for 1% FP rate (`NewBloomFilterPolicy(10)` or `9.9`).

**ASSESSMENT:** The 0.3-0.4 bit per key excess in RocksDB accounts for integer rounding, cache line alignment padding, and hash distribution variance. The theoretical optimal assumes fractional $k = m/n \ln 2 \approx 6.64$, whereas real implementations must round $k$ to an integer ($k=7$).

---

## Conclusion
No material contradictions discovered regarding core invariants:
- Zero false negatives (100% accuracy for absence).
- Non-zero false positive rate governed by $m$, $n$, $k$.
- Deletion impossibility without counters/cuckoo variants.
