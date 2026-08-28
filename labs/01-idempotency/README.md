# Lab 01: Idempotency

> **Problem**: Duplicate payment requests create multiple charges.

## Duplicate HTTP Requests

| Scenario | Description | Risk |
|----------|-------------|------|
| **Network timeout** | Request succeeds on server but response lost | Client retries, duplicate charge |
| **Frontend retry** | UI retry button clicked multiple times | Same as timeout |
| **Queue redelivery** | Message queue requeues unacknowledged message | Duplicate processing |
| **Payment webhook retry** | Gateway retries webhook delivery | Multiple webhook handlers |
| **Concurrent identical requests** | Race condition: requests start simultaneously | Both succeed |

## POST /orders/{id}/pay Example

## Scenario

A client sends a payment request. The server processes it, calls the payment gateway, and the gateway succeeds — but the response is lost due to network timeout. The client retries with the same intent. **Without idempotency, a second charge is created.**

### Timeline Example

```
Request: POST /orders/order-123/pay
Body: {"amount": 500000}

Timeline:
  T0  Client → POST /orders/order-123/pay {amount: 500000}
  T1  Server processes, calls gateway
  T2  Gateway success, response lost (timeout)
  T3  Client retries POST /orders/order-123/pay {amount: 500000}
  T4  Without idempotency: 2nd charge created  ❌ (duplicate 500,000)
  T4  With idempotency: 2nd charge prevented   ✅ (returns first result)
```

## Perbedaan Konsep

| Konsep | Definisi | Cara Kerja | Kelemahan |
|--------|----------|------------|-----------|
| **Idempotency** | Request yang diulang hasilnya sama | Key unik + replay response | Key collision, TTL habis |
| **Database Transaction** | Operasi yang bersifat atomic | BEGIN → COMMIT/ROLLBACK | Tidak cegah duplicate, hanya atomic |
| **Unique Constraint** | Database mencegah duplicate key | UNIQUE(index) error jika duplikat | Error 500/jika tidak ditangani dengan baik |
| **Deduplication** | Bersihkan data duplikat setelah | Proses batch/periodik | Data duplikat masuk dulu, lambatnya |

## Root Cause

HTTP retry patterns (timeout → retry) can cause duplicate processing. The server has no way to distinguish "retry of an existing request" from "new request."

## What To Look At

```bash
# Reproduce the bug (unsafe implementation)
go test ./labs/01-idempotency/unsafe/... -v -count=1

# Fix is verified
go test ./labs/01-idempotency/safe/... -v -count=1

# Full test suite with race detector
go test ./labs/01-idempotency/... -race -count=1

# Benchmark safety vs overhead
go test ./labs/01-idempotency/... -bench=. -benchmem -count=1

# Run integration tests
go test ./labs/01-idempotency/tests/... -v -count=1
```

## Solutions Implemented

### Safe implementation

1. **Idempotency-Key header** — Client provides unique key per request
2. **Payload hash** — Hash of request body tied to idempotency key
3. **Unique constraint** — `UNIQUE(idempotency_key)` on payments table
4. **Database transaction** — Idempotent insert within transaction
5. **Concurrent request handling** — Serialize by idempotency key
6. **Cached/replayed response** — Skip processing, return stored result
7. **Conflict detection** — Detect payload mismatch on same key

### Architecture

```
                    ┌─────────────────┐
  Idempotency-Key   │   Redis Cache   │
  ─────────────────▶│  (fast path)    │
                    └────────┬────────┘
                             │ cache miss
                             ▼
                    ┌─────────────────┐
  Idempotency Key   │   DB Lookup     │
  ─────────────────▶│  (source of    │
                    │   truth)       │
                    └────────┬────────┘
                             │ exists
                             ▼
                    ┌─────────────────┐
                    │  Replay cached  │
                    │    response     │
                    └─────────────────┘
```

## Trade-offs

