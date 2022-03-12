package lines

import (
	"fmt"
	"strings"
)

var blockIndent = strings.Repeat(" ", 4)

// Line represents a single line of the instruction visualization.
type Line struct {
	mark  string
	value string

	// block is number of basic block in the model this line refers to.
	// Negative value of block means that the line doesn't belong to any
	// block.
	block int
	// instr is the index of instruction in the basic block referred by this
	// line. Negative value of block means that the line doesn't belong to
	// any instruction.
	instr int
}

func newEmptyLine() Line {
	return Line{
		block: -1,
		instr: -1,
	}
}

func (l Line) String() string { return l.value }
func (l Line) Mark() string   { return l.mark }

func (l *Line) setMark(s string) {
	if l := len(s); l > MaxMarkLen {
		panic(fmt.Sprintf("mark is too long: %d > %d", l, MaxMarkLen))
	}

	l.mark = s
}
