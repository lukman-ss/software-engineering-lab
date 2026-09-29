# 03 — Claim Audit

## Claim 1
Claim: Same-Origin Policy (SOP) permits cross-origin writes (e.g., links, redirects, form submissions) but restricts cross-origin reads.
Location: `research/03-evidence.md:5-7`, `research/05-report.md:28-36`
Evidence Provided: Direct quotes from MDN SOP and PortSwigger CORS.
Source: MDN SOP (`https://developer.mozilla.org/en-US/docs/Web/Security/Defenses/Same-origin_policy`)
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Core principle of web architecture. Fully substantiated.

---

## Claim 2
Claim: CORS is a browser mechanism, not server-side firewall; server executes cross-origin requests even when the browser blocks the response.
Location: `research/03-evidence.md:21-27`, `research/05-report.md:46-54`
Evidence Provided: MDN CORS, PortSwigger CORS, WHATWG Fetch Standard.
Source: WHATWG Fetch Standard §3.3, PortSwigger CORS
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Crucial distinction between request dispatch and response read access.

---

## Claim 3
Claim: `Access-Control-Allow-Origin: *` cannot be used with credentials (`credentials: include`). The browser rejects response access if wildcard is provided with credentials.
Location: `research/03-evidence.md:39-45`, `research/05-report.md:46-54`
Evidence Provided: WHATWG Fetch specification and MDN CORS guidelines.
Source: WHATWG Fetch Standard, MDN CORS
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Prevents trivial universal credential theft via wildcard CORS headers.

---

## Claim 4
Claim: Simple requests (GET, HEAD, POST with safelisted Content-Type) do not trigger preflight OPTIONS and are dispatched immediately by modern browsers.
Location: `research/03-evidence.md:57-63`, `research/05-report.md:64-72`
Evidence Provided: WHATWG Fetch §2.2.1/§2.2.2 and MDN CORS.
Source: WHATWG Fetch Standard
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Key technical foundation for CSRF viability.

---

## Claim 5
Claim: Preflight requests (OPTIONS) are mandated for non-safelisted HTTP methods (PUT, DELETE, PATCH) or custom headers (e.g. `X-CSRF-Token`, `Authorization`).
Location: `research/03-evidence.md:75-81`, `research/05-report.md:64-72`
Evidence Provided: Fetch Standard §cors-preflight-fetch and MDN CORS.
Source: WHATWG Fetch Standard, MDN CORS
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Well-documented standard behavior.

---

## Claim 6
Claim: `SameSite=Lax` sends cookies on same-site requests and top-level GET navigations, but withholds cookies on cross-site POST.
Location: `research/03-evidence.md:93-99`, `research/05-report.md:82-93`
Evidence Provided: MDN Set-Cookie, web.dev SameSite Explained.
Source: MDN Set-Cookie, Google web.dev
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Accurately reflects cookie behavior across standard modern user agents.

---

## Claim 7
Claim: `SameSite=Strict` completely withholds cookies on any cross-site request, while `SameSite=None` requires the `Secure` attribute.
Location: `research/03-evidence.md:111-117`, `research/05-report.md:82-93`
Evidence Provided: MDN Set-Cookie, web.dev.
Source: MDN Set-Cookie
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Fully aligns with RFC 6265bis specifications.

---

## Claim 8
Claim: SameSite operates at registrable domain (eTLD+1) level ("site"), not "origin". Vulnerable sibling subdomains can compromise same-site protections.
Location: `research/03-evidence.md:125-127`, `research/04-contradictions.md:41-50`, `research/05-report.md:82-93`
Evidence Provided: MDN CSRF, OWASP CSRF Prevention, PortSwigger.
Source: MDN CSRF, OWASP
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Essential security nuance explicitly highlighted in research.

---

## Claim 9
Claim: Naive Double-Submit Cookie pattern is vulnerable to cookie injection (sibling domain, plaintext HTTP); Signed Double-Submit Cookie (HMAC bound to session/identity) or Synchronizer Token Pattern is required.
Location: `research/03-evidence.md:147-153`, `research/05-report.md:104-114`
Evidence Provided: OWASP CSRF Prevention Cheat Sheet.
Source: OWASP CSRF Prevention Cheat Sheet
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Accurately reflects modern OWASP best practices.

---

## Claim 10
Claim: Fetch Metadata headers (`Sec-Fetch-Site`) enable the backend to identify request context (`same-origin`, `same-site`, `cross-site`, `none`) and reject cross-site state-changing actions without client tokens.
Location: `research/03-evidence.md:201-207`, `research/05-report.md:123-131`
Evidence Provided: OWASP CSRF Prevention, MDN CSRF, W3C Fetch Metadata.
Source: OWASP CSRF Prevention, MDN CSRF
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Correctly notes that older browsers require fallback verification.
