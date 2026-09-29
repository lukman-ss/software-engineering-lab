# Claim Audit: Bloom Filters Research

## Claim 1
Claim: Bloom Filter tidak pernah menghasilkan False Negative (garansi 100% akurat tanpa false negative).  
Location: `05-report.md` (Executive Summary, Finding 1) & `03-evidence.md` (Evidence 1)  
Evidence Provided: Sifat dasar pengisian bit array. Elemen yang telah di-insert menjamin seluruh $k$ bit bernilai 1.  
Source: Burton H. Bloom (1970)  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: CRITICAL (Key fundamental property of Bloom Filters)  
Notes: Terverifikasi.  

---

## Claim 2
Claim: Penentuan $m$ dan $k$ optimal mengikuti rumus $m = - \frac{n \ln p}{(\ln 2)^2} \approx -1.4427 \cdot n \log_2 p$ dan $k = \frac{m}{n} \ln 2 \approx 0.6931 \cdot \frac{m}{n}$. Untuk $p = 0.01$, alokasi adalah $\sim 9.6$ bit/elemen dan $k = 7$ (konsumsi $\sim 1.2$ MB RAM per 1.000.000 elemen).  
Location: `05-report.md` (Executive Summary, Finding 1) & `03-evidence.md` (Evidence 2)  
Evidence Provided: Derivasi probabilitas $p \approx (1 - e^{-kn/m})^k$.  
Source: Bloom (1970), Kirsch & Mitzenmacher (2006)  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: HIGH  
Notes: Turunan matematis tepat secara eksak. Notasi $-1.4427 \cdot n \log_2 p$ menggunakan basis $\log_2$, di mana $\frac{1}{\ln 2} \approx 1.4427$, sehingga $- \frac{\ln p}{(\ln 2)^2} = - \frac{\log_2 p}{\ln 2} = - 1.4427 \log_2 p$. Perhitungan RAM per 1.000.000 elemen ($9.585 \times 10^6$ bit $/ 8 / 1024^2 \approx 1.143$ MB $\approx 1.2$ MB) akurat.  

---

## Claim 3
Claim: Kombinasi dua fungsi hash $g_i(x) = h_1(x) + i \cdot h_2(x) \pmod m$ mensimulasikan $k$ fungsi hash independen tanpa mengorbankan false positive rate secara asimtotik (Kirsch-Mitzenmacher Double Hashing).  
Location: `05-report.md` (Finding 2) & `03-evidence.md` (Evidence 3)  
Evidence Provided: Teorema Kirsch-Mitzenmacher 2006.  
Source: Kirsch & Mitzenmacher (ESA 2006)  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: HIGH  
Notes: Diadopsi secara luas pada mesin industri seperti RocksDB dan Google Guava.  

---

## Claim 4
Claim: Bloom Filter memotong 99% query nonexistent sebelum mencapai Cache atau Database, mencegah Cache Penetration.  
Location: `05-report.md` (Finding 3) & `03-evidence.md` (Evidence 4)  
Evidence Provided: Analisis arsitektural filter di depan query path database.  
Source: Martin Kleppmann (DDIA); RocksDB Wiki  
Source Actually Supports Claim: YES  
Classification: INTERPRETATION  
Severity: MEDIUM  
Notes: Angka 99% berlaku spesifik untuk konfigurasi target false positive rate $p = 0.01$ ($1 - 0.01 = 0.99$). Secara konsep tepat dan terdukung.  

---

## Claim 5
Claim: Storage engine berbasis LSM-Tree (RocksDB, Cassandra, Google Bigtable) mengandalkan Bloom Filter per-SSTable untuk menghindari disk seek pada point query.  
Location: `05-report.md` (Finding 4) & `03-evidence.md` (Evidence 5)  
Evidence Provided: Pengutipan langsung dari paper Bigtable (Chang et al., 2006) dan arsitektur RocksDB/Cassandra.  
Source: Chang et al. (2006); RocksDB Wiki; Apache Cassandra Docs  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: HIGH  
Notes: Terverifikasi komprehensif.  

---

## Claim 6
Claim: Standard Bloom Filter tidak mendukung operasi deletion (`delete`), karena pengosongan bit dapat menyebabkan false negative pada elemen lain.  
Location: `05-report.md` (Areas of Disagreement) & `03-evidence.md` (Evidence 6)  
Evidence Provided: Karakteristik bit sharing akibat hash collision.  
Source: Fan et al. (2014); Broder & Mitzenmacher (2004)  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: HIGH  
Notes: Sesuai dengan spesifikasi dasar Bloom Filter.  
