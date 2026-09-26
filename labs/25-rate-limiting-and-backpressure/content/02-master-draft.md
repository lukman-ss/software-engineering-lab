# Rate Limiting & Backpressure: Mencegah Runnurnya Sistem di Bawah Beban Lebih

## Problem

Sistem distribusi sering runtuh ketika laju permintaan masuk melebihi kapasitas pemrosesan. Tanpa mekanisme perlindungan, beberapa hal yang dapat terjadi:

- **Antrian tidak terbatas tumbuh eksis**, memicu kehabisan memori
- **Worker overload** terus-menerima beban melebihi kemampuan
- **Retry storm**: banyak klien mengula kembali setelah jeda serupa, menciptakan puncak beban saling menambah

## Mengapa Ini Penting

Rate limiting dan backpressure adalah konsep inti untuk menciptakan sistem yang **stabil, prediksi, dan dapat mengelola kembali beban**. Bukan sekadar "membatasi kecepatan" — melainkan **mengontrol aliran sehingga sistem tetap hidup**.

## Mental Model

Pertanyaan kunci untuk kesehatan sistem:

> **Apakah pekerjaan masuk lebih cepat daripada kemampuan sistem menyelesaikannya?**

Jika ya → Anda sedang menghadapi **masalah kapasitas** atau **masalah backpressure** yang memerlukan intervensi.

Little's Law (`L = λW`) menyatakan: kedalaman antrian (`L`) secara langsung proporsional dengan laju kedatangan (`λ`) dan tidak langsung proporsional dengan laju pemrosesan (`W`). Ini adalah dasar matematis untuk perencanaan kapasitas.

## Core Concept: Token Bucket

Token Bucket adalah algoritma rate limiting paling umum yang memungkinkan **burst kontrol** sekaligus **enforcing batas rata-rata**.

### Bagaimana Cara Kerjanya

1. Bucket berisi `capacity B` token
2. Token ditambah secara berkelanjutan dengan laju `R` token per detik
3. Setiap permintaan menghabiskan 1 token jika tersedia
4. Jika tidak ada token cukup → permintaan ditolak

Keuntungan utama: **izinkan burst hingga B**, lalu **paksa diperlambat sampai refill**.

Contoh: capacity=3, refill=5/token/detik.

```
Request #1: Diterima (Sisa token: 2.0)
Request #2: Diterima (Sisa token: 1.0)
Request #3: Diterima (Sisa token: 0.0)
Request #4: Ditolak (Sisa token: 0.0)
```

Setelah 300ms istirahat: token kembali 1,5, sehingga Request #5 menerima.

## Core Concept: Leaky Bucket

Leaky Bucket mengalirkan beban dengan laju konstan, **menolak burst instan**.

### Bagaimana Cara Kerjanya

1. Air (permintaan) dimasukkan ke dalam bucket dengan kapasitas `C`
2. Air keluar dengan laju `leak rate R` secara terus-menerus
3. Jika level air melebihi kapasitas → permintaan ditolak

Perbedaan dengan Token Bucket:
- **Token Bucket**: izinkan burst, refill secara berkelanjutan
- **Leaky Bucket**: luruskan aliran, tolak burst

## HTTP 429: Standar RFC 6585

HTTP 429 "Too Many Requests" adalah status standar untuk respon rate limiting.

### Implementasi di Middleware

Middleware mengekstrak tenant key dari header `X-API-Key` (atau "anonymous"). Setiap tenant memiliki bucket terpisah.

Jika limit terlampaui:

```http
HTTP/1.1 429 Too Many Requests
Content-Type: application/json
Retry-After: 1

{
  "error": "rate_limit_exceeded",
  "retry_after": 1
}
```

### Mengapa Bukan Hanya IP?

RFC 6598 (CGNAT) berarti banyak pengguna di balik NAT berbagi IP. Menggunakan IP sebagai satu-satunya identifier menghambat beberapa pengguna secara tidak adil. Solusi: gunakan **tenant key atau API key** untuk isolasi yang lebih presisi.

## Core Concept: Bounded Queue (Backpressure)

Backpressure mengizinkan sistem menerima sinyal bahwa beban telah melebihi kemampuan.

### Implementasi

```go
type BoundedQueue struct {
    capacity int
    queue    chan Job
    workers  int
}
```

`TrySubmit` menggunakan pola `select-default` untuk:

- **Menerima** pekerjaan jika channel masih memiliki ruang
- **MENOLAK CEPAT** (`ErrQueueFull`) jika penuh, tanpa blocking

