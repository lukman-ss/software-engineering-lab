# 04 — Contradictions Audit

## Contradiction 1: SameSite=Lax Default Rollout Date & Grace Period
- Statement A: PortSwigger states "Since 2021, Chrome enforces Lax SameSite restrictions by default." (`research/02-sources.md:91`)
- Statement B: Google web.dev documents Chrome 80 rollout starting February 2020 (`research/02-sources.md:68-71`). MDN Set-Cookie documents a temporary 2-minute POST grace period in Chromium's implementation.
- Type: SOURCE_CONFLICT (Chronological / Implementation Granularity)
- Impact: Minimal. Both confirm that Lax-by-default is currently standard behavior across modern Chromium browsers.
- Assessment: RESOLVED. The research explicitly noted this timing discrepancy in `04-contradictions.md` and `06-open-questions.md`.

---

## Contradiction 2: CORS Preflight and Credentials (TLS Client Certificates)
- Statement A: WHATWG Fetch §3.3.3 specifies credentials mode is 'same-origin' (no credentials sent) for preflight OPTIONS.
- Statement B: Chromium bug 775438 / MDN CORS enterprise notes indicate Chromium browsers send TLS client certificates in CORS preflight requests.
- Type: SPEC_VS_BROWSER_BEHAVIOR
- Impact: Low. Affects mutual TLS (mTLS) enterprise setups; does not affect cookie/bearer token semantics.
- Assessment: RESOLVED. Accurately documented in `04-contradictions.md`.

---

## Contradiction 3: CORS as CSRF Defense vs Preflight Side-Effect
- Statement A: "CORS does not prevent CSRF" (`research/03-evidence.md:23`).
- Statement B: Custom headers trigger CORS preflight, thereby blocking unauthorized cross-origin calls if preflight is rejected (`research/03-evidence.md:185`).
- Type: CONCEPTUAL_DELIMITATION
- Impact: Medium if misunderstood by junior engineers.
- Assessment: RESOLVED. The research clearly distinguishes between CORS proper (response reading control) and making requests non-simple via custom headers (leveraging preflight side-effect).

---

## Overall Consistency Assessment
No unresolved or material contradictions exist. All technical points are cross-corroborated between WHATWG, MDN, OWASP, Google web.dev, and PortSwigger.
