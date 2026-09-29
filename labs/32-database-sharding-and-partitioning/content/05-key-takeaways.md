# Key Takeaways

1. **Partitioning vs Sharding adalah Solusi Berbeda**: Logical Partitioning membagi tabel secara internal dalam satu database engine (optimal untuk query pruning dan fast drop partition tanpa VACUUM). Physical Sharding membagi dataset ke beberapa server independen untuk melampaui limit I/O, storage, dan CPU satu mesin.

2. **Shard Key Adalah Keputusan Terpenting**: Shard key yang optimal wajib memiliki kardinalitas tinggi (*high cardinality*), frekuensi rendah (*low frequency*), dan bersifat non-monotonik (seperti `user_id` atau `tenant_id`). Nilai berkardinalitas rendah (misal `continent`) membatasi jumlah maksimum shard.

3. **Hindari Monotonic Key Sebagai Range Shard Key**: Sharding key berbasis nilai yang bertambah secara sekuensial (`created_at`, auto_increment, ObjectId) memusatkan 100% write traffic pada satu shard aktif (*write hotspot*). Hashed Sharding adalah solusi standar untuk mendistribusikan write secara acak/merata.

4. **Consistent Hashing Meminimalkan Resharding Data Shuffle**: Saat menambah kapasitas node (misal dari $N$ ke $N+1$), Hash Modulo ($K \pmod N$) memaksa $\approx \frac{N}{N+1}$ ($\approx 75\% - 80\%$) key berpindah shard. Consistent Hashing dengan virtual nodes membatasi perpindahan key hanya sebesar $\approx \frac{1}{N+1}$ ($\approx 12\% - 20\%$).

5. **Gunakan Virtual Nodes untuk Distribusi Ring yang Seimbang**: Consistent Hashing murni tanpa virtual nodes rentan terhadap *skew* (distribusi tidak merata pada physical node yang sedikit). Menggunakan 50–150 virtual nodes per physical shard memastikan persebaran token di sepanjang ring hash merata.

6. **Query Non-Shard Key Membutuhkan GSI untuk Menghindari Scatter-Gather**: Query tanpa menyertakan shard key terpaksa melakukan broadcast paralel ke seluruh shard (*scatter-gather*). Menggunakan Global Secondary Index (Lookup Vindex) memungkinkan *point-lookup* terarah (1 node) dengan trade-off berupa biaya double-write saat insert/update.

7. **Auto-Increment Bawaan Database Gagal di Multi-Shard**: Sistem terdistribusi membutuhkan generator ID unik independen. RFC 9562 UUIDv7 memberikan ID 128-bit terurut waktu (*time-ordered*) yang ramah B-tree tanpa perlu koordinasi jaringan, sedangkan Sequence Block Allocation (Vitess-style) memberikan ID numerik berurutan dengan alokasi blok memori berkala.

8. **Two-Phase Commit (TwoPC) Mengorbankan Latensi dan Isolasi**: Transaksi terdistribusi lintas shard (TwoPC) menjamin atomisitas eksekusi namun memperkenalkan penalti latensi commit dan tidak menyediakan full ACID isolation tradisional (aplikasi berpotensi mengamati *fractured reads*). Sebisa mungkin co-locate data terkait dalam satu shard key yang sama.
