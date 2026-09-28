# Leader Election dengan Fencing Tokens

## Problem

Dalam sistem terdistribusi, hanya satu node yang boleh mengeksekusi operasi kritis (writer, job runner, migrator) pada satu waktu. Lease-based leader election menyelesaikan separuh masalahnya: lease memberikan otoritas berbatas waktu, node harus memperbaruinya secara berkala (heartbeat), dan bila node berhenti memperbarui, lease kedaluwarsa dan node lain dapat mengambil alih.

Masalah yang tersisa adalah **split-brain transient**. Seorang leader bisa mengalami stop-the-world GC pause atau network partition lebih lama dari lease TTL. Dalam jendela itu, coordinator menganggap lease telah kedaluwarsa dan memberikannya kepada kandidat lain. Node lama, saat terbangun, masih yakin ia memegang lease karena ia tidak menerima pemberitahuan apa pun — ia hanya berhenti karena pause. Dua node kini sama-sama percaya diri sebagai leader, dan tanpa mekanisme tambahan, tulisan dari leader lama dapat menimpa tulisan leader baru.

## Why This Matters

Skenario ini bukan teoretis. Pause GC panjang umum terjadi pada runtime dengan tracing garbage collector (JVM, .NET, dan sejenisnya). Research menunjukkan bahwa lease saja hanya memberikan *timed mutual exclusion* — cukup untuk optimisasi efisiensi (deduplikasi, rate limiting), tetapi tidak cukup untuk keamanan operasi yang tidak idempoten: transaksi finansial, pengurangan stok, migrasi skema database.

Fencing token menjembatani celah ini. Dengan fencing, tulisan yang datang dari pemegang lease yang sudah kedaluwarsa ditolak di lapisan penyimpanan berdasarkan nilai token, terlepas dari apakah pengirimnya masih merasa sah.

## Mental Model

Bayangkan pintu dengan kartu akses elektronik. Kartu hanya berlaku selama masa sewa (TTL). Jika Anda lupa memperpanjang, kartu mati dan orang lain mendapatkan kartu baru. Namun pemilik kartu lama tidak tahu kartunya sudah mati — ia hanya berhenti memperbarui. Ketika ia kembali dan mencoba membuka pintu dengan kartu lamanya, pintu menolak karena nomor kartu lama sudah kalah nomor dari kartu baru.

Fencing token adalah nomor kartu itu: nilai yang **selalu naik** setiap kali lease diberikan. Penyimpanan bersama menyimpan nomor terakhir yang diterima. Setiap tulisan harus menyertakan nomor kartunya; jika nomor itu `<=` nomor terakhir yang diterima, tulisan ditolak.

## Core Concept

Tiga komponen membentuk pola ini:

1. **Coordinator** — layanan koordinasi linearizable yang memberikan lease dengan TTL dan merilis fencing token yang monotonik naik (implementasi lab: field `revision` yang diincrement atomik di bawah mutex).

2. **Candidate** — node yang berkampanye untuk lease, memperbaruinya secara periodik lewat heartbeat (`KeepAlive` / `Renew`), dan mengeksekusi tulisan berfencing saat menjadi leader.

3. **Shared Resource dengan Fencing Check** — penyimpanan yang menolak tulisan dengan `incomingToken <= lastSeenToken` menggunakan check-and-set atomik.

Aturan penerimaan tulisan:

```go
if token <= s.lastSeenToken {
    return ErrStaleFencingToken
}
s.lastSeenToken = token
```

Ini adalah syarat ketat: token harus **lebih besar secara strict** dari token terakhir yang terlihat.

## Failure Scenario

Skenario yang didemonstrasikan di lab:

1. Node A memperoleh lease dengan fencing token 1.
2. Node A melakukan tulisan sukses dengan token 1.
3. Node A mengalami pause simulasi 350 ms — lebih besar dari TTL 200 ms.
4. Coordinator meniadakan lease Node A.
5. Node B berkampanye, memperoleh lease, dan menerima token 2.
6. Node B melakukan tulisan sukses dengan token 2 — diterima, `lastSeenToken` menjadi 2.
7. Node A terbangun, masih memegang token 1, dan mencoba menulis.
8. Shared storage menolak: `token 1 <= last seen 2`.

