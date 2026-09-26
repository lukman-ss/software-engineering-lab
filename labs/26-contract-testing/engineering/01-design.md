# Engineering Design

Target Lab: `labs/26-contract-testing`
Research Status: APPROVED

## Concept To Prove

1. Consumer-Driven Contracts (CDC) capture consumer obligations and expectations without testing internal provider implementation details.
2. Breaking changes (enum casing change `IN_PROGRESS` -> `in_progress`, field rename `customer.name` -> `customer.full_name`, primitive type mutation `total` int -> string) break consumer contract validation.
3. CI/CD Gate ("Can I Deploy"): contract verification fails when breaking changes are introduced on the provider, blocking deployment.
4. Safe API Evolution: introducing a V2 DTO / endpoint while preserving V1 compatibility allows consumers to migrate independently without breaking existing contracts.

## Expected Behavior

- **Consumer contract generation**: Mobile App consumer defines its minimal contract (needed fields: `id`, `status`, `customer.name`, `total`). Contract is serialized to JSON.
- **Provider contract verification (V1 compliant)**: Provider V1 endpoint satisfies mobile consumer contract. Verification passes, exit code 0.
- **Provider contract verification (Breaking provider)**: Breaking provider mutation fails verification with exact field diffs (missing `customer.name`, type mismatch `total`, enum mismatch `status`). Verification fails, exit code non-zero, deployment blocked.
- **Safe evolution (Dual DTO / V2)**: Provider introduces V2 DTO while maintaining V1 contract handler. V1 contract passes against `/v1/orders/{id}`, while V2 contract passes against `/v2/orders/{id}`.

## Failure Scenario

1. Mobile consumer relies on:
   - `status`: enum value `IN_PROGRESS`
   - `customer.name`: string
   - `total`: integer
2. Provider introduces three unannounced breaking changes:
   - `status`: changed to lowercase `in_progress`
   - `customer`: renamed `name` to `full_name`
   - `total`: changed type from integer `150000` to formatted string `"150000"`
3. Provider CI executes contract verification against the consumer contract file.
4. Contract runner detects:
   - Enum mismatch on `status`
   - Missing field `customer.name`
   - Data type mismatch on `total` (expected number/int, got string)
5. CI gate exits with failure, blocking deployment.

## Success Criteria

1. Standalone Go standard library implementation without heavyweight external C-bindings or external daemon dependencies (pure Go CDC runner).
2. Pure contract generation from consumer test doubles.
3. Provider verification engine verifying HTTP interactions against generated contracts.
4. Clear reporting of contract verification failures blocking invalid provider versions.
5. Verification passes for valid providers and dual-versioned providers.
6. Full test suite passing with race detector (`go test -race ./...`).
7. Demo CLI executing all 3 lifecycle stages:
   - Stage 1: Consumer writes contract & Provider V1 passes verification.
   - Stage 2: Breaking Provider triggers verification failure and deployment block.
   - Stage 3: Dual DTO Provider (V1 + V2) restores contract compatibility and passes gate.

## Architecture

```text
[Consumer (Mobile App Client)]
            │
            ▼ (generates)
 [contracts/mobile_order_v1.json]
            │
            ▼ (verifies against)
 [Contract Verification Engine]
            │
            ├─► [Provider V1 Service]        ===> PASS (Deployable)
            ├─► [Provider Breaking Service]  ===> FAIL (Blocked)
            └─► [Provider V2 Dual Service]   ===> PASS (Deployable)
```

## Components

1. `internal/model`: Order domain entities, DTOs (V1, Breaking, V2).
2. `internal/consumer`: Mobile order client and contract generator creating CDC specification.
3. `internal/provider`: HTTP handlers for Provider V1, Provider Breaking, and Provider Dual (V1 + V2).
4. `internal/contract`: Contract definition structures, JSON serializer, and verification runner comparing interactions (method, path, headers, status, response schema/types/enums).
5. `tests`: Unit and contract integration tests proving verification success, failure on breaking changes, and backward-compatible dual routing.
6. `cmd/demo`: Executable demonstration illustrating the full contract testing lifecycle and CI gate blocking.

## Test Strategy

- `TestConsumerContractGeneration`: Asserts consumer client generates expected interaction contract.
- `TestProviderV1_ContractVerification_Success`: Asserts Provider V1 satisfies Mobile consumer contract.
- `TestProviderBreaking_ContractVerification_Fails`: Asserts Breaking Provider fails contract verification with precise diffs.
- `TestProviderDual_ContractVerification_Success`: Asserts Dual Provider maintains V1 contract compatibility while exposing V2.
- `TestConcurrentVerification`: Verifies provider runner handles concurrent contract checks safely.

## Execution Plan

1. Setup `go.mod` for `labs/26-contract-testing`.
2. Implement `internal/contract` verification library.
3. Implement `internal/model` and `internal/provider` services.
4. Implement `internal/consumer` client and contract builder.
5. Write unit and integration tests under `tests/`.
6. Implement `cmd/demo/main.go`.
7. Execute `go test ./...`, `go test -race ./...`, `go run ./cmd/demo`.
8. Document implementation notes and execution results.

## Implementation Decisions

- *ponytail: Standard library HTTP test servers and custom minimal CDC runner over heavy external pact-go CLI binaries to eliminate external runtime dependencies while demonstrating exact CDC semantics.*
- Minimal subset rule: Consumer specifies only `id`, `status`, `customer.name`, and `total`. Extra provider fields (e.g. `created_at`, `notes`) are ignored during consumer contract verification.
