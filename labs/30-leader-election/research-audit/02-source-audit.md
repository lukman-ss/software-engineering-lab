# Source Audit: Leader Election Research

Target Lab: `labs/30-leader-election`  
Audit Date: September 28, 2026  

---

## Source 1

Claimed Title: In Search of an Understandable Consensus Algorithm  
Claimed Publisher: USENIX Annual Technical Conference (ATC '14)  
URL: https://raft.github.io/raft.pdf  

Reachable: YES  
Source Type: PRIMARY (Academic Peer-Reviewed Paper)  
Relevant: YES  
Supports Claimed Topic: YES  

Problems:
- None. URL serves the definitive peer-reviewed Raft paper by Diego Ongaro and John Ousterhout.

Assessment: PASS

---

## Source 2

Claimed Title: How to do distributed locking  
Claimed Publisher: University of Cambridge / Personal Technical Blog  
URL: https://martin.kleppmann.com/2016/02/08/how-to-do-distributed-locking.html  

Reachable: YES  
Source Type: PRIMARY / Reputable Technical Reference  
Relevant: YES  
Supports Claimed Topic: YES  

Problems:
- None. Martin Kleppmann's authoritative analysis of distributed locking, fencing tokens, and Redlock critique.

Assessment: PASS

---

## Source 3

Claimed Title: Is Redlock safe?  
Claimed Publisher: Redis / antirez blog  
URL: http://antirez.com/news/101  

Reachable: YES  
Source Type: PRIMARY (Original Author Technical Rebuttal)  
Relevant: YES  
Supports Claimed Topic: YES  

Problems:
- None. Salvatore Sanfilippo's official response and technical rebuttal to Martin Kleppmann.

Assessment: PASS

---

## Source 4

Claimed Title: How to conduct leader election in etcd cluster / etcd Concurrency API  
Claimed Publisher: etcd / CNCF  
URL: https://etcd.io/docs/v3.5/tutorials/how-to-conduct-elections/  

Reachable: YES  
Source Type: PRIMARY (Official CNCF Documentation)  
Relevant: YES  
Supports Claimed Topic: YES  

Problems:
- None. Canonical documentation on leader election tutorial using etcd v3 and concurrency API.

Assessment: PASS

---

## Source 5

Claimed Title: ZooKeeper Recipes and Solutions: Leader Election  
Claimed Publisher: The Apache Software Foundation  
URL: https://zookeeper.apache.org/doc/current/recipes.html#sc_leaderElection  

Reachable: YES  
Source Type: PRIMARY (Official Open Source Project Documentation / Standards Reference)  
Relevant: YES  
Supports Claimed Topic: YES  

Problems:
- None. The official Apache ZooKeeper recipes documentation detailing ephemeral sequential znode leader election and herd effect mitigation.

Assessment: PASS

---

## Source 6

Claimed Title: Designing Data-Intensive Applications: The Big Ideas Behind Reliable, Scalable, and Maintainable Systems  
Claimed Publisher: O'Reilly Media  
URL: https://dataintensive.net/  

Reachable: YES  
Source Type: PRIMARY / Reputable Technical Reference (Definitive Technical Textbook)  
Relevant: YES  
Supports Claimed Topic: YES  

Problems:
- None. Official book website for Martin Kleppmann's DDIA, widely recognized reference on distributed consensus, timing assumptions, and distributed locking.

Assessment: PASS

---

## Source Audit Summary

- Total Sources Checked: 6
- Reachable: 6 / 6 (100%)
- Authoritative / Primary: 6 / 6 (100%)
- Overall Source Quality: High-grade technical and peer-reviewed literature.
