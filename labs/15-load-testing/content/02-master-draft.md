# Load Testing dan Analisis Bottleneck Saturasi Sumber Daya

## Problem
Banyak tim rekayasa perangkat lunak mengandalkan metrik rata-rata waktu respons (*average response time*) untuk menilai performa sistem di bawah beban. Pendekatan ini menimbulkan ilusi stabilitas: rata-rata menyamarkan lonjakan latensi ekor (*tail latency outliers*). Ketika 95% request selesai dalam 20ms tetapi 5% sisanya tertahan selama 2000ms karena antrean sumber daya, rata-rata matematika tetap tampak dapat diterima. Akibatnya, degradasi sistem tidak terdeteksi hingga terjadi kegagalan fatal di lingkungan produksi.

## Why This Matters
Pengujian beban (*load testing*) bertujuan memetakan batas daya tahan sistem dan perilaku degradasinya, bukan sekadar membuktikan bahwa aplikasi berjalan cepat dalam kondisi ideal. Tanpa pengujian bertahap dan observasi persentil, tim tidak dapat memprediksi komponen downstream mana—seperti connection pool basis data atau thread pool aplikasi—yang pertama kali mencapai titik saturasi (*breaking point*).

## Mental Model
Model mental pengujian beban berpusat pada hubungan antara kapasitas sumber daya tetap dan pertumbuhan antrean konkuren:
- **Smoke Stage**: Beban konkurensi berada jauh di bawah kapasitas penanganan paralel. Tidak ada antrean; latensi murni ditentukan oleh durasi komputasi dasar. Rata-rata dan P95 bernilai identik.
- **Saturation Point**: Tingkat konkurensi menyamai batas penanganan sumber daya (misalnya ukuran koneksi basis data).
- **Stress Stage**: Tingkat konkurensi melampaui batas penanganan. Kelebihan permintaan dipaksa mengantre. Waktu tunggu bertambah secara kumulatif, memicu lonjakan eksponensial pada P95 dan P99, sementara rata-rata hanya meningkat moderat.

## Core Concept
1. **Inkremental Testing**: Pengujian wajib dimulai dari *smoke test* berkonkurensi rendah (2–5 Virtual Users) untuk memvalidasi integritas skrip dan konfigurasi, sebelum dielevasi ke beban normal (*load test*), beban puncak (*stress test*), lonjakan mendadak (*spike test*), atau durasi panjang (*soak test*).
2. **Distribusi Persentil**:
   - **P50 (Median)**: Pengalaman 50% pengguna tipikal.
   - **P90 / P95**: Pengalaman pengguna pada kuartil teratas; indikator utama degradasi antrean.
   - **P99**: Kasus ekstrem ekor distribusi; menunjukkan penumpukan antrean parah, garbage collection, atau locking.
3. **Korelasi Metrik Sisi Klien vs Sisi Server**: Metrik latensi, RPS, dan error rate dari generator beban harus selalu dikorelasikan dengan utilisasi CPU, memori, I/O, serta saturasi pool koneksi pada server target.

## Failure Scenario
Ketika endpoint transaksional menerima beban konkuren yang melebihi kapasitas *connection pool* basis data:
1. Slot koneksi terpakai penuh oleh sejumlah request pertama.
2. Request berikutnya terhenti di memori menunggu slot dibebaskan.
3. Waktu tunggu antrean terakumulasi ke total durasi pemrosesan.
4. Latensi P95 melonjak tajam melampaui ambang batas toleransi (misalnya melompat dari 20ms ke >200ms).
5. Jika waktu tunggu melampaui batas timeout klien, request gagal dengan error 5xx atau koneksi diputus.

## How It Works
Sistem pengujian beban pada lab ini mengimplementasikan dua bagian terpisah:
1. **Mock Server (`internal/server`)**: Menyediakan endpoint `POST /booking`. Kapasitas koneksi basis data dimodelkan menggunakan buffered channel (semafor) berukuran tetap (default 5 koneksi), di mana setiap transaksi menahan slot selama durasi tertentu (default 20ms).
2. **Load Runner (`internal/loadtest`)**: Mengorkestrasi *Virtual Users* (VUs) independen menggunakan goroutine, mengeksekusi request HTTP berulang kali selama durasi yang ditentukan, mengumpulkan durasi latensi ke dalam slice privat per-VU guna menghindari overhead mutex, dan menghitung ringkasan statistik persentil setelah seluruh goroutine selesai.

## Architecture
Komponen lab dirancang tanpa dependensi eksternal:

