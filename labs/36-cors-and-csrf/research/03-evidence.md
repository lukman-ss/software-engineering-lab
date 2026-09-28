# 03 — Bukti Riset

## Evidence 1

Claim: Same-Origin Policy (SOP) mengizinkan cross-origin writes (termasuk form submissions) tetapi memblokir cross-origin reads.

Evidence: MDN SOP: "Cross-origin writes are typically allowed. Examples are links, redirects, and form submissions." dan "Cross-origin reads are typically disallowed." Mekanisme pertahanan tulis lintas origin adalah CSRF token.

Source: MDN Same-origin policy

URL: https://developer.mozilla.org/en-US/docs/Web/Security/Defenses/Same-origin_policy

Confidence: HIGH

Corroborated By: MDN CSRF page, MDN CORS page, PortSwigger CORS page.

Notes: SOP memang dirancang agar form submission silang origin tetap berfungsi — ini alasan mengapa simple requests tidak memerlukan preflight.

---

## Evidence 2

Claim: CORS adalah mekanisme browser (bukan server) yang mengontrol apakah JavaScript dapat membaca respons lintas origin. Server tetap mengeksekusi request walaupun browser memblokir respons.

Evidence: MDN CORS: "CORS is an HTTP-header based mechanism that allows a server to indicate any origins other than its own from which a browser should permit loading resources." PortSwigger: "Cross-origin resource sharing (CORS) is a browser mechanism which enables controlled access to resources located outside of a given domain." MDN SOP: "To prevent cross-origin writes, check an unguessable token in the request — known as a CSRF token."

Source: MDN CORS, PortSwigger CORS, MDN SOP

URL: https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/CORS, https://portswigger.net/web-security/cors, https://developer.mozilla.org/en-US/docs/Web/Security/Defenses/Same-origin_policy

Confidence: HIGH

Corroborated By: Fetch Standard (CORS check hanya memengaruhi aksesibilitas respons di JavaScript, bukan pengiriman request).

Notes: Fetch spec mendefinisikan CORS check sebagai langkah yang menentukan apakah respons "shared" — request tetap dikirim ke server, server tetap memproses.

---

## Evidence 3

Claim: `Access-Control-Allow-Origin: *` tidak dapat digunakan bersama credentials (cookies). Browser memblokir akses respons ketika credential dikirim bersama wildcard ACAO.

Evidence: MDN CORS: "If a request includes a credential (most commonly a Cookie header) and the response includes an Access-Control-Allow-Origin: * header (that is, with the wildcard), the browser will block access to the response, and report a CORS error in the devtools console." Fetch spec: Access-Control-Allow-Origin cannot be `*` when request's credentials mode is "include".

Source: MDN CORS, Fetch Standard

URL: https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/CORS#credentialed_requests_and_wildcards, https://fetch.spec.whatwg.org/

Confidence: HIGH

Corroborated By: PortSwigger CORS ("server must not specify the * wildcard for Access-Control-Allow-Origin response-header value when Access-Control-Allow-Credentials: true").

Notes: Fakta kunci: wildcard ACAO memang mengizinkan siapa saja mengakses resource, TAPI tidak dengan cookies/credentials. Ini berarti ACAO: * tidak memberi akses penuh lintas origin dengan session cookies.

---

## Evidence 4

Claim: Simple requests (GET, HEAD, POST dengan content-type form/text) lolos preflight dan dikirim oleh browser ke server mana pun. Browser tidak memblokir pengiriman request ini.

Evidence: MDN CORS: "A simple request is one that meets all the following conditions: One of the allowed methods: GET, HEAD, POST." "The motivation is that the <form> element from HTML 4.0 can submit simple requests to any origin." Fetch spec: "A CORS-safelisted method is a method that is GET, HEAD, or POST."

Source: MDN CORS, Fetch Standard

URL: https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/CORS#simple_requests, https://fetch.spec.whatwg.org/

Confidence: HIGH

Corroborated By: MDN CSRF page (menjelaskan POST form sebagai simple request yang dapat memicu CSRF).

Notes: Ini alasan CSRF tetap mungkin meskipun CORS dikonfigurasi — simple requests tidak memicu preflight.

---

## Evidence 5

Claim: Preflight request (OPTIONS) dikirim oleh browser untuk method/header non-safelisted (PUT, DELETE, custom headers seperti Authorization).

Evidence: Fetch spec: "A CORS-preflight request is a CORS request that checks to see if the CORS protocol is understood. It uses OPTIONS as method and includes Access-Control-Request-Method." MDN CORS: "for HTTP request methods that can cause side-effects on server data (in particular, HTTP methods other than GET, or POST with certain MIME types), the specification mandates that browsers 'preflight' the request."

