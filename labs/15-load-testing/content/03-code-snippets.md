# Code Snippets

## Snippet 1 — Server Semaphore untuk Simulasi Connection Pool

Source File: `internal/server/server.go:44-71`

Purpose: Mensimulasikan batas koneksi database menggunakan buffered channel (semaphore) sehingga request melebihi kapasitas harus menunggu.

```go
func New(cfg Config) *Server {
	if cfg.MaxDBConnections <= 0 {
		cfg.MaxDBConnections = 5
	}
	if cfg.DBQueryDuration <= 0 {
		cfg.DBQueryDuration = 10 * time.Millisecond
	}
	return &Server{
		cfg:       cfg,
		semaphore: make(chan struct{}, cfg.MaxDBConnections),
	}
}

func (s *Server) handleBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	atomic.AddInt64(&s.activeReq, 1)
	defer atomic.AddInt64(&s.activeReq, -1)

	select {
	case s.semaphore <- struct{}{}:
	case <-r.Context().Done():
		return
	}
	defer func() { <-s.semaphore }()
	// ... query simulation
}
```

Explanation: Semaphore diinisialisasi dengan kapasitas `MaxDBConnections`. Setiap request coba mengirim struct kosong ke channel; jika channel penuh, goroutine block sampai slot melepaskan (defer). Ini meniru perilaku connection pool database: request menunggu saat maksimum tercapai.

---

## Snippet 2 — Percentile Calculator Exact Sort

Source File: `internal/loadtest/metrics.go:44-72`

Purpose: Menghitung P50, P90, P95, P99 melalui pengurutan slice latensi.

```go
func CalculateMetrics(latencies []time.Duration, errors int, totalDuration time.Duration) Result {
	total := len(latencies) + errors
	if total == 0 {
		return Result{}
	}

	res := Result{
		TotalRequests: total,
		SuccessCount:  len(latencies),
		ErrorCount:    errors,
		Duration:      totalDuration,
	}

	if totalDuration > 0 {
		res.RPS = float64(total) / totalDuration.Seconds()
	}

	if len(latencies) == 0 {
		return res
	}

	sorted := make([]time.Duration, len(latencies))
	copy(sorted, latencies)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	var sum time.Duration
	for _, lat := range sorted {
		sum += lat
	}

	res.MinLatency = sorted[0]
	res.MaxLatency = sorted[len(sorted)-1]
	res.AvgLatency = sum / time.Duration(len(sorted))
	res.P50Latency = percentile(sorted, 50)
	res.P90Latency = percentile(sorted, 90)
	res.P95Latency = percentile(sorted, 95)
	res.P99Latency = percentile(sorted, 99)

	return res
}

func percentile(sorted []time.Duration, pct float64) time.Duration {
	if len(sorted) == 0 {
		return 0
	}
	idx := int(float64(len(sorted)-1) * (pct / 100.0))
	return sorted[idx]
}
```

Explanation: Setelah mengurutkan latensi, indeks persentil dihitung dengan rumus `(n-1) * pct/100`. Metode ini akurat untuk sample count kecil (kurang dari 1M). Komentar "ponytail" menandakan bahwa exact sorting bukan scalable untuk benchmark berdurasi jam; upgrade ke streaming histogram (misal HdrHistogram) bila diperlukan.

---

## Snippet 3 — Load Runner dengan Per-VU Buffer untuk Menghindari Lock Contention

Source File: `internal/loadtest/runner.go:44-114`

Purpose: Menjalankan VUs konkuren, masing-masing menyimpan latensi lokal agar tidak butuh mutex; aggregasi sekali di akhir.

```go
func (r *Runner) Run(ctx context.Context) Result {
	var wg sync.WaitGroup
	start := time.Now()

	type vuResult struct {
		latencies []time.Duration
		errors    int
	}
	results := make([]vuResult, r.cfg.VUs)

	ctx, cancel := context.WithTimeout(ctx, r.cfg.Duration)
	defer cancel()

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
					req, err := http.NewRequestWithContext(ctx, r.cfg.Method, r.cfg.URL, bytes.NewReader(r.cfg.Body))
					if err != nil {
						if ctx.Err() == nil {
							errs++
						}
						continue
					}
					if r.cfg.ContentType != "" {
						req.Header.Set("Content-Type", r.cfg.ContentType)
					}

					reqStart := time.Now()
					resp, err := r.client.Do(req)
					if err != nil {
						if ctx.Err() == nil {
							errs++
						}
						continue
					}
					_, _ = io.Copy(io.Discard, resp.Body)
					_ = resp.Body.Close()

					if resp.StatusCode >= 400 {
						errs++
					} else {
						lats = append(lats, time.Since(reqStart))
					}
				}
			}
		}(i)
	}

	wg.Wait()
	totalDuration := time.Since(start)

	var allLatencies []time.Duration
	var totalErrs int
	for _, res := range results {
		allLatencies = append(allLatencies, res.latencies...)
		totalErrs += res.errors
	}

	return CalculateMetrics(allLatencies, totalErrs, totalDuration)
}
```

