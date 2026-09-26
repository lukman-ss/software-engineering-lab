# Research Report

## Research Question
Bagaimana SLI, SLO, dan Error Budget berperan dalam reliability engineering, apa best practice penetapannya, dan apa pitfall umum — khususnya untuk konteks Senior Software Engineer (Bahasa Indonesia) berdasarkan temuan laboratorium?

## Executive Summary
SLI, SLO, dan Error Budget adalah mekanisme Google SRE untuk mengubah "harus stabil" menjadi angka ambil-keputusan. Inti: jauhi target 100% (maksa biaya ekstrem, tapi tidak terlihat oleh user), ukur dari perspektif user (bukan CPU/RAM), pakai percentile bukan average, dan gunakan error budget sebagai alat negosiasi antara product vs engineering. Semua temuan didasari primaduga satu sumber (Google SRE Book + Workbook) konsisten secara internal.

## Findings

### Finding 1 — SLI: ukuran kuantitatif kualitas layanan, rekomendasi bentuk rasio good/total
Claim: SLI adalah indikator kuantitatif; SRE merekomendasikan SLI berupa `good events / total events` agar rentang 0–100% dan konsisten tooling.
Evidence: "An SLI is a service level indicator—a carefully defined quantitative measure..." + "we generally recommend treating the SLI as the ratio of two numbers: the number of good events divided by the total number of events."
Sources: SRE Book Ch.4 §Indicators in Practice; SRE Workbook Ch.2 "What to Measure".
Confidence: HIGH

### Finding 2 — SLO: target untuk SLI, struktur ≤ atau rentang
Claim: SLO adalah nilai target/rantau untuk sebuah SLI. Bisa multiple thresholds (90% < 100ms, 99% < 400ms) untuk menangkap distribusi user experience.
Evidence: "An SLO is a service level objective: a target value or range of values..." + contoh bertingkat P90/P99.
Sources: SRE Book Ch.4 §Objectives; SRE Workbook Ch.2 lat. 2-3 (multiple thresholds).
Confidence: HIGH

### Finding 3 — SLA ≠ SLO; akibat = kontrak
Claim: SLA mengandung konsekuensi (mis. pengembalian dana); SLO tidak punya konsekuensi eksplisit. Orang yang bilang "SLA violation" biasanya maksudnya SLO.
Evidence: "SLAs are service level agreements: an explicit or implicit contract... if there is no explicit consequence, then you are almost certainly looking at an SLO."
Sources: SRE Book Ch.4 §Agreements.
Confidence: HIGH

### Finding 4 — Empat sinyal emas (golden signals)
Claim: Fokus ukuran user-facing: Latency, Traffic, Errors, Saturation.
Evidence: "The four golden signals of monitoring are latency, traffic, errors, and saturation."
Sources: SRE Book Ch.6 "The Four Golden Signals".
Confidence: HIGH (per sumber) / MEDIUM (cross-check tidak lakukan ke sumber independen)

### Finding 5 — Gunakan percentile, bukan rata-rata
Claim: Mean tidakrepresentatif untuk latency; gunakan P50/P95/P99/P99.9 dan histogram bucketed.
Evidence: Fig 4-1: typical 50ms tapi 5% request 20× lambat; average menyembunyikan tail.
Sources: SRE Book Ch.4 §Aggregation; Ch.6 §Worrying About Your Tail.
Confidence: HIGH

### Finding 6 — 100% availability adalah target yang salah
Claim: 100% reliability (a) tidak realistis (kegagalan komponen nonzero), (b) tidak terlihat user (rantai panjang ke client), (c) cost naik tiap nines sekaligus, (d) change = outage #1 → keharusan update menjadi mustahil.
Evidence: "100% reliability is the wrong target" + 4 poin; SRE Book Ch.3: user di smartphone 99% tak bedakan 99.99% vs 99.999%.
Sources: SRE Workbook Ch.2 "Reliability Targets"; SRE Book Ch.3 "Embracing Risk".
Confidence: HIGH

### Finding 7 — Availability = nines: tabel downtime resmi
Claim: 99% ≈ 7j12m/bulan, 99.9% = 8j44m/yr (43.2 menit/bulan), 99.99% = 52.6 menit/yr (4.32 menit/bulan), 99.999% = 5.26 menit/yr.
Evidence: Appendix A Table 1-1.
Sources: SRE Book Appendix A.
Confidence: HIGH. Lihat Contradiction (selisih 6 menit vs klaim lab 7j18m) — karena asumsi bulan.

