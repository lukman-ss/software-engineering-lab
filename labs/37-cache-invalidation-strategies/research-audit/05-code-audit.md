# 05 — Code Audit

## PIPELINE OVERRIDE NOTICE

Code and implementation audit has been explicitly excluded from this pipeline stage.

Per the pipeline directive:
- "Audit research only."
- "Do not audit implementation/code in this stage."

No implementation, tests, or runnable code was inspected in this audit pass.

No build commands were executed.

## Scope

Code Audit Status: NOT APPLICABLE (Pipeline Override)

## Notes for Next Stage

The following implementation-relevant concerns identified during research audit should be verified when code audit is enabled:

1. **XFetch formula sign error**: Lab spec formula `Δ·β·ln(rand()) > TTL_remaining` is ALWAYS FALSE for `rand() ∈ (0,1)`. Any implementation using this formula verbatim will be functionally broken (early refresh never triggers). Correct form: `-Δ·β·ln(rand()) > TTL_remaining` or equivalently `Δ·β·(-ln(rand())) > TTL_remaining`.

2. **singleflight scope**: Go `singleflight.Group` is process-local; multi-replica deployments require a separate distributed coordination mechanism.

3. **Write-through atomicity**: No distributed ACID guarantee between DB and Redis writes; cache.Set failure after DB commit leaves a window for stale reads until TTL expiry.

4. **SWR background job trigger**: Application-layer stale-while-revalidate implementations should tie revalidation to incoming request context, not autonomous background goroutines.

5. **Redis TTL retrieval for XFetch**: The XFetch algorithm requires reading `expiry` time from the cache. Redis `TTL key` command returns integer seconds; code must convert this to an absolute timestamp for the formula.
