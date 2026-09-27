# 06 — Open Questions

## Unanswered Questions

1. **Industry-wide SLO adoption rates:** Berapa persen organisasi engineering yang secara luas mengimplementasikan SLO-based error budgeting? Tidak ada data kuantitatif publik universal yang tersedia.

2. **Optimasi burn rate thresholds:** Apakah ada standar otoritatif (Google, CNCF, NIST) untuk parameter burn rate selain rekomendasi Google 14.4×/6×/1×? Datadog menggunakan 1-6 (elevated) / 6+ (critical) pada 2h window.

3. **Statistik "70% outage dari change":** Tidak ada sumber independen (AWS, Azure, postmortemm data, paper akademik) yang diverifikasi yang mendukung angka 70%. Google hanya mengutip ini tanpa metodologi.

4. **Multi-service error budget allocation:** Bagaimana secara praktis mengalokasikan/fungsi error budget ke banyak layanan dengan SLO berbeda? Google SRE tidak menyebutkan framework untuk ini.

5. **SLO untuk data correctness/durability:** Bagaimana mengukur "benar/baik" sebagai SLI pada sistem storage/pipeline? Google menyebut penting, tapi tidak memberikan metode praktis.

## Weak Evidence Areas

1. **Burn rate numeric thresholds** — Hanya Google SRE Workbook + Datadog. Tidak ada sumber tidak-Google (AWS, Azure, Prometheus Cookbook) yang menyebut "burn rate 14.4x".

2. **"Cost 100x per nine" claim** — Google SRE Book menyatakan "an incremental improvement in reliability may cost 100x more" tanpa metodologi empiri yang dikutip. Tidak ada studi industri independen.

3. **User preference for lower variance** — Google menyatakan "people prefer slightly slower system to one with high variance" tanpa referensi studi pengguna yang dikutip.

4. **SLO status corrections formula** — Datadog formalisasi, Google hanya contoh Chubby planned outage. Tidak ada standar industri.

## Claims Needing Deeper Research

1. **Perbandingan window 28 hari vs 30 hari** — Apakah ada data A/B testing yang menunjukkan mana yang lebih baik untuk traffic dengan pola weekday/weekend yang berbeda?

2. **Error budget policy di enterprise kecil** — Bagaimana implementasi kebijakan budget error di organisasi tanpa dedicated SRE team?

3. **Statistik internal Google vs industry** — Apakah data Google tentang "70% outage from change" konsisten dengan data dari konferensi SREcon, post-mortem publik, atau studi DevOps Research?

4. **SLO vs Feature Velocity trade-off empiris** — Apakah ada data kuantitatif perusahaan yang menunjukkan correlation antara error budget consumption dan deployment frequency?

## Possible Next Research Directions

1. **Cross-company postmortem collection** — Kumpulkan postmortem publik (AWS Status, GCP Status, GitHub Status, Netflix TechBlog) untuk menganalisis distribusi penyebab outage.

2. **Vendor SLO comparison** — Bandingkan implementasi SLO di AWS CloudWatch (Service Health), GCP SLO, Azure Monitor, New Relic, VictorOps.

3. **OpenTelemetry integration** — Bagaimana SLI/SLO berlaku pada traces/metrics/logs? OpenTelemetry semantic conventions untuk service-level metrics apa yang tersedia?

4. **Case study industri** — Cari informasi publik dari:
   - Netflix Tech Blog (Chaos Engineering, SRE)
   - Uber Engineering (Metric-based SLOs)
   - Shopify Engineering (Error budget culture)
   - LinkedIn Engineering (SRE practices)

5. **Akademik research** — Cari di Google Scholar untuk "Service Level Objectives", "Error Budget SRE", "Reliability Engineering Metrics".

6. **Indonesian SLO adoption** — Eksplorasi adopsi SLO di perusahaan teknologi Indonesia (Gojek, Tokopedia, Traveloka, Bukalapak) melalui blog teknis mereka.

---

**Status:** Beberapa klaim utama telah diverifikasi (SLI/SLO/Error Budget definisi, kalkulasi, downtime, burn rate). Namun beberapa statistik industri (70% outages from change, 100x cost per nine) tetap bersifat klaim internal tanpa verifikasi independen. Error budget policy example dari Google dapat diadopsi, tetapi perlu disesuaikan dengan budaya organisasi masing-masing.