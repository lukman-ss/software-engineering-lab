# Docs vs Code Audit

Target Lab: labs/37-cache-invalidation-strategies

## README vs Implementation

### README Claim: Cache-Aside
> Checks cache, loads from DB on miss, populates cache. Invalidates cache on update.

Code (`patterns.go:22-47`): `Get` checks `cache.Get(key)` -> falls back to `db.Query` -> `cache.Set`. `Update` calls `db.Write` first, then `cache.Delete`.

Assessment: MATCH

---

### README Claim: Write-Through
> Updates DB and cache synchronously on write; subsequent reads hit cache.

Code (`patterns.go:78-89`): `Update` calls `db.Write` -> `cache.Set` synchronously. Read after update hits cache without DB query.

Assessment: MATCH

---

### README Claim: Write-Behind
> Updates cache immediately; flushes to DB asynchronously via background worker queue.

Code (`patterns.go:152-165`): `Update` calls `cache.Set` immediately then `select { case s.writeQueue <- req: default: }` for async queuing.

Assessment: MATCH

---

### README Claim: SingleFlight
> Coalesces concurrent cache miss requests on a hot key down to a single DB query using `golang.org/x/sync/singleflight`.

Code (`stampede.go:56-84`): Uses `s.flight.Do(key, ...)` with `singleflight.Group`.

Assessment: MATCH

---

### README Claim: XFetch
> Early refresh using formula `-Δ · β · ln(U) > TTL_remaining` where `U ~ Uniform(0,1)`.

Code (`stampede.go:125-136`): `ShouldRecompute` computes `expiryCompute := -deltaSec * beta * math.Log(u)` and returns `expiryCompute > ttlRemainingSec`.

Assessment: MATCH

---

### README Claim: Stale-While-Revalidate (SWR)
> Returns stale cached data immediately while asynchronously triggering background DB revalidation.

Code (`stampede.go:197-223`): Stale-but-within-window path calls `triggerRevalidate` then returns `item.Value`. Background goroutine performs `db.Query` and updates cache.

Assessment: MATCH

---

### README Claim: TTL Jitter
> Adds randomized offset `[0, maxJitter)` to base TTL to prevent synchronized key expiry.

Code (`store.go:79-86`): `TTLWithJitter` returns `base + rand.Int63n(int64(maxJitter))`. Range is `[base, base+maxJitter)`.

Assessment: MATCH

---

## Engineering Design vs Implementation

### Design Note: `jitter.go` listed as separate file

Engineering `01-design.md` line 47 references `jitter.go` as a separate file:
> `internal/cache/jitter.go`: TTL jitter calculation.

Actual code: `TTLWithJitter` lives in `store.go`, not a separate `jitter.go`.

Assessment: DOC_CODE_MISMATCH (LOW severity — minor reference error in design document; implementation is correct and function is findable)

---

### Design vs Implementation: SWR Revalidation Timeout

Design (`01-design.md`) does not explicitly specify a timeout for background SWR revalidations.

Implementation (`stampede.go:244`): Background goroutine uses `context.WithTimeout(context.Background(), 10*time.Second)`.

Assessment: NOT_CLAIMED — timeout is an implementation detail, not a contradiction, and is a safe design choice.

---

## Execution Result vs Actual Output

Engineering `03-execution-result.md` records demo output. Audited live run produced identical structure with only minor TTL jitter numeric variation (expected non-determinism). All demo patterns behave identically.

Assessment: MATCH

---

## Summary

| Check | Result |
|---|---|
| README vs Code | MATCH (all 7 features) |
| Engineering Notes vs Code | MINOR DOC_CODE_MISMATCH (jitter.go reference) |
| Engineering Notes vs Execution Results | MATCH |
| Demo Output vs Live Execution | MATCH |
