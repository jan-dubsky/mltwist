package view

import (
	"decomp/internal/console/internal/cursor"
	"decomp/internal/console/internal/lines"
	"fmt"
	"math"

	"golang.org/x/crypto/ssh/terminal"
)

const minHeight = 5

type View struct {
	l *lines.Lines
	c *cursor.Cursor

	format string
}

func New(l *lines.Lines, c *cursor.Cursor) *View {
	idFormat := fmt.Sprintf("%%%dd", numDigits(l.Len(), 10))
	markFormat := fmt.Sprintf("%%%ds", lines.MaxMarkLen)
	format := fmt.Sprintf("  %s  | %s | %%s\n", idFormat, markFormat)

	return &View{
		l: l,
		c: c,

		format: format,
	}
}

func (v *View) Print() error {
	_, screenHeight, err := terminal.GetSize(0)
	if err != nil {
		return fmt.Errorf("cannot get terminal size: %w", err)
	}

	// Clean the screen.
	fmt.Print("\033[H\033[2J")

	height := screenHeight - 3
	if height < minHeight {
		fmt.Printf("screen height is not sufficient: %d > %d", height, minHeight)
		return nil
	}

	offset := v.c.Value()
	begin := offset - int(math.Ceil(float64(height)/(math.Phi+1)))
	if begin < 0 {
		begin = 0
	}
	end := begin + height

	for i := begin; i < end; i++ {
		fmt.Print(v.Format(i))
	}

	return nil
}

func (v *View) Format(i int) string {
	l := v.l.Index(i)
	return fmt.Sprintf(v.format, i, l.Mark(), l.String())
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
