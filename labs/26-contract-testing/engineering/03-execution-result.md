# Execution Result

Target Lab: `labs/26-contract-testing`

## Build

Command:
```bash
go build ./...
```

Result:
```text
SUCCESS
```

## Tests

Command:
```bash
go test -v ./...
```

Result:
```text
?   	labs/26-contract-testing/cmd/demo	[no test files]
?   	labs/26-contract-testing/internal/consumer	[no test files]
?   	labs/26-contract-testing/internal/contract	[no test files]
?   	labs/26-contract-testing/internal/model	[no test files]
?   	labs/26-contract-testing/internal/provider	[no test files]
=== RUN   TestConsumerContractGeneration
--- PASS: TestConsumerContractGeneration (0.00s)
=== RUN   TestProviderV1_ContractVerification_Success
--- PASS: TestProviderV1_ContractVerification_Success (0.00s)
=== RUN   TestProviderBreaking_ContractVerification_Fails
--- PASS: TestProviderBreaking_ContractVerification_Fails (0.00s)
=== RUN   TestProviderDual_ContractVerification_Success
--- PASS: TestProviderDual_ContractVerification_Success (0.00s)
=== RUN   TestConcurrentContractVerification
--- PASS: TestConcurrentContractVerification (0.01s)
PASS
ok  	labs/26-contract-testing/tests	0.390s
```

## Race Detector

Command:
```bash
go test -race ./...
```

Result:
```text
ok  	labs/26-contract-testing/tests	1.425s
```

## Demo

Command:
```bash
go run ./cmd/demo
```

Result:
```text
=== Contract Testing Lab: Consumer-Driven Contracts & CI Verification ===

[Stage 1] Consumer generates contract:
Generated Contract (MobileApp -> OrderService):
{
  "consumer": "MobileApp",
  "provider": "OrderService",
  "interactions": [
    {
      "description": "A request for order details by ID",
      "provider_state": "Order ORD-123 exists and is IN_PROGRESS",
      "request": {
        "method": "GET",
        "path": "/v1/orders/ORD-123"
      },
      "response": {
        "status": 200,
        "headers": {
          "Content-Type": "application/json"
        },
        "body": {
          "customer": {
            "name": "Budi Santoso"
          },
          "id": "ORD-123",
          "status": "IN_PROGRESS",
          "total": 150000
        }
      }
    }
  ]
}

[Stage 2] Running Provider V1 Contract Verification:
Result: PASSED. Provider V1 satisfies Mobile consumer contract.
CI Deployment Gate: ALLOWED.

[Stage 3] Running Breaking Provider Contract Verification:
Result: BLOCKED! Breaking changes detected before deployment:
  1. [A request for order details by ID] path 'total': type mismatch (expected 150000 [json.Number], got 150000 [string])
  2. [A request for order details by ID] path 'status': value mismatch (expected "IN_PROGRESS", got "in_progress")
  3. [A request for order details by ID] missing expected field 'customer.name'
CI Deployment Gate: PREVENTED PRODUCTION OUTAGE.

[Stage 4] Running Dual Provider (V1 + V2) Verification:
Result: PASSED. Dual provider maintains backwards-compatible V1 contract.
CI Deployment Gate: ALLOWED for independent canary/migration.

=== Contract Testing Demonstration Complete ===
```

## Final Engineering Status

ENGINEERING_STATUS: READY_FOR_ENGINEERING_AUDIT
