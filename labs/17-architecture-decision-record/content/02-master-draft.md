# Mengelola Keputusan Arsitektur Menggunakan Architecture Decision Record (ADR) dan Validasi Struktural Otomatis

## Problem

Dalam siklus pengembangan perangkat lunak, tim rekayasa perangkat lunak sering kali kehilangan konteks historis di balik keputusan teknis penting. Ketika anggota tim baru bergabung atau sistem bertumbuh, sering muncul pertanyaan mengapa suatu pola arsitektur atau teknologi dipilih. 

Tanpa dokumentasi yang memadai mengenai *mengapa* suatu pilihan diambil dan batasan apa yang memengaruhinya pada saat itu, tim cenderung mengalami dua masalah berulang:
1. Membalikkan keputusan masa lalu secara sembarangan tanpa memahami batasan awal yang memicunya.
2. Mempertahankan keputusan lama secara buta meskipun kondisi sistem dan skala beban telah berubah drastis.

Ketiadaan rekaman keputusan ini memicu perdebatan teknis berulang (*cyclical technical debates*) yang menghabiskan waktu tim dan memperlambat laju pengiriman sistem.

## Why This Matters

Dokumentasi arsitektur tradisional sering kali berwujud dokumen spesifikasi desain besar yang terpisah dari repositori kode sumber (misalnya pada wiki internal atau sistem manajemen dokumen pihak ketiga). Pendekatan ini memiliki kelemahan mendasar:
- **Dokumentasi Usang (*Documentation Drift*):** Kode terus berevolusi melalui commit dan pull request, sedangkan dokumen eksternal tidak ikut diperbarui.
- **Konteks Terfragmentasi:** Pengembang yang membaca kode tidak memiliki visibilitas langsung terhadap latar belakang keputusan desain tersebut.
- **Ketidakmampuan Melacak Evolusi:** Sulit menentukan kapan dan mengapa keputusan tertentu digantikan oleh arsitektur yang baru.

Menurut pengamatan Michael Nygard (2011) dan panduan *AWS Prescriptive Guidance*, solusi dari masalah ini adalah mencatat keputusan yang signifikan secara arsitektural (*Architecturally Significant Decisions*) langsung di dalam repositori kode menggunakan format yang terstruktur, bernomor urut monotonik, dan memiliki siklus hidup yang tidak dapat diubah (*immutable*).

## Mental Model

Architecture Decision Record (ADR) dapat dipahami melalui prinsip-prinsip kunci berikut:

1. **Ko-lokasi dengan Kode Sumber:** ADR disimpan langsung dalam direktori proyek Git bersama kode (misalnya pada direktori `docs/adr/` atau `doc/arch/`). ADR berevolusi melalui mekanisme *pull request* dan *code review* yang sama dengan kode aplikasi.
2. **Kekekalan (*Immutability*):** Setelah sebuah ADR disetujui (*Accepted*), isinya tidak boleh disunting secara retroaktif untuk mengubah arah keputusan. Perubahan desain harus dituangkan ke dalam ADR baru.
3. **Silsilah Keputusan Berarah (*Directed Acyclic Graph / DAG*):** Keputusan baru yang mengubah keputusan lama tidak menghapus catatan lama, melainkan menandai ADR lama sebagai *Superseded* dan memberikan tautan eksplisit ke ADR baru yang menggantikannya (*Supersedes*).
4. **Fokus pada "Mengapa", Bukan Hanya "Bagaimana":** ADR tidak menduplikasi rincian implementasi tingkat rendah, melainkan menangkap konteks, gaya dorong/kendala (*forces*), alternatif yang dievaluasi, dan konsekuensi (positif, negatif, dan netral).

## Core Concept

### 1. Batasan Keputusan yang Signifikan (*Architecturally Significant*)
Tidak semua perubahan kode memerlukan ADR. Berdasarkan kriteria Richards & Ford (2020) yang diadopsi dalam literatur ADR, sebuah keputusan dinilai bernilai arsitektural jika berdampak langsung pada:
- **Struktur (*Structure*):** Pemisahan subsistem, batasan domain, atau pola modularitas.
- **Persyaratan Non-Fungsional (*Non-Functional Requirements / NFR*):** Skalabilitas, ketersediaan, performa, dan toleransi kegagalan.
- **Ketergantungan (*Dependencies*):** Pemilihan kerangka kerja utama, basis data, atau protokol integrasi eksternal.
- **Antarmuka (*Interfaces*):** Kontrak API, protokol komunikasi antarlayanan.
- **Teknik Konstruksi (*Construction Techniques*):** Standar pengujian, mekanisme penanganan transaksi terdistribusi, atau aturan batas paket.

