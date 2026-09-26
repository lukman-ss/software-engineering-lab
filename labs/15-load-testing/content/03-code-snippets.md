## Snippet 1 — Server Connection Pool Semaphore

Source File: internal/server/server.go
Purpose: Menyimulasikan kapasitas terbatas database connection pool menggunakan semafor kanal Go untuk memaksa antrean permintaan saat melampaui ambang batas.

```go
func (s *Server) handleBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	atomic.AddInt64(&s.activeReq, 1)
	defer atomic.AddInt64(&s.activeReq, -1)

	// Acquire DB connection slot (simulates DB connection pool limit)
	select {
	case s.semaphore <- struct{}{}:
	case <-r.Context().Done():
		return
	}
	defer func() { <-s.semaphore }()

	dur := s.cfg.DBQueryDuration
	if atomic.LoadInt64(&s.activeReq) > int64(s.cfg.MaxDBConnections) {
		if rand.Float32() < 0.10 {
			dur = s.cfg.DBQueryDuration * 25
		}
	}
	t := time.NewTimer(dur)
	defer t.Stop()

	select {
	case <-t.C:
	case <-r.Context().Done():
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(BookingResponse{
		Status:    "confirmed",
		BookingID: "BK-1001",
	})
}
```

Explanation: Semafor memblokir pemrosesan masuk jika batas `MaxDBConnections` telah terisi, meniru antrean sumber daya. Context cancellation ditangani dengan `select` untuk menghindari goroutine terjebak. Penundaan query disimulasikan menggunakan `time.Timer` (bukan `time.Sleep` langsung) untuk mendukung context cancellation. Saat request terakumulasi di atas kapasitas pool, server menambahkan 10% probabililitas menunda query 25x lebih lama (simulasi latency variasi real-world). Slowdown 10% ini adalah *penambah (amplifier)* di atas mekanisme antrean utama; bahkan tanpanya, stress test tetap menunjukkan degradasi P95 murni dari penumpukan antrean semafor.

## Snippet 2 — Lock-Free Concurrent Load Runner

Source File: internal/loadtest/runner.go
Purpose: Menggunakan goroutine independen yang mencatat hasil ke lokasinya sendiri dalam slice pre-alokasi untuk mencegah kontensi resource atau sinkronisasi pelambatan internal pada level generator.

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
				req, err := http.NewRequestWithContext(ctx, r.cfg.Method, r.cfg.URL, bytes.NewReader(r.cfg.Body))
				if err != nil {
					errs++
					continue
				}
				if r.cfg.ContentType != "" {
					req.Header.Set("Content-Type", r.cfg.ContentType)
				}

				reqStart := time.Now()
				resp, err := r.client.Do(req)
				if err != nil {
					// Only count as error if not a context cancellation
					if ctx.Err() == nil {
						errs++
					}
					continue
				}
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
```

Explanation: Hasil dari tiap Virtual User (VU) tidak dikunci (lock-free) dengan menempatkannya di array `results` berbasis index unik `vuID`. Data ini baru diagregasikan usai fungsi generator tuntas, memastikan alat ukur tidak menjadi *bottleneck* akibat `sync.Mutex`. **Catatan**: Latency hanya direkam untuk request berhasil (HTTP 2xx); request dengan error transport atau HTTP >= 400 tidak menghasilkan entri di slice `lats`, tetapi dihitung sebagai `errs`. Ini berarti metrik P50/P95/P99 mencerminkan latensi permintaan berhasil saja.

## Snippet 3 — Percentile Calculation

Source File: internal/loadtest/metrics.go
Purpose: Menghitung distribusi statistik persentil menggunakan sorting matematis linear.

```go
func CalculateMetrics(latencies []time.Duration, errors int, totalDuration time.Duration) Result {
// ...
	// ponytail: slice sorting fine for <1M samples; upgrade to streaming histogram if memory constrained
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
```

Explanation: Data durasi respon direplikasi dan disortir (*O(N log N)*). Pengambilan persentil dilakukan secara deterministik dari index kalkulasi proporsional. Metode ini disengaja untuk skenario lab ringan dan dapat ditukar dengan histogram streaming (HdrHistogram) jika perlu menangani beban produksi yang berat memori.

**Catatan implementasi**: Fungsi `CalculateMetrics` mengembalikan `Result` dengan tujuh metrik persentil (P50/P90/P95/P99 dan min/avg/max). Namun, `cmd/demo/main.go` memilih mencetak hanya subset P50/P95/P99 pada `printResults`; P90 tetap dihitung dan tersedia di `Result.P90Latency` jika dibutuhkan.
