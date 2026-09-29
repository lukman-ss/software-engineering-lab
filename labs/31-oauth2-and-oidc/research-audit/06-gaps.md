# Research Gap Analysis — OAuth 2.0 & OIDC Research

Target Lab: `labs/31-oauth2-and-oidc`

---

## Gap 1 — Overgeneralization: Access Token Storage Prohibition

Type:
OVERGENERALIZATION

Severity:
MEDIUM

Location:
`research/05-report.md` §Finding 6, §Lab Implementation Guidance (line 80)

Problem:
The report asserts "hindari localStorage" without documenting the exception noted in `draft-ietf-oauth-browser-based-apps-27` Section 8.5: persistent token storage (e.g., localStorage) is acceptable if the tokens are sender-constrained (e.g., DPoP) or encrypted. While `04-contradictions.md` acknowledges this nuance, it is omitted from the main findings and actionable guidance.

Required Revision:
Update `05-report.md` Section "Finding 6" and "Lab Implementation Guidance" to state: "LocalStorage is prohibited for bearer tokens; sender-constrained tokens (DPoP / mTLS) or encrypted tokens represent a scoped exception under draft-ietf-oauth-browser-based-apps-27 Section 8.5."

Can Be Approved Without Fix:
YES (Non-blocking for research foundation; caveat is recorded in contradiction file)

---

## Gap 2 — Source Tier Classification: OAuth 2.1 URL

Type:
SCOPE_ERROR

Severity:
LOW

Location:
`research/02-sources.md` §Source 7

Problem:
`https://oauth.net/2.1/` is classified as "Tier 1 (draft spec / summary)". While the content is accurate and widely accepted in the OAuth community, `oauth.net` is a community website maintained by individual contributors, not an official IETF working group publication. The official document is `draft-ietf-oauth-v2-1`.

Required Revision:
Change Source 7 classification to "Tier 2 (Community Educational Summary)" or replace URL with `https://datatracker.ietf.org/doc/draft-ietf-oauth-v2-1/`.

Can Be Approved Without Fix:
YES

---

## Gap 3 — OIDC Core Implicit Flow Deprecation Alert

Type:
MISSING_CASE

Severity:
LOW

Location:
`research/05-report.md` §Finding 5, §Areas of Disagreement

Problem:
`05-report.md` states that Implicit Flow is deprecated per RFC 9700 and OAuth 2.1, but does not provide an explicit note that developers referencing the primary OIDC Core 1.0 specification will still find Implicit and Hybrid flows listed as standard. Developers moving between the specs without this context may be confused.

Required Revision:
Add a brief callout in `05-report.md` Finding 5 clarifying that although OIDC Core 1.0 Sections 3.2 and 3.3 document Implicit and Hybrid flows, the current security consensus per RFC 9700 Section 2.1.2 supersedes them.

Can Be Approved Without Fix:
YES

---

## Gap 4 — Overgeneralization: Access Token Opacity

Type:
OVERGENERALIZATION

Severity:
LOW

Location:
`research/03-evidence.md` §Evidence 2, `research/05-report.md` §Finding 1

Problem:
Claim asserts "access tokens are typically opaque to the client" based on RFC 6749 Section 1.4 ("usually opaque"). However, structured access tokens (such as JWT profile for OAuth 2.0 Access Tokens per RFC 9068) are widely used in modern architectures. While RFC 6749 did not standardize JWT access tokens, RFC 9068 (October 2021) does. Mentioning only opacity without acknowledging RFC 9068 leaves out a major production pattern.

Required Revision:
Acknowledge RFC 9068 ("JSON Web Token (JWT) Profile for OAuth 2.0 Access Tokens") as the standard for non-opaque access tokens, while reinforcing that the client still treats them as opaque unless specifically configured otherwise.

Can Be Approved Without Fix:
YES

---

## Summary of Gaps

| Gap ID | Type | Severity | Can Be Approved Without Fix? |
|--------|------|----------|------------------------------|
| Gap 1 | OVERGENERALIZATION | MEDIUM | YES |
| Gap 2 | SCOPE_ERROR | LOW | YES |
| Gap 3 | MISSING_CASE | LOW | YES |
| Gap 4 | OVERGENERALIZATION | LOW | YES |

No CRITICAL or unresolvable HIGH gaps identified. All gaps are actionable recommendations for documentation refinement.
