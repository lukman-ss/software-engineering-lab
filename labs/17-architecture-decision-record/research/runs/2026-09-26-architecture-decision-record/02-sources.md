# Sources

Research date: 2026-09-26. All URLs opened and inspected (not snippets).

## Source 1

Title: Documenting Architecture Decisions
Publisher: Cognitect Blog (Michael Nygard)
URL: https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions
Published: 2011-11-15
Accessed: 2026-09-26
Source Tier: Tier 1 (original foundational post, primary historical evidence)
Relevance: Defines ADR anatomy (Title, Context, Decision, Status, Consequences), numbering, immutability via supersede, co-location in version control.

## Source 2

Title: Architectural decision record process
Publisher: AWS Prescriptive Guidance (Amazon Web Services)
URL: https://docs.aws.amazon.com/prescriptive-guidance/latest/architectural-decision-records/adr-process.html
Published: NOT VERIFIED (continuously updated guidance; page inspected 2026-09-26)
Accessed: 2026-09-26
Source Tier: Tier 1 (official cloud vendor guidance)
Relevance: Defines ADR scope (structure, NFRs, dependencies, interfaces, construction techniques), lifecycle (Proposed/Accepted/Rejected/Superseded), immutability, review/ownership process.

## Source 3

Title: Architectural Decision Records (ADRs) — Homepage of the ADR GitHub organization
Publisher: adr.github.io
URL: https://adr.github.io/
Published: NOT VERIFIED (continuously updated index)
Accessed: 2026-09-26
Source Tier: Tier 1 (central architectural knowledge management index)
Relevance: Definitions of AD / ASR / ADR / ADL / AKM; template index; background literature pointers (Nygard, Zdun et al., WICSA 2015 template comparison).

## Source 4

Title: Markdown Architectural Decision Records (MADR)
Publisher: adr/madr (MADR project site)
URL: https://adr.github.io/madr/
Published: Scientific publication 2018-04-03; template versions MADR 3.0.0 on 2022-10-09, MADR 4.0.0 on 2024-09-17 (per News section on same page)
Accessed: 2026-09-26
Source Tier: Tier 1 (standards-adjacent template + peer-reviewed publication reference)
Relevance: Full template anatomy (Context and Problem Statement, Decision Drivers, Considered Options, Decision Outcome, Consequences, Confirmation, Pros and Cons); filename convention NNNN-title-with-dashes.md; linting guidance.

## Source 5

Title: Maintain an architecture decision record (ADR)
Publisher: Microsoft Learn — Azure Well-Architected Framework (Architect role)
URL: https://learn.microsoft.com/en-us/azure/well-architected/architect-role/architecture-decision-record
Published: 2026-04-10 (ms.date in page metadata; updated 2026-04-13 per metadata)
Accessed: 2026-09-26
Source Tier: Tier 1 (official vendor framework)
Relevance: Append-only log rule; record confidence level; break multi-phase decisions into multiple records; avoid design-guide bloat; store openly with workload docs.

## Source 6

Title: Monolith First
Publisher: martinfowler.com (Martin Fowler)
URL: https://martinfowler.com/bliki/MonolithFirst.html
Published: 2015-06-03
Accessed: 2026-09-26
Source Tier: Tier 1 (authoritative practitioner source; Tier 2 by strict academic tiering but treated as expert technical article)
Relevance: Primary evidence for modular-monolith-first strategy: MicroservicePremium, YAGNI, boundary-discovery difficulty, prerequisite capabilities; supports ERP scenario rationale without claiming monolith always wins.

## Source 7

Title: Architecture decision record (ADR) — repository README and template collection
Publisher: GitHub — architecture-decision-record/architecture-decision-record (Joel Parker Henderson, community-maintained)
URL: https://github.com/joelparkerhenderson/architecture-decision-record
Published: NOT VERIFIED (living repository; star/fork counts volatile, not used as evidence)
Accessed: 2026-09-26
Source Tier: Tier 2 (reputable community reference aggregating Nygard, MADR, Tyree/Akerman, arc42 templates)
Relevance: Good-ADR characteristics (rationale, single-decision scope, timestamps, immutability), Context/Consequences writing guidance, when to / not to raise an ADR; used only as corroboration, never sole evidence.
