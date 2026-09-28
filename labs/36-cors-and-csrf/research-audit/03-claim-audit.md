# 03 — Claim Audit

## Claim 1
Claim: Same-Origin Policy (SOP) permits cross-origin writes (e.g., HTML form submissions) while disallowing cross-origin reads.
Location: `research/03-evidence.md:5-17`, `research/05-report.md:28-43`
Evidence Provided: MDN Same-origin policy documentation.
Source: MDN Web Docs, Fetch Standard.
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Core architectural foundation of browser security model.

---

## Claim 2
Claim: CORS is a browser-enforced mechanism controlling JavaScript response readability, not a server-side execution firewall. Servers execute requests regardless of CORS header checks on the browser.
Location: `research/03-evidence.md:21-36`, `research/05-report.md:46-62`
Evidence Provided: MDN CORS, PortSwigger CORS, WHATWG Fetch Standard.
Source: MDN, PortSwigger, WHATWG.
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Distinguishes browser client policy from server processing semantics.

---

## Claim 3
Claim: `Access-Control-Allow-Origin: *` does not allow credentialed access; browsers reject responses when credentials mode is `include` and wildcard origin is returned.
Location: `research/03-evidence.md:39-54`, `research/05-report.md:46-62`
Evidence Provided: MDN CORS credentialed requests specification, WHATWG Fetch §3.3.
Source: MDN, WHATWG.
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Fully aligns with normative W3C/WHATWG specification.

---

## Claim 4
Claim: Simple requests (`GET`, `HEAD`, `POST` with safelisted Content-Types) bypass preflight `OPTIONS` and are dispatched directly to the server by browsers.
Location: `research/03-evidence.md:57-72`, `research/05-report.md:64-79`
Evidence Provided: WHATWG Fetch Standard §2.2.1-2.2.2, MDN CORS simple requests.
Source: WHATWG, MDN.
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Direct citation of normatively defined safelisted methods and headers.

---

## Claim 5
Claim: Preflight `OPTIONS` request is triggered when non-safelisted methods (e.g., `PUT`, `DELETE`) or custom request headers are used.
Location: `research/03-evidence.md:75-90`, `research/05-report.md:64-79`
Evidence Provided: WHATWG Fetch Standard, MDN CORS.
Source: WHATWG, MDN, OWASP.
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Used by defenses (e.g., custom header CSRF defense).

---

## Claim 6
Claim: `SameSite=Lax` restricts cookie inclusion on cross-site subrequests (such as `POST`), but allows cookies on top-level navigations using safe methods (`GET`).
Location: `research/03-evidence.md:93-108`, `research/05-report.md:82-101`
Evidence Provided: MDN Set-Cookie, web.dev, PortSwigger.
Source: MDN, Google web.dev.
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Nuance of Chromium 2-minute POST grace period noted.

---

## Claim 7
Claim: `SameSite=Strict` prevents cookie transmission on any cross-site request, while `SameSite=None` requires the `Secure` attribute.
Location: `research/03-evidence.md:111-126`, `research/05-report.md:82-101`
Evidence Provided: MDN Set-Cookie, web.dev.
Source: MDN, web.dev, RFC 6265bis.
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Standards-compliant definition.

---

## Claim 8
Claim: Synchronizer Token Pattern requires unique, secret, and unpredictable per-session or per-request tokens generated via CSPRNG on the server side.
Location: `research/03-evidence.md:129-144`, `research/05-report.md:104-120`
Evidence Provided: OWASP CSRF Prevention Cheat Sheet, MDN CSRF.
Source: OWASP, MDN.
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Canonical CSRF prevention reference.

---

## Claim 9
Claim: Naive Double-Submit Cookie pattern is vulnerable to cookie injection (e.g. via sibling subdomains); Signed Double-Submit Cookie (HMAC bound to session) is required.
Location: `research/03-evidence.md:147-162`, `research/05-report.md:104-120`
Evidence Provided: OWASP CSRF Prevention Cheat Sheet.
Source: OWASP.
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Highlights vulnerability in naive implementations.

---

## Claim 10
Claim: CORS is not a replacement for server-side security controls.
Location: `research/03-evidence.md:165-180`, `research/05-report.md:46-62`
Evidence Provided: PortSwigger CORS, MDN CORS.
Source: PortSwigger, MDN.
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Rebuts common developer misconception.

---

## Claim 11
Claim: Custom request headers trigger CORS preflight, allowing API endpoints to verify the presence of custom headers as a CSRF defense.
Location: `research/03-evidence.md:183-198`, `research/05-report.md:64-79`
Evidence Provided: OWASP CSRF Prevention, MDN CSRF, WHATWG Fetch.
Source: OWASP, MDN, WHATWG.
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Applicable to XHR/Fetch APIs, not traditional HTML form submissions.

---

## Claim 12
Claim: Fetch Metadata request headers (`Sec-Fetch-Site`) indicate request initiation context (`same-origin`, `same-site`, `cross-site`, `none`) and cannot be manipulated by client-side JS.
Location: `research/03-evidence.md:201-216`, `research/05-report.md:123-137`
Evidence Provided: OWASP CSRF Prevention, MDN CSRF, W3C Fetch Metadata.
Source: OWASP, MDN, W3C.
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Reliable modern browser defense signal.

---

## Claim 13
Claim: Dynamically reflecting the `Origin` header into `Access-Control-Allow-Origin` without validation creates severe data exposure vulnerabilities.
Location: `research/03-evidence.md:219-234`, `research/05-report.md:46-62`
Evidence Provided: PortSwigger CORS.
Source: PortSwigger, OWASP.
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Well-documented CORS misconfiguration pattern.

---

## Claim 14
Claim: CSRF attack feasibility requires three conditions: relevant action, cookie-based session handling, and no unpredictable request parameters.
Location: `research/03-evidence.md:237-252`
Evidence Provided: PortSwigger CSRF, MDN CSRF.
Source: PortSwigger, MDN.
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Accurate condition trilemma.

---

## Claim 15
Claim: `SameSite` operates at the registrable domain level ("site" / eTLD+1), not the strict "origin" level, allowing cross-subdomain attacks.
Location: `research/03-evidence.md:255-269`, `research/05-report.md:82-101`
Evidence Provided: MDN CSRF, PortSwigger CSRF, OWASP.
Source: MDN, PortSwigger, OWASP.
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Important boundary clarification between site and origin.
