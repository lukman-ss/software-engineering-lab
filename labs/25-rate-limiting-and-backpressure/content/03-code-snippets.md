# Code Snippets

## Snippet 1 — Token Bucket Core Logic

Source File: `internal/ratelimit/bucket.go`
Purpose: Implementasi thread-safe token bucket dengan refill terus-menerus dan metode `RetryAfterSeconds`.

```go
type TokenBucket struct {
	mu         sync.Mutex
	capacity   float64
	tokens     float64
	refillRate float64
	lastRefill time.Time
}

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

Explanation: Token bucket menggunakan mutex untuk keamanan konkurensi. Token diisi secara berkala berdasarkan waktu elapsed. `RetryAfterSeconds` menghitung waktu tunggu hingga token cukup tersedia — digunakan middleware HTTP untuk header `Retry-After`.

---

## Snippet 2 — Leaky Bucket Core Logic

Source File: `internal/ratelimit/bucket.go`
Purpose: Implementasi leaky bucket untuk traffic smoothing dengan drain rate konstan.

```go
type LeakyBucket struct {
	mu       sync.Mutex
	capacity float64
	water    float64
	leakRate float64
	lastLeak time.Time
}

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

Explanation: Leaky bucket mengalirkan beban dengan laju konstan. `water` mewakili "air" dalam bucket yang terus berkurang seiring waktu. Permintaan ditolak jika menambahkan 1 unit akan melebihi kapasitas.

---

## Snippet 3 — Multi-tenant Registry

Source File: `internal/ratelimit/registry.go`
Purpose: Registry per-tenant token bucket untuk isolasi (menggunakan API key/tenant ID, bukan IP).

```go
type Registry struct {
	mu          sync.RWMutex
	buckets     map[string]*TokenBucket
	capacity    float64
	refillRate  float64
}

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

Explanation: Double-checked locking pattern — baca dengan `RLock`, tulis dengan `Lock` hanya saat bucket belum ada. Ini memastikan tenant A dan B punya kuota mandiri.

---

## Snippet 4 — Bounded Queue dengan Fast Rejection

Source File: `internal/backpressure/queue.go`
Purpose: Worker pool dengan channel bounded, menolak cepat saat penuh (load shedding).

```go
var ErrQueueFull = errors.New("backpressure: queue capacity exceeded")

type BoundedQueue struct {
	capacity int
	queue    chan Job
	workers  int
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
}

func (bq *BoundedQueue) TrySubmit(job Job) error {
	select {
	case bq.queue <- job:
		bq.accepted.Add(1)
		return nil
	default:
		bq.rejected.Add(1)
		return ErrQueueFull
	}
}
```

Explanation: Pola `select-default` pada channel Go memastikan `TrySubmit` tidak blocking. Jika channel penuh, segera mengembalikan `ErrQueueFull`. Worker loop memproses job dari channel dan menghitung `processed`.

---

## Snippet 5 — AWS Jitter Backoff Strategies

Source File: `internal/retry/backoff.go`
Purpose: Implementasi 4 varian exponential backoff dengan jitter mengikuti spesifikasi AWS (Marc Brooker).

```go
type BackoffStrategy string

const (
	NoJitter           BackoffStrategy = "NoJitter"
	FullJitter         BackoffStrategy = "FullJitter"
	EqualJitter        BackoffStrategy = "EqualJitter"
	DecorrelatedJitter BackoffStrategy = "DecorrelatedJitter"
)

func ComputeBackoff(strategy BackoffStrategy, attempt int, cfg Config, prevSleep time.Duration) time.Duration {
	baseFloat := float64(cfg.Base)
	capFloat := float64(cfg.Cap)

	expBackoff := baseFloat * math.Pow(2, float64(attempt))
	temp := math.Min(capFloat, expBackoff)

	switch strategy {
	case NoJitter:
		return time.Duration(temp)

	case FullJitter:
		if temp <= 0 {
			return 0
		}
		sleep := rand.Float64() * temp
		return time.Duration(sleep)

	case EqualJitter:
		half := temp / 2.0
		sleep := half + rand.Float64()*half
		return time.Duration(sleep)

	case DecorrelatedJitter:
		prevFloat := float64(prevSleep)
		if prevFloat < baseFloat {
			prevFloat = baseFloat
		}
		rangeMax := prevFloat * 3.0
		sleep := baseFloat + rand.Float64()*(rangeMax-baseFloat)
		return time.Duration(math.Min(capFloat, sleep))

	default:
		return time.Duration(temp)
	}
}
```

Explanation: Setiap varian mengikuti formula AWS:
- **FullJitter**: `random(0, min(cap, base * 2^attempt))` — distribusi paling merata
- **EqualJitter**: setengah tetap + setengah acak
- **DecorrelatedJitter**: menggunakan `prevSleep * 3` sebagai upper bound baru

---

## Snippet 6 — HTTP 429 Middleware (RFC 6585)

Source File: `internal/httputil/middleware.go`
Purpose: Middleware HTTP yang menegakkan rate limit tenant dan mengembalikan 429 standar.

```go
type RateLimitResponse struct {
	Error      string `json:"error"`
	RetryAfter int    `json:"retry_after"`
}

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
			w.WriteHeader(http.StatusTooManyRequests)

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

