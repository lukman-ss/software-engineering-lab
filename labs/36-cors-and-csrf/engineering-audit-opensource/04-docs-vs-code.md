# Docs vs Code Audit

## 1. README vs Implementation

### Claim 1: Architecture Structure
- **README**: Mentions `internal/cors`, `internal/csrf`, `internal/bank`, `cmd/demo`, and `tests`.
- **Code**: All five modules exist and correspond directly to the claimed packages and files.
- **Assessment**: PASS

### Claim 2: CORS Middleware Spec Compliance
- **README**: Claims `OPTIONS` preflight, allowed origins, method/header safelists, and credential checks.
- **Code**: `internal/cors/middleware.go` implements all of these behaviors including disallowing wildcard `*` with credentials.
- **Assessment**: PASS

### Claim 3: Anti-CSRF Defense Mechanisms
- **README**: Mentions HMAC-SHA256 signed session-bound tokens, Fetch Metadata (`Sec-Fetch-Site`), and API custom header middleware.
- **Code**: All three mechanisms are fully implemented in `internal/csrf/`.
- **Assessment**: PASS

### Claim 4: Test Commands
- **README**: Lists `go test -v ./...` and `go test -race ./...`.
- **Execution**: Both commands run and pass without errors.
- **Assessment**: PASS

### Claim 5: Demo Command
- **README**: Lists `go run ./cmd/demo`.
- **Execution**: Runs cleanly and completes with exit code 0.
- **Assessment**: PASS

---

## 2. Research Claims vs Implementation

### Research Claim 1: CORS Does Not Prevent CSRF
- **Research**: Concludes CORS is an origin-based access control policy enforced by browsers to protect responses, NOT a server-side firewall preventing request processing.
- **Code**: Proved in `tests/integration_test.go:54` and in demo step 2. Request from `https://evil.com` executes state mutation on the server regardless of CORS configuration.
- **Assessment**: PASS

### Research Claim 2: Signed Session-Bound CSRF Token Defeats Simple Requests
- **Research**: Advocates for session-bound cryptographically signed CSRF tokens (HMAC-SHA256).
- **Code**: `internal/csrf/token.go` implements `TokenManager` which signs `sessionID:ts:nonce` with HMAC-SHA256 and validates with `subtle.ConstantTimeCompare`.
- **Assessment**: PASS

### Research Claim 3: Modern Defense-in-Depth (Sec-Fetch-Site & Custom Headers)
- **Research**: Recommends `Sec-Fetch-Site` inspection and custom headers for API endpoints.
- **Code**: `FetchMetadataMiddleware` and `RequireCustomHeaderMiddleware` implement these exact strategies.
- **Assessment**: PASS

---

## Identified Discrepancies

- No `DOC_CODE_MISMATCH` detected.
- No `TEST_CLAIM_MISMATCH` detected.
- No `RESEARCH_IMPLEMENTATION_MISMATCH` detected.
