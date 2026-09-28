# Technical Content Audit: SLI/SLO Error Budget Lab

## Executive Summary

This audit evaluates the technical content in `labs/24-slo-sli-error-budget/content/` against the approved research (`research/`) and engineering (`engineering/` source, tests, execution results). Overall, the content demonstrates strong fidelity to the core concepts with proper sourcing and transparency about limitations.

However, audit revealed **three significant accuracy issues** in specific examples and one diagram, alongside several minor observations. The core conceptual material (SLI/SLO/Error Budget definitions, multi-window alerting, burn rate concepts, endpoint criticality) remains correct and well-supported.

---

## Detailed Findings

### HIGH SEVERITY INACCURACIES

#### 1. Inaccurate Burn Rate Calculation Example
- **Location**: `content/03-code-snippets.md`, Snippet 11 explanation (line 442)
- **Claim**: "90 request sukses + 10 request gagal (status 500, latency 500ms) menghasilkan 10% error rate. Pada SLO 99.9% (allowed 0.1%), burn rate = 100x"
- **Reality**: 
  - The incident batch has 10/100 = 10% error rate (correct)
  - But AlertEngine computes burn rate on cumulative data: Phase 1 baseline (1000 good requests) + Phase 2 incident (90 good, 10 bad) = 1100 total requests, 10 bad
  - Actual error rate = 10/1100 ≈ 0.91%
  - Burn rate = 0.91% / 0.1% = **9.09x** (not 100x)
  - Demo output and execution result confirm: `ShortBurn: 9.09x | LongBurn: 9.09x (Threshold: 6.00x)`
- **Severity**: HIGH – Misstates a core quantitative result
- **Evidence**: 
  - Demo source: `cmd/demo/main.go` lines 75, 102-103, 112-113
  - Execution result: `engineering/03-execution-result.md` lines 62, 103
  - Burn rate formula: `internal/alerting/engine.go` lines 51-61

#### 2. Inconsistent Burn Rate Example in Master Draft
- **Location**: `content/02-master-draft.md`, lines 169-170
- **Content**: Shows generic formula with comment `// mis. 10% / 1000` and `// 10% / 0.1% = 100x`
- **Issue**: While the Case Study section (line 284) correctly calculates 10/1100 = 0.91% → 9.1×, the initial formula example implies the incident calculation yields 100x burn rate, creating confusion.
- **Severity**: MEDIUM – Risk of reader misunderstanding the actual calculation
- **Mitigation**: Case Study section later corrects this (line 284)

#### 3. Incorrect FALSE POSITIVE Illustration in Diagram D4
- **Location**: `content/04-diagrams.md`, Diagram D4: Multi-Window Alert: True Positive vs False Positive
- **Issue**: 
  - The diagram's "FALSE POSITIVE (no alert)" column shows:
    - Short Window: 90 OK + 10 ERR → 10% error → 100x burn
    - Long Window: 90 OK + 10 ERR → 10% error → 100x burn
  - Then states: "Short ≥ 6.0 BUT Long < 6.0 → NO ALERT"
  - **Contradiction**: If both windows show 100x burn rate, BOTH are ≥ 6.0 (threshold), so this is a TRUE POSITIVE (alert should fire), not a false positive.
  - The correct false positive scenario (transient spike) has:
    - Short window: High error (e.g., 100x burn)
    - Long window: Low error (e.g., 0.1x burn, below threshold)
- **Correct Implementation**: Verified in `tests/slo_test.go` lines 133-151 (negative test for transient spikes)
- **Severity**: HIGH – Diagram teaches incorrect alerting logic
- **Evidence**: 
  - Test: `tests/slo_test.go` lines 133-151
  - AlertEngine logic: `internal/alerting/engine.go` line 73 (`shortBurn >= threshold && longBurn >= threshold`)

### MEDIUM SEVERITY OBSERVATIONS

