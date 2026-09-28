# 02-sources.md

## Source 1

Title: Space/Time Trade-offs in Hash Coding with Allowable Errors  
Publisher: Communications of the ACM (Vol 13, No 7)  
URL: https://dl.acm.org/doi/10.1145/362686.362692  
Published: 1970-07-01  
Accessed: 2026-09-28  
Source Tier: Tier 1  
Relevance: Primary paper establishing Bloom Filters, bit array structure, and allowable error trade-offs.

## Source 2

Title: Less Hashing, Same Performance: Building a Better Bloom Filter  
Publisher: Harvard University / ESA 2006  
URL: https://www.eecs.harvard.edu/~michaelm/postscripts/rsa2006.pdf  
Published: 2006-09-01  
Accessed: 2026-09-28  
Source Tier: Tier 1  
Relevance: Kirsch & Mitzenmacher paper proving double hashing technique $g_i(x) = h_1(x) + i \cdot h_2(x) \pmod m$ yields asymptotic equivalent false positive rate.

## Source 3

Title: RocksDB Bloom Filter Documentation  
Publisher: Meta / Facebook RocksDB Wiki  
URL: https://github.com/facebook/rocksdb/wiki/RocksDB-Bloom-Filter  
Published: 2021-09-01  
Accessed: 2026-09-28  
Source Tier: Tier 1  
Relevance: Industrial LSM-Tree usage, cache-local full bloom filters, Ribbon filters, and bits-per-key tuning in storage engine.

## Source 4

Title: Redis Bloom Filter Data Type Documentation  
Publisher: Redis Ltd.  
URL: https://redis.io/docs/latest/develop/data-types/probabilistic/bloom-filter/  
Published: 2024-05-01  
Accessed: 2026-09-28  
Source Tier: Tier 1  
Relevance: Production RedisBloom documentation detailing memory formulas, sub-filter scaling, command complexity, and cache penetration use-cases.

## Source 5

Title: Bloom filter — Wikipedia  
Publisher: Wikimedia Foundation  
URL: https://en.wikipedia.org/wiki/Bloom_filter  
Published: 2024-05-01  
Accessed: 2026-09-28  
Source Tier: Tier 2  
Relevance: Mathematical derivations of false positive probability, optimal $k$ and $m$, item count estimation formulas.
