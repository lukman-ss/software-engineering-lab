# Key Takeaways: Backward Compatibility — Expand → Migrate → Contract Pattern

## 1. Pola Parallel Change Memecah Breaking Change Menjadi Tiga Tahap Non-Breaking
Setiap perubahan struktural dapat dibagi menjadi Expand (tambah baru bersama lama), Migrate (transfer data dan traffic secara gradual), dan Contract (hapus lama hanya ketika tidak ada yang menggunakannya lagi). Ketiga fase masing-masing independently deployable.

## 2. Dual-Write Menjamin Konsistensi Data Selama Transisi
Menulis ke kedua representasi (legacy dan modern) secara simultan dalam satu transaksi memastikan bahwa jika terjadi crash atau rollback, data tidak hilang di salah satu sisi. Dual-write mengorbankan sedikit latensi untuk keamanan data.

## 3. Backfill Harus Idempotent dan Resumable
Worker migrasi batch harus dapat dihentikan dan dilanjutkan tanpa menghasilkan duplikasi data. Cek duplikasi sebelum insert dan checkpoint (`last_processed_id`) adalah komponen penting untuk keandalan migrasi di produksi.

## 4. Fallback Read Mencegah Data Starvation
Mengalihkan bacaan ke skema baru sebelum backfill selesai akan menyebabkan data kosong untuk record yang belum terbackfill. Mode `ReadFallback` (baca modern dulu, fallback ke legacy jika kosong) adalah mekanisme penting selama masa transisi.

## 5. Kontrak Hanya Boleh Diterapkan Ketika Trafic Legacy = 0
Penerapan kontrak (drop kolom, hapus endpoint) harus dijaga oleh penghitung metrik. Jika masih ada klien lama yang mengakses, kontrak akan memutus akses mereka dan menyebabkan error.

## 6. Keamanan Rollback Hanya Berlaku Selama Fase Dual-Write
Jika aplikasi dirollback dari V2 ke V1 selama mode dual-write aktif, data tetap aman karena V1 masih menulis ke kolom legacy. Namun, jika dual-write telah dihentikan (WriteNewOnly) lalu rollback dilakukan, data yang ditulis dalam mode new-only akan hilang dari perspektif V1.

## 7. Feature Flags Memisahkan Deployment dari Aktivasi
Pengontrolan mode tulis dan baca via fitur flag memungkinkan organisasi untuk melakukan deploy kode baru ke produksi tanpa langsung mengaktifkan perilaku baru. Ini memungkinkan canary rollout 1% → 10% → 100% dan instant rollback tanpa perlu redeploy.

## 8. Header Deprecation dan Sunset Memberikan Pemberitahuan Formal kepada Klien
Header `Deprecation` (RFC 9224) dan `Sunset` (RFC 8594) memberikan waktu yang terukur kepada klien eksternal untuk migrasi. Setelah sunset date, endpoint mengembalikan `410 Gone`. Ini komunikasi teknis yang terstandarisasi (catatan: README menyatukan kedua RFC sebagai "RFC 8594" — nuansa kosmetik, bukan kesalahan fungsional).

## 9. Observasi dan Metrik adalah Panduan untuk Mengambil Keputusan Pada Setiap Fasa
Pengukuran traffic lama vs baru, jumlah backfill, dan deteksi drift data memberikan data nyata untuk menentukan kapan aman untuk berpindah dari satu fase ke fase berikutnya. Tanpa observasi, keputusan migrasi didasarkan pada perkiraan, bukan data.

## 10. Implementasi Demonstrasi ini Menggunakan Penyederhanaan yang Disadari
Penyimpanan dalam-memori (bukan PostgreSQL aktual) dan dual-write via mutex (bukan transaksi database sebenarnya) dipilih untuk menghilangkan dependensi infrastruktur luar sambil mempertahankan mekanika relasional yang akurat. Dalam produksi, penyimpanan harus diganti dengan `database/sql` atau ORM dan dual-write harus diimplementasikan dalam transaksi yang sesungguhnya atau menggunakan pola outbox.
