# 04 — Kontradiksi

## Contradiction 1: Perilaku SameSite=Lax Default Chrome vs Explisit SameSite=Lax

**SOURCE A (web.dev / Chromium):** "Caution: Chrome's default behavior is slightly more permissive than an explicit SameSite=Lax, because it lets sites send some cookies on top-level POST requests. ... This is intended as a temporary mitigation."

**SOURCE B (MDN Set-Cookie):** "Some browsers use Lax as the default value if SameSite is not specified. Note: When Lax is applied as a default, a more permissive version is used. In this more permissive version, cookies are also included in POST requests, as long as they were set no more than two minutes before the request was made."

**SOURCE C (PortSwigger CSRF):** "Since 2021, Chrome enforces Lax SameSite restrictions by default." (Implies standard Lax behavior)

**ASSESSMENT:** Tidak ada kontradiksi material — ketiga sumber sepakat bahwa default Lax Chrome 80+ lebih longgar dari eksplisit SameSite=Lax. PortSwigger merujuk pada standar proposal ("Lax-by-default"), sedangkan web.dev dan MDN mendokumentasikan implementasi aktual Chrome (Lax-with-2-minute-POST-grace). **Tidak ada kontradiksi yang tersembunyi** — semuanya jelas tentang perilaku "lebih permisif".

---

## Contradiction 2: Apakah CORS Preflight Mengirim Credentials?

**SOURCE A (MDN CORS):** "CORS-preflight requests must never include credentials. The response to a preflight request must specify Access-Control-Allow-Credentials: true to indicate that the actual request can be made with credentials."

**SOURCE B (Fetch Standard §3.3.3):** "For a CORS-preflight request, request's credentials mode is always 'same-origin', i.e., it excludes credentials, but for any subsequent CORS requests it might not be."

**SOURCE C (MDN CORS - Enterprise note):** "Some enterprise authentication services require that TLS client certificates be sent in preflight requests, in contravention of the Fetch specification. Firefox 87 allows this non-compliant behavior... Chromium-based browsers currently always send TLS client certificates in CORS preflight requests."

**ASSESSMENT:** Standar (Fetch, MDN) menyatakan preflight tidak boleh mengirim credentials. Namun Chromium-based browser mengirim TLS client certificates di preflight (bug 775438). Firefox menonaktifkan ini secara default tapi bisa diaktifkan. **Kontradiksi praktis antara standar dan implementasi Chromium** — untuk TLS client certificates saja, bukan cookies. Perlu catatan di laporan.

---

## Contradiction 3: Apakah CORS Bisa Melindungi dari CSRF?

**SOURCE A (Topic Spec / Lab):** "Mengapa setting Access-Control-Allow-Origin: * tidak mencegah serangan CSRF, dan mengapa CORS bukan mekanisme untuk melindungi data di server?"

**SOURCE B (PortSwigger CORS):** "CORS is not a protection against cross-origin attacks such as cross-site request forgery (CSRF)." "CORS defines browser behaviors and is never a replacement for server-side protection of sensitive data."

**SOURCE C (MDN CORS - Simple Requests section):** "The motivation is that the <form> element from HTML 4.0 can submit simple requests to any origin, so anyone writing a server must already be protecting against cross-site request forgery (CSRF)."

**SOURCE D (OWASP CSRF Prevention - Custom Headers + CORS):** "This defense relies on the CORS preflight mechanism which sends an OPTIONS request to verify CORS compliance... When the API verifies that the custom header is there, you know that the request must have been preflighted if it came from a browser."

**ASSESSMENT:** **Tidak ada kontradiksi.** Semua sumber sepakat: CORS sendiri bukan proteksi CSRF. Namun, *membuat request menjadi non-simple* (via custom header) memaksa preflight — dan server yang *tidak* mengizinkan CORS origin attacker akan memblokir preflight. Ini adalah *pemanfaatan side-effect CORS preflight* sebagai defense-in-depth, bukan CORS itu sendiri melindungi CSRF. OWASP menjelaskan ini dengan benar: defense "mengandalkan CORS preflight mechanism". Harus dibedakan dengan tegas di laporan.

---

## Contradiction 4: SameSite Protection Level — Origin vs Site

**SOURCE A (MDN CSRF):** "Another problem with the SameSite attribute is that it protects you from requests from a different site, not a different origin. This is a looser protection, because (for example) https://foo.example.org and https://bar.example.org are considered the same site, although they are different origins."

**SOURCE B (OWASP CSRF Prevention):** "SameSite can be used for session cookies but be careful to NOT set a cookie specifically for a domain. This action introduces a security vulnerability because all subdomains of that domain will share the cookie, and this is particularly an issue if a subdomain has a CNAME to domains not in your control."

**SOURCE C (PortSwigger CSRF - Bypassing SameSite):** "Bypassing restrictions via vulnerable sibling domains" — attacker bisa kompromikan subdomain lain untuk mengirim cookie.

**ASSESSMENT:** **Konsisten.** Semua sumber setuju SameSite beroperasi pada level "site" (registrable domain + suffix), bukan "origin" (scheme+host+port). Ini adalah *limitasi desain*, bukan kontradiksi.

---

## Contradiction 5: Double-Submit Cookie — Naive vs Signed

**SOURCE A (OWASP):** "The Naive Double-Submit Cookie pattern is bypassable... For new code, use the Signed Double-Submit Cookie pattern above. The naive pattern is documented for reference only." (DISCOURAGED)

**SOURCE B (Topic Spec / Lab):** "Anti-CSRF Token (Double Submit Cookie)" — menyebut sebagai pertahanan modern tanpa membedakan naive vs signed.

**ASSESSMENT:** **Perbedaan granularitas.** Lab spec menyebut "Double Submit Cookie" sebagai kategori, sedangkan OWASP membedakan naive (DISCOURAGED) vs signed (RECOMMENDED). Ini **bukan kontradiksi faktual**, tapi perbedaan presisi terminologi. Di laporan perlu memisahkan secara eksplisit.

---

## No material contradictions discovered beyond the documented differences above.

All significant claims are corroborated by multiple Tier 1 sources. Minor implementation differences (Chrome Lax default permissiveness, TLS certs in preflight) are explicitly noted and sourced.