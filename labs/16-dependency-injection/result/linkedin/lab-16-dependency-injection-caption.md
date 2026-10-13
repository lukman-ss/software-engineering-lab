Lab 16 kali ini membahas pemisahan konstruksi objek dari logika bisnis melalui Dependency Injection dan Inversion of Control di Go.

Masalah teknis yang sering terjadi adalah dependensi infrastruktur diinstansiasi langsung di dalam method domain, membuat kode terkunci pada provider tertentu dan pengujian unit terpaksa bergantung pada jaringan eksternal.

Lab ini membuktikan tiga hal konkret:
1. Constructor Injection: Processor menerima interface PaymentGateway secara eksplisit sehingga state objek selalu lengkap dan dependensi tampak jelas.
2. Service Locator Anti-Pattern: BadProcessor menerima Container generik, menyembunyikan dependensi nyata di balik query method runtime.
3. Fast Isolated Testing: MockGateway memungkinkan 6 test assertions berjalan instan tanpa koneksi jaringan, membuktikan error propagation serta validasi input yang menolak transaksi sebelum memanggil gateway.

Mental model inti: Externalize object creation. Biarkan assembler luar merakit object graph, sementara logika domain hanya fokus menjalankan kontrak bisnis. Nilai data murni seperti struct Money diinstansiasi langsung tanpa melalui container DI.

Satu batasan penting: Lab ini menggunakan manual wiring di Go standar dan tidak mendemonstrasikan automated DI framework (seperti reflection atau code-gen), serta tidak mengelola container lifecycle (singleton, scoped, transient).

Pelajari implementasi lengkap dan test assertion di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/16-dependency-injection

#softwareengineering #golang #systemdesign #backend #architecture #designpatterns
