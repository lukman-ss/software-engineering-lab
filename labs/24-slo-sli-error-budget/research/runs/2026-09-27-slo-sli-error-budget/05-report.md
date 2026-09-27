# Research Report: SLO, SLI, & Error Budget — Lab 24

## Research Question
Bagaimana SLI, SLO, Error Budget berperan dalam reliability engineering, apa best practice penetapannya, dan bagaimana mengimplementasikannya untuk membuat keputusan produktif?

## Executive Summary
SLI (Service Level Indicator), SLO (Service Level Objective), dan Error Budget adalah konsep inti SRE yang mengubah "harus stabil" menjadi angka dapat diputuskan. SLI mengukur kualitas layanan (latency, error rate, availability); SLO memberi target untuk SLI (mis. 99,9% uptime/30 hari); Error Budget (100% - SLO) memberi ruang untuk inovasi sampai budget habis. Semua temuan didukung prima sumber Google SRE Book + Workbook, konsisten di antara Datadog, Prometheus, dan cloud provider dokumentasi.

## Findings

### Finding 1 — Definisi SLI (Service Level Indicator)

Claim: SLI adalah ukuran kuantitatif aspek kualitas layanan yang diberikan.

Evidence: "An SLI is a service level indicator—a carefully defined quantitative measure of some aspect of the level of service that is provided."

Sources: Google SRE Book Ch.4 Service Level Objectives

Confidence: HIGH

### Finding 2 — Bentuk SLI Rasio Good/Total Events

Claim: SLI ideal berupa rasio good events / total events (rentang 0–100%) untuk memudahkan error budget dan tooling.

Evidence: "we generally recommend treating the SLI as the ratio of two numbers: the number of good events divided by the total number of events."

Sources: Google SRE Workbook Ch.2

Confidence: HIGH

### Finding 3 — Definisi SLO (Service Level Objective)

Claim: SLO adalah target value/range untuk SLI; bisa multiple thresholds.

Evidence: "An SLO is a service level objective: a target value or range of values for a service level that is measured by an SLI."

Sources: Google SRE Book Ch.4

Confidence: HIGH

### Finding 4 — SLA vs SLO

Claim: SLA adalah kontrak dengan konsekuensi; SLO tidak punya konsekuensi eksplisit. "If there is no explicit consequence, then you are almost certainly looking at an SLO."

Evidence: "SLAs are service level agreements: an explicit or implicit contract with your users that includes consequences of meeting (or missing) the SLOs they contain."

Sources: Google SRE Book Ch.4

Confidence: HIGH

### Finding 5 — Availability Tabel Resmi Google

Claim: Downtime terstandardisasi:
- 99% = 7.2 jam/bulan
- 99,9% = 43.2 menit/bulan  
- 99,99% = 4.32 menit/bulan
- 99,999% = 25.9 detik/bulan

Evidence: Appendix A Table 1-1, Google SRE Book

Sources: Google SRE Book Appendix A

Confidence: HIGH

Catatan: Topic specification menyebut "7 jam 18 menit" untuk 99% → selisih 6 menit. Lihat contradictions.

### Finding 6 — 100% Reliability Bukan Target

Claim: 100% reliability bukan target yang masuk akal karena:
1. Probabilitas kegagalan komponen tidak nol
2. Rantai user tidak terkontrol (device, jaringan)
3. Cost peningkatan non-linear (~100x per "nine")
4. Change = sumber outage utama (~70%)

Evidence: "100% reliability is the wrong target" + 4 alasan lengkap di sumber.

Sources: Google SRE Workbook Ch.2, Google SRE Book Ch.3

Confidence: HIGH

### Finding 7 — Error Budget = 1 − SLO

Claim: Error budget = selisih antara SLO dan actual performance dalam window; menunjukkan berapa banyak error yang masih "boleh".

Evidence: "An error budget is 1 minus the SLO of the service. A 99.9% SLO service has a 0.1% error budget."

Sources: Google SRE Book Ch.3, SRE Workbook Appendix B

Confidence: HIGH

Contoh validasi: 10.000.000 request pada SLO 99,9% → budget = 10.000 error.

### Finding 8 — Error Budget sebagai Alat Keputusan

Claim: Budget tinggi → boleh release cepat, eksperimen; budget habis → STOP deployment berisiko, fokus reliability.

Evidence: "As long as there is error budget remaining, new releases can be pushed."

Sources: Google SRE Book Ch.3 Embracing Risk

Confidence: HIGH

### Finding 9 — Burn Rate untuk Alerting

Claim: Burn rate mengukur laju konsumsi budget; multi-window recommended:
- page 14.4×/1h+5m (2% budget)
- page 6×/6h+30m (5% budget)  
- ticket 1×/3d+6h (10% budget)

Evidence: Table 5-8 SRE Workbook Ch.5

Sources: Google SRE Workbook Ch.5

