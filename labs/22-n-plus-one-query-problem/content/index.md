# N+1 Query Problem

## Apa Itu N+1 Query Problem?

N+1 query problem adalah **anti-pattern performa** yang terjadi pada aplikasi yang menggunakan ORM (Object-Relational Mapper) ketika mengakses data relasional secara berulang.

Nama "N+1" berasal dari pola eksekusi querynya:

- **1 query** untuk mengambil koleksi utama (misalnya: semua author)
- **N query** tambahan untuk mengakses relasi pada setiap item (misalnya: post untuk setiap author)

Jika kita mengambil 100 author, maka total query yang dijalankan adalah **101 query** (1 untuk author + 100 untuk post masing-masing author).

## Kenapa N+1 Terjadi?

### Lazy Loading sebagai Default

Di hampir semua ORM, **lazy loading adalah perilaku default**. Ini berarti:

- Relasi (relationship) tidak dimuat saat objek utama di-query
- Relasi baru dimuat **pada saat diakses** (lazy = "dimuat hanya ketika dibutuhkan")

```php
// Contoh dengan Laravel Eloquent
$invoices = Invoice::all(); // 1 query: SELECT * FROM invoices

foreach ($invoices as $invoice) {
    echo $invoice->customer->name; // N query: SELECT * FROM customers WHERE id = ?
}

// Total: 1 + N queries
```

### N+1 Tertjebak di Development

N+1 sering **tidak terlihat di development atau testing** karena dataset kecil. Dengan 3 author, hanya 4 query — terlihat cepat. Tapi di production dengan 10.000 author:

```text
1 + 10.000 = 10.001 queries
```

Ini bisa menyebabkan:
- **Latensi tinggi** (menunggu puluhan ribu query berselang-seluang gilirannya)
- **Connection pool exhaustion** (setiap koneksi database bisa hanya melayani satu query pada satu waktu)
- **Timeout** pada request yang melebihi batas waktu

## Bagaimana Cara Mendeteksi N+1?

N+1 tidak bisa dideteksi dengan **code review biasa** atau testing dengan data kecil. Perlu tooling:

| ORM | Tool / Fitur |
|-----|-------------|
| Rails | `strict_loading`, `rails dev:cache`, Rails Mini Profiler |
| Django | `connection.queries`, `django-debug-toolbar` |
| Laravel | Laravel Debugbar, query log |
| SQLAlchemy | `echo=True`, event listeners, `raiseload()` |
| EF Core | Query logging, interceptor |

```python
# Django - melihat query yang dijalankan
from django.db import connection

def my_view(request):
    blogs = Blog.objects.all()
    for blog in blogs:
        print(blog.entry_set.all())  # Ini trigger lazy loading

    # Cek query count
    print(f"Total queries: {len(connection.queries)}")
```

## Cara Memperbaiki N+1

### 1. Eager Loading (Solusi Utama)

Eager loading memuat relasi **sekali saja** dalam jumlah query yang terbatas, bukan per objek.

```php
// Laravel - Menggunakan with() untuk eager loading
$invoices = Invoice::with('customer')->get();

foreach ($invoices as $invoice) {
    echo $invoice->customer->name; // Tidak ada query tambahan
}

// Total: 2 query (1 untuk invoices + 1 untuk customers dengan WHERE id IN (...))
```

Berbeda ORM memiliki API berbeda tapi konsepnya sama:

| ORM | Eager Loading Method | Mekanisme |
|-----|---------------------|-----------|
| Rails | `includes()` | 2 query dengan `WHERE id IN (...)` |
| Django | `prefetch_related()` | 2 query dengan `WHERE id IN (...)` |
| Django | `select_related()` | 1 query dengan JOIN (untuk FK/O2O) |
| Laravel | `with()` | 2 query dengan `WHERE id IN (...)` |
| SQLAlchemy | `selectinload()` | 2 query dengan `WHERE IN (...)` |
| EF Core | `Include()` | JOIN atau split query |

### 2. Aggregation: Gunakan withCount() Ketika Hanya Butuh Jumlah

Seringkali kita hanya butuh **jumlah** record terkait, bukan objeknya:

```php
// N+1 - buruk
$posts = Post::all();
foreach ($posts as $post) {
    echo "Comments: " . $post->comments->count(); // N query tambahan
}

// Solusi: withCount()
$posts = Post::withCount('comments')->get();
foreach ($posts as $post) {
    echo "Comments: " . $post->comments_count; // N+1 dihindari
}

// Total: 2 query (1 untuk posts + 1 untuk COUNT)
```

### 3. Column Selection: Hanya Pilih Kolom yang Dibutuhkan

Jika kita tidak butuh seluruh objek, gunakan `select()`, `pluck()`, atau `values()`:

```php
// Hanya butuh nama author, bukan seluruh objek
$authors = Invoice::pluck('customer.name'); // Laravel
// atau
$authorNames = Invoice::values_list('customer__name', flat=True) # Django
```

### 4. Pagination: Batasi Jumlah Data per Request

```python
# Django - pagination
blogs = Blog.objects.all()[:50]  # LIMIT 50
```

Pagination tidak memperbaiki N+1, tapi **membatasi dampaknya** ketika dataset besar.

## Trade-off: Jangan Eager Loading Apa Saja

Eager loading juga punya risiko:

1. **Cartesian Product Explosion**: JOIN dengan koleksi (one-to-many, many-to-many) bisa menghasilkan hasil yang dikalikan (10 author × 5 post = 50 baris gabungan).

2. **Memory Overhead**: Memuat seluruh objek besar ketika hanya butuh sebagian field.

3. **Over-fetching**: Memuat relasi yang tidak dipakai di view tertentu.

```text
Contoh:
- Author: 10.000 baris
- Post: relasi one-to-many (rata-rata 10 post/author)
- Eager load ALL with JOIN → 100.000 baris gabungan
- Memory: 100x lipat dibanding lazy loading per-author!
```

## Workflow yang Direkomendasikan

Berdasarkan dokumentasi Django ("Profile first") dan konsensus dari semua ORM:

1. **Ukur dulu**: Hitung query count, lihat query apa yang paling mahal
2. **Identifikasi pola N+1**: Apakah ada lazy loading berulang?
3. **Perbaiki root cause**: Eager loading, aggregation, column selection
4. **Baru pertimbangkan caching**: Cache hanya menyembunikan masalah, bukan memperbaikinya

## Ringkasan

N+1 query problem adalah anti-pattern yang:
- **Definition**: 1 query + N query untuk N relasi
- **Root cause**: Lazy loading sebagai default di semua ORM
- **Solusi**: Eager loading, aggregation (`withCount`), column selection, pagination
- **Deteksi**: Query logging, strict loading mode, profiling tools
- **Trade-off**: Jangan eager loading sembarangan — bisa menyebabkan over-fetching

### Catatan tentang Angka Spesifik

Angka-angka spesifik di skenario lab (misalnya "712 query = 2.4 detik") adalah **contoh ilustratif**, bukan benchmark universal. Performa sebenarnya bergantung pada:

- Engine database (PostgreSQL, MySQL, Oracle, dsb.)
- Schema dan index
- Network latency
- Konfigurasi connection pool
- Beban (concurrent requests)

Angka 180ms sebagai target performa juga spesifik ke skenario ini, bukan aturan universal.
