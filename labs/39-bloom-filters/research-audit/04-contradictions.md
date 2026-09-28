# Contradictions

## Contradiction 1

Statement A:
`research/05-report.md` Finding 6 sources BitFunnel confirmation with "Source 4 (Google Percolator analogy)" alongside Source 5 (Wikipedia BitFunnel).

Location:
`research/05-report.md`, Finding 6

Statement B:
`research/02-sources.md` Source 4 is the Percolator paper, which covers Google's incremental indexing, not Bing's BitFunnel architecture. Percolator and BitFunnel are unrelated systems.

Location:
`research/02-sources.md`, Source 4

Type:
INTERNAL

Impact:
The parenthetical "(Google Percolator analogy)" implies Percolator supports BitFunnel claims. These are independent systems. Cross-citing Percolator as corroboration for BitFunnel is logically unsound.

Assessment:
The error is in cross-reference justification. Does not invalidate the BitFunnel claim itself (supported by Source 5), but weakens citation traceability.

---

## Contradiction 2

Statement A:
`research/05-report.md` Executive Summary and `research/03-evidence.md` Evidence 11 list "Ripple filters" as a deletion-supporting extension.

Location:
`research/05-report.md`, Executive Summary

Statement B:
The actual extension in the current literature is "Ribbon filters" (not "Ripple"). Ribbon filters originate from the 2022 paper "Ribbon filter: practically smaller than Bloom and Xor" by Peter C. Dillinger & Stefan Walzer. Ribbon filters trade space for construction complexity and are not primarily deletion-supporting.

Location:
`research/03-evidence.md`, Evidence 11

Type:
INTERNAL / CODE_DOC_MISMATCH

Impact:
Incorrect naming may confuse lab participants looking for references to "Ripple filters" in the literature.

Assessment:
The claim about deletion support for Counting Bloom and Cuckoo filters is correct. The "Ripple filter" name is erroneous. The research uses a non-existent name, but corroborating deletion evidence does not depend on it.

---

No other material contradictions detected. All sources agree on:
- Core algorithm: bit array + k hash functions.
- No false negatives guarantee.
- Standard false-positive formula.
- Optimal k formula.
- LSM-tree per-SST Bloom filter use.
