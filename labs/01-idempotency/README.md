# Lab 01 — Idempotency pada API

> **Mental Model**: Jangan pernah menganggap satu user action menghasilkan tepat satu request. Request bisa dikirim berulang karena network timeout, retry, atau double-click.

---

## 1. Masalah

Anda klik tombol "Bayar" sekali. Gateway sukses, tapi response hilang di jaringan (timeout). Client melakukan *retry* otomatis. 

Hasil tanpa proteksi: **Dua kali charge untuk satu pesanan.**

---

## 2. Mengapa duplicate request terjadi

Penyebab umum:
- **Network Timeout**: Server proses request, response gagal sampai ke client.
- **Client Retry**: Client mengira request gagal, mengirim ulang payload yang sama.
- **Frontend Double-Click**: User menekan tombol berkali-kali dengan cepat.
- **Queue Redelivery**: Message broker mengirim ulang pesan karena acknowledgement gagal.

**Akar Masalah**: Server secara default melihat tiap request sebagai request baru yang independen.

---

## 3. Contoh double payment

```
POST /orders/order-123/pay
Body: {"amount": 500000}

T0: Client kirim request bayar 500.000
T1: Server proses & charge gateway sukses (pay_1)
T2: Response timeout di client
T3: Client retry request yang sama (pay_2 dibuat, total 1.000.000 ter-charge)
```

---

## 4. Apa itu idempotency?

**Idempotency** adalah sifat di mana sebuah operasi dapat dilakukan berulang kali tanpa mengubah hasil akhir setelah eksekusi pertama.

- `GET` secara inheren idempotent.
- `POST` **tidak** idempotent secara default, butuh implementasi tambahan.

---

## 5. Idempotency-Key

Header HTTP wajib untuk transaksi kritikal:
```http
Idempotency-Key: 01JXYZabc123def456ghi789
```
atau:
```text
Idempotency-Key: vendor-payment-opl-1029-01JXYZabc123
```

**Idempotency Key ≠ Business ID**

`Idempotency-Key` berbeda dengan business identifier seperti `order-123`, `invoice-123`, atau `vendor-payment-opl-1029`:

| Business ID | Idempotency Key |
|-------------|-----------------|
| Mengidentifikasi resource (order, invoice) | Mengidentifikasi satu logical operation |
| Sifatnya **statis** | Sifatnya **dinamis per request** |
| Tidak ada batas waktu | Bisa memiliki TTL |
| Dapat diubah/rekonsiliasi | Harus immutable untuk duration operasi |

### Mengapa Business ID tidak cukup?
- Resource dapat memiliki operation baru di masa depan (refund, correction, adjustment)
- Payment dapat dikoreksi atau direfund
- Workflow operasional dapat berubah
- Partial payment mungkin ditambahkan kemudian

Gunakan format yang mencakup UUID sebagai operation identifier, bukan hanya resource identifier.

---

## 6. Satu key = satu logical operation

Satu `Idempotency-Key` mewakili tepat **satu logical operation**. 

Contoh Kasus Pembayaran Invoice #123:
```text
Request pertama (gagal karena timeout):
Idempotency-Key: 550e8400-e29b-41d4-a716-446655440000

Retry (harus menggunakan key yang sama persis):
Idempotency-Key: 550e8400-e29b-41d4-a716-446655440000
```

### ❌ Kesalahan Fatal yang Sering Terjadi:
```text
Request pertama:
Idempotency-Key: abc

Retry (SALAH - mengganti key saat retry):
Idempotency-Key: xyz
```
Walaupun payload dan tujuannya sama persis, backend akan menganggap `xyz` sebagai operasi baru yang berbeda dari `abc`, sehingga **double charge tetap terjadi**.

### Rule Inti:
```text
Same logical operation  → Same idempotency key
New logical operation   → New idempotency key
```
*Catatan: Gunakan UUIDv4 atau ULID yang di-generate di client sebagai key yang sehat. Tidak perlu membuat subsystem kompleks untuk key generation.*

---

## 7. Retry harus reuse key

