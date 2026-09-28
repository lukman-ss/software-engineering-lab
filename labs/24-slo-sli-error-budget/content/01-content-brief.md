# Content Brief

Topic:
Penerapan SLI, SLO, Error Budget, dan Burn Rate Alerting berdasarkan prinsip Google SRE, diimplementasikan dalam Go.

Target Reader:
Insinyur perangkat lunak dan SRE pemula yang ingin memahami dan mengimplementasikan error budgeting untuk pengambilan keputusan keandalan layanan.

Problem:
Tim keandalan seringkali mengukur "kesehatan" layanan dengan metrik infrastructure (CPU, memory) yang tidak berhubungan langsung dengan pengalaman pengguna. Akibatnya, tim tidak dapat membedakan antara "layanan sebenarnya error" vs "server over-provisioned". Error budget memberi kerangka angka untuk memutuskan kapan harus berhenti merilis dan fokus pada keandalan.

Core Mental Model:
Keandalan ≠ 100%. SLO mendefinisikan target yang realistis (<100%). Error budget (= 1 - SLO) adalah "uang" yang bisa "habiskan" untuk inovasi. Burn rate mengukur laju pengeluaran budget; jika laju tinggi, berhenti dan perbaiki.

Approved Research Status:
APPROVED

Approved Engineering Status:
APPROVED — internal engineering audit: no non-blocking findings; open-source engineering audit: APPROVED with 2 non-blocking warnings (LatencyThreshold dead field, per-rule window fields unimplemented) — disclosed below. The two audits differ in scope/perspective; both are approved.

Main Concepts:
1. SLI (Service Level Indicator): ukuran kuantitatif pengalaman pengguna (mis. rasio good/total events).
2. SLO (Service Level Objective): target nilai untuk SLI (mis. 99.9% dalam 30 hari).
3. Error Budget: 1 - SLO. Selisih antara target dan actual.
4. Burn Rate: laju konsumsi error budget (actual_error_rate / allowed_error_rate).
5. Multi-Window Multi-Burn-Rate Alerting: periksa jendela pendek dan panjang secara bersamaan untuk mengurangi false positive.
6. Criticality Bucketing: endpoint kritis (Payment 99.9%) vs non-kritis (Reports 95.0%).

Verified Behaviors:
- WindowTracker merekam event dengan timestamp out-of-order dan mengembalikan agregasi total/good/bad yang benar.
- WindowTracker mengeviksi bucket di luar jendela rolling secara otomatis.
- Evaluator menghitung SLI = good / total; SLI = 1.0 ketika tidak ada trafik.
- Evaluator menghitung totalErrorBudget = (1 - TargetSLO) * total; budgetRemaining = totalErrorBudget - bad.
- Evaluator mengembalikan CanDeploy = false ketika budgetRemaining <= 0 (dan total > 0).
- AlertEngine menghitung burn rate = actual_error_rate / allowed_error_rate.
- AlertEngine memicu alert ketika burn rate >= threshold pada BOTH short dan long window.
- Concurrency: 20 goroutine × 100 request tidak menyebabkan race condition.
- Demo: baseline 1000 request 100% sukses → budget remaining positif; insiden 100 request 10 error (10%) menghabiskan budget, CanDeploy = false, slow burn alert (6.0x) terpicu.

Available Case Studies:
- Demo Phase 1: Traffic baseline 1000 request, 0 error → SLI 100%, budget tersisa.
- Demo Phase 2: 10 error dalam 100 request insiden ditambah 1000 baseline → total 1100 request, 10 bad (error rate 0.91% terhadap total) → burn rate ~9.09x, budget habis, CanDeploy = false.
- Demo Phase 3: Burn rate alert (TICKET 6.0x) terpicu pada short & long window.
- Demo Phase 4: Payment (99.9%) vs Reports (95.0%) dengan traffic error sama → kedua budget habis, tetapi Reports memiliki toleransi 5% yang jauh lebih lebar.
- Test: Transient spike hanya pada short window (100x) tapi long window clean (0.1x) → alert TIDAK terpicu.

Warnings:
- Semua sumber riset berasal dari ekosistem Google SRE, Datadog, Prometheus, GCP; tidak ada sumber independen (AWS, Azure, CNCF) untuk burn rate thresholds spesifik.
- Angka "70% outages dari change" adalah observasi internal Google tanpa verifikasi independen.
- Cost kenaikan ~100x per nine adalah heuristic, bukan hukum matematis.
- Implementasi menggunakan penyimpanan in-memory; metrics hilang saat proses restart.
- Burn rate thresholds (14.4x, 6.0x) adalah rekomendasi Google, bukan standar universal.
- Window demo menggunakan kompresi waktu: 30 menit mensimulasikan 30 hari (bukan jendela 30 hari sebenarnya); 28 hari/4 minggu rekomendasi Google keduanya valid dengan trade-off.
- Latency tidak diukur sebagai histogram penuh; hanya boolean good/bad predicate (WindowTracker.Bucket hanya TotalCount/GoodCount/BadCount).
- Zero traffic: Evaluator mengembalikan SLI=1.0, CanDeploy=true untuk hindari false freeze.
