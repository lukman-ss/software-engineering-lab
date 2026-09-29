# Bloom Filter: Struktur Data Probabilistik untuk Pengurangan I/O dan Pertahanan terhadap Cache Penetration

## Masalah
Sistem backend modern sering menghadapi dua masalah performa terkait: (1) disk I/O yang berlebihan saat melakukan pencarian kunci yang tidak ada di database berbasis LSM‑Tree, dan (2) cache penetration yang terjadi ketika request untuk kunci tidak existing melintasi cache dan menganiaya backend. Masalah ini menyebabkan latensi tinggi, penggunaan resource yang berlebihan, dan kerentanan terhadap serangan sederhana yang memanfaatkan pola query acak.

## Mengapa Hal Ini Penting
Disk I/O adalah operasi yang paling mahal dalam basis data, sementara backend yang terlalu banyak dipanggil dapat menyalahi kapasitas layanan dan meningkatkan biaya operasional. Dengan mengurangi I/O yang tidak perlu dan membatasi akses backend hanya untuk kunci yang mungkin ada, sistem dapat mencapai skalabilitas yang lebih baik, responsivitas yang lebih tinggi, dan ketahanan yang lebih kuat terhadap beban lalu lintas yang tidak terduga.

## Model Mental
Bloom Filter adalah struktur data probabilistik yang memberikan dua jaminan asymmetric:
1. **Definitely Not in Set**: Jika filter mengembalikan false, kunci tersebut pasti tidak pernah ditambahkan ke filter (zero false negative).
2. **Possibly in Set**: Jika filter mengembalikan true, kunci mungkin ada dalam set, tetapi dapat menghasilkan false positive dengan probabilitas yang dapat diatur.

Dengan sifat ini, Bloom Filter berfungsi sebagai filter admission yang efisien: kunci yang diketahui tidak ada dapat ditolak segera tanpa perlu mengakses sumber data yang lebih lambat.

## Konsep Inti
Bloom Filter terdiri dari:
- **Bit array** berukuran `m` bit, awalnya semua bernilai 0.
- **k fungsi hash** yang masing‑masing memetakan sebuah kunci ke posisi dalam bit array.

Saat menambahkan kunci:
1. Hitung k hash dari kunci.
2. Set bit pada posisi yang sesuai menjadi 1.

Saat memeriksa keanggotaan kunci:
1. Hitung k hash dari kunci.
2. Jika ada satu bit bernilai 0, kunci **definitely not in set**.
3. Jika semua bit bernilai 1, kunci **possibly in set** (mungkin merupakan false positive).

Parameter optimal dapat dihitung dari jumlah elemen yang diharapkan (`n`) dan target false positive rate (`ε`):
- `m = ceil(−n × ln(ε) / (ln 2)²)`
- `k = round((m/n) × ln 2)`

Untuk `ε = 0.01`, diperlukan sekitar 9,6 bit per elemen dan 7 fungsi hash.

## Skenario Kegagalan
Tanpa Bloom Filter:
- Setiap query kunci tidak ada di LSM‑Tree akan memindai semua segment (SSTable), menghasilkan satu disk read per segment.
- Request untuk kunci tidak existing dalam serangan cache penetration akan melewati cache dan langsung memanggil backend untuk setiap query.

Dengan Bloom Filter:
- Kun kunci tidak ada dapat ditolak oleh filter sebelum mengakses segment LSM‑Tree, mengurangi disk read secara signifikan.
- Cache yang dilengkapi dengan Bloom Filter sebagai admission gate akan memblokir kunci tidak known sebelum mencapai backend.

## Cara Kerja
Bloom Filter mengandalkan dua sifat:
1. **Linear bit setting**: Setiap kunci yang ditambahkan akan mengatur beberapa bit menjadi 1; bit yang sama dapat diset oleh banyak kunci.
2. **Hash collision → false positive**: Probabilitas bahwa semua k hash dari kunci tidak ada menunjuk ke bit yang telah di‑set oleh kunci lain memberi peluang false positive.

Karena bit hanya dapat di‑set dari 0 ke 1 (dan tidak pernah kembali ke 0 dalam standar Bloom Filter), suatu kunci yang pernah ditambahkan akan selalu menghasilkan semua bit bernilai 1 saat diperiksa—oleh karena itu, **tidak ada false negative**.

## Arsitektur Implementasi Lab
Lab ini terdiri dari tiga komponen utama:
1. **`bloom.Filter`**: Implementasi inti dengan bit array `[]uint64` dan double hashing Kirsch‑Mitzenmacher.
2. **`store.LSMStore`**: Simulasi LSM‑Tree dengan segment immutable; setiap segment dapat memiliki Bloom Filter untuk pruning disk read.
3. **`store.Cache`**: Cache in‑memory dengan optional Bloom Filter admission gate untuk mengurangi panggilan backend.

Alur data:
- **LSM Lookup**: Jika segment memiliki filter dan `Check()` mengembalikan false, segment dilewati tanpa meningkatkan counter disk read.
- **Cache Lookup**: Jika cache memiliki filter dan `Check()` mengembalikan false, backend tidak dipanggil.

