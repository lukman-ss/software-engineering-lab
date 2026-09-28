# Key Takeaways

1. **Lease Memberikan Timed Mutual Exclusion, Bukan Safety Mutlak**  
   Lease membatasi otoritas berdasarkan waktu (TTL), namun tidak mencegah leader yang mengalami pause panjang (stop-the-world GC, network delay) untuk bangun dan mengeksekusi operasi dengan asumsi masih sah.

2. **Fencing Token Mengubah Otoritas Temporal Menjadi Validasi Monotonik**  
   Setiap lease baru harus menerbitkan token yang strictly monotonically increasing. Hal ini memungkinkan resource layer menentukan urutan otoritas secara independen dari jam atau status lokal klien.

3. **Gerbang Fencing Berada di Shared Resource, Bukan di Klien**  
   Klien yang stale tidak tahu bahwa dirinya sudah tidak sah. Hanya shared resource (database, storage layer) yang dapat menolak tulisan dengan aturan check-and-set atomik: `incomingToken <= lastSeenToken → REJECT`.

4. **Transient Dual-Leader Dapat Terjadi Saat Pause Window**  
   Node lama dan node baru bisa sama-sama berstatus LEADER secara bersamaan jika node lama belum sadar lease-nya kedaluwarsa. Ini bukan bug, melainkan karakteristik sistem terdistribusi asinkron; integritas data dijamin oleh penolakan fencing di storage layer.

5. **Gunakan Monotonic Clock untuk Semua Logika Durasi**  
   Pengecekan lease TTL dan interval renewal wajib menggunakan monotonic clock (`time.Now()` di Go, `clock_gettime(CLOCK_MONOTONIC)` di C) guna mencegah anomali akibat lompatan wall-clock (NTP step).

6. **Bedakan Kebutuhan: Efisiensi vs Korektivitas**  
   Untuk optimisasi efisiensi (deduplikasi job, pencegahan pengiriman email ganda), lock sederhana (misalnya Redis `SET NX EX`) sudah memadai. Untuk operasi korektivitas kritis (transaksi finansial, mutasi inventori, migrasi skema), wajib menggunakan sistem strongly-consistent (etcd, ZooKeeper) dengan fencing token.

7. **Nilai Waktu Lab Bersifat Ilustratif**  
   Parameter TTL 200 ms dan renew 50 ms pada lab ini digunakan untuk kecepatan eksekusi demo dan unit test. Pada lingkungan produksi, TTL dan interval harus dikalibrasi berdasarkan p99 durasi GC pause dan round-trip time (RTT) jaringan nyata.
