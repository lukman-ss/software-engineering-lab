# Read Replica dan Replication Lag: Strategi Menjaga Read-Your-Own-Writes

## Problem

Ketika arsitektur database menggunakan satu primary dan beberapa replica untuk skala baca, sebuah anomali umum muncul setelah user melakukan write: read berikutnya bisa memuat data lama atau bahkan melaporkan key tidak ditemukan. Anomali ini disebut **stale read** (read basi) dan terjadi karena write yang baru committing ke primary belum mengalir penuh ke replica — sebuah jarak waktu yang disebut **replication lag**.

Dalam dunia nyata, PostgreSQL streaming replication, MySQL semi-sync/async replication, AWS RDS/Aurora read replicas, Azure read replicas, dan MongoDB secondaries semuanya berjalan asynchronous secara default. Data pada replica ter-update dengan *eventual consistency* — bukan linearizable. Tanpa strategi aplikasi, user yang baru saja mengubah profilnya bisa disuguhi versi lama jika query dialihkan ke replica yang masih tertinggal.

## Why This Matters

Read-your-own-writes adalah keharusan UX yang sederhana: setelah kamu menyimpan perubahan, hal pertama yang kamu lihat setelah itu harus mencerminkan perubahan tersebut. Mengabaikan ini berakibat pada bug yang sulit dilacak — user menganggap sistem "patah", padahal yang terjadi adalah asinkroni replikasi yang tidak tertutup.

Dampak bisnisnya meluas:
- User tidak melihat order yang baru dibuat.
- Saldo rekening yang baru diupdate terbaca nol lagi.
- Feature flag terbaru tidak aktif di session saat itu.

Kunci masalahnya bukan sekadar "memilih replica yang tepat". Masalahnya lebih dalam: **bagaimana aplikasi melacak urutan kejadian write sehingga read berikutnya terjamin menyertakan write itu sendiri**.

## Mental Model

Bayangkan primary sebagai juru tulis yang mencatat setiap transaksi di buku besar dengan nomor urut — LSN (*Log Sequence Number*). Setiap tulis baru diberi LSN naik (1, 2, 3, …). Pada replica, ada worker yang membaca LSN satu per satu dan menerapkan perubahan ke data lokal.

```
Primary LSN:  [3] ← kamu menulis order:101
Replica 1 LSN: 2  ← belum sampai batch ke-3
Replica 2 LSN: 3  ← sudah apply batch ke-3
```

Jika read `order:101` dikirim ke Replica 1, hasil akan kosong atau lama. Jika ke Replica 2, hasilnya fresh. Permasalahannya menjadi pertanyaan operasional: **bagaimana aplikasi mengetahui mana replica yang sudah mencapai LSN tertentu, sebelum melayani read?**

## Core Concepts

Terdapat empat pendekatan yang telah teruji, masing-masing menyeimbangkan freshness, latensi, dan keterlibatan primary dengan cara berbeda.

### 1. Asynchronous vs Synchronous Replication

PostgreSQL membedakan dua mode dengan `synchronous_commit`:
- **Async (default)**: write commit secepatnya di primary; WAL stream dikerjakan worker secara background.
- **Sync (`remote_apply`)**: write tidak commit sampai semua replica telah menerapkan transaksi tersebut — menjamin freshness instan, tapi menelan biaya latensi write sebesar RTT × jumlah replica.

MySQL mengenal **semi-synchronous**: write menunggu replica menerima WAL (*relay log*), tapi tidak sampai apply — jadi durability terjamin, visibility belum. Vitess menyebutkan hal serupa secara eksplisit.

### 2. Time-Based Sticky Routing

Pendekatan heuristik: setelah write, alihkan semua read session tersebut ke primary selama jendela waktu `T_sticky`. Sesudah TTL habis, baru kembalikan ke pool replica. Default industri yang lazim dikutip adalah 5 detik, tapi ini konvensi, bukan jaminan engine.

### 3. Causal Token / Minimum LSN Routing

Pendekatan deterministik: setelah write, simpan LSN commit. Untuk read berikutnya, arahkan ke replica yang `appliedLSN >= minLSN`. Jika belum ada replica yang cukup fresh, tunggu dengan timeout, lalu fallback ke primary. Ini adalah implementasi langsung dari guarantee **read-your-own-writes** yang diformalkan oleh Terry et al. (Microsoft Research).

### 4. Lag-Aware Dynamic Routing & Primary Fallback

Untuk read-bebas (bukan read-your-own-writes), gunakan threshold ΔLSN (`MaxLSNDiff`). Hanya replica dengan selisih LSN ≤ threshold yang ikut rotasi beban; semua replica yang terlalu jauh dikecualikan. Jika semuanya terfilter, seluruh read jatuh ke primary — trade-off antara throughput dan konsistensi.

