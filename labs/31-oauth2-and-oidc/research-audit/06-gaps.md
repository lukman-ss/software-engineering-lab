# Research Gap Analysis: OAuth 2.0 & OIDC Research

Target Lab: `labs/31-oauth2-and-oidc`
Audit Scope: Research Completeness and Gaps

---

## Gap 1

Type:
SCOPE_ERROR

Severity:
LOW

Location:
`05-report.md` Section: Limitations, Item 2 & 3

Problem:
JWKS dynamic key rotation/caching nuances and native app custom scheme vs claimed universal links (iOS Universal Links / Android App Links) are mentioned at a high level but not deeply elaborated.

Required Revision:
None blocking. The lab design phase can detail JWKS parsing and callback handling in code.

Can Be Approved Without Fix:
YES

---

## Gap 2

Type:
UNVERIFIED_CLAIM

Severity:
LOW

Location:
`05-report.md` Finding 4

Problem:
Browser storage mechanisms (HttpOnly cookies / BFF) are not direct normative clauses in RFC 9700, though they directly derive from RFC 6819 and OAuth 2.0 security considerations.

Required Revision:
The report has already accurately classified this under "Limitations" as derived best practice.

Can Be Approved Without Fix:
YES
