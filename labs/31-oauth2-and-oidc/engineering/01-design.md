# Engineering Design

Target Lab: `labs/31-oauth2-and-oidc`
Research Status: APPROVED (Audit Date: 2026-09-28)

## Concept To Prove

1. **OAuth 2.0 (Authorization) vs. OIDC (Authentication)**:
   - OAuth 2.0 Access Token provides delegated authorization to protected resources (scopes, opaque/bearer token).
   - OIDC ID Token (signed JWT) provides verified identity claims (`iss`, `sub`, `aud`, `exp`, `iat`, `nonce`) about the authenticated user.
2. **Authorization Code Flow with PKCE (RFC 7636 / RFC 9700)**:
   - Dynamic `code_verifier` and SHA-256 `code_challenge` (S256).
   - Prevents authorization code interception attack.
3. **ID Token Verification (OIDC Core 1.0 Section 3.1.3.7)**:
   - Strict cryptographic signature verification (HMAC-SHA256 / RSA).
   - Claims validation: `iss` match, `aud` match, `exp` validity, `nonce` match.
4. **Refresh Token Rotation (RFC 9700 Section 4.14)**:
   - Old refresh tokens invalidated upon single use.
   - Token reuse detection revokes entire token family.

## Expected Behavior

- **Authorization Request**: Client sends `response_type=code`, `client_id`, `redirect_uri`, `scope="openid profile"`, `state`, `code_challenge`, `code_challenge_method=S256`, `nonce`.
- **Authorization Code Issuance**: Authorization Server (AS) authenticates resource owner, stores challenge and nonce bound to the one-time code.
- **Token Exchange**: Client sends `grant_type=authorization_code`, `code`, `redirect_uri`, `client_id`, and plain `code_verifier`.
- **Token Response**: AS validates `code_verifier` hash against stored `code_challenge`, returning an Access Token, ID Token (JWT), and Refresh Token.
- **Protected Resource Access**: Resource server accepts Access Token and checks scopes. Rejects using Access Token directly as an identity assertion.
- **Token Refresh**: Using a Refresh Token returns a new Access Token and a rotated Refresh Token; attempting to reuse the previous Refresh Token fails and triggers family revocation.

## Failure Scenario

1. **PKCE Mismatch / Interception**: An attacker intercepting the authorization code cannot exchange it without the original `code_verifier`.
2. **Forged / Tampered ID Token**: Modifying payload claims or tampering with signature causes verification rejection.
3. **Expired ID Token / Wrong Audience / Wrong Issuer**: ID Token validation rejects tokens with past `exp`, mismatched `aud`, or invalid `iss`.
4. **Refresh Token Replay**: Re-using an already consumed refresh token triggers immediate rejection and invalidates active session tokens.

## Success Criteria

- 100% Go standard library implementation (no external heavyweight dependencies; crypto/hmac/sha256/jwt built cleanly with standard library).
- Unit tests cover authorization code flow, PKCE S256 verification, plain rejection/support, ID token claims validation, and refresh token rotation.
- Concurrency test passes with Go race detector (`go test -race ./...`).
- CLI demo runs end-to-end showing authorization, authentication, verification, interception failure, and refresh token rotation.

## Architecture

```text
+-------------+
|   Client    |----+ (1) Auth Request (code_challenge, nonce)
| Application |    |
+-------------+    v
      |         +-----------------------+
      |         | Authorization Server  | (Handles AS + OIDC Provider)
      |         +-----------------------+
      |            | (2) Auth Code
      |<-----------+
      |
      | (3) Token Exchange (code + code_verifier)
      |------------------->+-----------------------+
      |<-------------------| AS validates PKCE     |
      | (4) Tokens         | Issues:               |
      |  - Access Token    | - Access Token        |
      |  - ID Token (JWT)  | - ID Token (signed)   |
      |  - Refresh Token   | - Refresh Token       |
      |                    +-----------------------+
      v
+-------------+
|  Resource   | (5) Request with Bearer Access Token
|   Server    |------------------------------------> [Returns Protected Resource]
+-------------+
```

## Components

1. `pkg/pkce`: PKCE code verifier and S256/plain challenge generation and validation.
2. `pkg/oidc`: ID Token structure, JWT signing/verification, and strict OIDC claims validator (`iss`, `aud`, `exp`, `nonce`).
3. `pkg/server`: In-memory Authorization Server supporting Auth Code grant, PKCE validation, ID Token generation, and Refresh Token rotation with replay detection.
4. `pkg/client`: Client helper demonstrating correct flow parameters, verifier storage, and token validation.
5. `cmd/demo`: Interactive/automated terminal walkthrough executing valid flows and attack failure scenarios.

## Test Strategy

- `pkce_test.go`: PKCE S256 vector verification, invalid verifier length/charset, challenge computation.
- `oidc_test.go`: JWT signature verification, expired token rejection, wrong audience/issuer/nonce detection.
- `server_test.go`: End-to-end authorization code grant, code one-time use, PKCE interception rejection, refresh token rotation, refresh replay revocation.
- `race_test.go`: Concurrent token requests and refresh operations verifying thread safety.

## Execution Plan

1. Create `go.mod` (Go 1.22+).
2. Implement `pkg/pkce`, `pkg/oidc`, `pkg/server`, and `pkg/client`.
3. Implement `tests/` and run `go test -v -race ./...`.
4. Implement `cmd/demo/main.go` and run demo.
5. Record execution results in `engineering/03-execution-result.md` and notes in `engineering/02-implementation-notes.md`.

## Implementation Decisions

- Standard library `crypto/hmac`, `crypto/sha256`, and `encoding/json` used for zero-dependency JWT (HS256) implementation, keeping the lab lightweight and self-contained.
- Refresh tokens track token families with lineage IDs to demonstrate RFC 9700 Section 4.14 token reuse detection.
