# Content Revision Record

## Target Lab
labs/15-load-testing

## Audit Source
engineering-audit-opensource/05-gaps.md — G5 (MEDIUM): Latency not recorded for failed requests.

## Issue Summary
- **Severity**: MEDIUM
- **Finding**: Engineering code-audit Finding 3 (G5) identified that the load runner only records latency for successful (HTTP 2xx) requests. Transport errors and HTTP >= 400 responses are counted as errors but their durations are not captured. The content previously implied all request latencies were collected for percentile computation, without explicitly clarifying this scoping.

## Changes Made

### 01-content-brief.md
- **Line 24 (Warnings)**: Added new warning bullet clarifying that percentile metrics (P50/P90/P95/P99) only cover successful requests (HTTP 2xx); failed requests are counted as errors but their latency is not measured.

### 02-master-draft.md
- **Line 34 (How It Works)**: Clarified that the load runner collects latency only for successful requests (HTTP 2xx). Added explicit statement that failed requests (transport errors or HTTP >= 400) are not recorded in latency metrics but are counted as errors.
- **Line 170 (Code Walkthrough Snippet 2)**: Added note after the explanation that latency is only recorded for HTTP 2xx responses, and that error responses are counted separately as `errs`.
- **Line 218 (Production Considerations)**: Added bullet noting that in production, error response latencies should also be captured—not just error count—to reveal timeout or dependency failure patterns.

### 03-code-snippets.md
- **Line 100 (Snippet 2 Explanation)**: Added explicit note that latency is only recorded for successful (HTTP 2xx) requests; transport errors and HTTP >= 400 responses do not produce entries in the `lats` slice but are counted as `errs`. Clarified that P50/P95/P99 metrics reflect only successful-request latency.

### 05-key-takeaways.md
- **Line 7 (new takeaway)**: Added point #7 clarifying that percentile metrics only cover HTTP 2xx responses; error latencies are not measured (only error count is tracked). Warned this can mislead analysis if failures carry high latency.

## Files Not Modified (Out of Scope per Pipeline Override)
- `engineering/01-design.md` — DOC_CODE_MISMATCH (iterations vs bounded duration) — engineering doc, not content file.
- `engineering/02-implementation-notes.md` — Tail latency "strictly queuing" claim — engineering doc, not content file. Content master-draft correctly describes the 10% random slowdown.
- `engineering/03-execution-result.md` — Stale test list — engineering doc, not content file.
- `cmd/demo/main.go` — P90 not printed — code file, not content.

## Verification
- `go test -v ./...` — all 12 tests pass
- `go test -race ./...` — no races detected
- Content reviewed against code audit findings from engineering-audit-opensource/
- Content brief updated to include the new limitation warning

## Status
READY_FOR_CONTENT_ADAPTER
