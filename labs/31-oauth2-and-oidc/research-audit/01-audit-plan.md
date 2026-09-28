# Audit Plan: OAuth 2.0 & OIDC Research

Target Lab: `labs/31-oauth2-and-oidc`
Audit Scope: Research artifacts only (Pipeline Override active)
Audit Date: 2026-09-28

## Target Lab
`labs/31-oauth2-and-oidc`

## Files Reviewed
- `labs/31-oauth2-and-oidc/research/01-plan.md`
- `labs/31-oauth2-and-oidc/research/02-sources.md`
- `labs/31-oauth2-and-oidc/research/03-evidence.md`
- `labs/31-oauth2-and-oidc/research/04-contradictions.md`
- `labs/31-oauth2-and-oidc/research/05-report.md`
- `labs/31-oauth2-and-oidc/research/06-open-questions.md`

## Claims To Verify
1. OAuth 2.0 is an authorization protocol issuing `access_token` for resource access, not an authentication protocol asserting user identity.
2. OpenID Connect (OIDC) layers authentication onto OAuth 2.0 via the ID Token data structure formatted as JWT.
3. Using an OAuth 2.0 `access_token` for authentication is insecure and considered an anti-pattern.
4. Authorization Code Flow with PKCE (RFC 7636) mitigates authorization code interception and injection attacks; `S256` is Mandatory-To-Implement (MTI).
5. RFC 9700 mandates PKCE for public clients and recommends it for confidential clients while formally deprecating Implicit Flow and ROPC.
6. ID Tokens require validation of signature, issuer (`iss`), audience (`aud`), expiration (`exp`), and nonce if present.
7. Token storage recommendations against `localStorage` due to XSS, and refresh token rotation requirements for public clients under RFC 9700 Sec 2.2.2 and Sec 4.14.

## Code To Execute
- None: Pipeline override specifies "Audit research only. Do not audit implementation/code in this stage." No runnable codebase or tests exist yet in this lab phase.

## Primary Risks
1. Over-generalization of browser storage guidance (e.g. claiming RFCs directly mandate HttpOnly cookies or BFF when RFCs focus on bearer token threats).
2. Unverified or unreachable RFC URLs / specification anchors.
3. Verification of whether RFC 9700 is accurately cited as published January 2025 (BCP 240).

## Audit Strategy
1. Validate external sources against authoritative standards registries (IETF Datatracker, OpenID Foundation).
2. Audit each major claim across findings and evidence files against primary RFC citations.
3. Check internal consistency between research plan, sources, evidence, contradictions, and final report.
4. Document gaps, over-generalizations, and issue final evidence-based verdict.
