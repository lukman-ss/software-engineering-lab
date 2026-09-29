# Changes Made

## Revision 1

Audit Issue:
Medium — Token Storage Overgeneralization (Gap 1)
`05-report.md` Finding 6 and Lab Implementation Guidance (line 80) stated "hindari localStorage" without noting the draft-ietf-oauth-browser-based-apps-27 Section 8.5 exception for sender-constrained (DPoP/mTLS) or encrypted tokens.

Files Changed:
- `research/05-report.md`

Action:
- Added qualifying Note to Finding 6 explaining that localStorage is prohibited for bearer tokens, but draft-ietf-oauth-browser-based-apps-27 Sec 8.5 permits persistent storage for sender-constrained or encrypted tokens.
- Updated Lab Implementation Guidance item 6 to include the scoped exception explicitly.

Verification:
- Claim now aligned with draft-ietf-oauth-browser-based-apps-27 Sec 8.5 as cited in audit/06-gaps.md and audit/04-contradictions.md.

Status:
RESOLVED

---

## Revision 2

Audit Issue:
Low — Source Tier Mismatch (Gap 2)
Source 7 `oauth.net/2.1/` was classified as "Tier 1 (draft spec / summary)" but is a community-maintained website, not the IETF working group's authoritative document.

Files Changed:
- `research/02-sources.md`
- `research/05-report.md` (Finding 5 source reference updated from the oauth.net URL to draft name)

Action:
- Changed Source 7 to reference `https://datatracker.ietf.org/doc/draft-ietf-oauth-v2-1/` as the primary URL with Publisher: "IETF OAuth Working Group".
- Changed Source Tier from "Tier 1 (draft spec / summary)" to "Tier 1 (draft spec)".
- Retained the community summary URL as a secondary informational reference.
- Updated Finding 5 sources to cite "draft-ietf-oauth-v2-1" instead of the oauth.net URL.

Verification:
- Source now correctly points to the official IETF datatracker document for OAuth 2.1.

Status:
RESOLVED

---

## Revision 3

Audit Issue:
Low — OIDC Core Implicit Flow Deprecation Alert (Gap 3)
Finding 5 did not alert readers that OIDC Core 1.0 Sections 3.2 and 3.3 still document Implicit/Hybrid flows, which may confuse developers reading both specs.

Files Changed:
- `research/05-report.md`

Action:
- Added a Note to Finding 5 clarifying that OIDC Core 1.0 Sec 3.2/3.3 documents Implicit/Hybrid flows historically, and that RFC 9700 Sec 2.1.2 / draft OAuth 2.1 supersede this; Authorization Code Flow + PKCE is mandatory.

Verification:
- Claim now provides necessary context to prevent developer confusion between OIDC Core historical documentation and current security BCP.

Status:
RESOLVED

---

## Revision 4

Audit Issue:
Low — Access Token Opacity Overgeneralization (Gap 4)
`05-report.md` Finding 1 and `03-evidence.md` Evidence 2 presented "access tokens are opaque" without acknowledging RFC 9068 (JWT Profile for OAuth 2.0 Access Tokens), which standardizes non-opaque JWT access tokens widely used in production.

Files Changed:
- `research/05-report.md`
- `research/03-evidence.md`

Action:
- Added a Note to Finding 1 in `05-report.md` explaining that RFC 6749's "usually opaque" applies to bearer tokens; RFC 9068 standardizes structured JWT access tokens; clients SHOULD still treat them as opaque unless the deployment profile requires otherwise.
- Extended Evidence 2 Notes in `03-evidence.md` with the same RFC 9068 qualification.

Verification:
- Claim now consistent with RFC 9068 (https://www.rfc-editor.org/rfc/rfc9068.html) scope.

Status:
RESOLVED
