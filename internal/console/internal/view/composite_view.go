package view

import (
	"fmt"
)

type CompositeView struct {
	elements []Element
}

func NewCompositeView(elements ...Element) *CompositeView {
	return &CompositeView{
		elements: elements,
	}
}

func (v *CompositeView) Print(lines int) error {
	minLines := v.MinLines()
	remainingLines := lines - v.MinLines()
	if remainingLines < 0 {
		return fmt.Errorf("not enough lines to render: %d < %d", lines, minLines)
	}

	lineCnts := v.distributeLines(remainingLines)
	for i, e := range v.elements {
		if i != 0 {
			fmt.Printf("\n")
		}

		err := e.Print(lineCnts[i])
		if err != nil {
			return fmt.Errorf("cannot print element %d/%d: %w", i,
				len(v.elements), err)
		}
	}

	return nil
}

// elementSpaces returns number of spaces in between individual elements. As
// every 2 consecutive elements are split by an empty line, the number returned
// is number of elements minus one.
func (v *CompositeView) elementSpaces() int { return len(v.elements) - 1 }

func (v *CompositeView) MinLines() int {
	var h int
	for _, p := range v.elements {
		h += p.MinLines()
	}

	return h + v.elementSpaces()
}

func (v *CompositeView) MaxLines() int {
	var h int
	for _, p := range v.elements {
		m := p.MaxLines()
		if m < 0 {
			return -1
		}

		h += m
	}

	return h + v.elementSpaces()
}

func (v *CompositeView) distributeLines(remLines int) map[int]int {
	mins := v.mins()

	lineCnts := make(map[int]int, len(mins))
	for i, min := range mins {
		lineCnts[i] = min
	}

	diffs := v.diffLinesMax(mins, remLines)
	positiveDiffs := countPositive(diffs)

	for i := 0; remLines > 0 && positiveDiffs > 0; i = (i + 1) % len(v.elements) {
		if diffs[i] <= 0 {
			continue
		}

		lineCnts[i]++
		diffs[i]--
		if diffs[i] == 0 {
			positiveDiffs--
		}
	}

	return lineCnts
}

func (v *CompositeView) mins() map[int]int {
	mins := make(map[int]int, len(v.elements))
	for i, e := range v.elements {
		val := e.MinLines()
		if val < 0 {
			val = 0
		}
		mins[i] = val
	}

	return mins
}

func (v *CompositeView) diffLinesMax(mins map[int]int, maxDiff int) map[int]int {
	diffs := make(map[int]int, len(v.elements))
	for i, e := range v.elements {
		min, max := mins[i], e.MaxLines()
		// Invalid min and max.
		if max < min {
			diffs[i] = 0
			continue
		}

		diff := max - min
		if max < 0 || diff > maxDiff {
			diff = maxDiff
		}

		diffs[i] = diff
	}

	return diffs
}

func countPositive(m map[int]int) int {
	var cnt int
	for _, v := range m {
		if v > 0 {
			cnt++
		}
	}
	return cnt
}
