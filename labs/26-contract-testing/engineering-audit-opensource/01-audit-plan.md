# Engineering Audit Plan

Target Lab: labs/26-contract-testing
Implementation Files:
- internal/contract/verifier.go (CDC engine, Verify, diffValues)
- internal/consumer/client.go (MobileOrderClient.FetchOrder, GenerateMobileContract)
- internal/provider/server.go (ProviderV1, ProviderBreaking, ProviderDual)
- internal/model/order.go (Order, OrderResponseV1/Breaking/V2 DTOs)
- cmd/demo/main.go (4-stage demo)
- tests/contract_test.go (7 tests)
Tests: tests/contract_test.go (TestConsumerContractGeneration, TestProviderV1_ContractVerification_Success, TestProviderBreaking_ContractVerification_Fails, TestProviderDual_ContractVerification_Success, TestConcurrentContractVerification, TestVerifier_HeaderValidation_And_ErrorBranches, TestProviderDual_V2Endpoint_DirectAssertion)
Executable/Demo: cmd/demo (httptest-backed, no external deps)
Approved Research Inputs: NOT AUDITED (pipeline override — implementation and tests only)
Main Claims To Verify:
1. Consumer generates minimal contract (MobileApp -> OrderService, GET /v1/orders/ORD-123, subset body)
2. Verifier checks status + headers + body-subset against live provider
3. Breaking provider (status casing, customer.name->full_name, total int->string) FAILS verification with >=3 errors and breaks FetchOrder
4. Dual provider keeps V1 passing while exposing V2
5. Concurrent verification is race-safe
6. Demo output is real (PASSED / BLOCKED / PASSED across 3 providers)
Commands To Run:
- go test -v ./... (from labs/26-contract-testing)
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Subset matcher (diffValues) too lax/strict (arrays unsupported, number compare by string)
- Breaking-change detection asserted only by error count, not per-field
- V2 endpoint asserted by status only, schema unchecked
- Shared http.Client under concurrency
