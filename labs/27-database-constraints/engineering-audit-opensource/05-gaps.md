# Gap Analysis

Gap types permitted: MISSING_TEST, BROKEN_IMPLEMENTATION, DOC_CODE_MISMATCH, RACE_CONDITION, UNHANDLED_ERROR, MISSING_EDGE_CASE, IMPLEMENTATION_OVERCLAIM, RESEARCH_MISMATCH, FAKE_DEMO, FAKE_BENCHMARK, UNVERIFIED_RESULT.

Only actual gaps observed (no fabrication).

## GAP 1 — MISSING_TEST: MapToDomainError output never asserted

Location: No test calls `dberr.MapToDomainError(err)` and compares the returned string to an expected message.

Impact: While the constructor tests (`TestErrorClassification`) prove that each `New*` yields a `*ConstraintError` with the correct `SQLState`, the translation layer used by `SafeStore` is not verified by an automated test. A regression in the `switch` (e.g., typo in message, missing case) would not be caught by existing tests.

Severity: MEDIUM. Reason: The mapper is exercised by the demo and by `SafeStore` calls (which log or return the error), but an automated assertion is missing. Could lead to silent drift in user-facing messages.

Suggested fix: Add a table-driven test that feeds each `SQLState`/`Detail` into `MapToDomainError` and asserts the string contains the expected phrase ("mandatory field is missing", "outside permissible boundary", etc.).

## GAP 2 — DOC_CODE_MISMATCH: engineering/01-design.md package tree does not match reality

Location: `labs/27-database-constraints/engineering/01-design.md`, Section 31 Architecture and 39 Components.

Claims code is laid out as:
- `internal/db`
- `internal/domain`
- `internal/service`
- `internal/errors`

Actual layout:
- `internal/engine`
- `internal/model`
- `internal/store`
- `internal/dberr`

Impact: A reader following the design doc would create wrong import paths and be confused. The implementation notes (02) correctly list the real packages, so this is an internal doc drift not affecting the shipped code, but still a documentation inaccuracy.

Severity: MEDIUM. Reason: Misleading for future contributors trying to understand the intended structure from the design artifact.

Suggested fix: Either update 01-design.md to reflect the true package names, or annotate that the names are logical placeholders and map to the actual paths given in 02-implementation-notes.md.

## GAP 3 — DOC_CODE_MISMATCH: engineering/03-execution-result.md omits one test

Location: `labs/27-database-constraints/engineering/03-execution-result.md`, "Tests" section.

Lists 7 tests, omits `TestConcurrentRegistration_Unsafe_SuffersRaceCondition`.

Impact: The execution record under-reports the test suite. A technical writer relying on this file would think one of the success criteria (concurrency race demonstration) is missing verification, when in fact the test exists and passes. This misaligns traceability.

Severity: MEDIUM. Reason: Omits a core piece of evidence (the unsafe race test) from the engineering log, potentially causing downstream audits to question completeness.

Suggested fix: Amend 03-execution-result.md to include the omitted test in the list and show its PASS status (same as the others). The result block already shows `PASS` for the internal/store package — adding the test name costs nothing and improves fidelity.

## GAP 4 — MISSING_EDGE_CASE: Mixed registration paths (full-UNIQUE vs partial-UNIQUE) not tested / documented

Location: `internal/engine/engine.go` — two independent uniqueness indexes (`emailIndex` and `activeEmails`) updated only by their respective call sites (`InsertUser(...,false)` vs `InsertUser(...,true)`).

Impact: If a caller registers the same email via `RegisterUserPartial` (active → goes to `activeEmails`) and then later via `RegisterUser` (full unique → goes to `emailIndex`), the second registration will **not** see a conflict despite there being an active row with that email (only visible in the other index). This creates a silent duplicate-active scenario, violating the intent that "email is unique among active rows".

Severity: MEDIUM. Reason: The demo and tests never mix the two registration methods, so this latent issue is not exercised. It contradicts the design goal of a global uniqueness guarantee for active emails (claimed in 01-design.md: "UNIQUE / PRIMARY KEY constraints enforce single occurrence across rows").

Note: This is a design simplification; the lab presents the two paths as separate demos. But the claim of UNIQUE preventing duplicates across rows is over-broad given the implementation.

Suggested fix: Either
(a) Change the engine to maintain a single logical uniqueness decision (e.g., always check both indexes, or have one source of truth), **or**
(b) Document that `RegisterUser` and `RegisterUserPartial` are mutually exclusive namespaces and must not be mixed for the same logical key, **or**
(c) Add a test demonstrating the conflict (or lack thereof) and label it a known limitation.

