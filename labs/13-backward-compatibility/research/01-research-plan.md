# Research Plan

## Research Topic
Backward Compatibility — Cara Mengembangkan Sistem Lama Tanpa Merusak Production

## Objective
Riset bagaimana engineer mengembangkan sistem lama yang sudah berjalan di production tanpa merusak client, service, database consumer, API consumer, worker, mobile app, reporting, dan integration yang masih bergantung pada contract lama.

## Research Questions
1. Apa itu backward compatibility?
2. Apa bedanya dengan forward compatibility?
3. Apa yang membuat sebuah perubahan menjadi breaking change?
4. Bagaimana Expand → Migrate → Contract bekerja?
5. Bagaimana melakukan schema migration tanpa downtime?
6. Kapan dual read dibutuhkan?
7. Apa risiko dual write?
8. Bagaimana melakukan backfill secara aman?
9. Bagaimana menjaga compatibility saat rolling deployment?
10. Bagaimana mengetahui consumer lama masih aktif?
11. Kapan field, column, atau endpoint lama boleh dihapus?
12. Bagaimana observability membantu migration?
13. Bagaimana feature flag membantu rollout dan rollback?
14. Bagaimana rollback memengaruhi desain migration?
15. Apa failure mode paling umum dalam backward-compatible migration?
16. Bagaimana membuat migration resumable dan idempotent?
17. Apa yang harus diuji sebelum contract lama dihapus?
18. Apa perbedaan database compatibility dan API compatibility?

## Search Strategy
Mencari sumber terpercaya tier-1 dan tier-2 (Martin Fowler, dokumentasi Stripe) untuk mengumpulkan best practice dalam menangani API versioning, database schema evolution, backward compatibility, dan expand/contract pattern.

## Expected Primary Sources
- Martin Fowler's "Evolutionary Database Design"
- Martin Fowler's "Parallel Change"
- Stripe API Versioning Blog Post
- Official documentation vendor-specific untuk zero-downtime database migration

## Risks / Unknowns
- Kinerja dual write pada database berskala sangat besar.
- Waktu yang dibutuhkan untuk menghentikan consumer lama.
- Perbedaan dukungan schema changes antara PostgreSQL, MySQL, dan database lainnya (NoSQL vs SQL).
