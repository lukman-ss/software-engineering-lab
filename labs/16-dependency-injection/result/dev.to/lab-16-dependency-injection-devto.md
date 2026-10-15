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

:  () :  () :  ()

## Review Manual Melewatkan Hidden Dependencies

:  () :  () :  ()

## Separation of Configuration from Use

:  () :  ()

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

:  () :  () :  ()

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

:  () :  () :  ()

## Test Suite 02: Propagasi Kegagalan Gateway

:  () :  () :  ()

## Test Suite 03: Short-Circuit Invalid Amount

:  () :  () :  ()

## Desain Rasional: Value Object Bebas DI

:  () :  () :  ()

## Bahaya Service Locator: PSR-11 Warning

:  () :  () :  ()

## Checklist Adopsi Dependency Injection

Deklarasikan dependensi sebagai interface, bukan tipe struct konkret. Gunakan Constructor Injection sebagai standar default; tolak setter opsional. Pusatkan instansiasi objek konkret hanya pada Composition Root (main.go). Jangan pernah mengoper Container atau Service Locator ke dalam domain logic. Instansiasi struct data murni (Value Objects) langsung tanpa melalui DI. Pastikan setiap method bisnis memiliki unit test dengan test double terisolasi.

## Tiga Anti-Pattern Umum Seputar DI

:  () :  () :  ()

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
