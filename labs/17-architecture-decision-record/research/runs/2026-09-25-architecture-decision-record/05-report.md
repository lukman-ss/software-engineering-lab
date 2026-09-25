# Research Report

## Research Question
How should an engineering team document and maintain architectural decisions through Architectural Decision Records (ADRs) to prevent lost context and redundant debate, specifically when choosing between a Modular Monolith and Microservices under strict organizational constraints?

## Executive Summary
An Architecture Decision Record (ADR) is a lightweight document capturing an architecturally significant design choice, its context, evaluated alternatives, decision rationale, trade-offs, and observable review triggers. ADRs must reside directly in the source code repository under monotonic numbering, remaining immutable once accepted. Revising an architecture requires issuing a subsequent ADR that marks the prior decision as superseded, preserving historical context. For early-stage systems with small teams and evolving requirements, modular monoliths minimize distributed systems overhead while enforcing domain boundaries.

## Findings

### Finding 1: Architecture Decisions Require Context, Not Just Solutions
Claim: Documenting what was chosen without recording why under specific constraints guarantees recurring debates.
Evidence: Michael Nygard (2011) observed that new team members either blindly accept or blindly reverse decisions without understanding forces. AWS Prescriptive Guidance emphasizes focusing on why a decision was made rather than how it was implemented.
Sources: Cognitect Blog (Nygard), AWS Prescriptive Guidance.
Confidence: HIGH

### Finding 2: Immutability and Lifecycle Management
Claim: ADRs must not be retroactively edited when architecture evolves; they must follow an explicit lifecycle (Proposed, Accepted, Superseded, Deprecated, Rejected).
Evidence: Both Nygard and AWS guidelines mandate that changes generate a new ADR. The existing ADR is linked and marked `Superseded`.
Sources: Cognitect Blog, AWS Prescriptive Guidance.
Confidence: HIGH

### Finding 3: Scope of Architecturally Significant Decisions
Claim: ADRs must be reserved for decisions that are hard to reverse, affect structure, non-functional requirements, dependencies, interfaces, or construction techniques.
Evidence: AWS Prescriptive Guidance (citing Richards & Ford 2020) and adr.github.io define Architecturally Significant Requirements (ASRs) as the qualification boundary for ADRs, explicitly excluding trivial code-level refactorings.
Sources: AWS Prescriptive Guidance, adr.github.io.
Confidence: HIGH

### Finding 4: Modular Monolith vs Microservices Trade-off Context
Claim: For an early-stage SaaS ERP with a small team (5 engineers) and a 3-month deadline, a Modular Monolith delivers necessary domain boundaries without the operational, networking, and distributed transaction complexity of microservices.
Evidence: Architectural principles show that premature distributed systems introduce network failure modes, observability burdens, and data consistency issues before product-market fit is established. Modular monoliths preserve boundaries, allowing future service extraction when observable review triggers are met.
Sources: Cognitect Blog, AWS Prescriptive Guidance.
Confidence: HIGH

## Areas of Agreement
- ADRs must be co-located with code in Git.
- ADRs must state positive, negative, and neutral consequences.
- Monotonically numbered markdown files provide the optimal balance of readability and durability.
- Decisions are contextual; an earlier decision is not "wrong" simply because it is superseded later.

## Areas of Disagreement
- Inclusion of the `Rejected` state: Nygard omits it (favoring uncommitted PRs or informal notes), whereas AWS Prescriptive Guidance formally tracks `Rejected` ADRs to stop recurring debates on previously discarded options.

## Limitations
- Automated validation is restricted to structural linting (verifying presence of required headings, valid status keywords) and cannot evaluate semantic soundness or factual rigor of rationales.

## Conclusion
ADRs preserve engineering intent and historical context. By combining contextual constraints, evaluated alternatives, explicit trade-offs, and observable review triggers, teams prevent architectural drift and eliminate cyclical debates.
