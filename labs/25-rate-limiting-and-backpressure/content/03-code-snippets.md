# Code Snippets

Semua snippet diambil verbatim dari implementasi yang lulus Engineering Audit (APPROVED).

## Snippet 1 — Token Bucket: Allow dengan Replenish + Debit

Source File:
`internal/ratelimit/bucket.go`

Purpose:
Menunjukkan replenishment token berbasis elapsed time (monotonic), cap ke capacity, lalu debit satu token.

```go
func (tb *TokenBucket) AllowN(n float64) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	tb.tokens = tb.tokens + elapsed*tb.refillRate
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}
	tb.lastRefill = now

	if tb.tokens >= n {
		tb.tokens -= n
		return true
	}
	return false
}
```

Explanation:
Interval antar-panggilan dikonversi ke detik dan dikalikan `refillRate` untuk menambah token. Token dibatasi maksimum `capacity` sehingga idle panjang tidak menghasilkan burst tanpa batas. Jika token cukup, dikurangi dan request diizinkan; jika tidak, request ditolak tanpa blocking.

## Snippet 2 — Retry-After untuk HTTP 429

Source File:
`internal/ratelimit/bucket.go`

Purpose:
Menghitung detik tunggu minimum agar `n` token tersedia, dipakai middleware untuk header `Retry-After`.

```go
func (tb *TokenBucket) RetryAfterSeconds(n float64) int {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	tokens := tb.tokens + elapsed*tb.refillRate
	if tokens > tb.capacity {
		tokens = tb.capacity
	}

	if tokens >= n {
		return 0
	}
	needed := n - tokens
	secs := needed / tb.refillRate
	if secs <= 0 {
		return 1
	}
	rounded := int(secs)
	if float64(rounded) < secs {
		rounded++
	}
	return rounded
}
```

Explanation:
Kekurangan token dibagi `refillRate` lalu dibulatkan ke atas, sehingga `Retry-After` selalu membumikan ke waktu tunggu yang memadai. Per audit, method ini membutuhkan `refillRate > 0`; nilai nol menyebabkan panic (dokumentasi precondition, non-blocking gap).

## Snippet 3 — Leaky Bucket: Drain Konstan + Tolak Saat Penuh

Source File:
`internal/ratelimit/bucket.go`

Purpose:
Menunjukkan perilaku smoothing: air dikeluarkan konstan, request baru ditolak ketika level air mencapai kapasitas.

```go
func (lb *LeakyBucket) Allow() bool {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(lb.lastLeak).Seconds()
	lb.water = lb.water - elapsed*lb.leakRate
	if lb.water < 0 {
		lb.water = 0
	}
	lb.lastLeak = now

	if lb.water+1.0 <= lb.capacity {
		lb.water += 1.0
		return true
	}
	return false
}
```

Explanation:
Air berkurang sebesar `elapsed × leakRate` pada setiap panggilan. Request hanya diterima bila setelah penambahan air masih `<= capacity`, memaksa laju keluaran tetap konstan berbeda dari token bucket yang mentolerir burst.

## Snippet 4 — Registry Per-Tenant (Double-Checked Locking)

Source File:
`internal/ratelimit/registry.go`

Purpose:
Menyediakan satu token bucket per tenant key agar pembatasan tidak berbagi state antar-tenant (hindari tabrakan IP di balik CGNAT/RFC 6598).

```go
func (r *Registry) Get(tenantKey string) *TokenBucket {
	r.mu.RLock()
	tb, exists := r.buckets[tenantKey]
	r.mu.RUnlock()
	if exists {
		return tb
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if tb, exists = r.buckets[tenantKey]; exists {
		return tb
	}
	tb = NewTokenBucket(r.capacity, r.refillRate)
	r.buckets[tenantKey] = tb
	return tb
}
```

Explanation:
Pembacaan pertama memakai `RLock` agar jalur normal murah. Jika key belum ada, fallback ke `Lock` dengan pemeriksaan ulang sehingga hanya satu bucket dibuat per key walau ada akses konkuren.

## Snippet 5 — Bounded Queue: Fast Rejection (Backpressure)

Source File:
`internal/backpressure/queue.go`

Purpose:
Menerima job bila buffer tersedia, menolak seketika dengan `ErrQueueFull` bila penuh — tanpa memblokir pemanggil.

```go
func (bq *BoundedQueue) TrySubmit(job Job) error {
	bq.stopMu.RLock()
	defer bq.stopMu.RUnlock()

	if bq.stopped {
		return ErrQueueStopped
	}

	select {
	case <-bq.ctx.Done():
		return ErrQueueStopped
	case bq.queue <- job:
		bq.accepted.Add(1)
		return nil
	default:
		bq.rejected.Add(1)
		return ErrQueueFull
	}
}
```

Explanation:
Pola `select-default` memberi penolakan non-blocking (zero-allocation) saat kanal penuh, yang merupakan inti backpressure: sinyal overload dikembalikan ke pemanggil alih-alih antrean tumbuh tanpa batas.

## Snippet 6 — Worker Loop dengan Shutdown via Context

Source File:
`internal/backpressure/queue.go`

Purpose:
Konsumer yang memproses job sampai context dibatalkan, dengan penghitung `processed`.

```go
func (bq *BoundedQueue) workerLoop() {
	defer bq.wg.Done()
	for {
		select {
		case <-bq.ctx.Done():
			return
		case job, ok := <-bq.queue:
			if !ok {
				return
			}
			_ = job(bq.ctx)
			bq.processed.Add(1)
		}
	}
}
```

Explanation:
`select` dua-arah antara cancellation dan job masuk memungkinkan shutdown bersih. Per audit, return error dari `job` dibuang (`_ = job(...)`) dan tidak ada pemisahan panic — semantik worker dicatat sebagai known gap, bukan perilaku terverifikasi.

