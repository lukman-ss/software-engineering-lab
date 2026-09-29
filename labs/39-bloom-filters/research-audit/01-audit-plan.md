# Audit Plan: Bloom Filters Research

## Target Lab
`labs/39-bloom-filters/research/`

## Files Reviewed
- `labs/39-bloom-filters/research/01-plan.md`
- `labs/39-bloom-filters/research/02-sources.md`
- `labs/39-bloom-filters/research/03-evidence.md`
- `labs/39-bloom-filters/research/04-contradictions.md`
- `labs/39-bloom-filters/research/05-report.md`
- `labs/39-bloom-filters/research/06-open-questions.md`

## Claims To Verify
1. **Zero False Negative Guarantee**: Elemen yang dimasukkan tidak pernah menghasilkan false negative (100% true negative guarantee jika bit bernilai 0).
2. **Optimal Bit Array Size and Hash Count Formulas**:
   - $m = - \frac{n \ln p}{(\ln 2)^2} \approx -1.4427 \cdot n \log_2 p$
   - $k = \frac{m}{n} \ln 2 \approx 0.6931 \cdot \frac{m}{n}$
   - Rasio untuk $p = 0.01 \implies \sim 9.6$ bits/elemen, $k = 7$.
3. **Kirsch-Mitzenmacher Double Hashing Optimization**:
   - Formula $g_i(x) = h_1(x) + i \cdot h_2(x) \pmod m$ mensimulasikan $k$ fungsi hash independen tanpa degradasi false positive rate asimtotik.
4. **Cache Penetration Mitigation**: Bloom Filter memangkas query nonexistent sebelum membebani disk/cache/database.
5. **LSM-Tree SSTable Disk I/O Pruning**: Penggunaan Bloom Filter pada Google Bigtable, RocksDB, dan Cassandra untuk melewati pembacaan SSTable yang tidak memuat kunci.
6. **Deletion Incompatibility**: Standard Bloom Filter tidak mendukung operasi deletion tanpa merusak elemen lain yang berbagi bit.

## Primary Risks
- Validitas URL dan kredibilitas sumber akademis/industri.
- Ketepatan perumusan matematika dan penulisan notasi logaritma/konstanta.
- Generalisasi berlebihan pada pencegahan cache penetration atau implementasi LSM storage engine.

## Audit Strategy
1. Audit kredibilitas dan jangkauan 8 sumber pada `02-sources.md`.
2. Validasi keselarasan klaim (`03-evidence.md` & `05-report.md`) dengan sumber primer.
3. Evaluasi analisis kontradiksi (`04-contradictions.md`) dan gap/pertanyaan terbuka (`06-open-questions.md`).
4. Berikan penilaian independen dan status kualitas riset.
