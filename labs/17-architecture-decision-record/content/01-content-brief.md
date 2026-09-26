# Content Brief

Topic: Architecture Decision Record (ADR) dan validasi struktural otomatis

Target Reader: Software architects, developers, technical leads

Problem: Bagaimana mendokumentasikan keputusan arsitektur yang signifikan secara konsisten dan terlacak, serta mencegah perdebatan teknis berulang akibat hilangnya konteks keputusan.

Core Mental Model: Keputusan direkam bersama konteks, konsekuensi, dan status lifecycle; keputusan lama tidak dihapus tetapi digantikan melalui silsilah timbal-balik (Superseded/Supersedes) yang dapat divalidasi secara otomatis.

Approved Research Status: APPROVED

Approved Engineering Status: APPROVED

Main Concepts:
- Batasan keputusan architecturally significant (structure, NFR, dependencies, interfaces, construction techniques)
- Siklus hidup ADR: Proposed, Accepted, Superseded, Deprecated, Rejected
- Immutability riwayat keputusan (tidak menyunting ADR Accepted secara retroaktif)
- Silsilah penggantian timbal-balik dan penomoran monotonik
- Ko-lokasi ADR dengan kode dalam satu repositori

Verified Behaviors:
- Parser mengekstrak ID, Title, Status, SupersededBy, Supersedes dari Markdown via regex (tests/parser_test.go)
- Parser menolak judul hilang, status hilang, status tidak valid, serta bagian Context/Decision/Consequences yang hilang atau kosong (TestParse_Invalid)
- Linter menolak referensi penggantian yang tidak ada, tautan sepihak, penomoran non-monotonik (TestLinter_BrokenReferences)
- Linter menolak self-supersession, duplikasi ID, dan siklus supersession (TestLinter_BrokenReferences)
- Validasi konkuren aman race (go test -race PASS, TestLinter_ConcurrencyStress 100 record)
- Demo cmd/demo lolos integrity check pada skenario 3 ADR (engineering/03-execution-result.md)

Available Case Studies:
- SaaS ERP evolution: Modular Monolith (ADR 1) → Microservices (ADR 2) → Rejected Event Sourcing (ADR 3), divalidasi end-to-end via cmd/demo/main.go

Warnings:
- Thesis "monolith-first" bergantung pada pengamatan anekdot Martin Fowler, bukan studi empiris (research-audit/07-verdict.md, research/05-report.md Limitations)
- Tidak ada data kuantitatif dampak ADR terhadap velocity/onboarding (research-audit/07-verdict.md)
- Literatur yang dikaji terfokus pada tim kecil-menengah; skala 100+ engineer belum tercakup (research/06-open-questions.md)
- Parser format sangat spesifik (`# 1. Title`, `Status: ...`), bukan parser Markdown umum (engineering/02-implementation-notes.md)
- Desain menyebut "Fake File System / In-Memory Repo" tetapi implementasi memakai string in-memory; keterbatasan terdokumentasi (engineering-audit-opensource/05-gaps.md, LOW)
- Uji case-insensitive status, input slice kosong, serta status Proposed/Deprecated tidak diuji eksplisit (engineering-audit-opensource/05-gaps.md, LOW)
- Lab tidak mendemonstrasikan integrasi git pre-commit hook maupun pembuatan berkas ADR (engineering/02-implementation-notes.md)
- Skenario kasus studi adalah contoh kontekstual, bukan rekomendasi universal atau benchmark
