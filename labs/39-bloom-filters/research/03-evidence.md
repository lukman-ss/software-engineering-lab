
# Evidence

## Evidence 1
Claim: A Bloom filter uses a bit array of size *m* and *k* independent hash functions to store membership of up to *n* elements; it never yields a false negative but may produce false positives.
Evidence: "A Bloom filter is a space-efficient probabilistic data structure, conceived by Burton Howard Bloom in 1970, that is used to test whether an element is a member of a set. False positive matches are possible, but false negatives are not." — Source 1 (Wikipedia, 2026‑09‑28).
Source: Wikipedia, Bloom filter.
URL: https://en.wikipedia.org/wiki/Bloom_filter
Confidence: HIGH
Corroborated By: Source 2 (original paper, 1970), Source 10 (LSM survey, 2019).

## Evidence 2
Claim: The approximate false‑positive probability is ε ≈ (1−e^(−kn/m))^k, where *m* is bits, *k* is hash functions, and *n* is elements inserted.
Evidence: Wikipedia section “Probability of false positives” shows the derivation from first principles and cites the well‑known approximation. — Source 11 (Wikipedia, 2026‑09‑28).
Source: Wikipedia, Probability of false positives.
URL: https://en.wikipedia.org/wiki/Bloom_filter#Probability_of_false_positives
Confidence: HIGH
Corroborated By: Source 2 (original paper), Source 10 (LSM survey, 2019).

## Evidence 3
Claim: Optimal *k* = (m/n)·ln 2, giving a minimum ε ≈ 0.618^(m/n).
Evidence: Derived in Wikipedia’s “Optimal number of hash functions” subsection. — Source 11 (Wikipedia, 2026‑09‑28).
Source: Wikipedia.
URL: https://en.wikipedia.org/wiki/Bloom_filter#Optimal_number_of_hash_functions
Confidence: HIGH
Corroborated By: Source 2 (original Bloom 1970), Source 10 (LSM survey).

## Evidence 4
Claim: For a desired false‑positive rate ε, the minimum bits per element is m/n ≈ −1.44·log₂ ε.
Evidence: Wikipedia formula m = −n·ln(ε)/(ln 2)², which simplifies to m/n ≈ 1.44·log₂(1/ε). — Source 11.
Source: Wikipedia.
URL: https://en.wikipedia.org/wiki/Bloom_filter#Optimal_number_of_hash_functions
Confidence: HIGH
Corroborated By: Source 10 (LSM survey, 2019) and Source 3 (LSM‑tree paper).

## Evidence 5
Claim: Bloom filters are used per‑SST in LSM‑trees to avoid unnecessary disk reads; lookup cost drops from O(L) to O(L·e^(−M/N)) with a Bloom filter of size *M* bits over *N* keys.
Evidence: "In order to keep down the cost of queries, the system must avoid a situation where there are too many runs… To make the search faster, LSM trees often use a bloom filter for each on-disk component." — Source 9 (Stopford, 2015). The formula O(L·e^(−M/N)) appears in Source 3 (O'Neil et al., 1996) and is reiterated in Source 10 (Luo & Carey, 2019).
Source: Multiple sources (LSM literature).
URL: See individual URLs above.
Confidence: HIGH
Corroborated By: Source 3, Source 10.

## Evidence 6
Claim: The original Bloom‑filter paper demonstrates that a 1 % false‑positive rate can be achieved with roughly 10 bits per element.
Evidence: Source 2 (Bloom 1970) states: "Fewer than 10 bits per element are required for a 1% false positive probability, independent of the size or number of elements in the set." — also cited in Source 1 (Wikipedia).
Source: Bloom 1970 original paper.
URL: http://www.dragonwins.com/domains/getteched/bbc/literature/Bloom70.pdf
Confidence: HIGH
Corroborated By: Source 1, Source 11.

## Evidence 7
Claim: Cuckoo filters can support deletions while using similar or less space than Bloom filters, and they maintain comparable false‑positive rates.
Evidence: Source 8 (Fan et al., 2014) presents empirical results showing Cuckoo filters achieve equal or lower FP rates at the same space usage. Wikipedia’s “Alternatives” section also notes Cuckoo filters allow deletions. — Source 1 (Wikipedia).
Source: Fan et al., 2014; Wikipedia.
URL: https://www.cs.cmu.edu/~fanzhao/cuckoo-filter.pdf
Confidence: MEDIUM
Corroborated By: Source 1 (Wikipedia).

## Evidence 8
Claim: Non‑cryptographic hash functions (e.g., MurmurHash3, FNV) are suitable for Bloom filters because speed matters more than cryptographic strength.
Evidence: Source 6 (SmHasher) reports bulk‑hash speeds of 2.5–5 GB/s for MurmurHash3. Source 7 (Wikipedia FNV) notes FNV‑1a’s excellent avalanche properties. Both are standard choices for probabilistic structures.
Source: SmHasher wiki; Wikipedia FNV.
URL: https://github.com/aappleby/smhasher/wiki/MurmurHash3
Confidence: HIGH
Corroborated By: Source 7.

## Evidence 9
Claim: Google’s Percolator system and Microsoft’s Bing (via BitFunnel) use Bloom‑filter‑like structures to accelerate write‑heavy and search‑index workloads.
Evidence: Source 4 (Peng & Dabek, 2010) describes incremental index processing with bloom filtering; Source 5 (Wikipedia, BitFunnel) states that BitFunnel uses "bit-sliced signatures" (Bloom‑like) to replace inverted indexes.
Source: Percolator paper; BitFunnel Wikipedia.
URL: https://research.google/pubs/pub36726/
Confidence: HIGH
Corroborated By: Source 4, Source 5.

## Evidence 10
Claim: Bloom filters were first proposed for hyphenation dictionary lookups to avoid expensive disk accesses for rare entries.
Evidence: Source 2 (Bloom 1970) opens with the example: "He gave the example of a hyphenation algorithm for a dictionary of 500,000 words, out of which 90% follow simple hyphenation rules, but the remaining 10% require expensive disk accesses."
Source: Bloom 1970.
URL: http://www.dragonwins.com/domains/getteched/bbc/literature/Bloom70.pdf
Confidence: HIGH
Corroborated By: Source 1 (Wikipedia summary).

## Evidence 11
Claim: Counting Bloom filters and Ripple filters are extensions that allow deletions while preserving low false‑positive rates.
Evidence: Source 1 (Wikipedia) lists Counting Bloom filters, Scalable Bloom filters, and Ripple filters as extensions that address deletions and dynamic size growth.
Source: Wikipedia.
URL: https://en.wikipedia.org/wiki/Bloom_filter
Confidence: MEDIUM
Corroborated By: Source 1 only (no independent academic source found yet).

## Evidence 12
Claim: The Rigorous upper bound for finite Bloom filters (Goel & Gupta, 2007) proves the standard approximation is within a small factor.
Evidence: Source 11 (Wikipedia) cites: "Goel and Gupta, however, give a rigorous upper bound that makes no approximations … ε ≤ (1 − e^(−k(n+0.5)/(m−1)))^k."
Source: Wikipedia citing Goel & Gupta 2007.
URL: https://en.wikipedia.org/wiki/Bloom_filter#Probability_of_false_positives
Confidence: HIGH
Corroborated By: Source 11 only (primary source would be the Goel & Gupta paper, not directly accessed here).
