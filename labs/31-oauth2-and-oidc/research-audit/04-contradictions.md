# Contradictions Audit

Target Lab: `labs/31-oauth2-and-oidc`

---

## Contradiction 1

Statement A:
OIDC Core 1.0 Sec 3.2 / 3.3 defines Implicit Flow (`response_type=id_token token` / `id_token`) and Hybrid Flow (`code id_token`, etc.) as standard flows.

Location:
`research/04-contradictions.md` Contradiction 1 (citing OIDC Core 1.0)

Statement B:
RFC 9700 Sec 2.1.2 and OAuth 2.1 deprecate / omit the Implicit Grant.

Location:
`research/04-contradictions.md` Contradiction 1 (citing RFC 9700 Sec 2.1.2 and OAuth 2.1)

Type:
SOURCE_CONFLICT

Impact:
Implementers may be confused about whether Implicit Flow is acceptable for new browser-based OIDC implementations.

Assessment:
The research correctly identifies the timeline and normative hierarchy: OIDC Core 1.0 was finalized in 2014 (errata 2023) prior to RFC 9700 (BCP 240, Jan 2025). RFC 9700 is the current authoritative Best Current Practice for OAuth 2.0 / OIDC security. Authorization Code Flow + PKCE (`response_type=code` with `scope=openid`) supersedes Implicit Flow for all new deployments. The research resolves this contradiction cleanly and accurately.

---

## Contradiction 2

Statement A:
OIDC Core 1.0 Sec 12.2 presents a sample refresh token response that contains a new `refresh_token`, but does not use the term "rotation" or specify that the previous token must be revoked.

Location:
`research/04-contradictions.md` Contradiction 2 (citing OIDC Core 1.0 Sec 12.2)

Statement B:
RFC 9700 Sec 2.2.2 / 4.14.2 explicitly mandates refresh token rotation or sender-constraining for public clients and requires that previous refresh tokens are invalidated.

Location:
`research/04-contradictions.md` Contradiction 2 (citing RFC 9700 Sec 2.2.2 and 4.14.2)

Type:
SOURCE_CONFLICT

Impact:
Implementers might assume that returning the same refresh token or issuing a new one without revoking the previous one is sufficient.

Assessment:
The research correctly notes that OIDC Core 1.0 provides an illustrative response format, whereas RFC 9700 provides prescriptive security hardening requirements. RFC 9700 BCP governs the security posture. No material contradiction exists, and the resolution is sound.

---

## Contradiction 3

Statement A:
OIDC Core 1.0 Sec 2 allows `none` as an `alg` value for ID Tokens under restricted conditions (when no ID Token is returned via authorization endpoint and client explicitly registered for `none`).

Location:
`research/04-contradictions.md` Contradiction 3 (citing OIDC Core 1.0 Sec 2)

Statement B:
RFC 8725 Sec 3.1 / 3.2 specifies that JWT libraries MUST perform algorithm verification and SHOULD NOT use or accept `none` unless explicitly requested, with cryptographic transport protection.

Location:
`research/04-contradictions.md` Contradiction 3 (citing RFC 8725 Sec 3.1/3.2)

Type:
INTERNAL / SOURCE_CONFLICT

Impact:
Implementers might inadvertently allow `alg: "none"` without proper channel security.

Assessment:
Both specifications agree that `none` requires transport-layer security or out-of-band integrity and must not be accepted by default. The research's guidance to strictly pin algorithms and reject `none` by default per RFC 8725 is accurate and secure.

---

## Contradiction 4

Statement A:
RFC 9700 Sec 4.2 / 4.3 outlines credential leakage vectors via referer headers and browser history.

Location:
`research/04-contradictions.md` Contradiction 4 (citing RFC 9700)

Statement B:
`draft-ietf-oauth-browser-based-apps-27` Sec 8.5 states applications MUST NOT use persistent token storage (`localStorage`) unless tokens are sender-constrained or encrypted, while recommending BFF or in-memory storage.

Location:
`research/04-contradictions.md` Contradiction 4 (citing browser-based-apps draft Sec 8.5)

Type:
SOURCE_CONFLICT / REFINEMENT

Impact:
Developer confusion regarding whether localStorage is ever permissible for OAuth tokens.

Assessment:
The research captures the nuance: bearer tokens in localStorage are vulnerable to XSS; only sender-constrained (e.g. DPoP) or server-held tokens (BFF) provide robust defenses. The resolution is sound.

---

## Summary
All 4 identified contradictions are real technical tensions between older foundation RFCs and modern Best Current Practice (BCP) documents. Each contradiction has been clearly explained and resolved according to normative precedence.