## Failure Scenario

User A melakukan `PUT /profile` dengan payload baru. Server melakukan write ke primary, mendapat LSN 47. Dalam 200 ms berikutnya, user A membuka dashboard yang menjalankan `SELECT * FROM profile WHERE id = A`. Load balancer pengarah read ke replica yang sedang ber-LSN 44. Query berhasil, tapi payload sudah kadaluarsa — atau Worse, key tidak ada di replica, sehingga aplikasi mengembalikan 404 ke user.

## How It Works — Arsitektur Lab

Lab ini mensimulasikan cluster primary-replica dalam memori Go, memodelkan WAL entry, LSN monotonik, dan dua mode replikasi. Di atasnya ada `Router` dengan empat method read:

| Method | Jaminan | Mekanisme |
|---|---|---|
| `ReadNaive` | none | Round-robin ke replica; bisa stale |
| `ReadWithStickySession` | session-level | Primary selama `T_sticky` sejak last-write session |
| `ReadWithToken(minLSN)` | causal token | Pilih replica ≥ minLSN; blocking wait + timeout; fallback primary |
| `ReadLagAware` | threshold-based | Filter replica dengan ΔLSN ≤ MaxLSNDiff; fallback primary jika kosong |

Komponen utama:
- `internal/cluster/cluster.go`: cluster in-memory, `Node` dengan data KV + `appliedLSN` + channel WAL, worker goroutine per replica, sync vs async path di `Cluster.Write`.
- `internal/router/router.go`: routing logic + session map (thread-safe `sync.Map`).
- `cmd/demo/main.go`: CLI end-to-end yang menjalankan keempat skenario.
- `tests/replication_test.go`: 7 test case yang memverifikasi setiap perilaku.

## Implementation

### Simulasi Cluster

Cluster diinisialisasi dengan jumlah replica dan default lag (durasi delay worker replica). Setiap replica menjalankan goroutine `replicaWorker` yang consume dari `walChannel` — channel dengan buffer 1024 entri. Saat write terjadi:

- Primary write + increment `currentLSN` dilakukan secara atomic.
- Pada **async mode**, WAL entry diseleksi ke semua replica via channel; worker replica menunggu `lagDuration` sebelum apply.
- Pada **sync mode**, write diblokir sampai semua replica diterapkan secara sinkron dalam loop.

### Sticky Session

Router menyimpan di `sync.Map` berkey sessionID:
```go
type SessionState struct {
    LastWriteTime time.Time
    LastWriteLSN  uint64
}
```
Saat `ReadWithStickySession` dipanggil, router memeriksa apakah waktu sekarang masih dalam `StickyDuration` sejak `LastWriteTime`; jika ya, read diarahkan ke primary tanpa kecuali.

### LSN Token Wait

`ReadWithToken` mengiterasi replica, memilih pertama yang `AppliedLSN() >= minLSN`. Jika tidak ada, ia memanggil `WaitForLSN(ctx, minLSN)` pada replica pertama — utilitas berbasis `sync.Cond` yang menunggu until kondisi terpenuhi atau context timeout. Setelah waktu tunggu habis, fallback ke primary.

### Lag-Aware Filter

`ReadLagAware` menghitung `diff = primaryLSN - replicaAppliedLSN` untuk setiap replica. Hanya yang `diff ≤ MaxLSNDiff` yang masuk pool `eligible`. Jika kosong, read langsung ke primary dengan label `(fallback-lag)`.

## Code Walkthrough

Berikut cuplikan inti dari setiap method routing.

### Sticky Routing
```go
func (r *Router) ReadWithStickySession(sessionID, key string) (string, string, uint64, error) {
    if stateVal, ok := r.sessions.Load(sessionID); ok {
        state := stateVal.(SessionState)
        if time.Since(state.LastWriteTime) < r.config.StickyDuration {
            val, lsn, err := r.cluster.Primary().Read(key)
            return val, r.cluster.Primary().ID(), lsn, err
        }
    }
    return r.ReadLagAware(key)
}
```
Setelah TTL habis, flow jatuh ke `ReadLagAware` — komposisi dua strategi.

### LSN Token Wait
```go
func (r *Router) ReadWithToken(ctx context.Context, minLSN uint64, key string) (string, string, uint64, error) {
    for _, replica := range r.cluster.Replicas() {
        if replica.AppliedLSN() >= minLSN {
            return replica.Read(key)
        }
    }
    // blocking wait dengan timeout, fallback ke primary
}
```

### Async vs Sync di sisi Cluster
```go
if c.replMode == SyncReplication {
    for _, replica := range c.replicas {
        if lag > 0 { time.Sleep(lag) }
        replica.data[key] = value; replica.appliedLSN = lsn
    }
} else {
    for _, replica := range c.replicas {
        select {
        case replica.walChannel <- entry:
        default: // drop jika penuh
        }
    }
}
```

