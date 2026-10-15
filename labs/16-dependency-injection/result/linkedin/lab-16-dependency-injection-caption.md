Software Engineering Lab 16 — Dependency Injection di Go ini membandingkan dua pendekatan pengelolaan ketergantungan: Constructor Injection eksplisit versus Service Locator yang menyembunyikan dependency.

Di mana sebuah komponen payment system menginstansiasi gateway jaringan secara langsung di dalam method-nya sendiri, setiap perubahan vendor memaksa modifikasi pada domain logic. Lab ini memindahkan keputusan instansiasi itu ke composition root dengan menerima PaymentGateway interface pada konstruktor.

Empat poin kunci:
1. Konstruktor dengan parameter interface membuat dependensi terlihat di signature publik — bukan sekedar memudahkan mock, tetapi membuat kebutuhan eksternal dapat dibaca oleh siapa pun.
2. Service Locator menyamarkan ketergantungan riil di balik method container.Get(...), sehingga engineer harus mem-trace indirection ekstra dan test terpaksa mem-mock container.
3. Nil guard dan lifecycle management bukan tanggung jawab pola DI — mereka adalah tanggung jawab implementasi yang masih perlu ditambahkan secara eksplisit.
4. Race detector bersih pada 6 test cases (0.00s) hanya membuktikan jalur yang diuji aman, bukan garantee semua lifecycle production.

Lab ini tidak memakai framework DI eksternal. Manual wiring di main.go membekalkan compile-time type safety, sementara BadProcessor berperilaku fungsional sama tetapi menambah indirection tak perlu. Simulasi in-memory mock bukan pengganti load test atau integration test produksi.

Lihat thread carousel di slide 1 dan link repo di profil untuk selengkapnya.
