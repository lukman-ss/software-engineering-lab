package router

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"labs/33-read-replicas-and-replication-lag/internal/cluster"
)

var (
	ErrNoReplicaAvailable = errors.New("no replica available within acceptable lag threshold")
	ErrLSNWaitTimeout     = errors.New("timeout waiting for replica to reach required LSN")
)

type Config struct {
	StickyDuration time.Duration
	MaxAllowedLag  time.Duration
	MaxLSNDiff     uint64
	WaitTimeout    time.Duration
}

func DefaultConfig() Config {
	return Config{
		StickyDuration: 5 * time.Second,
		MaxAllowedLag:  2 * time.Second,
		MaxLSNDiff:     5,
		WaitTimeout:    1 * time.Second,
	}
}

type SessionState struct {
	LastWriteTime time.Time
	LastWriteLSN  uint64
}

type Router struct {
	cluster     *cluster.Cluster
	config      Config
	sessions    sync.Map
	rrIndex     atomic.Uint64
}

func NewRouter(c *cluster.Cluster, cfg Config) *Router {
	return &Router{
		cluster: c,
		config:  cfg,
	}
}

func (r *Router) Write(sessionID, key, value string) (uint64, error) {
	lsn, err := r.cluster.Write(key, value)
	if err != nil {
		return 0, err
	}

	r.sessions.Store(sessionID, SessionState{
		LastWriteTime: time.Now(),
		LastWriteLSN:  lsn,
	})

	return lsn, nil
}

func (r *Router) ReadNaive(key string) (string, string, uint64, error) {
	replicas := r.cluster.Replicas()
	if len(replicas) == 0 {
		val, lsn, err := r.cluster.Primary().Read(key)
		return val, r.cluster.Primary().ID(), lsn, err
	}

	idx := r.rrIndex.Add(1) % uint64(len(replicas))
	target := replicas[idx]
	val, lsn, err := target.Read(key)
	return val, target.ID(), lsn, err
}

func (r *Router) ReadWithStickySession(sessionID, key string) (string, string, uint64, error) {
	if stateVal, ok := r.sessions.Load(sessionID); ok {
		state := stateVal.(SessionState)
		if time.Since(state.LastWriteTime) < r.config.StickyDuration {
			val, lsn, err := r.cluster.Primary().Read(key)
			return val, r.cluster.Primary().ID(), lsn, err
		}
	}

	return r.ReadLagAware(key)
}

func (r *Router) ReadWithToken(ctx context.Context, minLSN uint64, key string) (string, string, uint64, error) {
	replicas := r.cluster.Replicas()
	
	for _, replica := range replicas {
		if replica.AppliedLSN() >= minLSN {
			val, lsn, err := replica.Read(key)
			return val, replica.ID(), lsn, err
		}
	}

	if len(replicas) > 0 {
		targetReplica := replicas[0]
		waitCtx, cancel := context.WithTimeout(ctx, r.config.WaitTimeout)
		defer cancel()

		if err := targetReplica.WaitForLSN(waitCtx, minLSN); err == nil {
			val, lsn, err := targetReplica.Read(key)
			return val, targetReplica.ID(), lsn, err
		}
	}

	val, lsn, err := r.cluster.Primary().Read(key)
	return val, r.cluster.Primary().ID() + " (fallback)", lsn, err
}

func (r *Router) ReadLagAware(key string) (string, string, uint64, error) {
	replicas := r.cluster.Replicas()
	primaryLSN := r.cluster.CurrentLSN()

	var eligible []*cluster.Node
	for _, replica := range replicas {
		appliedLSN := replica.AppliedLSN()
		if primaryLSN >= appliedLSN {
			diff := primaryLSN - appliedLSN
			if diff <= r.config.MaxLSNDiff {
				eligible = append(eligible, replica)
			}
		}
	}

	if len(eligible) == 0 {
		val, lsn, err := r.cluster.Primary().Read(key)
		return val, r.cluster.Primary().ID() + " (fallback-lag)", lsn, err
	}

	idx := r.rrIndex.Add(1) % uint64(len(eligible))
	target := eligible[idx]
	val, lsn, err := target.Read(key)
	return val, target.ID(), lsn, err
}
