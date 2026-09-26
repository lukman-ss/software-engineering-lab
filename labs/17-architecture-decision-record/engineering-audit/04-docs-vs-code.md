# Docs vs Code Audit

- README vs Code:
  - Lists internal/adr files accurately.
  - Lists `cmd/demo/main.go` accurately.
  - Lists commands accurately (`go test -v ./...`, `go test -race ./...`, `go run ./cmd/demo`).
- Engineering Notes vs Code:
  - Matches design expectations and implementation decisions.
- Demo Output vs Code:
  - Verified real execution output from `cmd/demo/main.go` matches claimed architectural evolution scenario.
