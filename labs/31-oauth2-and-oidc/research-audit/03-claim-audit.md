# Claim Audit — OAuth 2.0 & OIDC Research

Target Lab: `labs/31-oauth2-and-oidc`

---

## Claim 1

Claim:
OAuth 2.0 is an authorization (delegated access) framework, not an authentication protocol.

Location:
`research/03-evidence.md` §Evidence 1; `research/05-report.md` §Finding 1

Evidence Provided:
RFC 6749 Abstract + Section 1: "The OAuth 2.0 authorization framework enables a third-party application to obtain limited access to an HTTP service..."
OIDC Core Sec 1: "without profiling OAuth 2.0, it is incapable of providing information about the authentication of an End-User."

Source:
RFC 6749; OIDC Core 1.0 Sec 1

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Claim is grounded directly in the RFC 6749 Abstract and Section 1 definitions. OIDC Core Sec 1 corroborates it explicitly. No issues.

---

## Claim 2

Claim:
Access tokens are opaque to the client and do not carry standardized authentication semantics.

Location:
`research/03-evidence.md` §Evidence 2; `research/05-report.md` §Finding 1

Evidence Provided:
RFC 6749 Section 1.4: "An access token is a string representing an authorization issued to the client. The string is usually opaque to the client."

Source:
RFC 6749 Sec 1.4

Source Actually Supports Claim:
PARTIAL

Classification:
FACT

Severity:
LOW

Notes:
RFC 6749 says tokens are "usually opaque" — not always. The research overstates this slightly by not qualifying "usually". The core conclusion (access tokens lack standardized authentication semantics) is still strongly supported by OIDC Core Sec 1 corroboration. Not a material issue, but the word "usually" is present in the spec and should not be silently universalized.

---

## Claim 3

Claim:
OIDC adds an identity layer on top of OAuth 2.0; user identity is returned as an ID Token JWT; required claims: iss, sub, aud (must include client_id), exp, iat.

Location:
`research/03-evidence.md` §Evidence 3, 4; `research/05-report.md` §Finding 2

Evidence Provided:
OIDC Core Abstract/Sec 1.3; OIDC Core Sec 2 (required claims list with explicit "REQUIRED" language)

Source:
OIDC Core 1.0 Secs 1, 1.3, 2, 3, 5.3

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Directly supported by canonical specification. No issues.

---

## Claim 4

Claim:
Relying Party MUST validate ID Token via a 13-step procedure including: issuer exact match, audience contains `client_id`, JWS signature validation using issuer JWKS, expiration check, and nonce anti-replay.

Location:
`research/03-evidence.md` §Evidence 5; `research/05-report.md` §Finding 3

Evidence Provided:
OIDC Core Sec 3.1.3.7 "13-step list"; also RFC 8725 Sec 3.1/3.3/3.8/3.9

Source:
OIDC Core Sec 3.1.3.7; RFC 8725

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
OIDC Core Sec 3.1.3.7 is publicly available and well-known to enumerate these validation steps. The specific section URL `#IDTokenValidation` is cited in `03-evidence.md:43`, which is the correct anchor. No issues.

---

## Claim 5

Claim:
Using an OAuth 2.0 access token to authenticate a user (login) is unsafe because access tokens carry no standardized authentication semantics or audience binding to the Relying Party (RP).

Location:
`research/03-evidence.md` §Evidence 6; `research/05-report.md` §Finding 1, §Executive Summary

Evidence Provided:
OIDC Core Sec 1 (OAuth 2.0 cannot provide authentication info without profiling); RFC 8725 Sec 2.7/2.8/3.9/3.12 (token substitution attacks, cross-JWT confusion, audience binding)

Source:
OIDC Core Sec 1; RFC 8725

Source Actually Supports Claim:
YES

Classification:
FACT / INTERPRETATION

Severity:
LOW

