# Claim Audit

## Claim 1
Claim: A Bloom filter uses a bit array of size *m* and *k* independent hash functions to store membership of up to *n* elements; it never yields a false negative but may produce false positives.

Location:
`research/03-evidence.md`: Evidence 1; `research/05-report.md`: Finding 1

Evidence Provided:
Wikipedia quote on definition and Bloom 1970 paper.

Source:
Source 1 (Wikipedia), Source 2 (Bloom 1970).

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Standard universally verified mathematical property.

---

## Claim 2
Claim: The approximate false-positive probability is $\epsilon \approx (1 - e^{-kn/m})^k$.

Location:
`research/03-evidence.md`: Evidence 2; `research/05-report.md`: Finding 2

Evidence Provided:
Wikipedia derivation and Bloom 1970.

Source:
Source 11 (Wikipedia), Source 2 (Bloom 1970).

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Approximation assumes hash function independence and $(1 - 1/m)^{kn} \approx e^{-kn/m}$.

---

## Claim 3
Claim: Optimal $k = (m/n) \ln 2$, giving minimum $\epsilon \approx 0.618^{m/n}$.

Location:
`research/03-evidence.md`: Evidence 3; `research/05-report.md`: Finding 2

Evidence Provided:
Wikipedia section and calculus derivation.

Source:
Source 11 (Wikipedia), Source 2 (Bloom 1970).

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
$(1/2)^{\ln 2 \cdot (m/n)} = (0.5^{\ln 2})^{m/n} \approx 0.6185^{m/n}$. Validated.

---

## Claim 4
Claim: For a desired false-positive rate $\epsilon$, the minimum bits per element is $m/n \approx -1.44 \log_2 \epsilon$.

Location:
`research/03-evidence.md`: Evidence 4; `research/05-report.md`: Executive Summary

Evidence Provided:
Mathematical simplification $m/n = -\ln \epsilon / (\ln 2)^2 \approx 1.4427 \log_2(1/\epsilon)$.

Source:
Source 11 (Wikipedia).

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
$1 / (\ln 2)^2 \approx 2.0813$ in natural log base; converted to $\log_2$, $1 / \ln 2 \approx 1.4427$. Formula is verified.

---

## Claim 5
Claim: Bloom filters are used per-SST in LSM-trees to avoid unnecessary disk reads; lookup cost drops from $O(L)$ to $O(L \cdot e^{-M/N})$.

Location:
`research/03-evidence.md`: Evidence 5; `research/05-report.md`: Finding 4

Evidence Provided:
Stopford (2015), O'Neil et al. (1996), Luo & Carey (2019).

Source:
Source 3, Source 9, Source 10.

Source Actually Supports Claim:
YES

Classification:
FACT / IMPLEMENTATION-SPECIFIC

Severity:
LOW

Notes:
Standard architectural pattern in RocksDB, LevelDB, Cassandra, HBase.

---

## Claim 6
Claim: The original Bloom-filter paper demonstrates that a 1% false-positive rate can be achieved with roughly 10 bits per element.

Location:
`research/03-evidence.md`: Evidence 6; `research/05-report.md`: Finding 3

Evidence Provided:
Direct quote from Bloom 1970: "Fewer than 10 bits per element are required for a 1% false positive probability...".

Source:
Source 2 (Bloom 1970).

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
$m/n = 1.4427 \cdot \log_2(100) \approx 9.585$ bits. Verified.

---

## Claim 7
Claim: Cuckoo filters can support deletions while using similar or less space than Bloom filters, and maintain comparable false-positive rates.

Location:
`research/03-evidence.md`: Evidence 7; `research/05-report.md`: Finding 8

Evidence Provided:
Fan et al. 2014 slides / presentation.

Source:
Source 8 (Fan et al., 2014).

Source Actually Supports Claim:
PARTIAL

Classification:
INTERPRETATION

Severity:
MEDIUM

Notes:
The URL provided for Source 8 returned 404 Not Found. While the claim itself is widely supported in peer-reviewed literature (Fan et al., ACM CoNEXT 2014), the provided citation link in the research document is dead.

---

## Claim 8
Claim: Non-cryptographic hash functions (e.g., MurmurHash3, FNV-1a) are suitable for Bloom filters because speed matters more than cryptographic strength.

Location:
`research/03-evidence.md`: Evidence 8; `research/05-report.md`: Finding 7

Evidence Provided:
SmHasher benchmarks (2.5–5 GB/s) and FNV avalanche documentation.

Source:
Source 6, Source 7.

Source Actually Supports Claim:
YES

Classification:
FACT / BEST PRACTICE

Severity:
LOW

Notes:
Well-established in data structure literature.

---

## Claim 9
Claim: Google's Percolator system uses Bloom-filter-like structures to prevent cache penetration.

Location:
`research/03-evidence.md`: Evidence 9; `research/05-report.md`: Finding 5

Evidence Provided:
Peng & Dabek (2010).

Source:
Source 4 (Percolator).

Source Actually Supports Claim:
PARTIAL

Classification:
INTERPRETATION / IMPLEMENTATION-SPECIFIC

Severity:
MEDIUM

Notes:
Percolator is built on top of Bigtable. Bigtable itself uses Bloom filters on SSTables to avoid disk seeks for row/column lookups. Percolator does not implement an application-layer cache penetration filter; attributing this feature specifically to Percolator rather than underlying Bigtable/LSM architecture is imprecise.

---

## Claim 10
Claim: Bloom filters were first proposed for hyphenation dictionary lookups.

Location:
`research/03-evidence.md`: Evidence 10

Evidence Provided:
Bloom 1970 paper text on 500,000 word dictionary.

Source:
Source 2 (Bloom 1970).

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Historically accurate motivation described by Burton Bloom.

---

## Claim 11
Claim: Counting Bloom filters and "Ripple filters" are extensions that allow deletions.

Location:
`research/03-evidence.md`: Evidence 11; `research/05-report.md`: Executive Summary / Finding 8

Evidence Provided:
Wikipedia reference.

Source:
Source 1 (Wikipedia).

Source Actually Supports Claim:
PARTIAL

Classification:
INTERPRETATION

Severity:
MEDIUM

Notes:
Typo or misnaming: the research refers to "Ripple filters". The standard literature and Facebook/RocksDB technology is "Ribbon filter" (based on coupling Bloom/cuckoo ideas with linear systems solving), not "Ripple filter". Ribbon filters optimize space but are not primarily designed for deletions. Counting Bloom filters and Cuckoo filters support deletions.

---

## Claim 12
Claim: Rigorous upper bound for finite Bloom filters (Goel & Gupta, 2007) proves the standard approximation is within a small factor.

Location:
`research/03-evidence.md`: Evidence 12; `research/05-report.md`: Limitations

Evidence Provided:
Wikipedia section citing Goel & Gupta 2007.

Source:
Source 11 (Wikipedia).

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Acknowledged as relying on secondary Wikipedia quote in limitations.
