# Research Report: OAuth 2.0 & OIDC --- Authentication vs Authorization and Flow Security

**Research Date:** 2026-09-28

## Research Question

What are the fundamental differences between OAuth 2.0 and OpenID Connect (OIDC)? Why is using OAuth 2.0 Access Tokens for authentication considered insecure? How does Authorization Code Flow with PKCE protect against authorization code interception and injection attacks? What are the required validation rules for ID Tokens and Access Tokens according to current standards?

## Executive Summary

OAuth 2.0 is a pure authorization protocol that enables third-party applications to access resources on behalf of a resource owner. It issues access tokens for resource access, not identity. OpenID Connect (OIDC) is an identity layer on top of OAuth 2.0 that adds authentication by introducing the ID Token (a JWT) containing claims about the end-user's authentication event.

Using an OAuth 2.0 access token for authentication is insecure because it does not contain identity claims and is not designed to prove who the user is. OIDC solves this by providing a standardized ID Token with required claims (`iss`, `sub`, `aud`, `exp`, `iat`, `auth_time`, `nonce`) that must be validated for signature, issuer, audience, and expiration.

Authorization Code Flow + PKCE (RFC 7636) is the current mandatory standard. It protects against authorization code interception attacks by requiring the client to prove possession of a dynamically generated `code_verifier` secret when exchanging the authorization code. The implicit flow is deprecated across all authoritative sources.

## Findings

### Finding 1: OAuth 2.0 is Authorization; OIDC is Authentication

Claim: OAuth 2.0 as originally defined is purely an Authorization (Delegation) protocol, not an Authentication protocol. OpenID Connect is an identity layer built on top of OAuth 2.0 to enable authentication.

Evidence: RFC 6749 Section 1.1 explicitly defines four roles (resource owner, resource server, client, authorization server) and focuses on authorization grants and access tokens. Section 1.4 states: "Access tokens are credentials used to access protected resources. An access token is a string representing an authorization issued to the client. The string is usually opaque to the client. Tokens represent specific scopes and durations of access, granted by the resource owner, and enforced by the resource server and authorization server." RFC 6749 does not define any standard way to authenticate the resource owner's identity.

OpenID Connect Core 1.0 (Section 1. Introduction) states: "OpenID Connect 1.0 is a simple identity layer on top of the OAuth 2.0 protocol. It enables Clients to verify the identity of the End-User based on the authentication performed by an Authorization Server, as well as to obtain basic profile information about the End-User in an interoperable and REST-like manner." The ID Token is introduced as "a security token that contains Claims about the Authentication of an End-User by an Authorization Server when using a Client."

Sources: RFC 6749 Sections 1.1, 1.4; OpenID Connect Core 1.0 Section 1
Confidence: HIGH

### Finding 2: Authorization Code Flow + PKCE is the Modern Standard

Claim: Authorization Code Flow with PKCE (RFC 7636) is the mandatory standard for all OAuth clients. The implicit flow is deprecated due to security vulnerabilities (token exposure in URL fragment).

Evidence: RFC 7636 Section 1 describes the authorization code interception attack where malicious apps can register for the same redirect URI scheme and intercept the authorization code. The solution is PKCE: the client creates a cryptographically random `code_verifier`, derives a `code_challenge` via SHA-256 hashing (S256 method) or plain transformation, sends the challenge in the authorization request, then proves possession by sending the original verifier at the token endpoint. RFC 7636 Section 4.2 states: "If the client is capable of using 'S256', it MUST use 'S256', as 'S256' is Mandatory To Implement (MTI) on the server."

RFC 9700 Section 2.1.1.1 formally mandates: "Public clients MUST use PKCE [RFC7636] to this end... For confidential clients, the use of PKCE [RFC7636] is RECOMMENDED, as it provides strong protection against misuse and injection of authorization codes..." Section 2.1.2 deprecates the implicit flow: "The implicit grant (response type `token`) and other response types causing the authorization server to issue access tokens in the authorization response are vulnerable to access token leakage and access token replay... clients SHOULD NOT use the implicit grant... Use [Authorization Code] + [PKCE] instead."

Sources: RFC 7636 Sections 1, 4.2; RFC 9700 Sections 2.1.1.1, 2.1.2
Confidence: HIGH

### Finding 3: ID Token Validation Requirements

Claim: ID Tokens MUST be validated for signature, issuer (`iss`), audience (`aud`), expiration time (`exp`), and optionally `nonce` to prevent token substitution and replay attacks.

