# Code Audit

Target Lab: `labs/26-contract-testing`

## Finding 1

Location: `internal/contract/verifier.go:58-115`
Claimed Behavior: Verifier queries provider HTTP endpoint, decodes JSON using `decoder.UseNumber()`, and diffs required subset against actual body.
Observed Implementation: Uses standard `http.Client`, checks status codes, properly closes response bodies (`_ = resp.Body.Close()`), and delegates recursive schema matching to `diffValues`.
Assessment: PASS
Severity: LOW
Notes: Robust implementation using Go stdlib.

## Finding 2

Location: `internal/contract/verifier.go:117-176`
Claimed Behavior: Accurate diffing of objects, numbers, and types for consumer subset validation.
Observed Implementation: Handles nested map keys, reports missing fields, explicitly handles `json.Number` numeric equality, and validates primitive type identities.
Assessment: PASS
Severity: LOW
Notes: Handles non-breaking additive fields correctly by asserting consumer subset expectations only.

## Finding 3

Location: `internal/consumer/client.go:34-79`
Claimed Behavior: Consumer client parses provider response into `MobileOrderSummary` and validates contract expectations locally.
Observed Implementation: Validates response status code, unmarshals JSON, verifies required field presence (`customer.name`), and enforces enum status constraints (`IN_PROGRESS` or `COMPLETED`).
Assessment: PASS
Severity: LOW
Notes: Correctly models consumer-side parsing strictness.

## Finding 4

Location: `internal/provider/server.go:11-131`
Claimed Behavior: Three provider implementations representing baseline V1, breaking change scenario, and dual-routing backwards-compatible evolution.
Observed Implementation:
- `ProviderV1`: returns baseline order matching V1 contract.
- `ProviderBreaking`: modifies status casing (`in_progress`), renames `name` to `full_name`, and changes `total` to string.
- `ProviderDual`: serves V1 format on `/v1/orders/` and V2 format on `/v2/orders/`.
Assessment: PASS
Severity: LOW
Notes: Clear separation of provider behaviors matching architectural test cases.
