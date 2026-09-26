# Dependency Injection dan Inversion of Control: Memisahkan Konfigurasi dari Penggunaan

## Problem

Ketika sebuah kelas menginstansiasi dependensinya secara langsung — misalnya `Processor` yang memanggil `new RealGateway()` di dalam method-nya — kode tersebut menjadi tightly coupled. Dependensi konkret terkunci di dalam logika bisnis. Akibatnya, mengganti payment gateway, mengganti provider notifikasi, atau bahkan sekadar menguji `Processor` tanpa memanggil jaringan menjadi sulit. Setiap perubahan pada infrastruktur eksternal memaksa modifikasi di banyak tempat, dan setiap test terpaksa menjadi integration test yang lambat serta flaky.

Lab ini mendemonstrasikan masalah tersebut melalui kontras dua desain: `Processor` yang menerima dependensi dari luar versus `BadProcessor` yang meminta dependensi dari container.

## Why This Matters

Coupling yang keras merusak evolvability sistem jangka panjang. Dampaknya konkret: pengujian menjadi lambat karena bergantung pada status eksternal (database, HTTP API, payment provider); perubahan vendor memicu modifikasi menyeluruh di lapisan logika bisnis; dan arsitektur menjadi kaku karena setiap consumer terikat pada implementasi spesifik, bukan pada kontrak interface. Pemisahan construction dari use adalah investasi untuk testability dan kemampuan sistem berevolusi tanpa rewrite.

## Mental Model

**Separation of Configuration from Use** — prinsip yang ditegaskan Fowler sebagai hal yang "lebih penting daripada pilihan antara Service Locator dan DI itu sendiri."

Bayangkan dua fase yang terpisah: fase konfigurasi (di mana objek dibuat dan dihubungkan) dan fase penggunaan (di mana logika bisnis berjalan). Modul logika bisnis tidak boleh mencari atau membangun dependensinya sendiri. Entitas eksternal — assembler, composition root, atau container — yang bertanggung jawab merakit object graph dan menyuntikkannya ke consumer. Hollywood Principle berlaku: "Don't call us, we'll call you" — framework memanggil kode aplikasi, bukan sebaliknya.

## Core Concept

### Inversion of Control (IoC)

IoC adalah prinsip desain umum: custom code menerima flow of control dari external framework. Contohnya meliputi UI main loop, callback, scheduler, dan template method. DI adalah bentuk spesifik IoC yang menginversi kontrol atas implementasi dependensi — bukan keseluruhan control flow. Fowler memperkenalkan istilah "Dependency Injection" pada 2004 justru karena "Inversion of Control" terlalu generik dan membingungkan.

### Dependency Injection (DI)

Definisi yang disepakati lintas sumber — Wikipedia, Spring 7.x, Microsoft .NET, dan Laravel 12.x/13.x — konsisten: sebuah objek menerima objek lain yang dibutuhkannya dari kode eksternal (injector), alih-alih menciptakannya secara internal. DI memisahkan concern konstruksi objek dari penggunaannya, menghasilkan loosely coupled programs.

DI melibatkan empat peran: service (implementasi), client (consumer), interface (kontrak), dan injector (assembler/container yang menghubungkan service ke client). Injector tidak boleh menjadi client itu sendiri untuk menghindari circular dependency.

### Tiga Bentuk DI

Secara historis terdapat tiga bentuk: Constructor Injection (Type 3), Setter Injection (Type 2), dan Interface Injection (Type 1). Framework modern — Spring 7.x, .NET, Laravel — hanya mendokumentasikan dua bentuk pertama; Interface Injection dianggap obsolete (pernah digunakan di Avalon). Praktik modern: constructor injection sebagai default, setter sebagai fallback.

### DI vs Service Locator

Keduanya men-decoupling client dari implementasi konkret, tetapi mekanismenya berbeda. Dengan Service Locator, setiap consumer memiliki dependensi eksplisit terhadap locator dan secara aktif meminta service (`locator.get("gateway")`). Dengan DI, service muncul di consumer melalui injection — tidak ada request eksplisit, itulah inversion of control. PSR-11 menegaskan: "Users SHOULD NOT pass a container into an object" — container sebagai Service Locator secara umum discouraged. Kata kunci RFC 2119 di sini adalah SHOULD NOT (rekomendasi kuat, bukan larangan absolut).

