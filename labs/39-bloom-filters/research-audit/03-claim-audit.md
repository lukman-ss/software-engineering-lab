# Claim Audit: Bloom Filters Research

## Claim 1
Claim: A Bloom filter uses a bit array of size $m$ and $k$ independent hash functions; false negatives are impossible, but false positives may occur.
Location: `research/03-evidence.md:5-10`, `research/05-report.md:20-25`
Evidence Provided: Bloom 1970, Wikipedia Bloom filter section.
Source: Source 1, Source 2, Source 10.
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Fundamental invariant of standard Bloom filters.

## Claim 2
Claim: The approximate false-positive probability is $\varepsilon \approx (1 - e^{-kn/m})^k$.
Location: `research/03-evidence.md:12-18`, `research/05-report.md:26-30`
Evidence Provided: Standard mathematical derivation assuming independent bit setting probabilities.
Source: Source 2, Source 10, Source 11.
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Verified against analytical derivations and asymptotic approximations.

## Claim 3
Claim: Optimal number of hash functions is $k = (m/n) \ln 2$, giving minimum false positive probability $(1/2)^k \approx (0.6185)^{m/n}$.
Location: `research/03-evidence.md:20-26`, `research/05-report.md:13`
Evidence Provided: Calculus derivative of $(1 - e^{-kn/m})^k$ with respect to $k$.
Source: Source 2, Source 10, Source 11.
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Exact minimum $(1/2)^{(m/n)\ln 2} = 2^{-(m/n)\ln 2} = e^{-(m/n)(\ln 2)^2} \approx 0.6185^{m/n}$.

## Claim 4
Claim: For a target false positive rate $\varepsilon$, the required bit ratio is $m/n \approx -1.44 \log_2 \varepsilon$ (or $\sim 9.6$ bits/element for $1\%$ FP rate).
Location: `research/03-evidence.md:28-34`, `research/05-report.md:32-36`
Evidence Provided: Bloom 1970 calculation: $m/n = -\ln(\varepsilon)/(\ln 2)^2 \approx 1.4427 \log_2(1/\varepsilon)$.
Source: Source 2, Source 11.
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: For $\varepsilon = 0.01$, $m/n = 1.442695 \times \log_2(100) \approx 9.585$ bits/element with $k = \lceil 9.585 \ln 2 \rceil = 7$.

## Claim 5
Claim: LSM-tree storage engines use per-SST Bloom filters to reduce point-lookup disk read costs from $O(L)$ to $O(L \cdot e^{-M/N})$.
Location: `research/03-evidence.md:36-42`, `research/05-report.md:38-43`
Evidence Provided: O'Neil et al. (1996) LSM-tree paper, Luo & Carey (2019) survey.
Source: Source 3, Source 9, Source 10.
Source Actually Supports Claim: YES
Classification: FACT / IMPLEMENTATION-SPECIFIC
Severity: LOW
Notes: Universal across modern LSM engines (RocksDB, Cassandra, LevelDB).

## Claim 6
Claim: Google Percolator and Microsoft Bing (BitFunnel) utilize Bloom-filter principles to eliminate redundant reads or accelerate search queries.
Location: `research/03-evidence.md:68-74`, `research/05-report.md:44-54`
Evidence Provided: Peng & Dabek (2010), Goodwin et al. (2017).
Source: Source 4, Source 5.
Source Actually Supports Claim: YES
Classification: FACT / EXAMPLE
Severity: LOW
Notes: Properly contextualized in research notes.

## Claim 7
Claim: Fast non-cryptographic hash functions (MurmurHash3, FNV-1a) provide superior throughput over cryptographic hashes while maintaining low collision correlation for Bloom filters.
Location: `research/03-evidence.md:60-66`, `research/05-report.md:56-60`
Evidence Provided: SmHasher benchmarks (2.5–5 GB/s throughput) and FNV specification.
Source: Source 6, Source 7, Source 12.
Source Actually Supports Claim: YES
Classification: FACT / INTERPRETATION
Severity: LOW
Notes: Industry consensus standard for in-memory Bloom filter implementations.

## Claim 8
Claim: Standard Bloom filters do not support deletion; Counting Bloom filters and Cuckoo filters enable deletions with distinct trade-offs.
Location: `research/03-evidence.md:52-58, 84-90`, `research/05-report.md:62-67`
Evidence Provided: Fan et al. (CoNEXT 2014) Cuckoo filter analysis.
Source: Source 1, Source 8.
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Accurately contrasts structural limitations of classic Bloom bit-arrays against fingerprint/bucket structures.
