
# Contradictions

No material contradictions discovered among the consulted sources.

All examined references agree on:
- The fundamental definition of a Bloom filter (probabilistic, no false negatives).
- The approximate false‑positive formula ε ≈ (1−e^(−kn/m))^k.
- The optimal k = (m/n)·ln 2 and corresponding bits‑per‑element m/n ≈ 1.44·log₂(1/ε).
- Practical usage in LSM‑tree databases, caching layers, and search indexes.

Minor differences:
- Source 8 (Cuckoo filter presentation) vs. Source 1 (Wikipedia): Cuckoo filters claim lower FP rates in some scenarios, while Wikipedia lists them as an alternative without asserting superiority. This is not a contradiction but a matter of emphasis.
- Source 10 (LSM survey) vs. Source 3 (original LSM paper): Survey reports O(L·e^(−M/N)) lookup cost; original paper gives the same formula but with different notation. No conflict.
