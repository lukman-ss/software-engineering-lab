# Gap Analysis

| # | Gap Type | Location | Severity | Description |
|---|---|---|---|---|
| 1 | DOC_CODE_MISMATCH | engineering/01-design.md vs internal/compat/handler.go | LOW | Design doc describes header values as `@epoch`; code emits `Deprecation: true` + concrete Sunset date string. Intent met; format shorthand inaccurate. |
| 2 | DOC_CODE_MISMATCH | README.md (RFC citation) | LOW | README labels both `Deprecation` and `Sunset` as "RFC 8594"; `Deprecation` header is RFC 9224. Cosmetic/accuracy note. |
| 3 | IMPLEMENTATION_OVERCLAIM | engineering/01-design.md + service.go:229 | LOW | "legacy traffic dropped to zero" guard is implemented as a cumulative counter with no reset/sliding window. In prod this blocks the guard after any legacy hit; lab uses documented `force` escape. Fidelity gap, not a bug. |
| 4 | MISSING_EDGE_CASE / MISSING_TEST | internal/compat/handler.go:20-22 | LOW | HTTP handlers discard `strconv.Atoi` error for `id`; invalid/missing id silently becomes 0 -> 404 instead of 400 Bad Request. No test covers this path. |
| 5 | MISSING_EDGE_CASE | store.go GetUserIDs:184-189 | LOW | O(n^2) bubble sort used instead of stdlib `sort.Ints`. Correct output but needlessly inefficient; trivially replaceable. |
| 6 | MISSING_TEST | internal/compat/store.go:147 (SavePhoneEntry) | LOW | Idempotency at storage-method level (return existing entry on duplicate number) not unit-tested; only indirectly asserted via backfill rerun count. |
| 7 | MISSING_TEST | internal/compat/service.go (WriteNewOnly branch) | LOW | `CreateUser` `WriteNewOnly` path not exercised end-to-end via service in tests (only WriteDual/LegacyOnly + direct CreateModern). |

Summary: 0 CRITICAL, 0 HIGH, 0 MEDIUM, 7 LOW. All gaps are non-blocking polish/coverage/accuracy items. No fabricated results; no unproven core behavior; no data-race failures.
