# Panduan Lengkap Property-Based Testing (PBT): Menemukan Edge-Case Tersembunyi Melalui Universal Invariant dan Counterexample Shrinking

## Problem

Dalam pengembangan perangkat lunak, pendekatan pengujian paling umum adalah *Example-Based Testing* (unit test konvensional). Developer menentukan sebuah set input tertentu, mengeksekusi fungsi, dan mencocokkannya dengan output yang diharapkan (`assert f(x) == y`).

Kelemahan fundamental dari pengujian berbasis contoh adalah:
1. **Bias Developer:** Pengembang hanya menguji skenario yang mereka bayangkan (sering kali hanya *happy path* atau edge case yang disadari).
2. **Kerapuhan Skenario Nyata:** Bug di lingkungan produksi umumnya dipicu oleh kombinasi data yang tidak terduga: bilangan negatif, presisi floating-point yang bocor, string kosong, array dengan urutan acak, atau boundary values ekstrim.
3. **False Confidence:** Rangkaian unit test bisa mencatat 100% *code coverage* dan seluruh test berstatus hijau (PASS), namun aplikasi langsung *crash* atau menghasilkan data korup begitu menerima input di luar contoh yang ditulis.

## Why This Matters

Menguji seluruh kemungkinan input secara manual mustahil karena ruang domain data bersifat tak hingga. Menulis ratusan baris unit test untuk variasi input statis membebani *maintenance* tanpa menjamin kebenaran logika.

*Property-Based Testing* (PBT) mengubah paradigma ini: alih-alih memikirkan contoh spesifik, engineer mendefinisikan aturan universal (*property/invariant*) yang harus selalu berlaku untuk setiap input yang valid. Mesin PBT kemudian mengeksekusi ratusan hingga ribuan kombinasi input acak secara otomatis untuk membuktikan apakah invariant tersebut dapat dirusak.

## Mental Model

PBT memisahkan peran antara manusia dan komputer:
- **Tugas Engineer:** Mendefinisikan *domain input* dan *universal invariant* yang harus selalu benar (`∀x ∈ Domain: Property(f(x)) == true`).
- **Tugas Mesin PBT:** Berperan sebagai *adversary* yang secara agresif mencari counterexample (input yang menggagalkan invariant) melalui *random generation*, lalu memperkecil input tersebut (*shrinking*) hingga diperoleh reproduksi bug paling minimal.

```text
+-------------------------------------------------------------+
|                     PROPERTY-BASED TESTING                  |
|                                                             |
|   1. GENERATION       2. INVARIANT CHECK    3. SHRINKING    |
|   +-------------+     +-----------------+   +-------------+ |
|   | Domain spec | --> | ∀x: Prop(f(x))  |-->| Minimal Repro| |
|   | & Boundary  |     |   Valid / Fail? |   | (e.g. [-1]) | |
|   +-------------+     +-----------------+   +-------------+ |
+-------------------------------------------------------------+
```

## Core Concept & Invariant Categories

PBT mengandalkan empat pola invariant kanonikal yang berlaku lintas domain:

### 1. Roundtrip Invariant (`Decode(Encode(x)) == x`)
Setiap data yang diubah ke representasi lain (misal: serialisasi JSON, formatting string, kompresi, enkripsi simetris) harus kembali ke bentuk aslinya tanpa kehilangan informasi jika diproses balik.

### 2. Idempotence Invariant (`f(f(x)) == f(x)`)
Menjalankan operasi berulang kali harus menghasilkan state yang sama dengan menjalankannya satu kali (misal: normalisasi format, sorting, merge interval, update state pembersihan).

### 3. Equivalence / Test Oracle (`f_optimized(x) == f_reference(x)`)
Hasil dari algoritma baru yang kompleks/teroptimasi harus identik dengan algoritma referensi yang sederhana namun terbukti benar (*naive implementation*).

### 4. Hard-to-Prove, Easy-to-Verify
Menemukan solusi mungkin memerlukan komputasi berat, namun memvalidasi keabsahan hasilnya sangat murah (misal: hasil *shortest path* harus berupa jalur kontinu valid tanpa siklus).

