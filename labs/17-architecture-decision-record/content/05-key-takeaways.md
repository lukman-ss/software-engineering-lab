# Key Takeaways

1. **Catat "Mengapa", Bukan Hanya "Apa":** Architecture Decision Records (ADR) berfokus pada motivasi, gaya dorong (*forces*), batasan, dan kompromi dari suatu keputusan, guna mencegah perdebatan berulang (*cyclical technical debates*) yang memperlambat tim.
2. **Kekekalan Rekaman (*Immutability*):** Keputusan yang telah disetujui (*Accepted*) tidak boleh diubah isinya untuk mencerminkan rancangan baru. Sejarah harus tetap utuh. Perubahan arsitektur mengharuskan penerbitan ADR baru.
3. **Simpan Usulan yang Ditolak (*Rejected*):** Jangan menghapus atau membuang usulan arsitektur yang tidak terpilih. Menyimpannya sebagai ADR berstatus `Rejected` memberikan dokumen bukti bagi tim masa depan tentang alasan sebuah opsi ditolak.
4. **Ko-lokasi dengan Kode Mengurangi *Documentation Drift*:** Dengan meletakkan berkas Markdown ADR berdampingan dengan direktori kode sumber dalam kontrol versi, siklus hidup keputusan berjalan paralel dengan evaluasi *pull request* kode sumber.
5. **Penomoran Harus Monotonik:** Identifikasi ADR wajib menggunakan bilangan bulat berurutan tanpa jeda atau pengulangan untuk memastikan tidak ada konteks sejarah yang hilang dari *repository*.
6. **Integritas Silsilah Timbal-Balik:** Validator struktural menjamin relasi penggantian harus konsisten di kedua belah pihak: ADR lama menunjuk penggantinya (`Superseded by <ID>`), dan ADR baru menyatakan keputusan yang ia gantikan (`Supersedes: <ID>`).
7. **Batas Keputusan Arsitektural:** Hanya keputusan yang secara signifikan memengaruhi struktur sistem, ketersediaan, dependensi, pola integrasi, atau batasan domain (*Architecturally Significant Requirements*) yang membutuhkan ADR, bukan untuk refaktorisasi internal level rendah.
8. **Trade-off Bersifat Kontekstual:** Pilihan desain tidak bersifat benar atau salah secara universal. Kasus uji seperti "Monolith vs Microservices" merupakan contoh penerapan yang sangat bergantung pada skala tim dan tahap umur produk saat keputusan dibuat.
