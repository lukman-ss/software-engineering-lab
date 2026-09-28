# Code Audit

Target Lab: labs/26-contract-testing

## Finding 1

Location: internal/contract/verifier.go:46-115
Claimed Behavior: Verifier executes HTTP requests against target provider and performs contract assertion.
Observed Implementation:
- Uses `http.Client{}` to issue requests defined in `Interaction.Request`.
- Properly closes response body (`_ = resp.Body.Close()`).
- Uses `decoder.UseNumber()` to prevent float64 coercion issues on numbers.
- Handles HTTP errors, status code mismatches, invalid JSON payloads, and value/type differences.
Assessment: PASS
Severity: LOW
Notes: Robust implementation using standard library HTTP and JSON tools.

## Finding 2

Location: internal/contract/verifier.go:117-176
Claimed Behavior: Recursive diff algorithm detects missing fields, value mismatches, and primitive type mismatches without requiring full schema parity (supports subset matching / Postel's Law).
Observed Implementation:
- Map traversal verifies expected fields exist in actual map while ignoring extra provider fields (consumer subset verification).
- Compares `json.Number` values strictly and flags type mismatch if one side is a string or non-number.
- Correctly reports nested field paths (e.g. `customer.name`).
Assessment: PASS
Severity: LOW
Notes: Matches consumer-driven contract design principles where consumer dictates required subset.

## Finding 3

Location: internal/consumer/client.go:21-79
Claimed Behavior: Mobile consumer client parses expected subset and fails on contract violations.
Observed Implementation:
- Enforces strict parsing and validates required fields (`customer.name`) and enum bounds (`IN_PROGRESS` or `COMPLETED`).
- Directly returns typed `MobileOrderSummary`.
Assessment: PASS
Severity: LOW
Notes: Validates consumer behavior when interacting with compliant vs breaking payloads.

## Finding 4

Location: internal/provider/server.go:11-131
Claimed Behavior: Providers implement compliant V1, breaking schema modifications, and dual V1+V2 backward-compatible handlers.
Observed Implementation:
- `ProviderV1` outputs canonical V1 JSON format.
- `ProviderBreaking` outputs lowercase enum (`in_progress`), renamed field (`full_name`), and stringified total (`"150000"`).
- `ProviderDual` multiplexes `/v1/orders/` and `/v2/orders/` paths seamlessly.
Assessment: PASS
Severity: LOW
Notes: Concrete and clean demonstration of breaking mutations versus evolutionary versioning.