Refaktorisasi internal rutin, penamaan variabel, atau perbaikan *bug* lokal tidak memenuhi kualifikasi ADR.

### 2. Siklus Hidup ADR (*ADR Lifecycle*)
Status ADR merepresentasikan kondisi terkini dari keputusan tersebut:
- **Proposed:** Keputusan sedang diajukan dan dalam tahap peninjauan tim.
- **Accepted:** Keputusan telah disetujui oleh tim dan menjadi panduan arsitektur yang aktif.
- **Superseded:** Keputusan tidak lagi aktif karena telah digantikan oleh keputusan arsitektur berikutnya.
- **Deprecated:** Keputusan dinonaktifkan atau ditinggalkan tanpa ada pengganti langsung.
- **Rejected:** Usulan keputusan ditolak setelah evaluasi; status ini dicatat agar argumen penolakan terdokumentasi dan mencegah tim mendiskusikan kembali opsi yang sama di masa depan tanpa bukti baru.

## Failure Scenario

Tanpa validasi integritas otomatis, pemeliharaan ADR berbasis Markdown manual rentan terhadap kegagalan struktural berikut:

1. **Tautan Penggantian Patah (*Dangling / Broken Supersession Link*):** ADR lama diberi status `Superseded by 99`, padahal ADR nomor 99 tidak pernah dibuat atau tidak ada di repositori.
2. **Penggantian Sepihak (*One-sided Supersession Link*):** ADR baru mengklaim `Supersedes: 1`, namun berkas ADR nomor 1 tidak diperbarui statusnya menjadi `Superseded by 2`, mengakibatkan inkonsistensi status operasional.
3. **Penomoran Tidak Monotonik (*Non-monotonic Numbering*):** Penomoran berkas melompati urutan (misalnya dari ADR 1 langsung ke ADR 3), yang memicu ambiguitas apakah ada keputusan yang hilang dari riwayat Git.
4. **Status Tidak Valid (*Invalid Lifecycle Status*):** Pengembang memasukkan status informal (seperti "Draft", "In Review", atau "Pending") yang tidak diakui dalam skema siklus hidup formal.

## How It Works

Untuk mencegah anomali struktural tersebut, sistem validasi (linter) bertindak sebagai gerbang pengujian kualitas pada siklus CI/CD atau pengujian lokal:

1. **Parsing Berkas:** Setiap berkas ADR dibaca dan diekstraksi atribut kuncinya melalui pola pencocokan ekspresi reguler (*regular expression*):
   - Nomor ID dan Judul ADR.
   - Status siklus hidup saat ini.
   - Metadata hubungan silsilah (`Superseded by <ID>` atau `Supersedes: <ID>`).
2. **Pemeriksaan Monotonik:** Semua ID dikumpulkan dan diurutkan untuk memastikan rangkaian ID utuh tanpa celah atau duplikasi (`1, 2, 3, ...`).
3. **Validasi Graf Silsilah Timbal-Balik (*Bidirectional Lineage Validation*):**
   - Jika berkas $A$ berstatus `Superseded by B`, validator memverifikasi keberadaan berkas $B$ dan memastikan berkas $B$ memiliki deklarasi `Supersedes: A`.
   - Sebaliknya, jika berkas $B$ mendeklarasikan `Supersedes: A`, validator memastikan berkas $A$ ada dan berstatus `Superseded by B`.
4. **Eksekusi Paralel:** Validasi simpul-simpul dalam graf keputusan dapat dieksekusi secara konkuren memanfaatkan konkurensi bawaan (Goroutines) dengan sinkronisasi muteks untuk agregasi galat.

## Architecture

Arsitektur sistem validasi ADR pada lab ini terdiri dari modul-modul berikut:

