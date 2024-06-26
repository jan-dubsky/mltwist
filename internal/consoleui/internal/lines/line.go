package lines

import (
	"fmt"
	"mltwist/internal/deps"
	"strings"
)

// instrMaxLen is maximal expected length of instruction textual representation.
const instrMaxLen = 24

// instrLineFormat is cached format string to produce an instruction line value.
var instrLineFormat = fmt.Sprintf("%4s %%-%ds | %%s", "", instrMaxLen)

// Line represents a single Line of the instruction visualization.
type Line struct {
	value string
	mark  Mark

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

func newBlockLine(b deps.Block) Line {
	return Line{
		// We number blocks from 1 as Block zero doesn't look good to
		// humans.
		value: fmt.Sprintf("Block %d: 0x%x", b.Idx()+1, b.Begin()),
		block: b.Idx(),
		instr: -1,
	}
}

func byteStr(bs []byte) string {
	var sb strings.Builder
	sb.Grow(3*len(bs) - 1)

	for i, b := range bs {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%02X", b))
	}

	return sb.String()
}

func newInstrLine(b deps.Block, ins deps.Instruction) Line {
	return Line{
		value: fmt.Sprintf(instrLineFormat, ins.String(), byteStr(ins.Bytes())),
		block: b.Idx(),
		instr: ins.Idx(),
	}
}

func (l *Line) setMark(m Mark) {
	if l := len(m); l > MaxMarkLen {
		panic(fmt.Sprintf("mark is too long: %d > %d", l, MaxMarkLen))
	}

	l.mark = m
}

func (l Line) String() string           { return l.value }
func (l Line) Mark() Mark               { return l.mark }
func (l Line) Block() (int, bool)       { return l.block, l.block >= 0 }
func (l Line) Instruction() (int, bool) { return l.instr, l.instr >= 0 }
