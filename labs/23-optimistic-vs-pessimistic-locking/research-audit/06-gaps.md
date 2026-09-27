# Research Gap Analysis: Optimistic vs Pessimistic Locking

## Gap 1

Type: UNVERIFIED_CLAIM / MISSING_SOURCE

Severity: MEDIUM

Location: `research/03-evidence.md:31-38`, `research/04-contradictions.md:20-27`, `research/06-open-questions.md:5-10`

Problem: MySQL/InnoDB documentation was not directly accessible during the initial research session due to HTTP 403. While behavior is widely known, official vendor reference was not directly captured.

Required Revision: Keep transparently documented or consult accessible MySQL mirrors if vendor-specific locking quirks (e.g. Next-Key locks / gap locks) are critical to subsequent engineering benchmarks.

Can Be Approved Without Fix: YES

---

## Gap 2

Type: MISSING_CASE / OVERGENERALIZATION

Severity: LOW

Location: `research/05-report.md:83-95`

Problem: Stating that atomic `UPDATE ... WHERE condition` eliminates race conditions is valid for single-row updates, but should explicitly note limitations when business invariants span multiple tables or aggregate constraints (e.g. account transfer or bank-wide balance limits).

Required Revision: Ensure downstream engineering content notes that single-row atomic updates do not replace multi-row transactions when invariants span multiple entities.

Can Be Approved Without Fix: YES

---

## Gap 3

Type: WEAK_SOURCE

Severity: LOW

Location: `research/06-open-questions.md:51-54`

Problem: Absence of quantitative benchmarks (throughput, latency, conflict rates) in primary theoretical literature.

Required Revision: Acknowledge that empirical performance figures will be measured directly in the runnable implementation/benchmarking phase.

Can Be Approved Without Fix: YES
