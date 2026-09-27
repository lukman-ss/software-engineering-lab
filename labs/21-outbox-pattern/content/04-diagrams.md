# Diagrams

## Diagram 1 — System Architecture (matches engineering/01-design.md and README.md)

```text
[ Client ] -> [ Order Service ]
                    |
                    v (Single DB Tx: orders + outbox_events)
      +-----------------------------+
      |  In-Memory Tx DB            |
      |  - orders table             |
      |  - outbox_events table      |
      +-----------------------------+
                    ^
                    | Poll (GetPendingOutbox)
            [ Outbox Relay ]
                    |
                    v Publish
            [ Message Broker ]
                    |
                    v Deliver
           [ Idempotent Consumer ] -> (message_log / processedIDs map)
```

Components (all exist in lab):
- `internal/outbox/db.go` — In-Memory Tx DB
- `internal/outbox/service.go` — Order Service
- `internal/outbox/relay.go` — Outbox Relay
- `internal/outbox/broker.go` — Message Broker (mock)
- `internal/outbox/consumer.go` — Idempotent Consumer

## Diagram 2 — Dual-Write Failure Sequence (verified by TestDualWriteProblem_Failure + demo Scenario 1)

```text
Client -> Service.CreateOrderDualWriteNaive
  Service -> DB.Tx.Commit : SaveOrder("o-bug") OK
  Service -> Broker.Publish : FAIL (broker unavailable)
  Result: DB contains order, broker contains 0 messages
  => INCONSISTENT STATE
```

## Diagram 3 — Outbox Happy Path Sequence (verified by TestTransactionalOutbox_HappyPath + demo Scenario 2)

```text
Client -> Service.CreateOrderWithOutbox
  Service -> DB.Tx : SaveOrder + SaveOutbox (staged)
  Service -> DB.Tx.Commit : both flushed atomically
  Relay (ticker) -> DB.GetPendingOutbox : [evt-o-1 PENDING]
  Relay -> Broker.Publish(evt-o-1) : OK
  Relay -> DB.MarkOutboxProcessed(evt-o-1) : PROCESSED
  Consumer.Handle(evt-o-1) : accepted=true
```

## Diagram 4 — Relay Retry After Broker Failure (verified by TestTransactionalOutbox_RelayRetryAfterBrokerFailure)

```text
Service -> DB : CreateOrderWithOutbox OK (PENDING, order persisted)
Relay poll #1 -> Broker.Publish : FAIL (SetFailNext)
  => outbox stays PENDING, no status change
Relay poll #2 -> Broker.Publish : OK
  => MarkOutboxProcessed -> PROCESSED
Consumer.Handle : accepted=true
```

## Diagram 5 — Duplicate Delivery and Idempotency (verified by TestTransactionalOutbox_Idempotent_DuplicateDelivery + demo Scenario 3)

```text
Relay (crash before marking PROCESSED) -> Broker : publish evt-X
Relay (retry) -> Broker : publish evt-X again (duplicate)
Consumer.Handle(evt-X) #1 : processedIDs miss -> store, return true
Consumer.Handle(evt-X) #2 : processedIDs hit -> return false, ignored
Received count stays 1
```

## Diagram 6 — Outbox Lifecycle States

```text
[ PENDING ] --Publish OK + MarkOutboxProcessed--> [ PROCESSED ] --PurgeProcessedOutbox--> [ DELETED ]
[ PENDING ] --Publish FAIL--> [ PENDING ] (retry next poll)
```
