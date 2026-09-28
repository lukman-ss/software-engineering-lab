# 05 — Laporan Riset: CORS & CSRF — Mengamankan Backend dari Eksploitasi Browser dan Cross-Origin Attacks

## Research Question

Menjawab dua pertanyaan mendasar dari spesifikasi Lab #36:
1. Mengapa `Access-Control-Allow-Origin: *` tidak mencegah serangan CSRF, dan mengapa CORS bukan mekanisme untuk melindungi data di server?
2. Bagaimana Same-Origin Policy (SOP), CORS, dan CSRF Protection saling melengkapi dan bagaimana pertahanan modern (SameSite, Anti-CSRF Token, Fetch Metadata) bekerja?

---

## Executive Summary

Riset ini memvalidasi dasar teoretis dan teknis untuk modul lab `labs/36-cors-and-csrf`. Hasil riset mengonfirmasi:

1. **CORS adalah kebijakan relaksasi SOP di level browser, BUKAN firewall/pengaman backend**: Browser tetap mengirimkan permintaan (*write requests* seperti form `POST`), dan backend tetap mengeksekusinya, terlepas dari apakah header CORS cocok atau ditolak oleh browser. CORS hanya mengontrol apakah JavaScript pada origin pemanggil diizinkan *membaca respons*.
2. **`Access-Control-Allow-Origin: *` tidak relevan terhadap CSRF**: Serangan CSRF mengandalkan pengiriman permintaan yang mengeksekusi aksi di server menggunakan cookie otentikasi otomatis milik korban. Attacker tidak perlu membaca isi respons untuk melancarkan CSRF. Selain itu, browser memblokir pembacaan respons ber-kredensial jika ACAO disetel `*`.
3. **Simple vs Non-Simple Requests**: Permintaan menggunakan method `GET`/`HEAD`/`POST` dengan `Content-Type` formulir standar (`application/x-www-form-urlencoded`, `multipart/form-data`, `text/plain`) lolos tanpa preflight (*simple request*). Permintaan dengan method selain itu atau dengan header kustom akan memicu preflight `OPTIONS`.
4. **Pertahanan CSRF modern multi-layer**:
   - `SameSite=Lax` (default browser modern seperti Chrome 80+) memblokir pengiriman cookie pada cross-site `POST`, namun tetap mengirim pada top-level navigasi `GET`.
   - `SameSite=Strict` memblokir pengiriman cookie pada seluruh cross-site request.
   - **Anti-CSRF Token** (Synchronizer Token Pattern / Signed Double-Submit Cookie) wajib untuk operasi yang memerlukan integritas tinggi dan backward compatibility.
   - **Fetch Metadata (`Sec-Fetch-Site`)** dan **Custom Headers** memberikan perlindungan modern berbasis API.

---

## Findings

### Finding 1: Perbedaan Fundamental SOP, CORS, dan CSRF

Claim: Same-Origin Policy (SOP) membatasi pembacaan respons lintas origin tetapi secara historis mengizinkan penulisan (form submission); CORS adalah mekanisme untuk melonggarkan pembatasan pembacaan respons tersebut; CSRF mengeksploitasi izin penulisan lintas origin.

Evidence:
- **MDN SOP**: "Cross-origin writes are typically allowed... Cross-origin reads are typically disallowed... Use CORS to allow cross-origin access... To prevent cross-origin writes, check an unguessable token (CSRF token)."
- **PortSwigger CORS**: "The same-origin policy generally allows a domain to issue requests to other domains, but not to access the responses... CORS is not a protection against cross-origin attacks such as cross-site request forgery (CSRF)."
- **Fetch Standard §3.3**: CORS checks hanya menentukan apakah suatu respons HTTP boleh di-share ke JavaScript di origin pemanggil.

Sources:
- MDN Same-origin policy: https://developer.mozilla.org/en-US/docs/Web/Security/Defenses/Same-origin_policy
- PortSwigger CORS: https://portswigger.net/web-security/cors
- WHATWG Fetch Standard: https://fetch.spec.whatwg.org/

Confidence: HIGH

---

### Finding 2: CORS Bukan Pelindung Sisi Server dan ACAO: * Tidak Mencegah CSRF

Claim: Menyetel header CORS (termasuk `Access-Control-Allow-Origin: *` atau menolak Origin) tidak mencegah server mengeksekusi request berbahaya karena eksekusi terjadi sebelum/tanpa bergantung pada evaluasi CORS di browser.

