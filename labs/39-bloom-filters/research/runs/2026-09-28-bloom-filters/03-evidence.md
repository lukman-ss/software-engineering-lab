# 03-evidence.md

Research Date: 2026-09-28

---

## Evidence 1

**Claim:** A Bloom filter is a space-efficient probabilistic data structure that returns "definitely not in set" (no false negatives) or "possibly in set" (may have false positives). Elements can be added but not removed in the standard form.

**Evidence:** "In computing, a Bloom filter is a space-efficient probabilistic data structure, conceived by Burton Howard Bloom in 1970, that is used to test whether an element is a member of a set. False positive matches are possible, but false negatives are not – in other words, a query returns either 'possibly in set' or 'definitely not in set'."

**Source:** Wikipedia — Bloom filter  
**URL:** https://en.wikipedia.org/wiki/Bloom_filter  
**Publication Date:** Ongoing (based on Bloom 1970)  
**Confidence:** HIGH

**Corroborated By:**
- RocksDB Wiki (Source 3): "this bit array may be used to determine if the key *may exist* or *definitely does not exist* in the key set."
- Redis docs (Source 4): "A Bloom filter can guarantee the absence of an item from a set, but it can only give an estimation about its presence."

---

## Evidence 2

**Claim:** The false positive probability formula is approximately $\varepsilon = (1 - e^{-kn/m})^k$, where $m$ = number of bits, $k$ = hash functions, $n$ = items inserted.

**Evidence:** "The probability of false positives decreases as m (the number of bits in the array) increases, and increases as n (the number of inserted elements) increases… $\varepsilon \approx (1 - e^{-kn/m})^k$."

**Source:** Wikipedia — Bloom filter  
**URL:** https://en.wikipedia.org/wiki/Bloom_filter  
**Publication Date:** Ongoing (mathematical derivation from Bloom 1970)  
**Confidence:** HIGH

**Corroborated By:**
- Redis docs (Source 4): mentions the same formula indirectly with "The required number of bits per item, given the desired error_rate and the optimal number of hash functions, is `-ln(error_rate) / ln(2)^2`."

---

## Evidence 3

**Claim:** The optimal number of hash functions is $k = \frac{m}{n} \ln 2$. The minimum number of bits needed is $m = -\frac{n \ln \varepsilon}{(\ln 2)^2}$, giving approximately $-2.08 \ln \varepsilon$ bits per element.

**Evidence:** "Putting this constraint aside, for a given m and n, the value of k that minimizes the false positive probability is $k = \frac{m}{n} \ln 2$. This results in: $m = -\frac{n \ln \varepsilon}{\ln(2)^2}$. So the optimal number of bits per element is $m/n \approx -2.08 \ln \varepsilon$."

**Source:** Wikipedia — Bloom filter  
**URL:** https://en.wikipedia.org/wiki/Bloom_filter  
**Publication Date:** Ongoing (mathematical derivation)  
**Confidence:** HIGH

**Corroborated By:**
- Redis docs (Source 4): "The optimal number of hash functions is `ceil(-ln(error_rate) / ln(2))`. The required number of bits per item is `-ln(error_rate) / ln(2)^2`."
- Both sources agree: 1% error rate needs ~9.585 bits/item and 7 hash functions.

---

## Evidence 4

**Claim:** Memory comparison: for 1% false positive rate a Bloom filter needs ~9.6 bits per element regardless of element size. A Redis set for IP addresses requires ~320 bits per item.

**Evidence:**
- Wikipedia: "A Bloom filter with a 1% error and an optimal value of k, in contrast, requires only about 9.6 bits per element, regardless of the size of the elements."
- Redis docs: "For a set of IP addresses, for example, we would have around 40 bytes (320 bits) per item - considerably higher than the 19.170 bits we need for a Bloom filter with a 0.01% false positives rate."

**Source:** Wikipedia (Source 5) + Redis docs (Source 4)  
**URL:**  
- https://en.wikipedia.org/wiki/Bloom_filter  
- https://redis.io/docs/latest/develop/data-types/probabilistic/bloom-filter/  
**Confidence:** HIGH

**Notes:** At 1% FP rate: ~9.6 bits/item. At 0.01%: ~19.2 bits/item. Both sources independently calculate from the same formula.

---

## Evidence 5

**Claim:** For practical implementation, the Kirsch-Mitzenmacher double hashing technique allows generating $k$ independent hash positions using only two base hash functions: $g_i(x) = h_1(x) + i \cdot h_2(x) \pmod m$, without significant performance degradation.

