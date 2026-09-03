# Lab 09 — Code Review: Cara Senior Engineer Menemukan Bug Sebelum Masuk Production

## Tujuan

Melatih kemampuan **code review** berbasis risiko dengan memanalisis implementasi checkout yang memiliki multiple masalah serta membandingkan dengan versi yang sudah direview.

Senior engineer tidak hanya mengecek formatting atau naming. Mereka memfokuskan review pada:

- **Business Logic Correctness**
- **Concurrency & Race Conditions**
- **Data Integrity**
- **Error Handling**
- **Security**
- **Performance**

## Business Case: Checkout Sistem

Berikut pseudo-code PHP yang menggambarkan masalah utama:

```php
public function checkout(Request $request)
{
    $cart = Cart::where('user_id', auth()->id())->get();

    foreach ($cart as $item) {
        $product = Product::find($item->product_id);
        $product->stock -= $item->qty;  // BUG: non-atomic, race condition
        $product->save();               // BUG: N+1 query
    }

    Order::create(...);  // BUG: no transaction, partial state on failure

    return response()->json(['success' => true]);  // BUG: no idempotency, error swallowed
}
```

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

## Severity Classification

| Level | Deskripsi | Contoh |
|-------|-----------|--------|
| BLOCKER | Risiko data korupsi, kehilangan data, atau security critical | Race condition overselling |
| MAJOR | Masalah bisa menyebabkan error di production | No transaction |
| MINOR | Masalah performa atau user experience | N+1 query |
| NIT | Formatting, naming yang tidak konsisten | Inconsistent naming |

## Business Logic Review

### [ ] Apakah business logic benar?

- Stock deduction harus atomic
- Validasi kuantitas sebelum proses
- Total harga harus akurat

### [ ] Apakah invariant tetap terjaga?

- Stock >= 0 (tidak boleh negatif)
- Order hanya dibuat jika semua validasi lolos

## Concurrency Review

### [ ] Apakah concurrent request aman?

Run test: `go test -race ./...` untuk mendeteksi race condition.

**Naive:** Tidak aman - dapat overselling
**Improved:** Atomic stock update, satu request berhasil satu gagal

## Security Review

### Authentication vs Authorization

| Aspek | Deskripsi |
|-------|-----------|
| **Authentication** | Siapa yang login? (user ID) |
| **Authorization** | Apakah user boleh akses resource tersebut? (cart milik user yang sama) |

Implementation harus memastikan:
- User hanya dapat checkout cart miliknya
- Tidak ada akses ke cart orang lain

```go
// Benar: menggunakan userID dari konteks
cartItems, _ := c.cartSource.GetCart(ctx, userID)
```

## Performance Review

### N+1 Query Problem

Naive implementation melakukan `GetProduct` dalam loop:

```go
for _, item := range cartItems {
    product, _ := c.products.GetProduct(ctx, item.ProductID)  // N query
    // ...
}
```

**Improved:** Pre-fetch semua product sekaligus atau gunakan batch query.

## Error Handling & Logging

### [ ] Apakah error ditangani dengan baik?

- [ ] Empty cart
- [ ] Product tidak ditemukan
- [ ] kuantitas tidak valid
- [ ] Stock tidak cukup

### [ ] Apakah logs punya context?

Logging harus mencakup:
- User ID
- Product ID
- Quantity
- Hasil operasi

## Improved Implementation

Lihat `checkout_improved.go` untuk implementasi yang sudah direview dengan:

1. **Validasi input** - empty cart, product existence, quantity, stock
2. **Idempotency** - mencegah duplicate request
3. **Atomic stock update** - cek dan update dalam satu operasi
4. **Good error propagation** - mengembalikan error yang sesuai
5. **Notification setelah transaction** - tidak ada side effect di dalam transaksi

## Failure Scenario Tests

### Test: Stock Tidak Cukup

Jika stock habis, checkout harus gagal tanpa mengubah stock.

### Test: Duplicate Request

Client dapat mengirim request yang sama dua kali. Implementasi harus:
- Mengembalikan hasil yang sama (response replay)
- Tidak mengubah state dua kali

### Test: Concurrent Checkout

Dua request checkout untuk item dengan stock=1 yang dijalankan bersamaan.
- Hanya satu yang berhasil
- Satu yang gagal dengan error
- Stock tetap >= 0

## Risk-Based Review

Ukuran PR bukan ukuran risiko. Fokus pada area berbahaya:

| Area | Risk Level |
|------|------------|
| CSS/HTML change | Low |
| Reporting feature | Medium |
| Authentication | High |
| Payment / Inventory | **Critical** |

**Contoh matrix risiko:**

| Risiko | Dampak | Probabilitas | Tindakan |
|--------|--------|--------------|----------|
| Stock negative | Critical | Low | Harus dicek di unit test |
| Duplicate order | High | Medium | Perlu idempotency key |
| Log error | Medium | High | Check logging coverage |

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
- Test `TestNativeCheckout_RaceCondition` skip (membuktikan bug)
- Tidak ada race condition yang terdeteksi pada improved implementation

## Senior Engineer Takeaways

1. **Code review bukan formatting check** - fokus pada correctness dan production risk
2. **Race condition sering tidak terlihat di development** - butuh stress test atau race detector
3. **Transaction tidak cukup** - transaction menjamin atomicity, concurrency control tetap diperlukan
4. **Idempotency penting untuk distributed system** - retry bukan workaround, harus dirancang
5. **Business invariant harus selalu dipertahankan** - stock >= 0 adalah invariant kritikal