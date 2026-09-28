# Chaos Engineering & Fault Injection: Membuktikan Ketahanan Sistem Secara Empiris

## Problem
Layanan terdistribusi sering mengasumsikan komponen eksternal (seperti payment gateway, database, atau microservice downstream) selalu tersedia. Saat kegagalan nyata terjadi di production, ketiadaan perlindungan seperti timeout, circuit breaker, dan fallback memicu *cascading failures* yang melumpuhkan seluruh sistem. Pengujian konvensional (unit test atau integration test standar) hanya menguji jalur normal atau mock statis, gagal membuktikan apakah sistem benar-benar dapat bertahan di bawah kondisi turbulen.

## Why This Matters
Kegagalan di sistem terdistribusi tidak terelakkan. Tanpa pengujian disengaja (*fault injection*), ketahanan sistem hanyalah tebakan teoretis. Chaos Engineering mentransformasikan ketahanan dari asumsi pasif menjadi bukti empiris yang terukur, memastikan arsitektur dapat menyerap gangguan tanpa merusak pengalaman pengguna.

## Mental Model
Chaos Engineering adalah eksperimen ilmiah pada perangkat lunak:
1. **Steady State**: Tentukan perilaku normal sistem melalui metrik output eksternal (misal error rate dan throughput), bukan metrik internal komponen.
2. **Hipotesis**: "Kami menduga steady state akan tetap terjaga meskipun gangguan X disuntikkan."
3. **Eksperimen**: Suntikkan gangguan terkontrol (latensi atau error).
4. **Verifikasi & Blast Radius**: Pantau metrik secara kontinu. Jika threshold terlewati, eksperimen wajib dihentikan otomatis (*auto-abort*) untuk melindungi pengguna.

## Core Concept
- **Fault Injection**: Menyisipkan latensi atau error buatan secara terprogram pada titik batas layanan downstream.
- **Circuit Breaker**: State machine (`Closed`, `Open`, `Half-Open`) yang mendeteksi lonjakan kegagalan dan mencegah pemanggilan berulang ke layanan yang sedang down.
- **Graceful Degradation**: Mekanisme fallback yang memberikan respons alternatif (misal antrean atau cache) ketika pemanggilan utama gagal, menjaga metrik error rate pengguna tetap prima.
- **Blast Radius Control**: Pembatasan dampak eksperimen melalui pemantauan metrik ketat dan pembatalan otomatis (*abort*).

## Failure Scenario
Tanpa Circuit Breaker, Fallback, dan Auto-Abort:
- Kegagalan downstream menyebabkan thread atau koneksi habis (*resource exhaustion*), memicu latensi kumulatif dan error berantai ke layanan hulu.
- Eksperimen chaos yang berjalan tanpa batas waktu atau pemantauan metrik akan terus merusak trafik, melanggar batas *blast radius*.

## How It Works
Arsitektur laboratorium ini terdiri dari empat komponen utama di dalam direktori `internal/`:
1. **`Injector` (`internal/fault/injector.go`)**: Menyisipkan latensi (menggunakan `time.After` dengan dukungan pembatalan konteks `ctx.Done()`) atau error paksa (`ErrInjectedFault`) secara aman dan thread-safe.
2. **`Monitor` (`internal/monitor/monitor.go`)**: Melacak total permintaan, permintaan sukses, dan permintaan gagal secara atomik (`sync/atomic`), serta mengevaluasi apakah error rate melampaui batas maksimum yang diizinkan.
3. **`CircuitBreaker` (`internal/circuitbreaker/circuitbreaker.go`)**: Mengelola transisi state (`Closed` → `Open` → `Half-Open` setelah cooldown) dan mengeksekusi fungsi `fallback` saat state terbuka atau terjadi error.
4. **`Experiment` (`internal/experiment/runner.go`)**: Menjalankan siklus eksperimen, memeriksa kesehatan sistem secara periodik via ticker, dan melakukan auto-abort serta netralisasi injeksi seketika saat steady-state breach terdeteksi.

## Architecture
```
[Client Request] 
      │
      ▼
[Resilient Client] ──(Circuit Breaker)──> [Downstream Service]
      │                                           │
      ├─(On Failure / Open)──> [Fallback]        (Fault Injector)
      │                                           │
      ▼                                           ▼
[Steady State Monitor] <──(Records Success/Fail)──┘
      │
      ├─(If Error Rate > Threshold)──> [Auto-Abort Experiment & Clear Injector]
```

## Implementation
Seluruh komponen diimplementasikan murni menggunakan pustaka standar Go (Go 1.22+) tanpa ketergantungan pihak ketiga, menjamin kesederhanaan, portabilitas, dan performa tinggi yang diverifikasi oleh race detector (`go test -race`).

## Code Walkthrough
Berikut adalah cuplikan kode penting dari implementasi:

1. **Injeksi Fault (`internal/fault/injector.go`)**:
Menggunakan `sync.RWMutex` untuk membaca konfigurasi secara aman dan mendukung pembatalan konteks saat latensi disuntikkan.
2. **State Machine Circuit Breaker (`internal/circuitbreaker/circuitbreaker.go`)**:
Mengatur ambang batas kegagalan (`threshold`) dan waktu jeda (`cooldown`) sebelum memasuki mode `Half-Open`.
3. **Monitor Steady-State (`internal/monitor/monitor.go`)**:
Menghitung rasio kesalahan secara atomik dengan syarat minimum sampel (5 request) sebelum evaluasi kesehatan.

## What the Tests Prove
Pengujian unit dan konurensi (`tests/chaos_test.go`) membuktikan secara empiris:
- `TestFaultInjector`: Injector aktif, menyuntikkan latensi dan error, serta kembali bersih setelah `Clear()`.
- `TestCircuitBreakerStateTransitions`: Transisi state dari `Closed` ke `Open` setelah threshold kegagalan tercapai, pemanggilan cepat menghasilkan `ErrCircuitOpen`, cooldown memicu `Half-Open`, dan eksekusi sukses mengembalikan ke `StateClosed`.
- `TestCircuitBreakerGracefulDegradation`: Fallback berhasil menelan error dan mengembalikan respons alternatif.
- `TestExperimentAutoAbortOnSteadyStateViolation`: Peningkatan error rate di atas threshold memicu pembatalan eksperimen otomatis dan menetralkan injector seketika.
- `TestConcurrencyAndRace`: Seluruh operasi thread-safe di bawah pengujian `-race` Go.

## Recovery / Rollback
Ketika eksperimen chaos selesai atau di-abort, fungsi `terminate()` langsung memanggil `injector.Clear()` secara sinkron untuk menetralkan gangguan. Trafik normal pulih seketika dan Circuit Breaker otomatis kembali ke state `Closed` setelah menerima permintaan sukses pada mode `Half-Open`.

## Production Considerations
- **Sliding Window Metrics**: Lab menggunakan penghitungan kumulatif sederhana; di production, gunakan *sliding time window* (misal metrik 1 menit terakhir) agar pemulihan sistem tidak terbebani riwayat kegagalan lama.
- **Metrik Percentile Latency**: Selain error rate, pantau latensi p95/p99 untuk mendeteksi *incipient degradation* sebelum error meledak.
- **Canary & Blast Radius**: Mulai eksperimen dari subset kecil instance atau canary deployment, bukan langsung ke seluruh armada production.

## Common Mistakes
- **Menjalankan Chaos Tanpa Auto-Abort**: Menyuntikkan gangguan tanpa pemantauan metrik otomatis dan tombol abort darurat.
- **Mengukur Kondisi Internal**: Mengukur metrik CPU/memory node alih-alih metrik kepuasan/perilaku pengguna akhir (*steady state output*).
- **Mengabaikan Fallback**: Menggunakan Circuit Breaker tanpa fallback, sehingga kegagalan downstream tetap langsung berdampak pada pemanggil.

## Case Study
Demonstrasi interaktif (`cmd/demo/main.go`) mensimulasikan layanan pembayaran:
1. **Baseline**: 5 request sukses (`CLOSED`, error rate 0%).
2. **Chaos dengan Mitigasi**: 10 request di bawah gangguan downstream dengan Resilient Client (Circuit Breaker + Fallback). Hasil: Berhasil memproses fallback ("Payment Queued"), error rate 0.00%, state CB beralih dari CLOSED ke OPEN.
3. **Chaos Tanpa Mitigasi**: Injeksi error murni memicu pelanggaran threshold (lab example: 20%). Eksperimen mendeteksi pelanggaran, melakukan `ABORT`, dan menetralkan injector secara otomatis.
4. **Recovery**: Trafik pasca-eksperimen kembali normal (`CLOSED`, sukses penuh).

## Checklist
- [x] Tentukan metrik steady-state (error rate).
- [x] Rumuskan hipotesis eksperimen.
- [x] Batasi blast radius dengan durasi dan auto-abort.
- [x] Bungkus pemanggilan dengan Circuit Breaker dan Fallback.
- [x] Verifikasi dengan unit test dan race detector (`go test -race ./...`).

## Key Takeaways
1. Chaos Engineering adalah ilmu eksperimen terstruktur, bukan pengujian acak.
2. Steady-state harus mengukur output perilaku sistem, bukan komponen internal.
3. Circuit Breaker mencegah cascading failures dengan memutuskan panggilan saat gagal.
4. Fallback menyediakan graceful degradation untuk menjaga pengalaman pengguna.
5. Auto-abort adalah pengaman wajib untuk membatasi blast radius eksperimen.

## Sources
- Principles of Chaos Engineering (`https://principlesofchaos.org/`)
- AWS Well-Architected Framework - Reliability Pillar (`https://docs.aws.amazon.com/wellarchitected/latest/reliability-pillar/chaos-engineering.html`)
- Netflix TechBlog (`https://netflixtechblog.com/tag/chaos-engineering`)
- Google SRE Book (`https://sre.google/sre-book/table-of-contents/`)
- Lab Implementation & Tests (`labs/34-chaos-engineering/`)