| Aspect | Trade-off |
|--------|-----------|
| Storage | Stores response payload (small overhead) |
| Complexity | Concurrent request handling adds locking complexity |
| Cache invalidation | TTL-based; Redis failure falls back to DB |
| Key rotation | Client must generate new key per logical operation |

## Production Considerations

1. **Key format**: Use UUIDv4 or ULID (monotonic, sortable)
2. **TTL**: Default 24-72h (covers most client retry windows)
3. **Response size**: Consider storing only result ID, not full payload
4. **Cleanup**: Background job to remove expired keys
5. **Monitoring**: Alert on high idempotency hit rate (indicates client issues)
6. **Scope design**: Use composite `UNIQUE(scope, key)` rather than global unique
   - Example scopes: `tenant_id`, `user_id`, `endpoint:orders/pay`

## Mengapa Global Key Kurang Ideal untuk Multi-User/Multi-Endpoint

Dengan `UNIQUE(key)` saja (tanpa `scope`), ada 3 masalah:

1. **Key Collision antar User/Tenant**:
   - Jika dua user berbeda yang kebetulan mengirim request dengan idempotency key yang sama
     (misal: client yang menggunakan generator timestamp yang sama-sama menghasilkan ID serupa),
     salah satu request bisa saling bentrok dan mengakibatkan salah satu gagal atau malah
     me-return data dari user lain.

2. **Multi-Endpoint Namespace Overlap**:
   - Idempotency key `pay-001` untuk endpoint `/orders/pay` bisa berbeda konteks dengan
     `pay-001` untuk endpoint `/subscriptions/pay`. Tanpa scope, kedua key ini akan
     bentrok dan menghasilkan perilaku yang tidak terduga.

3. **Security & Isolation**:
   - Pada lingkungan multi-tenant, satu tenant harus tidak dapat menebak atau menimpa
     idempotency key milik tenant lain. Scope memastikan isolasi antar tenant.

## What Can Still Fail

1. **Key collision**: Client generates duplicate keys for different requests
   - Mitigation: Include payload hash validation
2. **Mid-processing crash**: Key stored but response not persisted
   - Mitigation: Use database transactions for atomic insert
3. **Redis failover**: Cache lost during failover
   - Mitigation: Always check DB after cache miss
4. **Large payload**: Storing full response in Redis exhausts memory
   - Mitigation: Store compact result identifiers only
5. **Clock drift**: TTL expiration while request in flight
   - Mitigation: Use reasonable TTLs, handle edge cases gracefully

## References

