# Code Audit: Lab 26 Contract Testing

## Finding 1

Location: `internal/contract/verifier.go:58-115`
Claimed Behavior: Verifies HTTP interactions against running provider, flagging status code mismatches, invalid JSON, and missing/mismatched contract fields.
Observed Implementation: Correctly constructs requests, executes via `http.Client`, validates status code, reads response body, uses `json.NewDecoder.UseNumber()`, and calls recursive `diffValues` to detect schema differences.
Assessment: PASS
Severity: LOW
Notes: Properly closes response bodies and reports interaction failure details with interaction descriptions.

## Finding 2

Location: `internal/contract/verifier.go:117-176`
Claimed Behavior: Recursive comparison validates expected subset against actual payload, distinguishing between objects, numbers, and strings.
Observed Implementation: Implements recursive map traversal, checks presence of expected keys, verifies `json.Number` values and strings, and checks reflection types.
Assessment: PASS
Severity: LOW
Notes: Covers JSON objects and primitives cleanly. Does not currently handle JSON array/slice diffing, but sufficient for the order single-resource contract scope implemented.

## Finding 3

Location: `internal/consumer/client.go:82-111`
Claimed Behavior: Generates Consumer-Driven Contract specifying only fields Mobile consumer needs (`id`, `status`, `customer.name`, `total`).
Observed Implementation: Matches claimed contract generation structure. Omits provider-internal fields like `notes` and `customer.id`.
Assessment: PASS
Severity: LOW
Notes: Correctly enforces Postel's Law / Consumer-Driven Contracts principle.

## Finding 4

Location: `internal/provider/server.go:12-44`, `47-80`, `83-131`
Claimed Behavior: Implements three provider variants: Provider V1 (conforming), Provider Breaking (3 contract violations), Provider Dual (V1 backwards-compatible with V2 evolution).
Observed Implementation:
- Provider V1 serves conforming schema on `/v1/orders/`.
- Provider Breaking serves lower-cased status `in_progress`, renamed customer field `full_name`, and stringified total `"150000"` on `/v1/orders/`.
- Provider Dual routes `/v1/orders/` with V1 schema and `/v2/orders/` with V2 schema.
Assessment: PASS
Severity: LOW
Notes: Handlers are stateless, thread-safe, and handle unsupported methods/paths with appropriate HTTP error codes.

## Finding 5

Location: `internal/consumer/client.go:34-79`
Claimed Behavior: Consumer client parses and validates response, returning error if contract fields are missing or invalid.
Observed Implementation: Strictly unmarshals and validates `customer.name` non-empty and `status` in expected enum set (`IN_PROGRESS`, `COMPLETED`).
Assessment: PASS
Severity: LOW
Notes: Properly surfaces contract violations at client runtime when communicating with breaking provider.
