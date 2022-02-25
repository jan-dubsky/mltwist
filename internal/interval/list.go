package interval

import (
	"decomp/internal/addr"
	"fmt"
	"sort"
)

type List struct {
	intervals []interval
}

func NewList(ints []Interval) (*List, error) {
	if len(ints) == 0 {
		return &List{}, nil
	}

	intervals := make([]interval, len(ints))
	for i, in := range ints {
		intervals[i] = newInterval(in)
	}

	less := func(i, j int) bool { return intervals[i].begin < intervals[j].begin }
	sort.Slice(intervals, less)

	for i := range intervals[1:] {
		if e, b := intervals[i].end, intervals[i+1].begin; e > b {
			return nil, fmt.Errorf(
				"block %d (ending 0x%x) and %d (starting 0x%x) overlap",
				i, e, i+1, b)
		}
	}

	return &List{intervals: intervals}, nil
}

func (l *List) Addr(addr addr.Address) Interval {
	idx := sort.Search(len(l.intervals), func(i int) bool {
		return l.intervals[i].end > addr
	})

	if idx == len(l.intervals) || l.intervals[idx].begin > addr {
		return nil
	}

	return l.intervals[idx].interval
}
