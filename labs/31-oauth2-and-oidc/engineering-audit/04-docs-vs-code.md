# Docs vs Code Audit

Target Lab: labs/31-oauth2-and-oidc

## Documentation Comparison

### 1. Structure and Package Layout
- **Claimed in README.md**:
  - `pkg/pkce`: PKCE code verifier and S256 challenge generation and verification.
  - `pkg/oidc`: ID Token creation, HMAC-SHA256 signature verification, and standard claims validation.
  - `pkg/server`: In-memory Authorization Server supporting Auth Code grant, PKCE, OIDC ID Tokens, and Refresh Token Rotation.
  - `pkg/client`: Client helper orchestrating requests and token parsing.
  - `cmd/demo`: Executable walkthrough of legitimate flows and attack defenses.
  - `tests`: Unit tests and race detector tests.
- **Observed in Implementation**: All packages and paths exist and match documented structure exactly.

### 2. Standards and Specifications
- **Claimed in README.md**:
  - OAuth 2.0 delegated authorization grants.
  - OIDC ID Tokens (signed JWTs) with `iss`, `sub`, `aud`, `exp`, `nonce`.
  - PKCE (RFC 7636 / RFC 9700) with `S256`.
  - Refresh Token Rotation (RFC 9700 Section 4.14) with token family revocation.
- **Observed in Code**:
  - `pkg/pkce/pkce.go` implements RFC 7636 S256/plain generation & verification.
  - `pkg/oidc/oidc.go` implements standard JWT signing and claim validation.
  - `pkg/server/server.go` implements family tracking and invalidation upon replay.

### 3. Execution Commands
- **Claimed in README.md**:
  - `go test -v ./...`
  - `go test -race ./...`
  - `go run ./cmd/demo`
- **Observed**: All commands execute successfully and produce the expected test passes and demonstration steps.

## Mismatch Detection

- `DOC_CODE_MISMATCH`: None detected.
- `TEST_CLAIM_MISMATCH`: None detected.
- `RESEARCH_IMPLEMENTATION_MISMATCH`: None detected.
