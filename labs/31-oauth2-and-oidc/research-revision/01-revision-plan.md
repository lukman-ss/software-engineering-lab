# Revision Plan

Target Lab: `labs/31-oauth2-and-oidc`

Previous Audit Status: APPROVED_WITH_WARNINGS

## Blocking Issues

None.

## Non-Blocking Issues (Gaps from Audit)

1. **Gap 1 (Medium - Token Storage Qualification)**: `05-report.md` states a blanket rule ("hindari localStorage") without noting the draft-ietf-oauth-browser-based-apps-27 Sec 8.5 exception for sender-constrained (e.g. DPoP) or encrypted tokens.
2. **Gap 2 (Low - Source Tier Mismatch)**: Source 7 (`oauth.net/2.1/`) in `02-sources.md` is listed as Tier 1 draft spec, but is a community summary. The direct IETF draft is `draft-ietf-oauth-v2-1`.
3. **Gap 3 (Low - Implicit Flow Context in OIDC Core)**: `05-report.md` Finding 5 does not explicitly alert readers that OIDC Core 1.0 Sec 3.2/3.3 still documents Implicit Flow, despite RFC 9700 deprecation.
4. **Gap 4 (Low - Structured Access Token Context)**: `05-report.md` Finding 1 and `03-evidence.md` Evidence 2 omit mention of RFC 9068 (JWT profile for access tokens), presenting access token opacity without noting modern standardized non-opaque forms.

## Files To Modify

- `labs/31-oauth2-and-oidc/research/02-sources.md`
- `labs/31-oauth2-and-oidc/research/03-evidence.md`
- `labs/31-oauth2-and-oidc/research/05-report.md`

## Verification Plan

- Verify source URLs and draft names.
- Ensure research files match audit recommendations and preserve all existing valid claims.
