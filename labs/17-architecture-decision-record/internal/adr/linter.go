package adr

import (
	"fmt"
	"sort"
	"sync"
)

type Linter struct{}

func NewLinter() *Linter {
	return &Linter{}
}

func (l *Linter) Validate(records []*Record) []error {
	var errs []error
	var mu sync.Mutex

	recordMap := make(map[int]*Record)
	for _, r := range records {
		if _, exists := recordMap[r.ID]; exists {
			errs = append(errs, fmt.Errorf("duplicate ADR ID: %d", r.ID))
		}
		recordMap[r.ID] = r
	}

	// Validate IDs are monotonic
	ids := make([]int, 0, len(records))
	for id := range recordMap {
		ids = append(ids, id)
	}
	sort.Ints(ids)

	for i := 0; i < len(ids); i++ {
		if ids[i] != i+1 {
			errs = append(errs, fmt.Errorf("non-monotonic numbering, expected %d but got %d", i+1, ids[i]))
			break // Only report once
		}
	}

	// Validate graph in parallel
	var wg sync.WaitGroup
	for _, r := range records {
		wg.Add(1)
		go func(rec *Record) {
			defer wg.Done()
			var localErrs []error

			if rec.Status == StatusSuperseded {
				if rec.SupersededBy == 0 {
					localErrs = append(localErrs, fmt.Errorf("ADR %d is superseded but missing superseded_by reference", rec.ID))
				} else {
					replacement, exists := recordMap[rec.SupersededBy]
					if !exists {
						localErrs = append(localErrs, fmt.Errorf("ADR %d superseded by non-existent ADR %d", rec.ID, rec.SupersededBy))
					} else if replacement.Supersedes != rec.ID {
						localErrs = append(localErrs, fmt.Errorf("ADR %d superseded by ADR %d, but ADR %d does not declare it supersedes ADR %d", rec.ID, rec.SupersededBy, rec.SupersededBy, rec.ID))
					}
				}
			}

			if rec.Supersedes != 0 {
				old, exists := recordMap[rec.Supersedes]
				if !exists {
					localErrs = append(localErrs, fmt.Errorf("ADR %d supersedes non-existent ADR %d", rec.ID, rec.Supersedes))
				} else if old.Status != StatusSuperseded || old.SupersededBy != rec.ID {
					localErrs = append(localErrs, fmt.Errorf("ADR %d supersedes ADR %d, but ADR %d is not properly marked as superseded by ADR %d", rec.ID, rec.Supersedes, rec.Supersedes, rec.ID))
				}
			}

			if len(localErrs) > 0 {
				mu.Lock()
				errs = append(errs, localErrs...)
				mu.Unlock()
			}
		}(r)
	}

	wg.Wait()

	return errs
}
