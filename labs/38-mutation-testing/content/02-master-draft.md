# Mutation Testing: Menguji Kualitas Pengujian di Balik Ilusi 100% Code Coverage

## Problem

Banyak tim rekayasa perangkat lunak mengandalkan metrik *code coverage* (line coverage, statement coverage, branch coverage) sebagai indikator utama kesehatan dan keandalan sistem pengujian. Target seperti 80% hingga 100% cakupan baris kode kerap dijadikan syarat wajib dalam *pull request review* atau *CI quality gate*.

Namun, metrik cakupan tradisional menyimpan kelemahan mendasar: cakupan kode hanya mengukur baris atau pernyataan apa yang dieksekusi saat pengujian dijalankan. Metrik ini sama sekali tidak memvalidasi apakah pengujian yang berjalan memiliki asersi (*assertion*) yang sanggup mendeteksi kerusakan logika (*fault detection*).

Akibatnya timbul ilusi keamanan (*false confidence*). Sebuah *test suite* dapat mengeksekusi 100% baris kode tanpa menyertakan asersi substantif—misalnya hanya memastikan fungsi mengembalikan nilai tanpa error atau memeriksa `result > 0`. Ketika ada kesalahan operator, batasan nilai (*boundary*), atau logika boolean yang terbalik pada kode produksi, seluruh pengujian tetap lolos (*green*), dan *bug* kritis lolos ke *production*.

## Why This Matters

Mengandalkan *code coverage* saja memunculkan celah kualitas yang berbahaya:

1. **Ilusi Kematangan Pengujian**: Laporan CI menunjukkan 100% *coverage*, memberi rasa aman semu kepada manajemen dan *reviewer*, padahal *test suite* tumpul terhadap mutasi logika.
2. **Regresi Bisnis Tersembunyi**: Perubahan operator relational sederhana (misalnya batas diskon berubah dari `>= 100.0` menjadi `> 100.0`, atau diskon kupon berubah dari `+` menjadi `-`) sering lolos dari *test suite* yang asersinya longgar.
3. **Biaya Pemeliharaan**: Tim menghabiskan energi menulis tes demi menaikkan angka cakupan, namun tes tersebut tidak memberikan perlindungan nyata saat dilakukan *refactoring*.

*Mutation testing* hadir untuk menguji pengujian itu sendiri: alih-alih menguji kode aplikasi menggunakan tes, teknik ini menyuntikkan kesalahan (*fault injection*) ke kode aplikasi untuk melihat apakah *test suite* mampu menangkap kesalahan tersebut.

## Mental Model

Pembeda konseptual antara cakupan tradisional dan *mutation testing* adalah:

> **Code coverage** mengukur apa yang **dieksekusi** (*what is executed*).  
> **Mutation testing** mengukur apa yang **diverifikasi** (*what is verified*).

Analogi praktis dari Stryker: bayangkan sepotong roti yang diolesi selai. *Code coverage* melaporkan bahwa 80% permukaan roti telah tertutup selai. Namun, *mutation testing* memeriksa apakah selai tersebut benar-benar selai cokelat asli atau zat lain yang keliru.

Jika kode aplikasi diubah sedikit secara sengaja (*mutant* dibuat) dan *test suite* gagal mendeteksi perubahan tersebut (semua tes tetap *pass*), mutan tersebut dinyatakan selamat (*survived*). Mutan yang bertahan adalah bukti nyata adanya lubang asersi pada *test suite*.

## Core Concept

### 1. Definisi Mutan dan Skor Mutasi
*Mutant* adalah salinan kode sumber yang telah disisipi perubahan sintaksis kecil yang meniru kesalahan manusiawi yang lazim terjadi saat pemrograman.

Kualitas dari suatu *test suite* dievaluasi menggunakan **Mutation Score**:

$$\text{Mutation Score} = \left(\frac{\text{Killed Mutants}}{\text{Total Mutants}}\right) \times 100\%$$

