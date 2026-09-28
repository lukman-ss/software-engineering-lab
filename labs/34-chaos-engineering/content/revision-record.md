# Revision Record — labs/34-chaos-engineering

Date: Mon Sep 28 2026
Auditor Findings Addressed: content-audit/09-verdict.md (APPROVED_WITH_WARNINGS)

## Changes Made

### 1. content/02-master-draft.md:29 (Line 29)
- **Issue**: Technical inaccuracy — text described Injector using `time.Sleep`, but actual implementation (`internal/fault/injector.go:66-70`) uses `select { case <-time.After(latency): case <-ctx.Done(): }`.
- **Fix**: Changed "menggunakan `time.Sleep` dengan pengecualian `ctx.Done()`" → "menggunakan `time.After` dengan dukungan pembatalan konteks `ctx.Done()`".
- **Impact**: Text now matches code snippet (Snippet 1) and actual implementation.

### 2. content/02-master-draft.md:87 (Line 87)
- **Issue**: Threshold context — mentions specific percentage (20%) without clarifying it's a lab-specific example value.
- **Fix**: Changed "threshold 20%" → "threshold (lab example: 20%)".
- **Impact**: Readers now understand 20% is an illustrative lab value, not a production recommendation (covered in Production Considerations).

### 3. content/06-source-map.md
- **Issue**: Gap #2 flagged line variance for runner.go methods.
- **Verification**: Current runner.go `terminate` function at lines 91-97 matches source-map reference. No edit needed — variance was from prior edit, now aligned.

## Verification
- All code snippets (content/03-code-snippets.md) already correctly reflect `time.After` implementation.
- No research or engineering files modified — content-only revisions.
- Zero new dependencies or claims introduced.

---

READY_FOR_CONTENT_ADAPTER