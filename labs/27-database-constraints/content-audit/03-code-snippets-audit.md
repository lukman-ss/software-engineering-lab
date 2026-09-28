# Code Snippets Audit

## File Reviewed
`content/03-code-snippets.md` — 11 snippets, 394 lines

## Method
Each snippet compared verbatim against source file + line range cited.

| # | Snippet | Cited Source | Actual Source | Match |
|---|---|---|---|---|
| 1 | SQLSTATE constants | dberr/errors.go:10-18 | errors.go:10-18 (7 consts) | PASS verbatim |
| 2 | ConstraintError struct | dberr/errors.go:20-37 | errors.go:20-37 | PASS verbatim |
| 3 | NewNotNullViolation + NewForeignKeyViolation | dberr/errors.go:39-46, 66-73 | errors.go:39-46, 66-73 | PASS verbatim (two funcs combined) |
| 4 | MapToDomainError | dberr/errors.go:83-100 | errors.go:83-100 | PASS verbatim |
| 5 | Engine struct | engine/engine.go:12-30 | engine.go:12-30 | PASS verbatim |
| 6 | InsertUser atomic | engine/engine.go:45-99 | engine.go:45-99 | PASS verbatim (incl. comments, constraint names) |
| 7 | InsertOrder | engine/engine.go:122-148 | engine.go:122-148 | PASS verbatim |
| 8 | UnsafeStore.RegisterUser | store/store.go:23-47 | store.go:23-47 | PASS verbatim |
| 9 | SafeStore.RegisterUser | store/store.go:58-66 | store.go:58-66 | PASS verbatim |
| 10 | SoftDeleteUser | engine/engine.go:101-120 | engine.go:101-120 | PASS verbatim |
| 11 | TestConcurrentRegistration_Safe | store/store_test.go:162-208 | store_test.go:162-208 | PASS verbatim |

## Issues

### W-04 LOW — Snippet 10 title typo
Header reads "Soft-Drop User" should be "SoftDelete User" (function name is SoftDeleteUser). Cosmetic.

### W-05 LOW — Snippet 3 line range notation
Cites "39-46, 66-73" as single snippet; two non-contiguous ranges merged. Not inaccurate but non-standard citation.

### W-06 INFO — HTTP mapping in explanation
Snippet 4 explanation adds "409 Conflict for unique, 400 for not null, 422 for check" — not in source code; reasonable domain mapping but not implemented as HTTP layer in lab. Documented as interpretation, not claimed as code.

## Truncation / Omission
None. No logic hidden behind ellipsis that would change semantics.

## Hallucinations
None. All code exists at cited paths.

## Verdict for this file
VERIFIED_WITH_WARNINGS (cosmetic warnings only)
