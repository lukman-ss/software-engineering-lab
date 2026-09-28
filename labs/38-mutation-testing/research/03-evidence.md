# Evidence

## Evidence 1

Claim: Mutation testing was first proposed by Richard Lipton in 1971 and published by DeMillo, Lipton, and Sayward in 1978.
Evidence: "Mutation testing was originally proposed by Richard Lipton as a student in 1971, and first developed and published by DeMillo, Lipton and Sayward." — Wikipedia (citing [1]).
Source: Wikipedia (referencing DeMillo, Lipton, Sayward 1978)
URL: https://en.wikipedia.org/wiki/Mutation_testing
Confidence: HIGH
Corroborated By: Martin Fowler's bliki entry also references the historical timeline of mutation testing.
Notes: Original paper NOT directly accessed. Bibliographic reference from Wikipedia: R. A. DeMillo, R. J. Lipton, F. G. Sayward. Hints on test data selection: Help for the practicing programmer. IEEE Computer, 11(4):34-41, April 1978.

## Evidence 2

Claim: The two core hypotheses of mutation testing are the "competent programmer hypothesis" and the "coupling effect hypothesis".
Evidence: "Mutation testing is based on two hypotheses. The first is the competent programmer hypothesis. This hypothesis states that competent programmers write programs that are close to being correct. The second hypothesis is called the coupling effect. The coupling effect asserts that simple faults can cascade or couple to form other emergent faults." — Wikipedia (citing [6] and [12], [13]).
Source: Wikipedia
URL: https://en.wikipedia.org/wiki/Mutation_testing
Confidence: HIGH
Corroborated By: Not independently verified from original papers, but consistent with description across tool documentation.
Notes: Source [6] is Ammann & Offutt (2008). [12] is Offutt (1992). [13] is Acree et al. (1979).

## Evidence 3

Claim: The mutation score measures test quality by the percentage of mutants killed.
Evidence: "The value of a test suite is measured by the percentage of mutants that it kills." — Wikipedia (citing [6]). Also: "Each run makes a small modification to the code, such as reversing a conditional or removing a line. We then run the test suite. If the tests pass, then we've found a problem." — Martin Fowler.
Source: Wikipedia + Martin Fowler
URL: https://en.wikipedia.org/wiki/Mutation_testing / https://martinfowler.com/bliki/MutationTesting.html
Confidence: HIGH
Corroborated By: PIT FAQ: "The quality of your tests can be gauged from the percentage of mutations killed."
Notes: The formula Mutation Score = (Killed / Total) × 100% is standard convention; Wikipedia states it as a ratio.

## Evidence 4

