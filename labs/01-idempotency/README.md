# Lab 01 — Idempotency pada API

> **Mental Model**: Jangan pernah menganggap satu user action menghasilkan tepat satu request. Request bisa dikirim berulang karena network timeout, retry, atau double-click.

---

## 1. Masalah

Anda klik tombol "Bayar" sekali. Gateway sukses, tapi response hilang di jaringan (timeout). Client melakukan *retry* otomatis. 

Hasil tanpa proteksi: **Dua kali charge untuk satu pesanan.**

---

## 2. Mengapa duplicate request terjadi

Penyebab umum:
- Koneksi lambat / timeout
- User double-click
- Frontend retry
- HTTP client retry
- Queue redelivery (message broker)
- Payment gateway mengirim webhook berulang
- Server berhasil memproses tetapi response hilang

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

Salah satu pola umum untuk membuat mutation API aman terhadap retry adalah menggunakan `Idempotency-Key`:
```http
Idempotency-Key: 01JXYZabc123def456ghi789
```
atau:
```text
Idempotency-Key: vendor-payment-opl-1029-01JXYZabc123
```

| Business ID | Idempotency Key |
|-------------|-----------------|
| Mengidentifikasi resource (order, invoice) | Mengidentifikasi satu logical operation |
| Sifatnya **statis** | Unik per logical operation |
| Tidak ada batas waktu | Bisa memiliki TTL |
| Dapat diubah/rekonsiliasi | Harus immutable selama operation |

> Idempotency key harus unik untuk satu logical operation dan tetap sama selama seluruh retry operation tersebut.

### Mengapa Business ID tidak cukup?
- Resource dapat memiliki operation baru di masa depan (refund, correction, adjustment)
- Payment dapat dikoreksi atau direfund
- Workflow operasional dapat berubah

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

### Rule:
```text
same logical operation → same idempotency key
new logical operation → new idempotency key
```

---

## 7. Retry harus reuse key

```go
// Benar saat retry: gunakan key yang sama
key := cachedKeyOrGenerate()
client.Post("/pay", WithHeader("Idempotency-Key", key))
```

---

## 8. Unsafe implementation

```go
func ProcessPayment(req PaymentRequest) (PaymentResult, error) {
    return gateway.Charge(req) // tidak ada proteksi duplikat
}
```

Setiap request diproses sebagai transaksi baru.

---

## 9. Safe implementation

Flow pemrosesan dengan idempotency key:

```
Request masuk
      ↓
Validasi Idempotency-Key
      ↓
Hitung request fingerprint
      ↓
┌──────────────────────────────┐
│ key belum ada                │
│ → proses request             │
│ → simpan hasil               │
├──────────────────────────────┤
│ key ada + completed          │
│ → replay response            │
├──────────────────────────────┤
│ key ada + processing         │
│ → return 409 Conflict        │
├──────────────────────────────┤
│ key sama + payload berbeda   │
│ → 409 Conflict               │
└──────────────────────────────┘
```

---

## 10. Race condition sederhana

```
Request A → cek key → tidak ada
Request B → cek key → tidak ada

Request A → create payment → success
Request B → create payment → ??? (race condition)
```

Application-level check saja tidak cukup. Database harus memiliki **unique constraint**.

> **Catatan**: Topik concurrency mendalam dibahas lebih dalam pada **Lab 11**.

---

## 11. Unique constraint

```sql
CREATE UNIQUE INDEX idx_idempotency_key ON payments(idempotency_key);
```

Hanya satu request yang berhasil melakukan *insert* pertama.

---

## 12. Payload fingerprint

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
import (
    "crypto/sha256"
    "encoding/hex"
)

