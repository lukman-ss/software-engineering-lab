
# Sources

## Source 1
Title: Bloom filter (general)
Author(s): Wikipedia community
Date: 2026‑09‑28 (accessed)
URL: https://en.wikipedia.org/wiki/Bloom_filter
Type: Wikipedia
Summary: Broad overview of Bloom filters—algorithm description, false‑positive probability, space/time trade‑offs, and multiple extensions (counting Bloom, Scalable Bloom, Ribbon filter, etc.).

## Source 2
Title: Space/Time Trade-offs in Hash Coding with Allowable Errors
Author: Burton H. Bloom
Date: 1970
URL: http://www.dragonwins.com/domains/getteched/bbc/literature/Bloom70.pdf
Type: Original paper (PDF)
Summary: Introduces Bloom filters. Presents the core bit‑array algorithm, derives the false‑positive probability bound, and discusses early memory‑efficiency calculations.

## Source 3
Title: The log-structured merge-tree (LSM-tree)
Authors: Patrick O'Neil, Edward Cheng, Dieter Gawlick, Elizabeth O'Neil
Date: 1996
URL: https://www.cs.umb.edu/~poneil/lsmtree.pdf
Type: Academic paper
Summary: Classic LSM‑tree architecture. Shows how per‑SST Bloom filters reduce disk reads in key‑value stores. Includes a formula for expected point‑lookup cost with Bloom filters: O(L·e^(−M/N)).

## Source 4
Title: Large-scale Incremental Processing Using Distributed Transactions and Notifications (Percolator)
Authors: Daniel Peng, Frank Dabek
Date: 2010
URL: https://research.google/pubs/pub36726/
Type: USENIX Symposium paper
Summary: Google research on distributed transactions built on Bigtable. Bloom filter usage for avoiding unnecessary disk reads resides in the underlying Bigtable SSTable layer (per-SST filters), not as an application-level cache-penetration barrier in Percolator itself.

## Source 5
Title: BitFunnel – Search Engine Indexing Algorithm
Authors: Bob Goodwin et al.
Date: 2017 (original); Wikipedia accessed 2026‑09‑28
URL: https://en.wikipedia.org/wiki/BitFunnel
Type: Wikipedia
Summary: Describes BitFunnel's use of Bloom‑like signatures to replace inverted indexes in Bing search. Confirms the algorithmic equivalence and provides practical performance motivation.

## Source 6
Title: MurmurHash3 performance and design notes
Author: Austin Appleby / SmHasher project
Date: 2016 (last update on wiki)
URL: https://github.com/aappleby/smhasher/wiki/MurmurHash3
Type: Project documentation
Summary: Details MurmurHash3 variants (x86_32, x86_128, x64_128). Provides benchmark speeds (e.g., ~2.5–5 GB/s bulk hashing). Useful for selecting hash functions in a Bloom‑filter implementation.

## Source 7
Title: Fowler–Noll–Vo (FNV) hash function
Author: Wikipedia community
Date: 2026‑07‑30 (accessed)
URL: https://en.wikipedia.org/wiki/Fowler%E2%80%93Noll%E2%80%93Vo_hash_function
Type: Wikipedia
Summary: Explains FNV‑1 and FNV‑1a designs. Highlights that FNV‑1a has better avalanche characteristics and provides the offset basis and prime parameters for 32/64/128‑bit variants.

## Source 8
Title: Cuckoo filters – replacing Bloom filters with deletions
Authors: Fan, Zhao, et al.
Date: 2014 (presentation)
URL: https://www.cs.cmu.edu/~dga/papers/cuckoo-conext2014.pdf
Type: Conference presentation
Summary: Shows how Cuckoo filters remove the deletion limitation of classic Bloom filters while keeping similar space efficiency. Contains empirical FP‑rate comparison (Figure 3).
Corrected URL: https://dl.acm.org/doi/10.1145/2674005.2674994 (ACM DL)

## Source 9
Title: An overview of log-structured merge trees
Author: Ben Stopford
Date: 2015‑02‑14
URL: http://www.benstopford.com/2015/02/14/log-structured-merge-trees/
Type: Technical blog
Summary: Clear explanation of LSM‑tree internals. States that each SST file contains a Bloom filter to avoid unnecessary disk reads, reinforcing Source 3.

## Source 10
Title: LSM-based storage techniques: a survey
Authors: Chen Luo, Michael J. Carey
Date: July 2019
URL: https://arxiv.org/abs/1812.07527
Type: Academic survey (arXiv)
Summary: Comprehensive review of LSM‑tree variants and their use of Bloom filters across levels. Provides the false‑positive probability formula and discusses read/write trade‑offs.

## Source 11
Title: Bloom filter false-positive probability analysis
Author: Wikipedia (Probability of false positives section)
Date: 2026‑09‑28 (accessed)
URL: https://en.wikipedia.org/wiki/Bloom_filter#Probability_of_false_positives
Type: Wikipedia
Summary: Derives the standard approximation ε ≈ (1−e^(−kn/m))^k and the optimal k = (m/n)·ln 2. Gives the exact formula with Stirling numbers and references Goel & Gupta (2007) rigorous bound.

## Source 12
Title: FNV hash – practical notes
Author: Landon Curt Noll et al. (ISTHE.COM)
Date: 2020 (cited in Wikipedia)
URL: https://www.isthe.com/chongo/tech/comp/fnv/index.html
Type: Reference page
Summary: Authoritative source for FNV parameters. Explains why FNV‑1a is preferred for Bloom filters (better avalanche behavior) and provides offset basis and prime tables.
