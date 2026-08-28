# Lab 03 — Distributed Transaction: Kenapa Database Transaction Saja Tidak Cukup?

> **Mental Model**: Database transactions (ACID) guarantee atomicity **hanya dalam satu transactional resource (satu database)**. Ketika workflow melibatkan sistem eksternal (API, webhook, messaging broker, notifier), database transaction tidak dapat mengelilingi semua perubahan secara atomik.

---

## 1. Fondasi: Local Database Transaction

Bayangkan skenario pembayaran:
1. Buat pembayaran di tabel `payments`
2. Update status order di tabel `orders`
3. Buat transaksi di dompet digital `wallet_transactions`

### A. Unsafe Local Transaction (Partial State)

Tanpa transaksi, jika error terjadi di tengah (misal setelah `INSERT payment` tapi sebelum `INSERT wallet_transaction`), sistem meninggalkan **partial state** yang tidak konsisten.

### B. Safe Local Transaction (ACID Rollback)

Dengan `BEGIN`, `COMMIT`, `ROLLBACK` (atau helper `WithTx`), seluruh perubahan dalam satu transaksi dibatalkan bersama jika ada error.

> **Kesimpulan**: Database transaction efektif **jika dan hanya jika** seluruh perubahan berada dalam transactional resource yang sama.

---

## 2. Masalah Distributed Transaction: External Side Effect

Di dunia nyata, transaksi sering melibatkan:
- Payment gateway (HTTP API)
- Notifikasi WhatsApp/Email
- Webhook ke sistem ERP
- Publishing event ke message broker

**Masalah**: `DB Rollback` tidak dapat "membatalkan" efek samping yang sudah terjadi di luar database.

---

## 3. Anti-Pattern: HTTP Call Di Dalam Transaction

### Contoh Flow
```
BEGIN
  ├── UPDATE invoice SET status = 'paid'
  ├── HTTP call ke payment gateway (lambat! 50ms-5detik)
  └── COMMIT
```

### Risiko
- **Transaction duration menjadi lama** → Connection pool tertahan
- **Row lock bertahan** → Contention & deadlock
- **Downstream latency ikut mempengaruhi DB** → Timeout
- **Throughput menurun** → Semua transaksi menunggu

### Bukti (Test `TestHTTPCallInsideTransactionLatency`)
- Transaksi bertahan selama 50ms-100ms karena `http.Ping()` di dalam transaksi.

---

## 4. Masalah Dual-Write: DB Write + Message Publish

### Contoh Flow
```
BEGIN
  ├── UPDATE invoice SET status = 'paid'
  └── COMMIT (invoice = 'paid' di database)
  ├── Publish 'InvoicePaid' event ke message broker
  └── SIMULASI CRASH → Event tidak pernah dikirim!
```

### Hasil
- **Database**: `invoice = paid`
- **Event**: **Tidak pernah terkirim** → Sistem downstream tidak tahu invoice berubah

### Bukti (Test `TestDualWriteProblemEventLost`)
- Invoice berstatus 'paid' tetapi tidak ada event yang dipublikasikan ke broker.

---

## 5. Solusi: Transactional Outbox Pattern

### Prinsip
Simpan "intent to publish" sebagai baris di database dalam **sitransaksi yang sama** dengan perubahan business state.

### Schema Tabel Outbox
```sql
CREATE TABLE outbox_events (
    id            UUID PRIMARY KEY,
    event_type    TEXT,
    aggregate_id  BIGINT,
    payload       JSON,
    status        TEXT CHECK (status IN ('pending', 'published')),
    attempts      INT DEFAULT 0,
    created_at    TIMESTAMPTZ,
    published_at  TIMESTAMPTZ NULL
);
```

### Flow dengan Outbox Pattern
```
BEGIN TRANSACTION
  ├── UPDATE invoice SET status = 'paid' WHERE order_id = 101
  ├── INSERT INTO outbox_events (event_type='InvoicePaid', payload=..., status='pending')
  └── COMMIT
```

Jika commit berhasil → **KEDUA** berada di database.
Jika transaction gagal → **KEDUA** tidak tersimpan.

### Bukti (Test `TestTransactionalOutboxPatternAtomicity`)
- Invoice berubah menjadi 'paid'
- Satu outbox event dengan status 'pending' disimpan secara atomik

---

## 6. Ringkasan Perbandingan

| Aspek | Dual-Write Bug | Transactional Outbox |
|-------|---------------|---------------------|
| Status DB | `paid` | `paid` |
| Event | **Tidak ada** (hilang!) | **Ada satu** (status pending) |
| Konsistensi | tidak terjamin | terjamin |
| Recovery | manual cleanup | worker retry otomatis |
| Cara kerja | DB → publish berurutan | DB → outbox di satu transaksi |

---

## Menjalankan Semua Test

```bash
go test -v ./labs/03-database-transaction/... -count=1
```

### Hasil Test
1. **TestUnsafeLocalTransaction**: Bukti partial state corruption tanpa transaksi.
2. **TestSafeLocalTransaction**: Bukti ACID rollback yang bersih.
3. **TestDistributedTransactionExternalSideEffectLimitation**: Bukti bahwa `DB rollback ≠ external rollback` (WhatsApp tetap terkirim).
4. **TestHTTPCallInsideTransactionLatency**: Bukti transaction bertahan lama karena HTTP call.
5. **TestDualWriteProblemEventLost**: Bukti event hilang setelah DB commit.
6. **TestTransactionalOutboxPatternAtomicity**: Bukti business state + outbox disimpan atomik.
7. **TestTransactionalOutboxRollback**: Bukti konsistensi terjaga.

---

## Next Steps

Untuk mengimplementasikan worker yang membaca outbox dan mengirimkan event ke broker, lihat **Lab 07 – Outbox Pattern**.