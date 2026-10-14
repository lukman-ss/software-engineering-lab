# Ringkasan Deliverables Publikasi — Lab 31: OAuth 2.0 & OpenID Connect (OIDC)

Paket artefak publikasi untuk **Lab 31: OAuth 2.0 & OpenID Connect (OIDC)** telah berhasil dibangun sesuai standar mutu teknis dan visual.

---

## 1. Tabel Unduhan & Inventaris Artefak (8 Deliverables Lengkap + Python Source)

| No | Tipe Artefak | Jalur Berkas Relatif | Status & Verifikasi Ukuran | Format / Target Platform |
|:---:|:---|:---|:---|:---|
| **0** | **Slide & Meta Definition** | `result/content_def.py` | Definisi 24 Slide & Metadata Lab | Python Source (DO NOT DELETE) |
| **1** | **Slide Deck (Editable)** | `result/linkedin/lab-31-oauth2-and-oidc-carousel.pptx` | 113.2 KB (24 slide + Speaker Notes) | Microsoft PowerPoint / Keynote (16:9) |
| **2** | **Slide Deck (PDF)** | `result/linkedin/lab-31-oauth2-and-oidc-linkedin.pdf` | 718.0 KB (24 halaman selectable text) | LinkedIn Carousel Document (16:9) |
| **3** | **Artikel Teknis (EN)** | `result/dev.to/lab-31-oauth2-and-oidc-devto.md` | 14.7 KB (1.792 kata, Frontmatter lengkap) | Dev.to Markdown |
| **4** | **Artikel Mendalam (ID)** | `result/medium/lab-31-oauth2-and-oidc-medium.html` | 24.3 KB (2.896 kata, 9 section wajib) | Medium Standalone HTML |
| **5** | **Metadata Artikel** | `result/medium/lab-31-oauth2-and-oidc-metadata.md` | Lengkap (SEO Title, Desc, 5 Tag) | Panduan Pengaturan Penerbitan |
| **6** | **Cover Visual (ID)** | `result/image-cover/lab-31-oauth2-and-oidc-cover-id.png` | 762.5 KB (1920x1080 PNG) | Cover Publikasi Medium & Sosial ID |
| **7** | **Cover Visual (EN)** | `result/image-cover/lab-31-oauth2-and-oidc-cover-en.png` | 762.0 KB (1920x1080 PNG) | Cover Publikasi Dev.to & Global EN |
| **8** | **Caption LinkedIn** | `result/linkedin/lab-31-oauth2-and-oidc-caption.md` | 232 kata (tanpa bold markdown) | Siap Dipublikasikan ke LinkedIn |

---

## 2. Catatan Rekayasa & Verifikasi Teknis

- **Repositori Target**: `https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/31-oauth2-and-oidc`
- **Commit SHA Diverifikasi**: `cf0cd90745a547032e8fb3a2cb838f40db1de30a`
- **Hasil Pengujian Unit**: `17/17 PASS` (`go test -v ./...` selesai dalam 0.317s).
- **Pengujian Balapan Konkurensi**: `PASS` (`go test -race ./...` selesai dalam 1.366s tanpa data race).
- **Hasil Eksekusi Demo**: Berhasil membuktikan alur legitimasi, penolakan intersepsi kode otorisasi via PKCE S256, rotasi refresh token, dan pencabutan keluarga token saat replay terdeteksi (`go run ./cmd/demo`).
- **Batasan Implementasi Utama**:
  1. *In-Memory Mock Storage*: Penyimpanan state kode dan metadata token berjalan di memori tanpa persistensi database relasional.
  2. *Symmetric Signing (HMAC-SHA256)*: Pemilihan HMAC-SHA256 menjaga implementasi mandiri (zero-dependency) tanpa dependensi JWKS endpoint eksternal.
  3. *State Parameter Server-Side Check*: Validasi state dijalankan pada helper client tanpa penyimpanan sesi state di server otorisasi.

---

## 3. Salinan Teks Caption LinkedIn Siap Pakai

```text
Lab 31 kali ini membahas pemisahan mendasar antara otorisasi OAuth 2.0 dan autentikasi OpenID Connect (OIDC) serta mekanisme pertahanan token modern di Go.

Kesalahan umum arsitektur modern adalah memperlakukan access token sebagai kartu identitas pengguna. Access token hanya memberikan izin akses delegasi terhadap sumber daya, bukan membuktikan siapa pemilik sesi. Mengabaikan batas ini membuka celah pembajakan akun, intersepsi kode otorisasi pada browser, dan pencurian sesi melalui refresh token statis.

Lab ini menguji empat mekanisme pertahanan:
• Pertahanan Intersepsi Kode: Proof Key for Code Exchange (PKCE) S256 memastikan hanya client pemegang verifier asli yang dapat menukar kode otorisasi.
• Verifikasi Identitas Kriptografis: ID Token berformat JWT divalidasi tanda tangannya dengan HMAC-SHA256 serta diperiksa kecocokan issuer, audience, kedaluwarsa, dan nonce anti-replay.
• Rotasi Refresh Token Sekali Pakai: Setiap pembaruan sesi menerbitkan pasangan token baru dan menghanguskan token lama.
• Deteksi Replay Satu Keluarga: Penggunaan ulang token yang sudah dicabut memicu pencabutan seluruh token aktif dalam keluarga yang sama.
• Concurrency Safety: Seluruh mutasi token dilindungi sync.RWMutex dan lolos pengujian race detector.

Mental model inti: OAuth 2.0 grants delegated access authority while OIDC asserts cryptographically verified user identity. PKCE binds authorization codes to the initiator, and refresh token rotation with family revocation neutralizes token replay.

Batasan implementasi: Pengujian berjalan in-memory tanpa persistensi database dan menggunakan symmetric signing HMAC-SHA256 untuk kesederhanaan demonstrasi.

Pelajari kode sumber lengkap dan pengujian di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/31-oauth2-and-oidc

#softwareengineering #golang #security #oauth2 #oidc #systemdesign
```
