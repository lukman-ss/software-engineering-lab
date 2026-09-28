## Snippet 1 — Context Deadline Budget Execution

Source File: `internal/deadline/deadline.go:13-27`

Purpose: Membungkus function `fn` dalam child context dengan budget timeout lokal, namun tetap mengikuti deadline parent yang lebih ketat.

```go
func ExecuteWithBudget(ctx context.Context, budget time.Duration, fn WorkerFunc) error {
	childCtx, cancel := context.WithTimeout(ctx, budget)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- fn(childCtx)
	}()

	select {
	case <-childCtx.Done():
		return childCtx.Err()
	case err := <-done:
		return err
	}
}
```

Explanation: `childCtx` memiliki deadline `budget` lebih kecil dari parent. Pemilihan `select` memastikan eksekusi dibatalkan ketika salah satu dari `childCtx.Done()` terjadi atau `fn` selesai. Buffered channel (`chan error, 1`) mencegah goroutine leak jika `childCtx.Done()` terjadi sebelum `done <- fn(childCtx)` dieksekusi.

---

## Snippet 2 — Full Jitter Backoff Calculation

Source File: `internal/retry/retry.go:35-48`

Purpose: Menghitung durasi sleep dengan exponential backoff dan Full Jitter — random uniform dalam rentang [0, max].

```go
func (r *Retrier) CalculateBackoff(attempt int) time.Duration {
	if attempt <= 0 {
		return 0
	}
	multiplier := 1 << uint(attempt-1)
	temp := float64(r.cfg.BaseBackoff) * float64(multiplier)
	maxVal := float64(r.cfg.MaxBackoff)
	if temp > maxVal {
		temp = maxVal
	}
	sleep := rand.Float64() * temp
	return time.Duration(sleep)
}
```

Explanation: Bit-shift `1 << uint(attempt-1)` menghasilkan 2^(attempt-1). Setiap retry, rentang sleep meningkat secara eksponensial. `rand.Float64()` menghasilkan nilai uniform [0,1), sehingga `sleep` berada dalam [0, temp]. Batas maksimum `MaxBackoff` mencegah jitter menjadi terlalu lama.

---

## Snippet 3 — Circuit Breaker State Machine

Source File: `internal/circuit/circuit.go:108-126`

Purpose: Menghitung kegagalan dan mengatur transisi state CLOSED → OPEN → HALF_OPEN.

```go
func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.state == StateHalfOpen {
		b.state = StateOpen
		b.failures = 0
		b.successes = 0
		b.lastStateChg = time.Now()
	} else if b.state == StateClosed {
		b.failures++
		if b.failures >= b.cfg.FailureThreshold {
			b.state = StateOpen
			b.failures = 0
			b.successes = 0
			b.lastStateChg = time.Now()
		}
	}
}
```

Explanation: State HALF_OPEN menerima request uji. `SuccessThreshold` sukses berturut-turut diperlukan untuk kembali CLOSED. Di HALF_OPEN, satu kegagalan langsung mengembalikan OPEN. State CLOSED menghitung gagal berturut-turut; setelah mencapai `FailureThreshold`, circuit OPEN. Semua state checks dan write mutex (`mu.Lock()`) memastikan thread-safety di banyak goroutine.

---

## Snippet 4 — Idempotency Store dengan Lazy TTL Eviction

Source File: `internal/idempotency/idempotency.go:29-42`

Purpose: Mengambil response yang tersimpan berdasarkan key, atau set baru jika belum ada atau expired.

```go
func (s *Store) Get(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rec, ok := s.records[key]
	if !ok {
		return "", false
	}
	if time.Since(rec.CreatedAt) > s.ttl {
		delete(s.records, key)
		return "", false
	}
	return rec.Response, true
}
```

Explanation: `Get` menggunakan write lock (`mu.Lock()`) karena melakukan mutasi map (`delete`). Lazy eviction: expired key dihapus saat diakses, bukan oleh goroutine periodik. `CreatedAt` digunakan untuk memeriksa TTL — tidak mengandalkan sistem time tick.

---

## Snippet 5 — Demo Execution: Deadline Propagation

Source File: `cmd/demo/main.go:18-31`

Purpose: Demonstrasi parent deadline lebih kecil dari budget child.

```go
ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
defer cancel()

err := deadline.ExecuteWithBudget(ctx, 100*time.Millisecond, func(childCtx context.Context) error {
	select {
	case <-time.After(80 * time.Millisecond):
		return nil
	case <-childCtx.Done():
		return childCtx.Err()
	}
})
fmt.Printf("Deadline propagation result: %v\n", err)
```

Explanation: Parent deadline 50ms lebih kecil dari child budget 100ms. `time.After(80ms)` melebihi parent deadline, sehingga `childCtx.Done()` terpicu dan mengembalikan `context.DeadlineExceeded`.

---

## Snippet 6 — Demo Execution: Circuit Breaker State Transitions

Source File: `cmd/demo/main.go:51-74`

Purpose: Demonstrasi CLOSED → OPEN → HALF_OPEN → CLOSED state transitions.

```go
cb := circuit.NewBreaker(circuit.Config{
	FailureThreshold: 2,
	SuccessThreshold: 1,
	Cooldown:         50 * time.Millisecond,
})

cb.Execute(func() error { return dummyErr })
cb.Execute(func() error { return dummyErr })
fmt.Printf("State after 2 failures: %s\n", cb.State())

err = cb.Execute(func() error { return nil })
fmt.Printf("Execution attempt while OPEN: %v\n", err)

time.Sleep(60 * time.Millisecond)
fmt.Printf("State after cooldown: %s\n", cb.State())

err = cb.Execute(func() error { return nil })
fmt.Printf("Execution attempt in HALF_OPEN (success): %v -> New State: %s\n", err, cb.State())
```

Explanation: Dua kegagalan berturut-turut dalam state CLOSED membuat circuit OPEN. Request berikutnya langsung ditolak dengan `ErrCircuitOpen`. Setelah cooldown 60ms (> 50ms), state berubah menjadi HALF_OPEN. Satu keberhasilan dalam HALF_OPEN cukup untuk kembali ke CLOSED karena `SuccessThreshold=1`.

---

## Snippet 7 — Demo Execution: Idempotency Deduplication

Source File: `cmd/demo/main.go:75-95`

Purpose: Demonstrasi request retry dengan idempotency key hanya mengeksekusi sekali.

```go
store := idempotency.NewStore(5 * time.Minute)
idKey := "req-tx-99231"

processPayment := func(key string, amount int) (string, error) {
	if res, ok := store.Get(key); ok {
		return fmt.Sprintf("%s (DEDUPLICATED)", res), nil
	}
	result := fmt.Sprintf("Charged $%d successfully", amount)
	store.Set(key, result)
	return result, nil
}

res1, _ := processPayment(idKey, 100)
fmt.Printf("First request execution: %s\n", res1)

res2, _ := processPayment(idKey, 100)
fmt.Printf("Retried request execution: %s\n", res2)
```

Explanation: Request kedua dengan key yang sama langsung menemukan response yang tersimpan. Tidak ada "Charged $100 successfully" kedua — hanya satu eksekusi aktual dan response kedua ditandai "(DEDUPLICATED)."