## Implementasi Inti
### Bit Array dan Hashing
Filter menggunakan `[]uint64` untuk representasi bit array yang word‑aligned, mengurangi overhead per bit. Dua hash dasar (`h1`, `h2`) di‑generate menggunakan FNV‑1a dan transformasi rotate‑XOR untuk memberikan hash independen sesuai teknik Kirsch‑Mitzenmacher:
```go
func baseHashes(key []byte) (uint64, uint64) {
    h := fnv.New64a()
    _, _ = h.Write(key)
    h1 := h.Sum64()
    h2 := h1 ^ (h1>>17 | h1<<47) // rotate‑xorshift
    if h2 == 0 {
        h2 = 0xdeadbeefdeadbeef
    }
    return h1, h2
}
```
Posisi bit ke‑i dihitung sebagai `doubleHash(h1, h2, i, m) = (h1 + i*h2) % m`.

### Penambahan dan Pengecekan
```go
func (f *Filter) Add(key []byte) {
    h1, h2 := baseHashes(key)
    for i := uint(0); i < f.k; i++ {
        pos := doubleHash(h1, h2, uint64(i), uint64(f.m))
        f.bits[pos/64] |= 1 << (pos % 64)
    }
}

func (f *Filter) Check(key []byte) bool {
    h1, h2 := baseHashes(key)
    for i := uint(0); i < f.k; i++ {
        pos := doubleHash(h1, h2, uint64(i), uint64(f.m))
        if f.bits[pos/64]&(1<<(pos%64)) == 0 {
            return false
        }
    }
    return true
}
```

## Penjelasan Kode
### Snippet 1 — Inisialisasi Filter Optimal
Source File: `internal/bloom/bloom.go:22-31`
Purpose: Menghitung ukuran bit array dan jumlah hash berdasarkan target FP rate dan jumlah elemen.
```go
func New(n uint, epsilon float64) *Filter {
    m := optimalM(n, epsilon)
    k := optimalK(m, n)
    words := (m + 63) / 64
    return &Filter{
        bits: make([]uint64, words),
        m:    m,
        k:    k,
    }
}
```
Explanation: Konstruktor menggunakan rumus tertutup dari paper Bloom (1970) dan Kirsch‑Mitzenmacher (2006) untuk menentukan `m` dan `k` yang optimal.

### Snippet 2 — Double Hashing Kirsch‑Mitzenmacher
Source File: `internal/bloom/bloom.go:79-83`
Purpose: Menghasilkan k posisi bit dari dua hash dasar untuk mengurangi overhead komputasi.
```go
func doubleHash(h1, h2, i, m uint64) uint64 {
    return (h1 + i*h2) % m
}
```
Explanation: Teknik ini membuktikan bahwa penggunaan `h1 + i*h2 (mod m)` memberikan distribusi yang serupa dengan k hash independen tanpa memanggil fungsi hash berulang kali.

### Snippet 3 — Pruning Segment LSM‑Tree
Source File: `internal/store/lsm.go:56-66`
Purpose: Melewati segment bila filter menyatakan kunci tidak ada, sehingga tidak increment counter disk read.
```go
if seg.filter != nil && !seg.filter.Check(keyBytes) {
    // Bloom filter: definitely not in this segment; skip disk read!
    continue
}
```
Explanation: Jika filter mengembalikan false, kunci tersebut tidak dapat ada di segment tersebut (zero false negative), sehingga operasi baca data (yang mensimulasikan disk I/O) dapat di‑skip sepenuhnya.

## Apa yang Dibuktikan oleh Tes
- **Zero False Negatives**: `TestNoFalseNegatives` menambahkan 10 000 kunci dan memastikan semua kembali true.
- **False Positive Rate**: `TestFalsePositiveRate` mengukur FP rate ~0,996% untuk target 1,00% dengan 50 000 query kunci absent.
- **Optimal Sizing**: `TestOptimalSizing` memverifikasi bahwa `m` dan `k` sesuai dengan rumus teoritis (±1 ULP).
- **Disk I/O Reduction**: `TestLSMStoreWithAndWithoutFilter` menunjukkan disk read turun dari 25 000 (tanpa filter) menjadi 276 (dengan filter) untuk 5 000 query kunci absent—penurunan >99 %.
- **Cache Penetration Mitigation**: `TestCachePenetrationMitigation` menunjukkan backend call turun dari 1000 (tanpa filter) menjadi 0 (dengan filter) untuk 1 000 query kunci absent dalam skenario serangga.
- **Keamanan Konkuren**: `TestSyncFilterConcurrency` menjalankan 8 penulis dan 16 pembaca secara bersamaan tanpa race detector mendeteksi konflik.

