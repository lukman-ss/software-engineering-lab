# Contradictions Audit — OAuth 2.0 & OIDC Research

Target Lab: `labs/31-oauth2-and-oidc`

---

## Contradiction 1 — Implicit Flow Documented vs Deprecated

Statement A:
OIDC Core 1.0 defines Implicit Flow (response_type=`id_token token` / `id_token`) and Hybrid Flow as valid authentication flows in Sections 3.2 and 3.3, with a comparison table listing them alongside Authorization Code Flow.

Location:
`research/04-contradictions.md` §Contradiction 1 → attributed to OIDC Core 1.0 Sec 3.2/3.3

Statement B:
RFC 9700 Sec 2.1.2: "Clients SHOULD NOT use the implicit grant (response type token)... Clients SHOULD instead use the response type code."
OAuth 2.1: "The Implicit grant (response_type=token) is omitted from this specification."

Location:
`research/04-contradictions.md` §Contradiction 1 → attributed to RFC 9700 Sec 2.1.2 / OAuth 2.1

Type:
SOURCE_CONFLICT

Impact:
Moderate. The research correctly identifies this contradiction and resolves it with appropriate priority attribution (RFC 9700 as current BCP). However, the main findings section (`05-report.md`) does not explicitly alert readers to this split between the OIDC Core spec (which still formally includes Implicit) and the current BCP (which deprecates it). A developer reading only the findings chapter might not be aware that OIDC Core still formally defines Implicit Flow.

Assessment:
CORRECTLY IDENTIFIED by research. Resolution is sound. Findings chapter could benefit from an explicit alert.

---

## Contradiction 2 — Refresh Token Rotation: OIDC Core Example vs RFC 9700 Mandate

Statement A:
OIDC Core 1.0 Sec 12.2 example response includes a `refresh_token` field, implying issuance of a new token, but the spec text does not use the term "rotation" and does not require revocation of the old token.

Location:
`research/04-contradictions.md` §Contradiction 2 → attributed to OIDC Core 1.0 Sec 12.2

Statement B:
RFC 9700 Sec 4.14.2 explicitly mandates rotation: "The previous refresh token is invalidated... one of them will present an invalidated refresh token, which will inform the authorization server of the breach... it will revoke the active refresh token."

Location:
`research/04-contradictions.md` §Contradiction 2 → attributed to RFC 9700 Sec 4.14.2

Type:
SOURCE_CONFLICT

Impact:
Low. The research accurately identifies that OIDC Core is silent on the security mandate and correctly attributes the MUST requirement exclusively to RFC 9700. No misleading claim resulted.

Assessment:
CORRECTLY IDENTIFIED. Assessment in research is sound: OIDC Core is compatible with rotation but not prescriptive; RFC 9700 is the applicable security mandate.

---

## Contradiction 3 — OIDC Core `none` Algorithm Exception vs RFC 8725

Statement A:
OIDC Core Sec 2: "ID Tokens MUST NOT use `none` as the `alg` value unless the Response Type used returns no ID Token from the Authorization Endpoint... and the Client explicitly requested the use of `none` at Registration time."

Location:
`research/04-contradictions.md` §Contradiction 3 → attributed to OIDC Core Sec 2

Statement B:
RFC 8725 Sec 3.2: "`none` algorithm should only be used when the JWT is cryptographically protected by other means."

Location:
`research/04-contradictions.md` §Contradiction 3 → attributed to RFC 8725 Sec 3.2

Type:
SOURCE_CONFLICT

Impact:
Low. The research correctly identifies that both documents align (both allow `none` only under exceptional, explicitly-consented conditions; neither contradicts the other). Labeling it "Contradiction 3" is slightly misleading since the research's own Assessment concludes "No contradiction." An accurate characterization would be "complementary restrictions."

Assessment:
CORRECTLY RESOLVED. Research conclusion ("no contradiction but RFC 8725 is stricter on algorithm pinning") is technically accurate.

---

## Contradiction 4 — Token Storage: Absolute Prohibition vs Conditional Prohibition

Statement A:
`05-report.md` §Finding 6, §Lab Implementation Guidance: "hindari localStorage" (avoid localStorage), implying a blanket prohibition.
README.md (implementation context, not audited per pipeline override) uses similar framing.

Location:
`research/05-report.md` §Finding 6 (line 80); `research/05-report.md` §Lab Implementation Guidance (line 80)

Statement B:
`research/04-contradictions.md` §Contradiction 4 (line 46): draft-ietf-oauth-browser-based-apps-27 Sec 8.5 says: "applications MUST NOT use persistent token storage (e.g., localStorage) unless the tokens are sender-constrained or encrypted."

Location:
`research/04-contradictions.md` §Contradiction 4

Type:
INTERNAL

Impact:
MEDIUM. The main findings chapter states "hindari localStorage" as a blanket rule but the contradiction file correctly records the sender-constrained/encrypted exception from the authoritative draft. These two statements are internally inconsistent. A reader relying only on `05-report.md` would receive an oversimplified, slightly overstated security rule. The exception is real and important for implementers using DPoP or encrypted tokens.

Assessment:
INTERNAL CONTRADICTION. The research files contain this correctly in `04-contradictions.md` but the summary report (`05-report.md`) does not propagate the qualification. Needs revision.

---

## Additional Checks — Not Identified in Research

### Audit Check: OAuth 2.1 Source Tier Mismatch
The research (`02-sources.md` §Source 7) classifies `https://oauth.net/2.1/` as "Tier 1 (draft spec / summary)". The oauth.net page is not an IETF publication; it is a community-maintained informational summary. This is a source classification error. The actual draft is at `datatracker.ietf.org/doc/draft-ietf-oauth-v2-1/`. The content accuracy is likely high but the tier attribution is misleading.

Type:
INTERNAL

Impact:
LOW. Does not affect core conclusions (the OAuth 2.1 claims are also independently supported by RFC 9700 Sec 2.1.2, RFC 7636).

Assessment:
MINOR CLASSIFICATION ERROR — should be listed as Tier 2 / SECONDARY community source.