## Failure Scenario

Lab ini memverifikasi dua skenario kegagalan:

**1. Invalid input (amount <= 0).** Baik `Processor` maupun `BadProcessor` melakukan validasi `if amount <= 0` sebelum memanggil gateway. Jika amount tidak valid, method mengembalikan `errors.New("invalid amount")` tanpa pernah memanggil `Charge`. Behavior ini terbukti oleh `TestProcessor_InvalidAmount` dan `TestBadProcessor_InvalidAmount` — keduanya menegaskan bahwa `mock.ChargedMoney.Amount` tetap 0, artinya gateway tidak pernah dipanggil.

**2. Gateway error.** Ketika `MockGateway.ShouldFail = true`, method `Charge` mengembalikan `errors.New("gateway unavailable")`. Baik `Processor` maupun `BadProcessor` mempropagasikan error tersebut ke caller tanpa menelan atau mengubahnya. Terbukti oleh `TestProcessor_GatewayError` dan `TestBadProcessor_GatewayError`.

Tanpa validasi pre-check, sistem akan melakukan panggilan jaringan yang sia-sia atau bahkan memproses transaksi dengan nilai tidak valid. DI tidak menghilangkan kebutuhan validasi — ia hanya memastikan bahwa validasi dan pemanggilan infrastruktur dapat diuji secara terpisah.

## How It Works

Alur kerja DI dalam lab ini:

1. Modul bisnis (`Processor`) mendeklarasikan dependensi sebagai interface `PaymentGateway` melalui parameter constructor `NewProcessor(g PaymentGateway)`.
2. Lapisan konfigurasi (`cmd/demo/main.go`) menginisialisasi implementasi konkret `&di.RealGateway{}`.
3. Composition root menyuntikkan implementasi tersebut: `processor := di.NewProcessor(realGateway)`.
4. Saat `processor.ProcessPayment(100)` dipanggil, `Processor` membuat value object `Money{Amount: 100, Currency: "USD"}` secara langsung dan mendelegasikan ke `gateway.Charge(m)`.
5. Dalam pengujian, `MockGateway` disuntikkan menggantikan `RealGateway` — tanpa perubahan pada kode `Processor`.

Untuk Service Locator, alurnya berbeda: `main.go` membungkus `RealGateway` ke dalam `SimpleContainer`, lalu `BadProcessor` menerima container tersebut dan di dalam `ProcessPayment` memanggil `p.container.GetPaymentGateway().Charge(m)` — satu langkah indirection tambahan yang menyembunyikan dependensi sebenarnya.

## Architecture

| Komponen | Lokasi | Peran |
|---|---|---|
| `PaymentGateway` | `internal/di/gateway.go` | Interface yang mengabstraksi payment provider eksternal |
| `RealGateway` | `internal/di/gateway.go` | Implementasi konkret yang mensimulasikan panggilan jaringan |
| `Money` | `internal/di/gateway.go` | Value object — data murni tanpa dependensi infrastruktur |
| `Processor` | `internal/di/processor.go` | Service dengan Constructor Injection |
| `Container` | `internal/di/locator.go` | Interface Service Locator |
| `BadProcessor` | `internal/di/locator.go` | Service dengan Service Locator anti-pattern |
| `SimpleContainer` | `cmd/demo/main.go` | Container minimal untuk demo |
| `MockGateway` | `tests/processor_test.go` | Test double untuk pengujian terisolasi |

Dependency graph untuk Constructor Injection: `main.go` → `RealGateway` + `Processor(PaymentGateway)` → `Money`. Untuk Service Locator: `main.go` → `RealGateway` → `SimpleContainer` → `BadProcessor(Container)` → `Container.GetPaymentGateway()` → `PaymentGateway`.

## Implementation

