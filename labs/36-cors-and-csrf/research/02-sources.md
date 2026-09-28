# 02 — Sumber Riset

## Source 1

Title: Same-origin policy - Security | MDN
Publisher: Mozilla Developer Network (MDN)
URL: https://developer.mozilla.org/en-US/docs/Web/Security/Defenses/Same-origin_policy
Published: 2026-09-17
Accessed: 2026-09-28
Source Tier: Tier 1 (dokumentasi otoritatif browser)
Relevance: Definisi SOP, pengklasifikasian cross-origin writes/reads/embeds, CSRF token sebagai mekanisme pertahanan tulis lintas origin.

## Source 2

Title: Cross-origin request forgery (CSRF) - Security | MDN
Publisher: MDN
URL: https://developer.mozilla.org/en-US/docs/Web/Security/Attacks/CSRF
Published: 2026-06-08
Accessed: 2026-09-28
Source Tier: Tier 1
Relevance: Mekanisme serangan CSRF, simple vs non-simple request, SameSite cookies, CSRF token, Fetch metadata, hubungan CORS dengan CSRF.

## Source 3

Title: Cross-Origin Resource Sharing (CORS) - HTTP | MDN
Publisher: MDN
URL: https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/CORS
Published: 2026-09-04
Accessed: 2026-09-28
Source Tier: Tier 1
Relevance: Simple request vs preflight, `Access-Control-Allow-Origin`, credentials dan wildcard, CORS tidak menggantikan server-side security.

## Source 4

Title: Fetch Standard (WHATWG Living Standard)
Publisher: WHATWG
URL: https://fetch.spec.whatwg.org/
Published: 2026-09-21 (Living Standard)
Accessed: 2026-09-28
Source Tier: Tier 1 (standar web)
Relevance: Definisi normatif CORS-safelisted method, CORS-safelisted request-header, preflight fetch, credential handling, Access-Control-Allow-Origin semantics.

## Source 5

Title: Cross-Site Request Forgery Prevention Cheat Sheet | OWASP
Publisher: OWASP
URL: https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html
Published: 2024 (latest revision as of access date)
Accessed: 2026-09-28
Source Tier: Tier 1 (standar keamanan industri)
Relevance: Synchronizer token, double-submit cookie, Fetch Metadata headers (Sec-Fetch-Site), SameSite cookie attribute, custom request header + CORS, keterbatasan SameSite.

## Source 6

Title: Set-Cookie header - HTTP | MDN
Publisher: MDN
URL: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Set-Cookie
Published: 2026-09-01
Accessed: 2026-09-28
Source Tier: Tier 1
Relevance: SameSite attribute (Strict, Lax, None), perilaku default Lax di beberapa browser, cookie prefix (__Host-, __Secure-).

## Source 7

Title: SameSite cookies explained | web.dev (Google)
Publisher: Google / Chromium
URL: https://web.dev/articles/samesite-cookies-explained
Published: 2019-05-07
Accessed: 2026-09-28
Source Tier: Tier 1 (dokumentasi Chrome/Chromium)
Relevance: Default SameSite=Lax di Chrome 80+, definisi SameSite=Strict/Lax/None, cookie recipes, keterbatasan SameSite.

## Source 8

Title: What is CORS (cross-origin resource sharing)? | PortSwigger Web Security Academy
Publisher: PortSwigger
URL: https://portswigger.net/web-security/cors
Published: 2026 (active page)
Accessed: 2026-09-28
Source Tier: Tier 2 (publikasi keamanan industri)
Relevance: "CORS is not a protection against CSRF", CORS sebagai browser mechanism, server-side security bukan pengganti CORS, kerentanan misconfiguration CORS.

## Source 9

Title: What is CSRF (Cross-site request forgery)? | PortSwigger Web Security Academy
Publisher: PortSwigger
URL: https://portswigger.net/web-security/csrf
Published: 2026 (active page)
Accessed: 2026-09-28
Source Tier: Tier 2
Relevance: Kondisi CSRF, cookie-based session handling, CSRF token, SameSite, Referer-based defense, Chrome Lax-by-default sejak 2021.

## Source 10

Title: Samesite cookies explained (Chrome Lax-by-default announcement context)
Publisher: web.dev / Chromium
URL: https://web.dev/articles/samesite-cookies-explained#changes-to-the-default-behavior-without-samesite
Published: 2019-05-07
Accessed: 2026-09-28
Source Tier: Tier 1
Relevance: Chrome 80+ mengubah default tanpa SameSite attribute menjadi Lax, Edge 86, referensi Incrementally Better Cookies IETF draft.
