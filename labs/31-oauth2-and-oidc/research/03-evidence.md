# Evidence

## Evidence 1
Claim: OAuth 2.0 is an authorization (delegated access) framework, not an authentication protocol.
Evidence: "The OAuth 2.0 authorization framework enables a third-party application to obtain limited access to an HTTP service, either on behalf of a resource owner... or by allowing the third-party application to obtain access on its own behalf."
Source: RFC 6749, Abstract + Section 1
URL: https://www.rfc-editor.org/rfc/rfc6749.html
Confidence: HIGH
Corroborated By: OIDC Core 1.0 Introduction ("without profiling OAuth 2.0, it is incapable of providing information about the authentication of an End-User"), RFC 9700
Notes: Foundation for AuthN vs AuthZ distinction.

## Evidence 2
Claim: OAuth 2.0 abstract flow issues an access token for protected-resource access; access tokens are typically opaque to the client.
Evidence: Steps (A)-(F) protocol flow (authorization grant -> access token -> protected resource); "An access token is a string representing an authorization issued to the client. The string is usually opaque to the client."
Source: RFC 6749, Sections 1.2, 1.4
URL: https://www.rfc-editor.org/rfc/rfc6749.html
Confidence: HIGH
Corroborated By: RFC 9700 Sec 2.3 (audience/scope restriction of access tokens)
Notes: Supports claim that access token answers "what may be accessed". "Usually opaque" reflects the RFC 6749 default; RFC 9068 (October 2021, "JSON Web Token (JWT) Profile for OAuth 2.0 Access Tokens", https://www.rfc-editor.org/rfc/rfc9068.html) standardizes structured JWT access tokens widely deployed in modern systems. Even when structured, clients SHOULD treat access tokens as opaque unless the deployment profile explicitly requires introspection.

## Evidence 3
Claim: OIDC is an identity layer on top of OAuth 2.0; authentication is signaled by openid scope and returned as an ID Token JWT; profile via UserInfo Endpoint.
Evidence: "OpenID Connect 1.0 is a simple identity layer on top of the OAuth 2.0 protocol. It enables Clients to verify the identity of the End-User based on the authentication performed by an Authorization Server... Information about the authentication performed is returned in a JWT called an ID Token." Abstract 5-step flow: RP request -> OP authenticates -> OP responds with ID Token (+ usually Access Token) -> RP calls UserInfo -> UserInfo returns Claims.
Source: OpenID Connect Core 1.0, Sections 1, 1.3, 2, 5.3
URL: https://openid.net/specs/openid-connect-core-1_0.html
Confidence: HIGH
Corroborated By: RFC 9700 (references OpenID.Core), RFC 7519 (JWT structure of ID Token)
Notes: Directly supports "OAuth -> access_token (what), OIDC -> id_token (who)" framing.

## Evidence 4
Claim: ID Tokens MUST contain iss, sub, aud, exp, iat; aud must include the RP client_id.
Evidence: "iss REQUIRED... sub REQUIRED... aud REQUIRED... MUST contain the OAuth 2.0 client_id of the Relying Party... exp REQUIRED... iat REQUIRED."
Source: OIDC Core 1.0, Section 2
URL: https://openid.net/specs/openid-connect-core-1_0.html
Confidence: HIGH
Corroborated By: RFC 7519 Sec 4.1 (registered claims definitions), RFC 8725 Sec 3.8/3.9 (must validate iss/sub/aud)
Notes: Basis for signature + aud/iss/exp verification exercise.

## Evidence 5
Claim: RP MUST validate ID Token via multi-step check: issuer match, audience contains client_id, signature via issuer keys, exp, nonce replay check.
Evidence: 13-step list in Section 3.1.3.7: iss exact match; aud contains client_id, reject if not listed; validate signature per JWS using issuer keys; current time before exp; nonce claim must equal request nonce; optional acr/auth_time checks.
Source: OIDC Core 1.0, Section 3.1.3.7
URL: https://openid.net/specs/openid-connect-core-1_0.html#IDTokenValidation
Confidence: HIGH
Corroborated By: RFC 8725 Sec 3.1/3.3/3.8/3.9 (algorithm verification, validate all crypto ops, validate iss/sub/aud)
Notes: Retrieved via section-specific extraction from fetched full text.

## Evidence 6
Claim: Using an OAuth access token for login is unsafe because access tokens carry no standardized authentication semantics or audience binding to the RP.
Evidence: OIDC Core: OAuth defines "mechanisms to obtain and use Access Tokens... but do not define standard methods to provide identity information. Notably, without profiling OAuth 2.0, it is incapable of providing information about the authentication of an End-User." JWT BCP: substitution attacks (token intended for one recipient reused at another) and cross-JWT confusion require aud/iss validation and mutually-exclusive validation rules per JWT kind.
Source: OIDC Core 1.0 Sec 1; RFC 8725 Sec 2.7/2.8/3.9/3.12
URL: https://openid.net/specs/openid-connect-core-1_0.html / https://www.rfc-editor.org/rfc/rfc8725.html
Confidence: HIGH
Corroborated By: RFC 9700 Sec 2.3 (access tokens should be audience-restricted to resource servers, not RPs)
Notes: Interpretation (why access-token-login fails) built from two HIGH-confidence facts.

## Evidence 7
Claim: PKCE prevents authorization-code interception: client creates code_verifier (43-128 unreserved chars, >=256-bit entropy), sends code_challenge (S256 = BASE64URL(SHA256(verifier))); server binds challenge to code and verifies verifier at token endpoint.
Evidence: RFC 7636 Sec 1.1 flow (A-D), Sec 4.1/4.2/4.6 formulas; "If the values are equal, the token endpoint MUST continue... If not equal, invalid_grant."
Source: RFC 7636
URL: https://www.rfc-editor.org/rfc/rfc7636.html
Confidence: HIGH
Corroborated By: RFC 9700 Sec 2.1.1 + 4.5.3.1 (PKCE as code-injection countermeasure for all clients)
Notes: Includes RFC 7636 Appendix B worked example (verifier dBjftJeZ... -> challenge E9Melhoa...).

## Evidence 8
Claim: S256 is mandatory-to-implement; plain must only be fallback and must not be used in new implementations.
Evidence: "If the client is capable of using S256, it MUST use S256, as S256 is Mandatory To Implement on the server... plain SHOULD NOT be used and exists only for compatibility."
Source: RFC 7636, Sec 4.2, 7.2
URL: https://www.rfc-editor.org/rfc/rfc7636.html
Confidence: HIGH
Corroborated By: RFC 9700 Sec 2.1.1 ("clients SHOULD use PKCE methods that do not expose the verifier... Currently, S256 is the only such method")
Notes: Justifies lab exercise using SHA256 challenge.

## Evidence 9
Claim: Implicit flow (token in authorization response / URL fragment) is deprecated; clients should use response_type=code + PKCE instead.
Evidence: RFC 9700 Sec 2.1.2: "clients SHOULD NOT use the implicit grant (response type token)... Clients SHOULD instead use the response type code"; OAuth 2.1: "The Implicit grant (response_type=token) is omitted."
Source: RFC 9700; https://oauth.net/2.1/
URL: https://www.rfc-editor.org/rfc/rfc9700.html / https://oauth.net/2.1/
Confidence: HIGH
Corroborated By: Browser-based-apps draft Sec 7.2 (implicit grant threat analysis); OIDC Core flow table (implicit exposes tokens to user agent)
Notes: OIDC Core still documents implicit/hybrid historically — see contradictions file.

## Evidence 10
Claim: Refresh tokens for public clients MUST be sender-constrained or rotated; rotation invalidates the previous token so reuse signals breach and active token is revoked.
Evidence: RFC 9700 Sec 2.2.2 + Sec 4.14.2: "Refresh tokens for public clients MUST be sender-constrained or use refresh token rotation"; rotation text: "issues a new refresh token with every access token refresh response. The previous refresh token is invalidated... it will revoke the active refresh token."
Source: RFC 9700
URL: https://www.rfc-editor.org/rfc/rfc9700.html
Confidence: HIGH
Corroborated By: OAuth 2.1 summary (same requirement); OIDC Core Sec 12.2 example response containing a new refresh_token alongside access_token
Notes: OIDC Core Sec 12 itself does not use the word "rotation" — rotation mandate comes from RFC 9700.

## Evidence 11
Claim: JWT validation must pin the algorithm (never trust alg header blindly; reject none unless explicitly configured), validate all nested crypto operations, and use mutually-exclusive rules per JWT kind (typ/aud/iss/key separation).
Evidence: RFC 8725 Sec 3.1 ("MUST NOT use any other algorithms... each key MUST be used with exactly one algorithm"), Sec 3.3 (reject if any operation fails), Sec 3.12 + 3.11 (explicit typing, distinct aud/iss/keys per JWT kind); threat Sec 2.1 (alg=none swap, RS256->HS256 confusion).
Source: RFC 8725
URL: https://www.rfc-editor.org/rfc/rfc8725.html
Confidence: HIGH
Corroborated By: OIDC Core Sec 2 (ID Tokens MUST be JWS-signed; MUST NOT use none except code-flow + explicit registration)
Notes: Supports "verify signature + aud" trap in lab.

## Evidence 12
Claim: Tokens must not live in browser JS-reachable storage; recommended pattern is Backend-for-Frontend (BFF): confidential client holds tokens server-side in a cookie session and proxies resource requests.
Evidence: Browser-based-apps draft Sec 6.1: BFF "interacts with the authorization server as a confidential OAuth client; manages OAuth access and refresh tokens in the context of a cookie-based session, avoiding direct exposure of any tokens to the browser-based application; forwards all requests to a resource server, augmenting them with the correct access token." Sec 5: malicious JS has same privileges as app code; can steal from localStorage/IndexedDB.
Source: draft-ietf-oauth-browser-based-apps-27, Sec 5, 6.1, 8
URL: https://datatracker.ietf.org/doc/html/draft-ietf-oauth-browser-based-apps-27
Confidence: MEDIUM (authoritative draft, pre-RFC; corroborated by RFC 9700 token-leakage sections 4.2/4.3 on referer/history leakage)
Corroborated By: RFC 9700 Sec 4.2/4.3 (credential leakage via referer headers, browser history)
Notes: Draft status (expires Jan 2027); published as RFC 10017 — content may have shifted slightly.