Source: Fetch Standard, MDN CORS

URL: https://fetch.spec.whatwg.org/#cors-preflight-fetch, https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/CORS#preflighted_requests

Confidence: HIGH

Corroborated By: OWASP CSRF Prevention (menggunakan custom request header sebagai defense — karena custom header memicu preflight).

Notes: OWASP merekomendasikan X-CSRF-Token atau header kustom lainnya — ini secara efektif mengubah semua request menjadi non-simple, sehingga preflight diperlukan.

---

## Evidence 6

Claim: `SameSite=Lax` (default Chrome 80+) mengirim cookie hanya pada navigasi top-level GET, menghentikan sebagian besar serangan CSRF berbasis form POST cross-site.

Evidence: MDN Set-Cookie: "Lax: Send the cookie only for requests originating from the same site that set the cookie, and for cross-site requests that meet both of the following criteria: (1) The request is a top-level navigation. (2) The request uses a safe method: in particular, this excludes POST, PUT, and DELETE." PortSwigger CSRF: "Since 2021, Chrome enforces Lax SameSite restrictions by default." web.dev: "Cookies without a SameSite attribute are treated as SameSite=Lax."

Source: MDN Set-Cookie, PortSwigger CSRF, web.dev

URL: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Set-Cookie, https://portswigger.net/web-security/csrf, https://web.dev/articles/samesite-cookies-explained

Confidence: HIGH

Corroborated By: MDN CSRF page juga menjelaskan perilaku Lax.

Notes: Chrome 80+ (Februari 2020) mengubah default. MDN juga mencatat ada "more permissive version" yang mengizinkan POST ≤ 2 menit setelah cookie diset.

---

## Evidence 7

Claim: `SameSite=Strict` tidak mengirim cookie pada cross-site request apapun, termasuk navigasi top-level. `SameSite=None` mengirim cookie pada semua konteks (memerlukan `Secure`).

Evidence: MDN Set-Cookie: "Strict: Send the cookie only for requests originating from the same site that set the cookie." "None: Send the cookie with both cross-site and same-site requests. The Secure attribute must also be set when using this value." web.dev: "SameSite=Strict ... cookie can only be sent in a first-party context."

Source: MDN Set-Cookie, web.dev

URL: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Set-Cookie, https://web.dev/articles/samesite-cookies-explained

Confidence: HIGH

Corroborated By: OWASP CSRF Prevention (SameSite section).

Notes: SameSite beroperasi pada level "site" (registrable domain), bukan "origin" — subdomain dianggap same-site.

---

## Evidence 8

Claim: Anti-CSRF Token (Synchronizer Token Pattern) adalah mekanisme utama pertahanan CSRF. Token harus unik per sesi, rahasia, dan unpredictable.

Evidence: OWASP: "CSRF tokens should be generated on the server-side and they should be generated only once per user session or each request." "CSRF tokens prevent CSRF because without a CSRF token, an attacker cannot create valid requests to the backend server." MDN CSRF: "the server embeds an unpredictable value in the page, called the CSRF token."

Source: OWASP CSRF Prevention, MDN CSRF

URL: https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html, https://developer.mozilla.org/en-US/docs/Web/Security/Attacks/CSRF

Confidence: HIGH

Corroborated By: PortSwigger CSRF ("CSRF tokens are a unique, secret, and unpredictable value").

Notes: OWASP juga merekomendasikan Signed Double-Submit Cookie (HMAC-bound to session) sebagai alternatif stateless.

---

## Evidence 9

Claim: Double-Submit Cookie pattern (naive) rentan terhadap cookie injection — Signed Double-Submit Cookie (HMAC) direkomendasikan.

Evidence: OWASP: "The Naive Double-Submit Cookie pattern is bypassable by an attacker who can write cookies on the target domain (e.g., via a vulnerable sibling subdomain, DNS takeover, or plaintext-HTTP cookie injection)." "Signed Double-Submit Cookie (RECOMMENDED): Always bind the CSRF token explicitly to session-specific data."

Source: OWASP CSRF Prevention

URL: https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html

Confidence: HIGH

Corroborated By: PortSwigger CSRF (bypass SameSite via sibling domains).

Notes: Double-submit bekerja karena attacker tidak dapat menulis cookie di origin target — tapi ini asumsi yang bisa dilanggar.

---

## Evidence 10

Claim: CORS bukan pengganti keamanan server-side. "CORS defines browser behaviors and is never a replacement for server-side protection of sensitive data."

Evidence: PortSwigger CORS: "CORS is not a protection against cross-origin attacks such as cross-site request forgery (CSRF)." dan "CORS defines browser behaviors and is never a replacement for server-side protection of sensitive data — an attacker can directly forge a request from any trusted origin."

