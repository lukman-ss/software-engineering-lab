# Lab 09 — Code Review: Cara Senior Engineer Menemukan Bug Sebelum Masuk Production

## Tujuan

Melatih kemampuan **code review** berbasis risiko dengan memanalisis implementasi checkout yang memiliki multiple masalah serta membandingkan dengan versi yang sudah direview.

Senior Engineer biasanya tidak hanya membaca kode secara linear. Mereka memprioritaskan area dengan risiko production tertinggi terlebih dahulu.

Pada banyak tim engineering, code review menjadi bagian signifikan dari pekerjaan seorang Senior Engineer. Mereka memfokuskan review pada:

- **Business Logic Correctness**
- **Concurrency & Race Conditions**
- **Data Integrity**
- **Error Handling**
- **Security**
- **Performance**
- **Testing Adequacy**

## Pendekatan Code Review

### Alur Analisis Masalah
Dalam code review, jangan langsung memilih pattern atau teknologi. Urutannya:

1. Code
2. Failure mode
3. Impact
4. Perbaikan paling sederhana yang cukup

Contoh: Jangan langsung berkata "Gunakan distributed lock." Mulai dari: "Dua checkout bersamaan dapat membaca stok yang sama dan menyebabkan overselling." Baru setelah risiko jelas, tentukan mekanisme paling sederhana yang menyelesaikannya.

### Review Perubahan, Bukan Hanya File

Reviewer perlu memahami tujuan PR sebelum menilai implementasinya:

1. Requirement
2. Diff
3. Behavior yang berubah
4. Risiko

Pertanyaan Senior Reviewer:
> *"Apakah perubahan ini benar-benar menyelesaikan masalah yang diminta tanpa mengubah behavior lain yang tidak diperlukan?"*

Code yang secara teknis benar tetap bisa salah jika tidak sesuai requirement.

- Perubahan CSS/HTML → low risk
- Perubahan Payment/Inventory → high risk
- Jumlah baris kode ≠ tingkat risiko.

## Business Case: Checkout Sistem

Berikut pseudo-code PHP yang menggambarkan masalah utama:

```php
public function checkout(Request $request)
{
    $cart = Cart::where('user_id', auth()->id())->get();

    foreach ($cart as $item) {
        $product = Product::find($item->product_id); // BUG: N+1 query lookup
        $product->stock -= $item->qty;               // BUG: non-atomic, race condition
        $product->save();                            // BUG: multiple individual UPDATEs
    }

    Order::create(...);  // BUG: no transaction, partial state on failure

    return response()->json(['success' => true]);  // BUG: no idempotency, error swallowed
}
```

> **Catatan Terminology:** Masalah N+1 query klasik terjadi pada baris `Product::find(...)` di dalam loop (melakukan query `SELECT` satu per satu untuk setiap item cart). Sementara baris `$product->save()` menghasilkan banyak eksekusi individual `UPDATE`, bukan contoh N+1 relation lookup, namun tetap tidak efisien dibanding bulk/batch operations.

Masalah di atas terdapat implementasi Go di folder `checkout_naive.go` dan `checkout_improved.go`.

## Temukan Masalahnya

Sebelum melihat implementasi, **cobalah temukan minimal 10 issue** pada kode berikut:

### CheckoutNaive (Intentionally Broken)

<details>
<summary>Klik untuk melihat kode naive</summary>

```go
func (c *CheckoutNaive) Checkout(ctx context.Context, userID string) (*CheckoutResponse, error) {
    cartItems, _ := c.cartSource.GetCart(ctx, userID)

    for _, item := range cartItems {
        product, _ := c.products.GetProduct(ctx, item.ProductID)  // BUG: N+1 query
        product.Stock -= item.Quantity                            // BUG: race condition
        _, _ = c.products.UpdateStock(ctx, item.ProductID, -item.Quantity)  // BUG: non-atomic
        total += int64(item.Quantity) * int64(item.UnitPrice)
    }

    order := &Order{...}
    _ = c.repo.CreateOrder(ctx, order)    // BUG: no transaction
    _ = c.notify.SendOrderConfirmation(ctx, userID, orderID)  // BUG: side effect

    return &CheckoutResponse{Success: true}, nil  // BUG: always succeeds
}
```

