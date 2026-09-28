# Content Brief

Topic: Read Replicas and Replication Lag — Menjaga read-your-own-writes saat membaca dari replica

Target Reader: Backend engineer / DB engineer yang sudah memahami konsep primary-replica database dan ingin memahami strategi mitigasi replication lag secara konkret.

Problem: Ketika aplikasi mengarahkan read ke replica untuk skala baca, user bisa melihat data basi (stale read) setelah menulis data sendiri. Replication lag membuat read-your-own-writes tidak terjamin.

Core Mental Model: Replication lag adalah jarak antara LSN primary dan applied-LSN replica. Setiap strategi mitigasi (sticky routing, causal token, lag-aware routing, synchronous replication) adalah cara menutup jarak itu: arahkan read ke primary, tunggu replica catch-up, atau batasi read hanya ke replica yang sudah cukup segar.

Approved Research Status: APPROVED

Approved Engineering Status: APPROVED

Main Concepts:
- Asynchronous vs Synchronous Replication (`remote_apply`)
- Stale Read Detection & Mitigation
- Time-Based Sticky Routing (T_sticky)
- Causal Token / Minimum LSN Routing (blocking wait untuk catch-up)
- Lag-Aware Dynamic Routing & Primary Fallback

Verified Behaviors:
- TestNaiveReplicationLag_StaleRead: read naive ke replica lagging menghasilkan `ErrNotFound` (stale read), lalu berhasil setelah lag terlewati.
- TestStickySessionRouting: read dalam jendela sticky diarahkan ke primary; setelah TTL habis, read jatuh ke replica.
- TestReadWithToken_LSN: `ReadWithToken` menunggu ~200ms (WAL catch-up) sebelum melayani read dari replica dengan `readLSN >= lsn`.
- TestReplicaLagThreshold_Fallback: replica dengan ΔLSN > `MaxLSNDiff` dikecualikan; read jatuh ke `primary (fallback-lag)`.
- TestSynchronousReplication_Freshness: write sync langsung terlihat di replica (`readLSN >= lsn`).
- TestConcurrentAccess_RaceFree: 15 goroutine × 20 operasi lolos `go test -race`.
- TestWaitForLSN_ContextTimeout: `WaitForLSN` mengembalikan `context.DeadlineExceeded` saat timeout.

Available Case Studies: Demo CLI (`cmd/demo`) yang menampilkan 4 skenario: anomaly async lag, sticky session read-your-own-writes, causal token wait, dan sync replication freshness-vs-latency.

Warnings:
- Jendela sticky 5 detik adalah konvensi praktis / default konfigurabel, BUKAN jaminan engine (research-audit Gap 1).
- Cek LSN per-query menambah round-trip overhead (research-audit Gap 2, open question 3).
- p99 lag multi-region tidak tersedia di dokumentasi publik (research-audit Gap 3) — di luar scope lab.
- Lab memakai simulasi cluster in-memory (key-value), bukan database SQL produksi; tidak ada persistence, tidak ada topologi cascading, tidak ada split-brain election.
- Demo memakai durasi pendek (sticky 500ms, lag 200ms) untuk kepentingan deterministik — angka ilustratif, bukan rekomendasi universal.
