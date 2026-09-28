# 01 — Riset Plan

## Research Topic

CORS & CSRF — Mengamankan Backend dari Eksploitasi Browser dan Cross-Origin Attacks
(Target lab: `labs/36-cors-and-csrf`, lab #36, kategori Backend Engineering / Security & Web Fundamentals)

## Objective

Menyediakan bukti terverifikasi (bukan opini) untuk spesifikasi topik lab #36, khususnya
menguji dua klaim inti:

1. `Access-Control-Allow-Origin: *` tidak mencegah serangan CSRF, dan CORS bukan mekanisme
   pelindung data di server.
2. Peran Same-Origin Policy (SOP) sebagai dasar, serta pemisahan peran CORS (kebijakan
   akses browser) vs CSRF Protection (integritas permintaan).

Riset ini **hanya mengumpulkan dan memverifikasi bukti**. Implementasi kode, test, dan
konten publikasi ditangani agen lain.

## Research Questions

- **RQ1**: Apa yang diatur Same-Origin Policy (SOP), dan apa batasnya terhadap pengiriman
  permintaan (bukan pembacaan respons)?
- **RQ2**: Apa yang dilakukan CORS terhadap respons, dan siapa yang menegakkannya
  (browser vs server)?
- **RQ3**: Kapan browser mengirim Preflight `OPTIONS`, dan method/header apa yang termasuk
  "simple request" (lolos preflight)?
- **RQ4**: Apakah CORS (dengan `Access-Control-Allow-Origin: *` maupun tanpa CORS) mencegah
  CSRF? Apakah server tetap mengeksekusi request yang diblokir browser?
- **RQ5**: Bagaimana mekanisme `SameSite` cookie (`Strict`, `Lax`, `None`) bekerja, termasuk
  perilaku default dan perilaku navigasi top-level?
- **RQ6**: Bagaimana Anti-CSRF Token (termasuk Double Submit Cookie) bekerja, dan syarat
  keamanannya?
- **RQ7**: Apa miskomunisi umum yang perlu diluruskan di lab (mitos vs fakta)?

## Search Strategy

1. **Tier 1 (primer/standar)** — sumber rujukan utama:
   - WHATWG Fetch Standard (definisi CORS, preflight, safelisted method/header)
   - MDN Web Docs (SOP, CORS, CSRF, Cookie SameSite) — dokumentasi otoritatif Mozilla
   - RFC 6265bis / RFC 6265 (cookie, atribut SameSite)
   - OWASP Cheat Sheet Series (CSRF Prevention, CORS)
   - PortSwigger Web Security Academy (laboratorium + penjelasan teknis)
2. **Tier 2** — literatur akademik & publikasi teknis terkemuka:
   - Barth, Jackson, Mitchell — "Robust Defenses for Cross-Site Request Forgery" (SOSP 2008)
   - Sotirov & Bowden — "The Most Dangerous Code in the World" (USENIX Security 2008)
   - Chromium Project / blog Google (perilaku default SameSite Lax-by-default)
3. Verifikasi: setiap klaim penting dibuka langsung (bukan dari snippet hasil pencarian),
   lalu dicocokkan ke sumber kedua yang independen.
4. Bila bukti tidak ditemukan: tulis **NOT VERIFIED** (tanpa mengarang URL/tanggal/statistik).

## Expected Primary Sources

- https://fetch.spec.whatwg.org/ (CORS protocol, preflight, safelisted methods)
- https://developer.mozilla.org/en-US/docs/Web/Security/Same-origin_policy
- https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS
- https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Set-Cookie#samesitesamesite-value
- https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html
- https://portswigger.net/web-security/cors , https://portswigger.net/web-security/csrf
- RFC terkait cookie (6265 / 6265bis)
- Barth et al. 2008 (SOSP), Sotirov & Bowden 2008 (USENIX)

## Risks / Unknowns

- Perilaku default `SameSite` berubah antar versi browser (awalnya tidak ada atribut,
  lalu Lax-by-default di Chrome 80+, dsb.) → perlu tanggal & sumber primer perubahan.
- Definisi "simple request" terkait spesifikasi yang masih hidup (Fetch) dan dapat berubah
  → catat versi/tanggal akses.
- Perbedaan terminologi: CORS *memblokir aksesibilitas respons di JS*, bukan mencegah
  request dikirim → sumber harus eksplisit soal "server tetap memproses request".
- Some lab claims (mis. angka atau contoh spesifik) tidak punya data empiris → tandai
  NOT VERIFIED bila tidak ada sumber.
- Bahasa laporan: mengikuti bahasa dominan instruksi (Inggris spesifikasi) dengan
  istilah teknis dipertahankan apa adanya; target lab berbahasa Indonesia tetapi artefak
  riset mengikuti struktur template laporan.