- **Killed Mutant**: *Test suite* mendeteksi adanya mutasi (setidaknya satu tes gagal/merah). Ini hasil yang diharapkan.
- **Survived Mutant**: Seluruh pengujian tetap lolos (hijau) meskipun kode telah dirusak. Ini mengindikasikan asersi pengujian lemah atau input tes tidak memadai.
- **Equivalent Mutant**: Mutan yang secara semantik menghasilkan perilaku identik dengan kode asli, sehingga secara teoretis mustahil untuk digagalkan oleh tes apa pun.

### 2. Fondasi Teoretis
*Mutation testing* bertumpu pada dua hipotesis fundamental (DeMillo et al., 1978; Offutt, 1992):
- **Competent Programmer Hypothesis**: *Programmer* kompeten umumnya menulis kode yang hampir benar; kesalahan logika yang timbul biasanya berupa kekeliruan sintaksis kecil (seperti salah operator atau tanda pembanding).
- **Coupling Effect**: Kesalahan-kesalahan sederhana yang saling terkait dapat bertingkat (*cascade*) membentuk kesalahan kompleks (*emergent faults*). Jika pengujian mampu menangkap kesalahan sederhana, pengujian tersebut juga efektif menangkap kegagalan yang lebih rumit.

### 3. Model RIP (Reach, Infect, Propagate)
Agar sebuah tes berhasil membunuh mutan (*killed*), tiga kondisi dalam model RIP (Offutt & Untch, 2000) harus terpenuhi secara berurutan:
1. **Reach**: Eksekusi pengujian harus mencapai baris kode atau pernyataan yang dimutasi. (Ini setara dengan *code coverage*).
2. **Infect**: Input data pengujian harus menyebabkan *state* internal program menjadi salah (*infected*) akibat mutasi tersebut.
3. **Propagate**: *State* yang salah tersebut harus menjalar hingga ke nilai keluaran (*output*) program dan diperiksa secara eksplisit oleh asersi pengujian.

*Weak assertion testing* gagal pada tahap **Propagate**: tes mencapai kode (*Reach*) dan mutasi mengubah variabel (*Infect*), tetapi asersi tes tidak memeriksa nilai spesifik tersebut sehingga mutan tetap selamat (*Survived*).

## Failure Scenario

Bayangkan sebuah modul *e-commerce* yang menghitung diskon dan ongkos kirim. Salah satu aturan bisnis menyatakan bahwa pesanan dengan nilai akhir di atas atau sama dengan 200 berhak atas *free shipping*:

```go
freeShipping := (order.Tier == TierVIP) || (finalAmount >= 200.0)
```

Sebuah *test suite* yang lemah ditulis sebagai berikut:

```go
res := CalculateDiscount(Order{TotalAmount: 1200.0, ItemCount: 3, Tier: TierVIP, HasCoupon: true})
if res.FinalAmount <= 0 {
    t.Errorf("Expected positive final amount, got %v", res.FinalAmount)
}
```

Uji coba di atas mengeksekusi seluruh baris dalam fungsi dan menghasilkan 100% *statement coverage*. Namun perhatikan apa yang terjadi jika operator `||` diubah menjadi `&&`, atau `>= 200.0` diubah menjadi `> 200.0`:

- Input tes memiliki total 1200.0 dan nilai akhir yang jauh lebih besar dari nol.
- Pemeriksaan `res.FinalAmount <= 0` tidak pernah terpicu.
- Asersi sama sekali tidak memeriksa nilai `res.FreeShipping`.
- Kerusakan aturan bisnis lolos tanpa terdeteksi.

## How It Works

Alur kerja dasar mesin *mutation testing* berbasis evaluasi:

```text
[Source Code] ───(Parse AST)───> [Identifikasi Node BinaryExpr / BasicLit]
                                            │
                                            ▼
                               [Generate Mutation Plans]
                                            │
                                            ▼
                          [Iterasi Mutan (Paralel / Terisolasi)]
                                ┌───────────┴───────────┐
                                ▼                       ▼
                        [Terapkan Mutasi]       [Render Kode Baru]
                                │                       │
                                └───────────┬───────────┘
                                            ▼
                                 [Eksekusi Test Suite]
                                            │
                           ┌────────────────┴────────────────┐
                           ▼                                 ▼
                     [Test Gagal]                       [Test Lolos]
                           │                                 │
                           ▼                                 ▼
                   Status: KILLED                    Status: SURVIVED
                                            │
                                            ▼
                              [Kalkulasi Mutation Score]
```

