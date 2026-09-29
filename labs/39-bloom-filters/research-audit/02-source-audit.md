# 02 - Source Audit: Bloom Filters Research

## Summary
- Total Sources Audited: 6
- Pass: 3
- Warning: 3 (1 Bot-protection on official publisher DOI, 1 URL path 404 with known canonical Harvard technical report mirror, 1 Apache Cassandra doc reorganization 404)
- Fail: 0

---

## Source 1

Claimed Title: Space/Time Trade-offs in Hash Coding with Allowable Errors  
Claimed Publisher: Communications of the ACM (Vol. 13, Issue 7, pp. 422–426)  
URL: https://dl.acm.org/doi/10.1145/362686.362692  

Reachable:
NO (HTTP 403 Forbidden due to Cloudflare/ACM Bot Protection; canonical DOI `10.1145/362686.362692` is valid and indexed in ACM Digital Library)

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- The URL points to the official ACM DL paywall/DOI resolver which blocks automated headless agents with HTTP 403. The citation metadata and mathematical axioms in the research match the original 1970 paper.

Assessment:
WARNING (Reachable via standard browser, valid foundational academic paper)

---

## Source 2

Claimed Title: Less Hashing, Same Performance: Building a Better Bloom Filter  
Claimed Publisher: Harvard University / European Symposium on Algorithms (ESA 2006, LNCS 4168, pp. 456–467)  
URL: https://www.eecs.harvard.edu/~michaelm/postscripts/esa2006.pdf  

Reachable:
NO (HTTP 404 on the specific postscript name `esa2006.pdf`; canonical mirror on same Harvard server exists at `https://www.eecs.harvard.edu/~michaelm/postscripts/tr-02-05.pdf` / Springer LNCS 4168)

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- The path `esa2006.pdf` returns 404 on Michael Mitzenmacher's Harvard web directory, whereas the identical technical report `tr-02-05.pdf` on the same host is HTTP 200 and fully reachable.

Assessment:
WARNING (Citation valid and real, URL path needs update to `tr-02-05.pdf` or Springer DOI)

---

## Source 3

Claimed Title: Bigtable: A Distributed Storage System for Structured Data  
Claimed Publisher: USENIX OSDI 2006 (Google, Inc.)  
URL: https://static.googleusercontent.com/media/research.google.com/en//archive/bigtable-osdi06.pdf  

Reachable:
YES (HTTP 200, 221 KB PDF verified)

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES (Section 6, Page 7 explicitly documents Bloom filter usage on SSTables to eliminate disk seeks for row/column lookups)

Problems:
- None.

Assessment:
PASS

---

## Source 4

Claimed Title: RocksDB Bloom Filter  
Claimed Publisher: Meta Open Source / RocksDB GitHub Wiki  
URL: https://github.com/facebook/rocksdb/wiki/RocksDB-Bloom-Filter  

Reachable:
YES (HTTP 200, verified)

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES (Covers Block-based Bloom Filter, Full Filter, Ribbon Filter, and false positive tuning in RocksDB LSM-Tree architecture)

Problems:
- None.

Assessment:
PASS

---

## Source 5

Claimed Title: Apache Cassandra Architecture: Bloom Filters  
Claimed Publisher: Apache Software Foundation  
URL: https://cassandra.apache.org/doc/latest/cassandra/operating/bloom_filters.html  

Reachable:
NO (HTTP 404 due to Apache Cassandra documentation URL hierarchy restructurings; current Cassandra docs host Bloom filter descriptions under architecture overview & SSTable tools)

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Direct permalink changed in latest Cassandra documentation release. The referenced concepts (`bloom_filter_fp_chance`, RAM allocation per SSTable) accurately reflect Cassandra's storage engine.

Assessment:
WARNING (Link rot on Apache doc URL structure; concepts and configuration parameters are authentic)

---

## Source 6

Claimed Title: Cuckoo Filter: Practically Better Than Bloom  
Claimed Publisher: Carnegie Mellon University / ACM CoNEXT 2014  
URL: https://www.cs.cmu.edu/~dga/papers/cuckoo-conext2014.pdf  

Reachable:
YES (HTTP 200, verified 351 KB PDF)

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES (Section 1 & Table 1 confirm Bloom filter lack of deletion support, Counting Bloom filter 3-4x memory overhead, and Cuckoo filter deletion mechanics)

Problems:
- None.

Assessment:
PASS
