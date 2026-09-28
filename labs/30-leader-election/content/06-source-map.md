# Source Map

## Leader Election with Fencing Tokens — Overview

Research:
- research/01-plan.md
- research/02-sources.md
- research/03-mechanisms-and-leases.md
- research/04-fencing-and-redlock.md
- research/05-system-comparison-and-best-practices.md

Implementation:
- internal/coordinator/coordinator.go
- internal/candidate/candidate.go
- internal/storage/storage.go
- cmd/demo/main.go

Tests:
- tests/election_test.go

## Problem and Mental Model

Research:
- research/03-mechanisms-and-leases.md (section "Risk: Unbounded Pauses and Clock Skew")
- research/04-fencing-and-redlock.md (section "The Problem: Paused Clients and Expired Leases")

Implementation:
- internal/candidate/candidate.go (SimulatePause)

## Lease and Heartbeat Mechanism

Research:
- research/03-mechanisms-and-leases.md (section "Lease-Based Leader Election (etcd Model)")
- research/02-sources.md (source #4 etcd documentation)

Implementation:
- internal/coordinator/coordinator.go (Acquire, Renew)
- internal/candidate/candidate.go (runElectionLoop)

## Fencing Token Generation

Research:
- research/04-fencing-and-redlock.md (section "Generation in Various Systems")
- research/02-sources.md (source #2 Kleppmann, source #4 etcd docs)

Implementation:
- internal/coordinator/coordinator.go (Acquire — `c.revision++`)

## Fencing Check in Shared Storage

Research:
- research/04-fencing-and-redlock.md (section "Fencing Solution: Monotonically Increasing Tokens")
- research/04-fencing-and-redlock.md (section "Requirements for Correct Fencing Implementation")

Implementation:
- internal/storage/storage.go (Write)

Tests:
- tests/election_test.go (TestFencedStorageRejectsStaleTokens)

## Failover on Lease Expiration

Research:
- research/03-mechanisms-and-leases.md (section "Core Principle: Time-Bounded Exclusivity")
- research/02-sources.md (source #4 etcd, source #5 ZooKeeper)

Implementation:
- internal/coordinator/coordinator.go (Renew — expired path)
- internal/candidate/candidate.go (runElectionLoop — follower path)

Tests:
- tests/election_test.go (TestLeaderElectionFailoverAndSplitBrainDefense)

## Split-Brain Prevention

Research:
- research/03-mechanisms-and-leases.md (section "Mitigation: Fencing Tokens")
- research/04-fencing-and-redlock.md (section "Fencing Is Not Optional for Correctness-Critical Operations")

Implementation:
- internal/storage/storage.go (Write — check-and-set)
- internal/candidate/candidate.go (PerformFencedWrite)

Tests:
- tests/election_test.go (TestLeaderElectionFailoverAndSplitBrainDefense)

## Concurrent Election Safety

Research:
- research/03-mechanisms-and-leases.md (section "Consensus-Based Leader Election (Raft Model)")
- research/02-sources.md (source #1 Raft paper)

Implementation:
- internal/coordinator/coordinator.go (Acquire under mutex)

Tests:
- tests/election_test.go (TestConcurrentElectionRace)

## Redlock Debate and System Comparisons

Research:
- research/04-fencing-and-redlock.md (section "The Redlock Safety Debate: Kleppmann vs. antirez")
- research/05-system-comparison-and-best-practices.md

## Production Best Practices

Research:
- research/05-system-comparison-and-best-practices.md (section "General Production Best Practices")
- research/02-sources.md

## Demo Execution

Implementation:
- cmd/demo/main.go

Engineering:
- engineering/01-design.md
- engineering/02-implementation-notes.md
- engineering/03-execution-result.md

## Audit and Verdicts

Research-Audit:
- research-audit/07-verdict.md (APPROVED)
- research-audit/06-gaps.md (minor gap noted, non-blocking)
- research-revision/03-revision-result.md (revision history)

Engineering-Audit:
- engineering-audit/06-verdict.md (APPROVED)
- engineering-audit/05-gaps.md (no gaps)

Engineering-Audit-OpenSource:
- engineering-audit-opensource/06-verdict.md (APPROVED, noted transient dual-leader state mitigated by fencing)
- engineering-audit-opensource/05-gaps.md (non-blocking note only)

## Source of Truth

Semua klaim pada artikel ini bersumber dari:
- dokumen riset yang sudah di-approve (`research/`)
- dokumen engineering dan execution result (`engineering/`)
- implementasi aktual (`internal/`, `cmd/`)
- test suite (`tests/election_test.go`)
- hasil audit (`research-audit/`, `engineering-audit/`, `engineering-audit-opensource/`)

Tidak ada klaim faktual eksternal ditambahkan. Angka TTL dan interval pada demo dan test bersifat ilustratif sesuai penjelasan pada `engineering/02-implementation-notes.md` dan research `05-system-comparison-and-best-practices.md`.