Explanation: Setiap VU goroutine menyimpan `latencies` dan `errors` di struktur lokal. Hanya setelah `wg.Wait()` semua slice digabungkan. Pola ini menghindari `sync.Mutex` contention pada path kritis dan memastikan akurasi pengukuran latency tanpa overhead locking. Hanya latency untuk response HTTP 2xx direkam; error transportasi dan response HTTP >= 400 meningkatkan `errs` tanpa menambah entry ke `lats`, sehingga P50/P95/P99 hanya merefleksikan latency request sukses.

---

## Snippet 4 — Custom HTTP Transport untuk Menghindari Client-Side Bottleneck

Source File: `internal/loadtest/runner.go:26-41`

Purpose: Meningkatkan `MaxIdleConns` dan `MaxIdleConnsPerHost` agar load generator tidak menjadi bottleneck.

```go
func NewRunner(cfg Config) *Runner {
	if cfg.VUs <= 0 {
		cfg.VUs = 1
	}
	tr := &http.Transport{
		MaxIdleConns:        1000,
		MaxIdleConnsPerHost: 1000,
	}
	return &Runner{
		cfg: cfg,
		client: &http.Client{
			Transport: tr,
			Timeout:   5 * time.Second,
		},
	}
}
```

Explanation: Default HTTP client Go memiliki limit koneksi idle yang rendah, yang bisa membatasi throughput load generator secara tidak adil. Custom transport dengan pool besar memastikan sisi klien tidak melimit beban.

---

## Snippet 5 — Demo Runner: Smoke vs Stress Contrast

Source File: `cmd/demo/main.go:16-58`

Purpose: Menjalankan dua skenario bersebelahan: smoke test (VUs < kapasitas) dan stress test (VUs >> kapasitas).

```go
func main() {
	cfg := server.Config{
		MaxDBConnections: 5,
		DBQueryDuration:  20 * time.Millisecond,
	}
	srv := server.New(cfg)
	ts := httptest.NewServer(srv.Routes())
	defer ts.Close()

	fmt.Println("Starting Load Test Demo (Booking Bengkel)")
	fmt.Printf("Server simulated DB connections: %d\n", cfg.MaxDBConnections)
	fmt.Printf("Simulated DB query duration: %s\n\n", cfg.DBQueryDuration)

	testDuration := 2 * time.Second

	smokeCfg := loadtest.Config{
		URL:         ts.URL + "/booking",
		Method:      http.MethodPost,
		Body:        []byte(`{"vehicle_id":"V123","service":"oil_change"}`),
		ContentType: "application/json",
		VUs:         2, // Below DB max connections
		Duration:    testDuration,
	}
	fmt.Println("--- Running Smoke Test (2 VUs) ---")
	smokeRunner := loadtest.NewRunner(smokeCfg)
	smokeRes := smokeRunner.Run(context.Background())
	printResults(smokeRes)

	stressCfg := loadtest.Config{
		URL:         ts.URL + "/booking",
		Method:      http.MethodPost,
		Body:        []byte(`{"vehicle_id":"V123","service":"oil_change"}`),
		ContentType: "application/json",
		VUs:         50, // 10x DB max connections, will cause queuing
		Duration:    testDuration,
	}
	fmt.Println("\n--- Running Stress Test (50 VUs) ---")
	stressRunner := loadtest.NewRunner(stressCfg)
	stressRes := stressRunner.Run(context.Background())
	printResults(stressRes)
}
```

Explanation: Skenario demo mensimulasikan aplikasi Booking Bengkel dengan kapasitas koneksi DB 5. Smoke test (2 VUs) berjalan tanpa antrian; stress test (50 VUs) memaksa antrian panjang, memperlihatkan perbedaan P95 yang drastis.

