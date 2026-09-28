# Changes Made

## Revision 1 — Fix Incomplete Research Deliverables (CRITICAL)

Audit Issue: CRITICAL — Gap 1 / Verdict Blocking 1: Only `research/01-plan.md` exists. 5 research questions unanswered, no synthesis documents.
Severity: CRITICAL

Files Changed:
- `research/01-plan.md` (updated — marked risks resolved, deliverables listed)
- `research/03-mechanisms-and-leases.md` (created)
- `research/04-fencing-and-redlock.md` (created)
- `research/05-system-comparison-and-best-practices.md` (created)

Action:
- Created 03-mechanisms-and-leases.md answering RQ1 (Raft vs lease-based mechanisms, comparison matrix) and RQ2 (lease/heartbeat split-brain prevention, clock/pause dependencies, fencing necessity)
- Created 04-fencing-and-redlock.md answering RQ3 (Redis/etcd/ZooKeeper/Consul tradeoffs) and RQ4 (fencing tokens, GC pause scenario, Kleppmann vs antirez deep dive)
- Created 05-system-comparison-and-best-practices.md answering RQ5 (production best practices, etcd/ZK/Consul/Redis guidance, architecture diagram)
- Updated 01-plan.md to reflect completed deliverables and resolved unknowns

Verification:
- All 5 research questions now have dedicated sections with cited sources
- Structural contradiction (04-contradictions.md) resolved — research directory now complete

Status: RESOLVED

---

## Revision 2 — Fix Missing Validated Sources (HIGH)

Audit Issue: HIGH — Gap 2 / Source Audit FAIL: 6 generic placeholders, 0 concrete URLs/DOIs, 0 reachable verifications. Verdict Blocking 2.
Severity: HIGH

Files Changed:
- `research/02-sources.md` (created)

Action:
- Replaced 6 placeholders with 6 verified sources, each with title, publisher, exact URL, reachable status, type, relevance, extraction notes:
  1. Raft paper — https://raft.github.io/raft.pdf (peer-reviewed)
  2. Kleppmann distributed locking — https://martin.kleppmann.com/2016/02/08/how-to-do-distributed-locking.html
  3. antirez rebuttal — http://antirez.com/news/101
  4. etcd election — https://etcd.io/docs/v3.5/tutorials/how-to-conduct-elections/
  5. ZooKeeper recipes — https://zookeeper.apache.org/doc/current/recipes.html#sc_leaderElection
  6. DDIA — https://dataintensive.net/
- Verified reachability via webfetch (6/6 YES) on 2026-09-28
- All research docs now cite sources inline with exact URLs

Verification:
- Source Integrity gate: FAIL → PASS (6/6 reachable, publisher/title verified)
- URLs re-fetched and content confirmed to support claimed extractions

Status: RESOLVED

---

## Revision 3 — Fix Unsupported Claims & Unanalyzed Risks (HIGH)

Audit Issue: HIGH — Claims 1+2 unsupported (split-brain claim, Redlock debate claim NOT VERIFIED). Gap 3: fencing tokens, clock skew, lease expiration unanalyzed. Verdict Blocking 3.
Severity: HIGH (Claim 2), MEDIUM (Claim 1)

Files Changed:
- `research/03-mechanisms-and-leases.md`
- `research/04-fencing-and-redlock.md`

Action:
- Narrowed Claim 1 ("Leader election avoids split-brain") → qualified: avoids split-brain only with quorum intersection (Raft) or fencing tokens + bounded clock/pause assumptions; added GC pause and clock-jump failure modes with sources
- Grounded Claim 2 (Redlock debate) → cited both Kleppmann (2016) and antirez (2016) primary URLs, extracted 4 Kleppmann arguments and 5 antirez counter-arguments, provided resolution table
- Added dedicated sections: fencing token mechanics, monotonic counter table (etcd revision/ZK zxid/Raft term), GC pause walkthrough (10s pause vs 5s TTL), clock drift bounds
- All numeric examples labeled (T=5s, timeout 150-300ms) with source-anchored values

Verification:
- Claim Support gate: FAIL → PASS (2/2 claims now supported with primary sources)
- Source/claim mismatch eliminated — every technical claim traceable to 02-sources.md entry

Status: RESOLVED

---

## Summary

Total Audit Issues Addressed: 7 (3 blocking + 2 claim + 1 contradiction + 1 source audit)
Resolved: 7
Partially Resolved: 0
Unresolved: 0
Files Created: 4 (02-sources.md, 03-mechanisms-and-leases.md, 04-fencing-and-redlock.md, 05-system-comparison-and-best-practices.md)
Files Modified: 1 (01-plan.md)
