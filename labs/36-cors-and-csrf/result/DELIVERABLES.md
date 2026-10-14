# Deliverables Rekapitulasi: Lab 36 — CORS & CSRF Defense

Paket deliverable publikasi lengkap untuk **Software Engineering Lab 36: CORS & CSRF — Securing Backend from Cross-Origin Attacks** oleh **Lukman / lukman-ss**.

---

## Tabel Unduhan & Berkas Publikasi

| No | Tipe Deliverable | File Output | Status / Ukuran | Deskripsi |
| :---: | :--- | :--- | :--- | :--- |
| **0** | **Python Definition** | `result/content_def.py` | Source Code | Definisi 24 slide, konstanta warna, metadata, dan mental model (DILARANG DIHAPUS) |
| **1** | **Dev.to Article** | `result/dev.to/lab-36-cors-and-csrf-devto.md` | ~11.3 KB (1.402 kata) | Artikel teknis Bahasa Inggris lengkap dengan frontmatter dan 22 subbagian |
| **2** | **Medium Article** | `result/medium/lab-36-cors-and-csrf-medium.html` | ~14.9 KB (1.530 kata) | Standalone HTML Bahasa Indonesia dengan 22 section dan styling Georgia/Code |
| **3** | **Medium Metadata** | `result/medium/lab-36-cors-and-csrf-metadata.md` | Markdown | Metadata SEO, slug, social title, dan 5 tags Medium |
| **4** | **Cover Image (EN)**| `result/image-cover/lab-36-cors-and-csrf-cover-en.png` | 797 KB (1920x1080) | Cover monitor CRT retro pixel-art dengan subtitle Bahasa Inggris |
| **5** | **Cover Image (ID)**| `result/image-cover/lab-36-cors-and-csrf-cover-id.png` | 772 KB (1920x1080) | Cover monitor CRT retro pixel-art dengan subtitle Bahasa Indonesia |
| **6** | **LinkedIn PDF** | `result/linkedin/lab-36-cors-and-csrf-linkedin.pdf` | 818 KB (24 slide) | Carousel dokumen PDF 16:9 24 halaman (Cover Dark Navy, 22 Editorial Cream, Closing Dark Navy) |
| **7** | **LinkedIn PPTX** | `result/linkedin/lab-36-cors-and-csrf-carousel.pptx` | 115 KB (24 slide) | File presentasi PowerPoint 24 slide dengan layout 16:9 dan Speaker Notes |
| **8** | **LinkedIn Caption**| `result/linkedin/lab-36-cors-and-csrf-caption.md` | 228 kata | Teks caption LinkedIn tanpa bold markdown siap salin |

---

## Verifikasi Kualitas & Catatan Teknis

1. **Commit Snapshot**: `7ced81ed`
2. **Status Pengujian**:
   - `go test -v ./...` -> 13/13 PASS (0.42s)
   - `go test -race ./...` -> PASS (Race-Clean Verified Go 1.22)
   - `go run ./cmd/demo` -> Berhasil mendemonstrasikan serangan pada endpoint rentan vs blokade pada endpoint terlindungi
3. **Batas Implementasi Utama**:
   - Menggunakan secret key in-memory untuk demonstrasi pedagogis; produksi memerlukan Secret Manager/Vault.
   - Menggunakan delimiter titik dua (`:`) pada parsing token; produksi memerlukan serialisasi terstruktur (JWS/Protobuf).
   - Lingkup SOP, CORS, dan CSRF berada di browser; client non-browser tidak dibatasi oleh SOP.
   - Kerentanan XSS pada origin yang sama dapat melumpuhkan seluruh pertahanan anti-CSRF.

---

## Salinan Teks Caption LinkedIn

```text
Lab 36 kali ini membahas miskonsepsi fatal dalam keamanan web modern bahwa konfigurasi CORS dapat melindungi backend dari mutasi data berbahaya. Pada kenyataannya, CORS hanya mengatur izin pembacaan respons oleh JavaScript di browser, sedangkan serangan CSRF mengeksploitasi pengiriman request write lintas origin yang diizinkan secara historis oleh Same-Origin Policy.

Banyak engineer mengira jika origin penyerang tidak terdaftar di whitelist CORS, server otomatis menolak request tersebut. Padahal, browser tetap mengirimkan simple POST langsung ke server beserta cookie otentikasi korban, sehingga mutasi database selesai sebelum browser memeriksa respons.

Eksperimen Go standard library pada lab ini membuktikan:
• CORS tidak mencegah eksekusi backend: simple POST dari evil-attacker.com memotong saldo rekening korban dari $1000 menjadi $700 meskipun header Access-Control-Allow-Origin tidak diberikan.
• Signed Double-Submit Token HMAC-SHA256 yang terikat ke Session ID dan timestamp TTL berhasil memblokir serangan CSRF tanpa token dengan HTTP 403 Forbidden.
• Fetch Metadata (Sec-Fetch-Site) menolak mutasi cross-site di layer gateway secara mandiri.
• Custom Header X-Requested-With menetralisasi eksploitasi form HTML klasik secara total.
• Konfigurasi cookie SameSite=Lax menahan pengiriman cookie pada request mutasi lintas situs.

Mental model: SOP restricts cross-origin read, CSRF exploits cross-origin write. CORS relaxes read, it never defends write.

Batasan implementasi: Lab ini menggunakan HMAC secret in-memory dan string delimiter untuk parsing token; sistem produksi memerlukan secret vault dan serialisasi data terstruktur.

Pelajari kode lengkap dan pengujian di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/36-cors-and-csrf

#softwareengineering #golang #cybersecurity #webdevelopment #systemdesign #backend
```
