# Database Connection Pooling: Kapasitas, Overhead, dan Mencegah Kebocoran

## Problem
Meningkatkan ukuran *connection pool* (jumlah koneksi database) sering dianggap sebagai solusi saat aplikasi mengalami error antrean. Ini adalah kesimpulan yang salah. Database berbasis *process-per-connection* seperti PostgreSQL memakan banyak resource untuk setiap koneksi (TCP handshake, otentikasi TLS, alokasi memori internal). Saat aplikasi didistribusikan secara horizontal, total koneksi bisa menembus batas maksimal database, menyebabkan penolakan koneksi dan kegagalan sistem.

## Why This Matters
Membuka koneksi baru untuk setiap *request* memicu lonjakan latensi yang parah. Sebaliknya, memperbesar ukuran *pool* melebihi batas perangkat keras (CPU dan I/O) akan menghancurkan *throughput* transaksi. Pemahaman tentang batasan dan pengelolaan *pool* secara ketat memastikan reliabilitas batas akhir aplikasi Anda.

## Mental Model
Kapasitas *connection pool* ditentukan oleh batas maksimum perangkat keras database (jumlah *core* dan arsitektur disk), bukan oleh jumlah konkurensi di level aplikasi. Lebih dari batas fisik ini, sistem hanya melakukan *context switching* yang memperlambat semua proses. 

## Core Concept
1. **Connection Overhead**: Membuka koneksi memakan waktu. Memakai ulang koneksi yang sudah terbuka (pooling) jauh lebih cepat dan murah.
2. **Kapasitas Pool Saturated**: Rumus dasar untuk ukuran *pool* pada *storage* modern umumnya mendekati `(core_count * 2)`.
3. **Exhaustion (Kehabisan Koneksi)**: Ketika ukuran *pool* di aplikasi dikalikan dengan jumlah *instance* (skala horizontal) melebihi parameter server (`max_connections`), server akan menolak permintaan koneksi baru.
4. **Kebocoran Koneksi (Leaks)**: Terjadi jika kode aplikasi melakukan operasi I/O eksternal lambat (seperti panggilan API) di dalam blok transaksi aktif, menahan koneksi dan membuat *request* lain kelaparan (*pool starvation*).

## Failure Scenario
Aplikasi memiliki *pool* berukuran 50 koneksi namun database hanya mengizinkan maksimal 15 koneksi server. Saat beban aplikasi tinggi dan konkurensi meningkat, 15 koneksi pertama diterima, sedangkan sisanya langsung ditolak (*server overloaded*). 

Skenario kegagalan lainnya terjadi ketika operasi eksternal memblokir *thread* namun masih mengunci koneksi database. Koneksi tidak kembali ke *pool*, menyebabkan panggilan dari pengguna lain mengalami *timeout*.

## How It Works
*Connection pool* menyimpan koneksi yang sudah terbuka dan tervalidasi. Saat *request* datang, aplikasi "meminjam" koneksi, mengeksekusi *query*, lalu mengembalikannya ke *pool* (bukan menutupnya). 
Dalam arsitektur terdistribusi yang aman, sering digunakan *proxy* tambahan (misalnya PgBouncer dalam *transaction mode*) di depan database untuk menerjemahkan ribuan koneksi aplikasi menjadi belasan koneksi nyata ke database.

## Implementation
Simulasi lab membuktikan hal ini melalui:
- Driver database *mock* yang menerapkan limitasi (`maxConnections`) dan latensi inisiasi.
- Simulasi koneksi tak di-*pool* vs di-*pool* untuk membuktikan *overhead*.
- Layanan pesanan (*OrderService*) yang mendemonstrasikan proses peminjaman koneksi yang aman vs yang bocor.

## Code Walkthrough
- **Direct Connection Overhead**: Kode mensimulasikan latensi *handshake* database (misalnya 10ms). Mengatur `MaxIdleConns(0)` menonaktifkan pemakaian ulang, membuat koneksi selalu baru. Sebaliknya, dengan pemakaian ulang, waktu eksekusi turun drastis.
- **Oversized Pool Exhaustion**: Klien membuka *pool* berkapasitas 50 sementara server diatur maksimal 15. Ketika 30 koneksi paralel dieksekusi, 15 akan berhasil dan 15 sisanya akan ditolak oleh *driver*.
- **Connection Leak**: Fungsi `ProcessOrderUnsafeLeak` meminjam koneksi lalu melakukan jeda I/O eksternal (misal: memanggil pihak ketiga). Karena *pool* dikunci kecil, *request* lain yang menunggu koneksi di fungsi `ProcessOrderSafe` akan mengalami kehabisan waktu (*context deadline exceeded*).

## What the Tests Prove
1. **Penalti Overhead**: *Pool* koneksi dengan *idle connections* lebih cepat (mikrodetik) dibandingkan membuka koneksi baru tiap saat (milidetik).
2. **Pool Exhaustion**: Database (lewat simulasi driver) dengan tegas akan memblokir dan menolak permintaan koneksi jika *pool* sisi klien terlalu besar dan melebihi batasan server.
3. **Starvation**: Melakukan panggilan eksternal yang lambat tanpa mengembalikan koneksi database akan secara absolut menyebabkan *pool* habis dan menyebabkan kegagalan pada operasi berurutan (terverifikasi dari pesan error *context deadline exceeded*).
4. **Validasi Skala Konkurensi**: *Pool* terbatas bisa melayani lebih banyak permintaan paralel secara aman jika koneksi segera dilepas pasca-*query*.

## Production Considerations
- Selalu pisahkan operasi I/O jaringan eksternal (seperti panggil API pembayaran) dari operasi yang membutuhkan transaksi / *lock* koneksi database.
- Hindari pengaturan *pool* sisi aplikasi yang terlalu besar. Gunakan pemantauan seperti metrik waktu akuisisi koneksi (*connection acquisition time*) dan metrik database seperti `pg_stat_activity` (perhatikan status `idle in transaction`).
- Untuk arsitektur terdistribusi, pertimbangkan *proxy pooling* tingkat infrastruktur.

## Warnings
- Efek latensi spesifik dari *dynamic scaling pool* atau *context switching* level kernel tidak disimulasikan sepenuhnya, lab disederhanakan melalui *mock memory structure*.
- Dampak alokasi `work_mem` riil pada batasan memori tidak diamati secara empiris dalam lab.

## Checklist
- [ ] Ukuran maksimal *pool* tidak melebihi CPU *core* server secara irasional.
- [ ] Operasi API jaringan tidak berada di tengah blok transaksi database.
- [ ] Limit `max_connections` server database sudah memperhitungkan skala horizontal klien.
- [ ] Terdapat limit *timeout* ketika aplikasi mencoba mengakuisisi koneksi.

## Key Takeaways
Silakan lihat dokumen `05-key-takeaways.md`.

## Sources
Silakan lihat dokumen `06-source-map.md`.