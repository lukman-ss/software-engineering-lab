# 03 - Claim Audit

## Claim 1

Claim:
Backward-incompatible changes can be achieved safely through the Parallel Change (Expand -> Migrate -> Contract) pattern.

Location:
`research/06-expand-migrate-contract.md`, `research/11-final-research.md` (Q4)

Evidence Provided:
Danilo Sato and Martin Fowler's "Parallel Change" bliki post.

Source:
Source 1

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Perfect alignment with Source 1. Both conceptual phases and practical implementations are detailed.

---

## Claim 2

Claim:
Stripe achieves API backward compatibility via internal version change modules where responses are transformed sequentially backward through time, enabling core logic to run on the latest schema.

Location:
`research/05-api-compatibility.md` (Strategy 3)

Evidence Provided:
Stripe Engineering Blog (APIs as infrastructure).

Source:
Source 2

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Accurately maps to Stripe's documented "Version change modules" and rolling API versioning system.

---

## Claim 3

Claim:
During rolling deployments, Version N and Version N+1 concurrently access the same database. Thus, the database schema must be simultaneously compatible with both application versions.

Location:
`research/07-deployment-and-rollback.md` (Section 1), `research/11-final-research.md` (Q9)

Evidence Provided:
Derived from the requirements of zero-downtime database deployment and Parallel Change.

Source:
Source 1, Source 3

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Source 1 explicitly covers Deployments (BlueGreen, Canary) with mixed versions requiring compatibility. Source 3 describes the multi-step deployment strategy to ensure zero disruption.

---

## Claim 4

Claim:
Dual-writing introduces risks of increased latency, partial data inconsistencies, and drift if operations are not atomic or fully synchronized.

Location:
`research/04-database-migration.md` (Section 4), `research/08-failure-modes.md` (Failure Mode 2), `research/11-final-research.md` (Q7)

Evidence Provided:
Engineering risks associated with the migrate phase.

Source:
Source 1, Source 3

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Source 3 explicitly details writing to both structures during the Migrate/Expand phase and emphasizes ensuring data correctness before cutting over.

---

## Claim 5

Claim:
Database backfill of large datasets must use chunking, batch processing, throttling (sleep), and resume/checkpoint logic (idempotency) to avoid exhaustive locking and replication lag.

Location:
`research/04-database-migration.md` (Section 3), `research/11-final-research.md` (Q8)

Evidence Provided:
General database reliability engineering patterns.

Source:
Source 3

Source Actually Supports Claim:
PARTIAL

Classification:
IMPLEMENTATION-SPECIFIC

Severity:
MEDIUM

Notes:
While Source 3 details data migration scripts, it does not explicitly delve into batching, throttling limits, or sleep mechanisms in the provided text. These are universal database engineering truths, but they extend slightly beyond the explicit text of the cited Prisma guide. They are technically correct but under-sourced.

---

## Claim 6

Claim:
Legacy interfaces should only be removed once usage metrics drop to zero. A 30-day consistent zero-usage period is cited as a guideline but explicitly labeled as `NOT VERIFIED` as a universal standard.

Location:
`research/08-failure-modes.md` (Failure Mode 4), `research/11-final-research.md` (Q11)

Evidence Provided:
Industry heuristic with explicit disclaimer.

Source:
N/A (Derived operational heuristic)

Source Actually Supports Claim:
YES (As a documented caveat)

Classification:
HYPOTHESIS

Severity:
LOW

Notes:
The research properly flags the 30-day time window as an unverified heuristic, avoiding treating arbitrary numbers as absolute engineering law.

---

## Claim 7

Claim:
Adding default constraints on new columns on large tables in older database versions forces full table rewrites and heavy locking.

Location:
`research/04-database-migration.md` (Section 2)

Evidence Provided:
General PostgreSQL/MySQL operational behavior.

Source:
N/A (General knowledge)

Source Actually Supports Claim:
PARTIAL

Classification:
IMPLEMENTATION-SPECIFIC

Severity:
LOW

Notes:
This is a true database engine quirk (e.g. PostgreSQL < 11 rewrote the entire table for default values). Although Source 3 does not discuss engine-specific locking internals, the claim is factually accurate and appropriately constrained by "in older database versions."
