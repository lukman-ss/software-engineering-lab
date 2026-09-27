# Content Brief

Topic: Load testing untuk aplikasi Booking Bengkel — praktik terbaik, tools, metrics, dan strategi identifikasi bottleneck pada aplikasi software.
Target Reader: Software engineer yang bertanggung jawab untuk performa sistem, SRE, dan tim QA yang perlu memastikan aplikasi tetap stabil di bawah beban nyata.
Problem: Tim sering menganggap aplikasi "siap produksi" hanya karena lolos functional test, tanpa memverifikasi perilaku saat banyak pengguna mengakses secara bersamaan. Ini menghasilkan degradasi performa, timeout, atau kegagalan sistem pada saat kritis seperti promo besar atau go-live.
Core Mental Model: Load testing bukan sekadar memastikan server tidak crash; ia mengekspos bagaimana resource terbatas (koneksi database, thread, memori) membuat latensi meningkat secara non-linear — rata-rata menutupi lonjakan ekstrim pada persentil tinggi (P95, P99) karena antrian di balik batas kapasitas.
Approved Research Status: APPROVED (research-audit/07-verdict.md)
Approved Engineering Status: APPROVED (engineering-audit/06-verdict.md)
Main Concepts:
- Enam jenis tes performa: smoke, load (average-load), stress, spike, soak/endurance, breakpoint
- Metrics kunci: P50/P95/P99 response time, error rate, RPS, resource utilization (CPU, memory, koneksi database)
- Metodologi identifikasi bottleneck: monitoring metrics tiap layer (aplikasi, database, eksternal API) dan korelasikan dengan penurunan response time
- Tool pilihan: k6 (JavaScript), JMeter (GUI/XML), Locust (Python), Gatling (Scala/JVM) — pilih sesuai kebutuhan tim
- Pitfall umum: hanya menguji endpoint /health, data dummy terlalu sedikit, tidak memantau server, tidak menentukan target performa
- Timing SDLC: sebelum go-live, sebelum promosi besar, setelah optimasi besar, setelah perubahan infrastruktur penting
Verified Behaviors:
- Pada beban rendah (VUs < kapasitas koneksi database), semua metrik latensi (min, avg, P50, P95, P99) berada di kisaran normal dan berkisar dekat satu sama lain.
- Pada beban tinggi (VUs >> kapasitas), P95 dan P99 latency naik jauh lebih drastis dibanding rata-rata karena antrian pembatasan resource (misalnya: koneksi database habis).
- Calculator metrik mengukur persentil secara akurat melalui pengurutan latensi (menggunakan sort standar library Go).
- Harness beban menggenerate trafik konkuren tanpa bottleneck sendiri melalui custom HTTP transport dengan MaxIdleConns tinggi.
- Implementasi server mensimulasikan batas koneksi database dengan buffered channel (semaphore) sehingga ketika slot habis, request baru menunggu dalam antrian.
- Semua integrasi test lolos tanpa race condition (`go test -race ./...` bersih).
Available Case Studies:
- Demo aplikasi Booking Bengkel dalam `cmd/demo/main.go`: menunjukkan kontras jelas antara smoke test (2 VUs) dan stress test (50 VUs) terhadap server dengan kapasitas koneksi 5.
- Laporan eksekusi engineering mencatat contoh output demo: pada smoke test P95 ≈ 21ms, pada stress test P95 ≈ 1.35s (naik 64x) walaupun rata-rata hanya naik 35x.
Warnings:
- Metrik dalam demo bersifat ilustratif; nilai aktual tergantung pada hardware host, jadwal CPU, dan pause GC.
- Formula hitung pengguna konkuren (sessions per jam × durasi rata-rata sesi / 3600) membutuhkan adaptasi think-time untuk alur kerja multistep seperti Booking Bengkel (login → pilih cabang → bayar → konfirmasi WhatsApp).
- Tidak ada ambang batas universal untuk P95 atau error rate; target harus diturunkan dari SLA bisnis dan penelitian pengalaman pengguna.
- Implementasi lab disederhanakan: menggunakan `time.Sleep` dan buffered channel sebagai pengganti kontensi database nyata atau degradasi CPU.
- Dokumentasi resmi JMeter (Source 16) mengalami timeout saat verifikasi jaringan, meskipun charakteristik arsitekturnya secara luas dikenal.
- Metrik persentil (P50/P90/P95/P99) hanya mencakup request sukses (HTTP 2xx); request gagal dihitung sebagai error tanpa latensi tercatat.
- Demo mencetak subset P50/P95/P99, walau `Result` menghitung P90 juga.