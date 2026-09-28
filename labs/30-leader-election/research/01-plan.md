# Research Plan

## Research Topic
Leader Election --- Menghindari Split-Brain dan Eksekusi Ganda di Multi-Node

## Objective
Investigate leader election mechanisms in distributed systems, focusing on lease/heartbeat patterns, split-brain prevention, and practical implementations using Redis, etcd, and consensus algorithms like Raft.

## Research Questions
1. What are the fundamental mechanisms for leader election in distributed systems?
2. How do lease/heartbeat mechanisms work to prevent split-brain scenarios?
3. What are the trade-offs between different implementations (Redis Redlock, etcd, Consul, ZooKeeper)?
4. How do fencing tokens/generation IDs prevent split-brain during GC pauses or network partitions?
5. What are the best practices for implementing leader election in production systems?

## Search Strategy
- Search for authoritative documentation on leader election patterns
- Review academic papers on distributed consensus (Raft, Paxos)
- Examine official documentation for Redis Redlock, etcd, Consul, ZooKeeper
- Search for split-brain prevention techniques and fencing mechanisms
- Look for production case studies and best practices

## Expected Primary Sources
- Redis official documentation (Redlock algorithm)
- etcd documentation (lease mechanism)
- Raft consensus algorithm paper
- ZooKeeper documentation
- Consul documentation
- Distributed systems textbooks (e.g., "Designing Data-Intensive Applications")

## Risks / Unknowns
- Conflicting advice on Redis Redlock safety (Martin Kleppmann vs Redis authors)
- Practical performance characteristics under network partitions
- Exact implementation details for fencing tokens in various systems
- How different systems handle clock skew in lease expiration