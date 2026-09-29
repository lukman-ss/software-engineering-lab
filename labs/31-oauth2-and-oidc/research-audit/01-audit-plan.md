# Audit Plan — OAuth 2.0 & OIDC Research Audit

## Target Lab
`labs/31-oauth2-and-oidc`

## Audit Scope Notice
Per PIPELINE OVERRIDE instructions:
- Audit research only (`research/` files).
- Do not audit implementation/code in this stage.
- Do not modify research files.
- Output directory: `labs/31-oauth2-and-oidc/research-audit/`.

## Files Reviewed
- `research/01-plan.md`
- `research/02-sources.md`
- `research/03-evidence.md`
- `research/04-contradictions.md`
- `research/05-report.md`
- `research/06-open-questions.md`

## Claims To Verify
1. OAuth 2.0 is strictly an authorization (delegated access) framework, not an authentication protocol (RFC 6749).
2. Access tokens are typically opaque to clients and lack standardized authentication semantics (RFC 6749, RFC 9700).
3. OIDC extends OAuth 2.0 by providing an identity layer returning ID Tokens as JWTs with required claims `iss`, `sub`, `aud`, `exp`, `iat` (OpenID Connect Core 1.0, RFC 7519).
4. Relying Party must perform a 13-step verification on ID Tokens, including signature, issuer, audience (`aud` matching `client_id`), expiration, and nonce (OIDC Core 1.0 Sec 3.1.3.7, RFC 8725).
5. Using an OAuth 2.0 Access Token for user authentication / login is unsafe (OIDC Core Sec 1, RFC 8725).
6. PKCE (`S256`) is required for public clients and prevents authorization code interception attacks (RFC 7636, RFC 9700).
7. `S256` is mandatory-to-implement for PKCE, while `plain` is deprecated (RFC 7636).
8. Implicit flow is deprecated; Authorization Code Flow + PKCE is the modern standard (RFC 9700, OAuth 2.1).
9. Refresh tokens for public clients MUST be sender-constrained or rotated with automatic breach revocation (RFC 9700 Sec 4.14).
10. Tokens MUST NOT be stored in browser-accessible storage like `localStorage` due to XSS risks; BFF pattern is recommended (draft-ietf-oauth-browser-based-apps-27 / RFC 10017).

## Code To Execute
None. Code/implementation execution and auditing are explicitly excluded per PIPELINE OVERRIDE.

## Primary Risks
- Cited RFC URLs or draft URLs may be broken or mistyped.
- Citations from RFC 9700 or draft specs may over-attribute requirements or misquote dates.
- Internal inconsistencies between OIDC Core 1.0 (2014/2023) and RFC 9700 (2025) might be overlooked or misclassified.
- Draft RFC (e.g. `draft-ietf-oauth-browser-based-apps-27` / RFC 10017) status might be incorrectly cited as fully finalized RFC without noting publication status.

## Audit Strategy
1. Source verification: Audit all 8 sources cited in `02-sources.md` for validity, tier, reachability, and relevance.
2. Claim extraction & verification: Match all claims in `03-evidence.md` and `05-report.md` against standard specs (RFC 6749, 7636, 7519, 8725, 9700, OIDC Core 1.0, OAuth 2.1).
3. Contradiction audit: Review identified contradictions in `04-contradictions.md` for accuracy and completeness.
4. Gap analysis: Evaluate missing cases, overgeneralizations, and unverified assumptions across research output.
5. Final Verdict: Issue evidence-based final status (`APPROVED`, `APPROVED_WITH_WARNINGS`, `NEEDS_REVISION`, or `REJECTED`) in `07-verdict.md`.
