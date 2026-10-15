---
title: "Dependency Injection: Separating Configuration from Use for Testability and Evolvability"
published: false
description: "Explore Dependency Injection in Go: mechanics, architecture, caveats, and testable invariants."
tags: go, testing, performance, softwareengineering
---

# Dependency Injection: Separating Configuration from Use for Testability and Evolvability

A production system must remain resilient under load. Lab 16 makes that behavior visible with a focused Go simulation and deterministic test assertions.

> Hard-coded dependencies fuse domain logic to external plumbing, destroying testability.

## Infrastruktur Terkunci di Logika Bisnis

Tightly Coupled Code: Processor menginstansiasi RealGateway langsung di method; logika bisnis terkunci ke implementasi vendor tunggal. Pengujian Lambat & Flaky: Unit test terpaksa membuka koneksi jaringan nyata atau memerlukan mock REST API server yang lambat dan rapuh. Ripple Effect Perubahan: Mengganti vendor payment gateway memaksa refactoring dan rekompilasi menyeluruh di seluruh lapisan consumer.

## Review Manual Melewatkan Hidden Dependencies

Kontrak Opaque: Signature NewBadProcessor(c Container) tidak menyatakan butuh PaymentGateway; dependensi baru terlihat di dalam baris method. Mocking Menjadi Rumit: Test tidak bisa menyuntikkan MockGateway langsung, melainkan harus membuat MockContainer yang membungkus MockGateway. Ledakan Runtime: Jika container lupa meregistrasikan gateway, compiler tetap lolos namun aplikasi meledak panic saat dipanggil di production.

## Separation of Configuration from Use

Fase Konfigurasi (Composition Root) (main.go / Assembler): Bertanggung jawab membuat objek konkret (RealGateway), merakit relasi, dan menginjeksikan ke consumer. Fase Penggunaan (Domain Execution) (Processor (Business Service)): Hanya menerima PaymentGateway via interface; mengeksekusi ProcessPayment tanpa tahu kelas konkret di baliknya.

## Pipeline Alur Eksekusi Dependency Injection

cmd/demo/main.go | Composition Root inisialisasi RealGateway di.NewProcessor(gateway) | Konstruktor menerima kontrak PaymentGateway processor.ProcessPayment(100) | Validasi input amount <= 0 & buat Money gateway.Charge(Money) | Delegasi eksekusi ke implementasi konkret

## Constructor Injection: Valid State at Birth

```go
type Processor struct {
    gateway PaymentGateway
}

func NewProcessor(g PaymentGateway) *Processor {
    return &Processor{gateway: g}
}
```

## Failure Path: Pencegahan Eksekusi I/O Sia-sia

Input Transaksi | amount = -10 (invalid input) Processor Pre-check | if amount <= 0 { return error } Short-Circuit Exit | Return errors.New("invalid amount") Infrastruktur Aman | gateway.Charge() tidak pernah dipanggil (Amount = 0)

## Tabel Evaluasi Mekanisme Dependensi

Kriteria | Hard-Coded New | Service Locator | Constructor Injection
|---|---|---|---|
Transparansi Kontrak | Tertutup di method | Tertutup di container | Terbuka di constructor
Validasi Compiler | Kuat (tipe konkret) | Lemah (runtime get) | Sangat Kuat (interface)
Kompleksitas Test | Sangat tinggi (mock I/O) | Sedang (mock container) | Minimal (mock interface)
Standar Industri | Anti-pattern monolit | PSR-11: SHOULD NOT | Rekomendasi Fowler / Default

## Tiga Keputusan Desain Dependensi

Program to Interface: Definisikan PaymentGateway Charge(Money) error untuk mengabstraksi seluruh varian vendor pembayaran. Explicit Dependencies: Tolak Service Locator; oper PaymentGateway langsung ke konstruktor NewProcessor. Exempt Pure Data: Instansiasi Money{Amount, Currency} langsung tanpa DI container karena tidak memiliki efek I/O.

## Output Demo: Eksekusi Dua Pendekatan

```sh
$ go run ./cmd/demo
--- Running Constructor Injection ---
RealGateway charging 100 USD
--- Running Service Locator ---
RealGateway charging 200 USD
```

## Batasan Bukti: Cakupan Verifikasi Lab

MockGateway mengisolasi Processor 100% tanpa jaringan. Validasi pre-check mencegah pemanggilan gateway. Error gateway dipropagasikan utuh ke caller. Race detector bersih tanpa data race. Overhead performa container vs direct instantiation. Korelasi penurunan angka bug/defect di tim engineering. Lifecycle manajemen (Singleton, Scoped, Transient). Efektivitas DI framework otomatis (Wire, Dig) skala besar.

## Trade-off: Kontrol Bersih vs Biaya Boilerplate

Isolasi total: unit test berjalan cepat dalam orde mikrodetik. Vendor agnostic: ganti gateway tanpa sentuh logika bisnis. Immutabilitas: dependensi tidak dapat diganti sembarangan di runtime. Manual wiring menambah baris inisialisasi di composition root. Tracing alur kode memerlukan navigasi antarmuka interface. Risiko constructor bloat jika kelas melanggar Single Responsibility.

