---
title: "Saga Pattern: Distributed Transaction Management Without Two-Phase Commit in Microservices"
published: false
description: "Explore Saga Pattern in Go: mechanics, architecture, caveats, and testable invariants."
tags: go, testing, performance, softwareengineering
---

# Saga Pattern: Distributed Transaction Management Without Two-Phase Commit in Microservices

A production system must remain resilient under load. Lab 29 makes that behavior visible with a focused Go simulation and deterministic test assertions.

> When services own their data, distributed transactions need local independence plus compensation.

## Ketika Database Terpisah Tidak Bisa Berkomitmen Bersama

In modern microservice architectures, each service manages its own database. When a business flow involves order creation, payment processing, and inventory reservation sequentially, classical ACID transactions relying on a single database boundary break down. Two-phase commit emerges as a theoretical solution, but it locks all resources synchronously until every participant votes commit — causing high latency and vulnerability to network partitions.

## Mengapa Pendekatan Naif Gagal Menangani Transaksi Multi-Layanan

The naive implementation attempts sequential service calls without cross-domain rollback mechanisms. If payment fails after the order is created, the order database records a PENDING status while payment records failure. The system enters an inconsistent state with no automatic recovery path. This manual approach also hinders observability since there is no centralized execution log recording each step's status.

```go
func (o *OrderService) CreateFullOrder(orderID string) error {
    o.CreateOrder(orderID)     // DB-commit #1
    p.ProcessPayment(orderID)  // DB-commit #2 — jika gagal?
    i.ReserveInventory(...)    // DB-commit #3 — jika gagal?
    o.ApproveOrder(orderID)    // DB-commit #4
}
```

## Mental Model: Estafet Transaksi dengan Kompensasi Mundur

Imagine a relay race where each runner completes their leg and immediately records results to a local database. If the third runner stumbles, the first and second runners run backward to restore their original positions. Each step stands alone as a local transaction with immediate commit — no cross-service locking exists. When failure occurs, the orchestrator triggers a Last-In-First-Out chain of compensations to restore consistency.

## Pipeline Orkestrasi: Dari Langkah Eksekusi Hingga Kompensasi

The orchestration pipeline begins when a client submits a request to the orchestrator. The orchestrator iterates through steps sequentially, executing each into the respective local service database. On each success, the step is recorded to the executed stack. If any step fails, the loop breaks and the orchestrator calls compensate(), which unwinds backward from the stack — ensuring refund happens before cancel order. Each step's status log is recorded atomically under mutex for concurrency consistency.

## Implementasi Inti: Orchestrator Step & LIFO Compensate

The Step data structure defines an Execute-Compensate function pair for each local transaction. The orchestrator stores this Step slice and executes them sequentially. The compensate() function iterates backward from the last index to zero — this is the core of LIFO rollback. If compensate() itself fails on a step, the error is aggregated into a compErrors slice and reported alongside the primary step error.

```go
type Step struct {
    Name       string
    Execute    func(ctx context.Context) error
    Compensate func(ctx context.Context) error
}

func (o *Orchestrator) compensate(ctx context.Context, executed []Step) error {
    for i := len(executed) - 1; i >= 0; i-- {
        step := executed[i]
        if step.Compensate != nil {
            err := step.Compensate(ctx)  // LIFO unwind
            ...
        }
    }
}
```

## Failure Path: ReserveInventory Gagal — Urutan Kompensasi LIFO

This scenario demonstrates a realistic failure path where inventory reservation fails due to insufficient stock. The order has been created and payment processed before the failure is detected. The orchestrator then triggers LIFO compensation: first processing payment refund, then cancelling the order. Through this order, the system returns to its pre-saga state — payment refund occurs before order cancellation, maintaining financial sequence consistency.

## Tabel Perbandingan: Mekanisme Kompensasi vs Gejala Sistem

This table maps the three saga transaction types from Azure Architecture Center to specific lab implementations. The compensable transaction is most common — every step except ApproveOrder has a Compensate function. Pivot transactions represent irreversible turning points, while retryable transactions guarantee success post-pivot. In this lab, ApproveOrder serves as the final terminal state step without compensation.

