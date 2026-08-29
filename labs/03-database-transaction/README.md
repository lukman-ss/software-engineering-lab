# Lab 03 — Distributed Transaction: Kenapa Database Transaction Saja Tidak Cukup?

> **Mental Model**: Database transactions (ACID) menjamin atomicity **hanya pada transactional resource yang berpartisipasi dalam transaction tersebut**. External HTTP APIs, email, WhatsApp, object storage, dan kebanyakan external systems tidak ikut berpartisipasi dalam local transaction database aplikasi.

---

## 1. Local Database Transaction

### Flow yang Benar
```
[Local Transaction]
     │
     ├── Payment
     ├── Inventory
     └── Wallet Tx
```

Jika ada error → **ROLLBACK** → Semua pembatalan bersama.

---

## 2. Transaction Boundary

**Tidak menghabiskan semua di satu database juga bukan solusi terbaik jika ada external HTTP call.**

```
[Unsafe Distributed Flow]
     │
     ├── Payment
     ├── Inventory
     ├── [BEGIN]
     ├── ...business logic...
     ├── [HTTP WhatsApp] ────► External World
     ├── ERP sync
     └── [COMMIT/Rollback]
            │
            X process crash

[Result]
ROLLBACK DB          ← Payment & Inventory dibatalkan
WhatsApp Tetap Terkirim  ← TIDAK dapat di-rollback!
```

---

## 3. Anti-Pattern: HTTP Call Di Dalam Transaction

```
[BEGIN TRANSACTION]
     │
     ├── UPDATE invoice SET status = 'paid'
     │
     ├── HTTP call ke payment gateway (50ms-5detik)
     │
     └── [COMMIT]
```

**Risiko**: Transaction bertahan lama → Connection pool tertahan → Throughput menurun.

---

## 4. Dual-Write Problem

```
[Database]                    [Message Queue]
   │                               │
   ├── INSERT payment            │
   │                               │
   └── COMMIT ────► success       │
            │                      │
            X process crash        │
            │                      │
            └─────────────────────► Event TIDAK pernah dipublish
```

**Hasil**: Invoice = Paid, tapi **event tidak pernah sampai ke consumers**.

---

## 5. Transactional Outbox Pattern

```
[BEGIN TRANSACTION]
     │
     ├── UPDATE invoice SET status = 'paid'
     │
     ├── INSERT INTO outbox_events (...)
     │
     └── [COMMIT]
              │
              ▼
[Outbox Dispatcher]
     │
     ├── Read pending event
     ├── Publish ke broker
     └── UPDATE status = 'published'
```

```
[Dispatcher Flow]
         │
         ▼
   Event Broker ◄────────────────────┐
         │                           │
         ├─► [Consumer: Inventory]    │
         │                           │
         ├─► [Consumer: Commission] │
         │                           │
         ├─► [Consumer: ERP]         │
         │                           │
         └─► [Consumer: Notification]
```

---

## 6. Saga & Compensation

```
[Saga Flow]
   │
   ▼
Step A: Reserve Budget
   │
   ▼
Step B: Process Payment (❌ gagal)
   │
   ▼
[Compensate B: Refund]
   │
   ▼
[Compensate A: Release Budget]
```

---

## 7. Event Naming Convention (Poin 25)

Gunakan **past-tense naming** untuk events yang merepresentasikan fakta yang sudah terjadi:

✅ **Baik** (past tense / fact-based):
- `InvoicePaid`
- `VendorPaymentCompleted`
- `PaymentRecorded`
- `InventoryUpdated`
- `CommissionGenerated`

❌ **Jangan** (command-style):
- `DoInventory`
- `SendWhatsappNow`
- `ProcessPayment`

**Alasan**: Event harus mencerminkan **state yang sudah ada**, bukan **perintah untuk melakukan sesuatu**. Consumers harus dapat memproses event dua kali (duplicate delivery) dan hasilnya tetap konsisten.

---

## 8. Queue ≠ Exactly-Once Delivery (Poin 26)

