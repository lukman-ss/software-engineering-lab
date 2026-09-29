# Contradictions Analysis

Target Lab: labs/38-mutation-testing
Audit Scope: Research Files Only (PIPELINE OVERRIDE)

---

## Analysis of Document Consistency

### 1. Research vs. Sources
- `research/02-sources.md` defines 13 sources.
- `research/03-evidence.md` cross-references these sources consistently (Evidence 1 through 14 match Sources 1 through 13).
- `research/05-report.md` synthesizes all 14 evidence statements across 11 clear findings.
- No discrepancy between source claims and reported findings.

### 2. Theoretical Undecidability vs. Empirical Tool Claims
- Theoretical fact: Proving equivalence of arbitrary programs is mathematically undecidable (reduction to the Halting Problem).
- Practical tool behavior: Tools use heuristics (e.g., PIT avoiding enum constructors, Meta ACH using LLM classifiers with static analysis preprocessing) to detect equivalence with practical precision (up to 0.95).
- Audit Assessment: This represents the classic distinction between formal computability limits and engineering heuristics. The research explicitly highlights this nuance in `research/04-contradictions.md` (Point 2) and `research/05-report.md` (Finding 6). No contradiction exists.

### 3. "Mutation Testing Requires Existing Tests" vs. "LLM-Guided Test Generation"
- Traditional concept: Mutation testing evaluates an existing test suite by mutating code and measuring test failure.
- Modern expansion (Meta ACH): Mutation testing identifies uncaught mutants, which then serve as targeted prompts for LLMs to generate new tests.
- Audit Assessment: This is an evolutionary expansion rather than a logical conflict. The research analyzes this distinction in `research/04-contradictions.md` (Point 1). No contradiction exists.

### 4. Computational Cost Claims
- Earlier historical context: Full mutation testing was computationally prohibitive (recompiling and rerunning all tests for every mutant).
- Modern tools: Bytecode manipulation in-memory (PIT), AST-based mutating (Stryker), and incremental testing on changed code only allow execution in minutes.
- Audit Assessment: The research accurately describes how modern tools mitigate computational cost without claiming the problem is eliminated. No contradiction exists.

### 5. Status of Martin Fowler's Bliki
- Note: Martin Fowler's bliki carries an explicit "This is a draft entry" banner.
- Audit Assessment: The researcher resolved earlier audit feedback by consistently annotating all in-text citations as `Martin Fowler (pre-publication draft; carries 'This is a draft entry' notice)`. This eliminates any risk of treating pre-publication commentary as definitive finalized doctrine.

---

## Verdict

No material contradictions found.
