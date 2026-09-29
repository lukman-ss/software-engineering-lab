# Source Audit

Target Lab: `/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/38-mutation-testing`
Audit Date: 2026-09-29

---

## Source 1
Claimed Title: Hints on test data selection: Help for the practicing programmer
Claimed Publisher: IEEE Computer (DeMillo, Lipton, Sayward)
URL: NOT VERIFIED by research agent (DOI not opened)
Reachable: NO (URL not provided by research agent; bibliographic only via Wikipedia)
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES (founding mutation testing paper)
Problems: Research agent correctly disclosed this as "NOT VERIFIED for direct quotations." All claims attributed to this paper are sourced via Wikipedia references. No fabrication detected. The research file is transparent about secondary attribution.
Assessment: WARNING (secondary-attributed; acceptable with disclosure as documented)

---

## Source 2
Claimed Title: Mutation testing
Claimed Publisher: Wikipedia
URL: https://en.wikipedia.org/wiki/Mutation_testing
Reachable: YES
Source Type: SECONDARY
Relevant: YES
Supports Claimed Topic: YES
Problems: Page last edited 6 September 2026 — confirmed by live fetch. Contents verified:
- History section confirms: "Mutation testing was originally proposed by Richard Lipton as a student in 1971" and "first developed and published by DeMillo, Lipton and Sayward."
- RIP model section present with exactly the 3 conditions cited in research (Reach, Infect, Propagate).
- Competent programmer hypothesis and coupling effect confirmed.
- Mutation operators section confirmed (statement deletion, arithmetic replacement, relational replacement, boolean replacement).
- Subsumed mutants section confirmed (definition and example match research).
- Statement/value/decision mutation classification confirmed.
- Equivalent mutant problem confirmed as "one of the biggest obstacles."
Assessment: PASS

---

## Source 3 (REMOVED/UNREACHABLE)
Claimed Title: Mutation Testing (bliki entry)
Claimed Publisher: Martin Fowler
URL: https://martinfowler.com/bliki/MutationTesting.html
Reachable: NO (HTTP 404 confirmed by direct fetch)
Source Type: SECONDARY
Relevant: PARTIAL
Supports Claimed Topic: N/A — unreachable
Problems: URL returns HTTP 404. Research agent correctly marked this as "REMOVED (UNREACHABLE)" in `research/02-sources.md`. However, residual citation appears in `research/05-report.md` Finding 2's Sources field: "Martin Fowler (pre-publication draft; https://martinfowler.com/bliki/MutationTesting.html)". The research notes this claim rests fully on Wikipedia and PIT. Research-revision/02-changes-made.md documents this as verified non-blocking since claim support is multi-source. Finding 2 claim is independently supported by Wikipedia and PIT without Fowler.
Assessment: WARNING (dead URL; correctly disclosed and decoupled from all claim support; residual citation in report is labeled as pre-publication and unreachable)

---

## Source 4
Claimed Title: PIT Mutation Testing (homepage)
Claimed Publisher: PIT Project
URL: https://pitest.org/
Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: Live fetch confirmed:
- "PIT is a state of the art mutation testing system, providing gold standard test coverage for Java and the jvm."
- Confirmed: "Traditional test coverage... measures only which code is executed by your tests. It does not check that your tests are actually able to detect faults in the executed code."
- Workflow description confirmed: "Faults (or mutations) are automatically seeded into your code, then your tests are run. If your tests fail then the mutation is killed, if your tests pass then the mutation lived. The quality of your tests can be gauged from the percentage of mutations killed."
- Confirmed: Java and Kotlin (via Arcmutate) are supported languages; no Go.
Assessment: PASS

---

## Source 5
Claimed Title: PIT FAQ
Claimed Publisher: PIT Project
URL: https://pitest.org/faq/
Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: Live fetch confirmed:
- "Mutation testing is a computationally expensive process and can take quite some time depending on the size of your codebase..."
- Mutator groups DEFAULTS, STRONGER, ALL confirmed.
- "Currently supported languages are Java, Kotlin (via the Arcmutate kotlin plugin)" — no Go support confirmed.
- Determinism claims confirmed.
- "The most effective way to use mutation testing is usually to limit analysis to the code that you are changing."
- Note: Fetch confirmed PIT mutates bytecode in memory and never writes mutations to disk.
Assessment: PASS

---

## Source 6
Claimed Title: What is mutation testing? (Introduction)
Claimed Publisher: Stryker Mutator
URL: https://stryker-mutator.io/docs/
Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: Live fetch confirmed:
- Definition of killed/survived mutants confirmed: "Bugs, or mutants, are automatically inserted into your production code. Your tests are run for each mutant. If your tests fail then the mutant is killed. If your tests passed, the mutant survived."
- Sandwich/paste analogy confirmed verbatim: "Imagine a sandwich covered with paste. Code coverage would tell you the bread is 80% covered with paste. Mutation testing, on the other hand, would tell you it is actually chocolate paste."
- `user.age >= 18` worked example with 4 mutants confirmed exactly as cited.
- Mutator names BinaryOperator and RemoveConditionals confirmed in example output.
Assessment: PASS

---

## Source 7
Claimed Title: Welcome to the RoboCoasters — An introduction to mutation testing
Claimed Publisher: Stryker Mutator
URL: https://stryker-mutator.io/docs/General/example/
Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: Live fetch confirmed:
- Page titled "Welcome to the RoboCoasters 🤖🎢"
- Tagline confirmed: "How code coverage of 100% could mean only 60% is tested."
- The demonstration project (RoboCoasters) with 100% code coverage and ~60% mutation score confirmed.
- Instructions include: "Review the 100% code coverage score" → "Review ~60% mutation score."
Assessment: PASS

---