## What the Tests Prove

| Test | Bukti |
|---|---|
| `TestNaiveReplicationLag_StaleRead` | Baca round-robin ke replica dengan lag 500ms setelah write langsung menghasilkan `ErrNotFound`. Setelah 600ms, baca berhasil. |
| `TestStickySessionRouting` | Baca dalam 300ms window sticky → primary. Baca session lain → stale. Baca setelah 550ms → replica, fresh. |
| `TestReadWithToken_LSN` | `ReadWithToken` menunggu ~200ms (WAL catch-up) sebelum mengembalikan nilai dari replica. `readLSN >= lsn` terbukti. |
| `TestReplicaLagThreshold_Fallback` | Replica di-set lag 10 detik; `ReadLagAware` mengembalikan nilai dari `primary (fallback-lag)` bukan replica. |
| `TestSynchronousReplication_Freshness` | Setelah write sync, `ReadNaive` langsung menemukan nilai di replica; `readLSN >= lsn` terbukti. |
| `TestConcurrentAccess_RaceFree` | 15 goroutine writer+reader × 20 operasi = 300 operasi konkuren; `go test -race` tidak mendeteksi race. |
| `TestWaitForLSN_ContextTimeout` | `WaitForLSN` dengan context timeout 50ms dan target LSN jauh (999) mengembalikan `context.DeadlineExceeded`. |

Demo CLI additionally menampilkan:
- Write primary, naive read replica → key not found.
- Sticky read dalam 500ms → primary, fresh.
- Sticky read setelah 600ms → replica, fresh.
- Token read dengan LSN = 2 → replica, fresh setelah menunggu catch-up.
- Sync write + naive read → replica fresh, dengan write duration ~100ms.

## Recovery / Rollback

Jika read dari replica menghasilkan stale data:
- **Sticky session**: ulang read dalam jendela sticky, atau perpanjang `StickyDuration`.
- **Token wait**: perluas `WaitTimeout` jika replica memang lambat catching up; atau perkecil `MaxLSNDiff` untuk filter lebih ketat.
- **Fallback**: read selalu aman ke primary — tapi pergunakan selektif karena menekan primary.

Tidak ada operasi rollback data: lab adalah simulasi, bukan transaksi database nyata.

## Production Considerations

**Keuntungan pendekatan lab:**
- Sticky routing sangat murah secara komputasi (cuma cek waktu).
- Token/LSN routing adalah jaminan causal yang tepat secara matematis.
- Lag-aware routing memberi perlindungan graceful degradation, bukan fail total.

**Kekurangan / trade-off:**
- Sticky window 5 detik (default kode `StickyDuration`) bersifat heuristik. Di lingkungan dengan latency tinggi atau traffic spike, jendela bisa terlalu pendek (stale read berulang) atau terlalu panjang (primary terbebani). Demo memakai 500ms khusus agar deterministik — bukan default kode.
- Memanggil `pg_last_wal_receive_lsn` / `SHOW REPLICA STATUS` per-query menambah round-trip; alternatif yang lebih baik adalah caching heartbeat LSN di connection pool.
- Session map di router adalah state yang perlu didistribusikan pada arsitektur stateless horizontal scaling;JWT atau Redis perlu dipakai di produksi.
- Sync replication meningkatkan write latency sebesar RTT × N (N = jumlah replica). Tidak cocok untuk write-heavy workload.

**Rekomendasi opsional (di luar scope lab):**
- Gunakan middleware sudah teruji: GORM DBResolver, Vitess VTGate, ProxySQL.
- Pantau metrik native: `pg_last_wal_receive_lsn` (PostgreSQL), `Seconds_Behind_Source` (MySQL), `Read Replica Lag` (AWS RDS).
- Set alert threshold berdasarkan p99 replica lag, bukan rata-rata.

## Common Mistakes

1. **Mencampuradukkan semi-sync dengan synchronous.** Semi-sync menjamin durability (replica punya relay log), tapi bukan visibility — read dari replica bisa tetap stale.
2. **Mempercayai "replica selalu baru" tanpa cek lag.** Asumsi ini salah di semua database relational populer.
3. **Menggunakan jendela sticky terlalu pendek (< 1s).** Bisa gagal jika replica memang tertinggal 2–3 detik.
4. **Tidak menyiapkan fallback primary.** Kalau semua replica ter-filter oleh lag, read akan error jika tidak ada fallback.
5. **Menganggap konfigurasi default cukup.** `MaxLSNDiff=5`, `WaitTimeout=1s` adalah angka default di lab — perlu penyesuaian terhadap latensi jaringan produksi.

## Case Study — Demo End-to-End