- [Stripe Idempotency](https://stripe.com/docs/api/idempotent_requests)
- [AWS Idempotency](https://aws.amazon.com/builders-library/idempotency/)
- [RFC 7231 - Idempotency](https://datatracker.ietf.org/doc/html/rfc7231#section-4.2.2)

## Alternatif: Sistem Partial Payment (Cicilan)

> **Catatan**: Tahap ini hanya dokumentasi untuk referensi jika sistem mendukung cicilan. Implementasi tidak dilakukan pada lab ini untuk menghindari pencampuran konsep.

### Perbedaan Business Rule

Dalam sistem normal (full payment):
```
UNIQUE(order_id) pada tabel payments
```

Dalam sistem partial payment (cicilan):
```
order_id tidak unique
Rule: SUM(payment.amount) <= order.outstanding_amount
```

### Race Condition pada Concurrent Partial Payment

```
Timeline:
  T0  Worker A: SELECT SUM(amount) FROM payments WHERE order_id='O123' → returns 800,000
  T1  Worker B: SELECT SUM(amount) FROM payments WHERE order_id='O123' → returns 800,000
  T2  Order O123: total=1,000,000 (outstanding=200,000 remaining)
  T3  Worker A: INSERT payment {amount: 300,000} → SUM=1,100,000 > 1,000,000 ❌ OVERSPEND
  T4  Worker B: INSERT payment {amount: 300,000} → SUM=1,100,000 > 1,000,000 ❌ OVERSPEND
```

**Masalah**: Kedua cicilan bertentangan, total payment melebihi outstanding amount.

### Solusi yang Diperlukan

1. **Pessimistic Lock** pada order row:
   ```sql
   SELECT ... FOR UPDATE
   INSERT payments...
   ```

2. **Database Constraint** dengan trigger:
   - Trigger menghitung total payment sebelum insert
   - Abort transaction jika melebihi outstanding

3. **Distributed Lock** (Redis) sebelum proses:
   ```redis
   SET lock:order:O123 "worker-1" NX EX 30
   ```

4. **Two-Phase Commit** dengan verifikasi di akhir

---

## Idempotency Expiration Policy

### Konfigurasi TTL (Time To Live)

TTL bervariasi tergantung jenis operasi:

| Entity Type | TTL Default | Alasan |
|-------------|-------------|--------|
| **Payment** | 72 hours | Prioritaskan auditability untuk financial record. Klien mungkin butuh waktu 1-2 hari untuk merespon timeout. |
| **Invoice Creation** | 168 hours (7 days) | Invoice mungkin direview oleh tim keuangan sebelum dikirim ke klien. |
| **File Upload** | 24 hours | User biasanya menunggu konfirmasi upload. Ukuran file kecil. |
| **Generic Command** | 1 hour | Operasi ephemeral seperti webhook callback. |

### Perbedaan TTL untuk Financial Record

Untuk payment, kita **tidak** memaksimalkan cleanup karena:

1. **Audit Trail**: Payment harus dapat direproduksi selama compliance period
2. **Dispute Resolution**: Customer support mungkin butuh lihat request lama
3. **Reconciliation**: Bank statements perlu cocok dengan idempotency keys
4. **Replay Security**: Key yang habis TTL bisa diaploitasi untuk replay attack

### Cleanup Mechanism

```go
// Periodic cleanup job (run every hour)
func cleanupJob() {
    removed, _ := store.Cleanup(ctx)
    log.Printf("Cleanup: removed %d expired idempotency keys", removed)
}

// Monitoring expired keys
expiredCount := store.ExpiredKeysCount(ctx)
if expiredCount > 10000 {
    alert("High expired idempotency key count: may indicate config issue")
}
```

---

## Metrics (Prompt 023)

| Metric | Type | Labels | Purpose |
|--------|------|--------|---------|
| `idempotency_requests_total` | Counter | `operation`, `status` | Track request volume and outcomes |
| `idempotency_replays_total` | Counter | `operation` | Measure cache hit rate |
| `idempotency_conflicts_total` | Counter | `operation`, `reason` | Track payload mismatches |
| `idempotency_processing_total` | Gauge | `operation` | Current in-flight requests |
| `idempotency_failures_total` | Counter | `operation`, `type` | Track error patterns |

**Label constraints**: Never include user IDs, order IDs, or other business identifiers as labels. This would create unbounded cardinality in Prometheus TSDB.

---

## Benchmark (Prompt 024)

### Test Scenarios

| Scenario | Description | Expected Behavior |
|----------|-------------|-------------------|
| `BenchmarkNormalRequest` | Fresh key per request | Full gateway cost each iteration |
| `BenchmarkIdempotentNewRequest` | Same key, same payload | First: gateway, rest: cached |
| `BenchmarkIdempotentReplay` | Replay existing record | No gateway calls, pure cache lookup |
| `BenchmarkConcurrentRequestsSameKey` | Concurrent access | Serialized by in-flight guard |

### Limitations

**These benchmarks measure in-process overhead only. Do NOT draw production conclusions without:**

1. **Real database latency**: Map lookups (~1μs) ≠ SQL queries (~1-10ms)
2. **Redis round-trip**: Network latency adds 0.5-2ms per cache operation
3. **Concurrent load**: Single-threaded benchmarks hide lock contention
4. **Memory pressure**: Large payload storage affects GC pauses
5. **Gateway latency**: Real payment gateway = 100-500ms (dominates idempotency cost)

**Rule of thumb**: Idempotency overhead is <15% for new requests, 95%+ reduction for replays.