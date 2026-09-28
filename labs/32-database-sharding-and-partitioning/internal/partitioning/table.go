package partitioning

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrNoMatchingPartition = errors.New("no matching partition found for record")
)

type Record struct {
	ID        string
	CreatedAt time.Time
	Payload   string
}

type PartitionRange struct {
	Name  string
	Start time.Time // inclusive
	End   time.Time // exclusive
}

type Partition struct {
	Range PartitionRange
	mu    sync.RWMutex
	rows  []Record
}

func (p *Partition) Insert(rec Record) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.rows = append(p.rows, rec)
}

func (p *Partition) Count() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.rows)
}

func (p *Partition) Rows() []Record {
	p.mu.RLock()
	defer p.mu.RUnlock()
	copied := make([]Record, len(p.rows))
	copy(copied, p.rows)
	return copied
}

// Table represents a logically partitioned table residing on a single engine.
type Table struct {
	Name       string
	mu         sync.RWMutex
	partitions []*Partition
}

func NewTable(name string, ranges []PartitionRange) *Table {
	t := &Table{
		Name:       name,
		partitions: make([]*Partition, len(ranges)),
	}
	for i, r := range ranges {
		t.partitions[i] = &Partition{
			Range: r,
			rows:  make([]Record, 0),
		}
	}
	return t
}

func (t *Table) Insert(rec Record) error {
	t.mu.RLock()
	defer t.mu.RUnlock()

	for _, p := range t.partitions {
		if (rec.CreatedAt.Equal(p.Range.Start) || rec.CreatedAt.After(p.Range.Start)) && rec.CreatedAt.Before(p.Range.End) {
			p.Insert(rec)
			return nil
		}
	}
	return fmt.Errorf("%w: timestamp %s", ErrNoMatchingPartition, rec.CreatedAt.Format(time.RFC3339))
}

type QueryResult struct {
	Records         []Record
	PartitionsScanned int
	TotalPartitions   int
}

// QueryRange demonstrates partition pruning: only partitions overlapping [start, end) are scanned.
func (t *Table) QueryRange(start, end time.Time) QueryResult {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var records []Record
	scanned := 0

	for _, p := range t.partitions {
		// Overlap condition: start < p.Range.End && end > p.Range.Start
		if start.Before(p.Range.End) && end.After(p.Range.Start) {
			scanned++
			p.mu.RLock()
			for _, r := range p.rows {
				if (r.CreatedAt.Equal(start) || r.CreatedAt.After(start)) && r.CreatedAt.Before(end) {
					records = append(records, r)
				}
			}
			p.mu.RUnlock()
		}
	}

	return QueryResult{
		Records:         records,
		PartitionsScanned: scanned,
		TotalPartitions:   len(t.partitions),
	}
}

func (t *Table) DropPartition(name string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	for i, p := range t.partitions {
		if p.Range.Name == name {
			t.partitions = append(t.partitions[:i], t.partitions[i+1:]...)
			return true
		}
	}
	return false
}

func (t *Table) PartitionCount() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.partitions)
}
