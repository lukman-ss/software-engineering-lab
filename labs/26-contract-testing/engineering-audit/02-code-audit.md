# Code Audit

Target Lab: labs/26-contract-testing

## Finding 1

Location: `internal/contract/verifier.go:58-115`
Claimed Behavior: Pure Go contract verifier testing HTTP interactions against baseURL with strict diff validation on status code, response headers, and response payload types/values.
Observed Implementation: Verifier handles request dispatch, status code checks, JSON decoding using `json.Number`, and deep diff recursion for map keys, types, and numbers.
Assessment: PASS
Severity: LOW
Notes: `diffValues` correctly handles nested object traversal, type assertion comparison, and number formatting checks.

## Finding 2

Location: `internal/model/order.go:13-53`
Claimed Behavior: Models support original schema (V1), breaking modification schema (enum casing, rename, type change), and evolutionary dual schema (V2).
Observed Implementation: Struct tags and fields explicitly define `OrderResponseV1`, `OrderResponseBreaking`, and `OrderResponseV2` aligning with failure scenario specifications.
Assessment: PASS
Severity: LOW
Notes: Clear separation of models cleanly models breaking changes.

## Finding 3

Location: `internal/consumer/client.go:34-110`
Claimed Behavior: Consumer implements client fetching order subset and generates minimal CDC contract specification.
Observed Implementation: `FetchOrder` unmarshals into minimal struct (`id`, `status`, `customer.name`, `total`) and verifies validation rules. `GenerateMobileContract` produces expected interaction JSON data structure.
Assessment: PASS
Severity: LOW
Notes: Demonstrates consumer-driven specification pattern accurately.

## Finding 4

Location: `internal/provider/server.go:11-131`
Claimed Behavior: Provider servers implement V1, Breaking, and Dual (V1 + V2) HTTP routes.
Observed Implementation: Handlers return respective model DTOs with JSON serialization and handle 404/405 conditions.
Assessment: PASS
Severity: LOW
Notes: Handlers are stateless and safe for concurrent requests.

## Finding 5

Location: `internal/contract/verifier.go:47-55`
Claimed Behavior: Concurrency safety during verification runs.
Observed Implementation: `Verifier` uses `http.Client` which is goroutine-safe. `Verify` creates new requests and local verification state per invocation without shared mutable state.
Assessment: PASS
Severity: LOW
Notes: Concurrency test confirms race-free execution under `go test -race`.