```go
// Benar saat retry: gunakan key yang sama dari state sebelumnya
key := cachedKeyOrGenerate()
client.Post("/pay", WithHeader("Idempotency-Key", key))
```

---

## 8. Unsafe implementation

Implementasi tanpa idempotency key:
```go
func ProcessPayment(req PaymentRequest) (PaymentResult, error) {
    // Langsung tembak gateway tanpa cek duplikasi
    return gateway.Charge(req) 
}
```
Setiap request diproses sebagai transaksi baru.

---

## 9. Safe implementation

Flow lengkap pemrosesan dengan idempotency key:

```
Request masuk
      ↓
Cek Idempotency-Key
      ↓
┌──────────────────────────────┐
│ key belum ada                │
│ → proses request             │
├──────────────────────────────┤
│ key ada + completed          │
│ → replay response            │
├──────────────────────────────┤
│ key ada + masih processing   │
│ → return 409/503 (tunggu)  │
├──────────────────────────────┤
│ key sama + payload berbeda   │
│ → 409 Conflict               │
└──────────────────────────────┘
```

**Catatan**: Request duplikat dapat datang sebelum request pertama selesai (*race condition*). Mekanisme "still processing" mencegah proses ganda namun tidak perlu membuat lease, worker ownership, atau distributed lock untuk Lab dasar.

---

## 10. Race condition sederhana

Dua request identik datang secara bersamaan (hampir paralel):

```
Request A → cek key → tidak ada
Request B → cek key → tidak ada

Request A → create payment
Request B → create payment
```

**Hasil: Duplicate payment**

Application-level check saja tidak cukup. Database harus memiliki **unique constraint** untuk race condition ini.

> **Catatan**: Topik concurrency mendalam (mutex, pessimistic/optimistic locking, distributed lock, isolation level, deadlock) dibahas lebih dalam pada **Lab 02**.

---

## 11. Unique constraint

Solusi database untuk race condition:
```sql
CREATE UNIQUE INDEX idx_idempotency_key ON payments(idempotency_key);
```
Hanya satu request yang berhasil melakukan *insert* pertama, sisanya mendeteksi duplikat atau *conflict*.

---

## 12. Payload fingerprint

`Idempotency-Key` yang sama **tidak boleh** digunakan untuk request dengan intent berbeda:

```json
// Request pertama
{ "amount": 500000 }
```

kemudian:

```json
// Request kedua dengan key yang sama → HARUS ditolak!
{ "amount": 800000 }
```

**Implementasi sederhana:**

```go
// HTTP method + normalized payload → SHA-256
func hashRequest(method, path string, body []byte) string {
    normalized := fmt.Sprintf("%s %s %s", method, path, string(body))
    return sha256.Sum256([]byte(normalized))
}
```

**Aturan:**
- `same key + same payload` → accepted/replayed
- `same key + different payload` → 409 Conflict

**Penting**: Side effect (gateway charge) **tidak boleh dijalankan** untuk request yang konflik payload.

---

## 13. Request masih processing

Jika request pertama masih berjalan di gateway dan request second masuk dengan key yang sama:
- Jangan jalankan ulang.
- Kembalikan status bahwa request sedang diproses atau tunggu hasilnya.

---

## 14. Response replay

Jika request dengan key tersebut sudah sukses sebelumnya:
- Ambil response tersimpan dari database.
- Kirim kembali response asli ke client tanpa memanggil gateway lagi.

---

## 15. Idempotency vs database transaction vs unique constraint

| Konkrepansi | Tujuan | Contoh |
|-------------|--------|--------|
| **Idempotency** | Mencegah logical request yang sama dieksekusi berulang | Request `POST /pay` dengan key yang sama hanya diproses sekali |
| **Database Transaction** | Memastikan beberapa perubahan database commit atau rollback bersama | `BEGIN → create payment + update OPL + create cash out → COMMIT` |
| **Unique Constraint** | Database invariant menjadi protection terakhir terhadap duplicate identifier | `UNIQUE(idempotency_key)` mencegah duplikasi key |

### Perbedaan kunci