Notes:
The underlying claim is strongly supported by OIDC Core's explicit statement and RFC 8725's threat analysis. The research correctly identifies this as an "interpretation built from two HIGH-confidence facts" (`03-evidence.md:55`). This is an appropriate epistemic attribution. No issues.

---

## Claim 6

Claim:
PKCE prevents authorization-code interception: client creates `code_verifier` (43–128 unreserved chars, ≥256-bit entropy); sends `code_challenge` (S256 = BASE64URL(SHA256(verifier))); server binds challenge to code and verifies verifier at token endpoint.

Location:
`research/03-evidence.md` §Evidence 7; `research/05-report.md` §Finding 4

Evidence Provided:
RFC 7636 Sec 1.1 flow (A–D); Sec 4.1/4.2/4.6 (formulas); Sec 4.6 "If the values are equal, the token endpoint MUST continue... If not equal, invalid_grant."
RFC 9700 Sec 2.1.1 + 4.5.3.1

Source:
RFC 7636; RFC 9700

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
All specifics (character range 43–128, ≥256-bit entropy, S256 formula) are normatively defined in RFC 7636 Sec 4.1 and 4.2. The Appendix B worked example is also cited correctly (verifier → challenge transformation is deterministic and verifiable). No issues.

---

## Claim 7

Claim:
S256 is Mandatory-To-Implement (MTI) for PKCE; `plain` SHOULD NOT be used and is deprecated.

Location:
`research/03-evidence.md` §Evidence 8; `research/05-report.md` §Finding 4

Evidence Provided:
RFC 7636 Sec 4.2: "If the client is capable of using S256, it MUST use S256, as S256 is Mandatory To Implement on the server..."
RFC 9700 Sec 2.1.1: "clients SHOULD use PKCE methods that do not expose the verifier... Currently, S256 is the only such method"

Source:
RFC 7636 Sec 4.2, 7.2; RFC 9700 Sec 2.1.1

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Verbatim RFC 7636 language is reproduced correctly. RFC 9700 corroborates the S256 preference. No issues.

---

## Claim 8

Claim:
Implicit flow (response_type=token / id_token) is deprecated; clients SHOULD use Authorization Code Flow + PKCE instead.

Location:
`research/03-evidence.md` §Evidence 9; `research/05-report.md` §Finding 5

Evidence Provided:
RFC 9700 Sec 2.1.2: "clients SHOULD NOT use the implicit grant..."
OAuth 2.1: "The Implicit grant (response_type=token) is omitted."
Browser-based-apps draft Sec 7.2 (threat analysis)

Source:
RFC 9700; https://oauth.net/2.1/; draft-ietf-oauth-browser-based-apps-27

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
The research correctly notes in `04-contradictions.md` that OIDC Core 1.0 still documents Implicit Flow historically (Contradiction 1), with the accurate assessment that RFC 9700 is the current BCP. No issues.

---

## Claim 9

Claim:
Refresh tokens for public clients MUST be sender-constrained or use refresh token rotation; rotation invalidates the previous token; reuse signals breach and triggers revocation of the active refresh token.

Location:
`research/03-evidence.md` §Evidence 10; `research/05-report.md` §Finding 7

Evidence Provided:
RFC 9700 Sec 2.2.2 + 4.14.2 verbatim text (MUST, rotation, revocation on reuse)

Source:
RFC 9700

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
The research correctly attributes the rotation mandate to RFC 9700 and acknowledges that OIDC Core Sec 12 is structurally compatible but does not impose the MUST requirement (`03-evidence.md:91`). This is epistemically sound.

---

## Claim 10

Claim:
JWT validation MUST pin the algorithm (never trust `alg` header blindly; reject `none` unless explicitly configured); validate all nested crypto operations; use mutually exclusive rules per JWT kind.

Location:
`research/03-evidence.md` §Evidence 11; `research/05-report.md` §Finding 6

