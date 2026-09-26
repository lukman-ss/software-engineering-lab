# Engineering Audit Plan

Target Lab:
  /Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/19-database-connection-pooling

Implementation Files:
  - go.mod                                  (module + go directive)
  - internal/pool/mockdb.go                 (MockDriver / mockConn / MockConnector / OpenDB)
  - internal/pool/service.go                (OrderService: ProcessOrderSafe, ProcessOrderUnsafeLeak)
  - cmd/demo/main.go                        (interactive CLI demo, 3 demos)

Tests:
  - tests/pool_test.go                      (4 tests: overhead, exhaustion, leak/starvation, concurrent safe)

Executable/Demo:
  - cmd/demo/main.go  (run via `go run ./cmd/demo`)

Approved Research Inputs:
  (Out of scope for this engineering stage per pipeline override — content/research not audited here.)

Main Claims To Verify:
  1. Code compiles (`go build ./...`).
  2. Test suite passes (`go test -v ./...`).
  3. Race detector is clean (`go test -race -v ./...`).
  4. Demo runs and produces the documented output (`go run ./cmd/demo`).
  5. README accurately lists components and run commands.
  6. Implementation matches the "safe vs unsafe connection use" narrative:
     - Safe path: external I/O performed outside of the DB connection lifetime.
     - Unsafe path: connection held during external I/O (intentional leak demonstration).
  7. Pool exhaustion / server overload is observable under concurrency.
  8. Connection starvation under an oversized leak is observable with a timeout.

Commands To Run:
  - go version
  - go build ./...
  - go vet ./...
  - go test -v ./...
  - go test -race -v ./...
  - go run ./cmd/demo

Primary Risks:
  - Timing-based test assertions (flakiness on slow/loaded machines).
  - Check-then-act (TOCTOU) on MockDriver max-connection cap enforcement.
  - Declared-but-unused error sentinel (ErrAcquireTimeout).
