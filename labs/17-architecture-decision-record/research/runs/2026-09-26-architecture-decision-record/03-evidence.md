# Evidence

## Evidence 1

Claim: An ADR must record the context, forces, a single decision, and full consequences (positive, negative, and neutral) rather than just implementation specs.

Evidence: Michael Nygard states: "Each record describes a set of forces and a single decision in response to those forces... This section describes the resulting context, after applying the decision. All consequences should be listed here, not just the 'positive' ones."

Source: Documenting Architecture Decisions
URL: https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions
Confidence: HIGH
Corroborated By: 
- AWS Prescriptive Guidance: "At a minimum, each ADR should define the context of the decision, the decision itself, and the consequences of the decision... focuses on the reason for the decision rather than how the team implemented it."
- Microsoft Azure Well-Architected: "Each entry captures the context, justifications, and implications of a decision... Avoid hiding consequences of decisions intentionally or accidentally."
- MADR full template: Contains explicit "Context and Problem Statement", "Decision Outcome", and sections for "Consequences" split into Good/Bad.

Notes: Critical for explaining why technical debates restart if rationale is absent; distinguishes ADR from implementation documentation.

## Evidence 2

Claim: ADRs must be immutable once accepted; changes require a new ADR that supersedes or deprecates the predecessor.

Evidence: AWS Prescriptive Guidance: "When the team accepts an ADR, it becomes immutable. If new insights require a different decision, the team proposes a new ADR. When the team accepts the new ADR, it supersedes the previous ADR."

Source: Architectural decision record process
URL: https://docs.aws.amazon.com/prescriptive-guidance/latest/architectural-decision-records/adr-process.html
Confidence: HIGH
Corroborated By: 
- Michael Nygard (2011): "If a decision is reversed, we will keep the old one around, but mark it as superseded."
- Microsoft Azure Well-Architected: "An ADR should be started for brownfield workloads... If a decision changes, write a new record that supersedes the original and link the two together."
- MADR: Status field supports "{proposed | rejected | accepted | deprecated | … | superseded by ADR-0123}".

Notes: Historical continuity is essential; editing accepted ADRs distorts historical context and prevents understanding why a decision was made.

## Evidence 3

Claim: ADRs belong inside version control alongside source code.

Evidence: Nygard states: "We will keep ADRs in the project repository under doc/arch/adr-NNN.md... keeping these in version control with the code makes them less accessible... In practice, our projects almost all live in GitHub... it looks just as friendly as any wiki page would."

Source: Documenting Architecture Decisions
URL: https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions
Confidence: HIGH
Corroborated By: 
- adr.github.io: "The aim of the [GitHub adr organization] is to: ... Strengthen the tooling around ADRs, in support of agile practices as well as iterative and incremental engineering processes."
- Microsoft Azure Well-Architected: "This log should be readily available and stored openly with the workload's documentation."
- MADR project: Instructions for "Create folder `docs/decisions` in your project" implying version control.

Notes: Prevents decision drift; couples architectural decisions to code history so future engineers can see both what changed and why.

## Evidence 4

Claim: Architectural significance applies to structure, NFRs, dependencies, interfaces, and construction techniques, not low-level refactoring.

Evidence: AWS Prescriptive Guidance: "Project members should create an ADR for every architecturally significant decision that affects the software project or product, including the following ([Richards and Ford 2020]): + Structure (for example, patterns such as microservices) + Non-functional requirements (security, high availability, and fault tolerance) + Dependencies (coupling of components) + Interfaces (APIs and published contracts) + Construction techniques (libraries, frameworks, tools, and processes)"

Source: Architectural decision record process
URL: https://docs.aws.amazon.com/prescriptive-guidance/latest/architectural-decision-records/adr-process.html
Confidence: HIGH
Corroborated By: 
- Michael Nygard: "We will keep a collection of records for 'architecturally significant' decisions: those that affect the structure, non-functional characteristics, dependencies, interfaces, or construction techniques."
- Microsoft Azure Well-Architected: "Your architecture is the accumulation of its decisions, so the ADR is effectively a record of how and why the system came to be its current shape. An architecture decision record (ADR) is one of the most important deliverables of a solution architect. ... Only include choices that affect the system's structure, key quality attributes, or are difficult to reverse."
- MADR: Focus on "Architectural Decision (AD) is a justified software design choice that addresses a functional or non-functional requirement of architectural significance."

