# Evidence

## Evidence 1
Claim: An ADR must record the motivation, forces, single decision, and full consequences (positive, negative, and neutral) rather than just implementation specs.
Evidence: Michael Nygard specifies: "Each record describes a set of forces and a single decision in response to those forces... This section describes the resulting context, after applying the decision. All consequences should be listed here, not just the 'positive' ones."
Source: Documenting Architecture Decisions
URL: https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions
Confidence: HIGH
Corroborated By: AWS Prescriptive Guidance ("At a minimum, each ADR should define the context of the decision, the decision itself, and the consequences of the decision... focuses on the reason for the decision rather than how the team implemented it.")
Notes: Crucial for explaining why technical debates restart if rationale is absent.

## Evidence 2
Claim: ADRs must be immutable once accepted; changes require a new ADR that supersedes or deprecates the predecessor.
Evidence: AWS Prescriptive Guidance: "When the team accepts an ADR, it becomes immutable. If new insights require a different decision, the team proposes a new ADR. When the team accepts the new ADR, it supersedes the previous ADR."
Source: Architectural Decision Record Process
URL: https://docs.aws.amazon.com/prescriptive-guidance/latest/architectural-decision-records/adr-process.html
Confidence: HIGH
Corroborated By: Michael Nygard (2011) ("If a decision is reversed, we will keep the old one around, but mark it as superseded.")
Notes: Historical continuity is essential; editing accepted ADRs distorts historical context.

## Evidence 3
Claim: ADRs belong inside version control alongside source code.
Evidence: Nygard states: "We will keep ADRs in the project repository under doc/arch/adr-NNN.md... keeping these in version control with the code makes them less accessible... In practice, our projects almost all live in GitHub... it looks just as friendly as any wiki page would."
Source: Documenting Architecture Decisions
URL: https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions
Confidence: HIGH
Corroborated By: adr.github.io
Notes: Prevents decision drift and maintains coupling between code history and architectural history.

## Evidence 4
Claim: Architectural significance applies to structure, NFRs, dependencies, interfaces, and construction techniques, not low-level refactoring.
Evidence: AWS Prescriptive Guidance / Richards and Ford (2020): "Project members should create an ADR for every architecturally significant decision... Structure, Non-functional requirements, Dependencies, Interfaces, Construction techniques."
Source: Architectural Decision Record Process
URL: https://docs.aws.amazon.com/prescriptive-guidance/latest/architectural-decision-records/adr-process.html
Confidence: HIGH
Corroborated By: adr.github.io ("Architecturally Significant Requirement (ASR) is a requirement that has a measurable effect on the architecture and quality...")
Notes: Filters trivial implementation details (e.g., helper rename, minor CRUD addition).

## Evidence 5
Claim: ADR status transitions follow: Proposed -> Accepted -> Superseded / Deprecated (or Rejected).
Evidence: Nygard / AWS: States define the lifecycle. Proposed ADRs are reviewed; Accepted become active and immutable; changes introduce new records which mark older records as Superseded with a reference to the replacement.
Source: Documenting Architecture Decisions / AWS Prescriptive Guidance
URL: https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions
Confidence: HIGH
Corroborated By: AWS Prescriptive Guidance
Notes: Rejected status is also preserved in AWS process to document negative decisions and eliminate repeated debates.
