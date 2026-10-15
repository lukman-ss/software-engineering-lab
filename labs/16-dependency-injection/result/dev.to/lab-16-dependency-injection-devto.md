---
title: "Dependency Injection & Inversion of Control: Separating Configuration from Use for Isolated Testing in Go"
published: false
description: "Explore Dependency Injection & Inversion of Control in Go: mechanics, architecture, caveats, and testable invariants."
tags: go, testing, performance, softwareengineering
---

# Dependency Injection & Inversion of Control: Separating Configuration from Use for Isolated Testing in Go

A production system must remain resilient under load. Lab 16 makes that behavior visible with a focused Go simulation and deterministic test assertions.

> Direct object instantiation inside business logic locks domain code to external infrastructure.

## Instansiasi Langsung di Logika Bisnis Mengunci Dependensi

Hardcoding concrete service instantiations inside business methods couples domain logic directly to external infrastructure. The code may run successfully initially; this does not mean the architecture can evolve cleanly when vendor integrations change.

## Service Locator Terlihat Bersih Namun Menyembunyikan Keterikatan

Teams often adopt Service Locators believing they eliminate coupling since direct instantiations disappear. In practice, injecting a container does not automatically decouple code when underlying dependencies remain hidden behind runtime lookups.

## Pemisahan Konfigurasi dari Penggunaan: Dua Fase Terisolasi

Fowler's foundational principle establishes two separate phases: object graph assembly at the composition root and business execution within domain logic. Only along this path do business components remain completely insulated from infrastructure details.

## Alur Eksekusi: Composition Root ke Processor dan Gateway

Execution initiates with concrete assembly at the composition root, followed by local input guards prior to delegating to the gateway interface. This architecture cleanly isolates logic, though it has not yet proven resiliency under high network latency.

## Constructor Injection: Dependensi Eksplisit dan Valid Saat Lahir

