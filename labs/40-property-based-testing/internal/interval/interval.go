package interval

import "sort"

// Interval represents a closed interval [Start, End].
type Interval struct {
	Start, End int
}

// NaiveMerge merges intervals but fails on unsorted inputs.
// Bug: assumes input is already sorted.
func NaiveMerge(intervals []Interval) []Interval {
	if len(intervals) == 0 {
		return nil
	}
	out := []Interval{intervals[0]}
	for _, iv := range intervals[1:] {
		last := &out[len(out)-1]
		if iv.Start <= last.End {
			if iv.End > last.End {
				last.End = iv.End
			}
		} else {
			out = append(out, iv)
		}
	}
	return out
}

// RobustMerge merges intervals correctly regardless of input order.
func RobustMerge(intervals []Interval) []Interval {
	if len(intervals) == 0 {
		return nil
	}
	sorted := make([]Interval, len(intervals))
	copy(sorted, intervals)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Start != sorted[j].Start {
			return sorted[i].Start < sorted[j].Start
		}
		return sorted[i].End < sorted[j].End
	})

	out := []Interval{sorted[0]}
	for _, iv := range sorted[1:] {
		last := &out[len(out)-1]
		if iv.Start <= last.End+1 {
			if iv.End > last.End {
				last.End = iv.End
			}
		} else {
			out = append(out, iv)
		}
	}
	return out
}

func intervalsEqual(a, b []Interval) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
