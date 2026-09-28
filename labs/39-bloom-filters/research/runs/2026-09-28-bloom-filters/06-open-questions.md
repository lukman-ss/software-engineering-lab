# 06-open-questions.md

## Unanswered Questions

### 1. Kirsch-Mitzenmacher Primary Paper Verification
The Kirsch & Mitzenmacher (2006) paper was not directly fetched. The claim that $g_i(x) = h_1(x) + i \cdot h_2(x) \pmod m$ produces asymptotically equivalent FP rate was verified only via Wikipedia and engineering secondary sources.

**Recommended next step:** Fetch and inspect https://www.eecs.harvard.edu/~michaelm/postscripts/rsa2006.pdf directly.

---

### 2. MurmurHash3 vs. FNV-1a Properties
The lab exercise specifies using MurmurHash3 or FNV as the underlying hash. The research has not verified:
- Whether MurmurHash3 (128-bit output) provides sufficient avalanche effect and uniformity for Bloom Filter probing.
- Whether FNV-1a generates adequately independent bit fields when seeded to derive $k$ hash positions.

**Recommended next step:** Fetch MurmurHash3 original documentation (Austin Appleby, 2011) and benchmark-level analysis.

---

### 3. Bloom Filter 1970 Original Paper Access
Burton H. Bloom's original 1970 CACM paper was not directly fetched (paywalled at ACM). Wikipedia and all downstream sources summarize it, but direct verification of the exact formulation was not done.

**Status:** NOT VERIFIED (primary source)  
**Recommended next step:** Fetch via https://dl.acm.org/doi/10.1145/362686.362692 or Open Access mirror.

---

### 4. False Positive Rate Behavior Beyond Design Capacity
The research documents that false positive rate rises as $n$ exceeds design capacity, but does not establish the empirical inflection point (e.g., at $n = 1.5 \times n_\text{design}$, what is $\varepsilon$?).

**Recommended calculation:**
- If $n = 1.5 n_\text{design}$ and $k = 7$: $\varepsilon = (1 - e^{-7 \cdot 1.5 n / m})^7$ — substituting actual values.

---

### 5. Cuckoo Filter vs. Bloom Filter Practical Performance Trade-offs
The research documents that Cuckoo Filters support deletion and approach space optimality. However, no primary source was fetched comparing real-world throughput, insertion worst-case latency, or eviction chain behavior under high load factors.

**Recommended next step:** Fetch Cuckoo Filter paper (Fan et al., 2014): "Cuckoo Filter: Practically Better Than Bloom."

---

### 6. Counting Bloom Filter Memory Overhead Quantification
Counting Bloom Filters support deletion via counter slots. The research notes this incurs higher memory but does not compute the exact overhead (e.g., for 4-bit counters: 4x bits vs 1 bit per slot).

**Recommended next step:** Verify the counter overflow probability for $k=7$, $n = 10^6$ elements.

---

### 7. Distributed / Partitioned Bloom Filters
Systems spanning multiple nodes (e.g., distributed cache, sharded databases) may use partitioned Bloom Filters. No research was done on inter-node synchronization of Bloom state, false positive rate impact of partitioning, or relevance to systems like Apache Cassandra's per-SSTable bloom filters.

**Recommended next step:** Review Cassandra source or documentation on per-SSTable Bloom Filter configuration.
