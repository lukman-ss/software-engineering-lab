# Open Questions

## Unanswered Questions

### 1. Optimal ADR Review Cadence
**Question**: What cadence should teams review existing ADRs against changing context?  
**Current Evidence**: Sources acknowledge review triggers exist but do not specify standardized review schedules.  
**Next Research Direction**: Survey engineering organizations on ADR review practices and effective trigger thresholds.

### 2. ADR Tooling Impact at Scale
**Question**: How do ADR practices differ between small teams (5 engineers) vs large organizations (100+ engineers)?  
**Current Evidence**: All sources address small-to-medium teams. No evidence on scaling effects.  
**Next Research Direction**: Investigate ADR practices at companies with >50 engineering teams.

### 3. Quantitative ADR ROI
**Question**: What measurable impact do ADRs have on team velocity, onboarding time, or decision quality?  
**Current Evidence**: No sources provide metrics or controlled studies.  
**Next Research Direction**: Collect survey data from teams using ADRs vs not using ADRs.

### 4. ADR Validation Automation
**Question**: What validation patterns ensure ADR quality beyond manual review?  
**Current Evidence**: Lab spec asks for Go-based validator but sources don't establish best practices.  
**Next Research Direction**: Research existing ADR linting tools (adr-tools, MADR linter) and their rule sets.

### 5. Modular Monolith Migration Path
**Question**: What are the concrete migration strategies for extracting services from a modular monolith?  
**Current Evidence**: Fowler acknowledges "gradual peel-off" works but notes "substantial monolith remains."  
**Next Research Direction**: Study cases of successful monolith-to-microservices extraction with explicit boundaries.

## Weak Evidence Areas

### ADR Granularity Threshold
- Sources define "architecturally significant" criteria but provide no quantitative thresholds
- Different organizations will have different tolerance levels
- Risk: Either too many trivial ADRs or too few critical ones

### ADR vs Regular Documentation Boundary
- Sources state ADRs are not for "full API documentation" or "database schema documentation"
- No clear boundary on what constitutes "architecturally significant" for specific decision types
- Risk: Over-documentation or gaps in critical decisions

### Team Culture Prerequisites
- ADRs require cultural adoption (teams actually read and write them)
- Sources assume willing participation
- Risk: ADRs become shelfware without enforcement mechanisms

## Claims Needing Deeper Research

### Microservices vs Monolith Metrics
- Need empirical data comparing developer productivity at different stages
- Fowler's claims are based on anecdotes, not controlled comparison

### Decision Fatigue in ADR Writing
- Does requiring ADRs for all architecture decisions slow teams?
- No evidence on the optimal number of ADRs per project

### ADR as Hiring Filter
- Can ADRs serve as interview material to assess new engineer context understanding?
- No documented evidence of this practice being used effectively

## Research Date
- Research conducted: 2026-09-26
- Information may change as cloud architecture practices evolve (especially container orchestration, service meshes, and serverless architectures).