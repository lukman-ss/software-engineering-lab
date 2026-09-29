# Sources

## Source 1
Title: RFC 6749 — The OAuth 2.0 Authorization Framework
Publisher: IETF
URL: https://www.rfc-editor.org/rfc/rfc6749.html
Published: October 2012
Accessed: 2026-09-29
Source Tier: Tier 1 (RFC Standard)
Relevance: Defines OAuth 2.0 as authorization framework; roles (resource owner, client, authorization server, resource server); access tokens, refresh tokens, authorization code grant, implicit grant

## Source 2
Title: RFC 7636 — Proof Key for Code Exchange by OAuth Public Clients (PKCE)
Publisher: IETF
URL: https://www.rfc-editor.org/rfc/rfc7636.html
Published: September 2015
Accessed: 2026-09-29
Source Tier: Tier 1 (RFC Standard)
Relevance: PKCE mechanism (code_verifier, code_challenge, S256); protects against authorization code interception

## Source 3
Title: OpenID Connect Core 1.0 incorporating errata set 2
Publisher: OpenID Foundation
URL: https://openid.net/specs/openid-connect-core-1_0.html
Published: December 2023
Accessed: 2026-09-29
Source Tier: Tier 1 (Specification Standard)
Relevance: OIDC as identity layer on top of OAuth 2.0; ID Token definition; required claims (iss, sub, aud, exp, iat); 13-step ID Token validation; Authorization Code Flow steps; UserInfo Endpoint

## Source 4
Title: RFC 7519 — JSON Web Token (JWT)
Publisher: IETF
URL: https://www.rfc-editor.org/rfc/rfc7519.html
Published: May 2015
Accessed: 2026-09-29
Source Tier: Tier 1 (RFC Standard)
Relevance: JWT claims (iss, sub, aud, exp, nbf, iat, jti); JWS/JWE structure; compact serialization format

## Source 5
Title: RFC 8725 — JSON Web Token Best Current Practices
Publisher: IETF
URL: https://www.rfc-editor.org/rfc/rfc8725.html
Published: February 2020
Accessed: 2026-09-29
Source Tier: Tier 1 (BCP)
Relevance: Algorithm verification; validate issuer/subject/audience; cross-JWT confusion prevention; use explicit typing

## Source 6
Title: RFC 9700 — Best Current Practice for OAuth 2.0 Security
Publisher: IETF
URL: https://www.rfc-editor.org/rfc/rfc9700.html
Published: January 2025
Accessed: 2026-09-29
Source Tier: Tier 1 (BCP)
Relevance: Implicit flow deprecation; PKCE mandatory for all public clients; refresh token rotation for public clients; exact redirect URI matching; token privilege restriction (audience restriction)

## Source 7
Title: OAuth 2.1 (draft-ietf-oauth-v2-1)
Publisher: IETF OAuth Working Group
URL: https://datatracker.ietf.org/doc/draft-ietf-oauth-v2-1/
Accessed: 2026-09-29
Source Tier: Tier 1 (draft spec)
Relevance: Consolidation of RFCs; PKCE required for all clients; implicit grant omitted; ROPC omitted; refresh token rotation for public clients (see also community summary at https://oauth.net/2.1/)

## Source 8
Title: RFC 10017 / draft-ietf-oauth-browser-based-apps-27 — OAuth 2.0 for Browser-Based Applications
Publisher: IETF
URL: https://datatracker.ietf.org/doc/html/draft-ietf-oauth-browser-based-apps-27
Published: July 2026 (draft)
Accessed: 2026-09-29
Source Tier: Tier 1 (BCP draft)
Relevance: BFF pattern; token storage in browser; XSS threats; malicious JS scenarios; BFF as confidential client hiding tokens from browser
