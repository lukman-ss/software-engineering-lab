# Research Plan: SLO, SLI & Error Budget

## Research Topic

Service Level Indicators (SLI), Service Level Objectives (SLO), and Error Budget in software reliability engineering — particularly in the context of practical implementation, calculation methods, and decision-making frameworks.

## Objective

Provide a comprehensive, evidence-based understanding of:

1. Definition and distinction between SLI, SLO, and Error Budget
2. Practical calculation methods and formulas
3. Industry best practices for setting reliability targets
4. How Error Budget enables engineering product and engineering stability tradeoffs
5. Common pitfalls and anti-patterns in SLO implementation

## Research Questions

1. What is the authoritative definition of SLI, SLO, and Error Budget in SRE practice?
2. How are SLI and SLO mathematically calculated in production systems?
3. What frameworks exist for determining appropriate reliability targets based on business impact?
4. How does Error Budget function as a decision-making tool for deployment risk?
5. What are the most common implementation errors, and what evidence supports them?

## Search Strategy

- Primary: Google SRE books and documentation
- Secondary: Industry publications (Netflix, AWS, Datadog, Prometheus)
- Academic: Any peer-reviewed papers on SLO frameworks
- Community: Relevant discussions (with verification from primary sources)

## Expected Primary Sources

1. *Google SRE Book* (O'Reilly)
2. Google SRE documentation
3. Cloud vendor SLO guides (AWS, GCP, Azure)
4. CNCF projects (Prometheus, OpenTelemetry) documentation
5. OpenTelemetry semantic conventions for service level metrics

## Risks / Unknowns

- Language barrier: Topic specification provided in Bahasa Indonesia, but authoritative sources are primarily in English
- Some claims in topic specification (e.g., specific uptime calculations) require verification
- Practical implementation examples may be scattered across blog posts — must prioritize authoritative sources