#### 4. Ambiguous Comment in Formula Example
- **Location**: `content/02-master-draft.md`, line 169
- **Issue**: Comment `// mis. 10% / 1000` mixes a percentage (10%) with a count (1000). Should clarify as `// mis. actualErrorRate = 2% (20 errors / 1000 total)` or similar.
- **Severity**: LOW – Minor readability issue

#### 5. Over-Simplified Release Policy Description
- **Location**: `content/02-master-draft.md`, line 128
- **Claim**: "Budget tinggi → boleh release cepat, eksperimen; budget habis → STOP deployment berisiko, fokus reliability."
- **Issue**: While aligned with the core concept, omits the full Google SRE error budget policy details (postmortem wajib + P0 action item for single incident >20% budget, halt all changes except P0/security when budget exhausted, per `research/03-evidence.md` Evidence 15 and `research/05-report.md` Finding 14).
- **Severity**: LOW – Conceptually correct but omits policy nuance

#### 6. Limited Demonstration of Endpoint Criticality Differentiation
- **Location**: Multiple content files referencing Phase 4 (e.g., `content/02-master-draft.md` line 299, `content/05-key-takeaways.md` line 15)
- **Issue**: 
  - Phase 4 feeds identical 10% error rate to both endpoints (Payment 99.9%, Reports 95.0%)
  - Both exhaust their budgets simultaneously (Payment: -8.90, Reports: -5.00)
  - While Reports does have wider tolerance (5% vs 0.1%), the identical error rate means both fail regardless of tolerance width
  - Content correctly shows both CanDeploy=false but could clarify that with lower error rates (e.g., 2%), Reports would remain deployable while Payment might not
- **Severity**: LOW – Concept accurate but demonstration doesn't fully showcase the differentiation benefit

### CORRECTLY DOCUMENTED LIMITATIONS & SIMPLIFICATIONS

The content properly discloses the following limitations and simplifications:

1. **Burn rate thresholds (14.4×, 6.0×) are Google recommendations** (`content/05-key-takeaways.md` line 23, `content/01-content-brief.md` line 48)
2. **In-memory metrics reset on process restart** (`content/02-master-draft.md` line 51, `content/05-key-takeaways.md` line 19)
3. **Latency not measured as full histogram** (boolean isGood predicate used) (`content/02-master-draft.md` line 54, `content/05-key-takeaways.md` line 21)
4. **Zero traffic edge case handled** (SLI=1.0, CanDeploy=true) (`content/02-master-draft.md` line 102, `content/05-key-takeaways.md` line 18)
5. **"70% outages from change" is Google-internal observation** (`content/05-key-takeaways.md` line 12, `content/01-content-brief.md` line 49)
6. **Cost ~100x per nine is heuristic, not mathematical law** (`content/05-key-takeaways.md` line 11, `content/01-content-brief.md` line 50)
7. **Window demo uses time compression** (30 min ≡ 30 days) (`content/02-master-draft.md` line 53, `content/01-content-brief.md` line 53)

### RESEARCH & ENGINEERING FIDELITY ASSESSMENT

#### ✅ Strong Alignment Areas
- **SLI Definition**: Content correctly states "SLI = good events ÷ total events" (ratio 0–100%), matching research (Evidence 3, Finding 2) and implementation (evaluator.go:44-47)
- **Error Budget Formula**: Content correctly defines as `1 − SLO` (Finding 7), with implementation `totalErrorBudget = (1.0 − target) * float64(total)` (evaluator.go:49-50)
- **Multi-window Alerting**: Content correctly requires BOTH short AND long window to exceed threshold (engine.go:73), matching research (Evidence 9, Finding 9) and test (TestAlertEngineBurnRate)
- **Endpoint Criticality**: Content correctly shows different SLOs per endpoint (Payment 99.9% vs Reports 95.0%), matching research (Evidence 12, Finding 12) and demo
- **Thread Safety**: Content correctly notes mutex protection (content/05-key-takeaways.md line 17), verified in tracker.go:22-28,46-48 and TestConcurrencyMetrics
- **Research Sourcing**: Content properly attributes definitions to Google SRE Book/Workbook and acknowledges vendor concentration (content/05-key-takeaways.md line 11-12)

