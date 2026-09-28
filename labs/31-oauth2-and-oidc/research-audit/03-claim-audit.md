# Claim Audit: OAuth 2.0 & OIDC Research

Target Lab: `labs/31-oauth2-and-oidc`
Audit Scope: Major Research Claims

---

## Claim 1

Claim:
OAuth 2.0 is purely an Authorization (Delegation) protocol issuing `access_token` for resource access, not an Authentication protocol. OpenID Connect is an identity layer on top of OAuth 2.0 that adds authentication via ID Token.

Location:
`05-report.md` (Finding 1) & `03-evidence.md` (Evidence 1 & 2)

Evidence Provided:
Cites RFC 6749 Section 1.1 & Section 1.4 defining token purpose; OIDC Core 1.0 Section 1 & Section 2 defining ID Token and claims.

Source:
RFC 6749, OpenID Connect Core 1.0

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW (Fully verified)

Notes:
Accurate distinction according to official specifications.

---

## Claim 2

Claim:
Using an OAuth 2.0 Access Token for authentication is insecure because access tokens do not contain identity assertions for the client and are subject to token leakage/replay vulnerabilities when issued via front-channel flows.

Location:
`05-report.md` (Finding 1 & Finding 4) & `03-evidence.md` (Evidence 3)

Evidence Provided:
Cites RFC 6749 Section 1.4 (access tokens are authorization credentials, opaque to client) and RFC 9700 Section 2.1.2 (implicit grant / front-channel access tokens vulnerable to leakage and replay).

Source:
RFC 6749 Section 1.4, RFC 9700 Section 2.1.2

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW (Fully verified)

Notes:
Accurately reflects security threats documented in RFC 6749 and RFC 9700.

---

## Claim 3

Claim:
Authorization Code Flow with PKCE (RFC 7636) prevents authorization code interception and injection attacks; PKCE is mandatory for public clients and recommended for confidential clients under RFC 9700; S256 is Mandatory-To-Implement (MTI).

Location:
`05-report.md` (Finding 2) & `03-evidence.md` (Evidence 4)

Evidence Provided:
Cites RFC 7636 Section 1, Section 4.2 (S256 MTI), Section 7.1/7.2; RFC 9700 Section 2.1.1.1 (MUST for public, RECOMMENDED for confidential).

Source:
RFC 7636, RFC 9700

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW (Fully verified)

Notes:
Quotes and normative requirements accurately match RFC text.

---

## Claim 4

Claim:
Implicit Flow (response type `token`) is formally deprecated across modern standards due to access token leakage and replay risks.

Location:
`05-report.md` (Finding 2, Executive Summary) & `03-evidence.md` (Evidence 3 & 4)

Evidence Provided:
Cites RFC 9700 Section 2.1.2.

Source:
RFC 9700 Section 2.1.2

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW (Fully verified)

Notes:
RFC 9700 explicitly establishes BCP 240 guidance deprecating implicit grant.

---

## Claim 5

Claim:
ID Tokens MUST be validated for signature against the issuer's keys (JWKS), issuer (`iss`), audience (`aud` containing client_id), expiration time (`exp`), and matching `nonce` if present.

Location:
`05-report.md` (Finding 3) & `03-evidence.md` (Evidence 5)

Evidence Provided:
Cites OpenID Connect Core 1.0 Section 2 and Section 3.1.3.7; RFC 7519 Section 4.1.

Source:
OpenID Connect Core 1.0 Section 3.1.3.7, RFC 7519 Section 4.1

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW (Fully verified)

Notes:
Step-by-step validation matches Section 3.1.3.7 of OIDC Core 1.0.

---

## Claim 6

Claim:
Storing tokens in browser `localStorage` is vulnerable to XSS; tokens should be stored using secure server-side/BFF patterns or HttpOnly cookies.

Location:
`05-report.md` (Finding 4, Limitations) & `03-evidence.md` (Evidence 6) & `04-contradictions.md` (Nuance 3)

Evidence Provided:
Cites RFC 6819 Section 10.3, RFC 9700 Section 2.1.2, oauth.net implicit flow docs.

Source:
RFC 6819 Section 10.3, RFC 9700 Section 2.1.2, RFC 9700 Section 2.6.5

Source Actually Supports Claim:
PARTIAL (Direct RFCs mandate bearer token protection; specific browser mechanisms like "HttpOnly SameSite cookies" / "BFF" are industry architectural best practices derived from RFC threat models).

Classification:
INTERPRETATION

Severity:
LOW (The research explicitly identifies this as an interpretation in `04-contradictions.md` and `05-report.md` under Limitations).

Notes:
The research author properly scoped and qualified this claim without misrepresenting it as a direct RFC normative RFC MUST keyword.

---

## Claim 7

Claim:
Refresh tokens for public clients MUST be sender-constrained or use refresh token rotation.

Location:
`05-report.md` (Finding 4) & `03-evidence.md` (Evidence 7)

Evidence Provided:
Cites RFC 9700 Section 2.2.2 and Section 4.14.

Source:
RFC 9700 Section 2.2.2, Section 4.14

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW (Fully verified)

Notes:
Exact normative requirement from RFC 9700 (BCP 240).
