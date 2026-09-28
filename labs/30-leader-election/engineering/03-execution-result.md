# Execution Result

Target Lab: `labs/30-leader-election`  
Date: September 28, 2026  

## Build
Command:
```bash
go build ./...
```
Result:
```text
SUCCESS (0 errors, 0 warnings)
```

## Tests
Command:
```bash
go test -v ./...
```
Result:
```text
=== RUN   TestCoordinatorLeaseAcquisitionAndRenewal
--- PASS: TestCoordinatorLeaseAcquisitionAndRenewal (0.25s)
=== RUN   TestFencedStorageRejectsStaleTokens
--- PASS: TestFencedStorageRejectsStaleTokens (0.00s)
=== RUN   TestLeaderElectionFailoverAndSplitBrainDefense
--- PASS: TestLeaderElectionFailoverAndSplitBrainDefense (0.28s)
=== RUN   TestConcurrentElectionRace
--- PASS: TestConcurrentElectionRace (0.10s)
PASS
ok  	labs/30-leader-election/tests	1.170s
```

## Race Detector
Command:
```bash
go test -race ./...
```
Result:
```text
PASS
ok  	labs/30-leader-election/tests	2.008s
```

## Demo
Command:
```bash
go run ./cmd/demo
```
Result:
```text
=== Distributed Leader Election with Fencing Tokens Demo ===

[1] Starting Node-A and Node-B election campaigns...
Status: Node-A=LEADER, Node-B=FOLLOWER
[2] Leader elected: Node-A (Fencing Token: 1)
[3] Node-A performing legitimate fenced write...
Write SUCCESS by Node-A (Token 1)

[4] Simulating Stop-the-World GC Pause / Network Partition on Node-A (350ms > TTL 200ms)...
[5] Status after TTL expiration: Node-A=LEADER, Node-B=LEADER
[6] New Leader promoted: Node-B (Fencing Token: 2)
[7] New Leader Node-B performing fenced write...
Write SUCCESS by Node-B (Token 2)

[8] Old Leader Node-A wakes up from GC pause (still holding stale token 1)...
[9] Old Leader Node-A attempts write to shared storage using stale token 1...
Write REJECTED by FencedStorage: stale fencing token: rejected to prevent split-brain write: token 1 <= last seen 2 (author: Node-A)
Result: Split-brain prevented! Old leader's stale write was safely rejected.

[10] Shared Storage History:
  1. Token: 1 | Author: Node-A | Value: record-from-initial-leader
  2. Token: 2 | Author: Node-B | Value: record-from-failover-leader

=== Demo completed successfully ===
```

## Final Engineering Status
ENGINEERING_STATUS: READY_FOR_ENGINEERING_AUDIT
