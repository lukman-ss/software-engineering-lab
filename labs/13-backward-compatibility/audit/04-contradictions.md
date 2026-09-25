# 04 - Contradictions

## Contradiction 1

Statement A:
"Berdasarkan panduan Google AIP-180 dan praktek Stripe: ... Gunakan HTTP Deprecation Header (mis. `Deprecation: true`, `Sunset: ...`)" (Strategies for API Backward Compatibility)

Location:
`research/05-api-compatibility.md:4-17`

Statement B:
Stripe maintains backward compatibility by pinning API versions via the `Stripe-Version` header in requests, maintaining older versions indefinitely, and issuing explicit major/minor releases, rather than forcing standard HTTP Sunset headers.

Location:
Source 2 (Stripe API Versioning)

Type:
SOURCE_CONFLICT

Impact:
Misrepresents Stripe's actual API versioning architecture. Stripe does not rely on a sunset transition window for their public API; they maintain continuous version pinning to guarantee old clients never break.

Assessment:
Medium impact. The research correctly identifies deprecation headers as an industry pattern, but incorrectly attributes this specific sunsetting strategy to Stripe's core versioning model.

---

No other material contradictions found.
