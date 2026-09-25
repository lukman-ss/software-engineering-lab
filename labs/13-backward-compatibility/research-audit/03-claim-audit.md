# Claim Audit

## Claim 1

Claim:
The Expand -> Migrate -> Contract (Parallel Change) pattern allows breaking incompatible changes into non-breaking, incremental phases: expand interface/schema, migrate usages/data, contract legacy interfaces.

Location:
`research/06-expand-migrate-contract.md`, Section "Overview"
`research/11-final-research.md`, Section "Answers to Core Research Questions: #4"

Evidence Provided:
Danilo Sato / Martin Fowler blog on Parallel Change; Prisma Data Guide.

Source:
Source 1 (Parallel Change, martinfowler.com), Source 3 (Prisma Data Guide)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Accurately reflects the principles of evolutionary architecture.

---

## Claim 2

Claim:
Stripe's API versioning maintains backward compatibility using internal version change modules that walk backward through time to translate modern response models to client-pinned versions.

Location:
`research/05-api-compatibility.md`, Section "Strategy 3: Internal Transformation Pipelines"
`research/02-sources.md`, Section "Source 2"

Evidence Provided:
Brandur Leach, Stripe Engineering Blog: "APIs as infrastructure: future-proofing Stripe with versioning".

Source:
Source 2 (Stripe Blog)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
The description matches Stripe's published architecture and implementation DSL.

---

## Claim 3

Claim:
During rolling deployments where versions N and N+1 run concurrently, the database schema must simultaneously support both versions without breaking.

Location:
`research/07-deployment-and-rollback.md`, Section "1. Rolling Deployment Realities"
`research/11-final-research.md`, Section "#9"

Evidence Provided:
Prisma Data Guide and Fowler's Parallel Change (deployment applications).

Source:
Source 1 & Source 3

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Fundamental requirement in continuous deployment environments.

---

## Claim 4

Claim:
Removing or renaming fields, columns, or endpoints constitutes an immediate breaking change.

Location:
`research/03-core-concepts.md`, Section "3. Breaking vs. Non-Breaking Changes"
`research/08-failure-modes.md`, Section "1. Immediate Destructive Schema Alteration"

Evidence Provided:
Fowler (Parallel Change) and Stripe API versioning article.

Source:
Source 1 & Source 2

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Confirmed across all cited sources.

---

## Claim 5

Claim:
Field, column, or endpoint should only be removed after usage metrics show zero hits for 30 consecutive days.

Location:
`research/11-final-research.md`, Section "#11"
`research/08-failure-modes.md`, Section "#4"

Evidence Provided:
None cited for the specific 30-day threshold.

Source:
Uncited / Operational Heuristic

Source Actually Supports Claim:
PARTIAL

Classification:
RECOMMENDATION

Severity:
MEDIUM

Notes:
The "30 consecutive days" metric is an arbitrary operational heuristic, not backed by an authoritative cited source. While common in practice, it should be marked as an operational guideline rather than a definitive standard.
