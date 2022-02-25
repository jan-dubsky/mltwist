package interval

import (
	"decomp/internal/addr"
	"fmt"
)

type Interval interface {
	Begin() addr.Address
	End() addr.Address
}

type interval struct {
	interval Interval
	begin    addr.Address
	end      addr.Address
}

func newInterval(i Interval) interval {
	return interval{
		interval: i,
		begin:    i.Begin(),
		end:      i.End(),
	}
}

func (i interval) Begin() addr.Address { return i.begin }
func (i interval) End() addr.Address   { return i.end }

func CheckOverlaps(ints []Interval) error {
	for i := range ints[1:] {
		e := ints[i].End()
		b := ints[i+1].Begin()
		if e <= b {
			continue
		}

		return fmt.Errorf("block %d (ending 0x%x) and %d (starting 0x%x) overlap",
			i, e, i+1, b)
	}

	return nil
}
