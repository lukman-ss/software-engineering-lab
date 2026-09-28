# Research Plan: Bloom Filters

## Research Topic
Bloom Filters — Probabilistic Data Structure for Disk I/O Reduction & Cache Penetration Prevention.

## Objective
Investigate Bloom Filter core mathematics, optimal bit array ($m$) and hash count ($k$) formulas, false positive rate dynamics, MurmurHash3/FNV hash strategies (Kirsch-Mitzenmacher optimization), real-world systems use cases (LSM-tree in RocksDB/Cassandra, Redis Bloom, Cache Penetration mitigation), and practical memory trade-offs.

## Research Questions
1. What are the mathematical formulas for optimal array size ($m$) and number of hash functions ($k$) given item count ($n$) and target false positive probability ($p$)?
2. How does the false positive rate behave when elements exceed expected capacity $n$?
3. How does Kirsch-Mitzenmacher optimization allow generating $k$ hash values using only 2 hash functions ($h_1(x) + i \cdot h_2(x) \pmod m$)?
4. How do LSM-Tree databases (RocksDB, Apache Cassandra, LevelDB) and Caching layers use Bloom Filters to prevent disk read amplification and cache penetration?
5. What are the limitations of standard Bloom Filters (e.g., deletion impossibility) and what alternatives exist (Counting Bloom Filter, Cuckoo Filter)?

## Search Strategy
- Target Tier 1 authoritative computer science papers (Burton H. Bloom 1970, Kirsch & Mitzenmacher 2006).
- Target Tier 1 official documentation (RocksDB wiki, RedisBloom docs, Apache Cassandra docs, Google Guava BloomFilter docs).
- Cross-check false positive rate math and double-hashing technique proofs.

## Expected Primary Sources
- Burton H. Bloom (1970) "Space/Time Trade-offs in Hash Coding with Allowable Errors", CACM.
- Adam Kirsch and Michael Mitzenmacher (2006) "Less Hashing, Same Performance: Building a Better Bloom Filter".
- RocksDB Wiki: Leveled Compaction / Bloom Filter architecture.
- RedisBloom & Apache Cassandra Documentation.

## Risks / Unknowns
- Numerical precision variations across articles for optimal $k = (m/n) \ln 2$.
- Differences between standard Bloom filters and Blocked/Spatial Bloom filters used in hardware/CPU cache optimization.
