# Content Brief

**Topic:** OAuth 2.0 & OpenID Connect (OIDC) — Authentication vs Authorization, PKCE Protection, ID Token Validation, Refresh Token Rotation.

**Target Reader:** Software engineer or security practitioner who needs to understand how OAuth 2.0 and OIDC fit together in a production system, why the standards enforce specific defenses, and how the lab verifies each mechanism.

**Problem:** OAuth 2.0 provides delegated authorization (access tokens), but it does not define user identity. OIDC adds an identity layer on top of OAuth 2.0 through signed ID Tokens, while PKCE prevents authorization-code interception and refresh-token rotation detects replayed stolen tokens. The common failure mode is confusing authorization with authentication, relying on deprecated flows (implicit flow), skipping ID Token verification, or reusing refresh tokens.

**Core Mental Model:**
- OAuth 2.0 = "What can this client access?" → access token with scopes.
- OIDC = "Who is the user?" → ID Token (JWT) with signed identity claims.
- PKCE = proof-of-possession binding the authorization code to the original client request.
- Refresh Token Rotation = single-use tokens with family revocation for replay detection.

**Approved Research Status:** APPROVED (Research Audit Verdict 2026-09-29; major claims reviewed: 12; contradictions analyzed: 4; gaps: 3 LOW).

**Approved Engineering Status:** APPROVED (Engineering Audit Verdict 2026-09-29; tests passed: 17/17; race detector: PASS; demo: PASS; warnings: 3 LOW severity).

**Main Concepts:**
1. OAuth 2.0 Authorization vs OIDC Authentication.
2. Authorization Code Flow + PKCE (RFC 7636 / RFC 9700).
3. ID Token JWT structure and multi-step validation.
4. Refresh Token Rotation and Token Family Revocation (RFC 9700 Section 4.14).
5. Common pitfalls: implicit flow, missing signature/audience checks, local storage XSS, static refresh tokens.

**Verified Behaviors:**
- PKCE S256 challenge generation, verifier length bounds, mismatch rejection.
- ID Token HMAC-SHA256 signature verification, issuer/audience/expiration/nonce validation.
- Full authorization code grant flow with client validation.
- Refresh token rotation with single-use invalidation.
- Replay detection that revokes the entire token family.
- Concurrent request handling without race conditions.

**Available Case Studies:**
- Demo walkthrough `cmd/demo/main.go` showing legitimate flows and attack failures (interception, replay detection, family revocation).

**Warnings:**
- Implementation uses in-memory storage (tokens/codes reset on restart).
- Uses symmetric HMAC-SHA256 instead of RSA/JWKS for lab simplicity.
- PKCE comparison uses string equality (not constant-time) per audit LOW warning.
- Implicit flow is deprecated; this lab intentionally excludes it.
