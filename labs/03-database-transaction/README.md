# Lab 03 — Distributed Transaction: Kenapa Database Transaction Saja Tidak Cukup?

> **Mental Model**: Database transactions (ACID) garantikan atomicity **hanya dalam satu transactional resource (satu database)**. Ketika workflow melibatkan sistem eksternal (API, webhook, messaging broker, notifier), database transaction tidak dapat mengelilingi semua perubahan secara atomik.

---

## 1. Fondasi: Local Database Transaction

1. Buat pembayaran di tabel `payments`
2. Update status order di tabel `orders`
3. Buat transaksi di dompet digital `wallet_transactions`

### A. Unsafe Local Transaction (Partial State)
Tanpa transaksi, jika error terjadi di tengah, sistem meninggalkan **partial state** yang tidak konsisten.

### B. Safe Local Transaction (ACID Rollback)
Dengan `BEGIN`, `COMMIT`, `ROLLBACK` (atau helper `WithTx`), seluruh perubahan dalam satu transaksi dibatalkan bersama jika ada error.

> **Kesimpulan**: Database transaction efektif **jika dan hanya jika** seluruh perubahan berada dalam transactional resource yang sama.

---

## 2. External Side Effect - Tidak Bisa Diretryakan Database-nya

```
BEGIN
  ├── INSERT INTO payments ...
  ├── UPDATE invoices SET status = 'paid'
  ├── Send WhatsApp notification (SUCCESS – side effect terjadi!)
  └── ERP service ERROR → ROLLBACK
```

**Hasil**:
- Database: `ROLLBACK` → 0 payment, invoice tetap unpaid.
- External: WhatsApp tetap terkirim → `DB rollback ≠ external rollback`.

---

## 3. Anti-Pattern: HTTP Call Di Dalam Transaction (§4)

```
BEGIN TRANSACTION
  ├── UPDATE invoice SET status = 'paid'
  ├── HTTP call ke payment gateway (lambat! 50ms-5detik)
  └── COMMIT
```

### Risiko
- Transaction duration bertambah → Connection pool tertahan
- Row lock bertahan lama → Contention & deadlock
- Downstream latency ikut mempengaruhi DB → Timeout
- Throughput menurun secara sistemik

---

## 4. Masalah Dual-Write: DB Write + Message Publish (§5)

```
BEGIN
  ├── UPDATE invoice SET status = 'paid' WHERE order_id = 101
  └── COMMIT (invoice = 'paid' di database)
  ├── Publish 'InvoicePaid' event ke message broker
  └── SIMULASI CRASH → Event TIDAK pernah dikirim!
```

**Hasil**: Database `invoice = paid`, tetapi **event tidak pernah terkirim**.

---

## 5. Transactional Outbox Pattern (§6)

### Schema Tabel Outbox
```sql
CREATE TABLE outbox_events (
    id            TEXT PRIMARY KEY,
    event_type    TEXT,
    aggregate_id  TEXT,
    payload       TEXT,
    status        TEXT CHECK (status IN ('pending', 'published')),
    attempts      INT DEFAULT 0,
    created_at    TIMESTAMPTZ,
    published_at  TIMESTAMPTZ NULL
);
```

### Flow dengan Outbox
```
BEGIN TRANSACTION
  ├── UPDATE invoices SET status = 'paid'
  ├── INSERT INTO outbox_events (event_type='InvoicePaid', status='pending', ...)
  └── COMMIT
```

Jika commit berhasil → **KEDUA** berada di database.
Jika transaction gagal → **KEDUA** tidak tersimpan.

---

## 6. Outbox Dispatcher & Retry Mechanism (§10)

```
Dispatcher:
  1. Baca event dengan status = 'pending'
  2. Publish ke broker
  3. Jika gagal → increment attempts, retry hingga max_attempts
  4. Jika sukses → UPDATE status = 'published'
```

### Konfigurasi
- `maxAttempts`: batas maksimum retry (misal: 3).
- Tidak ada retry loop tanpa batas.

---

## 7. Duplicate Delivery (§8) – At-Least-Once Delivery

### Skenario
```
1. Dispatcher publish event berhasil
2. Process crash SEBELUM event ditandai 'published'
3. Event tetap 'pending' di database
```

### Hasil
Event `InvoicePaid(event-123)` dipublish **dua kali**.

> **Catatan**: Transactional Outbox menghasilkan **at-least-once delivery**, bukan exactly-once!

---

## 8. Idempotent Consumer (§9) – Menghubungkan dengan Lab #01 Idempotency

### Mekanisme Deduplication
Buat tabel `processed_events`:
```sql
CREATE TABLE processed_events (
    consumer_name TEXT,
    event_id      TEXT,
    processed_at  TIMESTAMPTZ,
    PRIMARY KEY (consumer_name, event_id)
);
```

### Flow Consumer
```
Maju terima Event InvoicePaid
    ↓
Cek: sudah pernah diproses?
    ↓
Jika BELUM:
    ├── Hitung komisi sekali
    └── INSERT INTO processed_events
    ↓
Jika SUDAH:
    └── SKIP (idempotent!)
```

