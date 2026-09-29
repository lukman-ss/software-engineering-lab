# Key Takeaways: Property-Based Testing

Poin-poin esensial dari riset dan implementasi Property-Based Testing pada Lab 40:

1. **Eksplorasi Domain Mengalahkan Contoh Manual:** Example-based testing hanya memvalidasi asumsi pengembang yang terbatas, sedangkan Property-Based Testing (PBT) mengeksplorasi ribuan kombinasi input pada ruang domain untuk membongkar bug laten.
2. **Universal Invariant sebagai Kontrak Formal:** Menguji invariant (seperti `Decode(Encode(x)) == x` atau `f(f(x)) == f(x)`) bertindak sebagai spesifikasi kebenaran logika yang dapat dieksekusi secara otomatis oleh komputer.
3. **Deteksi Kebocoran Presisi Floating-Point:** Unit test berbasis contoh dapat memberikan ilusi aman pada representasi moneter berbasis `float64`, sementara PBT membuktikan kegagalan fatal pada roundtrip parsing dan formatting.
4. **Validasi State Idempotent:** Operasi pembersihan dan penggabungan interval harus tahan terhadap eksekusi berulang tanpa mengubah hasil setelah iterasi pertama (`Merge(Merge(x)) == Merge(x)`).
5. **Shrinking Mengisolasi Akar Masalah:** Saat menemukan kegagalan dari puluhan atau ratusan elemen input acak, mesin PBT secara sistematis memangkas data hingga menghasilkan *minimal counterexample* (misalnya dari slice 10 elemen acak menjadi `[-1]`).
6. **Biased Generation Krusial untuk Edge Cases:** Generator PBT efektif bila dirancang condong ke nilai-nilai kritis seperti nol, bilangan negatif, nilai maksimum/minimum tipe data, dan batas koleksi kosong.
7. **PBT dan Unit Test Bersifat Komplementer:** PBT tidak menggantikan unit test berbasis contoh yang menguji aturan bisnis spesifik, melainkan melengkapinya pada level ketahanan algoritma dan integritas data.
8. **Dukungan Standard Library Go:** Go menyediakan paket `testing/quick` bawaan tanpa dependensi pihak ketiga, yang dapat langsung digunakan untuk implementasi PBT dan generator custom.