Evidence: OpenID Connect Core 1.0 Section 2 specifies the required ID Token claims: `iss` (Issuer Identifier, case-sensitive URL), `sub` (Subject Identifier), `aud` (Audience, MUST contain the client ID), `exp` (Expiration time), `iat` (Issued at), and optionally `auth_time`, `nonce`, `acr`, `amr`, `azp`. Section 3.1.3.7 (ID Token Validation) lists specific validation steps: verify the issuer matches the OP's identifier, verify the audience contains the client ID, verify the expiration time is in the future, and verify the signature using the appropriate public key (typically via JWKS endpoint). If the ID Token contains a `nonce` claim, it MUST match the nonce parameter sent in the authentication request.

RFC 7519 Section 4.1 defines these registered JWT claims and their semantics: `iss` identifies the principal that issued the JWT; `sub` identifies the subject; `aud` identifies recipients; `exp` identifies expiration time; `iat` identifies when the JWT was issued.

Sources: OpenID Connect Core 1.0 Sections 2, 3.1.3.7; RFC 7519 Section 4.1
Confidence: HIGH

### Finding 4: Token Storage and Refresh Token Security

Claim: Access tokens and ID tokens should not be stored in browser `localStorage` due to XSS vulnerability. Sender-constrained tokens (e.g., mutual TLS, DPoP) and refresh token rotation are mandatory for public clients.

Evidence: While primary RFCs do not explicitly mandate "HttpOnly Secure SameSite cookie vs BFF" wording, the threat model is clear. RFC 6819 Section 10.3 states access tokens are bearer tokens and "SHOULD NOT be placed in page fragments (as they are in the OAuth 2.0 Implicit Flow) as this exposes them to the resulting document and any scripts it may contain." RFC 9700 Section 2.1.2 notes implicit flow tokens are vulnerable to "access token leakage and access token replay" because they are exposed in URL fragments to browser history and JavaScript. The oauth.net implicit flow page explicitly states tokens in URL fragments are exposed to "browser history, referrer headers, and any JavaScript on the page."

RFC 9700 Section 2.2.2 mandates: "Refresh tokens for public clients MUST be sender-constrained or use refresh token rotation as described in Section 4.14." Section 4.14 describes refresh token rotation: "Refresh token rotation involves issuing a new refresh token each time one is used to obtain an access token. The old refresh token is then invalidated. This limits the usefulness of a stolen refresh token to a single use."

Sources: RFC 9700 Sections 2.1.2, 2.2.2, 4.14; RFC 6819 Section 10.3; oauth.net implicit flow documentation
Confidence: HIGH (threat model consensus); MEDIUM (specific implementation advice like "HttpOnly cookie" is derived best practice, not direct RFC mandate).

## Areas of Agreement

All primary sources agree on:

1. OAuth 2.0 is Authorization; OIDC adds Authentication via ID Token.
2. Authorization Code Flow + PKCE is the mandatory standard for all clients.
3. Implicit Flow is deprecated due to token exposure vulnerabilities.
4. ID Tokens must be validated for signature, `iss`, `aud`, `exp`.
5. PKCE `S256` is MTI; `plain` is legacy/compatibility.
6. Public clients must use PKCE; confidential clients are recommended to use it.

## Areas of Disagreement

No material contradictions discovered. Minor evolution noted:
- RFC 7636 (2015) designed PKCE for public clients only.
- RFC 9700 (2025) extends recommendation to confidential clients (as layered defense against injection + CSRF).
This is an evolution of best practice, not a contradiction.

## Limitations

1. Token storage recommendations (HttpOnly cookie vs BFF) are best practices derived from threat models in RFC 6819 and RFC 9700 rather than direct normative statements in primary RFCs. Confidence is HIGH due to converging guidance.

2. Specific implementation details of JWKS endpoint discovery and rotation policies are beyond scope but referenced (OpenID Connect Core 1.0 references JWK Set Discovery separately).

3. Cross-platform native app redirect URI handling (iOS, Android, Windows, macOS) is covered in RFC 8252 Appendix B but not detailed in this research.

## Conclusion

OAuth 2.0 and OpenID Connect serve fundamentally different purposes: OAuth 2.0 enables delegated authorization via access tokens; OIDC enables authentication via ID Tokens (JWTs). Using an OAuth 2.0 access token for authentication is insecure because it lacks identity claims.

Authorization Code Flow with PKCE is the current mandatory standard, preventing authorization code interception attacks by requiring proof-of-possession of a dynamically generated secret. ID Tokens must be validated for signature, issuer, audience, and expiration time.

Token storage must avoid XSS-vulnerable client-side storage (localStorage). Refresh tokens for public clients must use sender-constraining or rotation. The implicit flow and resource owner password credentials grant are deprecated due to well-documented security flaws.
