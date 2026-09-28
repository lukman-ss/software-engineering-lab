# Engineering Audit Plan

Target Lab: labs/24-slo-sli-error-budget
Implementation Files:
- internal/alerting/engine.go
- internal/metrics/tracker.go
- internal/slo/evaluator.go
- cmd/demo/main.go
Tests:
- tests/slo_test.go
Executable/Demo:
- cmd/demo/main.go (go run ./cmd/demo)
Approved Research Inputs:
- research/01-plan.md
- research/02-sources.md
- research/03-evidence.md
- research-audit/ (source/claim audits)
- engineering/01-design.md, 02-implementation-notes.md, 03-execution-result.md
Main Claims To Verify:
1. SLO evaluation computes error budget consumption correctly.
2. SLI/metrics tracking records errors and success counts accurately.
3. Alerting engine fires/removes alerts at correct SLO thresholds/burn rates.
4. Demo runs without error and reflects README behavior.
Commands To Run:
- cd labs/24-slo-sli-error-budget && go build ./...
- go test ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Race conditions in concurrent metric ingestion/alert evaluation.
- Mismatched error budget math vs research claims.
- Demo output fabricated or diverging from README.
- Insufficient failure/edge-case coverage in tests.