</details>

### Review Findings

| No | Issue | Severity |
|----|-------|----------|
| 1 | Stock dapat menjadi negatif (overselling) | BLOCKER |
| 2 | Race condition pada concurrent checkout | BLOCKER |
| 3 | Tidak ada transaction boundary | MAJOR |
| 4 | N+1 query pattern | MINOR |
| 5 | Partial state jika order creation gagal | MAJOR |
| 6 | Tidak ada idempotency / retry protection | MAJOR |
| 7 | Product tidak ditemukan tidak ditangani | MAJOR |
| 8 | Empty cart tidak divalidasi | MINOR |
| 9 | Error handling buruk (selalu return success) | MAJOR |
| 10 | Logging/context tidak memadai | NIT |
| 11 | Notification sebagai side effect berbahaya | MINOR |
| 12 | Duplicate code berlebihan pada validasi | NIT |
| 13 | Hardcoded environment value | MINOR |

## Severity Classification

> **Prinsip Utama:** `Severity != jumlah baris kode`.
> Severity ditentukan oleh: **impact × likelihood × blast radius**.
> - **3 baris** perubahan payment/stock logic → **Critical / Blocker**
> - **500 baris** perubahan styling halaman reporting internal → **Medium / Low**

| Level | Deskripsi | Contoh |
|-------|-----------|--------|
| BLOCKER | Risiko data korupsi, kehilangan data, atau security critical | Race condition overselling, double charge |
| MAJOR | Masalah bisa menyebabkan error di production | No transaction boundary, missing rollback |
| MINOR | Masalah performa atau user experience | N+1 query lookup, sub-optimal UI feedback |
| NIT | Formatting, naming yang tidak konsisten | Inconsistent naming |

## Business Logic Review

### [ ] Apakah business logic benar?

- Stock deduction harus atomic
- Validasi kuantitas sebelum proses
- Total harga harus akurat

**Mendeteksi Missing State Handling:**
Jika sebuah endpoint menandai status pesanan menjadi `PAID` atau `COMPLETED`, reviewer harus bertanya: *"Bagaimana kalau statusnya sudah PAID?"* Code harus melindungi state tersebut. Response HTTP-nya bergantung *contract* sistem (bisa 200 OK Idempotent, bisa 409 Conflict, dll). Yang penting, reviewer mendeteksi ketiadaan *state handling* tersebut.

### [ ] Apakah invariant tetap terjaga?

- Stock >= 0 (tidak boleh negatif)
- Order hanya dibuat jika semua validasi lolos

## Advanced Review Finding — Atomicity Gap Between Business Transaction and Idempotency Finalization

Saat proses multi-step, sering terjadi failure window antara penyelesaian transaksi database dengan persistensi record idempotency (misal ke Redis).

Flow:
1. `BEGIN TRANSACTION`
2. `ReserveStock()`
3. `CreateOrder()`
4. `COMMIT`
5. ↓ **CRASH/NETWORK ERROR**
6. `MarkCompleted(idempotency)` gagal

Jika langkah ke-6 gagal, *business transaction sudah sukses*, namun *idempotency record masih PROCESSING*. Jika client retry, server mungkin merespons dengan error `ErrDuplicateRequest` atau memproses ulang jika TTL expired. 

**Opsi Production:**
- Simpan record idempotency ke dalam satu *transactional boundary* yang sama dengan business transaction (seperti pola Outbox).
- Gunakan *durable reconciliation* via worker.
- Terapkan *unique business constraint* di database (misalnya `idempotency_key` menjadi `UNIQUE` column di table orders) untuk mencegah duplikasi sejati, mengizinkan recovery *stale* PROCESSING record secara aman.

