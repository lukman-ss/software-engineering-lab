# Research Plan: OAuth 2.0 & OIDC --- Authentication vs Authorization and Flow Security

## Research Topic
OAuth 2.0 & OpenID Connect (OIDC), Authorization Code Flow + PKCE, Token Types (`access_token`, `id_token`, `refresh_token`), and Security Best Practices.

## Objective
Investigate the architectural and security differences between OAuth 2.0 (Authorization) and OIDC (Authentication), evaluate the necessity of Authorization Code Flow + PKCE over deprecated flows (such as Implicit Flow), and establish authoritative evidence regarding token validation (`iss`, `aud`, `exp`, signature/JWKS) and token storage/handling.

## Research Questions
1. What are the fundamental differences between OAuth 2.0 (`access_token`) and OpenID Connect (`id_token`)?
2. Why is using an OAuth 2.0 Access Token for authentication considered an anti-pattern or security risk?
3. How does Authorization Code Flow with PKCE (RFC 7636) mitigate authorization code interception and injection attacks for public and confidential clients?
4. What are the validation rules and security considerations for ID Tokens and Access Tokens according to RFC 7519, OpenID Connect Core 1.0, and RFC 9700?
5. What are the current security best practices regarding token storage (e.g., avoiding `localStorage` due to XSS) and refresh token rotation?

## Search Strategy
- Query IETF standards (RFC 6749, RFC 7636, RFC 8252, RFC 9700, RFC 7519) and OpenID Foundation specifications (OpenID Connect Core 1.0).
- Cross-reference findings across primary and secondary authoritative sources.

## Expected Primary Sources
- RFC 6749: The OAuth 2.0 Authorization Framework
- OpenID Connect Core 1.0 incorporating errata set 2
- RFC 7636: Proof Key for Code Exchange by OAuth Public Clients (PKCE)
- RFC 8252: OAuth 2.0 for Native Apps
- RFC 9700: Best Current Practice for OAuth 2.0 Security
- RFC 7519: JSON Web Token (JWT)

## Risks / Unknowns
- Potential variations in library implementations of PKCE (`plain` vs `S256`).
- Evolution of browser cookie security policies (SameSite, Secure, HttpOnly) and BFF (Backend-for-Frontend) architecture patterns.