Evidence:
- **MDN CSRF**: Menjelaskan alur CSRF di mana form `POST` dikirim langsung ke server target; server mengeksekusi perubahan state dan mengembalikan respons. Evaluasi CORS di sisi browser hanya terjadi jika request dibuat via XHR/fetch dan bertujuan membaca respons.
- **MDN CORS (Credentialed requests and wildcards)**: "If a request includes a credential (most commonly a Cookie header) and the response includes an Access-Control-Allow-Origin: * header, the browser will block access to the response." Attacker CSRF tidak peduli jika respons diblokir; mutasi data di database server sudah selesai dilakukan.
- **PortSwigger CORS**: "CORS defines browser behaviors and is never a replacement for server-side protection of sensitive data — an attacker can directly forge a request from any trusted origin."

Sources:
- MDN CSRF: https://developer.mozilla.org/en-US/docs/Web/Security/Attacks/CSRF
- MDN CORS: https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/CORS#credentialed_requests_and_wildcards
- PortSwigger CORS: https://portswigger.net/web-security/cors

Confidence: HIGH

---

### Finding 3: Mekanisme Preflight OPTIONS dan Simple Requests

Claim: Browser hanya mengirimkan preflight `OPTIONS` untuk request non-simple (menggunakan method di luar GET/HEAD/POST, atau custom headers, atau Content-Type selain yang disafelist). Simple request langsung dikirim ke server.

Evidence:
- **WHATWG Fetch Standard §2.2.1, §2.2.2**: CORS-safelisted methods adalah `GET`, `HEAD`, `POST`. CORS-safelisted headers mencakup `Accept`, `Accept-Language`, `Content-Language`, `Content-Type` (hanya `application/x-www-form-urlencoded`, `multipart/form-data`, `text/plain`), serta `Range`.
- **MDN CORS (Preflighted requests)**: "Unlike simple requests, for 'preflighted' requests the browser first sends an HTTP request using the OPTIONS method... in order to determine if the actual request is safe to send."
- **OWASP CSRF Prevention**: Mengonfirmasi bahwa request dengan custom header otomatis memicu preflight, sehingga API berbasis JSON yang menolak Simple Content-Type atau mewajibkan custom header terlindungi dari exploit via `<form>`.

Sources:
- WHATWG Fetch Standard: https://fetch.spec.whatwg.org/#cors-preflight-fetch
- MDN CORS: https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/CORS#preflighted_requests
- OWASP CSRF Prevention: https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html

Confidence: HIGH

---

### Finding 4: SameSite Cookie Attribute (Strict, Lax, None) dan Batasannya

Claim: Atribut `SameSite` mengontrol pengiriman cookie pada request lintas situs. `SameSite=Lax` mengizinkan cookie pada navigasi top-level `GET` dan memblokirnya pada `POST`; `SameSite=Strict` memblokir semua request lintas situs; namun SameSite beroperasi pada level *site* (eTLD+1), bukan *origin*, sehingga subdomain yang rentan dapat membypass perlindungan ini.

Evidence:
- **MDN Set-Cookie**:
  - `Strict`: Cookie hanya dikirim untuk request dari situs yang sama persis.
  - `Lax`: Cookie dikirim untuk same-site dan cross-site top-level navigation dengan safe method (GET).
  - `None`: Dikirim pada semua request, wajib disertai atribut `Secure`.
- **web.dev (Google/Chromium)**: Chrome 80+ memberlakukan default `SameSite=Lax` jika atribut tidak disetel.
- **OWASP & PortSwigger CSRF**: Menyoroti bahwa SameSite tidak melindungi dari cross-origin same-site attacks (misalnya subfolder atau sibling domain `attacker.example.com` terhadap `app.example.com`).

Sources:
- MDN Set-Cookie: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Set-Cookie#samesitesamesite-value
- web.dev SameSite Explained: https://web.dev/articles/samesite-cookies-explained
- OWASP CSRF Prevention: https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html
- PortSwigger CSRF: https://portswigger.net/web-security/csrf/bypassing-samesite-restrictions

Confidence: HIGH

---

### Finding 5: Pola Anti-CSRF Token: Synchronizer Token vs Double-Submit Cookie

Claim: Synchronizer Token Pattern (server-side session bound) adalah standar baku. Jika menggunakan Double Submit Cookie (stateless), varian Naive tidak aman dan harus menggunakan Signed Double-Submit Cookie (HMAC bound to session).

Evidence:
- **OWASP CSRF Prevention**:
  - *Synchronizer Token Pattern*: Token unik, rahasia, kriptografis acak (CSPRNG), disimpan di sesi server dan disisipkan di form/header.
  - *Signed Double-Submit Cookie (RECOMMENDED)*: Token di-generate menggunakan HMAC dengan rahasia server dan session ID: `HMAC(secret, sessionID + randomValue)`.
  - *Naive Double-Submit Cookie (DISCOURAGED)*: "The Naive Double-Submit Cookie pattern is bypassable by an attacker who can write cookies on the target domain (e.g., via a vulnerable sibling subdomain, DNS takeover, or plaintext-HTTP cookie injection)."