Message queue **bukan** solusi magic untuk exactly-once delivery.

```
[Queue]
   │
   ▼
[Consumer]
   │
   ├── Receive message
   ├── Process business logic
   └── ??? crash sebelum ACK ???
   │
   ▼
Queue kirim message lagi (at-least-once)
```

**Implikasi**: Consumer harus **idempotent**. Jangan andalkan queue untuk deduplikasi!

---

## 9. Event-Driven Trade-offs (Poin 23)

### Keuntungan (+)
- **Resilience**: Proses dapat melanjutkan setelah crash
- **Loose Coupling**: Consumers tidak harus tahu satu sama lain
- **Retryable**: Gagal bisa dicoba ulang tanpa membatalkan transaksi utama
- **Scalable Independently**: Consumer bisa scale terpisah

### Risiko / Kekurangan (-)
- **Eventual Consistency**: Data tidak konsisten secara serentak
- **Operational Complexity**: Monitoring harus mencakup multiple services
- **Duplicate Messages**: At-least-once delivery berarti ada duplikasi
- **Harder Debugging**: Flow menyebar di banyak service
- **Ordering Problems**: Event mungkin sampai out-of-order
- **Observability Requirement**: Diperlukan distributed tracing

> **Senior Engineer Mindset**: Pilih event-driven architecture bila memang membutuhkan retry, compensation, atau loose coupling. Jangan gunakan hanya karena "modern".

---

## 10. Eventual Consistency (Poin 24)

### Contoh Scenario
```
Time 0:  Invoice status = PAID (di database)

Time 50ms: Commission belum generated
Time 100ms: ERP belum synced  
Time 200ms: WhatsApp belum sent
```

**Itu bukan bug!** Itu adalah **eventual consistency** yang normal pada sistem async.

### Kapan Acceptable?
✅ Notification services (WhatsApp, Email)
✅ Analytics / Reporting
✅ Audit logs
✅ Search index updates

### Kapan Butuh Strong Consistency?
❌ Payment processing
❌ Inventory stok fisik
❌ Balance / akuntansi
❌ SLA-critical business logic

---

## 11. No External Infrastructure Required (Poin 27)

Lab ini **tidak memerlukan**:
- RabbitMQ
- Kafka
- Redis
- WhatsApp/Google credential
- Midtrans/Stripe credential
- ERP system
- Docker cluster

Semua diimplementasikan menggunakan **in-memory fake broker** yang deterministik untuk testing.

### Cara Menjalankan:
```bash
cd labs/03-database-transaction
go test -v -count=1
```

---

## Menjalankan Semua Test

```bash
go test -v ./labs/03-database-transaction/... -count=1
```

| Test | Bukti |
|------|-------|
| `TestUnsafeLocalTransaction` | Partial state corruption tanpa transaksi |
| `TestSafeLocalTransaction` | ACID rollback yang bersih |
| `TestDistributedTransactionExternalSideEffectLimitation` | DB rollback ≠ WhatsApp rollback |
| `TestHTTPCallInsideTransactionLatency` | Transaction berlama 50ms karena HTTP latency |
| `TestDualWriteProblemEventLost` | Invoice 'paid' tapi event tidak pernah terkirim |
| `TestTransactionalOutboxPatternAtomicity` | Business state + outbox disimpan atomik |
| `TestOutboxDispatcherPublishesPending` | Dispatcher publikasi event pending |
| `TestOutboxDuplicateDeliveryAtLeastOnce` | Duplicate delivery terjadi |
| `TestIdempotentConsumerDeduplication` | Idempotent consumer deduplikasi |
| `TestTransientFailureSuccessAfterRetry` | Retry mekanisme bekerja |
| `TestDeadLetterQueue` | Event gagal masuk DLQ |
| `TestSagaPaymentWithCompensatingAction` | Compensasi dijalankan saat step gagal |

---

## Navigasi

- **Previous**: [Lab 02 — Database Index](../02-database-index/)
- **Next**: [Lab 04 — Caching](../04-caching/)