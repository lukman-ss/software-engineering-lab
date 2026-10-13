The new online-booking flow is ready to deploy. Does that mean every workshop customer must use it immediately?

Lab 12 separates those decisions. Deployment puts code in place; release exposes its behavior to users. The `online_booking` flag controls that exposure through internal targeting, percentage rollout, and a global OFF switch.

> Deploy does not have to mean release.

The lab is a small Go simulation, not a production flag platform. Reading its boundaries is part of understanding what a feature flag actually guarantees.

## Keep the decision at a clear boundary

`FeatureService` owns flag configuration and evaluation. `BookingService.CreateBooking` makes one routing decision rather than scattering targeting rules throughout business logic:

```go
// booking_service.go: complete method
func (s *BookingService) CreateBooking(req BookingRequest) BookingResponse {
    if s.featureService.IsEnabled("online_booking", req.UserID) {
        return s.newBookingFlow(req)
    }
    return s.legacyBookingFlow(req)
}
```

Both flows return simulated responses. The new flow succeeds unless `simulateFail` is true; the legacy flow always succeeds. There is no database write, mechanic reservation, or real vendor call. `Vendor Error (Simulated)` is deliberately injected behavior.

This distinction matters: the lab demonstrates choosing a path, not implementing the booking transaction itself.

## Evaluation order is part of the contract

`IsEnabled` applies rules in this order:

| Order | Rule |
| --- | --- |
| 1 | Unknown flag returns false. |
| 2 | `Enabled == false` returns false for everyone. |
| 3 | A matching ID in `InternalUsers` returns true. |
| 4 | Percentage <= 0 returns false for ordinary users; >= 100 returns true. |
| 5 | For percentages 1–99, compare the user's bucket with the threshold. |

Global OFF wins over internal targeting. Conversely, `Enabled: true` with a zero percentage can release to internal users while leaving ordinary users on legacy.

The `User` struct has an `IsInternal` field, but evaluation does not use it. Internal access comes from matching `userID` against the configured list. Do not infer targeting behavior from a field that is never read.

The percentage comment says 0–100, but there is no rejecting validation. Negative values behave like zero; values above 100 behave like 100 for ordinary users.

## Deterministic cohorts, not random requests

Randomizing every request could send the same user through different flows on consecutive calls. The implemented calculation is deterministic:

```go
// feature_service.go: bucketing excerpt
hashInput := fmt.Sprintf("%s:%s", userID, featureKey)
hash := sha256.Sum256([]byte(hashInput))
bucketVal := binary.BigEndian.Uint64(hash[:8])
bucket := int(bucketVal % 100)
return bucket < flag.RolloutPercentage
```

The exact input includes a colon separator. SHA-256 produces the digest; the first eight bytes are interpreted as a big-endian uint64; modulo 100 gives a bucket from 0 to 99.

At 10%, buckets 0–9 qualify. At 50%, buckets 0–49 qualify. With the same IDs, feature key, algorithm, and targeting rules, the original 10% cohort remains included as the threshold increases.

For `online_booking`, recalculating the source algorithm puts `user-3` in bucket 25 and `user-4` in bucket 83. The former enters at 50%; the latter remains excluded until a higher threshold or the 100% shortcut. These are derived examples, not recorded Go execution output.

Consistency is conditional. Changing configuration is supposed to change eligibility. Changing the user ID or feature key can also change the bucket.

## A 10% threshold does not promise exactly ten requests

The CLI generates 100 users: three internal IDs, then `user-3` through `user-99`. It issues one request per user in each phase. Internal users remain allowlisted through the percentage phases.

Re-evaluating this exact fixture with the source hash algorithm gives:

| Phase | New flow | Legacy |
| --- | ---: | ---: |
| OFF | 0 | 100 |
| Internal only | 3 | 97 |
| 10% plus internal allowlist | 15 | 85 |
| 50% plus internal allowlist | 48 | 52 |
| 100% | 100 | 0 |

I calculated these counts in Python from the repository's fixture and algorithm; I did not run the Go demo. At 10%, 12 of the 97 ordinary users qualify, plus all three internal users. At 50%, 45 ordinary users qualify, plus the three internal users.

Two effects explain the difference from the threshold label: a finite sample need not contain exactly 10% of users in buckets 0–9, and internal targeting overrides the percentage gate.

User exposure is also not necessarily traffic exposure. Users can send unequal request volumes. The demo makes those quantities align by sending exactly one request per user per phase. A 10% cohort is therefore not a hard ceiling on every possible production impact.

## Observe before increasing exposure

