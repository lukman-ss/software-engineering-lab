# Claim Audit

## Claim 1

Claim:
Backward compatibility means new systems accept old inputs; Forward compatibility means old systems accept new inputs.

Location:
`03-core-concepts.md` (Evidence 1 & 2)

Evidence Provided:
Quotes from Google AIP-180, Stripe API Versioning, Confluent Schema Registry.

Source:
Confluent Schema Registry (schema-evolution.html), Google AIP-180

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Accurately reflects industry-standard definitions of message evolution.

---

## Claim 2

Claim:
Parallel Change pattern divides breaking changes into three non-breaking phases: Expand, Migrate, Contract.

Location:
`03-core-concepts.md` (Evidence 4), `06-expand-migrate-contract.md`

Evidence Provided:
Joshua Kerievsky / Martin Fowler "Parallel Change" bliki.

Source:
Martin Fowler (Parallel Change)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
The claim accurately reflects the source material, including the `Grid` and `Coordinate` Java examples.

---

## Claim 3

Claim:
Zero-downtime database migration relies on additive non-rewriting changes, such as PostgreSQL `ALTER TABLE ADD COLUMN` with a constant default and `CREATE INDEX CONCURRENTLY`.

Location:
`04-database-migration.md`, `03-core-concepts.md` (Evidence 5)

Evidence Provided:
PostgreSQL documentation on `ALTER TABLE` and indexing lock behaviors.

Source:
PostgreSQL ALTER TABLE Docs

Source Actually Supports Claim:
YES

Classification:
IMPLEMENTATION-SPECIFIC

Severity:
LOW

Notes:
The research correctly identifies these traits as PostgreSQL-specific and explicitly warns against treating them as universal relational database features.

---

## Claim 4

Claim:
During the migrate phase, dual writes to old and new schemas risk data drift and split-brain if not executed atomically.

Location:
`03-core-concepts.md` (Evidence 7), `08-failure-modes.md`

Evidence Provided:
References to 2PC limitations and dual-write unreliability described in Microservices.io Transactional Outbox pattern.

Source:
Microservices.io Transactional Outbox

Source Actually Supports Claim:
PARTIAL

Classification:
INTERPRETATION

Severity:
MEDIUM

Notes:
The source specifically discusses outbox patterns for message brokers, not dual-table writes. However, the architectural physics (distributed/dual write atomicity failure) are identical, and the inference is sound.

---

## Claim 5

Claim:
Rolling and Blue-Green deployments require separating schema changes from application upgrades, forcing the schema to simultaneously support Version N and Version N+1.

Location:
`07-deployment-and-rollback.md` (Fundamental Rule)

Evidence Provided:
Martin Fowler's guidance on decoupling DB refactoring from Blue Green deploy.

Source:
Martin Fowler (Blue Green Deployment)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Source verbatim states: "first apply a database refactoring to change the schema to support both the new and old version".

---

## Claim 6

Claim:
Deprecation heuristic requires removing legacy code when usage metrics remain at 0 for 30 days.

Location:
`08-failure-modes.md` (Failure Mode 4)

Evidence Provided:
Inferred best practice. GitHub API uses a 24-month window.

Source:
None (Marked NOT VERIFIED in research)

Source Actually Supports Claim:
NO

Classification:
HYPOTHESIS

Severity:
MEDIUM

Notes:
The research appropriately flags the 30-day metric as an unverified heuristic and contrasts it with GitHub's 24-month hard policy.

---

## Claim 7

Claim:
Renaming an API component is semantically equivalent to a breaking 'remove and add' operation.

Location:
`05-api-compatibility.md` (Removing atau Renaming Components)

Evidence Provided:
Google AIP-180 rules.

Source:
Google AIP-180

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Matches AIP-180 verbatim.
