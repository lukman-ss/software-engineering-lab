
# Research Plan

**Topic**: Bloom Filters – Probabilistic Data Structures for Reducing Disk I/O and Preventing Cache Penetration

**Objective**: Provide a clear, concise overview of Bloom filter fundamentals, applications, and performance characteristics.

**Research Questions**
1. What is a Bloom filter and how does it work?
2. What are false positives, false negatives, and the false‑positive probability formula?
3. How does filter size, number of hash functions, and element count affect space usage and error rate?
4. Which real‑world systems use Bloom filters (e.g., caches, databases, distributed stores, search engines)?
5. How do Bloom filters compare to alternative probabilistic structures (e.g., Cuckoo filters, Quotient filters)?

**Search Strategy**
- Primary tier: Original Bloom 1970 paper, FNV‑FNV, RocksDB, Percolator, and BitFunnel docs.
- Secondary tier: Wikipedia, academic surveys, and recent conference papers (SIGIR 2019 LSM survey, etc.).
- Use websearch for "Bloom filter false positive probability" and "Bloom filter use cases".

**Expected Primary Sources**
- Bloom, B. H. (1970). *Space/Time Trade‑offs in Hash Coding with Allowable Errors*.
- O'Neil et al. (1996). *The log‑structured merge‑tree (LSM‑tree)* – LSM uses Bloom filters.
- Google Percolator (USENIX 2010) – describes Bloom filter usage for cache penetration.
- RocksDB documentation – Bloom filter configuration.
- BitFunnel Wikipedia – Bloom filter application in search indexing.
- RFC/whitepaper on Bloom filters in Redis, Cassandra, etc.

**Deliverables**
- `01-plan.md` (research plan – already drafted).
- `02-sources.md` (list of sources with metadata).
- `03-evidence.md` (claims, evidence, confidence, corroboration).
- `04-contradictions.md` (none expected; placeholder).
- `05-report.md` (structured research report with executive summary, findings, and limitations).
- `06-open-questions.md` (any gaps or ambiguous findings).

**Timeline**
- Day 1: Gather primary sources, verify URLs, and store metadata.
- Day 2: Extract key claims, calculate false‑positive formulas, and capture performance numbers.
- Day 3: Write evidence, check cross‑sources, and draft report.
- Day 4: Review, finalize, and populate open‑questions.

**Notes**
- All claims will be backed by at least two independent sources where possible.
- Confidence levels: HIGH for Bloom‑1970 and O'Neil‑1996; MEDIUM for industry implementations; LOW for speculative performance numbers.
