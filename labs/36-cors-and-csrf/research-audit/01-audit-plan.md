# 01 — Audit Plan

Target Lab: `labs/36-cors-and-csrf`
Audit Type: Research Audit (Pipeline Override: Implementation/Code audit excluded)
Audit Date: 2026-09-28

## Files Reviewed
- `research/01-plan.md`
- `research/02-sources.md`
- `research/03-evidence.md`
- `research/04-contradictions.md`
- `research/05-report.md`
- `research/06-open-questions.md`

## Claims To Verify
1. Same-Origin Policy (SOP) allows cross-origin writes (form submissions) while restricting cross-origin reads.
2. CORS is a browser-enforced mechanism for response accessibility in JavaScript, not a server-side execution firewall.
3. `Access-Control-Allow-Origin: *` does not prevent CSRF attacks and is rejected by browsers when credentials are included.
4. Simple requests (`GET`, `HEAD`, `POST` with standard form MIME types) bypass preflight `OPTIONS` requests.
5. Preflight `OPTIONS` requests are triggered by non-simple methods or custom HTTP headers.
6. `SameSite=Lax` blocks cross-site form POST requests while allowing top-level GET navigations.
7. Naive Double-Submit Cookie pattern is vulnerable to cookie injection; Signed Double-Submit Cookie (HMAC) is recommended.
8. Fetch Metadata (`Sec-Fetch-Site`) headers provide modern context-aware CSRF protection.

## Primary Risks
- Over-reliance on secondary/community sources without verifying primary specifications (WHATWG Fetch, W3C Fetch Metadata, RFC 6265bis).
- Misrepresenting CORS preflight behavior as a server security boundary rather than a browser pre-check.
- Chronological inaccuracies regarding browser defaults (e.g., Chrome SameSite Lax rollout timeline).

## Audit Strategy
1. Audit all 10 cited sources for reachability, relevance, correct classification, and scope boundaries.
2. Verify all 15 major claims against source evidence and check for attribution accuracy.
3. Inspect internal consistency and potential unaddressed contradictions between research documents.
4. Categorize research gaps and evaluate impact on technical validity.
5. Formulate final verdict based on evidence-based quality gates.
