# Code Audit

Target Lab: labs/26-contract-testing

## Finding 1

Location: `internal/contract/verifier.go:58-115`
Claimed Behavior: Verifier executes HTTP request against baseURL and compares returned payload against contract expectations, failing when contract expectations are violated.
Observed Implementation: Verifier handles request creation, status code validation, JSON decoding with `decoder.UseNumber()`, and field-by-field recursive diffing via `diffValues`.
Assessment: PASS
Severity: LOW
Notes: `io.ReadAll` and response body close handled properly.

## Finding 2

Location: `internal/contract/verifier.go:117-176`
Claimed Behavior: Detects structural missing keys, primitive type mismatches, and value mismatches while tolerating additive provider fields not declared in consumer contract.
Observed Implementation: `diffValues` checks key presence recursively. Extra fields on `actual` map not present in `expected` map are ignored, adhering to consumer-driven contract principles. Type differences and value differences are captured with descriptive path strings.
Assessment: PASS
Severity: LOW
Notes: Correctly utilizes `json.Number` comparisons to prevent int/float float64 coercion quirks.

## Finding 3

Location: `internal/consumer/client.go:34-79`
Claimed Behavior: Mobile client fetches order and fails when breaking provider schema (missing `customer.name`, unexpected status, type failure) is returned.
Observed Implementation: Strict unmarshalling and runtime validation of required contract constraints (`raw.Customer.Name == ""` and status enums) ensures consumer-side failure matches verifier failure.
Assessment: PASS
Severity: LOW
Notes: Client behaves identically to production consumer expectations.

## Finding 4

Location: `internal/provider/server.go:12-131`
Claimed Behavior: Provides three HTTP handlers representing V1 compliant provider, breaking change provider, and dual-version evolutionary provider.
Observed Implementation:
- `ProviderV1`: returns contract-matching JSON (`status: "IN_PROGRESS"`, `total: 150000`, `customer.name: "Budi Santoso"`).
- `ProviderBreaking`: introduces casing mismatch (`"in_progress"`), type mismatch (`"150000"` string), and field rename (`full_name`).
- `ProviderDual`: serves V1 contract compliant payload on `/v1/orders/` and new schema on `/v2/orders/`.
Assessment: PASS
Severity: LOW
Notes: Handlers are stateless, robust, and return standard Content-Type and HTTP status codes.

## Finding 5

Location: `tests/contract_test.go:96-114`
Claimed Behavior: Concurrent verification runs safely without race conditions.
Observed Implementation: Spawns 20 parallel goroutines invoking `verifier.Verify(srv.URL, c)` concurrently against `httptest.Server`. Verified with `go test -race ./...`.
Assessment: PASS
Severity: LOW
Notes: No shared mutable state in Verifier. Safe for concurrent CI runner simulations.
