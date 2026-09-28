# Code Audit

Target Lab: labs/26-contract-testing

## Finding 1

Location: `internal/contract/verifier.go:60-128`
Claimed Behavior: Verifier inspects method, path, headers, status codes, and performs subset JSON schema/field validation against provider endpoints.
Observed Implementation: `Verify` builds HTTP requests, validates response headers, status codes, and parses response JSON into recursive diff comparator `diffValues` (`decoder.UseNumber()` preserved for accurate type comparisons).
Assessment: PASS
Severity: LOW
Notes: Correctly handles nested paths, maps, and numbers. Correctly closes response body.

## Finding 2

Location: `internal/contract/verifier.go:130-189`
Claimed Behavior: Recursive comparison of JSON subset detects missing keys, primitive type mismatches, and value mismatches.
Observed Implementation: `diffValues` checks missing keys in target objects, distinguishes `json.Number` vs strings, and checks reflection types. Correctly implements consumer-driven contract subset semantics (extra provider fields are ignored; only fields declared in the contract are enforced).
Assessment: PASS
Severity: LOW
Notes: Accurately reports missing fields (`customer.name`), type mismatches (`total`: string vs number), and value/enum mismatches (`status`: `in_progress` vs `IN_PROGRESS`).

## Finding 3

Location: `internal/provider/server.go:11-44`, `46-80`, `82-131`
Claimed Behavior: Independent HTTP handlers for Provider V1 (compliant), Provider Breaking (non-compliant), and Provider Dual (backward-compatible V1 + new V2).
Observed Implementation: Distinct structs `ProviderV1`, `ProviderBreaking`, and `ProviderDual` implementing `http.Handler` via `ServeHTTP`. Method checking handles non-GET with 405. Route checking handles 404.
Assessment: PASS
Severity: LOW
Notes: Handlers are stateless, thread-safe, and self-contained.

## Finding 4

Location: `internal/consumer/client.go:37-82`, `84-114`
Claimed Behavior: Mobile consumer defines its required contract and implements client parsing enforcing its contract assertions.
Observed Implementation: `GenerateMobileContract()` builds contract data structure specifying required interactions. `FetchOrder()` handles HTTP response and explicitly verifies contract integrity on unmarshalled fields.
Assessment: PASS
Severity: LOW
Notes: Standard library HTTP client with 5-second timeout.

## Finding 5

Location: `internal/contract/verifier.go:52-58`
Claimed Behavior: Thread-safe verifier execution under concurrent loads.
Observed Implementation: `Verifier` uses standard `*http.Client` which is goroutine-safe. The `Verify` method creates local request and response structures per invocation without shared mutable state.
Assessment: PASS
Severity: LOW
Notes: Validated by `TestConcurrentContractVerification` with `go test -race ./...`.
