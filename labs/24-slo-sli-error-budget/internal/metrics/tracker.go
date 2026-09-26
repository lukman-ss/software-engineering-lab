package metrics

import (
	"sync"
	"time"
)

type Event struct {
	Timestamp  time.Time
	Duration   time.Duration
	StatusCode int
	Endpoint   string
}

type Bucket struct {
	StartTime  time.Time
	TotalCount int64
	GoodCount  int64
	BadCount   int64
}

type WindowTracker struct {
	mu           sync.RWMutex
	windowSize   time.Duration
	bucketSize   time.Duration
	buckets      []Bucket
	isGoodEvent  func(e Event) bool
}

func NewWindowTracker(windowSize time.Duration, bucketSize time.Duration, isGood func(e Event) bool) *WindowTracker {
	if bucketSize <= 0 {
		bucketSize = time.Second
	}
	numBuckets := int(windowSize / bucketSize)
	if numBuckets < 1 {
		numBuckets = 1
	}
	return &WindowTracker{
		windowSize:  windowSize,
		bucketSize:  bucketSize,
		buckets:     make([]Bucket, 0, numBuckets),
		isGoodEvent: isGood,
	}
}

func (w *WindowTracker) Record(e Event) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.evictStaleLocked(e.Timestamp)

	bucketStart := e.Timestamp.Truncate(w.bucketSize)
	good := w.isGoodEvent(e)

	n := len(w.buckets)
	if n > 0 && w.buckets[n-1].StartTime.Equal(bucketStart) {
		w.buckets[n-1].TotalCount++
		if good {
			w.buckets[n-1].GoodCount++
		} else {
			w.buckets[n-1].BadCount++
		}
		return
	}

	// If timestamp belongs to an earlier existing bucket or should be inserted in order
	if n > 0 && bucketStart.Before(w.buckets[n-1].StartTime) {
		for i := 0; i < n; i++ {
			if w.buckets[i].StartTime.Equal(bucketStart) {
				w.buckets[i].TotalCount++
				if good {
					w.buckets[i].GoodCount++
				} else {
					w.buckets[i].BadCount++
				}
				return
			}
			if w.buckets[i].StartTime.After(bucketStart) {
				b := Bucket{
					StartTime:  bucketStart,
					TotalCount: 1,
				}
				if good {
					b.GoodCount = 1
				} else {
					b.BadCount = 1
				}
				w.buckets = append(w.buckets[:i], append([]Bucket{b}, w.buckets[i:]...)...)
				return
			}
		}
	}

	b := Bucket{
		StartTime:  bucketStart,
		TotalCount: 1,
	}
	if good {
		b.GoodCount = 1
	} else {
		b.BadCount = 1
	}
	w.buckets = append(w.buckets, b)
}

func (w *WindowTracker) evictStaleLocked(now time.Time) {
	cutoff := now.Add(-w.windowSize)
	idx := 0
	for idx < len(w.buckets) && w.buckets[idx].StartTime.Before(cutoff) {
		idx++
	}
	if idx > 0 {
		w.buckets = w.buckets[idx:]
	}
}

func (w *WindowTracker) Summary(now time.Time) (total int64, good int64, bad int64) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.evictStaleLocked(now)
	for _, b := range w.buckets {
		total += b.TotalCount
		good += b.GoodCount
		bad += b.BadCount
	}
	return total, good, bad
}
