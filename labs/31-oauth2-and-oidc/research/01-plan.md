# Research Plan

## Research Topic
OAuth 2.0 & OIDC — Perbedaan Authentication vs Authorization dan Mekanisme Keamanan Flow

## Objective
Provide evidence-based understanding of:
1. Apa perbedaan mendasar antara OAuth 2.0 (Authorization) dan OIDC (Authentication)?
2. Mengapa access token OAuth 2.0 tidak boleh digunakan untuk login (authentication)?
3. Bagaimana Authorization Code Flow + PKCE bekerja?
4. Bagaimana ID Token diverifikasi secara kriptografis?
5. Apa jebakan umum dan bagaimana mitigasinya?

## Research Questions
- RQ1: Bagaimana OAuth 2.0 didefinisikan secara resmi dalam RFC 6749?
- RQ2: Bagaimana OIDC menambahkan identity layer di atas OAuth 2.0?
- RQ3: Mengapa implicit flow dideprecate dan PKCE diwajibkan?
- RQ4: Apa saja claim yang wajib ada dalam ID Token dan bagaimana validasinya?
- RQ5: Mengapa access token tidak boleh digunakan untuk authentication?
- RQ6: Bagaimana refresh token rotation bekerja?

## Search Strategy
- Tier 1: IETF RFC (6749, 7636, 7519, 8725, 9700, 10017)
- Tier 1: OpenID Connect Core 1.0 spec
- Tier 2: OAuth.net summary (OAuth 2.1 differences)
- Verified by cross-checking multiple RFCs for consistency

## Expected Primary Sources
- RFC 6749 (OAuth 2.0)
- RFC 7636 (PKCE)
- RFC 7519 (JWT)
- RFC 8725 (JWT BCP)
- RFC 9700 (OAuth 2.0 Security BCP)
- OpenID Connect Core 1.0

## Risks / Unknowns
- Some OIDC Core sections are very long; evidence extracted via section-specific retrieval
- OAuth 2.1 is still a draft (not yet RFC)
- BFF pattern details depend on deployment; not prescriptive in standards
