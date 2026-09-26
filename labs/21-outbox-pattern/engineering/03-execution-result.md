# Execution Result

## Build
Command:
```bash
go build ./...
```
Result:
```text
(no output - compiled successfully)
```

## Tests
Command:
```bash
go test ./...
```
Result:
```text
?   	github.com/software-engineering-lab/labs/21-outbox-pattern/cmd/demo	[no test files]
?   	github.com/software-engineering-lab/labs/21-outbox-pattern/internal/outbox	[no test files]
ok  	github.com/software-engineering-lab/labs/21-outbox-pattern/tests	0.574s
```

## Race Detector
Command:
```bash
go test -race ./...
```
Result:
```text
?   	github.com/software-engineering-lab/labs/21-outbox-pattern/cmd/demo	[no test files]
?   	github.com/software-engineering-lab/labs/21-outbox-pattern/internal/outbox	[no test files]
ok  	github.com/software-engineering-lab/labs/21-outbox-pattern/tests	1.526s
```

## Demo
Command:
```bash
go run ./cmd/demo
```
Result:
```text
=== Lab 21: Transactional Outbox Pattern Demo ===

[Scenario 1: The Dual-Write Problem]
Direct write failed: failed to publish to broker after DB commit: broker unavailable
State Inconsistency: Order in DB = true, Broker Message Count = 0

[Scenario 2: Transactional Outbox Solution]
Creating order with transactional outbox...
Order and Outbox record atomically saved to DB.
Broker received messages: 1
 - Event ID: evt-order-outbox-success, Type: OrderCreated, Payload: {"ID":"order-outbox-success","CustomerID":"cust-2","Amount":300,"Status":"CREATED"}
 - Consumer processing initial message: accepted=true

[Scenario 3: At-Least-Once Delivery & Idempotent Consumer]
Simulating duplicate delivery to consumer...
Consumer processing duplicate delivery: accepted=false (Duplicate safely skipped!)
Total events processed by consumer: 1

=== Demo Complete ===
```

## Final Engineering Status
READY_FOR_ENGINEERING_AUDIT