```text
[ Markdown ADR Strings / Files ]
               │
               ▼
      [ adr.Parser ]  ── (Regex extraction: ID, Title, Status, References)
               │
               ▼
     []*adr.Record (In-Memory Structs)
               │
               ▼
      [ adr.Linter ]  ── (Monotonic Sequence Check: 1..N)
               │
               ├─ Parallel Validation (Goroutines)
               │   ├── Validasi reciprocal SupersededBy -> Supersedes
               │   └── Validasi reciprocal Supersedes -> SupersededBy
               ▼
     [ Validation Result: Slice of Errors or Nil ]
```

### Komponen Utama:
- `internal/adr/models.go`: Mendefinisikan tipe data `Status`, daftar status yang valid, dan struktur data `Record`.
- `internal/adr/parser.go`: Bertanggung jawab memindai teks berkas Markdown baris demi baris, mengekstrak nomor urut, judul, status siklus hidup, dan tautan penggantian.
- `internal/adr/linter.go`: Mesin validasi berbasis memori yang memeriksa keutuhan urutan penomoran dan memverifikasi integritas referensial timbal-balik secara konkuren.
- `cmd/demo/main.go`: Skenario percontohan yang memuat dan memvalidasi riwayat evolusi arsitektur berbasis ADR.

## Implementation

Implementasi validator ditulis dalam bahasa Go menggunakan pustaka standar tanpa dependensi pihak ketiga (`zero-dependency`).

Struktur data representasi ADR dimodelkan pada `internal/adr/models.go`:

```go
type Status string

const (
	StatusProposed   Status = "Proposed"
	StatusAccepted   Status = "Accepted"
	StatusSuperseded Status = "Superseded"
	StatusDeprecated Status = "Deprecated"
	StatusRejected   Status = "Rejected"
)

type Record struct {
	ID           int
	Title        string
	Status       Status
	SupersededBy int // 0 jika tidak digantikan
	Supersedes   int // 0 jika tidak menggantikan
	Content      string
}
```

Metode `IsValid()` pada tipe `Status` memastikan bahwa hanya status yang terdaftar secara resmi yang dapat diproses oleh parser.

Logika parsing pada `internal/adr/parser.go` menggunakan ekspresi reguler berikut:
- Judul: `^#\s+(\d+)\.\s+(.+)$`
- Status: `(?i)^Status:\s*([A-Za-z]+)(?:\s+by\s+(\d+))?`
- Supersedes: `(?i)^Supersedes:\s*(\d+)`

Logika linter pada `internal/adr/linter.go` memverifikasi integritas graf referensi keputusan dengan melakukan pemetaan seluruh catatan ke dalam peta memori `map[int]*Record`, memeriksa urutan monotonik, lalu memvalidasi tautan timbal balik menggunakan konkurensi Goroutine.

## Code Walkthrough

Berikut adalah penelusuran implementasi validasi pada `internal/adr/linter.go`:

### 1. Pemeriksaan Penomoran Monotonik
Linter memastikan bahwa penomoran ADR berurutan tanpa celah:

```go
// Validate IDs are monotonic
ids := make([]int, 0, len(records))
for id := range recordMap {
	ids = append(ids, id)
}
sort.Ints(ids)

for i := 0; i < len(ids); i++ {
	if ids[i] != i+1 {
		errs = append(errs, fmt.Errorf("non-monotonic numbering, expected %d but got %d", i+1, ids[i]))
		break // Only report once
	}
}
```

Jika terdapat berkas bernomor `1` dan `3` tanpa keberadaan ADR `2`, linter akan menangkap kegagalan `non-monotonic numbering, expected 2 but got 3`.

### 2. Validasi Konkuren Tautan Penggantian Timbal-Balik
Validasi hubungan silsilah dilakukan di dalam perulangan konkuren:

```go
if rec.Status == StatusSuperseded {
	if rec.SupersededBy == 0 {
		localErrs = append(localErrs, fmt.Errorf("ADR %d is superseded but missing superseded_by reference", rec.ID))
	} else {
		replacement, exists := recordMap[rec.SupersededBy]
		if !exists {
			localErrs = append(localErrs, fmt.Errorf("ADR %d superseded by non-existent ADR %d", rec.ID, rec.SupersededBy))
		} else if replacement.Supersedes != rec.ID {
			localErrs = append(localErrs, fmt.Errorf("ADR %d superseded by ADR %d, but ADR %d does not declare it supersedes ADR %d", rec.ID, rec.SupersededBy, rec.SupersededBy, rec.ID))
		}
	}
}

if rec.Supersedes != 0 {
	old, exists := recordMap[rec.Supersedes]
	if !exists {
		localErrs = append(localErrs, fmt.Errorf("ADR %d supersedes non-existent ADR %d", rec.ID, rec.Supersedes))
	} else if old.Status != StatusSuperseded || old.SupersededBy != rec.ID {
		localErrs = append(localErrs, fmt.Errorf("ADR %d supersedes ADR %d, but ADR %d is not properly marked as superseded by ADR %d", rec.ID, rec.Supersedes, rec.Supersedes, rec.ID))
	}
}
```

