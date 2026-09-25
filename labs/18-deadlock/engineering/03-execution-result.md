# Execution Result

## Build
Command:
```bash
go build ./cmd/demo
```
Result:
```text
Exit Code: 0
```

## Tests
Command:
```bash
go test ./...
```
Result:
```text
?   	labs/18-deadlock/cmd/demo	[no test files]
?   	labs/18-deadlock/internal/bank	[no test files]
?   	labs/18-deadlock/internal/transfer	[no test files]
ok  	labs/18-deadlock/tests	0.717s
```

## Race Detector
Command:
```bash
go test -race ./...
```
Result:
```text
?   	labs/18-deadlock/cmd/demo	[no test files]
?   	labs/18-deadlock/internal/bank	[no test files]
?   	labs/18-deadlock/internal/transfer	[no test files]
ok  	labs/18-deadlock/tests	1.715s
```

## Demo
Command:
```bash
go run ./cmd/demo
```
Result:
```text
=== Deadlock Simulation Demo ===

1. Simulating Naive Concurrent Transfers (Circular Wait)...
Transfer ACC-2 -> ACC-1: deadlock victim
Transfer ACC-1 -> ACC-2: <nil>

2. Simulating Lock Ordered Transfers (Deadlock Prevention)...
TransferOrdered ACC-2 -> ACC-1: <nil>
TransferOrdered ACC-1 -> ACC-2: <nil>
Balances: ACC-1=1100, ACC-2=900

3. Simulating Retry Mechanism (Deadlock Recovery)...
TransferWithRetry ACC-1 -> ACC-2: <nil>
TransferWithRetry ACC-2 -> ACC-1: <nil>
Final Balances: ACC-1=1100, ACC-2=900
```

## Final Engineering Status
READY_FOR_ENGINEERING_AUDIT
