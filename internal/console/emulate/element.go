package emulate

import (
	"fmt"
	"mltwist/internal/console/internal/view"
	"mltwist/internal/state"
	"mltwist/pkg/expr"
	"sort"
	"strings"
)

const (
	valuesWidth = 80
	regsPerLine = 2
)

var _ view.Element = &valuesElement{}

type valuesElement struct {
	state *state.State
}

func newValuesElement(state *state.State) *valuesElement {
	return &valuesElement{
		state: state,
	}
}

func (v *valuesElement) lines() int {
	regCnt := len(v.state.Regs)
	if _, ok := v.state.Regs[expr.IPKey]; ok {
		regCnt--
	}

	lines := regCnt / regsPerLine
	if regCnt%regsPerLine > 0 {
		lines++
	}

	return lines
}

func (v *valuesElement) MinLines() int { return v.lines() }
func (v *valuesElement) MaxLines() int { return v.lines() }

func regKeys(regs state.RegMap) []expr.Key {
	keys := make([]expr.Key, 0, len(regs))
	for k := range regs {
		// This one will be visualized by cursor position.
		if k == expr.IPKey {
			continue
		}

		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
}

func (v *valuesElement) printLine(line []expr.Key) error {
	regs := make([]string, len(line))
	for i, k := range line {
		val := v.state.Regs[k].(expr.Const)

		bs := make([]byte, val.Width())
		copy(bs, val.Bytes())
		revertBytes(bs)

		regs[i] = fmt.Sprintf("%s: 0x%x", k, bs)
	}

	var maxWidth int
	for _, r := range regs {
		if l := len(r); maxWidth < l {
			maxWidth = l
		}
	}

	rem := valuesWidth - regsPerLine*maxWidth
	spaces := rem / (regsPerLine + 1)
	if spaces < 1 {
		return fmt.Errorf(
			"not enough space to render %d values with width %d: %d",
			len(regs), maxWidth, rem)
	}

	for _, r := range regs {
		fmt.Print(strings.Repeat(" ", spaces))
		fmt.Print(r)
		if len(r) < maxWidth {
			fmt.Print(strings.Repeat(" ", maxWidth-len(r)))
		}
	}
	fmt.Printf("\n")

	return nil
}

func (v *valuesElement) Print(n int) error {
	keys := regKeys(v.state.Regs)
	for i := 0; i < len(keys); i += regsPerLine {
		line := keys[i:]

		l := regsPerLine
		if len(line) < l {
			l = len(line)
		}

		err := v.printLine(line[:l])
		if err != nil {
			return fmt.Errorf("cannot print line %d: %w", i/regsPerLine, err)
		}
	}

	return nil
}

func revertBytes(bs []byte) {
	for i := 0; i < len(bs)/2; i++ {
		bs[i], bs[len(bs)-i-1] = bs[len(bs)-i-1], bs[i]
	}
}