Explanation: Ekstrak tenant dari header `X-API-Key` (fallback "anonymous"). Jika bucket menolak, tulis status 429, header `Retry-After`, dan body JSON standar. Klien dapat membaca `Retry-After` untuk menentukan kapan retry.

---

## Snippet 7 — Demo Output (Verified Behavior)

Source File: `cmd/demo/main.go`
Purpose: Menampilkan perilaku nyata dari keempat komponen dalam satu eksekusi.

```go
func main() {
	fmt.Println("=== 1. Token Bucket Burst & Rate Limiting ===")
	tb := ratelimit.NewTokenBucket(3, 5)
	for i := 1; i <= 5; i++ {
		allowed := tb.Allow()
		fmt.Printf("Request #%d: Allowed=%v (Remaining Tokens: %.1f)\n", i, allowed, tb.Tokens())
	}
	time.Sleep(300 * time.Millisecond)
	fmt.Printf("After 300ms pause: Allowed=%v (Remaining Tokens: %.1f)\n", tb.Allow(), tb.Tokens())

	fmt.Println("\n=== 2. Leaky Bucket Traffic Smoothing ===")
	lb := ratelimit.NewLeakyBucket(3, 10)
	for i := 1; i <= 5; i++ {
		allowed := lb.Allow()
		fmt.Printf("Request #%d: Allowed=%v (Current Water Level: %.1f)\n", i, allowed, lb.Water())
	}

	fmt.Println("\n=== 3. Bounded Queue Backpressure (Load Shedding) ===")
	bq := backpressure.NewBoundedQueue(3, 1)
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
	// ... stats output

	fmt.Println("\n=== 4. AWS Retry Backoff Strategies (Attempts 0..3) ===")
	cfg := retry.Config{Base: 100 * time.Millisecond, Cap: 1000 * time.Millisecond}
	for attempt := 0; attempt < 4; attempt++ {
		noJitter := retry.ComputeBackoff(retry.NoJitter, attempt, cfg, 0)
		fullJitter := retry.ComputeBackoff(retry.FullJitter, attempt, cfg, 0)
		equalJitter := retry.ComputeBackoff(retry.EqualJitter, attempt, cfg, 0)
		fmt.Printf("Attempt %d -> NoJitter: %-6v | FullJitter: %-6v | EqualJitter: %-6v\n",
			attempt, noJitter.Round(time.Millisecond), fullJitter.Round(time.Millisecond), equalJitter.Round(time.Millisecond))
	}
}
```

Explanation: Demo menunjukkan: (1) token bucket burst 3 lalu tolak, refill setelah jeda; (2) leaky bucket tolak burst; (3) bounded queue tolak job ke-4+; (4) perbandingan jitter variants.

---

## Snippet 8 — Test Verifikasi Bounds Jitter

Source File: `internal/retry/backoff_test.go`
Purpose: Test memastikan semua varian jitter menghasilkan sleep dalam batas yang ditetapkan.

```go
func TestComputeBackoff_Bounds(t *testing.T) {
	cfg := Config{
		Base: 100 * time.Millisecond,
		Cap:  2 * time.Second,
	}

	for attempt := 0; attempt < 10; attempt++ {
		sleepFull := ComputeBackoff(FullJitter, attempt, cfg, 0)
		if sleepFull < 0 || sleepFull > cfg.Cap {
			t.Fatalf("FullJitter sleep out of bounds: %v", sleepFull)
		}

		sleepEqual := ComputeBackoff(EqualJitter, attempt, cfg, 0)
		if sleepEqual < 0 || sleepEqual > cfg.Cap {
			t.Fatalf("EqualJitter sleep out of bounds: %v", sleepEqual)
		}

		sleepNo := ComputeBackoff(NoJitter, attempt, cfg, 0)
		if sleepNo < cfg.Base || sleepNo > cfg.Cap {
			t.Fatalf("NoJitter sleep out of bounds: %v", sleepNo)
		}
	}
}
```

Explanation: Test memverifikasi bahwa untuk 10 percobaan, FullJitter berada dalam `[0, cap]`, EqualJitter dalam `[0, cap]`, NoJitter dalam `[base, cap]`. Race detector tidak menemukan masalah.