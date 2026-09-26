# Content Audit Report

## Overview

| Field | Value |
|-------|-------|
| Target | `labs/23-optimistic-vs-pessimistic-locking/content/` |
| Audit Date | 2026-09-26 |
| Engineer Audit Status | APPROVED |
| Research Status | APPROVED |
| Tests | PASS (all 6 tests pass, race detector clean) |
| Demo | PASS (all 5 scenarios execute as described) |

---

## File-by-File Findings

### File 1: `01-content-brief.md`

**Status: APPROVED**

| Check | Result | Notes |
|-------|--------|-------|
| Topic accuracy | PASS | Correctly identifies lost update anomaly and three prevention strategies |
| Approved research status | PASS | Stated as APPROVED — matches research audit verdict |
| Approved engineering status | PASS | Stated as APPROVED — matches engineering audit verdict |
| Core mental model | PASS | Pessimistic=prevent, Optimistic=detect, Atomic=eliminate window — aligns with research Finding 5, 6 |
| Verified behaviors (6 items) | PASS | All 6 verified behavior claims match actual test/demo execution: 50 goroutine naive → 99, 50 pessimistic → 50, 20 optimistic direct → 1 success/19 conflict, 20 retry → 20 success/stok 80, 50 atomic → 50, race detector passes |
| Warnings section | PASS | Correctly notes in-memory simulation with sync.Mutex (not real RDBMS), artificial time.Sleep delays, nondeterministic outcomes |
| Available case studies | PASS | Demo CLI and test suite accurately described |

No issues found. The brief accurately frames the problem, target audience, and verified behaviors.

---

### File 2: `02-master-draft.md`

**Status: APPROVED_WITH_WARNINGS**

| Check | Result | Notes |
|-------|--------|-------|
| Problem section (lost update) | PASS | Correctly describes lost update with 100→99 example. However, the example says "stok akhir = 99, padahal harus = 98" which is accurate for the 2-thread example. |
| Why This Matters | PASS | Business impact of lost update correctly identified |
| Mental Model diagram | PASS | ASCII diagram accurately shows two threads reading same value, overwriting |
| Pessimistic locking section | PASS | SQL example with `SELECT ... FOR UPDATE` is correct and standard |
| Optimistic locking section | PASS | Version guard + affected_rows check correctly described |
| Atomic section | PASS | `UPDATE ... SET stock = stock - 3 WHERE stock >= 3` matches research Finding 6 |
| Architecture diagram | PASS | Structure (cmd/demo, tests/locking_test.go, internal/inventory/{model,store,service}) matches actual file layout |
| Implementation section | **WARNING** | Minor simplification: omits `qty <= 0` validation and `atomic.AddInt64` counter increments in code descriptions. Not an error but slightly incomplete. |
| Code walkthrough (lines 122-186) | **WARNING** | The Go code snippets in the walkthrough omit several implementation details present in actual code: `if qty <= 0 { return ErrInvalidQuantity }` guard and `atomic.AddInt64(&s.NaivelyDrawn, ...)` counter increments are stripped. This reduces fidelity but does not change core logic. |
| Test results table (lines 190-197) | PASS | All test outcomes match actual execution results |
| Recovery / Rollback (lines 201-205) | PASS | Deadlock handling and optimistic retry/409 correctly described |
| Production Considerations (lines 207-212) | PASS | Lock-during-HTTP warning, retry requirement, 64-bit version recommendation, isolation level nuances all align with research |
| Common Mistakes (lines 214-220) | PASS | All 5 anti-patterns match research Finding 9 |
| Case Study (lines 222-245) | **WARNING** | The case study section omits the "Total Attempted Conflicts Retried" value (shown as 61 in engineering execution results, 41-44 in live demo runs). The master draft only shows "Successful Deductions: 20" and "Actual Final Stock: 80" for scenario 4, omitting the conflict-retry count metric that the demo actually outputs. This is an omission of demo output detail, not a factual error. |
| Checklist (lines 247-257) | PASS | Comprehensive and accurate |
| Key Takeaways (lines 259-270) | PASS | Aligns with research and engineering findings |
| Sources (lines 272-279) | PASS | Research sources correctly cited |

**Issues found:**
1. **Minor omission** (line 106): Typo "acuisisi" should be "akuisisi" (acquire).
2. **Code simplification** (walkthrough lines 122-186): Snippets omit validation guards and atomic counter lines for brevity. Reduces code fidelity but logic is preserved.
3. **Demo output omission** (Case Study): Scenario 4 ([4]) omits "Total Attempted Conflicts Retried" and "Elapsed Time" fields that are present in the actual demo output.
4. **Explanation inaccuracy** (line 194 of 03-code-snippets.md, referenced in master draft): States "attempt 0 = ~0ms" but actual backoff is `1<<0 = 1ms` minimum, so attempt 0 sleep is ~1-6ms, not ~0ms.

**No factual errors found. All core claims verified against implementation.**

---

### File 3: `03-code-snippets.md`

**Status: APPROVED_WITH_WARNINGS**

