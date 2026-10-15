---
title: "Dependency Injection & Inversion of Control: Separating Configuration from Use in Go"
published: false
description: "Explore Dependency Injection & Inversion of Control in Go: mechanics, architecture, caveats, and testable invariants."
tags: go, testing, performance, softwareengineering
---

# Dependency Injection & Inversion of Control: Separating Configuration from Use in Go

A production system must remain resilient under load. Lab 16 makes that behavior visible with a focused Go simulation and deterministic test assertions.

> Hard-coded dependencies couple business logic to infrastructure, breaking testability and evolvability.

## Hard-Coded Dependencies Mengunci Logika Bisnis pada Infrastruktur

When a class builds its own dependencies, business logic becomes locked to external infrastructure. Unit testing becomes difficult because it must rely on non-deterministic live network states.

## Code Review Manual Tidak Menangkap Dependensi Tersembunyi

Manual code review relies on human attention to catch coupling patterns. This approach does not automatically guarantee that every external dependency has been injected through constructors.

## Separation of Configuration from Use

The core principle of Dependency Injection is separating object construction from its use. The composition root assembles the object graph, while business modules only consume dependencies that are provided.

## Alur Kerja Dependency Injection Manual

Manual DI pipeline starts at the composition root instantiating concrete implementations. The constructor stores dependencies as fields, business methods validate input then delegate to the gateway.

## Constructor Injection: NewProcessor Menerima Interface

Constructor injection ensures Processor only knows the PaymentGateway interface, not the concrete RealGateway implementation. Dependencies are made explicit in the constructor signature.

```go
type Processor struct {
    gateway PaymentGateway
}

func NewProcessor(g PaymentGateway) *Processor {
    return &Processor{gateway: g}
}

func (p *Processor) ProcessPayment(amount int) error {
    if amount <= 0 {
        return errors.New("invalid amount")
    }
    m := Money{Amount: amount, Currency: "USD"}
    return p.gateway.Charge(m)
}
```

## Service Locator Menyembunyikan Dependensi Sesungguhnya

Service Locator hides actual dependencies behind a container interface. This reduces readability and makes mocking harder because tests must prepare a container instead of the dependency directly.

## Kontras Constructor Injection vs Service Locator

This table maps the fundamental differences between two approaches. Constructor injection yields explicit, testable dependencies; service locator introduces an extra indirection layer that complicates isolation.

Aspek | Constructor Injection | Service Locator
|---|---|---|
Dependensi | Eksplisit di constructor | Tersembunyi di container
Mocking | Langsung inject mock | Harus mock container + get()
Kekekalan | Field immutable setelah konstruksi | Bisa berubah jika container berubah
Kompleksitas | Mudah dipahami | Menambah indirection layer
Standar Industri | Diproleh sebagai best practice | Disebabkan sebagai anti-pattern

## Keputusan Desain DI pada Lab Ini

This lab chooses constructor injection as the default pattern because it provides immutability and dependency clarity. Service Locator is used only to demonstrate the anti-pattern, not as a solution.

## 6 Test Case Lolos dengan Race Detector

Terminal output proves all six test cases pass including race detector. Tests successfully isolate Processor from real gateway, verifying happy path, error propagation, and input validation.

```go
=== RUN   TestProcessor_Success
--- PASS: TestProcessor_Success (0.00s)
=== RUN   TestProcessor_GatewayError
--- PASS: TestProcessor_GatewayError (0.00s)
=== RUN   TestProcessor_InvalidAmount
--- PASS: TestProcessor_InvalidAmount (0.00s)
=== RUN   TestBadProcessor_Success
--- PASS: TestBadProcessor_Success (0.00s)
=== RUN   TestBadProcessor_GatewayError
--- PASS: TestBadProcessor_GatewayError (0.00s)
=== RUN   TestBadProcessor_InvalidAmount
--- PASS: TestBadProcessor_InvalidAmount (0.00s)
PASS
ok      lab16/tests   (cached)
```

