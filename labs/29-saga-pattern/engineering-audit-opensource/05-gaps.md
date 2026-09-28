# Gaps

1. Lock possibly remains after failed saga (not explicitly released on compensation unless Cancel called).
   Type: MISSING_EDGE_CASE
   Severity: MEDIUM

2. Compensation error is aggregated but not prominently surfaced; final error message combines step and compensation errors but compensation error may be lost in %v formatting.
   Type: UNHANDLED_ERROR
   Severity: MEDIUM

3. Design docs reference pkg/ rather than internal/; minor inconsistency.
   Type: DOC_CODE_MISMATCH
   Severity: LOW

4. Design references Delivery step missing from code/demo (uses ApproveOrder instead).
   Type: RESEARCH_MISMATCH
   Severity: LOW

5. No benchmark despite concurrency claims; tests cover concurrency but without performance validation.
   Type: UNVERIFIED_RESULT
   Severity: LOW