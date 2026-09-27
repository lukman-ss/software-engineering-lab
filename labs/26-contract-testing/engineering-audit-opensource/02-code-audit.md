# Engineering Audit — Code Audit

## Finding 1

Location: internal/contract/verifier.go:51-54, 58
Claimed Behavior: Verifier executes contract HTTP interactions against provider base URL and reports precise diffs.
Observed Implementation: `Verifier.Client` is `&http.Client{}` with no timeout. `Verify` iterates interactions, builds request, sets request headers from contract, checks status code, decodes response with `json.Number`, recursively diffs expected vs actual.
Assessment: PASS with WARNING
Severity: MEDIUM
Notes: No per-request timeout on the HTTP client. For a CI gate this is acceptable over trusted localhost, but a hung provider stalls CI indefinitely. Recommend `Timeout`.

## Finding 2

Location: internal/contract/verifier.go:30, 70-72
Claimed Behavior (implementation-notes): "Minimal Subset Verification: ... validates ... headers".
Observed Implementation: `ResponseDefinition.Headers` is declared and serialized, but `Verify` never reads `Response.Headers`. Request headers (line 70-72) are SET on the outgoing request, but response headers are NOT validated.
Assessment: FAIL
Severity: HIGH
Notes: Implementation notes explicitly claim response header validation ("comparing expected contract interactions, headers, status codes"). The README structure comment for `internal/contract` says "comparing ... headers". Headers field is populated in demo contract and tests but never asserted. Breaking-change detection omits header mismatches.

## Finding 3

Location: internal/contract/verifier.go:147-161
Claimed Behavior: Type mismatch detected for `total` (int expected vs string actual).
Observed Implementation: `diffValues` checks `json.Number` vs non-Number via `isExpNum || isActNum`. Contract stores `total` as `json.Number("150000")`; provider breaking returns `"150000"` (string). Decoder uses `UseNumber()` on both expected and actual, so actual `"150000"` decodes to string, not json.Number. Correctly flagged as type mismatch.
Assessment: PASS
Severity: LOW
Notes: Correct; note that the type-mismatch branch (line 164-168) is a fallback only reached for non-number primitives.

## Finding 4

Location: internal/contract/verifier.go:124-144
Claimed Behavior: Nested missing field `customer.name` detected.
Observed Implementation: Recursive map walk over expected keys; if key absent in actual, emits "missing expected field". For breaking provider, `customer.full_name` present but expected `name`, so `customer.name` reported missing. Correct.
Assessment: PASS
Severity: LOW

## Finding 5

Location: internal/contract/verifier.go:170-173
Claimed Behavior: Exact enum case comparison `IN_PROGRESS` vs `in_progress`.
Observed Implementation: Both decode to `string`; `expected != actual` triggers value mismatch with `%q` formatting. Exact case sensitivity enforced.
Assessment: PASS
Severity: LOW

## Finding 6

Location: internal/provider/server.go:89-131
Claimed Behavior: Dual provider exposes V2 endpoint.
Observed Implementation: `ProviderDual.ServeHTTP` handles `/v1/orders/` and `/v2/orders/` with distinct DTOs (V1 compliant, V2 with full_name/total-string/currency).
Assessment: PASS with WARNING
Severity: MEDIUM
Notes: V2 endpoint is routed but the consumer contract + tests never exercise `/v2/orders/{id}`. V2 path is reachable but its schema is not asserted by any test or the demo. This is "safe evolution" in principle but the V2 route is unverified. Per YAGNI/success-criteria #5, V2 should be demonstrably verified.

## Finding 7

Location: internal/consumer/client.go:34-79
Claimed Behavior: Mobile client enforces contract at runtime (status enum, customer.name presence, total type).
Observed Implementation: `FetchOrder` unmarshals into strict struct, validates `Customer.Name != ""`, validates status in {IN_PROGRESS, COMPLETED}, returns wrapped errors. `total` is `int64` — unmarshalling breaking provider's `"150000"` string yields json error, surfacing contract violation.
Assessment: PASS
Severity: LOW

## Finding 8

Location: internal/consumer/client.go:82-110
Claimed Behavior: Consumer generates minimal contract with exact fields.
Observed Implementation: Interaction declares path `/v1/orders/ORD-123`, GET, status 200, Content-Type header, body with id/status/customer.name/total. Matches design success-criteria minimal subset.
Assessment: PASS
Severity: LOW

## Finding 9

Location: internal/model/order.go:14-49
Claimed Behavior: V1, Breaking, V2 DTOs reflect distinct schemas.
Observed Implementation: V1 int total + name; Breaking string total + full_name + in_progress; V2 string total + full_name + currency. Enum comment on V2 says "in_progress". JSON tags consistent with provider handlers.
Assessment: PASS
Severity: LOW

## Finding 10

Location: internal/contract/verifier.go:120-130
Claimed Behavior: Minimal subset — extra provider fields ignored.
Observed Implementation: `diffValues` only iterates EXPECTED keys; extra actual keys are never visited. Extra `notes` field on V1 ignored. Correct minimal-subset semantics.
Assessment: PASS
Severity: LOW

## Finding 11

Location: internal/contract/verifier.go:98-102
Claimed Behavior: Non-object JSON response handled.
Observed Implementation: `decoder.Decode(&actualBody map[string]interface{})` — if response body is a JSON array or primitive, decode into map fails, reported as "response is not valid JSON object". Contract expects object, so this is acceptable.
Assessment: PASS
Severity: LOW
