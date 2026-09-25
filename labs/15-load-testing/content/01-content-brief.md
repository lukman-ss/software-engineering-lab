# Content Brief

Topic: Load Testing, Analisis Persentil Latensi, dan Deteksi Bottleneck Saturasi Sumber Daya
Target Reader: Software Engineer, Backend Engineer, Platform / DevOps Engineer
Problem: Rata-rata latency (average response time) menyamarkan degradasi performa pada ekor distribusi (tail latency). Di bawah beban tinggi atau keterbatasan connection pool, antrean request menyebabkan lonjakan P95/P99 yang fatal bagi sebagian pengguna meski angka rata-rata tampak normal.
Core Mental Model: Beban konkurensi bertahap (Smoke vs Stress) memetakan batas saturasi sumber daya downstream (database connection pool); persentil latensi (P95, P99) mengekspos degradasi antrean tak linier yang disamarkan oleh rata-rata.
Approved Research Status: APPROVED
Approved Engineering Status: APPROVED
Main Concepts:
- Perbedaan pengujian inkremental: Smoke Test vs Stress Test
- Kegagalan metrik rata-rata (average dilution) vs keandalan persentil (P50, P90, P95, P99)
- Peniruan saturasi connection pool basis data menggunakan semafor
- Korelasi metrik sisi klien (latency, RPS, error rate) dengan metrik sisi server (utilisasi connection pool)
Verified Behaviors:
- Pada smoke test (2 VU vs 5 DB connections), request diproses tanpa antrean dengan P95 mendekati rata-rata (~22ms).
- Pada stress test (50 VU vs 5 DB connections), request mengantre di balik semafor, menyebabkan lonjakan P95/P99 signifikan (>200ms) dan membuktikan degradasi tail latency non-linier.
- Runner konkurensi mengukur metrik thread-safe tanpa race conditions (`go test -race ./...` lolos).
Available Case Studies:
- Sistem "Booking Bengkel" (`POST /booking`) dengan simulasi saturasi pool koneksi database 5 koneksi dan durasi kueri 20ms.
Warnings:
- Penghitungan persentil menggunakan sorting slice (`sort.Slice`), cocok untuk dataset lab (<10.000 sampel), namun butuh histogram streaming (misal HdrHistogram) untuk beban jutaan sampel jangka panjang.
- Penundaan kueri disimulasikan menggunakan `time.Sleep` dan semafor in-memory, bukan engine database nyata dengan lock contention sebenarnya.
- Belum ada automated unit test khusus untuk akumulasi status HTTP 4xx/5xx pada runner, meski fungsi terverifikasi bekerja via review manual.
