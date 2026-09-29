# Docs vs Code Audit

## Document Inventory

- `README.md`: Overview, Structure, Running Tests and Demo
- `engineering/01-design.md`: Architecture, data structures, workflows
- `engineering/02-implementation-notes.md`: Implementation decisions and notes
- `engineering/03-execution-result.md`: Recorded execution output

---

## Comparison Matrix

| Claim in Docs | Location in Docs | Code Implementation | Status |
|---|---|---|---|
| Access Tokens represent delegated authorization grants | README.md:7 | `pkg/server/server.go:160-165`, `tokens map` | MATCH |
| ID Tokens are signed JWTs with `iss`, `sub`, `aud`, `exp`, `nonce` | README.md:8 | `pkg/oidc/oidc.go:29-38` | MATCH |
| PKCE (RFC 7636 / RFC 9700) S256 protects against code interception | README.md:9 | `pkg/pkce/pkce.go:53-56`, `pkg/server/server.go:149-151` | MATCH |
| Refresh Token Rotation (RFC 9700 §4.14) with family revocation | README.md:10 | `pkg/server/server.go:229-238` | MATCH |
| Structure: `pkg/pkce`, `pkg/oidc`, `pkg/server`, `pkg/client`, `cmd/demo`, `tests` | README.md:14-19 | Files exist as described | MATCH |
| Demo outputs 8 steps verifying flow and attack blocks | `cmd/demo/main.go` vs `engineering/03-execution-result.md` | Actual execution matches recorded output | MATCH |

---

## Discrepancies Found

### 1. DOC_CODE_MISMATCH: Incomplete State CSRF handling
- **Claim**: Client stores `c.State` in `pkg/client/client.go:47`.
- **Reality**: `c.State` is never passed to `server.Authorize` or checked anywhere. The README does not claim state/CSRF validation, but internal design notes might imply full OAuth 2.0 redirect security.
- **Classification**: Minor. README does not claim CSRF state protection; scope is explicitly Authentication vs Authorization, PKCE, ID Token validation, and Refresh Token Rotation.

### 2. TEST_CLAIM_MISMATCH: None
- All test claims in README and implementation notes correspond to real test cases in `tests/oauth_test.go`.

### 3. RESEARCH_IMPLEMENTATION_MISMATCH: None
- PKCE S256, ID Token HMAC-SHA256, Refresh Token Family Revocation match RFC 7636, OIDC Core 1.0, and RFC 9700 §4.14 as specified in `research/`.

---

## Demo Verification

Actual demo output matches recorded output in `engineering/03-execution-result.md`:
- PKCE Challenge generated correctly
- Auth code issued with 5-minute expiry
- Access token, Refresh token, ID token issued upon exchange
- ID token claims (iss, sub, aud, nonce) validated successfully
- Resource Server accepts valid access token
- Interception attack (wrong verifier) blocked with `pkce verification failed: code_verifier does not match code_challenge`
- Refresh token rotated to new tokens
- Attacker replaying original refresh token blocked with `refresh token reuse detected: entire token family revoked`
- Active session using rotated token immediately revoked due to family invalidation

No fake demo output exists. Output is generated live by `go run ./cmd/demo`.