Dalam lab ini, error ini didemonstrasikan dengan mengembalikan error khusus: `ErrIdempotencyFinalize`.

## Duplicate Code & YAGNI

Reviewer sering menganggap duplicate code otomatis harus di-*extract* menjadi method/service/helper. Ini bisa menjadi *premature abstraction*. Duplicate code adalah **signal untuk diperiksa**, bukan perintah otomatis membuat abstraction.

- Extract jika *behavior* memang sama dan kemungkinan berubah bersama.
- Jangan membuat abstraction hanya karena code kebetulan terlihat mirip secara struktur.
- Pilih solusi paling sederhana yang cukup.

Dalam hal **Maintainability**: Buat *intent* (tujuan) kode jelas melalui penamaan yang baik dan *readability*. Jangan sembunyikan setiap logika percabangan di dalam method baru jika itu hanya membuat developer harus melompat antar file. Jangan *abstract without need*.

Terapkan YAGNI secara ketat.

## Hardcode Configuration

Environment-specific value tidak boleh tersebar sebagai hardcoded value di application code. Arahkan via configuration, misalnya `config('services.payment.url')`, dan jadikan environment variable sebagai *sumber nilai* configuration.

## Concurrency Review

Reviewer harus mengenali *failure mode* terlebih dahulu. Jangan jadikan Transaction, Idempotency, dan Lock sebagai *checklist* yang semuanya harus dipakai sekaligus. Solusi bisa berupa:
- transaction,
- atomic update,
- row lock,
- idempotency,
- atau kombinasi yang memang diperlukan.

Fokuslah pada kemampuan menemukan risiko.

### [ ] Apakah concurrent request aman?

Run test: `go test -race ./...` untuk mendeteksi data race (concurrent memory access tanpa sinkronisasi yang benar).

**Penting:** Go `race` detector **tidak** secara otomatis mendeteksi logical/business race condition (seperti dua checkout sukses padahal stock cuma satu). Data race vs Logical Race Condition:
- **Data Race**: Dua goroutine read/write *memory address* yang sama di saat bersamaan. Terdeteksi oleh `go test -race`.
- **Logical Race Condition**: Interleaving urutan eksekusi merusak *business invariant* (contoh: Overselling). Tidak terdeteksi oleh data race detector, **harus** diuji menggunakan *concurrency/invariant test* khusus yang memastikan invarian bisnis dipertahankan.

**Naive:** Tidak aman - dapat overselling.
**Improved:** Aman - transaksi terisolasi dan atomic stock reservation menjaga invarian.

### [ ] Keterbatasan MockTransactionManager

Mock transaksi dalam lab ini menggunakan Global Mutex, yang secara efektif men-serialize eksekusi di test. Ini dilakukan untuk tujuan **demonstrasi yang deterministic**. 

**Di Production:** Database memiliki Concurrency Control (MVCC) untuk mencapai atomicity. Jangan menganggap bahwa "berada dalam DB Transaction" otomatis mencegah race condition. Untuk safe stock reservation di production:

```sql
UPDATE products
SET stock = stock - $1
WHERE id = $2
  AND stock >= $1;
```
Lalu program mengecek *affected rows*. Jika 0, maka stock tidak cukup. Pendekatan lain adalah dengan `SELECT ... FOR UPDATE` (Pessimistic Locking).

## Security Review

### Authentication vs Authorization

| Aspek | Deskripsi |
|-------|-----------|
| **Authentication** | Siapa yang login? (`Principal{UserID}`) |
| **Authorization** | Apakah user boleh checkout cart milik `CartOwnerID`? |

```go
// Benar: memakai Principal sebagai source of truth, lakukan ownership check
if principal.UserID != cmd.CartOwnerID {
    return nil, ErrForbidden
}
cartItems, _ := c.cartSource.GetCart(ctx, cmd.CartOwnerID)
```

