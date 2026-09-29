# 02 — Source Audit

## Source 1
Claimed Title: Same-origin policy - Security | MDN
Claimed Publisher: Mozilla Developer Network (MDN)
URL: https://developer.mozilla.org/en-US/docs/Web/Security/Defenses/Same-origin_policy
Reachable: YES
Source Type: PRIMARY (Authoritative Web Standards Documentation)
Relevant: YES
Supports Claimed Topic: YES
Problems: None. MDN provides direct authoritative definitions of SOP write/read rules.
Assessment: PASS

---

## Source 2
Claimed Title: Cross-origin request forgery (CSRF) - Security | MDN
Claimed Publisher: MDN
URL: https://developer.mozilla.org/en-US/docs/Web/Security/Attacks/CSRF
Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: None.
Assessment: PASS

---

## Source 3
Claimed Title: Cross-Origin Resource Sharing (CORS) - HTTP | MDN
Claimed Publisher: MDN
URL: https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/CORS
Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: None. Detailed breakdown of simple requests, preflight OPTIONS, and wildcard rules.
Assessment: PASS

---

## Source 4
Claimed Title: Fetch Standard (WHATWG Living Standard)
Claimed Publisher: WHATWG
URL: https://fetch.spec.whatwg.org/
Reachable: YES
Source Type: PRIMARY (Normative Web Standard)
Relevant: YES
Supports Claimed Topic: YES
Problems: None. Normative specification for CORS, preflight fetch algorithm, safelisted methods and headers.
Assessment: PASS

---

## Source 5
Claimed Title: Cross-Site Request Forgery Prevention Cheat Sheet | OWASP
Claimed Publisher: OWASP
URL: https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html
Reachable: YES
Source Type: PRIMARY (Industry Security Standard)
Relevant: YES
Supports Claimed Topic: YES
Problems: None. Definitive guide for Synchronizer Token, Double Submit Cookie (Naive vs Signed), SameSite, and Fetch Metadata.
Assessment: PASS

---

## Source 6
Claimed Title: Set-Cookie header - HTTP | MDN
Claimed Publisher: MDN
URL: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Set-Cookie
Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: None. Exact specification of SameSite attributes (Strict, Lax, None) and cookie security prefixes.
Assessment: PASS

---

## Source 7
Claimed Title: SameSite cookies explained | web.dev (Google)
Claimed Publisher: Google / Chromium
URL: https://web.dev/articles/samesite-cookies-explained
Reachable: YES
Source Type: PRIMARY (Browser Engine Implementation Doc)
Relevant: YES
Supports Claimed Topic: YES
Problems: None. Documents Chrome 80 Lax-by-default behavior and top-level navigation context.
Assessment: PASS

---

## Source 8
Claimed Title: What is CORS (cross-origin resource sharing)? | PortSwigger Web Security Academy
Claimed Publisher: PortSwigger
URL: https://portswigger.net/web-security/cors
Reachable: YES
Source Type: SECONDARY (Reputable Security Authority)
Relevant: YES
Supports Claimed Topic: YES
Problems: None. Authoritative technical reference on CORS misconceptions and exploit vectors.
Assessment: PASS

---

## Source 9
Claimed Title: What is CSRF (Cross-site request forgery)? | PortSwigger Web Security Academy
Claimed Publisher: PortSwigger
URL: https://portswigger.net/web-security/csrf
Reachable: YES
Source Type: SECONDARY
Relevant: YES
Supports Claimed Topic: YES
Problems:
- PortSwigger text mentions "Since 2021, Chrome enforces Lax SameSite restrictions by default" whereas Chromium/web.dev notes Chrome 80 released in early 2020. This discrepancy is properly captured and flagged in `04-contradictions.md` and `06-open-questions.md`.
Assessment: PASS

---

## Source 10
Claimed Title: Samesite cookies explained (Chrome Lax-by-default announcement context)
Claimed Publisher: web.dev / Chromium
URL: https://web.dev/articles/samesite-cookies-explained#changes-to-the-default-behavior-without-samesite
Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: None. Sub-anchor to Source 7, verified as genuine URL and content.
Assessment: PASS
