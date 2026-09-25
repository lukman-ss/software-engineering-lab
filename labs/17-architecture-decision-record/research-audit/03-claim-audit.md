# Claim Audit: Architecture Decision Record Research

## Claim 1

Claim:
An ADR must record the motivation, forces, single decision, and full consequences (positive, negative, and neutral) rather than just implementation specs.

Location:
`03-evidence.md` (Evidence 1), `05-report.md` (Finding 1)

Evidence Provided:
Michael Nygard (2011) emphasizes that each record describes forces and a single decision, listing all consequences. AWS Prescriptive Guidance instructs focusing on reasons over implementation mechanisms.

Source:
- Cognitect Blog (Michael Nygard)
- AWS Prescriptive Guidance

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Accurately reflects the core definition in original and contemporary literature.

---

## Claim 2

Claim:
ADRs must be immutable once accepted; changes require a new ADR that supersedes or deprecates the predecessor.

Location:
`03-evidence.md` (Evidence 2), `05-report.md` (Finding 2)

Evidence Provided:
AWS Prescriptive Guidance states that accepted ADRs are immutable and modifications require proposing a new ADR that supersedes the prior one. Nygard confirms keeping older decisions marked as superseded.

Source:
- AWS Prescriptive Guidance
- Cognitect Blog

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Foundational consensus across all ADR governance frameworks.

---

## Claim 3

Claim:
ADRs belong inside version control co-located with source code in markdown format with monotonic numbering.

Location:
`03-evidence.md` (Evidence 3), `05-report.md` (Executive Summary & Areas of Agreement)

Evidence Provided:
Nygard specifies storing ADRs under `doc/arch/adr-NNN.md` in version control, numbered sequentially and monotonically without number reuse. adr.github.io corroborates git-based AKM practices.

Source:
- Cognitect Blog
- adr.github.io

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Well-established standard practice across industry and open-source ecosystems.

---

## Claim 4

Claim:
Architectural significance applies to structure, NFRs, dependencies, interfaces, and construction techniques, excluding low-level refactoring.

Location:
`03-evidence.md` (Evidence 4), `05-report.md` (Finding 3)

Evidence Provided:
AWS Prescriptive Guidance cites Richards and Ford (2020) defining the five dimensions of architectural significance. adr.github.io defines Architecturally Significant Requirements (ASRs).

Source:
- AWS Prescriptive Guidance
- adr.github.io

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Precise boundary distinguishing architectural decisions from implementation tasks.

---

## Claim 5

Claim:
ADR status transitions follow: Proposed -> Accepted -> Superseded / Deprecated (or Rejected).

Location:
`03-evidence.md` (Evidence 5), `05-report.md` (Finding 2 & Areas of Disagreement)

Evidence Provided:
Nygard defines Proposed, Accepted, Deprecated, Superseded. AWS Prescriptive Guidance introduces the Rejected state to archive non-viable architectural choices.

Source:
- Cognitect Blog
- AWS Prescriptive Guidance

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Discrepancy between Nygard and AWS regarding the `Rejected` state is properly identified and documented.

---

## Claim 6

Claim:
For an early-stage SaaS ERP with a small team (5 engineers) and a 3-month deadline, a Modular Monolith delivers necessary domain boundaries without the operational, networking, and distributed transaction complexity of microservices.

Location:
`05-report.md` (Finding 4)

Evidence Provided:
Architectural trade-off principles highlighting that premature distributed systems introduce distributed transactions, observability burdens, and network failure modes before product-market fit.

Source:
Cognitect Blog, AWS Prescriptive Guidance (cited generally as background)

Source Actually Supports Claim:
PARTIAL

Classification:
INTERPRETATION

Severity:
MEDIUM

Notes:
While sound as an architectural heuristic for early-stage engineering, the cited ADR process documents (Nygard, AWS ADR Process) do not specifically evaluate or benchmark "Laravel Modular Monolith vs Go Microservices for 5 engineers". The claim is an application scenario / case study interpretation, not an empirical finding directly derived from the two process citations. Needs clear qualification as a contextual scenario analysis rather than a direct citation claim.
