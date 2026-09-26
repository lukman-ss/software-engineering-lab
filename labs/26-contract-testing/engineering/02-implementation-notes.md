# Implementation Notes

Target Lab: `labs/26-contract-testing`

## Files Added

- `go.mod`: Go module declaration.
- `internal/contract/verifier.go`: Minimal CDC verification engine comparing expected contract interactions, headers, status codes, and subset JSON schemas against live HTTP services.
- `internal/model/order.go`: Domain models and DTO representations for V1, Breaking, and V2 schemas.
- `internal/provider/server.go`: HTTP servers representing Provider V1, Breaking Provider, and Backward-Compatible Dual Provider.
- `internal/consumer/client.go`: Mobile consumer client and CDC contract generator.
- `tests/contract_test.go`: Comprehensive test suite verifying contract generation, CI gate verification failure on breaking changes, dual provider compatibility, and concurrent execution safety.
- `cmd/demo/main.go`: Interactive CLI demonstrating full lifecycle: generation, V1 verification pass, breaking change block, and dual DTO recovery.
- `engineering/01-design.md`: Technical lab design and architecture.
- `engineering/02-implementation-notes.md`: Engineering choices, trade-offs, and coverage.
- `engineering/03-execution-result.md`: Build, test, race, and demo execution logs.
- `README.md`: Lab documentation explaining runnable behavior and usage.

## Core Design Decisions

1. **Pure Go Implementation**: Built without external heavy daemon dependencies (e.g. Ruby pact binaries or external pact mock service processes) to ensure reproducible, zero-friction local and CI execution using Go's `net/http` and `httptest`.
2. **Minimal Subset Verification**: Consumer contracts assert only fields explicitly required by the consumer (`id`, `status`, `customer.name`, `total`). Additional provider fields (`notes`, `created_at`) are intentionally omitted from validation.
3. **Strict Validation on Declared Fields**: Validates type safety (`json.Number` vs `string`), exact enum case (`IN_PROGRESS` vs `in_progress`), and nested field existence (`customer.name`).

## Implementation-Specific Choices

- *ponytail: Used JSON number decoding and recursive map comparison rather than full Pact specification AST to achieve minimal footprint with exact semantic equivalence.*

## Known Limitations

1. Does not implement a remote HTTP Pact Broker client (uses local/in-memory contract exchange).
2. Provider states in this lab are handled via deterministic endpoint routing rather than state-setup callbacks.

## Trade-offs

- Custom verifier vs full Pact framework: Custom verifier avoids heavy multi-platform C/Ruby bindings and external process dependencies while clearly demonstrating CDC principles and breaking change detection.

## What Is Demonstrated

- Consumer contract generation defining minimal expectations.
- Verification pass on compliant Provider V1.
- Verification failure with 3 exact breaking change diffs on Breaking Provider:
  - Enum casing mismatch (`IN_PROGRESS` vs `in_progress`)
  - Missing field (`customer.name`)
  - Type mismatch (`int64` vs `string`)
- Safe API evolution using dual DTO routing (`/v1` and `/v2`).
- Concurrency safety under race detector.

## What Is Not Demonstrated

- Asynchronous message broker contract testing (e.g., Kafka / RabbitMQ pacts).
- Contract version matrix publishing to Pact Broker cloud service.
