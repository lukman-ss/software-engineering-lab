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
