# 03 - Claim Audit: Bloom Filters Research

## Claim 1: Zero False Negative Guarantee
- **Claim**: Bloom filter menjamin tidak ada false negative (0% false negative rate), namun dapat menghasilkan false positive ($p > 0$).
- **Location**: `03-evidence.md` (Evidence 1); `05-report.md` (Finding 1).
- **Evidence Provided**: Setting bit 1 saat insert menjamin pengecekan $k$ hash posisi selalu mengembalikan 1 jika elemen ada.
- **Source**: Burton H. Bloom (1970).
- **Source Actually Supports Claim**: YES.
- **Classification**: FACT.
- **Severity**: LOW.
- **Notes**: Asas fundamental struktur data Bloom Filter.

---

## Claim 2: Mathematical Formulae for Optimal Bit Size and Hash Count
- **Claim**:
  $$m = - \frac{n \ln p}{(\ln 2)^2} \approx -1.4427 \cdot n \log_2 p$$
  $$k = \frac{m}{n} \ln 2 \approx 0.6931 \cdot \frac{m}{n}$$
  Untuk $p = 0.01$ (1%), rasio $m/n \approx 9.585$ bit/elemen dan $k \approx 7$.
- **Location**: `03-evidence.md` (Evidence 2); `05-report.md` (Finding 1).
- **Evidence Provided**: Derivasi probabilitas false positive $p \approx (1 - e^{-kn/m})^k$ dengan minimasi terhadap $k$.
- **Source**: Burton H. Bloom (1970); Kirsch & Mitzenmacher (2006).
- **Source Actually Supports Claim**: YES.
- **Classification**: FACT.
- **Severity**: LOW.
- **Notes**: Kalkulasi matematis diverifikasi akurat. $-1 / (\ln 2)^2 \approx -1 / (0.693147^2) \approx 2.08136 \cdot \ln(1/p) \approx 1.4427 \cdot \log_2(1/p)$. Untuk $p=0.01$, $- \ln(0.01) / (\ln 2)^2 \approx 4.60517 / 0.480453 \approx 9.585$ bits/elemen. $k = 9.585 \times \ln 2 \approx 6.64 \approx 7$.

---

## Claim 3: Kirsch-Mitzenmacher Double Hashing Optimization
- **Claim**: Kombinasi linear dua fungsi hash independen $g_i(x) = h_1(x) + i \cdot h_2(x) \pmod m$ cukup untuk mensimulasikan $k$ fungsi hash independen tanpa degradasi false positive rate asimtotik.
- **Location**: `03-evidence.md` (Evidence 3); `05-report.md` (Finding 2).
- **Evidence Provided**: Teorema Kirsch-Mitzenmacher (ESA 2006) membuktikan kesetaraan distribusi keanggotaan bit asimtotik.
- **Source**: Kirsch & Mitzenmacher (2006).
- **Source Actually Supports Claim**: YES.
- **Classification**: FACT.
- **Severity**: LOW.
- **Notes**: Standar de facto industri (Google Guava, RocksDB, Redis).

---

## Claim 4: Cache Penetration Mitigation
- **Claim**: Bloom filter mencegah masalah "Cache Penetration" dengan memfilter query untuk kunci nonexistent sebelum query menyentuh cache atau disk database, memotong hingga ~99% beban disk/database untuk nonexistent keys pada $p=0.01$.
- **Location**: `03-evidence.md` (Evidence 4); `05-report.md` (Finding 3).
- **Evidence Provided**: Kunci nonexistent yang tidak ada dalam filter dieliminasi 100% saat bit 0 ditemukan, hanya meloloskan fraksi $p$ false positive.
- **Source**: Martin Kleppmann (DDIA); RocksDB Wiki.
- **Source Actually Supports Claim**: YES.
- **Classification**: INTERPRETATION.
- **Severity**: LOW.
- **Notes**: Valid pattern in distributed backend architectures.

---

## Claim 5: LSM-Tree Disk I/O Pruning per SSTable
- **Claim**: Mesin basis data berbasis Log-Structured Merge-Tree (LSM-Tree) seperti Bigtable, RocksDB, dan Cassandra menggunakan Bloom Filter pada setiap SSTable untuk menghindari pembacaan disk yang tidak perlu.
- **Location**: `03-evidence.md` (Evidence 5); `05-report.md` (Finding 4).
- **Evidence Provided**: Bigtable paper (Chang et al., OSDI 2006, p. 7): "A Bloom filter allows us to ask whether an SSTable might contain any data for a specified row/column pair... Drastically reduces the number of disk seeks required for read operations."
- **Source**: Fay Chang et al. (2006).
- **Source Actually Supports Claim**: YES.
- **Classification**: FACT.
- **Severity**: LOW.
- **Notes**: Kutipan diverifikasi langsung pada dokumen orisinal Bigtable.

---

## Claim 6: Standard Bloom Filter Inability to Support Deletion
- **Claim**: Standard Bloom Filter tidak mendukung operasi penghapusan (`delete`). Menghapus bit 1 menjadi 0 dapat menyebabkan false negative pada elemen lain akibat hash collisions.
- **Location**: `03-evidence.md` (Evidence 6); `05-report.md` (Limitations & Disagreement).
- **Evidence Provided**: Cuckoo Filter paper (Fan et al., ACM CoNEXT 2014) Table 1 & Section 1.
- **Source**: Fan et al. (2014).
- **Source Actually Supports Claim**: YES.
- **Classification**: FACT.
- **Severity**: LOW.
- **Notes**: Properti intrinsik bit array terbagi.
