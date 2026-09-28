# Source Audit

## Source 1
Claimed Title: Bloom filter (general)
Claimed Publisher: Wikipedia community
URL: https://en.wikipedia.org/wiki/Bloom_filter

Reachable:
YES

Source Type:
COMMUNITY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Primary reliance on community source for theoretical algorithms and guarantees where original literature should be cited.

Assessment:
PASS

---

## Source 2
Claimed Title: Space/Time Trade-offs in Hash Coding with Allowable Errors
Claimed Publisher: Burton H. Bloom / Communications of the ACM (1970)
URL: http://www.dragonwins.com/domains/getteched/bbc/literature/Bloom70.pdf

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. PDF is accessible and matches the original CACM publication.

Assessment:
PASS

---

## Source 3
Claimed Title: The log-structured merge-tree (LSM-tree)
Claimed Publisher: Patrick O'Neil et al. / Acta Informatica (1996)
URL: https://www.cs.umb.edu/~poneil/lsmtree.pdf

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Original paper detailing LSM-tree components.

Assessment:
PASS

---

## Source 4
Claimed Title: Large-scale Incremental Processing Using Distributed Transactions and Notifications (Percolator)
Claimed Publisher: Daniel Peng, Frank Dabek / USENIX OSDI 2010
URL: https://research.google/pubs/pub36726/

Reachable:
YES (redirects to HTTPS canonical URL)

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
PARTIAL

Problems:
- The research claims Percolator uses Bloom filters to prevent "cache penetration". However, Percolator's primary mechanism for avoiding disk checks relies on Bigtable's internal Bloom filters rather than an application-layer cache penetration guard.

Assessment:
WARNING

---

## Source 5
Claimed Title: BitFunnel – Search Engine Indexing Algorithm
Claimed Publisher: Bob Goodwin et al. / Wikipedia / SIGIR 2017
URL: https://en.wikipedia.org/wiki/BitFunnel

Reachable:
YES

Source Type:
COMMUNITY / SECONDARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Cited URL is Wikipedia rather than the primary SIGIR 2017 paper ("BitFunnel: Revisiting Bit-sliced Signatures for Search Engines").

Assessment:
PASS

---

## Source 6
Claimed Title: MurmurHash3 performance and design notes
Claimed Publisher: Austin Appleby / SmHasher project
URL: https://github.com/aappleby/smhasher/wiki/MurmurHash3

Reachable:
YES

Source Type:
SECONDARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. SmHasher wiki provides canonical MurmurHash3 speed benchmarks (~2.5-5 GB/s).

Assessment:
PASS

---

## Source 7
Claimed Title: Fowler–Noll–Vo (FNV) hash function
Claimed Publisher: Wikipedia community
URL: https://en.wikipedia.org/wiki/Fowler%E2%80%93Noll%E2%80%93Vo_hash_function

Reachable:
YES

Source Type:
COMMUNITY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None.

Assessment:
PASS

---

## Source 8
Claimed Title: Cuckoo filters – replacing Bloom filters with deletions
Claimed Publisher: Fan, Zhao, et al. / CMU (2014)
URL: https://www.cs.cmu.edu/~fanzhao/cuckoo-filter.pdf

Reachable:
NO (Returns HTTP 404 Not Found)

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
UNVERIFIED

Problems:
- Link is broken (404). Canonical CoNEXT 2014 paper URL or ACM DL link was not provided.

Assessment:
FAIL

---

## Source 9
Claimed Title: An overview of log-structured merge trees
Claimed Publisher: Ben Stopford
URL: http://www.benstopford.com/2015/02/14/log-structured-merge-trees/

Reachable:
YES (redirects to HTTPS)

Source Type:
COMMUNITY / SECONDARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Personal technical blog. Good for context, but secondary to academic LSM literature.

Assessment:
PASS

---

## Source 10
Claimed Title: LSM-based storage techniques: a survey
Claimed Publisher: Chen Luo, Michael J. Carey / arXiv (2019)
URL: https://arxiv.org/abs/1812.07527

Reachable:
YES

Source Type:
SECONDARY (Academic Survey)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Abstract landing page does not directly display the term Bloom without downloading full PDF text, but document content verifies per-component Bloom filter survey.

Assessment:
PASS

---

## Source 11
Claimed Title: Bloom filter false-positive probability analysis
Claimed Publisher: Wikipedia (Probability of false positives section)
URL: https://en.wikipedia.org/wiki/Bloom_filter#Probability_of_false_positives

Reachable:
YES

Source Type:
COMMUNITY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Cites Goel & Gupta (2007) for rigorous upper bounds without inspecting the primary paper.

Assessment:
PASS

---

## Source 12
Claimed Title: FNV hash – practical notes
Claimed Publisher: Landon Curt Noll et al. (ISTHE.COM)
URL: https://www.isthe.com/chongo/tech/comp/fnv/index.html

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Authoritative reference for FNV parameters and avalanche analysis.

Assessment:
PASS
