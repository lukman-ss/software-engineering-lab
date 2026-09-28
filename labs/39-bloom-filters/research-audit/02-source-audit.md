# Source Audit: Bloom Filters Research

## Source 1
Claimed Title: Bloom filter (general)
Claimed Publisher: Wikipedia community
URL: https://en.wikipedia.org/wiki/Bloom_filter
Reachable: YES
Source Type: SECONDARY
Relevant: YES
Supports Claimed Topic: YES
Problems: Secondary community encyclopedia; must be backed by original literature for core theorems.
Assessment: PASS

## Source 2
Claimed Title: Space/Time Trade-offs in Hash Coding with Allowable Errors
Claimed Publisher: Burton H. Bloom (CACM 1970)
URL: http://www.dragonwins.com/domains/getteched/bbc/literature/Bloom70.pdf
Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: Third-party hosted PDF mirror, but contents match original 1970 Communications of the ACM paper.
Assessment: PASS

## Source 3
Claimed Title: The log-structured merge-tree (LSM-tree)
Claimed Publisher: Patrick O'Neil et al. (Acta Informatica 1996)
URL: https://www.cs.umb.edu/~poneil/lsmtree.pdf
Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: None. Canonical paper on LSM-trees and on-disk Bloom filters.
Assessment: PASS

## Source 4
Claimed Title: Large-scale Incremental Processing Using Distributed Transactions and Notifications (Percolator)
Claimed Publisher: Daniel Peng, Frank Dabek (USENIX OSDI 2010)
URL: https://research.google/pubs/pub36726/
Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: PARTIAL
Problems: Percolator uses Bigtable SSTable Bloom filters beneath its transaction layer; it does not introduce a standalone application-level Bloom filter cache. The research notes clarify this distinction.
Assessment: PASS

## Source 5
Claimed Title: BitFunnel – Search Engine Indexing Algorithm
Claimed Publisher: Bob Goodwin et al. / Wikipedia
URL: https://en.wikipedia.org/wiki/BitFunnel
Reachable: YES
Source Type: SECONDARY
Relevant: YES
Supports Claimed Topic: YES
Problems: Primary paper is Goodwin et al. (SIGIR 2017); Wikipedia page accurately summarizes bit-sliced signature application.
Assessment: PASS

## Source 6
Claimed Title: MurmurHash3 performance and design notes
Claimed Publisher: Austin Appleby / SmHasher wiki
URL: https://github.com/aappleby/smhasher/wiki/MurmurHash3
Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: Project wiki documentation; accurately reflects standard hash speeds and non-cryptographic hash selection rationale.
Assessment: PASS

## Source 7
Claimed Title: Fowler–Noll–Vo (FNV) hash function
Claimed Publisher: Wikipedia community
URL: https://en.wikipedia.org/wiki/Fowler%E2%80%93Noll%E2%80%93Vo_hash_function
Reachable: YES
Source Type: SECONDARY
Relevant: YES
Supports Claimed Topic: YES
Problems: Secondary reference; corroborated by primary author source (Source 12).
Assessment: PASS

## Source 8
Claimed Title: Cuckoo filters – replacing Bloom filters with deletions
Claimed Publisher: Fan, Zhao, et al. (ACM CoNEXT 2014)
URL: https://www.cs.cmu.edu/~dga/papers/cuckoo-conext2014.pdf
Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: Contains canonical comparison between Bloom filters and Cuckoo filters regarding deletions and space efficiency.
Assessment: PASS

## Source 9
Claimed Title: An overview of log-structured merge trees
Claimed Publisher: Ben Stopford
URL: http://www.benstopford.com/2015/02/14/log-structured-merge-trees/
Reachable: YES
Source Type: COMMUNITY
Relevant: YES
Supports Claimed Topic: YES
Problems: Informal blog post; serves as practical explanation rather than primary academic authority. Corroborates Source 3 and Source 10.
Assessment: PASS

## Source 10
Claimed Title: LSM-based storage techniques: a survey
Claimed Publisher: Chen Luo, Michael J. Carey (ACM Computing Surveys / arXiv 2019)
URL: https://arxiv.org/abs/1812.07527
Reachable: YES
Source Type: PRIMARY / SECONDARY (Peer-reviewed survey)
Relevant: YES
Supports Claimed Topic: YES
Problems: None. Authoritative peer-reviewed survey of modern LSM engines (RocksDB, Cassandra, AsterixDB).
Assessment: PASS

## Source 11
Claimed Title: Bloom filter false-positive probability analysis
Claimed Publisher: Wikipedia
URL: https://en.wikipedia.org/wiki/Bloom_filter#Probability_of_false_positives
Reachable: YES
Source Type: SECONDARY
Relevant: YES
Supports Claimed Topic: YES
Problems: Sub-section of Source 1; provides exact Stirling numbers and Azuma-Hoeffding concentration references.
Assessment: PASS

## Source 12
Claimed Title: FNV hash – practical notes
Claimed Publisher: Landon Curt Noll et al. (ISTHE.COM)
URL: https://www.isthe.com/chongo/tech/comp/fnv/index.html
Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: Canonical reference page maintained by FNV authors.
Assessment: PASS
