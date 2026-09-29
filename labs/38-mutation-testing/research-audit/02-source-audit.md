# Source Audit

Target Lab: labs/38-mutation-testing
Audit Scope: Research Files Only (PIPELINE OVERRIDE)

---

## Source 1

Claimed Title: Hints on test data selection: Help for the practicing programmer
Claimed Publisher: IEEE Computer (DeMillo, Lipton, Sayward)
URL: NOT VERIFIED (original paper not opened; DOI not confirmed)
Reachable: NO (Explicitly noted as unaccessed by researcher)
Source Type: PRIMARY (Academic paper, but treated as secondary citation)
Relevant: YES
Supports Claimed Topic: YES
Problems:
- The original 1978 text was not directly inspected.
- The research explicitly and honestly disclaims this: "Paper text NOT accessed directly. All claims attributed to this paper in this research come from Wikipedia's reference list and summaries."
Assessment: WARNING (Accurate attribution of limitation, but primary source unverified directly)

---

## Source 2

Claimed Title: Mutation testing
Claimed Publisher: Wikipedia
URL: https://en.wikipedia.org/wiki/Mutation_testing
Reachable: YES
Source Type: SECONDARY (Crowdsourced encyclopedia)
Relevant: YES
Supports Claimed Topic: YES
Problems:
- Secondary source subject to crowdsourced edits, but accurately reflects established academic consensus (DeMillo 1978, Offutt 1992, Ammann & Offutt 2008).
Assessment: PASS

---

## Source 3

Claimed Title: Mutation Testing (bliki entry)
Claimed Publisher: Martin Fowler
URL: https://martinfowler.com/bliki/MutationTesting.html
Reachable: YES
Source Type: COMMUNITY / EXPERT (Draft bliki entry)
Relevant: YES
Supports Claimed Topic: YES
Problems:
- Explicitly carries "This is a draft entry" on the page.
- The research properly qualifies every inline reference as "(pre-publication draft; carries 'This is a draft entry' notice)".
Assessment: WARNING (Draft status noted; reliable author commentary)

---

## Source 4

Claimed Title: PIT Mutation Testing (homepage)
Claimed Publisher: PIT Project
URL: https://pitest.org/
Reachable: YES (Verified via WebFetch)
Source Type: PRIMARY (Official tool documentation)
Relevant: YES
Supports Claimed Topic: YES
Problems:
- None. Text matches verified quotes regarding line coverage limitations and mutation testing principles.
Assessment: PASS

---

## Source 5

Claimed Title: PIT FAQ
Claimed Publisher: PIT Project
URL: https://pitest.org/faq/
Reachable: YES
Source Type: PRIMARY (Official tool documentation)
Relevant: YES
Supports Claimed Topic: YES
Problems:
- None. Authoritative FAQ on computational cost, Java/Kotlin support, and mutator sets.
Assessment: PASS

---

## Source 6

Claimed Title: What is mutation testing? (Introduction)
Claimed Publisher: Stryker Mutator
URL: https://stryker-mutator.io/docs/
Reachable: YES (Verified via WebFetch)
Source Type: PRIMARY (Official tool documentation)
Relevant: YES
Supports Claimed Topic: YES
Problems:
- None. Accurately contains the sandwich/chocolate paste analogy, age verification code example, and mutator reporting format.
Assessment: PASS

---

## Source 7

Claimed Title: Welcome to the RoboCoasters — An introduction to mutation testing
Claimed Publisher: Stryker Mutator
URL: https://stryker-mutator.io/docs/General/example/
Reachable: YES
Source Type: PRIMARY (Official tool documentation)
Relevant: YES
Supports Claimed Topic: YES
Problems:
- None. Concrete example repository demonstrating 100% code coverage with 60% mutation score.
Assessment: PASS

---

## Source 8

Claimed Title: LLMs Are the Key to Mutation Testing and Better Compliance
Claimed Publisher: Engineering at Meta (by Mark Harman)
URL: https://engineering.fb.com/2025/09/30/security/llms-are-the-key-to-mutation-testing-and-better-compliance/
Reachable: YES (Verified via WebFetch)
Source Type: PRIMARY (Official engineering blog from industry author)
Relevant: YES
Supports Claimed Topic: YES
Problems:
- Single-vendor corporate engineering publication detailing internal trials.
Assessment: PASS

---

## Source 9

Claimed Title: An Analysis and Survey of the Development of Mutation Testing (bibliographic record)
Claimed Publisher: IEEE Transactions on Software Engineering (Jia & Harman, 2009)
URL: doi:10.1109/TSE.2010.62
Reachable: NO (Bibliographic metadata only via Wikipedia citation; paper not directly read)
Source Type: PRIMARY (Academic survey, treated as secondary-attributed)
Relevant: YES
Supports Claimed Topic: YES
Problems:
- Original survey text was not fetched or inspected directly.
- Transparently acknowledged by researcher as unverified direct text.
Assessment: WARNING (Secondary attribution acknowledged)

---

## Source 10

Claimed Title: go-mutesting — Mutation testing for Go source code
Claimed Publisher: zimmski (GitHub repository)
URL: https://github.com/zimmski/go-mutesting
Reachable: YES
Source Type: COMMUNITY (Open-source project)
Relevant: YES
Supports Claimed Topic: YES
Problems:
- Community-maintained tool; smaller ecosystem adoption compared to PIT/Stryker. Accurately characterized as such.
Assessment: PASS

---

## Source 11

Claimed Title: gremlins — A mutation testing tool for Go
Claimed Publisher: go-gremlins (GitHub organization)
URL: https://github.com/go-gremlins/gremlins
Homepage: https://gremlins.dev
Reachable: YES
Source Type: COMMUNITY (Open-source project)
Relevant: YES
Supports Claimed Topic: YES
Problems:
- Emerging tool in active development. Properly framed as less mature than JVM/JS equivalents.
Assessment: PASS

---

## Source 12

Claimed Title: Mutation-Guided LLM-based Test Generation at Meta (arXiv preprint)
Claimed Publisher: arXiv (cs.SE, cs.AI, cs.LG)
URL: https://arxiv.org/abs/2501.12862
Reachable: YES (Verified via WebFetch)
Source Type: PRIMARY (Academic preprint, submitted to FSE 2025 Industry Track)
Relevant: YES
Supports Claimed Topic: YES
Problems:
- Abstract text verified; binary PDF text not extracted. Abstract contains all claimed numeric data points (9,095 mutants, 571 tests, 10,795 classes, 7 platforms, 0.79/0.47 to 0.95/0.96 precision/recall, 73% accepted, 36% privacy-relevant).
Assessment: PASS

---

## Source 13

Claimed Title: Mutation testing (Wikipedia — subsumed mutants section)
Claimed Publisher: Wikipedia
URL: https://en.wikipedia.org/wiki/Mutation_testing
Reachable: YES
Source Type: SECONDARY (Crowdsourced encyclopedia)
Relevant: YES
Supports Claimed Topic: YES
Problems:
- None. Correctly defines subsumed mutants.
Assessment: PASS
