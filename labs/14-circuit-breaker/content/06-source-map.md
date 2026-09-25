# Source Map (Peta Sumber Konten)

Dokumen ini memetakan setiap bagian artikel utama (`02-master-draft.md`) dengan basis bukti riset, file implementasi arsitektur terverifikasi, serta pengujian yang disahkan.

---

## Masalah (Problem) & Kegagalan Beruntun (Cascade Failure)

**Research:**
- `research/04-cascade-failure.md`
- `README.md` (Bagian: Problem, Cascade Failure)

**Implementation:**
- `internal/payment/fake_server.go` (Mode `ModeSlow` dan efek `time.Sleep` pada *thread pool*)
- `engineering/01-design.md` (Bagian: Failure Scenario)

**Tests:**
- `cmd/demo/main.go` (Skenario 1: Tanpa Circuit Breaker, melacak durasi dan blokade eksekusi)

---

## Mental Model, Konsep Inti (State), dan Arsitektur

**Research:**
- `research/03-core-concepts.md`
- `research/05-circuit-states.md`
- `README.md` (Bagian: Circuit Breaker Mental Model, CLOSED, OPEN, HALF-OPEN, Architecture)

**Implementation:**
- `internal/circuitbreaker/circuit_breaker.go` (Definisi struct, `State` constants, `Execute` switch-case logika)
- `internal/checkout/service.go` (Penerapan proxy pada layanan)

**Tests:**
- `internal/circuitbreaker/circuit_breaker_test.go` (Unit test transisi mesin status 1-10)

---

## Timeout vs Retry vs Circuit Breaker

**Research:**
- `research/06-timeout-retry-backoff.md`
- `README.md` (Bagian: Timeout vs Retry vs Circuit Breaker)
- `research/02-sources.md` (Ref AWS Builders Library)

**Implementation:**
- `internal/payment/client.go` (Penetapan HTTP `Timeout` klien sebesar 50ms)

---

## Pertimbangan Produksi (Production Considerations)

**Research:**
- `research/08-observability.md`
- `research-audit/07-verdict.md` (Peringatan non-blocking terkait klaim)
- `engineering/02-implementation-notes.md` (Bagian: Known Limitations, Trade-offs)

**Implementation:**
- `internal/circuitbreaker/circuit_breaker.go` (Penggunaan *simple consecutive count* via `failureCount`, dan sinkronisasi `sync.Mutex`)

---

## Bukti Pengujian Terverifikasi (What the Tests Prove)

**Engineering / Tests:**
- `engineering/03-execution-result.md` (Output bukti log eksekusi, Unit, dan Race Test)
- `internal/circuitbreaker/circuit_breaker_test.go` (`TestCircuitBreaker/11._concurrency_and_race_safety` memastikan tidak ada *data race*)
- `tests/integration_test.go` (`TestCircuitBreakerIntegration` dan pembuktian durasi penolakan fail-fast)
- `cmd/demo/main.go` (Durasi rekam nyata: >100ms Timeout vs <1µs Open State Fail-Fast)

---

## Fallback dan Studi Kasus (CMMS & PPOB)

**Research:**
- `research/07-fallback-bulkhead.md`
- `research/10-final-research.md` (Catatan arsitektur dan mitigasi isolasi/bulkhead)
- `README.md` (Bagian: Fallback, Bulkhead, CMMS Example, PPOB Example)

**Audit Validation:**
- `research-audit/07-verdict.md` (Peringatan agar `silent fallback` pada kondisi krusial ditangani secara khusus dan tidak digeneralisasi keliru).

---

## Kesalahan Umum (Common Mistakes)

**Research:**
- `research/09-failure-modes.md`
- `README.md` (Bagian: Failure Modes - ambang terlalu rendah, probe berlebih, error keliru).
