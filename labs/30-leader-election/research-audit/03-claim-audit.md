# Claim Audit: Leader Election Research

Target Lab: `labs/30-leader-election`  
Audit Date: September 28, 2026  

---

## Claim 1

Claim: In Raft, randomized election timeouts (typically 150ms–300ms) prevent split votes and split-brain when electing leaders.  
Location: `03-mechanisms-and-leases.md:15-28`  
Evidence Provided: Cited Ongaro & Ousterhout (USENIX ATC '14), Sections 5.2–5.3.  
Source: `https://raft.github.io/raft.pdf`  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Accurate representation of the Raft consensus mechanism.

---

## Claim 2

Claim: Lease/heartbeat leader election alone provides time-bounded mutual exclusion, but cannot guarantee safety against stop-the-world GC pauses or long network delays without fencing tokens.  
Location: `03-mechanisms-and-leases.md:83-115`, `04-fencing-and-redlock.md:93-147`  
Evidence Provided: Kleppmann (2016) "How to do distributed locking" and DDIA Chapter 8.  
Source: `https://martin.kleppmann.com/2016/02/08/how-to-do-distributed-locking.html`  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Supported by Gray & Cheriton (1989), Burrows (OSDI 2006), and Kleppmann (2016).

---

## Claim 3

Claim: Redis Redlock algorithm does not generate monotonically increasing fencing tokens and relies on semi-synchronous timing assumptions (bounded drift, bounded network delay, bounded pauses).  
Location: `04-fencing-and-redlock.md:15-27`, `152-164`  
Evidence Provided: Kleppmann (2016) and antirez (2016).  
Source: `https://martin.kleppmann.com/2016/02/08/how-to-do-distributed-locking.html`, `http://antirez.com/news/101`  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: antirez acknowledges Redlock returns a random token and relies on semi-synchronous timing bounds with drift bounds, while Kleppmann highlights the lack of monotonic fencing.

---

## Claim 4

Claim: etcd uses `clientv3/concurrency` where `Revision` numbers serve as strictly monotonic fencing tokens, and ZooKeeper uses `zxid` / znode versions as fencing tokens.  
Location: `04-fencing-and-redlock.md:40-43`, `59-60`, `117-124`  
Evidence Provided: etcd concurrency documentation and Apache ZooKeeper recipes.  
Source: `https://etcd.io/docs/v3.5/tutorials/how-to-conduct-elections/`, `https://zookeeper.apache.org/doc/current/recipes.html#sc_leaderElection`  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Consistent with official etcd and ZooKeeper architectures.

---

## Claim 5

Claim: ZooKeeper leader election avoids herd effects by having candidates watch only their immediate predecessor znode in the sequence (`ELECTION/guid-n_j` where j < i).  
Location: `04-fencing-and-redlock.md:53-54`, `05-system-comparison-and-best-practices.md:80-86`  
Evidence Provided: Apache ZooKeeper Recipes documentation.  
Source: `https://zookeeper.apache.org/doc/current/recipes.html#sc_leaderElection`  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Exactly matches the algorithm described in ZooKeeper Recipes section "Leader Election".

---

## Claim 6

Claim: Distributed locks without fencing (such as single-node Redis `SET NX EX`) are suitable for efficiency optimizations (e.g. deduplication of background jobs) where idempotency exists, but unsafe for correctness-critical operations (e.g., financial transactions).  
Location: `04-fencing-and-redlock.md:133-142`, `194-198`, `05-system-comparison-and-best-practices.md:19-26`  
Evidence Provided: Mike Burrows (Chubby, OSDI 2006) and Kleppmann (2016).  
Source: `https://martin.kleppmann.com/2016/02/08/how-to-do-distributed-locking.html`  
Source Actually Supports Claim: YES  
Classification: FACT / BEST_PRACTICE  
Severity: LOW  
Notes: Widely recognized architectural taxonomy first formalized by Burrows.

---

## Claim 7

Claim: Setting lease TTL as `TTL = 3 × (max GC pause + max RTT)` provides a reliable empirical heuristic for lease sizing.  
Location: `05-system-comparison-and-best-practices.md:62-69`  
Evidence Provided: etcd and ZooKeeper operational recommendations.  
Source: `https://etcd.io/docs/v3.5/tutorials/how-to-conduct-elections/`, `https://zookeeper.apache.org/doc/current/recipes.html`  
Source Actually Supports Claim: YES  
Classification: RECOMMENDATION  
Severity: LOW  
Notes: Nuanced engineering rule of thumb based on 99th percentile measurements.
