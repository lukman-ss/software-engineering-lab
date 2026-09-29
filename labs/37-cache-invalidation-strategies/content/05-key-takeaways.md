# Key Takeaways — Cache Invalidation Strategies

1. **Cache-aside adalah default pragmatis untuk read-heavy workload.** Pola lazy-load on miss hemat memori, tetapi read-after-write freshness tidak dijamin — writer harus menulis DB dulu, baru menghapus cache.

2. **Write-through memberikan freshness sinkron dengan harga throughput.** DB dan cache ditulis berurutan; reader berikutnya mendapat nilai terbaru tanpa query DB. "Same write operation" berarti urutan aplikasi, bukan transaksi ACID lintas sistem.

3. **Write-behind mengorbankan durability demi throughput tulis.** Cache di-update segera, DB diflush async. Jika process crash sebelum flush, write yang belum tereksekusi hilang. Queue overflow pada demo menggunakan silent drop — bukan pola produksi.

4. **Cache stampede adalah cascading failure yang terprediksi.** Pada expiry hot key di bawah high concurrency, banyak goroutine melakukan recompute sekaligus; load pada DB meningkat linear dan bisa menyebabkan congestion collapse hingga hit rate nol.

5. **Singleflight intra-process adalah mitigasi stampede paling sederhana.** `golang.org/x/sync/singleflight.Group.Do(key, fn)` menjamin satu eksekusi per key; waiter menerima hasil yang sama. Efektif untuk single-process Go service; multi-process deployment membutuhkan distributed lock.

6. **XFetch menggeser puncak rebuild ke forward-time secara probabilistik.** Formula `-Δ · β · ln(U) > TTL_remaining` (dengan U ~ Uniform(0,1)) memicu early refresh ketika traffic tinggi (Δ besar) atau dekat expiry. Implementasi harus menyertakan guard `u <= 0 || u >= 1`. Tanda minus pada logaritma sangat penting — versi tanpa minus (`Δ · β · ln(U)`) selalu false untuk U∈(0,1).

7. **Stale-while-revalidate (SWR) memanfaatkan window staleness.** Saat TTL expired tetapi masih dalam `staleDelta`, client menerima nilai stale segera sambil pemicu revalidasi async request-triggered. Deduplication per key mencegah revalidation storm. RFC 5861 berstatus Informational; penerapan di app-layer adalah analogi, bukan kepatuhan spesifikasi HTTP.

8. **Jitter deskronisasi TTL cross-key, bukan single-key stampede.** Menambahkan `[0, maxJitter)` pada base TTL mencegah banyak key expire serentak (mis. setelah deploy), tetapi tidak mengurangi jumlah concurrent miss pada satu hot key yang expired.

9. **Laboratorium ini pedagogis, bukan benchmark produksi.** In-memory cache dan `MockDB` digunakan agar lab dapat dijalankan standalone tanpa Redis/PostgreSQL. Beban 20 goroutines dan `queryDelay=20ms` adalah parameter sintetik; angka 10.000 RPS dalam konteks pedagogis bukan ukuran lapangan. ponytail: replace MockDB with Redis client for production; instrument Prometheus metrics.

10. **Observability wajib sebelum production.** Counter `queryCount` / `writeCount` / `revalCount` pada demo harus digantikan metrics sesungguhnya (Prometheus histogram/counter) untuk memantau pattern hit/miss, revalidation rate, dan queue depth write-behind.