Notes: Filters trivial implementation details (e.g., helper rename, minor CRUD addition, variable rename) to keep ADR log signal-to-noise ratio high.

## Evidence 5

Claim: ADR status transitions follow: Proposed -> Accepted -> Superseded / Deprecated (or Rejected).

Evidence: Nygard / AWS: States define the lifecycle. Proposed ADRs are reviewed; Accepted become active and immutable; changes introduce new records which mark older records as Superseded with a reference to the replacement.

Source: Documenting Architecture Decisions / AWS Prescriptive Guidance
URL: https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions
Confidence: HIGH
Corroborated By: 
- AWS Prescriptive Guidance (explicit "Rejected" state): "If the team finds action points to improve the ADR, the state of the ADR stays 'Proposed'... The team can also decide to reject the ADR. In this case, the ADR owner adds a reason for the rejection to prevent future discussions on the same topic."
- Microsoft Azure Well-Architected: Lists status options as "*Proposed*, *Accepted*, or *Superseded*" but accepts "Rejected" via process description.
- MADR: Status field in template explicitly includes rejected state.

Notes: Rejected status preserves institutional memory of considered-and-discarded options to prevent re-litigation.

## Evidence 6

Claim: Modular monoliths reduce operational complexity while preserving domain boundaries for early-stage systems compared to microservices.

Evidence: Martin Fowler's "Monolith First" article states: "Microservices are a useful architecture, but even their advocates say that using them incurs a significant MicroservicePremium, which means they are only useful with more complex systems. This premium, essentially the cost of managing a suite of services, will slow down a team, favoring a monolith for simpler applications." Further: "By building a monolith first, you can figure out what the right boundaries are, before a microservices design brushes a layer of treacle over them. It also gives you time to develop the MicroservicePrerequisites you need for finer-grained services."

Source: Monolith First
URL: https://martinfowler.com/bliki/MonolithFirst.html
Confidence: HIGH
Corroborated By:
- AWS Prescriptive Guidance discussion of operational complexity: While not directly quoted in the page, the guidance inherently acknowledges trade-offs by focusing on architecturally significant decisions that affect NFRs and structure, implying simplicity is a valid consideration.
- Microsoft Azure Well-Architected: Implicit in the advice to avoid hiding consequences and to include justifications; if choosing microservices introduced unacceptable operational overhead for the team's context, that would be a valid consequence to document.
- Joel Parker Henderson repo (corroboration only): Templates include space for consequences/trade-offs which would capture this context.

Notes: Directly supports the lab scenario (small team, evolving requirements, short deadline) by showing microservices introduce overhead that may be premature before product-market fit.

## Evidence 7

Claim: Measurable criteria such as deployment cadence divergence, resource contention, or team ownership shifts constitute valid "Review Triggers" to re-evaluate architectural decisions.

Evidence: Microsoft Azure Well-Architected Framework: Implicit in the recommendation to document confidence levels and to avoid hiding consequences; a decision made with low confidence or one whose consequences have changed (e.g., reporting deployment cadence diverging from core ERP) would trigger re-evaluation.

Source: Maintain an architecture decision record (ADR) - Microsoft Azure Well-Architected Framework
URL: https://learn.microsoft.com/en-us/azure/well-architected/architect-role/architecture-decision-record
Confidence: MEDIUM (explicit examples not in source but strongly implied by confidence logging and consequence tracking)
Corroborated By:
- Joel Parker Henderson repo (Review Triggers section in many ADR examples): While not inspected as primary source, the aggregation shows community practice of including review triggers.
- AWS Prescriptive Guidance: The review process itself implies that changed conditions warrant new ADRs.

Notes: More explicit examples would require inspecting actual ADRs in the wild; the principle is well-established that ADRs should record when they should be revisited.