Lab 39 kali ini membahas struktur data probabilistik Bloom filter, mekanisme eliminasi operasi I/O disk yang berlebihan pada sistem penyimpanan berbasis LSM-Tree, serta proteksi backend terhadap serangan cache penetration.

Masalah teknis yang kerap melumpuhkan arsitektur data adalah biaya pencarian kunci yang sebenarnya tidak ada. Pada basis data LSM-Tree, kueri data non-existent memaksa sistem membaca seluruh file SSTable segment di disk secara berurutan. Di lapisan caching, banjir permintaan untuk identifier yang tidak valid melompati memori cache dan langsung menghantam basis data hilir, memicu kelelahan connection pool dan lonjakan latensi p99.

Lab ini membuktikan empat hal konkret:
• Zero False Negatives: Pengujian 10.000 kunci membuktikan seluruh elemen yang terdaftar selalu dikembalikan true tanpa kegagalan.
• Presisi False Positive: Pengujian 50.000 kueri kunci absen menghasilkan tingkat error empiris 0.996 persen, konsisten terhadap target teoretis 1.00 persen.
• Reduksi 99 Persen I/O Disk: Simulasi 10.000 kueri pada 10 segment memangkas pembacaan disk dari 100.000 kali menjadi hanya 988 kali.
• Perlindungan Cache Mutlak: Pemasangan admission gate menahan 100 persen kueri penetrasi berbahaya sehingga 0 panggilan tembus ke database.

Mental model inti: Zero false negatives is a feature, not a bug. Trade certainty for compression and prune before you peek.

Satu batasan penting: Bloom filter standar tidak mendukung operasi penghapusan dan membutuhkan estimasi kapasitas di muka untuk mencegah saturasi bit array.

Pelajari implementasi lengkap dan unit test di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/39-bloom-filters

#softwareengineering #golang #database #performance #systemdesign #backend