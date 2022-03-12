package view

import (
	"decomp/internal/console/internal/lines"
	"fmt"

	"golang.org/x/crypto/ssh/terminal"
)

const minHeight = 5

type View struct {
	l *lines.Lines

	format string
}

func New(l *lines.Lines) *View {
	idFormat := fmt.Sprintf("%%%dd", numDigits(l.Len(), 10))
	markFormat := fmt.Sprintf("%%%ds", lines.MaxMarkLen)

	format := fmt.Sprintf("  %s  | %s | %%s\n", idFormat, markFormat)
	return &View{
		l:      l,
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

	offset := v.l.Offset()
	lines := v.l.Lines(height)

	for i, l := range lines {
		fmt.Printf(v.format, i+offset, l.Mark(), l.String())
	}

	return nil
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
