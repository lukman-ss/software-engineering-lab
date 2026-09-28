# Source Audit: OAuth 2.0 & OIDC Research

Target Lab: `labs/31-oauth2-and-oidc`
Audit Scope: Research Sources

---

## Source 1

Claimed Title: The OAuth 2.0 Authorization Framework
Claimed Publisher: IETF (Internet Engineering Task Force)
URL: https://datatracker.ietf.org/doc/html/rfc6749

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Standard baseline specification defining OAuth 2.0 core roles, grant flows, access tokens, and refresh tokens.

Assessment:
PASS

---

## Source 2

Claimed Title: OpenID Connect Core 1.0 incorporating errata set 2
Claimed Publisher: OpenID Foundation
URL: https://openid.net/specs/openid-connect-core-1_0.html

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Official specification defining identity layer, ID Token claims (`iss`, `sub`, `aud`, `exp`, `iat`, `nonce`), validation rules (Sec 3.1.3.7), and UserInfo endpoint.

Assessment:
PASS

---

## Source 3

Claimed Title: Proof Key for Code Exchange by OAuth Public Clients (PKCE)
Claimed Publisher: IETF (Internet Engineering Task Force)
URL: https://datatracker.ietf.org/doc/html/rfc7636

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Standard specification for PKCE, defining `code_verifier`, `code_challenge`, `S256` MTI requirements, and mitigation against authorization code interception attacks.

Assessment:
PASS

---

## Source 4

Claimed Title: OAuth 2.0 for Native Apps (BCP 212)
Claimed Publisher: IETF (Internet Engineering Task Force)
URL: https://datatracker.ietf.org/doc/html/rfc8252

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Authoritative Best Current Practice detailing external browser usage, web-view deprecation, and mandatory PKCE for native apps.

Assessment:
PASS

---

## Source 5

Claimed Title: Best Current Practice for OAuth 2.0 Security (BCP 240 / RFC 9700)
Claimed Publisher: IETF (Internet Engineering Task Force)
URL: https://datatracker.ietf.org/doc/html/rfc9700

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Published January 2025 as RFC 9700 / BCP 240, superseding older drafts and formalizing deprecation of Implicit Flow and ROPC, mandating PKCE for public clients and refresh token rotation/sender-constraining.

Assessment:
PASS

---

## Source 6

Claimed Title: JSON Web Token (JWT)
Claimed Publisher: IETF (Internet Engineering Task Force)
URL: https://datatracker.ietf.org/doc/html/rfc7519

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Authoritative standard for JWT structure and registered claim definitions (`iss`, `sub`, `aud`, `exp`, `iat`, `nbf`, `jti`).

Assessment:
PASS

---

## Secondary Cited Sources (in Evidence)

Claimed Title: OAuth 2.0 Threat Model and Security Considerations (RFC 6819) / OAuth.net Implicit Grant Docs
Claimed Publisher: IETF / OAuth Community
URL: https://datatracker.ietf.org/doc/html/rfc6819#section-10.3, https://oauth.net/2/grant-types/implicit/

Reachable:
YES

Source Type:
PRIMARY (RFC 6819) / COMMUNITY (OAuth.net)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Accurately labeled and used to corroborate threat model rationale.

Assessment:
PASS
