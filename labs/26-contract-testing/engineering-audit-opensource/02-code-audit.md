# Code Audit

## Finding 1 — CDC verifier engine correctness
Location: `internal/contract/verifier.go:58` (`Verify`)
Claimed Behavior: Verifier executes consumer contract interactions against a running HTTP provider, checking status code, JSON object body, and a minimal-subset diff of expected vs. actual fields.
Observed Implementation: `Verify` builds a request per interaction, sets expected request headers, sends via shared `*http.Client`, checks status code equality, decodes the response body as a JSON object using `json.Decoder` with `UseNumber`, then runs `diffValues` to assert every consumer-expected field is present with matching type/value. Missing fields, type mismatches (`json.Number` vs `string`, primitive type changes), and value mismatches (enum casing) are all reported as errors.
Assessment: PASS
Severity: LOW
Notes: Body diff is intentionally a minimal-subset check (consumer-side contract): extra provider fields are ignored, which matches CDC semantics.

## Finding 2 — Breaking change detection (3 claimed cases)
Location: `internal/provider/server.go` (`ProviderBreaking.ServeHTTP`) and verifier `diffValues` (`verifier.go:164-173`)
Claimed Behavior: Three breaking changes are detectable: enum casing (`IN_PROGRESS` vs `in_progress`), field rename (`customer.name` -> `customer.full_name`), and primitive type mutation (`total` int -> string).
Observed Implementation: `ProviderBreaking` serves `status="in_progress"`, `customer.full_name`, and `total="150000"`. The demo output and `TestProviderBreaking_ContractVerification_Fails` both report exactly these three error classes. Verified at runtime (see Step 3 evidence):
```
1. [A request for order details by ID] missing expected field 'customer.name'
2. [A request for order details by ID] path 'total': type mismatch (expected 150000 [json.Number], got 150000 [string])
3. [A request for order details by ID] path 'status': value mismatch (expected "IN_PROGRESS", got "in_progress")
```
Assessment: PASS
Severity: LOW
Notes: None.

## Finding 3 — Safe API evolution (V1 + V2 dual provider)
Location: `internal/provider/server.go` (`ProviderDual.ServeHTTP`)
Claimed Behavior: Dual provider preserves the V1 contract while also exposing a V2 schema (new fields, renamed `full_name`, string `total`, `currency`).
Observed Implementation: `ProviderDual` serves `/v1/orders/*` with compliant `OrderResponseV1` and `/v2/orders/*` with `OrderResponseV2`. The V1 path is verified to satisfy the MobileApp contract; the V2 path adds fields without breaking the V1 contract. Tests confirm V1 path passes; demo confirms V1 verification passes.
Assessment: PASS
Severity: LOW
Notes: V2 path itself is not directly contract-tested, but V2 is not part of the MobileApp consumer contract (out of scope).

## Finding 4 — Concurrency safety of verifier
Location: `internal/contract/verifier.go:47-55` (`Verifier` struct + `NewVerifier`)
Claimed Behavior: Verifier is safe under concurrent parallel verification.
Observed Implementation: `Verifier` exposes a single `*http.Client` field. `http.Client` is documented as safe for concurrent use by multiple goroutines, and the per-call `VerificationResult`/`diffs` are stack locals with no shared mutable state. `diffValues` is a pure function.
Assessment: PASS
Severity: LOW
Notes: The client has no `Timeout` configured, but against `httptest` servers this is acceptable; not a concurrency hazard.

## Finding 5 — Response-header verification not implemented
Location: `internal/contract/verifier.go:74-111`
Claimed Behavior (structural): `ResponseDefinition.Headers` is modeled in the contract, implying response headers could be validated.
Observed Implementation: In `Verify`, the request-side `Request.Headers` are applied to the outbound request, but the response-side `ResponseDefinition.Headers` is never compared against the actual response headers. Only status and body are verified.
Assessment: WARNING
Severity: LOW
Notes: No test or README explicitly claims response-header validation, and the demo does not reference response headers. This is an uncovered-but-declared schema field; flag as a latent gap rather than a false claim.

## Finding 6 — Consumer client contract enforcement
Location: `internal/consumer/client.go:34` (`FetchOrder`)
Claimed Behavior: Consumer client enforces minimal schema (required `customer.name`, enumerated `status`).
Observed Implementation: `FetchOrder` validates `raw.Customer.Name != ""` and `raw.Status` in `{IN_PROGRESS, COMPLETED}`, returning a `contract violation` error otherwise. Against the breaking provider this yields an error (proven by `TestProviderBreaking_ContractVerification_Fails`).
Assessment: PASS
Severity: LOW
Notes: None.

## Finding 7 — Module is self-contained (stdlib only)
Location: `go.mod`, imports
Claimed Behavior: Lab builds with no hidden/fabricated dependencies.
Observed Implementation: `go.mod` declares `module labs/26-contract-testing` / `go 1.22` with no `require` block. All imports are `std lib` (`net/http`, `encoding/json`, `httptest`, etc.). `go build ./...`, `go vet ./...`, `go test -race -count=1 ./...` all succeed (see Step 3).
Assessment: PASS
Severity: LOW
Notes: None.
