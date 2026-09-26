# Research Gaps: Optimistic vs Pessimistic Locking

## Gap 1

Type:
WEAK_SOURCE

Severity:
LOW

Location:
`research/02-sources.md`, Source 7, 8, 9

Problem:
Direct access to `dev.mysql.com` returned HTTP 403 (bot blocking) during automated fetch. Verification was completed via official Oracle CDN documentation mirrors (`docs.oracle.com/cd/E17952_01/mysql-8.0-en/...`).

Required Revision:
None required for approval; content was verified on official Oracle CDN.

Can Be Approved Without Fix:
YES

---

## Gap 2

Type:
SCOPE_ERROR

Severity:
LOW

Location:
`research/06-open-questions.md`, OQ-3

Problem:
Potential integer overflow of 32-bit version counters in long-lived optimistic locking systems is noted but not deeply explored in vendor documentation.

Required Revision:
Include 64-bit integer (`BIGINT`) or timestamp-based versioning recommendations in future iterations if high update rates are expected.

Can Be Approved Without Fix:
YES

---

## Gap 3

Type:
UNVERIFIED_CLAIM

Severity:
LOW

Location:
`research/06-open-questions.md`, OQ-4

Problem:
Quantitative empirical performance comparison (latency/throughput numbers under high contention) between pessimistic locking, optimistic retries, and atomic single-statement updates is described qualitatively without benchmark numbers.

Required Revision:
Perform TPC-C style benchmark testing during lab execution phase to capture concrete performance graphs.

Can Be Approved Without Fix:
YES

---

## Gap 4

Type:
UNVERIFIED_CLAIM

Severity:
MEDIUM

Location:
`research/05-report.md`, Finding 9; `research/06-open-questions.md`, OQ-2

Problem:
Asserting that using distributed locks (Redis) when data lives in a single database is an anti-pattern relies on architectural best-practice inference rather than an explicit Tier 1 benchmark source.

Required Revision:
Explicitly document boundary conditions (e.g., cross-service coordination vs single relational DB bounds) in content creation stage.

Can Be Approved Without Fix:
YES
