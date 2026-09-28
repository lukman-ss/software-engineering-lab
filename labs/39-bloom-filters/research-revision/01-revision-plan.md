# Revision Plan

Target Lab: labs/39-bloom-filters
Previous Audit Status: APPROVED_WITH_WARNINGS

## Blocking Issues

1. **Broken Cuckoo Filter URL**: Source 8 URL (`https://www.cs.cmu.edu/~fanzhao/cuckoo-filter.pdf`) returns 404. Replace with ACM DL canonical link (`https://dl.acm.org/doi/10.1145/2674005.2674994`) and public paper (`https://www.cs.cmu.edu/~dga/papers/cuckoo-conext2014.pdf`).
2. **Erroneous "Ripple filter" name**: References to "Ripple filter" in `03-evidence.md`, `04-contradictions.md`, and `05-report.md` do not exist in data structure literature. Replace with "Ribbon filter" (Dillinger & Walzer, 2022) with accurate attributes, or use standard Counting Bloom/Cuckoo filters for deletion support claim.

## Non-Blocking Issues

1. **Percolator vs Bigtable attribution**: Percolator's Bloom filtering is inherited from underlying Bigtable SSTable lookups rather than an application-level cache-penetration barrier. Clarifying in `03-evidence.md` and `05-report.md`.
2. **Wikipedia over-reliance for theoretical claims**: Add primary source (Bloom 1970) wherever possible, supplement Wikipedia citations.
3. **Cross-citation disconnect (BitFunnel/Percolator)**: Finding 6 mentions "(Google Percolator analogy)" for BitFunnel claims. Percolator and BitFunnel are independent systems. Remove cross-reference.
4. **BitFunnel primary source**: Cite SIGIR 2017 paper alongside Wikipedia in Source 5.
5. **Goel & Gupta (2007) primary source**: Add arxiv/DOI link if retrievable.

## Files To Modify

- `research/02-sources.md`
- `research/03-evidence.md`
- `research/04-contradictions.md`
- `research/05-report.md`

## Verification Plan

- source verification (HTTP reachability / URL check)
- documentation consistency (grep for "Ripple")
- final audit comparison
- re-read modified research files for internal consistency
