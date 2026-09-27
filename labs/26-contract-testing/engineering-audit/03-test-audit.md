# Engineering Test Audit

Target Lab: labs/26-contract-testing

## Test Suite Overview

Test file: `tests/contract_test.go`

Coverage includes:
1. `TestContractVerification_V1Provider_Pass`: Verifies baseline V1 provider passes contract check.
2. `TestContractVerification_BreakingProvider_Fails`: Verifies casing change, missing field, and type mutation cause contract failure.
3. `TestContractVerification_DualProvider_Pass`: Verifies backwards-compatible dual provider passes V1 contract.
4. `TestConsumerClient_Integration`: Tests full consumer HTTP fetch against running server.
5. `TestContractVerifier_ConcurrentExecutions`: Tests concurrent verification requests to ensure race safety.

## Required Execution Results

### Unit / Integration Tests (`go test -v ./...`)
```text
=== RUN   TestContractVerification_V1Provider_Pass
--- PASS: TestContractVerification_V1Provider_Pass (0.00s)
=== RUN   TestContractVerification_BreakingProvider_Fails
--- PASS: TestContractVerification_BreakingProvider_Fails (0.00s)
=== RUN   TestContractVerification_DualProvider_Pass
--- PASS: TestContractVerification_DualProvider_Pass (0.00s)
=== RUN   TestConsumerClient_Integration
--- PASS: TestConsumerClient_Integration (0.00s)
=== RUN   TestContractVerifier_ConcurrentExecutions
--- PASS: TestContractVerifier_ConcurrentExecutions (0.00s)
PASS
ok  	github.com/software-engineering-lab/labs/26-contract-testing/tests	0.279s
```

### Race Detector (`go test -race ./...`)
```text
PASS
ok  	github.com/software-engineering-lab/labs/26-contract-testing/tests	0.384s
```

### Demo Execution (`go run ./cmd/demo`)
```text
=== Contract Testing Lab: Consumer-Driven Contracts & CI Verification ===

[Stage 1] Consumer generates contract:
Generated Contract (MobileApp -> OrderService):
...
[Stage 2] Running Provider V1 Contract Verification:
Result: PASSED. Provider V1 satisfies Mobile consumer contract.
CI Deployment Gate: ALLOWED.

[Stage 3] Running Breaking Provider Contract Verification:
Result: BLOCKED! Breaking changes detected before deployment:
...
[Stage 4] Running Dual Provider (V1 + V2) Verification:
Result: PASSED. Dual provider maintains backwards-compatible V1 contract.
CI Deployment Gate: ALLOWED for independent canary/migration.

=== Contract Testing Demonstration Complete ===
```

## Assessment

All test scenarios pass cleanly under race detection. Demo output matches claimed behavior exactly.