Processor stores its gateway interface in an unexported field, accepting it strictly during initialization. Relying on interfaces does not imply zero-cost abstractions; it deliberately accepts interface pointer overhead to achieve complete domain isolation.

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
    return p.gateway.Charge(m)
}
```

## Kegagalan Service Locator: Keterikatan Tersembunyi Container

The Service Locator masks actual dependencies inside method bodies, risking runtime panics when resolution fails. Only along this path are unit tests forced to construct irrelevant container mocks instead of directly supplying mock services.

## Matriks Komparasi: Constructor Injection vs Service Locator

This matrix contrasts Constructor Injection against the Service Locator pattern. Constructor injection guarantees compile-time type safety and explicit dependencies, whereas service locators defer missing-dependency detection to runtime.

Dimensi Evaluasi | Constructor Injection (Processor) | Service Locator (BadProcessor)
|---|---|---|
Deklarasi Dependensi | Eksplisit di parameter NewProcessor | Tersembunyi di balik Container
Type Safety & Validitas | Compile-time check, valid at birth | Runtime resolution, rawan nil panic
Beban Mocking Pengujian | Minimal: langsung injeksi MockGateway | Ganda: wajib MockContainer + MockGateway
Kopling Domain Logic | Nol ke framework / container | Terikat langsung ke interface Container
Rekomendasi Standar | Best practice industri modern | Discouraged (PSR-11: SHOULD NOT)

## Keputusan Arsitektur: Memilih Pola Perakitan Dependensi

The architectural sequence establishes constructor injection as default while directly instantiating plain value objects. Choosing manual wiring in main.go does not reflect anti-framework dogma; it reflects a pragmatic commitment to Go simplicity.

## Eksekusi Demo cmd/demo dan Test Suite Lolos Bersih

Verification demonstrates all 6 unit test scenarios passing deterministically in sub-second execution. The demo run confirms both injection mechanisms process payments successfully without triggering data races.

```go
=== RUN   TestProcessor_Success
--- PASS: TestProcessor_Success (0.00s)
=== RUN   TestProcessor_GatewayError
--- PASS: TestProcessor_GatewayError (0.00s)
=== RUN   TestProcessor_InvalidAmount
--- PASS: TestProcessor_InvalidAmount (0.00s)
=== RUN   TestBadProcessor_Success / GatewayError / InvalidAmount
--- PASS: TestBadProcessor_* (0.00s)
PASS: ok lab16/tests 1.356s (race detector: clean)
--- Running Constructor Injection: RealGateway charging 100 USD
--- Running Service Locator: RealGateway charging 200 USD
```

## Batasan Bukti Lab: Yang Dibuktikan vs Belum

This lab proves isolated unit testing and input guard execution prior to external infrastructure calls. It does not prove automated lifecycle resolution or empirical claims regarding defect rate reduction in production.

## Trade-off Penerapan Dependency Injection

Adopting DI introduces cognitive indirection because concrete implementations are detached from invocation sites. Introducing interfaces and assemblers does not automatically improve architecture if the domain does not require implementation swapping.

## Keamanan Konkurensi pada Komponen Stateless

The Processor struct remains stateless post-construction, rendering it inherently safe for concurrent execution across goroutines. While this simulation is not equivalent to distributed connection pools under network load, it passes race checks cleanly.

```go
// Processor stateless: gateway bersifat read-only setelah konstruksi
type Processor struct {
    gateway PaymentGateway
}
// Aman dipanggil lintas goroutine selama gateway thread-safe
func (p *Processor) ProcessPayment(amount int) error {
    if amount <= 0 { return errors.New("invalid amount") }
    m := Money{Amount: amount, Currency: "USD"}
    return p.gateway.Charge(m)
}
```

## Transformasi Desain: Hardcoded vs Constructor Injection

Prior to DI, business workflows directly instantiated concrete network gateways, making offline unit testing impossible. Refactoring to constructor injection does not mean adding unnecessary complexity; it frees domain models from external details.

## Interface Minimal dan Pengecualian untuk Value Object

The PaymentGateway interface declares a single Charge method while Money is instantiated directly as a value object. Scoping DI strictly to services rather than value objects reflects design boundaries validated within this lab.

```go
// Money adalah value object data murni: instansiasi langsung
type Money struct {
    Amount   int
    Currency string
}
// PaymentGateway mendefinisikan kontrak minimal abstraction
type PaymentGateway interface {
    Charge(m Money) error
}
```

## Pengujian Terisolasi Menggunakan Test Double MockGateway

MockGateway simulates network outages and records transaction payloads without touching HTTP sockets. Passing unit tests with mocks does not automatically guarantee third-party integration integrity if production API schemas change.

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

## Validasi Pre-Check Mencegah Panggilan Jaringan Sia-sia

Tests confirm negative amounts are rejected at the domain boundary prior to contacting the gateway. Only along this pathway does the mock assert that Amount remains 0, confirming wasteful network invocations were prevented.

```go
func TestProcessor_InvalidAmount(t *testing.T) {
    mock := &MockGateway{}
    proc := di.NewProcessor(mock)
    err := proc.ProcessPayment(-10)
    if err == nil { t.Fatalf("expected error for negative amount") }
    if mock.ChargedMoney.Amount != 0 {
        t.Errorf("gateway should not have been called, got %d",
            mock.ChargedMoney.Amount)
    }
}
```

## Propagasi Error Transparan dari Lapisan Infrastruktur

The Processor bubbles gateway errors transparently back to the caller. Propagating raw errors does not imply the system lacks retry policies; it confirms interface boundaries accurately convey downstream failure states.

```go
func TestProcessor_GatewayError(t *testing.T) {
    mock := &MockGateway{ShouldFail: true}
    proc := di.NewProcessor(mock)
    err := proc.ProcessPayment(50)
    if err == nil {
        t.Fatalf("expected error from gateway, got nil")
    }
    // Verifikasi error gateway sampai ke pemanggil
}
```

## Composition Root di cmd/demo: Satu Titik Perakitan Objek

The main entrypoint functions as the sole composition root where concrete structs are instantiated and wired. This minimal structure avoids third-party DI reflection frameworks, with its scope strictly matching lab verification needs.

```go
func main() {
    // Instansiasi dependensi konkret
    realGateway := &di.RealGateway{}
    // Constructor Injection eksplisit
    processor := di.NewProcessor(realGateway)
    _ = processor.ProcessPayment(100)
    // Service Locator anti-pattern untuk demonstrasi
    container := &SimpleContainer{gateway: realGateway}
    badProcessor := di.NewBadProcessor(container)
    _ = badProcessor.ProcessPayment(200)
}
```

## Checklist Penerapan Dependency Injection di Go

This checklist enforces robust DI practices: interface-driven contracts, explicit constructor injection, centralized composition root, container prohibition in domain logic, and isolated mock testing.

## Anti-Pattern Umum dalam Penerapan Dependency Injection

Passing containers into domain logic conceals actual dependencies, while parameter-bloated constructors signal SRP violations. Managing simple value objects inside DI containers does not make code cleaner; it adds pointless ceremony.

## Reproduksi Mandiri di Mesin Lokal Menggunakan Go

The entire test harness relies solely on the pure Go standard library with zero third-party dependencies. While this validates architectural isolation locally, production systems may still require code-generation tools like wire at enterprise scale.

```sh
cd labs/16-dependency-injection
go test -v ./...    # validasi unit test terisolasi Processor dan BadProcessor
go test -race ./... # verifikasi konkurensi stateless bebas race condition
go run ./cmd/demo   # eksekusi komparasi Constructor Injection vs Service Locator
```

## Prinsip Inti: Pisahkan Konfigurasi, Eksplisitkan Kontrak

Dependency Injection fundamentally separates configuration concerns from operational business logic. Leveraging interfaces and constructors does not automatically fix architectural debt if domain boundaries are poorly defined initially.

## Source

[Software Engineering Lab 16 - Dependency Injection & Inversion of Control](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/16-dependency-injection), snapshot commit `5563838`. Source: README, implementation, tests, and execution records. Code excerpts belong in their original package context.