---

## 9. Dead Letter Queue (DLQ) – Event yang Konsisten Gagal (§11)

Jika event gagal publish meski sudah mencapai `maxAttempts`, kirim ke **dead letter queue**.

### Schema DLQ
```sql
CREATE TABLE dead_letter_events (
    event_id    TEXT PRIMARY KEY,
    event_type  TEXT,
    payload     TEXT,
    reason      TEXT,
    attempts    INT,
    failed_at   TIMESTAMPTZ
);
```

### Flow
```
Dispatcher:
  1. Publish attempt 1 → Gagal
  2. Publish attempt 2 → Gagal
  3. Publish attempt 3 → Gagal (maxAttempts tercapai)
  4. MOVE to DLQ + log alasan
```

DLQ memungkinkan operasi manual / analisis terhadap event yang secara otomatis tidak dapat diproses.

---

## 10. Saga Pattern & Compensating Transactions (§12)

### Saga yang Sederhana: Vendor Payment

```
Step 1: Reserve Cash
    ↓
Step 2: Process Payment (gateway)
    ↓
Step 3: Generate Journal Entry
    ↓
Step 4: Sync to ERP
```

Jika Step 4 gagal:
- **Bukan rollback database** – semua sudah dicommit!
- **Compensating Action** diperlukan:
  ```
  Step 4 GAGAL
    ↓
  Compensate: Refund Cash
  Step 3 GAGAL → Compensate: Reverse Journal
  Step 2 GAGAL → Compensate: Cancel Payment
  ```

### Contoh Implementasi
```go
saga := NewSaga().
    Then(ReserveStep).
    Then(ProcessStep).
    Then(JournalStep).
    Then(ERPStep)

err := saga.Execute(ctx)
if err != nil {
    // Compensation otomatis dijalankan
}
```

> **Poin Penting**: Saga menggunakan **compensating actions** untuk mengembalikan keadaan, bukan `ROLLBACK` lintas service.

---

## 11. Choreography vs Orchestration (§13)

### Saga Choreography
```
InvoiceCreated
     ↓
  Inventory listens → InventoryUpdated
     ↓
  Commission listens → CommissionUpdated
```

**Kelebihan**:
- Loose coupling – setiap consumer independen
- Tidak ada titik tunggal kegagalan

**Risiko**:
- Flow sulit dilacak jika terlalu rumit (debugging)
- Tidak ada view global atas proses

### Saga Orchestration
```
PaymentSaga (orchestrator)
  ├── inventory.Process()
  ├── commission.Process()
  ├── accounting.Process()
  └── erp.Sync()
```

**Kelebihan**:
- Flow eksplisit, mudah dipantau
- Kontrol timeout & retry terpusat

**Risiko**:
- Orchestrator menjadi single point of failure
- Penambahan langkah baru memerlukan perubahan di orchestrator

> **Praktik Terbaik**: Pilih choreographed saga untuk sistem yang kompleks dan loose coupling diinginkan. Gunakan simple orchestration (seperti contoh di atas) ketika flow kontrol diperlukan.

---

## 12. Menjalankan Semua Test

```bash
go test -v ./labs/03-database-transaction/... -count=1
```

### Semua Test
| Test | Bukti |
|------|-------|
| `TestUnsafeLocalTransaction` | Partial state corruption tanpa transaksi |
| `TestSafeLocalTransaction` | ACID rollback yang bersih |
| `TestDistributedTransactionExternalSideEffectLimitation` | DB rollback ≠ WhatsApp rollback |
| `TestHTTPCallInsideTransactionLatency` | Transaction bertahan 50ms karena HTTP latency |
| `TestDualWriteProblemEventLost` | Invoice 'paid' tapi event tidak pernah terkirim |
| `TestOutboxDispatcherWithRetry` | Retry 3 kali: 2 gagal, 1 sukses |
| `TestOutboxDuplicateDeliveryAtLeastOnce` | Event dikirim 2x → at-least-once delivery |
| `TestIdempotentConsumerDeduplication` | Duplicate event diproses hanya 1x |
| `TestDeadLetterQueue` | Event gagal masuk DLQ setelah max attempts |
| `TestSagaPaymentWithCompensatingAction` | Compensation dijalankan saat step gagal |

---

## Ringkasan Konsep

| Masalah | Solusi |
|---------|--------|
| Partial state | Local transaksi ACID |
| External side effect tidak dapat di-rollback | JANGAN panggil HTTP di dalam transaksi |
| Dual-write tidak atomic | Transactional Outbox |
| Duplicate delivery | Idempotent consumer dengan deduplication |
| Event gagal sampai max retry | Dead Letter Queue (DLQ) |
| Cascading failure di luar DB | Saga + Compensating Transactions |
| Konektivitas lemah antar service | Choreography |
| Kebutuhan kontrol proses | Orchestration |