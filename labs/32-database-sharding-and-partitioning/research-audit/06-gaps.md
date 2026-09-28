# Research Gap Analysis

## Gap 1
Type: MISSING_CASE
Severity: LOW
Location: `research/05-report.md` (Finding 3), `research/06-open-questions.md`
Problem: While consistent hashing ring properties are thoroughly explained, virtual node count tuning strategies (e.g., standard ratio of 100-300 vnodes per physical node to control variance) are mentioned only qualitatively.
Required Revision: Detail specific virtual node configuration heuristics for balancing distribution variance against memory overhead.
Can Be Approved Without Fix: YES

---

## Gap 2
Type: WEAK_SOURCE
Severity: LOW
Location: `research/06-open-questions.md`
Problem: The transition threshold from vertical partitioning to distributed sharding is noted as lacking cross-vendor benchmark standards, which is accurately flagged in open questions.
Required Revision: Document empirical sizing boundaries commonly used in practice (e.g., table size exceeding RAM / buffer pool working set, write IOPS saturation).
Can Be Approved Without Fix: YES

---

## Gap 3
Type: SCOPE_ERROR
Severity: LOW
Location: `research/05-report.md` (Finding 4)
Problem: Cross-shard distributed transaction isolation degradation (e.g., read phenomena during 2PC phase transitions) is noted conceptually via Vitess TwoPC but not detailed across distributed concurrency control protocols (e.g., Percolator, Spanner TrueTime, MVCC).
Required Revision: Acknowledge distributed isolation levels in the subsequent lab design.
Can Be Approved Without Fix: YES
