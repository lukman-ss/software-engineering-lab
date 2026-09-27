# Engineering Code Audit

Target Lab: labs/26-contract-testing

## Finding 1

Location: internal/contract/verifier.go:58-115
Claimed Behavior: Verifier issues HTTP requests, validates status code, decodes JSON with `UseNumber()`, and performs recursive subset comparison against consumer expectations.
Observed Implementation: `Verify` checks status code, checks response body JSON formatting, and calls `diffValues` to verify that all fields expected by the consumer match both type and value.
Assessment: PASS
Severity: LOW
Notes: Correctly handles nested JSON objects, primitive mismatches, and numeric representations via `json.Number`.

## Finding 2

Location: internal/contract/verifier.go:117-176
Claimed Behavior: Verification diff detects missing expected fields, type mismatches, and value mismatches while ignoring extra fields produced by provider.
Observed Implementation: `diffValues` iterates over keys in expected map; extra keys present in `actMap` but absent from `expMap` are ignored, satisfying the consumer-driven contract principle.
Assessment: PASS
Severity: LOW
Notes: Correctly identifies case changes (`IN_PROGRESS` vs `in_progress`), type changes (`number` vs `string`), and missing objects (`customer.name`).

## Finding 3

Location: internal/provider/server.go:1-68
Claimed Behavior: HTTP server handlers provide V1 baseline, Breaking provider, and backwards-compatible Dual provider implementations.
Observed Implementation: Clean standard library `net/http` implementations. Uses `http.Handler` routing for `/v1/orders/` and `/v2/orders/`.
Assessment: PASS
Severity: LOW
Notes: No unneeded dependencies; conforms strictly to Go standard library conventions.

## Finding 4

Location: internal/consumer/client.go:1-73
Claimed Behavior: Mobile consumer defines contract schema and provides client implementation consuming provider responses.
Observed Implementation: `BuildContract` explicitly declares required fields (`id`, `status`, `total`, `customer.name`). `GetOrder` parses response into `MobileOrderDetail`.
Assessment: PASS
Severity: LOW
Notes: Correct implementation demonstrating consumer-driven schema definition.
