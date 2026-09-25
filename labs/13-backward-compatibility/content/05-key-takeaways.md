# Key Takeaways

1. **Parallel Change (Expand-Migrate-Contract)**: Pisahkan perubahan schema dan contract API menjadi tiga fase independen untuk menjamin zero downtime dan keandalan sistem.
2. **Backward Compatibility API**: Tambahkan field baru tanpa menghilangkan field lama. Client lama tidak akan rusak jika endpoint tetap mengirimkan format yang diandalkannya.
3. **Dual-Write Consistency**: Dalam masa transisi, pastikan operasi write mengubah representasi lama dan baru secara atomik.
4. **Fallback Read Mechanism**: Terapkan logika fallback saat membaca data dari schema baru untuk melindungi integritas respons jika data historis belum tersinkronisasi (belum selesai di-backfill).
5. **Idempotent Backfill**: Proses sinkronisasi data harus dibuat idempoten dan batched agar dapat di-resume dengan aman di tengah jalan tanpa menimbulkan duplikasi data atau me-lock database secara masif.
6. **Observability adalah Syarat Kontrak**: Jangan pernah menjalankan penghapusan schema atau interface lama tanpa observabilitas yang memastikan traffic ke versi lama telah mencapai nilai nol.
7. **Instant Rollback Safety**: Jika sistem di-rollback saat migrasi berjalan, ketersediaan data di schema lama memastikan fungsionalitas berjalan normal tanpa data loss.
