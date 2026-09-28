# Diagrams

## Diagram 1 — Orchestration Saga Flow

```text
[ Client / API ]
       │
       ▼
┌──────────────────────────────────────────────┐
│           Orchestrator (Coordinator)         │
└──────┬──────────────┬──────────────┬─────────┘
       │              │              │
       ▼ (1)          ▼ (2)          ▼ (3)
┌──────────────┐┌──────────────┐┌──────────────┐
│ OrderService ││PaymentService││InventorySvc  │
│ (CreateOrder)││ (ProcessPay) ││ (ReserveInv) │
└──────────────┘└──────────────┘└──────────────┘
```

## Diagram 2 — Failure Rollback (LIFO Compensation)

```text
Step 1: CreateOrder ────► [ SUCCESS (Logged) ]
Step 2: ProcessPay  ────► [ SUCCESS (Logged) ]
Step 3: ReserveInv  ────► [ FAILED  (Out of stock) ]
                              │
                              ▼ Trigger Compensations (LIFO)
                      ┌──────────────────────┐
                      │ Refund Payment       │ (Reverses Step 2)
                      └──────────────────────┘
                              │
                              ▼
                      ┌──────────────────────┐
                      │ Cancel Order         │ (Reverses Step 1)
                      └──────────────────────┘
```

## Diagram 3 — Choreography Event-Driven Flow

```text
[ OrderService ] ──( OrderCreated )──► [ EventBus ] ──► [ PaymentService ]
                                                              │
                                                       ( PaymentCompleted )
                                                              │
                                                              ▼
[ InventoryService ] ◄──( ReserveInventory )─────────── [ EventBus ]
```
