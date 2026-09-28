# Content Revision Record

Target Lab: `labs/28-timeouts-and-deadlines`
Revision Date: 2026-09-28
Revision Scope: Content only (per pipeline override — research/code untouched)
Audit Reference: `content-audit/08-content-audit.md`, `content-audit/09-verdict.md` (NEEDS_REVISION)

## Blocking Fixes

- **B1** `04-diagrams.md:14` — Parent deadline expiry corrected: `context.Canceled` → `context.DeadlineExceeded` (cancels only explicit `cancel()`; matches `TestExecuteWithBudget_ParentTimeoutInherited`).
- **B2** `02-master-draft.md:186` — Truncated heading restored: `"supaya ama` → `"supaya aman."` (also fixed second truncated instance at line 227 and stray `**`).
- **B3** `02-master-draft.md:228`, `05-key-takeaways.md:5` — Hallucinated statistic removed: "lebih dari 50%" → "secara substantial" per AWS source wording in `research/03-evidence.md` Evidence 4 ("substantial decrease in client work and server load").

## Non-Blocking Fixes

- **N1** `02-master-draft.md:84` — Typo: "percayaan" → "percobaan".
- **N2** `03-code-snippets.md:84` — HALF_OPEN overstatement corrected: no longer claims "hanya menerima satu request uji"; states `SuccessThreshold` sukses diperlukan kembali CLOSED (matches `RecordSuccess`, test uses `SuccessThreshold: 2`).
- **N3** Backoff formula unified to `2^(attempt-1)` (matches `1 << uint(attempt-1)` in code):
  - `01-content-brief.md:13`: `base * 2^attempt` → `base * 2^(attempt-1)`
  - `02-master-draft.md:58`: already correct (`2^(attempt-1)`), unchanged
  - `05-key-takeaways.md:5`: `base×2^attempt` → `base×2^(attempt-1)`
  - `02-master-draft.md:228`: `2^attempt` → `2^(attempt-1)`
- **N5** `06-source-map.md:113` — Stale `REJECTED` status replaced with `NEEDS_REVISION` (matches current verdict) and revised-file entry added.

## Additional Correction (found during revision)

- `02-master-draft.md:230` — Backwards inheritance fixed: "parent timeout ≤ child timeout" → "child timeout ≤ parent timeout" (child = min(parent, budget), never exceeds parent).

## Files Modified

| File | Changes |
|------|---------|
| `01-content-brief.md` | N3 |
| `02-master-draft.md` | B2, B3, N1, N3, inequality fix |
| `03-code-snippets.md` | N2 |
| `04-diagrams.md` | B1 |
| `05-key-takeaways.md` | B3, N3 |
| `06-source-map.md` | N5 |

Not edited per audit: N4 (`03-code-snippets.md` Snippet 5/6 — audit verdict "Tanpa perubahan"), N6 (PASS, no change needed).

## Final Status

READY_FOR_CONTENT_REAUDIT
