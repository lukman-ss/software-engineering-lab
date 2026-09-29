# Research Gap Analysis

## Gap 1
Type: WEAK_SOURCE
Severity: LOW
Location: `research/02-sources.md:Source 9`, `research/03-evidence.md:Evidence 5`
Problem: The foundational paper for consistent hashing (Karger et al. 1997) is gated behind ACM paywall; research relied on the abstract/landing page corroborated by Wikipedia's mathematical restatement.
Required Revision: None required for architectural validity, as the mathematical properties ($1/n$ remap, virtual nodes) are standard and correctly stated.
Can Be Approved Without Fix: YES

## Gap 2
Type: MISSING_CASE
Severity: LOW
Location: `research/06-open-questions.md:Unanswered Questions`
Problem: Quantitative threshold for transitioning from single-instance table partitioning to multi-instance physical sharding is not strictly quantified with universal benchmarks across engines.
Required Revision: None. This is inherently workload, hardware, and schema dependent; the research correctly flags this in open questions.
Can Be Approved Without Fix: YES

## Gap 3
Type: SCOPE_ERROR
Severity: LOW
Location: `research/05-report.md:Limitations`
Problem: Modern distributed NewSQL systems (CockroachDB, Google Cloud Spanner, YugabyteDB) that manage range-based sharding and consensus under the hood are excluded from deep comparative analysis.
Required Revision: None. The research topic is explicitly scoped to application/proxy-level sharding and partitioning fundamentals (PostgreSQL, MongoDB, Vitess).
Can Be Approved Without Fix: YES
