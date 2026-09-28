# Implementation Notes

## Files Added

- `go.mod`: Module definition (`labs/31-oauth2-and-oidc`) targeting Go 1.22+.
- `pkg/pkce/pkce.go`: PKCE verifier and challenge generation (`S256` and `plain`) and verification per RFC 7636.
- `pkg/oidc/oidc.go`: ID Token JWT generation, HMAC-SHA256 signature verification, and standard claims validation (`iss`, `aud`, `exp`, `nonce`) per OpenID Connect Core 1.0 Section 3.1.3.7.
- `pkg/server/server.go`: In-memory Authorization Server supporting authorization code grant, PKCE validation, ID Token minting, Bearer Access Token validation, and Refresh Token rotation with token family replay detection (RFC 9700 Section 4.14).
- `pkg/client/client.go`: Client helper orchestrating PKCE generation, authorization request building, code exchange, ID token validation, and token refresh.
- `cmd/demo/main.go`: End-to-end runnable demo demonstrating the full authentication and authorization lifecycle, PKCE attack mitigation, and refresh token replay defense.
- `tests/oauth_test.go`: Unit tests and concurrent race tests covering PKCE, ID Token validation, code flow, and refresh token rotation with race condition checks.

## Core Design Decisions

1. **Pure Go Standard Library**:
   - Implemented JWT HS256 signing and base64url encoding using standard library (`crypto/hmac`, `crypto/sha256`, `encoding/base64`, `encoding/json`) without external third-party dependencies.
2. **Token Family Tracking for Replay Detection**:
   - Implemented RFC 9700 Section 4.14 refresh token rotation with family identifiers (`FamilyID`). When an already-consumed refresh token is replayed, the server flags the family as compromised and revokes all active tokens in that family.
3. **Mandatory PKCE Verification**:
   - Authorization requests require `code_challenge` and `code_challenge_method`. Token exchange enforces cryptographic proof-of-possession before issuing tokens.

## Implementation-Specific Choices

- HMAC-SHA256 (symmetric) used for ID Token signing in the standalone lab environment instead of RSA/JWKS network endpoints to keep the lab self-contained and reproducible.
- Thread safety enforced via mutexes across the in-memory Authorization Server.

## Known Limitations

- In-memory store: Tokens and codes are stored in memory and reset on process restart.
- No UI/HTML login prompt: Demonstrates the underlying protocol and token exchange programmatically.

## Trade-offs

- Symmetric JWT (HS256) vs Asymmetric (RS256/JWKS): Symmetric was selected to avoid external crypto key generation / JWKS web server complexity while preserving 100% cryptographic validation logic.

## What Is Demonstrated

- Authorization Code Flow with PKCE (`S256`).
- Separation of concerns: Access Token (opaque bearer for resource access) vs. ID Token (JWT with identity claims).
- ID Token cryptographic signature and claims verification (`iss`, `aud`, `exp`, `nonce`).
- Protection against authorization code interception attacks.
- Refresh Token Rotation and automatic Token Family Revocation on replay.
- Concurrent execution without race conditions.

## What Is Not Demonstrated

- Deprecated Implicit Flow or Resource Owner Password Credentials Flow.
- Complex distributed distributed JWKS key rotation policies.
