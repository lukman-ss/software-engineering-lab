# Source Audit — OAuth 2.0 & OIDC Research

Target Lab: `labs/31-oauth2-and-oidc`

## Source 1
Claimed Title: RFC 6749 — The OAuth 2.0 Authorization Framework  
Claimed Publisher: IETF  
URL: https://www.rfc-editor.org/rfc/rfc6749.html  

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Canonical authoritative standard for OAuth 2.0.

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
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Standard specification defining code_verifier, code_challenge, and S256 method.

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
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Canonical authoritative standard for OpenID Connect Core 1.0.

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
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Definitive specification for JWT claim representation and formatting.

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
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Official BCP document detailing cryptographic validation, algorithm confusion pitfalls, and claim checking.

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
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Official security BCP published January 2025 (obsoleting RFC 6819), establishing mandatory PKCE, implicit deprecation, and refresh token rotation.

Assessment:
PASS

---

## Source 7
Claimed Title: OAuth 2.1 (draft summary)  
Claimed Publisher: OAuth Working Group / oauth.net  
URL: https://oauth.net/2.1/  

Reachable:
YES

Source Type:
SECONDARY (Informational community summary of ongoing draft specification)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Marked as Tier 1 in `02-sources.md`, but `oauth.net/2.1/` is a community reference/educational summary maintained by Aaron Parecki / Okta / community rather than the direct datatracker IETF draft (`draft-ietf-oauth-v2-1`). Classification is more accurately SECONDARY, though accurate in content.

Assessment:
WARNING

---

## Source 8
Claimed Title: RFC 10017 / draft-ietf-oauth-browser-based-apps-27 — OAuth 2.0 for Browser-Based Applications  
Claimed Publisher: IETF  
URL: https://datatracker.ietf.org/doc/html/draft-ietf-oauth-browser-based-apps-27  

Reachable:
YES

Source Type:
PRIMARY (IETF Working Group Draft)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Title mixes "RFC 10017" with draft name "draft-ietf-oauth-browser-based-apps-27". In IETF tracking, this was an active internet-draft progressing towards RFC status; citing prospective RFC numbers before publication should be carefully scoped as draft status. The research explicitly noted this limitation in `05-report.md` Section "Limitations", mitigating the severity.

Assessment:
PASS
