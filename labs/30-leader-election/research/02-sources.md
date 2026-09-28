# Source Compendium: Leader Election and Distributed Locking

Target Lab: `labs/30-leader-election`  
Last Verified: September 28, 2026  

---

## Overview

This source compendium provides verified citations, URLs, publisher details, and extraction notes for all primary and secondary literature referenced in the leader election and distributed locking research documents.

---

## Source Inventory

### 1. Raft Consensus Algorithm Paper
- **Title**: In Search of an Understandable Consensus Algorithm
- **Authors**: Diego Ongaro and John Ousterhout
- **Publisher**: USENIX Annual Technical Conference (ATC '14)
- **URL**: [https://raft.github.io/raft.pdf](https://raft.github.io/raft.pdf)
- **Reachable**: YES
- **Source Type**: Academic Peer-Reviewed Paper / Primary Engineering Source
- **Relevance**: Definitive guide to leader election via randomized timers and term-based log replication in consensus systems (Raft).
- **Extraction Notes**: Establishes how leader election uses randomized election timeouts (150ms–300ms) and heartbeats to ensure single-leader invariants per term, preventing split-brain when network partitions isolate a minority.

### 2. Martin Kleppmann's Redlock Critique
- **Title**: How to do distributed locking
- **Author**: Martin Kleppmann
- **Publisher**: University of Cambridge / Personal Technical Blog
- **URL**: [https://martin.kleppmann.com/2016/02/08/how-to-do-distributed-locking.html](https://martin.kleppmann.com/2016/02/08/how-to-do-distributed-locking.html)
- **Reachable**: YES
- **Source Type**: Reputable Technical Reference / Distributed Systems Expert Analysis
- **Relevance**: Evaluates fault-tolerant distributed locking (leases) on Redis (Redlock) and highlights safety violations caused by clock jumps, GC pauses, and network delays without fencing tokens.
- **Extraction Notes**: Proves that distributed locks used for correctness require monotonic clocks and monotonically increasing *fencing tokens* (generation numbers checked by storage backends) because stop-the-world GC pauses or delayed packets can cause two clients to concurrently believe they hold the lock.

### 3. Salvatore Sanfilippo (antirez) Redlock Rebuttal
- **Title**: Is Redlock safe?
- **Author**: Salvatore Sanfilippo (antirez)
- **Publisher**: Redis / antirez blog
- **URL**: [http://antirez.com/news/101](http://antirez.com/news/101)
- **Reachable**: YES
- **Source Type**: Original Author Technical Rebuttal / Primary Engineering Source
- **Relevance**: Provides counter-analysis to Kleppmann's critique of Redlock, discussing the semi-synchronous system model, drift bounds, and check-time-again guard steps in Redlock.
- **Extraction Notes**: Acknowledges that monotonic time APIs (`clock_gettime`) should be used to avoid wall-clock jumps, but maintains that Redlock's drift bound assumption is practically satisfiable in well-administered datacenters and that fencing tokens assume a linearizable storage layer.

### 4. etcd Concurrency and Leader Election Documentation
- **Title**: How to conduct leader election in etcd cluster / etcd Concurrency API
- **Publisher**: etcd / CNCF
- **URL**: [https://etcd.io/docs/v3.5/tutorials/how-to-conduct-elections/](https://etcd.io/docs/v3.5/tutorials/how-to-conduct-elections/)
- **Reachable**: YES
- **Source Type**: Official Vendor / CNCF Documentation
- **Relevance**: Practical guide to building strongly consistent leader election using etcd v3 leases, revisions, and the `clientv3/concurrency` package.
- **Extraction Notes**: Demonstrates how etcd uses Raft log revision numbers as natural monotonic fencing tokens and lease TTL heartbeats to automatically revoke leadership when a leader crashes or loses network connectivity.

### 5. Apache ZooKeeper Recipes and Leader Election
- **Title**: ZooKeeper Recipes and Solutions: Leader Election
- **Publisher**: The Apache Software Foundation
- **URL**: [https://zookeeper.apache.org/doc/current/recipes.html#sc_leaderElection](https://zookeeper.apache.org/doc/current/recipes.html#sc_leaderElection)
- **Reachable**: YES
- **Source Type**: Official Open Source Project Documentation / Standards Reference
- **Relevance**: Reference recipe for distributed leader election using ephemeral sequential nodes (`SEQUENCE | EPHEMERAL`).
- **Extraction Notes**: Details how creating ephemeral sequential znodes (`/election/guid-n_000000001`) allows the client with the lowest sequence number to claim leadership, while watching only the immediately preceding znode prevents herd effects during failover.

### 6. Designing Data-Intensive Applications (DDIA)
- **Title**: Designing Data-Intensive Applications: The Big Ideas Behind Reliable, Scalable, and Maintainable Systems
- **Author**: Martin Kleppmann
- **Publisher**: O'Reilly Media
- **URL**: [https://dataintensive.net/](https://dataintensive.net/)
- **Reachable**: YES
- **Source Type**: Definitive Technical Textbook
- **Relevance**: Comprehensive textbook covering distributed data systems, consensus, replication, transactions, and distributed locking edge cases (Chapters 8 and 9).
- **Extraction Notes**: Explains partial synchrony, asynchronous networks, clock skew, GC pause hazards, and why coordination services require consensus and fencing tokens for correctness.

---

## Verification Summary

- **Total Sources Cited**: 6
- **Reachable URLs Verified**: 6 / 6
- **Unsupported Claims Resolved**: All citations are now explicitly tied to peer-reviewed papers, official CNCF/Apache documentation, or recognized distributed systems literature.
