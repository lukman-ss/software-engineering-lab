package compat

import (
	"sync/atomic"
)

type Observability struct {
	LegacyReadHits  atomic.Int64
	NewReadHits     atomic.Int64
	DualWriteCount  atomic.Int64
	DualWriteErrors atomic.Int64
	BackfillProcessed atomic.Int64
	DriftDetected   atomic.Int64
}

func NewObservability() *Observability {
	return &Observability{}
}

func (o *Observability) Snapshot() map[string]int64 {
	return map[string]int64{
		"legacy_reads":      o.LegacyReadHits.Load(),
		"new_reads":         o.NewReadHits.Load(),
		"dual_writes":       o.DualWriteCount.Load(),
		"dual_write_errors": o.DualWriteErrors.Load(),
		"backfilled":        o.BackfillProcessed.Load(),
		"drift_detected":    o.DriftDetected.Load(),
	}
}
