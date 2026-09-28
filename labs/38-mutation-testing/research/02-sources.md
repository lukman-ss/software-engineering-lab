# Sources

## Source 1

Title: Hints on test data selection: Help for the practicing programmer
Publisher: IEEE Computer (DeMillo, Lipton, Sayward)
URL: NOT VERIFIED (original paper not opened; DOI not confirmed)
Published: April 1978 (per Wikipedia reference list)
Accessed: 2026-09-28 (bibliographic metadata only, via Wikipedia citation)
Source Tier: Tier 1 (foundational academic paper)
Relevance: Original paper proposing mutation testing; source of mutation score and hypotheses.
Notes: Paper text NOT accessed directly. All claims attributed to this paper in this research come from Wikipedia's reference list and summaries. NOT VERIFIED for direct quotations from the paper itself.

## Source 2

Title: Mutation testing
Publisher: Wikipedia
URL: https://en.wikipedia.org/wiki/Mutation_testing
Published: Last edited 6 September 2026
Accessed: 2026-09-28
Source Tier: Tier 2 (crowdsourced encyclopedia; cross-referenced with primary sources)
Relevance: Overview of definitions, goals, history, hypotheses, RIP model, operators, equivalent mutant problem. Cites primary literature (DeMillo 1978; Ammann & Offutt 2008; Jia & Harman 2009; Offutt 1992; Deng et al. 2013; Madeyski et al. 2014).

## Source 3

Title: Mutation Testing (bliki entry)
Publisher: Martin Fowler
URL: https://martinfowler.com/bliki/MutationTesting.html
Published: Draft (no date shown; page carries "This is a draft entry" notice)
Accessed: 2026-09-28
Source Tier: Tier 1 (authoritative expert)
Relevance: Explains motivation (probing test quality), cost history (Jester era), speedups in modern tools (Pitest, Stryker), manual mutation testing practice, LLM-era relevance.
Notes: Page is explicitly marked DRAFT by the author; treat as pre-publication source.

## Source 4

Title: PIT Mutation Testing (homepage)
Publisher: PIT Project
URL: https://pitest.org/
Published: not stated (site copyright 2023 on FAQ; site current as of access)
Accessed: 2026-09-28
Source Tier: Tier 1 (official tool documentation)
Relevance: Definition of mutation workflow (seed → run tests → killed/lived), critique of line coverage ("measures only which code is executed... does not check that your tests are actually able to detect faults"), mutation coverage vs line coverage reports.

## Source 5

Title: PIT FAQ
Publisher: PIT Project
URL: https://pitest.org/faq/
Published: not stated (site current as of access)
Accessed: 2026-09-28
Source Tier: Tier 1 (official tool documentation)
Relevance: Computational cost guidance, test selection strategy, mutator groups (DEFAULTS/STRONGER/ALL), supported languages (Java, Kotlin), determinism, recommended usage (limit analysis to changed code).

## Source 6

Title: What is mutation testing? (Introduction)
Publisher: Stryker Mutator
URL: https://stryker-mutator.io/docs/
Published: not stated (site current as of access)
Accessed: 2026-09-28
Source Tier: Tier 1 (official tool documentation)
Relevance: Definition of killed/survived mutants, code coverage critique (sandwich/paste analogy), worked example (`user.age >= 18` with `>`, `<`, `false`, `true` mutants), mutator names observed in output (BinaryOperator, RemoveConditionals).

## Source 7

Title: Welcome to the RoboCoasters — An introduction to mutation testing
Publisher: Stryker Mutator
URL: https://stryker-mutator.io/docs/General/example/
Published: not stated (site current as of access)
Accessed: 2026-09-28
Source Tier: Tier 1 (official tool documentation)
Relevance: Concrete demonstration that a project can have 100% code coverage and ~60% mutation score simultaneously ("How code coverage of 100% could mean only 60% is tested"), plus reproducible example repository.

## Source 8

Title: LLMs Are the Key to Mutation Testing and Better Compliance
Publisher: Engineering at Meta (by Mark Harman)
URL: https://engineering.fb.com/2025/09/30/security/llms-are-the-key-to-mutation-testing-and-better-compliance/
Published: 30 September 2025
Accessed: 2026-09-28
Source Tier: Tier 1 (primary industry source for ACH; author is a mutation testing researcher)
Relevance: Five barriers to mutation testing at scale; ACH system; LLM equivalence detection (0.79/0.47 precision/recall; 0.95/0.96 with preprocessing); Meta trial stats (73% acceptance, 36% privacy-relevant); links to arXiv:2501.12862.

## Source 9

Title: An Analysis and Survey of the Development of Mutation Testing (bibliographic record)
Publisher: IEEE Transactions on Software Engineering (Jia & Harman, 2009)
URL: doi:10.1109/TSE.2010.62 (DOI string as printed in Wikipedia references; not opened)
Published: September 2009
Accessed: 2026-09-28 (bibliographic metadata only, via Wikipedia citation)
Source Tier: Tier 1 (peer-reviewed academic survey)
Relevance: Standard survey of mutation testing development.
Notes: Paper NOT opened directly; cite as secondary-attributed only.
