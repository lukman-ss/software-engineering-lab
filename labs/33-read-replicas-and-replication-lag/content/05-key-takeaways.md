# Key Takeaways

1. **Replication lag adalah asinkronitas default.** PostgreSQL, MySQL, MongoDB, AWS RDS, dan Azure semuanya menjalankan async replication sebagai default. Tanpa strategi aplikasi, read dari replica bisa basi.

2. **Read-your-own-writes hanya terjamin dari primary ATAU replica yang sudah catch-up.** Guarantee ini tidak bisa diandalkan dari replica secara acak — perlu koordinasi eksplisit di level aplikasi.

3. **Sticky routing (time-based) adalah approach paling sederhana.** Simpan `lastWriteTime` per session; baca ke primary selama `T_sticky`. Tanpa dependensi LSN query. Default 5 detik bersifat heuristik — perlu tuning.

4. **Causal token (minLSN) adalah pendekatan paling presisi.** `ReadWithToken` secara deterministik menunggu atau memilih replica yang sudah apply write. Jaminan matematis, tapi butuh `WaitForLSN` yang menambah latensi.

5. **Lag-aware filtering memberikan graceful degradation.** Threshold ΔLSN memfilter replica terlalu tertinggal. Jika semua ter-filter, read jatuh ke primary — degradasi terukur, bukan kegagalan total.

6. **Synchronous replication menghilangkan lag dengan harga latency write.** Write diblokir sampai semua replica menerapkan entry. Cocok untuk freshness, tidak cocok untuk write-heavy workload.

7. **State session perlu didistribusikan untuk horizontal scaling.** Router state disimpan di `sync.Map` (single process). Di produksi stateless, gunakan Redis atau JWT untuk menyimpan sticky timestamp / LSN token.

8. **Jangan mencampuradukkan semi-sync dengan synchronous.** MySQL semi-sync menjamin durability (relay log diterima) tapi bukan visibility (belum apply). `remote_apply` PostgreSQL adalah yang menjamin visibility.

9. **Monitoring replica lag adalah kewajiban opsional, bukan sekadar keinginan.** AWS `Read Replica Lag`, Azure `physical_replication_delay_in_seconds`, dan PostgreSQL `pg_last_wal_receive_lsn` harus dipantau dan di-alert.

10. **Lab ini menggunakan simulasi in-memory — angka numerik bersifat ilustratif.** Sticky 500ms, lag 200ms, write latency 104ms bukan benchmark produksi. Tujuannya adalah memverifikasi perilaku (stale read terjadi → mitigate), bukan memvalidasi performa.