## Pipeline Validasi & Thread-Safety

Goroutine Request | Panggilan concurrent ke ProcessPayment(amount) Pure Pre-Validation | amount <= 0 dicek tanpa mutasi shared state Value Object Build | Money{Amount, Currency} dialokasikan per request Thread-Safe Charge | Gateway dieksekusi secara aman (Race Clean)

## Sebelum vs Sesudah Penerapan DI

Processor memanggil new RealGateway() secara internal. Unit test terpaksa menghubungi server gateway nyata. Pergantian vendor menuntut perubahan di logika bisnis. Sulit menguji skenario kegagalan jaringan secara deterministik. Processor menerima PaymentGateway via NewProcessor. Unit test menyuntikkan MockGateway tanpa sentuh jaringan. Pergantian vendor hanya mengubah satu baris di main.go. Skenario kegagalan gateway dapat diuji deterministik 100%.

## Test Suite 01: Verifikasi Happy Path Terisolasi

Substitusi Bersih: mock := &MockGateway{} disuntikkan ke NewProcessor(mock); tidak ada panggilan TCP atau setup database. Validasi Status Mutasi: mock.ChargedMoney mencatat persis Amount: 50 dan Currency: 'USD' setelah ProcessPayment(50) berhasil. Zero-Latency Run: Eksekusi test selesai dalam 0.00s deterministik, bebas dari network flakiness atau rate-limit vendor.

## Test Suite 02: Propagasi Kegagalan Gateway

Simulasi Down Time: MockGateway dengan ShouldFail: true mensimulasikan kondisi 'gateway unavailable' secara instan. Error Forwarding: Processor mengembalikan error gateway langsung ke caller tanpa menelan (swallow) pesan kesalahan. Audit Integritas: Caller menerima sinyal kegagalan deterministik untuk memicu mekanisme retry atau fallback otomatis.

## Test Suite 03: Short-Circuit Invalid Amount

Input Negatif / Nol: ProcessPayment(-10) menghasilkan error 'invalid amount' seketika pada boundary validasi. Assertion Gateway Call: Assertion mock.ChargedMoney.Amount == 0 membuktikan Charge() tidak pernah disentuh sama sekali. Cost & Traffic Saving: Mencegah overhead biaya API call vendor dan transmisi data yang ditakdirkan gagal sejak awal.

## Desain Rasional: Value Object Bebas DI

Data Murni Tanpa Efek: Struct Money{Amount, Currency} tidak berinteraksi dengan database, network, atau state eksternal. Bahaya Over-Injection: Memasukkan struct data ke dalam container DI menambah layer indirection tanpa keuntungan testability. Heuristik Desain: Gunakan DI untuk Service ber-efek samping (I/O); gunakan direct instantiation untuk data value murni.

## Bahaya Service Locator: PSR-11 Warning

Pelanggaran PSR-11: Standar PHP-FIG PSR-11 menegaskan: 'Users SHOULD NOT pass a container into an object'. Double-Mocking Overhead: Pengujian BadProcessor mewajibkan pembuatan MockContainer hanya untuk mengembalikan MockGateway. Kopling ke API Container: Kelas bisnis menjadi terikat pada method GetPaymentGateway() milik container, bukan antarmuka generik.

## Checklist Adopsi Dependency Injection

Deklarasikan dependensi sebagai interface, bukan tipe struct konkret. Gunakan Constructor Injection sebagai standar default; tolak setter opsional. Pusatkan instansiasi objek konkret hanya pada Composition Root (main.go). Jangan pernah mengoper Container atau Service Locator ke dalam domain logic. Instansiasi struct data murni (Value Objects) langsung tanpa melalui DI. Pastikan setiap method bisnis memiliki unit test dengan test double terisolasi.

## Tiga Anti-Pattern Umum Seputar DI

Constructor Bloat: Konstruktor dengan belasan parameter adalah sinyal kuat pelanggaran Single Responsibility (SRP). Service Locator Disguise: Membungkus dependensi ke dalam satu Map atau Container serbaguna demi menyederhanakan signature ctor. Over-Abstracting Everything: Membuat interface untuk setiap kelas tanpa pernah memiliki lebih dari satu implementasi atau kebutuhan mock.

## Reproduksi Mandiri Suite Pengujian

```sh
cd labs/16-dependency-injection
go test -v ./...       # Verifikasi 6 skenario test suite
go test -race ./...    # Verifikasi thread-safety zero data race
go run ./cmd/demo      # Eksekusi demo komparasi Ctor vs Locator
```

## Prinsip Arsitektur: Hollywood Principle

Inversion of Control: Logika bisnis fokus mengeksekusi use case; framework/assembler yang menyediakan instrumen. Program to Interfaces: Modul tingkat tinggi tidak boleh bergantung pada modul tingkat rendah. Evolvability: Kemampuan menukar provider atau infrastruktur tanpa memicu modifikasi pada logika inti domain.

## Source

[Software Engineering Lab 16 - Dependency Injection](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/16-dependency-injection), snapshot commit `8678a90`. Source: README, implementation, tests, and execution records. Code excerpts belong in their original package context.
