Target Lab: labs/18-deadlock

Files Reviewed:
- research/runs/2026-09-25-deadlock/01-plan.md
- research/runs/2026-09-25-deadlock/02-sources.md
- research/runs/2026-09-25-deadlock/03-evidence.md
- research/runs/2026-09-25-deadlock/04-contradictions.md
- research/runs/2026-09-25-deadlock/05-report.md
- research/runs/2026-09-25-deadlock/06-open-questions.md

Claims To Verify:
1. Deadlock adalah siklus dependensi (circular wait) permanen.
2. Penanganan deadlock otomatis melalui aborsi/rollback (Deadlock Victim).
3. Urutan akses (Lock Ordering) mencegah mayoritas deadlock.
4. Durasi transaksi yang panjang memperbesar probabilitas deadlock.
5. Deadlock dapat ditangani dengan mekanisme Retry di aplikasi.

Code To Execute:
None. PIPELINE OVERRIDE: Audit research only. Do not audit implementation/code in this stage.

Primary Risks:
- Generalisasi berlebihan (menganggap satu perilaku RDMS berlaku universal).
- Kesimpulan dari satu vendor (Microsoft / PostgreSQL).
- Asumsi tentang "pencegahan absolut".

Audit Strategy:
- Verifikasi setiap source memastikan kutipan ada dan relevan.
- Validasi claim berdasarkan source.
- Catat gap jika klaim melampaui referensi.
