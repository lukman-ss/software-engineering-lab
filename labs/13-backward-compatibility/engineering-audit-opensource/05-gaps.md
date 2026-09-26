# Gap Analysis

Scope: Labs/13-backward-compatibility. Gaps are scoped to the in-memory demonstration (documented `ponytail:` boundary), not fabrications.

## Gap Inventory

1. **IDEMPOTENT_BACKFILL_CHECK_DUPLICATE_NUMBER** — Backfill idempotency relies on `SavePhoneEntry`'s duplicate-number check rather than a `WHERE new_field IS NULL`-style migration predicate. Functionally correct; minor redundancy with backfill's `len(phones)==0` guard. Not a defect.

2. **UNHANDLED_ERROR: extra-phone `SavePhoneEntry` discard** — `internal/compat/service.go:80`: `_, _ = s.store.SavePhoneEntry(u.ID, extra, false)`. Extra (non-primary) phones written during dual-write swallow errors. In-memory store will not fail, but the "atomic dual-write" claim is weakened (Finding 14).

3. **UNHANDLED_ERROR: `id, _ := strconv.Atoi(idStr)`** — `handler.go:22,46`. Invalid/missing id yields `id=0` and a NotFound rather than a 400. Cosmetic for demo; acceptable for lab.

4. **O_N_SQUARE_SORT_GETUSERIDS** — `store.go:184-189` uses manual insertion/bubble sort instead of `sort.Ints`. Performance only; not a concurrency issue.

5. **CUMULATIVE_CONTRACT_GUARD_COUNTER** — `LegacyReadHits` never reset, so the real-time zero-traffic window cannot be measured over a sliding period (research Q11). Guard is stricter than the heuristic, which over-approximates caution; no correctness regression.

6. **IN_MEMORY_PERSISTENCE** — Documented `ponytail:` boundary: data & flag state lost on restart; backfill resume only works within process. Not scoped for persistence; documented in Known Limitations.

7. **MISSING_TEST: CONTRACT_PRECONDITION_ERROR_PATH** — No unit test where `GetModernUser` is called when `IsContractApplied()` (returns ErrLegacyUnavailable only via V1 path; modern path unaffected but untested post-contract extra-phone).

8. **MISSING_TEST: GETUSERIDS_ORDERING** — Backfill test does not assert ordering correctness; only total migrated count. Not a behavioral gap (checkpoint is monotonic by construction).

## Missing Tests (vs design Success Criteria coverage)
- Batch boundaries: covered (TestBackfillIdempotentAndResumable batchSize=3).
- Empty records: `CreateModern` rejects empty phones (ErrMissingRequired, covered).
- Concurrency: covered under `-race`.
- Failure path: contract violation, drift, unsafe rollback — all covered.
- Negative cases: dual-write errors (injected via guard), legacy-unavailable after contract covered.

## Research Mismatches
- None. Implementation faithfully reflects research questions 6, 7, 8, 9, 13, 14, 16, 17.

## Verdict on Fabrication
- No fake benchmarks.
- No fake demo output (re-executed verbatim).
- No fake results; metrics snapshot matches code.

## Severity Summary
- None CRITICAL.
- None HIGH.
- Medium: #2 (silent extra-phone error discard) — weakens "atomic" dual-write phrasing.
- Low: #3, #4, #5, #6 (documented scope), #7, #8.

## Conclusion
No blocking gaps against the Expand-Migrate-Contract spec. Warnings are implementation-bounded simplifications explicitly documented in `engineering/02-implementation-notes.md` (Known Limitations) and `design.md` (`ponytail:`), not fabrication.
