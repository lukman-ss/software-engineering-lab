---
title: "Architecture Decision Record: Parser and Linter for Structural Validation of Architectural Decisions"
published: false
description: "Explore Architecture Decision Record in Go: mechanics, architecture, caveats, and testable invariants."
tags: go, testing, performance, softwareengineering
---

# Architecture Decision Record: Parser and Linter for Structural Validation of Architectural Decisions

A production system must remain resilient under load. Lab 17 makes that behavior visible with a focused Go simulation and deterministic test assertions.

> Architectural context without structural records becomes tribal knowledge.

## Konteks Keputusan Arsitektur Hilang Tanpa Jejak

Putaran architect merangkap code review: pertanyaan berulang 'kenapa kita pilih X?' di every PR. Keputusan diambil via chat atau meeting, lalu tenggelam dalam commit history. ADR lama tidak diperbarui saat sistem berubah; status dokumentasi drift dari kode aktual. Re-debate arsitektur lama memakan waktu engineering sprint tanpa solusi baru. Keputusan yang salah diterapkan karena engineer tidak mengetahui ADR penggantian.

## Review Markdown Manual Lupa Relasi Penggantian

ADR 1 tetap berstatus 'Accepted' padahal sudah diganti ADR 2 — lupa di-update. ADR 2 menulis 'Supersedes: 1', tapi ADR 1 tidak pernah dicatat 'Superseded by 2'. Penomoran loncat: ADR 1, ADR 3 langsung — engineer lupa insert ADR 2. Status informal seperti 'Draft', 'In Review' diterima tanpa penolakan. ADR 1 ber-status Superseded by 99, padahal ADR 99 belum pernah dibuat.

## ADR Bukan Dokumen, Tapi Record yang Bisa Dicek

ADR diolah sebagai data, bukan sekadar teks human-readable. Parser mengekstrak ID, Title, Status, SupersededBy, Supersedes, Content. Linter memvalidasi aturan: monotonik ID, valid status, relasi timbal-balik. Keputusan arsitektur tetap ditulis manusiawi; mesin hanya memeriksa format.

## Pipeline: Markdown → Parser → Records → Linter

[Markdown ADR Text] → adr.Parser (Regex extraction: ID, Title, Status, refs) → []*adr.Record (In-Memory Structs) → adr.Linter.Validate(records)   ├─ Monotonic numbering check: sort + sequential scan   ├─ Cycle detection: DFS with visiting/visited states   └─ Bidirectional lineage: goroutines + sync.Mutex per-record validation → []error (empty = valid, non-empty = integrity breach)

## Struct Record: Definisi Inti ADR

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
    SupersededBy int // 0 if not superseded
    Supersedes   int // 0 if does not supersede
    Content      string
}
```

## Supersesi Patah: ADR 1 Merujuk ke ADR 99

TestLinter_BrokenReferences / superseded points to non-existent → PASS TestLinter_BrokenReferences / supersedes points to non-existent → PASS

```go
superseded points to non-existent:
  ADR 1: Status Superseded by 99
  → Error: superseded by non-existent ADR 99

supersedes points to non-existent:
  ADR 1: Status Accepted, Supersedes: 99
  → Error: supersedes non-existent ADR 99
```

## Biaya Tersembunyi: Lineage Drift Tanpa Deteksi

Time to detect | Cause | Impact
|---|---|---|
Instant (CI) | Linter reject di PR | Build gagal, engineer feedback loop cepat
Mingguan | Manual audit dokumentasi | Engineer kebingungan siapa yang aktif
Bulanan | Production incident review | Tim kehilangan konteks keputusan penting
Selamanya | Keputusan salah diterapkan | Sistem berkembang dengan asumsi keliru

## Setup Eksperimen: 3 ADR Evolusi Arsitektur

ADR 0001: Use Modular Monolith for Core SaaS ERP | Status: Superseded by 2 ADR 0002: Extract Notification Service to Microservice | Status: Accepted, Supersedes: 1 ADR 0003: Reject Event Sourcing for Order Management | Status: Rejected  Test suite: parser_test.go (valid + superseded + supersedes + invalid cases) Test suite: linter_test.go (valid sequence + 7 broken reference scenarios + concurrency)

## Output Demo: Parsing Berhasil, Linter Lolos

```sh
$ go run ./cmd/demo
=== Architecture Decision Record (ADR) Lab ===
Successfully parsed 3 ADRs:
  - [0001] Use Modular Monolith for Core SaaS ERP | Status: Superseded -> Superseded by ADR 2
  - [0002] Extract Notification Service to Microservice  | Status: Accepted   -> Supersedes ADR 1
  - [0003] Reject Event Sourcing for Order Management    | Status: Rejected  

