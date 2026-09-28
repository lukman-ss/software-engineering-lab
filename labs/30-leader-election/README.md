# Leader Election with Fencing Tokens

Runnable implementation and verification lab for lease-based leader election, failure detection via lease TTL expiration, and split-brain defense using monotonically increasing fencing tokens.

## Architecture

- `internal/coordinator`: Thread-safe linearizable key-lease coordinator generating monotonically increasing fencing tokens (`revision`).
- `internal/candidate`: Node campaign runner with periodic keep-alive renewals and GC pause / partition simulation.
- `internal/storage`: Fenced storage engine enforcing check-and-set semantics (`incomingToken > lastSeenToken`).

## Component Structure

```text
labs/30-leader-election/
├── cmd/
│   └── demo/
│       └── main.go
├── internal/
│   ├── candidate/
│   │   └── candidate.go
│   ├── coordinator/
│   │   └── coordinator.go
│   └── storage/
│       └── storage.go
├── tests/
│   └── election_test.go
├── engineering/
│   ├── 01-design.md
│   ├── 02-implementation-notes.md
│   └── 03-execution-result.md
├── go.mod
└── README.md
```

## Running the Lab

### Run Tests
```bash
go test -v ./...
```

### Run Race Detector
```bash
go test -race ./...
```

### Run Demo
```bash
go run ./cmd/demo
```
