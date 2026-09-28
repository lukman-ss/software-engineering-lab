# Source Map Audit

## File Reviewed
`content/06-source-map.md` — 109 lines

## Verification of Line References

| Reference | Actual Location | Status |
|---|---|---|
| store.go:24-47 (UnsafeStore.RegisterUser) | store.go:24-47 | PASS |
| store_test.go:210-239 (Unsafe test) | store_test.go:210-239 | PASS |
| store_test.go:162-208 (Safe concurrency test) | store_test.go:162-208 | PASS |
| store_test.go:114-160 (PartialUniqueIndex) | store_test.go:114-160 | PASS |
| store_test.go:133-157 (soft delete flow) | store_test.go:133-157 | PASS |
| store_test.go:241-261 (ErrorClassification) | store_test.go:241-261 | PASS |
| store_test.go:16-159 (constraint tests) | store_test.go:16-159 | PASS |
| engine.go:45-99 (InsertUser) | engine.go:45-99 | PASS |
| engine.go:45-148 (combined) | engine.go:45-148 | PASS |
| engine.go:71-77, 101-120 (partial index) | engine.go:71-77, 101-120 | PASS |
| errors.go:83-100 (MapToDomainError) | errors.go:83-100 | PASS |
| cmd/demo/main.go:59-94 (race demo) | main.go:59-94 | PASS |
| cmd/demo/main.go:43-56 (partial index demo) | main.go:43-56 | PASS |
| cmd/demo/main.go:51-56 (soft delete demo) | main.go:51-56 | PASS |

## Research Reference Verification

| Reference | Actual | Status |
|---|---|---|
| research/05-report.md Finding 1 (race) | present | PASS |
| research/05-report.md Finding 10 (NOT VALID) | present | PASS |
| research/05-report.md Limitations, Finding 5, Finding 4 | present | PASS |
| research/02-sources.md "all 10 sources" | 10 accessible sources (1-10) + blocked | PASS |
| research/04-contradictions.md (MySQL, partial index) | present | PASS |
| research-audit/07-verdict.md APPROVED | APPROVED | PASS |
| engineering-audit/06-verdict.md APPROVED | APPROVED | PASS |
| engineering-audit-opensource findings | 3 LOW issues documented | PASS |

## Issues
None. All 14 line references and 8 research references resolve correctly.

## Verdict
VERIFIED
