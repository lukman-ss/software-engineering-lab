# Timeouts and Deadlines: Mencegah Satu Dependency Lambat Menghancurkan Sistem

## Problem

Dalam sistem terdistribusi, **satu dependency yang lambat jauh lebih berbahaya daripada satu dependency yang gagal**.

Ketika sebuah service mengembalikan `ECONNREFUSED`, worker thread langsung melepaskan resource dan tersedia untuk request baru. Namun ketika dependency melambat — dari 300ms menjadi 60s — request in-flight terakumulasi di thread pool, connection pool, dan memory. Menggunakan Little's Law:

```
L = λW
```

dimana *L* = jumlah request concurrently in-flight, *λ* = request rate, *W* = waktu pemrosesan (latency). Jika *W* naik 200 kali lipat (300ms → 60s), *L* meledak. Worker pool 100 thread mampu menangani 100 QPS dengan latency 1s, namun hanya 1.6 QPS dengan latency 60s. Request untuk endpoint yang **tidak terkait** — `/products`, `/profile` — tidak mendapatkan thread dan gagal dengan timeout. Satu dependency lambat menjadi system-wide outage.

Timeout dan deadline adalah mekanisme proteksi resource yang menghindarkan skenario ini.

## Mental Model

Timeout bukan angka tebakan yang dibesar-besarkan "supaya aman." Timeout adalah **resource protection guardrail** — batas waktu yang harus dibuat berdasarkan distribusi latency aktual (P95/P99) dan budget waktu end-to-end.

Kunci mental model:

1. **Timeout ≠ failure, timeout = indeterminate.** Timeout berarti caller tidak menerima response tepat waktu. Tidak ada yang tahu apakah operasi benar-benar gagal atau telah berhasil di server.
2. **Deadline harus diwariskan, tidak dikunci di awal.** Jika client membatalkan request pada detik ke-2, downstream yang mengira masih punya 15 detik tetap mengerjakan request yang sudah tidak ada yang menunggu hasilnya.
3. **Retri tanpa jitter sama dengan retry storm.** Banyak client yang gagal pada waktu yang sama akan retry pada waktu yang sama pula, memperparah beban server.
4. **Database dan worker juga butuh timeout.** API timeout tidak melindungi database; query yang berjalan selamanya atau lock yang tidak dilepas tetap menghabiskan connection pool.

## Core Concept

### 1. Context Deadline Propagation

`context.WithTimeout` di Go membuat batas waktu yang diwariskan secara hierarkis. Child context memiliki deadline yang adalah minimum dari parent deadline dan budget lokal — tidak pernah lebih besar dari parent.

```go
func ExecuteWithBudget(ctx context.Context, budget time.Duration, fn WorkerFunc) error {
    childCtx, cancel := context.WithTimeout(ctx, budget)
    defer cancel()
    done := make(chan error, 1)
    go func() { done <- fn(childCtx) }()
    select {
    case <-childCtx.Done(): return childCtx.Err()
    case err := <-done:      return err
    }
}
```

`chan error` berkapasitas 1 memastikan goroutine tidak bocor ketika `childCtx.Done()` terjadi sebelum `fn` selesai — goroutine menulis ke channel yang sudah dibaca oleh `select`.

**Yang dibuktikan test:**
- `TestExecuteWithBudget_Timeout`: function dengan durasi 100ms dibatalkan oleh budget 20ms, mengembalikan `context.DeadlineExceeded`
- `TestExecuteWithBudget_ParentTimeoutInherited`: parent dengan deadline 20ms membatalkan execution meskipun budget lokal 500ms

### 2. Exponential Backoff dengan Full Jitter

Tanpa jitter, semua client yang gagal pada waktu yang sama akan retry pada waktu yang sama pula (thundering herd). Full Jitter mendistribusikan waktu retry secara uniform:

```
sleep = random(0, min(maxBackoff, baseBackoff × 2^(attempt-1)))
```

Implementasi di lab:

```go
func (r *Retrier) CalculateBackoff(attempt int) time.Duration {
    multiplier := 1 << uint(attempt-1)
    temp := float64(r.cfg.BaseBackoff) * float64(multiplier)
    if temp > float64(r.cfg.MaxBackoff) {
        temp = float64(r.cfg.MaxBackoff)
    }
    sleep := rand.Float64() * temp
    return time.Duration(sleep)
}
```

