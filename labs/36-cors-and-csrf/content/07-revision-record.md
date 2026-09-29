# Revision Record

Tanggal: 2026-09-29
Lab: `labs/36-cors-and-csrf`
Target: Memperbaiki temuan Content Audit (`content-audit/report.md`).

## 1. Temuan Audit yang Ditangani

1. **Pipeline Misrepresentation**:
   - *Masalah*: Dokumentasi menggambarkan pipeline endpoint terproteksi (`/api/transfer/protected`) melalui rantai CORS → Fetch-Metadata → Custom-Header → CSRF, padahal dalam implementasi aktual endpoint tersebut hanya dibungkus CORS dan CSRF middleware. Middleware Fetch-Metadata dan Custom-Header diekspos pada rute terpisah (`/api/transfer/fetch-metadata` dan `/api/transfer/custom-header`).
   - *Tindakan*: Memperbaiki teks dan diagram arsitektur di `02-master-draft.md` dan `04-diagrams.md`. Menjelaskan secara akurat bahwa alur `/api/transfer/protected` menggunakan CORS + CSRF Token Validator, sementara Fetch Metadata dan Custom Header merupakan middleware modular pendukung yang diuji pada endpoint terpisah.

2. **Layered Defense Claim Overstatement**:
   - *Masalah*: Klaim pertahanan multi-layer pada protected endpoint dilebih-lebihkan seolah-olah protected endpoint secara default dieksekusi dengan semua filter sekaligus.
   - *Tindakan*: Memperjelas deskripsi modularitas komponen pertahanan di `02-master-draft.md` dan `04-diagrams.md`. Menegaskan bahwa pertahanan berlapis adalah konsep arsitektural yang dapat dikomposisikan, sedangkan lab menguji komponen-komponen tersebut secara terisolasi dan spesifik pada endpoint masing-masing.

3. **Token Delimiter Caveat Omitted (GAP-02)**:
   - *Masalah*: Dokumentasi belum mencatat limitasi delimiter titik dua (`:`) pada parsing payload token jika session ID memuat karakter `:`.
   - *Tindakan*: Menambahkan poin peringatan teknis di `02-master-draft.md` (seksi *Production Considerations*, poin 3) dan memperbarui seksi *Warnings* pada `01-content-brief.md`.

## 2. File yang Diubah
- `labs/36-cors-and-csrf/content/02-master-draft.md`
- `labs/36-cors-and-csrf/content/04-diagrams.md`
- `labs/36-cors-and-csrf/content/01-content-brief.md`
- `labs/36-cors-and-csrf/content/07-revision-record.md` (file ini)
