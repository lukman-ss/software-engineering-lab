# Case Studies

## Case Study A — Customer Phone
Old: `phone` (string). New: `customer_phone_numbers` table.
1. Expand: Create `customer_phone_numbers` table.
2. App Migrate (Write): Setiap simpan user baru/update, tulis juga ke tabel `customer_phone_numbers` (dual write/CDC).
3. Data Backfill: Sinkronisasi data lama dari `phone` ke tabel `customer_phone_numbers`.
4. App Migrate (Read): Baca dari `customer_phone_numbers`. Fallback ke `phone` jika tidak ditemukan.
5. Contract: Hapus `phone` (setelah masa transisi selesai).

## Case Study B — CMMS Invoice Mechanics
Old: `mechanic_id` di `invoice`. New: `invoice_mechanics` table.
Tetap biarkan kolom `mechanic_id` ada, tapi buat nullable. Sistem sync `invoice_mechanics` ke `mechanic_id` untuk mekanik pertama sebagai bentuk backward compatibility untuk app reporting lawas, lalu deprecate perlahan.

## Case Study C — Multi Currency
Old: `price` di `products`. New: `product_prices`.
Expand `product_prices`, isi backfill data dari `price`. API masih me-return "price" dari "product_prices" default currency. Client baru memanggil list dari `product_prices`.