Tanpa fencing, langkah 7 akan menimpa tulisan Node B dan menyebabkan korupsi data.

Catatan penting dari audit: selama jendela pause, status dual-leader dapat terjadi secara transient (Node-A dan Node-B sama-sama berstatus LEADER). Ini adalah perilaku desain, bukan bug — keselamatan data tetap dijaga oleh fencing token.

## How It Works

### Lease dan Heartbeat

Kandidat memperoleh lease dengan `Acquire(key, holderID, ttl)`. Jika lease masih berlaku dan dipegang node lain, `Acquire` mengembalikan `ErrLeaseHeld`. Saat leader, kandidat memanggil `Renew` secara periodik (interval `renewRate`). Bila `Renew` gagal (lease kedaluwarsa atau pemegang berubah), kandidat kembali berstatus `FOLLOWER` dan `currentLease` di-set `nil`.

Masa berlaku lease dihitung dengan `time.Now()` (monotonic di Go), konsisten dengan best practice riset bahwa operasi durasi harus memakai monotonic clock untuk menghindari lompatan wall-clock (NTP step).

### Deteksi Gagal yang Terbatas

Lease TTL menyediakan *bounded failure detection*. Tidak perlu failure detector kompleks: bila leader berhenti renew sebelum TTL habis, lease otomatis kedaluwarsa dan node berikutnya dapat mengambil alih. Dalam demo, failover terjadi dalam skala TTL + renew interval.

### Fencing di Shared Storage

Penyimpanan menyimpan `lastSeenToken` dan menerapkan check-and-set atomik di bawah mutex. Tiga syarat kebenaran fencing dari riset dipenuhi di lab:

- **Monotonik**: token berasal dari counter `revision` yang hanya bertambah.
- **Validasi di resource layer**: perbandingan `token <= lastSeenToken` dilakukan di `FencedStorage.Write`.
- **Atomik**: perbandingan dan update berada dalam satu sekuens kunci mutex.

## Architecture

```text
┌─────────────────┐       ┌─────────────────┐
│ Candidate /     │       │ Candidate /     │
│ Node A          │       │ Node B          │
└────────┬────────┘       └────────┬────────┘
         │                         │
         │ Acquire / KeepAlive     │ Acquire / KeepAlive
         ▼                         ▼
  ┌─────────────────────────────────────────┐
  │         Coordinator (In-Memory)         │
  │ - Linearizable key-lease mapping        │
  │ - Monotonic revision / fencing token    │
  │ - Lease TTL & expiration monitoring     │
  └────────────────────┬────────────────────┘
                       │
                       │ Fenced Write (token, payload)
                       ▼
  ┌─────────────────────────────────────────┐
  │         Shared Resource (Storage)       │
  │ - Tracks last_seen_fencing_token        │
  │ - Rejects write if token <= last_seen   │
  │ - Accepts and updates state otherwise   │
  └─────────────────────────────────────────┘
```

## Implementation

Implementasi murni Go standard library, tanpa dependensi eksternal (etcd/Redis server tidak diperlukan untuk menjalankan test atau demo).

### internal/coordinator

`Coordinator` adalah struct thread-safe dengan `sync.Mutex`, map lease, dan counter `revision`. Tiga operasi utama:

- `Acquire` — mengembalikan lease bila kosong atau sudah kedaluwarsa; jika holder sama, memperpanjang expiry tanpa mengganti token; jika holder lain, `ErrLeaseHeld`. Saat lease baru dibuat, `revision++` dan token diambil dari `revision`.
- `Renew` — menolak bila lease tidak ada (`ErrLeaseExpired`), kedaluwarsa, atau pemanggil bukan pemegang dengan token yang cocok (`ErrNotLeaseOwner`).
- `Release` — menghapus lease hanya jika pemegang dan token cocok.

### internal/candidate

`Node` menjalankan election loop dalam goroutine. Setiap tick `renewRate`:

- Jika `now < pauseUntil`, loop melewati tick (simulasi GC pause / partition).
- Jika berstatus `LEADER`, memanggil `Renew`; kegagalan mengembalikan status ke `FOLLOWER`.
- Jika `FOLLOWER`, mencoba `Acquire`; keberhasilan menaikkan status ke `LEADER`.

