# Engineering Revision Plan

Target Lab: `labs/20-zero-downtime-deployment`
Previous Verdict: `APPROVED`

## Blocking Issues
None.

## Non-Blocking Issues
None.

## Files To Change
- `internal/server/server.go`
- `tests/server_test.go`

## Tests To Add/Modify
- Modify `tests/server_test.go` (`TestServerMultiRequestDrain`) to track unhandled write errors cleanly without false negative assertion failures.

## Validation Commands
```bash
go test -v ./...
go test -race ./...
go run ./cmd/demo
```
