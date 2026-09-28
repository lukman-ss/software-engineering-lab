# Master Draft Audit

## File Reviewed
`content/02-master-draft.md` (228 lines)

## Verification Result
VERIFIED_WITH_WARNINGS

## Cross-Check vs Research

| Claim in Master Draft | Research Source | Status |
|---|---|---|
| Check-then-act fails under concurrency | 05-report Finding 1, 02-sources 5.5 + 13.3 | PASS |
| UNIQUE uses B-tree index atomically under ROW EXCLUSIVE | 05-report Finding 2, 02-sources 5.5.3 + 11.6 + 13.3 | PASS |
| Quote: adding unique constraint creates unique B-tree index | 02-sources 5.5.3 verbatim | PASS |
| SQLite verifies NOT NULL/CHECK only on INSERT/UPDATE | 02-sources SQLite section, 03-evidence 6 | PASS |
| Defense-in-depth mental model (app UX vs DB correctness) | 05-report Executive Summary + Finding 15 | PASS |
| SQLSTATE codes 23502/23503/23505/23514 Class 23 | 02-sources Appendix A, 05-report Finding 7 | PASS |
| CHECK row-scoped, cannot reference other tables | 05-report Finding 5, 02-sources 5.5.1 | PASS |
| Partial unique index WHERE deleted_at IS NULL | 05-report Finding 6, 02-sources 11.8 | PASS |
| FK does not auto-index referencing columns | 05-report Finding 4, 02-sources 5.5.5 | PASS |
| SERIALIZABLE for multi-row invariants (context) | 05-report Finding 8 | PASS |
| NOT VALID + VALIDATE CONSTRAINT not implemented, research-only | 05-report Finding 10, 02-sources ALTER TABLE | PASS |
| NULL distinct in UNIQUE | 05-report Findings 3/11, 02-sources 5.5.3 + 11.6 | PASS |
| Partitioned table limitation (not in master body) | 05-report Finding 11 | NOT COVERED — correctly omitted from implemented scope |

## Cross-Check vs Engineering

| Claim | Code File | Status |
|---|---|---|
| UnsafeStore scan → sleep → InsertUserUnsafe | store.go:24-47 | PASS verbatim |
| SafeStore delegates to InsertUser(false) | store.go:59-66 | PASS |
| SafeStore 20 goroutine test 1/19, demo 50 goroutine 1/49 | store_test.go:162-208, demo/main.go:59-94, execution-result.md | PASS |
| Engine.InsertUser Lock + NOT NULL → CHECK → UNIQUE → PK → commit | engine.go:45-99 | PASS |
| Engine.InsertOrder NOT NULL → CHECK → FK | engine.go:122-148 | PASS |
| Concrete rules (age>=18, status IN, total_cents>0, etc.) 9 mappings | engine.go + dberr/errors.go | PASS exact constraint names |
| Two maps emailIndex / activeEmails | engine.go:12-21 | PASS |
| MapToDomainError 23505→conflict 23502→invalid 23514→validation 23503→reference | dberr/errors.go:83-100 | PASS |
| Simulator coarse-grained lock not B-tree latch, zero-dependency | 02-implementation-notes Trade-offs | PASS, disclosed 3× |
| Demo output exactly 1 success 49 rejections | execution-result.md demo section | PASS |
| 8 tests pass + -race | execution-result.md, store_test.go | PASS |

## Issues

### W-01 LOW — Takeaways count mismatch
Master draft line "Lihat 05-key-takeaways.md untuk 8 poin ringkas." but 05-key-takeaways.md contains 10 points. Inconsistency between brief (8) and actual file (10).

### W-02 LOW — Duplicated sentence
Lines 77 and 85 identical: "Implementasi lab memakai coarse-grained table lock, bukan B-tree page latch — trade-off yang didokumentasikan demi readability dan zero-dependency." Duplicated verbatim.

### W-03 LOW — Evaluation order overstatement
Section Implementation states "PostgreSQL mengacu pada urutan pemeriksaan: NOT NULL sebelum CHECK, CHECK diurutkan secara alfabetis" and implies FK/UNIQUE ordering. Research 04-contradictions Open Question 3 notes FK vs UNIQUE order not explicitly documented; content presents assumed order as confirmed.

## Missing / Not Covered
- Partitioned table UNIQUE limitation (Finding 11) not discussed in master body — acceptable as peripheral, documented in research only.

## Hallucinations
None. No invented facts, no platform bias (MySQL correctly marked unverified).

## Recommendation
Accurate and complete for implemented scope. Fix W-01 count reference.