---

## Failure Scenario

Mari telaah dua skenario kegagalan nyata yang sering lolos dari pengujian konvensional:

### Kasus 1: Presisi Moneter (Floating-Point Rounding Leak)
Fungsi format mata uang sederhana yang menggunakan `float64` tampak bekerja sempurna pada contoh unit test standar seperti `$10.50`, `$99.99`, dan `$1.00`. Namun, ketika dihadapkan pada nilai dengan sub-cent atau representasi biner pecahan IEEE 754, parsing roundtrip gagal total karena representasi desimal tidak eksak.

### Kasus 2: Penggabungan Interval (Interval Merging Assumption)
Fungsi `Merge(intervals)` yang mengasumsikan input selalu terurut (*sorted*) akan lolos pengujian berbasis contoh yang rapi. Namun saat menerima irisan waktu yang acak, fungsi menghasilkan partisi yang salah atau kehilangan interval penting.

---

## How It Works & Architecture

Lab ini mengimplementasikan PBT pada Go menggunakan paket standar `testing/quick` serta generator dan shrinker custom.

```text
labs/40-property-based-testing/
├── cmd/demo/main.go           # CLI runner perbandingan Example vs PBT
├── internal/
│   ├── currency/             # Roundtrip invariant (Naive Float vs Robust Cents)
│   │   ├── currency.go
│   │   └── currency_test.go
│   ├── interval/             # Idempotence & Oracle invariants
│   │   ├── interval.go
│   │   └── interval_test.go
│   └── shrinker/             # Binary chunk & element reduction engine
│       ├── shrinker.go
│       └── shrinker_test.go
├── go.mod
└── README.md
```

---

## Code Walkthrough & What the Tests Prove

### 1. Verifikasi Roundtrip Currency (`internal/currency`)

Implementasi naive menggunakan `float64`:

```go
type NaiveCurrency struct{}

func (n NaiveCurrency) Format(dollars float64) string {
	return fmt.Sprintf("$%.2f", dollars)
}

func (n NaiveCurrency) Parse(s string) (float64, error) {
	s = strings.TrimPrefix(s, "$")
	return strconv.ParseFloat(s, 64)
}
```

Implementasi robust menggunakan integer cents (`int64`):

```go
type RobustAmount struct {
	Cents int64
}

func (r RobustAmount) Format() string {
	absCents := r.Cents
	sign := ""
	if absCents < 0 {
		sign = "-"
		absCents = -absCents
	}
	dollars := absCents / 100
	cents := absCents % 100
	return fmt.Sprintf("%s$%d.%02d", sign, dollars, cents)
}
```

#### Hasil Pengujian:
- **Example-Based Test:** Lolos 100% pada 4 contoh statis (`10.50, 99.99, 1.00, 5.25`).
- **Property-Based Test:** Gagal 100% pada 1000 iterasi random input `float64` karena `%.2f` memotong desimal dan representasi IEEE 754 tidak presisi.
- **Robust Roundtrip:** `ParseRobust(amount.Format()) == amount` lolos 1000/1000 iterasi di rentang nilai negatif, nol, hingga multi-miliar sen.

---

### 2. Verifikasi Idempotence & Oracle Interval Merge (`internal/interval`)

Implementasi robust memastikan pengurutan sebelum penggabungan:

```go
func RobustMerge(intervals []Interval) []Interval {
	if len(intervals) == 0 {
		return nil
	}
	sorted := make([]Interval, len(intervals))
	copy(sorted, intervals)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Start != sorted[j].Start {
			return sorted[i].Start < sorted[j].Start
		}
		return sorted[i].End < sorted[j].End
	})

	out := []Interval{sorted[0]}
	for _, iv := range sorted[1:] {
		last := &out[len(out)-1]
		if iv.Start <= last.End+1 {
			if iv.End > last.End {
				last.End = iv.End
			}
		} else {
			out = append(out, iv)
		}
	}
	return out
}
```