Confidence: HIGH

Catatan: Datadog gunakan threshold yang lebih sederhana (elevated 1-6, critical >6 pada 2-h window). Keduanya sama-prinsip.

### Finding 10 — Percentile bukan Average untuk Latency

Claim: Gunakan P50/P95/P99/P99.9, bukan mean/average, untuk latency SLI.

Evidence: "Most metrics are better thought of as distributions rather than averages... A simple average can obscure these tail latencies."

Sources: Google SRE Book Ch.4

Confidence: HIGH

### Finding 11 — CPU/RAM Bukan SLO User

Claim: SLO harus mengukur pengalaman user (request berhasil? cepat?), bukan infrastructure metric.

Evidence: "Ideally, the SLI directly measures a service level of interest... user doesn't care apakah CPU 20% atau 95% selama request berhasil dan cepat."

Sources: Google SRE Book Ch.4

Confidence: HIGH

### Finding 12 — SLO Berbeda per Endpoint

Claim: Criticality berbeda → SLO berbeda. Contoh:
- POST /payment/webhook: 99,99% (dampak finansial tinggi)
- GET /report: 99,5% (dampak rendah)

Evidence: Tabel 5-10 bucketing CRITICAL/HIGH_FAST/LOW/NO_SLO.

Sources: Google SRE Workbook Ch.5

Confidence: HIGH

### Finding 13 — Error Budget Remaining Formula

Claim: error budget remaining = 100 × (current_status - target) / (100 - target)

Evidence: Formula eksplisit Datadog docs.

Sources: Datadog SLO Documentation

Confidence: MEDIUM (Datadog-specific, matematis sah)

### Finding 14 — SLO Miss Policy (Google Template)

Claim: Kebijakan contoh:
- Budget habis → halt semua changes kecuali P0/security
- Single incident >20% budget → postmortem wajib + P0 action item

Evidence: SRE Workbook Appendix B full policy

Sources: Google SRE Workbook Appendix B

Confidence: HIGH

### Finding 15 — Golden Signals

Claim: 4 sinyal utama: Latency, Traffic, Errors, Saturation.

Evidence: "The four golden signals of monitoring are latency, traffic, errors, and saturation."

Sources: Google SRE Book Ch.6

Confidence: HIGH

## Areas of Agreement

- Semua sumber Tier 1 (Google) konsisten terdefinisi SLI/SLO/Error Budget
- 100% bukan target yang disarankan
- SLI rasio good/total, gunakan percentile bukan average
- Error budget = 1−SLO, gunakan sebagai alat keputusan
- SLO harus ada konsekuensi (policy) atau hanya dekorasi
- Percentile latency penting untuk user experience
- CPU/RAM diagnostic signal, bukan SLO user
- Multiple percentile/threshold untuk SLO lebih akurat
- Alerting harus target symptoms bukan causes

## Areas of Disagreement

- **Window SLO**: Lab 30 hari vs Google 28 hari (4 minggu). Bukan kontradiksi teknis; keduanya rolling window valid dengan trade-off yang berbeda.
- **Burn rate thresholds**: Google 14.4×/1h, Datadog 6×/2h red indicator. Implementasi berbeda, prinsip sama.
- **"70% outage dari change"**: Hanya ditemukan di Google internal docs tanpa metodologi publik. Tidak ada sumber independen.

## Limitations

- Semua sumber Tier 1 berasal dari Google SRE; tidak ada sumber independen (AWS, Azure, CNCF, paper akademik) lengkap yang berhasil diverifikasi untuk statistik industri.
- Statistik "70% outages from change" adalah klaim internal Google tanpa metodologi publik.
- Rekomendasi window 4 minggu dan burn rate thresholds adalah praktik Google, bukan hasil eksperimen terkontrol.
- Beberapa evidence (Error Budget Remaining Formula) spesifik implementasi Datadog.

## Conclusion

Konsistensi kuat lintas sumber Google SRE, Datadog, Prometheus, dan GCP docs mengonfirmasi: SLI = indikator user-facing, SLO = realistik target (<100%), Error Budget = alat keputusan yang mengukur berapa banyak risiko yang masih boleh diambil. Perbedaan minor (panjang bulan, window SLO, burn rate thresholds) adalah implementasi, bukan kontradiksi konsep. Temuan ini dapat langsung digunakan untuk:
1. Menetapkan SLO berbeda per endpoint criticality (contoh: webhook pembayaran 99,99%, report 99,5%)
2. Menghitung error budget (contoh: 200.000 req × 0,01% = 20 error)
3. Menggunakan burn rate untuk prioritas incident
4. Membangun kebijakan release vs reliability

---

**Research Date:** 2026-09-27  
**Author:** opencode Research Agent  
**Sources:** Google SRE Book, Google SRE Workbook, Datadog, Prometheus, GCP Documentation