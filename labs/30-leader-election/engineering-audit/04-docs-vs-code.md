# Docs vs Code Audit

## Comparisons

### 1. README.md vs Code
- README lists commands `go test -v ./...`, `go test -race ./...`, `go run ./cmd/demo`. All commands execute cleanly and pass.
- Component structure matches the exact file tree in repository.
- Architecture summary matches package responsibilities.
- Assessment: PASS.

### 2. Engineering Notes vs Code
- Core design choices (monotonic tokens, resource-layer fencing, standard library) match implementation.
- Limitations accurately document in-memory scope and lack of disk WAL / Raft distribution.
- Assessment: PASS.

### 3. Research Claims vs Implementation
- Research claim: Leases provide timed mutual exclusion. Implemented in `coordinator.Coordinator` using TTL and monotonic clock.
- Research claim: Fencing tokens guard against STW GC pauses. Implemented and proven via `storage.FencedStorage` and `candidate.Node.SimulatePause`.
- Assessment: PASS.

### 4. Demo Output vs Execution
- `cmd/demo/main.go` output demonstrates real-time state changes and split-brain rejection matching logged results in `engineering/03-execution-result.md`.
- Assessment: PASS.