#### Hasil Pengujian:
- **Oracle Discrepancy:** `NaiveMerge` menghasilkan ketidaksesuaian sebesar 85% pada 100 input acak tidak terurut jika dibandingkan dengan `RobustMerge`.
- **Idempotence Property:** `RobustMerge(RobustMerge(x)) == RobustMerge(x)` terverifikasi valid pada 1.000 iterasi acak.
- **Non-Overlapping Property:** Seluruh interval keluaran terbukti tidak saling tumpang tindih.

---

### 3. Mekanisme Counterexample Shrinking (`internal/shrinker`)

Ketika PBT menemukan kegagalan pada array besar acak (misalnya 10–500 elemen), membaca log tersebut sangat sulit. *Shrinking* menyederhanakan input kegagalan ke bentuk paling mendasar.

Algoritma `FindAndShrink` bekerja dalam tiga fase:
1. **Binary Sectioning / Chunk Removal:** Mencoba memotong setengah bagian kiri/kanan array untuk melihat apakah kegagalan tetap terjadi.
2. **Individual Element Removal:** Menghapus elemen satu per satu.
3. **Value Reduction:** Memperkecil magnitudo nilai menuju 0 atau nilai batas terdekat.

#### Jejak Eksekusi Shrinking (Demo):
```text
Initial failing input (10 elements): [137 18 100 -27 95 107 126 78 -7 79]
  Step   1 remove-right-half         [5 elements] [137 18 100 -27 95] -> FAIL
  Step   3 remove-left-half          [3 elements] [100 -27 95]        -> FAIL
  Step   5 remove-left-half          [2 elements] [-27 95]            -> FAIL
  Step   6 remove-right-half         [1 elements] [-27]               -> FAIL
  Step   9 reduce-val-idx-0-to--13   [1 elements] [-13]               -> FAIL
  Step  12 reduce-val-idx-0-to--6    [1 elements] [-6]                -> FAIL
  Step  15 reduce-val-idx-0-to--3    [1 elements] [-3]                -> FAIL
  Step  18 reduce-val-idx-0-to--1    [1 elements] [-1]                -> FAIL

Minimal counterexample (1 elements): [-1]
```
Hanya dalam 21 langkah evaluasi, array 10 elemen tereduksi menjadi counterexample paling mendasar: `[-1]`.

---

## Production Considerations & Best Practices

1. **Biased Generators:** Pastikan generator tidak hanya menghasilkan distribusi seragam, namun condong ke nilai batas (*boundary bias*): `0`, `-1`, `1`, `MaxInt`, `MinInt`, string kosong, atau karakter kontrol.
2. **Determinisme & Seed Reproducibility:** Setiap kegagalan PBT harus mencatat nilai `seed` generator sehingga kegagalan dapat direproduksi 100% pada mesin lokal developer.
3. **Execution Time Budget:** Jalankan PBT dengan 100–1.000 iterasi pada *local pre-commit/PR check*, dan naikkan hingga 10.000–50.000 iterasi pada *nightly CI build*.

## Common Mistakes

- **Menulis Ulang Logika Implementasi di Dalam Assertion:** Assertion property tidak boleh menduplikasi kode fungsi yang diuji; fokuslah pada sifat relasional (*invariants*).
- **Generator Terlalu Sempit:** Membatasi generator hanya pada input yang validasi kodenya sudah diketahui bekerja, sehingga mengalahkan esensi penemuan bug.
- **Mengabaikan Determinisme:** Menggunakan random generator tanpa kemampuan menyuntikkan seed tetap saat debugging.

---

## Key Takeaways

1. **Example-Based Testing Menguji Asumsi; PBT Menguji Domain:** Unit test konvensional memverifikasi jalur yang diketahui, sedangkan PBT mencari celah pada domain input secara keseluruhan.
2. **Invariant Melindungi Desain:** Invariant seperti Roundtrip dan Idempotence bertindak sebagai spesifikasi formal yang dapat dieksekusi (*executable specifications*).
3. **Shrinking Menghemat Waktu Debugging:** Dari ribuan elemen acak, mesin PBT mengisolasi satu baris data minimal penyebab kegagalan.
4. **Zero Third-Party Dependency di Go:** Standar library `testing/quick` bawaan Go sudah cukup untuk membangun pengujian berbasis properti yang andal.
