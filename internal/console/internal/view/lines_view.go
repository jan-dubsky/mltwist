package view

import (
	"fmt"
	"math"
	"mltwist/internal/console/internal/cursor"
	"mltwist/internal/console/internal/lines"
)

type LinesView struct {
	l *lines.Lines
	c *cursor.Cursor

	format string
}

func NewLinesView(l *lines.Lines, c *cursor.Cursor) *LinesView {
	idFormat := fmt.Sprintf("%%%dd", numDigits(l.Len(), 10))
	markFormat := fmt.Sprintf("%%%ds", lines.MaxMarkLen)
	format := fmt.Sprintf("%%1s %s  | %s | %%s\n", idFormat, markFormat)

	return &LinesView{
		l: l,
		c: c,

		format: format,
	}
}

func (*LinesView) MinLines() int   { return 5 }
func (v *LinesView) MaxLines() int { return v.l.Len() }

func (v *LinesView) Print(n int) error {
	offset := v.c.Value()

	// Golden ratio calculation.
	begin := offset - int(math.Floor(float64(n)/(math.Phi+1)))
	if begin < 0 {
		begin = 0
	}
	end := begin + n

	for i := begin; i < end; i++ {
		fmt.Print(v.Format(i))
	}

	return nil
}

func (v *LinesView) Format(i int) string {
	var cursor string
	if i == v.c.Value() {
		cursor = ">"
	}

	l := v.l.Index(i)
	return fmt.Sprintf(v.format, cursor, i, l.Mark(), l.String())
}

func numDigits(num int, base int) int {
	if num == 0 {
		return 1
	}

	var cnt int

	// The minus (-) sign.
	if num < 0 {
		cnt++
	}

	for ; num != 0; num /= base {
		cnt++
	}

	return cnt
}
