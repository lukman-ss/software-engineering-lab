# Research Audit Plan

## Target Lab
`labs/31-oauth2-and-oidc`

## Audit Scope & Override
- Pipeline override active: Research audit only.
- Implementation/code/test audit excluded for this phase.
- No modifications to research files.

## Files Reviewed
- `labs/31-oauth2-and-oidc/research/01-plan.md`
- `labs/31-oauth2-and-oidc/research/02-sources.md`
- `labs/31-oauth2-and-oidc/research/03-evidence.md`
- `labs/31-oauth2-and-oidc/research/04-contradictions.md`
- `labs/31-oauth2-and-oidc/research/05-report.md`
- `labs/31-oauth2-and-oidc/research/06-open-questions.md`

## Claims To Verify
1. OAuth 2.0 is delegated authorization only, not authentication (RFC 6749, OIDC Core 1.0 Sec 1).
2. Access tokens are typically opaque to clients; JWT access tokens standardized by RFC 9068 (RFC 6749 Sec 1.4, RFC 9068).
3. OIDC introduces identity layer via signed ID Token JWT + UserInfo endpoint (OIDC Core 1.0 Sec 1, 2, 5.3).
4. Mandatory claims in ID Token: `iss`, `sub`, `aud` (RP client_id), `exp`, `iat` (OIDC Core 1.0 Sec 2).
5. 13-step ID Token validation requirements including issuer match, audience validation, signature verification, expiry, and nonce checking (OIDC Core 1.0 Sec 3.1.3.7).
6. Insecurity of using OAuth access tokens as proof of authentication (OIDC Core 1.0 Sec 1, RFC 8725, RFC 9700 Sec 2.3).
7. PKCE mechanism: `code_verifier` (43-128 chars, >=256-bit entropy), `code_challenge` via S256 (`BASE64URL(SHA256(verifier))`), S256 MTI, plain deprecated (RFC 7636 Sec 4.1/4.2/4.6/7.2, RFC 9700 Sec 2.1.1).
8. Deprecation of Implicit Grant; mandate for Authorization Code Flow + PKCE (RFC 9700 Sec 2.1.2, OAuth 2.1 draft).
9. Refresh Token rotation mandate for public clients, invalidating old refresh tokens upon use, and revoking active family upon reuse detection (RFC 9700 Sec 2.2.2/4.14.2).
10. JWT security best practices: algorithm pinning, rejecting `none` algorithm unless explicit, verifying all crypto layers, explicit typing (RFC 8725, OIDC Core Sec 2).
11. Browser-based application storage risks: avoiding localStorage for bearer tokens, BFF pattern with httpOnly cookies, constrained storage exceptions (RFC 9700 Sec 4.2/4.3, draft-ietf-oauth-browser-based-apps-27 / RFC 10017).

## Primary Risks
- Discrepancies between older foundation standards (RFC 6749, OIDC Core 1.0 2014/2023) and modern BCPs (RFC 9700 Jan 2025, OAuth 2.1 draft).
- Draft references: RFC 10017 citation in `02-sources.md` vs draft status of browser-based apps.
- Distinguishing universal RFC normative requirements vs deployment architecture recommendations (e.g., BFF vs in-memory vs DPoP).

## Audit Strategy
1. Inspect all 8 cited sources across IETF and OpenID Foundation standards.
2. Evaluate claim-to-source fidelity, ensuring normative RFC keywords (`MUST`, `SHOULD`, `MTI`) match citations.
3. Validate contradiction handling (historical OIDC implicit flow vs modern RFC 9700 BCP).
4. Record research gaps and assign quality gate status.