Lab ini sengaja menggunakan manual wiring tanpa DI framework (seperti `google/wire` atau `uber/dig`). Keputusan ini mengikuti prinsip YAGNI: Go yang sederhana sudah cukup membuktikan konsep struktural tanpa menambah dependensi framework. Manual wiring di `main.go` memang menambah sedikit boilerplate dibanding auto-wiring, tetapi menjamin type-safety dan compile-time validation.

Container dalam lab ini dibuat minimal — hanya interface dengan satu method `GetPaymentGateway() PaymentGateway` — untuk mengisolasi perilaku anti-pattern tanpa memerlukan framework berat. Lab ini tidak mendemonstrasikan reflection-based atau code-gen based DI container, dan tidak mengelola lifecycle singleton/scoped/transient.

## Code Walkthrough

### Interface dan Value Object — `internal/di/gateway.go`

`PaymentGateway` mendefinisikan kontrak tunggal `Charge(m Money) error`. `Money` adalah struct dengan dua field `Amount int` dan `Currency string` — value object yang diinstansiasi langsung tanpa melalui DI, sesuai prinsip bahwa DI dicadangkan untuk service yang berinteraksi dengan infrastruktur eksternal.

`RealGateway.Charge` hanya mencetak `RealGateway charging %d %s` dan mengembalikan `nil` — simulasi panggilan jaringan tanpa efek samping nyata.

### Constructor Injection — `internal/di/processor.go`

`Processor` menyimpan `gateway PaymentGateway` sebagai field unexported. `NewProcessor` menerima interface tersebut sebagai parameter dan mengembalikannya dalam keadaan valid. `ProcessPayment` memvalidasi `amount <= 0` terlebih dahulu, lalu membuat `Money` secara langsung dan mendelegasikan ke `p.gateway.Charge(m)`. Tidak ada import container, tidak ada lookup dinamis.

### Service Locator — `internal/di/locator.go`

`Container` mendefinisikan `GetPaymentGateway() PaymentGateway`. `BadProcessor` menyimpan `container Container` alih-alih `PaymentGateway` langsung. Konstruktor `NewBadProcessor(c Container)` menerima container, dan `ProcessPayment` harus memanggil `p.container.GetPaymentGateway().Charge(m)` — dependensi sebenarnya (`PaymentGateway`) tersembunyi di balik container.

### Composition Root — `cmd/demo/main.go`

`SimpleContainer` mengimplementasikan `Container` dengan menyimpan `gateway di.PaymentGateway`. Fungsi `main` mendemonstrasikan kedua pola secara berurutan: pertama membuat `realGateway` dan `di.NewProcessor(realGateway)` untuk Constructor Injection (memproses 100 USD), kemudian membungkus gateway yang sama ke `SimpleContainer` dan `di.NewBadProcessor(container)` untuk Service Locator (memproses 200 USD). Kedua jalur menghasilkan output `RealGateway charging ...` yang identik secara fungsional.

## What the Tests Prove

Test suite di `tests/processor_test.go` berisi 6 kasus — 3 untuk `Processor` dan 3 untuk `BadProcessor` — semuanya PASS termasuk dengan race detector (`go test -race ./...`).

**Yang terbukti terverifikasi:**

- **Isolasi tanpa infrastruktur nyata.** `MockGateway` menggantikan `RealGateway` sepenuhnya. Tidak ada panggilan jaringan, tidak ada setup database. Test berjalan dalam 0.00s per kasus.
- **Happy path.** `TestProcessor_Success` (amount 50) dan `TestBadProcessor_Success` (amount 75) memverifikasi bahwa `MockGateway.ChargedMoney` mencatat nilai dan currency yang benar setelah `ProcessPayment`.
- **Error propagation.** `TestProcessor_GatewayError` dan `TestBadProcessor_GatewayError` dengan `ShouldFail: true` memverifikasi bahwa error dari gateway dipropagasikan ke caller.
- **Validasi mencegah panggilan sia-sia.** `TestProcessor_InvalidAmount` (amount -10) dan `TestBadProcessor_InvalidAmount` (amount -5) memverifikasi dua hal: error dikembalikan, dan `mock.ChargedMoney.Amount` tetap 0 — gateway tidak pernah dipanggil.
- **Race safety.** `go test -race ./...` PASS tanpa warning.