## Architecture

Pada lab ini, arsitektur mesin mutasi dirancang berbasis Abstract Syntax Tree (AST) murni menggunakan pustaka standar Go (`go/ast`, `go/parser`, `go/token`, `go/format`) tanpa dependensi eksternal:

```text
labs/38-mutation-testing/
├── cmd/demo/main.go            # Entry point CLI runner untuk demonstrasi
├── internal/
│   ├── engine/                 # Inti mesin mutation testing
│   │   ├── types.go            # Model data: Mutant, MutationPlan, Report
│   │   ├── mutator.go          # Traversal AST dan generator mutasi
│   │   └── runner.go           # Eksekusi konkuren terisolasi per mutan
│   └── service/                # Domain logika bisnis
│       ├── discount.go         # Logika diskon & shipping target mutasi
│       ├── discount_weak_test.go   # Test suite 100% coverage (asersi lemah)
│       └── discount_strong_test.go # Test suite komprehensif (asersi kuat)
└── tests/
    └── engine_test.go          # Pengujian unit mesin mutasi
```

Komponen utama sistem:
- **`ASTMutator` (`internal/engine/mutator.go`)**: Melakukan inspeksi pohon sintaksis kode Go, mendeteksi node operator biner dan literal konstanta, lalu membentuk rencana mutasi (*MutationPlan*) yang dilengkapi fungsi *undo*.
- **`Runner` (`internal/engine/runner.go`)**: Mengelola eksekusi mutan secara konkuren. Setiap *goroutine* mem-parsing salinan AST baru (*fresh file set*) untuk menjamin isolasi mutasi, menjalankan fungsi tes (`TestFunc`), dan mengembalikan laporan statistik (*Report*).
- **Target Bisnis (`internal/service/discount.go`)**: Fungsi `CalculateDiscount` yang menggabungkan operasi relasional, logika boolean majemuk, kalkulasi aritmetika, dan nilai batas.

## Implementation

Mesin mutasi pada lab ini menerapkan empat kategori operator mutasi yang umum:

1. **Relational Operator Replacement**: Mengubah operator relasional untuk menguji batas kondisi:
   - `>=` menjadi `>`
   - `>` menjadi `>=`
   - `==` menjadi `!=`
2. **Boolean Operator Flip**: Membalik logika percabangan:
   - `||` menjadi `&&`
   - `&&` menjadi `||`
3. **Arithmetic Operator Replacement**: Membalik kalkulasi matematis:
   - `*` menjadi `/`
   - `-` menjadi `+`
4. **Boundary Value Shift**: Menggeser konstanta bilangan bulat positif sebesar $+1$ (misal `2` menjadi `3`).

Setiap perubahan dienkapsulasi ke dalam struktur `MutationPlan` dengan *closure* pembalik (*undo*):

```go
type MutationPlan struct {
    Type        MutationType
    Description string
    Line        int
    Original    string
    Mutated     string
    Apply       func(file *ast.File) func()
}
```

Hal ini memastikan AST dapat dikembalikan ke keadaan semula atau dimanipulasi per *goroutine* tanpa *data race*.

## Code Walkthrough

### 1. Pembangkitan Mutasi melalui Inspeksi AST (`internal/engine/mutator.go`)

Fungsi `GenerateMutations` menelusuri AST menggunakan `ast.Inspect`. Ketika menemukan `*ast.BinaryExpr`, mutator memeriksa token operator:

```go
switch expr.Op {
case token.GEQ:
    plans = append(plans, MutationPlan{
        Type:        RelationalOpReplace,
        Description: "Replace '>=' with '>'",
        Line:        pos.Line,
        Original:    ">=",
        Mutated:     ">",
        Apply: func(f *ast.File) func() {
            expr.Op = token.GTR
            return func() { expr.Op = token.GEQ }
        },
    })
```

Untuk konstanta bilangan bulat (`*ast.BasicLit` dengan `token.INT`), mutator menaikkan nilai sebesar 1:

