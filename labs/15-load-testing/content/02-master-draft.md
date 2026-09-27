# Load Testing: Menemukan Batas Sistem Sebelum Pengguna Menemukannya

## Problem

Tim sering menganggap aplikasi "siap produksi" hanya karena lolos functional test, tanpa memverifikasi perilaku saat banyak pengguna mengakses secara bersamaan. Ini menyebabkan degradasi performa tak terduga, timeout, atau kegagalan sistem pada saat kritis seperti promo besar, periode pembayaran gaji, atau go-live. Tanpa data nyata tentang bagaimana sistem berperilaku di bawah beban, keputusan deploy bersifat asumsi.

## Why This Matters

Load testing bukan sekadar proses tautologi — ia menjadi satu-satunya cara untuk memperkirakan dampak ekonomi dari kegagalan sistem. Sebuah aplikasi booking bengkel yang dapat menangani 99% permintaan dalam 50ms pada ragam normal, tapi mengalami 5 detik latency pada antrian 1000 permintaan sekaligus, dapat kehilangan ribuan pelanggan dan merusak reputasi. Persentil tinggi (P95, P99) memberi gambaran tentang pengalaman pengguna yang paling buruk, bukan rata-rata yang dapat menutupi masalah.

## Mental Model

Fokus pada tiga tonggak: (1) Baseline: beban rendah, semua resource cukup, latency stabil. (2) Saturation: beban melebihi kapasitas resource terbatas (misalnya: koneksi database), request baru menunggu dalam antrian. (3) Failure: resource habis total, menghasilkan error atau timeout. Penting untuk membedakan antara beban yang mengisi antrian (menyebabkan P95 meningkat) dan beban yang benar-benar crash sistem (menyebabkan error rate meningkat).

## Core Concept

### Enam Jenis Performance Test

Standar industri mengakui enam jenis utama performance test dengan definisi yang konsisten:

1. **Smoke Test**: Tes beban minimal (2-20 VUs, detik hingga menit) dijalankan setiap kali script dibuat/diperbarui. Tujuannya memvalidasi kebenaran skrip dan mengumpulkan baseline metrics.
2. **Load/Average Test**: Menyimulasikan traffic produksi normal dengan pola ramp-up/plateau/ramp-down. Mengukur kinerja sistem pada kondisi terduga sehari-hari.
3. **Stress Test**: Beban di atas rata-rata untuk menguji batas sistem. Harus dijalankan setelah load test lolos. Level beban bergantung pada profil risiko sistem (rush hour, payday, akhir minggu), bukan persentase tetap seperti 50% atau 100%.
4. **Spike Test**: Lonjakan tiba-tiba dengan minimal/no ramp-up. Digunakan untuk flash sale, peluncuran produk, atau kejutan musiman.
5. **Soak/Endurance Test**: Load test rata-rata yang diperpanjang menghari-hari (3-72 jam) untuk mendeteksi memory leak, resource leak, atau kehabisan storage.
6. **Breakpoint Test**: Peningkatan beban bertahap hingga sistem gagal untuk mengidentifikasi kapasitas maksimum dan titik kegagalan.

### Metrics Kunci

Senior engineer memantau kombinasi metrics berikut sebagai standar industri:

- **Response Time Percentiles**: P50, P95, P99 (bukan rata-rata, karena rata-rata menutupi tail latency spike)
- **Error Rate**: Persentase request yang gagal (HTTP 5xx, timeout, error aplikasi)
- **Requests Per Second (RPS)**: Throughput yang dihasilkan sistem
- **Resource Utilization**: CPU, memori, disk I/O, jaringan, serta koneksi database

Azure Well-Architected Performance Efficiency Pillar menekankan definisi performance targets yang mencakup response time, throughput, resource usage, dan stability. ISO/IEC 25010 (Performance Efficiency subcharacteristics) mengonfirmasi Time behaviour dan Resource utilisation sebagai atribut kualitas perangkat lunak.

## How It Works

### Identifikasi Bottleneck

Untuk membedakan bottleneck aplikasi vs database vs external API, metode adalah memonitor component-level metrics secara paralel dan mencari korelasi dengan degradasi response time:

- Peningkatan `http_req_waiting` (Time To First Byte) menunjukkan bottleneck pada server-side (aplikasi atau database)
- Peningkatan `http_req_connecting` menunjukkan network atau connection issues
- Latensi pada third-party API yang menunjukkan dependency eksternal

k6 memecah `http_req_duration` menjadi komponen-mikro: blocked, connecting, TLS handshaking, sending, waiting, receiving. Analisis ini memungkinkan isolasi lapisan mana yang menjadi penanggung jawab degradasi.

