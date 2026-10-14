---
title: "Contract Testing: Consumer-Driven Contracts and CI Gates in Go"
published: false
description: "How to prevent breaking API changes between microservices using Consumer-Driven Contracts and automated CI verification gates in Go."
tags: "golang, microservices, testing, architecture"
---

Integration failures routinely occur in production despite all individual unit and integration tests passing green. A familiar scenario: Consumer Service reads a `name` field from a JSON payload. Provider Service renames that field to `full_name`. The consumer test suite passes because its mock fixtures remain stale. The provider test suite passes because internal business logic succeeds. When deployed together, the consumer receives `null` or empty strings, unmarshaling fails, and transactions crash.

This breakdown stems from an unverified contract: an implicit assumption about message formats across process boundaries that was never executed against the live provider implementation prior to deployment.

![Cover Image](https://raw.githubusercontent.com/lukman-ss/software-engineering-lab/main/labs/26-contract-testing/result/image-cover/lab-26-contract-testing-cover-en.png)

## The Core Problem: Why E2E and Unit Tests Fall Short

Traditional testing pyramids struggle at service boundaries:

1. **Unit tests with mocks** verify internal business logic against synthetic expectations. When downstream APIs change, stale mocks continue to return obsolete schemas, masking integration drift.
2. **End-to-End (E2E) integration environments** are brittle, slow, expensive to provision, and difficult to isolate in concurrent CI pipelines. Flaky network hops and cross-service test data contamination cause false negatives.
3. **Pure JSON Schemas** validate structural syntax, but fail to enforce HTTP semantics, specific status codes, strict enum casing (`IN_PROGRESS` vs `in_progress`), or consumer-specific subset needs.

Consumer-Driven Contract (CDC) testing solves this boundary problem. The consumer declares an executable specification containing strictly the fields, headers, and HTTP status codes it requires. The provider executes these contracts against its actual HTTP handlers during CI.

```text
[Consumer (MobileApp Client)]
             │
             │ (Generates minimal contract specification)
             ▼
[contract.json — Minimal Consumer Expectations]
             │
             │ (Verifier runner executes contract interactions)
             ▼
[Contract Verifier Engine — Verify(baseURL, contract)]
             │
      ┌───────┼────────────────────────┐
      ▼       ▼                        ▼
[Provider V1] [Breaking Provider]      [Dual Provider (/v1 & /v2)]
      │               │                        │
      ▼               ▼                        ▼
    PASS        FAIL (CI Blocked)            PASS (V1 route preserved)
```

## Anatomy of a Consumer-Driven Contract

In Consumer-Driven Contracts, the consumer authors executable interactions. The contract captures:
- HTTP Method and Target Endpoint Path
- Expected Request and Response Headers (e.g., `Content-Type: application/json`)
- Expected HTTP Status Code
- Required Response Body Subset (ignoring extra unneeded provider fields)

Here is how the consumer contract is generated in Go using standard library primitives:

```go
// internal/consumer/client.go:17-48
func GenerateMobileContract() *contract.Contract {
	return &contract.Contract{
		Consumer: "MobileApp",
		Provider: "OrderService",
		Interactions: []contract.Interaction{
			{
				Description:   "A request for order details by ID",
				ProviderState: "Order ORD-123 exists and is IN_PROGRESS",
				Request: contract.RequestDefinition{
					Method: http.MethodGet,
					Path:   "/v1/orders/ORD-123",
				},
				Response: contract.ResponseDefinition{
					Status: http.StatusOK,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
					Body: map[string]interface{}{
						"id":     "ORD-123",
						"status": "IN_PROGRESS",
						"customer": map[string]interface{}{
							"name": "Budi Santoso",
						},
						"total": json.Number("150000"),
					},
				},
			},
		},
	}
}
```

Notice the use of `json.Number("150000")`. In dynamic JSON deserialization, raw floating-point or generic number decoders can mask type mutations. Explicit number handling guarantees that string-to-numeric or integer-to-float mutations trigger verification failures.

## Three Common Breaking Changes

Consider an Order Service provider that introduces three simultaneous unannounced mutations:

```go
// internal/model/order.go:21-27
type OrderResponseBreaking struct {
	ID       string                  `json:"id"`
	Status   string                  `json:"status"`    // mutated: "in_progress" (lowercase enum)
	Customer CustomerResponseBreaking `json:"customer"` // mutated: field renamed to full_name
	Total    string                  `json:"total"`     // mutated: int converted to string
}
```

1. **Enum Casing Mutation**: Changing `"IN_PROGRESS"` to `"in_progress"`. Syntactically valid JSON, but mobile app enum parsers crash or default to unknown states.
2. **Field Rename**: Changing `customer.name` to `customer.full_name`. The original key disappears, leaving consumer deserializers with empty string values.
3. **Primitive Type Mutation**: Changing `total` from integer `150000` to string `"150000"`. A consumer performing numeric arithmetic experiences deserialization panics.

When the contract verifier runs against this breaking provider, all three discrepancies are captured deterministically:

```text
[ProviderBreaking] Verification failed:
 - Missing expected field 'customer.name' at 'customer.name'
 - Value mismatch at 'status': expected IN_PROGRESS, got in_progress
 - Type mismatch at 'total': expected json.Number, got string
```

The CI verification runner halts the deployment pipeline, preventing the faulty code from reaching production.

## Verifier Engine: Recursive Subset Matching

The verifier executes each declared interaction against the target provider, validating status codes, response headers, and payload structures recursively.

```go
// internal/contract/verifier.go:81-118
func (v *Verifier) Verify(baseURL string, c *Contract) VerificationResult {
	result := VerificationResult{Passed: true}

	for _, interaction := range c.Interactions {
		targetURL := strings.TrimRight(baseURL, "/") + interaction.Request.Path
		req, err := http.NewRequest(interaction.Request.Method, targetURL, nil)
		if err != nil {
			result.Passed = false
			result.Errors = append(result.Errors, fmt.Sprintf("[%s] Request creation failed: %v", interaction.Description, err))
			continue
		}

		resp, err := v.Client.Do(req)
		if err != nil {
			result.Passed = false
			result.Errors = append(result.Errors, fmt.Sprintf("[%s] HTTP request failed: %v", interaction.Description, err))
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != interaction.Response.Status {
			result.Passed = false
			result.Errors = append(result.Errors, fmt.Sprintf("[%s] Status mismatch: expected %d, got %d",
				interaction.Description, interaction.Response.Status, resp.StatusCode))
		}

		for k, expectedVal := range interaction.Response.Headers {
			actualVal := resp.Header.Get(k)
			if actualVal != expectedVal {
				result.Passed = false
				result.Errors = append(result.Errors, fmt.Sprintf("[%s] Header mismatch for %q: expected %q, got %q",
					interaction.Description, k, expectedVal, actualVal))
			}
		}

		// Body subset validation via recursive map traversal
		diffs := diffValues("", interaction.Response.Body, actualBody)
		if len(diffs) > 0 {
			result.Passed = false
			for _, d := range diffs {
				result.Errors = append(result.Errors, fmt.Sprintf("[%s] %s", interaction.Description, d))
			}
		}
	}

	return result
}
```

The subset comparison logic ensures backward compatibility: provider responses can include extra fields (e.g., `created_at`, `notes`, or telemetry tracking) without failing consumer tests. Only properties explicitly demanded by the consumer contract are enforced.

## Safe API Evolution: The Expand-Contract Pattern

When schema enhancements are required, teams must avoid in-place breaking changes. By adopting the Expand-Contract (Parallel Change) pattern, the provider introduces a new route or backward-compatible schema without disturbing existing consumers.

In our lab implementation, `ProviderDual` serves both `/v1` and `/v2` endpoints:

```go
// internal/provider/server.go:61-98
func (p *ProviderDual) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case strings.HasPrefix(r.URL.Path, "/v1/orders/"):
		// Preserves V1 contract compatibility
		resp := model.OrderResponseV1{
			ID:       "ORD-123",
			Status:   "IN_PROGRESS",
			Customer: model.CustomerResponseV1{Name: "Budi Santoso"},
			Total:    150000,
			Notes:    "Dual provider v1 route",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)

	case strings.HasPrefix(r.URL.Path, "/v2/orders/"):
		// Exposes enhanced V2 schema with Currency
		resp := model.OrderResponseV2{
			ID:       "ORD-123",
			Status:   "IN_PROGRESS",
			Customer: model.CustomerResponseV2{FullName: "Budi Santoso"},
			Total:    150000,
			Currency: "IDR",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	default:
		http.NotFound(w, r)
	}
}
```

When CI executes verification against `ProviderDual`, the `/v1` contract passes completely, allowing production deployment while consumers migrate to `/v2` at their own pace.

## Test Verification

The test suite in `tests/contract_test.go` verifies 7 concrete test scenarios:

```bash
$ go test -v -race ./...
=== RUN   TestConsumerContractGeneration
--- PASS: TestConsumerContractGeneration (0.00s)
=== RUN   TestProviderV1_ContractVerification_Success
--- PASS: TestProviderV1_ContractVerification_Success (0.00s)
=== RUN   TestProviderBreaking_ContractVerification_Fails
--- PASS: TestProviderBreaking_ContractVerification_Fails (0.00s)
=== RUN   TestProviderDual_ContractVerification_Success
--- PASS: TestProviderDual_ContractVerification_Success (0.00s)
=== RUN   TestConcurrentContractVerification
--- PASS: TestConcurrentContractVerification (0.00s)
=== RUN   TestVerifier_HeaderValidation_And_ErrorBranches
--- PASS: TestVerifier_HeaderValidation_And_ErrorBranches (0.00s)
=== RUN   TestProviderDual_V2Endpoint_DirectAssertion
--- PASS: TestProviderDual_V2Endpoint_DirectAssertion (0.00s)
PASS
ok  	github.com/lukman-ss/software-engineering-lab/labs/26-contract-testing/tests	0.184s
```

## Implementation Boundaries & Production Architecture

This laboratory uses a pure Go standard library implementation to demonstrate contract mechanics without external operational dependencies. In enterprise production systems, several components extend this foundation:

1. **Centralized Contract Broker**: Replace local in-memory structs with Pact Broker or central artifact stores to coordinate contract versions and matrix verification (`can-i-deploy` gates).
2. **Dynamic Provider States**: Real-world providers require state fixture hooks (e.g., seeding a database with an order ID before verifying the HTTP handler).
3. **Event-Driven Messaging**: Message Pact applies the same consumer-driven contract principles to asynchronous Kafka, RabbitMQ, and SQS payloads.
4. **Avoid Over-Specification**: Contracts must assert message structures and boundary errors, never complex domain validation rules (e.g., asserting string length rules belongs in provider unit tests).

## Source Code

Complete source code, verification suites, and demo binaries are available in the repository:
[https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/26-contract-testing](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/26-contract-testing)
