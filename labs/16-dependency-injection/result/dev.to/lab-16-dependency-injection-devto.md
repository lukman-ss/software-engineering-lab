---
title: "Dependency Injection & Inversion of Control: Separating Configuration from Use: Constructor Injection vs Service Locator"
published: false
description: "Explore Dependency Injection & Inversion of Control in Go: mechanics, architecture, caveats, and testable invariants."
tags: go, testing, performance, softwareengineering
---

# Dependency Injection & Inversion of Control: Separating Configuration from Use: Constructor Injection vs Service Locator

A production system must remain resilient under load. Lab 16 makes that behavior visible with a focused Go simulation and deterministic test assertions.

> Hard-coded dependencies fuse business logic directly to infrastructure.

## Coupling Keras Merusak Testability & Evolusi

:  () :  () :  ()

## Mengapa Global State & Manual Swap Gagal

:  () :  () :  ()

## Pemisahan Konfigurasi dari Penggunaan

:  () :  () :  ()

## Alur Dependensi: Composition Root ke Gateway

cmd/demo/main.go | Composition root: inisialisasi RealGateway & Processor PaymentGateway Interface | Kontrak abstrak: Charge(Money) error Processor.ProcessPayment | Validasi domain amount > 0 & bungkus Money MockGateway / RealGateway | Eksekusi jaringan nyata atau mock test terisolasi

## Constructor Injection: Valid at Birth

```go
type Processor struct {
    gateway PaymentGateway // unexported field
}

func NewProcessor(g PaymentGateway) *Processor {
    return &Processor{gateway: g}
}

func (p *Processor) ProcessPayment(amount int) error {
    if amount <= 0 { return errors.New("invalid amount") }
    return p.gateway.Charge(Money{Amount: amount, Currency: "USD"})
}
```

## Failure Path: Validasi Dini Mencegah I/O Sia-sia

```go
Client Input: ProcessPayment(-10)
  -> Guard Clause: if amount <= 0 (VALIDATION FAILURE)
  -> Return: errors.New("invalid amount")
  -> Short-Circuit: Gateway.Charge() TIDAK PERNAH DIPANGGIL

Gateway Down: MockGateway.ShouldFail = true
  -> Input Valid: ProcessPayment(50) -> Gateway.Charge(50 USD)
  -> Infrastructure Error: "gateway unavailable" dipropagasi utuh
```

## Tabel Komparasi: Tiga Pola Pengelolaan Dependensi

Dimensi | Hard-Coded (new) | Service Locator | Constructor Injection
|---|---|---|---|
Visibilitas Dependensi | Terkunci internal | Opaque / Tersembunyi | Eksplisit di constructor
Validasi Compile-Time | Ya (tapi kaku) | Tidak (runtime lookup) | Ya (interface-checked)
Isolasi Unit Test | Sangat Sulit | Perlu Mock Container | Mudah (Mock Interface)
Kopling ke Framework | Nol | Tinggi (Container API) | Nol (POPO / Standard Go)

## Decision Sequence: Kapan Menggunakan Pola Tertentu

:  () :  () :  ()

## Bukti Eksekusi: Demo Wiring & Test Suite

```sh
$ go run ./cmd/demo
--- Running Constructor Injection ---
RealGateway charging 100 USD
--- Running Service Locator ---
RealGateway charging 200 USD

$ go test -v -count=1 ./...
PASS: TestProcessor_Success (0.00s) | TestProcessor_GatewayError (0.00s)
PASS: TestProcessor_InvalidAmount (0.00s) | TestBadProcessor_Success (0.00s)
ok   lab16/tests   0.387s (all 6 tests PASS, race-clean)
```

## Analisis Batasan Bukti: Fakta vs Ekstrapolasi

Mock injection mengisolasi test tanpa network call (0.00s). Constructor injection menjamin valid-at-birth compile-time. Guard clause mencegah pemanggilan gateway saat input invalid. Value object (Money) dapat dipakai langsung tanpa DI. Lifecycle container (singleton vs scoped vs transient). Performa overhead auto-wiring framework (wire / dig). Klaim penurunan defect rate sistem production secara empiris. Batas 12 parameter constructor (hanya heuristik lab).

## Trade-off Arsitektur: Eksplisit vs Fleksibilitas

