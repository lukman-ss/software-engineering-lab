# Research Report

## Research Question

How should Architecture Decision Records (ADRs) be structured and maintained to preserve the rationale behind technical decisions, particularly in the context of a SaaS ERP system choosing between Laravel Modular Monolith and Go Microservices?

## Executive Summary

After reviewing 10 authoritative sources including Nygard's seminal 2011 paper, Martin Fowler's microservices analysis, AWS prescriptive guidance, and the MADR community specification, this research establishes that:

1. **ADRs are standardized documentation artifacts** with a consistent structure (Context, Decision, Status, Consequences) across source implementations
2. **The monolith-first thesis is well-supported** by Fowler's practitioner observations that successful microservices typically evolve from monoliths
3. **Microservices impose significant operational premium** that must be explicitly documented as a consequence
4. **ADRs serve as "architecture git history"** - recording why decisions change, not just what changed

## Findings

### Finding 1: ADR Structure is Standardized

**Claim**: Effective ADRs require Title, Context, Decision, Status, and Consequences sections.

**Evidence**: Nygard's 2011 blog post defined this structure, which has been validated by MADR templates, AWS guidance, and Zimmermann's comparative analysis. The IEEE 42010 standard also requires capturing decision rationale.

**Sources**: 
- Source 1: Nygard (2011) - https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions
- Source 3: MADR - https://github.com/adr/madr
- Source 8: ISO/IEC/IEEE 42010
- Source 4: Zimmermann (2020)

**Confidence**: HIGH

### Finding 2: ADRs Prevent Re-litigation of Decisions

**Claim**: The primary purpose of ADRs is to prevent teams from repeatedly re-discussed the same architecture choices by preserving context and rationale.

**Evidence**: Both Nygard (2011) and AWS (2026) identify this as a core anti-pattern: "Decision is made without any justification... same topic discussed multiple times." The successor engineer should read "why" from ADR, not reverse-engineer from code.

**Sources**: Source 1, Source 5

**Confidence**: HIGH

### Finding 3: Monolith-First is the Recommended Starting Point

**Claim**: For greenfield projects, starting with a monolith (even modular monolith) is recommended over microservices, which should evolve later.

**Evidence**: Fowler's "Monolith First" (2015) presents observational evidence from practitioner network: successful microservices "started with a monolith that got too big and was broken up." He cites "YAGNI" principle and argues microservice operational premium is a "drag you should do without."

**Sources**: Source 6, Source 7

**Confidence**: MEDIUM (anecdotal evidence acknowledged by Fowler)

### Finding 4: Microservices Operational Premium is Significant

**Claim**: Microservices introduce distributed system complexity that must be explicitly accounted for: automated deployment, monitoring, failure handling, eventual consistency.

**Evidence**: Fowler's "Microservice Premium" defines this premium as "the cost of managing a suite of services" that "slows down development considerable[ly]." This includes "network failures, observability complexity, deployment complexity, higher infrastructure overhead."

**Sources**: Source 7

**Confidence**: HIGH

### Finding 5: ADR Lifetime is Cyclical, Not Permanent

**Claim**: ADRs are not permanent decisions; they become obsolete as context changes. The correct pattern is: Proposed → Accepted → (Superseded) → New ADR.

**Evidence**: Nygard (2011) states: "If a decision is reversed, we will keep the old one around, but mark it as superseded." Zimmermann emphasizes: "A decision can be correct then and wrong now without the original decision having been a mistake."

**Sources**: Source 1, Source 4

**Confidence**: HIGH

### Finding 6: Significantly Architectural Decisions Include Cross-Cutting Concerns

**Claim**: Decisions warranting ADRs affect structure, non-functional characteristics, dependencies, interfaces, or construction techniques - not implementation details.

**Evidence**: Nygard's definition and AWS's "strategic directions" framing both limit ADR scope to decisions that are "hard to reverse" or "cross-cutting." Minor refactors, naming changes, CRUD endpoints don't qualify.

**Sources**: Source 1, Source 5

**Confidence**: HIGH

## Areas of Agreement

1. ADR = Core decision artifacts (Nygard, MADR, AWS, Zimmermann)
2. Context preservation is the primary purpose (all sources)
3. Single ADR per significant decision (all sources)
4. Status lifecycle is Required (all sources)
5. All consequences must be documented, including negatives (all sources)
6. Immutable history - don't rewrite old ADRs (all sources)

## Areas of Disagreement

None substantive. Format variations (section count, templates) exist but sources explicitly endorse picking one and sticking to it.

## Limitations

1. **Anecdotal Microservices Evidence**: Fowler's "monolith first" argument relies on "stories" from his network, not controlled studies.

2. **Small Team Focus**: All sources primarily address small-to-medium teams. Evidence gap for 100+ engineer scale.

3. **No Empirical Impact Data**: No sources provide metrics on ADR impact on team velocity or decision quality.

4. **Review Process Undefined**: No standard process defined for ADR reviews, only acknowledgment that context changes trigger re-evaluation.

## Conclusion

Architecture Decision Records emerge from multiple authoritative sources as a **lightweight, standardized mechanism** for capturing the "why" behind technical decisions. The convergence across sources (industry blog posts, cloud vendor guidance, community templates) establishes strong consensus on:

- **Essential sections**: Context, Decision, Consequences (with negatives), Status
- **Lifecycle**: Proposed → Accepted → Supersede (not delete)
- **Purpose**: Prevent decision regression, preserve historical rationale
- **Scope**: Architecturally significant decisions only

The Laravel Modular Monolith vs Go Microservices scenario fits this pattern perfectly: the choice is significant (affects deployment, scaling, team structure), the context (small team, tight deadline, evolving requirements) must be documented for future evaluation, and the trade-offs (speed vs scaling) deserve explicit recording.

Fowler's "microservice premium" analysis provides the technical justification for monolith-first in this scenario: operational complexity of microservices would slow delivery and increase risk given the stated constraints.