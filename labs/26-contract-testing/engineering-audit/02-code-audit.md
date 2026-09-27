# Code Audit

Target Lab: `labs/26-contract-testing`

## Finding 1

Location: `internal/contract/verifier.go:diffValues`
Claimed Behavior: Detects breaking schema changes (missing fields, enum/value differences, data type mismatches) when comparing actual JSON responses against consumer expectations.
Observed Implementation: Decodes JSON using `json.Number` via `decoder.UseNumber()` to prevent silent float conversions and compares nested object trees recursively. Detects type mismatches across strings vs numbers vs objects, missing fields, and exact value/enum differences.
Assessment: PASS
Severity: LOW
Notes: Correctly handles Go's `json.Number` abstraction to differentiate between numeric representations and JSON strings or object keys.

## Finding 2

Location: `internal/contract/verifier.go:Verify`
Claimed Behavior: Verifies HTTP interactions (method, path, headers, status code, response body subset) without mutating provider state or failing silently.
Observed Implementation: Makes HTTP calls to the target baseURL + interaction path, forwards request headers, checks status code, and validates expected response body as a consumer subset. Response body streams are properly closed using `defer` or explicit close after reading.
Assessment: PASS
Severity: LOW
Notes: Minimal consumer subset matching ensures extra provider fields (like `notes` in `ProviderV1`) do not fail contract verification, adhering strictly to Postel's Law / Consumer-Driven Contracts.

## Finding 3

Location: `internal/provider/server.go:ProviderBreaking`
Claimed Behavior: Simulates unannounced breaking changes in provider API.
Observed Implementation: Provides endpoints introducing lowercase status enum (`in_progress`), renamed customer field (`full_name` instead of `name`), and stringified total (`"150000"` instead of integer `150000`).
Assessment: PASS
Severity: LOW
Notes: Conforms directly to the failure scenarios documented in design notes.

## Finding 4

Location: `internal/provider/server.go:ProviderDual`
Claimed Behavior: Supports safe API evolution by providing both backward-compatible V1 endpoints and newly formatted V2 endpoints.
Observed Implementation: Correctly serves `/v1/orders/{id}` using `OrderResponseV1` and `/v2/orders/{id}` using `OrderResponseV2`.
Assessment: PASS
Severity: LOW
Notes: Demonstrates safe migration path in CI gates without breaking existing consumer contracts.

## Finding 5

Location: `internal/consumer/client.go:FetchOrder`
Claimed Behavior: Consumer client parses and validates required fields, failing when contracts are violated.
Observed Implementation: Deserializes response into consumer DTO, explicitly validating that `customer.name` is present and `status` conforms to expected uppercase enum values.
Assessment: PASS
Severity: LOW
Notes: Error propagation is clean with wrapped context.