Tipe Langkah | Eksekusi Sukses | Eksekusi Gagal | Kompensasi
|---|---|---|---|
Compensable Transaction | Local DB commit persisten | Rollback via Compensate() | Refund, Cancel, Release
Pivot Transaction | Titik balik tidak bisa dibatalkan | Saga dijamin mencapai final state | Tidak ada kompensasi
Retryable Transaction | Post-pivot, dijamin sukses | Idempotent retry | Tidak diperlukan
CreateOrder Step | Status PENDING + semantic lock | Rollback: delete lock + CANCELLED | ✓
ProcessPayment Step | Payment recorded + idempotency key | Rollback: delete payment entry | ✓ Refund
ReserveInventory Step | Stock decremented | Rollback: stock restored via Release() | ✓
ApproveOrder Step | Status APPROVED + lock released | Tanpa kompensasi (final step) | —

## Keputusan Arsitektur: Orkestrasi vs Koreografi

The choice between orchestration and choreography is not binary. Orchestration excels when step ordering has tight dependencies and centralized logging is needed, such as in e-commerce checkout flows. Choreography suits event-driven systems where services react independently to events without knowing about a central orchestrator. This lab implements both patterns to demonstrate each trade-off.

## Bukti Eksekusi: Semua Test Suite Lolos Bersih

All nine test cases run cleanly without a single failure. The happy path proves four steps succeed, changing order status to APPROVED with stock reduced by 2 units. The failure test proves LIFO rollback: order returns to CANCELLED, payment refunded, stock restored to 1. The concurrency test with 10 simultaneous goroutines ensures mutex atomicity remains intact.

```go
=== RUN   TestOrchestrator_HappyPath
    --- PASS: TestOrchestrator_HappyPath (0.00s)
=== RUN   TestOrchestrator_FailureCompensatesLIFO
    --- PASS: TestOrchestrator_FailureCompensatesLIFO (0.00s)
=== RUN   TestPayment_Idempotency
    --- PASS: TestPayment_Idempotency (0.00s)
=== RUN   TestSemanticLock
    --- PASS: TestSemanticLock (0.00s)
=== RUN   TestOrchestrator_Concurrency
    --- PASS: TestOrchestrator_Concurrency (0.00s)
=== RUN   TestChoreography_Flow
    --- PASS: TestChoreography_Flow (0.00s)
=== RUN   TestChoreography_FailureCompensates
    --- PASS: TestChoreography_FailureCompensates (0.00s)
=== RUN   TestOrchestrator_CompensationErrorPropagated
    --- PASS: TestOrchestrator_CompensationErrorPropagated (0.00s)
=== RUN   TestOrchestrator_ContextCancellation
    --- PASS: TestOrchestrator_ContextCancellation (0.00s)
PASS: ok  labs/29-saga-pattern/tests  0.003s
```

## Batasan Bukti Lab: Yang Dibuktikan vs Belum

This lab proves computational correctness and concurrency safety in an in-memory environment. However, strict limitations must be acknowledged: there is no network partition simulation, no real database crash, and no persistence layer. In production, compensation failure can occur due to service unavailability — scenarios not covered by the current test suite. Dead-letter queues and manual reconciliation are countermeasures needed beyond this lab's scope.

## Trade-off Arsitektur: Availability vs Isolation

Saga trades high isolation for high availability. Each step commits immediately without cross-service locking, so there is no blocking between services. However, a consistency window exists between step execution and compensation — an order may appear APPROVED while payment remains unprocessed. Developers must write compensation logic explicitly, adding application complexity compared to monolithic solutions.

Dimensi | ACID (Monolitik) | Saga Pattern | Dampak Produksi
|---|---|---|---|
Isolasi (I dalam ACID) | serializable penuh | tidak ada isolasi lintas service | dirty reads, lost updates mungkin terjadi
Ketersediaan | blocking 2PC | non-blocking, each step commit segera | saga jauh lebih available di bawah load tinggi
Latensi | satu round-trip ke semua peserta | N round-trips sekuensial + kompensasi | total latency lebih tinggi per transaksi
Konsistensi | strong consistency | eventual consistency | window inconsistency antar step execute sampai kompensasi
Complexity | database engine handles it | application-level orchestration | developer wrote compensation logic manually
Observability | transaction logs | step-level logs + compensation trace | lebih transparan, lebih mudah debug

## Pipeline Konkurensi: Goroutine → Validate → Mutex → Append Error