Implementasi memastikan:
- User tidak dapat checkout cart milik user lain (`TestImprovedCheckout_CannotCheckoutAnotherUsersCart`)
- Idempotency key *namespaced* per-user: `checkout:{userID}:{idempotencyKey}` sehingga dua user dengan key sama tidak bertabrakan (`TestImprovedCheckout_IdempotencyKeyIsScopedPerUser`)

## Performance Review

### N+1 Query Problem

Naive implementation melakukan `GetProduct` dalam loop:

```go
for _, item := range cartItems {
    product, _ := c.products.GetProduct(ctx, item.ProductID)  // N query
    // ...
}
```

**Improved:** Mengumpulkan unique product ID dan batch query `GetProducts` sekali saja:

```go
seen := make(map[string]struct{})
productIDs := make([]string, 0, len(cartItems))
for _, item := range cartItems {
    if _, ok := seen[item.ProductID]; !ok {
        seen[item.ProductID] = struct{}{}
        productIDs = append(productIDs, item.ProductID)
    }
}
productsMap, err := c.products.GetProducts(ctx, productIDs) // 1 batch query
```
Dibuktikan via call counter di test `TestNaiveCheckout_NPlusOneProductLookups` vs `TestImprovedCheckout_BatchLoadsProducts` (dan `TestImprovedCheckout_BatchLoadUsesUniqueProductIDs`).

## Test Coverage Review

Saat mereview, pertanyaan bukan:
"Apakah ada test?"

Tetapi:
"Apakah behavior yang berubah dilindungi test?"

Contoh minimal (pada checkout):
- checkout normal
- stok tidak cukup
- product tidak ditemukan
- order gagal dibuat

Tidak semua kemungkinan harus memiliki test. Prioritaskan behavior yang berubah dan failure path yang memiliki risiko.

## Error Handling & Logging

### [ ] Apakah error ditangani dengan baik?

- [ ] Empty cart
- [ ] Product tidak ditemukan
- [ ] Kuantitas tidak valid
- [ ] Stock tidak cukup

Reviewer harus membedakan kegagalan:
- **Transient** (misal: network timeout) → mungkin aman untuk retry/idempotency
- **Permanent** (misal: invalid input, stock habis) → fail fast, jangan retry

### [ ] Apakah logs punya context?

Logging harus terstruktur dan membantu investigasi, bukan sekadar noise.

```go
// BURUK: Sulit diinvestigasi
c.logger.Error(ctx, "checkout failed", "error", err)

// BAIK: Membawa context
c.logger.Error(ctx, "checkout transaction failed", 
    "userID", principal.UserID, 
    "idempotencyKey", scopedKey, 
    "error", err)
```

## Improved Implementation

Lihat `checkout_improved.go` untuk implementasi yang sudah direview dengan:

1. **Validasi input** - empty cart, quantity validation, product validation
2. **Atomic Idempotency Claim** - `Claim` dengan state (PROCESSING/COMPLETED), conflict check via hash, auto-release pada error
3. **Transaction Boundary & Rollback** - unit of work `TransactionManager` memastikan stock reservation dan `CreateOrder` atomic (rollback jika salah satu gagal)
4. **Conditional Atomic Stock Operation** - validasi dan pemotongan stock berada dalam critical section transaksional (invarian `stock >= 0` terjaga)
5. **Good error propagation** - mengembalikan error yang sesuai tanpa merusak state idempotency; infra error tidak disamakan dengan domain error
6. **Best-Effort Post-Commit Notification** - notification dijalankan secara synchronous di luar transaksi setelah commit sukses. Kegagalan kirim dicatat tanpa membatalkan order yang sudah committed.
   > **Production Failure Window:** Jika server crash tepat setelah commit namun sebelum notification dikirim, notifikasi akan hilang. Di sistem production berskala besar, gunakan **Transactional Outbox Pattern** (lihat Lab Distributed Transactions): event dimasukkan ke database dalam transaksi yang sama, lalu worker background mengirimkan notifikasi secara asinkron dan terjamin (*at-least-once*).