| Check | Result | Notes |
|-------|--------|-------|
| Snippet 1 (model.go) | PASS | Error definitions and Product struct match actual model.go exactly |
| Snippet 2 (Store core) | PASS | Store struct fields and Seed() method match actual store.go. Snippet adds `Get(id)` method reference not in actual code but matches store.go line 53-61. |
| Snippet 3 (NaiveDeduct) | **WARNING** | Code snippet omits `if qty <= 0` validation at top and `atomic.AddInt64(&s.NaivelyDrawn, ...)` at the end. Core logic (read, sleep, stale write) is correct. |
| Snippet 4 (PessimisticDeduct) | **WARNING** | Code snippet omits `if qty <= 0` validation and `atomic.AddInt64(&s.Pessimistically, ...)` line. Core logic is correct. |
| Snippet 5 (OptimisticDeduct) | **WARNING** | Code snippet omits `if qty <= 0` validation and `atomic.AddInt64(&s.Optimistically, ...)` line. Core logic is correct. |
| Snippet 6 (DeductOptimisticWithRetry) | PASS | Matches actual service.go lines 28-45 exactly |
| Snippet 7 (AtomicDeduct) | **WARNING** | Code snippet omits `if qty <= 0` validation and `atomic.AddInt64(&s.Atomically, ...)` line. Core logic is correct. |
| Snippet 8 (TestNaiveLostUpdate) | PASS | Test code matches actual locking_test.go lines 10-38 (adapted as standalone snippet) |
| Snippet 9 (TestPessimisticLocking) | PASS | Test code matches actual locking_test.go lines 40-67 |
| Snippet 10 (TestOptimisticLockingConflict) | PASS | Test code matches actual locking_test.go lines 85-121 |
| Snippet 11 (Demo CLI) | **WARNING** | Shows truncated demo with `// ... (Optimistic, Optimistic+Retry, Atomic sama polanya)` comment. Does not show full scenarios 3-5. This is a reasonable simplification for a snippet but reduces completeness. |
| Explanation text for Snippet 6 | **ISSUE** | States "attempt 0 = ~0ms, 1 = ~2ms, 2 = ~4ms" — this is INCORRECT. The actual code `time.Duration(1<<attempt)*time.Millisecond` means: attempt 0 = 1ms base, attempt 1 = 2ms base, attempt 2 = 4ms base. Plus 0-5ms jitter. The explanation should say "attempt 0 = ~1-6ms, 1 = ~2-7ms, 2 = ~4-9ms". The "~0ms" claim is factually wrong. |

**Issues found:**
1. **Explanation inaccuracy** (Snippet 6 explanation, line 194): Backoff timing description for attempt 0 is incorrect (~0ms stated vs ~1-6ms actual).
2. **Code simplification**: 5 of 7 implementation snippets omit validation guards and atomic counter increments for brevity, reducing code fidelity.
3. **Demo snippet truncation**: Scenario 3-5 shown as abbreviated comments.

**The backoff timing inaccuracy is a minor factual error that should be corrected.**

---

### File 4: `04-diagrams.md`

**Status: APPROVED**

| Check | Result | Notes |
|-------|--------|-------|
| Diagram 1 (Lost Update Flow) | PASS | Accurately shows two threads reading same value, both writing, result is overwritten |
| Diagram 2 (Pessimistic Locking) | PASS | Shows lock acquisition, blocking, sequential execution — matches PessimisticDeduct implementation |
| Diagram 3 (Optimistic with Conflict) | PASS | Version guard + affected_rows=0 + retry/409 — matches OptimisticDeduct implementation and research Finding 4 |
| Diagram 4 (Atomic Update Flow) | PASS | Single lock-check-update-unlock flow matches AtomicDeduct implementation |
| Diagram 5 (Retry Convergence) | PASS | Shows conflict detection, backoff, retry, version increment, convergence — matches actual behavior |
| Diagram 6 (Store Architecture) | PASS | Shows mu (global), rowLocks, products — matches actual Store struct fields |
| Diagram 7 (Version Guard SQL) | PASS | SQL pattern matches research Finding 4 and content brief description |
| Diagram 8 (Service Retry Loop) | PASS | Retry loop flowchart with exponential backoff and jitter — matches actual implementation |
| Diagram 9 (Error Flow) | PASS | Error paths for naive and optimistic correctly shown |
| Diagram 10 (Counter Invariants) | PASS | `InitialStock - TotalDrawn == FinalStock` invariant matches research and test assertions |

**No issues found. All diagrams are accurate representations of the implementation.**

---

### File 5: `05-key-takeaways.md`

**Status: APPROVED**

