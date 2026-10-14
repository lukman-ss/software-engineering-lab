Lab 22 kali ini membahas anti-pattern N+1 query problem, mekanisme terjadinya latensi akumuler pada lapisan akses data, serta mitigasi deterministik menggunakan eager loading berbasis batching.

Masalah teknis yang sering terjadi adalah aplikasi mengeksekusi satu query awal untuk mengambil koleksi data induk, lalu menjalankan satu query tambahan di dalam loop untuk setiap entitas anak yang terasosiasi. Pada dataset pengujian kecil, setiap query individu selesai dalam waktu sub-milidetik sehingga lolos tanpa terdeteksi oleh slow query log. Namun ketika volume data meningkat di lingkungan produksi, pola ini memicu ratusan round-trip jaringan, lonjakan overhead protokol basis data, serta potensi kehabisan connection pool.

Lab ini membuktikan empat hal konkret:
• Loop-Driven Anti-Pattern: Pengambilan relasi secara naif terbukti mengeksekusi 4 query untuk memuat 3 entitas author beserta postingannya.
• Batched Eager Loading: Pengumpulan seluruh identifier relasi ke dalam klausa IN tunggal memangkas eksekusi menjadi tepat 2 query konstan.
• Deep Equivalence: Struktur data hasil batching terbukti identik secara semantik dengan hasil metode naif melalui pengujian refleksi mendalam.
• Zero Race Condition: Pengujian konkuren dengan race detector membuktikan pembaruan query counter berjalan aman secara thread-safe.

Mental model inti: N+1 is a pattern, not a layer. Batch all related identifiers into a single IN query instead of executing per-row lookups.

Satu batasan penting: Lab ini adalah simulasi konseptual berbasis in-memory mock store di tingkat aplikasi, bukan pengujian terhadap basis data relasional fisik, serta berfokus murni pada query-level eager loading.

Pelajari implementasi lengkap dan unit test di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/22-n-plus-one-query-problem

#softwareengineering #golang #database #performance #systemdesign #backend