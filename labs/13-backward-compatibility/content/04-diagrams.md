# Diagrams — Backward Compatibility Lab

## Diagram 1 — Overall Architecture

```
[ Client V1 (Legacy) ]     [ Client V2 (Modern) ]
          │                         │
          ▼                         ▼
┌──────────────────────────────────────────────────┐
│              HTTP / Service Layer                │
│  - APIHandler: /api/v1/users, /api/v2/users      │
│  - CompatService: GetLegacyUser, GetModernUser   │
│  - FeatureFlags: WriteMode, ReadMode             │
│  - Observability: atomic counters                │
└─────────┬────────────────────────────────┬───────┘
          │                                │
          ▼                                ▼
┌──────────────────┐             ┌─────────────────┐
│ Legacy Storage   │◄──Backfill──┤ Modern Storage  │
│ users(id,        │   Worker    │ user_phones(id, │
│  name, phone*)   │  (batched,  │  user_id,       │
│                  │  checkpoint)│  number,        │
│                  │             │  is_primary)    │
└──────────────────┘             └─────────────────┘

* phone = nil after Contract
```

Source: `engineering/01-design.md` Architecture, `internal/compat/store.go`, `internal/compat/service.go`, `internal/compat/handler.go`

## Diagram 2 — Expand → Migrate → Contract Lifecycle

```
Baseline (N)          Expand (N+1)           Migrate                Contract
users:                users:                 users:                 users:
 id | name | phone     id | name | phone     id | name | phone      id | name | phone
 1  | Alice| +62..     1  | Alice| +62..     1  | Alice| +62..      1  | Alice| NULL*
 2  | Bob  | +62..     2  | Bob  | +62..     2  | Bob  | +62..      2  | Bob  | NULL*
                      3  | Charl| +62..     3  | Charl| +62..      3  | Charl| NULL*
                      user_phones:           user_phones:           user_phones:
                       (new table)            1 | 1 | +62.. (bf)     1 | 1 | +62..
                                              2 | 2 | +62.. (bf)     2 | 2 | +62..
                                              3 | 3 | +62.. (dual)   3 | 3 | +62..
                                              4 | 3 | +62.. (dual)   4 | 3 | +62..

WriteMode: LegacyOnly  → DualWrite           → DualWrite → NewOnly  → NewOnly
ReadMode:  LegacyOnly  → LegacyOnly          → Fallback  → NewOnly  → NewOnly
Backfill:  —           → not started         → batched+ckpt        → done
Contract:  —           → false               → false               → true (drop column)

* legacy column dropped
```

Source: `internal/compat/flags.go`, `internal/compat/store.go:204-212`, `tests/migration_test.go:10-94`, `cmd/demo/main.go`

## Diagram 3 — Write Path by WriteMode

```
CreateUser(name, phone, extra...) 
        │
        ├── WriteLegacyOnly ──► store.CreateLegacy() ──► users.phone only
        │                                                  obs: -
        │
        ├── WriteDual ──► store.CreateDual() ──► users.phone + user_phones (primary)
        │                    + SavePhoneEntry(extra, isPrimary=false)
        │                                                  obs: DualWriteCount++
        │
        └── WriteNewOnly ──► store.CreateModern() ──► user_phones only (phone=nil)
                                                           obs: -
```

Source: `internal/compat/service.go:55-95`, `internal/compat/store.go:32-118`

## Diagram 4 — Read Path with Fallback

```
GetUser(id) 
   ├── ReadLegacyOnly:  return users.phone
   │
   ├── ReadFallback:
   │     if len(user_phones) > 0 → return primary number
   │     else if users.phone != nil → return phone + SavePhoneEntry(lazy backfill)
   │     else → return ""
   │
   └── ReadNewOnly:     if len(user_phones) > 0 → return primary number
                        else → return ""
```

Source: `internal/compat/service.go:97-148`

## Diagram 5 — Backfill Worker (Resumable Batch)

```
RunAll()
  loop:
    RunBatch(ctx)
      ├── checkpoint.IsComplete? → done
      ├── GetUserIDs(after: LastProcessedID, limit: batchSize)
      │     └── empty → IsComplete=true, done
      └── for each id:
            if users.phone != nil && user_phones empty:
                SavePhoneEntry(phone, isPrimary=true) // idempotent check inside
                obs.BackfillProcessed++
            LastProcessedID = id
      └── len(ids) < batchSize → IsComplete=true
```

Source: `internal/compat/backfill.go:35-101`, `internal/compat/store.go:147-171`

## Diagram 6 — Contract Guard

```
ApplyContract(force bool)
   ├── if !force && LegacyReadHits > 0 → ErrContractViolation
   ├── SetWriteMode(WriteNewOnly)
   ├── SetReadMode(ReadNewOnly)
   ├── SetContractApplied(true)
   └── store.ApplyContractDropLegacyColumn() → users.phone = nil, legacyDropped=true
        └── subsequent GetLegacyUser → ErrLegacyUnavailable → HTTP 410 Gone
```

Source: `internal/compat/service.go:228-245`, `internal/compat/handler.go:20-41`, `internal/compat/store.go:204-212`

## Diagram 7 — Observability Counters

```
Observability (atomic.Int64)
  ├── LegacyReadHits  ← GetLegacyUser
  ├── NewReadHits     ← GetModernUser
  ├── DualWriteCount  ← CreateUser(WriteDual)
  ├── DualWriteErrors ← CreateDual failure
  ├── BackfillProcessed ← BackfillWorker.RunBatch
  └── DriftDetected   ← ReconcileData
         │
         ▼ Snapshot() → map[string]int64
              used by: ApplyContract guard, demo metrics, tests
```

Source: `internal/compat/metrics.go:7-28`, `internal/compat/service.go`

## Diagram 8 — HTTP Deprecation Flow

```
GET /api/v1/users?id=1              GET /api/v2/users?id=1
        │                                     │
        ▼                                     ▼
  GetLegacyUser                         GetModernUser
        │                                     │
   ┌────┴────┐                          ┌──────┴──────┐
   │contract?│                          │   success   │
   └──┬───┬──┘                          └──────┬──────┘
   yes│   │no                                  │ 200 OK
      ▼   ▼                                   ▼ JSON
  410 Gone  200 OK + headers:            ModernConsumerDTO
  {"error":  Deprecation: true
   "legacy   Sunset: Mon, 31 Dec 2026
    removed"} Content-Type: application/json
             LegacyConsumerDTO
```

Source: `internal/compat/handler.go:20-56`, `internal/compat/service.go:150-184`
