# 04 — Contradiction Audit

## Contradiction 1: SameSite=Lax Default Behavior in Chromium vs Explisit SameSite=Lax
Statement A: Chromium default SameSite=Lax includes a 2-minute grace period allowing top-level POST requests without SameSite attribute.
Location: `research/04-contradictions.md:5-12`, web.dev, MDN Set-Cookie.
Statement B: Standards specify Lax blocks all cross-site POST requests.
Location: RFC 6265bis, WHATWG.
Type: SOURCE_CONFLICT (Implementation nuance vs standard specification)
Impact: None on core security recommendations; research explicitly notes the nuance and insists on explicit `SameSite` configuration.
Assessment: RESOLVED (Documented accurately with implementation caveats).

---

## Contradiction 2: TLS Client Certificates in Preflight Requests
Statement A: WHATWG Fetch §3.3.3 specifies credentials mode is 'same-origin' for preflight (no credentials sent).
Location: `research/04-contradictions.md:16-24`, WHATWG Fetch.
Statement B: Chromium sends TLS client certificates in CORS preflight (Chrome bug 775438).
Location: MDN CORS notes.
Type: SOURCE_CONFLICT (Browser implementation divergence)
Impact: Edge case limited to mTLS/TLS client auth, does not affect cookie or Bearer token auth.
Assessment: RESOLVED (Documented accurately as browser implementation anomaly).

---

## Contradiction 3: CORS as CSRF Defense vs Side-Effect Preflight Defense
Statement A: "CORS is not a protection against CSRF."
Location: PortSwigger CORS, MDN SOP.
Statement B: "Custom request headers rely on CORS preflight to prevent CSRF."
Location: OWASP CSRF Prevention Cheat Sheet.
Type: INTERNAL / TERMINOLOGY_NUANCE
Impact: Potential confusion if developers equate preflight side-effects with CORS origin policy.
Assessment: RESOLVED (The research explicitly differentiates between CORS response filtering and the preflight side-effect triggered by custom headers).

---

## Contradiction 4: Double-Submit Cookie Granularity
Statement A: General descriptions often cite "Double-Submit Cookie" without qualification.
Location: Lab specification context.
Statement B: OWASP explicitly discourages Naive Double-Submit and mandates Signed Double-Submit (HMAC).
Location: OWASP CSRF Prevention Cheat Sheet.
Type: SOURCE_CONFLICT (Precision difference)
Impact: High risk if developers implement the naive pattern.
Assessment: RESOLVED (The research report explicitly mandates the Signed Double-Submit pattern and warns against the naive variant).

---

## Overall Assessment
No unaddressed material contradictions found across the research documents. All nuanced divergences between specs and browser implementations are properly contextualized.