Claim: Code coverage can be 100% while test quality is poor; mutation testing exposes this discrepancy.
Evidence: "Traditional test coverage (i.e line, statement, branch, etc.) measures only which code is executed by your tests. It does not check that your tests are actually able to detect faults in the executed code." — PIT. Stryker analogy: "Imagine a sandwich covered with paste. Code coverage would tell you the bread is 80% covered with paste. Mutation testing, on the other hand, would tell you it is actually chocolate paste."
Source: PIT (https://pitest.org/) and Stryker (https://stryker-mutator.io/docs/General/example/)
Confidence: HIGH
Corroborated By: Consistent with the central thesis of both PIT and Stryker documentation.
Notes: This directly validates the lab's key claim about false confidence from 100% code coverage.

## Evidence 5

Claim: Mutation operators include statement deletion, arithmetic operator replacement, relational operator replacement, Boolean replacement, and statement insertion.
Evidence: "Many mutation operators have been explored by researchers. Here are some examples... Statement deletion; Statement duplication or insertion, e.g. goto fail; Replacement of Boolean subexpressions with true and false; Replacement of some arithmetic operations with others, e.g. + with *, - with /; Replacement of some Boolean relations with others, e.g. > with >=, == and <=" — Wikipedia (citing Hamimoune & Falah, 2016).
Source: Wikipedia (referencing Hamimoune & Falah, 2016)
URL: https://en.wikipedia.org/wiki/Mutation_testing
Confidence: HIGH
Corroborated By: PIT FAQ lists mutator groups DEFAULTS, STRONGER, ALL. Stryker example output shows "Mutator: BinaryOperator" and "Mutator: RemoveConditionals".
Notes: These are called "traditional mutation operators". MuJava also defines class-level and method-level operators.

## Evidence 6

Claim: Equivalent mutants are a major obstacle in mutation testing because determining equivalence is mathematically undecidable.
Evidence: "Equivalent mutants detection is one of the biggest obstacles to practical usage of mutation testing. The effort needed to check if mutants are equivalent or not can be very high, even for small programs." — Wikipedia (citing [18], [19]). Meta's ACH paper: "Determining whether a mutant is equivalent or not is known to be mathematically undecidable."
Source: Wikipedia + Meta Engineering Blog (https://engineering.fb.com/2025/09/30/security/llms-are-the-key-to-mutation-testing-and-better-compliance/)
Confidence: HIGH
Corroborated By: PIT FAQ confirms they filter out certain code patterns.
Notes: A 2014 systematic literature review (Madeyski et al.) identified 17 techniques across 22 articles to address the Equivalent Mutant Problem.

## Evidence 7

Claim: Strong mutation testing requires all three RIP conditions: Reach, Infect, and Propagate.
Evidence: "For the test to kill this mutant, the following three conditions should be met: 1. A test must reach the mutated statement. 2. Test input data should infect the program state by causing different program states for the mutant and the original program. 3. The incorrect program state must propagate to the program's output and be checked by the test. These conditions are collectively called the RIP model." — Wikipedia (citing [9]).
Source: Wikipedia (referencing Offutt & Untch, 2000 — Mutation 2000)
URL: https://en.wikipedia.org/wiki/Mutation_testing
Confidence: HIGH
Corroborated By: Foundational in the mutation testing literature; [9] links to cs.gmu.edu/~offutt/rsrch/papers/mut00.pdf.
Notes: Weak mutation testing only requires Reach + Infect (first two conditions). PIT implements strong mutation.

## Evidence 8

Claim: All major mutation testing tools follow the same pattern: generate mutants → run tests → classify as killed/survived.
Evidence: "Mutation testing involves making small changes to the program being tested. Each changed version is called a mutant. A test detects, and therefore rejects, a mutant upon test failure –– failure indicating that the test successfully discerned that the behaviour of the mutant differs from the behaviour of the original code. Rejection is called killing the mutant." — Wikipedia (citing DeMillo et al., 1978). PIT states: "PIT runs your unit tests against automatically modified versions of your application code. When the application code changes, it should produce different results and cause the unit tests to fail."
Source: Wikipedia + PIT
URL: https://en.wikipedia.org/wiki/Mutation_testing / https://pitest.org/
Confidence: HIGH
Corroborated By: Stryker docs: "Bugs, or mutants, are automatically inserted into your production code. Your tests are run for each mutant. If your tests fail then the mutant is killed. If your tests passed, the mutant survived."
Notes: PIT mutates bytecode in memory (never writes to disk). Stryker may work differently at the AST level.

## Evidence 9

Claim: Mutation testing faces five major barriers to industrial adoption: scalability, unrealistic mutants, equivalent mutants, computational cost, and overstretching testing efforts.
Evidence: "Traditional mutation testing generates a very large number of mutants, making it computationally expensive and difficult to scale... Even though mutation testing cannot exist on its own (it requires a test to already exist), it helps engineers and developers identify weak assertions and encourages them to write tests that truly validate code behavior instead of just executing it." — Meta Engineering Blog (by Mark Harman, September 2025).
Source: Meta Engineering Blog
URL: https://engineering.fb.com/2025/09/30/security/llms-are-the-key-to-mutation-testing-and-better-compliance/
Confidence: HIGH
Corroborated By: PIT FAQ: "Mutation testing is a computationally expensive process and can take quite some time depending on the size of your codebase and the quality and speed of your test suite."
Notes: The five barriers are explicitly listed in the Meta article. The claim that "mutation testing cannot exist on its own" is consistent with PIT's description that it evaluates existing tests.

## Evidence 10

Claim: Mutation testing is classified as white-box testing and its purpose is to help develop effective regression tests.
Evidence: "Mutation testing is a form of white-box testing. Its purpose is to help the tester develop effective regression tests by locating weaknesses in the test data used to test the program and discovering sections of the tested program's code that are seldom or never accessed during execution." — Wikipedia (citing [3], [4]).
Source: Wikipedia (citing Ostrand 2002; Misra 2003)
URL: https://en.wikipedia.org/wiki/Mutation_testing
Confidence: HIGH
Corroborated By: Standard textbook definition (Ammann & Offutt, 2008).
Notes: Unlike black-box testing, mutation testing requires source code access.

## Evidence 11

Claim: Mutation testing includes three types based on what aspect is modified: statement mutation, value mutation, and decision mutation.
Evidence: "There are three types of mutation testing; Statement mutation... Value mutation... Decision mutation..." — Wikipedia, describing each category with code examples.
Source: Wikipedia
URL: https://en.wikipedia.org/wiki/Mutation_testing
Confidence: MEDIUM (synthesized from multiple academic citations in Wikipedia)
Corroborated By: PIT mutator groups map to these categories: METHOD_DELETION → statement, MATH/NAMESPACE → value, CONDITIONALS_BOUNDARY → decision.
Notes: This taxonomy is pedagogically useful; industrial tools use different categorizations.

## Evidence 12

Claim: The Meta ACH system uses LLMs to generate targeted mutants and tests, achieving 73% engineer acceptance and 36% privacy relevance in trials.
Evidence: "By leveraging LLMs we've been able to overcome the barriers that have prevented mutation testing from being efficiently deployed at scale." — Mark Harman, Meta Engineering. "Over thousands of mutants and hundreds of generated tests, privacy engineers at Meta accepted 73% of the generated tests, with 36% judged as privacy relevant."
Source: Meta Engineering Blog
URL: https://engineering.fb.com/2025/09/30/security/llms-are-the-key-to-mutation-testing-and-better-compliance/
Confidence: MEDIUM (industry report, single source, not yet peer-reviewed)
Corroborated By: arXiv preprint https://arxiv.org/pdf/2501.12862 referenced in the blog (PDF was fetched but content was binary; abstract text cited in blog).
Notes: ACH combines mutation-guided test generation with LLMs. LLM equivalence detector: 0.79/0.47 precision/recall → 0.95/0.96 with preprocessing. This is the most recent development in the field (September 2025).
