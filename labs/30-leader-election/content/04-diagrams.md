# Diagrams

## Diagram 1 — Component Architecture

Sumber: `engineering/01-design.md` dan `README.md`.

```text
┌─────────────────────┐       ┌─────────────────────┐
│   Candidate /       │       │   Candidate /       │
│   Node A            │       │   Node B            │
│   (FOLLOWER/LEADER) │       │   (FOLLOWER/LEADER) │
└────────┬────────────┘       └────────┬────────────┘
         │                             │
         │   Acquire / Renew / Release │
         ▼                             ▼
  ┌──────────────────────────────────────────────┐
  │           Coordinator (In-Memory)            │
  │                                              │
  │  •leases map[string]*Lease                   │
  │  •revision int64 (monotonic fencing token)   │
  │  •sync.Mutex (linearizable operations)       │
  │  •TTL expiration via time.Now() monotonic     │
  └──────────────────────┬───────────────────────┘
                         │
                         │  Write(token, payload)
                         ▼
  ┌──────────────────────────────────────────────┐
  │         Shared Resource (FencedStorage)      │
  │                                              │
  │  •lastSeenToken: int64                       │
  │  •records: []Record                          │
  │  •sync.RWMutex (check-and-set atomik)        │
  │  •Tolak: incomingToken <= lastSeenToken      │
  │  •Terima: incomingToken > lastSeenToken      │
  └──────────────────────────────────────────────┘
```

## Diagram 2 — Leader Lifecycle

Sumber: `internal/candidate/candidate.go`.

```text
                    ┌─────────────────────┐
                    │   Candidate Start   │
                    │   State = FOLLOWER  │
                    └──────────┬──────────┘
                               │
                               ▼
                    ┌─────────────────────┐
           ┌───── │  Election Loop Tick  │ ◄──────────────┐
           │      └──────────┬───────────┘                │
           │                 │                            │
           │    ┌────────────┴────────────┐               │
           │    ▼                         ▼               │
           │  ┌─────────┐          ┌───────────┐          │
           │  │FOLLOWER │          │  LEADER   │          │
           │  └────┬────┘          └─────┬─────┘          │
           │       │                     │                │
           │       ▼                     ▼                │
           │  ┌──────────┐         ┌──────────┐           │
           │  │ Acquire  │         │  Renew   │           │
           │  │(try lock)│         │(heartbeat│           │
           │  └────┬─────┘         └────┬─────┘           │
           │       │                    │                 │
           │  ┌────┴────┐          ┌────┴────┐            │
           │  │ Success │          │ Success │            │
           │  │→ LEADER │          │→ tetap  │──────┐     │
           │  └─────────┘          └─────────┘      │     │
           │       │                    │           │     │
           │       │                    ▼           │     │
           │       │            ┌──────────────┐   │     │
           │       │            │  Kegagalan   │   │     │
           │       │            │  (ErrExpired │   │     │
           │       │            │  /ErrHeld)   │   │     │
           │       │            └──────┬───────┘   │     │
           │       │                   │           │     │
           │       ▼                   ▼           │     │
           │  ┌───────────────────────────────┐   │     │
           │  │     State = FOLLOWER          │   │     │
           │  │     currentLease = nil        │   │     │
           │  └───────────────┬───────────────┘   │     │
           │                  │                   │     │
           │                  └───────────────────┘     │
           │                                            │
           │   ┌───────────────────────┐                │
           │   │ now < pauseUntil ?    │                │
           │   │ YES → skip tick       │── tick dilewati
           │   └───────────────────────┘                │
           │                                            │
           │   ┌───────────────────────┐                │
           │   │ ctx.Done()?           │                │
           │   │ YES → return          │                │
           │   └───────────────────────┘                │
           └────────────────────────────────────────────┘
```

## Diagram 3 — Fencing Token Flow (Split-Brain Defense)

Sumber: `engineering/01-design.md` dan `engineering/03-execution-result.md`.

