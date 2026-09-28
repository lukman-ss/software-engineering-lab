# Audit Plan: Bloom Filters Research

## Target Lab
`labs/39-bloom-filters`

## Files Reviewed
- `labs/39-bloom-filters/research/01-plan.md`
- `labs/39-bloom-filters/research/02-sources.md`
- `labs/39-bloom-filters/research/03-evidence.md`
- `labs/39-bloom-filters/research/04-contradictions.md`
- `labs/39-bloom-filters/research/05-report.md`
- `labs/39-bloom-filters/research/06-open-questions.md`

## Claims To Verify
1. Bloom filters use *m* bits and *k* hash functions; guarantee zero false negatives with possible false positives.
2. False-positive probability approximation: $\epsilon \approx (1 - e^{-kn/m})^k$.
3. Optimal number of hash functions: $k = (m/n) \ln 2$, leading to $\epsilon \approx 0.618^{m/n}$.
4. Space requirement: $m/n \approx -1.44 \log_2 \epsilon$; ~10 bits/element for ~1% false positive rate.
5. LSM-tree SST lookup disk I/O reduction and lookup cost formula $O(L \cdot e^{-M/N})$.
6. Use in Percolator (Google) and BitFunnel (Bing).
7. Hash function suitability: non-cryptographic MurmurHash3 and FNV-1a.
8. Deletion handling limitations and alternatives (Counting Bloom filters, Cuckoo filters).

## Code To Execute
None. Per pipeline override, this is a research-only audit stage. No code exists or is audited in this stage.

## Primary Risks
- Over-reliance on Wikipedia as primary evidence for theoretical and empirical claims.
- Broken or dead URLs (e.g. academic homepage moves).
- Misattribution or over-simplification of distributed systems architectures (e.g. Percolator vs Bigtable Bloom filters).
- Misleading naming or claims regarding extensions (e.g., "Ripple filter" vs Ribbon filter).

## Audit Strategy
1. Perform HTTP reachability and content verification on all 12 cited sources in `02-sources.md`.
2. Cross-check all 12 claims in `03-evidence.md` and 8 findings in `05-report.md` against original publications.
3. Validate mathematical formulas and numeric thresholds.
4. Record contradictions, unsupported statements, or weak citations.
5. Formulate audit gap report and final verdict.
