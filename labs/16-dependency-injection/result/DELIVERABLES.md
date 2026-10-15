# Deliverables Paket Publikasi: Lab 16 - Dependency Injection

**Commit Hash:** `d7fb619`  
**Target Lab:** `labs/16-dependency-injection`  
**Seri:** Software Engineering Lab | Lab 16  
**Penulis:** Lukman / lukman-ss  

---

## 1. Daftar Berkas Deliverables Lengkap (8 Files + Content Def)

| No | Tipe Deliverable | Lokasi Berkas | Ukuran | Status Verifikasi |
|---|---|---|---|---|
| 0 | Python Source Definition | `result/content_def.py` | 35.0 KB | Verified (24 slides, DO NOT DELETE) |
| 1 | Dev.to Article (EN) | `result/dev.to/lab-16-dependency-injection-devto.md` | ~1.918 kata | Verified (22 H2 subsections, YAML frontmatter) |
| 2 | Medium Article (ID) | `result/medium/lab-16-dependency-injection-medium.html` | 12.4 KB | Verified (22 H2 subsections, clean semantic HTML) |
| 3 | Medium Metadata | `result/medium/lab-16-dependency-injection-metadata.md` | 1.1 KB | Verified (SEO, canonical URL, tags) |
| 4 | Cover Image (EN) | `result/image-cover/lab-16-dependency-injection-cover-en.png` | 898 KB | Verified (1920x1080 PNG via Chrome headless) |
| 5 | Cover Image (ID) | `result/image-cover/lab-16-dependency-injection-cover-id.png` | 892 KB | Verified (1920x1080 PNG via Chrome headless) |
| 6 | LinkedIn Carousel PDF | `result/linkedin/lab-16-dependency-injection-linkedin.pdf` | 2.72 MB | Verified (24 pages 16:9, selectable vector text) |
| 7 | LinkedIn Carousel PPTX | `result/linkedin/lab-16-dependency-injection-carousel.pptx` | 108 KB | Verified (24 slides + Speaker Notes via python-pptx) |
| 8 | LinkedIn Post Caption | `result/linkedin/lab-16-dependency-injection-caption.md` | 203 kata | Verified (Tanpa **bold**, alur ringkas & profesional) |

---

## 2. Catatan Integritas & Verifikasi Teknis

- **Implemented vs Simulated Behavior:**
  - *Implemented:* `Processor` dengan Constructor Injection (`NewProcessor(PaymentGateway)`), `BadProcessor` dengan Service Locator (`NewBadProcessor(Container)`), value object `Money`, dan manual composition root di `cmd/demo/main.go`.
  - *Simulated:* `RealGateway` mencetak log stdout (`RealGateway charging ...`), `MockGateway` mencatat mutasi state di memori dengan toggle `ShouldFail`.
- **Verified Metrics & Invariants:**
  - 6/6 test cases passing (`0.00s` per test case).
  - Go data race detector (`go test -race ./...`) lolos bersih tanpa data race warning.
  - Boundary check: input nominal `<= 0` mengembalikan error sebelum gateway dieksekusi (`ChargedMoney.Amount` tetap 0).
- **Boundaries & Caveats (Do Not Overclaim):**
  - Lab ini sengaja tidak menggunakan framework DI (Google Wire, Uber Dig, dsb.).
  - Durasi test `0.00s` adalah metrik eksekusi in-memory lokal, bukan estimasi latensi jaringan atau production SLO.
  - Test bebas race hanya membuktikan jalur sequential yang diuji, bukan garansi thread-safety concurrency pada skala production.
  - Ketiadaan explicit nil-check pada `NewProcessor` berarti runtime panic tetap mengancam jika caller mengirim interface nil.

---

## 3. Caption LinkedIn Siap Pakai (Ready to Copy)

```text
Software Engineering Lab 16 — Dependency Injection di Go ini membandingkan dua pendekatan pengelolaan ketergantungan: Constructor Injection eksplisit versus Service Locator yang menyembunyikan dependency.

Di mana sebuah komponen payment system menginstansiasi gateway jaringan secara langsung di dalam method-nya sendiri, setiap perubahan vendor memaksa modifikasi pada domain logic. Lab ini memindahkan keputusan instansiasi itu ke composition root dengan menerima PaymentGateway interface pada konstruktor.

Empat poin kunci:
1. Konstruktor dengan parameter interface membuat dependensi terlihat di signature publik — bukan sekedar memudahkan mock, tetapi membuat kebutuhan eksternal dapat dibaca oleh siapa pun.
2. Service Locator menyamarkan ketergantungan riil di balik method container.Get(...), sehingga engineer harus mem-trace indirection ekstra dan test terpaksa mem-mock container.
3. Nil guard dan lifecycle management bukan tanggung jawab pola DI — mereka adalah tanggung jawab implementasi yang masih perlu ditambahkan secara eksplisit.
4. Race detector bersih pada 6 test cases (0.00s) hanya membuktikan jalur yang diuji aman, bukan garantee semua lifecycle production.

Lab ini tidak memakai framework DI eksternal. Manual wiring di main.go membekalkan compile-time type safety, sementara BadProcessor berperilaku fungsional sama tetapi menambah indirection tak perlu. Simulasi in-memory mock bukan pengganti load test atau integration test produksi.

Lihat thread carousel di slide 1 dan link repo di profil untuk selengkapnya.
```
