# Safe Deployment Sequence & Rollback

## Urutan Deployment yang Aman
1. Deploy schema migration baru (Additive only, e.g. add new column).
2. Deploy App version N+1 (Menulis ke dual/struktur baru).
3. Jalankan Backfill script di background.
4. Pindahkan Read ke kolom/tabel baru (App N+2).
5. Observability monitoring (Tunggu legacy idle).
6. Hapus akses legacy (App N+3).
7. Drop kolom lama di Database.

## Rolling Deployment & Rollback
- App N (lama) dan App N+1 (baru) akan hidup bersamaan saat rolling deployment. Karenanya database **wajib** kompatibel untuk kedua versi.
- Jika App N+1 bug dan di-rollback ke App N, App N harus tetap bisa berjalan dengan schema database baru. Inilah mengapa schema migration bersifat "Expand" (Additive) terlebih dahulu.