```go
val, err := strconv.Atoi(expr.Value)
if err == nil && val > 0 {
    orig := expr.Value
    plans = append(plans, MutationPlan{
        Type:        BoundaryValueMutate,
        Description: "Shift integer constant +1",
        Line:        pos.Line,
        Original:    orig,
        Mutated:     strconv.Itoa(val + 1),
        Apply: func(f *ast.File) func() {
            expr.Value = strconv.Itoa(val + 1)
            return func() { expr.Value = orig }
        },
    })
}
```

### 2. Eksekusi Paralel yang Terisolasi (`internal/engine/runner.go`)

Agar eksekusi tidak saling mengontaminasi state AST, `Runner.Run` membuat satu `token.NewFileSet()` baru dan mem-parsing ulang `sourceCode` pada setiap *goroutine*:

```go
for i, plan := range plans {
    wg.Add(1)
    go func(idx int, p MutationPlan) {
        defer wg.Done()

        fset := token.NewFileSet()
        parsedFile, err := parser.ParseFile(fset, "source.go", sourceCode, 0)
        if err != nil { /* tangani parse error */ return }

        undo := p.Apply(parsedFile)
        mutatedSrc, err := renderSource(fset, parsedFile)
        undo()

        killed := testFn(mutatedSrc)
        status := StatusSurvived
        if killed {
            status = StatusKilled
        }
        results[idx] = MutantResult{ ... }
    }(i, plan)
}
wg.Wait()
```

Hasil disimpan langsung ke *slice* `results` yang telah dialokasikan terlebih dahulu berdasarkan indeks (`results[idx]`), menghilangkan kebutuhan *mutex* saat penulisan.

### 3. Kontras Antara Asersi Lemah dan Asersi Kuat

Target pengujian adalah fungsi `CalculateDiscount` pada `internal/service/discount.go`:

```go
func CalculateDiscount(order Order) DiscountResult {
    rate := 0.0
    if order.Tier == TierVIP || order.TotalAmount >= 1000.0 {
        rate = 0.20
    } else if order.Tier == TierPremium && order.TotalAmount >= 500.0 {
        rate = 0.10
    } else if order.TotalAmount >= 100.0 {
        rate = 0.05
    }

    if order.HasCoupon && order.ItemCount > 2 {
        rate = rate + 0.05
    }

    discountAmount := order.TotalAmount * rate
    finalAmount := order.TotalAmount - discountAmount
    freeShipping := (order.Tier == TierVIP) || (finalAmount >= 200.0)

    return DiscountResult{ ... }
}
```

Pada `internal/service/discount_weak_test.go`, pengujian mengeksekusi 4 kasus yang melintasi setiap cabang logika:
- Kasus 1: VIP customer (`TotalAmount: 1200.0`, `ItemCount: 3`)
- Kasus 2: Premium customer (`TotalAmount: 600.0`, `ItemCount: 1`)
- Kasus 3: Standard customer (`TotalAmount: 150.0`, `ItemCount: 1`)
- Kasus 4: Standard customer (`TotalAmount: 50.0`, `ItemCount: 4`)

Namun asersi yang digunakan sangat longgar:
```go
if res1.FinalAmount <= 0 { t.Errorf(...) }
if res2.DiscountTotal < 0 { t.Errorf(...) }
if res3.DiscountRate < 0 { t.Errorf(...) }
if res4.OriginalTotal != 50.0 { t.Errorf(...) }
```
Semua asersi ini tetap benar meskipun diskon salah hitung atau batas diskon bergeser.

Sebaliknya, pada `internal/service/discount_strong_test.go`, pengujian menggunakan pendekatan *table-driven test* dengan 11 skenario presisi tinggi, mencakup nilai tepat di batas (`100.0`, `500.0`, `1000.0`), tepat di bawah batas (`99.99`, `499.99`), ambang *coupon item count* (`> 2` vs `<= 2`), ambang *free shipping* (`200.0` vs `190.0`), serta kasus degenerasi (`0.0`). Setiap pengujian memvalidasi `expectedRate`, `expectedDiscount`, `expectedFinal`, dan `expectedFreeShipping` secara eksak.

