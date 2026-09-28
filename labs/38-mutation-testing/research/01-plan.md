# Research Plan

## Research Topic

Mutation Testing — Menguji Kualitas Test Suite dan Mengungkap False Confidence 100% Coverage

## Objective

To investigate mutation testing methodologies, mutation operators, tools, and practical applications to prevent false confidence in test suites despite achieving 100% code coverage, focusing on identifying weak tests and improving regression test effectiveness.

## Research Questions

1. What is mutation testing and how does it differ from traditional code coverage metrics?
2. What are the core theoretical foundations and hypotheses behind mutation testing?
3. What are the common mutation operators and categories (statement, value, decision mutation)?
4. What tools are available for mutation testing (PIT, Stryker, MuJava, ACH) and how do they work?
5. How is the mutation score calculated and what does it indicate about test quality?
6. What are the challenges in practical mutation testing (equivalent mutants, scalability, computational cost)?
7. How does mutation testing fit into CI/CD pipelines and modern development practices?
8. What is the relationship between code coverage, test quality, and mutation testing?
9. How do LLMs enhance mutation testing (Meta ACH approach)?
10. What is the historical development and academic foundation of mutation testing?

## Search Strategy

- Search for official documentation of mutation testing tools (PIT, Stryker, MuJava).
- Seek authoritative academic papers (DeMillo, Lipton, Sayward 1978; Jia & Harman 2009).
- Check for engineering blog posts from industry leaders (Meta, ThoughtWorks, PIT authors).
- Use search terms: "mutation testing", "mutation score", "killed vs survived mutants", "equivalent mutants", "mutation operators", "PIT", "Stryker".
- Cross-reference Wikipedia with primary sources.

## Expected Primary Sources

- DeMillo, Lipton, and Sayward (1978) IEEE Computer paper — original mutation testing proposal
- PIT documentation and FAQ — https://pitest.org/
- Stryker documentation — https://stryker-mutator.io/docs/
- Wikipedia Mutation Testing — https://en.wikipedia.org/wiki/Mutation_testing
- Meta Engineering blog on ACH — https://engineering.fb.com/2025/09/30/security/llms-are-the-key-to-mutation-testing-and-better-compliance/
- Martin Fowler bliki on Mutation Testing
- ACM/IEEE academic papers on mutation testing improvements

## Risks / Unknowns

- Mutation testing theory is over 45 years old but practical tooling only became viable recently.
- Equivalent mutant detection remains a mathematically undecidable problem.
- Tool coverage varies by language (PIT for JVM, Stryker for JS/TS, limited Go support).
- Academic papers may be paywalled; need to rely on open-access sources and public implementations.
- The topic specification mentions Go but mutation testing tools for Go are limited compared to JVM/JS ecosystems.
- Meta ACH approach (Sep 2025) is very recent; peer review status unknown.
