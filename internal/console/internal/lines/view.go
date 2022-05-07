package lines

import (
	"fmt"
	"math"
	"mltwist/internal/console/internal/cursor"
	"mltwist/internal/deps"
)

type View struct {
	Lines  *Lines
	Cursor *cursor.Cursor

	format string
}

func NewView(p *deps.Program) *View {
	lns := newLines(p)

	idFormat := fmt.Sprintf("%%%dd", numDigits(lns.Len(), 10))
	markFormat := fmt.Sprintf("%%%ds", MaxMarkLen)
	format := fmt.Sprintf("%%1s %s  | %s | %%s\n", idFormat, markFormat)

	return &View{
		Lines:  lns,
		Cursor: cursor.New(lns.Len()),

		format: format,
	}
}

func (*View) MinLines() int   { return 5 }
func (v *View) MaxLines() int { return v.Lines.Len() }

func (v *View) Print(n int) error {
	offset := v.Cursor.Value()

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

func (v *View) Format(i int) string {
	var cursor string
	if i == v.Cursor.Value() {
		cursor = ">"
	}

	l := v.Lines.Index(i)
	return fmt.Sprintf(v.format, cursor, i, l.Mark(), l.String())
}

func (v *View) ShiftCursor(offset int) error { return v.Cursor.Set(v.Cursor.Value() + offset) }

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
