# Code Audit

Target Lab: labs/26-contract-testing
Date: 2026-09-28
Scope: implementation + tests only. No code modified.

## Finding 1

Location: internal/contract/verifier.go:61-128 (`Verify`)
Claimed Behavior: Verifies status, headers, subset JSON schema against live provider.
Observed Implementation: Builds request from interaction, checks status code, case-insensitive header compare (`EqualFold`), decodes body with `UseNumber`, delegates to `diffValues` subset check. All error branches append errors and set `Passed=false`.
Assessment: PASS
Severity: LOW
Notes: Stateless per call. No shared mutable state.

## Finding 2

Location: internal/contract/verifier.go:130-189 (`diffValues`)
Claimed Behavior: Minimal-subset match; exact enum case; type-strict numbers; missing-field detection.
Observed Implementation: Nil expected returns no diffs. Maps checked subset-only (extra provider fields ignored by design). `json.Number` vs `json.Number` compared by string; number-vs-non-number flagged type mismatch. Other primitives compared by `reflect.TypeOf` then value.
Assessment: PASS
Severity: LOW
Notes: Extra provider fields (`notes`) correctly ignored. Breaking triple (status value, `customer.name` missing, `total` type) each produces distinct diff. Matches demo output observed.

## Finding 3

Location: internal/contract/verifier.go:84-106 (status-mismatch path)
Claimed Behavior: Status mismatch fails verification.
Observed Implementation: On status mismatch records error but continues to read and JSON-decode body, potentially appending a second noise error (e.g. empty 500 body yields "not valid JSON object").
Assessment: WARNING
Severity: LOW
Notes: Gate outcome unaffected (`Passed=false` either way). Noise only. No repair made.

## Finding 4

Location: internal/consumer/client.go:37-114
Claimed Behavior: Fetches `/v1/orders/{id}`, enforces `customer.name` presence and `IN_PROGRESS`/`COMPLETED` enum, builds consumer contract.
Observed Implementation: 5s HTTP timeout present. Non-200, read, and parse errors wrapped. Empty `customer.name` rejected. Unknown status rejected. `GenerateMobileContract` emits 1 interaction with `id`, `status`, `customer.name`, `total` as `json.Number`.
Assessment: PASS
Severity: LOW
Notes: Revision GAP-4 (bare `http.Client{}`) verified resolved.

## Finding 5

Location: internal/provider/server.go:18-131
Claimed Behavior: V1 compliant; Breaking carries 3 mutations; Dual preserves V1 plus V2.
Observed Implementation: V1 returns `IN_PROGRESS`, `customer.name`, int `150000`. Breaking returns `in_progress`, `customer.full_name`, string `"150000"`. Dual serves identical V1 on `/v1/orders/` and V2 (`full_name`, string total, `currency`) on `/v2/orders/`. All handlers reject non-GET, set JSON content type.
Assessment: PASS
Severity: LOW
Notes: Stateless handlers. Safe under concurrent use.

## Finding 6

Location: internal/model/order.go
Claimed Behavior: V1, Breaking, V2 DTOs as documented.
Observed Implementation: Matches claims exactly, including `Notes omitempty` on V1/Breaking and `Currency` on V2.
Assessment: PASS
Severity: LOW
Notes: None.

## Finding 7

Location: cmd/demo/main.go
Claimed Behavior: 4-stage lifecycle demo with CI gate messaging.
Observed Implementation: Real `httptest` servers, real `Verify` calls. Stage 2 expects pass, Stage 3 expects fail (exits 1 only on unexpected pass), Stage 4 expects pass. No fabricated output; live run reproduced claimed transcript modulo diff ordering.
Assessment: PASS
Severity: LOW
Notes: Gate "blocking" is represented by log + branch logic, exit code stays 0 on the expected-block path. See 04-docs-vs-code.md.
