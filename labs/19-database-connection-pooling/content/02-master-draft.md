# Database Connection Pooling: Overhead, Batas Server, dan Pencegahan Kebocoran

## Problem

Aplikasi menghasilkan error koneksi (`FATAL: sorry, too many clients already` pada PostgreSQL) sementara CPU dan memory database terlihat sehat. Respons pertama tim sering kali salah: memperbesar *connection pool* atau menaikkan `max_connections` server. Hasilnya justru latency memburuk, bukan membaik.

Dua kesalahan berulang yang ditangani lab ini:

1. Membuka koneksi baru untuk setiap request tanpa memakai ulang koneksi yang sudah ada (*unpooled*).
2. Memperbesar pool melebihi kapasitas server, atau menahan koneksi database selama operasi I/O eksternal yang lambat — sehingga koneksi tidak pernah kembali ke pool (*leak*) dan request lain kelaparan (*pool starvation*).

## Why This Matters

`max_connections` adalah batas slot keras, bukan ambang performa. PostgreSQL mengalokasikan memori backend per koneksi; Azure mencatat bahwa "each connection, regardless of whether it's idle or active, consumes a significant amount of resources". Koneksi pendek (< 60 detik) juga menaikkan CPU karena overhead pembukaan/pemutusan berulang.

Di sisi lain, melewati titik jenuh (*knee*) pada kurva koneksi-vs-throughput membuat performa menurun — bukan naik. PostgreSQL Wiki menyebut enam mekanismenya: disk contention, tekanan RAM (`work_mem` × jumlah koneksi), lock contention, context switching, cache line contention, dan struktur internal berskala O(N²). Artinya: pertanyaannya bukan "sebesar apa pool yang aman?", melainkan "sekecil apa pool yang masih memenuhi kebutuhan throughput?"

Kesalahan sizing juga terakumulasi lintas instance. Empat instance × 16 worker × 10 koneksi = 640 potensi koneksi terhadap `max_connections = 200`. Pool harus dihitung di level deployment, bukan per instance.

## Mental Model

Kapasitas koneksi ditentukan oleh hardware database dan konfigurasi server — bukan oleh jumlah goroutine atau rasio koneksi-per-thread di aplikasi.

- **Pool kecil & tepat** → antrean pendek, latensi terprediksi, koneksi selalu valid.
- **Pool terlalu besar** → server menolak koneksi baru, atau throughput jatuh karena perebutan resource.
- **Koneksi ditahan terlalu lama** → pool kosong, request sehat ikut timeout.

Analogi: koneksi database adalah meja di restoran. Menambah meja melewati kapasitas ruangan tidak menambah kapasitas dapur — hanya memperpanjang antrian dan memperlambat semua orang.

## Core Concept

### 1. Overhead koneksi baru itu nyata

Setiap koneksi baru melewati TCP handshake, otentikasi, dan alokasi resource server. Lab ini mensimulasikan latency tersebut dengan `connectDelay` pada driver mock: 10ms pada demo, 5ms pada test. Hasil demo (angka bervariasi antar runtime):

```text
Unpooled (5 requests): 54.593875ms
Pooled (5 requests): 6.333µs
```

Yang diverifikasi test bukan nilai absolutnya, melainkan urutannya: pooled selalu lebih cepat dari unpooled (`TestDirectConnectionOverhead`), dan jumlah koneksi yang dibuat oleh pooled lebih sedikit (`TestTotalCreatedPoolReuse`).

### 2. Rumus sizing sebagai titik awal

Formula yang paling banyak dikutip (PostgreSQL Wiki dan HikariCP secara independen menyebut formula identik):

```text
connections = ((core_count * 2) + effective_spindle_count)
```

- `core_count` tidak termasuk thread hyperthreading.
- `effective_spindle_count` mendekati 0 bila dataset aktif sepenuhnya cache di RAM; mendekati jumlah disk fisik saat cache hit rate turun.
- Contoh: server 4-core dengan 1 disk → `((4×2)+1) = 9`, dibulatkan 10.

