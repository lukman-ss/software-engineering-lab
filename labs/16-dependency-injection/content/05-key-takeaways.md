# Key Takeaways

1. **Pemisahan Configuration dari Use**: Konstruksi objek dilakukan di composition root (`main.go`), bukan di dalam kelas bisnis — ini memungkinkan swap implementasi tanpa menyentuh logika domain.

2. **Constructor Injection sebagai Default**: Menjamin objek valid saat lahir (`valid at birth`) dan dependensi bersifat eksplisit di signature — immutable dan tipe-safe secara compile-time.

3. **Service Locator adalah Anti-Pattern**: Menyuntikkan container ke dalam kelas menyembunyikan dependensi sebenarnya dan mengikat consumer pada API kontainer — PSR-11 merekomendasikan untuk tidak melakukannya (`SHOULD NOT`).

4. **Program to Interfaces**: Consumer bergantung pada `PaymentGateway` interface, bukan `RealGateway` konkret — ini memungkinkan swap implementasi (multiple payment gateways, PPOB providers, notification providers) tanpa mengubah consumer code.

5. **Mock Injection untuk Pengujian Terisolasi**: Karena dependensi diinjeksi via interface, test dapat mensubstitusi `MockGateway` tanpa setup jaringan — 6 test cases PASS termasuk race detector.

6. **Value Objects Tidak Perlu DI**: Struktur data murni seperti `Money` yang tidak berinteraksi dengan infrastruktur eksternal sebaiknya diinstansiasi langsung — bukan dikelola oleh container.

7. **Over-Injection adalah SRP Violation**: Constructor dengan terlalu banyak parameter menandakan kelas yang terlalu sibuk — pecah menjadi beberapa service yang lebih fokus.

8. **DI Memiliki Biaya**: Konfigurasi tambahan, tracing lebih sulit (behavior terpisah dari construction), dan potensi framework lock-in — terapkan hanya ketika manfaatnya sebanding.

9. **Heuristik Lab Bukan Standar Industri**: Angka "12 parameter" dan daftar value object (`DateTime`, `Money`, `Address`) adalah panduan internal lab ini, bukan definisi universal dari sumber primer.
