# Research Report

## Research Question
How should an engineering team document and maintain architectural decisions through Architecture Decision Records (ADRs) to prevent lost context and redundant debate, specifically when choosing between a Modular Monolith and Microservices under strict organizational constraints?

## Executive Summary
An Architecture Decision Record (ADR) is a lightweight, version-controlled document capturing an architecturally significant design choice, its forces/context, evaluated alternatives, decision rationale, and full consequences (positive, negative, neutral). ADRs must be immutable once accepted and stored in version control alongside code; evolving architecture issues a new ADR that marks the predecessor as Superseded or Deprecated, preserving historical context. ADRs are reserved for decisions affecting structure, NFRs, dependencies, interfaces, or construction techniques — not routine implementation details. For early-stage systems with small teams, evolving requirements, and short deadlines (such as the SaaS ERP scenario), a modular monolith minimizes distributed-system operational overhead while preserving explicit domain boundaries that permit future service extraction.

## Findings

### Finding 1: An ADR is a single-decision record capturing forces, decision, and consequences

Claim: An ADR must document a single architectural decision in response to a set of forces, including full consequences (positive, negative, neutral), not merely an implementation specification.

Evidence: Michael Nygard (Cognitect, 2011): "Each record describes a set of forces and a single decision in response to those forces" and "All consequences should be listed here, not just the 'positive' ones."

Sources:
- Cognitect Blog: https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions
- AWS Prescriptive Guidance: https://docs.aws.amazon.com/prescriptive-guidance/latest/architectural-decision-records/adr-process.html
- Microsoft Azure Well-Architected: https://learn.microsoft.com/en-us/azure/well-architected/architect-role/architecture-decision-record
- MADR full template: https://adr.github.io/madr/

Confidence: HIGH

### Finding 2: ADRs must be immutable in version control; evolution happens by supersedure

Claim: Once accepted, an ADR is immutable. Any change of direction produces a new ADR that references the predecessor (Superseded / Deprecated). Accepted records must not be edited to make history appear consistent with the present.

Evidence: AWS Prescriptive Guidance (2026-09-26 inspection): "When the team accepts an ADR, it becomes immutable. If new insights require a different decision, the team proposes a new ADR. When the team accepts the new ADR, it supersedes the previous ADR." Nygard (2011): "If a decision is reversed, we will keep the old one around, but mark it as superseded." Microsoft Azure Well-Architected: "If a decision changes, write a new record that supersedes the original and link the two together."

Sources:
- Cognitect Blog: https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions
- AWS Prescriptive Guidance: https://docs.aws.amazon.com/prescriptive-guidance/latest/architectural-decision-records/adr-process.html
- Microsoft Azure Well-Architected: https://learn.microsoft.com/en-us/azure/well-architected/architect-role/architecture-decision-record

Confidence: HIGH

### Finding 3: Scope is architecturally significant decisions, not implementation details

Claim: ADRs belong on decisions affecting structure, non-functional requirements, dependencies, interfaces, construction techniques, or choices that are hard to reverse — explicitly excluding routine refactoring, renames, CRUD endpoints, or minor changes.

Evidence: AWS Prescriptive Guidance (2026-09-26 inspection): lists five categories requiring ADRs: structure (e.g., microservices), NFRs (security, HA, fault tolerance), dependencies, interfaces, construction techniques. Microsoft Azure Well-Architected: "Only include choices that affect the system's structure, key quality attributes, or are difficult to reverse." Nygard: "architecturally significant decisions: those that affect the structure, non-functional characteristics, dependencies, interfaces, or construction techniques."

Sources:
- AWS Prescriptive Guidance: https://docs.aws.amazon.com/prescriptive-guidance/latest/architectural-decision-records/adr-process.html
- Microsoft Azure Well-Architected: https://learn.microsoft.com/en-us/azure/well-architected/architect-role/architecture-decision-record
- Cognitect Blog: https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions

Confidence: HIGH

### Finding 4: ADR lifecycle (status states and transitions)

Claim: ADR status follows the transitions Proposed → Accepted (and optionally Rejected) → Superseded / Deprecated. These states are mutually recognized across Nygard, AWS, MADR, and Azure Well-Architected, with Azure AWS adding an explicit Rejected state to capture discarded proposals.

Evidence: Nygard defines Proposed, Accepted, Deprecated, Superseded. AWS adds Rejected ("the ADR owner adds a reason for the rejection to prevent future discussions"). MADR template status field lists "{proposed | rejected | accepted | deprecated | … | superseded by ADR-0123}". Azure Well-Architected lists "*Proposed*, *Accepted*, or *Superseded*".

