# Key Takeaways

1. **N+1 query problem** menyebabkan latensi akumuler tinggi karena satu query awal diikuti N query terpisah untuk setiap item, sering lolos dari slow query logs.

2. Fenomena **N+1 muncul di dua tempat**: pada lapisan database (ORM lazy loading default) dan lapisan jaringan (GraphQL resolvers, microservice API calls).

3. **Solusi utama**: *Request batching* atau *query-level eager loading* — kumpulkan semua ID, jalankan satu query `IN (...)`, gabungkan hasil di memori.

4. **Mapping-level eager loading** (JPA `FetchType.EAGER`, dll) adalah anti-pattern yang menyebabkan memory bloat karena semua relasi dimuat tanpa seleksi.

5. **Cartesian explosion** bisa terjadi pada eager loading via JOIN untuk relasi one-to-many/many-to-many — perlu `Result.unique()` atau strategi terpisah.

6. **Column selection** (`select`, `pluck`, `values`, `withCount`) sering lebih efisien daripada eager loading penuh ketika hanya nilai scalar dibutuhkan.

7. **Profiling-first workflow** (ukur → identifikasi pola query → eliminasi query tidak perlu) adalah best practice yang didokumentasikan semua ORM besar.

8. **Unit test dengan query counter** memverifikasi reduksi N+1 secara otomatis dan mencegah regresi.

9. Angka performa pada lab bersifat **illustratif**, tidak universal — tergantung volume data, schema, dan infrastruktur.

10. Lab ini hanya mendemonstrasikan **query-level eager loading** — tidak memuat mapping-level behavior atau kondisi OOM.
