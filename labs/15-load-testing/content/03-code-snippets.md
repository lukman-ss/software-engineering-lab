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

Explanation: Semafor memblokir pemrosesan masuk jika batas `MaxDBConnections` telah terisi, meniru antrean sumber daya. Waktu tidur tambahan (`time.Sleep`) menyimulasikan durasi eksekusi kueri yang menghalangi slot koneksi dari pembebasan instan.

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

Explanation: Hasil dari tiap Virtual User (VU) tidak dikunci (lock-free) dengan menempatkannya di array `results` berbasis index unik `vuID`. Data ini baru diagregasikan usai fungsi generator tuntas, memastikan alat ukur tidak menjadi *bottleneck* akibat `sync.Mutex`.

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

func percentile(sorted []time.Duration, pct float64) time.Duration {
	if len(sorted) == 0 {
		return 0
	}
	idx := int(float64(len(sorted)-1) * (pct / 100.0))
	return sorted[idx]
}
```

Explanation: Data durasi respon direplikasi dan disortir (*O(N log N)*). Pengambilan persentil dilakukan secara deterministik dari index kalkulasi proporsional. Metode ini disengaja untuk skenario lab ringan dan dapat ditukar dengan histogram streaming (HdrHistogram) jika perlu menangani beban produksi yang berat memori.