- **MDN CSRF**: Menjelaskan framework modern (seperti Django) menggunakan embedded token tag untuk memvalidasi form `POST`.

Sources:
- OWASP CSRF Prevention Cheat Sheet: https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html#alternative-using-a-double-submit-cookie-pattern
- MDN CSRF: https://developer.mozilla.org/en-US/docs/Web/Security/Attacks/CSRF#csrf_tokens

Confidence: HIGH

---

### Finding 6: Pertahanan Berbasis Fetch Metadata (`Sec-Fetch-Site`)

Claim: Header Fetch Metadata (`Sec-Fetch-Site`) memungkinkan backend secara langsung mengetahui konteks inisiasi request (`same-origin`, `same-site`, `cross-site`, `none`) dan memblokir cross-site state-changing request tanpa memerlukan token di sisi client.

Evidence:
- **OWASP CSRF Prevention**: "Sec-Fetch-Site — the primary signal for CSRF protection... treat cross-site as untrusted for state-changing actions (reject non-safe methods POST/PUT/DELETE when `Sec-Fetch-Site: cross-site`)."
- **MDN CSRF**: Mendokumentasikan middleware Express yang memeriksa `req.headers["sec-fetch-site"]`.
- **W3C Fetch Metadata Spec**: Menjamin header `Sec-Fetch-*` di-set oleh user agent dan tidak dapat dimanipulasi oleh JavaScript frontend.

Sources:
- OWASP CSRF Prevention: https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html#fetch-metadata-headers
- MDN CSRF: https://developer.mozilla.org/en-US/docs/Web/Security/Attacks/CSRF#fetch_metadata

Confidence: HIGH

---

## Areas of Agreement

Semua sumber (MDN, WHATWG, OWASP, PortSwigger, Google web.dev) secara bulat menyepakati:
1. CORS bukan mekanisme pertahanan server-side dan tidak mencegah eksekusi request di backend.
2. Wildcard `Access-Control-Allow-Origin: *` otomatis ditolak oleh browser jika kredensial disertakan (`credentials: include`).
3. Form HTML submission menghasilkan Simple Request yang tidak memicu preflight `OPTIONS`.
4. Atribut cookie `SameSite=Lax` secara signifikan mengurangi vektor CSRF default, namun pertahanan komprehensif tetap membutuhkan Anti-CSRF Token atau Fetch Metadata.
5. Permintaan `GET` tidak boleh digunakan untuk operasi perubahan state (state-changing actions).

---

## Areas of Disagreement / Nuances

1. **Perilaku Default Lax di Chromium**: Standard RFC 6265bis mengusulkan Strict Lax-by-default, namun Chromium mengimplementasikan "Lax + POST 2-minute grace period" (Blink-dev announcement) untuk kompatibilitas SSO. Hal ini tidak mengubah kesimpulan keamanan bahwa `SameSite` eksplisit tetap wajib.
2. **Kirim TLS Client Cert di Preflight**: WHATWG Fetch melarang credential di preflight, namun Chromium selalu melampirkan TLS client certificate pada preflight (Chrome bug 775438).

---

## Limitations

1. **Lingkup Browser-Only**: Pembahasan SOP, CORS, dan CSRF berlaku khusus pada konteks Web Browser. Client non-browser (cURL, Postman, backend-to-backend API calls) tidak mematuhi SOP maupun CORS.
2. **Ketergantungan XSS**: Semua mitigasi CSRF (termasuk Anti-CSRF Token, SameSite, dan Fetch Metadata) dapat dilumpuhkan jika target memiliki kerentanan Cross-Site Scripting (XSS) pada origin yang sama.

---

## Conclusion

Spesifikasi Topik Lab #36 memiliki landasan teoretis dan standar industri yang sangat kuat dan valid:
1. Pernyataan *"CORS itu fitur security backend"* adalah **Mitos/Salah**; CORS adalah relaksasi SOP di browser untuk membaca respons.
2. Pernyataan *"Kalau sudah pasang CORS, pasti bebas dari CSRF"* adalah **Mitos/Salah**; request simple tetap dikirim dan dieksekusi server.
3. Desain kurikulum lab yang mempraktikkan simulasi Cookie Auth + Attacker Site, konfigurasi Preflight OPTIONS, penerapan `SameSite=Lax/Strict`, serta Anti-CSRF Token (Signed Double Submit) sepenuhnya selaras dengan standar OWASP dan Fetch Standard.
