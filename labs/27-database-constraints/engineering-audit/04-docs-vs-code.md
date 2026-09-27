# Docs vs Code Audit

## Comparison Matrix

| Subject | README Claim | Code & Test Implementation | Status |
|---|---|---|---|
| NOT NULL (`23502`) | Ensures required columns accept no null values | Implemented in `engine.go:50-56`, tested in `store_test.go:16-38`, demonstrated in `main.go:25-27` | MATCH |
| CHECK (`23514`) | Evaluates boolean predicates (`age >= 18`, `status IN (...)`, `total_cents > 0`) | Implemented in `engine.go:58-68,132-135`, tested in `store_test.go:40-67`, demonstrated in `main.go:29-36` | MATCH |
| UNIQUE (`23505`) | Prevents duplicate rows & race conditions | Implemented in `engine.go:78-83`, tested in `store_test.go:69-86,162-208`, demonstrated in `main.go:58-95` | MATCH |
| FOREIGN KEY (`23503`) | Enforces referential integrity | Implemented in `engine.go:137-140`, tested in `store_test.go:88-112`, demonstrated in `main.go:37-41` | MATCH |
| PARTIAL UNIQUE INDEX | Enables soft delete re-registration (`WHERE deleted_at IS NULL`) | Implemented in `engine.go:71-77,102-120`, tested in `store_test.go:114-160`, demonstrated in `main.go:42-57` | MATCH |
| Test Commands | `go test -v ./...`, `go test -race ./...`, `go run ./cmd/demo` | Executed and confirmed passing with zero race conditions | MATCH |

## Audit Discrepancy Findings
- DOC_CODE_MISMATCH: None.
- TEST_CLAIM_MISMATCH: None.
- RESEARCH_IMPLEMENTATION_MISMATCH: None.