Source: PortSwigger CORS

URL: https://portswigger.net/web-security/cors

Confidence: HIGH

Corroborated By: MDN CORS, Fetch Standard (CORS hanya memengaruhi sharing respons ke JavaScript, bukan pengiriman request).

Notes: Ini adalah poin kritis yang sering disalahpahami — CORS melindungi data respons dari akses JavaScript lintas origin, bukan mencegah server memproses request.

---

## Evidence 11

Claim: Custom request header (misal X-CSRF-Token) memicu preflight, sehingga CSRF dapat dicegah tanpa token — hanya perlu server memeriksa keberadaan header.

Evidence: OWASP: "This defense relies on the CORS preflight mechanism which sends an OPTIONS request to verify CORS compliance with the destination server. All modern browsers designate requests with custom headers as 'to be preflighted'." MDN CSRF: "setting a custom header on the request will prevent it being treated as a simple request."

Source: OWASP CSRF Prevention, MDN CSRF

URL: https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html, https://developer.mozilla.org/en-US/docs/Web/Security/Attacks/CSRF

Confidence: HIGH

Corroborated By: Fetch spec (CORS-safelisted request headers tidak termasuk custom headers).

Notes: Pendekatan ini hanya efektif untuk API yang diakses via fetch/XHR, bukan untuk form submissions tradisional.

---

## Evidence 12

Claim: Fetch Metadata headers (Sec-Fetch-Site) dapat digunakan sebagai mekanisme pertahanan CSRF modern.

Evidence: OWASP: "Fetch Metadata request headers provide extra information about the context from which an HTTP request was made." "Sec-Fetch-Site — the primary signal for CSRF protection. It indicates the relationship between the request initiator's origin and its target's origin: same-origin, same-site, cross-site, or none." MDN CSRF: "Fetch metadata is a collection of HTTP request headers, added by the browser, that provide extra information about the context of an HTTP request."

Source: OWASP CSRF Prevention, MDN CSRF

URL: https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html, https://developer.mozilla.org/en-US/docs/Web/Security/Attacks/CSRF

Confidence: HIGH

Corroborated By: W3C Fetch Metadata specification.

Notes: OWASP mencatat: "a fallback to standard origin verification headers is a mandatory requirement for any Fetch Metadata implementation" karena browser lama tidak mengirim header ini.

---

## Evidence 13

Claim: Kerentanan CORS utama terjadi saat server memantulkan Origin header tanpa validasi (reflects origin in Access-Control-Allow-Origin).

Evidence: PortSwigger CORS: "Some applications take the easy route of effectively allowing access from any other domain. One way to do this is by reading the Origin header from requests and including a response header stating that the requesting origin is allowed."

Source: PortSwigger CORS

URL: https://portswigger.net/web-security/cors

Confidence: HIGH

Corroborated By: OWASP CORS-related misconfigurations.

Notes: Ini adalah salah satu kerentanan CORS paling umum — attacker menggunakan XMLHttpRequest dengan credentials: include untuk membaca data sensitif dari origin target.

---

## Evidence 14

Claim: Kondisi yang diperlukan untuk serangan CSRF: (1) action yang relevan, (2) cookie-based session handling, (3) tidak ada parameter request yang unpredictable.

Evidence: PortSwigger CSRF: "For a CSRF attack to be possible, three key conditions must be in place: A relevant action. Cookie-based session handling. No unpredictable request parameters."

Source: PortSwigger CSRF

URL: https://portswigger.net/web-security/csrf

Confidence: HIGH

Corroborated By: MDN CSRF (kondisi serupa dijelaskan).

Notes: Ketiga kondisi ini harus terpenuhi secara bersamaan.

---

## Evidence 15

Claim: SameSite Lax memiliki keterbatasan: attacker dapat memicu top-level navigation GET, dan beberapa framework mendukung "method override" yang mengubah POST menjadi GET.

Evidence: MDN CSRF: "An attacker can trigger a top-level navigation. For example, the attacker submits a form to the target: this is considered a top-level navigation. If the form were submitted using GET, then the request would still include cookies with SameSite=Lax." "some web frameworks support 'method override': this enables an attacker to send a request using GET but have it appear to the server as if it used POST."

Source: MDN CSRF

URL: https://developer.mozilla.org/en-US/docs/Web/Security/Attacks/CSRF

Confidence: HIGH

Corroborated By: PortSwigger CSRF (bypass SameSite restrictions section).

Notes: OWASP juga mencatat keterbatasan: SameSite beroperasi pada level "site", bukan "origin" — semua subdomain dianggap same-site.
