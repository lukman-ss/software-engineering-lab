# Engineering Audit Plan

Target Lab: labs/26-contract-testing
Implementation Files:
  - internal/model/order.go
  - internal/contract/verifier.go
  - internal/consumer/client.go
  - internal/provider/server.go
Tests:
  - tests/contract_test.go
Executable/Demo:
  - cmd/demo/main.go
Approved Research Inputs:
  - research/01-plan.md
  - research/02-sources.md
  - research/03-evidence.md
  - research/04-contradictions.md
  - research/05-report.md
  - research/06-open-questions.md
Main Claims To Verify:
  1. Consumer generates contract with minimal required schema expectations.
  2. Verifier checks provider endpoints against contract without external daemon/C-bindings.
  3. Breaking changes (enum casing, field rename, primitive type mutation) fail contract verification with detailed error reporting.
  4. Safe API evolution (dual V1/V2 handlers) preserves V1 contract compliance.
  5. Implementation passes tests under race detector (`go test -race ./...`).
  6. Executable demo runs clean lifecycle matching README.
Commands To Run:
  - go test -count=1 ./...
  - go test -race -count=1 ./...
  - go run ./cmd/demo
Primary Risks:
  - Incomplete diff logic in verifier missing subtle type or structure mismatches.
  - Race conditions during concurrent verifier calls.
  - Discrepancy between README claims and actual code behavior.
