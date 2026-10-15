---
title: "Dependency Injection & Inversion of Control: Separating Object Construction from Domain Logic in Go"
published: false
description: "Explore Dependency Injection & Inversion of Control in Go: mechanics, architecture, caveats, and testable invariants."
tags: go, testing, performance, softwareengineering
---

# Dependency Injection & Inversion of Control: Separating Object Construction from Domain Logic in Go

A production system must remain resilient under load. Lab 16 makes that behavior visible with a focused Go simulation and deterministic test assertions.

> Hardcoded instantiation locks infrastructure into domain rules.

## Instansiasi Langsung Mengunci Logika Domain

Hardcoding RealGateway inside Processor binds domain logic to infrastructure. This does not mean instant failure, but forces every test into a brittle, network-dependent integration run.

## Global Singleton dan Flag Mengaburkan Kontrak

Global singletons do not automatically solve coupling. They obscure dependencies from the caller signature and introduce state pollution across concurrent test executions.

## Pisahkan Perakitan Objek dari Eksekusi Bisnis

Separation of configuration from use isolates assembly from runtime logic. Components are wired at the composition root and execute strictly against interfaces.

## IoC adalah Arah Kendali, DI adalah Bentuk Spesifik

IoC inverts control flow broadly, while DI specifically externalizes dependency injection. Adopting DI does not imply handing total application control to a framework.

## Abstraksi Gateway dan Value Object Tanpa DI

PaymentGateway defines a minimal contract while Money remains a simple value object. Defining an interface isolates the consumer, but does not prove network resiliency on its own.

```go
type Money struct { Amount int; Currency string }

type PaymentGateway interface {
    Charge(m Money) error
}

type RealGateway struct{}
func (g *RealGateway) Charge(m Money) error {
    fmt.Printf("RealGateway charging %d %s\n", m.Amount, m.Currency)
    return nil
}
```

## Constructor Injection Membuat Dependensi Terlihat

NewProcessor explicitly demands its PaymentGateway dependency at birth. While the dependency is transparent, an interface parameter does not automatically guard against nil values.

```go
type Processor struct { gateway PaymentGateway }

func NewProcessor(g PaymentGateway) *Processor {
    return &Processor{gateway: g}
}

func (p *Processor) ProcessPayment(amount int) error {
    if amount <= 0 { return errors.New("invalid amount") }
    m := Money{Amount: amount, Currency: "USD"}
    return p.gateway.Charge(m)
}
```

## Validasi Domain Berjalan Sebelum Panggilan Gateway

DI does not replace input validation. Pre-checking amount <= 0 ensures infrastructure is never invoked on invalid data, guarding against unnecessary network trips.

## Service Locator Menyembunyikan Dependensi Nyata

BadProcessor takes a generic Container instead of PaymentGateway. The code compiles, but this does not mean it is decoupled: actual dependencies are hidden behind service lookups.

```go
type Container interface { GetPaymentGateway() PaymentGateway }

type BadProcessor struct { container Container }

func NewBadProcessor(c Container) *BadProcessor {
    return &BadProcessor{container: c}
}

func (p *BadProcessor) ProcessPayment(amount int) error {
    if amount <= 0 { return errors.New("invalid amount") }
    m := Money{Amount: amount, Currency: "USD"}
    return p.container.GetPaymentGateway().Charge(m)
}
```

## Perbandingan Struktural DI vs Service Locator

Constructor DI makes component dependencies transparent at compile-time, whereas Service Locator forces tests to construct artificial mock containers to satisfy queries.

Dimensi | Constructor Injection | Service Locator
|---|---|---|
Kontrak Signature | Eksplisit (PaymentGateway) | Opaque (Container)
Ketergantungan API | Interface Domain Murni | Terikat Container Interface
Mocking di Test | Inject Mock Langsung | Wajib Mock Container & Service
Deteksi Compile | Error jika interface hilang | Bisa lolos tapi lookup runtime

## Composition Root Menjadi Titik Tunggal Perakitan

main.go serves as the application composition root. This simulation is not equivalent to an enterprise DI container; it demonstrates pure, explicit manual wiring.