## What the Tests Prove

Berdasarkan eksekusi langsung pada lab:

1. **Cakupan Baris Penuh pada Pengujian Lemah**:
   Menjalankan perintah:
   ```bash
   go test -run=TestCalculateDiscount_Weak -coverprofile=weak_cov.out ./internal/service
   go tool cover -func=weak_cov.out
   ```
   Membuktikan bahwa fungsi `CalculateDiscount` mencapai **100.0% statement coverage**.

2. **Skor Mutasi Pengujian Lemah Adalah Nol**:
   Saat `cmd/demo/main.go` dieksekusi terhadap 15 mutan sintaksis yang dihasilkan dari `discount.go`:
   - Mutan dihasilkan: 8 relasional, 4 boolean, 2 aritmetika, 1 batasan nilai (total 15 mutan).
   - Pengujian lemah membunuh: **0 / 15 mutan (Mutation Score: 0.00%)**.
   - Seluruh 15 mutan selamat (*Survived*), membuktikan bahwa 100% *code coverage* tidak menjamin kemampuan deteksi kesalahan (*zero fault detection*).

3. **Skor Mutasi Pengujian Kuat Mencapai 100%**:
   - Pengujian komprehensif membunuh: **15 / 15 mutan (Mutation Score: 100.00%)**.
   - Setiap modifikasi operator atau batas langsung tertangkap oleh asersi eksak.

4. **Keamanan Konkurensi**:
   Eksekusi `go test -race ./...` lolos tanpa *data race* (0 race conditions), memvalidasi desain isolasi AST per *goroutine*.

## Production Considerations

Penerapan *mutation testing* pada skala produksi memerlukan pertimbangan arsitektur dan efisiensi:

1. **Beban Komputasi (*Computational Cost*)**:
   Jika sebuah basis kode memiliki 1.000 pengujian dan menghasilkan 2.000 mutan, mengeksekusi seluruh pengujian untuk setiap mutan berarti menjalankan pengujian sebanyak $2.000 \times 1.000 = 2.000.000$ kali.
   - **Solusi Praktis**: Gunakan *Incremental / Differential Mutation Testing* (seperti pada PIT dan Stryker) yang hanya memutasi baris kode yang berubah dalam suatu *Pull Request* atau *commit*.

2. **Masalah Equivalent Mutants**:
   Secara komputasi, mendeteksi apakah mutan bersifat ekivalen (tidak mengubah semantik output program) merupakan masalah *undecidable*. Mutan ekivalen akan selalu *survived* dan menurunkan skor mutasi secara artifisial. Alat bantu produksi menggunakan heuristik untuk menghindari mutasi pada konstruktor statis atau pola kode tertentu.

3. **Integrasi CI/CD sebagai Quality Gate**:
   *Mutation testing* dapat diintegrasikan ke dalam pipeline CI/CD. Namun, karena belum ada standar baku industri mengenai ambang batas nilai (*threshold*), tim disarankan menetapkan target secara bertahap dan empiris berdasarkan profil risiko domain bisnis (misal: fokus 90%+ pada modul kalkulasi finansial, tetapi lebih fleksibel pada lapisan antarmuka).

4. **Lanskap Perkakas Bahasa Go**:
   Ekosistem Go memiliki perkakas komunitas seperti `go-mutesting` dan `gremlins`. Namun kematangannya (fitur integrasi CI, pelaporan visual, inkrementalitas) masih berkembang jika dibandingkan perkakas pada ekosistem JVM (PIT) atau JavaScript/TypeScript (Stryker).

## Common Mistakes

1. **Menyamakan 100% Code Coverage dengan Kualitas Pengujian**: Mengasumsikan kode sudah bebas regresi hanya karena semua baris dilewati pengujian.
2. **Asersi Kosong atau Longgar (*Assertion-free Testing*)**: Menulis tes yang hanya mengecek `err == nil` atau `res != nil` tanpa memeriksa validitas mutasi data.
3. **Memaksa Target Skor Mutasi 100% pada Seluruh Proyek**: Berupaya membunuh seluruh mutan tanpa memfilter *equivalent mutants* atau kode *boilerplate*, yang berujung pada kelelahan tim (*burnout*).
4. **Menjalankan Full Mutation Test pada Setiap Commit Kecil**: Tidak menggunakan mode inkremental pada proyek berukuran besar sehingga siklus CI menjadi sangat lambat.