Penggunaan `sync.Mutex` menjamin agregasi galat ke dalam iris (*slice*) `errs` bersifat aman terhadap *race condition*.

## What the Tests Prove

Rangkaian pengujian otomatis pada `tests/parser_test.go` dan `tests/linter_test.go` membuktikan invarian sistem berikut:

1. **Parsing Format Positif dan Negatif (`tests/parser_test.go`):**
   - Berhasil mem-parsing judul, nomor ID, dan status valid (`Accepted`, `Superseded`, `Rejected`).
   - Berhasil mendeteksi klausa `Superseded by <ID>` dan `Supersedes: <ID>`.
   - Mengembalikan galat secara tepat ketika format judul hilang (`missing title`), status tidak dicantumkan (`missing status`), atau nilai status berada di luar skema valid (`invalid status: Draft`).

2. **Validitas Rangkaian Positif (`TestLinter_ValidSequence`):**
   - Membuktikan bahwa rangkaian ADR yang sah (ADR 1 digantikan oleh ADR 2, dan ADR 3 yang ditolak) lolos validasi tanpa galat.

3. **Deteksi Kerusakan Integritas Referensial (`TestLinter_BrokenReferences`):**
   - Membuktikan penolakan ketika ADR yang digantikan menunjuk ke target yang tidak ada (`superseded by non-existent ADR 99`).
   - Membuktikan penolakan ketika ADR baru mengklaim menggantikan ADR yang tidak ada (`supersedes non-existent ADR 99`).
   - Membuktikan penolakan pada tautan sepihak di mana ADR baru lupa mencantumkan klausa `Supersedes: 1` padahal ADR lama telah ditandai `Superseded by 2`.
   - Membuktikan penolakan pada urutan penomoran yang tidak monotonik (`non-monotonic numbering, expected 2 but got 3`).

4. **Keamanan Konkurensi:**
   - Eksekusi pengujian dengan bendera `-race` (`go test -race ./...`) lulus tanpa peringatan *data race*.

## Recovery / Rollback

Jika peninjauan CI mendeteksi kegagalan linter ADR:

1. **Kasus Tautan Sepihak:**
   - Jika PR baru menambahkan ADR pengganti (misalnya ADR 2), periksa apakah berkas ADR lama (ADR 1) telah diperbarui atribut statusnya menjadi `Status: Superseded by 2`. Perbarui berkas ADR 1 pada cabang kerja yang sama.
2. **Kasus Penomoran Loncat / Konflik:**
   - Jika dua pengembang membuat ADR dengan nomor urut yang sama secara bersamaan pada cabang terpisah, linter akan menandai adanya duplikasi atau ketidaksesuaian urutan saat proses integrasi. Pengembang dengan PR yang masuk belakangan harus mengubah nomor urut ADR miliknya ke nomor urut berikutnya yang tersedia.
3. **Pembatalan Keputusan Arsitektur:**
   - Jangan pernah menghapus berkas ADR yang telah berstatus *Accepted* dari repositori Git.
   - Buat ADR baru dengan nomor urut berikutnya yang menjelaskan konteks pembalikan keputusan tersebut, lalu tandai ADR lama sebagai *Superseded* atau *Deprecated*.

## Production Considerations

Penerapan ADR dan validasinya dalam skala tim produksi memerlukan beberapa pertimbangan:

1. **Integrasi CI/CD dan Pre-commit Hook:**
   - Validasi linter sebaiknya dijalankan secara otomatis pada setiap *pull request* yang mengubah direktori dokumentasi arsitektur.
   - Dapat dipasang pada *git pre-commit hook* lokal untuk memberi umpan balik instan kepada pengembang sebelum perubahan di-commit.