## Source 8
Claimed Title: LLMs Are the Key to Mutation Testing and Better Compliance
Claimed Publisher: Engineering at Meta (by Mark Harman)
URL: https://engineering.fb.com/2025/09/30/security/llms-are-the-key-to-mutation-testing-and-better-compliance/
Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: Live fetch confirmed:
- Title and author (Mark Harman) confirmed.
- Published September 30, 2025 confirmed.
- Five barriers confirmed with exact names matching research.
- Equivalence detection quote confirmed: "rising to 0.95 and 0.96 with simple preprocessing."
- ACH trial statistics paragraph confirmed: "privacy engineers at Meta accepted 73% of the generated tests, with 36% judged as privacy relevant."
- "Determining whether a mutant is equivalent or not is known to be mathematically undecidable" confirmed.
- Note: Article cites https://arxiv.org/pdf/2501.12862 (ACH preprint) internally.
Assessment: PASS

---

## Source 9
Claimed Title: An Analysis and Survey of the Development of Mutation Testing (bibliographic record)
Claimed Publisher: IEEE Transactions on Software Engineering (Jia & Harman, 2009)
URL: doi:10.1109/TSE.2010.62 (not directly accessible)
Reachable: NOT VERIFIED (DOI string only; not fetched)
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: Research agent correctly disclosed this as "Paper NOT opened directly; cite as secondary-attributed only." Wikipedia's reference section links to an archived PDF at Semantic Scholar confirming the paper is real and authored by Jia & Harman. Mark Harman (co-author of Source 9) is also the author of Source 8 (Meta blog post), reinforcing credibility. No fabrication detected. Claim is properly disclosed as secondary-attributed.
Assessment: WARNING (secondary-attributed; acceptable with disclosure as documented)

---

## Source 10
Claimed Title: go-mutesting — Mutation testing for Go source code
Claimed Publisher: zimmski (GitHub repository)
URL: https://github.com/zimmski/go-mutesting
Reachable: YES
Source Type: COMMUNITY
Relevant: YES
Supports Claimed Topic: YES
Problems: Live fetch confirmed:
- Repository exists and is active (677 stars, 59 forks as of fetch).
- Description confirmed: "go-mutesting is a framework for performing mutation testing on Go source code."
- Mutators confirmed: branch (if, else, case), expression (comparison, remove), statement (remove) — matching research's "branch, expression, and statement levels."
- Exec-based mutation workflow confirmed.
- Mutation score calculation confirmed: "The mutation score is 0.750000 (6 passed, 2 failed, 0 skipped, total is 8)."
- Minor discrepancy noted: research calls this an "AST-level mutation" tool but the README describes an exec-based file-replacement workflow (replacing original source, running tests). go-mutesting does parse Go AST, but the execution mechanism modifies files on disk — the research's description of "exec-based workflow" is accurate; "AST-level mutation" at generation phase is also accurate.
Assessment: PASS

---

## Source 11
Claimed Title: gremlins — A mutation testing tool for Go
Claimed Publisher: go-gremlins (GitHub organization)
URL: https://github.com/go-gremlins/gremlins
Reachable: YES
Source Type: COMMUNITY
Relevant: YES
Supports Claimed Topic: YES
Problems: Live fetch confirmed:
- Repository exists (433 stars, 48 forks as of fetch).
- Tagline confirmed: "A mutation testing tool for Go."
- Tool targets "smallish Go modules, for example microservices" — relevant limitation not mentioned in research.
- Gremlins is in 0.x.x releases with explicit warning: "Gremlins is still in its 0.x.x release, which, as per SemVer, doesn't guarantee backward compatibility." This supports research's claim about inferior maturity vs PIT/Stryker.
- Documentation at https://gremlins.dev confirmed (referenced in README).
- Research's claim that gremlins "doesn't work very well on very big Go modules" not mentioned in research; this was found in the official README. This is an additional nuance the research omitted.
Assessment: PASS (with minor omission: research did not note gremlins is explicitly 0.x.x pre-stable and designed for small modules only)

---

## Source 12
Claimed Title: Mutation-Guided LLM-based Test Generation at Meta (arXiv preprint)
Claimed Publisher: arXiv (cs.SE, cs.AI, cs.LG)
URL: https://arxiv.org/abs/2501.12862
Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: Live fetch confirmed abstract text verbatim:
- "In total, ACH was applied to 10,795 Android Kotlin classes in 7 software platforms deployed by Meta, from which it generated 9,095 mutants and 571 privacy-hardening test cases."
- "ACH also deploys an LLM-based equivalent mutant detection agent that achieves a precision of 0.79 and a recall of 0.47 (rising to 0.95 and 0.96 with simple pre-processing)."
- "ACH was used by Messenger and WhatsApp test-a-thons where engineers accepted 73% of its tests, judging 36% to privacy relevant."
- Authors confirmed: Christopher Foster, Abhishek Gulati, Mark Harman, et al.
- "Submitted to FSE 2025 Industry Track" confirmed (submitted 22 Jan 2025).
- All statistics cited in research exactly match the abstract text.
Assessment: PASS

---

## Source 13
Claimed Title: Mutation testing (Wikipedia — subsumed mutants section)
Claimed Publisher: Wikipedia
URL: https://en.wikipedia.org/wiki/Mutation_testing
Reachable: YES (same as Source 2)
Source Type: SECONDARY
Relevant: YES
Supports Claimed Topic: YES
Problems: Subsumed mutants section confirmed during Source 2 verification. Quote confirmed: "In addition to equivalent mutants, there are subsumed mutants which are mutants that exist in the same source code location as another mutant, and are said to be 'subsumed' by the other mutant. Subsumed mutants are not visible to a mutation testing tool, and do not contribute to coverage metrics."
Assessment: PASS