**Caveat yang wajib dipertahankan:** kedua sumber mengakui formula ini divalidasi terutama pada benchmark berbasis HDD (artikel wiki terakhir disunting 2014). HikariCP menulis eksplisit: "There hasn't been any analysis so far regarding how well the formula works with SSDs." Argumen HikariCP untuk SSD — operasi tanpa seek berarti lebih sedikit blocking sehingga *lebih sedikit* thread yang lebih baik — bersifat teoretis, tanpa benchmark SSD spesifik. Formula adalah titik awal yang wajib divalidasi lewat load test, bukan angka universal.

Contoh "3000 front-end users at 6000 TPS dengan 10 koneksi" berasal dari wiki HikariCP dengan redaksi hedged ("we'd wager") — ilustrasi berdasarkan pengalaman, bukan hasil benchmark terpublikasi.

### 3. Exhaustion bukan masalah performa database

Saat total koneksi seluruh instance melebihi `max_connections`, server menolak koneksi baru. PostgreSQL default 100; Azure menghitung `max_connections - (reserved_connections + superuser_reserved_connections)` sebagai kapasitas pengguna efektif dan menyisihkan 15 slot untuk physical replication dan monitoring. PostgreSQL wiki menyarankan `max_connections` sedikit lebih besar dari total pool agar selalu ada slot untuk maintenance dan monitoring.

Cloud provider konsisten menyarankan pooling eksternal, bukan menaikkan `max_connections` tanpa batas:

- Azure: gunakan PgBouncer dalam *transaction mode*, mulai dari 2–5× vCores; "we advise against" menaikkan `max_connections`.
- AWS: RDS Proxy; formula `max_connections = LEAST(DBInstanceClassMemory/9531392, 5000)`.
- Google Cloud: contoh kode dengan HikariCP/SQLAlchemy dan fitur Managed Connection Pooling.

Catatan akurasi: klaim formula AWS berasal dari dokumentasi RDS yang ditelusuri selama riset; confidence MEDIUM karena halaman tidak diambil ulang penuh dalam sesi riset.

### 4. Kebocoran koneksi tidak terlihat oleh monitoring infrastruktur

Gejala leak: response time naik, waktu tunggu koneksi naik, utilisasi pool naik, request timeout, error `Too Many Connections` — sementara CPU/memory normal. HikariCP menyediakan `leakDetectionThreshold` (default 0 = nonaktif, minimum 2000ms untuk mengaktifkan) yang mencetak stack trace ketika koneksi keluar pool melebihi threshold.

Mekanisme ini diverifikasi pada lab lewat pola kodenya, bukan lewat HikariCP: `ProcessOrderUnsafeLeak` menahan koneksi selama `externalCall` berjalan, dan test membuktikan efeknya pada request lain.

### 5. Formula anti-deadlock pool-locking

Ketika satu thread membutuhkan beberapa koneksi secara bersamaan, jumlah minimum agar tidak deadlock:

```text
pool_size = Tn × (Cm - 1) + 1
```

`Tn` = jumlah maksimum thread, `Cm` = koneksi maksimum yang dibutuhkan satu thread. Ini **minimum** untuk menghindari deadlock, bukan ukuran optimal. Contoh dokumentasi HikariCP: 3 thread × 4 koneksi → 10; 8 thread × 3 koneksi → 17.

## Failure Scenario

**Skenario A — Pool klien melebihi limit server.** Server dibatasi 15 koneksi, pool klien diizinkan 50, 30 request berjalan konkuren. Lima belas pertama diterima, lima belas berikutnya ditolak `ErrServerOverloaded` (demo menghasilkan `Client attempted: 30, Succeeded: 15, Server Rejected: 15`; test `TestOversizedPoolExhaustsServerConnections` menegaskan minimal 10 dari 20 percobaan gagal pada server max 10 / pool 20).

**Skenario B — Leak menstarvasi pool.** Pool dikunci 1 koneksi. Satu request menjalankan `ProcessOrderUnsafeLeak` yang menahan koneksi selama 100ms di dalam `externalCall`. Request kedua — dengan pola aman sekalipun — mengantre, lalu timeout dengan `context deadline exceeded`. Satu kode bocor menghancurkan request yang kode-nya benar.

**Skenario C — Pool-locking deadlock.** Dengan pool size 1, mengakuisisi koneksi kedua sambil koneksi pertama masih dipegang menghasilkan timeout (`TestPoolLockingDeadlock`).

