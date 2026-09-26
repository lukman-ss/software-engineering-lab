# Key Takeaways

1. N+1 menyebabkan latensi agregat tinggi dan sering lolos dari *slow query logs*.
2. Fenomena N+1 terjadi pada ORM (SQL) dan lapisan API jaringan (REST/GraphQL).
3. Solusi utama: *Request batching* atau *query-level eager loading*.
4. *Mapping-level eager fetching* (seperti JPA `FetchType.EAGER`) adalah *anti-pattern* penyebab *memory bloat*.
5. Verifikasi penghematan secara konkrit melalui *unit tests* yang menghitung jumlah total eksekusi query.
