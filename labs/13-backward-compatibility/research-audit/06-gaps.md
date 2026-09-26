# 06 - Research Gap Analysis

## Gap 1

Type:
WEAK_SOURCE

Severity:
LOW

Location:
`research/04-database-migration.md` (Section 3: Data Backfill Strategies)

Problem:
Batch processing (`LIMIT`/`OFFSET` or ID range queries), sleep throttling, and checkpoint idempotency are standard production database patterns, but are not explicitly discussed in Source 3 (Prisma Data Guide).

Required Revision:
Include a dedicated Tier 1 source for large-scale production schema migrations and background migrations (e.g., GitLab Batched Background Migrations or GitHub `gh-ost` documentation).

Can Be Approved Without Fix:
YES

---

## Gap 2

Type:
UNVERIFIED_CLAIM

Severity:
LOW

Location:
`research/08-failure-modes.md` (Section 4), `research/11-final-research.md` (Q11)

Problem:
The 30-day sustained observation window before removing legacy contracts is an operational rule of thumb rather than an empirically verified constant across all industries. While correctly labeled as `NOT VERIFIED`, it could benefit from contextualizing with periodic business processes (e.g., quarterly accounting runs or monthly billing cycles).

Required Revision:
Clarify that observation windows must match or exceed the longest consumer business cycle (e.g., quarterly reconciliation).

Can Be Approved Without Fix:
YES

---

## Gap 3

Type:
MISSING_CASE

Severity:
LOW

Location:
`research/04-database-migration.md` (Section 4), `research/10-open-questions.md`

Problem:
The research predominantly discusses synchronous application-level dual writing. Asynchronous Change Data Capture (CDC via Kafka/Debezium or logical replication) is only mentioned as an open question, leaving an alternative dual-writing pattern under-explored.

Required Revision:
Add a brief architectural comparison between application-level synchronous dual writing and asynchronous CDC-based replication in subsequent iterations.

Can Be Approved Without Fix:
YES