Bit-shift `1 << uint(attempt-1)` menghasilkan 2^(attempt-1). Setiap retry, sleep berada dalam rentang [0, temp], sehingga retry tersebar secara uniform dan tidak menumpuk di satu titik waktu.

**Yang dibuktikan test:**
- `TestRetrier_RetryUntilSuccess`: 3 attempt dengan transient error pada attempt 1-2, sukses di attempt 3
- `TestRetrier_ExceedMaxAttempts`: max 2 attempt, persistent error, mengembalikan `ErrMaxRetriesExceeded`
- `TestRetrier_ContextCanceled`: context timeout 20ms membatalkan loop retry
- `TestRetrier_JitterBoundsAndZeroConfig`: backoff tidak negatif dan tidak melebihi `MaxBackoff`; zero config menghasilkan default (3 attempts, 10ms base, 100ms max)

**Yang dibuktikan demo:**
3 attempt dijalankan, 2 percayaan transient error, attempt ke-3 sukses tanpa error.

### 3. Circuit Breaker — Tiga State

Circuit breaker mencegah hammering ke downstream yang sudah diketahui gagal.

| State | Behavior |
|-------|----------|
| `CLOSED` | Normal, menghitung gagal berturut-turut |
| `OPEN` | Semua request langsung ditolak, tidak ada panggilan ke downstream |
| `HALF_OPEN` | Setelah cooldown, izinkan percobaan; sukses → kembali CLOSED, gagal → langsung kembali OPEN |

Transisi dijaga oleh `sync.RWMutex` sehingga aman diakses dari banyak goroutine.

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

**Yang dibuktikan test:**
- `TestCircuitBreaker_StateTransitions`: CLOSED → 2 gagal → OPEN → `Allow()` tolak → tunggu cooldown → HALF_OPEN → 2 sukses → CLOSED
- `TestCircuitBreaker_HalfOpenFailureTripsOpen`: OPEN → cooldown → HALF_OPEN → 1 gagal → langsung kembali OPEN

### 4. Idempotency Key Deduplication

Ketika timeout terjadi pada operasi non-idempoten (misalnya `POST /payment`), timeout **bukan bukti kegagalan**. Server mungkin sudah memproses charge, namun ACK atau response hilang di jaringan. Retry tanpa proteksi menghasilkan double charge.

Solusinya: idempotency key. Client mengirim `Idempotency-Key` header. Server menyimpan response yang dihasilkan berdasarkan key tersebut; request retry dengan key yang sama mengembalikan response yang tersimpan tanpa eksekusi ulang.

```go
func (s *Store) Get(key string) (string, bool) {
    s.mu.Lock()
    defer s.mu.Unlock()
    rec, ok := s.records[key]
    if !ok { return "", false }
    if time.Since(rec.CreatedAt) > s.ttl {
        delete(s.records, key)
        return "", false
    }
    return rec.Response, true
}
```

Lazy eviction: expired key dihapus saat diakses (`Get`), bukan oleh goroutine periodik. `Get` mengambil write lock (`mu.Lock()`, bukan `RLock`) karena melakukan mutasi map.

**Yang dibuktikan test:**
- `TestStore_GetSet`: set, dapat, tunggu TTL lewat → expired, key dihapus dari map
- `TestStore_ConcurrentAccess`: 50 goroutine melakukan Get/Set secara konkuren — race detector PASS
- `TestStore_LazyEvictionOnGet`: key expired terhapus dari map setelah Get

**Yang dibuktikan demo:**
`req-tx-99231` diproses sebagai "Charged $100 successfully", retry mengembalikan "Charged $100 successfully (DEDUPLICATED)" tanpa eksekusi ulang.

### 5. Integrasi: Retry + Circuit Breaker + Idempotency

Integration test `TestIntegration_RetryWithCircuitBreaker` menggabungkan retrier (4 attempt, backoff 5-20ms) dan circuit breaker (failure threshold 2, cooldown 50ms). Karena backend selalu gagal, circuit breaker terbuka setelah 2 attempt pertama, sisa attempt langsung ditolak tanpa menyentuh "downstream."

`TestIntegration_IdempotentRetry` menjalankan 1 eksekusi nyata + retry loop (3 attempt) dengan idempotency store — variabel `actualExecutions` tetap bernilai 1, membuktikan deduplication.

## Architecture

