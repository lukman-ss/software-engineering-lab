# 03 - Claim Audit

## Claim 1

Claim:
Expand-Migrate-Contract (Parallel Change) separates adding new capabilities from removing legacy contracts to avoid breaking consumers.

Location:
`research/06-expand-migrate-contract.md:3-22`, `research/11-final-research.md:19-23`

Evidence Provided:
Danilo Sato & Martin Fowler publication on Parallel Change.

Source:
Source 1 (Martin Fowler)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Accurately reflects the 3-phase methodology outlined by Fowler and Kerievsky.

---

## Claim 2

Claim:
Breaking changes classify into Wire incompatibility, Semantic incompatibility, and Source incompatibility.

Location:
`research/03-core-concepts.md:7-13`, `research/11-final-research.md:16-18`

Evidence Provided:
Google Cloud AIP-180 standard.

Source:
Source 4 (Google Cloud AIP-180)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
AIP-180 explicitly defines these three dimensions of compatibility.

---

## Claim 3

Claim:
Additive fields (adding properties or response fields) are safe because consumers following Postel's Law ignore unknown fields.

Location:
`research/05-api-compatibility.md:5`, `research/11-final-research.md:13-14`

Evidence Provided:
General reference to AIP-180 and Postel's Law.

Source:
Source 4 (Google Cloud AIP-180)

Source Actually Supports Claim:
PARTIAL

Classification:
INTERPRETATION

Severity:
MEDIUM

Notes:
While generally additive changes are safe, AIP-180 lists several critical caveats: new required fields, default serialization shifts, enum additions that crash client switch statements, and naming conflicts in generated code stubs. The claim should state that additive changes are safe only when optional and backward-tolerated.

---

## Claim 4

Claim:
Zero-downtime column removal requires a 3-release process: Release M (Ignore Column), Release M+1 (Drop Column via post-deployment migration), Release M+2 (Remove Ignore rule).

Location:
`research/04-database-migration.md:8-13`, `research/11-final-research.md:24-26`

Evidence Provided:
GitLab database migration guidelines.

Source:
Source 3 (GitLab Documentation)

Source Actually Supports Claim:
YES

Classification:
IMPLEMENTATION-SPECIFIC

Severity:
MEDIUM

Notes:
The 3-release requirement is specifically dictated by Rails ActiveRecord's schema cache behavior and GitLab's monthly release cadence. In frameworks without dynamic boot-time schema caching (such as typed Go/Java SQL mappers with explicit column projections), Release M+2 is often unnecessary.

---

## Claim 5

Claim:
Backfill must use small batches, throttling, and idempotent operations (`ON CONFLICT DO NOTHING`) to avoid exclusive table locks and replication lag.

Location:
`research/04-database-migration.md:19-23`, `research/08-failure-modes.md:13-16`

Evidence Provided:
GitLab Batched Background Migrations documentation.

Source:
Source 3 (GitLab Documentation)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Directly matches standard database reliability engineering practices.

---

## Claim 6

Claim:
During rolling deployment, versions N and N+1 run simultaneously against a single database; the database must remain compatible with version N until rollout completes.

Location:
`research/04-database-migration.md:24-27`, `research/07-deployment-and-rollback.md:3-12`, `research/11-final-research.md:36-38`

Evidence Provided:
Fowler Parallel Change & GitLab avoiding downtime.

Source:
Source 1 & Source 3

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Fundamental requirement of zero-downtime rolling updates.

---

## Claim 7

Claim:
Dual-write introduces risks of partial data inconsistency if one operation fails without atomic transaction handling, as well as latency and race conditions.

Location:
`research/04-database-migration.md:16`, `research/08-failure-modes.md:9-12`, `research/11-final-research.md:30-32`

Evidence Provided:
Parallel Change failure modes analysis.

Source:
Source 1 (Martin Fowler)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Accurately identifies the primary risk in dual-write architectures.