2. **Keterbatasan Format Parser:**
   - Implementasi parser pada lab ini bergantung pada pola regex yang ketat (format `# 1. Judul` dan `Status: ...`). Dalam produksi, tim harus mendokumentasikan templat penulisan ADR secara eksplisit atau memanfaatkan parser AST Markdown lengkap jika mendukung format penulisan yang lebih fleksibel.
3. **Repositori Majemuk (*Poly-repo vs Multi-repo*):**
   - Lab ini mengasumsikan ADR disimpan dalam satu repositori yang sama dengan kode (*monorepo* atau *single service repo*).
   - *Catatan Audit:* Pendekatan pengelolaan keputusan arsitektur lintas banyak repositori (*multi-repo*) belum dicakup dalam cakupan lab ini dan memerlukan repositori arsip pusat atau mekanisme sinkronisasi terpisah.
4. **Validasi Semantik vs Struktural:**
   - Linter otomatis hanya dapat memverifikasi integritas struktural (keberadaan berkas, tautan referensi, dan kata kunci status). Linter tidak dapat memvalidasi kualitas penalaran teknis, keakuratan data, atau kelayakan trade-off arsitektural. Penilaian konten tetap membutuhkan evaluasi manusia melalui *peer review*.

## Common Mistakes

Berdasarkan audit riset dan rekayasa, berikut adalah kekeliruan umum dalam penerapan ADR:

1. **Menyunting ADR yang Telah Disetujui Secara Retroaktif:** Mengubah isi keputusan pada ADR lama tanpa membuat ADR baru merusak catatan sejarah dan memanipulasi konteks masa lalu.
2. **Mencatat Detail Implementasi Rendah:** Menuliskan nama variabel, skema tabel privat, atau detail pustaka mikro yang mudah berubah ke dalam ADR. ADR harus dikhususkan untuk keputusan yang signifikan secara arsitektural (*ASR*).
3. **Menghilangkan Konsekuensi Negatif:** Menuliskan keputusan hanya dari sudut pandang keuntungan (*positive consequences*) tanpa mencatat kerugian atau beban operasional (*negative consequences*). Setiap keputusan arsitektur selalu memiliki kompromi (*trade-off*).
4. **Siklus Silsilah Tertutup (*Self-supersession / Cyclic Supersession*):** Mengarahkan tautan penggantian kembali ke nomor itu sendiri atau ke nomor yang lebih rendah. (*Catatan Audit: Implementasi lab saat ini belum membatasi arah temporal bahwa `SupersededBy` harus bernilai lebih besar dari `ID`*).
5. **Menghapus Usulan yang Ditolak:** Menghapus berkas usulan yang ditolak dari riwayat. Menyimpan berkas dengan status `Rejected` mencegah tim mengulang perdebatan yang sama tanpa argumen atau data baru.

## Case Study

Kasus percontohan yang dimodelkan pada `cmd/demo/main.go` mengilustrasikan evolusi arsitektur pada sistem SaaS ERP:

> **Peringatan Audit Kontekstual:** Skenario kasus studi ini merupakan contoh penerapan kontekstual untuk menguji alur siklus hidup ADR, bukan rekomendasi mutlak atau tolok ukur (*benchmark*) industri yang berlaku untuk semua proyek.

1. **ADR 0001: Penggunaan Modular Monolith untuk Inti SaaS ERP**
   - **Konteks:** Tim tahap awal (5 teknisi) dengan target peluncuran 3 bulan. Batasan domain bisnis masih berkembang pesat.
   - **Keputusan:** Membangun sistem sebagai *Modular Monolith* dengan isolasi paket domain yang ketat dan antarmuka *in-process*.
   - **Trade-off:**
     - *Positif:* Iterasi cepat, penerapan terpadu (*unified deployment*), integritas transaksi ACID lokal, tanpa latensi jaringan antardomain.
     - *Negatif:* Artefak penggelaran tunggal, basis data bersama (*shared database*).
   - **Pemicu Tinjauan (*Review Triggers*):** Ukuran tim melampaui 20 teknisi atau kebutuhan frekuensi rilis independen untuk domain dengan beban lonjakan tinggi.