**Yang tidak terbukti (tidak ada test yang mendukung):**

- Lifecycle management singleton/scoped/transient — tidak didemonstrasikan dalam implementasi.
- Performance overhead container vs direct `new` — tidak diukur.
- Pengurangan defect rate akibat DI — tidak ada studi empiris yang ditemukan dalam riset.

## Recovery / Rollback

Lab ini tidak melibatkan state persisten atau transaksi yang memerlukan rollback. Namun pola DI sendiri mendukung recovery: karena dependensi diinjeksi melalui interface, mengganti implementasi yang bermasalah (misalnya gateway yang error) hanya memerlukan perubahan di composition root (`main.go`) tanpa menyentuh `Processor`. Jika `RealGateway` gagal di production, rollback ke implementasi sebelumnya atau swap ke provider alternatif dilakukan dengan mengganti satu baris binding, bukan mengubah logika bisnis.

Untuk pengujian, recovery dari kegagalan gateway disimulasikan melalui `MockGateway{ShouldFail: true}` — caller menerima error dan dapat memutuskan retry, fallback, atau pelaporan sesuai kebutuhan.

## Production Considerations

- **Framework vs manual wiring.** Lab ini menggunakan manual wiring yang ideal untuk codebase kecil-menengah di Go. Pada aplikasi berskala besar dengan object graph yang dalam, pertimbangkan DI framework (seperti `wire` untuk code generation atau `dig` untuk reflection) dengan trade-off: auto-wiring mengurangi boilerplate tetapi menambah kompleksitas tracing dan potensi framework lock-in.
- **Constructor Injection sebagai default.** Mulailah dengan constructor injection untuk menjamin valid object at birth dan immutable fields. Beralih ke setter injection hanya ketika menghadapi banyak parameter, kombinasi konstruksi yang valid berbeda-buat, atau inheritance constructor explosion (Fowler).
- **Hindari Service Locator di domain logic.** Service Locator mengurangi keterbacaan (dependensi tersembunyi) dan mempersulit mocking — setiap test harus mem-mock container, bukan dependensi langsung.
- **Jangan injeksikan value objects.** Struktur data murni seperti `Money` yang tidak berinteraksi dengan infrastruktur eksternal sebaiknya diinstansiasi langsung. Menginjeksikannya menambah abstraksi tanpa manfaat.
- **Perhatikan over-injection.** Constructor dengan banyak parameter adalah signal pelanggaran Single Responsibility Principle. Riset mencatat bahwa angka "12 parameter" adalah heuristik spesifik lab ini, bukan standar industri — Fowler hanya menyebut "a lot of parameters" secara kualitatif.
- **Biaya DI.** DI menambah beban konfigurasi, mempersulit tracing (behavior terpisah dari construction), dan dapat menghambat tooling IDE jika menggunakan reflection. Evaluasi apakah kompleksitas DI sebanding dengan kebutuhan sebelum menerapkannya — "prefer to avoid unless needed" (Fowler).

## Common Mistakes

**1. Menggunakan Service Locator sebagai DI.** Menyuntikkan container ke dalam kelas bisnis dan memanggil `container.Get(...)` di dalam method. Tanda: signature constructor hanya menerima `Container` tetapi method memanggil berbagai service berbeda. Solusi: injeksikan interface spesifik yang dibutuhkan.

**2. Mengabstraksi objek sederhana.** Memasukkan DTO, value object, atau primitif ke dalam container DI. `Money`, `DateTime`, `Address` dalam lab ini adalah contoh — ketiganya adalah heuristik lab, bukan daftar universal, tetapi prinsipnya jelas: jika objek tidak memiliki dependensi eksternal dan tidak memerlukan interface contract, buat langsung.

**3. Constructor bloat.** Konstruktor dengan belasan parameter menandakan kelas yang terlalu sibuk. Pecah kelas tersebut menjadi beberapa service yang lebih fokus, masing-masing dengan tanggung jawab tunggal.