```text
+-------------------------------------------------------------+
|                     Load Test Runner                        |
|                                                             |
|  [VU 1] -----> HTTP POST /booking                           |
|  [VU 2] -----> HTTP POST /booking                           |
|  ...                                                        |
|  [VU N] -----> HTTP POST /booking                           |
|                                                             |
|  (Tiap VU mencatat latensi ke memory slice independen)     |
+------------------------------+------------------------------+
                               |
                               v
+-------------------------------------------------------------+
|                  HTTP Server (/booking)                     |
|                                                             |
|        +-------------------------------------------+        |
|        | Semafor Penampung Koneksi (Max = 5 Slot)  |        |
|        +-------------------------------------------+        |
|            | Slot 1 | Slot 2 | Slot 3 | Slot 4 | Slot 5     |
|                                                             |
|   (Request ke-6 dan seterusnya tertahan mengantre)          |
|   (Pemrosesan simulasi query memakan durasi 20ms)           |
+-------------------------------------------------------------+
```

## Implementation
Struktur modul lab terdiri dari:
- `internal/server/server.go`: Mengimplementasikan server HTTP dengan semafor kanal Go.
- `internal/loadtest/runner.go`: Mengimplementasikan runner beban konkuren dengan transport HTTP kustom (`MaxIdleConns: 1000`) untuk mencegah limitasi pooling sisi klien menyamarkan bottleneck server.
- `internal/loadtest/metrics.go`: Mengimplementasikan fungsi penghitungan min, max, avg, P50, P90, P95, dan P99 dari slice latensi terurut.
- `cmd/demo/main.go`: Menjalankan perbandingan Smoke Test (2 VU) dan Stress Test (50 VU) melawan server berkapasitas 5 koneksi.

## Code Walkthrough

### 1. Pembatasan Kapasitas Menggunakan Semafor (`internal/server/server.go`)
```go
func (s *Server) handleBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	atomic.AddInt64(&s.activeReq, 1)
	defer atomic.AddInt64(&s.activeReq, -1)

	// Acquire DB connection slot (simulates DB connection pool limit)
	s.semaphore <- struct{}{}
	time.Sleep(s.cfg.DBQueryDuration)
	<-s.semaphore

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(BookingResponse{
		Status:    "confirmed",
		BookingID: "BK-1001",
	})
}
```
Ketika 5 goroutine mengisi `s.semaphore`, goroutine berikutnya akan memblokir (*block*) pada baris `s.semaphore <- struct{}{}` sampai ada request sebelumnya yang membaca kanal pada baris `<-s.semaphore`.

### 2. Eksekusi Beban Tanpa Kontensi Mutex (`internal/loadtest/runner.go`)
```go
for i := 0; i < r.cfg.VUs; i++ {
	wg.Add(1)
	go func(vuID int) {
		defer wg.Done()
		var lats []time.Duration
		var errs int

		for {
			select {
			case <-ctx.Done():
				results[vuID] = vuResult{latencies: lats, errors: errs}
				return
			default:
				// HTTP request execution and latency tracking
```
Setiap virtual user menulis data ke slice miliknya sendiri (`lats`). Agregasi metrik hanya dilakukan satu kali di thread utama setelah `wg.Wait()` selesai, memastikan eksekusi bebas dari kontensi kunci sinkronisasi.

### 3. Perhitungan Persentil (`internal/loadtest/metrics.go`)
```go
sorted := make([]time.Duration, len(latencies))
copy(sorted, latencies)
sort.Slice(sorted, func(i, j int) bool {
	return sorted[i] < sorted[j]
})
// ...
res.P95Latency = percentile(sorted, 95)
res.P99Latency = percentile(sorted, 99)
```
Data latensi disortir secara ascending, lalu indeks persentil diambil berdasarkan posisi peringkat:
`idx := int(float64(len(sorted)-1) * (pct / 100.0))`.

## What the Tests Prove
Pengujian otomatis pada `tests/loadtest_test.go` dan `internal/loadtest/metrics_test.go` membuktikan:
1. **Akurasi Penghitungan Statistik**: `TestCalculateMetrics` memastikan kalkulasi Min, Max, Average, P50, P90, P95, dan P99 menghasilkan nilai eksak sesuai distribusi sampel yang diuji.
2. **Degradasi P95 pada Kondisi Stres**: `TestLoadTest_SmokeVsStress` membuktikan secara deterministik bahwa latensi P95 saat stress test (50 VU) lebih besar daripada smoke test (2 VU) ketika berhadapan dengan limit 5 koneksi server (`assert stressRes.P95Latency > smokeRes.P95Latency`).
3. **Keamanan Konkurensi**: Seluruh eksekusi lolos deteksi race detector (`go test -race ./...`).

Hasil eksekusi riil dari `cmd/demo`:
- **Smoke Test (2 VUs, kapasitas 5)**:
  - Average: ~21.7ms
  - P50: ~21.5ms
  - P95: ~22.6ms
  - P99: ~23.3ms
  - Error: 0
- **Stress Test (50 VUs, kapasitas 5)**:
  - Average: ~201.2ms
  - P50: ~200.7ms
  - P95: ~213.9ms
  - P99: ~216.2ms
  - Error: 0

Angka ini membuktikan bahwa saat beban melampaui kapasitas pool sebesar 10 kali lipat, waktu antrean melonjak tajam (naik ~10x lipat dari ~21ms ke >210ms).