**Evidence:** Wikipedia references "Dillinger & Manolios (2004b) show the effectiveness of deriving the k indices using enhanced double hashing and triple hashing." The Kirsch-Mitzenmacher (2006) result is referenced extensively in literature. Wikipedia states: "For a good hash function with a wide output, there should be little if any correlation between different bit-fields of such a hash, so this type of hash can be used to generate multiple 'different' hash functions by slicing its output into multiple bit fields."

**Source:** Wikipedia — Bloom filter (references Kirsch & Mitzenmacher 2006)  
**URL:** https://en.wikipedia.org/wiki/Bloom_filter  
**Confidence:** MEDIUM (the paper exists and is widely cited but the paper URL was not directly fetched; Wikipedia and engineering literature universally cite this result)

**Notes:** The paper "Less Hashing, Same Performance: Building a Better Bloom Filter" by Kirsch & Mitzenmacher proves asymptotically the FP rate of double hashing equals that of $k$ independent hash functions.

---

## Evidence 6

**Claim:** Google Bigtable, Apache HBase, Apache Cassandra, ScyllaDB, and PostgreSQL use Bloom filters to reduce disk lookups for non-existent rows/columns, considerably increasing database query performance.

**Evidence:** "Google Bigtable, Apache HBase, Apache Cassandra, ScyllaDB and PostgreSQL use Bloom filters to reduce the disk lookups for non-existent rows or columns. Avoiding costly disk lookups considerably increases the performance of a database query operation."

