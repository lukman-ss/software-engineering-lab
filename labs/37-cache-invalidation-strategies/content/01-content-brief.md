# Content Brief

**Topic:** Cache Invalidation Strategies — patterns untuk menulis-cache dan mitigasi cache stampede (thundering herd) di backend.

**Target Reader:** Software engineer yang membangun atau memelihara backend Go dengan cache in-process (termasuk Redis sebagai target akhir); pembaca harus sudah paham dasar HTTP caching dan concurrency Go.

**Problem:** Saat item cache populer kadaluarsa serempak, banyak goroutine/klien secara bersamaan melakukan recompute dari database — menyebabkan database kewalahan (congestion collapse), latency melonjak, dan hit rate turun mendekati nol. Pola tulis (write-through vs write-behind vs cache-aside) juga memilih antara freshnes read-after-write dan throughput tulis.

**Core Mental Model:**
1. Cache-aside = "minta dulu, simpan kalau miss" — lazy loading, invalidasi eksplisit setelah write.
2. Write-through = sinkron DB + cache — freshness tinggi, latency tulis lebih mahal.
3. Write-behind = tulis cache dulu, flush DB async — throughput tinggi, risiko data loss sebelum flush.
4. Stampede mitigation = singleflight (coalesce), probabilistic early expiration (XFetch), stale-while-revalidate (SWR), jitter pada TTL.

**Approved Research Status:** APPROVED (research-audit/07-verdict.md, date 2026-09-28).

**Approved Engineering Status:** APPROVED (engineering-audit/06-verdict.md dan engineering-audit-opensource/06-verdict.md, date 2026-09-29).

**Main Concepts:**
- Cache-Aside, Write-Through, Write-Behind / Write-Back
- Cache Stampede / Thundering Herd (cache-miss dog-piling)
- Single-flight coalescing via `golang.org/x/sync/singleflight`
- XFetch (probabilistic early expiration, formula `-Δ · β · ln(U) > TTL_remaining`)
- Stale-While-Revalidate (SWR, RFC 5861 analog)
- TTL Jitter

**Verified Behaviors (dari demo + tests):**
- Cache-aside: miss → DB query; hit → return cached; write DB → delete cache; next read reloads DB.
- Write-through: write DB + update cache secara berurutan; subsequent reads hit cache tanpa DB query.
- Write-behind: write cache + enqueue async flush; immediate reads hit cache; unflushed writes hilang jika process crash sebelum flush worker selesai.
- Stampede: 20 goroutine konkuren pada key expired menyebabkan 20 DB queries di path naive; dengan singleflight menjadi 1 query dan semua goroutine mendapat nilai yang sama.
- XFetch: formula `-delta * beta * math.Log(u) > ttlRemaining` benar secara matematis (u ∈ (0,1)); guard mencegah u ≤ 0 atau u ≥ 1.
- SWR: saat TTL expired tapi masih dalam stale window, client menerima value stale segera; revalidate dijalankan async; fresh value tersedia setelah revalidasi selesai.
- Jitter: `base + rand[0, maxJitter)` — memastikan TTL tidak sinkron across keys.

**Available Case Studies:**
- Demonstrasi synthetic load 20 goroutine konkuren dengan `queryDelay = 20 ms` menunjukkan difference query count 20 vs 1 (src: `tests/cache_test.go:109-158`, `cmd/demo/main.go:80-126`).
- XFetch simulated deterministic rand draw (`1e-9` trigger early refresh; `0.9999` no early refresh) — src: `cmd/demo/main.go:128-154`.
- SWR lifecycle 20ms TTL + 300ms stale window — src: `cmd/demo/main.go:156-176`, `tests/cache_test.go:186-221`.

**Warnings:**
- Lab menggunakan in-memory cache dan `MockDB` — bukan Redis/PostgreSQL asli; hasilnya bersifat pedagogis, bukan produksi.
- Singleflight ini proses-lokal; multi-process deployment membutuhkan distributed lock.
- Write-behind queue overflow diamkan (drop writes) — bukan rekomendasi produksi; production butuh WAL/Kafka.
- XFetch optimality proof belum diverifikasi dari paper primer (PDF unparsable) — algoritma diterima dari Wikipedia summary dan DOI metadata.
- Demo workload (20 goroutines, 20ms delay) dan klaim "10.000 RPS / P99 comparison" bersifat synthetik — bukan benchmark industri.
- RFC 5861 berstatus *Informational* (Independent Submission), bukan standards-track.
