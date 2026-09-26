# Evidence

## Evidence 1: ADR Structure Definition

**Claim**: ADRs must contain Title, Context, Decision, Status, and Consequences sections to be effective.

**Evidence**: Michael Nygard's original blog post explicitly defines ADR structure: "The format has just a few parts. Title: These documents have names that are short noun phrases. Context: This section describes the forces at play, including technological, political, social, and project local. Decision: This section describes our response to these forces. Status: A decision may be 'proposed' if the project stakeholders haven't agreed with it yet, or 'accepted' once it is agreed. Consequences: This section describes the resulting context, after applying the decision."

**Source**: Nygard (2011), "Documenting Architecture Decisions"  
**URL**: https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions  
**Published**: 2011-11-15  
**Confidence**: HIGH  
**Corroborated By**: MADR template (Source 3), AWS Prescriptive Guidance (Source 5)

**Notes**: Multiple independent sources converge on same section names. This is the canonical ADR format.

---

## Evidence 2: ADR Prevents Repeated Debates

**Claim**: ADRs prevent teams from repeatedly re-discussing the same architecture decisions by preserving decision rationale.

**Evidence**: AWS documentation states: "Three major anti-patterns often emerge when making architectural decisions: A decision is made without any justification, and people don't understand why it was made. This results in the same topic being discussed multiple times."

**Source**: AWS Prescriptive Guidance (Kunce, Goby)  
**URL**: https://docs.aws.amazon.com/prescriptive-guidance/latest/architectural-decision-records/introduction.html  
**Published**: Recent (2026)  
**Confidence**: HIGH  
**Corroborated By**: Nygard (2011) - "A new person coming on to a project may be perplexed, baffled, delighted, or infuriated by some past decision."

**Notes**: Both sources identify the same failure mode: without documented rationale, decision context is lost and re-litigated.

---

## Evidence 3: ADR Records "Architecturally Significant" Decisions

**Claim**: ADRs should document only architecturally significant decisions, not all technical choices.

**Evidence**: Nygard writes: "We will keep a collection of records for 'architecturally significant' decisions: those that affect the structure, non-functional characteristics, dependencies, interfaces, or construction techniques."

**Source**: Nygard (2011)  
**URL**: https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions  
**Published**: 2011-11-15  
**Confidence**: HIGH  
**Corroborated By**: Zimmermann (2020) - "Focus on the architecturally significant requirements and decisions — the ones that matter, the ones that are hard and costly to change."

**Notes**: Both sources emphasize filtering criteria: significance, cost to change, structural impact.

---

## Evidence 4: ADR Lifecycle States

**Claim**: ADRs have lifecycle states (Proposed, Accepted, Deprecated, Superseded) and old ADRs are preserved even when superseded.

**Evidence**: Nygard explicitly states: "A decision may be 'proposed' if the project stakeholders haven't agreed with it yet, or 'accepted' once it is agreed. If a later ADR changes or reverses a decision, it may be marked as 'deprecated' or 'superseded' with a reference to its replacement." Also: "If a decision is reversed, we will keep the old one around, but mark it as superseded."

**Source**: Nygard (2011)  
**URL**: https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions  
**Published**: 2011-11-15  
**Confidence**: HIGH  
**Corroborated By**: MADR templates define status field with these states.

**Notes**: Lifecycle management prevents history rewriting and enables audit trail.

---

## Evidence 5: Microservices Have Operational Premium

**Claim**: Microservices architecture introduces significant operational complexity premium compared to monoliths.

**Evidence**: Martin Fowler: "microservices introduce complexity on their own account. This adds a premium to a project's cost and risk - one that often gets projects into serious trouble." Also: "When you use microservices you have to work on automated deployment, monitoring, dealing with failure, eventual consistency, and other factors that a distributed system introduces."

**Source**: Fowler (2015), "Microservice Premium"  
**URL**: https://martinfowler.com/bliki/MicroservicePremium.html  
**Published**: 2015-05-13  
**Confidence**: HIGH  
**Corroborated By**: Fowler (2015) "Monolith First" - "the premium of microservices is a drag you should do without" for early-stage projects.