2. **ADR 0002: Ekstraksi Layanan Notifikasi Menjadi Microservice**
   - **Konteks:** Tim berkembang menjadi 25 teknisi. Mesin notifikasi mengalami beban lalu lintas tinggi yang mulai memengaruhi performa transaksi inti ERP.
   - **Keputusan:** Mengubah status ADR 0001 menjadi `Superseded by 2`, dan menerbitkan ADR 0002 dengan status `Accepted` yang menyatakan `Supersedes: 1`. Modul notifikasi diekstraksi menjadi layanan independen berbasis antrean pesan asinkron.
   - **Trade-off:**
     - *Positif:* Penskalaan dan penerapan independen, isolasi domain kegagalan dari inti ERP.
     - *Negatif:* Beban operasional sistem terdistribusi, konsistensi eventual (*eventual consistency*), peningkatan kompleksitas observabilitas.

3. **ADR 0003: Penolakan Event Sourcing untuk Manajemen Pesanan**
   - **Konteks:** Usulan penulisan ulang modul pesanan menggunakan pola *Event Sourcing* penuh demi kebutuhan audit historis.
   - **Keputusan:** Menetapkan status `Rejected`. Kebutuhan audit hukum dipenuhi melalui tabel catatan audit relasional standar tanpa menambah beban migrasi skema kejadian (*event schema*).
   - **Trade-off:**
     - *Positif:* Mempertahankan kueri relasional sederhana dan semantik transaksi basis data yang sudah teruji.
     - *Negatif:* Logika pencatatan audit tambahan perlu dipelihara di lapisan aplikasi.

Hasil eksekusi program percontohan membuktikan bahwa silsilah keputusan ini diurai secara lengkap dan diverifikasi tanpa galat integritas.

## Checklist

Sebelum menggabungkan (*merging*) ADR ke cabang utama:

- [ ] ADR diberi nomor urut berikutnya secara kontigu tanpa melompati nomor.
- [ ] Berkas disimpan di direktori yang sama dengan kode dalam format Markdown.
- [ ] Bagian konteks menjelaskan batasan dan gaya dorong (*forces*) saat keputusan dibuat.
- [ ] Bagian konsekuensi mencantumkan dampak positif, negatif, dan netral secara objektif.
- [ ] Pemicu peninjauan (*review triggers*) terdefinisi untuk menentukan kapan keputusan harus dievaluasi kembali.
- [ ] Jika menggantikan keputusan lama, status ADR lama telah diubah menjadi `Superseded by <ID>` dan ADR baru mencantumkan `Supersedes: <ID>`.
- [ ] Jika usulan ditolak, status dicatat sebagai `Rejected` beserta alasan penolakannya.
- [ ] Linter struktural dijalankan dan dinyatakan lulus tanpa peringatan referensi putus.

## Key Takeaways

1. **ADR Menjaga Konteks Teknis:** ADR mendokumentasikan motif, kendala, dan kompromi dari keputusan arsitektur, mencegah perdebatan berulang dan keputusan sembarangan.
2. **Immutability Menjamin Keaslian Sejarah:** Keputusan yang telah diterima tidak disunting; perubahan selalu menghasilkan rekaman baru dengan silsilah penggantian yang jelas.
3. **Validasi Struktural Mencegah Tautan Patah:** Validasi otomatis menjamin urutan penomoran monotonik dan keutuhan referensi timbal-balik (`Supersedes` dan `Superseded by`).
4. **Skenario Nyata Adalah Kontekstual:** Keputusan arsitektur (seperti Monolith vs Microservices) bergantung pada kendala tim dan produk pada titik waktu tertentu, bukan aturan universal.
5. **Rekam Juga Penolakan:** Status `Rejected` sama pentingnya dengan `Accepted` untuk mendokumentasikan opsi yang telah dianalisis dan dinyatakan tidak sesuai.

## Sources

1. **Documenting Architecture Decisions**
   - Penulis: Michael Nygard
   - Publikasi: Cognitect Blog (2011-11-15)
   - URL: https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions
2. **Architectural Decision Record Process**
   - Publikasi: AWS Prescriptive Guidance, Amazon Web Services (2023-04-14)
   - URL: https://docs.aws.amazon.com/prescriptive-guidance/latest/architectural-decision-records/adr-process.html
3. **Architectural Decision Records (ADRs)**
   - Publikasi: adr.github.io (2024-11-10)
   - URL: https://adr.github.io/