### Kenapa Ini Baik?

Menghindari:
- Buffer overflow / kehabisan memori
- Latency akumulasi (job menunggu berjam-jam)
- Timeout cascade di lapisan atas

## Exponential Backoff dengan Jitter

Tanpa jitter, banyak klien akan **mengirim retry bersamaan** setelah jeda serupa → *retry storm*.

### AWS Formulas (Marc Brooker)

Berbagai varian jitter:

| Varian | Formula |
|--------|---------|
| NoJitter | `min(cap, base * 2^attempt)` |
| FullJitter | `random(0, min(cap, base * 2^attempt))` |
| EqualJitter | `min(cap, base * 2^attempt) / 2 + random(0, min(cap, base * 2^attempt) / 2)` |
| DecorrelatedJitter | `min(cap, random(base, prevSleep * 3))` |

Demo menunjukkan perbedaan distribusi:

```
Attempt 0 -> NoJitter: 100ms | FullJitter: 99ms | EqualJitter: 92ms
Attempt 1 -> NoJitter: 200ms | FullJitter: 109ms | EqualJitter: 183ms
Attempt 2 -> NoJitter: 400ms | FullJitter: 111ms | EqualJitter: 202ms
Attempt 3 -> NoJitter: 800ms | FullJitter: 125ms | EqualJitter: 709ms
```

Catatan: FullJitter menyebarkan retry secara merata dalam rentang `[0, cap]`.

## Architecture

```
[Incoming Request]
        │
        ▼
[HTTP Middleware / Tenant Extractor]
        │ (RFC 6598 Tenant Key)
        ▼
[Token Bucket / Rate Limiter]
        │ (Gagal? → HTTP 429 + Retry-After)
        ▼ (Lulus)
[Bounded Queue / Backpressure Channel]
        │ (Penuh? → 503 Overloaded)
        ▼ (Dienqueue)
[Worker Pool / Consumer]
        │ (Little's Law L = λW steady state)
        ▼
[Client dengan Full Jitter Retry]
```

## Test That Proves It Works

### Token Bucket Tests

- `TestTokenBucket_BurstAndRefill`: Burst 3 token, refill ~2 token setelah 200ms
- `TestTokenBucket_ConcurrencyRace`: 50 goroutine mengecek race condition bersih

### Leaky Bucket Tests

- `TestLeakyBucket_LeakRate`: Verifikasi penggunaan kapasitas dan proses *leak*

### Registry Tests

- `TestRegistry_TenantIsolation`: Tenant A dan B punya kuota mandiri

### Bounded Queue Tests

- `TestBoundedQueue_RejectionUnderLoad`: Queue kapasitas 5, item ke-6 ditolak dengan `ErrQueueFull`
- `TestBoundedQueue_ConcurrencySafety`: 30 goroutine bersaing, semua tercatat sebagai terima atau tolak

### Backoff Tests

- `TestComputeBackoff_Bounds`: Semua varian jitter berada dalam batas yang ditetapkan
- `TestDecorrelatedJitter_Bounds`: Decorrelated Jitter mengacu pada `prevSleep` sebelumnya

### Middleware Tests

- `TestRateLimitMiddleware_RFC6585`: Verifikasi status 429 dan header `Retry-After`

## What the Tests Prove

1. **Thread-safety**: Semua komponen aman untuk goroutine bersamaan (race detector: PASS)
2. **Akurasi token**: Token bucket benar-benar melacak dan menolak sesuai kapasitas
3. **Fast rejection**: Bounded queue menolak tanpa blocking
4. **Jitter range**: Backoff berada dalam batas teori AWS
5. **HTTP compliance**: Middleware mengembalikan 429 standar RFC 6585

## Production Considerations

### Capacity Planning dengan Little's Law

Jika Anda memiliki:
- Laju kedatangan λ = 5.000.000 produk/jam
- Laju pemrosesan R = 2.000 produk/detik

Maka waktu minimum untuk menyelesaikan backlog:
```
5.000.000 / 2.000 = 2.500 detik ≈ 41 menit 40 detik
```

**Tidak ada konfigurasi queue yang dapat mengubah batas fisik ini** — hanya meningkatkan kapasitas pemrosesan atau mengurangi laju masuk.

### Queue Depth vs Oldest Job Age

10.000 job dalam antrian **bisa** baik jika usia paling tua hanya 10 detik. Masalah dimulai ketika:
- Usia job menumpuk > menit
- Rasio kedatangan > pemrosesan terus bertahan