## How It Works

`sql.DB` pada Go sudah mengelola *connection pool* secara internal. Pengaturan perilaku dilakukan lewat `SetMaxOpenConns` (batas koneksi terbuka serentak) dan `SetMaxIdleConns` (berapa koneksi dipertahankan untuk dipakai ulang). Dengan `SetMaxIdleConns(0)`, koneksi ditutup setelah dipakai — perilaku *unpooled* yang memaksa pembukaan koneksi baru di tiap query.

Alur pooling:

```text
Request → pinjam koneksi dari pool → eksekusi query → kembalikan ke pool (bukan tutup)
                │
                └─ pool kosong → antre hingga context timeout
```

Dalam arsitektur multi-instance, layer proxy seperti PgBouncer menerjemahkan ribuan koneksi aplikasi menjadi belasan koneksi nyata ke database. Tiga mode PgBouncer punya kompatibilitas berbeda: *session* (fitur lengkap, pooling paling lemah, default, `default_pool_size = 20`), *transaction* (dianjurkan cloud provider, tetapi memutus `SET/RESET`, `LISTEN`, `WITH HOLD CURSOR`, `PREPARE/DEALLOCATE`, advisory lock level sesi, `LOAD`, dan temp tables), dan *statement* (selain itu melarang transaksi multi-statement).

## Architecture

```text
[Aplikasi konkuren]
       │
[Client Pool — sql.DB: MaxOpenConns / MaxIdleConns] ── timeout saat kosong
       │
[MockDriver / Proxy / Database]
   maxConnections, connectDelay
       │
[Backend DB: max_connections + reserved slots]
```

Lab memakai mock `database/sql/driver` alih-alih PostgreSQL sungguhan — keputusan desain eksplisit di `engineering/01-design.md` untuk menghilangkan dependensi eksternal dan mensimulasikan latency handshake serta batas `max_connections` secara deterministik. Batas dan antrean diwakili struktur memori (`int32` dengan operasi atomik via `sync/atomic`), bukan proses Postgres nyata.

## Implementation

Empat komponen lab:

| File | Peran |
|---|---|
| `internal/pool/mockdb.go` | Mock driver: `maxConnections`, `connectDelay`, penghitungan `activeConns`/`totalCreated` |
| `internal/pool/service.go` | `OrderService`: `ProcessOrderSafe` vs `ProcessOrderUnsafeLeak` |
| `tests/pool_test.go` | 10 test: overhead, exhaustion, starvation, deadlock, error propagation, reuse |
| `cmd/demo/main.go` | Demo CLI tiga bagian: overhead, rejection, starvation |

## Code Walkthrough

Detail tiap snippet dan sumbernya ada di `03-code-snippets.md`. Ringkasan alur:

1. **`MockDriver.Open`** — menahan `connectDelay`, lalu menolak koneksi bila `activeConns >= maxConnections` (`ErrServerOverloaded`), selain itu menaikkan hitungan dan mengembalikan `mockConn`.
2. **`mockConn.Close`** — flag `closed` ber-double-checked di bawah mutex menjamin dekrement hanya sekali (divalidasi `TestMockConnDoubleClose`).
3. **`ProcessOrderSafe`** — `externalCall` dijalankan *sebelum* meminjam koneksi; koneksi dipakai hanya untuk `ExecContext` singkat, lalu `defer conn.Close()`.
4. **`ProcessOrderUnsafeLeak`** — meminjam koneksi, eksekusi UPDATE, lalu menjalankan `externalCall` sementara koneksi masih terbuka.
5. **`demoDirectOverhead`** — membandingkan `SetMaxIdleConns(0)` terhadap pool yang sudah di-*warm up*.
6. **`demoOversizedPool`** — 30 goroutine terhadap server 15 / pool 50.
7. **`demoConnectionLeak`** — dua unsafe order menghabiskan pool size 2; order ketiga (safe, timeout 100ms) gagal.

## What the Tests Prove

Eksekusi aktual (`go test -v ./...` dan `go test -race -v ./...`, keduanya PASS 10/10, tanpa data race; `go build ./...` sukses):