## Pemulihan / Rollback
Bloom Filter standar tidak mendukung operasi delete karena penghapusan bit dapat menyebabkan false negative untuk kunci lain yang berbagi bit yang sama. Dalam lab ini:
- Segment LSM‑Tree bersifat immutable; filter tidak pernah diperbarui setelah dibuat.
- Cache hanya menambahkan kunci baru ke filter ketika data di‑fetch dari backend; tidak ada penghapusan eksplisit.
Jika diperlukan penghapusan, solusi yang direkomendasikan adalah rebuild filter atau menggunakan varian seperti Counting Bloom Filter atau Cuckoo Filter (lihat bagian Pertimbangan Produksi).

## Pertimbangan Produksi
- **Static Sizing**: Filter di‑size pada estimasi `n` dan `ε`. Jika jumlah elemen melebihi `n`, false positive rate akan meningkat. Untuk sistem dengan pertumbuhan tidak terduga, pertimbangkan Scalable Bloom Filter atau mekanisme rebuild periodik.
- **Tidak Mendukung Delete**: Standar Bloom Filter tidak dapat menghapus elemen tanpa membangun ulang seluruh bit array. Jika delete diperlukan, gunakan Counting Bloom Filter (memori ×3‑4) atau Cuckoo Filter.
- **Hash Seed Deterministik**: Untuk filter yang disimpan dan dibaca kembali (misalnya di SSTable), pastikan hash function menggunakan seed tetap agar hasil konsisten antar proses atau node.
- **Overhead CPU**: Double hashing mengurangi panggilan hash, tetapi setiap probe masih membutuhkan operasi bitwise. Dalam kriteria latency ultra‑rendah, pertimbangkan varian block‑based atau Ribbon Filter yang mengurangi cache miss.
- **Keamanan terhadap Hash Flooding**: FNV‑1a dapat rentan terhadap collision yang disengaja jika input dikendalikan oleh attacker. Dalam konteks adversarial, pertimbangkan hash kriptografik yang lebih kuat (misalnya SipHash) atau gunakan tabel hash yang diacak per instance.

## Kesalahan Umum
- **Mengira false positive rate aman untuk semua n**: Pengguna sering lupa bahwa FP rate bersifat komposisional; jika jumlah elemen aktual melebihi estimasi awal, rate dapat naik drastis.
- **Mengharapkan filter dapat memberi nilai asosiasi**: Bloom Filter hanya menyertai keanggotaan himpunan, bukan mapping kunci‑nilai.
- **Mengabaikan kostbarnya bit array yang terlalu besar**: Memilih `ε` yang sangat kecil (misalnya 0,0001) dapat mengakibatkan alokasi memori yang tidak proporsional bila tidak diperhitungkan dengan bijak.
- **Menggunakan filter yang tidak diskronkan dengan data**: Dalam sistem terdistribusi, filter harus selalu mencerminkan set kunci yang sebenarnya; filter usus dapat menyebabkan peningkatan false positive atau false negative jika tidak diperbarui.

## Studi Kasus
Demonstrasi lab menunjukkan tiga manfaat kuat:
1. **Verifikasi Matematika**: Untuk 10 000 elemen dengan `ε=0,01`, filter mengalokasikan 95 851 bit (≈9,59 bit/elemen) dan menggunakan 7 hash. Tidak ada false negative; false positive terukur 0,996%.
2. **Pengurangan I/O LSM‑Tree**: Dengan 10 segment masing‑masing 2 000 kunci, 10 000 query kunci absent menghasilkan 100 000 disk read tanpa filter, tetapi hanya 988 read dengan filter—penurunan I/O 99,01%.
3. **Pencegahan Cache Penetration**: 5 000 request kunci absent menghasilkan 5 000 panggilan backend tanpa filter, tetapi 0 panggilan dengan filter—penetrasi berkurang ke 0%.

## Daftar Periksa
- [ ] Pastikan estimasi `n` tidaksepuluhnya underestimated; gunakan padding atau mekanisme rebuild.
- [ ] Verifikasi bahwa penggunaan thread‑safe (`SyncFilter`) diperlukan bila ada konkuren tulis/baca.
- [ ] Konfirmasi bahwa hash function memberikan distribusi seragam untuk kunci target.
- [ ] Ukur false positive rate empiris pada beban produksi sebelum menetapkan `ε` di lingkungan latihan.
- [ ] Dokumentasikan asumsi tentang immutable atau append‑only nature data bila menggunakan filter pada segment penyimpanan.

## Poin Kunci
- Bloom Filter memberikan jaminan **zero false negative** dan false positive rate yang dapat diatur dengan tepat.
- Memori yang diperlukan sangat rendah: ~9,6 bit per elemen untuk 1% false positive rate.
- Dalam LSM‑Tree, filter per‑segment dapat mengurangi disk read untuk kunci absent lebih dari 99 %.
- Sebagai admission gate pada cache, filter dapat mencegah cache penetration hampir sepenuhnya.
- Implementasi lab menggunakan hanya standar library Go (hash/fnv, math, sync) tanpa dependensi pihak ketiga.
- Filter bersifat statis dan tidak mendukung delete; pertimbangkan solusi alternatif bila diperlukan.
- Semua klaim didukung oleh tes unit dan integrasi yang lolos dengan race detector.

## Sumber
Berikut pemetaan bagian ke sumber yang disetujui: