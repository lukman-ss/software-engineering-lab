# Research Gap Analysis

Target Lab: `labs/31-oauth2-and-oidc`

---

## Gap 1

Type:
WEAK_SOURCE

Severity:
LOW

Location:
`research/02-sources.md` Source 8

Problem:
Source 8 is titled `RFC 10017 / draft-ietf-oauth-browser-based-apps-27 — OAuth 2.0 for Browser-Based Applications`, referencing an active draft URL (`draft-ietf-oauth-browser-based-apps-27`) with publication year noted as `July 2026 (draft)`. RFC 10017 has been allocated in the datatracker as the published version of this specification. While the content is authoritative and verified directly against the draft text, mixing the draft identifier and future RFC number in the title without an explicit RFC URL is slightly imprecise.

Required Revision:
Include both the RFC 10017 direct RFC-Editor reference and the working draft historical reference for clarity.

Can Be Approved Without Fix:
YES

---

## Gap 2

Type:
MISSING_SOURCE

Severity:
LOW

Location:
`research/03-evidence.md` Evidence 2 & `research/05-report.md` Finding 1

Problem:
RFC 9068 ("JSON Web Token (JWT) Profile for OAuth 2.0 Access Tokens", October 2021) is cited inline in notes under Evidence 2 and Finding 1 to explain modern structured JWT access tokens versus opaque access tokens, but RFC 9068 is not formally enumerated in `02-sources.md`.

Required Revision:
Add RFC 9068 as an official Tier 1 source entry in `02-sources.md` for completeness.

Can Be Approved Without Fix:
YES (The standard is cited inline and correctly summarized; omission from the formal list does not invalidate any core claims).

---

## Gap 3

Type:
SCOPE_ERROR

Severity:
LOW

Location:
`research/05-report.md` Section "Limitations" & `research/06-open-questions.md` RQ7

Problem:
The research focuses predominantly on browser-based client applications (SPAs and BFF architectures). Native mobile applications (iOS / Android) and machine-to-machine (M2M) confidential client flows are only mentioned tangentially.

Required Revision:
Clarify in the research plan and report scope that the primary focus of this lab is browser-based and web application flows (Auth Code + PKCE + OIDC), with native app specifics left as an open question / future lab topic.

Can Be Approved Without Fix:
YES (The target lab topic is `oauth2-and-oidc` foundation and browser security flows; scope focus is appropriate).