Concurrency safety is achieved through a snapshot-and-release strategy. When Execute is called, the orchestrator copies the steps slice to a local variable while releasing the mutex, avoiding blocking during step execution. Log append and read operations are protected by individual mutexes to ensure read consistency. The concurrency test runs 10 simultaneous goroutines with the Go race detector reporting a clean result with no data races.

```go
func (o *Orchestrator) AddStep(step Step) {
    o.mu.Lock()
    defer o.mu.Unlock()
    o.steps = append(o.steps, step)
}

func (o *Orchestrator) Execute(ctx context.Context) error {
    o.mu.Lock()
    steps := make([]Step, len(o.steps))
    copy(steps, o.steps)  // snapshot tanpa lock holder
    o.mu.Unlock()
    // ...execute steps without holding mu during step calls
    o.mu.Lock()
    o.logs = append(o.logs, StepLog{Name: step.Name, Status: ...})
    o.mu.Unlock()
}
```

## Before vs After: Dari Manual Calls ke Saga Orkestrasi

The before-after comparison shows a fundamental transformation from manual state management to structured orchestration. Before saga, a failure in one service leaves persistent state inconsistency requiring manual intervention. After saga, the orchestrator automatically logs status and triggers sequential compensation. Each step has an audit trail enabling fast debugging without combining logs from multiple services.

## Service Layer: Semantic Lock, Idempotency, dan Atomic Stock Mutation

The service layer is not merely a storage abstraction — each service contains its own protection mechanisms. OrderService uses semantic locking to prevent double order creation within a single saga instance. PaymentService relies on idempotency keys to handle repeated retry on the payment step. InventoryService guarantees atomic stock mutation through a single mutex. Without these protections, saga compensation itself could introduce race conditions.

## Choreography Pattern: Event Bus Decoupled via Pub/Sub

The choreography pattern eliminates the central coordinator entirely. Each service listens to events from other services and reacts according to registered handlers. The event bus uses RWMutex for concurrency safety: Subscribe holds write locks, Publish holds read locks while copying the handler slice. This pattern yields loose coupling but reduces observability since there is no centralized execution log.

```go
type EventBus struct {
    mu       sync.RWMutex
    handlers map[EventType][]EventHandler
}

func (b *EventBus) Subscribe(t EventType, h EventHandler) {
    b.mu.Lock()
    b.handlers[t] = append(b.handlers[t], h)
    b.mu.Unlock()
}

func (b *EventBus) Publish(ctx context.Context, e Event) {
    b.mu.RLock()
    handlers := append([]EventHandler(nil), b.handlers[e.Type]...)
    b.mu.RUnlock()
    for _, h := range handlers { h(ctx, e) }
}
```

## Bukti Invarian: Idempotency Payment dan Semantic Lock Order

The idempotency test proves that the payment service accepts repeated calls with the same ID without double-charging. The semantic lock test proves that the order service rejects duplicate order creation while the saga instance is still active — preventing race conditions between concurrent saga instances using the same orderID.

```go
func TestPayment_Idempotency(t *testing.T) {
    paymentSvc := services.NewPaymentService()
    id := "tx-1"
    // First call: records payment
    paymentSvc.ProcessPayment(id, 100, false)
    // Second call with same ID: idempotent no-op
    err := paymentSvc.ProcessPayment(id, 100, false)
    // Must succeed (nil error) — double charge prevented
    if err != nil { t.Fatalf(...) }
}

func TestSemanticLock(t *testing.T) {
    orderSvc := services.NewOrderService()
    id := "ord-lock"
    orderSvc.CreateOrder(id)  // sets locks[id] = true
    err := orderSvc.CreateOrder(id)  // should FAIL
    if err == nil { t.Fatal("second create should fail due to semantic lock") }
}
```

## Bukti LIFO Log: Verifikasi Urutan Status Step pada Kompensasi

This test verifies three things simultaneously: state recovery (order CANCELLED, payment refunded), stock integrity (stock remains 1), and correct log sequence. Six log entries — two EXECUTED, one FAILED, two COMPENSATED — form a sequence proving LIFO rollback works as designed. This order is critical for debugging: engineers can reconstruct the exact failure chronology.