### Konfigurasi Semaphore untuk Simulasi Connection Pool

Server mock menggunakan buffered channel (semaphore) untuk mensimulasikan batas koneksi database. Kapasitas ditentukan oleh ukuran channel; ketika penuh, request baru menunggu sampai slot melepaskan. Penambahan `time.Sleep` pada duration query DB mensimulasikan latency operasi database sebenarnya.

## Architecture

```
Booking Server (constrained connection pool)
    ↓ POST /booking
Load Tester (multiple VUs)
    ↓ concurrent HTTP requests
Result Aggregator (collects latencies per VU)
    ↓ stats calculation
Metrics Output (Min, Max, Avg, P50, P90, P95, P99, RPS)
```

Komponen utama:
1. `BookingServer`: HTTP handler dengan concurrency limit (semaphore) mensimulasikan database connection pool
2. `LoadTester`: Concurrency orchestrator menghasilkan trafik HTTP dengan VUs yang ditentukan untuk durasi yang dipertimbangkan
3. `MetricsAggregator`: Buffer latensi per-VU dikumpulkan ke `CalculateMetrics` untuk menghitung persentil tanpa lock

## Implementation

Implementasi menggunakan standar library Go (`net/http`, `sync`, `time`) tanpa dependensi pihak ketiga. Pendekatan "ponytail" menggunakan exact sorting untuk persentil karena skala tes (kurang dari 10.000 sampel) masih terukur. Untuk benchmark ratusan ribu RPS atau durasi berjam-jam, histogram streaming seperti HdrHistogram lebih cocok.

Server menerima konfigurasi `MaxDBConnections` (default 5) dan `DBQueryDuration` (default 10ms). Setiap request mengakuisisi slot pada semaphore sebagai gantinya koneksi database sebenaranya. Jika banyak request masuk sekaligus, yang terakhir menunggu sampai ada slot yang melepaskan (defer). Penambahan `time.Sleep` pada duration query DB mensimulasikan latency operasi database sebenarnya. Dalam server.go:74-78 terdapat peningkatan acak 10% pada durasi query untuk mencerminkan beban sistem.

Load generator menggunakan `http.Transport` custom dengan `MaxIdleConns` dan `MaxIdleConnsPerHost` tinggi (1000) untuk memastikan bariknya bukanlah limit koneksi HTTP klien. Timeout klien ditetapkan 5 detik. Catatan: hanya request sukses (HTTP 2xx) yang merekam latency; request dengan error >= 400 ditulis sebagai error tanpa latency tercatat.

## Code Walkthrough

`internal/server/server.go`:
- `New(cfg Config)` membuat server dengan semaphore berukuran `MaxDBConnections`
- `handleBooking` mengakuisisi semaphore, menunggu konteks dibatalkan jika terlalu lama, kemudian menunggu `DBQueryDuration` sebagai simulasi query database
- `ActiveConnections()` mengembalikan jumlah slot yang sedang digunakan untuk monitoring durasi live

`internal/loadtest/runner.go`:
- `Run(ctx)` membuat multiple goroutine sebanyak `VUs`, masing-masing mengirim request secara loop hingga konteks selesai
- Setiap goroutine menyimpan latensi dalam slice terpisah untuk menghindari lock contention
- Setelah semua selesai, semua latensi dikumpulkan dan diteruskan ke `CalculateMetrics`

`internal/loadtest/metrics.go`:
- `percentile(sorted, pct)` menghitung indeks aray yang dapatan persentil dengan rumus `len(sorted)-1 * pct/100`
- `CalculateMetrics` mengembalikan struct dengan semua statistik yang dibutuhkan

## What the Tests Prove

Tests memverifikasi tiga hal utama:

1. **Akurasi persentil**: `TestCalculateMetrics` menggunakan 100 sampel terurut 1-100ms dan memverifikasi P50=50ms, P95=95ms, P99=99ms sesuai rumus index-based.

2. **Perbedaan smoke vs stress**: `TestLoadTest_SmokeVsStress` menunjukkan bahwa P95 stress test jauh lebih tinggi daripada P95 smoke test, membuktikan bahwa tail latency naik drastis ketika resource terebut.

3. **Invarian statistik**: `TestCalculateMetrics_Invariants` memverifikasi bahwa Min ≤ P50 ≤ P90 ≤ P95 ≤ P99 ≤ Max selalu terpenuhi.

4. **Tidak ada race condition**: Semua test lulus dengan `go test -race ./...` bersih, menunjukkan implementasi thread-safe.

## Failure Scenario

