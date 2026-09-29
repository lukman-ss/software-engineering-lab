# Research Report: Bloom Filters

## Research Question
Bagaimana struktur data probabilistik Bloom Filter memangkas disk I/O, mencegah fenomena Cache Penetration, dan bagaimana implementasi matematis serta arsitektural yang optimal untuk sistem backend berkinerja tinggi?

## Executive Summary
Bloom Filter adalah struktur data probabilistik berbasis bit array yang dirancang oleh Burton H. Bloom pada tahun 1970 untuk memverifikasi keanggotaan himpunan (*set membership*) dengan penggunaan ruang memori minimal. Bloom Filter memiliki garansi asimetris:
1. **Pasti Tidak Ada (Definitely Not in Set)**: 100% akurat tanpa kemungkinan false negative.
2. **Mungkin Ada (Possibly in Set)**: Memiliki probabilitas kecil false positive ($p$), di mana elemen yang tidak ada terdeteksi seolah-olah ada karena tabrakan bit (*hash collision*).

Pada sistem skala besar, karakteristik ini dimanfaatkan untuk mengeliminasi disk lookup yang mahal pada database (seperti LSM-Tree SSTables di RocksDB dan Cassandra) dan mencegah *Cache Penetration* (situasi di mana attacker meminta kunci nonexistent secara masif, memaksa query langsung ke database). Dengan alokasi ~9.6 bit per elemen dan 7 fungsi hash, sistem dapat mencapai 99% akurasi penolakan query nonexistent hanya dengan konsumsi RAM ~1.2 MB per 1.000.000 elemen (dibandingkan ~50-100 MB jika menggunakan Hash Set konvensional dalam memori, di mana overhead pointer, bucket table, dan string key object umumnya mengonsumsi rata-rata 48-96 byte per entri di runtime seperti Go map atau Java HashSet).

---

## Findings

### Finding 1: Asimetri Kesalahan dan Rumus Optimal
- **Claim**: Bloom Filter tidak pernah menghasilkan False Negative. Ukuran bit array optimal ($m$) dan jumlah fungsi hash ($k$) dapat dihitung secara eksak dari target False Positive Rate ($p$) dan jumlah elemen yang diharapkan ($n$).
- **Evidence**:
  - Probabilitas bit bernilai 0 setelah memasukkan $n$ elemen: $p_0 \approx e^{-kn/m}$.
  - Probabilitas false positive: $p \approx (1 - e^{-kn/m})^k$.
  - Ukuran bit array optimal: $m = - \frac{n \ln p}{(\ln 2)^2} \approx -1.4427 \cdot n \log_2 p$.
  - Jumlah hash optimal: $k = \frac{m}{n} \ln 2 \approx 0.6931 \cdot \frac{m}{n}$.
- **Sources**: Burton H. Bloom (1970); Kirsch & Mitzenmacher (2006).
- **Confidence**: HIGH.

### Finding 2: Optimasi Komputasi dengan Kirsch-Mitzenmacher Double Hashing
- **Claim**: Menjalankan $k$ algoritma hash berbeda menimbulkan overhead CPU tinggi. Menggunakan teknik Kirsch-Mitzenmacher $g_i(x) = h_1(x) + i \cdot h_2(x) \pmod m$ mensimulasikan $k$ fungsi hash independen dengan kinerja asimtotik yang identik.
- **Evidence**: Membagi hash 128-bit (misalnya dari MurmurHash3 atau xxHash) menjadi dua 64-bit unsigned integer ($h_1$ dan $h_2$) memungkinkan pembuatan $k$ indeks bit secara instan di CPU register tanpa perlu memanggil fungsi hashing berulang kali.
- **Sources**: Kirsch & Mitzenmacher (ESA 2006); Google Guava Library; RocksDB.
- **Confidence**: HIGH.

### Finding 3: Pencegahan Cache Penetration
- **Claim**: Bloom Filter memotong 99% query nonexistent sebelum mencapai Cache atau Database.
- **Evidence**: Dalam skenario penyerangan atau traffic acak kunci kosong, arsitektur backend tradisional mengalami Cache Miss dan langsung mengeksekusi Disk I/O ke Database. Dengan meletakkan Bloom Filter di depan jalur pembacaan, request yang menghasilkan nilai 0 pada bit array langsung dihentikan dan me-return response 404 / empty tanpa membebani thread pool database.
- **Sources**: Martin Kleppmann (Designing Data-Intensive Applications); RocksDB Wiki.
- **Confidence**: HIGH.

### Finding 4: Disk I/O Pruning pada LSM-Tree
- **Claim**: Mesin basis data berbasis Log-Structured Merge-Tree (LSM-Tree) mengandalkan Bloom Filter per-SSTable untuk menghindari disk seek pada point query.
- **Evidence**: Karena data baru pada LSM-Tree ditulis secara append-only dan terdistribusi di berbagai file SSTable berjenjang, pencarian kunci tanpa filter mengharuskan pembacaan setiap SSTable dari memori/disk. Bloom Filter yang disimpan di header SSTable atau RAM memvalidasi apakah file tersebut memuat kunci target. Jika negatif, pembacaan file SSTable dilewati sepenuhnya.
- **Sources**: Google Bigtable Paper (Chang et al., OSDI 2006); Apache Cassandra Docs; RocksDB Wiki.
- **Confidence**: HIGH.

---

## Areas of Agreement
- **Zero False Negative**: Seluruh sumber akademis dan dokumentasi sistem sepakat bahwa elemen yang telah di-insert tidak akan pernah menghasilkan status "Not in Set".
- **Space Efficiency**: Kebutuhan memori konstan per elemen untuk target error rate tertentu, terlepas dari ukuran data string asli (karena hanya nilai hash yang di-petakan ke bit array).
- **Double Hashing Equivalence**: Penerimaan luas terhadap teknik Kirsch-Mitzenmacher dalam implementasi sistem produksi.

## Areas of Disagreement & Trade-offs
- **Dukungan Operasi Delete**: Standard Bloom Filter tidak mendukung penghapusan elemen. Sistem yang memerlukan pembaruan dinamis dengan penghapusan harus memilih antara:
  - Me-rebuild filter secara berkala (pendekatan umum pada SSTable immutable).
  - Menggunakan Counting Bloom Filter (biaya memori meningkat 3-4x per counter 4-bit).
  - Menggunakan Cuckoo Filter (mendukung deletion dan lookup locality, tetapi proses insertion berpotensi gagal jika kapasitas penuh).

## Limitations
1. **Kapasitas Statis**: Begitu jumlah elemen aktual melebihi estimasi $n$, rasio bit 1 akan mendominasi dan false positive rate ($p$) melonjak drastis mendekati 100%.
2. **Tidak Menyimpan Nilai Asli**: Bloom Filter hanya menjawab pertanyaan keanggotaan himpunan (*membership query*), bukan media penyimpanan key-value.

## Conclusion
Bloom Filter memberikan efisiensi ruang dan waktu ekstrem untuk sistem backend berskala tinggi. Penerapan formula optimal $m$ dan $k$ serta double hashing MurmurHash3/xxHash memberikan pertahanan efektif terhadap cache penetration dan disk I/O bottleneck pada database storage engine.