```go
func main() {
    realGateway := &di.RealGateway{}

    // Constructor Injection (Bersih & Eksplisit)
    processor := di.NewProcessor(realGateway)
    _ = processor.ProcessPayment(100)

    // Service Locator (Anti-Pattern Demo)
    container := &SimpleContainer{gateway: realGateway}
    badProcessor := di.NewBadProcessor(container)
    _ = badProcessor.ProcessPayment(200)
}
```

## Test Double Menggantikan Infrastruktur Jaringan

MockGateway stands in for the concrete gateway during unit tests. The ShouldFail flag simulates transport failures deterministically within the scoped test conditions.

```go
type MockGateway struct {
    ChargedMoney di.Money
    ShouldFail   bool
}

func (m *MockGateway) Charge(money di.Money) error {
    if m.ShouldFail {
        return errors.New("gateway unavailable")
    }
    m.ChargedMoney = money
    return nil
}
```

## Enam Kasus Uji PASS Termasuk Race Detector

All 6 test cases executed in 0.00s without external infrastructure. Race detector confirmed zero data races across the simulated test scenarios.

```go
=== RUN   TestProcessor_Success           --- PASS (0.00s)
=== RUN   TestProcessor_GatewayError      --- PASS (0.00s)
=== RUN   TestProcessor_InvalidAmount     --- PASS (0.00s)
=== RUN   TestBadProcessor_Success        --- PASS (0.00s)
=== RUN   TestBadProcessor_GatewayError   --- PASS (0.00s)
=== RUN   TestBadProcessor_InvalidAmount  --- PASS (0.00s)
PASS
ok  github.com/lukman-ss/.../tests  1.034s
```

## Batas Bukti Lab: Yang Dibuktikan vs Belum

This lab proves interface decoupling and test double substitution. It has not proven automated DI container overhead or production lifecycle management strategies.

## Value Object Diinstansiasi Langsung Tanpa DI

Injecting pure data structures creates unnecessary abstraction layers. Value objects should be instantiated directly; this does not violate core decoupling principles.

## Parameter Terlalu Banyak Menandakan SRP Violation

A bloated constructor signature signals Single Responsibility violations. Hiding parameters inside a container does not automatically fix design flaws; decomposing the class is required.

## Rekomendasi Industri: Hindari Menginjeksi Container

PSR-11 explicitly discourages passing containers into objects via SHOULD NOT. This is a strong architectural guideline against coupling domain entities to runtime lookup tools.

## Race Clean di Test Bukan Bukti Lifecycle Production Aman

A race-free unit test verifies current execution paths. However, production systems still require rigorous concurrency validation whenever dependencies share pooled connections.

## Biaya Tersembunyi: Kapan DI Tidak Diperlukan?

DI introduces architectural overhead and complicates code tracing. Introducing abstractions does not automatically provide value when components have only one stable implementation.

## Checklist Penerapan Dependency Injection di Go

This checklist enforces DI best practices: consume interfaces, require explicit constructors, localize wiring in the composition root, and bypass DI for value types.

## Anti-Pattern Umum dalam Penerapan DI

Treating the container as a universal parameter destroys encapsulation. Creating 1:1 interfaces for every struct does not mean clean code; it produces severe interface pollution.

## Reproduksi Mandiri di Mesin Lokal

This lab is self-contained with zero third-party dependencies. Running these commands validates happy paths, error propagation, input pre-checks, and manual wiring mechanics.

```sh
cd labs/16-dependency-injection
go test -v ./...    # validasi 6 skenario pengujian
go test -race ./... # verifikasi race detector
go run ./cmd/demo   # simulasi manual wiring & locator
```

## Kesimpulan Inti: Konstruksi Terpisah dari Penggunaan

DI cleanly decouples object construction from domain execution. While it enables deterministic mock isolation, DI does not automatically guarantee overall runtime correctness.

## Source

[Software Engineering Lab 16 - Dependency Injection & Inversion of Control](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/16-dependency-injection), snapshot commit `9c271ae`. Source: README, implementation, tests, and execution records. Code excerpts belong in their original package context.
