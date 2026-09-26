# Content Brief

**Topic:** Circuit Breaker Pattern — Implementasi dan Demonstrasi di Go

**Target Reader:** Backend Engineer, SRE, dan Software Architect yang mengimplementasikan resilience patterns di sistem terdistribusi.

**Problem:** Panggilan remote across jaringan gagal atau timeout. Tanpa proteksi, caller memblokir menunggu timeout, memegang thread, socket, dan memory hingga sumber daya sistem terdepleksi dan kegagalan menyebar ke seluruh sistem (cascade failure).

**Core Mental Model:** Circuit Breaker memantau kegagalan downstream, mendorong (trip) OPEN untuk memblokir semua pemanggilan ketika threshold tercapai (fail-fast dalam microsecond), lalu mengirim probe terbatas (probe) sebelum mengembalikan trafik ke normal.

**Approved Research Status:** APPROVED (research-audit/07-verdict.md)

**Approved Engineering Status:** APPROVED (engineering-audit/06-verdict.md)

**Main Concepts:**
- Tiga state machine: CLOSED, OPEN, HALF_OPEN
- Fail-fast behavior ketika OPEN (tidak ada jaringan call)
- Cooldown period dengan configurable timeout
- Probe calls terbatas saat HALF_OPEN untuk menguji recovery
- Thread-safe state transitions menggunakan sync.Mutex
- Consecutive failure counting (bukan sliding window)
- Generation-based invalidation untuk trailing in-flight requests
- Panic safety pada saat execute

**Verified Behaviors:**
- CLOSED: semua request routing ke downstream; failure increment counter; success mereset counter; transit ke OPEN saat failures >= FailureThreshold
- OPEN: semua request langsung gagal dengan `ErrCircuitOpen` tanpa network call; memulai cooldown timer
- HALF_OPEN: mengizinkan probe terbatas (HalfOpenMaxCalls); success → CLOSED; failure → OPEN
- Zero downstream requests dieksekusi ketika circuit OPEN
- Concurrency-safe di bawah race detector (16 unit + 2 integration test, 0 data race)
- Panic tidak corrupt state; counters diupdate before re-panic
- Trailing in-flight request (generasi lama) tidak corrupt state baru
- Default config: FailureThreshold=3, OpenTimeout=300ms, HalfOpenMaxCalls=1
- Demo output verified: fail-fast ~40-125ns, downstream_calls terhenti setelah OPEN

**Available Case Studies:**
- Skenario 1: Tanpa Circuit Breaker — slow dependency menyebabkan blocking timeout (~100ms/request)
- Skenario 2: Dengan Circuit Breaker — fail-fast setelah 3 failures, request berikutnya < 1µs
- Skenario 3: Recovery — HALF_OPEN → CLOSED saat probe sukses
- Skenario 4: Failed Recovery — HALF_OPEN → OPEN saat probe gagal
- Skenario tambahan: slow dependency timeout (integration test)

**Warnings:**
- Nilai timeout demo (100ms HTTP, 300ms cooldown) bersifat illustratif untuk testing cepat; produksi harus tune ke SLA dan recovery profile
- Implementasi pakai consecutive failure counting saja, bukan sliding window atau error rate
- Circuit breaker tidak membedakan error 4xx vs 5xx — error predicate tidak ada; semua non-nil error increment counter (dokumentasi educational limitation)
- Observability metrics (circuit_state, circuit_open_count, dll.) adalah rekomendasi arsitektur, tidak di-export di kode
- DefaultConfig OpenTimeout = 300ms (bukan 5s) — sesuai kode sumber

**Approved Research Sources (Tier 1):**
- Martin Fowler, *Circuit Breaker* (2014)
- Microsoft Azure Architecture Center, *Circuit Breaker Pattern* (2025-02)
- Microsoft Azure, *Retry Pattern* (2024-07)
- Microsoft Azure, *Bulkhead Pattern* (2026-03)
- Resilience4j, *CircuitBreaker*
- Google SRE Book, *Handling Overload* (2017)
- Go Standard Library, *net/http.Client*
- AWS Builders Library, *Timeouts, retries, and backoff with jitter* — canonical https://builder.aws.com/content/3EumjoZascWd1oZiEgL8ORlv3qE/timeouts-retries-and-backoff-with-jitter (301 from https://aws.amazon.com/builders-library/timeouts-retries-and-backoff-with-jitter/)

**Approved Engineering Artifacts:**
- engineering/01-design.md — design dan success criteria
- engineering/02-implementation-notes.md — keputusan implementasi
- engineering/03-execution-result.md — hasil eksekusi demo dan test
- engineering-revision/02-changes-made.md — revisi: generation tracking, panic safety, slow dependency test, concurrency test, demo formatting