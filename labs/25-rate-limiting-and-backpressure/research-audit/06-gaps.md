# Research Gaps

## Gap 1: Reliance on Tertiary Wikipedia Sources for Core Algorithms
- **Type**: WEAK_SOURCE
- **Severity**: LOW
- **Location**: `research/02-sources.md` (Sources 2, 3, 6, 8, 9)
- **Problem**: 5 out of 10 sources cited in `02-sources.md` are Wikipedia articles rather than primary academic publications or standard textbooks (e.g., Kurose & Ross, John D.C. Little's 1961 Operations Research paper, or the official Reactive Streams specification repository).
- **Required Revision**: Supplement Wikipedia citations with primary academic or specification references where rigorous mathematical or protocol authority is needed.
- **Can Be Approved Without Fix**: YES (The cited Wikipedia formulations accurately reflect canonical definitions and formulas).

---

## Gap 2: Limited Empirical Benchmarks for High-Concurrency Rate Limiting
- **Type**: MISSING_CASE
- **Severity**: LOW
- **Location**: `research/05-report.md:185`, `research/06-open-questions.md:35`
- **Problem**: The research identifies the difference between Sliding Window Log, Sliding Window Counter, and Token Bucket, but lacks quantitative throughput/memory benchmark comparisons under heavy concurrent loads.
- **Required Revision**: Include empirical benchmark numbers (e.g. Redis memory per million keys across algorithms) in future deep dives.
- **Can Be Approved Without Fix**: YES (Noted in `06-open-questions.md` as open research areas).
