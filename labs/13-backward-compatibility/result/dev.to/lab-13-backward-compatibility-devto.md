---
title: "Backward Compatibility Must Survive the Transition"
published: false
description: "Explore a phone-to-phones migration in Go with additive contracts, dual writes, backfill, reconciliation, and explicit rollback limits."
tags: go, architecture, testing, softwareengineering
---

# Backward Compatibility Must Survive the Transition

The old system stores one phone number in `users.phone`. The new requirement is a list of numbers per user. Replacing a JSON string with an array is easy to code, but an old consumer still expects that string.

Lab 13 uses **Expand → Migrate → Contract** to let old and new representations coexist. The important question is not only whether the final design works. It is whether readers and writers can survive every intermediate state.

> Compatibility must survive the transition, not just the final state.

The implementation uses an in-memory Go store. It teaches a migration sequence; it does not demonstrate a live database migration or establish zero downtime.

## Expand the contract before retiring the old representation

The enriched response preserves `phone` as a string and adds a separate `phones` array:

```go
// model.go: complete response type
type UserResponse struct {
    ID     int          `json:"id"`
    Name   string       `json:"name"`
    Phone  string       `json:"phone,omitempty"`
    Phones []PhoneEntry `json:"phones"`
}
```

`PhoneEntry` contains an ID, user ID, number, and primary marker. `LegacyConsumerDTO` expects the scalar field; `ModernConsumerDTO` expects the array.

`TestSerializationBackwardCompatibility` creates a dual-written user, marshals one enriched payload, and unmarshals it into both DTOs. Its assertions check the legacy number, two modern entries, and the primary number. That is useful evidence for these two Go consumers and this fixture, not a universal claim about all strict decoders or generated clients.

Empty representations also matter. The scalar has `omitempty`, so an empty string may disappear from the enriched payload. `CreateUser` in legacy-only mode leaves `Phones` nil, which can serialize as `null`; the store's empty read slice can serialize as `[]`. Those are different payload shapes.

## API version and storage mode are separate choices

The V1 handler returns `LegacyConsumerDTO`; the V2 handler returns `ModernConsumerDTO`. They do not both return the same enriched payload used by the serialization test.

The service also has independent write and read modes:

| Mode | Behavior |
| --- | --- |
| `WriteLegacyOnly` | Store the scalar; ignore extra phone arguments. |
| `WriteDual` | Store scalar and modern primary, then add extras to modern storage. |
| `WriteNewOnly` | Store modern entries with a nil legacy phone. |
| `ReadLegacyOnly` | Derive the response's scalar from legacy storage. |
| `ReadFallback` | Prefer modern primary/first entry, otherwise use legacy. |
| `ReadNewOnly` | Derive the scalar from modern entries. |

There is a subtle boundary here: `GetUser` always fetches modern entries before selecting the scalar's source. `ReadLegacyOnly` does not hide the enriched response's `Phones` array. `GetModernUser` returns that array.

## Dual-write protects what the old reader still needs

`MemoryStore.CreateDual` creates the user and modern primary phone under one mutex. The old scalar stays populated, allowing a legacy read of records created during the dual-write phase.

Additional numbers exist only in the modern representation. Backward compatibility preserves the old view; it does not give an old client all new capabilities.

The atomicity claim needs a precise scope. `Service.CreateUser` saves `extraPhones` through separate `SavePhoneEntry` calls after `CreateDual`, ignoring errors from those extra calls. The complete multi-phone operation is not one database transaction. The store is a collection of maps, while `schema.sql` is a schema illustration rather than an executed storage adapter.

## Backfill historical data in bounded batches

Dual-write covers creates passing through that mode. It does not migrate historical users by itself.

`BackfillWorker` scans IDs after `LastProcessedID`. It copies a nonempty legacy phone only when modern entries are absent. Existing modern entries are skipped, so backfill is not a repair mechanism for every mismatch.

The storage helper prevents duplicate numbers for the same user:

```go
// store.go: SavePhoneEntry excerpt
for _, p := range s.userPhones[userID] {
    if p.Number == number {
        return p, nil
    }
}
```

Returning an existing entry does not update its primary marker. The behavior is deduplication by user and number, not a general merge policy.

The worker keeps `LastProcessedID`, `TotalMigrated`, and `IsComplete` in memory. Its test processes ten users as batches of three, three, then four remaining records, and expects zero from rerunning the completed worker.

“Resumable” therefore means continuing the same worker object. There is no durable checkpoint or restart-recovery demonstration. Once `IsComplete` is set, that worker stops scanning; newly created legacy-only records are not automatically revisited. A completed scan does not establish that all writers have switched.

## A fallback read can repair storage without repairing this response

`ReadFallback` prefers modern data. If the modern list is empty and the legacy pointer exists, it sets the scalar from legacy and performs a lazy `SavePhoneEntry`.

