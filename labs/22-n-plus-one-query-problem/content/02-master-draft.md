# N+1 Query Problem: Analisis dan Solusi

## Problem

N+1 query problem adalah anti-pattern performa yang terjadi ketika aplikasi mengambil koleksi data dengan satu query (1), lalu untuk setiap item dalam koleksi tersebut, mengeluarkan query terpisah untuk mengambil data terkait (N). Akibatnya, operasi yang diharapkan butuh 1 query akhirnya butuh N+1 query — secara linear berlipat dengan ukuran koleksi.

Dalam lab ini, kita mengambil 3 author, kemudian untuk masing-masing author, kita mengambil postingannya satu per satu. Hasilnya adalah 1 query untuk semua author ditambah 3 query untuk postingan = 4 query total. Pada skenario nyata dengan ribuan record, pola ini mengakibatkan ribuan query ekstra.

## Why This Matters

Setiap query individu pada dataset kecil terlihat cepat — sering sub-milidetik. Karena itu, N+1 masih lolos dari deteksi pada development dan *slow query logs*. Degradasi performa terjadi secara diam-diam: latensi jaringan akumuler, overhead koneksi database berulang, hingga kehabisan koneksi pool (connection pool exhaustion) pada beban produksi.

Dengan demikian, N+1 adalah masalah "invisible until it hurts" — kode berfungsi sempurna pada data uji coba kecil, tetapi mogok di produksi ketika volume naik.

## Mental Model

Hindari eksekusi query di dalam loop. Alihkan ke pola *batch*: kumpulkan semua identifier relasi yang dibutuhkan, kirimkan satu query berisi klausa `IN (...)`, kemudian gabungkan hasilnya di memori aplikasi.

## Core Concept

- **Naif (N+1):** Ambil semua author → iterasi tiap author → ambil post untuk setiap author.
- **Solusi (Eager Loading / Batching):** Ambil semua author → kumpulkan ID author → ambil semua post dengan satu query `IN (id1, id2, ...)` → kelompokkan post per author di memori → susun struktur akhir.

Kunci perubahan: mengganti N query individual dengan 1 query batch. Kompleksitas query berkurang dari O(N) menjadi O(1) relatif terhadap jumlah record (berganti menjadi O(k) di mana k = jumlah tipe relasi).

## Architecture & Implementation

Lab ini memakai arsitektur tiga layer:

### Store — Mock Database dengan Query Counter

`internal/blog/store.go` mensimulasikan database in-memory yang thread-safe. Setiap metode retrieval (`GetAllAuthors`, `GetPostsByAuthorID`, `GetPostsByAuthorIDs`) men-increment `queryCount`. Desain ini memungkinkan verifikasi kuantitatif — bukan estimasi — atas jumlah query yang dieksekusi.

### Repository — Data Access Layer

`internal/blog/repository.go` menyajikan dua implementasi:
1. **`GetAuthorsWithPostsNPlusOne`** — anti-pattern: lazy loading per item.
2. **`GetAuthorsWithPostsEager`** — solusi: batch loading via klausa IN.

### Demo — CLI Runner

`cmd/demo/main.go` mengeksekusi kedua pendekatan dan mencetak statistik query untuk verifikasi visual.

## How It Works

### Langkah-Langkah N+1 (Anti-Pattern)

1. `GetAllAuthors()` → 1 query, kembalikan 3 author.
2. Iterasi tiap author:
   - `GetPostsByAuthorID(1)` → query #2
   - `GetPostsByAuthorID(2)` → query #3
   - `GetPostsByAuthorID(3)` → query #4
3. Total: **4 query** (1 + N).

### Langkah-Langkah Eager Loading (Solusi)

1. `GetAllAuthors()` → 1 query, kembalikan 3 author.
2. Ekstrak semua author ID ke slice `[1, 2, 3]`.
3. `GetPostsByAuthorIDs([1,2,3])` → 1 query batch, kembalikan semua post.
4. Kelompokkan post per author ID di `map[int][]Post` di memori.
5. Bangun hasil akhir dari map.
6. Total: **2 query** (1 + 1 batch), terlepas dari berapa banyak author.

## What the Tests Prove

Unit test di `internal/blog/repository_test.go` membuktikan perilaku secara kuantitatif:

- `TestGetAuthorsWithPostsNPlusOne`: 3 author → tepat **4 query** (1 + 3)
- `TestGetAuthorsWithPostsEager`: 3 author → tepat **2 query** (1 + 1 batch)
- `TestGetAuthorsWithPostsEager` juga memastikan hasil data identik dengan N+1 via `reflect.DeepEqual` — bukti bahwa batching tidak mengubah semantik, hanya mengurangi I/O.
- `TestEmptyStore`: edge case kosong menghasilkan hasil yang konsisten dari kedua metode.

Hasil eksekusi test:
```
=== RUN   TestGetAuthorsWithPostsNPlusOne
--- PASS: TestGetAuthorsWithPostsNPlusOne (0.00s)
=== RUN   TestGetAuthorsWithPostsEager
--- PASS: TestGetAuthorsWithPostsEager (0.00s)
=== RUN   TestEmptyStore
--- PASS: TestEmptyStore (0.00s)
PASS
```

Race detector juga lolos tanpa data race.

## Recovery / Rollback

Jika solusi eager loading menyebabkan masalah (misalnya memory bloat karena over-fetching — lihat Common Mistakes), prosedur recovery:

1. **Profile kembali**: Ukur query count dan memory usage.
2. **Rollback ke batch kecil**: Ganti satu batch gede dengan beberapa batch berukuran 100–500 (chunked batching).
3. **Gunakan column limiting**: Ambil hanya kolom yang dibutuhkan, bukan seluruh entitas.
4. **Pertimbangkan caching**: Cache hasil relasi yang relatif statis.

Namun, pada lab ini tidak ada mekanisme rollback yang diimplementasikan — ini adalah demonstrasi konseptual, bukan sistem produksi.

## Production Considerations

### Skala N
- N+1 sangat terlihat ketika N ratusan atau ribuan. Pada lab ini N=3 hanya untuk demonstrasi; prinsipnya identik.

### Connection Pool
- Pada produksi, N query berurut menghabiskan slot koneksi pool secara sequential. Under concurrent load, ini bisa menguras pool. *Connection pool exhaustion cascade* adalah konsekuensi logis dari kenyataan bahwa tiap request N+1 mengkonsumsi banyak koneksi berurutan — namun nilai ambangnya dependen pada konfigurasi pool size, jumlah concurrent request, dan latency masing-masing query.

### ORM-Specific Notes
- Semua ORM besar (Django, Rails, Laravel, SQLAlchemy, EF Core) memiliki mekanisme eager loading masing-masing (`select_related`, `includes`, `with`, `selectinload`, `Include`).
- Perbedaan implementasi antar ORM adalah pada API naming dan strategi join vs split query, bukan pada konsep dasar N+1 vs eager loading.

### Trade-off: Cartesian Explosion
- Eager loading via JOIN bisa menyebabkan *cartesian product explosion* pada relasi one-to-many atau many-to-many, di mana satu baris parent dikali jumlah row child. SQLAlchemy, Django, dan EF Core semuanya mendokumentasikan trade-off ini.

## Common Mistakes

### Mapping-level eager loading vs Query-level eager loading

**Mapping-level eager loading** (misalnya `FetchType.EAGER` pada JPA) secara otomatis memuat semua relasi yang dideklarasikan — seringkali melebihi kebutuhan. Ini adalah anti-pattern yang menyebabkan memory bloat karena entitas tidak relevan ikut ter-load.

**Query-level eager loading** (misalnya `select_related()`/`prefetch_related()` di Django, `includes()` di Rails, atau batch query di repository layer) memilih relasi yang dibutuhkan per-query. Ini adalah solusi yang benar.

Lab ini mendemonstrasikan query-level eager loading. Memori bloat dan OOM akibat mapping-level eager loading tidak ditunjukkan secara eksplisit pada lab ini — dataset sengaja kecil. Namun prinsip ini terverifikasi oleh dokumentasi resmi SQLAlchemy, Django, dan EF Core.

### Membungkitkan hasil query count menjadi rekomendasi universial

Angka-angka pada lab (4 query, 2 query, 3 author) adalah ilustratif untuk dataset tertentu. Pada produksi, query count bergantung pada kompleksitas skema, jumlah relasi, dan keputusan batching. Jangan mensyokongkan pola ini sebagai aturan absolut tanpa profiling.