Transparansi total: semua kebutuhan terlihat di signature. Immutabilitas: dependensi disetel saat inisialisasi dan unexported. Compile-time safety: compiler Go menolak jika dependensi kurang. Boilerplate konstruktor bertambah di composition root (main.go). Perubahan dependensi memerlukan modifikasi signature konstruktor. Tracing alur kode memerlukan navigasi manual tanpa magic wiring.

## Pipeline Validasi & Keamanan Konkurensi

```sh
Goroutine Client Request -> NewProcessor(mock)
  -> ProcessPayment(amount) [Stateless Validation]
  -> if amount <= 0 -> Instant Error (No I/O)
  -> Value Object: Money{Amount, USD} [Immutable, Stack-allocated]
  -> Gateway.Charge(Money) [Thread-Safe Interface Call]
Race Detector Verification: go test -race ./... -> 0 warnings
```

## Sebelum vs Sesudah Constructor Injection

Processor memanggil new RealGateway() atau container.Get(). Test unit wajib menjalankan server mock HTTP atau mock container. Signature NewBadProcessor(Container) menyembunyikan kebutuhan gateway. Missing dependency meledak di runtime saat method dipanggil. Processor hanya bergantung pada interface PaymentGateway. Test unit menyuntikkan MockGateway in-memory tanpa network (0.00s). Signature NewProcessor(PaymentGateway) eksplisit dan terverifikasi. Missing dependency gagal di compile-time jika argumen kosong.

## Bedah Anatomi: Service Locator Anti-Pattern

```go
type BadProcessor struct {
    container Container // Anti-Pattern: Mengikat ke container API
}

func (p *BadProcessor) ProcessPayment(amount int) error {
    // Opaque lookup: pembaca kode tidak tahu dependensi asli
    return p.container.GetPaymentGateway().Charge(m)
}
```

## Value Objects Tidak Boleh Menggunakan DI

:  () :  () :  ()

## Strategi Mocking: Pengujian Cepat & Terisolasi

```go
type MockGateway struct {
    ChargedMoney di.Money
    ShouldFail   bool
}

func (m *MockGateway) Charge(money di.Money) error {
    if m.ShouldFail { return errors.New("gateway unavailable") }
    m.ChargedMoney = money
    return nil
}
```

## Batas Integritas: Apa yang Tidak Diuji

Ketahanan koneksi TCP/HTTP saat RealGateway menghadapi latency spike. Penyelesaian cyclic dependency kompleks antar puluhan service. Manajemen memori lifecycle singleton vs transient under high load. Kecepatan eksekusi compile-time vs reflection-based DI containers.

## Risiko Produksi: Constructor Bloat & Over-Injection

:  () :  () :  ()

## Checklist Adopsi Dependency Injection

Deklarasikan dependensi sebagai interface minimal di dekat consumer. Gunakan Constructor Injection sebagai mekanisme default (valid-at-birth). Pusatkan instansiasi konkret di Composition Root (main.go). Larang passing Container atau Service Locator ke domain logic. Instansiasi Value Object (Money, DTO) secara langsung tanpa DI. Pastikan semua unit test berjalan terisolasi menggunakan Test Double.

## Anti-Pattern: Jebakan Umum dalam Praktik DI

Framework Lock-in: Memasukkan decorator / tag framework ke domain core. Service Locator Terselubung: Mengoper map[string]any sebagai dependency bag. Interface Bloat: Membuat interface 1-to-1 untuk struct yang tidak pernah di-swap. Ambient Context: Menggunakan global singleton yang dimutasi saat testing. Inject mock yang mengembalikan error -> verifikasi penanganan error. Kirim amount <= 0 -> verifikasi gateway mock tidak menerima hit sama sekali.

## Reproduksi Mandiri & Command Runner

```sh
cd labs/16-dependency-injection
go test -v ./...       # Menjalankan 6 test cases (mock isolation)
go test -race ./...    # Verifikasi data race pada goroutine
go run ./cmd/demo      # Menjalankan demo perbandingan kedua pola
```

## Prinsip Inti: Program to Interfaces, Not Implementations

:  () :  () :  ()

## Source

[Software Engineering Lab 16 - Dependency Injection & Inversion of Control](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/16-dependency-injection), snapshot commit `9dd0553`. Source: README, implementation, tests, and execution records. Code excerpts belong in their original package context.