Sources:
- Cognitect Blog: https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions
- AWS Prescriptive Guidance: https://docs.aws.amazon.com/prescriptive-guidance/latest/architectural-decision-records/adr-process.html
- MADR: https://adr.github.io/madr/
- Microsoft Azure Well-Architected: https://learn.microsoft.com/en-us/azure/well-architected/architect-role/architecture-decision-record

Confidence: HIGH

### Finding 5: Modular monolith is appropriate for early-stage SaaS ERP given the scenario constraints

Claim: The SaaS ERP scenario (5 engineers, 3-month delivery window, rapidly evolving requirements, no dedicated SRE/platform team, existing Laravel expertise) favors a Laravel Modular Monolith over Go Microservices, because microservices introduce a MicroservicePremium (network failure modes, distributed transactions, observability burden, deployment orchestration) without yet delivering the scaling or team-decoupling benefits that would offset it.

Evidence: Martin Fowler ("Monolith First", 2015-06-03): "Microservices are a useful architecture, but even their advocates say that using them incurs a significant MicroservicePremium... they are only useful with more complex systems." "By building a monolith first, you can figure out what the right boundaries are" and "it also gives you time to develop the MicroservicePrerequisites you need for finer-grained services." AWS scope definition: structure choices (monolith vs microservices) and construction techniques are explicitly ADR-worthy — confirming the decision merits documentation.

Sources:
- Martin Fowler: https://martinfowler.com/bliki/MonolithFirst.html
- AWS Prescriptive Guidance: https://docs.aws.amazon.com/prescriptive-guidance/latest/architectural-decision-records/adr-process.html

Confidence: HIGH

### Finding 6: ADR filename convention and index discovery pattern

Claim: ADRs use monotonically increasing numeric prefix with descriptive kebab-case lowercase filename, e.g., ADR-001-use-modular-monolith.md, and are discoverable via a decision log index (README.md in docs/adr/).

Evidence: Nygard: "ADRs will be numbered sequentially and monotonically. Numbers will not be reused." Nygard: keep under "doc/arch/adr-NNN.md". MADR: "NNNN-title-with-dashes.md". Community convention (Joel Parker Henderson repo): filename index with descriptive title. MADR and adr.github.io describe decision-log README indexes as the discovery mechanism.

Sources:
- Cognitect Blog: https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions
- MADR: https://adr.github.io/madr/
- joelparkerhenderson/architecture-decision-record: https://github.com/joelparkerhenderson/architecture-decision-record

Confidence: HIGH

## Areas of Agreement
- ADRs co-located with source code in Git (Nygard, AWS, Azure Well-Architected, MADR all concur).
- ADRs must record positive, negative, and neutral consequences.
- Monotonically numbered markdown files provide the optimal balance of readability and durability.
- Decisions are contextual; a previously correct decision is not "wrong" simply because it is later superseded.
- Architectural significance is scoped to structure, NFRs, dependencies, interfaces, and construction techniques.

## Areas of Disagreement
- **Rejected status**: Nygard omits a formal `Rejected` state (left to informal tracking); AWS and MADR explicitly support `Rejected` to prevent re-litigation of discarded options. The union set `{Proposed, Accepted, Rejected, Deprecated, Superseded}` resolves this for the lab.
- **Strict immutability vs living document**: Tier 1 sources mandate strict append-only immutability; a Tier 2 community reference reports teams preferring date-stamped edits within the same record. Lab follows Tier 1 to preserve auditable history (matching the topic spec's prohibition on rewriting old ADRs).
- **Template granularity**: Nygard (5 sections) vs MADR/AWS (10+ sections). Lab adopts the richer anatomy (Context, Drivers, Alternatives, Decision, Rationale, Consequences, Risks, Review Triggers) to satisfy the topic spec's explicit required sections, while noting Nygard's minimalism as the canonical floor.

## Limitations
- The lab scenario (SaaS ERP, Laravel vs Go microservices) is illustrative guidance context — not a field study. Real performance, scaling, and team-velocity claims for Laravel vs Go would require production benchmarks (out of scope; topic spec explicitly excludes "benchmark Laravel vs Go").
- Automated ADR validation (section presence, status keyword) is structural; semantic soundness of rationales or trade-off correctness cannot be machine-verified.
- MicroservicePremium magnitude is described qualitatively (per Fowler) without quantified empirical baselines.

## Conclusion
ADRs preserve engineering intent and historical context by combining contextual constraints, evaluated alternatives, explicitly accepted trade-offs, and observable review triggers. For the given SaaS ERP scenario, a modular monolith records the right contextual boundary: it minimizes premature distributed-systems complexity while preserving domain separation that permits measured service extraction when review triggers fire. This reflects the senior-engineering mental model: architecture is a trade-off under constraints, and a decision can be correct for its time without being a permanent decree.