7. **Idempotency finalization failure di-handle** - `MarkCompleted` error dikembalikan sebagai `ErrIdempotencyFinalize` untuk indikasi atomicity gap.
8. **Idempotency payload hashing yang benar** - Hash dibuat dari request payload (CheckoutCommand), bukan dari mutable state seperti cart content.
9. **Validasi idempotency key** - Menolak key kosong atau whitespace.
10. **Release error handling** - `Release` error digabungkan dengan error asal via `errors.Join`

### Order ID Generation di Lab vs Production

Pada `CheckoutImproved`, kita menggunakan `atomic.Int64` semata-mata untuk **demonstrasi ID unik dalam satu proses/test suite lokal**.

> **Peringatan Production:** `atomic.Int64` **tidak aman** untuk sistem terdistribusi multi-replica:
> - `Instance A` akan membuat ID: `order-1`
> - `Instance B` juga akan membuat ID: `order-1` (terjadi ID collision)
> - Setiap kali aplikasi restart, counter kembali ke `0`.
> 
> Di production, selalu gunakan strategi distributed-safe ID:
> - **Database Sequence / Auto Increment** (jika single writer)
> - **UUID v4 / v7** atau **ULID** (lexicographically sortable)
> - **Snowflake / KSUID** untuk ID terdistribusi berperforma tinggi.

## Transaction Boundary

Unit of work `TransactionManager` menjamin atomicity pada multi-step checkout. Jika salah satu operasi gagal, seluruh perubahan state dalam unit of work di-rollback.

```go
txManager.WithinTransaction(ctx, func(tx CheckoutTx) error {
    for _, item := range cartItems {
        if err := tx.ReserveStock(ctx, item.ProductID, item.Quantity); err != nil {
            return err
        }
    }
    if err := tx.CreateOrder(ctx, order); err != nil {
        return err
    }
    return nil
})
```

## Atomic Stock Reservation

Check dan mutasi stock digabung dalam operasi atomik untuk menghindari race condition.

```go
// Mock meniru critical section yang menjamin invariant `stock >= 0`
func (r *MockProductRepository) ReserveStock(ctx context.Context, productID string, quantity int) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    stock, exists := r.stock[productID]
    if !exists {
        return ErrProductNotFound
    }
    if stock < quantity {
        return ErrInsufficientStock
    }
    r.stock[productID] = stock - quantity
    return nil
}
```

## Idempotency State Machine

State idempotency berkembang melalui state yang jelas:
- **PROCESSING**: request sedang dijalankan
- **COMPLETED**: request sukses, response tersimpan untuk replay

Pada error, record di-release agar retry memungkinkan.

**Catatan: Idempotency Terhadap Mutable State**
Sangat penting bahwa Hash Idempotency dan verifikasinya **hanya bergantung pada payload asli dari client**, bukan pada mutable state internal (seperti isi Cart yang diambil dari DB pada saat request).
Jika hash dibuat dari state Cart, dan kemudian client mengosongkan/mengubah Cart-nya sebelum melakukan retry (karena timeout), hash-nya akan berubah dan sistem keliru menganggapnya sebagai request berbeda (atau terjadi conflict), alih-alih mengembalikan cache response original yang sudah sukses! Dalam `ImprovedCheckout`, pengecekan dilakukan _sebelum_ memuat cart, dan yang di-hash adalah `CheckoutCommand` (client payload).

## Deterministic Concurrency Test

Test konkurensi (seperti `TestImprovedCheckout_ConcurrentSameIdempotencyKeyRunsOnce`) memakai coordinated concurrent start dengan `sync.WaitGroup` (`startGate`) untuk mensimulasikan concurrent requests diwaktu bersamaan. Sedangkan test yang mendemostrasikan race condition pada naive implementation (`TestNaiveCheckout_DeterministicallyDemonstratesOverselling`) memakai hook dengan read/write barrier exact interleaving (`sync.WaitGroup`) untuk memaksakan overselling secara deterministic.

## Rollback Demonstration

