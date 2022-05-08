package interval

import (
	"sort"

	"golang.org/x/exp/constraints"
)

// Map represents a set of non-overlapping intervals.
type Map[T constraints.Integer] struct {
	intvs []Interval[T]
}

// NewMap creates new intervals from intvs. Intervals which follow one another
// are merged into a single interval. Overlaps on intervals in intvs are
// ignored.
func NewMap[T constraints.Integer](intvs ...Interval[T]) Map[T] {
	sort.Slice(intvs, func(i, j int) bool { return intvs[i].begin < intvs[j].begin })

	j := 0
	for i := 0; i < len(intvs); i, j = i+1, j+1 {
		if i == 0 || intvs[j-1].end < intvs[i].begin {
			intvs[j] = intvs[i]
			continue
		}

		if intvs[j-1].end >= intvs[i].end {
			j--
			continue
		}

		intvs[j-1].end = intvs[i].end
		j--
	}

	return Map[T]{intvs: intvs[:j]}
}

// Idx returns idxth element in i. This function panics if idx < 0 or idx <
// i.Len().
func (m Map[T]) Index(idx int) Interval[T] { return m.intvs[idx] }

// Len returns number of intervals stored in i.
func (m Map[T]) Len() int { return len(m.intvs) }

// Intervals returns sorted array of intervals in i. This array has to be
// treated as readonly. Modification of the array returned can tresult in
// undefined behaviour.
func (m Map[T]) Intervals() []Interval[T] { return m.intvs }

// Equal compares m1 to m2 and informs if they are identical. Interval maps are
// indencital if they contain the same intervals.
func (m1 Map[T]) Equal(m2 Map[T]) bool {
	if m1.Len() != m2.Len() {
		return false
	}

	for i := 0; i < m1.Len(); i++ {
		if m1.Index(i) != m2.Index(i) {
			return false
		}
	}

	return true
}

func addInterval[T constraints.Integer](is []Interval[T], i Interval[T]) []Interval[T] {
	if len(is) == 0 {
		return append(is, i)
	}

	last := is[len(is)-1]
	if last.end < i.begin {
		return append(is, i)
	}

	if i.end > last.end {
		is[len(is)-1].end = i.end
		return is
	}

	return is
}

// MapUnion produces new Map spanning both ranges from i1 and i2. Overlapping
// intervals are ignored, but continuous intervals are merged into a single one.
func MapUnion[T constraints.Integer](i1 Map[T], i2 Map[T]) Map[T] {
	added := make([]Interval[T], 0, i1.Len()+i2.Len())

	i, j := 0, 0
	for i < i1.Len() && j < i2.Len() {
		var intv Interval[T]
		if intv1, intv2 := i1.Index(i), i2.Index(j); intv1.begin < intv2.begin {
			intv = intv1
			i++
		} else {
			intv = intv2
			j++
		}

		added = addInterval(added, intv)
	}

	if i < i1.Len() {
		for _, intv := range i1.intvs[i:] {
			added = addInterval(added, intv)
		}
	} else {
		for _, intv := range i2.intvs[j:] {
			added = addInterval(added, intv)
		}
	}

	return Map[T]{intvs: added}
}

func sub[T constraints.Integer](intv Interval[T], sub []Interval[T]) ([]Interval[T], int) {
	intvs := make([]Interval[T], 0, 1)

	var cnt int
	for _, s := range sub {
		cnt++

		if s.end <= intv.begin {
			continue
		}

		if intv.end < s.begin {
			break
		}

		if intv.begin < s.begin {
			intvs = append(intvs, New(intv.begin, s.begin))
		}

		if s.end < intv.end {
			intv.begin = s.end
		} else {
			return intvs, cnt - 1
		}
	}

	intvs = append(intvs, intv)
	return intvs, cnt - 1
}

// MapIntersect creates new Map spanning those ranges from i1 which are not
// present in i2.
func MapIntersect[T constraints.Integer](i1 Map[T], i2 Map[T]) Map[T] {
	subtracted := make([]Interval[T], 0, i1.Len())

	j := 0
	for i := 0; i < i1.Len(); i++ {
		var subList []Interval[T]
		if j < i2.Len() {
			subList = i2.intvs[j:]
		}

		intvs, cnt := sub(i1.Index(i), subList)
		j += cnt

		subtracted = append(subtracted, intvs...)
	}

	return Map[T]{intvs: subtracted}
}
