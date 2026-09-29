# Content Brief

Topic: Mengamankan Backend dari Eksploitasi Browser dan Serangan Cross-Origin (CORS & CSRF)
Target Reader: Backend engineers, security engineers, API architects, dan full-stack developers.
Problem: Miskonsepsi luas bahwa CORS adalah firewall/mekanisme keamanan backend untuk mencegah unauthorized write atau CSRF, serta ketidaktahuan cara membangun arsitektur pertahanan anti-CSRF multi-layer yang tepat.
Core Mental Model:
- Same-Origin Policy (SOP) membatasi akses pembacaan (*read*) respons lintas origin di level browser, bukan membatasi pengiriman (*write*).
- CORS adalah mekanisme relaksasi SOP agar client JavaScript diizinkan membaca respons; CORS bukan pelindung mutasi data di sisi server.
- CSRF mengeksploitasi mekanisme browser yang otomatis melampirkan kredensial (cookie sesi) saat mengirim request lintas origin.
- Pertahanan CSRF efektif wajib diterapkan di sisi backend melalui kombinasi HMAC-SHA256 Signed Double-Submit Token, Fetch Metadata (`Sec-Fetch-Site`), Custom Header validation, dan atribut cookie `SameSite`.

Approved Research Status: APPROVED (Research Audit: 0 blocking issues, 10 major claims verified).
Approved Engineering Status: APPROVED (Engineering Audit & OpenSource Audit: 0 failures, 100% tests pass with race detector).

Main Concepts:
- Same-Origin Policy (SOP) boundary.
- Simple Request vs Preflight `OPTIONS` (WHATWG Fetch standard).
- Anatomy of Cross-Site Request Forgery (CSRF).
- Misconceptions around `Access-Control-Allow-Origin: *` & CORS vs CSRF.
- Signed Double-Submit CSRF Token (HMAC-SHA256 bound to Session ID & Expiry).
- Fetch Metadata policy (`Sec-Fetch-Site`).
- Custom Header enforcement (`X-Requested-With` / `X-CSRF-Protection`).
- Cookie `SameSite` behavior (`Strict`, `Lax`, `None`).

Verified Behaviors:
1. Cross-origin simple POST request dieksekusi oleh server dan berhasil memutasi state/saldo akun meskipun origin pemanggil tidak diizinkan oleh konfigurasi CORS.
2. CORS middleware menolak preflight request (`OPTIONS`) dari origin yang tidak terdaftar dengan status HTTP 403 Forbidden.
3. CORS middleware mengembalikan header `Access-Control-Allow-Origin`, `Access-Control-Allow-Credentials`, `Access-Control-Allow-Methods`, dan `Access-Control-Allow-Headers` yang sesuai untuk origin yang terdaftar.
4. Protected transfer endpoint menolak request CSRF tanpa token atau dengan invalid HMAC signature (HTTP 403 Forbidden) sehingga saldo korban tetap aman.
5. HMAC token yang valid untuk Session A ditolak saat digunakan oleh Session B (mencegah cross-session token reuse).
6. Expired HMAC token otomatis ditolak oleh validator CSRF.
7. Middleware Fetch Metadata memblokir request non-safe method (POST) jika `Sec-Fetch-Site: cross-site` dan mengizinkan `Sec-Fetch-Site: same-origin`.
8. Middleware Custom Header memblokir request form-based HTML tanpa header custom kustom (`X-Requested-With`) dan meloloskan request AJAX/API valid.
9. Token generator dan validator aman dari race condition (teruji via `go test -race`).

Available Case Studies:
- Simulasi Vulnerable Banking Service (`/api/vulnerable-transfer`) dieksploitasi oleh Attacker Site (`http://attacker.com`) via Cookie-authenticated POST.
- Migrasi dan perlindungan Banking Service menggunakan Multi-layer Defense (`/api/protected-transfer` dengan Signed Token, Fetch Metadata, dan Custom Header).

Warnings:
- Lingkup SOP, CORS, dan CSRF murni berada pada konteks Web Browser. HTTP client non-browser (cURL, script Python, Postman) tidak dibatasi oleh SOP/CORS.
- Mitigasi CSRF (termasuk Anti-CSRF Token, SameSite, dan Fetch Metadata) tidak dapat melindungi aplikasi jika aplikasi tersebut memiliki kerentanan Cross-Site Scripting (XSS) pada origin yang sama.
- Implementasi Signed Double-Submit Token pada lab menggunakan delimiter internal dan secret key demo; pada produksi pastikan secret key dikelola melalui KMS/secret manager serta sanitasi input session.
- `SameSite=Lax` default browser membantu mengurangi risiko, namun navigasi top-level `GET` tetap mengirim cookie, sehingga endpoint `GET` tidak boleh melakukan mutasi state (*must be safe/idempotent*).
