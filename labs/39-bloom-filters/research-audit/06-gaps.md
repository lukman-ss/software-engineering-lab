# 06 - Research Gaps: Bloom Filters Research

## Gap 1

Type: WEAK_SOURCE

Severity: MEDIUM

Location: `research/03-evidence.md` (Evidence 4); `research/05-report.md` (Finding 3)

Problem:
"Cache Penetration Mitigation" di-source ke RocksDB Wiki dan Martin Kleppmann (DDIA). Kleppmann (DDIA) tidak dikutip secara formal dalam `02-sources.md` (tidak ada Source 7). Referensi ke buku DDIA disebutkan di `03-evidence.md` tetapi tidak masuk dalam daftar sumber resmi.

Required Revision:
Tambahkan DDIA (Martin Kleppmann, O'Reilly Media) sebagai Source resmi, termasuk nomor halaman/chapter yang relevan (Chapter 3, Bab Encoding dan Evolution pada Storage Engines atau BAB 3 LSM-Trees & SSTables).

Can Be Approved Without Fix:
YES

---

## Gap 2

Type: MISSING_SOURCE

Severity: MEDIUM

Location: `research/02-sources.md` (referensi corroborating di Evidence 2)

Problem:
"Broder & Mitzenmacher (2004) 'Network Applications of Bloom Filters'" disebutkan sebagai koroborasi di `03-evidence.md` (Evidence 2), tetapi tidak terdaftar sebagai sumber resmi di `02-sources.md`. Survey Broder & Mitzenmacher (Internet Mathematics, 2004) adalah sumber sekunder yang sering dikutip untuk formula matematis.

Required Revision:
Daftarkan "Network Applications of Bloom Filters: A Survey" (Andrei Broder & Michael Mitzenmacher, Internet Mathematics 1(4), 2004) sebagai Source resmi, atau hapus referensi corroborating dari `03-evidence.md`.

Can Be Approved Without Fix:
YES

---

## Gap 3

Type: MISSING_SOURCE

Severity: LOW

Location: `research/06-open-questions.md` §1.2

Problem:
Scalable Bloom Filter (Almeida et al. 2007) dikutip sebagai referensi Open Question tetapi tidak terdaftar sebagai sumber resmi. Karena ini berada di Open Questions (bukan klaim utama), dampak minimal.

Required Revision:
Tambahkan catatan kaki atau referensi pendukung untuk Almeida et al. (2007).

Can Be Approved Without Fix:
YES

---

## Gap 4

Type: WEAK_SOURCE

Severity: LOW

Location: `research/02-sources.md` (Source 1)

Problem:
URL Source 1 (`https://dl.acm.org/doi/10.1145/362686.362692`) tidak dapat diakses melalui automated agent (HTTP 403 Cloudflare protection). Ini bukan link rot — DOI valid dan makalah dapat diakses melalui browser browser normal. Tidak ada indikasi konten palsu, namun URL tidak dapat diverifikasi secara programatik.

Required Revision:
Tambahkan URL alternatif atau PDF mirror (e.g., via Semantic Scholar atau ResearchGate) jika akses harus terverifikasi programatik. Ini opsional.

Can Be Approved Without Fix:
YES

---

## Gap 5

Type: MISSING_SOURCE (URL Broken)

Severity: MEDIUM

Location: `research/02-sources.md` (Source 2)

Problem:
URL Source 2 (`https://www.eecs.harvard.edu/~michaelm/postscripts/esa2006.pdf`) mengembalikan HTTP 404. Canonical technical report terkait (`tr-02-05.pdf`) tersedia di URL berbeda pada host yang sama. Klaim yang di-source dari Kirsch & Mitzenmacher valid (makalah nyata, bukti formal valid), namun URL yang tercantum tidak dapat diakses.

Required Revision:
Update URL Source 2 ke: `https://www.eecs.harvard.edu/~michaelm/postscripts/tr-02-05.pdf` (HTTP 200, verified) atau ke Springer LNCS 4168 DOI.

Can Be Approved Without Fix:
YES

---

## Gap 6

Type: MISSING_SOURCE (URL Broken)

Severity: MEDIUM

Location: `research/02-sources.md` (Source 5)

Problem:
URL Source 5 (`https://cassandra.apache.org/doc/latest/cassandra/operating/bloom_filters.html`) mengembalikan HTTP 404. Apache Cassandra docs mengalami restrukturisasi URL. Konsep `bloom_filter_fp_chance` dan panduan tuning yang diklaim valid dalam Cassandra, tetapi URL spesifik perlu diperbarui.

Required Revision:
Update URL Source 5 ke halaman Cassandra docs yang valid dengan konten Bloom filter (perlu verifikasi URL baru dari situs Apache Cassandra langsung).

Can Be Approved Without Fix:
YES

---

## Gap 7

Type: UNVERIFIED_CLAIM

Severity: MEDIUM

Location: `research/05-report.md` (Executive Summary)

Problem:
Klaim "~50-100 MB jika menggunakan Hash Set biasa" untuk 1.000.000 elemen tidak memiliki sumber empiris yang dikutip. Ini adalah estimasi heuristik yang tidak dapat diverifikasi dari sumber terdaftar.

Required Revision:
Tambahkan bukti atau perhitungan eksplisit (rata-rata 50-100 bytes overhead per entry di Java HashMap / Go map), atau ubah menjadi range ilustratif dengan catatan yang jelas sebagai approximation.

Can Be Approved Without Fix:
YES

---

## Gap 8

Type: UNVERIFIED_CLAIM

Severity: LOW

Location: `research/06-open-questions.md` §2.1

Problem:
Penelitian sendiri mengakui bahwa klaim perbandingan empiris MurmurHash3 vs xxHash vs FNV "belum di-cross-check dengan benchmark primer terpercaya." Ini adalah gap yang sudah diidentifikasi oleh peneliti sendiri.

Required Revision:
Open question bisa diterima sebagai area future work.

Can Be Approved Without Fix:
YES