**Transaction tanpa idempotency** masih memungkinkan:
```text
Transaction A (create payment #1) — commit
Transaction B (create payment #2) — commit  ← duplicate request yang dianggap valid
```

**Transaction + Idempotency**:
```text
Request retry → key ada di cache → return replay response (tidak buat transaksi baru)
```

**Transaction + Unique Constraint**:
```text
Request A → insert → success
Request B → insert → UNIQUE violation → error 409
```

---

## 16. Scope idempotency key

`scope + key` menjadi identity dari idempotent operation.

Contoh scope:
```text
tenant_id          → multi-tenant isolation
user_id            → per-user operation
API operation      → vendor-payment:create, vendor-payment:refund
endpoint           → /payments, /refunds
```

Contoh:
```text
scope = vendor-payment:create
key   = abc123
→ identity: "Pembayaran vendor dengan key abc123"
```

Contoh SQL scope + key:
```sql
CREATE UNIQUE INDEX idx_idempotency_scoped 
ON payments(scope, idempotency_key);
```

**Basic example** untuk satu-tenant sederhana:
```sql
UNIQUE(idempotency_key)
```

> Produksi membutuhkan scope untuk isolation yang lebih baik.

---

## 17. TTL (Time To Live)

Idempotency key tidak perlu disimpan selamanya. Nilai TTL merupakan **keputusan bisnis dan operasional**:

- **Payment**: Sesuai kebutuhan audit (umumnya 72+ jam)
- **Create invoice**: 24–72 jam
- **Upload**: Beberapa jam
- **Generic request**: Sesuai window retry klien (misal 1 jam)

**Trade-off TTL:**
- **Terlalu pendek**: Retry dari klien setelah TTL lewat akan dianggap sebagai operasi baru (*duplicate operation*).
- **Terlalu panjang**: Kapasitas penyimpanan database membengkak seiring waktu.

---

## 18. HTTP Method dan Idempotency

| Method | Idempotent? | Keterangan |
|--------|-------------|------------|
| **GET** | ✅ | Mengembalikan state server tanpa mengubahnya |
| **PUT** | ✅ | `PUT /users/123` dengan body yang sama hasilnya tetap |
| **DELETE** | ✅ | `DELETE /files/123` kali pertama → dihapus, kali berikutnya → 404 |
| **POST** | ❌ | Non-idempotent secara default, tiap pemanggilan berpotensi menghasilkan state baru |

**POST /payments** biasanya memiliki side effect (charge gateway). Untuk membuatnya idempotent, application menambahkan:
```http
Idempotency-Key: 01JXYZabc123
```

**Idempotency** dalam HTTP berarti: efek akhir terhadap intended server state sama ketika request yang sama dijalankan berulang kali.

Contoh:
```http
POST /payments
Idempotency-Key: pay-2024-01-15-abc123
Content-Type: application/json
{ "amount": 500000 }
```
Request dengan key yang sama → hasil selalu sama (replay response), tidak ada charge ganda.

---

## 19. Batas Idempotency di Sistem Eksternal

Perhatikan skenario berikut:

```text
Backend menerima payment request
        ↓
Backend memanggil payment gateway
        ↓
Gateway berhasil
        ↓
Backend crash sebelum menyimpan hasil idempotency
        ↓
Client melakukan retry
```

> **Warning**: Idempotency record di database aplikasi tidak otomatis menjamin external side effect hanya terjadi satu kali jika kegagalan terjadi di celah antara pemanggilan pihak ketiga dan commit lokal.

Sistem production biasanya menggunakan idempotency mechanism atau reconciliation yang juga disediakan pada boundary sistem eksternal.

---

## 20. Kesalahan umum

