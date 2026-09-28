# Content Audit — Timeouts and Deadlines

Target Lab: `labs/28-timeouts-and-deadlines`
Audit Scope: Content only (research/code not audited per pipeline override)
Audit Date: 2026-09-28

## Files Reviewed

- `content/01-content-brief.md`
- `content/02-master-draft.md`
- `content/03-code-snippets.md`
- `content/04-diagrams.md`
- `content/05-key-takeaways.md`
- `content/06-source-map.md`

## Cross-Checks Performed

- Code snippets vs `internal/deadline/deadline.go`, `internal/retry/retry.go`, `internal/circuit/circuit.go`, `internal/idempotency/idempotency.go`, `cmd/demo/main.go`
- Test claims vs `*_test.go` and `tests/integration_test.go`
- Demo claims vs `engineering/03-execution-result.md`
- Verdict/status claims vs `research-audit/07-verdict.md`, `engineering-audit/06-verdict.md`
- Factual claims vs `research/03-evidence.md`, `research/06-open-questions.md`

## Verdict of Preceding Gates (Verified Accurate)

- Research: APPROVED, 0 unsupported claims, 3 LOW gaps — matches `research-audit/07-verdict.md`
- Engineering: APPROVED, 0 failures, race PASS, demo PASS — matches `engineering-audit/06-verdict.md`

---

## Issues Found

### BLOCKING

**B1. Diagram: salah menampilkan error parent deadline**
File: `content/04-diagrams.md:14`
```
- parent timeout (500ms tercapai) → context.Canceled
```
SALAH. Deadline parent yang kedaluwarsa menghasilkan `context.DeadlineExceeded`, bukan `context.Canceled`. `context.Canceled` hanya untuk `cancel()` eksplisit. Ini kontradiksi internal: draft sendiri menyatakan `TestExecuteWithBudget_ParentTimeoutInherited` menghasilkan `context.DeadlineExceeded`, dan test memang meng-assert `context.DeadlineExceeded`. Diagram mengajarkan hal yang salah ke pembaca.

**B2. Teks terpotong (truncated) pada Common Mistakes**
File: `content/02-master-draft.md:186`
```
**1. Timeout angka besar "supaya ama**
```
Judul terpotong di tengah kata (harusnya `"supaya aman."`). Kalimat tidak utuh — defect formatting yang merusak keterbacaan bagian utama artikel.

**B3. Klaim kuantitatif "lebih dari 50%" tidak didukung sumber**
File: `content/02-master-draft.md:228`, `content/05-key-takeaways.md:5`
```
Full Jitter ... menurunkan server contention lebih dari 50% dibanding unjittered.
```
Sumber AWS (`research/03-evidence.md` Evidence 4) hanya menyatakan "substantial decrease in client work and server load" — TANPA angka 50%. Angka spesifik ini adalah hallusinasi kuantitatif. B1-B3 melanggar checklist draft sendiri ("Tidak ada benchmark palsu", "Semua klaim faktual berasal dari sumber yang teridentifikasi").

### NON-BLOCKING

**N1. Salah ketik: "percayaan"**
File: `content/02-master-draft.md:84`
```
3 attempt dijalankan, 2 percayaan transient error
```
Harusnya "percobaan" (atau "attempt"). Kata "percayaan" tidak bermakna di konteks ini.

**N2. Penjelasan Snippet 3 melebih-lebihkan perilaku HALF_OPEN**
File: `content/03-code-snippets.md:84`
```
State HALF_OPEN hanya menerima satu request uji.
```
Tidak didukung kode. `RecordSuccess` menghitung hingga `SuccessThreshold` (bisa >1); `Allow()` tidak membatasi jumlah request di HALF_OPEN. Test menggunakan `SuccessThreshold: 2`. Formulasi akurat: "cukup `SuccessThreshold` sukses untuk kembali CLOSED".

**N3. Formula backoff inkonsisten antar file**
- Brief (`01-content-brief.md:13`): `base * 2^attempt`
- Draft (`02-master-draft.md:58`): `baseBackoff × 2^(attempt-1)`
- Key takeaway (`05-key-takeaways.md:5`): `base×2^attempt`

Kode benar: `2^(attempt-1)`. Dua file menyatakan `2^attempt` — salah untuk attempt yang di-1-index-kan. Perbaiki ke `2^(attempt-1)` agar konsisten dengan kode dan draft.

**N4. Snippet 5 & 6: tujuan/demo spesifikasi tidak konsisten dengan label**
File: `content/03-code-snippets.md:119` — Snippet 5 "parent deadline lebih kecil dari budget child" benar (50 < 100), tetapi explanation baris 136 menyebut `time.After(80ms)` melebihi parent deadline; 80ms memang > 50ms, akurat. Tanpa perubahan.

**N5. Source map menunjuk artefak audit yang belum ada**
File: `content/06-source-map.md:113`
```
content-audit/09-verdict.md — REJECTED (no content to audit sebelum draft ini dibuat)
```
Status REJECTED ini basi; audit ini akan menulis verdict baru. Harus diperbarui setelah verdict terbit (catatan, bukan defect konten).

**N6. Ringkasan demo retry**
File: `content/02-master-draft.md:84` — menyatakan "3 attempt... attempt ke-3 sukses" — cocok dengan `engineering/03-execution-result.md` (Attempt #1-#3, result `<nil>`). PASS.

---

## Verified Accurate (no issues)

- Semua 7 code snippet identik dengan kode aktual (termasuk guard `attempt <= 0`, komentar Full Jitter, `errors.Join(ErrMaxRetriesExceeded, lastErr)`).
- Semua nama test yang disebut ada dan akurat: `TestExecuteWithBudget_*` (3), `TestRetrier_*` (5), `TestCircuitBreaker_*` (3), `TestStore_*` (3), `TestIntegration_*` (2).
- Klaim demo 1-4 cocok persis dengan output `engineering/03-execution-result.md`.
- Integrasi: "retry 4 attempt, threshold 2, cooldown 50ms, sisa attempt ditolak" — cocok dengan `tests/integration_test.go`.
- `TestIntegration_IdempotentRetry` → `actualExecutions == 1` — akurat.
- Little's Law arithmetic: 100 QPS @ 1s = 100 workers; 100/60 ≈ 1.6 QPS — akurat.
- Lazy eviction + `mu.Lock()` pada `Get` — akurat (mutasi map).
- Child deadline = min(parent, budget) — akurat.
- Database/worker guardrails ditandai jelas sebagai "konteks riset, tidak diimplementasikan" — tidak menyesatkan.
- Tidak ada platform-specific bias (klaim gRPC vs HTTP disampaikan netral, dengan caveat "belum standar tunggal").
- Tidak ada benchmark/inventaris angka lain yang dibuat-buat selain isu B3.
- Tidak ada incident cerita karangan.

## Summary

- Blocking: 3 (1 factual/diagram error, 1 formatting truncation, 1 hallucinated statistic)
- Non-blocking: 5
- Akurasi kode & test: tinggi; kelemahan utama pada diagram deadline, teks terpotong, dan angka "50%" yang tidak bersumber.

Verdict: NEEDS_REVISION
