# Documentation vs Code

## README.md vs implementation
README claims:
  - internal/pool/mockdb.go implements a driver simulating max_connections + connect latency.
  - internal/pool/service.go provides safe vs unsafe (leak) operations.
  - tests/pool_test.go covers direct overhead, pool exhaustion, leaks, concurrency.
  - cmd/demo/main.go is an interactive CLI demo.
Code reality: all four components exist exactly as described. PASS.

README commands:
  - `go test -v ./...`  → verified, 4 PASS.
  - `go test -race -v ./...` → verified, 4 PASS, 0 races.
  - `go run ./cmd/demo` → verified, exit 0, output shown below.
All documented commands reproduce. PASS.

## Engineering notes (engineering/03-execution-result.md) vs actual execution

| Step            | Claimed         | Actual (re-executed)      | Match? |
|-----------------|-----------------|---------------------------|--------|
| go build ./...  | Success, no output | Success, no output      | PASS |
| go test -v ./.. | 4 PASS          | 4 PASS                    | PASS |
| go test -race ..| 4 PASS, 0 races | 4 PASS, 0 races           | PASS |
| Demo output     | 3 sections      | 3 sections                | PASS |

Demo output — claimed (engineering notes):
  Unpooled (5 requests): 54.593875ms
  Pooled (5 requests): 6.333µs
  Client attempted: 30, Succeeded: 15, Server Rejected: 15
  Order 3 Failed: context deadline exceeded

Demo output — actual (this audit run):
  Unpooled (5 requests): 54.242792ms
  Pooled (5 requests): 2.833µs
  Client attempted: 30, Succeeded: 15, Server Rejected: 15
  Order 3 Failed: context deadline exceeded

The integer-valued behavioral assertions (15/15 split, context deadline exceeded, all demo
sections present) reproduce exactly. The sub-millisecond timing values differ because they
are wall-clock dependent (54.24ms vs 54.59ms; 2.8µs vs 6.3µs). This is expected variance
on a live system, NOT fabrication: the engineering notes recorded one real run and the
structure/counts are reproducible verbatim. No DOC_CODE_MISMATCH or RESEARCH_IMPLEMENTATION_MISMATCH.

## Findings

## Finding 1

Location: README.md vs code
Claimed Behavior: Demo prints "Client attempted: 30, Succeeded: 15, Server Rejected: 15".
Observed: matches exactly.
Assessment: PASS
Severity: N/A
Notes: 15 successes / 15 rejections because server max_connections=15 and the demo fires 30
  requests with 10ms holds. Matches README narrative.

## Finding 2

Location: engineering/03-execution-result.md final status line
Claimed Behavior: Final Engineering Status = READY_FOR_ENGINEERING_AUDIT
Observed Implementation: Status field is prose, not enforced by code/tests — but all
  preceding commands (build/test/race/demo) were verified to pass independently in this
  audit, so the readiness claim is substantiated.
Assessment: PASS
Severity: N/A
Type: None (no mismatch; status is corroborated by re-execution).