## Batasan Bukti Lab: Yang Dibuktikan vs Belum

This lab proves isolation and error propagation in simple scenarios. It does not prove lifecycle management, performance overhead, or empirical impact on software engineering metrics.

## Trade-off Penerapan DI Manual vs Framework

Manual DI implementation requires a clear composition root. At large scale, boilerplate overhead can become significant, but this trade-off is often accepted for clarity and testability.

## Race Safety pada Implementasi DI

Go's race detector checks concurrent memory accesses. This DI implementation uses no goroutines or shared state, so tests run cleanly without race condition warnings.

```go
func TestProcessor_Success(t *testing.T) {
    mock := &MockGateway{}
    proc := di.NewProcessor(mock)
    err := proc.ProcessPayment(50)
    if err != nil {
        t.Fatalf("expected nil error, got %v", err)
    }
    // Race detector passes; no shared mutable state
}
```

## Transformasi: Dari Hard-Coded ke Constructor Injection

Before DI, business classes directly instantiate concrete dependencies, locking to one implementation. After DI, dependencies come from outside, enabling substitution and isolated testing.

## Interface PaymentGateway: Kontrak Tunggal

The PaymentGateway interface defines a single contract for all gateway implementations. Money is a value object representing currency with amount and currency fields.

```go
type Money struct {
    Amount   int
    Currency string
}

type PaymentGateway interface {
    Charge(m Money) error
}
```

## Money: Value Object yang Diinstansiasi Langsung

Money is created directly inside the method because it is a pure data structure without infrastructure dependencies. Using DI for value objects would add abstraction without benefit.

```go
func (p *Processor) ProcessPayment(amount int) error {
    if amount <= 0 {
        return errors.New("invalid amount")
    }
    m := Money{Amount: amount, Currency: "USD"}
    return p.gateway.Charge(m)
}
```

## MockGateway: Test Double untuk Isolasi Pengujian

MockGateway records Charge calls and can be configured to fail. This enables testing error propagation without real network connections.

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

## Composition Root: main.go sebagai Titik Penyambungan

The composition root is the main function responsible for assembling the object graph. It is the only place that imports the concrete RealGateway implementation.

```go
func main() {
    realGateway := &di.RealGateway{}
    processor := di.NewProcessor(realGateway)
    fmt.Println("--- Running Constructor Injection ---")
    err := processor.ProcessPayment(100)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}
```

## Validasi Input Sebelum Delegasi ke Gateway

Input validation occurs before delegating to external dependencies. This prevents unnecessary network calls and ensures only valid amounts are processed.

```go
func (p *Processor) ProcessPayment(amount int) error {
    if amount <= 0 {
        return errors.New("invalid amount")
    }
    // ... create Money and call gateway
}
```

## Checklist Penerapan Dependency Injection

This checklist summarizes the DI principles that must be met to keep code testable, decoupled, and easily evolving.

## Anti-Pattern Umum dalam Dependency Injection

Service Locator is an anti-pattern because it hides dependencies. Value objects should not be injected. Constructor bloat indicates a violation of the Single Responsibility Principle.

## Reproduksi Mandiri di Mesin Lokal Menggunakan Go

All code uses the Go standard library with zero external dependencies. The test suite covers isolation, error propagation, and input validation.

```sh
cd labs/16-dependency-injection
go test -v ./...    # jalankan semua test case
go test -race ./... # verifikasi race safety
go run ./cmd/demo   # eksekusi demo constructor injection & service locator
```

## Prinsip Inti: Dependency Harus Eksplisit

DI principles emphasize dependency clarity. By separating configuration from use, systems become easier to test, understand, and evolve.

## Source

[Software Engineering Lab 16 - Dependency Injection & Inversion of Control](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/16-dependency-injection), snapshot commit `3d936ba`. Source: README, implementation, tests, and execution records. Code excerpts belong in their original package context.