```
┌──────────────────────────────────────────────┐
│                    Caller                     │
│                                              │
│  context.WithTimeout(ctx, budget)            │
│         │                                    │
│         ▼                                    │
│  ┌─────────────┐    ┌──────────────────┐    │
│  │  Retrier    │───▶│  Circuit Breaker │    │
│  │  (jittered  │    │  (3 states)      │    │
│  │   backoff)  │    └──────────────────┘    │
│  └──────┬──────┘                             │
│         │                                    │
│         ▼                                    │
│  ┌──────────────────┐                        │
│  │  Idempotency     │                        │
│  │  Store (TTL,     │                        │
│  │  mutex-guarded)  │                        │
│  └──────────────────┘                        │
└──────────────────────────────────────────────┘
```

## Common Mistakes

**1. Timeout angka besar "supaya ama**

`context.WithTimeout(ctx, 30*time.Second)` untuk semua request. Dengan 100 QPS dan latency normal 200ms, hanya 20 request in-flight. Saat latency naik ke 30s, jumlah in-flight naik menjadi 3000 — melebihi kapasitas thread pool, mengubah slow dependency menjadi system-wide outage.

**2. Retry tanpa jitter dan budget.**
Client yang gagal menunggu waktu yang sama → retry bersamaan → load berlipat ganda → server makin down → retry makin banyak. Siklus yang memperparah outage alih-alih memperbaikinya.

**3. Menganggap timeout = gagal.**
Timeout hanya berarti response tidak diterima tepat waktu. Charge mungkin sudah terjadi. Retry tanpa idempotency key = double charge.

**4. Hanya timeout di layer HTTP.**
Query PostgreSQL tanpa `statement_timeout` tetap menahan connection meskipun HTTP timeout sudah terjadi di layer atas.

## Production Considerations

| Concern | Lab | Production |
|---------|-----|------------|
| Idempotency store | In-memory, hilang saat crash | DB transaksional atau Redis dengan TTL |
| Deadline propagation | Proses tunggal, context Go | gRPC `grpc-timeout` header, HTTP via `Request-Timeout` atau OpenTelemetry baggage |
| Adaptive timeout | Static budget | Moving-window P99 + MAD (open question) |
| Database guards | Tidak diimplementasikan | `statement_timeout`, `lock_timeout`, `idle_in_transaction_session_timeout` |
| Queue worker | Tidak diimplementasikan | Execution timeout, max retries, dead-letter queue |

**Catatan dari research audit (LOW severity):**
- Header deadline propagation HTTP/REST belum ada standar tunggal — heterogen antar framework
- Dynamic adaptive timeout berbasis real-time observability masih memerlukan infrastruktur observability tambahan

## Checklist

- [x] Research approved (0 unsupported claims, 3 LOW gaps)
- [x] Engineering approved (0 failures, race detector PASS, demo PASS)
- [x] Semua klaim faktual berasal dari sumber yang teridentifikasi dalam research
- [x] Code snippet sesuai dengan implementation aktual
- [x] Test claims sesuai dengan test yang sebenarnya
- [x] Caveats dipertahankan
- [x] Tidak ada benchmark palsu
- [x] Tidak ada incident yang dibuat-buat

## Key Takeaways

1. Slow dependencies lebih berbahaya daripada failed dependencies — Little's Law menjelaskan mengapa.
2. Timeout harus dibuat dari distribusi latency P95/P99, bukan angka "supaya ama."
3. Full Jitter (`rand(0, base × 2^attempt)`) menurunkan server contention lebih dari 50% dibanding unjittered.
4. Timeout adalah *state indeterminate*, bukan kegagalan — idempotency key melindungi operasi mutasi.
5. Context deadline diwariskan ke child; parent timeout ≤ child timeout.
6. Circuit breaker memotong chain retry segera setelah downstream diketahui gagal.
7. `statement_timeout` dan `lock_timeout` PostgreSQL melindungi connection pool dari query runaway — HTTP timeout saja tidak cukup.

## Sources

| Sumber | Kategori |
|--------|----------|
| Google SRE Book — Ch.22 Addressing Cascading Failures | Cascading failure, Little's Law, deadline propagation |
| AWS Architecture Blog — Exponential Backoff And Jitter | Backoff algorithms, Full Jitter math |
| gRPC Documentation — Deadlines | Deadline propagation, `grpc-timeout`, clock skew |
| PostgreSQL 18 Documentation — Ch.19.11 Client Connection Defaults | `statement_timeout`, `lock_timeout`, `idle_in_transaction_session_timeout` |
| Stripe API Documentation — Error Handling & Idempotency | Timeout ambiguity, idempotency key |
