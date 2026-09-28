# Sources

## Source 1
Title: The OAuth 2.0 Authorization Framework
Publisher: IETF (Internet Engineering Task Force)
URL: https://datatracker.ietf.org/doc/html/rfc6749
Published: October 2012
Accessed: 2026-09-28
Source Tier: Tier 1
Relevance: Defines roles, authorization code flow, implicit flow, access tokens, and refresh tokens.

## Source 2
Title: OpenID Connect Core 1.0 incorporating errata set 2
Publisher: OpenID Foundation
URL: https://openid.net/specs/openid-connect-core-1_0.html
Published: December 15, 2023
Accessed: 2026-09-28
Source Tier: Tier 1
Relevance: Defines OIDC identity layer on top of OAuth 2.0, ID Token format (JWT), claims (`iss`, `sub`, `aud`, `exp`, `iat`, `nonce`), validation rules, and UserInfo endpoint.

## Source 3
Title: Proof Key for Code Exchange by OAuth Public Clients (PKCE)
Publisher: IETF (Internet Engineering Task Force)
URL: https://datatracker.ietf.org/doc/html/rfc7636
Published: September 2015
Accessed: 2026-09-28
Source Tier: Tier 1
Relevance: Defines `code_verifier`, `code_challenge`, transformation methods (`S256`, `plain`), and authorization code interception mitigation.

## Source 4
Title: OAuth 2.0 for Native Apps (BCP 212)
Publisher: IETF (Internet Engineering Task Force)
URL: https://datatracker.ietf.org/doc/html/rfc8252
Published: October 2017
Accessed: 2026-09-28
Source Tier: Tier 1
Relevance: Outlines use of external user-agents (browsers), deprecation of embedded web-views, loopback interface and custom URI redirection, and mandatory PKCE.

## Source 5
Title: Best Current Practice for OAuth 2.0 Security (BCP 240 / RFC 9700)
Publisher: IETF (Internet Engineering Task Force)
URL: https://datatracker.ietf.org/doc/html/rfc9700
Published: January 2025
Accessed: 2026-09-28
Source Tier: Tier 1
Relevance: Formally deprecates Implicit Flow and Resource Owner Password Credentials (ROPC); mandates PKCE for public clients and recommends it for confidential clients; specifies exact redirect URI matching, mix-up defenses, sender-constrained tokens, and refresh token rotation.

## Source 6
Title: JSON Web Token (JWT)
Publisher: IETF (Internet Engineering Task Force)
URL: https://datatracker.ietf.org/doc/html/rfc7519
Published: May 2015
Accessed: 2026-09-28
Source Tier: Tier 1
Relevance: Standardizes JWT claims representation, format, header (`typ`, `alg`), registered claims (`iss`, `sub`, `aud`, `exp`, `nbf`, `iat`, `jti`), and cryptographic verification requirements.
