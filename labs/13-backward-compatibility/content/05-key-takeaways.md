# Key Takeaways

1. **Paralel perubangan, bukan perubahan atom.** Pecahkan satu migrasi skema menjadi tiga fase terpisah (Expand → Migrate → Contract) agar setiap fase dapat di-deploy, diuji, dan di-rollback secara mandiri tanpa downtime.

2. **Skema harus kompatibel dengan N dan N+1 secara bersamaan.** Akibat rolling deployment, basis data tidak pernah berhenti melayani versi lama dan baru. Jangan pernah lakukan perubahan destruktif (`DROP COLUMN`) sampai kedua versi selesai deploy.

3. **Dual-write adalah jembatan, bukan solusi permanen.** Tulis ke skema lama dan baru secara atomik selama fase transisi. Gunakan fitur flag untuk canary rollout (`WriteDual`). Stop dual-write hanya setelah semua instance berada di versi baru.

4. **Backfill harus idempotent dan resumable.** Gunakan checkpoint (`last_processed_id`) sehingga jika proses terhenti, ia melanjutkan dari terakhir. Cek duplikat sebelum insert untuk mencegah replikasi data saat retry.

5. **Fallback read mencegah data starvation selama migrasi.** Baca dari skema baru; jika data belum ada, ambil dari skema lama dan larutkan secara lazy. Ini menghapus kebutuahan menunggu backfill selesai 100% sebelum beralih read path.

6. **Jangan hapus skema lama sampai traffic legacy = nol.** Gunakan observabilitas (counter, log, header) untuk memverifikasi. Fase kontrak hanya boleh diluncurkan setelah periode observasi bertahan nol traffic. Periode 30 hari adalah heuristik, bukan aturan universal.

7. **Rollback selama dual-write aman; rollback setelah dual-write berhenti berisiko kehilangan data.** Konsep dual-write memastikan data terbaru tersedia di kedua skema. Jika dual-write dihentikan, data baru hanya ada di skema baru dan tidak terlihat oleh aplikasi lama.

8. **API harus menjadi additive.** Tambah field baru dengan `omitempty`; jangan pernah ubah nama, tipe, atau hapus field yang sudah pernah dirilis ke klien eksternal yang tidak dapat dipaksa update.

9. **Gunakan header `Deprecation` dan `Sunset` (RFC 8594).** Sinyalkan ke klien lama bahwa endpoint/field sudah tidak didukung dan kapan akan dihapus. Ini memberi klien waktu untuk bermigrasi tanpa kejutan.

10. **Reconciliation mendeteksi drift sebelum kontrak.** Bandingkan nilai primary di skema baru dengan skema lama secara berkala. Drift mendeteksi dual-write yang gagal atomik sebelum fase kontrak menghapus jalur escape.