## GAP 5 — IMPLEMENTATION_OVERCLAIM: PRIMARY KEY constraint not demonstrated or tested

Location: `labs/27-database-constraints/engineering/01-design.md`, Expected Behavior #2: "UNIQUE / PRIMARY KEY constraints reject duplicate keys (23505)" and Architecture bullet 32: "...primary & unique index checking...".

Impact: The code models PRIMARY KEY implicitly via a monotonically increasing sequence (`userSeq`, `orderSeq`). There is no way to insert a duplicate fixed PK value because the engine always generates the ID if zero. There is no test asserting rejection of a user-supplied duplicate PK (e.g., manually setting `u.ID = 1` twice). If one tried to insert a pre-existing PK, the engine would simply overwrite the row (since `e.users[u.ID] = u` is an unconditional map assign under the lock — **a bug!**).

Wait: let’s double-check `InsertUser`:

```go
// Primary Key generation
if u.ID == 0 {
    u.ID = e.userSeq.Add(1)
}
// Commit row and indexes
e.users[u.ID] = u
```

If `u.ID != 0` (caller supplied), we skip the sequence and write to `e.users[u.ID]`. This **overwrites** whatever was there. There is no uniqueness check on the PK! That means:

- The engine does NOT enforce PRIMARY KEY uniqueness when the caller supplies an ID.
- There is no test for this scenario.

Claim in design: PK enforces uniqueness. Reality: only true if you always let the engine assign the ID (zero). If you supply your own, you can silently replace rows.

Severity: MEDIUM to HIGH depending on threat model. Reason: The docs (design) claim PK gives you uniqueness; the code only gives it if you obey the convention (ID==0). A caller who mistakenly supplies an existing ID causes data loss (row replacement) without any error.

Note: This is benign in the current usage because all callers (`UnsafeStore`, `SafeStore`) always pass zero ID. But the claimed guarantee is broader than the implementation supports.

Suggested fix: Either
(a) Add a check: if `u.ID != 0 && e.users[u.ID] != nil` then return a PK violation (SQLSTATE 23P01? or treat as UNIQUE on PK), **or**
(b) Downgrade the claim: document that PK identity is generated by the engine, not validated, **or**
(c) Add a test that attempts to insert a duplicate PK and expects an error (and then fix the engine to return that error). Given the lab's focus is on UNIQUE/FK/NOT NULL/CHECK, option (b) may be simplest: adjust the design doc to say PK is *generated* not *validated*.

Given that the README and implementation notes never mention PK as a demonstrated constraint (only the design doc does), I lean to treating this as a design-doc overclaim rather than a broken implementation. Still, it is a gap between claimed and actual behavior.

For GAP logging, I will use `IMPLEMENTATION_OVERCLAIM` with severity MEDIUM.

## GAP 6 — RACE_CONDITION: UnsafeStore race dependent on `time.Sleep` (flakiness risk if removed)

Location: `internal/store/store.go:38` (`time.Sleep(1 * time.Millisecond)`).

Impact: The test `TestConcurrentRegistration_Unsafe_SuffersRaceCondition` relies on this sleep to widen the window so that the unsafe path *reliably* loses (count > 1). If the sleep were removed or reduced, the test could become flaky (sometimes passing when it should fail). The sleep is an intentional crutch for demonstration, but it makes the test less pure and ties correctness to a timing assumption.

Severity: LOW (because the sleep is part of the production code, not just the test). However, if someone attempted to "clean up" the unsafe store by removing the sleep (thinking it’s dead code), the test might start failing intermittently and hide the race. So it documents a latent risk.

Note: The auditor observed that the test passes reliably with the sleep present. No flakiness witnessed in repeated runs.

Suggested fix: None required; just flag that the race demonstration leans on a deliberate delay. In production-grade test one would use a barrier or channel to synchronize the check-phase across goroutines before releasing them to insert. Acceptable for a teaching lab.

## Summary of Gears (non-exhaustive, only actual)

List in required format:

- MISSING_TEST: MapToDomainError output never asserted
- DOC_CODE_MISMATCH: engineering/01-design.md package tree does not match reality
- DOC_CODE_MISMATCH: engineering/03-execution-result.md omits TestConcurrentRegistration_Unsafe_SuffersRaceCondition
- MISSING_EDGE_CASE: Mixed registration paths (full-UNIQUE vs partial-UNIQUE) not tested / documented
- IMPLEMENTATION_OVERCLAIM: PRIMARY KEY constraint not demonstrated (PK overwrites on caller-supplied ID)
- RACE_CONDITION: UnsafeStore race dependent on `time.Sleep` (flakiness risk if removed)