#### ✅ Engineering Implementation Accuracy
All code snippets in `content/03-code-snippets.md` are verbatim extracts from the approved source, as stated in the file header and verified against:
- `internal/metrics/tracker.go`
- `internal/slo/evaluator.go` 
- `internal/alerting/engine.go`
- `cmd/demo/main.go`

---

## Root Cause Analysis of Issues

1. **Burn Rate Calculation Errors**: Appear to stem from presenting the incident batch error rate (10%) as the cumulative window error rate, neglecting the Phase 1 baseline traffic already present in the tracking windows.

2. **D4 Diagram Error**: The FALSE POSITIVE column appears to be a copy-paste error where it duplicates the TRUE POSITIVE column's data, with incorrect labeling that creates internal contradiction ("Short ≥ 6.0 BUT Long < 6.0" while showing both at 100x).

3. **Formula Comment Ambiguity**: The `mis. 10% / 1000` comment mixes rate and count units, likely intended as a hypothetical example but poorly annotated.

---

## Recommendations for Correction

### For `content/03-code-snippets.md` Snippet 11:
**Change**: 
```
90 request sukses + 10 request gagal (status 500, latency 500ms) menghasilkan 10% error rate. Pada SLO 99.9% (allowed 0.1%), burn rate = 100x.
```
**To**:
```
90 request sukses + 10 request gagal (status 500, latency 500ms) memiliki 10% error rate pada batch insiden.
Akumulasi dengan baseline Fase 1 (1000 request sukses) menghasilkan 1100 total request dengan 10 error (0.91% error rate).
Pada SLO 99.9% (allowed error rate 0.1%), burn rate = 0.91% / 0.1% = 9.09x.
Evaluator mendeteksi budget habis (CanDeploy = false). AlertEngine memicu slow burn alert (6.0x) karena kedua window melewati threshold (9.09x ≥ 6.0x).
```

### For `content/02-master-draft.md` Lines 169-170:
**Change**:
```
actualErrorRate := float64(bad) / float64(total)   // mis. 10% / 1000
allowedErrorRate := 1.0 - targetSLO                 // 0.1%
burnRate := actualErrorRate / allowedErrorRate      // 10% / 0.1% = 100x
```
**To**:
```
// Contoh: Jika ada 20 error dari 1000 total request (error rate 2%):
actualErrorRate := float64(bad) / float64(total)
allowedErrorRate := 1.0 - targetSLO
burnRate := actualErrorRate / allowedErrorRate      // 2% / 0.1% = 20x
```

### For `content/04-diagrams.md` Diagram D4:
**Fix the FALSE POSITIVE column** to correctly show a transient spike scenario:
- **Short Window**: 90 OK + 10 ERR → 10% error → 100x burn (above threshold)
- **Long Window**: 9999 OK + 1 ERR → 0.01% error → 0.1x burn (below threshold)
- **Result**: Short ≥ 6.0x (YES) BUT Long ≥ 6.0x (NO) → NO ALERT ✓ (transient spike filtered)

Alternatively, relabel this column as "TRUE NEGATIVE" (short above threshold, long below → no alert) and ensure the existing TRUE NEGATIVE column shows a different scenario (e.g., both windows below threshold).

### For Phase 4 Demonstration Clarity:
Add a note that with lower error rates (e.g., 2%), the Reports endpoint (95.0% SLO) would remain deployable while Payment (99.9% SLO) might exhaust its budget, demonstrating the practical value of criticality-based SLO differentiation.

---

## Summary

The content demonstrates **strong conceptual fidelity** with research and engineering, properly sourcing definitions and disclosing limitations. The **three HIGH severity issues** involve:
1. Incorrect burn rate calculation examples (claiming 100x vs actual 9.09x)
2. Misleading diagram illustrating false positive alerting logic
3. Ambiguous formula comment mixing units

These are specific, correctable errors that do not undermine the core conceptual value of the material. With the recommended corrections, the content would achieve full alignment.

Verdict: **APPROVED_WITH_WARNINGS** (see `09-verdict.md`).