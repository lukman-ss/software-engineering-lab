# Contradictions

## No Material Contradictions Discovered

The primary sources (Nygard 2011, MADR specification, AWS guidance, Zimmermann 2020) show strong convergence on core ADR concepts:

- **Structure**: All agree on Context, Decision, Consequences as mandatory
- **Purpose**: All identify preserving rationale and preventing re-debate as primary goal
- **Lifecycle**: All support Proposed/Accepted/Superseded states with immutable history
- **Scope**: All emphasize "architecturally significant" decisions only

## Areas of Divergence (Not Contradictions)

### Template Section Count
- **Nygard**: 5 sections (Title, Context, Decision, Status, Consequences)
- **MADR Minimal**: 4 sections (Title, Context, Decision, Consequences)
- **Y-Statement (Zimmermann)**: 6 parts in single sentence format
- **arc42 Section 9**: 9 tips for ADRs, not a template per se

**Assessment**: These are format variations on the same core elements, not contradictory guidance. All sources explicitly state "pick one and stick to it."

### Microservices Adoption Strategy
- **Fowler (2015)**: Strongly advocates monolith-first for greenfield projects
- **Fowler (2015)**: Acknowledges "by no means unanimous" - some argue for starting with microservices to build organizational capability
- **AWS ADR guidance**: Neutral on architecture style, focuses on documenting whatever decision is made

**Assessment**: Fowler explicitly acknowledges dissenting views and states evidence is sparse. This is honest uncertainty, not contradiction. The counter-argument (start with microservices to build distributed-system muscle) is noted but Fowler leans monolith-first based on observed patterns.

### ADR Granularity
- **Nygard**: "One ADR describes one significant decision"
- **Zimmermann**: Warns "an AD log with more than 100 entries will probably put your readers... to sleep"
- **AWS**: Focuses on "strategic directions for a project or product"

**Assessment**: Consistent guidance on filtering for significance. No contradiction on principle; difference only in threshold calibration.

## Uncertainty Areas Requiring Judgment

1. **What counts as "architecturally significant"?** - Sources define criteria (structure, non-functional, dependencies, interfaces, construction) but threshold is contextual.

2. **How to handle ADR reviews?** - AWS suggests "defining a process," Nygard mentions "changes in project's context" trigger review, but no standardized review cadence or process is specified.

3. **ADR tooling vs. plain files** - MADR emphasizes tool support, Nygard/AWS emphasize simplicity and markdown. No source claims tooling is required.

4. **Team size scaling** - Sources primarily address small-to-medium teams. No evidence on ADR practices at 100+ engineer scale.

These are gaps in prescriptive guidance, not contradictions between sources.