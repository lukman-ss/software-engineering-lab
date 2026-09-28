1. Slow dependencies (high latency) lebih berbahaya daripada failed dependencies karena Little's Law (L = λW): in-flight concurrency meningkat proporsional dengan latency.

2. Timeout harus dibuat berdasarkan distribusi P95/P99 latency dan budget end-to-end, bukan angka arbitrer besar. Timeout besar mengubah bimodal latency menjadi system-wide outage.

3. Full Jitter (`sleep = random(0, min(cap, base×2^attempt))`) menurunkan server contention lebih dari 50% dengan mendistribusikan retry secara uniform, mencegah thundering herd.

4. Timeout adalah *state indeterminate*, bukan kegagalan. Retry non-idempotent operations tanpa idempotency key berisiko double charge — server mungkin sudah memproses request sebelum response hilang.

5. Context deadline diwariskan ke child: child timeout ≤ parent timeout. Deadline propagation mencegah downstream melakukan work yang sudah tidak ada yang menunggu.

6. Circuit breaker CLOSED → OPEN → HALF_OPEN → CLOSED (atau kembali OPEN jika gagal) menghentikan retry chain ke downstream yang diketahui gagal.

7. Idempotency store dengan TTL dan lazy eviction menggunakan write lock saat Get untuk melakukan mutasi map — thread-safe dan aman diakses dari banyak goroutine.

8. Database layer perlu `statement_timeout`, `lock_timeout`, `idle_in_transaction_session_timeout` karena HTTP timeout tidak melindungi query yang berjalan tanpa batas.

9. Queue worker perlu execution timeout, max retries, dan dead-letter queue (DLQ) untuk mencegah infinite loop atau hung worker mengunci thread secara permanen.

10. Lab tidak mengimplementasikan koordinasi antar-proses (header `grpc-timeout` atau HTTP deadline propagation), in-memory idempotency store tidak persisten, dan adaptive timeout masih open question.