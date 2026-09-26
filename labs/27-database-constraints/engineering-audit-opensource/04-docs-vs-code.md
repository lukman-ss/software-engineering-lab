# Docs vs Code

## Summary

The README accurately describes the five implemented constraints and the SQLSTATE taxonomy. The demo and engineering execution-record reproduce the README's claims exactly. Core behavior is verified.

## Matches (PASS)

| README / Doc claim | Implementation | Verified |
|---|---|---|
| NOT NULL (23502) | engine.go:33-43, 110-117 | go test PASS, go test -race PASS |
| CHECK (23514) | engine.go:46-55, 119-122 | go test PASS |
| UNIQUE (23505) | engine.go:65-69 | go test + concurrency test PASS |
| FOREIGN KEY (23503) | engine.go:124-127 | go test PASS |
| Partial UNIQUE INDEX | engine.go:58-64, 79-83, 102-106 | go test PASS |
| 50 concurrent → 1 success | cmd/demo/main.go:59-94; store_test.go:162-208 | both PASS |
| `go test -v ./...` | README:17-19 | matches actual output |
| `go test -race ./...` | README:23-25 | matches actual (clean) |
| `go run ./cmd/demo` | README:31-33 | matches actual output |
| engineering/03-execution-result.md | cmd/demo output | byte-for-byte match |

## Doc vs Code Mismatches

### DOC_CODE_MISMATCH 1 (MEDIUM): UnsafeStore documented but misleading
- engineering/02-implementation-notes.md:9 ("Application store demonstrating safe enforcement using engine constraints vs vulnerable patterns")
- internal/store/store.go:12-13 ("UnsafeStore lacks database constraints ... vulnerable to read-then-write race conditions")
- Reality: UnsafeStore calls InsertUser which STILL enforces uniqueness (Finding 1 of code-audit). The comment "bypassing unique constraints" (store.go:40) is false. README itself does not mention UnsafeStore, so README-vs-code is fine; this is an implementation-note vs code mismatch.

### TEST_CLAIM_MISMATCH 1 (MEDIUM): UnsafeStore untested
- Notes claim "Comprehensive unit and concurrency race condition tests verifying each constraint behavior" (implementation-notes.md:9).
- UnsafeStore is not exercised by any test. "vulnerable patterns" are not verified.

### DOC_CODE_MISMATCH 2 (LOW): SQLSTATE taxonomy verification is synthetic
- cmd/demo/main.go:98 constructs a standalone `dberr.NewUniqueViolation` to assert `IsConstraintViolation`, rather than deriving the value from an actual failed registration. The line demonstrates the predicate function, not the store's error pipeline.
- engineering/03-execution-result.md:82 echoes the demo output faithfully (not a fabrication).

### RESEARCH_DEMO_MISMATCH 1 (LOW): Row-level lock mechanism differs
- research/05-report.md:50-55 ("row-level locks acquired during INSERT/UPDATE ... B-tree index ... row-exclusive lock").
- Implementation uses a single coarse-grained `sync.RWMutex` (engine.go:14). Documented as a trade-off in implementation-notes.md:25. Behaviorally equivalent; mechanism differs. No falsification.

## Conclusion

README matches code and tests. The only real mismatches are (a) the UnsafeStore's false "bypassing" comment and (b) the gap in asserting SQLSTATE on store-returned errors. Neither invalidates the README's core claims, which are all reproduced by execution.