```go
func TestOrchestrator_FailureCompensatesLIFO(t *testing.T) {
    // ...setup: CreateOrder, ProcessPayment, ReserveInventory(fail)...
    err := orch.Execute(ctx)
    // Verify state recovery
    if orderSvc.GetOrder(orderID) != services.OrderCancelled {
        t.Fatalf("expected order cancelled, got %v", orderSvc.GetOrder(orderID))
    }
    if paymentSvc.HasPayment(orderID) {
        t.Fatal("expected payment refunded")
    }
    // Verify log order (5 entries exactly)
    logs := orch.Logs()
    expectedOrder := []saga.StepStatus{
        saga.StatusExecuted, saga.StatusExecuted, saga.StatusFailed,
        saga.StatusCompensated, saga.StatusCompensated,
    }
    for i, l := range logs {
        if l.Status != expectedOrder[i] { ... }
    }
}
```

## Bukti Context Cancellation: Kompensasi Saat Saga Dibatalkan

This test validates saga behavior when context is cancelled mid-execution. Context cancellation is detected via the select case in the Execute loop, halting iteration before the next step runs. Despite context cancellation, compensate() is still called with context.Background() — ensuring compensation is not aborted by cancellation propagation. Every executed step is guaranteed compensation, maintaining the atomicity invariant.

```go
func TestOrchestrator_ContextCancellation(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    orch := saga.NewOrchestrator()
    step1Compensated := false
    orch.AddStep(saga.Step{
        Name: "Step1",
        Execute: func(ctx context.Context) error {
            cancel() // cancel context before step 2
            return nil
        },
        Compensate: func(ctx context.Context) error {
            step1Compensated = true
            return nil
        },
    })
    orch.AddStep(saga.Step{
        Name: "Step2",
        Execute: func(ctx context.Context) error {
            t.Fatal("step 2 should not execute when context cancelled")
            return nil
        },
    })
    err := orch.Execute(ctx)
    if !step1Compensated { t.Fatal("expected step 1 compensated") }
}
```

## Checklist Implementasi Saga Pattern yang Benar

This checklist summarizes eight critical requirements for production-ready saga pattern implementation. The most commonly overlooked point is compensation idempotency — if compensate() is called twice, the system must remain stable. The sixth point about dead-letter queues is also crucial: permanent compensation failures require manual intervention mechanisms that must be planned before deployment.

## Anti-Pattern Saga: Yang Harus Dihindari di Production

These four anti-patterns are most commonly found in production saga implementations. Non-idempotent compensate is the primary cause of double-refund that harms financial reconciliation. Circular compensation dependency creates infinite loops that paralyze the orchestrator. Missing compensate on the final step means approval state persists permanently even when other steps are compensated. Blocking long-running compensate locks the orchestrator goroutine and prevents subsequent steps from executing.

## Takeaway Praktis: Command & Runbook untuk Saga Pattern

The practical takeaway for engineers adopting saga pattern in production: the commands above are the starting point for independent reproduction. The three operational principles at the bottom are monitoring guidelines that must be implemented before go-live. Compensation failures entering the DLQ must be processed manually by on-call engineers through the admin console — there is no fully automatic resolution for all cases.

```sh
# Jalankan seluruh test suite
go test -v ./...

# Verifikasi race condition
go test -race -v ./...

# Demo interaktif (happy path + rollback)
go run ./cmd/demo

# Production runbook:
# 1. Monitor orch.Logs() untuk STATUS_FAILED / STATUS_COMPENSATE_FAILED
# 2. Alert on COMPENSATE_FAILED → trigger DLQ processing
# 3. Manual reconciliation via admin console for stuck compensations
# 4. Regular healthcheck: verify all in-flight sagas reached terminal state
```

## Prinsip Arsitektural Jangka Panjang: Eventual Consistency sebagai Kontrak

These architectural principles extend far beyond the lab implementation. Sagas are not merely a technical pattern — they represent a distributed systems design philosophy that accepts eventual consistency as a fundamental contract. Each principle above guides architectural decisions: from per-service database selection to monitoring and alerting system design. Engineers who understand this will avoid the temptation to revert to monolithic ACID.

## Source

[Software Engineering Lab 29 - Saga Pattern](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/29-saga-pattern), snapshot commit `ba4b6216`. Source: README, implementation, tests, and execution records. Code excerpts belong in their original package context.
