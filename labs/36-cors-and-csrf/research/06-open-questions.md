# 06 — Pertanyaan Terbuka

## Unanswered Questions & Weak Evidence

1. **Detail implementasi preflight per-browser (Chrome, Safari, Firefox)** — Berapa persentase pengguna aktual yang menjalankan browser versi yang mengimplementasikan `SameSite=Lax-by-default` secara penuh vs yang masih menggunakan perilaku legacy? Perlu data statistik dari MDN Browser Compatibility atau W3Counter.

2. **Perilaku edge-case preflight dengan redirect** — Fetch spec sudah tidak mewajibkan browser memblokir redirect setelah preflight, tetapi beberapa browser masih melakukannya. Sumber primer mengenai persentase browser yang sudah mengimplementasikan perubahan spec ini belum ditemukan dengan kuat.

3. **Perilaku SameSite di mobile webviews (Android/iOS)** — Bagaimana WebView aplikasi mobile memperlakukan SameSite dan CORS? Ini relevan untuk lab yang juga menangani mobile.

4. **Kapan persis Chrome 80 mengubah default SameSite?** — Web.dev menyebutkan Chrome 80 dan Edge 86, namun tanggal rilis persis belum dicatat dalam sumber yang dibaca. Dapat ditelusuri ke Chromium release notes.

5. **Apakah `document.domain` masih relevan untuk SOP modern?** — MDN menyatakan mekanisme ini deprecated. Berapa proporsi situs yang masih memakainya dan apakah ini masih menimbulkan kerentanan?

6. **Statistik eksploitasi CSRF di dunia nyata** — Data empiris tentang seberapa umum CSRF masih terjadi meskipun SameSite Lax default, persentase situs yang masih menggunakan naive double-submit cookie, dsb. Belum ditemukan di sumber primer (butuh data dari CVE/NVD, SANS, atau OWASP Top 10 terbaru).

7. **Apakah Fetch Metadata headers dapat dipalsukan oleh proxy/cache?** — OWASP mencatat: "Intermediaries (proxies, gateways, load balancers) may remove or modify Origin and Sec-* headers." Detail teknis seberapa sering ini terjadi dan bagaimana mitigasinya perlu penelitian lebih lanjut.

## Claims Needing Deeper Research

- Klaim tentang "Chrome Lax-by-default since 2021" di PortSwigger: web.dev menyebut Chrome 80 (2020), bukan 2021. Perlu verifikasi kronologi yang tepat.
- Dampak spesifik dari "POST 2-minute grace period" Chrome terhadap keamanan sesi login dan CSRF — belum ada analisis dampak keamanan yang mendalam dari sumber primer.

## Possible Next Research Directions

1. **Audit implementasi CORS** — Riset terhadap N=50 website populer (dari Alexa/Chrome-Stats) untuk menentukan seberapa sering `Access-Control-Allow-Origin: *` digunakan secara tidak aman dengan credential.
2. **Perbandingan kebijakan cookie SameSite lintas negara/regulasi** — Bagaimana GDPR/ePrivacy atau regulasi lain berinteraksi dengan kebijakan SameSite dan cookie tracking.
3. **Riset akademik terbaru** — Cari paper terbaru (post-2020) tentang XS-Leaks dan bagaimana SameSite/CORS dapat dimanfaatkan sebagai saluran bocornya informasi (side-channel).
4. **Post-quantum considerations for CSRF tokens** — Bagaimana quantum computing memengaruhi keamanan HMAC yang digunakan dalam Signed Double-Submit Cookie (masih theoretical).
5. **Integrasi CORS dengan CSP (Content Security Policy)** — Bagaimana CSP directives berinteraksi dengan CORS headers dan apakah kombinasi keduanya memberikan perlindungan lebih kuat terhadap CSRF dan data exfiltration.
