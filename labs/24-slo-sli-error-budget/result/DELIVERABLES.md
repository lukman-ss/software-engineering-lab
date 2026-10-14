# Publication Deliverables — Software Engineering Lab 24

## Tabel File Deliverable

| No | Tipe Deliverable | Lokasi File | Deskripsi |
|---|---|---|---|
| 1 | PowerPoint Carousel | [linkedin/lab-24-slo-sli-error-budget-carousel.pptx](linkedin/lab-24-slo-sli-error-budget-carousel.pptx) | File .pptx 24 slide 16:9 editable dengan Speaker Notes (via python-pptx) |
| 2 | LinkedIn Carousel PDF | [linkedin/lab-24-slo-sli-error-budget-linkedin.pdf](linkedin/lab-24-slo-sli-error-budget-linkedin.pdf) | File PDF 24 slide rasio 16:9 (Cover CRT + 22 slide krem editorial + Penutup navy) |
| 3 | Dev.to Article | [dev.to/lab-24-slo-sli-error-budget-devto.md](dev.to/lab-24-slo-sli-error-budget-devto.md) | Artikel Markdown Bahasa Inggris lengkap dengan frontmatter |
| 4 | Medium Article | [medium/lab-24-slo-sli-error-budget-medium.html](medium/lab-24-slo-sli-error-budget-medium.html) | Artikel HTML lengkap Bahasa Indonesia (~1.310 kata) |
| 5 | Cover Image (ID) | [image-cover/lab-24-slo-sli-error-budget-cover-id.png](image-cover/lab-24-slo-sli-error-budget-cover-id.png) | Cover PNG 1920×1080 retro CRT terminal (Bahasa Indonesia) |
| 6 | Cover Image (EN) | [image-cover/lab-24-slo-sli-error-budget-cover-en.png](image-cover/lab-24-slo-sli-error-budget-cover-en.png) | Cover PNG 1920×1080 retro CRT terminal (Bahasa Inggris) |
| 7 | Medium Metadata | [medium/lab-24-slo-sli-error-budget-metadata.md](medium/lab-24-slo-sli-error-budget-metadata.md) | Metadata SEO, slug, tags, social preview, & hitungan karakter |
| 8 | Caption LinkedIn | [linkedin/lab-24-slo-sli-error-budget-caption.md](linkedin/lab-24-slo-sli-error-budget-caption.md) | Caption teks 243 kata siap salin (tanpa markdown bold) |

## Catatan Verifikasi & Keterbatasan
- **Commit SHA**: `ffb1d0c42b78d71a36d1e3b9769d01bd8f458f67`
- **Status Test**: 6/6 test cases PASS (`go test -race -v ./...` durasi 2.904s, 0 data races).
- **Batas Implementasi**: Lab ini mengimplementasikan pelacakan metrik in-memory Go dengan window tracker time-bucketed untuk mendemonstrasikan penghitungan error budget matematis, penegakan deployment freeze, dan pencegahan false alarm multi-window alerting secara deterministik. Lab tidak menyimpan riwayat data ke database deret waktu eksternal (Prometheus/TSDB), tidak mengukur kuantil histogram latency lengkap (P95/P99), dan menyimulasikan kompresi skala waktu agar eksekusi demonstrasi berjalan cepat. Di lingkungan produksi, integrasi TSDB persisten dan konfigurasi rolling window 28 hari direkomendasikan.

## Caption LinkedIn (Siap Salin)

```text
Lab 24 kali ini membahas Service Level Indicators (SLI), Service Level Objectives (SLO), kalkulasi matematis Error Budget, dan mesin peringatan Multi-Window Multi-Burn-Rate berbasis prinsip Google SRE di Go.

Masalah teknis yang sering terjadi adalah tim keandalan mengukur stabilitas layanan memakai metrik infrastruktur seperti utilisasi CPU dan memori host. Metrik tersebut tidak mencerminkan pengalaman pengguna, sehingga kegagalan transaksi HTTP 500 tersembunyi di balik dasbor server yang tampak sehat. Selain itu, alert berbasis satu jendela waktu memicu alarm palsu pada gangguan jaringan sesaat atau terlambat mendeteksi pemadaman total.

Lab ini membuktikan empat pilar keandalan konkret:
• Indikator SLI Kuantitatif: Mengukur rasio kejadian sukses terhadap total transaksi melalui predikat evaluasi caller, bukan metrik hardware mentah.
• Bucketing Waktu Geser: WindowTracker mengelompokkan request ke interval tetap dan memangkas bucket kadaluarsa untuk efisiensi memori konkuren.
• Penegakan Freeze Otomatis: Evaluator menghitung sisa error budget secara absolut dan memblokir izin deploy saat kuota toleransi kesalahan habis.
• Alerting Multi-Jendela: Mesin alarm mengevaluasi jendela pendek dan jendela panjang secara simultan, menekan alarm palsu hingga nol pada lonjakan sesaat.

Mental model inti: Reliability is not 100%. Target 100% adalah ilusi teknis; error budget mendefinisikan kuota kegagalan yang dapat dihabiskan untuk inovasi dan rilis fitur terukur.

Satu batasan penting: Lab ini menggunakan penyimpanan metrik in-memory Go yang tereset saat proses restart dan menyimulasikan kompresi waktu jendela, sehingga lingkungan produksi memerlukan basis data deret waktu seperti Prometheus atau VictoriaMetrics.

Pelajari implementasi lengkap dan verifikasi pengujian di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/24-slo-sli-error-budget

#softwareengineering #golang #sre #devops #observability #systemdesign
```