## Failure Scenario

Scenario: Tim meluncurkan endpoint `GET /api/work-orders` yang mengembalikan 100 work order. Setiap work order memiliki 5 item terkait. Pendekatan N+1 yang tidak disadari menghasilkan 101 query (1 untuk work order + 100 untuk items). Pada data uji coba dengan 5 work order, hanya 6 query — terlihat cepat. Di produksi dengan 100 work order, request memakan waktu 2.4 detik (angka ilustratif dari spesifikasi lab), connection pool ter-exhaust, dan halaman menjadi lambat.

## Checklist

- [ ] Profiling query count dilakukan sebelum deploy ke production
- [ ] Lazy loading dalam loop di-audit dan diganti dengan batch/eager loading
- [ ] Eager loading dilakukan per-query, bukan per-mapping (hindari FetchType.EAGER global)
- [ ] Cartesian explosion dicek untuk relasi one-to-many/many-to-many
- [ ] Column selection (`SELECT col1, col2`) digunakan ketika full entity tidak dibutuhkan
- [ ] Pagination diterapkan untuk hasil yang besar
- [ ] Query count diuji otomatis (unit/integration test)
- [ ] Strict loading mode dipertimbangkan untuk ORM yang mendukungnya

## Key Takeaways

1. N+1 query problem menyebabkan latensi akumuler tinggi dan sering lolos dari slow query logs karena setiap query individu cepat.
2. Fenomena N+1 terjadi pada lapisan database (ORM lazy loading) dan lapisan jaringan (GraphQL resolvers, HTTP per-object requests).
3. Solusi utama: batch semua relasi dalam 1-2 query dengan klausa `IN (...)` atau eager loading level query.
4. Mapping-level eager fetching (JPA `FetchType.EAGER`) adalah anti-pattern yang menyebabkan memory bloat. Perhatikan perbedaan dengan query-level eager loading.
5. Eager loading via JOIN dapat menyebabkan cartesian product explosion pada relasi koleksi.
6. Column selection (`pluck`, `values`, `select`) dan aggregation (`withCount`, `annotate`) sering lebih baik daripada eager loading penuh ketika hanya scalar dibutuhkan.
7. Profiling-first workflow — ukur jumlah query dulu, baru optimasi — adalah best practice yang didokumentasikan Django dan semua ORM.
8. Unit test dengan query counter dapat memverifikasi reduksi N+1 secara otomatis dan regresi-mencegah.
9. Angka performa pada lab (4 → 2 query untuk 3 author) bersifat ilustratif; faktor nyata bergantung pada volume data, schema, dan infrastruktur.

## Sources

### Research
- `research/runs/2026-09-25-n-plus-one-query-problem/05-report.md` — Finding 1 (Definition), Finding 2 (Eager Loading), Finding 3 (Trade-offs), Finding 4 (Column Selection), Finding 5 (Lazy Loading as Root Cause), Finding 6 (Detection Tooling), Finding 7 (Profiling-First), Finding 8 (Network/API N+1)
- `research-audit/07-verdict.md` — APPROVED

### Engineering
- `engineering/01-design.md` — Architecture design
- `engineering/02-implementation-notes.md` — Implementation decisions
- `engineering/03-execution-result.md` — Verified test output and demo output
- `engineering-audit/06-verdict.md` — APPROVED

### Implementation
- `internal/blog/models.go` — Domain models
- `internal/blog/store.go` — Mock store dengan query counter
- `internal/blog/repository.go` — Repository layer N+1 dan eager loading
- `internal/blog/repository_test.go` — Unit tests verifikasi query count
- `cmd/demo/main.go` — CLI demo

### External References (from research)
- Django Database Access Optimization: https://docs.djangoproject.com/en/5.1/topics/db/optimization/
- Rails Active Record Query Interface: https://guides.rubyonrails.org/active_record_querying.html
- Laravel Eloquent Relationships: https://laravel.com/docs/13.x/eloquent-relationships
- SQLAlchemy Relationship Loading Techniques: https://docs.sqlalchemy.org/en/21/orm/queryguide/relationships.html
- EF Core Eager Loading: https://learn.microsoft.com/en-us/ef/core/querying/related-data/eager
