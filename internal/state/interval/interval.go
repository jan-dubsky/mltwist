package interval

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

// Interval represents a single interval [begin, end). Interval is on purpose
// designed to be immutable.
type Interval[T constraints.Integer] struct {
	begin T
	end   T
}

// New creates a new interval spanning range [begin, end). This function panics
// if begin is greater than end.
func New[T constraints.Integer](begin, end T) Interval[T] {
	if begin > end {
		panic(fmt.Sprintf("begin is greater than end: %d > %d", begin, end))
	}

	return Interval[T]{
		begin: begin,
		end:   end,
	}
}

// Begin returns inclusive start of i.
func (i Interval[T]) Begin() T { return i.begin }

// End returns exclusive end of i.
func (i Interval[T]) End() T { return i.end }

// Len returns length of i as End() - Begin().
func (i Interval[T]) Len() T { return i.end - i.begin }

// Containts checks if val belongs to the interval.
func (i Interval[T]) Containts(val T) bool { return i.begin <= val && val < i.end }