```text
┌──────────┐                    ┌──────────────┐                    ┌──────────────────┐
│  Node A  │                    │ Coordinator  │                    │ FencedStorage    │
│ (Leader) │                    │              │                    │                  │
└────┬─────┘                    └──────┬───────┘                    └────────┬─────────┘
     │                                 │                                     │
     │  Acquire(key, ttl=200ms)        │                                     │
     │────────────────────────────────▶│                                     │
     │                                 │ revision++ = 1                      │
     │  ◀──────────────────────────────│ return Lease{token:1, ttl:200ms}    │
     │                                 │                                     │
     │  Write(value, token=1)          │                                     │
     │────────────────────────────────────────────────────────────────────────▶│
     │                                 │              token 1 > lastSeen 0    │
     │  ◀─────────────────────────────────────────────────────────────────────│
     │                                 │              lastSeenToken = 1      │
     │                                 │              [SUCCES]               │
     │                                 │                                     │
     │  *** PAUSE 350ms (> TTL) ***    │                                     │
     │────────────────── X ───────────▶│  (heartbeat terlewat)              │
     │                                 │                                     │
     │              Lease kedaluwarsa  │                                     │
     │              (200ms lewat)      │                                     │
     │                                 │                                     │
     │                                 │                                     │
┌────┴─────┐                           │                                     │
│  Node B  │                           │                                     │
│(FOLLOWER)│                           │                                     │
└────┬─────┘                           │                                     │
     │  Acquire(key, ttl=200ms)        │                                     │
     │────────────────────────────────▶│                                     │
     │                                 │ revision++ = 2                      │
     │  ◀──────────────────────────────│ return Lease{token:2, ttl:200ms}    │
     │                                 │                                     │
     │  Write(value, token=2)          │                                     │
     │────────────────────────────────────────────────────────────────────────▶│
     │                                 │              token 2 > lastSeen 1    │
     │  ◀─────────────────────────────────────────────────────────────────────│
     │                                 │              lastSeenToken = 2      │
     │                                 │              [SUCCES]               │
     │                                 │                                     │
┌────┴─────┐                           │                                     │
│  Node A  │                           │                                     │
│(terbangun)│                          │                                     │
│token=1   │                           │                                     │
└────┬─────┘                           │                                     │
     │  Write(value, token=1)          │                                     │
     │────────────────────────────────────────────────────────────────────────▶│
     │                                 │              token 1 <= lastSeen 2   │
     │  ◀─────────────────────────────── ErrStaleFencingToken ───────────────│
     │              [DITOLAK: split-brain dicegah]                            │
     │                                 │              lastSeenToken tetap 2  │
```

## Diagram 4 — Lease State Transitions

Sumber: `internal/coordinator/coordinator.go`.

```text
         ┌──────────────────────────────────────────────────────────┐
         │                     COORDINATOR LEASE MAP                │
         │                                                          │
         │  Key "leader-lock"                                       │
         │  ┌─────────────────────────────────────────────────┐     │
         │  │                                                  │     │
         │  │  ┌───────────┐    Acquire (holder A, ttl)       │     │
         │  │  │  TIDAK    │ ──────────────────────────────▶  │     │
         │  │  │  ADA      │    revision++ = 1                │     │
         │  │  │           │    Lease{token:1, expiry: now+TTL}│    │
         │  │  └───────────┘                                   │     │
         │  │                                                  │     │
         │  │  ┌───────────┐    Acquire (holder B, ttl)       │     │
         │  │  │  ADA      │ ──────────────────────────────▶  │     │
         │  │  │ (valid)   │    ErrLeaseHeld                   │     │
         │  │  │           │    (lease masih valid)            │     │
         │  │  └───────────┘                                   │     │
         │  │                                                  │     │
         │  │  ┌───────────┐    Acquire (holder B, ttl)       │     │
         │  │  │  KEDALU-  │ ──────────────────────────────▶  │     │
         │  │  │  WARSA    │    revision++ = 2                │     │
         │  │  │           │    Lease{token:2, expiry: now+TTL}│    │
         │  │  └───────────┘                                   │     │
         │  │                                                  │     │
         │  │  ┌───────────┐    Release (holder B, token 2)   │     │
         │  │  │  ADA      │ ──────────────────────────────▶  │     │
         │  │  │ (valid)   │    delete(leases[key])            │     │
         │  │  │           │    (lease dihapus)               │     │
         │  │  └───────────┘                                   │     │
         │  └─────────────────────────────────────────────────┘     │
         └──────────────────────────────────────────────────────────┘
```
