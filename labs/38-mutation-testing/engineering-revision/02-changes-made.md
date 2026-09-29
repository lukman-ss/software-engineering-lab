## Revision 1

Audit Issue: Unused `mu sync.Mutex` in `Runner` struct
Severity: LOW
Files Changed: `internal/engine/runner.go`
Action: Removed unused `mu sync.Mutex` field from `Runner` struct declaration.
Verification: `go test ./...` and `go test -race ./...` pass.
Status: RESOLVED

## Revision 2

Audit Issue: Missing degenerate input and below-boundary test cases
Severity: LOW
Files Changed: `internal/service/discount_strong_test.go`
Action: Added 4 test cases covering zero total amount and zero item count, below-boundary 499.99 for Premium tier, below-boundary 99.99 for Standard tier, and zero item count with coupon enabled.
Verification: `go test -v ./...` passes all 11 sub-cases of `TestCalculateDiscount_Strong`.
Status: RESOLVED
