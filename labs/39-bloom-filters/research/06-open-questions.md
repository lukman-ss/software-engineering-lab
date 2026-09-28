
# Open Questions

- **Dynamic resizing**: How do modern production systems handle bloom‑filter growth when the element count exceeds the initially reserved capacity? (Scalable Bloom filters are mentioned but not detailed.)
- **Exact false‑positive benchmarks**: What empirical FP rates have been observed in large‑scale deployments (e.g., Cassandra, RocksDB) under realistic key distributions?
- **Hash‑function independence**: In practice, many implementations derive multiple hash functions from a single base hash (e.g., double‑hashing). How much does this relax the independence assumption and affect FP rate?
- **Counting Bloom vs. Cuckoo filter choice**: Under what workload characteristics (insert/delete ratio, key distribution) should an engineer prefer one over the other?
- **Cryptographic hash usage**: Are there scenarios where a cryptographic hash (e.g., SHA‑256 truncated) is justified in a bloom filter, despite the speed penalty?
- **Parallel / GPU implementations**: How do GPU‑accelerated bloom‑filter checks scale compared to CPU implementations?
- **Bloom filter in SQL databases**: Which RDBMS (PostgreSQL, MySQL) include built‑in bloom‑filter support, and how is it exposed to users?
