# Database Migration Tanpa Downtime

## Expand → Migrate → Contract (Parallel Change)
Target: Database migration untuk tabel skala besar tanpa lock down time dan breaking changes.

### Transition Phase (Database Compatibility)
Evidence: Menggunakan View atau Triggers selama periode transisi (migrate phase). Database refactoring kecil dan tidak merusak harus menjadi prioritas.
Source: Evolutionary Database Design
URL: https://martinfowler.com/articles/evodb.html
Confidence: HIGH

### Dual Write / Dual Read
- **Dual Write Risiko**: Kegagalan pada satu tabel (lama atau baru) bisa menyebabkan in-consistency (membutuhkan distributed transaction atau idempotent background sync).
- **Data Backfill Aman**: Dijalankan secara batch processing kecil (throttling) menggunakan script idempotent dan resumable.
