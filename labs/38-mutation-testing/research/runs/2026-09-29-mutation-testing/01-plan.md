# Research Plan: Mutation Testing

## Research Topic
Mutation Testing — Menguji Kualitas Test Suite dan Mengungkap False Confidence 100% Coverage

## Objective
Investigate mutation testing principles, tools, empirical evidence comparing mutation score vs code coverage, practical challenges (equivalent mutants, cost), and available Go ecosystem tools — to produce a structured research report for Lab 38 engineering implementation.

## Research Questions
1. **Definition & History**: What is mutation testing? Who invented it? What are the core hypotheses (competent programmer, coupling effect)?
2. **Mechanics**: How do mutation operators work? What are standard operator categories (arithmetic, relational, conditional, statement)? What is RIP model (Reach, Infect, Propagate)?
3. **Mutation Score**: How is it calculated? What do killed/survived/equivalent/not-viable mutants mean?
4. **Coverage vs Mutation Score**: What empirical evidence shows high code coverage can coexist with low mutation score? What specific test weakness patterns exist?
5. **Tools Landscape**: What mature mutation testing tools exist per language (Java/PIT, JS/Stryker, Python/mutmut, Go/gremlins/ooze/go-mutesting, Rust/cargo-mutants, PHP/Infection, C++/Mull)?
6. **Go Ecosystem**: Which Go mutation tools are actively maintained? What mutators do they support? How do they integrate with `go test`?
7. **Practical Challenges**: Equivalent mutant problem — prevalence, detection techniques. Cost/runtime overhead. CI integration strategies. False positives.
8. **Best Practices**: How to strengthen assertions to kill survived mutants. When to use blacklisting/annotations. Mutation score thresholds.

## Search Strategy
- Primary sources: Official tool documentation (PIT, Stryker, Gremlins, Ooze, go-mutesting), Wikipedia with references, seminal papers (DeMillo 1978, Ammann & Offutt 2008, Jia & Harman 2011).
- Secondary: Expert blogs (Pedro Rijo, Martin Fowler), tool READMEs, GitHub topics.
- Cross-check: Verify tool features against multiple sources (e.g., Go tools listed in awesome-mutation-testing, GitHub topic search).

## Expected Primary Sources (Tier 1)
- DeMillo, Lipton, Sayward (1978) "Hints on test data selection"
- Ammann & Offutt (2008) "Introduction to Software Testing" — textbook
- PIT documentation (pitest.org)
- Stryker Mutator documentation (stryker-mutator.io)
- Go tool repos: go-gremlins/gremlins, gtramontina/ooze, avito-tech/go-mutesting
- Jia & Harman (2011) "An Analysis and Survey of the Development of Mutation Testing" IEEE TSE

## Risks / Unknowns
- Empirical studies directly comparing coverage % vs mutation score % in industry — may be scarce.
- Equivalent mutant prevalence statistics — mostly academic, limited industry data.
- Go mutation tool maturity — most are pre-1.0 (0.x versions).
- Cost benchmarks for mutation testing at scale (CI time multipliers) — tool-specific, not standardized.