## Case Study

Berikut hasil demonstrasi aktual perbandingan pengujian pada modul perhitungan diskon pelanggan (`internal/service/discount.go`):

| Metrik | Weak Test Suite | Strong Test Suite | Keterangan |
| :--- | :---: | :---: | :--- |
| **Statement Coverage** | **100.0%** | **100.0%** | Keduanya mengeksekusi seluruh baris kode |
| **Total Mutan Terbentuk** | 15 | 15 | 8 Relasional, 4 Boolean, 2 Aritmetika, 1 Batasan |
| **Mutan Terbunuh (Killed)** | **0** | **15** | Pengujian lemah gagal mendeteksi semua mutasi |
| **Mutan Selamat (Survived)** | **15** | **0** | Pengujian kuat menangkap seluruh 15 perubahan |
| **Mutation Score** | **0.00%** | **100.00%** | Kontras absolut efektivitas asersi |

Ketika operator relasional pada baris 35 diubah dari `order.TotalAmount >= 100.0` menjadi `order.TotalAmount > 100.0`:
- Pengujian lemah dengan input `TotalAmount: 150.0` tetap bernilai benar dan asersi `res3.DiscountRate < 0` tetap terpenuhi. Mutan bertahan hidup.
- Pengujian kuat dengan input batas eksak `TotalAmount: 100.0` langsung mendeteksi bahwa diskon menjadi 0% padahal seharusnya 5%. Asersi gagal dan mutan langsung terbunuh.

## Checklist

- [ ] Evaluasi pengujian kritis tidak hanya bersandar pada laporan *line coverage*.
- [ ] Pastikan asersi pengujian memeriksa keluaran kalkulasi secara eksak, bukan sekadar nilai positif atau bukan nol.
- [ ] Masukkan pengujian batas (*boundary value tests*) untuk operator perbandingan (`>`, `<`, `>=`, `<=`).
- [ ] Periksa kembali percabangan logika boolean majemuk (`&&`, `||`) dengan skenario kombinatorial yang relevan.
- [ ] Jalankan *mutation testing* secara berkala atau bertahap pada *domain logic* bertaraf risiko tinggi.

## Key Takeaways

1. *Code coverage* mengukur baris yang dieksekusi, sedangkan *mutation testing* mengukur bug yang berhasil dideteksi.
2. 100% *code coverage* dengan asersi lemah menghasilkan *mutation score* 0%.
3. *Mutation testing* bekerja dengan menyuntikkan kesalahan terarah (operator relasional, logika boolean, operasi aritmetika, pergeseran batasan) menggunakan analisis AST.
4. Model RIP (Reach, Infect, Propagate) menegaskan bahwa kesalahan harus dieksekusi, merusak state, dan menjalar ke asersi tes untuk dapat dibunuh.
5. Mutan yang selamat (*survived*) mengindikasikan ketiadaan asersi atau data uji yang tidak menyentuh kondisi batas.
6. Pada sistem skala besar, *mutation testing* idealnya dijalankan secara diferensial/inkremental untuk efisiensi CI/CD.

## Sources

- Research Report & Evidence: `labs/38-mutation-testing/research/05-report.md`, `03-evidence.md`
- Research Audit Verdict: `labs/38-mutation-testing/research-audit/07-verdict.md`
- Engineering Implementation & Notes: `labs/38-mutation-testing/engineering/01-design.md`, `02-implementation-notes.md`
- Engineering Execution Result & Demo Output: `labs/38-mutation-testing/engineering/03-execution-result.md`, `cmd/demo/main.go`
- Source Code & Tests: `internal/engine/`, `internal/service/`, `tests/engine_test.go`
- Referensi Konseptual Terverifikasi: PIT Mutation Testing (pitest.org), Stryker Mutator (stryker-mutator.io), DeMillo et al. (1978), Offutt & Untch (2000).