| Test | Perilaku yang dibuktikan |
|---|---|
| `TestDirectConnectionOverhead` | Pool dengan idle connections lebih cepat daripada membuka koneksi baru tiap request |
| `TestTotalCreatedPoolReuse` | Pool membuat lebih sedikit koneksi daripada unpooled untuk 5 query |
| `TestOversizedPoolExhaustsServerConnections` | Server max 10 + pool 20 → minimal 10 percobaan gagal |
| `TestConnectionStarvationDueToLeak` | Satu koneksi ditahan saat external call → request berikutnya `context deadline exceeded` |
| `TestSafeProcessingConcurrently` | 20 request aman pada pool 5 dengan I/O eksternal singkat → semua sukses |
| `TestPoolLockingDeadlock` | Koneksi kedua pada pool size 1 → timeout |
| `TestMockConnDoubleClose` | `Close` ganda tidak panik dan tidak double-decrement |
| `TestExternalCallErrorPropagation` | Error `externalCall` diteruskan pada pola safe & unsafe; `activeConns` kembali 0 |
| `TestPreCancelledContextProcessOrderSafe` | Context yang sudah dibatalkan → error, bukan sukses |
| `TestUnsafeLeakExecContextFailure` | Pool terkunci → `ProcessOrderUnsafeLeak` gagal; setelah `db.Close()`, `activeConns` = 0 |

Yang **tidak** dibuktikan lab ini (dan karena itu tidak boleh diklaim): degradasi performa query akibat lock contention internal database, perilaku PgBouncer, alokasi memori per backend process yang menyebabkan OOM pada koneksi tinggi, serta `pg_stat_activity` waiting states.

## Recovery / Rollback

Kondisi pool terserang leak atau kelebihan beban:

1. **Identifikasi lewat pool-level metrics**, bukan CPU/memory: active/idle/pending connections, connection acquisition wait time, dan `pg_stat_activity` (perhatikan status `idle in transaction`).
2. **Kembalikan pola kodenya**: pindahkan I/O eksternal ke luar blok yang memegang koneksi (ikuti `ProcessOrderSafe`). `defer conn.Close()` pada setiap jalur — pola ini yang menjamin `activeConns` kembali 0 pada test error-propagation.
3. **Perkecil pool ke baseline**, jangan perbesar. Validasi dengan load test; gunakan formula di atas hanya sebagai titik awal.
4. **Tambahkan deteksi**: `leakDetectionThreshold` (HikariCP) atau setara; pada lab, hitungan `ActiveConnections()`/`TotalCreated()` memperlihatkan kebocoran langsung.
5. **Skala horizontal dengan proxy**: naikkan instance, bukan `max_connections`; arahkan ke PgBouncer (transaction mode) atau RDS Proxy, dan periksa dulu kompatibilitas fitur sesi aplikasi Anda.

## Production Considerations

- Hitung total potensi koneksi seluruh deployment: `instance_count × pool_size`, termasuk reserved slots server (PostgreSQL superuser default 3; Azure 15).
- Sediakan context timeout pada akuisisi koneksi; timeout yang gagal adalah gejala, bukan solusi.
- Pantau wait time akuisisi koneksi, utilisasi pool, dan durasi hold — bukan hanya CPU/memory database.
- Jangan menahan koneksi selama panggilan API eksternal, pembuatan file, atau pemrosesan lama.
- `maxLifetime` pada HikariCP default 30 menit dan harus beberapa detik lebih pendek dari batas lifetime yang dipasang database/infrastruktur; `keepaliveTime` default 2 menit menjaga koneksi tidak diputus perangkat jaringan.
- Pertimbangkan PgBouncer *transaction mode* hanya setelah memverifikasi aplikasi tidak memakai fitur sesi yang rusak (lihat daftar pada Core Concept).
- Statement cache di layer pool adalah anti-pattern menurut HikariCP: cache per koneksi mengalikan rencana eksekusi (250 query × 20 koneksi = 5000 rencana); biarkan driver/database yang menanganinya.

## Common Mistakes