However, the local `phones` slice was loaded before that write and is not reloaded afterward. The first fallback response can contain the correct scalar and an empty array, even though storage contains the new entry by the time the method returns. `GetModernUser` can therefore receive an empty array on that first read.

`TestFallbackRead` checks the scalar and inspects storage after the call. It does not assert the first response's array. That difference is exactly why reading assertions matters more than inferring guarantees from a test name.

## Reconciliation detects drift; it does not fix it

`ReconcileData` compares a populated legacy phone with the modern primary number. Missing entries and mismatched primary numbers count as drift. It skips users whose legacy pointer is nil and returns immediately after contract has been applied.

The function returns the number of findings in that scan. `DriftDetected`, by contrast, is cumulative: rescanning the same mismatch increments it again. A clean current scan and a nonzero historical counter can both be correct.

There is no automatic repair or periodic scheduler. The demo and tests invoke reconciliation explicitly. Zero drift within this comparison is not proof that every aspect of modern data is correct.

## Rollback safety depends on write history

`TestRollbackScenarios` contrasts two states:

| Record creation mode | Result after returning to legacy reads |
| --- | --- |
| Dual-write | The primary phone remains visible. |
| New-only | The scalar is empty, without a read error. |

The modern number may still exist. What fails is visibility through the old representation. Restoring an old application mode does not populate data that the old reader needs.

This rollback boundary appears before the column is dropped: merely retaining a legacy column is insufficient if new records have stopped writing it. The demo changes flags in one process; it does not deploy and roll back separate binaries against a shared database.

## Contract needs evidence that the current guard does not gather

The initial guard in `ApplyContract` is:

```go
// service.go: ApplyContract excerpt
if !force && s.obs.LegacyReadHits.Load() > 0 {
    return fmt.Errorf("%w: recorded %d legacy reads",
        ErrContractViolation, s.obs.LegacyReadHits.Load())
}
```

`LegacyReadHits` is a cumulative counter, not a traffic rate or observation window. The guard has no managed reset or time window. `GetLegacyUser` increments it before checking errors, so failed attempts count too.

The demo calls `ApplyContract(true)`, bypassing this guard. That does not prove a zero-traffic period. Even without force, a zero counter can allow contract before historical data has been backfilled: the method does not check worker completion, drift, or old writers.

It then changes write mode, changes read mode, marks contract applied, and clears legacy pointers in storage. Those steps do not form one atomic service-wide transition.

After contract, legacy service reads return an error and the V1 handler responds with `410 Gone`. But the enriched `UserResponse` type still has `Phone`, and new-only reads can populate it from modern entries. This is retirement of the legacy route and storage value, not removal of the enriched field from source code.

## Do not promote illustrative SQL or headers into production guarantees

The example SQL keeps `users.phone TEXT NOT NULL` during expansion. The in-memory `CreateModern` creates users with a nil legacy phone. The SQL file does not implement the complete sequencing needed to reconcile new-only inserts with the still-present database constraint.

No test applies the DDL to a database or measures database locks. The example is not evidence of an online, zero-downtime schema migration.

The V1 handler emits hardcoded `Deprecation` and `Sunset` values. Its test checks `Deprecation == "true"` and a nonempty `Sunset`, not standards compliance or time-triggered retirement. `410` follows `ApplyContract`, not a sunset scheduler. This article therefore makes no RFC-compliance claim for those values.

## Read concurrency coverage narrowly

The concurrency test overlaps creates, legacy/modern reads, backfill, and reconciliation. Many returned errors are ignored, and its final explicit assertion checks only `DualWriteErrors`.

That does not establish final consistency of every record or safe concurrent contract transitions. Running a race detector is valuable evidence for executed paths, but is not a substitute for semantic assertions about the migration.

To reproduce the lab from its own module directory:

```sh
cd labs/13-backward-compatibility
go test -v ./...
go test -race ./...
go run ./cmd/demo
```

The module declares Go 1.26.7. Go is unavailable in this authoring environment, so I inspected the source and tests rather than executing them. Execution reports committed in the repository are historical author reports, not a fresh run performed for this article.

## Preserve the transition, not just the destination

Expand keeps both representations available. Migrate moves writes, historical data, and consumers with explicit checks. Contract retires the old path after its dependencies are actually gone.

The pattern supplies an order of work, not an automatic safety guarantee. Compatibility still depends on the real reader behavior, write history, checkpoint durability, and retirement gates implemented along that path.

## Source

[Software Engineering Lab 13](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/13-backward-compatibility), commit `82af988b89439fa92b1200406c5afbd59ce72a56`: `internal/compat/*.go`, `tests/*.go`, `cmd/demo/main.go`, `README.md`, `schema.sql`, and `go.mod`. Code excerpts belong in their existing package context.