Running ADR Integrity Linter...
Integrity check passed! Decision lineage and lifecycle invariants are intact.
```

## Apa Artinya Linter Passed di Lab Ini?

Format parsing: judul, status, dan tautan supersesi terstruktur dengan benar. Monotonik: ID 1, 2, 3 berurutan tanpa celah. Bidirectional: ADR 1 ↔ ADR 2 memiliki tautan timbal-balik yang lengkap. Cycle-free: tidak ada sirkuit dalam graf supersesi. Kebenaran keputusan arsitektur di dalam field Content. Apakah ADR 1 seharusnya benar-benar digantikan ADR 2 di dunia nyata. Kesiapan sistem untuk production deployment.

## Trade-off: Enforced Structure vs Dokumentasi Fleksibel

Regex-based parsing memaksa format fixed (## Context, ## Decision, ## Consequences). Status wajib dari 5 enum resmi; status informal ditolak parser. Supersedes/SupersededBy wajib konsisten secara timbal-balik. Tidak mendukung ADR yang memerlukan struktur sub-section kustom di luar 3 section wajib. Tidak memvalidasi kualitas argumen dalam section Decision. Tidak menangani kasus ADR multi-author atau decision yang melibatkan beberapa sub-system.

## Validasi Paralel: Goroutine + Mutex untuk Agregasi Error

```go
var wg sync.WaitGroup
for _, r := range records {
    wg.Add(1)
    go func(rec *Record) {
        defer wg.Done()
        var localErrs []error
        // ... validate rec, append localErrs ...
        if len(localErrs) > 0 {
            mu.Lock(); errs = append(errs, localErrs...); mu.Unlock()
        }
    }(r)
}
wg.Wait()
```

## Before vs After: Otomasi Linter Menghilangkan Human Error

Lineage break: ADR 1 Superseded by 2, ADR 2 tidak memiliki Supersedes: 1 → TIDAK TERDETEKSI Non-monotonic ID: ADR 1, ADR 3 meloncat → TIDAK TERDETEKSI selagi tidak inspect semua file Invalid status 'Draft' masuk ke repo → TIDAK TERDETEKSI oleh regex sederhana Self-supersession: ADR 1 Superseded by 1 → TIDAK TERDETEKSI Semua kasus di atas ditolak oleh TestLinter_BrokenReferences dalam <1ms CI gate menolak merge jika linter menemukan error Concurrent stress test membuktikan tidak ada race pada agregasi error

## Batasan Lab: Validasi Format, Bukan Kualitas Keputusan

Parsing teks Markdown ke struct Record dengan regex. Memvalidasi urutan ID monotonik dan uniqueness. Memeriksa konsistensi timbal-balik SupersededBy ↔ Supersedes. Mendeteksi siklus dalam graf supersesi. Tidak menilai kualitas argumen di dalam section Decision. Tidak memverifikasi apakah ADR cocok dengan implementasi kode aktual. Tidak memeriksa apakah engineer sudah mendapatkan approval stakeholder sebelum Accepted. Simulasi ini tidak sama dengan production ADR management system seperti ADRtool atau Architect.js.

## Edge Case: Keputusan Dibatalkan Tanpa ADR Baru

ADR X berstatus Accepted, tapi tim memutuskan berhenti menggunakan fitur tersebut. Tanpa ADR baru, status ADR X tetap Accepted —误导 engineer lain以为 masih relevan. Solusi: buat ADR baru berstatus Deprecated merujuk ke alasan penghentian. Kegagalan: engineer menghapus file ADR langsung → lineage history hilang. Wajibkan policy tim: setiap penghentian fitur memerlukan ADR Deprecated baru. Linter bisa diperluas untuk mendeteksi ADR Accepted yang tidak memiliki successor dalam 12 bulan.

## Apa yang Dibuktikan Oleh Test Suite

TestParse_Valid: ID=1, Status=Accepted → PASS TestParse_Superseded: SupersededBy=2 berhasil diekstrak → PASS TestParse_Supersedes: Supersedes=1 berhasil diekstrak → PASS TestParse_Invalid: missing_title, missing_status, invalid_status, empty sections → semua PASS TestLinter_ValidSequence: 3 ADR valid, zero errors → PASS TestLinter_BrokenReferences: 7 sub-case (non-existent ref, mismatch, non-monotonic, self-super, duplicate ID, cyclical) → semua PASS TestLinter_ConcurrencyStress: concurrent validation, no data race → PASS

## Apa yang Tidak Dibuktikan Oleh Test Suite

Kebenaran konten arsitektur di dalam section Decision dan Consequences. Adopsi ADR oleh seluruh anggota tim sesuai format yang ditetapkan. Integrasi linter dengan CI/CD pipeline (hanya test Go murni, bukan GitHub Actions). Skalabilitas linter pada ribuan ADR (hanya 3 record dalam demo, 3-7 dalam test). Parser bekerja pada string literal, bukan file system read dari ribuan dokumen Markdown.

## Risiko Produksi: False Confidence dari Linter Hijau

Tim percaya ADR collection-nya 'bersih' hanya karena linter passing, padahal argumen Decision lemah. Engineer bisa mengisi Content dengan placeholder kosong selama section ada → linter tetap passing. Policy tim yang tidak enforcing bisa di-skip tanpa konsekuensi jika gate CI tidak dipasang. Tambah human review checkpoint: ADR Accepted wajib ditinjau architecture board. Validasi minimum content length pada section Context dan Decision.

## Implementation Checklist: Menerapkan ADR Linter di Tim

[ ] Define ADR template Markdown dengan 3 section wajib: Context, Decision, Consequences. [ ] Implement parser regex yang mengekstrak ID, Title, Status, SupersededBy, Supersedes. [ ] Tambahkan status enum valid (Proposed, Accepted, Superseded, Deprecated, Rejected). [ ] Implement linter dengan validasi monotonik + bidirectional supersession + cycle detection. [ ] Integrasikan ke CI/CD: go test ./... atau custom binary di PR gate. [ ] Dokumentasikan policy tim: setiap perubahan keputusan wajib ADR baru, bukan edit retroaktif.

## Anti-Pattern: 5 Kebiasaan Buruk Saat Mengadopsi ADR

Mengklik ADR untuk segala hal: refactor kecil, penamaan variabel → noise berlebihan. Mengedit ADR yang sudah Accepted secara silent, tanpa membuat ADR baru (retroactive edit). Menulis ADR hanya sebagai catatan meeting tanpa section Decision yang eksplisit. Status 'Draft' bertahan terlalu lama tanpa transisi ke Accepted atau Rejected. Melupakan update SupersededBy pada ADR lama saat ADR baru dibuat.

## Reproduksi Mandiri: Jalankan Demo dan Test di Lokal

TestParse_*: semua PASS TestLinter_*: semua PASS (14 sub-case total) go test -race: race detector clean Demo: 3 ADR parsed, Integrity check passed

```sh
cd labs/17-architecture-decision-record
go test -v ./...
go test -race ./...
go run ./cmd/demo
```

## ADR Diubah Jadi Data yang Bisa Dicek, Bukan Hanya Dokumen

ADR yang diparse menjadi struct Record memungkinkan validasi otomatis yang tidak bisa dilakukan humans dengan andal. Bidirectional supersession link adalah invarian yang paling sering rusak; linter menangkapnya di detik pertama. Konkurensi goroutine + mutex pattern untuk error aggregation adalah reusable pattern untuk validation framework lain. Format standar membantu review; kualitas keputusan tetap ditentukan isi ADR.

## Source

[Software Engineering Lab 17 - Architecture Decision Record](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/17-architecture-decision-record), snapshot commit `def0f5a`. Source: README, implementation, tests, and execution records. Code excerpts belong in their original package context.