`SimulatePause(duration)` mengatur `pauseUntil` ke waktu monotonik masa depan — inilah yang meniru stop-the-world pause atau disconnect jaringan.

### internal/storage

`FencedStorage.Write(author, token, value)` adalah satu-satunya gerbang tulisan. Ia mengembalikan error yang membungkus `ErrStaleFencingToken` dengan detail token dan author. `History()` mengembalikan salinan record untuk verifikasi.

### cmd/demo

Demo CLI menjalankan lifecycle lengkap: election → tulisan sah → simulasi pause 350 ms (> TTL 200 ms) → failover → tulisan Node baru → penolakan tulisan stale Node lama → cetak history penyimpanan.

## Code Walkthrough

Bagian kritis dari fencing check (sumber: `internal/storage/storage.go`):

```go
func (s *FencedStorage) Write(author string, token int64, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if token <= s.lastSeenToken {
		return fmt.Errorf("%w: token %d <= last seen %d (author: %s)", ErrStaleFencingToken, token, s.lastSeenToken, author)
	}

	s.lastSeenToken = token
	s.records = append(s.records, Record{
		FencingToken: token,
		Author:       author,
		Value:        value,
		Timestamp:    time.Now(),
	})
	return nil
}
```

Dan penerbitan token monotonik (sumber: `internal/coordinator/coordinator.go`):

```go
c.revision++
lease := &Lease{
	Key:          key,
	HolderID:     holderID,
	FencingToken: c.revision,
	ExpiresAt:    now.Add(ttl),
	TTL:          ttl,
}
```

Kedua potongan ini adalah inti pola: satu counter yang naik, satu gerbang yang membandingkan. Snippet lengkap dengan sumber file tercatat di `content/03-code-snippets.md`.

## What the Tests Prove

Empat test berjalan bersih (`go test -v ./...` dan `go test -race ./...`):

1. **`TestCoordinatorLeaseAcquisitionAndRenewal`** — akuisisi lease pertama menghasilkan token 1; akuisisi kedua oleh node lain menghasilkan `ErrLeaseHeld`; renew menjaga token tetap 1; setelah TTL lewat, renew menghasilkan `ErrLeaseExpired`; node kedua kemudian memperoleh token 2.

2. **`TestFencedStorageRejectsStaleTokens`** — token 10 lalu 11 diterima; pengulangan token 10 dan token duplikat 11 keduanya menghasilkan `ErrStaleFencingToken`; history berisi tepat 2 record.

3. **`TestLeaderElectionFailoverAndSplitBrainDefense`** — dua node; satu menjadi leader dan menulis; leader di-pause 250 ms (TTL 150 ms); kandidat kedua mengambil alih dengan token lebih besar; tulisan leader kedua sukses; tulisan leader lama dengan token lama ditolak `ErrStaleFencingToken`.

4. **`TestConcurrentElectionRace`** — lima node berkampanye bersamaan; setelah 100 ms, tepat satu node berstatus leader.

Apa yang **tidak** dibuktikan test ini: fault jaringan nyata, clock skew antar mesin, persistensi disk, dan fault Byzantine. Jangan mengklaim sebaliknya.

## Recovery / Rollback

Failover terjadi otomatis: ketika leader berhenti renew, lease kedaluwarsa, kandidat berikutnya memperolehnya dengan token lebih tinggi. Tidak ada intervensi manual pada skenario lab. `ForceStepDown` dan `Stop` tersedia untuk melepas lease secara sukarela, dengan `Release` hanya menghapus lease bila pemegang dan token cocok — token lama tidak bisa melepas lease milik pemegang baru.

## Production Considerations

Dari riset yang sudah di-approve (bukan dari pengujian lab ini):

- **Cocokkan jaminan dengan use case.** Lock efisiensi (deduplikasi, rate limiting) cukup dengan operasi atomik single-node. Operasi kritis-korektivitas memerlukan sistem linearizable (etcd, ZooKeeper) plus fencing di resource layer.
- **Fencing di resource layer, bukan di klien.** antirez dan Kleppmann sepakat bahwa fencing milik lapisan resource; lock service saja tidak cukup.
- **Gunakan monotonic clock** untuk semua perhitungan durasi terkait lease.
- **Atur timeout berdasarkan pengukuran**, bukan tebakan: GC pause p99 dan RTT antar node.
- **Hindari herd effect** saat failover dengan watch pada predecessor langsung.
- **Amati fencing rejection** sebagai sinyal dini pause atau masalah sistemik.

