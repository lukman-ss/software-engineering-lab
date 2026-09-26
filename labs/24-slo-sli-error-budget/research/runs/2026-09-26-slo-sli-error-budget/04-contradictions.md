# Contradictions

## Contradiction 1: Downtime per Month untuk 99% — Selisih 6 Menit
- **Topic Lab klaim:** 99% → ~7 jam 18 menit/bulan
- **Google Appendix A klaim:** 99% → 7.2 jam/bulan = 7 jam 12 menit/bulan
- **Assessment:** Selisih 6 menit (~1.4%). Penyebab: asumsi panjang bulan berbeda. Lab 7j18m ≈ 438 menit → 0.6% dari 30.44 hari (bulan rata-rata 365/12). Google 7.2 jam = 432 menit → 1% dari 30 hari tepat (43,200 menit). Bukan kontradiksi substantif; keduanya benar di bawah asumsi berbeda. Untuk konsistensi, lab perlu eksplisit: "30 hari" vs "bulan kalender rata-rata".
- **Impact:** LOW. Tidak mengubah keputusan engineering.

## Contradiction 2: 99.99% per Bulan — 4m23s vs 4m19s
- **Topic Lab:** 99.99% → ~4 menit 23 detik/bulan
- **Google Appendix A:** 99.99% → 4.32 menit = 4 menit 19.2 detik/bulan
- **Assessment:** Selisih 3.8 detik (~1.5%). Penyebab sama: basis hitung bulan. Lab sedikit lebih longgar (mungkin pembulatan).
- **Impact:** LOW.

## Contradiction 3: Window SLO — 30 Hari (Lab) vs 28 Hari (Rekomendasi Google)
- **Topic Lab contoh:** SLO 99.9% per 30 hari (rolling)
- **SRE Workbook Ch.2:** "four-week rolling window to be a good general-purpose interval" (28 hari) + weekly/quarterly pelengkap
- **Assessment:** Bukan kontradiksi; keduanya rolling window. 28 hari dipilih agar tiap window punya jumlah weekend identik (hindari bias weekday/weekend). 30 hari lebih intuitif untuk konteks bisnis bulanan. Keduanya valid; trade-off: LAB lebih umum dipahami, GOOGLE lebih stabil untuk traffic mingguan.
- **Impact:** LOW. Perlu dokumentasikan window yang dipakai.

## Contradiction 4: Tidak Ada Kontradiksi Material Lain
- Klaim inti (SLI=indikator, SLO=target, Error Budget=1−SLO, hindari CPU sebagai SLO, percentile > average, burn rate untuk alerting) konsisten lintas Ch.3/Ch.4/Ch.6/Workbook Ch.2/Ch.5.
- Tidak ada sumber yang merekomendasikan 100% sebagai target atau average latency sebagai SLI utama.

## Verifikasi Hitungan Latihan Lab
- Payment Webhook SLO 99.99% dari 200.000 request → budget = 0.01% × 200.000 = **20 request boleh gagal**.
- 15 gagal = 15/20 = **75% error budget habis**.
- Cross-check: SRE Book aggregate availability formula (successful/total) menghasilkan angka identik. Konsisten.
