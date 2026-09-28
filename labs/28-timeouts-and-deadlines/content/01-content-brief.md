# Content Brief

Topic: Timeouts and Deadlines — deadline propagation, timeout budget, exponential backoff dengan Full Jitter, circuit breaker, dan deduplikasi idempotency di Go
Target Reader: Backend engineer yang membangun service terdistribusi dengan dependency downstream yang dapat melambat
Problem: Satu dependency lambat menahan worker/connection pool, memicu cascading failure ke endpoint yang tidak terkait
Core Mental Model: Timeout bukan angka arbitrer, melainkan guardrail proteksi resource dan budget deadline yang diwariskan ke downstream; `timeout == unknown`, bukan `timeout == failed`
Approved Research Status: APPROVED (research-audit/07-verdict.md, 2026-09-28, 0 unsupported claims, 3 LOW gaps)
Approved Engineering Status: APPROVED (engineering-audit/06-verdict.md, 0 failures, `go test -race ./...` PASS, demo PASS)
Main Concepts:
- Resource exhaustion via Little's Law (L = λW)
- Granular network timeout vs total timeout
- Timeout budgeting dari distribusi P95/P99
- Retry storm dan Full Jitter (`sleep = rand_between(0, min(cap, base * 2^(attempt-1)))`)
- Timeout ambiguity dan idempotency key
- Deadline propagation (relatif, bukan timestamp absolut)
- Guardrail database (`statement_timeout`, `lock_timeout`, `idle_in_transaction_session_timeout`) dan worker (execution timeout, DLQ) sebagai konteks riset — tidak diimplementasikan di lab ini
Verified Behaviors:
- `ExecuteWithBudget` membatalkan kerja saat parent atau child context kedaluwarsa
- `Retrier.Do` sukses setelah transient error, berhenti di `MaxAttempts`, dan abort saat context dibatalkan
- `Breaker` CLOSED → OPEN → HALF_OPEN → CLOSED, dan gagal di HALF_OPEN langsung kembali OPEN
- `Store` Get/Set/TTL expiry dan aman diakses 50 goroutine konkuren (race detector PASS)
- Integrasi retry+circuit membuka circuit; retry+idempotency hanya mengeksekusi 1 kali
- Demo `cmd/demo` 4 skenario berjalan sesuai `engineering/03-execution-result.md`
Available Case Studies:
- Demo 1: budget 50ms vs kerja 80ms → `context deadline exceeded`
- Demo 2: 3 attempt dengan Full Jitter → sukses
- Demo 3: 2 failure → OPEN → cooldown 50ms → HALF_OPEN → sukses → CLOSED
- Demo 4: `req-tx-99231` dieksekusi sekali, retry ter-deduplikasi
Warnings:
- Store idempotency in-memory, hilang saat crash; produksi butuh DB transaksional atau Redis dengan TTL
- Lab tidak mendemonstrasikan koordinasi antar-proses/jaringan (header `grpc-timeout`, HTTP deadline header)
- Propagasi deadline HTTP/REST heterogen (`Request-Timeout`, baggage OpenTelemetry) vs gRPC native — belum standar tunggal
- Adaptive timeout berbasis P99 moving-window masih open question
- Nilai timeout di lab (20ms, 50ms, 100ms) ilustratif untuk lab, bukan rekomendasi produksi
- URL Stripe kanonis adalah `https://docs.stripe.com/error-handling`