Nilai waktu di lab (TTL 200 ms, renew 50 ms) bersifat ilustratif untuk demo — bukan rekomendasi universal. Risiko research secara eksplisit menyatakan nilai seperti Raft election timeout 150–300 ms dan contoh TTL adalah angka berjangkar sumber, bukan rekomendasi umum.

## Common Mistakes

- **Menganggap lease saja cukup.** Lease memberikan timed mutual exclusion; tanpa fencing, leader yang di-pause masih bisa menulis.
- **Menerima token yang sama atau lebih kecil.** Aturan strict `incoming > lastSeen`. Accepting `token == lastSeen` membuka celah duplikasi.
- **Fencing hanya di klien.** Jika resource tidak memvalidasi token, fencing tidak ada artinya.
- **Token non-monotonik.** UUID atau token acak tidak bisa dipakai untuk fencing karena tidak memiliki ordering (riset: collision risk, no ordering).
- **Wall-clock untuk durasi.** Lompatan NTP membuat keputusan kedaluwarsa salah.

## Case Study

Riset mengulas empat sistem nyata untuk mengokohkan mental model:

- **etcd** — lease + `Revision` sebagai fencing token; revision monotonik berkat Raft.
- **ZooKeeper** — ephemeral sequential znode; leader = sequence terendah; `zxid` dan node version sebagai fencing token.
- **Consul** — Raft internal; session TTL plus health check.
- **Redis Redlock** — tidak menghasilkan fencing token monotonik; diperdebatkan Kleppmann vs antirez. Konsensus riset: Redlock tanpa fencing tidak menjamin mutual exclusion di bawah semua mode gagal realistis (lompatan jam, GC pause), tetapi dapat dipakai untuk efisiensi bila duplikasi dapat ditoleransi.

Debat Redlock adalah pengingat bahwa di mana token monotonik tidak tersedia dari lock service, fencing harus dibangun di resource layer.

## Checklist

- [ ] Lease diberikan dengan TTL yang diukur (bukan ditebak).
- [ ] Leader memperbarui lease sebelum TTL habis, dengan margin.
- [ ] Setiap kali lease diberikan, token monotonik naik.
- [ ] Setiap tulisan menyertakan token lease terakhir yang sah.
- [ ] Resource menolak `token <= lastSeenToken` secara atomik.
- [ ] Operasi durasi memakai monotonic clock.
- [ ] Fencing rejection di-log dan dimonitor.
- [ ] Keputusan timeout berbasis data pengukuran runtime.
- [ ] Pahami bahwa lease saja cukup untuk efisiensi, bukan untuk korektivitas.

## Key Takeaways

Lihat `content/05-key-takeaways.md`.

## Sources

Riset lab ini dikutip dari sumber terverifikasi (lihat `research/02-sources.md`):

1. Ongaro & Ousterhout, "In Search of an Understandable Consensus Algorithm", USENIX ATC '14 — https://raft.github.io/raft.pdf
2. Martin Kleppmann, "How to do distributed locking" (2016) — https://martin.kleppmann.com/2016/02/08/how-to-do-distributed-locking.html
3. Salvatore Sanfilippo (antirez), "Is Redlock safe?" (2016) — http://antirez.com/news/101
4. etcd v3.5 documentation, "How to conduct leader election" — https://etcd.io/docs/v3.5/tutorials/how-to-conduct-elections/
5. Apache ZooKeeper Recipes: Leader Election — https://zookeeper.apache.org/doc/current/recipes.html#sc_leaderElection
6. Martin Kleppmann, *Designing Data-Intensive Applications* — https://dataintensive.net/

Semua perilaku klaim pada artikel ini berasal dari implementasi dan test lab yang sudah di-approve; klaim konseptual berasal dari sumber di atas sebagaimana dicatat dalam `content/06-source-map.md`.