Test seperti `TestImprovedCheckout_OrderCreationFailureRollsBackStock` dan `TestImprovedCheckout_MultiProductOrderCreationFailureRollsBackStock` membuktikan bahwa kegagalan `CreateOrder` mengembalikan stock ke nilai semula tanpa menciptakan order.

## Authorization

```go
if principal.UserID != cmd.CartOwnerID {
    return nil, ErrForbidden
}
```

Dibuktikan via test `TestImprovedCheckout_CannotCheckoutAnotherUsersCart`.

## N+1 Demonstration

Naive tetap melakukan `GetProduct` sebanyak jumlah item. Improved menggunakan `GetProducts` batch.

Dibuktikan via call counter:
- `TestNaiveCheckout_NPlusOneProductLookups`
- `TestImprovedCheckout_BatchLoadsProducts`


## Failure Scenario Tests

### Test: Stock Tidak Cukup

Jika stock habis/tidak mencukupi, checkout harus gagal tanpa mengubah stock (`TestImprovedCheckout_InsufficientStock`).

### Test: Duplicate Request

Client mengirim request yang sama dengan payload identik. Mengembalikan response replay tanpa eksekusi transaksi ulang (`TestImprovedCheckout_DuplicateRequest`).

### Test: Concurrent Checkout

Dua request checkout berbeda atau sama yang dijalankan bersamaan:
- Coordinated concurrent start mencegah duplicate execution dengan key sama (`TestImprovedCheckout_ConcurrentSameIdempotencyKeyRunsOnce`).
- Atomic transaction mencegah overselling pada stock terbatas (`TestImprovedCheckout_NoRaceCondition`, `TestConcurrentCheckout_InvariantPreserved`).

### Test: Rollback on Order Creation Failure

Jika `CreateOrder` gagal saat transaksi:
- Stock yang di-reserve di-rollback kembali utuh.
- Tidak ada order yang tercatat.
- Notification tidak terkirim.
- Idempotency key di-release sehingga retry dapat diproses kembali (`TestImprovedCheckout_OrderCreationFailureRollsBackStock`, `TestImprovedCheckout_MultiProductOrderCreationFailureRollsBackStock`).

## Risk-Based Review & Blast Radius

Ukuran PR bukan ukuran risiko. Saat Senior Engineer mereview sebuah PR, pertanyaan kuncinya adalah:

**Jika bug ini lolos ke production:**
1. Berapa banyak user/transaksi yang terdampak?
2. Apakah menyebabkan kehilangan uang (*financial loss*)?
3. Apakah menyebabkan korupsi data (*data corruption*)?
4. Apakah kegagalannya *reversible* (bisa di-rollback/dikompensasi dengan mudah)?
5. Seberapa cepat kita bisa mendeteksi dampaknya di production (*observability*)?

Fokus review harus diprioritaskan pada area berisiko tinggi:

| Area | Risk Level | Blast Radius & Pertimbangan |
|------|------------|-----------------------------|
| CSS/HTML styling | Low | Visual glitch minor, reversible, tidak ada data loss |
| Internal reporting | Medium | Query lambat atau report salah, blast radius terbatas ke internal |
| Authentication / Authz | High | Kebocoran data user, impersonation |
| Payment / Inventory | **Critical** | Kehilangan uang, overselling, irreversible, dampak hukum |

**Contoh matrix risiko:**

| Risiko | Dampak | Probabilitas | Tindakan |
|--------|--------|--------------|----------|
| Stock negative | Critical | Low-Medium | Wajib verifikasi via concurrent stress test / atomicity |
| Duplicate order | High | Medium | Wajib idempotency key protection |
| Log error format | Low | High | Standardisasi logger format |

## Code Review Checklist

- [ ] Business logic benar?
- [ ] Invariant tetap terjaga (stock >= 0)?
- [ ] Authorization benar (hanya cart milik user)?
- [ ] Concurrent request aman?
- [ ] Retry aman (idempotency)?
- [ ] Transaction boundary benar?
- [ ] Query pattern efisien (tidak N+1)?
- [ ] Error ditangani dengan baik?
- [ ] Logs punya context?
- [ ] Failure state konsisten?
- [ ] Test mengcover happy path?
- [ ] Test mengcover failure path?
- [ ] Test mengcover concurrency?
- [ ] Naming/readability cukup jelas?

