# Key Takeaways

1. Stabilitas sistem ditentukan satu pertanyaan: apakah laju masuk melebihi laju pemrosesan — bukan sekadar panjang antrean.
2. Token Bucket mengizinkan burst hingga kapasitas B, lalu menegakkan refill rate R jangka panjang.
3. Leaky Bucket meratakan aliran pada leak rate konstan dan menolak burst instan.
4. HTTP 429 (RFC 6585) adalah respons standar rate limit; selalu sertakan header `Retry-After`.
5. Isolasi per-tenant memakai API key / tenant ID, bukan IP saja — menghindari unfair throttling di balik CGNAT (RFC 6598).
6. Bounded queue dengan `TrySubmit` non-blocking menolak cepat (`ErrQueueFull`) sebelum memori habis.
7. Retry tanpa jitter menimbulkan retry storm; Full Jitter menyebar retry pada `[0, min(cap, base * 2^attempt)]` (AWS / Marc Brooker).
8. Little's Law (L = λW) mematok batas fisik backlog: 5M item pada 2000/detik butuh ~41 menit 40 detik, tak bisa diakali konfigurasi antrean.
9. Pantau oldest job age selain queue depth; umur antrean memberi sinyal degradasi lebih dini.
10. Lab ini single-node in-memory; rate limiting terdistribusi (Redis/shared state) di luar cakupan yang terverifikasi.
