---
title: "Architecture Decision Record: Parser and Linter for Structural Validation of Architectural Decisions"
published: false
description: "Explore Architecture Decision Record in Go: mechanics, architecture, caveats, and testable invariants."
tags: go, testing, performance, softwareengineering
---

# Architecture Decision Record: Parser and Linter for Structural Validation of Architectural Decisions

A production system must remain resilient under load. Lab 17 makes that behavior visible with a focused Go simulation and deterministic test assertions.

> Architectural context without structural records becomes tribal knowledge.

## Konteks Arsitektur Hilang Tanpa Jejak

Keputusan arsitektur diambil via chat atau meeting tanpa catatan formal. Pertanyaan berulang muncul di review kode: mengapa sistem dibangun begini? ADR lama tidak diperbarui, sehingga status dokumen berbeda dari kode aktual.

## Review Manual Melewatkan Relasi Supersesi

ADR baru menyatakan menggantikan ADR lama, tetapi ADR lama tidak diubah. ADR lama tetap bertanda aktif padahal keputusannya sudah ditinggalkan. Tautan referensi menunjuk ke nomor yang salah tanpa disadari.

## ADR Bukan Dokumen, Tapi Record

ADR diperlakukan sebagai data terstruktur, bukan sekadar teks bebas. Field kunci diekstrak: ID, Judul, Status, SupersededBy, dan Supersedes. Mesin memvalidasi konsistensi; manusia tetap menulis esensi keputusannya.

## Pipeline Validasi Empat Tahap

Markdown Files | Format teks mentah adr.Parser | Ekstraksi struct Record adr.Linter | Aturan & integritas graf Laporan Error | Deteksi inkonsistensi

## Struct Record: Fondasi Validasi Data

```go
type Record struct {
    ID           int
    Title        string
    Status       Status
    SupersededBy int // 0 jika aktif
    Supersedes   int // 0 jika tidak menggantikan
    Content      string
}
```

## Deteksi Tautan Supersesi Patah

```go
ADR 1: Status Superseded by 99
→ Error: superseded by non-existent ADR 99

ADR 1: Status Accepted, Supersedes: 99
→ Error: supersedes non-existent ADR 99
```

## Biaya Tersembunyi: Lineage Drift Menumpuk

Tahap | Deteksi | Dampak
|---|---|---|
CI Gate | Seketika | PR ditolak, feedback loop cepat
Review Manual | Mingguan | Tim bingung dokumen mana yang aktif
Insiden Sistem | Berbulan-bulan | Keputusan baru berpijak pada asumsi keliru

## Skenario Demo: Tiga ADR Evolusi Arsitektur

ADR 0001: Modular Monolith (Superseded by ADR 2) ADR 0002: Notification Microservice (Accepted • Supersedes ADR 1) ADR 0003: Event Sourcing (Rejected • Biaya operasional tinggi)

## Output Demo: Validasi Integritas Lolos

```sh
$ go run ./cmd/demo
Successfully parsed 3 ADRs:
  - [0001] Modular Monolith -> Superseded by ADR 2
  - [0002] Notification Microservice -> Supersedes ADR 1
  - [0003] Event Sourcing -> Rejected
Running ADR Integrity Linter... Integrity check passed!
```

## Lolos Linter Bukan Berarti Keputusan Benar

Nomor ID berurutan tanpa celah. Relasi timbal-balik utuh tanpa tautan gantung. Tidak ada siklus ketergantungan antar-dokumen. Kebenaran argumen teknis di bagian Decision. Apakah rancangan arsitektur tepat untuk kapasitas tim. Evaluasi dampak performa dan biaya operasional sistem.

## Trade-off: Struktur Baku vs Fleksibilitas

Wajib tiga bagian: Context, Decision, Consequences. Status dibatasi lima enum resmi. Memungkinkan CI gate otomatis dan deterministik. Tidak mendukung format narasi paragraf bebas. Memerlukan disiplin tim mengisi field relasi. Perubahan format membutuhkan penyesuaian parser.

## Validasi Paralel dengan Mutex

```go
for _, r := range records {
    wg.Add(1)
    go func(rec *Record) {
        defer wg.Done()
        localErrs := validate(rec)
        mu.Lock()
        errs = append(errs, localErrs...)
        mu.Unlock()
    }(r)
}
```

## Sebelum vs Sesudah Validasi Otomatis

Tautan supersesi sepihak sering terlewat. Nomor meloncat sulit disadari saat banyak file. Status non-standar masuk tanpa sengaja. Semua inkonsistensi ditolak otomatis di CI. Nomor ID wajib monotonik dan teratur. Relasi timbal-balik divalidasi deterministik.

## Batasan Lab: Hanya Validasi Format

Ekstraksi regex dari string Markdown. Pemeriksaan nomor urut dan relasi timbal-balik. Deteksi siklus pada graf supersesi. Tidak menilai kebenaran argumen keputusan. Simulasi in-memory, bukan pembacaan ribuan berkas disk.

## Kasus Batas: Menghentikan Fitur Tanpa ADR

Fitur dimatikan tetapi ADR tetap bertanda Accepted. Engineer baru mengira fitur masih menjadi standar aktif. Menghapus berkas justru menghilangkan riwayat konteks. Solusi: wajibkan ADR baru berstatus Deprecated.

## Bukti yang Divalidasi oleh Test Suite

Parser menolak status ilegal dan format tanpa judul. Linter menangkap 7 skenario referensi patah dan tautan gantung. Race detector bersih pada eksekusi pengujian konkuren.

## Batas Integritas: Apa yang Tidak Diuji

Kebenaran konten teknis di dalam bagian Decision. Disiplin tim dalam memperbarui ADR saat arsitektur bergeser. Kesiapan sistem monolit atau microservice di production. Performa I/O disk pada ribuan berkas dokumen.

## Risiko Kepalsuan: Linter Hijau Bukan Aman

Dokumen lolos linter meski bagian Decision hanya teks formalitas. Tim mengira arsitektur sudah matang hanya karena status CI hijau. Mitigasi: pertahankan review arsitek untuk menilai substansi.

## Panduan Implementasi Linter ADR di Tim

Tetapkan template berkas: Context, Decision, Consequences. Gunakan enum status resmi (Accepted, Superseded, Rejected). Pasang parser dan linter pada CI gate pull request. Larang edit retroaktif: keputusan baru wajib nomor baru.

## Anti-Pattern: Kesalahan Adopsi ADR

Menulis ADR untuk hal sepele seperti penamaan variabel. Mengedit ADR Accepted diam-diam alih-alih membuat pengganti. Membiarkan status Draft berbulan-bulan tanpa kejelasan. Lupa mencatat nomor ADR pengganti saat membuat keputusan baru.

## Reproduksi Mandiri di Mesin Lokal

```sh
cd labs/17-architecture-decision-record
go test -v ./...    # validasi parser & linter
go test -race ./... # verifikasi race-free
go run ./cmd/demo   # jalankan simulasi 3 ADR
```

## ADR Diubah Jadi Data yang Bisa Dicek

ADR terstruktur memungkinkan verifikasi otomatis relasi dokumen. Linter menutup kelalaian mekanis yang sering luput dari review manusia. Format standar membantu review; kualitas keputusan tetap ditentukan isi ADR.

## Source

[Software Engineering Lab 17 - Architecture Decision Record](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/17-architecture-decision-record), snapshot commit `def0f5a`. Source: README, implementation, tests, and execution records. Code excerpts belong in their original package context.
