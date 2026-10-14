Lab 33 kali ini membahas anomali replication lag pada database read replicas, kegagalan pola pemisahan baca-tulis konvensional, serta empat strategi perutean cerdas untuk mempertahankan konsistensi read-your-own-writes.

Masalah teknis terjadi saat aplikasi mengarahkan kueri baca ke replika sekunder untuk skala beban. Karena replikasi basis data berjalan asinkron di latar belakang, terdapat selisih waktu antara commit di primary dengan penerapan data di replika. Pengguna yang baru saja memperbarui profil atau checkout pesanan dapat disajikan data usang atau menerima pesan data tidak ditemukan saat kueri masuk ke replika yang tertinggal.

Lab ini membuktikan empat hal konkret:
• Stale Read Anomaly: Kueri naif ke replika ber-lag 500ms langsung memicu ErrNotFound sesaat setelah transaksi tulis committed di primary.
• Time-Based Sticky Routing: Penahanan kueri baca ke primary selama jendela 300ms berhasil melindungi sesi sebelum beralih ke replika.
• Causal Token Routing: Perutean token LSN menunggu secara efisien via sync.Cond hingga replika mencapai posisi log yang dipersyaratkan.
• Synchronous Freshness: Replikasi sinkron menjamin ketersediaan data seketika pada replika dengan penalti waktu penulisan sebesar 102ms.

Mental model inti: Replication lag is the distance between primary commit and replica apply. Match your routing strategy to the cost of reading stale state.

Satu batasan penting: Lab ini menggunakan simulasi kluster in-memory dengan urutan LSN atomik, bukan pengujian basis data fisik multi-region, serta berfokus murni pada logika perutean dan sinkronisasi sesi.

Pelajari implementasi lengkap dan unit test di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/33-read-replicas-and-replication-lag

#softwareengineering #golang #database #distributeddatabase #systemdesign #backend