**Notes**: Consistent position across two articles by same author. Trade-off framing, not absolute claim.

---

## Evidence 6: Successful Microservices Typically Evolve from Monoliths

**Claim**: Most successful microservice implementations started as monoliths that later decomposed, not greenfield microservices.

**Evidence**: Fowler: "Almost all the successful microservice stories have started with a monolith that got too big and was broken up. Almost all the cases where I've heard of a system that was built as a microservice system from scratch, it has ended up in serious trouble."

**Source**: Fowler (2015), "Monolith First"  
**URL**: https://martinfowler.com/bliki/MonolithFirst.html  
**Published**: 2015-06-03  
**Confidence**: MEDIUM  
**Corroborated By**: Qualitative evidence only (anecdotal reports from practitioner network).

**Notes**: This is practitioner observation, not empirical study. Stated as strong pattern but acknowledged as based on anecdotes.

---

## Evidence 7: Early-Stage Projects Should Prioritize Speed

**Claim**: Early-stage projects with uncertain requirements should prioritize development speed over architectural sophistication.

**Evidence**: Fowler: "When you begin a new application, how sure are you that it will be useful to your users?... During this first phase you need to prioritize speed (and thus cycle time for feedback), so the premium of microservices is a drag you should do without."

**Source**: Fowler (2015), "Monolith First"  
**URL**: https://martinfowler.com/bliki/MonolithFirst.html  
**Published**: 2015-06-03  
**Confidence**: HIGH  
**Corroborated By**: Consistent with YAGNI principle (referenced in same article).

**Notes**: Aligns with agile principles of iterative development and validation-first approach.

---

## Evidence 8: ADR Template Comparison

**Claim**: Multiple ADR templates exist with varying sections, but core elements remain consistent across templates.

**Evidence**: WICSA 2015 paper compares 7 ADR templates. Zimmermann's article shows Y-statement template (6 lines) vs Nygard's 5-section template vs arc42 template, noting: "Many templates exist, just pick one and stick to it."

**Source**: Zimmermann (2020), "Architectural Decisions — The Making Of"  
**URL**: https://ozimmer.ch/practices/2020/04/27/ArchitectureDecisionMaking.html  
**Published**: 2020-04-27 (Updated 2026)  
**Confidence**: HIGH  
**Corroborated By**: MADR (Source 3) provides multiple template variants (full, minimal, bare).

**Notes**: Convergence on core concept (context/decision/consequences) despite structural variation.

---

## Evidence 9: ADR as Architecture Git History

**Claim**: ADRs serve as "architecture git history" - recording why architecture changed, complementing source control that records how code changed.

**Evidence**: Nygard: "The motivation behind previous decisions is visible for everyone, present and future. Nobody is left scratching their heads to understand, 'What were they thinking?' and the time to change old decisions will be clear from changes in the project's context."

**Source**: Nygard (2011)  
**URL**: https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions  
**Published**: 2011-11-15  
**Confidence**: MEDIUM  
**Corroborated By**: AWS (Source 5) - ADRs provide "hand-over experience and reference documentation."

**Notes**: "Architecture Git History" is an interpretation/analogy, not direct quote from sources.

---

## Evidence 10: ADR Quality Anti-patterns

**Claim**: Common ADR anti-patterns include: no context, fake alternatives, documentation dumps, missing consequences, no review conditions.

**Evidence**: AWS documentation lists anti-patterns: "A decision is made without any justification... The decision isn't captured in an architectural decision repository." Zimmermann identifies bad justifications: "'Everybody does it.' 'We have always done it like that.' 'Experience with this... will look fantastic on my resume.'"

**Source**: AWS (Source 5), Zimmermann (Source 4)  
**URLs**: https://docs.aws.amazon.com/prescriptive-guidance/latest/architectural-decision-records/introduction.html, https://ozimmer.ch/practices/2020/04/27/ArchitectureDecisionMaking.html  
**Published**: AWS (2026), Zimmermann (2020)  
**Confidence**: HIGH  
**Corroborated By**: Multiple sources identify overlapping anti-pattern categories.

**Notes**: Convergent guidance from industry (AWS) and practitioner/academic (Zimmermann) sources.