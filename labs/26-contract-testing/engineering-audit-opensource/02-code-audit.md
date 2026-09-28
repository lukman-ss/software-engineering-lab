# Code Audit

Target Lab: labs/26-contract-testing
Date: 2026-09-28

## Finding 1

Location: internal/contract/verifier.go:117-176 (`diffValues`)
Claimed Behavior: Minimal subset verification; detects enum casing, field rename, type mutation; ignores extra provider fields.
Observed Implementation: Recursive map-subset walk; numbers compared as `json.Number` strings with strict type gate (`json.Number` vs `string` → type mismatch); strings compared by type + exact value; missing keys reported; `nil` expected short-circuits to no diff.
Assessment: PASS
Severity: LOW
Notes: Strict `json.Number` string equality means `150000` vs `150000.0` would flag mismatch. No impact here (contract `json.Number("150000")`, provider encodes `int64(150000)` → identical). Breaking-change detection confirmed live via demo (3/3 diffs).

## Finding 2

Location: internal/contract/verifier.go:58-115 (`Verify`)
Claimed Behavior: Verifies method, path, headers, status, response schema/types/enums (per 01-design.md §Components, 02-implementation-notes.md).
Observed Implementation: Verifies method+path (via request execution), status code, body subset. Request headers sent; **response headers declared in contract (`Content-Type`) are never checked**.
Assessment: WARNING
Severity: LOW
Notes: DOC_CODE_MISMATCH (recorded in 04-docs-vs-code.md). No test covers a wrong-Content-Type provider passing verification. Fix (out of scope): enforce header subset or document header-agnostic intent.

## Finding 3

Location: internal/contract/verifier.go:51-55, internal/consumer/client.go:26-31
Claimed Behavior: HTTP verification and consumer fetch.
Observed Implementation: Both use bare `&http.Client{}` — no timeout, no context.
Assessment: WARNING
Severity: LOW
Notes: Against `httptest` servers no hang observed; real-provider use could block indefinitely. Optional hardening only.

## Finding 4

Location: internal/provider/server.go:89-131 (`ProviderDual.ServeHTTP`)
Claimed Behavior: V1 path contract-compliant; V2 path new schema.
Observed Implementation: `/v1/orders/` prefix checked first, serves `OrderResponseV1`; `/v2/orders/` serves `OrderResponseV2`; no path overlap; non-GET → 405; unknown → 404. V1 e2e + verification pass confirmed by test + live demo.
Assessment: PASS
Severity: LOW
Notes: `/v2` shape has no dedicated assertion (test covers V1 path only). Recorded as MISSING_TEST in 05-gaps.md.

## Finding 5

Location: internal/consumer/client.go:34-79 (`FetchOrder`)
Claimed Behavior: Consumer parses provider payload; fails loudly on contract violation.
Observed Implementation: Non-200 → error; body read errors wrapped; JSON unmarshal errors wrapped; empty `customer.name` → error; unknown `status` → error. Breaking-provider client failure proven by `TestProviderBreaking_ContractVerification_Fails` (expects non-nil err).
Assessment: PASS
Severity: LOW
Notes: Status allowlist (`IN_PROGRESS`, `COMPLETED`) is stricter than verifier (exact-match on `IN_PROGRESS` only for this fixture) — consistent, no conflict.

## Finding 6

Location: internal/provider/server.go:39,75,108,126; internal/contract/verifier.go:88
Claimed Behavior: JSON encode/decode with error handling.
Observed Implementation: Server-side `json.NewEncoder(w).Encode` errors discarded (`_ =`); verifier side decode errors recorded as verification errors; body-close error discarded.
Assessment: PASS
Severity: LOW
Notes: Encoding a static struct to `httptest` ResponseWriter cannot realistically fail; verifier (the asserting side) handles its errors. Acceptable for lab scope.

## Finding 7

Location: internal/model/order.go
Claimed Behavior: V1 / Breaking / V2 DTOs with exactly the 3 documented breaking deltas.
Observed Implementation: `OrderResponseBreaking` differs from V1 in exactly `status` value casing source (`in_progress`), `customer.full_name` rename, `total string` type; V2 adds `currency`. Matches design §Failure Scenario.
Assessment: PASS
Severity: LOW
Notes: None.

## Finding 8

Location: cmd/demo/main.go:30-68
Claimed Behavior: CI gate — exit non-zero on unexpected V1 fail or unexpected breaking pass.
Observed Implementation: `os.Exit(1)` on both inverted outcomes; breaking-fail path prints diffs and exits 0 (gate correctly *blocks deploy* while demo itself succeeds). Live run completed all 4 stages, exit 0.
Assessment: PASS
Severity: LOW
Notes: Demo output real, reproduced independently; matches engineering/03-execution-result.md modulo map-ordering of the 3 breaking diffs (Go map iteration — nondeterministic order, tests correctly assert count not order).
