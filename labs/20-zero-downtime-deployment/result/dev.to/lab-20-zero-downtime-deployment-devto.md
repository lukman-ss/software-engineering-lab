---
title: "Zero‑Downtime Deployment: Graceful Shutdown, Traffic Draining, & Expand‑Contract DB"
published: false
description: "Explore Zero‑Downtime Deployment in Go: mechanics, architecture, caveats, and testable invariants."
tags: go, testing, performance, softwareengineering
---

# Zero‑Downtime Deployment: Graceful Shutdown, Traffic Draining, & Expand‑Contract DB

A production system must remain resilient under load. Lab 20 makes that behavior visible with a focused Go simulation and deterministic test assertions.

> In‑flight traffic dan background work harus selesai sebelum pod dimatikan.

## Deploy tanpa draining memutus request aktif, jobs terhenti, dan schema berubah.

Interrupt HTTP: Pengguna menerima 502/504 selama pod terminta Job loss: Background worker terhenti meninggalkan state parsial Schema clash: Versi lama tidak kompatibel dengan kolom baru

## Pengujian sequential tidak meniru interleaving runtime.

Tidak drain: Server menolak request baru tapi tidak menunggu yang berjalan Tidak sync: Tidak ada preStop, load balancer masih mengirim traffic DB lock: Schema berubah tiba‑tiba menyebabkan error

## Start New Before Stopping Old

New version launches (Ready): Load balancer routes to new pods Old version drains (PreStop): Finishes in‑flight requests & jobs Old version terminates (Terminate): Database schema contract completed

## Alur Zero‑Downtime Deployment

Init (Init): Component start, readiness false Ready (Ready): Health probe 200, traffic starts SIGTERM (PreStop): PreStop delay, drain HTTP & worker Terminate (Exit): Pod dies, new version continues

## Server graceful shutdown dengan preStop

```go
func (s *Server) Shutdown(ctx context.Context) error {
    s.SetReady(false)
    if s.preStop > 0 {
        select { case <-time.After(s.preStop): case <-ctx.Done(): return ctx.Err() }
    }
    return s.srv.Shutdown(ctx)
}
```

## Tanpa preStop, load balancer masih mengirim request → 502.

1. SIGTERM: Pod dies instantly 2. In‑flight request: Client receives connection reset 3. Worker job: Job terpotong, state inconsistency

## Graceful vs Naïve Shutdown

Aspect | Graceful (preStop) | Naïve
|---|---|---|
Traffic | Drained | Dropped
Jobs | Completed | Aborted
DB Schema | Compatible | Conflict

## Tiga langkah mitigasi

1. Health Probe: Gunakan readiness untuk mengisolasi pod 2. PreStop delay: Setidaknya latency routing table 3. Expand‑Contract DB: Baca/tulis versi lama & baru simultan

## Terminal output demo

```sh
go run ./cmd/demo
... SIGTERM received, graceful shutdown completed
```

## Batasan bukti

Dibuktikan: Graceful shutdown, preStop delay, worker drain Tidak dibuktikan: Distributed load balancer latency >1s, real DB DDL lock contention

## PreStop vs Fast termination

PreStop: Safety, extra latency Fast kill: Zero latency but risk 502/504

## Worker drain dengan WaitGroup

```go
func (w *Worker) Stop(timeout time.Duration) {
    w.enqueueMu.Lock(); w.stopped.Store(true); close(w.jobChan); w.enqueueMu.Unlock()
    done := make(chan struct{})
    go func(){ w.wg.Wait(); close(done) }()
    select { case <-done: case <-time.After(timeout): w.cancel() }
}
```

## Sebelum vs Sesudah graceful shutdown

Sebelum: Request drop, job loss, schema error Sesudah: All in‑flight completed, no DB conflict

## Analisis: Health Probe Timing

Observation: ... Implication: ...

## Analisis: PreStop Duration Impact

Observation: ... Implication: ...

## Analisis: Worker Drain Guarantees

Observation: ... Implication: ...

## Analisis: DB Expand‑Contract Compatibility

Observation: ... Implication: ...

## Analisis: End‑to‑End Latency

Observation: ... Implication: ...

## Checklist produksi Zero‑Downtime

Implement readiness probe returning 503 until warmup Configure preStop hook >= external routing latency Use Expand‑Contract DB schema with backward compatible reads/writes Ensure worker can drain with timeout and context cancellation

## Mengabaikan preStop & readiness

Force kill: Mengakibatkan 502/504 Schema breaking: Deploy kolom baru tanpa fallback Job abort: Background worker terhenti abrupt

## Perintah demo lokal

```sh
cd labs/20-zero-downtime-deployment
go test -v ./...
go run ./cmd/demo
```

## Desain untuk koeksistensi versi

Versioned APIs: Maintain backward compatibility Graceful shutdown: Health probes + preStop Schema evolution: Expand‑Contract pattern

## Source

[Software Engineering Lab 20 - Zero‑Downtime Deployment](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/20-zero-downtime-deployment), snapshot commit `<auto>`. Source: README, implementation, tests, and execution records. Code excerpts belong in their original package context.