**4. Menggeneralisasi heuristik lab sebagai aturan universal.** Angka 12 dan daftar value object spesifik tidak berasal dari sumber primer eksternal. Jangan mengutipnya sebagai standar industri.

**5. Mengabaikan trade-off DI.** Menerapkan DI di setiap kelas tanpa mempertimbangkan apakah kelas tersebut benar-benar memiliki dependensi yang perlu di-swap atau di-mock.

## Case Study

Skenario payment processor dalam lab ini merepresentasikan kasus nyata yang selaras dengan temuan riset: sistem yang perlu mendukung multiple payment gateways (atau PPOB/notification provider) yang dapat ditukar tanpa mengubah consumer code.

`PaymentGateway` sebagai interface memungkinkan `Processor` tetap tidak berubah ketika `RealGateway` diganti dengan implementasi lain — misalnya gateway untuk provider berbeda. Konsumen hanya bergantung pada `Charge(Money) error`, bukan pada detail HTTP client, credential, atau endpoint spesifik provider. Pergantian dilakukan di `main.go` (atau service provider/binding configuration pada framework seperti Laravel atau Spring) tanpa menyentuh `Processor`.

Test double `MockGateway` merepresentasikan manfaat yang sama di sisi pengujian: verifikasi perilaku `Processor` tanpa memerlukan kredensial gateway nyata, tanpa risiko charge ganda, dan dengan kemampuan mensimulasikan kegagalan secara deterministik.

## Checklist

- [ ] Dependensi dideklarasikan sebagai interface, bukan tipe konkret.
- [ ] Constructor injection digunakan sebagai default; semua dependensi eksplisit di signature.
- [ ] Composition root (`main.go` atau setara) adalah satu-satunya tempat yang membuat implementasi konkret.
- [ ] Tidak ada `container.Get(...)` di dalam logika domain.
- [ ] Value objects diinstansiasi langsung, tidak melalui container.
- [ ] Setiap service dapat diuji dengan mock tanpa setup infrastruktur.
- [ ] Validasi input dilakukan sebelum delegasi ke dependensi eksternal.
- [ ] Test mencakup happy path, error propagation, dan invalid input.
- [ ] `go test -race ./...` PASS tanpa warning.
- [ ] Heuristik lab (12 parameter, daftar value object) tidak dikutip sebagai standar universal.

## Key Takeaways

1. DI memisahkan konfigurasi dari penggunaan — objek menerima dependensi dari luar, bukan menciptakannya sendiri.
2. IoC adalah prinsip umum; DI adalah bentuk spesifik IoC untuk merakit dependensi.
3. Constructor injection adalah default modern — menjamin valid state dan dependensi eksplisit.
4. Service Locator menyembunyikan dependensi dan mengikat consumer pada API container — hindari di domain logic.
5. Program to interfaces — consumer bergantung pada kontrak, bukan implementasi.
6. Mock injection memungkinkan pengujian terisolasi tanpa infrastruktur nyata.
7. Value objects tanpa dependensi eksternal diinstansiasi langsung, bukan melalui DI.
8. Over-injection adalah signal SRP violation — pecah kelas yang terlalu besar.

## Sources

- Martin Fowler — "Inversion of Control Containers and the Dependency Injection pattern" (2004) — https://martinfowler.com/articles/injection.html
- Wikipedia — "Dependency injection" — https://en.wikipedia.org/wiki/Dependency_injection
- Wikipedia — "Inversion of control" — https://en.wikipedia.org/wiki/Inversion_of_control
- Laravel — "Service Container" (12.x/13.x) — https://laravel.com/docs/container
- Spring Framework 7.0.9 — "Introduction to the Spring IoC Container and Beans" — https://docs.spring.io/spring-framework/reference/core/beans/introduction.html
- PHP-FIG — "PSR-11: Container interface" — https://www.php-fig.org/psr/psr-11/
- Microsoft Learn — "Dependency Injection (.NET)" — https://learn.microsoft.com/en-us/dotnet/core/extensions/dependency-injection
- PHP Manual — "Object Interfaces" — https://www.php.net/manual/en/language.oop5.interfaces.php
