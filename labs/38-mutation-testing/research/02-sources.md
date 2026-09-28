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

## Source 10

Title: go-mutesting — Mutation testing for Go source code
Publisher: zimmski (GitHub repository)
URL: https://github.com/zimmski/go-mutesting
Published: README last updated per repo activity (ongoing project)
Accessed: 2026-09-28
Source Tier: Tier 2 (community open-source project)
Relevance: Go mutation testing framework with mutators for branch, expression, and statement levels. Implements mutation score calculation and mutation output.
Notes: Available at https://github.com/zimmski/go-mutesting. Provides an exec-based mutation workflow for Go.

## Source 11

Title: gremlins — A mutation testing tool for Go
Publisher: go-gremlins (GitHub organization)
URL: https://github.com/go-gremlins/gremlins
Homepage: https://gremlins.dev
Published: Ongoing project
Accessed: 2026-09-28
Source Tier: Tier 2 (community open-source project)
Relevance: Another Go mutation testing tool with CLI, documentation site, and support for Go mutation workflows.
Notes: Available at https://github.com/go-gremlins/gremlins with homepage at https://gremlins.dev.

## Source 12

Title: Mutation-Guided LLM-based Test Generation at Meta (arXiv preprint)
Publisher: arXiv (cs.SE, cs.AI, cs.LG)
URL: https://arxiv.org/abs/2501.12862
Published: Submitted 22 January 2025
Accessed: 2026-09-28
Source Tier: Tier 1 (primary academic preprint; peer-reviewed venue: FSE 2025 Industry Track)
Relevance: Full arXiv abstract corroborates Meta ACH trial statistics from the blog post; authors include Mark Harman, Christopher Foster, Abhishek Gulati, et al. Provides exact numbers: 9,095 mutants, 571 privacy-hardening test cases across 10,795 Android Kotlin classes on 7 software platforms; LLM equivalence detector 0.79 precision / 0.47 recall, rising to 0.95 / 0.96 with preprocessing; 73% engineer acceptance; 36% privacy-relevant.
Notes: Abstract text directly accessible at https://arxiv.org/abs/2501.12862. The PDF was not fetched directly (binary content); all statistics verified from the abstract page.

## Source 13

Title: Mutation testing (Wikipedia — subsumed mutants section)
Publisher: Wikipedia
URL: https://en.wikipedia.org/wiki/Mutation_testing
Published: Last edited 6 September 2026
Accessed: 2026-09-28
Source Tier: Tier 2 (crowdsourced encyclopedia)
Relevance: Defines subsumed mutants — mutants that exist at the same source location as another mutant and are "subsumed" by the other; they do not contribute to coverage metrics. Provides concrete example.
Notes: Used to address Gap 8 (missing subsumed mutants case).