| Check | Result | Notes |
|-------|--------|-------|
| Takeaway 1 (Lost update = silent corruption) | PASS | Aligns with research Finding 1 |
| Takeaway 2 (Pessimistic = block to prevent) | PASS | Aligns with research Finding 5 |
| Takeaway 3 (Optimistic = detect via version + affected_rows) | PASS | Aligns with research Finding 4 |
| Takeaway 4 (Atomic = statement-level) | PASS | Aligns with research Finding 6 |
| Takeaway 5 (Isolation level insufficient) | PASS | Aligns with research Finding 7 |
| Takeaway 6 (Selection criteria) | PASS | Aligns with research Finding 5, 9 |
| Takeaway 7 (Deadlock handling) | PASS | Aligns with research Finding 3 |
| Takeaway 8 (Retry requirement) | PASS | Aligns with research Finding 4 |
| Takeaway 9 (Race detector) | PASS | Verified — race detector passes |
| Takeaway 10 (Invariant verification) | PASS | Aligns with test assertions and execution results |

**No issues found. All 10 key takeaways are accurate and verified.**

---

### File 6: `06-source-map.md`

**Status: APPROVED**

| Check | Result | Notes |
|-------|--------|-------|
| Structure references | PASS | All source file references exist and match actual project structure |
| Line number accuracy | PASS | All line number references verified against actual source files:
  - `NaiveDeduct()` store.go lines 65-89 — **matches**
  - `TestNaiveLostUpdate()` tests/locking_test.go lines 10-38 — **matches**
  - `PessimisticDeduct()` store.go lines 93-116 — **matches**
  - `GetRowLock()` store.go lines 42-51 — **matches**
  - `TestPessimisticLocking()` lines 40-67 — **matches**
  - `TestPessimisticLockingInsufficientStock()` lines 69-83 — **matches**
  - `OptimisticDeduct()` store.go lines 120-152 — **matches**
  - `DeductOptimisticDirect()` service.go lines 24-26 — **matches**
  - `DeductOptimisticWithRetry()` service.go lines 28-45 — **matches**
  - `ErrOptimisticLock` model.go line 8 — **matches**
  - `TestOptimisticLockingConflict()` lines 85-121 — **matches**
  - `TestOptimisticLockingWithRetry()` lines 123-145 — **matches**
  - `AtomicDeduct()` store.go lines 155-173 — **matches**
  - `DeductAtomic()` service.go lines 47-49 — **matches**
  - `TestAtomicConditionalUpdate()` lines 147-170 (actual: 147-171) — **near match**, off by 1 line |
| Research status verification | PASS | References `research-audit/07-verdict.md — APPROVED` and `engineering-audit/06-verdict.md — APPROVED` |
| Execution verification | PASS | References `engineering/03-execution-result.md` and other engineering docs |
| Demo line references | PASS | cmd/demo/main.go scenario line references (16-34, 36-53, 55-74, 76-97, 99-116) all verified against actual demo code |

**Minor issue found:**
- `TestAtomicConditionalUpdate()` line range stated as 147-170, actual is 147-171 (off by one line at end).

**No factual errors. All references resolve to existing files and methods.**

---

## Summary of Findings

### Factual Errors (0 blocking)
1. **Snippet 6 explanation** (`03-code-snippets.md` line 194): States "attempt 0 = ~0ms" for exponential backoff, but the actual code `1<<attempt` yields 1ms minimum for attempt 0. Correct explanation should state "~1-6ms" (1ms base + 0-5ms jitter).

### Minor Inaccuracies/Omissions
1. **Walkthrough code snippets** (`02-master-draft.md` lines 122-186, `03-code-snippets.md`): 5 of 7 implementation snippets omit `if qty <= 0` validation guard and `atomic.AddInt64` counter increments for brevity. Core logic preserved; fidelity reduced.
2. **Case Study section** (`02-master-draft.md` lines 226-245): Scenario 4 omits "Total Attempted Conflicts Retried" and "Elapsed Time" output fields present in actual demo.
3. **Demo snippet** (`03-code-snippets.md` lines 387): Scenario 3-5 shown as abbreviated comment instead of full code.
4. **Typo** (`02-master-draft.md` line 106): "acuisisi" should be "akuisisi".
5. **Line range** (`06-source-map.md`): `TestAtomicConditionalUpdate()` stated as lines 147-170, actual is 147-171.

### Positive Findings
- All source map line references verified accurate (15/16 exact matches, 1 off-by-one).
- All test results match actual execution output.
- All demo output values match actual demo run.
- All 3 core strategies correctly described with matching SQL patterns.
- All anti-patterns, common mistakes, and production considerations align with research findings.
- All 10 key takeaways verified against engineering implementation.
- Race detector claim (no warnings) verified — `go test -race ./...` passes.
- Simulation caveat properly documented in content brief.

## Conclusion

The content is **factually accurate and reflects the actual engineering implementation**. All core concepts, test results, demo outputs, source file references, and line numbers are verified against the codebase.

The single factual inaccuracy (backoff timing description for attempt 0 being ~0ms instead of ~1-6ms) is minor and does not affect the core educational message. All code omissions in walkthrough snippets are stylistic simplifications that preserve the essential logic. All source map references resolve correctly.

**No hallucinated facts, no platform-specific biases, no misleading claims.**

---

## Final Verdict

**APPROVED_WITH_WARNINGS**
