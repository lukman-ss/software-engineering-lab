Ketika sebuah modul logika domain menginstansiasi dependensinya secara langsung — misalnya memanggil `new RealGateway()` di dalam method bisnis — boundary arsitektur runtuh. Dependensi konkret terkunci, pengujian terpaksa menjadi integration test lambat yang bergantung pada koneksi jaringan, dan setiap pergantian vendor memaksa modifikasi meluas di domain logic.

Di Lab 16 dari seri Software Engineering Lab, kami mengeksplorasi pemisahan konstruksi dari penggunaan (*Separation of Configuration from Use*) menggunakan Go murni tanpa framework DI pihak ketiga.

Sorotan utama lab ini:
- **Constructor Injection eksplisit**: Signature `NewProcessor(g PaymentGateway)` memaksa dependensi terlihat sejak lahir dan menjamin compile-time contract.
- **Service Locator Anti-Pattern dibongkar**: Menginjeksi kontainer (`Container.GetPaymentGateway()`) mengaburkan dependensi nyata dan mengikat domain pada API lookup runtime (sejalan dengan rekomendasi PSR-11 "SHOULD NOT").
- **Test Double deterministik**: `MockGateway` menggantikan infrastruktur jaringan 100%, membuktikan penolakan input invalid sebelum gateway dipanggil dan propagasi error deterministik.
- **Zero Race Data**: 6 test suite PASS dalam 0.00s dengan flag `-race`.

Penting dicatat: DI memperjelas kontrak dan mempermudah isolasi test, tetapi tidak otomatis menjamin kebenaran runtime jika constructor tidak memvalidasi nil pointer, dan race-clean di test kecil bukan bukti seluruh lifecycle singleton di production otomatis aman.

Repositori kode sumber dan hasil eksekusi lengkap:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/16-dependency-injection

#SoftwareEngineering #Golang #SystemDesign #SoftwareArchitecture #CleanArchitecture #Testing