## Recovery / Rollback
Ketika hasil uji beban mengidentifikasi bottleneck connection pool:
1. **Right-sizing Connection Pool**: Tingkatkan batas pool basis data sejauh memori dan CPU server database mendukung batas koneksi paralel tersebut.
2. **Rate Limiting / Load Shedding**: Terapkan rate limiter di API gateway atau HTTP middleware untuk menolak request berlebih dengan respons `429 Too Many Requests` atau `503 Service Unavailable` daripada membiarkannya menumpuk di antrean tak terbatas.
3. **Queue Timeouts**: Konfigurasikan batas timeout pengambilan koneksi (*pool acquire timeout*) pada driver basis data agar request lekas gagal daripada menyebabkan resource exhaustion menyeluruh.

## Production Considerations
- **Memori Perhitungan Persentil**: Pendekatan sorting slice (`sort.Slice`) dalam implementasi lab membutuhkan memori linier terhadap jumlah request ($O(N)$). Di lingkungan produksi dengan jutaan request, gunakan algoritma histogram streaming seperti `HdrHistogram` atau `t-digest` untuk menghemat memori.
- **Isolasi Lingkungan Uji**: Uji beban skala penuh tidak boleh dijalankan langsung di database produksi aktif tanpa isolasi data yang ketat.
- **Kapasitas Generator Beban**: Pastikan mesin runner tidak mengalami saturasi CPU, network socket exhaustion, atau pembatasan client connection pool (`MaxIdleConnsPerHost`) yang dapat menimbulkan hasil uji palsu (*false bottleneck*).

## Common Mistakes
1. **Mengabaikan Tahap Smoke Test**: Langsung menjalankan ratusan atau ribuan VU sehingga skrip yang salah konfigurasi memicu kegagalan tanpa mengetahui baseline yang benar.
2. **Hanya Mengukur Average Response Time**: Menganggap sistem sehat karena rata-rata latensi rendah, padahal sebagian pengguna mengalami antrean ekstrem.
3. **Mengabaikan Metrik Sisi Server**: Mencatat latensi tinggi dari sisi klien tanpa memantau metrik internal server (CPU, RAM, koneksi DB), sehingga penyebab pasti bottleneck tidak dapat didiagnosis.
4. **Batas Idle Connection Client HTTP**: Menggunakan klien HTTP default yang membatasi konkurensi koneksi keluar (seperti default `DefaultTransport.MaxIdleConnsPerHost = 2` pada Go), sehingga antrean terjadi di sisi klien, bukan di server target.

## Case Study
Dalam skenario sistem transaksional seperti **Booking Bengkel**:
- **Critical Endpoints**: Pengujian difokuskan pada endpoint yang memutasi status seperti `POST /booking`, `POST /payment`, dan `POST /invoice`, bukan hanya endpoint pembacaan data statis (`GET /branches`).
- **Tahapan Beban**: Dimulai dari 5–10 VU untuk validasi fungsional (smoke), lalu naik ke 50–100 VU untuk beban tipikal, hingga 800 VU untuk pengujian stres.
- **Diagnosis Lonjakan P95**: Jika latensi P95 melonjak tajam dari 300ms ke 2.5s pada 800 VU, investigasi dilakukan secara terstruktur:
  1. Periksa APM/distributed tracing untuk melihat span pemrosesan yang membengkak.
  2. Periksa apakah connection pool basis data telah mencapai kapasitas 100% atau mengalami table lock.
  3. Periksa apakah dependensi eksternal (payment gateway atau WhatsApp webhook) mengalami perlambatan atau *rate limiting*.

## Checklist
- [ ] Mulai pengujian dengan Smoke Test (2–5 VU).
- [ ] Catat metrik persentil: P50, P90, P95, P99; jangan hanya mengandalkan rata-rata.
- [ ] Amati metrik saturasi server (CPU, Memory, Connection Pool) bersamaan dengan metrik klien.
- [ ] Pastikan generator beban memiliki kapasitas koneksi idle yang memadai (`MaxIdleConns`).
- [ ] Hindari race conditions pada pengumpul metrik beban.
- [ ] Tentukan ambang batas SLA/SLO berbasis persentil untuk kriteria lulus/gagal (*pass/fail threshold*).

## Key Takeaways
- Rata-rata latensi mengaburkan kegagalan ekor distribusi; persentil P95 dan P99 wajib digunakan untuk mengevaluasi performa riil.
- Pengujian beban harus bertahap: mulai dari smoke test untuk verifikasi integritas, kemudian dinaikkan ke beban target.
- Antrean pada sumber daya terbatas (seperti database connection pool) menyebabkan peningkatan latensi yang eksponensial.
- Generator beban harus dirancang thread-safe tanpa kontensi lock internal yang dapat mengaburkan hasil pengukuran.

## Sources
- Types of load testing — Grafana Labs (k6 Documentation): https://k6.io/docs/test-types/
- What is Azure Load Testing? — Microsoft Learn: https://learn.microsoft.com/en-us/azure/load-testing/overview-what-is-azure-load-testing