### Finding 8 — Error Budget = 1 − SLO
Claim: Error budget adalah selisih SLO vs uptime aktual suatu window; sisa budget menjadi sinyal kapan boleh berisiko.
Evidence: "An error budget is 1 minus the SLO... The difference between these two numbers is the budget of how much unreliability is remaining."
Sources: SRE Book Ch.3 §Motivation for Error Budgets; SRE Workbook Appendix B.
Confidence: HIGH. Contoh valid: 1M request SLO 99.9% → budget 1000 errors.

### Finding 9 — Error budget sebagai alat keputusan (bukan sekadar dashboard)
Claim: Jika budget tersisa → release/experiment/refactor; jika hampir habis → hentikan deployment berisiko, perbaiki reliability, bayar tech debt.
Evidence: "As long as the uptime measured is above the SLO—in other words, as long as there is error budget remaining—new releases can be pushed."
Sources: SRE Book Ch.3 §Forming Your Error Budget; SRE Workbook Appendix B SLO Miss Policy.
Confidence: HIGH

### Finding 10 — Burn rate untuk alerting (multiwindow, multi-burn-rate)
Claim: Burn rate = laju konsumsi budget relatif SLO. Rekomendasi awal: page 14.4× (2% budget/1jam + 5menit), page 6× (5%/6jam + 30menit), ticket 1× (10%/3hari + 6jam). Jangan pakai `for: 1h` saja.
Evidence: Tabel 5-8 + contoh PromQL; short window = 1/12 long window.
Sources: SRE Workbook Ch.5 §Ways to Alert / §6: Multiwindow, Multi-Burn-Rate.
Confidence: HIGH. Peringatan: untuk 100% outage pada SLO 99.999%/monthly → 26 detik, perlu canarying.

### Finding 11 — SLO berbeda per endpoint berdasarkan criticality
Claim: Jangan semua endpoint dapat SLO 99.99%; `POST /payment/webhook` (efek finansial) > `GET /report` (boleh 99.5%).
Evidence: "Not all requests are considered equal... bucketing untuk menambah label ke SLI dan menerapkan SLO berbeda" + Ch.4: criticality berbeda (contoh SLO bertingkat).
Sources: SRE Book Ch.4 §Objectives; SRE Workbook Ch.5 "Grading Interaction Importance" + Tabel 5-10 (buckets CRITICAL/HIGH_FAST/LOW/NO_SLO).
Confidence: HIGH

### Finding 12 — CPU/RAM bukan SLO user
Claim: SLO baik ukur pengalaman user (request berhasil? cepat?). CPU/RAM/connection pool = diagnostic signal, bukan SLO.
Evidence: "Ideally, the SLI directly measures a service level of interest..." + "we recommend standardizing on common definitions..."
Sources: SRE Book Ch.4 §Indicators in Practice; Ch.3 contoh AdSense latency vs AdWords latency (cost trade-off).
Confidence: MEDIUM (interpretasi konsisten, tidak ada perbandingan eksplisit CPU bukan SLO di korpus, tapi prinsip user-experience sudah jelas).

## Areas of Agreement
- 100% bukan target.
- SLI rasio good/total, percentile bukan mean.
- Error budget = 1−SLO, jadi tools negosiasi.
- SLO harus ada konsekuensi (policy) atau hanya dekorasi.
- Perbedaan criticality endpoint → SLO berbeda.
- Burn rate alerting unggul pada detection time + reset time.

## Areas of Disagreement
- Beda asumsi hitung bulan: klaim lab (30.44 hari) vs tabel Google (30 hari tepat / 28 hari minggu penuh) — bukan kontradiksi teknis; lihat 04-contradictions §1-2.

## Limitations
- Semua sumber Tier 1 berasal dari Google SRE. Tidak ada sumber independen (mis. AWS, Azure, CNCF, studi akademik) yang berhasil diambil (404/redirect) dalam jendela riset ini.
- Statistik "70% outages from change" dan rekomendasi empiris (4-minggu window) adalah klaim internal Google tanpa metodologi publik.
- Konteks bahasa Indonesia; istilah resmi (SLI/SLO/SLA/burn rate) tetap Inggris di literatur teknis.

## Conclusion
Empirisasi reliability didesak oleh Google SRE: definisikan SLI dari perspektif user (bentuk rasio good/total, pakai percentile), tetapkan SLO realistis (<100%), ubah budget kegagalan jadi alat keputusan. Pustaka utama: Google SRE Book Ch.3-4,6,10; SRE Workbook Ch.2,5 + Appendix B. Pengeabstraksian angka (selisih 6 menit) tidak mengubah keputusan engineering.
