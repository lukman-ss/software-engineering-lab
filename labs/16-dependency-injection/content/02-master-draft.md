# Dependency Injection dan Inversion of Control

## Problem
Ketika sebuah kelas menginstansiasi dependensinya secara langsung (misalnya dengan pemanggilan langsung ke pustaka eksternal atau gateway konkret), kode tersebut menjadi terikat kuat (*tightly coupled*). Ketergantungan ini menyulitkan penggantian pustaka di masa mendatang dan membuat pengujian unit yang terisolasi menjadi tidak praktis, karena pengujian terpaksa melakukan panggilan infrastruktur nyata (seperti jaringan atau basis data).

## Why This Matters
Keterikatan yang kuat merusak pemeliharaan kode dalam jangka panjang:
- Pengujian menjadi lambat dan rapuh (*flaky*) karena ketergantungan pada status eksternal.
- Setiap perubahan pada penyedia eksternal memicu modifikasi menyeluruh di berbagai lapisan logika bisnis.
- Arsitektur sistem menjadi kaku dan sulit diekstensi.

## Mental Model
Pemisahan konfigurasi dari penggunaan (*Separation of Configuration from Use*). Sebuah modul logika bisnis tidak boleh dibebani tanggung jawab untuk mencari atau membangun dependensinya. Sebagai gantinya, entitas perakit eksternal (*assembler* atau kontainer) akan menyuntikkan dependensi yang diperlukan ke dalam modul tersebut.

## Core Concept
Dependency Injection (DI) adalah bentuk spesifik dari prinsip Inversion of Control (IoC). Logika bisnis hanya bergantung pada antarmuka (*interface*), bukan tipe data konkret.
Terdapat dua pendekatan utama yang sering dibandingkan:
1. **Constructor Injection:** Menyediakan seluruh dependensi yang diperlukan secara eksplisit melalui parameter konstruktor.
2. **Service Locator:** Menyuntikkan kontainer atau penyedia layanan langsung ke dalam kelas, di mana kelas tersebut secara aktif meminta dependensi yang dibutuhkannya.

Di samping itu, konsep penting lainnya adalah bahwa tidak semua objek memerlukan DI. *Value Objects* atau struktur data sederhana yang murni menyimpan status tanpa interaksi infrastruktur (seperti objek uang/mata uang) sebaiknya diinstansiasi secara langsung tanpa melewati mekanisme DI.

## How It Works
1. Modul bisnis mendeklarasikan antarmuka abstraksi yang ia butuhkan.
2. Lapisan pembungkus atau konfigurasi menginisialisasi implementasi konkret.
3. Lapisan perakit memasukkan implementasi konkret tersebut ke dalam modul bisnis saat modul dibuat.
4. Ketika pengujian dijalankan, implementasi tiruan (*mock*) disuntikkan menggantikan implementasi konkret untuk memvalidasi alur tanpa efek samping.

## Architecture
- `di.PaymentGateway`: Antarmuka yang mengabstraksi pemrosesan pembayaran eksternal.
- `di.RealGateway`: Implementasi konkret yang menyimulasikan panggilan eksternal.
- `di.Processor`: Layanan yang menerapkan *Constructor Injection* untuk menerima `PaymentGateway`.
- `di.BadProcessor`: Layanan yang menerapkan pola *Service Locator* dengan menerima antarmuka `Container`.
- `di.Money`: *Value object* sederhana yang dibuat langsung di dalam modul logika bisnis.

## Implementation & Code Walkthrough
Komponen `Processor` memastikan bahwa dependensinya terpenuhi secara eksplisit saat inisialisasi:
```go
type Processor struct {
	gateway PaymentGateway
}

func NewProcessor(g PaymentGateway) *Processor {
	return &Processor{gateway: g}
}
```
Pendekatan ini menjamin bahwa instansi `Processor` selalu berada dalam status valid dan dependensinya bersifat *immutable* (tidak dapat diubah sembarangan setelah inisialisasi).

Sebaliknya, pada implementasi anti-pattern `BadProcessor`, sebuah antarmuka kontainer disuntikkan:
```go
type BadProcessor struct {
	container Container
}

func NewBadProcessor(c Container) *BadProcessor {
	return &BadProcessor{container: c}
}
```
Metode ini menyembunyikan dependensi `PaymentGateway` yang sebenarnya dibutuhkan dari tanda tangan fungsi (*signature*), serta mengikat komponen logika pada API kontainer.

## What the Tests Prove
Pengujian unit pada `tests/processor_test.go` membuktikan manfaat nyata dari *Constructor Injection*:
1. **Pemisahan dari Infrastruktur Nyata:** Dengan memanfaatkan `MockGateway`, pengujian dapat dijalankan tanpa melibatkan jaringan atau sistem nyata.
2. **Pengujian Jalur Sukses:** Memvalidasi bahwa pemanggilan fungsi memproses data dengan benar dan parameter transaksi tercatat pada objek tiruan.
3. **Pengujian Kegagalan Eksternal:** Menyimulasikan kegagalan gateway eksternal secara terprediksi dan memverifikasi penanganan error yang sesuai.
4. **Validasi Input:** Memastikan bahwa input yang tidak valid ditolak sebelum pemanggilan infrastruktur dilakukan.

## Production Considerations
- Dalam aplikasi berskala besar, penyambungan (*wiring*) dependensi manual dapat menjadi kompleks, sehingga penggunaan kontainer IoC otomatis atau *dependency injection framework* dapat dipertimbangkan.
- Service Locator sebaiknya dihindari di dalam logika domain karena mengurangi keterbacaan kode dan mempersulit pembuatan mock dalam skenario pengujian.
- Hindari menyuntikkan *Value Objects* ke dalam rantai DI demi mencegah penumpukan abstraksi yang tidak perlu.

## Common Mistakes
- **Menggunakan Service Locator sebagai DI:** Memasukkan objek kontainer ke dalam kelas bisnis untuk mengambil dependensi secara dinamis.
- **Terlalu Banyak Mengabstraksi Objek Sederhana:** Memasukkan struktur data atau primitif (seperti DTO, entitas nilai) ke dalam kontainer DI.
- **Constructor Bloat:** Konstruktor yang menerima terlalu banyak parameter dependensi menandakan pelanggaran terhadap prinsip *Single Responsibility Principle* (SRP).

## Key Takeaways
- Gunakan *Constructor Injection* untuk membuat dependensi terlihat jelas dan menjamin keutuhan status objek.
- Hindari *Service Locator* karena menyembunyikan dependensi dan mengikat kode pada pustaka kontainer.
- Hanya gunakan DI untuk layanan atau komponen yang berinteraksi dengan infrastruktur eksternal, bukan untuk objek nilai sederhana.

## Sources
- Martin Fowler, "Inversion of Control Containers and the Dependency Injection pattern" (2004).
- Microsoft Learn, "Dependency injection - .NET" (2026).
- Spring Framework Documentation, "Introduction to the Spring IoC Container and Beans".
- PHP-FIG, "PSR-11: Container interface Meta Document" (2017).