Evidence Provided:
RFC 8725 Sec 3.1 (MUST NOT use any other algorithms; each key MUST be used with exactly one algorithm); Sec 3.3 (reject if any operation fails); Sec 3.12 + 3.11 (explicit typing, distinct aud/iss/keys per JWT kind); threat Sec 2.1 (alg=none swap, RS256→HS256 confusion)

Source:
RFC 8725

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
RFC 8725 directly contains these requirements. Corroborated by OIDC Core Sec 2 (ID Tokens MUST be JWS-signed). No issues.

---

## Claim 11

Claim:
Tokens MUST NOT be stored in browser JS-accessible storage (e.g., localStorage); BFF pattern (Backend-for-Frontend) is the recommended secure architecture.

Location:
`research/03-evidence.md` §Evidence 12; `research/05-report.md` §Finding 6

Evidence Provided:
draft-ietf-oauth-browser-based-apps-27 Sec 6.1 (BFF definition), Sec 5 (malicious JS same privileges as app code)
RFC 9700 Sec 4.2/4.3 (credential leakage via referer, browser history)

Source:
draft-ietf-oauth-browser-based-apps-27; RFC 9700

Source Actually Supports Claim:
PARTIAL

Classification:
INTERPRETATION

Severity:
MEDIUM

Notes:
The draft (Sec 8 browser-based-apps) does NOT say tokens MUST NEVER be stored in localStorage — it says "applications MUST NOT use persistent token storage (e.g., localStorage) unless the tokens are sender-constrained or encrypted" (`04-contradictions.md:46`). The evidence file assigns confidence MEDIUM (`03-evidence.md:107`) and notes draft status. However, the claim in `05-report.md:80` ("hindari localStorage") simplifies this nuance. The research documents this distinction in `04-contradictions.md` (Contradiction 4), which is good practice, but the main findings section `05-report.md` §Finding 6 asserts the claim without the sender-constrained qualification. This is an oversimplification that could mislead implementers.

---

## Claim 12

Claim:
OAuth 2.1 consolidates RFC 6749, 6750, 7636, and deprecates/removes Implicit Grant, ROPC Grant, and mandates PKCE for all clients.

Location:
`research/02-sources.md` §Source 7; `research/05-report.md` §Areas of Agreement

Evidence Provided:
https://oauth.net/2.1/ summary page (not the direct IETF datatracker draft)

Source:
https://oauth.net/2.1/

Source Actually Supports Claim:
PARTIAL

Classification:
INTERPRETATION

Severity:
LOW

Notes:
oauth.net/2.1/ is a community-maintained summary, not the authoritative draft specification. The substance aligns with what is known about `draft-ietf-oauth-v2-1`, but should be attributed as a community summary. The research notes OAuth 2.1 is still a draft (`01-plan.md:38`). Acceptable as corroborating evidence for already RFC-grounded claims.

---

## Summary

| Claim | Status | Severity |
|-------|--------|----------|
| 1. OAuth 2.0 = authorization framework only | SUPPORTED | LOW |
| 2. Access tokens are opaque | PARTIAL (over-universalized "usually") | LOW |
| 3. OIDC adds ID Token JWT + required claims | SUPPORTED | LOW |
| 4. 13-step ID Token validation | SUPPORTED | LOW |
| 5. Access token login is unsafe | SUPPORTED | LOW |
| 6. PKCE mechanism specifics | SUPPORTED | LOW |
| 7. S256 is MTI; plain deprecated | SUPPORTED | LOW |
| 8. Implicit flow deprecated | SUPPORTED | LOW |
| 9. Refresh token rotation requirement | SUPPORTED | LOW |
| 10. JWT algorithm pinning | SUPPORTED | LOW |
| 11. No localStorage; BFF recommended | PARTIAL (sender-constrained exception omitted from main findings) | MEDIUM |
| 12. OAuth 2.1 consolidation | PARTIAL (community summary as source) | LOW |
