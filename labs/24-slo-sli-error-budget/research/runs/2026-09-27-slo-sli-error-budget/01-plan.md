# 01 — Research Plan

Research Topic
SLO, SLI & Error Budget — Mengubah "Sistem Harus Stabil" Menjadi Angka yang Bisa Diputuskan (labs/24-slo-sli-error-budget, Senior Software Engineer Daily #24, Bahasa Indonesia).

Objective
Verifikasi setiap klaim teknis material di TOPIC SPECIFICATION terhadap sumber primer otoritatif. Pisahkan fakta terverifikasi vs klaim sumber vs interpretasi vs unknowns. Sediakan angka yang sudah di-cross-check untuk fase engineering (perhitungan budget, burn rate, downtime).

Research Questions
1. Definisi kanonis SLI, SLO, SLA, Error Budget?
2. Apakah contoh hitung SLI 999.200/1.000.000 = 99,92% benar?
3. Apakah angka downtime 99% / 99,9% / 99,99% per bulan di lab akurat vs tabel otoritatif?
4. Apakah Error Budget = 1 − SLO dan contoh 10.000.000 req → 10.000 budget benar?
5. Apakah jawaban latihan (webhook 200.000 req, SLO 99,99% → 20 gagal; 15 gagal = 75%) benar?
6. Apakah klaim "100% bukan target" didukung sumber primer? Alasannya?
7. Apakah panduan "CPU bukan SLO user", "percentile bukan average", "SLO beda per endpoint", "burn rate untuk alerting" didukung?
8. Rekomendasi window (30 hari vs 4 minggu) dan error budget policy apa yang didokumentasikan primer?

Search Strategy
1. Mulai dari sumber Tier 1: Google SRE Book (Ch.3 Embracing Risk, Ch.4 SLO, Appendix A Availability Table), SRE Workbook (Ch.2 Implementing SLOs, Ch.5 Alerting on SLOs, Appendix B Error Budget Policy).
2. Cross-check independen Tier 1/2: Google Cloud SLO monitoring docs, Prometheus alerting practices, Datadog SLO docs.
3. Buka tiap halaman penuh via WebFetch; tidak pakai snippet. Klaim AWS Builders Library dicoba tetapi URL 404 pada 2026-09-27 → tidak dipakai sebagai evidence.
4. Verifikasi aritmetika lab dengan hitung ulang manual.

Expected Primary Sources
- https://sre.google/sre-book/service-level-objectives/
- https://sre.google/sre-book/embracing-risk/
- https://sre.google/sre-book/availability-table/
- https://sre.google/workbook/implementing-slos/
- https://sre.google/workbook/alerting-on-slos/
- https://sre.google/workbook/error-budget-policy/
- https://cloud.google.com/stackdriver/docs/solutions/slo-monitoring
- https://prometheus.io/docs/practices/alerting/
- https://docs.datadoghq.com/service_level_objectives/

Risks / Unknowns
- Semua sumber Tier 1 berasal dari ekosistem Google; risiko bias vendor tunggal. Mitigasi: cross-check Datadog + Prometheus.
- Angka downtime bergantung asumsi panjang bulan (30 hari vs 30,44 hari) → selisih kecil diekspektasikan, didokumentasikan di contradictions.
- Statistik "70% outage dari change" hanya klaim internal Google tanpa metodologi publik → tandai LOW.
- Rekomendasi window dan threshold burn rate bersifat praktik, bukan hasil eksperimen terkontrol.
