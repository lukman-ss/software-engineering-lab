# Docs vs Code Audit

## Comparisons

### 1. Architecture Alignment
- `README.md` claims:
  - `internal/cors`: Spec-compliant CORS middleware (`OPTIONS` preflight, allowed origins, method/header safelists, credential checks).
  - `internal/csrf`: Anti-CSRF mechanisms including HMAC-SHA256 signed session-bound tokens, Fetch Metadata (`Sec-Fetch-Site`), and API custom header middleware.
  - `internal/bank`: Bank application service simulating cookie-authenticated balance inquiries, vulnerable transfer endpoints, and protected transfer endpoints.
  - `cmd/demo`: Runnable CLI program showcasing attacks against vulnerable vs. protected configurations.
  - `tests`: Integration test suite verifying cross-origin requests, preflight, race safety, and attack mitigation.
- Code observation: All packages exist, contain the described modules, and implement the features as listed.

### 2. Execution Claims
- `README.md` claims:
  - `go test -v ./...` runs tests. (PASS: verified)
  - `go test -race ./...` runs tests with race detector. (PASS: verified)
  - `go run ./cmd/demo` executes runnable demonstration. (PASS: verified)

### 3. Demo Output Veracity
- Demo output produced from `go run ./cmd/demo`:
  - Step 1: Initial state ($1000 victim, $50 attacker).
  - Step 2: Cross-origin attack on vulnerable endpoint succeeds with balance deduction ($600 victim, $450 attacker) despite missing CORS header.
  - Step 3: Cross-origin attack on protected endpoint returns 403, balances unchanged.
  - Step 4: Legitimate flow with signed CSRF token succeeds ($500 victim, $550 attacker).
- Code verification: Output matches execution trace line-by-line; no mocked print statements masquerading as live handlers.

### 4. Mismatches
- None detected. Implementation matches documentation cleanly.
