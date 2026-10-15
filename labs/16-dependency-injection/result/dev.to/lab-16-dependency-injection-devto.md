---
title: "Dependency Injection: Inversion of Control & Modular Code"
published: false
description: "Explore Dependency Injection in Go: mechanics, architecture, caveats, and testable invariants."
tags: go, testing, performance, softwareengineering
---

# Dependency Injection: Inversion of Control & Modular Code

A production system must remain resilient under load. Lab 16 makes that behavior visible with a focused Go simulation and deterministic test assertions.

> Tight coupling erases architectural context from code.

## Konteks Arsitektur Hilang Tanpa Jejak

Kelas membuat objek dependensinya sendiri → tight coupling. Ganti payment gateway memaksa ubah banyak tempat. Test menjadi integration test yang lambat.

## Mengapa Ini Penting

Test lambat karena bergantung pada status eksternal (API, database). Perubahan vendor memicu modifikasi menyeluruh lapisan bisnis. Arsitektur kaku karena tiap consumer terikat pada implementasi spesifik.

## Pemisahan Konfigurasi dari Penggunaan

Modul logika tidak boleh mencari atau membangun dependensinya. Assembler/Composition Root yang bertanggung jawab merakit object graph. Hollywood Principle: Don't call us, we'll call you.

## DI vs Service Locator

Dependensi muncul lewat parameter constructor. Tidak ada request eksplisit ke container. Jelas: Processor(PaymentGateway) Consumer punya dependensi pada locator. Memanggil container.Get(...) secara aktif. Tersembunyi: processor hanya lihat Container

## Interface & Value Object

```go
type Money struct {
    Amount   int
    Currency string
}
type PaymentGateway interface {
    Charge(m Money) error
}
```

## Constructor Injection

```go
type Processor struct {
    gateway PaymentGateway
}
func NewProcessor(g PaymentGateway) *Processor {
    return &Processor{gateway: g}
}
func (p *Processor) ProcessPayment(amount int) error {
    if amount <= 0 { return errors.New("invalid amount") }
    m := Money{Amount: amount, Currency: "USD"}
    return p.gateway.Charge(m) }
```

## Skenario Kegagalan: Invalid Input

ProcessPayment(-10): Result (Error: invalid amount, gateway NOT called) ProcessPayment(100): Result (Money created, gateway.Charge() invoked)

## Bukti: Isolated Testing Tanpa Infrastruktur

```go
mock := &MockGateway{}
proc := di.NewProcessor(mock)
err := proc.ProcessPayment(50)
if err != nil { t.Fatalf("expected nil, got %v", err) }
if mock.ChargedMoney.Amount != 50 { t.Errorf("wrong")
```

## Batasan Klaim

DI memudahkan testing tanpa infrastructure nyata. MockGateway menggantikan network call secara penuh. Tes berjalan dalam 0.00s per kasus. "bukan berarti" behaving di production sama dengan test. "tidak otomatis" performa atau defect rate menurun. "perlu bukti tambahan" untuk skala enterprise.

## Batas Integritas

Lifecycle management singleton/scoped/transient. Performance overhead container vs direct new. Pengurangan defect rate akibat DI.

## Manual Wiring Tanpa Framework

main.go menginisialisasi RealGateway. NewProcessor(realGateway) – Constructor Injection. Output: RealGateway charging 100 USD

## Payment Processor: Contrast

Processor: NewProcessor(realGateway) BadProcessor: NewBadProcessor(container) Output keduanya identik: RealGateway charging ...

## Review Kode: Anti-Pattern

```go
type BadProcessor struct {
    container Container
}
func (p *BadProcessor) ProcessPayment(amount int) error {
    return p.container.GetPaymentGateway().Charge(...)
}
```

## Checklist DI dalam Tim

Deklarasikan dependensi sebagai interface. Gunakan constructor injection sebagai default. Composition root SATU-satunya tempat buat implementasi konkret. Container.Get(...) di logic domain. Inject value object ke container. Constructor dengan banyak parameter (>12 sebagai sinyal).

## Kesalahan Umum

Inject container ke kelas; panggil container.Get() di method. Solusi: injeksikan interface spesifik yang dibutuhkan. Masukkan Money, DateTime, Address ke container. Solusi: buat langsung, tidak butuh interface.

## Biaya DI

Testing terisolasi tanpa infrastructure. Swap implementasi lewat composition root. Clear dependency contracts. Konfigurasi ekstra di main.go. Tracing kesulitan (behavior dipisah dari construction). Refactoring penting untuk memilih DI framework.

## Kesimpulan

Constructor injection default untuk testability. Avoid Service Locator di domain logic. Mock injection = testing isolation tanpa network. Value objects langsung dibuat, tidak lewat DI.

## 8 Poin Kunci

1: Sepetakan konfigurasi dari penggunaan () 2: IoC umum; DI spesifik untuk dependensi () 3: Constructor injection default () 4: Service Locator dilarang di domain () 5: Program to interfaces () 6: Mock injection untuk pengujian () 7: Value objects langsung dibuat () 8: Over-injection = SRP violation ()

## Komponen DI

Komponen | Lokasi | Peran
|---|---|---|
PaymentGateway | internal/di/gateway.go | Interface abstraksi payment provider
RealGateway | internal/di/gateway.go | Implementasi produksi
Money | internal/di/gateway.go | Value object
Processor | internal/di/processor.go | Service dengan Constructor Injection
Container | internal/di/locator.go | Interface Service Locator
BadProcessor | internal/di/locator.go | Service dengan Anti-pattern
MockGateway | tests/processor_test.go | Test double

## Cara Menduplikasi

```sh
cd labs/16-dependency-injection
go test -v ./...    # 6 test cases, semua PASS
go test -race ./... # race-detector clean
go run ./cmd/demo   # jalankan demo DI vs Service Locator
```

## Hasil Verifikasi

6: Test Cases (3 untuk Processor • 3 untuk BadProcessor) 0: Failures (Semua lolos) PASS: Race Detector (go test -race bersih)

## Sebelum vs Sesudah DI

Processor membuat gateway sendiri. Hard-coded pada RealGateway. Test lambat, butuh network. Processor menerima interface. Swappable lewat constructor. Test cepat, mock-only.

## Source

[Software Engineering Lab 16 - Dependency Injection](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/16-dependency-injection), snapshot commit `8772235`. Source: README, implementation, tests, and execution records. Code excerpts belong in their original package context.
