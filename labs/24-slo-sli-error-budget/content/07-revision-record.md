# Revision Record — 2026-09-28

Target: `labs/24-slo-sli-error-budget` | Scope: `content/` only

## Audit Basis
`content-audit/01-audit-report.md` + `content-audit/09-verdict.md` — APPROVED_WITH_WARNINGS, 5 LOW non-blocking issues (NB-1..NB-5).

## Changes

### NB-2 — Failure Scenario 5% → 10%, clarified as hypothetical
- File: `02-master-draft.md:68-75`
- Before: `error rate 5%` / `5 error → SLI 95%`, no hypothetical label
- After: `## Skenario Kegagalan *(Hipotetik)*`, `error rate 10%` / `10 error → SLI 90%`
- Reason: Align hypothetical scenario with demo's actual 10% error rate (90 good + 10 bad), avoid reader confusion.

### NB-3 — Phase 4 criticality contrast
- File: `02-master-draft.md:297-304`
- Before: Both Payment 99.9% and Reports 95% shown as CanDeploy=false at 10% error, no differentiation
- After: Added illustration: at 2% error rate Payment exceeds budget (false) while Reports stays within 5% tolerance (true), showing why criticality matters.

### NB-4 — BurnRateRule unused fields in D3
- File: `04-diagrams.md: D1 + D3`
- Before: Unused fields noted only in prose (`02-master-draft.md:138`), diagram silent
- After: D3 box now includes callout `⚠ PERINGATAN: BurnRateRule.LongWindow, ShortWindow, BudgetConsumedPct TIDAK DIBACA oleh Check()`; D1 component list notes same.

### NB-5 — Section headers English → Indonesian
- File: `02-master-draft.md`
- Before: `Problem`, `Why This Matters`, `Mental Model`, `Core Concept`, `Failure Scenario`, `How It Works`, `Architecture`, `Error Budget Calculation`, `Multi-Window Alerting`, `Implementation`, `Code Walkthrough`, `Event Recording`, `Burn Rate Alert`, `What the Tests Prove`, `Unit Tests`, `Concurrency Test`, `Alert Engine Test`, `Recovery / Rollback (Hypothetical...)`, `Production Considerations`, `Common Mistakes`, `Case Study`, `Demo Output`, `Checklist`, `Key Takeaways`, `Sources`
- After: `Masalah`, `Mengapa Ini Penting`, `Model Mental`, `Konsep Inti`, `Skenario Kegagalan *(Hipotetik)*`, `Cara Kerja`, `Arsitektur`, `Perhitungan Error Budget`, `Alerting Multi-Jendela`, `Implementasi`, `Tur Kode`, `Pencatatan Event`, `Alert Burn Rate`, `Yang Dibuktikan oleh Tes`, `Tes Unit`, `Tes Konkurensi`, `Tes Alert Engine`, `Recovery / Rollback (Prosedur Hipotetis — TIDAK ditunjukkan di demo)`, `Pertimbangan Produksi`, `Kesalahan Umum`, `Studi Kasus`, `Output Demo (Ringkasan)`, `Daftar Periksa`, `Poin Penting`, `Sumber`
- Note: Code identifiers (`WindowTracker`, `Evaluator`, `Burn Rate`, `CanDeploy`, `isGood`) kept verbatim.

### NB-1 — Code comment stale (main.go:75) — No content change
- Content already correct at 9.09x cumulative (`03-code-snippets.md:442`, `02-master-draft.md:117`). Code comment outside content scope, not edited per pipeline override.

## Verification
- 02-master-draft.md: no `## ` header remains fully English (excluding code terms).
- Failure Scenario now 10% consistent with demo Phase 2.
- Phase 4 note illustrates 2% differential without altering demo numbers.
- D3 diagram visually distinguishes ignored BurnRateRule fields.
- No research/ or code/ files edited.