CLI demo dijalankan dengan konfigurasi: 2 replica, default lag 200ms, sticky 500ms, waitTimeout 300ms.

**Skenario 1 — Anomali async lag:**
```
[Write Primary] key='user:profile:123', LSN=1
[Naive Read]   replica-2, appliedLSN=0, val='', Err: key not found
```
Read langsung setelah write ke replica ber-LSN 0 menghasilkan key tidak ditemukan — bukti langsung stale-read anomaly.

**Skenario 2 — Sticky session:**
```
[Sticky Read (within 500ms)]   primary, LSN=1, val='{"name":"Alice","tier":"premium"}', Err: <nil>
[Sticky Read (after TTL)]      replica-1, LSN=1, val='{"name":"Alice","tier":"premium"}', Err: <nil>
```
Dalam jendela sticky, read dijemput primary (fresh). Setelah 600ms, replica sudah catch-up → read ke replica tetap fresh.

**Skenario 3 — Causal token:**
```
[Write Primary] key='order:101', LSN=2
[Token Read (MinLSN=2)] replica-1, LSN=2, val='{"status":"completed"}', Err: <nil>
```
`ReadWithToken` menunggu hingga replica-1 mencapai LSN 2, lalu melayani read dengan fresh value.

**Skenario 4 — Sync replication:**
```
[Sync Write] LSN=1, Write Duration: 104ms
[Sync Naive Read] replica-2, LSN=1, val='{"name":"Bob"}', Err: <nil>
```
Write sync memakan 104ms (termasuk propagasi ke 2 replica), tapi setelah itu `ReadNaive` langsung mendapatkan nilai fresh di replica.

## Checklist

Gunakan checklist ini saat mengevaluasi kebutuhan mitigasi replication lag di proyek Anda:

- [ ] Apakah write Anda meng-update key yang segera dibaca lagi dalam session yang sama? → perlu sticky/token.
- [ ] Berapa p99 replication lag replica Anda saat traffic normal? → gunakan metrik native database.
- [ ] Apakah Anda memiliki fallback primary bila semua replica exceeded threshold? → wajib.
- [ ] Apakah state session router dapat bertahan pada scaling horizontal? → pertimbangkan Redis/JWT.
- [ ] Apakah latency write tambahan dari sync replication dapat diterima workload Anda? → benchmark dahulu.
- [ ] Apakah monitoring alert sudah dikonfigurasi untuk replica lag SLA breach? → minimal p99.

## Key Takeaways

1. Async replication adalah default di semua database populer; read dari replica tanpa strategi = potensi stale read.
2. Read-your-own-writes dijamin hanya dengan membaca primary ATAU membaca replica yang sudah menerapkan write tersebut.
3. Sticky routing (time-based) sederhana dan murah; token/LSN routing (causal) presisi secara matematis.
4. Lag-aware filtering + primary fallback memberikan perlindungan graceful degradation.
5. Synchronous replication menghilangkan lag tapi menaikkan write latency secara linear terhadap jumlah replica.
6. Default konfigurasi di lab (5s sticky, 5 LSN diff, 1s wait) adalah titik awal — sesuaikan dengan peta latensi produksi Anda.
7. Semua strategi yang dibahas dapat dikombinasikan: contoh — sticky untuk critical session reads, lag-aware untuk generic reads.

## Sources

Riset merujuk 10 sumber resmi, diringkas sebagai berikut:

| # | Sumber | Topik |
|---|---|---|
| 1 | PostgreSQL 18 Docs — Log-Shipping Standby | Async streaming replication, LSN lag monitoring, hot standby |
| 2 | AWS RDS — Read Replicas | Read-only replica, metrics lag, promotion |
| 3 | Azure PostgreSQL Flexible Server — Read Replicas | Async physical replication, lag metrics (bytes/seconds) |
| 4 | MongoDB — Read Isolation, Consistency, Recency | Causal sessions, read-your-writes guarantee |
| 5 | AWS Aurora — DB Clusters | Storage-compute separation, shared-volume replica |
| 6 | PostgreSQL 18 — Runtime Config (Replication) | `synchronous_commit`, `hot_standby`, `recovery_min_apply_delay` |
| 7 | Vitess — MySQL Replication | Async/semi-sync, GTID, semi-sync ≠ guarantee freshness |
| 8 | GORM DBResolver | Application-level read/write split middleware |
| 9 | Kleppmann — "Please stop calling databases CP or AP" | Linearizability, async follower staleness |
| 10 | Terry et al. — "Replicated Data Consistency Explained Through Baseball" | Formalisasi session guarantees (read-your-writes, monotonic reads) |

Catatan audit penelitian (Gaps 1–3) dicatat sebagai warning: jendela sticky 5s bersifat heuristik; overhead polling LSN per-query perlu perhatian; distribusi lag multi-region belum terkuantifikasi publik.
