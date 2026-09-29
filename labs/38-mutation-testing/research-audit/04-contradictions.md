# Contradictions Analysis Audit

## Summary

The research documents three potential tension points and concludes that no material contradictions exist. This audit confirms that the three points analyzed are genuine tensions in the field and that the Research Agent's reconciliations are logically sound.

No new material contradictions were discovered during the audit.

---

## Evaluation of Documented Tension Points

### Point 1: "Mutation testing requires a test to already exist" vs. "LLMs can generate tests from mutants"

Statements:
- Statement A: "Even though mutation testing cannot exist on its own (it requires a test to already exist)..." (Meta Engineering Blog)
- Statement B: ACH uses unkilled mutants to generate new tests via LLMs (Meta Engineering Blog / arXiv abstract)

Auditor Assessment: SOUND

Analysis: The Research Agent correctly distinguishes between traditional mutation testing (which requires tests to evaluate) and mutation-guided test generation (which uses mutation gaps to prompt test creation). ACH still requires initial tests to establish an oracle. The reconciliation correctly explains this as workflow extension, not contradiction.

---

### Point 2: Equivalent mutants as "mathematically undecidable" vs. tools claiming equivalence detection

Statements:
- Statement A: "Determining whether a mutant is equivalent or not is known to be mathematically undecidable." (Meta Engineering Blog / academic literature)
- Statement B: Meta ACH reports 0.95 precision / 0.96 recall for equivalent mutant detection with static preprocessing.

Auditor Assessment: SOUND

Analysis: Undecidability is a theoretical limit for arbitrary programs (Rice's Theorem / Halting Problem analog). Heuristic approximations (ML classifiers, static analysis) can achieve high empirical precision on typical human-written code without solving the general theoretical problem. The distinction between theoretical undecidability and practical heuristic efficacy is standard in computer science and is accurately represented.

---

### Point 3: Computational cost concern vs. tool claims of speed

Statements:
- Statement A: "Mutation testing is a computationally expensive process and can take quite some time..." (PIT FAQ)
- Statement B: Stryker claims "fast to run and easy to use."

Auditor Assessment: SOUND

Analysis: The speed difference reflects relative improvement over first-generation tools (1970s–1990s) that required full recompilation per mutant. Modern tools optimize execution via bytecode manipulation, incremental analysis, and test selection. However, running mutation testing on large codebases remains slow, which both tools acknowledge (e.g., PIT recommends targeting changed files only; Gremlins warns about runs taking hours). "Fast" is relative. The reconciliation is accurate.

---

## Audit of Undocumented Contradictions

### Potential Contradiction 1: Fowler draft bliki citation vs. availability

Issue: Research cites Martin Fowler's bliki entry (`https://martinfowler.com/bliki/MutationTesting.html`) as corroborating evidence in multiple places, but the URL returns a 404 HTTP status during audit.

Assessment:
- Type: SOURCE_AVAILABILITY_MISMATCH
- Impact: Martin Fowler is a prominent figure whose citation lends authority. The URL is unreachable, and the research correctly noted the page carried a "This is a draft entry" notice. If the page was removed or never made permanently public, using it as an authoritative Tier 1 source is problematic.
- Mitigating Factor: The core claims (mutation score definition, false confidence of code coverage, mutation operators) are independently supported by PIT, Stryker, and Wikipedia. The Fowler citation was used only as corroborating evidence, not as sole support.
- Action: Documented in `02-source-audit.md` (FAIL) and `06-gaps.md` (WEAK_SOURCE). Research should either replace this source or explicitly note that the draft page was removed.

---

## Conclusion

No internal contradictions between research documents. No source-to-source factual conflicts that invalidate findings. Tension points documented by the Research Agent were analyzed accurately and resolved with appropriate nuance.
