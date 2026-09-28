## Gap 1

Type:
MISSING_TEST

Severity:
LOW

Location:
tests/slo_test.go (WindowTracker edge cases)

Problem:
In‑memory buckets retain data only for window duration; reset on process restart not explicitly tested.

Required Revision:
Add test verifying zeroed state after full window eviction and after in‑memory reset (e.g., after 30‑day window expiry).

Can Be Approved Without Fix:
YES
