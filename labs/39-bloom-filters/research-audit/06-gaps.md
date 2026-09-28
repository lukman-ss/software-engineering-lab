# Research Gaps

## Gap 1

Type:
MISSING_SOURCE

Severity:
MEDIUM

Location:
`research/02-sources.md`: Source 8 (Cuckoo filter)

Problem:
The canonical URL for the Cuckoo filter paper (Fan et al., 2014) returns HTTP 404. Content of the paper was not verified.

Required Revision:
Replace URL with the ACM DL canonical link:
https://dl.acm.org/doi/10.1145/2674005.2674994
or public paper: https://www.cs.cmu.edu/~dga/papers/cuckoo-conext2014.pdf

Can Be Approved Without Fix:
NO

---

## Gap 2

Type:
WEAK_SOURCE / OVERGENERALIZATION

Severity:
MEDIUM

Location:
`research/03-evidence.md`: Evidence 9; `research/05-report.md`: Finding 5

Problem:
The claim that Google Percolator uses Bloom filters to prevent "cache penetration" is attributed directly to Percolator. However, the Bloom filter functionality resides in the underlying Bigtable storage layer, not in the Percolator application layer. Cache penetration is also not a primary framing of the Percolator paper.

Required Revision:
Clarify that the Bloom filter usage is at the Bigtable/SSTable level (which Percolator is built upon), not at the Percolator application level.

Can Be Approved Without Fix:
YES (doesn't invalidate the lab concept, but misattributes architectural responsibility)

---

## Gap 3

Type:
UNVERIFIED_CLAIM

Severity:
MEDIUM

Location:
`research/03-evidence.md`: Evidence 11; `research/05-report.md`: Executive Summary

Problem:
The research names "Ripple filter" as a deletion-supporting Bloom filter extension. No such structure by that name exists in mainstream academic literature. The likely intended reference is "Ribbon filter" (Dillinger & Walzer, 2022), which is a space-efficient replacement, not primarily deletion-based.

Required Revision:
Remove the name "Ripple filter". Replace with accurate alternatives:
- Counting Bloom filter (supports deletion by using counters instead of single bits).
- Cuckoo filter (supports deletion, better FP performance).
- (Optionally) Ribbon filter as a newer space-efficient structure that does not support deletions.

Can Be Approved Without Fix:
NO

---

## Gap 4

Type:
WEAK_SOURCE

Severity:
LOW

Location:
`research/02-sources.md`: Source 5 (BitFunnel)

Problem:
The BitFunnel source cites Wikipedia rather than the original SIGIR 2017 paper: "BitFunnel: Revisiting Bit-sliced Signatures for Search Engines" by Goodwin et al. Wikipedia satisfactorily confirms the claim, but a stronger citation would be the primary paper.

Required Revision:
Consider adding the SIGIR 2017 paper DOI as the primary citation alongside the Wikipedia entry.

Can Be Approved Without Fix:
YES

---

## Gap 5

Type:
MISSING_SOURCE

Severity:
MEDIUM

Location:
`research/03-evidence.md`: Evidence 12

Problem:
Goel & Gupta (2007) is cited via Wikipedia only. The primary paper "Towards Tighter Space Bounds for Counting Bloom Filters" (Goel & Gupta, 2007) was not retrieved or independently verified.

Required Revision:
Either verify via ACM DL or arxiv, or downgrade the confidence label from HIGH to MEDIUM and note the source is secondary.

Can Be Approved Without Fix:
YES (the formula is widely corroborated; absence does not invalidate the primary formulas)

---

## Gap 6

Type:
MISSING_CASE

Severity:
LOW

Location:
`research/06-open-questions.md`

Problem:
Open questions mention PostgreSQL and MySQL built-in Bloom filter support, but the research does not acknowledge that PostgreSQL 9.6+ does include a bloom filter index access method extension. This is a known, documented feature not covered by the research.

Required Revision:
Optional: add PostgreSQL bloom filter index documentation as an additional source.

Can Be Approved Without Fix:
YES

---

## Gap 7

Type:
WEAK_SOURCE

Severity:
LOW

Location:
`research/02-sources.md`: Sources 1, 7, 11

Problem:
Multiple claims that should be traced back to mathematical derivations or primary algorithm papers (e.g., Bloom 1970) are instead supported solely by Wikipedia citations. This is acceptable for corroboration but should not be the sole basis for theoretical claims.

Required Revision:
Ensure primary references accompany Wikipedia for all theoretical claims.

Can Be Approved Without Fix:
YES (primary source Bloom 1970 is separately available and accessible)
