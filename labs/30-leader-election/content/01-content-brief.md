# Content Brief

Topic: Leader Election with Fencing Tokens
Target Reader: Backend engineers and distributed systems practitioners implementing coordination-critical services in Go.
Problem: How to prevent split-brain and duplicate execution when a leader pauses longer than its lease TTL, allowing a second node to assume leadership while the old leader still believes it holds authority.
Core Mental Model: A lease grants time-bounded authority; a fencing token converts that temporal authority into a monotonic, verifiable ordering that the shared resource enforces — so a paused leader's stale write is rejected regardless of what it believes about its own status.
Approved Research Status: APPROVED
Approved Engineering Status: APPROVED
Main Concepts: Lease-based leader election, fencing tokens (monotonically increasing revision), TTL expiration and failover, check-and-set semantics, split-brain defense, Redlock debate context, monotonic clock usage.
Verified Behaviors: Single-leader invariant under concurrent campaigns; automatic failover on lease expiry; stale write rejection when token <= lastSeenToken; race-free concurrent election via Go `sync.Mutex`; demo confirms full lifecycle from election to stale-write rejection.
Available Case Studies: etcd lease-based election with revision fencing, ZooKeeper ephemeral sequential nodes with zxid fencing, Redis Redlock limitations (no native fencing), Raft term-based implicit fencing.
Warnings: Coordinator is in-memory only — no external etcd/Redis; timing values (TTL 200ms, renew 50ms) are lab-illustrative, not production recommendations; dual-leader state can transiently occur during pause windows but is mitigated by fencing; this lab demonstrates the pattern in isolation, not a production-grade deployment.
