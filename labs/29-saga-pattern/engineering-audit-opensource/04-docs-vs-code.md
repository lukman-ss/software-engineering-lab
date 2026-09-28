## Docs vs Code Comparison

### Claim: Orchestrator performs LIFO compensation on failure
- Code: orchestrator.compensate iterates executed slice reverse order (orchestrator.go:75-84). PASS.
- Docs: README and design state LIFO rollback. PASS.

### Claim: Idempotent payment processing
- Code: ProcessPayment checks processedID map and returns nil if already processed (services.go:86-89). PASS.
- Docs: README mentions idempotency keys. PASS.

### Claim: Semantic lock prevents concurrent order creation
- Code: CreateOrder checks locks map and sets lock (services.go:33-39). PASS.
- Docs: README and design list semantic lock countermeasure. PASS.

### Claim: Concurrency safety verified with race detector
- Code: Tests run go test -race; all pass. PASS.
- Docs: README shows go test -race command. PASS.

### Claim: Demo reflects successful and failure scenarios
- Code: demo prints scenario results matching test expectations. PASS.
- Docs: README describes demo showing happy path and rollback. PASS.

### Missing: Compensation error handling and context cancellation
- Docs do not mention these aspects, and code does not handle them. WARNING.
