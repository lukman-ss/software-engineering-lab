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

### 4. content/03-code-snippets.md:24
- **Issue**: Code snippet comment deviation — Snippet 1 (Fault Injector) included non-verbatim comment "(unused field per audit, kept for structure)" on `errorRate` field.
- **Fix**: Removed editorial annotation to make snippet verbatim with source code: `errorRate  float64 // 0.0 to 1.0`
- **Impact**: Code snippet now exactly matches internal/fault/injector.go:16

### 5. content/03-code-snippets.md (Snippet 3) and content/02-master-draft.md
- **Issue**: `Monitor.Metrics()` method not documented in snippet 3 and not mentioned in master draft's "Code Walkthrough" or "How It Works" sections.
- **Fix**: 
  - Added `Metrics()` method to Snippet 3 in content/03-code-snippets.md (lines 230-242)
  - Updated Code Walkthrough in content/02-master-draft.md (line 60) to mention `Metrics()` method
- **Impact**: Public method used in demo is now documented; readers can trace how `mon.Metrics()` produces output shown in case study.

### 6. content/06-source-map.md:67-68
- **Issue**: Source map reference error — cited `engineering/01-design.md:37-38` for "What Is Not Demonstrated" but actual location is `engineering/02-implementation-notes.md:34-36`.
- **Fix**: Changed reference from `engineering/01-design.md:37-38` to `engineering/02-implementation-notes.md:34-36`
- **Impact**: Source traceability corrected for one entry.

## Verification
- All code snippets (content/03-code-snippets.md) correctly reflect source code
- No research or engineering files modified — content-only revisions
- Zero new dependencies or claims introduced
- Addressed all warnings from content-audit/09-verdict.md

---

READY_FOR_CONTENT_ADAPTER