---
title: "Leader Election with Fencing Tokens: Preventing Split-Brain and Stale Writes"
published: false
description: "Explore lease-based leader election, failure detection via TTL, and monotonic fencing tokens for split-brain defense in distributed Go systems."
tags: go, distributed-systems, architecture, consistency, gocong
---

# Leader Election with Fencing Tokens: Preventing Split-Brain and Stale Writes

In distributed systems, unannounced execution pauses like stop-the-world GC or network hiccups can freeze a leader process for hundreds of milliseconds. When the lease TTL expires during this pause, the coordinator may promote a new leader. The old leader, unaware it has been superseded, can then write stale data to shared storage, causing silent corruption. Lab 30 demonstrates a Go implementation of TTL-bounded leases paired with monotonically increasing fencing tokens to defend against this split-brain scenario.

> Fencing tokens provide write safety; leases provide timed coordination.

## The split-brain hazard of stale leadership

When a leader pauses unexpectedly, its lease TTL may expire while the process believes it still holds authority. A new leader can be elected, but the old leader resumes execution and continues writing to shared storage without knowing its leadership has been revoked. This creates a window where two nodes believe they are legitimate writers, and the storage layer receives conflicting mutations.

## TOCTOU vulnerabilities in local state checks

Naive implementations check local state before every write:

```go
func (n *Node) DoWorkNaive(val string) error {
    if n.state != StateLeader { return ErrNotLeader }
    // GC pause occurs here for 350ms!
    return n.sharedStorage.Save(val) // Overwrites new leader's data!
}
```

The gap between state verification and the actual storage operation (the TOCTOU gap) enables split-brain corruption. Network-asynchronous conditions prevent local status checks from providing adequate guarantees.

## Architecture: Leases and fencing tokens

```go
type Lease struct {
    Key          string
    HolderID     string
    FencingToken int64         // strictly monotonic revision
    ExpiresAt    time.Time
    TTL          time.Duration
}
```

Leases grant temporal authority for a bounded duration, while fencing tokens establish a permanent ordering at the storage layer. The coordinator issues monotonically increasing revision numbers on each leadership transition, and storage rejects mutations with tokens less than or equal to the last seen value.

## Pipeline execution and stale write rejection

1. Node-A acquires lease with Token=1, TTL=200ms
2. Node-A writes successfully: incoming 1 > lastSeen 0
3. Node-A pauses for 350ms (exceeds TTL)
4. Lease expires; Node-B acquires lease with Token=2
5. Node-B writes successfully: incoming 2 > lastSeen 1
6. Node-A awakens, attempts write with Token=1
7. Storage rejects: token 1 <= lastSeen 2 → ErrStaleFencingToken

## Fenced storage write gate

```go
func (s *FencedStorage) Write(author string, token int64, value string) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    if token <= s.lastSeenToken {
        return fmt.Errorf("%w: token %d <= last seen %d (author: %s)",
            ErrStaleFencingToken, token, s.lastSeenToken, author)
    }
    s.lastSeenToken = token
    s.records = append(s.records, Record{FencingToken: token, Author: author, Value: value})
    return nil
}
```

The mutex ensures linearizable check-and-set semantics. Only strictly increasing tokens update storage state, turning unordered distributed writes into a monotonic stream.

## Evidence: Terminal execution and test verification

```
[1] Starting Node-A and Node-B election campaigns...
[2] Leader elected: Node-A (Fencing Token: 1)
[3] Write SUCCESS by Node-A (Token 1)
[4] Simulating GC Pause on Node-A (350ms > TTL 200ms)...
[5] Status: Node-A=LEADER, Node-B=LEADER (transient dual-leader)
[6] New Leader promoted: Node-B (Token 2)
[7] Write SUCCESS by Node-B (Token 2)
[8] Old Leader Node-A awakens, attempts write with Token 1
[9] Write REJECTED: stale fencing token: token 1 <= last seen 2
Result: Split-brain prevented!
```

Four unit tests verify lease lifecycle, storage rejection, failover, and concurrent leadership competition. All pass with Go's race detector enabled.

## Design decisions: monotonic clocks and revisions

1. Use monotonic clocks for TTL evaluation to resist NTP-induced clock jumps
2. Increment revision counters atomically on every new lease
3. Enforce `incomingToken > lastSeenToken` at the storage boundary

## Trade-offs: safety versus complexity

Fencing tokens add storage layer complexity—backend resources must support token tracking and atomic check-and-set. However, they provide absolute data safety by rejecting stale writes at the boundary. Failover latency is bounded by the lease TTL, and clients must handle `ErrStaleFencingToken` when writes are rejected.

## Implementation checklist for production

- Calibrate TTL based on empirical p99 GC pause and network RTT
- Set renewal rate well below TTL (renewRate <= TTL/3) with measured backoff
- Use monotonic duration calculations for all lease logic
- Store only tokens that strictly exceed the last recorded watermark
- Monitor stale write rejection rates as early warning of node stalls

## Lab constraints: what is and isn't proven

This lab demonstrates concepts in an in-process test harness. The coordinator is single-node in-memory, not a Raft cluster. Storage is in-memory sync.RWMutex, not PostgreSQL CAS. Timing constants (TTL 200ms, renew 50ms) are illustrative only. Real-world fault injection (disk failures, network partitions, Byzantine failures) is not tested.

## Industry patterns: etcd, ZooKeeper, and fencing

Production systems embed fencing mechanisms: etcd uses 64-bit Raft revisions as tokens; ZooKeeper provides zxid and sequential ephemeral nodes. The Redis Redlock debate established that distributed locks without monotonic tokens serve efficiency, not safety-critical correctness.

## Core principle

Time passes invisibly to a paused process; local belief of leadership is never sufficient proof. The shared resource, with its global perspective, becomes the ultimate arbiter of write eligibility.

## Reproducibility

```bash
cd labs/30-leader-election
go test -v -race ./...
go run ./cmd/demo
```

All builds use Go's standard library. Source: [github.com/lukman-ss/software-engineering-lab/tree/main/labs/30-leader-election](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/30-leader-election), commit `ba4b6216`.