Di bawah beban ekstrem, antrian di kanal semaphore berpanjang, menghasilkan latency ekor (P95/P99) yang melonjak hampir seratus kali lipat dibanding rata-rata. Pada contoh demo: smoke test menghasilkan P95 ≈ 21ms, sedangkan stress test (50 VUs vs 5 koneksi) menghasilkan P95 ≈ 1.35s — kenaikan 64x. Error rate tetap nol karena server tidak crash; request hanya menunggu. Ini ilustrasi klasik bahwa "sistem masih berfungsi" tidak berarti "sistem berfungsi baik" bagi pengguna akhir.

## Production Considerations

### Common Mistakes

1. **Hanya menguji endpoint /health**: Endpoint init dsb tidak menyimulasikan beban transaksional nyata.
2. **Data dummy terlalu sedikit**: Volume data yang tidak merepresentasikan produksi menghasilkan hasil yang tidak realistis.
3. **Tidak memantau server**: Monitoring CPU, memori, database metrics dilakukan paralel untuk mengidentifikasi bottleneck.
4. **Tidak menentukan target performa**: Tanpa thresholds yang didefinisikan (misal: P95 < 500ms), tidak ada kriteria kegagalan yang jelas.
5. **Menguji di laptop**: Lingkungan dev tidak mencerminkan spesifikasi produksi (contoh: MacBook M4 vs server rack).
6. **Workload tidak realistis**: Pola traffic tidak mencerminkan perilaku pengguna sebenar.

### Timing dalam SDLC

Load testing dilakukan pada titik-titik kritis:
- Sebelum go-live
- Sebelum promosi besar
- Setelah optimasi apa pun (kode, konfigurasi, infrastruktur)
- Setelah perubahan database
- Setelah migrasi cloud
- Setelah mengubah arsitektur penting

Praktik terbaik adalah melakukan load testing secara konsisten sejak awal development hingga pre-go-live, bukan sekadar sebelum deploy.

## Checklist

- [ ] Tentukan test type yang dibutuhkan (smoke, load, stress, spike, soak, atau breakpoint)
- [ ] Kumpulkan baseline metrics pada load rendah (smoke test)
- [ ] Identifikasi resource bottleneck utama (CPU, memori, koneksi database, API eksternal)
- [ ] Tentukan thresholds pass/fail berdasarkan SLA bisnis (bukan angka standar universal)
- [ ] Gunakan data volume dan pola traffic yang realistis
- [ ] Jalankan tes di lingkungan yang mendekati production
- [ ] Pantau server-side metrics secara paralel selama tes
- [ ] Analisis korelasi antara penurunan performance dan resource usage
- [ ] Dokumentasikan temuan dan rencanakan mitigasi

## Key Takeaways

1. Load testing mengungkap batas kapasitas sistem melalui enam jenis tes utama: smoke, load, stress, spike, soak, dan breakpoint.
2. Persentil (P95, P99) lebih penting daripada rata-rata karena mengidentifikasi pengalaman pengguna yang buruk.
3. Bottleneck terdeteksi dengan memonitor metrics tiap lapisan secara paralel dan mencari korelasi.
4. Tool pilihan harus ditentukan oleh kebutuhan tim (k6 untuk JavaScript, Locust untuk Python, JMeter untuk GUI/XML, Gatling untuk JVM).
5. Smoke test membuktikan baseline; stress test mengungkapkan efek queuing pada tail latency.
6. Persentil (P50/P90/P95/P99) hanya mencakup request sukses (HTTP 2xx); request gagal dihitung sebagai error tanpa latency tercatat.
7. Thresholds harus ditentukan upfront berdasarkan SLA bisnis. Contoh P95 < 500ms adalah ilustratif, bukan standar universal.
8. Jangan uji endpoint /health atau data fixture kecil. Workload harus merepresentasikan trafik produksi nyata.
9. Load test di lingkungan mirip production. Laptop dev tidak bisa merepresentasikan spesifikasi server.
10. Demo mencetak subset P50/P95/P99, walaupun `Result` menghitung P90 juga.
11. 10% peningkatan latency query di server adalah amplifier di atas antrean, bukan penyebab utama degradasi.
12. Tidak ada keajaiban satu benchmark. Semua angka kontekstual tergantung arsitektur, workload, dan resource yang dikonsumsi.

## Sources

- Grafana k6 Documentation: Load test types, Thresholds, Built-in metrics
- Microsoft Azure Well-Architected Framework: Performance Efficiency Pillar
- Google SRE Book: Testing for Reliability (Chapter 17)
- ISO/IEC 25010: Software Quality Model (Performance Efficiency)
- Apache JMeter User Manual
- Locust Official Documentation
- Gatling Documentation