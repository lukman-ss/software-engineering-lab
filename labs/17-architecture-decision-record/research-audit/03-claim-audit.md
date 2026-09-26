# Claim Audit

## Claim 1

Claim: An ADR must record the context, forces, a single decision, and full consequences (positive, negative, and neutral) rather than just implementation specs.
Location: `research/runs/2026-09-26-architecture-decision-record/03-evidence.md` (Evidence 1); `05-report.md` (Finding 1)
Evidence Provided: Quotes from Michael Nygard (2011), corroborated by AWS Prescriptive Guidance, Microsoft Azure Well-Architected, and MADR.
Source: Cognitect Blog (https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions)
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Core invariant across all canonical ADR literature.

## Claim 2

Claim: ADRs must be immutable once accepted; changes require a new ADR that supersedes or deprecates the predecessor.
Location: `research/runs/2026-09-26-architecture-decision-record/03-evidence.md` (Evidence 2); `05-report.md` (Finding 2)
Evidence Provided: AWS Prescriptive Guidance: "When the team accepts an ADR, it becomes immutable...", Nygard: "mark it as superseded", Azure Well-Architected: "The ADR serves as an append-only log."
Source: AWS Prescriptive Guidance / Cognitect Blog / Microsoft Learn
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Fully supported by all Tier 1 primary sources. The research agent appropriately documented community divergence (living document pattern) in contradictions without compromising canonical guidance.

## Claim 3

Claim: ADRs belong inside version control alongside source code.
Location: `research/runs/2026-09-26-architecture-decision-record/03-evidence.md` (Evidence 3); `05-report.md` (Executive Summary & Finding 2)
Evidence Provided: Nygard: "We will keep ADRs in the project repository under doc/arch/adr-NNN.md", Azure Well-Architected: "stored openly with the workload's documentation."
Source: Cognitect Blog / Microsoft Learn
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Strongly supported. Ensures architectural evolution is versioned alongside the code it governs.

## Claim 4

Claim: Architectural significance applies to structure, NFRs, dependencies, interfaces, and construction techniques, not low-level refactoring.
Location: `research/runs/2026-09-26-architecture-decision-record/03-evidence.md` (Evidence 4); `05-report.md` (Finding 3)
Evidence Provided: AWS Prescriptive Guidance categorizing significance into structure, NFRs, dependencies, interfaces, and construction techniques; echoed by Nygard and Richards & Ford (2020).
Source: AWS Prescriptive Guidance (https://docs.aws.amazon.com/prescriptive-guidance/latest/architectural-decision-records/adr-process.html)
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Essential scoping definition that prevents ADR logs from becoming bloated with routine refactoring notes.

## Claim 5

Claim: ADR status transitions follow: Proposed -> Accepted -> Superseded / Deprecated (or Rejected).
Location: `research/runs/2026-09-26-architecture-decision-record/03-evidence.md` (Evidence 5); `05-report.md` (Finding 4)
Evidence Provided: Nygard / AWS / MADR lifecycle definitions. AWS explicitly formalizes Rejected status.
Source: Cognitect Blog / AWS Prescriptive Guidance / MADR
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Correctly accounts for differences between Nygard (4 statuses) and AWS/MADR (5 statuses) via a clean superset union.

## Claim 6

Claim: Modular monoliths reduce operational complexity while preserving domain boundaries for early-stage systems compared to microservices.
Location: `research/runs/2026-09-26-architecture-decision-record/03-evidence.md` (Evidence 6); `05-report.md` (Finding 5)
Evidence Provided: Martin Fowler "Monolith First": MicroservicePremium explanation and premature boundary risks.
Source: Martin Fowler (https://martinfowler.com/bliki/MonolithFirst.html)
Source Actually Supports Claim: YES
Classification: INTERPRETATION
Severity: MEDIUM
Notes: Supported by Fowler's analysis. The research agent accurately presents this as context-sensitive engineering trade-off rather than an absolute rule, and notes the counter-argument for teams with preexisting microservice capabilities and stable boundaries.

## Claim 7

Claim: Measurable criteria such as deployment cadence divergence, resource contention, or team ownership shifts constitute valid "Review Triggers" to re-evaluate architectural decisions.
Location: `research/runs/2026-09-26-architecture-decision-record/03-evidence.md` (Evidence 7); `05-report.md` (Finding 6 / Open Questions)
Evidence Provided: Implicit in Azure Well-Architected confidence logging and consequence tracking; explicit in community templates (Joel Parker Henderson).
Source: Microsoft Learn / Community Practice
Source Actually Supports Claim: PARTIAL
Classification: HYPOTHESIS
Severity: MEDIUM
Notes: The principle of review triggers is widely acknowledged, but specific numeric threshold calibrations (e.g. exactly how much cadence divergence warrants service extraction) lack authoritative empirical constants in the sources. The research agent correctly classified this under "Open Questions" as a medium-priority gap rather than presenting arbitrary numbers.
