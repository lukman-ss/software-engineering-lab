# Source Audit

Target Lab: `labs/31-oauth2-and-oidc`
Scope: Verification of cited sources in `research/02-sources.md`.

---

## Source 1

Claimed Title: RFC 6749 — The OAuth 2.0 Authorization Framework
Claimed Publisher: IETF
URL: https://www.rfc-editor.org/rfc/rfc6749.html

Reachable:
YES

Source Type:
PRIMARY (IETF Standards Track RFC)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Accurately defines OAuth 2.0 authorization roles, protocol flows, and token concepts.

Assessment:
PASS

---

## Source 2

Claimed Title: RFC 7636 — Proof Key for Code Exchange by OAuth Public Clients (PKCE)
Claimed Publisher: IETF
URL: https://www.rfc-editor.org/rfc/rfc7636.html

Reachable:
YES

Source Type:
PRIMARY (IETF Standards Track RFC)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Accurately defines `code_verifier`, `code_challenge`, S256/plain transformation methods, and Appendix B test vectors.

Assessment:
PASS

---

## Source 3

Claimed Title: OpenID Connect Core 1.0 incorporating errata set 2
Claimed Publisher: OpenID Foundation
URL: https://openid.net/specs/openid-connect-core-1_0.html

Reachable:
YES

Source Type:
PRIMARY (OpenID Foundation Specification Standard)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Defines ID Token JWT structure, required claims (`iss`, `sub`, `aud`, `exp`, `iat`), and Section 3.1.3.7 13-step ID Token validation.

Assessment:
PASS

---

## Source 4

Claimed Title: RFC 7519 — JSON Web Token (JWT)
Claimed Publisher: IETF
URL: https://www.rfc-editor.org/rfc/rfc7519.html

Reachable:
YES

Source Type:
PRIMARY (IETF Standards Track RFC)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Defines standard JWT claims, serialization formats, and base validation rules.

Assessment:
PASS

---

## Source 5

Claimed Title: RFC 8725 — JSON Web Token Best Current Practices
Claimed Publisher: IETF
URL: https://www.rfc-editor.org/rfc/rfc8725.html

Reachable:
YES

Source Type:
PRIMARY (IETF Best Current Practice BCP 225)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Outlines JWT vulnerabilities (alg:none, RS256/HS256 confusion, substitution attacks, cross-JWT confusion) and mitigations (algorithm verification, explicit typing, audience/issuer validation).

Assessment:
PASS

---

## Source 6

Claimed Title: RFC 9700 — Best Current Practice for OAuth 2.0 Security
Claimed Publisher: IETF
URL: https://www.rfc-editor.org/rfc/rfc9700.html

Reachable:
YES

Source Type:
PRIMARY (IETF Best Current Practice BCP 240, Published Jan 2025)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Mandates PKCE for public clients, deprecates Implicit Grant, requires refresh token rotation/sender-constraining for public clients, and exact redirect URI matching.

Assessment:
PASS

---

## Source 7

Claimed Title: OAuth 2.1 (draft-ietf-oauth-v2-1)
Claimed Publisher: IETF OAuth Working Group
URL: https://datatracker.ietf.org/doc/draft-ietf-oauth-v2-1/

Reachable:
YES

Source Type:
PRIMARY (IETF Active Internet-Draft)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Document is currently an active Internet-Draft (`draft-ietf-oauth-v2-1-16`), not yet a published RFC. The research correctly identifies it as a draft specification and does not attribute false RFC status to it.

Assessment:
PASS

---

## Source 8

Claimed Title: RFC 10017 / draft-ietf-oauth-browser-based-apps-27 — OAuth 2.0 for Browser-Based Applications
Claimed Publisher: IETF
URL: https://datatracker.ietf.org/doc/html/draft-ietf-oauth-browser-based-apps-27

Reachable:
YES

Source Type:
PRIMARY (IETF Internet-Draft / Informational / BCP)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Source entry lists "RFC 10017 / draft-ietf-oauth-browser-based-apps-27". The URL points to `draft-ietf-oauth-browser-based-apps-27`. Datatracker shows draft-27 was published as RFC 10017. The content in draft-27 accurately covers BFF patterns, malicious JS threat models, and token storage constraints.

Assessment:
PASS
