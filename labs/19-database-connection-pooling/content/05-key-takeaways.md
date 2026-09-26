# Key Takeaways

1. **Lebih Sedikit Justru Lebih Cepat**: Pool ideal mendekati `(core_count × 2)`, bukan sesuai konkurensi aplikasi. Melampaui batas fisik database hanya menambah *context switching* dan *lock contention*.

2. **Overhead Koneksi Signifikan**: Setiap koneksi baru membutuhkan TCP handshake, otentikasi TLS, dan alokasi memori backend. Dengan 10ms handshake delay, 5 request unpooled memakan ~54ms vs 6μs untuk pooled — perbedaan ~8.000x dalam lab.

3. **Perkalian Skala Horisontal Berbahaya**: `total_connections = instance_count × pool_size`. Empat instance dengan pool 16 masing-masing = 64 koneksi terhadap `max_connections=200` — masih aman. Tapi 4 instance × 50 pool = 200, tidak ada ruang untuk sistem operasi atau cadangan.

4. **I/O Eksternal di Dalam Transaksi = Kebocoran**: Fungsi seperti `ProcessOrderUnsafeLeak` memegang koneksi selama panggilan jaringan lambat. Koneksi tidak kembali ke pool, menyebabkan *starvation* pada request lain.

5. **Cloud Provider Menyarankan Pooling Eksternal**: AWS (RDS Proxy), Azure (PgBouncer built-in), dan Google Cloud (HikariCP) tidak menyarankan menambah `max_connections` tanpa batas.

6. **Pantau Pool-Level Metrics**: CPU dan memory tidak menunjukkan masalah kebocoran. Gunakan `Wait Timer`, `ActiveConnections`, `IdleConnections`, dan `pg_stat_activity` dengan fokus pada status `idle in transaction`.

7. **PgBouncer Transaction Mode Menghemat Tapi Membatasi**: Mode ini melepaskan koneksi setelah transaksi (bukan sesi), tapi memutuskan fitur seperti `LISTEN`, `WITH HOLD CURSOR`, dan advisory locks.

8. **Formula Pool Sizing untuk SSD Belum Terbukti**: Formula `((core_count * 2) + effective_spindle_count)` dikembangkan untuk HDD. Untuk SSD, HikariCP menyarankan *fewer* connections lebih baik, tapi tidak ada data empiris yang mengkonfirmasi.

9. **Pool-Locking Deadlock Formula**: Jika satu thread butuh beberapa koneksi: `pool_size = Tn × (Cm - 1) + 1`. Ini minimum untuk mencegah deadlock, bukan ukuran optimal.

10. **Persiapan untuk Kegagalan**: Dengan `max_connections`, selalu sediakan reserved connections untuk superuser dan monitoring (PostgreSQL default: 3 superuser slots). Cloud providers (Azure) menyisihkan 15 slot untuk sistem.
