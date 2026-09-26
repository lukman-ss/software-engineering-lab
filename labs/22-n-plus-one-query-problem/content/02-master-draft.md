# N+1 Query Problem: Analisis dan Solusi

## Problem
Aplikasi menjalankan 1 query utama untuk mengambil koleksi data, lalu N query tambahan untuk setiap data individual.

## Why This Matters
Query individu sangat cepat, meloloskannya dari *slow query logs*. Degradasi performa terjadi secara diam-diam akibat akumulasi latensi jaringan dan interaksi database berulang.

## Mental Model
Hindari iterasi query di dalam *loop*. *Batch* data, kumpulkan identitas relasi, dan minta seluruh relasi dalam satu permintaan besar.

## Core Concept
- **Naif (N+1):** Ambil *author*. Iterasi tiap *author*. Ambil *post* untuk tiap *author*.
- **Solusi (Eager Loading / Batching):** Ambil *author*. Kumpulkan ID *author*. Ambil seluruh *post* dengan klausa `IN (id)`. Pasangkan *post* dengan *author* di memori aplikasi.

## Architecture & Implementation
*Eager loading* level query (seperti menggunakan `IN`) memotong interaksi database. Cukup 2 query untuk merangkai relasi yang kompleks:
1. Ambil parent.
2. Ambil relasi anak terkait.

## What the Tests Prove
Pengujian membuktikan solusi bekerja secara konsisten:
- Pendekatan N+1 untuk 3 *author* mengeksekusi tepat 4 query.
- Pendekatan *Eager loading* mengeksekusi tepat 2 query (reduksi linier).

## Common Mistakes
Menyamakan *Query-level eager loading* (direkomendasikan) dengan *Mapping-level eager loading* (berbahaya). Menggunakan konfigurasi *hardcoded* seperti `FetchType.EAGER` secara global memicu *memory bloat* karena data relasional yang tidak relevan ikut ditarik tanpa henti.

## Warning / Caveats
Berdasarkan audit:
- Jangan gunakan *eager loading* level *mapping*. Gunakan level spesifik pada query.
- Solusi belum menguji ketat kesetaraan *field-by-field* (*deep equal*).
- Transisi status pada *zero records* harus diseragamkan (e.g., saat database kosong, *return nil* vs *empty slice* wajib dijaga konsistensinya).
