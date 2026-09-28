# Execution Result

## Build
Command:
```bash
go build ./...
```
Result:
```text
Exit Code: 0 (Success)
```

## Tests
Command:
```bash
go test -v ./...
```
Result:
```text
=== RUN   TestNaiveReplicationLag_StaleRead
--- PASS: TestNaiveReplicationLag_StaleRead (0.60s)
=== RUN   TestStickySessionRouting
--- PASS: TestStickySessionRouting (0.55s)
=== RUN   TestReadWithToken_LSN
--- PASS: TestReadWithToken_LSN (0.20s)
=== RUN   TestReplicaLagThreshold_Fallback
--- PASS: TestReplicaLagThreshold_Fallback (0.00s)
=== RUN   TestSynchronousReplication_Freshness
--- PASS: TestSynchronousReplication_Freshness (0.10s)
=== RUN   TestConcurrentAccess_RaceFree
--- PASS: TestConcurrentAccess_RaceFree (0.67s)
PASS
ok  	labs/33-read-replicas-and-replication-lag/tests	2.131s
```

## Race Detector
Command:
```bash
go test -race ./...
```
Result:
```text
ok  	labs/33-read-replicas-and-replication-lag/tests	3.983s
```

## Demo
Command:
```bash
go run ./cmd/demo
```
Result:
```text
=== Lab 33: Read Replicas and Replication Lag Demo ===

--- 1. Demonstrating Async Replication Lag Anomaly ---
[Write Primary] Wrote key 'user:profile:123', Primary LSN: 1
[Naive Read] Node: replica-2, Applied LSN: 0, Val: '', Err: key not found

--- 2. Demonstrating Read-Your-Own-Writes with Sticky Session ---
[Sticky Read (within 500ms)] Node: primary, LSN: 1, Val: '{"name":"Alice","tier":"premium"}', Err: <nil>
Sleeping 600ms for sticky duration to expire and replica to catch up...
[Sticky Read (after TTL)] Node: replica-1, LSN: 1, Val: '{"name":"Alice","tier":"premium"}', Err: <nil>

--- 3. Demonstrating Causal Token / LSN Wait ---
[Write Primary] Wrote key 'order:101', LSN: 2
[Token Read (MinLSN=2)] Node: replica-1, LSN: 2, Val: '{"status":"completed"}', Err: <nil>

--- 4. Demonstrating Synchronous Replication (remote_apply) ---
[Sync Write] LSN: 1, Write Duration: 104.055166ms
[Sync Naive Read] Node: replica-2, LSN: 1, Val: '{"name":"Bob"}', Err: <nil>

=== Demo Complete ===
```

## Final Engineering Status
READY_FOR_ENGINEERING_AUDIT