---

## Snippet 6 — Test Memverifikasi Stress P95 Lebih Tinggi dari Smoke P95

Source File: `tests/loadtest_test.go:16-60`

Purpose: Memastikan implementasi benar-benar menampilkan degradasi tail latency saat resource terebut.

```go
func TestLoadTest_SmokeVsStress(t *testing.T) {
	srv := server.New(server.Config{
		MaxDBConnections: 2,
		DBQueryDuration:  10 * time.Millisecond,
	})
	ts := httptest.NewServer(srv.Routes())
	defer ts.Close()

	// Smoke test: 1 VU, zero queueing
	smokeCfg := loadtest.Config{
		URL:      ts.URL + "/booking",
		Method:   http.MethodPost,
		Body:     []byte(`{}`),
		VUs:      1,
		Duration: 500 * time.Millisecond,
	}
	smokeRes := loadtest.NewRunner(smokeCfg).Run(context.Background())

	if smokeRes.TotalRequests == 0 {
		t.Fatal("smoke test recorded zero requests")
	}
	if smokeRes.ErrorCount != 0 {
		t.Fatalf("expected 0 errors in smoke test, got %d", smokeRes.ErrorCount)
	}

	// Stress test: 10 VUs against 2 DB slots -> heavy queuing
	stressCfg := loadtest.Config{
		URL:      ts.URL + "/booking",
		Method:   http.MethodPost,
		Body:     []byte(`{}`),
		VUs:      10,
		Duration: 500 * time.Millisecond,
	}
	stressRes := loadtest.NewRunner(stressCfg).Run(context.Background())

	if stressRes.TotalRequests == 0 {
		t.Fatal("stress test recorded zero requests")
	}
	if stressRes.P95Latency <= smokeRes.P95Latency {
		t.Fatalf("expected stress P95 (%s) to exceed smoke P95 (%s)", stressRes.P95Latency, smokeRes.P95Latency)
	}
	if stressRes.P95Latency <= stressRes.AvgLatency {
		t.Fatalf("expected stress P95 (%s) to exceed stress Avg (%s) due to tail queuing", stressRes.P95Latency, stressRes.AvgLatency)
	}
}
```

Explanation: Test ini adalah bukti terstruktur bahwa implementasi benar-benar menunjukkan perilaku queueing. Dua assert kunci: (1) P95 stress > P95 smoke, (2) P95 stress > Avg stress — membuktikan rata-rata menutupi tail.

---

## Snippet 7 — Test Race-Free Meterika (Invarian Statistik)

Source File: `internal/loadtest/metrics_test.go:83-108`

Purpose: Memverifikasi relasi orde persentil selalu terpenuhi.

```go
func TestCalculateMetrics_Invariants(t *testing.T) {
	latencies := []time.Duration{
		45 * time.Millisecond, 12 * time.Millisecond, 100 * time.Millisecond,
		5 * time.Millisecond, 20 * time.Millisecond, 500 * time.Millisecond,
		15 * time.Millisecond, 22 * time.Millisecond, 18 * time.Millisecond,
		80 * time.Millisecond, 250 * time.Millisecond, 30 * time.Millisecond,
	}

	res := CalculateMetrics(latencies, 0, time.Second)

	if res.MinLatency > res.P50Latency {
		t.Errorf("invariant violated: Min (%s) > P50 (%s)", res.MinLatency, res.P50Latency)
	}
	if res.P50Latency > res.P90Latency {
		t.Errorf("invariant violated: P50 (%s) > P90 (%s)", res.P50Latency, res.P90Latency)
	}
	if res.P90Latency > res.P95Latency {
		t.Errorf("invariant violated: P90 (%s) > P95 (%s)", res.P90Latency, res.P95Latency)
	}
	if res.P95Latency > res.P99Latency {
		t.Errorf("invariant violated: P95 (%s) > P99 (%s)", res.P95Latency, res.P99Latency)
	}
	if res.P99Latency > res.MaxLatency {
		t.Errorf("invariant violated: P99 (%s) > Max (%s)", res.P99Latency, res.MaxLatency)
	}
}
```

Explanation: Test ini memastikan algoritma percentile tidak pernah melanggar urutan statistik dasar. Berguna sebagai regression test saat mengganti implementasi percentile (misal: migrasi ke HdrHistogram).