The README recommends OFF, internal, 10%, 50%, and 100%, with observation between increases. The CLI demonstrates that sequence, but its transitions are hardcoded. It does not wait for an observation window or inspect a threshold before advancing.

That makes it an illustration of rollout stages, not an automatic rollout controller.

The runner's `phase` function records metrics after calling `CreateBooking`:

- Total requests.
- New-flow and legacy-flow request counts.
- Aggregate booking successes and failures.
- Response-time samples and their mean.

Although `BookingService` accepts a metrics pointer, `CreateBooking` does not update those counters. Calling the service directly does not provide the CLI's instrumentation automatically.

The error-rate formula in `Metrics.Snapshot` is:

```go
// metrics.go: struct-literal field excerpt
ErrorRate: float64(m.BookingFailed) /
    float64(max(1, m.TotalRequests)) * 100,
```

The denominator is all requests, not just requests using the new flow. During a partial rollout, aggregate error rate can dilute a severe failure rate inside the new cohort. The current metrics split request volume by flow but combine success/failure outcomes. Per-flow error rates would be an extension, not an existing feature.

Metrics reset at each phase. `AvgResponseTime` is a mean, not a percentile. It measures a local simulated method with no database or network work, so its latency should not be presented as a production booking benchmark.

## A kill switch changes subsequent decisions

After reaching 100%, the CLI sets `simulateFail` to true. Every request in phase 6a selects the new flow and fails. Then it applies this configuration:

```go
// cmd/simulate/main.go: kill-switch configuration
service.SetFlag(featureflags.Flag{
    Key:     "online_booking",
    Enabled: false,
})
```

Requests in phase 6b select legacy and succeed, even though the failure simulation remains enabled. Turning off the feature avoids the failing path on subsequent evaluations that read the updated configuration.

It does not retry an already failed request through legacy. It does not cancel an evaluation that already read the old configuration, undo completed side effects, or remove deployed code. The distinction between rerouting future work and repairing past work is essential to the kill-switch claim.

## Runtime configuration has an operational boundary

`FeatureService` stores flags in an in-memory map protected by an `RWMutex`. `SetFlag` replaces the entire `Flag` value. It is not a partial-field update: leaving out `RolloutPercentage` and `InternalUsers` resets them to their zero values.

The CLI sends complete targeting configuration during rollout. A caller that wants to preserve targeting must do the same rather than assume omitted fields remain unchanged.

There is no persistent configuration store, cross-replica distribution, operator dashboard, or administration endpoint. “Without redeploy” means updating the object in this running process. The map's mutex also does not establish concurrency safety for every other field: `simulateFail` is a plain bool, and toggling it concurrently with booking reads is not synchronized by this implementation. The demo itself runs sequentially.

## A flag is not authorization

Authorization asks whether the user is allowed to create a booking. The feature flag asks which booking behavior is selected for that user.

The README explicitly distinguishes these questions. This implementation answers the second one; it does not implement an authorization check. An internal allowlist for release targeting should not be described as a permission system.

## Test the paths, then remove the temporary branch

The test file covers global OFF, zero percentage, internal inclusion, examples on both sides of a 10% gate, repeated evaluation, increased rollout, kill switch, and unknown flags. Booking tests check ON/OFF responses and a simulated failure followed by an OFF update and a successful legacy response.

The percentage test searches for one ON user and one OFF user. It does not prove an exact 10% distribution. There are no metric assertions or concurrent execution tests in this file. Test names should not be read as broader guarantees than their assertions provide.

Use the repository's Makefile targets:

```sh
# From the repository root
make lab-12-demo
make lab-12-test

# Or enter the nested lab module
cd labs/12-feature-flags
go run ./cmd/simulate
go test -v -count=1 ./...
```

The lab's own `go.mod` declares Go 1.26.7. Entering that module matches the Makefile; the README's root-relative `go run` alternative needs to account for the nested-module boundary. Go is unavailable in this authoring environment, so I inspected the tests and demo rather than executing them. No passing test result is claimed.

Finally, the README's lifecycle continues beyond 100%: observe stability, remove the flag, then remove unused legacy code. That cleanup is a recommended future change, not a refactor already performed in this lab.

A release flag is useful because it separates code delivery from user exposure. Its value depends on consistent targeting, meaningful evidence, tested fallback behavior, and a plan to retire the extra branch.

## Source

[Software Engineering Lab 12](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/12-feature-flags), inspected at commit `82af988b89439fa92b1200406c5afbd59ce72a56`: the README, feature service, booking service, metrics, tests, CLI simulation, lab `go.mod`, and Lab 12 targets in the root Makefile. Code excerpts belong in their existing package context.
