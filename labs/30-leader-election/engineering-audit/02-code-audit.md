# Code Audit

## Finding 1

Location: `internal/coordinator/coordinator.go:36-65`
Claimed Behavior: Acquire lease and issue monotonically increasing fencing tokens (`revision`).
Observed Implementation: Locks mutex, checks if key exists and unexpired. If unexpired and held by caller, extends TTL without incrementing revision. If held by another caller, returns `ErrLeaseHeld`. If expired or absent, increments `c.revision++` and assigns as `FencingToken`. Returns copied struct pointer.
Assessment: PASS
Severity: LOW
Notes: Correctly enforces monotonic fencing token generation on new grants.

## Finding 2

Location: `internal/coordinator/coordinator.go:67-90`
Claimed Behavior: Renew existing lease if valid caller and token match.
Observed Implementation: Locks mutex, verifies key presence, checks `now.After(ExpiresAt) || now.Equal(ExpiresAt)`, deletes expired lease and returns `ErrLeaseExpired`. Checks `HolderID` and `FencingToken`, updates `ExpiresAt` and `TTL`. Returns copied struct.
Assessment: PASS
Severity: LOW
Notes: Safely validates token ownership and expiration before extending.

## Finding 3

Location: `internal/storage/storage.go:33-49`
Claimed Behavior: Shared storage gate enforces check-and-set condition `token > lastSeenToken`.
Observed Implementation: Locks write mutex, compares `token <= s.lastSeenToken`. Rejects with wrapped `ErrStaleFencingToken`. Updates `s.lastSeenToken = token` and appends immutable copy to `s.records`.
Assessment: PASS
Severity: LOW
Notes: Atomically protects state against stale leader writes.

## Finding 4

Location: `internal/candidate/candidate.go:113-155`
Claimed Behavior: Continuous background loop for campaign, renewal, and pause handling.
Observed Implementation: Runs ticker on `renewRate`. Checks `now.Before(n.pauseUntil)` to skip tick (simulating STW pause). Manages leader renewal vs follower acquisition.
Assessment: PASS
Severity: LOW
Notes: Clean context-cancellable background loop simulating process unresponsiveness.