## Snippet 7 — Full Jitter (Formula AWS)

Source File:
`internal/retry/backoff.go`

Purpose:
Menghitung sleep `random(0, min(cap, base × 2^attempt))` — varian jitter yang distandardisasi AWS SDK.

```go
	// temp = min(cap, base * 2^attempt)
	expBackoff := baseFloat * math.Pow(2, float64(attempt))
	temp := math.Min(capFloat, expBackoff)

	case FullJitter:
		// sleep = random_between(0, min(cap, base * 2^attempt))
		if temp <= 0 {
			return 0
		}
		sleep := rand.Float64() * temp
		return time.Duration(sleep)
```

Explanation:
Bagian eksponensial dibatasi `cap` lebih dulu, lalu dikalikan `random(0,1)`. Hasilnya interval tersebar merata di `[0, temp]`, menghindari sinkronisasi retry. Angka default AWS (50ms/1000ms/20s) adalah pilihan spesifik AWS SDK, bukan aturan universal.

## Snippet 8 — Equal Jitter

Source File:
`internal/retry/backoff.go`

Purpose:
Varian kedua: separuh nilai eksponensial tetap, separuh lagi dirandom.

```go
	case EqualJitter:
		// sleep = min(cap, base * 2^attempt) / 2 + random_between(0, min(cap, base * 2^attempt) / 2)
		half := temp / 2.0
		sleep := half + rand.Float64()*half
		return time.Duration(sleep)
```

Explanation:
Menjamin delay minimum setengah dari eksponensial sekaligus memberi sebaran acak pada separuhnya, mengurangi tapi tidak menghilangkan korelasi antar klien dibanding Full Jitter.

## Snippet 9 — Decorrelated Jitter

Source File:
`internal/retry/backoff.go`

Purpose:
Varian berbasis sleep sebelumnya: `min(cap, random(base, prevSleep × 3))`.

```go
	case DecorrelatedJitter:
		// sleep = min(cap, random_between(base, prevSleep * 3))
		prevFloat := float64(prevSleep)
		if prevFloat < baseFloat {
			prevFloat = baseFloat
		}
		rangeMax := prevFloat * 3.0
		sleep := baseFloat + rand.Float64()*(rangeMax-baseFloat)
		return time.Duration(math.Min(capFloat, sleep))
```

Explanation:
Range acak tumbuh dari sleep sebelumnya (multiplier 3), bukan dari eksponensial tetap. Sifat statistik rantai jitter decorrelated tidak dibuktikan oleh test (gap LOW per audit).

## Snippet 10 — Middleware 429 RFC 6585 + Retry-After

Source File:
`internal/httputil/middleware.go`

Purpose:
Ekstraksi tenant via header, cek bucket, dan respons 429 standar dengan header `Retry-After` serta body JSON.

```go
func RateLimitMiddleware(registry *ratelimit.Registry, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenantKey := r.Header.Get("X-API-Key")
		if tenantKey == "" {
			tenantKey = "anonymous"
		}

		bucket := registry.Get(tenantKey)
		if !bucket.Allow() {
			retryAfter := bucket.RetryAfterSeconds(1.0)
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
			w.WriteHeader(http.StatusTooManyRequests) // 429

			_ = json.NewEncoder(w).Encode(RateLimitResponse{
				Error:      "rate_limit_exceeded",
				RetryAfter: retryAfter,
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
```

Explanation:
Mengikuti RFC 6585: status 429 dan header `Retry-After`. Nilai header dan field `retry_after` di body diisi dari sumber yang sama (`RetryAfterSeconds`), meski test belum memastikan keduanya identik (gap LOW per audit).

## Snippet 11 — Demo: Token Bucket Burst Lalu Tolak

Source File:
`cmd/demo/main.go`

Purpose:
Membuktikan perilaku burst hingga kapasitas lalu penolakan, serta pemulihan token setelah jeda.

```go
	tb := ratelimit.NewTokenBucket(3, 5) // Cap 3, refill 5/s
	for i := 1; i <= 5; i++ {
		allowed := tb.Allow()
		fmt.Printf("Request #%d: Allowed=%v (Remaining Tokens: %.1f)\n", i, allowed, tb.Tokens())
	}

	time.Sleep(300 * time.Millisecond)
	fmt.Printf("After 300ms pause: Allowed=%v (Remaining Tokens: %.1f)\n", tb.Allow(), tb.Tokens())
```

Explanation:
Dengan kapasitas 3 dan refill 5/s, tiga request pertama lolos, request ke-4 dan ke-5 ditolak, lalu setelah jeda 300ms refill menghasilkan token kembali (sample demo: `Allowed=true`, sisa token ~0.5).

## Snippet 12 — Demo: Backpressure Load Shedding

Source File:
`cmd/demo/main.go`

Purpose:
Menunjukkan penolakan cepat ketika buffer penuh dan statistik accepted/rejected/processed.

```go
	bq := backpressure.NewBoundedQueue(3, 1) // 3 queue cap, 1 worker
	defer bq.Stop()

	for i := 1; i <= 6; i++ {
		jobID := i
		err := bq.TrySubmit(func(ctx context.Context) error {
			time.Sleep(50 * time.Millisecond)
			return nil
		})
		if err != nil {
			fmt.Printf("Job #%d: REJECTED (Backpressure Shedding: %v)\n", jobID, err)
		} else {
			fmt.Printf("Job #%d: ACCEPTED into bounded buffer\n", jobID)
		}
	}
```

Explanation:
Buffer 3 + worker 1 yang memproses 50ms per job membuat sebagian job ditolak `backpressure: queue capacity exceeded` saat pemanggil menyalip kecepatan konsumer — demonstrasi langsung dari mekanisme backpressure.