1. **Hanya disable button frontend**: `button.disabled = true;` mencegah user double-click, tapi ini **UX protection, bukan server correctness**.
2. **Generate key baru setiap retry**: Mengubah key saat retry membuat server melihatnya sebagai operasi baru. **Retry harus menggunakan key logical operation yang sama.**
3. **Timestamp sederhana sebagai key**: Waktu tidak cukup unik. Potensi collision besar dan tidak konsisten saat retry.
4. **Hanya melakukan application check**: `if !exists(key) { insert() }` rentan terhadap race condition. Database uniqueness wajib.
5. **Key sama untuk payload berbeda**: Diterima sebagai update? **Salah, harus 409 Conflict**.
6. **Tidak menentukan scope**: Menyebabkan collision lintas operasi/user/tenant.
7. **TTL terlalu pendek**: Request lambat atau retry lama setelah TTL lewat dapat dieksekusi lagi sebagai operasi baru.

---

## 21. Kapan menggunakan idempotency?

- Pembayaran / penagihan (Payment Gateway).
- Pengiriman email/notifikasi penting.
- Pembuatan resource berharga (Order, Invoice).
- Webhook processing.

---

## 22. Senior engineer mindset

Junior:
> "Endpoint berhasil ketika saya test sekali."

Senior:
> "Apa yang terjadi jika request dikirim ulang, timeout, atau dua request yang sama masuk hampir bersamaan?"

Senior engineer mempertimbangkan:
- Duplicate request
- Retry
- Timeout
- Concurrent duplicate
- Payload mismatch
- Partial database failure

**Mental habit:** Selalu tanyakan apakah sebuah side effect aman ketika operation diulang.

---

## 23. Latihan Teknis

Jalankan test yang tersedia:
```bash
# Uji kode unsafe (membuktikan terjadinya double charge)
go test ./labs/01-idempotency/unsafe/... -v -count=1

# Uji kode safe (membuktikan idempotency mencegah duplikasi)
go test ./labs/01-idempotency/safe/... -v -count=1
```

---

## 24. Latihan OPL Pembayaran

Bayangkan skenario perpembayaran OPL:

```text
1. Admin memilih pekerjaan OPL
2. Admin memasukkan pembayaran
3. Sistem membuat bukti bayar
4. Sistem mengubah OPL menjadi paid
5. Sistem membuat cash out
```

**Pertanyaan:**

1. Apa yang terjadi jika tombol bayar ditekan dua kali?

2. Apa yang terjadi jika request pertama berhasil tetapi response timeout dan frontend retry?

3. Mengapa transaction saja tidak cukup?

4. Idempotency key mana yang harus digunakan?

5. Apa yang terjadi jika key sama digunakan dengan amount berbeda?

6. Constraint database apa yang membantu?

7. Jika request kedua masuk ketika request pertama masih processing, apakah payment boleh dijalankan lagi?

8. Jika payment gateway merupakan sistem eksternal, apakah local idempotency otomatis menjamin gateway tidak memproses dua kali?

**Jawaban singkat:**

1. Tanpa idempotency: dua bukti bayar, dua OPL menjadi paid, dua cash out.
2. Dengan idempotency: retry mengembalikan hasil pertama (replay response).
3. Transaction tidak cegah client mengirim dua request POST secara bersamaan.
4. Key = UUID unik per logical operation pembayaran OPL.
5. 409 Conflict – payload fingerprint harus konsisten.
6. `UNIQUE(idempotency_key)` pada tabel payments.
7. Tidak boleh – return 409 "request in progress" atau tunggu hasil.
8. Tidak otomatis. Backend harus callback idempotency key ke gateway.

---

## 25. Inti Pelajaran

> Jangan pernah percaya bahwa satu aksi UI menghasilkan tepat satu request ke server. Rancang API yang aman terhadap *retry storm*.

---

## 26. Deduplication vs Idempotency

| Konsep | Definisi | Kapan Digunakan |
|--------|----------|-----------------|
| **Idempotency** | Request yang diulang hasilnya sama | Request HTTP retry, payment API |
| **Deduplication** | Mengenali dan mengabaikan message/request yang sudah pernah diproses | Message queue, streaming data, batch processing |

Contoh deduplication di sistem messaging:
```text
event_id sudah pernah diproses?
→ ignore / return already processed
→ bukan → proses lagi
```

> **Catatan**: Deduplication pada messaging dapat menjadi materi lab tersendiri mencakup inbox pattern, Kafka transaction, atau exactly-once processing semantics.