## Cara Menjalankan Lab

```bash
cd labs/09-code-review
go test -v -race ./...
```

## Expected Result

- Test `TestImprovedCheckout_*` semua passing
- Test naive overselling (`TestNaiveCheckout_DeterministicallyDemonstratesOverselling`) membuktikan stock dapat negatif secara deterministik
- Rollback test mengembalikan seluruh state setelah order failure (`TestImprovedCheckout_OrderCreationFailureRollsBackStock`)
- Duplicate request mengembalikan cached response tanpa mengurangi stock dua kali (`TestImprovedCheckout_DuplicateRequestReturnsSameResult`)
- Test concurrent request mempertahankan invarian `stock >= 0` (`TestImprovedCheckout_ConcurrentStockReservationPreservesInvariant`)
- Tidak ada race condition yang terdeteksi pada improved implementation via `go test -race`

## Senior Engineer Takeaways

1. **Code review bukan formatting check** - fokus pada correctness dan production risk
2. **Race condition sering tidak terlihat di development** - butuh stress test atau race detector
3. **Transaction tidak cukup** - transaction menjamin atomicity, concurrency control tetap diperlukan
4. **Idempotency penting untuk distributed system** - retry bukan workaround, harus dirancang
5. **Business invariant harus selalu dipertahankan** - stock >= 0 adalah invariant kritikal

## What Senior Reviewer Would Ask

Sebelum approve PR checkout/payment/inventory, Senior Reviewer akan bertanya:

1. Apa business invariant-nya? (contoh: stock >= 0, order total akurat)
2. Apa yang terjadi saat request diulang (*retry*)?
3. Apa yang terjadi jika dua request bersamaan?
4. Apa yang terjadi jika process crash di tengah jalan?
5. Apa yang terjadi jika dependency timeout?
6. Apakah partial state mungkin terjadi?
7. Apakah authorization dilakukan pada resource (cart vs order)?
8. Apakah query scale ke ribuan data?
9. Apakah error dapat diobservasi dengan log/metrics yang cukup?
10. Bagaimana blast radius jika asumsi kita salah?

## Review the Improved Code

**Improved != Perfect**

Implementasi `checkout_improved.go` ini merupakan contoh yang sudah direview namun **tidak sempurna** untuk production. Tujuannya adalah memperlihatkan pola-polanya, bukan menyimpan semua edge-case.

**Pertanyaan untuk Mencari Production Concerns:**

- **Distributed Order ID**: `atomic.Int64` tidak aman di multi-replica. Ada kemungkinan ID collision.
- **Durable Notification**: Implementasi synchronous loss-tolerant. Haruskah gunakan Outbox Pattern?
- **Stale Idempotency PROCESSING**: Jika service crash setelah `MarkCompleted` gagal, ada window mana request baru dengan key yang sama akan diterima?
- **Database Isolation Semantics**: Apakah `ReserveStock` memerlukan tingkat isolasi tertentu?
- **Timeout/Cancellation**: Apa yang terjadi jika context dibatalkan di tengah transaction?
- **Transaction Retry/Deadlock**: Apakah ada deadlock detection? Retry mechanism?
- **Cart Mutation**: Kita sudah handle, tapi apa bila cart berubah secara *external* saat request diproses?
- **Metrics & Tracing**: Mapping request ke trace ID untuk observability.

Code review adalah proses **risk reduction**, bukan mencari kode yang secara absolut sempurna. Setiap trade-off harus didokumentasikan dan dikompromikan dengan durasi waktu serta risiko bisnis.

---

## Navigasi

- **Previous**: [Lab 08 — Database Isolation Level](../08-database-isolation-level/)
- **Next**: [Lab 10 — Project Estimation](../10-project-estimation/)