1. **Menjawab error koneksi dengan memperbesar pool.** Gejala menunjukkan kelebihan koneksi, bukan kekurangan.
2. **Menghitung pool per instance tanpa menjumlahkan seluruh deployment.**
3. **Manggil API eksternal di dalam blok yang memegang koneksi** — penyebab starvation terbukti di test.
4. **Menganggap koneksi idle gratis.** Azure dan PostgreSQL wiki keduanya menegaskan koneksi idle tetap memakan resource.
5. **Menyamakan persistent connection dengan pooling.** PostgreSQL wiki: persistent connection (mis. mod_php) "*are not pooling and still require a connection pool*".
6. **Menjadikan angka ilustrasi sebagai rekomendasi universal** (mis. "3000 users / 6000 TPS" atau "50x lebih cepat").
7. **Menaikkan `max_connections` tanpa mengurangi pool klien** — cloud provider justru menyarankan sebaliknya.

## Case Study

**Oracle: pengurangan koneksi 2048 → 96.** Menurut HikariCP wiki yang mengutip video Oracle Real-World Performance group, mengurangi ukuran pool semata menurunkan response time dari ~100ms ke ~2ms — "over 50x improvement".

**Status bukti: MEDIUM.** Ini demonstrasi vendor lewat satu video, bukan benchmark independen; tidak ada publikasi akademis yang mereplikasi angka spesifik tersebut. Prinsip yang mendasarinya (koneksi lebih sedikit = kontensi lebih sedikit) diverifikasi independen oleh analisis PostgreSQL wiki. Dalam lab ini, prinsip yang sama terlihat pada demo overhead dan starvation — bukan angka 50x.

**Catatan berlaku:** seluruh angka dalam artikel ini (durasi demo, 54ms vs 6μs, 15 dari 30 koneksi) adalah hasil simulasi lab dengan driver mock dan latensi yang dipilih sendiri (5–10ms). Ini ilustrasi perilaku, **bukan** benchmark produksi dan bukan rekomendasi nilai.

## Checklist

- [ ] Pool klien per instance diketahui, dan total `instance_count × pool_size` < `max_connections` server (dengan cadangan untuk reserved slots).
- [ ] Tidak ada panggilan I/O eksternal di dalam blok yang memegang koneksi database.
- [ ] Setiap jalur meminjam koneksi memiliki `defer conn.Close()` (atau pola pemakaian `db.QueryRow`/`Exec` yang otomatis mengembalikan koneksi).
- [ ] Context timeout diset pada akuisisi koneksi dan eksekusi query.
- [ ] Pool-level metrics dipantau (wait time, active/idle, pending) — bukan hanya CPU/memory.
- [ ] Deteksi leak aktif (mis. `leakDetectionThreshold`) bila memakai HikariCP atau setara.
- [ ] Ukuran pool divalidasi dengan load test, memakai rumus `((core_count * 2) + effective_spindle_count)` hanya sebagai titik awal.
- [ ] Bila memakai proxy pooling, mode dan kompatibilitas fitur sesi aplikasi sudah diperiksa.
- [ ] Peningkatan beban dijawab dengan penambahan instance + proxy, bukan kenaikan `max_connections` tanpa analisis.

## Key Takeaways

Lihat `05-key-takeaways.md`.

## Sources

Lihat `06-source-map.md` (pemetaan seksi → file riset, implementasi, dan test) serta `research/02-sources.md` untuk daftar sumber lengkap beserta URL dan tanggal akses (2026-09-26).

### Caveat yang wajib dibawa ke publikasi

1. Rumus `((core_count * 2) + effective_spindle_count)` belum terverifikasi empiris pada storage SSD; kedua sumber mengakuinya.
2. Klaim peningkatan 50x Oracle berasal dari demonstrasi vendor (satu video), bukan benchmark independen.
3. Lab memakai driver mock: kontensi lock internal database, context switching CPU, dan OOM akibat alokasi memori per backend **tidak** disimulasikan.
4. Angka hasil demo bervariasi antar runtime; klaim test adalah urutan dan kondisi kegagalan, bukan nilai absolut.
5. `go.mod` menyatakan `go 1.26.7` yang bukan versi Go nyata (GAP-001, metadata tidak akurat; tidak memengaruhi build/test).
6. Fokus riset adalah PostgreSQL; MySQL, SQL Server, dan RDBMS lain tidak dicakup. Deteksi leak diteliti lewat HikariCP; padanan pada DBCP2, Tomcat JDBC, c3p0 belum diteliti.