### Multi-tenant Fairness

Tanpa isolasi, tenant Enterprise, pengguna gratis, dan layanan internal bersaing tidak adil. Solusi:
- Batas kapasitas per-tenant (misal: maks 5 pekerjaan paralel)
- Fair queueing
- Batasan konkurensi

## Common Mistakes

| Kesalahan | Dampak | Solusi |
|-----------|--------|--------|
| Rate limiting hanya berbasis IP | Menghambat banyak pengguna di CGNAT | Gunakan API key/tenant ID |
| Tanpa jitter pada retry | *Retry storm* pada saat outage | Implementasikan FullJitter atau DecorrelatedJitter |
| Antrian tidak terbatas | OOM, latency tinggi | Gunakan bounded queue dengan penggunaan kembali memori |
| Menunggu sampai timeout | *Cascading failure* | Katakan "cukup" cepat lewat 429/503 |
| Mengandalkan satu instance | Single point of failure | Untuk distribusi: Redis/Memcached untuk state terbagi |

## Case Study: Demo CLI

`cmd/demo/main.go` menampilkan empat fase utama:

### 1. Token Bucket Burst & Rate Limiting
```
Request #1: Allowed=true (Remaining Tokens: 2.0)
Request #2: Allowed=true (Remaining Tokens: 1.0)
Request #3: Allowed=true (Remaining Tokens: 0.0)
Request #4: Allowed=false (Remaining Tokens: 0.0)
After 300ms pause: Allowed=true (Remaining Tokens: 0.5)
```

### 2. Leaky Bucket Traffic Smoothing
```
Request #1: Allowed=true (Current Water Level: 1.0)
Request #2: Allowed=true (Current Water Level: 2.0)
Request #3: Allowed=true (Current Water Level: 3.0)
Request #4: Allowed=false (Current Water Level: 3.0)
```

### 3. Bounded Queue Backpressure (Load Shedding)
```
Job #1: ACCEPTED into bounded buffer
Job #2: ACCEPTED into bounded buffer
Job #3: ACCEPTED into bounded buffer
Job #4: REJECTED (Backpressure Shedding: queue capacity exceeded)
Stats: Accepted=3, Rejected=3, Processed=1
```

### 4. AWS Retry Backoff Strategies
```
Attempt 0 -> NoJitter: 100ms | FullJitter: 99ms | EqualJitter: 92ms
Attempt 1 -> NoJitter: 200ms | FullJitter: 109ms | EqualJitter: 183ms
```

## Checklist

- [x] Memahami laju masuk vs kapasitas pemrosesan
- [x] Memilih algoritma rate limiting yang tepat (token bucket untuk burst control, leaky bucket untuk smoothing)
- [x] Mengimplementasikan tenant-based key untuk multi-tenant fairness
- [x] Menyediakan response 429 standar dengan `Retry-After` header
- [x] Menggunakan bounded queue untuk backpressure fast rejection
- [x] Menambahkan jitter pada exponential backoff
- [x] Memantau metrik: queue depth, oldest job age, arrival rate vs processing rate
- [x] Menjalankan race detector untuk thread-safety verification

## Key Takeaways

1. Rate limiting dan backpressure adalah **prinsip desain inti** untuk sistem yang dapat diandalkan.
2. **Token Bucket** dominates untuk kontrol burst; **Leaky Bucket** untuk traffic smoothing.
3. **HTTP 429** standar RFC 6585 harus disertakan `Retry-After` header.
4. **Bounded queue** mencegah kehabisan memori dengan `TrySubmit` non-blocking.
5. **Exponential backoff tanpa jitter** menciptakan retry storm.
6. Gunakan **FullJitter** untuk distribusi retry paling merata.
7. Isolate rate limiting dengan **tenant key** (bukan IP-only) untuk menghindari dampak CGNAT.
8. Little's Law (`L = λW`) adalah dasar matematis untuk perencanaan kapasitas.
9. Pantau **oldest job age**, bukan hanya queue depth.
10. Katakan "cukup" secara terukur sebelum sistem kehabisan sumber daya.

## Sources

- RFC 6585: https://datatracker.ietf.org/doc/html/rfc6585
- RFC 6598 (CGNAT): https://datatracker.ietf.org/doc/html/rfc6598
- AWS Architecture Blog: https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/
- Wikipedia: Token Bucket, Leaky Bucket, Rate Limiting, Little's Law, Reactive Streams