**Source:** Wikipedia — Bloom filter  
**URL:** https://en.wikipedia.org/wiki/Bloom_filter  
**Confidence:** HIGH (Tier 2 source cites known deployments; each system's documentation would confirm)

**Corroborated By:** RocksDB Wiki (Source 3): "In RocksDB, when the filter policy is set, every newly created SST file will contain a Bloom filter, which is used to determine if the file may contain the key we're looking for."

---

## Evidence 7

**Claim:** RocksDB implements Bloom filters per SST file using ~10 bits per key by default. The "Full Filter" format limits probe bits to a single CPU cache line for cache efficiency. False positive rate at 10 bits/key is approximately 1%.

**Evidence:**
- "This generates filters using about 10 bits of space per key, which works well for many workloads."
- "9.9 bits per key (1% false positive rate) is 99% as effective as 100 bits per key."
- "Full filter limits the probe bits for a key to be all within the same CPU cache line. This ensures fast lookups and maximizes CPU cache effectiveness by limiting the CPU cache misses to one per key (per filter)."

**Source:** RocksDB Wiki — Bloom Filter  
**URL:** https://github.com/facebook/rocksdb/wiki/RocksDB-Bloom-Filter  
**Publication Date:** 2021-09-01  
**Confidence:** HIGH

**Notes:** RocksDB's 2021 update (v6.6+, format_version=5) migrated to 64-bit hashing, allowing billions of keys per filter without hash collision degradation.

---

## Evidence 8

**Claim:** Akamai Technologies uses Bloom filters to prevent "one-hit-wonder" web objects from being stored in disk cache. Approximately three-quarters of their cached content was requested only once.

**Evidence:** "The servers of Akamai Technologies, a content delivery provider, use Bloom filters to prevent 'one-hit-wonders' from being stored in its disk caches. One-hit-wonders are web objects requested by users just once, something that Akamai found applied to nearly three-quarters of their caching infrastructure."

**Source:** Wikipedia — Bloom filter (cites Maggs & Sitaraman 2015)  
**URL:** https://en.wikipedia.org/wiki/Bloom_filter  
**Confidence:** MEDIUM (Wikipedia secondary citation; primary = Maggs & Sitaraman 2015 paper)

---

## Evidence 9

**Claim:** Redis implements Bloom filter as a module (`RedisBloom`) with commands BF.RESERVE, BF.ADD, BF.EXISTS, BF.MADD, BF.MEXISTS, all O(k) complexity. Supports auto-scaling sub-filters when capacity exceeded.

**Evidence:**
- "Insertion in a Bloom filter is O(K), where k is the number of hash functions."
- "BF.RESERVE {key} {error_rate} {capacity} [EXPANSION expansion] [NONSCALING]"
- "when capacity is reached, an additional sub-filter will be created… The size of the new sub-filter is the size of the last sub-filter multiplied by EXPANSION. The default expansion value is 2."

**Source:** Redis Bloom Filter Documentation  
**URL:** https://redis.io/docs/latest/develop/data-types/probabilistic/bloom-filter/  
**Confidence:** HIGH

---

## Evidence 10

**Claim:** A standard Bloom filter cannot delete elements. Counting Bloom Filters extend the structure with multi-bit counters per slot to support deletion, at the cost of ~4x memory overhead.

**Evidence:** "Removing an element from this simple Bloom filter is impossible because there is no way to tell which of the k bits it maps to should be cleared."

**Source:** Wikipedia — Bloom filter  
**URL:** https://en.wikipedia.org/wiki/Bloom_filter  
**Confidence:** HIGH

**Notes:** The Counting Bloom Filter variant uses counters instead of single bits; deletion decrements the counter. Mentioned in Wikipedia Extensions section.

---

## Evidence 11

**Claim:** Bloom filters use approximately 44% more space than a theoretically optimal data structure for set membership. The lower bound for any such structure is $\log_2(1/\varepsilon)$ bits per key, while Bloom filters use $1.44 \log_2(1/\varepsilon)$ bits per key.

**Evidence:** "Classic Bloom filters use $1.44 \log_2(1/\varepsilon)$ bits of space per inserted key... the space that is strictly necessary for any data structure playing the same role as a Bloom filter is only $\log_2(1/\varepsilon)$ per key. Hence Bloom filters use 44% more space than an equivalent optimal data structure."

**Source:** Wikipedia — Bloom filter  
**URL:** https://en.wikipedia.org/wiki/Bloom_filter  
**Confidence:** HIGH

**Notes:** Quotient filters and Cuckoo filters can approach the theoretical optimum.

---

## Evidence 12

**Claim:** Google Chrome previously used a Bloom filter to screen URLs against a malicious URL database; positive hits triggered a full server-side check.

**Evidence:** "The Google Chrome web browser previously used a Bloom filter to identify malicious URLs. Any URL was first checked against a local Bloom filter, and only if the Bloom filter returned a positive result was a full check of the URL performed."

**Source:** Wikipedia — Bloom filter  
**URL:** https://en.wikipedia.org/wiki/Bloom_filter  
**Confidence:** HIGH (well-documented engineering decision by Chrome team, circa 2011)

---

## Evidence 13

**Claim:** RocksDB's Ribbon filter (available since v6.15.0) provides 1% FP rate with only ~7 bits/key vs ~10 bits/key for Bloom filter — saving ~30% memory — at the cost of ~3-4x more CPU during filter construction.

**Evidence:** "A new Bloom filter alternative is available as a drop-in replacement (since version 6.15.0), saving about 30% of Bloom filter space (most importantly, memory) but using about 3-4x as much CPU on filters. `rocksdb::NewRibbonFilterPolicy(9.9)` has the same 1% FP rate as Bloom but only uses around 7 bits per key."

**Source:** RocksDB Wiki — Bloom Filter  
**URL:** https://github.com/facebook/rocksdb/wiki/RocksDB-Bloom-Filter  
**Publication Date:** 2021-09-01  
**Confidence:** HIGH

---

## Evidence 14

**Claim:** Cache Penetration attack: repeated queries for non-existent keys bypass cache (cache miss) and hit the database directly, exhausting connection pools.

**Evidence (Interpretation/Design Pattern):** Redis docs (Source 4) describes use case: "Has this username/email/domain name/slug already been used? A new user types in the desired username. The app checks if the username exists in the Bloom filter. If no → user is created. If yes → check the main database or reject." This indirectly demonstrates Bloom filter as cache-penetration guard.

Redis docs also note: "A Bloom filter can guarantee the absence of an item from a set… There are many cases out there where a negative answer will prevent more costly operations."

**Source:** Redis Bloom Filter Documentation (Source 4)  
**URL:** https://redis.io/docs/latest/develop/data-types/probabilistic/bloom-filter/  
**Confidence:** MEDIUM (the use case logic is sound and from authoritative source; explicit "cache penetration" terminology not used in this doc)

**Notes:** Cache penetration as a named attack pattern is widely documented in cloud/system design literature, consistent with this evidence.

---

## Evidence 15

**Claim:** Memory estimate for 10 million items at 1% FP rate: requires approximately 11.93 MB (≈ 95.85 Mbits), matching the lab's claim of "~12 MB."

**Evidence (Derived from formula):**  
- Formula: $m = n \cdot \frac{-\ln \varepsilon}{(\ln 2)^2}$  
- $n = 10{,}000{,}000$, $\varepsilon = 0.01$  
- $m = 10^7 \cdot \frac{-\ln(0.01)}{(\ln 2)^2} = 10^7 \cdot \frac{4.6052}{0.4805} = 10^7 \cdot 9.585 = 95{,}850{,}584 \text{ bits} \approx 11.93 \text{ MB}$

**Source:** Formula derived from Wikipedia (Source 5) and Redis docs (Source 4) — both independently give the same formula  
**Confidence:** HIGH

**Notes:** Lab topic states "Bloom Filter (FP rate ~1%) hanya butuh ~12 MB RAM!" — this is mathematically correct per the formula.