// HTTP method + normalized payload → SHA-256
func hashRequest(data []byte) string {
    sum := sha256.Sum256(data)
    return hex.EncodeToString(sum[:])
}
```

**Aturan:**
- `same key + same payload` → accepted/replayed
- `same key + different payload` → 409 Conflict

---

## 13. Request masih processing

Jika request ada status `processing`, request kedua harus return 409 Conflict.

---

## 14. Response replay

Jika key ada status `completed`:
- Ambil response tersimpan dari database.
- Kirim kembali response asli ke client tanpa memanggil gateway lagi.

---

## 15. Idempotency vs database transaction vs unique constraint

| Keterkaitan | Tujuan | Contoh |
|-------------|--------|--------|
| **Idempotency** | Mencegah logical request yang sama dieksekusi berulang | Request `POST /pay` dengan key yang sama hanya diproses sekali |
| **Database Transaction** | Memastikan beberapa perubahan database commit atau rollback bersama | `BEGIN → create payment → COMMIT` |
| **Unique Constraint** | Database invariant menjadi protection terakhir | `UNIQUE(idempotency_key)` mencegah duplikasi |

---

## 16. Scope idempotency key

`scope + key` menjadi identity dari idempotent operation.

```text
UNIQUE(scope, idempotency_key)
```

---

## 17. TTL (Time To Live)

Idempotency key tidak perlu disimpan selamanya.

**Trade-off TTL:**
- **Terlalu pendek**: Retry lama dapat dianggap operasi baru.
- **Terlalu panjang**: Kapasitas penyimpanan membengkak.

---

## 18. HTTP Method dan Idempotency

| Method | Idempotent? |
|--------|-------------|
| **GET** | ✅ |
| **PUT** | ✅ |
| **DELETE** | ✅ |
| **POST** | ❌ |

`POST /payments` membutuhkan Idempotency-Key untuk mencegah duplicate charge.

---

## 19. Batas idempotency di sistem eksternal

Idempotency record di database aplikasi **tidak otomatis** menjamin external side effect hanya terjadi satu kali jika kegagalan terjadi di antara pemanggilan pihak ketiga dan commit lokal.

Jika payment provider mendukung provider-side idempotency, gunakan stable idempotency key pada request ke provider. Jika tidak, sistem biasanya membutuhkan mekanisme duplicate prevention atau reconciliation lain.

---

## 20. Deduplication vs Idempotency

| Konsep | Definisi | Kapan Digunakan |
|--------|----------|-----------------|
| **Idempotency** | Request yang diulang hasilnya sama | Request HTTP retry, payment API |
| **Deduplication** | Mengenali dan mengabaikan message yang sudah pernah diproses | Message queue, streaming |

> **Catatan**: Deduplication pada messaging dapat menjadi materi lab tersendiri.

---

## 21. Kesalahan umum

1. **Hanya disable button frontend**: UX protection, bukan server correctness.
2. **Generate key baru setiap retry**: Harus gunakan key logical operation yang sama.
3. **Timestamp sederhana sebagai key**: Tidak unik, berisiko collision.
4. **Hanya application check**: Rentan race condition. Database uniqueness wajib.
5. **Key sama untuk payload berbeda**: Harus 409 Conflict.
6. **Tidak menentukan scope**: Bisa collision lintas user/tenant.
7. **TTL terlalu pendek**: Retry lama dianggap operasi baru.

---

## 22. Kapan menggunakan idempotency?

- Pembayaran / pembayaran vendor
- Top-up saldo
- Pembuatan invoice / purchase order
- Pencairan komisi
- Pengurangan stok
- Pembuatan booking
- Webhook payment gateway
- Pengiriman email/WhatsApp

Idempotency sangat penting ketika **duplicate execution dapat menghasilkan side effect yang merugikan, mahal, atau sulit dibatalkan**.

---

## 23. Senior engineer mindset

Junior:
> "Endpoint berhasil ketika saya test sekali."

Senior:
> "Apa yang terjadi jika request dikirim ulang, timeout, atau dua request yang sama masuk hampir bersamaan?"

Senior engineer mempertimbangkan:
- Duplicate request, Retry, Timeout
- Concurrent duplicate, Payload mismatch
- Partial database failure

**Mental habit**: Selalu tanyakan apakah side effect aman ketika operation diulang.

---

## 24. Implementasi In-Memory (Simulasi)

Implementasi safe pada lab menggunakan in-memory store (`map` + `sync.RWMutex`) agar fokus tetap pada konsep. Atomic uniqueness pada contoh ini mensimulasikan protection yang pada sistem production biasanya diberikan oleh database unique constraint.

---

## 25. Latihan Teknis

Jalankan test:
```bash
# Uji kode unsafe (membuktikan terjadinya double charge)
go test ./labs/01-idempotency/unsafe/... -v -count=1

# Uji kode safe (membuktikan idempotency mencegah duplikasi)
go test ./labs/01-idempotency/safe/... -v -count=1
```

---

## 26. Latihan OPL Pembayaran

Bayangkan skenario pembayaran OPL:

1. Apa yang terjadi jika tombol bayar ditekan dua kali?
2. Apa yang terjadi jika request timeout dan client retry?
3. Mengapa transaction saja tidak cukup?
4. Idempotency key mana yang harus digunakan?
5. Apa yang terjadi jika key sama digunakan dengan amount berbeda?
6. Constraint database apa yang membantu?
7. Jika request kedua masuk ketika request pertama masih processing, apakah payment boleh dijalankan lagi?
8. Jika payment gateway merupakan sistem eksternal, apakah local idempotency otomatis menjamin gateway tidak memproses dua kali?

**Jawaban:**
1. Tanpa idempotency: dua bukti bayar, dua OPL menjadi paid, dua cash out.
2. Dengan idempotency: retry mengembalikan hasil pertama (replay response).
3. Transaction tidak cegah client mengirim dua request POST secara bersamaan.
4. Key = UUID unik per logical operation pembayaran OPL.
5. 409 Conflict – payload fingerprint harus konsisten.
6. `UNIQUE(idempotency_key)` pada tabel payments.
7. Tidak boleh – return 409 "request in progress" atau tunggu hasil.
8. Tidak otomatis. Backend harus memakai provider-side idempotency atau reconciliation engine.

---

## 27. Inti Pelajaran

> Jangan pernah percaya bahwa satu aksi UI menghasilkan tepat satu request ke server. Rancang API yang aman terhadap duplicate request dan retry.
